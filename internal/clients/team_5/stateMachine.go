package team5Agent

// 0 = normal, 1 = conservative

func (t5 *team5Agent) UpdateAgentInternalState() {
	// Update agent's state based on energy level or other parameters
	t5.updateState()
	t5.updateReputationOfAllAgents()
}

func (t5 *team5Agent) updateState() {
	// Example condition: switch to conservative if energy is low
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
