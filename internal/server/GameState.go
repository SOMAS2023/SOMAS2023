package server

import (
	"SOMAS2023/internal/common/objects"
	"github.com/google/uuid"
)

func (s *Server) GetLootBoxes() map[uuid.UUID]*objects.LootBox {
	return s.lootBoxes
}

func (s *Server) GetMegaBikes() map[uuid.UUID]*objects.MegaBike {
	return s.megaBikes
}
