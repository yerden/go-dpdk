package bpf

/*
#include <rte_bpf.h>

void
insn_set_dst_reg(struct ebpf_insn *d, uint8_t reg)
{
	d->dst_reg = reg;
}

void
insn_set_src_reg(struct ebpf_insn *d, uint8_t reg)
{
	d->src_reg = reg;
}
*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// Insn is an eBPF instruction descriptor.
type Insn struct {
	Code   uint8
	DstReg uint8 // 4 bit
	SrcReg uint8 // 4 bit
	Off    int16
	Imm    int32
}

func (d *Insn) transform() (insn C.struct_ebpf_insn) {
	insn.code = C.uchar(d.Code)
	C.insn_set_dst_reg(&insn, C.uchar(d.DstReg))
	C.insn_set_src_reg(&insn, C.uchar(d.SrcReg))
	insn.off = C.int16_t(d.Off)
	insn.imm = C.int32_t(d.Imm)
	return
}

// Insns contains C-allocated list of eBPF instructions.
type Insns struct {
	ins   *C.struct_ebpf_insn
	nbIns C.uint32_t
}

// NewInsns allocates list of instructions to use in Prm.
func NewInsns(data []Insn) *Insns {
	insns := &Insns{}

	a := common.NewAllocatorSession(&common.StdAlloc{})
	d := C.uint32_t(len(data))
	common.CallocT(a, &insns.ins, int(d))
	insns.nbIns = d
	ins := unsafe.Slice(insns.ins, d)
	for i, is := range data {
		ins[i] = is.transform()
	}

	runtime.SetFinalizer(insns, func(insns *Insns) {
		a.Flush()
	})

	return insns
}
