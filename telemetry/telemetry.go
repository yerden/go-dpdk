package telemetry

/*
#include <rte_common.h>
#include <rte_telemetry.h>

extern int telemetryHandler(char *name, char *params, struct rte_tel_data *d);
telemetry_cb telemetry_handler_cb = (telemetry_cb)telemetryHandler;

*/
import "C"

import (
	"sync"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	_ "github.com/yerden/go-dpdk/eal"
)

type ValueType uint32

const (
	StringVal ValueType = C.RTE_TEL_STRING_VAL
	IntVal    ValueType = C.RTE_TEL_INT_VAL
	U64Val    ValueType = C.RTE_TEL_U64_VAL
	Container ValueType = C.RTE_TEL_CONTAINER
)

type Data C.struct_rte_tel_data

func cData(d *Data) *C.struct_rte_tel_data {
	return (*C.struct_rte_tel_data)(d)
}

func (d *Data) StartArray(t ValueType) error {
	return common.IntToErr(int64(C.rte_tel_data_start_array(cData(d), uint32(t))))
}

func (d *Data) StartDict() error {
	return common.IntToErr(int64(C.rte_tel_data_start_dict(cData(d))))
}

func (d *Data) SetString(s string) error {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return common.IntToErr(int64(C.rte_tel_data_string(cData(d), cs)))
}

func (d *Data) AddArrayString(s string) error {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return common.IntToErr(int64(C.rte_tel_data_add_array_string(cData(d), cs)))
}

func (d *Data) AddArrayInt(x int) error {
	return common.IntToErr(int64(C.rte_tel_data_add_array_int(cData(d), C.int(x))))
}

func (d *Data) AddArrayU64(x uint64) error {
	return common.IntToErr(int64(C.rte_tel_data_add_array_u64(cData(d), C.uint64_t(x))))
}

func cFlag(f bool) C.int {
	if f {
		return 1
	}
	return 0
}

func (d *Data) AddArrayContainer(x *Data, keep bool) error {
	return common.IntToErr(int64(C.rte_tel_data_add_array_container(
		cData(d), cData(x), cFlag(keep))))
}

func (d *Data) AddDictString(name, val string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cVal))
	return common.IntToErr(int64(C.rte_tel_data_add_dict_string(cData(d), cName, cVal)))
}

func (d *Data) AddDictInt(name string, x int) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return common.IntToErr(int64(C.rte_tel_data_add_dict_int(cData(d), cName, C.int(x))))
}

func (d *Data) AddDictU64(name string, x uint64) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return common.IntToErr(int64(C.rte_tel_data_add_dict_u64(cData(d), cName, C.uint64_t(x))))
}

func (d *Data) AddDictContainer(name string, x *Data, keep bool) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return common.IntToErr(int64(C.rte_tel_data_add_dict_container(
		cData(d), cName, cData(x), cFlag(keep))))
}

type Handler func(string, string, *Data) int

var callbacks struct {
	mtx      sync.Mutex
	handlers map[string]Handler
}

func init() {
	callbacks.handlers = map[string]Handler{}
}

//export telemetryHandler
func telemetryHandler(cCmd *C.char, cParams *C.char, data *C.struct_rte_tel_data) C.int {
	cmd := C.GoString(cCmd)
	params := C.GoString(cParams)
	callbacks.mtx.Lock()
	handler := callbacks.handlers[cmd]
	callbacks.mtx.Unlock()
	return C.int(handler(cmd, params, (*Data)(data)))
}

func RegisterCmd(cmd, help string, h Handler) int {
	callbacks.mtx.Lock()
	callbacks.handlers[cmd] = h
	callbacks.mtx.Unlock()

	cCmd := C.CString(cmd)
	defer C.free(unsafe.Pointer(cCmd))

	cHelp := C.CString(help)
	defer C.free(unsafe.Pointer(cHelp))

	return int(C.rte_telemetry_register_cmd(cCmd, C.telemetry_handler_cb, cHelp))
}
