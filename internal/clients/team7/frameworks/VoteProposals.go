package frameworks

import (
	"SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

// VoteProposals: This determines how our agent dritrubutes its vote on which lootbox the bike should navigate towards.
// It depends on both the opinion our agent has on each lootbox and our agent's agreeableness level.

type VoteOnProposalsHandler struct {
	IDecisionFramework[VoteOnLootBoxesInput, MapIdBool]
}

func NewVoteOnProposalsHandler() *VoteOnProposalsHandler {
	return &VoteOnProposalsHandler{}
}

func (voteHandler *VoteOnProposalsHandler) GetDecision(inputs VoteOnLootBoxesInput) voting.LootboxVoteMap {
	vote := make(voting.LootboxVoteMap)
	share := make(map[uuid.UUID]bool)

	// Whether we share depends on what our opinion is on each lootbox.
	// Low opinion => no share
	// High opinion => more of a share
	for _, loot_id := range inputs.LootBoxCandidates {
		bikerOpinionOfLootBox, hasData := inputs.MyOpinion[loot_id]
		if loot_id == inputs.MyDesired {
			share[loot_id] = true
		} else if hasData && bikerOpinionOfLootBox > 0.5 {
			share[loot_id] = true
		} else {
			share[loot_id] = false
		}
	}

	// We distribute votes across the proposed lootboxes we've deemed worthy
	// Low agreeableness => no sharing => we give ourselves 100% of the share
	// High agreeableness => all sharing => we give others a share
	agreeableness := inputs.MyPersonality.Agreeableness
	numcandidates := len(share)
	othershares := agreeableness / float64(numcandidates)
	myshare := 1 - (othershares * (float64(numcandidates) - 1))

	for _, loot_id := range inputs.LootBoxCandidates {
		bikerShare, hasData := share[loot_id]
		if loot_id == inputs.MyDesired {
			vote[loot_id] = myshare
		} else if hasData {
			if bikerShare {
				vote[loot_id] = othershares
			} else {
				vote[loot_id] = 0
			}
		} else {
			vote[loot_id] = 0
		}
	}
	return vote
}
