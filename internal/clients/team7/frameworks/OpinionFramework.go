package frameworks

import (
	"fmt"
	"math"
	"math/rand"
)

// Define the necessary types and structures here

type OpinionFrameworkInputs struct {
}

type OpinionFrameworkOutputs struct {
}

type OpinionFramework struct {
	IOpinionFramework[OpinionFrameworkInputs, SocialNetwork, OpinionFrameworkOutputs]
	inputs *OpinionFrameworkInputs
}

// Constructor for OpinionFramework
func NewOpinionFramework(of OpinionFrameworkInputs) *OpinionFramework {
	return &OpinionFramework{inputs: &of}
}

func (of *OpinionFramework) GetOpinion(inputs OpinionFrameworkInputs, socialNetwork SocialNetwork) OpinionFrameworkOutputs {
	i := 5
	j := 6

	μ := make([]float64, i)
	for idx := range μ {
		μ[idx] = rand.Float64()
	}

	O := make([][]float64, i)
	for idx := range O {
		O[idx] = make([]float64, j)
		for jdx := range O[idx] {
			O[idx][jdx] = rand.Float64()
		}
	}

	W := make([][]float64, i)
	for idx := range W {
		W[idx] = make([]float64, j)
		for jdx := range W[idx] {
			W[idx][jdx] = rand.Float64()
		}
	}

	A := make([][]float64, i)
	for idx := range A {
		A[idx] = make([]float64, j)
		for jdx := range A[idx] {
			A[idx][jdx] = 1.0 - math.Abs(O[idx][jdx]-μ[idx])/math.Max(μ[idx], 1.0-μ[idx])
		}
	}

	for idx := range W {
		for jdx := range W[idx] {
			W[idx][jdx] = W[idx][jdx] + W[idx][jdx]*A[idx][jdx]
		}
	}

	for idx := range W {
		rowSum := 0.0
		for _, val := range W[idx] {
			rowSum += val
		}
		for jdx := range W[idx] {
			W[idx][jdx] /= rowSum
		}
	}

	o := make([][]float64, i)
	for idx := range o {
		o[idx] = make([]float64, 1)
		for jdx := range o[idx] {
			sum := 0.0
			for kdx := range W[idx] {
				sum += W[idx][kdx] * O[idx][kdx]
			}
			o[idx][jdx] = sum
		}
	}

	return OpinionFrameworkOutputs{o}
}
