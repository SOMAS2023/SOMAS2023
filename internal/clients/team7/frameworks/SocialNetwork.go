package frameworks

import (
	"SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

type SocialConnection struct {
	connectionAge      int // Number of rounds the agent has been known
	trustLevel         float64
	isActiveConnection bool // is on current bike
}

type SocialConnectionInput struct {
	agentDecisions map[uuid.UUID]utils.Forces
}

type ISocialNetwork[T any] interface {
	GetSocialNetwork() map[uuid.UUID]SocialConnection
	UpdateSocialNetwork(agentIds []uuid.UUID, inputs T)
	UpdateActiveConnections(agentIds []uuid.UUID)
	DeactivateConnections(agentIds []uuid.UUID)
}

type SocialNetwork struct {
	ISocialNetwork[SocialConnectionInput]
	socialNetwork *map[uuid.UUID]SocialConnection
}

func NewSocialNetwork() *SocialNetwork {
	return &SocialNetwork{
		socialNetwork: &map[uuid.UUID]SocialConnection{},
	}
}

func (sn *SocialNetwork) GetSocialNetwork() map[uuid.UUID]SocialConnection {
	return *sn.socialNetwork
}

func updateTrustLevel(currentTrustLevel float64, forces utils.Forces) float64 {
	// TODO: Update trust level based on forces
	return currentTrustLevel
}

func (sn *SocialNetwork) UpdateSocialNetwork(agentIds []uuid.UUID, inputs SocialConnectionInput) {
	for _, agentId := range agentIds {
		connection := (*sn.socialNetwork)[agentId]
		connection.connectionAge += 1
		connection.trustLevel = updateTrustLevel(connection.trustLevel, inputs.agentDecisions[agentId])
		(*sn.socialNetwork)[agentId] = connection
	}
}

func (sn *SocialNetwork) UpdateActiveConnections(agentIds []uuid.UUID) {
	for _, agentId := range agentIds {
		connection := (*sn.socialNetwork)[agentId]
		connection.isActiveConnection = true
		(*sn.socialNetwork)[agentId] = connection
	}
}

func (sn *SocialNetwork) DeactivateConnections(agentIds []uuid.UUID) {
	for _, agentId := range agentIds {
		connection := (*sn.socialNetwork)[agentId]
		connection.isActiveConnection = false
		(*sn.socialNetwork)[agentId] = connection
	}
}
