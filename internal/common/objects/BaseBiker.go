package objects

import (
	utils "SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"
	"math"

	"math/rand"

	baseAgent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

type IBaseBiker interface {
	baseAgent.IAgent[IBaseBiker]

	DecideGovernance() utils.Governance
	DecideAction() BikerAction                                                  // ** determines what action the agent is going to take this round. (changeBike or Pedal)
	DecideForce(direction uuid.UUID)                                            // ** defines the vector you pass to the bike: [pedal, brake, turning]
	DecideJoining(pendinAgents []uuid.UUID) map[uuid.UUID]bool                  // ** decide whether to accept or not accept bikers, ranks the ones
	ChangeBike() uuid.UUID                                                      // ** called when biker wants to change bike, it will choose which bike to try and join
	ProposeDirection() uuid.UUID                                                // ** returns the id of the desired lootbox based on internal strategy
	FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap // ** stage 3 of direction voting
	DecideAllocation() voting.IdVoteMap                                         // ** decide the allocation parameters
	VoteForKickout() map[uuid.UUID]int
	VoteDictator() voting.IdVoteMap
	VoteLeader() voting.IdVoteMap

	// dictator functions
	DictateDirection() uuid.UUID                // ** called only when the agent is the dictator
	DecideKickOut() []uuid.UUID                 // ** decide which agents to kick out (dictator)
	DecideDictatorAllocation() voting.IdVoteMap // ** decide the allocation (dictator)

	// leader functions
	DecideWeights(action utils.Action) map[uuid.UUID]float64 // decide on weights for various actions

	GetForces() utils.Forces        // returns forces for current round
	GetColour() utils.Colour        // returns the colour of the lootbox that the agent is currently seeking
	GetLocation() utils.Coordinates // gets the agent's location
	GetBike() uuid.UUID             // tells the biker which bike it is on
	GetEnergyLevel() float64        // returns the energy level of the agent
	GetPoints() int
	GetBikeStatus() bool // returns whether the biker is on a bike or not

	SetBike(uuid.UUID)                     // sets the megaBikeID. this is either the id of the bike that the agent is on or the one that it's trying to join
	SetForces(forces utils.Forces)         // sets the forces (to be updated in DecideForces())
	UpdateColour(totColours utils.Colour)  // called if a box of the desired colour has been looted
	UpdatePoints(pointGained int)          // called by server
	UpdateEnergyLevel(energyLevel float64) // increase the energy level of the agent by the allocated lootbox share or decrease by expended energy
	UpdateGameState(gameState IGameState)  // sets the gameState field at the beginning of each round
	ToggleOnBike()                         // called when removing or adding a biker on a bike
	ResetPoints()

	GetReputation() map[uuid.UUID]float64 // get reputation value of all other agents
	QueryReputation(uuid.UUID) float64    // query for reputation value of specific agent with UUID
	SetReputation(uuid.UUID, float64)     // set reputation value of specific agent with UUID

	HandleKickoutMessage(msg KickoutAgentMessage)
	HandleReputationMessage(msg ReputationOfAgentMessage)
	HandleJoiningMessage(msg JoiningAgentMessage)
	HandleLootboxMessage(msg LootboxMessage)
	HandleGovernanceMessage(msg GovernanceMessage)
	HandleForcesMessage(msg ForcesMessage)
	HandleVoteGovernanceMessage(msg VoteGoveranceMessage)
	HandleVoteLootboxDirectionMessage(msg VoteLootboxDirectionMessage)
	HandleVoteRulerMessage(msg VoteRulerMessage)
	HandleVoteKickoutMessage(msg VoteKickoutMessage)

	GetAllMessages([]IBaseBiker) []messaging.IMessage[IBaseBiker]
}

type BikerAction int

const (
	Pedal BikerAction = iota
	ChangeBike
)

type BaseBiker struct {
	*baseAgent.BaseAgent[IBaseBiker]              // BaseBiker inherits functions from BaseAgent such as GetID(), GetAllMessages() and UpdateAgentInternalState()
	soughtColour                     utils.Colour // the colour of the lootbox that the agent is currently seeking
	onBike                           bool
	energyLevel                      float64 // float between 0 and 1
	points                           int
	forces                           utils.Forces
	megaBikeId                       uuid.UUID             // if they are not on a bike it will be 0
	gameState                        IGameState            // updated by the server at every round
	reputation                       map[uuid.UUID]float64 // record reputation for other agents in float
}

func (bb *BaseBiker) GetEnergyLevel() float64 {
	return bb.energyLevel
}

func (bb *BaseBiker) GetPoints() int {
	return bb.points
}

// the function will be called by the server to:
// - reduce the energy level based on the force spent pedalling (energyLevel will be neg.ve)
// - increase the energy level after a lootbox has been looted (energyLevel will be pos.ve)
func (bb *BaseBiker) UpdateEnergyLevel(energyLevel float64) {
	bb.energyLevel += energyLevel
	if bb.energyLevel > 1.0 {
		bb.energyLevel = 1.0
	}
}

func (bb *BaseBiker) GetColour() utils.Colour {
	return bb.soughtColour
}

// through this function the agent submits their desired allocation of resources
// in the MVP each agent returns 1 whcih will cause the distribution to be equal across all of them
func (bb *BaseBiker) DecideAllocation() voting.IdVoteMap {
	bikeID := bb.GetBike()
	fellowBikers := bb.gameState.GetMegaBikes()[bikeID].GetAgents()
	distribution := make(voting.IdVoteMap)
	for _, agent := range fellowBikers {
		if agent.GetID() == bb.GetID() {
			distribution[agent.GetID()] = 1.0
		} else {
			distribution[agent.GetID()] = 0.0
		}
	}
	return distribution
}

// the biker itself doesn't technically have a location (as it's on the map only when it's on a bike)
// in fact this function is only called when the biker needs to make a decision about the pedaling forces
func (bb *BaseBiker) GetLocation() utils.Coordinates {
	megaBikes := bb.gameState.GetMegaBikes()
	return megaBikes[bb.megaBikeId].GetPosition()
}

// returns the nearest lootbox with respect to the agent's bike current position
// in the MVP this is used to determine the pedalling forces as all agent will be
// aiming to get to the closest lootbox by default
func (bb *BaseBiker) nearestLoot() uuid.UUID {
	currLocation := bb.GetLocation()
	shortestDist := math.MaxFloat64
	var nearestBox uuid.UUID
	var currDist float64
	for _, loot := range bb.gameState.GetLootBoxes() {
		x, y := loot.GetPosition().X, loot.GetPosition().Y
		currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
		if currDist < shortestDist {
			nearestBox = loot.GetID()
			shortestDist = currDist
		}
	}
	return nearestBox
}

// in the MVP the biker's action defaults to pedaling (as it won't be able to change bikes)
// in future implementations this function will be overridden by the agent's specific strategy
// which will be used to determine whether to pedalor try to change bike
func (bb *BaseBiker) DecideAction() BikerAction {
	return Pedal
}

// determine the forces (pedalling, breaking and turning)
// in the MVP the pedalling force will be 1, the breaking 0 and the tunring is determined by the
// location of the nearest lootbox

// the function is passed in the id of the voted lootbox and the default base bikers steer to that lootbox.
func (bb *BaseBiker) DecideForce(direction uuid.UUID) {

	// NEAREST BOX STRATEGY (MVP)
	currLocation := bb.GetLocation()
	currentLootBoxes := bb.gameState.GetLootBoxes()

	// Check if there are lootboxes available and move towards closest one
	if len(currentLootBoxes) > 0 {
		targetPos := currentLootBoxes[direction].GetPosition()

		deltaX := targetPos.X - currLocation.X
		deltaY := targetPos.Y - currLocation.Y
		angle := math.Atan2(deltaY, deltaX)
		normalisedAngle := angle / math.Pi

		// Default BaseBiker will always
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: normalisedAngle - bb.gameState.GetMegaBikes()[bb.megaBikeId].GetOrientation(),
		}

		nearestBoxForces := utils.Forces{
			Pedal:   utils.BikerMaxForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		bb.SetForces(nearestBoxForces)
	} else { // otherwise move away from audi
		audiPos := bb.GetGameState().GetAudi().GetPosition()

		deltaX := audiPos.X - currLocation.X
		deltaY := audiPos.Y - currLocation.Y

		// Steer in opposite direction to audi
		angle := math.Atan2(deltaY, deltaX)
		normalisedAngle := angle / math.Pi

		// Steer in opposite direction to audi
		var flipAngle float64
		if normalisedAngle < 0.0 {
			flipAngle = normalisedAngle + 1.0
		} else if normalisedAngle > 0.0 {
			flipAngle = normalisedAngle - 1.0
		}

		// Default BaseBiker will always
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: flipAngle - bb.gameState.GetMegaBikes()[bb.megaBikeId].GetOrientation(),
		}

		escapeAudiForces := utils.Forces{
			Pedal:   utils.BikerMaxForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		bb.SetForces(escapeAudiForces)
	}
}

// decide which bike to go to. the base agent chooses a random bike
func (bb *BaseBiker) ChangeBike() uuid.UUID {
	megaBikes := bb.gameState.GetMegaBikes()
	i, targetI := 0, rand.Intn(len(megaBikes))
	// Go doesn't have a sensible way to do this...
	for id := range megaBikes {
		if i == targetI {
			return id
		}
		i++
	}
	panic("no bikes")
}

func (bb *BaseBiker) SetBike(bikeId uuid.UUID) {
	bb.megaBikeId = bikeId
}

func (bb *BaseBiker) GetBike() uuid.UUID {
	return bb.megaBikeId
}

// this is called when a lootbox of the desidered colour has been looted in order to update the sought colour
func (bb *BaseBiker) UpdateColour(totColours utils.Colour) {
	bb.soughtColour = utils.Colour(rand.Intn(int(totColours)))
}

// update the points at the end of a round
func (bb *BaseBiker) UpdatePoints(pointsGained int) {
	bb.points += pointsGained
}

func (bb *BaseBiker) GetForces() utils.Forces {
	return bb.forces
}

func (bb *BaseBiker) SetForces(forces utils.Forces) {
	bb.forces = forces
}

func (bb *BaseBiker) UpdateGameState(gameState IGameState) {
	bb.gameState = gameState
}

// default implementation returns the id of the nearest lootbox
func (bb *BaseBiker) ProposeDirection() uuid.UUID {
	return bb.nearestLoot()
}

func (bb *BaseBiker) ToggleOnBike() {
	bb.onBike = !bb.onBike
}

func (bb *BaseBiker) GetBikeStatus() bool {
	return bb.onBike
}

func (bb *BaseBiker) GetGameState() IGameState {
	return bb.gameState
}

// Returns the other agents on your bike :)
func (bb *BaseBiker) GetFellowBikers() []IBaseBiker {
	bikes := bb.gameState.GetMegaBikes()
	if _, ok := bikes[bb.GetBike()]; !ok {
		return []IBaseBiker{}
	}
	bike := bikes[bb.GetBike()]
	fellowBikers := bike.GetAgents()
	return fellowBikers
}

// GetReputation map from agent, need to check if nil when call this function
func (bb *BaseBiker) GetReputation() map[uuid.UUID]float64 {
	return bb.reputation
}

// QueryReputation of specific agent with given ID, if there is no record for given agentID then return 0
func (bb *BaseBiker) QueryReputation(agentId uuid.UUID) float64 {
	if bb.reputation == nil {
		return 0
	}
	return bb.reputation[agentId]
}

func (bb *BaseBiker) SetReputation(agentId uuid.UUID, reputation float64) {
	if bb.reputation == nil {
		bb.reputation = make(map[uuid.UUID]float64)
	}
	bb.reputation[agentId] = reputation
}

// an agent will have to rank the agents that are trying to join and that they will try to
func (bb *BaseBiker) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	decision := make(map[uuid.UUID]bool)
	for _, agent := range pendingAgents {
		decision[agent] = true
	}
	return decision
}

