package cbdsl

import (
	"github.com/filecoin-project/mir/pkg/modules/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

type cbModuleImpl struct {
	dsl.DslModule
	moduleId t.ModuleID
}

func NewModule(moduleId t.ModuleID) *cbModuleImpl {
	return &cbModuleImpl{
		DslModule: dsl.NewModule(),
		moduleId:  moduleId,
	}
}
