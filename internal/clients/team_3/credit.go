package team_3

type credit struct {
	// score
	credit                   float64
	consecutiveNegativeCount int

	// memory
	_credit float64
}

/*
func (cre *credit) updatecredit(biker SmartAgent, agentsOnBike []objects.IBaseBiker) {
	// update memory
	if *(biker.whether_need_leader(agentsOnBike)) == 0.0 {
		cre._credit = -1.0
	} else {
		cre._credit = 1.0
	}

	// update score
	cre.credit += cre._credit

	if cre._credit == -1.0 {
		cre.consecutiveNegativeCount++
	} else {
		cre.consecutiveNegativeCount = 0
	}
}
*/
