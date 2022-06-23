/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Package modules provides interfaces of modules that serve as building blocks of a Node.
// Implementations of those interfaces are not contained by this package
// and are expected to be provided by other packages.
package modules

import (
	t "github.com/filecoin-project/mir/pkg/types"
)

// Module generalizes the ActiveModule and PassiveModule types.
type Module interface {

	// ImplementsModule only serves the purpose of indicating that this is a Module and must not be called.
	ImplementsModule()
}

// The Modules structs groups the modules a Node consists of.
type Modules map[t.ModuleID]Module

// Forkable modules can create copies of themselves.
type Forkable interface {
	Module
	// Fork creates a new module with the given id. This module may or may not have the same state as the parent module,
	// depending on the implementation.
	Fork(newModuleID string) (newModule Module, err error)
}