func (bb *BaseBiker) DecideGovernance() utils.Governance {
	// Change behaviour here to return different governance
	return utils.Democracy
}

func (bb *BaseBiker) ResetPoints() {
	bb.points = 0
}

// this function will contain the agent's strategy on deciding which direction to go to
// the default implementation returns an equal distribution over all options
// this will also be tried as returning a rank of options
func (bb *BaseBiker) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	votes := make(voting.LootboxVoteMap)
	totOptions := len(proposals)
	normalDist := 1.0 / float64(totOptions)
	for _, proposal := range proposals {
		if val, ok := votes[proposal]; ok {
			votes[proposal] = val + normalDist
		} else {
			votes[proposal] = normalDist
		}
	}
	return votes
}

func (bb *BaseBiker) VoteForKickout() map[uuid.UUID]int {
	voteResults := make(map[uuid.UUID]int)
	bikeID := bb.GetBike()

	fellowBikers := bb.gameState.GetMegaBikes()[bikeID].GetAgents()
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		if agentID != bb.GetID() {
			// random votes to other agents
			voteResults[agentID] = rand.Intn(2) // randomly assigns 0 or 1 vote
		}
	}

	return voteResults
}

// defaults to voting for first agent in the list
func (bb *BaseBiker) VoteDictator() voting.IdVoteMap {
	votes := make(voting.IdVoteMap)
	fellowBikers := bb.GetFellowBikers()
	for i, fellowBiker := range fellowBikers {
		if i == 0 {
			votes[fellowBiker.GetID()] = 1.0
		} else {
			votes[fellowBiker.GetID()] = 0.0
		}
	}
	return votes
}

