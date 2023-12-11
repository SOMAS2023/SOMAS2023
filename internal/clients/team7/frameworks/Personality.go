package frameworks

/*
  A personality is a set of traits that define how an agent behaves. These traits can
  be used to modulate and control how agents make decisions, and how they perceive and
  interact with other agents.
*/
import (
	"math/rand"
	"time"
)

type Personality struct {
	SelfConfidence    float64
	Compassion        float64
	PositiveTrustStep float64
	NegativeTrustStep float64
	Trustworthiness   float64
	// The following four should add up to 1
	Egalitarian float64
	Selfish     float64
	Judgemental float64
	Utilitarian float64

	// Five Factor (OCEAN) model
	Openness          float64
	Conscientiousness float64
	Extraversion      float64
	Agreeableness     float64
	Neuroticism       float64
}

func NewDefaultPersonality() *Personality {
	p := &Personality{
		SelfConfidence:    1,
		Compassion:        0.5,
		PositiveTrustStep: 0.1,
		NegativeTrustStep: 0.1,
		Trustworthiness:   1,
		Openness:          0.5,
		Conscientiousness: 1, // Dependable by default
		Extraversion:      0.5,
		Agreeableness:     0.5,
		Neuroticism:       0.5,
	}
	randomizeTraits(p)
	return p
}

func randomizeTraits(p *Personality) {
	rand.Seed(time.Now().UnixNano())
	choice := rand.Intn(4)
	switch choice {
	case 0:
		p.Egalitarian = 1
	case 1:
		p.Selfish = 1
	case 2:
		p.Judgemental = 1
	case 3:
		p.Utilitarian = 1
	}
}
