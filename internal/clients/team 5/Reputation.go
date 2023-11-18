package reputation

type AgentID int
type FinalRep float64

// use this to calculate final reputation metric
type Reputation struct {
	PointsContribution  float64
	CooperativeBehavior float64
	SurvivalScore       float64 //Will be more reliant on you if they're dying
}

type ReputationSystem struct {
	agents map[AgentID]Reputation
}

func NewReputationSystem() *ReputationSystem {
	return &ReputationSystem{
		agents: make(map[AgentID]Reputation),
	}
}

// UpdateReputation updates the reputation for a given agent
func (rs *ReputationSystem) UpdateReputation(id AgentID, rep Reputation) {
	// Here you can add logic to normalize and update reputations
	rs.agents[id] = rep
}

// func main() {
// 	repSystem := NewReputationSystem()

// 	agentID := AgentID(1)
// 	newReputation := Reputation{
// 		PointsContribution:  10.0,
// 		CooperativeBehavior: 8.5,
// 		SurvivalScore:       9.0,
// 	}
// 	repSystem.UpdateReputation(agentID, newReputation)

// 	fmt.Printf("Agent %d's Reputation: %+v\n", agentID, repSystem.agents[agentID])
// }
