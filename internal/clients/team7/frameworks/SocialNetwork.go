package frameworks

import (
	"SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"
	"math"

	"github.com/google/uuid"
)

const maxTrustIterations int = 6

type SocialConnection struct {
	connectionAge      int       // Number of rounds the agent has been known
	trustLevels        []float64 // Trust level of the agent, dummy float64 for now.
	isActiveConnection bool      // Boolean indicating if connection is on current bike
}

type SocialNetworkUpdateInput struct {
	AgentIds           []uuid.UUID
	AgentDecisions     map[uuid.UUID]utils.Forces
	AgentResourceVotes map[uuid.UUID]voting.IdVoteMap
	AgentEnergyLevels  map[uuid.UUID]float64
	AgentColours       map[uuid.UUID]utils.Colour
	BikeTurnAngle      float64
	// bikeLeaderId       uuid.UUID
}

type ISocialNetwork[T any] interface {
	GetSocialNetwork() map[uuid.UUID]SocialConnection
	GetCurrentTrustLevels() map[uuid.UUID]float64
	UpdateSocialNetwork(agentIds []uuid.UUID, inputs T)
	UpdateActiveConnections(agentIds []uuid.UUID)
	DeactivateConnections(agentIds []uuid.UUID)
}

type SocialNetwork struct {
	ISocialNetwork[SocialNetworkUpdateInput]
	socialNetwork map[uuid.UUID]*SocialConnection
	personality   *Personality
	myId          uuid.UUID
}

func NewSocialNetwork(myId uuid.UUID, personality *Personality) *SocialNetwork {
	return &SocialNetwork{
		socialNetwork: map[uuid.UUID]*SocialConnection{},
		myId:          myId,
		personality:   personality,
	}
}

func (sc *SocialConnection) GetCurrentTrustLevels() float64 {
	return sc.trustLevels[len(sc.trustLevels)-1]
}

func (sc *SocialConnection) GetAverageTrustLevels() float64 {
	totalTrust := 0.0
	for _, trustLevel := range sc.trustLevels {
		totalTrust += trustLevel
	}
	return totalTrust / float64(len(sc.trustLevels))
}

func (sn *SocialNetwork) GetCurrentTrustLevels() map[uuid.UUID]float64 {
	trustLevels := make(map[uuid.UUID]float64)
	for agentId, connection := range sn.socialNetwork {
		trustLevels[agentId] = connection.trustLevels[len(connection.trustLevels)-1]
	}
	return trustLevels
}

func (sn *SocialNetwork) GetAverageTrustLevels() map[uuid.UUID]float64 {
	averageTrustLevels := make(map[uuid.UUID]float64)
	for agentId, connection := range sn.socialNetwork {
		totalTrust := 0.0
		for _, trustLevel := range connection.trustLevels {
			totalTrust += trustLevel
		}
		averageTrustLevels[agentId] = totalTrust / float64(len(connection.trustLevels))
	}
	return averageTrustLevels
}

func (sn *SocialNetwork) GetSocialNetwork() map[uuid.UUID]*SocialConnection {
	return sn.socialNetwork
}

func (sn *SocialNetwork) updateTrustLevels(input SocialNetworkUpdateInput) {

	agentIds := make([]uuid.UUID, 0)
	for _, agentId := range input.AgentIds {
		if agentId != sn.myId {
			agentIds = append(agentIds, agentId)
		}
	}

	DistributionPenaltyMap := make(map[uuid.UUID]float64)
	if len(input.AgentResourceVotes) > 0 {
		DistributionPenaltyMap = sn.CalcDistributionPenalties(input.AgentResourceVotes, input.AgentEnergyLevels, agentIds)
	}
	PedallingPenaltyMap := sn.CalcPedallingPenalties(input.AgentDecisions, input.AgentEnergyLevels, agentIds)
	OrientationPenaltyMap := sn.CalcTurningPenalties(input.AgentDecisions, input.BikeTurnAngle, agentIds)
	BrakingPenaltyMap := sn.CalcBrakingPenalties(input.AgentDecisions, agentIds)

	W_dp := 1.0
	W_op := 1.0
	W_bp := 1.0
	W_pp := 1.0

	for _, agentId := range agentIds {
		_, exists := sn.socialNetwork[agentId]
		if !exists {
			sn.socialNetwork[agentId] = &SocialConnection{
				connectionAge:      0,
				trustLevels:        []float64{},
				isActiveConnection: true,
			}
		}

		updatedTrust := (W_pp * PedallingPenaltyMap[agentId]) +
			(W_op * OrientationPenaltyMap[agentId]) +
			(W_bp * BrakingPenaltyMap[agentId])

		if len(input.AgentResourceVotes) > 0 {
			updatedTrust += (W_dp * DistributionPenaltyMap[agentId])
		}

		if updatedTrust > 1 {
			updatedTrust = 1
		}
		if updatedTrust < 0 {
			updatedTrust = 0
		}

		trustLevels := sn.socialNetwork[agentId].trustLevels
		if len(trustLevels) < maxTrustIterations {
			trustLevels = append(trustLevels, updatedTrust)
		} else {
			trustLevels = append(trustLevels[1:], updatedTrust)
		}
		sn.socialNetwork[agentId].trustLevels = trustLevels
	}
}

