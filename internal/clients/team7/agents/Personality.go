package agents

/*
  A personality is a set of traits that define how an agent behaves. These traits can
  be used to modulate and control how agents make decisions, and how they perceive and
  interact with other agents.
*/
type Personality struct {
	SelfConfidence    float64 // Measure of self-confidence
	Compassion        float64 // Measure of compassion which modulates decisions
	PositiveTrustStep float64 // Rate at which trust is gained
	NegativeTrustStep float64 // Rate at which trust is lost
	Trustworthiness   float64 // Measure of trustworthiness to control truth and lies
}

func NewDefaultPersonality() *Personality {
	return &Personality{
		SelfConfidence:    1,   // Confident in own decisions
		Compassion:        0.5, // Neutral compassion
		PositiveTrustStep: 0.1, // Trust is gained slowly
		NegativeTrustStep: 0.1, // Trust is lost slowly
		Trustworthiness:   1,   // Will never lie
	}
}
