package server

import (
	"SOMAS2023/internal/common/objects"
	"math/rand"

	"github.com/google/uuid"
)

func (s *Server) GetMegaBikes() map[uuid.UUID]objects.IMegaBike {
	return s.megaBikes
}

func (s *Server) GetLootBoxes() map[uuid.UUID]objects.ILootBox {
	return s.lootBoxes
}

func (s *Server) GetAudi() objects.IAudi {
	return s.audi
}

func (s *Server) SetBikerBike(biker objects.IBaseBiker, bikeId uuid.UUID) {
	// Remove the agent from the old bike, if it was on one
	if oldBikeId, ok := s.megaBikeRiders[biker.GetID()]; ok {
		s.megaBikes[oldBikeId].RemoveAgent(biker.GetID())
	}

	// Add the agent to the new bike
	s.megaBikes[bikeId].AddAgent(biker)
	s.megaBikeRiders[biker.GetID()] = bikeId
	biker.SetBike(bikeId)
}

// get a map of megaBikeIDs mapping to the ids of all Bikers that are trying to join it
func (s *Server) GetJoiningRequests() map[uuid.UUID][]uuid.UUID {
	// iterate over all agents, if their onBike is false add to the map their id in correspondance of that of their desired bike
	bikeRequests := make(map[uuid.UUID][]uuid.UUID)

	for agentID, agent := range s.GetAgentMap() {
		if !agent.GetBikeStatus() {
			bike := agent.GetBike()
			if ids, ok := bikeRequests[bike]; ok {
				bikeRequests[bike] = append(ids, agentID)
			} else {
				bikeRequests[bike] = []uuid.UUID{agentID}
			}
		}
	}
	return bikeRequests
}

// GetRandomBikeId returns the ID of a random bike.
func (s *Server) GetRandomBikeId() uuid.UUID {
	i, targetI := 0, rand.Intn(len(s.megaBikes))
	// Go doesn't have a sensible way to do this...
	for id := range s.megaBikes {
		if i == targetI {
			return id
		}
		i++
	}
	panic("no bikes")
}
