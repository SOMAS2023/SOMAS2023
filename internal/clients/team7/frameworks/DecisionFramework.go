package frameworks

/*
	This is a generic interface that can be used for any decision framework.
	The inputs and outputs are generic types.
*/
type IDecisionFramework[I, O any] interface {
	GetDecision(inputs I) O
}
