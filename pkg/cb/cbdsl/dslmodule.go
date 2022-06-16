package cbdsl

import (
	"github.com/filecoin-project/mir/pkg/modules/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

type CBModule struct {
	dsl.DslModule
	moduleId t.ModuleID
}

func NewModule(moduleId t.ModuleID) *CBModule {
	return &CBModule{
		DslModule: dsl.NewModule(),
		moduleId:  moduleId,
	}
}
