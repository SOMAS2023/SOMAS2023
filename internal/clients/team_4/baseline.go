package team_4

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/google/uuid"
)

type IBaselineAgent interface {
	objects.IBaseBiker

	//INCOMPLETE/NO STRATEGY FUNCTIONS
	CalculateReputation( /*choose*/ ) map[uuid.UUID]float64    //calculate reputation matrix
	CalculateHonestyMatrix( /*choose*/ ) map[uuid.UUID]float64 //calculate honesty matrix

	DecideAction() objects.BikerAction //determines what action the agent is going to take this round. (changeBike or Pedal)
	DecideForce(direction uuid.UUID)   //defines the vector you pass to the bike: [pedal, brake, turning]
	ChangeBike() uuid.UUID             //called when biker wants to change bike, it will choose which bike to try and join
	VoteForKickout() map[uuid.UUID]int

	//CURRENTLY UNUSED/NOT CONSIDERED FUNCTIONS
	DecideGovernance() voting.GovernanceVote //decide the governance system
	VoteDictator() voting.IdVoteMap
	VoteLeader() voting.IdVoteMap
	LeadDirection() uuid.UUID //called only when the agent is the leader

	//IMPLEMENTED FUNCTIONS
	ProposeDirection() uuid.UUID                                    //returns the id of the desired lootbox
	FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap //returns rank of proposed lootboxes
	DecideAllocation() voting.IdVoteMap                             //decide the allocation parameters
	DecideJoining(pendinAgents []uuid.UUID) map[uuid.UUID]bool      //decide whether to accept or not accept bikers, ranks the ones
	nearestLoot() uuid.UUID                                         //returns the id of the nearest lootbox
	DictateDirection() uuid.UUID                                    //called only when the agent is the dictator

	//HELPER FUNCTIONS
	UpdateDecisionData()           //updates all the data needed for the decision making process(call at the start of any decision making function)
	getHonestyAverage() float64    //returns the average honesty of all agents
	getReputationAverage() float64 //returns the average reputation of all agents

	rankFellowsReputation(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) //returns normal rank of fellow bikers reputation
	rankFellowsHonesty(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error)    //returns normal rank of fellow bikers honesty

	rankTargetProposals(proposedLootBox []objects.ILootBox) (map[uuid.UUID]float64, error) //returns ranking of the proposed lootboxes

	IncreaseHonesty(agentID uuid.UUID, increaseAmount float64)
	DecreaseHonesty(agentID uuid.UUID, decreaseAmount float64)

	//PRINT FUNCTIONS
	DisplayFellowsEnergyHistory()
	DisplayFellowsHonesty()
	DisplayFellowsReputation()
}
type BaselineAgent struct {
	objects.BaseBiker
	currentBike       *objects.MegaBike
	lootBoxColour     utils.Colour
	proposedLootBox   objects.ILootBox
	mylocationHistory []utils.Coordinates     //log location history for this agent
	energyHistory     map[uuid.UUID][]float64 //log energy level for all agents
	reputation        map[uuid.UUID]float64   //record reputation for other agents, 0-1
	honestyMatrix     map[uuid.UUID]float64   //record honesty for other agents, 0-1
}

func (agent *BaselineAgent) UpdateDecisionData() {
	//Initialize mapping if not initialized yet (= nil)
	if agent.energyHistory == nil {
		agent.energyHistory = make(map[uuid.UUID][]float64)
	}
	if len(agent.mylocationHistory) == 0 {
		agent.mylocationHistory = make([]utils.Coordinates, 0)
	}
	if agent.honestyMatrix == nil {
		agent.honestyMatrix = make(map[uuid.UUID]float64)
	}
	if agent.reputation == nil {
		agent.reputation = make(map[uuid.UUID]float64)
	}
	fmt.Println("")
	fmt.Println("Updating decision data ...")
	//update location history for the agent
	agent.mylocationHistory = append(agent.mylocationHistory, agent.GetLocation())
	//get current bike
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	//get fellow bikers
	fellowBikers := currentBike.GetAgents()
	//update energy history for each fellow biker
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		currentEnergyLevel := fellow.GetEnergyLevel()
		//Append bikers current energy level to the biker's history
		agent.energyHistory[fellowID] = append(agent.energyHistory[fellowID], currentEnergyLevel)
	}
	//call reputation and honesty matrix to calcuiate/update them
	//save updated reputation and honesty matrix
	agent.CalculateReputation()
	agent.CalculateHonestyMatrix()
	// agent.DisplayFellowsEnergyHistory()
	// agent.DisplayFellowsHonesty()
	// agent.DisplayFellowsReputation()
}

