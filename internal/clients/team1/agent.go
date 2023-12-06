// MAIN AGENT FILE

package team1

import (
	obj "SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	utils "SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"
	"fmt"
	"math"
	"sort"

	"github.com/google/uuid"
)

// agent specific parameters
const deviateNegative = -0.2         // trust loss on deviation
const deviatePositive = 0.1          // trust gain on non deviation
const effortScaling = 0.1            // scaling factor for effort, highr it is the more effort chages each round
const fairnessScaling = 0.1          // scaling factor for fairness, higher it is the more fairness changes each round
const votingAlignmentThreshold = 0.6 // threshold for voting alignment
const leaveThreshold = 0.0           // threshold for leaving
const kickThreshold = 0.0            // threshold for kicking
const trustThreshold = 0.7           // threshold for trusting (need to tune)
const fairnessConstant = 1           // weight of fairness in opinion
const joinThreshold = 0.8            // opinion threshold for joining if not same colour
const leaderThreshold = 0.95         // opinion threshold for becoming leader
const trustconstant = 1              // weight of trust in opinion
const effortConstant = 1             // weight of effort in opinion
const fairnessDifference = 0.5       // modifies how much fairness increases of decreases, higher is more increase, 0.5 is fair
const lowEnergyLevel = 0.3           // energy level below which the agent will try to get a lootbox of the desired colour
const leavingThreshold = 0.3         // how low the agent's vote must be to leave bike
const colorOpinionConstant = 0.2     // how much any agent likes any other of the same colour in the objective function

// Governance decision constants
const democracyOpinonThreshold = 0.5
const democracyReputationThreshold = 0.3
const leadershipOpinionThreshold = 0.7
const leadershipReputationThreshold = 0.5
const dictatorshipOpinionThreshold = 0.9
const dictatorshipReputationThreshold = 0.7

// Bike scoring constants

const majorityWeight = 3.0
const lootboxWeight = 0.5
const lootboxColourWeight = 0.5
const audiDistWeight = 0.5
const opinionWeight = 0.5
const nearbyBikeWeight = 0.5



type Biker1 struct {
	*obj.BaseBiker                       // BaseBiker inherits functions from BaseAgent such as GetID(), GetAllMessages() and UpdateAgentInternalState()
	recentVote     voting.LootboxVoteMap // the agent's most recent vote
	recentDecided  uuid.UUID             // the most recent decision
	dislikeVote    bool                  // whether the agent disliked the most recent vote
	opinions       map[uuid.UUID]Opinion
	desiredBike    uuid.UUID
}


// part 1:
// the biker itself doesn't technically have a location (as it's on the map only when it's on a bike)
// in fact this function is only called when the biker needs to make a decision about the pedaling forces
func (bb *Biker1) GetLocation() utils.Coordinates {
	gs := bb.GetGameState()
	bikeId := bb.GetBike()
	megaBikes := gs.GetMegaBikes()
	position := megaBikes[bikeId].GetPosition()
	if math.IsNaN(position.X) {
		fmt.Printf("agent %v has no position\n", bb.GetID())
	}
	return position
}

// -------------------DECISION FUNCTIONS----------------------------

func (bb *Biker1) ScoreBike(bike obj.IMegaBike) float64 {
	var majorityScore float64
	if bb.BikeOurColour(bike){
		majorityScore = 1.0
	}else{
		majorityScore = 0.0
	}
	boxCount, colourCount, bikeCount := bb.GetNearBikeObjects(bike)
	score := majorityWeight*majorityScore
	score += lootboxWeight*bb.GetLootboxScore(bike)
	score += lootboxColourWeight*bb.GetLootboxColourScore(bike)
	score += audiDistWeight*bb.DistanceFromAudi(bike)
	score += opinionWeight*bb.GetAverageOpinionOfBike(bike)
	score -= nearbyBikeWeight*bb.GetNearbyBikeScore(bike)
	
	return score
}

func (bb *Biker1) PickBestBike() uuid.UUID {
	gs := bb.GetGameState()
	allBikes := gs.GetMegaBikes()
	scoreMap := make(map[uuid.UUID]float64)
	for _, bike := range allBikes {
		if len(bike.GetAgents()) < utils.BikersOnBike || bike.GetID() == bb.GetBike() {
			scoreMap[bike.GetID()] = bb.ScoreBike(bike)
		}
	}
	if len(scoreMap) == 0 {
		return bb.GetBike()
	}
	bestBike := bb.GetBike()
	bestScore := scoreMap[bestBike]
	for id, score := range scoreMap {
		if score > bestScore {
			bestBike = id
			bestScore = score
		}
	}
	bb.desiredBike = bestBike
	return bestBike
}

