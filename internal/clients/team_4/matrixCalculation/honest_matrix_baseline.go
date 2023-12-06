package matrixCalculation

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"math"
)

type honest struct {
	// score, 0~1
	recentContribution  float64
	historyContribution float64
	//decisionSimilarity  float64
	sameDirection float64
	isSameColor   float64
	lootBoxGet    float64
	energyGain    float64
	energyRemain  float64

	// memory or counter
	pedalCnt          float64
	lastEnergyLevel   float64
	energyReceivedCnt float64
	lootBoxGetCnt     float64
}

func (rep *honest) updateScore(biker objects.IBaseBiker, preferredColor utils.Colour) {
	// update memory
	pedal := biker.GetForces().Pedal - biker.GetForces().Brake
	rep.pedalCnt += pedal
	if biker.GetEnergyLevel() > rep.lastEnergyLevel {
		gain := biker.GetEnergyLevel() - rep.lastEnergyLevel
		rep.energyReceivedCnt += gain
		rep.lootBoxGetCnt += 1
	}
	rep.lastEnergyLevel = biker.GetEnergyLevel()

	//rep.sameDirection = calculateDirectionValue(objects.GetMegaBike().GetOrientation())

	// update score
	rep.recentContribution = normalize(pedal)
	rep.historyContribution = normalize(rep.pedalCnt)
	rep.energyRemain = normalize(rep.lastEnergyLevel)
	rep.energyGain = normalize(rep.energyReceivedCnt)
	rep.lootBoxGet = normalize(rep.lootBoxGetCnt)
	rep.sameDirection = normalize(rep.sameDirection)
	if biker.GetColour() == preferredColor {
		rep.isSameColor = 0.5
	} else {
		rep.isSameColor = -0.5
	}
}

func normalize(input float64) (output float64) {
	output = 2.0/(1.0+math.Exp(-input)) - 1
	return
}

func calculateDirectionValue(targetDirection float64, preferDirection float64) (output float64) {
	output = math.Abs((targetDirection-preferDirection)/180) * (-1)
	return
}