func (agent *BaselineAgent) rankFellowsReputation(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) {
	totalsum := float64(0)
	rank := make(map[uuid.UUID]float64)

	for _, fellow := range agentsOnBike {
		fellowID := fellow.GetID()
		totalsum += agent.reputation[fellowID]
	}
	//normalize the reputation
	for _, fellow := range agentsOnBike {
		fellowID := fellow.GetID()
		rank[fellowID] = float64(agent.reputation[fellowID] / totalsum)
	}
	return rank, nil
}

func (agent *BaselineAgent) rankFellowsHonesty(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) {
	totalsum := float64(0)
	rank := make(map[uuid.UUID]float64)

	for _, fellow := range agentsOnBike {
		fellowID := fellow.GetID()
		totalsum += agent.honestyMatrix[fellowID]
	}
	//normalize the honesty
	for _, fellow := range agentsOnBike {
		fellowID := fellow.GetID()
		rank[fellowID] = float64(agent.honestyMatrix[fellowID] / totalsum)
	}
	return rank, nil
}

func (agent *BaselineAgent) getReputationAverage() float64 {
	sum := float64(0)
	//loop through all bikers find the average reputation
	for _, bike := range agent.GetGameState().GetMegaBikes() {
		for _, biker := range bike.GetAgents() {
			bikerID := biker.GetID()
			sum += agent.reputation[bikerID]
		}
	}
	return sum / float64(len(agent.reputation))
}
func (agent *BaselineAgent) getHonestyAverage() float64 {
	sum := float64(0)
	//loop through all bikers find the average honesty
	for _, bike := range agent.GetGameState().GetMegaBikes() {
		for _, biker := range bike.GetAgents() {
			bikerID := biker.GetID()
			sum += agent.honestyMatrix[bikerID]
		}
	}
	return sum / float64(len(agent.honestyMatrix))
}

func (agent *BaselineAgent) rankTargetProposals(proposedLootBox []objects.ILootBox) (map[uuid.UUID]float64, error) {
	rank := make(map[uuid.UUID]float64)
	ranksum := make(map[uuid.UUID]float64)
	totalsum := float64(0)
	distanceRank := float64(0)
	w1 := float64(0.7) //weight for distance
	w2 := float64(0.2) //weight for reputation
	w3 := float64(0.1) //weight for honesty
	totaloptions := len(proposedLootBox)
	minEnergyThreshold := 0.2 //if energy level is below this threshold, the agent will increase voting towards its colour lootbox

	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	fellowBikers := currentBike.GetAgents()
	//This is the relavtive reputation and honest for bikers my bike
	reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
	honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)
	//This is the absolute reputation and honest for bikers my bike
	// reputationRank  := agent.reputation
	// honestyRank  := agent.honestyMatrix
	if e1 != nil || e2 != nil {
		panic("unexpected error!")
	}
	//sort Proposed Loot Boxes by distance
	sort.Slice(proposedLootBox, func(i, j int) bool {
		return physics.ComputeDistance(currentBike.GetPosition(), proposedLootBox[i].GetPosition()) < physics.ComputeDistance(currentBike.GetPosition(), proposedLootBox[j].GetPosition())
	})

	for i, lootBox := range proposedLootBox {
		//loop through all fellow bikers and check if they have the same colour as the lootbox
		for _, fellow := range fellowBikers {
			distanceRank := float64(totaloptions - i)
			fellowID := fellow.GetID()
			if fellow.GetColour() == lootBox.GetColour() {
				weight := (w1 * distanceRank) + (w2 * reputationRank[fellowID]) + (w3 * honestyRank[fellowID])
				ranksum[lootBox.GetID()] += weight
				totalsum += weight
			}
		}

		if lootBox.GetColour() == agent.GetColour() {
			weight := distanceRank * w1 * 1.25
			//if energy level is below threshold, increase weighting towards own colour lootbox
			if agent.GetEnergyLevel() < minEnergyThreshold {
				weight *= 2
			}
			ranksum[lootBox.GetID()] += weight
			totalsum += weight
		}
		if ranksum[lootBox.GetID()] == 0 {
			weight := (distanceRank * w1 * 2.6)
			ranksum[lootBox.GetID()] = weight
			totalsum += weight
		}
	}
	for _, lootBox := range proposedLootBox {
		rank[lootBox.GetID()] = ranksum[lootBox.GetID()] / totalsum
	}

	return rank, nil
}

