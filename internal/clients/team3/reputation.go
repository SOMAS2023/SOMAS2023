package team3

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"math"
)

type reputation struct {
	// score, 0~1
	recentContribution  float64
	historyContribution float64
	opinionSimilarity   float64
	isSameColor         float64
	lootBoxGet          float64
	energyGain          float64
	energyRemain        float64
	recentGetEnergy     bool

	// memory or counter
	_lastEnergyCost    float64
	_energyCostCnt     float64
	_lastEnergyLevel   float64
	_recentEnergyGain  float64
	_energyReceivedCnt float64
	_lootBoxGetCnt     float64
	_sameOpinionCnt    float64
}

func (rep *reputation) updateScore(biker objects.IBaseBiker, preferredColor utils.Colour) {
	// update memory
	currentEnergy := biker.GetEnergyLevel()
	if math.IsNaN(currentEnergy) {
		currentEnergy = 0
	}
	rep.recentGetEnergy = currentEnergy > rep._lastEnergyLevel
	if currentEnergy > rep._lastEnergyLevel {
		rep._recentEnergyGain = currentEnergy - rep._lastEnergyLevel
		rep._energyReceivedCnt += rep._recentEnergyGain
		rep._lootBoxGetCnt += 1.0
	} else {
		rep._lastEnergyCost = rep._lastEnergyLevel - currentEnergy
		// if agent gain energy in this iter, assume it pedals with the same energy cost as last iter
	}
	rep._lastEnergyLevel = currentEnergy
	rep._energyCostCnt += rep._lastEnergyCost

	// update score
	rep.recentContribution = normalize(rep._lastEnergyCost)
	rep.historyContribution = normalize(rep._energyCostCnt)
	rep.energyRemain = normalize(rep._lastEnergyLevel)
	rep.energyGain = normalize(rep._energyReceivedCnt)
	rep.lootBoxGet = normalize(rep._lootBoxGetCnt)
	if biker.GetColour() == preferredColor {
		rep.isSameColor = 1
	} else {
		rep.isSameColor = 0
	}
}

func (rep *reputation) findSameOpinion() {
	rep._sameOpinionCnt += 1
	rep.opinionSimilarity = normalize(rep._sameOpinionCnt)
}

func normalize(input float64) (output float64) {
	output = 2.0/(1.0+math.Exp(-input)) - 1
	return
}