func (sn *SocialNetwork) UpdateSocialNetwork(agentIds []uuid.UUID, inputs SocialNetworkUpdateInput) {
	sn.updateActiveConnections(agentIds)
	sn.updateTrustLevels(inputs)
}

func (sn *SocialNetwork) updateActiveConnections(agentIds []uuid.UUID) {
	for agentId, connection := range sn.socialNetwork {
		agentIsOnBike := false
		for _, id := range agentIds {
			if agentId == id {
				agentIsOnBike = true
				connection.connectionAge++
			}
			continue
		}
		connection.isActiveConnection = agentIsOnBike
	}
}

// Retrieve agents on the current bike
func (sn *SocialNetwork) GetCurrentBikeNetwork() map[uuid.UUID]SocialConnection {
	activeConnections := map[uuid.UUID]SocialConnection{}
	for agentId, connection := range sn.socialNetwork {
		if connection.isActiveConnection {
			activeConnections[agentId] = *connection
		}
	}
	return activeConnections
}

// Implement individual calculation methods within the SocialNetwork

// Calc_Distribution_penalty calculates the penalty based on resources given
// and resources requested. This is a method of the SocialNetwork type.
func (sn *SocialNetwork) CalcDistributionPenalties(resourceDistribution map[uuid.UUID]voting.IdVoteMap, energyLevelMap map[uuid.UUID]float64, agentIds []uuid.UUID) map[uuid.UUID]float64 {

	idPenaltyMap := make(map[uuid.UUID]float64)

	for agentId, agentDistributionVote := range resourceDistribution {

		selfishPenalty := 0.0

		for recipientId, recipientDistribution := range agentDistributionVote {

			if recipientId == sn.myId {
				selfishPenalty = math.Abs(1 - recipientDistribution)
			}

		}

		overallPenalty := selfishPenalty * (1 + sn.personality.Neuroticism + (1 - sn.personality.Agreeableness))

		idPenaltyMap[agentId] = overallPenalty

	}

	for _, agentId := range agentIds {
		if _, ok := idPenaltyMap[agentId]; !ok {
			idPenaltyMap[agentId] = 0.005
		}
	}

	return idPenaltyMap // Return the calculated penalty map
}

// TODO: Find shift to account for forgiveness
func (sn *SocialNetwork) CalcPedallingPenalties(agentForces map[uuid.UUID]utils.Forces, energyLevelMap map[uuid.UUID]float64, agentIds []uuid.UUID) map[uuid.UUID]float64 {

	pedallingPenaltyMap := make(map[uuid.UUID]float64)

	for id, forces := range agentForces {
		agentPedallingForce := forces.Pedal

		Penalty := 1 - agentPedallingForce

		overallPenalty := Penalty * (1 + sn.personality.Conscientiousness)

		pedallingPenaltyMap[id] = overallPenalty
	}

	for _, agentId := range agentIds {
		if _, ok := pedallingPenaltyMap[agentId]; !ok {
			pedallingPenaltyMap[agentId] = 0.005
		}
	}

	return pedallingPenaltyMap
}

func (sn *SocialNetwork) CalcTurningPenalties(agentForces map[uuid.UUID]utils.Forces, proposedTurnAngle float64, agentIds []uuid.UUID) map[uuid.UUID]float64 {
	turningPenaltyMap := make(map[uuid.UUID]float64)

	for id, forces := range agentForces {
		turningDecision := forces.Turning
		// if id == leaderId {
		// 	if turningDecision.SteerBike && turningDecision.SteeringForce == proposedTurnAngle {
		// 		turningPenaltyMap[id] = -0.2
		// 	} else {
		// 		turningPenaltyMap[id] = 0.3
		// 	}
		// 	continue
		// }
		if turningDecision.SteerBike && turningDecision.SteeringForce == proposedTurnAngle {
			// Biker is steering in the right direction
			turningPenaltyMap[id] = -0.2
		} else if turningDecision.SteerBike && turningDecision.SteeringForce != proposedTurnAngle {
			// Biker is steering in the wrong direction
			turningPenaltyMap[id] = 0.5
		} else {
			// Biker is not steering
			turningPenaltyMap[id] = 0.2
		}
	}

	for _, agentId := range agentIds {
		if _, ok := turningPenaltyMap[agentId]; !ok {
			turningPenaltyMap[agentId] = 0.005
		}
	}

	return turningPenaltyMap
}

func (sn *SocialNetwork) CalcBrakingPenalties(agentForces map[uuid.UUID]utils.Forces, agentIds []uuid.UUID) map[uuid.UUID]float64 {
	brakingPenaltyMap := make(map[uuid.UUID]float64)
	for id, forces := range agentForces {
		if forces.Brake == 0 {
			brakingPenaltyMap[id] = 0.4
			continue
		}
	}

	for _, agentId := range agentIds {
		if _, ok := brakingPenaltyMap[agentId]; !ok {
			brakingPenaltyMap[agentId] = 0.005
		}
	}

	return brakingPenaltyMap
}
