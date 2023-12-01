package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

const bannedFunctionErrorMessage = "you're not allowed to call this"

func (o PhysicsObjectDump) SetPhysicalState(state utils.PhysicalState) {
	panic(bannedFunctionErrorMessage)
}

func (o PhysicsObjectDump) UpdateForce() {
	panic(bannedFunctionErrorMessage)
}

func (o PhysicsObjectDump) UpdateOrientation() {
	panic(bannedFunctionErrorMessage)
}

func (o PhysicsObjectDump) CheckForCollision(objects.IPhysicsObject) bool {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) GetAllMessages([]objects.IBaseBiker) []messaging.IMessage[objects.IBaseBiker] {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) UpdateAgentInternalState() {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) DecideGovernance() voting.GovernanceVote {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) DecideAction() objects.BikerAction {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) DecideForce(uuid.UUID) {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) DecideJoining([]uuid.UUID) map[uuid.UUID]bool {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) ChangeBike() uuid.UUID {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) ProposeDirection() uuid.UUID {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) FinalDirectionVote(map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) DecideAllocation() voting.IdVoteMap {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) VoteForKickout() map[uuid.UUID]int {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) VoteDictator() voting.IdVoteMap {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) VoteLeader() voting.IdVoteMap {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) DictateDirection() uuid.UUID {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) LeadDirection() uuid.UUID {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) SetBike(uuid.UUID) {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) SetForces(utils.Forces) {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) UpdateColour(utils.Colour) {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) UpdatePoints(int) {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) UpdateEnergyLevel(float64) {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) UpdateGameState(objects.IGameState) {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) ToggleOnBike() {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) QueryReputation(uuid.UUID) float64 {
	panic(bannedFunctionErrorMessage)
}

func (a AgentDump) SetReputation(uuid.UUID, float64) {
	panic(bannedFunctionErrorMessage)
}

func (b BikeDump) AddAgent(objects.IBaseBiker) {
	panic(bannedFunctionErrorMessage)
}

func (b BikeDump) RemoveAgent(uuid.UUID) {
	panic(bannedFunctionErrorMessage)
}

func (b BikeDump) UpdateMass() {
	panic(bannedFunctionErrorMessage)
}

func (b BikeDump) KickOutAgent() map[uuid.UUID]int {
	panic(bannedFunctionErrorMessage)
}

func (b BikeDump) SetGovernance(utils.Governance) {
	panic(bannedFunctionErrorMessage)
}

func (b BikeDump) SetRuler(uuid.UUID) {
	panic(bannedFunctionErrorMessage)
}

func (a AudiDump) UpdateGameState(objects.IGameState) {
	panic(bannedFunctionErrorMessage)
}
