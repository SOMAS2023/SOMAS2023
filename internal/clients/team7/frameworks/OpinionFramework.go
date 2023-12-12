package frameworks

import (
	"math"

	"github.com/google/uuid"
)

type OpinionFrameworkInputs struct {
	AgentOpinion map[uuid.UUID]float64
	Mindset      float64
}

type OpinionFramework struct {
	Inputs *OpinionFrameworkInputs
}

func NewOpinionFramework(of OpinionFrameworkInputs) *OpinionFramework {
	return &OpinionFramework{Inputs: &of}
}

func (of *OpinionFramework) GetOpinion() float64 {

	i := len(of.Inputs.AgentOpinion)
	μ := of.Inputs.Mindset

	O := make([]float64, i)
	for idx := range O {
		O[idx] = of.Inputs.AgentOpinion[uuid.UUID]
	}

	W := make([]float64, i)
	for idx := range W {
		W[idx] = 1
	}

	A := make([]float64, i)
	for idx := range A {
		A[idx] = 1.0 - math.Abs(O[idx]-μ)/math.Max(μ, 1.0-μ)
	}

	for idx := range W {
		W[idx] = W[idx] + W[idx]*A[idx]
	}

	rowSum := 0.0
	for _, val := range W {
		rowSum += val
	}

	for idx := range W {
		W[idx] /= rowSum
	}

	o := 0.0
	for idx := range W {
		o += W[idx] * O[idx]
	}

	return o
}