func (agent *BaselineAgent) FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap {
	agent.UpdateDecisionData()
	boxesInMap := agent.GetGameState().GetLootBoxes()
	boxProposed := make([]objects.ILootBox, len(proposals))
	for i, pp := range proposals {
		boxProposed[i] = boxesInMap[pp]
	}
	rank, e := agent.rankTargetProposals(boxProposed)
	if e != nil {
		panic("unexpected error!")
	}
	return rank
}

func (agent *BaselineAgent) DecideAllocation() voting.IdVoteMap {
	agent.UpdateDecisionData()
	distribution := make(map[uuid.UUID]float64) //make(voting.IdVoteMap)
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	fellowBikers := currentBike.GetAgents()
	totalEnergySpent := float64(0)
	totalAllocation := float64(0)

	reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
	honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)
	if e1 != nil || e2 != nil {
		panic("unexpected error!")
	}

	for _, fellow := range fellowBikers {
		w1 := 0.3 //weight for reputation
		w2 := 0.3 //weight for honesty
		w3 := 0.3 //weight for energy spent
		w4 := 0.1 //weight for energy level
		fellowID := fellow.GetID()
		energyLog := agent.energyHistory[fellowID]
		energySpent := energyLog[len(energyLog)-2] - energyLog[len(energyLog)-1]
		totalEnergySpent += energySpent
		// In the case where the I am the same colour as the lootbox
		if fellowID == agent.GetID() {
			w1 = 0.001
			w2 = 0.001
			w3 = 0.398
			w4 = 0.6
			if agent.lootBoxColour == agent.GetColour() {
				w1 = 0.001
				w2 = 0.001
				w3 = 0.6
				w4 = 0.6
			}
		}
		distribution[fellow.GetID()] = float64((w1 * reputationRank[fellowID]) + (w2 * honestyRank[fellowID]) + (w3 * energySpent) + (w4 * fellow.GetEnergyLevel()))
		// distribution[fellow.GetID()] = energySpent * rand.Float64() // random for now
		totalAllocation += distribution[fellow.GetID()]
	}

	//normalize the distribution
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		distribution[fellowID] = distribution[fellowID] / totalAllocation
	}

	return distribution
}

// Reputation and Honesty Matrix Teams Must Implement these or similar functions

