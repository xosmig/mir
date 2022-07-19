package deploytest

import (
	"fmt"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

// NewNodeIDs returns a slice of node ids of the given size suitable for testing.
func NewNodeIDs(nNodes int) []t.NodeID {
	return sliceutil.Generate(nNodes, func(i int) t.NodeID {
		return t.NodeID(fmt.Sprintf("Node%v", i))
	})
}
