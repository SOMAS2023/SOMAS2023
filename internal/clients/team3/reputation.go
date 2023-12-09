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
	//decisionSimilarity  float64
	isSameColor     float64
	lootBoxGet      float64
	energyGain      float64
	energyRemain    float64
	recentGetEnergy bool

	// memory or counter
	_lastPedal         float64
	_pedalCnt          float64
	_lastEnergyLevel   float64
	_recentEnergyGain  float64
	_energyReceivedCnt float64
	_lootBoxGetCnt     float64
}

func (rep *reputation) updateScore(biker objects.IBaseBiker, preferredColor utils.Colour) {
	// update memory
	rep.recentGetEnergy = biker.GetEnergyLevel() > rep._lastEnergyLevel
	if biker.GetEnergyLevel() > rep._lastEnergyLevel {
		rep._recentEnergyGain = biker.GetEnergyLevel() - rep._lastEnergyLevel
		rep._energyReceivedCnt += rep._recentEnergyGain
		rep._lootBoxGetCnt += 1.0
	} else {
		rep._lastPedal = rep._lastEnergyLevel - biker.GetEnergyLevel()
		// if agent gain energy in this iter, assume it pedals with the same energy cost as last iter
	}
	rep._lastEnergyLevel = biker.GetEnergyLevel()
	rep._pedalCnt += rep._lastPedal

	// update score
	rep.recentContribution = normalize(rep._lastPedal)
	rep.historyContribution = normalize(rep.historyContribution)
	rep.energyRemain = normalize(rep._lastEnergyLevel)
	rep.energyGain = normalize(rep._energyReceivedCnt)
	rep.lootBoxGet = normalize(rep._lootBoxGetCnt)
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
