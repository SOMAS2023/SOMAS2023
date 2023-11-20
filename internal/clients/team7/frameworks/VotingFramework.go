package frameworks

import (
	"fmt"
)

// This map can hold any type of data as the value
type Map map[string]interface{}

// Define VoteTypes
type VoteType int

const (
	VoteToKickAgent VoteType = iota
	VoteToAcceptNewAgent
)

type VoteInputs struct {
	decisionType   VoteType       // Type of vote that needs to be made
	choiceMap      map[string]int // Map of choices [Dummy map for now]
	voteParameters Map            // Parameters for the vote
}

type Vote struct {
	result map[string]interface{}
}

type VotingFramework struct {
	IDecisionFramework[VoteInputs, Vote]
}

func NewVotingFramework() *VotingFramework {
	return &VotingFramework{}
}

func (vf *VotingFramework) GetDecision(inputs VoteInputs) Vote {
	fmt.Println("VotingFramework: GetDecision called")
	fmt.Println("VotingFramework: Decision type: ", inputs.decisionType)
	fmt.Println("VotingFramework: Choice map: ", inputs.choiceMap)
	fmt.Println("VotingFramework: Vote parameters: ", inputs.voteParameters)

	voteResult := vf.deliberateVote(inputs)

	return voteResult
}

func (vf *VotingFramework) deliberateVote(voteInputs VoteInputs) Vote {
	if voteInputs.decisionType == VoteToKickAgent {
		// TODO: Deliberate on whether to kick an agent
		fmt.Println("Deliberating on whether to kick an agent")
	} else if voteInputs.decisionType == VoteToAcceptNewAgent {
		// TODO: Deliberate on whether to accept a new agent
		fmt.Println("Deliberating on whether to accept a new agent")
	} else {
		// TODO: Deliberate on something else
		fmt.Println("Deliberating on something else")
	}

	return Vote{result: Map{"decision": true}}
}
