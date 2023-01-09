package table

import (
	"testing"

	"github.com/yerden/go-dpdk/common"
)

func TestACLParams(t *testing.T) {
	assert := common.Assert(t, true)

	p := &ACLParams{
		Name:  "hello",
		Rules: 16,
		FieldFormat: []ACLFieldDef{
			{Type: 1}, // some shit
		},
	}

	cptr, dtor := p.Transform(alloc)
	assert(cptr != nil)
	assert(dtor != nil)
	defer dtor(cptr)
}