func (bb *BaseBiker) DictateDirection() uuid.UUID {
	nearest := bb.nearestLoot()
	return nearest
}

// defaults to voting for first agent in the list
func (bb *BaseBiker) VoteLeader() voting.IdVoteMap {
	votes := make(voting.IdVoteMap)
	fellowBikers := bb.GetFellowBikers()
	for i, fellowBiker := range fellowBikers {
		if i == 0 {
			votes[fellowBiker.GetID()] = 1.0
		} else {
			votes[fellowBiker.GetID()] = 0.0
		}
	}
	return votes
}

// defaults to an equal distribution over all agents for all actions
func (bb *BaseBiker) DecideWeights(action utils.Action) map[uuid.UUID]float64 {
	weights := make(map[uuid.UUID]float64)
	agents := bb.GetFellowBikers()
	for _, agent := range agents {
		weights[agent.GetID()] = 1.0
	}
	return weights
}

// only called when the agent is the dictator
func (bb *BaseBiker) DecideKickOut() []uuid.UUID {
	return (make([]uuid.UUID, 0))
}

// only called when the agent is the dictator
func (bb *BaseBiker) DecideDictatorAllocation() voting.IdVoteMap {
	bikeID := bb.GetBike()
	fellowBikers := bb.gameState.GetMegaBikes()[bikeID].GetAgents()
	distribution := make(voting.IdVoteMap)
	equalDist := 1.0 / float64(len(fellowBikers))
	for _, agent := range fellowBikers {
		distribution[agent.GetID()] = equalDist
	}
	return distribution
}

