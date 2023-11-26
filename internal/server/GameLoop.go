package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"fmt"

	"github.com/google/uuid"
)

func (s *Server) RunGameLoop() {

	// take care of agents that want to leave the bike and of the acceptance/ expulsion process
	s.RunBikeSwitch()

	// get the direction decisions and pedalling forces
	s.RunActionProcess()

	// The Audi makes a decision
	s.audi.UpdateGameState(s)

	// Move the mega bikes
	for _, bike := range s.GetMegaBikes() {
		// update mass dependent on number of agents on bike
		bike.UpdateMass()
		s.MovePhysicsObject(bike)
	}

	// Move the audi
	s.MovePhysicsObject(s.audi)

	// Lootbox Distribution
	s.LootboxCheckAndDistributions()

	// Punish bikeless agents
	s.punishBikelessAgents()

	// Check if agents died
	s.unaliveAgents()

	// Replenish objects
	s.replenishLootBoxes()
	s.replenishMegaBikes()
}

func (s *Server) RunBikeSwitch() {
	// check if agents want ot leave the bike on this round
	s.GetLeavingDecisions()
	// process joining requests from last round
	s.ProcessJoiningRequests()
}

func (s *Server) GetLeavingDecisions() {
	for agentId, agent := range s.GetAgentMap() {
		fmt.Printf("Agent %s updating state \n", agentId)
		agent.UpdateGameState(s)
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
			agent.SetBike(agent.ChangeBike())
			agent.ToggleOnBike()

			// the biker needs to be removed from the current bike as well
			// it will be added to the desired one (if accepted) at the beginning of next loop
			if oldBikeId, ok := s.megaBikeRiders[agent.GetID()]; ok {
				s.megaBikes[oldBikeId].RemoveAgent(agent.GetID())
				delete(s.megaBikeRiders, agent.GetID())
			}
		default:
			panic("agent decided invalid action")
		}
	}
}

func (s *Server) ProcessJoiningRequests() {

	// -------------------------- PROCESS JOINING REQUESTS -------------------------
	// 1. group agents that have onBike = false by the bike they are trying to join
	bikeRequests := s.GetJoiningRequests()
	// 2. pass to agents on each of the desired bikes a list of all agents trying to join
	for bike, pendingAgents := range bikeRequests {
		agents := s.megaBikes[bike].GetAgents()

		responses := make([](map[uuid.UUID]bool), len(agents)) // list containing all the agents' ranking
		for _, agent := range agents {
			responses = append(responses, agent.DecideJoining(pendingAgents))
		}
		// 3. accept agents based on the response outcome (it will have to be a ranking system, as only 8-n bikers can be accepted)
		acceptedRanked := voting.GetAcceptanceRanking(responses)
		totalSeatsFilled := len(s.megaBikes[bike].GetAgents())
		emptySpaces := 8 - totalSeatsFilled
		for _, accepted := range acceptedRanked[:emptySpaces] {
			s.GetAgentMap()[accepted].ToggleOnBike()
			s.SetBikerBike(s.GetAgentMap()[accepted], bike)
		}
	}
}

func (s *Server) GetDirectionProposals(agent objects.IBaseBiker, proposedDirections map[uuid.UUID][]uuid.UUID) {
	// --------------------- VOTING ROUTINE - STEP 1 --------------------------
	// pitch proposal (desired lootbox)
	bike := agent.GetBike()
	if ids, ok := proposedDirections[bike]; ok {
		proposedDirections[bike] = append(ids, agent.ProposeDirection())
	} else {
		proposedDirections[bike] = []uuid.UUID{agent.ProposeDirection()}
	}
}

