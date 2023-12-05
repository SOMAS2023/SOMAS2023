package team5Agent

// 0 = normal, 1 = conservative (Boris biker)

func (t5 *team5Agent) updateState() {
	// Example condition: switch to conservative (Boris biker) if energy is low
	if t5.GetEnergyLevel() < 0.2 {
		t5.state = 1
	} else {
		t5.state = 0
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
