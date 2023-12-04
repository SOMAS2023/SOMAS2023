package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"fmt"

	"github.com/google/uuid"
)

func (s *Server) RunRoundLoop() {
	// Capture dump of starting state
	gameState := s.NewGameStateDump()
	s.UpdateGameStates()

	// take care of agents that want to leave the bike and of the acceptance/ expulsion process
	s.RunBikeSwitch(gameState)

	// get the direction decisions and pedalling forces
	s.RunActionProcess()

	// The Audi makes a decision
	s.audi.UpdateGameState(gameState)

	// Move the mega bikes
	for _, bike := range s.megaBikes {
		// update mass dependent on number of agents on bike
		bike.UpdateMass()
		s.MovePhysicsObject(bike)
	}

	// Move the audi
	s.MovePhysicsObject(s.audi)
	// Check Audi collision
	s.AudiCollisionCheck()

	s.UpdateGameStates()

	// Lootbox Distribution
	s.LootboxCheckAndDistributions()

	// Punish bikeless agents
	s.punishBikelessAgents()

	// Check if agents died
	s.unaliveAgents()

	// Replenish objects
	if utils.ReplenishLootBoxes {
		s.replenishLootBoxes()
	}
	if utils.ReplenishMegaBikes {
		s.replenishMegaBikes()
	}
}

func (s *Server) RunBikeSwitch(gameState GameStateDump) {
	inLimbo := make([]uuid.UUID, 0)
	// check if agents want ot leave the bike on this round
	changeBike := s.GetLeavingDecisions(gameState)
	inLimbo = append(inLimbo, changeBike...)
	// update gamestate as it has changed
	s.UpdateGameStates()
	//process the kickout request
	kickedOff := s.HandleKickoutProcess()
	// update gamestate as it has changed
	s.UpdateGameStates()
	inLimbo = append(inLimbo, kickedOff...)
	// process the joining request
	s.ProcessJoiningRequests(inLimbo)
	// update gamestate as it has changed
	s.UpdateGameStates()
}

func (s *Server) HandleKickoutProcess() []uuid.UUID {
	allKicked := make([]uuid.UUID, 0)
	for _, bike := range s.GetMegaBikes() {
		agentsVotes := make([]uuid.UUID, 0)

		// the kickout process only happens democratically in level 0 and level 1
		switch bike.GetGovernance() {
		case utils.Democracy:
			// make map of weights of 1 for all agents on bike
			agents := bike.GetAgents()
			weights := make(map[uuid.UUID]float64)
			for _, agent := range agents {
				weights[agent.GetID()] = 1.0
			}

			// get which agents are getting kicked out
			agentsVotes = bike.KickOutAgent(weights)

		case utils.Leadership:
			// get the map of weights from the leader
			leader := s.GetAgentMap()[bike.GetRuler()]
			weights := leader.DecideWeights(utils.Kickout)
			// get which agents are getting kicked out
			agentsVotes = bike.KickOutAgent(weights)

		case utils.Dictatorship:
			// in level 2 only the ruler can kick out people
			dictator := s.GetAgentMap()[bike.GetRuler()]
			agentsVotes = dictator.DecideKickOut()
		}

		// perform kickout
		leaderKickedOut := false
		allKicked = append(allKicked, agentsVotes...)
		for _, agentID := range agentsVotes {
			fmt.Println("kicking out someone")
			s.RemoveAgentFromBike(s.GetAgentMap()[agentID])
			// if the leader was kicked out vote for a new one
			if agentID == bike.GetRuler() {
				leaderKickedOut = true
			}
		}
		if leaderKickedOut && bike.GetGovernance() == utils.Leadership {
			bike.SetRuler(s.RulerElection(bike.GetAgents(), utils.Leadership))

		}

	}
	return allKicked
}

func (s *Server) GetLeavingDecisions(gameState objects.IGameState) []uuid.UUID {
	leavingAgents := make([]uuid.UUID, 0)
	for agentId, agent := range s.GetAgentMap() {
		fmt.Printf("Agent %s updating state \n", agentId)
		agent.UpdateGameState(gameState)
		agent.UpdateAgentInternalState()
		switch agent.DecideAction() {
		case objects.Pedal:
			continue
		case objects.ChangeBike:
			// decide which bike the agent is going to try and go t
			// the bike id is set to be the desired bike and onbike is set to false
			// so by looking at the values of onBike and megaBikeID it will be known
			// whether the agent is trying to join a bike (and which one)

			// the request is handled at the beginning of the next round, so the moving
			// will only be finalised then
			leavingAgents = append(leavingAgents, agentId)
			s.RemoveAgentFromBike(agent)
		default:
			panic("agent decided invalid action")
		}
	}
	return leavingAgents
}