func (s *Server) RunActionProcess() {
	// map of the proposed lootboxes by bike (for each bike a list of lootbox proposals is made, with one lootbox proposed by each agent on the bike)
	proposedDirections := make(map[uuid.UUID][]uuid.UUID)
	for _, agent := range s.GetAgentMap() {
		// agents that have decided to stay on the bike (and that haven't been kicked off it)
		// will participate in the voting for the directions
		// ---------------------------VOTING ROUTINE - STEP 1 ---------------------
		if agent.GetBikeStatus() {
			s.GetDirectionProposals(agent, proposedDirections)
		}
	}

	// pass the pitched directions of a bike to all agents on that bike and get their final vote
	for bikeID, proposals := range proposedDirections {
		// ---------------------------VOTING ROUTINE - STEP 2 ---------------------
		finalVotes := s.GetProposalsDist(bikeID, proposals)

		// ---------------------------VOTING ROUTINE - STEP 3 --------------
		direction := s.GetWinningDirection(finalVotes)

		// get the force given the chosen voting strategy
		for _, agent := range s.megaBikes[bikeID].GetAgents() {
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

func (s *Server) GetProposalsDist(bikeID uuid.UUID, proposals []uuid.UUID) []voting.LootboxVoteMap {
	// ----------------------------- VOTING ROUTINE - STEP 2 -----------------
	// get second vote given everyone's proposal
	// the finalVote can either be a ranking of proposed directions or a map from proposal to vote (between 0,1)
	// we will try implementing both, the infrastructure should be the same
	agentsOnBike := s.megaBikes[bikeID].GetAgents()
	// server collates all vote distributions from each agent into a list of final votes
	finalVotes := make([]voting.LootboxVoteMap, len(agentsOnBike))

	for _, agent := range s.megaBikes[bikeID].GetAgents() {
		finalVotes = append(finalVotes, agent.FinalDirectionVote(proposals))
	}
	return finalVotes
}

func (s *Server) GetWinningDirection(finalVotes []voting.LootboxVoteMap) uuid.UUID {
	// get overall winner direction using chosen voting strategy

	// this allows to get a slice of the interface from that of the specific type
	// this way we can substitute agent.FInalDirectionVote with another function that returns
	// another type of voting type which still implements INormaliseVoteMap
	IfinalVotes := make([]voting.IVoter, len(finalVotes))
	for i, v := range finalVotes {
		IfinalVotes[i] = v
	}

	return voting.WinnerFromDist(IfinalVotes)
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
					allAllocations := make([]voting.IdVoteMap, totAgents)
					for _, agent := range agents {
						// the agents return their ideal lootbox split by assigning a number between 0 and 1 to
						// each biker on their bike (including themselves)
						allAllocations = append(allAllocations, agent.DecideAllocation())
					}
					if totAgents <= len(allAllocations) {
						allAllocations = allAllocations[totAgents:]
					}

					Iallocations := make([]voting.IVoter, len(allAllocations))
					for i, v := range allAllocations {
						Iallocations[i] = v
					}
					// TODO handle error
					winningAllocation, _ := voting.CumulativeDist(Iallocations)
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
			if bikeId, ok := s.megaBikeRiders[id]; ok {
				s.megaBikes[bikeId].RemoveAgent(id)
				delete(s.megaBikeRiders, id)
			}
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

func (s *Server) Start() {
	fmt.Printf("Server initialised with %d agents \n\n", len(s.GetAgentMap()))
	gameStates := make([]GameStateDump, 0, s.GetIterations())
	for i := 0; i < s.GetIterations(); i++ {
		fmt.Printf("Game Loop %d running... \n \n", i)
		fmt.Printf("Main game loop running...\n\n")
		s.RunGameLoop()
		gameStates = append(gameStates, s.NewGameStateDump())
		fmt.Printf("\nMain game loop finished.\n\n")
		fmt.Printf("Messaging session started...\n\n")
		s.RunMessagingSession()
		fmt.Printf("\nMessaging session completed\n\n")
		fmt.Printf("Game Loop %d completed.\n", i)
	}
	// file, err := os.Create("game_dump.json")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// if err := json.NewEncoder(file).Encode(gameStates); err != nil {
	// 	panic(err)
	// }
}
