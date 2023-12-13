package team8

import "SOMAS2023/internal/common/utils"

type GP struct {
	EnergyThreshold              float64
	DistanceThresholdForVoting   float64
	ThresholdForJoiningDecision  float64
	ThresholdForChangingMegabike float64
}

var GlobalParameters GP = GP{
	EnergyThreshold:              0.6,
	DistanceThresholdForVoting:   (utils.GridHeight + utils.GridWidth) / 4,
	ThresholdForJoiningDecision:  0.2,
	ThresholdForChangingMegabike: 0.5,
}