func (agent *BaselineAgent) CalculateReputation() {
	////////////////////////////
	//  As the program I used for debugging invoked "padal" and "break" with values of 0, I conducted tests using random numbers.
	// In case of an updated main program, I will need to adjust the parameters and expressions of the reputation matrix.
	// The current version lacks real data during the debugging process.
	////////////////////////////
	megaBikes := agent.GetGameState().GetMegaBikes()

	for _, bike := range megaBikes {
		// Get all agents on MegaBike
		fellowBikers := bike.GetAgents()

		// Iterate over each agent on MegaBike, generate reputation assessment
		for _, otherAgent := range fellowBikers {
			// Exclude self
			selfTest := otherAgent.GetID() //nolint
			if selfTest == agent.GetID() {
				agent.reputation[otherAgent.GetID()] = 1.0
			}

			// Monitor otherAgent's location
			// location := otherAgent.GetLocation()
			// RAP := otherAgent.GetResourceAllocationParams()
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Location:", location, "ResourceAllocationParams:", RAP)

			// Monitor otherAgent's forces
			forces := otherAgent.GetForces()
			energyLevel := otherAgent.GetEnergyLevel()
			ReputationForces := float64(forces.Pedal+forces.Brake+rand.Float64()) / energyLevel //CAUTION: REMOVE THE RANDOM VALUE
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Forces:", ReputationForces)

			// Monitor otherAgent's bike status
			bikeStatus := otherAgent.GetBikeStatus()
			// Convert the boolean value to float64 and print the result
			ReputationBikeShift := 0.2
			if bikeStatus {
				ReputationBikeShift = 1.0
			}
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Bike_Shift", float64(ReputationBikeShift))

			// Calculate Overall_reputation
			OverallReputation := ReputationForces * ReputationBikeShift
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Overall Reputation:", OverallReputation)

			// Store Overall_reputation in the reputation map
			agent.reputation[otherAgent.GetID()] = OverallReputation
		}
	}
	// for agentID, agentReputation := range agent.reputation {
	// 	print("Agent ID: ", agentID.String(), ", Reputation: ", agentReputation, "\n")
	// }
}

func (agent *BaselineAgent) CalculateHonestyMatrix() {
	// Copy the local honesty matrix values
	for _, bike := range agent.GetGameState().GetMegaBikes() {
		for _, biker := range bike.GetAgents() {
			bikerID := biker.GetID()
			agent.honestyMatrix[bikerID] = 1.0
		}
	}
	// for agentID, _ /*honestyValue*/ := range agent.honestyMatrix {
	// 	//honestyMatrixCopy[agentID] = honestyValue
	// 	agent.honestyMatrix[agentID] = 1.0
	// }

}

func (agent *BaselineAgent) DisplayFellowsEnergyHistory() {
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	fellowBikers := currentBike.GetAgents()
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		fmt.Println("")
		fmt.Println("Energy history for: ", fellowID)
		fmt.Print(agent.energyHistory[fellowID])
		fmt.Println("")
	}
}
func (agent *BaselineAgent) DisplayFellowsHonesty() {
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	fellowBikers := currentBike.GetAgents()
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		fmt.Println("")
		fmt.Println("Honesty Matrix for: ", fellowID)
		fmt.Print(agent.honestyMatrix[fellowID])
		fmt.Println("")
	}
}
func (agent *BaselineAgent) DisplayFellowsReputation() {
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	fellowBikers := currentBike.GetAgents()
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		fmt.Println("")
		fmt.Println("Reputation Matrix for: ", fellowID)
		fmt.Print(agent.reputation[fellowID])
		fmt.Println("")
	}
}

func (agent *BaselineAgent) ProposeDirection() uuid.UUID {
	agent.proposedLootBox = nil
	lootBoxes := agent.GetGameState().GetLootBoxes()
	agentLocation := agent.GetLocation() //agent's location
	shortestDistance := math.MaxFloat64

	for _, lootbox := range lootBoxes {
		lootboxLocation := lootbox.GetPosition()
		distance := physics.ComputeDistance(agentLocation, lootboxLocation)
		if agent.proposedLootBox == nil && distance < shortestDistance {
			shortestDistance = distance
			agent.proposedLootBox = lootbox
		}
		if distance < shortestDistance || agent.GetColour() == lootbox.GetColour() {
			shortestDistance = distance
			agent.proposedLootBox = lootbox
		}
	}
	return agent.proposedLootBox.GetID()
}

// DecideAction only pedal
func (agent *BaselineAgent) DecideAction() objects.BikerAction {
	fmt.Println("Team 4")
	return objects.Pedal
}

