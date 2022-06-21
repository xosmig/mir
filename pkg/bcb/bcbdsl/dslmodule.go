package bcbdsl

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

type cbModuleImpl struct {
	dsl.Module
	moduleId t.ModuleID
}

func NewModule(moduleId t.ModuleID) *cbModuleImpl {
	return &cbModuleImpl{
		Module:   dsl.NewModule(),
		moduleId: moduleId,
	}
}
