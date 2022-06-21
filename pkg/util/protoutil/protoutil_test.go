package protoutil

import (
	"github.com/filecoin-project/mir/pkg/util/protoutil/testpb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyOneofWrapper(t *testing.T) {
	assert.Nil(t, VerifyOneofWrapper[testpb.Foo_BarOption, *testpb.Bar]())
	assert.Nil(t, VerifyOneofWrapper[testpb.Foo_BazOption, uint64]())
}
