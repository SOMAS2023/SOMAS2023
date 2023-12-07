package team5Agent

// 0 = conservative, 1 = observer, 2 = esteem, 3 = ultristic

func (t5 *team5Agent) updateState() {
	currentEnergy := t5.GetEnergyLevel()
	// Example condition: switch to conservative (Boris biker) if energy is low
	if currentEnergy < 0.1 {
		t5.state = 0
	}
	if currentEnergy > 40 {
		t5.state = 2
	}
	if currentEnergy > 90 {
		t5.state = 3
	} else {

		t5.state = 1
	}
}

// func (t5 *team5Agent) DecideAction() objects.BikerAction {
// 	switch t5.currentState {
// 	case Conservative:
// 		// Define conservative action
// 		return objects.ConserveEnergy
// 	case Aggressive:
// 		// Define aggressive action
// 		return objects.PedalHard
// 	default:
// 		// Normal action
// 		return objects.Pedal
// 	}
//}
