package team4

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"fmt"

	"github.com/google/uuid"
)

func (agent *BaselineAgent) InitializeDecisionData() {
	//fmt.Println("Initializing decision data ...")
	//Initialize mapping if not initialized yet (= nil)
	if agent.energyHistory == nil {
		agent.energyHistory = make(map[uuid.UUID][]float64)
	}
	//fmt.Println("Energy history")
	if len(agent.mylocationHistory) == 0 {
		agent.mylocationHistory = make([]utils.Coordinates, 0)
	}
	//fmt.Println("Location history")
	if agent.honestyMatrix == nil {
		agent.honestyMatrix = make(map[uuid.UUID]float64)
	}
	//.Println("Honesty matrix")
	if agent.reputation == nil {
		agent.reputation = make(map[uuid.UUID]float64)
	}
	//fmt.Println("Reputation matrix")
}
func (agent *BaselineAgent) UpdateDecisionData() {
	//fmt.Println("Updating decision data ...")
	agent.InitializeDecisionData()

	messages := agent.GetAllMessages(agent.GetFellowBikers())
	// Process messages
	for _, msg := range messages {
		switch msg := msg.(type) {
		case objects.KickoutAgentMessage:
			agent.HandleKickoutMessage(msg)
		case objects.ReputationOfAgentMessage:
			agent.HandleReputationMessage(msg)
		case objects.JoiningAgentMessage:
			agent.HandleJoiningMessage(msg)
		case objects.ForcesMessage:
			agent.HandleForcesMessage(msg)
		case objects.GovernanceMessage:
			agent.HandleGovernanceMessage(msg)
		case objects.LootboxMessage:
			agent.HandleLootboxMessage(msg)
		case objects.VoteGoveranceMessage:
			agent.HandleVoteGovernanceMessage(msg)
		case objects.VoteKickoutMessage:
			agent.HandleVoteKickoutMessage(msg)
		case objects.VoteLootboxDirectionMessage:
			agent.HandleVoteLootboxDirectionMessage(msg)
		case objects.VoteRulerMessage:
			agent.HandleVoteRulerMessage(msg)
		}
	}

	agent.onBike = agent.GetBikeStatus()
	if agent.onBike {
		//update location history for the agent
		agent.mylocationHistory = append(agent.mylocationHistory, agent.GetLocation())
		//get fellow bikers
		fellowBikers := agent.GetFellowBikers()
		//update current bike for the agent
		agent.currentBike = agent.GetBike()
		//update the audi's current target bike ID
		agent.audiTarget = agent.GetGameState().GetAudi().GetTargetID()
		//update governance for the current bike
		agent.currentGovernance = agent.GetGameState().GetMegaBikes()[agent.currentBike].GetGovernance()
		//update ruler for the current bike
		agent.currentRuler = agent.GetGameState().GetMegaBikes()[agent.currentBike].GetRuler()
		//update capacity (number of agents on my bike)
		agent.capacity = len(fellowBikers)
		//update energy history for each fellow biker
		for _, fellow := range fellowBikers {
			fellowID := fellow.GetID()
			currentEnergyLevel := fellow.GetEnergyLevel()
			//Append bikers current energy level to the biker's history
			agent.energyHistory[fellowID] = append(agent.energyHistory[fellowID], currentEnergyLevel)
		}
	} else {
		agent.currentBike = uuid.Nil
	}
	//fmt.Println("Energy history for")
	//call reputation and honesty matrix to calcuiate/update them
	//save updated reputation and honesty matrix
	agent.CalculateReputation()
	agent.CalculateHonestyMatrix()
	//agent.DisplayFellowsEnergyHistory()
	//agent.DisplayFellowsHonesty()
	//agent.DisplayFellowsReputation()
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

// //////////////////////////// DISPLAY FUNCTIONS ////////////////////////////////////////
func (agent *BaselineAgent) DisplayFellowsEnergyHistory() {
	fellowBikers := agent.GetFellowBikers()
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		fmt.Println("")
		fmt.Println("Energy history for: ", fellowID)
		fmt.Print(agent.energyHistory[fellowID])
		fmt.Println("")
	}
}
func (agent *BaselineAgent) DisplayFellowsHonesty() {
	fellowBikers := agent.GetFellowBikers()
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		fmt.Println("")
		fmt.Println("Honesty Matrix for: ", fellowID)
		fmt.Print(agent.honestyMatrix[fellowID])
		fmt.Println("")
	}
}
func (agent *BaselineAgent) DisplayFellowsReputation() {
	fellowBikers := agent.GetFellowBikers()
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		fmt.Println("")
		fmt.Println("Reputation Matrix for: ", fellowID)
		fmt.Print(agent.reputation[fellowID])
		fmt.Println("")
	}
}