func (bb *Biker1) DecideAction() obj.BikerAction {
	bb.UpdateOpinions()
	fellowBikers := bb.GetFellowBikers()
	avg_opinion := 0.0
	for _, agent := range fellowBikers {
		avg_opinion = avg_opinion + bb.opinions[agent.GetID()].opinion
	}
	if len(fellowBikers) > 0 {
		avg_opinion /= float64(len(fellowBikers))
	} else {
		avg_opinion = 1.0
	}
	if (avg_opinion < leaveThreshold) || bb.dislikeVote {
		bb.dislikeVote = false
		newBike := bb.PickBestBike()
		if newBike != bb.GetBike() {
			return 1
		} else {
			return 0
		}
	} else {
		return 0
	}
}

// -------------------END OF DECISION FUNCTIONS---------------------
// ----------------CHANGE BIKE FUNCTIONS-----------------

func (bb *Biker1) BikeOurColour(bike obj.IMegaBike) bool {
	matchCounter := 0
	totalAgents := len(bike.GetAgents())
	for _, agent := range bike.GetAgents() {
		bbColour := bb.GetColour()
		agentColour := agent.GetColour()
		if agentColour != bbColour {
			matchCounter++
		}
	}
	if matchCounter > totalAgents/2 {
		return true
	} else {
		return false
	}
}

// decide which bike to go to
func (bb *Biker1) ChangeBike() uuid.UUID {
	return bb.desiredBike
}



// -------------------BIKER ACCEPTANCE FUNCTIONS------------------------
// an agent will have to rank the agents that are trying to join and that they will try to
func (bb *Biker1) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	//gs.GetMegaBikes()[bikeId].GetAgents()

	decision := make(map[uuid.UUID]bool)

	for _, agentId := range pendingAgents {
		//TODO FIX
		agent := bb.GetAgentFromId(agentId)

		bbColour := bb.GetColour()
		agentColour := agent.GetColour()
		if agentColour == bbColour {
			decision[agentId] = true
			sameColourReward := 1.05
			bb.UpdateOpinion(agentId, sameColourReward)
		} else {
			if bb.opinions[agentId].opinion > joinThreshold {
				decision[agentId] = true
				// penalise for accepting them without same colour
				penalty := 0.9
				bb.UpdateOpinion(agentId, penalty)
			}
			decision[agentId] = false
		}
	}

	for _, agentId := range pendingAgents {
		decision[agentId] = true
	}
	return decision
}

func (bb *Biker1) lowestOpinionKick() uuid.UUID {
	fellowBikers := bb.GetFellowBikers()
	lowestOpinion := kickThreshold
	var worstAgent uuid.UUID
	for _, agent := range fellowBikers {
		if bb.opinions[agent.GetID()].opinion < lowestOpinion {
			lowestOpinion = bb.opinions[agent.GetID()].opinion
			worstAgent = agent.GetID()
		}
	}
	// if we want to kick someone based on our threshold, return their id, else return nil
	if lowestOpinion < kickThreshold {
		return worstAgent
	}
	return uuid.Nil
}

func (bb *Biker1) DecideKick(agent uuid.UUID) int {
	if bb.opinions[agent].opinion < kickThreshold {
		return 1
	}
	return 0
}

func (bb *Biker1) VoteForKickout() map[uuid.UUID]int {
	voteResults := make(map[uuid.UUID]int)
	fellowBikers := bb.GetFellowBikers()
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		if agentID != bb.GetID() {
			// random votes to other agents
			voteResults[agentID] = bb.DecideKick(agentID)
		}
	}
	return voteResults
}

//--------------------END OF BIKER ACCEPTANCE FUNCTIONS-------------------





// -------------------INSTANTIATION FUNCTIONS----------------------------
func GetBiker1(colour utils.Colour, id uuid.UUID) *Biker1 {
	fmt.Printf("Creating Biker1 with id %v\n", id)
	return &Biker1{
		BaseBiker:   obj.GetBaseBiker(colour, id),
		opinions:    make(map[uuid.UUID]Opinion),
		dislikeVote: false,
	}
}

// -------------------END OF INSTANTIATION FUNCTIONS---------------------
