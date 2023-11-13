package server

import (
	"SOMAS2023/internal/common/objects"
	"github.com/google/uuid"
)

func (s *Server) GetMegaBikes() map[uuid.UUID]objects.IMegaBike {
	return s.megaBikes
}

func (s *Server) GetLootBoxes() map[uuid.UUID]objects.ILootBox {
	return s.lootBoxes
}
