/*
Package bpf wraps librte_bpf from DPDK project.

librte_bpf provides a framework to load and execute eBPF bytecode
inside user-space dpdk based applications. It supports basic set of
features from eBPF spec (https://www.kernel.org/doc/Documentation/networking/filter.txt).
*/
package bpf

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

/*
#include <rte_bpf.h>
*/
import "C"

// Arg implements rte_bpf_arg argument descriptor.
type Arg interface {
	transform() C.struct_rte_bpf_arg
}

// ArgRaw is the "raw" argument descriptor.
type ArgRaw struct {
	Size uintptr
}

func (a *ArgRaw) transform() (dst C.struct_rte_bpf_arg) {
	dst._type = C.RTE_BPF_ARG_RAW
	dst.size = C.size_t(a.Size)
	return
}

// ArgPtr is the pointer argument.
type ArgPtr struct {
	Size uintptr
}

func (a *ArgPtr) transform() (dst C.struct_rte_bpf_arg) {
	dst._type = C.RTE_BPF_ARG_PTR
	dst.size = C.size_t(a.Size)
	return
}

// ArgPtrMbuf is the rte_mbuf argument.
type ArgPtrMbuf struct {
	BufSize uintptr
}

func (a *ArgPtrMbuf) transform() (dst C.struct_rte_bpf_arg) {
	dst._type = C.RTE_BPF_ARG_PTR_MBUF
	dst.size = C.size_t(unsafe.Sizeof(C.struct_rte_mbuf{}))
	dst.buf_size = C.size_t(a.BufSize)
	return
}

var (
	// compile-time check
	_ = []Arg{
		&ArgRaw{},
		&ArgPtr{},
		&ArgPtrMbuf{},
	}
)

// BPF is the loaded BPF program.
type BPF C.struct_rte_bpf

// Destroy deallocates all memory used by this eBPF execution context.
func (b *BPF) Destroy() {
	C.rte_bpf_destroy((*C.struct_rte_bpf)(b))
}

// Prm specifies input parameters for loading eBPF code.
type Prm struct {
	Insns   *Insns
	XSyms   *XSyms
	ProgArg Arg
}

func (p *Prm) transform() C.struct_rte_bpf_prm {
	cp := C.struct_rte_bpf_prm{}

	if p.Insns != nil {
		cp.ins = p.Insns.ins
		cp.nb_ins = p.Insns.nbIns
	}

	if p.XSyms != nil {
		cp.xsym = p.XSyms.xsym
		cp.nb_xsym = p.XSyms.nbXSym
	}

	if p.ProgArg != nil {
		cp.prog_arg = p.ProgArg.transform()
	}
	return cp
}

// Load creates a new eBPF execution context and load given BPF code
// into it.
func Load(p *Prm) (*BPF, error) {
	cp := p.transform()
	if res := C.rte_bpf_load(&cp); res != nil {
		return (*BPF)(res), nil
	}
	return nil, common.RteErrno()
}

// ELFLoad creates a new eBPF execution context and load BPF code from given ELF
// file into it. Note that if the function will encounter EBPF_PSEUDO_CALL
// instruction that references external symbol, it will treat is as standard
// BPF_CALL to the external helper function.
func ELFLoad(p *Prm, fname, sname string) (*BPF, error) {
	cfname := C.CString(fname)
	defer C.free(unsafe.Pointer(cfname))
	csname := C.CString(sname)
	defer C.free(unsafe.Pointer(csname))

	cp := p.transform()
	if res := C.rte_bpf_elf_load(&cp, cfname, csname); res != nil {
		return (*BPF)(res), nil
	}
	return nil, common.RteErrno()
}

// Exec executes given BPF bytecode.
func (b *BPF) Exec(ctx unsafe.Pointer) uint64 {
	return uint64(C.rte_bpf_exec((*C.struct_rte_bpf)(b), ctx))
}
