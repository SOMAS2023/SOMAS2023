package team5Agent

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"github.com/google/uuid"
	"fmt"
)

type Team5Biker struct {
	*objects.BaseBiker
	bikeVote uuid.UUID // uuid.UUID of voted lootbox from omars output
}

//for testing use any box in targetLootBoxID

func (t5 *Team5Biker) DecideForce(targetLootBoxID uuid.UUID) {
	fmt.Println("team5Agent: GetBike: t5.BaseBiker.GetBike(): ", t5.BaseBiker.GetBike())
	return t5.BaseBiker.DecideForce(targetLootBoxID) // Pass the UUID to BaseBiker's DecideForce
}

// so this bassically adjusts the force depending on the energy of the agent
func (t5 *Team5Biker) calculatePedalForceBasedOnEnergy() float64 {
	ownEnergyLevel := t5.GetEnergyLevel()
	// ask the guys what number i want to put the own energy level and if it
	// should be adjusted based on energy level to sva eenergy
	if ownEnergyLevel < 0.5 {
		return ownEnergyLevel * utils.BikerMaxForce
	}
	return utils.BikerMaxForce
}

// here I can add implementation of stategy like:
//___________________________________________________________________________________________________________________________

// func (t5 *Team5Biker) calculateAverageEnergyOfBikeagents() float64 {
//     bikeagents := t5.GetGameState().GetMegaBikes()[t5.GetBike()].GetAgents()
//     var totalEnergy float64
//     var count float64

//     for _, agents := range bikeagents {
//         if agents.GetID() != t5.GetID() { // Exclude self
//             totalEnergy += agents.GetEnergyLevel()
//             count++
//         }
//     }

//     if count == 0 {
//         return 1 // If no other mates, return full energy or none and just bike hop maybe? see what others think
//     }
//     return totalEnergy / count
// }
//___________________________________________________________________________________________________________________________

// can add this to decide force

// add a function depends


// speed of other bikes
// and 
// position of other bikes and how fast to peddle depending on that
// so the lootbox is the direction but we may need to turn more if the bike doesnt turn enough.



//have a meeting with others discuss what other fns i can implement nd what helps others
// runs no issues