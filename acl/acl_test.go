package acl

import (
	"fmt"
	"syscall"
	"testing"

	"github.com/yerden/go-dpdk/eal"
)

func assert(t testing.TB, expected bool, args ...interface{}) {
	if !expected {
		t.Helper()
		t.Fatal(args...)
	}
}

func TestContext(t *testing.T) {
	eal.InitOnceSafe("test", 4)

	cfg := &Config{
		Categories: 2,
		MaxSize:    0x800000,
		Defs: []FieldDef{
			{
				Type:       FieldTypeBitmask,
				Size:       1,
				Offset:     9, // ipv4 header -> ipproto
				FieldIndex: 0,
				InputIndex: 0,
			}, {
				Type:       FieldTypeMask,
				Size:       4,
				Offset:     12, // ipv4 header -> srcaddr
				FieldIndex: 1,
				InputIndex: 1,
			}, {
				Type:       FieldTypeMask,
				Size:       4,
				Offset:     16, // ipv4 header -> srcaddr
				FieldIndex: 2,
				InputIndex: 2,
			},
		},
	}

	p := &Param{
		Name:       "hello",
		RuleSize:   RuleSize(len(cfg.Defs)),
		MaxRuleNum: 64,
		SocketID:   -1,
	}

	fmt.Println("rule size=", p.RuleSize)
	ctx, err := Create(p)
	assert(t, err == nil, err)
	assert(t, ctx != nil)

	ctx.Dump()

	ctx.Reset()

	otherCtx, err := FindExisting(p.Name)
	assert(t, err == nil)
	assert(t, ctx == otherCtx)

	_, err = FindExisting("some_shit")
	assert(t, err == syscall.ENOENT)

	err = ctx.AddRules([]Rule{
		{
			Data: RuleData{CategoryMask: 3, Priority: 1, Userdata: 1},
			Fields: []Field{
				{uint8(6), uint8(0xff)},
				{uint32(0), uint8(0)},
				{uint32(0), uint8(0)},
			},
		},
	})
	assert(t, err == nil, err)

	err = ctx.AddRules([]Rule{
		{
			Data: RuleData{CategoryMask: 3, Priority: 1, Userdata: 1},
			Fields: []Field{
				{uint8(17), uint8(0xff)},
				{uint32(0x01020304), uint8(24)},
				{uint32(0), uint8(0)},
			},
		},
	})
	assert(t, err == nil, err)

	err = ctx.Build(cfg)
	assert(t, err == nil, err)

	ListDump()

	ctx.Free()
}
