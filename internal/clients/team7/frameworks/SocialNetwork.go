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
	for agentId := range input.AgentDecisions {
		if agentId != sn.myId {
			agentIds = append(agentIds, agentId)
		}
	}

	DistributionPenaltyMap := make(map[uuid.UUID]float64)
	if len(input.AgentResourceVotes) > 0 {
		DistributionPenaltyMap = sn.CalcDistributionPenalties(input.AgentResourceVotes, input.AgentEnergyLevels)
	}
	PedallingPenaltyMap := sn.CalcPedallingPenalties(input.AgentDecisions, input.AgentEnergyLevels)
	OrientationPenaltyMap := sn.CalcTurningPenalties(input.AgentDecisions, input.BikeTurnAngle)
	BrakingPenaltyMap := sn.CalcBrakingPenalties(input.AgentDecisions)
	DifferentLootPenaltyMap := sn.CalcDifferentLootBoxPenalties(input.AgentColours)

	W_dp := 1.0
	W_op := 1.0
	W_bp := 1.0
	W_pp := 1.0
	W_dlp := 1.0

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
			(W_bp * BrakingPenaltyMap[agentId]) +
			(W_dlp * DifferentLootPenaltyMap[agentId])

		if len(input.AgentResourceVotes) > 0 {
			updatedTrust += (W_dp * DistributionPenaltyMap[agentId])
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
func (sn *SocialNetwork) CalcDistributionPenalties(resourceDistribution map[uuid.UUID]voting.IdVoteMap, energyLevelMap map[uuid.UUID]float64) map[uuid.UUID]float64 {
	bikerCount := len(resourceDistribution)

	idPenaltyMap := make(map[uuid.UUID]float64)

	expectedEgalitarianValue := 1.0 / float64(bikerCount)
	for agentId, agentDistributionVote := range resourceDistribution {
		utilitarianPenalty := 0.0
		egalitarianPenalty := 0.0
		selfishPenalty := 0.0
		judgementalPenalty := 0.0
		for recipientId, recipientDistribution := range agentDistributionVote {
			egalitarianPenalty += math.Abs(expectedEgalitarianValue - recipientDistribution)

			utilitarianPenalty += math.Abs(recipientDistribution - (1 - energyLevelMap[recipientId]))

			if recipientId == sn.myId {
				selfishPenalty = math.Abs(1 - recipientDistribution)
			}

			if recipientId == agentId {
				judgementalPenalty = recipientDistribution
			}
		}

		overallPenalty := egalitarianPenalty*sn.personality.Egalitarian +
			selfishPenalty*sn.personality.Selfish +
			judgementalPenalty*sn.personality.Judgemental +
			utilitarianPenalty*sn.personality.Utilitarian
		idPenaltyMap[agentId] = overallPenalty

	}
	return idPenaltyMap // Return the calculated penalty map
}

// TODO: Find shift to account for forgiveness
func (sn *SocialNetwork) CalcPedallingPenalties(agentForces map[uuid.UUID]utils.Forces, energyLevelMap map[uuid.UUID]float64) map[uuid.UUID]float64 {
	agentCount := len(agentForces)

	pedallingPenaltyMap := make(map[uuid.UUID]float64)

	totalPedalling := 0.0
	for _, forces := range agentForces {
		totalPedalling += forces.Pedal
	}
	expectedPedalValue := totalPedalling / float64(agentCount)

	for id, forces := range agentForces {
		agentPedallingForce := forces.Pedal
		egalitarianPenalty := math.Abs(expectedPedalValue - agentPedallingForce)
		utilitarianPenalty := (1-agentPedallingForce)*(energyLevelMap[id]) - (math.Pow(agentPedallingForce, 1.0/agentPedallingForce))*0.3
		judgementalPenalty := 1 - agentPedallingForce
		selfishPenalty := 1 - agentPedallingForce

		overallPenalty := egalitarianPenalty*sn.personality.Egalitarian +
			selfishPenalty*sn.personality.Selfish +
			judgementalPenalty*sn.personality.Judgemental +
			utilitarianPenalty*sn.personality.Utilitarian

		pedallingPenaltyMap[id] = overallPenalty
	}

	return pedallingPenaltyMap
}

func (sn *SocialNetwork) CalcTurningPenalties(agentForces map[uuid.UUID]utils.Forces, proposedTurnAngle float64) map[uuid.UUID]float64 {
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

	return turningPenaltyMap
}

func (sn *SocialNetwork) CalcBrakingPenalties(agentForces map[uuid.UUID]utils.Forces) map[uuid.UUID]float64 {
	brakingPenaltyMap := make(map[uuid.UUID]float64)
	for id, forces := range agentForces {
		if forces.Brake == 0 {
			brakingPenaltyMap[id] = 0.8
			continue
		}
	}
	return brakingPenaltyMap
}

func (sn *SocialNetwork) CalcDifferentLootBoxPenalties(agentColourMap map[uuid.UUID]utils.Colour) map[uuid.UUID]float64 {
	myColour := agentColourMap[sn.myId]
	colourPenaltyMap := make(map[uuid.UUID]float64)

	for _, colour := range agentColourMap {
		penalty := 0.0
		if colour == myColour {
			penalty = -0.2
		} else {
			penalty = 0.095
		}
		penalty = (1 - sn.personality.Egalitarian) * penalty
		colourPenaltyMap[sn.myId] = penalty
	}

	return colourPenaltyMap
}
