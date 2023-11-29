package team_3

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"math"
)

type reputation struct {
	// score, 0~1
	recentContribution  float64
	historyContribution float64
	//decisionSimilarity  float64
	isSameColor  float64
	lootBoxGet   float64
	energyGain   float64
	energyRemain float64

	// memory or counter
	pedalCnt          float64
	lastEnergyLevel   float64
	energyReceivedCnt float64
	lootBoxGetCnt     float64
}

func (rep *reputation) updateScore(biker objects.IBaseBiker, preferredColor utils.Colour) {
	// update memory
	pedal := biker.GetForces().Pedal - biker.GetForces().Brake
	rep.pedalCnt += pedal
	if biker.GetEnergyLevel() > rep.lastEnergyLevel {
		gain := biker.GetEnergyLevel() - rep.lastEnergyLevel
		rep.energyReceivedCnt += gain
		rep.lootBoxGetCnt += 1
	}
	rep.lastEnergyLevel = biker.GetEnergyLevel()

	// update score
	rep.recentContribution = normalize(pedal)
	rep.historyContribution = normalize(rep.pedalCnt)
	rep.energyRemain = normalize(rep.lastEnergyLevel)
	rep.energyGain = normalize(rep.energyReceivedCnt)
	rep.lootBoxGet = normalize(rep.lootBoxGetCnt)
	if biker.GetColour() == preferredColor {
		rep.isSameColor = 1
	} else {
		rep.isSameColor = 0
	}
}

func normalize(input float64) (output float64) {
	output = 2.0/(1.0+math.Exp(-input)) - 1
	return
}
