package ethdev

/*
#include "lsc_telemetry.h"
*/
import "C"
import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// RegisterCallbackLSC installs a callback for RTE_ETH_EVENT_INTR_LSC event
// which counts all occurred events.
func (p Port) RegisterCallbackLSC() error {
	return common.IntErr(int64(C.lsc_counters_callback_register(C.ushort(p))))
}

// UnregisterCallbackLSC removes callback for RTE_ETH_EVENT_INTR_LSC event
// which counts all occurred events.
func (p Port) UnregisterCallbackLSC() error {
	return common.IntErr(int64(C.lsc_counters_callback_unregister(C.ushort(p))))
}

// RegisterTelemetryLSC registers telemetry handlers for accessing LSC counters.
func RegisterTelemetryLSC(name string) {
	cname := C.CString(name)
	C.lsc_register_telemetry_cmd(cname)
	C.free(unsafe.Pointer(cname))
}
