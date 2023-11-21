package server

import (
	"SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

// this function will take in a list of maps from ids to their corresponding vote (yes/ no in the case of acceptance)
// and retunr a list of ids that can be accepted according to some metric (ie more than half voted yes)
// ranked according to a metric (ie overall number of yes's)
func GetAcceptanceRanking([]map[uuid.UUID]bool) []uuid.UUID {
	// TODO implement
	return make([]uuid.UUID, 0)
}

// returns the winner accoring to chosen voting strategy (assumes all the maps contain a voting between 0-1
// for each option, and that all the votings sum to 1)
func WinnerFromDist([]utils.INormaliseVoteMap) uuid.UUID {
	panic("not implemented")
}