func (s *Server) ProcessJoiningRequests(inLimbo []uuid.UUID) {

	// -------------------------- PROCESS JOINING REQUESTS -------------------------
	// 1. group agents that have onBike = false by the bike they are trying to join
	bikeRequests := s.GetJoiningRequests(inLimbo)
	// 2. pass to agents on each of the desired bikes a list of all agents trying to join
	for bikeID, pendingAgents := range bikeRequests {
		agents := s.megaBikes[bikeID].GetAgents()
		if len(agents) == 0 {
			for i, pendingAgent := range pendingAgents {
				if i <= utils.BikersOnBike {
					acceptedAgent := s.GetAgentMap()[pendingAgent]
					s.AddAgentToBike(acceptedAgent)
				} else {
					break
				}
			}
		} else {
			bike := s.GetMegaBikes()[bikeID]
			acceptedRanked := make([]uuid.UUID, 0)

			switch bike.GetGovernance() {
			case utils.Democracy:
				// make map of weights of 1 for all agents on bike
				weights := make(map[uuid.UUID]float64)
				for _, agent := range agents {
					weights[agent.GetID()] = 1.0
				}

				// get approval votes from each agent
				responses := make([](map[uuid.UUID]bool), len(agents)) // list containing all the agents' ranking
				for i, agent := range agents {
					responses[i] = agent.DecideJoining(pendingAgents)
				}

				// accept agents based on the response outcome (it will have to be a ranking system, as only 8-n bikers can be accepted)
				acceptedRanked = voting.GetAcceptanceRanking(responses, weights)
			case utils.Leadership:
				// get the map of weights from the leader
				leader := s.GetAgentMap()[bike.GetRuler()]
				weights := leader.DecideWeights(utils.Joining)

				// get approval votes from each agent
				responses := make([](map[uuid.UUID]bool), len(agents)) // list containing all the agents' ranking
				for i, agent := range agents {
					responses[i] = agent.DecideJoining(pendingAgents)
				}

				// accept agents based on the response outcome (it will have to be a ranking system, as only 8-n bikers can be accepted)
				acceptedRanked = voting.GetAcceptanceRanking(responses, weights)
			case utils.Dictatorship:
				dictator := s.GetAgentMap()[bike.GetRuler()]
				acceptedRankedMap := dictator.DecideJoining(pendingAgents)
				for agentID, accepted := range acceptedRankedMap {
					if accepted {
						acceptedRanked = append(acceptedRanked, agentID)
					}
				}
			}

			// run acceptance process
			totalSeatsFilled := len(agents)
			emptySpaces := utils.BikersOnBike - totalSeatsFilled

			for i := 0; i < min(emptySpaces, len(acceptedRanked)); i++ {
				accepted := acceptedRanked[i]
				acceptedAgent := s.GetAgentMap()[accepted]
				s.AddAgentToBike(acceptedAgent)
			}
		}
	}
}

func (s *Server) RunActionProcess() {

	for _, bike := range s.GetMegaBikes() {
		agents := bike.GetAgents()
		if len(agents) == 0 {
			continue
		}

		// get the direction for this round (either the voted on or what's decided by the leader/ dictator)
		// for now it's actually just the elected lootbox (will change to accomodate for other proposal types)
		var direction uuid.UUID
		electedGovernance := bike.GetGovernance()
		switch electedGovernance {
		case utils.Democracy:
			// make map of weights of 1 for all agents on bike
			weights := make(map[uuid.UUID]float64)
			for _, agent := range agents {
				weights[agent.GetID()] = 1.0
			}
			direction = s.RunDemocraticAction(bike, weights)
			for _, agent := range agents {
				agent.UpdateEnergyLevel(-utils.DeliberativeDemocracyPenalty)
			}
		case utils.Leadership:
			// get weights from leader
			leader := s.GetAgentMap()[bike.GetRuler()]
			weights := leader.DecideWeights(utils.Direction)
			direction = s.RunDemocraticAction(bike, weights)
			for _, agent := range agents {
				agent.UpdateEnergyLevel(-utils.LeadershipDemocracyPenalty)
			}
		case utils.Dictatorship:
			direction = s.RunRulerAction(bike)
		}

		for _, agent := range agents {
			agent.DecideForce(direction)
			// deplete energy
			energyLost := agent.GetForces().Pedal * utils.MovingDepletion
			agent.UpdateEnergyLevel(-energyLost)
		}
	}
}

func (s *Server) MovePhysicsObject(po objects.IPhysicsObject) {

	// Server requests to update their force and orientation based on agents pedaling
	po.UpdateForce()
	force := po.GetForce()
	po.UpdateOrientation()
	orientation := po.GetOrientation()
	// Obtains the current xstate (i.e. velocity, acceleration, position, mass)
	initialState := po.GetPhysicalState()

	// Generates a new state based on the force and orientation
	finalState := physics.GenerateNewState(initialState, force, orientation)

	// Sets the new physical state (i.e. updates gamestate)
	po.SetPhysicalState(finalState)
}

func (s *Server) GetWinningDirection(finalVotes map[uuid.UUID]voting.LootboxVoteMap, weights map[uuid.UUID]float64) uuid.UUID {
	// get overall winner direction using chosen voting strategy

	// this allows to get a slice of the interface from that of the specific type
	// this way we can substitute agent.FInalDirectionVote with another function that returns
	// another type of voting type which still implements INormaliseVoteMap
	IfinalVotes := make(map[uuid.UUID]voting.IVoter)
	for i, v := range finalVotes {
		IfinalVotes[i] = v
	}

	// TODO integrate voting functions from group 8
	return voting.WinnerFromDist(IfinalVotes, weights)
}