// This function updates all the messages for that agent i.e. both sending and receiving.
// And returns the new messages from other agents to your agent
func (bb *BaseBiker) GetAllMessages([]IBaseBiker) []messaging.IMessage[IBaseBiker] {
	// For team's agent add your own logic on chosing when your biker should send messages and which ones to send (return)
	wantToSendMsg := false
	if wantToSendMsg {
		reputationMsg := bb.CreateReputationMessage()
		kickoutMsg := bb.CreatekickoutMessage()
		lootboxMsg := bb.CreateLootboxMessage()
		joiningMsg := bb.CreateJoiningMessage()
		governceMsg := bb.CreateGoverenceMessage()
		forcesMsg := bb.CreateForcesMessage()
		voteGoveranceMessage := bb.CreateVoteGovernanceMessage()
		voteLootboxDirectionMessage := bb.CreateVoteLootboxDirectionMessage()
		voteRulerMessage := bb.CreateVoteRulerMessage()
		voteKickoutMessage := bb.CreateVotekickoutMessage()
		return []messaging.IMessage[IBaseBiker]{reputationMsg, kickoutMsg, lootboxMsg, joiningMsg, governceMsg, forcesMsg, voteGoveranceMessage, voteLootboxDirectionMessage, voteRulerMessage, voteKickoutMessage}
	}
	return []messaging.IMessage[IBaseBiker]{}
}

func (bb *BaseBiker) CreatekickoutMessage() KickoutAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return KickoutAgentMessage{
		BaseMessage: messaging.CreateMessage[IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     uuid.Nil,
		Kickout:     false,
	}
}

func (bb *BaseBiker) CreateReputationMessage() ReputationOfAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return ReputationOfAgentMessage{
		BaseMessage: messaging.CreateMessage[IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     uuid.Nil,
		Reputation:  1.0,
	}
}