// DecideForces randomly based on current energyLevel
func (agent *BaselineAgent) DecideForces(direction uuid.UUID) {
	//save the target lootbox
	agent.lootBoxColour = agent.GetGameState().GetLootBoxes()[direction].GetColour()

	energyLevel := agent.GetEnergyLevel() // 当前能量

	randomBreakForce := float64(0)
	randomPedalForce := rand.Float64() * energyLevel

	if randomPedalForce == 0 {
		// just random break force based on energy level, but not too much
		randomBreakForce += rand.Float64() * energyLevel * 0.5
	} else {
		randomBreakForce = 0
	}

	forces := utils.Forces{
		Pedal: randomPedalForce,
		Brake: randomBreakForce, // random for now
		Turning: utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: physics.ComputeOrientation(agent.GetLocation(), agent.GetGameState().GetMegaBikes()[direction].GetPosition()) - agent.GetGameState().GetMegaBikes()[agent.currentBike.GetID()].GetOrientation(),
		},
	}

	agent.SetForces(forces)
}

func (agent *BaselineAgent) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	agent.UpdateDecisionData()
	decision := make(map[uuid.UUID]bool)
	for _, pendingAgent := range pendingAgents {
		// energyLog := agent.energyHistory[pendingAgent]
		// energySpent := energyLog[len(energyLog)-2] - energyLog[len(energyLog)-1]
		w1 := 0.5
		w2 := 0.5
		reputation := agent.reputation[pendingAgent]
		honesty := agent.honestyMatrix[pendingAgent]
		//calculate the decision
		if (w1*reputation + w2*honesty) >= 0.55 {
			decision[pendingAgent] = true
		} else {
			decision[pendingAgent] = false
		}

	}
	return decision
}

func (agent *BaselineAgent) DecideGovernance() voting.GovernanceVote {
	rank := make(map[utils.Governance]float64)
	rank[utils.Democracy] = 1
	rank[utils.Leadership] = 0
	rank[utils.Dictatorship] = 0
	rank[utils.Invalid] = 0
	//for i := utils.Democracy; i <= utils.Invalid; i++ {
	//  rank[i] = 0.25
	//}
	return rank
}

func (agent *BaselineAgent) nearestLoot() uuid.UUID {
	currLocation := agent.GetLocation()
	shortestDist := math.MaxFloat64
	var nearestBox uuid.UUID
	var currDist float64
	for _, loot := range agent.GetGameState().GetLootBoxes() {
		x, y := loot.GetPosition().X, loot.GetPosition().Y
		currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
		if currDist < shortestDist {
			nearestBox = loot.GetID()
			shortestDist = currDist
		}
	}
	return nearestBox
}

func (agent *BaselineAgent) DictateDirection() uuid.UUID {
	agent.proposedLootBox = nil
	if agent.GetEnergyLevel() <= 30 { //prioritize survival, if low on energy, go towards bereast lootbox
		return agent.nearestLoot()
	} else {
		lootBoxes := agent.GetGameState().GetLootBoxes()
		agentLocation := agent.GetLocation() //agent's location
		shortestDistance := math.MaxFloat64

		for _, lootbox := range lootBoxes {
			lootboxLocation := lootbox.GetPosition()
			distance := physics.ComputeDistance(agentLocation, lootboxLocation)
			if agent.proposedLootBox == nil && distance < shortestDistance {
				shortestDistance = distance
				agent.proposedLootBox = lootbox
			}
			if distance < shortestDistance || agent.GetColour() == lootbox.GetColour() {
				shortestDistance = distance
				agent.proposedLootBox = lootbox
			}
		}
		return agent.proposedLootBox.GetID()
	}

}
func (agent *BaselineAgent) VoteDictator() voting.IdVoteMap {
	votes := make(voting.IdVoteMap)
	fellowBikers := agent.GetFellowBikers()
	for _, fellowBiker := range fellowBikers {
		if fellowBiker.GetColour() == agent.GetColour() { //if there is going to be a dictatorship, vote for agents with the same colour.
			votes[fellowBiker.GetID()] = 1.0
		} else {
			votes[fellowBiker.GetID()] = 0.0
		}
	}
	return votes
}
