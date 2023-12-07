package modules

// Constants related to the calculation of social capital
const (
	ReputationWeight  = 1.0 // Weight for trust in social capital calculation
	InstitutionWeight = 0.0 // Weight for institution affiliation in social capital calculation
	NetworkWeight     = 1.0 // Weight for network strength in social capital calculation
)

const (
	SocialEventWeight_AgentSentMsg = 1

	SocialEventValue_AgentSentMsg = 1
)

// Constants related to the calculation of Institution
const (
	InstitutionEventWeight_Adhereance = 0.0  // Weight for rule adhereance in institution calculation
	InstitutionEventWeight_Voting     = 0.0  // Weight for voting in institution calculation
	InstitutionEventWeight_KickedOut  = 0.0  // Weight for being kicked out of bike in institution calculation
	InstitutionEventWeight_Accepted   = 0.0  // Weight for being accepted to bike in institution calculation
	InstitutionEventWeight_VotedRole  = 0.0  // Weight for role assignment in institution calculation
	InstitutionKickoffEventValue      = -0.5 // Value for being kicked out of bike in institution calculation
	InstitutionAcceptedEventValue     = 0.2  // Value for being accepted to bike in institution calculation
)

const (
	forgivenessFactor = 0.5 // Factor used in trustworthiness update calculations
)
