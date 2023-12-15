package frameworks

/*
  Personality: Configures the agent's default personality based on OCEAN model.

  A personality is a set of traits that define how an agent behaves. These traits can
  be used to modulate and control how agents make decisions, and how they perceive and
  interact with other agents.
*/

type Personality struct {
	// Five Factor (OCEAN) model
	Openness          float64
	Conscientiousness float64
	Extraversion      float64
	Agreeableness     float64
	Neuroticism       float64
}

func NewDefaultPersonality() *Personality {
	p := &Personality{
		Openness:          0.5,
		Conscientiousness: 0.8, // Hard-working by default
		Extraversion:      0.5,
		Agreeableness:     1, // Cooperative by default
		Neuroticism:       0.5,
	}

	return p
}