func (s *Server) AudiCollisionCheck() {
	// Check collision for audi with any megaBike
	for bikeid, megabike := range s.GetMegaBikes() {
		if s.audi.CheckForCollision(megabike) {
			// Collision detected
			fmt.Printf("Collision detected between Audi and MegaBike %s \n", bikeid)
			for _, agentToDelete := range megabike.GetAgents() {
				fmt.Printf("Agent %s killed by Audi \n", agentToDelete.GetID())
				s.RemoveAgent(agentToDelete)
			}
			if utils.AudiRemovesMegaBike {
				fmt.Printf("Megabike %s removed by Audi \n", megabike.GetID())
				delete(s.megaBikes, megabike.GetID())
			}
		}
	}
}

func (s *Server) LootboxCheckAndDistributions() {

	// checks how many bikes have looted one lootbox to split it between them
	looted := make(map[uuid.UUID]int)
	for _, megabike := range s.GetMegaBikes() {
		for lootid, lootbox := range s.GetLootBoxes() {
			if megabike.CheckForCollision(lootbox) {
				if value, ok := looted[lootid]; ok {
					looted[lootid] = value + 1
				} else {
					looted[lootid] = 1
				}
			}
		}
	}
	for bikeid, megabike := range s.GetMegaBikes() {
		for lootid, lootbox := range s.GetLootBoxes() {
			if megabike.CheckForCollision(lootbox) {
				// Collision detected
				fmt.Printf("Collision detected between MegaBike %s and LootBox %s \n", bikeid, lootid)
				agents := megabike.GetAgents()
				totAgents := len(agents)

				if totAgents > 0 {
					fmt.Printf("Total agents: %d \n", totAgents)
					gov := s.GetMegaBikes()[bikeid].GetGovernance()
					var winningAllocation voting.IdVoteMap
					switch gov {
					case utils.Democracy:
						allAllocations := make(map[uuid.UUID]voting.IdVoteMap)
						for _, agent := range agents {
							// the agents return their ideal lootbox split by assigning a number between 0 and 1 to
							// each biker on their bike (including themselves)
							allAllocations[agent.GetID()] = agent.DecideAllocation()
						}

						Iallocations := make(map[uuid.UUID]voting.IVoter)
						for i, v := range allAllocations {
							Iallocations[i] = v
						}
						// TODO handle error
						// make weights of 1 for all agents
						weights := make(map[uuid.UUID]float64)
						for _, agent := range agents {
							weights[agent.GetID()] = 1.0
						}
						winningAllocation, _ = voting.CumulativeDist(Iallocations, weights)
					case utils.Leadership:
						// get the map of weights from the leader
						leader := s.GetAgentMap()[megabike.GetRuler()]
						weights := leader.DecideWeights(utils.Allocation)
						// get allocation votes from each agent
						allAllocations := make(map[uuid.UUID]voting.IdVoteMap)
						for _, agent := range agents {
							allAllocations[agent.GetID()] = agent.DecideAllocation()
						}
						Iallocations := make(map[uuid.UUID]voting.IVoter)
						for i, v := range allAllocations {
							Iallocations[i] = v
						}
						winningAllocation, _ = voting.CumulativeDist(Iallocations, weights)
					case utils.Dictatorship:
						// dictator decides the allocation
						leader := s.GetAgentMap()[megabike.GetRuler()]
						winningAllocation = leader.DecideDictatorAllocation()
					}

					bikeShare := float64(looted[lootid]) // how many other bikes have looted this box

					for agentID, allocation := range winningAllocation {
						fmt.Printf("total loot: %f \n", lootbox.GetTotalResources())
						lootShare := allocation * (lootbox.GetTotalResources() / bikeShare)
						agent := s.GetAgentMap()[agentID]
						// Allocate loot based on the calculated utility share
						fmt.Printf("Agent %s allocated %f loot \n", agent.GetID(), lootShare)
						agent.UpdateEnergyLevel(lootShare)
						// Allocate points if the box is of the right colour
						if agent.GetColour() == lootbox.GetColour() {
							agent.UpdatePoints(utils.PointsFromSameColouredLootBox)
						}
					}
				}
			}
		}
	}

	// despawn lootboxes that have been looted
	for id, loot := range looted {
		if loot > 0 {
			delete(s.lootBoxes, id)
		}
	}
}

func (s *Server) unaliveAgents() {
	for id, agent := range s.GetAgentMap() {
		if agent.GetEnergyLevel() < 0 {
			fmt.Printf("Agent %s got game ended\n", id)
			s.RemoveAgent(agent)
		}
	}
}

func (s *Server) punishBikelessAgents() {
	for id, agent := range s.GetAgentMap() {
		if _, ok := s.megaBikeRiders[id]; !ok {
			// Agent is not on a bike
			agent.UpdateEnergyLevel(utils.LimboEnergyPenalty)
		}
	}
}
