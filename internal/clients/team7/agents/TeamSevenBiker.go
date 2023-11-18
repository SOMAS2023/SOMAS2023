package agents

import (
	objects "SOMAS2023/internal/common/objects"
)

type ITeamSevenBiker interface {
	objects.IBaseBiker
}

type BaseTeamSevenBiker struct {
	*objects.BaseBiker
}
