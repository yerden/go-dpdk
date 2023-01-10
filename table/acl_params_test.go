package table

import (
	"testing"

	"github.com/yerden/go-dpdk/common"
)

func TestACLField(t *testing.T) {
	assert := common.Assert(t, true)

	{
		f1 := NewACLField8(0xAA, 0xFF)
		f2 := NewACLField8(0xAA, 0xFF)
		assert(f1 == f2)
	}

	{
		f1 := NewACLField16(0xAABB, 0xFFF0)
		f2 := NewACLField16(0xAABB, 0xFFF0)
		assert(f1 == f2)
	}

	{
		f1 := NewACLField32(0xAABBFFF0, 0xFFFFFF00)
		f2 := NewACLField32(0xAABBFFF0, 0xFFFFFF00)
		assert(f1 == f2)
	}

	{
		f1 := NewACLField64(0xAABBFFF0FFFFFF00, 0xFFFFFF0FFFFFF00)
		f2 := NewACLField64(0xAABBFFF0FFFFFF00, 0xFFFFFF0FFFFFF00)
		assert(f1 == f2)
	}

	newRule := ACLRuleAdd{}
	newRule.Priority = 2
	newRule.FieldValue[0] = NewACLField8(0xAA, 0xFF)
}
