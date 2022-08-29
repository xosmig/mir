package model

import (
	"fmt"
	"testing"
)

func TestTypeCast(t *testing.T) {
	root := &EventRootNode{}
	event := (Event)(root)

	_, ok := event.(*EventNode)
	fmt.Println(ok)

	_, ok = event.(*EventRootNode)
	fmt.Println(ok)
}
