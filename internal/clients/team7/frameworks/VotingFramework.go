package frameworks

import (
	"fmt"

	"github.com/google/uuid"
)

// This map can hold any type of data as the value
type Map map[uuid.UUID]interface{}

// Define VoteTypes
type VoteType int

const (
	VoteToKickAgent VoteType = iota
	VoteToAcceptNewAgent
	VoteOnProposals
	VoteOnAllocation
)

// Define Vote Parameters - the way we are expected to cast votes (ranking, yes/no, proportions, etc)
type VoteParameter int

const (
	Proportion VoteParameter = iota // Assign a proportion of your vote to each candidate
	YesNo                           // Say yes or no to each candidate
	// add new vote parameters as required by the environment
)

type VoteInputs struct {
	DecisionType   VoteType      // Type of vote that needs to be made
	Candidates     []uuid.UUID   // List of candidate choices
	VoteParameters VoteParameter // Parameters for the vote
}

type Vote struct {
	result map[uuid.UUID]interface{}
}

type VotingFramework struct {
	IDecisionFramework[VoteInputs, Vote]
}

func NewVotingFramework() *VotingFramework {
	return &VotingFramework{}
}

func (vf *VotingFramework) GetDecision(inputs VoteInputs) Vote {
	fmt.Println("VotingFramework: GetDecision called")
	fmt.Println("VotingFramework: Decision type: ", inputs.DecisionType)
	fmt.Println("VotingFramework: Choice map: ", inputs.Candidates)
	fmt.Println("VotingFramework: Vote parameters: ", inputs.VoteParameters)

	voteResult := vf.deliberateVote(inputs)

	return voteResult
}

func (vf *VotingFramework) deliberateVote(voteInputs VoteInputs) Vote {
	var vote Vote
	if voteInputs.DecisionType == VoteToKickAgent {
		// TODO: Deliberate on whether to kick an agent
		fmt.Println("Deliberating on whether to kick an agent")
		vote = VoteToKickWrapper(voteInputs)
	} else if voteInputs.DecisionType == VoteToAcceptNewAgent {
		// TODO: Deliberate on whether to accept a new agent
		fmt.Println("Deliberating on whether to accept a new agent")
	} else if voteInputs.DecisionType == VoteOnProposals {
		// TODO: Deliberate on how to vote on proposed directions
		fmt.Println("Deliberating on how to vote on proposals")
		//vote = VoteOnProposalsWrapper(voteInputs)
	} else if voteInputs.DecisionType == VoteOnAllocation {
		// TODO: Deliberate on how to vote on resource allocation
		fmt.Println("Deliberating on how to vote on resource allocation")
	} else {
		// TODO: Deliberate on something else
		fmt.Println("Deliberating on something else")
		//vote = Vote{result: Map{"decision": true}}
	}
	return vote
}
