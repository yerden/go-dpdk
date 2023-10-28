package bpf

/*
#include <rte_bpf.h>

typedef uint64_t (bpf_func)(uint64_t, uint64_t, uint64_t, uint64_t, uint64_t);

struct xsym_func {
	bpf_func *val;
	uint32_t nb_args;
	struct rte_bpf_arg args[EBPF_FUNC_MAX_ARGS];
	struct rte_bpf_arg ret;
};

struct xsym_var {
	void *val;
	struct rte_bpf_arg desc;
};

void
xsym_set_func(struct rte_bpf_xsym *sym, struct xsym_func *v)
{
	sym->type = RTE_BPF_XTYPE_FUNC;
	sym->func.val = v->val;
	sym->func.nb_args = v->nb_args;
	int i;
	for (i = 0; i < RTE_DIM(sym->func.args); i++)
		sym->func.args[i] = v->args[i];

	sym->func.ret = v->ret;
}

void
xsym_set_var(struct rte_bpf_xsym *sym, struct xsym_var *v)
{
	sym->type = RTE_BPF_XTYPE_VAR;
	sym->var.val = v->val;
	sym->var.desc = v->desc;
}
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// Func represent a C-function to use as a symbol in XSymFunc.
type Func C.bpf_func

// XSym is the definition for external symbols available in the BPF
// program.
type XSym interface {
	XSymName() string

	// transform itself into C representation, name is specified
	transform(*C.char) (instance C.struct_rte_bpf_xsym)
}

var (
	_ XSym = (*XSymFunc)(nil)
	_ XSym = (*XSymVar)(nil)
)

// XSymFunc is the function-type of a symbol.
type XSymFunc struct {
	Name string
	Val  *Func
	// only up to EBPF_FUNC_MAX_ARGS will be taken
	Args []Arg
	Ret  Arg
}

// XSymVar is the var-type of a symbol.
type XSymVar struct {
	Name string
	Val  unsafe.Pointer
	Desc Arg
}

// XSymName implements XSym interface.
func (sym *XSymFunc) XSymName() string {
	return sym.Name
}

func (sym *XSymFunc) transform(name *C.char) (instance C.struct_rte_bpf_xsym) {
	instance.name = name
	choice := C.struct_xsym_func{
		val:     (*C.bpf_func)(sym.Val),
		nb_args: C.uint32_t(len(sym.Args)),
	}

	if sym.Ret != nil {
		choice.ret = sym.Ret.transform()
	}

	for i := range choice.args {
		if i == len(sym.Args) {
			break
		}
		choice.args[i] = sym.Args[i].transform()
	}

	C.xsym_set_func(&instance, &choice)
	return
}

// XSymName implements XSym interface.
func (sym *XSymVar) XSymName() string {
	return sym.Name
}

func (sym *XSymVar) transform(name *C.char) (instance C.struct_rte_bpf_xsym) {
	instance.name = name
	choice := C.struct_xsym_var{
		val:  sym.Val,
		desc: sym.Desc.transform(),
	}
	C.xsym_set_var(&instance, &choice)
	return
}

// XSyms contains C-allocated list of XSym descriptors.
type XSyms struct {
	xsym   *C.struct_rte_bpf_xsym
	nbXSym C.uint32_t
}

// NewXSyms allocates list of XSym descriptors to use in Prm.
func NewXSyms(data []XSym) *XSyms {
	xsyms := &XSyms{}

	a := common.NewAllocatorSession(&common.StdAlloc{})
	d := C.uint32_t(len(data))
	common.CallocT(a, &xsyms.xsym, int(d))
	xsyms.nbXSym = d
	ins := unsafe.Slice(xsyms.xsym, d)
	for i, is := range data {
		s := (*C.char)(common.CString(a, is.XSymName()))
		ins[i] = is.transform(s)
	}

	runtime.SetFinalizer(xsyms, func(xsyms *XSyms) {
		a.Flush()
	})

	return xsyms
}
