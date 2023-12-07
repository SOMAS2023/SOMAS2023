package team5Agent

// 0 = conservative, 1 = observer, 2 = esteem, 3 = altristic

func (t5 *team5Agent) updateState() {
	currentEnergy := t5.GetEnergyLevel()
	// Example condition: switch to conservative (Boris biker) if energy is low
	if t5.roundCount <= 5 {
		t5.state = 1
	} else if currentEnergy < 0.2 {
		t5.state = 0
	} else if currentEnergy < 0.85 {
		t5.state = 2
	} else {
		t5.state = 3
	}
	// //** fmt.Println("Energy Level: ", currentEnergy, "State: ", t5.state, "Round: ", t5.roundCount)
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
