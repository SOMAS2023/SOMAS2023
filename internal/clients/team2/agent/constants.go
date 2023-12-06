package agent

// Constants related to the calculation of social capital
const (
	ReputationWeight  = 0.6 // Weight for reputation in social capital calculation
	InstitutionWeight = 0.2 // Weight for institution affiliation in social capital calculation
	NetworkWeight     = 0.2 // Weight for network strength in social capital calculation
)

const (
	SocialEventWeight_AgentSentMsg = 0.1

	SocialEventValue_AgentSentMsg = 0.05
)

// Constants related to the calculation of Institution
const (
	InstitutionEventWeight_Adhereance = 0.0 // Weight for rule adhereance in institution calculation
	InstitutionEventWeight_Voting     = 0.0 // Weight for voting in institution calculation
	InstitutionEventWeight_Kickoff    = 0.0 // Weight for being kicked out of bike in institution calculation
	InstitutionEventWeight_Accepted   = 0.0 // Weight for being accepted to bike in institution calculation
	InstitutionEventWeight_VotedRole  = 0.0 // Weight for role assignment in institution calculation

	InstitutionEventValue_Kickoff   = -0.5 // Value for being kicked out of bike in institution calculation
	InstitutionEventValue_VotedRole = 0.0  // Value for role assignment in institution calculation
	InstitutionEventValue_Accepted  = 0.2  // Value for being accepted to bike in institution calculation
)

const (
	forgivenessFactor = 0.5 // Factor used in trustworthiness update calculations
)

// Constants related to decisions.
const (
	ChangeBikeSocialCapitalThreshold = 0.5 // Threshold for deciding whether to change bike or not
)

// weights for each type of governance
const (
	democracyWeight = 0.5
)
