package frameworks

/*
	This is a generic interface that can be used for any decision framework.
	The inputs and outputs are generic types which can be defined per decision framework type.
	We can add more functions here when needed.
*/
type IDecisionFramework[I, O any] interface {
	GetDecision(inputs I) O
}
