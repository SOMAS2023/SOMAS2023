package server

import (
	"SOMAS2023/internal/common/objects"
	"github.com/google/uuid"
	"math/rand"
)

func (s *Server) GetMegaBikes() map[uuid.UUID]objects.IMegaBike {
	return s.megaBikes
}

func (s *Server) GetLootBoxes() map[uuid.UUID]objects.ILootBox {
	return s.lootBoxes
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
