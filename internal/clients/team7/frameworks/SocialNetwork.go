package frameworks

import (
	"SOMAS2023/internal/common/utils"
	"fmt"

	"github.com/google/uuid"
)

type SocialConnection struct {
	connectionAge      int     // Number of rounds the agent has been known
	trustLevel         float64 // Trust level of the agent, dummy float64 for now.
	isActiveConnection bool    // Boolean indicating if connection is on current bike
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

func (sn *SocialNetwork) updateTrustLevel(connection *SocialConnection, forces utils.Forces) {
	// TODO: Update trust level based on forces
	fmt.Println("SocialNetwork: UpdateTrustLevel called")
	fmt.Println("SocialNetwork: Current trust level: ", (*connection).trustLevel)
}

func (sn *SocialNetwork) UpdateSocialNetwork(agentIds []uuid.UUID, inputs SocialConnectionInput) {
	for _, agentId := range agentIds {
		connection := (*sn.socialNetwork)[agentId]
		connection.connectionAge += 1
		sn.updateTrustLevel(&connection, inputs.agentDecisions[agentId])
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

// Retrieve agents on the current bike
func (sn *SocialNetwork) GetCurrentBikeNetwork() map[uuid.UUID]SocialConnection {
	activeConnections := map[uuid.UUID]SocialConnection{}
	for agentId, connection := range *sn.socialNetwork {
		if connection.isActiveConnection {
			activeConnections[agentId] = connection
		}
	}
	return activeConnections
}