func (bb *BaseBiker) CreateJoiningMessage() JoiningAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return JoiningAgentMessage{
		BaseMessage: messaging.CreateMessage[IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     uuid.Nil,
		BikeId:      uuid.Nil,
	}
}
func (bb *BaseBiker) CreateLootboxMessage() LootboxMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return LootboxMessage{
		BaseMessage: messaging.CreateMessage[IBaseBiker](bb, bb.GetFellowBikers()),
		LootboxId:   uuid.Nil,
	}
}

func (bb *BaseBiker) CreateGoverenceMessage() GovernanceMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return GovernanceMessage{
		BaseMessage:  messaging.CreateMessage[IBaseBiker](bb, bb.GetFellowBikers()),
		BikeId:       uuid.Nil,
		GovernanceId: 0,
	}
}

func (bb *BaseBiker) CreateForcesMessage() ForcesMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return ForcesMessage{
		BaseMessage: messaging.CreateMessage[IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     uuid.Nil,
		AgentForces: utils.Forces{
			Pedal: 0.0,
			Brake: 0.0,
			Turning: utils.TurningDecision{
				SteerBike:     false,
				SteeringForce: 0.0,
			},
		},
	}
}

func (bb *BaseBiker) CreateVoteGovernanceMessage() VoteGoveranceMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return VoteGoveranceMessage{
		BaseMessage: messaging.CreateMessage[IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     make(voting.IdVoteMap),
	}
}

func (bb *BaseBiker) CreateVoteLootboxDirectionMessage() VoteLootboxDirectionMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return VoteLootboxDirectionMessage{
		BaseMessage: messaging.CreateMessage[IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     make(voting.IdVoteMap),
	}
}

func (bb *BaseBiker) CreateVoteRulerMessage() VoteRulerMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return VoteRulerMessage{
		BaseMessage: messaging.CreateMessage[IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     make(voting.IdVoteMap),
	}
}

func (bb *BaseBiker) CreateVotekickoutMessage() VoteKickoutMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return VoteKickoutMessage{
		BaseMessage: messaging.CreateMessage[IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     make(map[uuid.UUID]int),
	}
}

func (bb *BaseBiker) HandleKickoutMessage(msg KickoutAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// kickout := msg.Kickout
}

func (bb *BaseBiker) HandleReputationMessage(msg ReputationOfAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// reputation := msg.Reputation
}

func (bb *BaseBiker) HandleJoiningMessage(msg JoiningAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// bikeId := msg.BikeId
}

func (bb *BaseBiker) HandleLootboxMessage(msg LootboxMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// lootboxId := msg.LootboxId
}

func (bb *BaseBiker) HandleGovernanceMessage(msg GovernanceMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// bikeId := msg.BikeId
	// governanceId := msg.GovernanceId
}

func (bb *BaseBiker) HandleForcesMessage(msg ForcesMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// agentForces := msg.AgentForces

}

func (bb *BaseBiker) HandleVoteGovernanceMessage(msg VoteGoveranceMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap
}

func (bb *BaseBiker) HandleVoteLootboxDirectionMessage(msg VoteLootboxDirectionMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap
}

func (bb *BaseBiker) HandleVoteRulerMessage(msg VoteRulerMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap
}

func (bb *BaseBiker) HandleVoteKickoutMessage(msg VoteKickoutMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap
}

// this function is going to be called by the server to instantiate bikers in the MVP
func GetIBaseBiker(totColours utils.Colour, bikeId uuid.UUID) IBaseBiker {
	return &BaseBiker{
		BaseAgent:    baseAgent.NewBaseAgent[IBaseBiker](),
		soughtColour: utils.GenerateRandomColour(),
		onBike:       true,
		energyLevel:  1.0,
		points:       0,
	}
}

// this function will be used by GetTeamAgent to get the ref to the BaseBiker
func GetBaseBiker(totColours utils.Colour, bikeId uuid.UUID) *BaseBiker {
	return &BaseBiker{
		BaseAgent:    baseAgent.NewBaseAgent[IBaseBiker](),
		soughtColour: utils.GenerateRandomColour(),
		onBike:       true,
		energyLevel:  1.0,
		points:       0,
	}
}
