package ring

/*
#include "ring_telemetry.h"
*/
import "C"
import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

func telemetryRegisterCmd(path string, handler C.telemetry_cb, help string) C.int {
	sPath := C.CString(path)
	defer C.free(unsafe.Pointer(sPath))
	sHelp := C.CString(help)
	defer C.free(unsafe.Pointer(sHelp))
	return C.rte_telemetry_register_cmd(sPath, (C.telemetry_cb)(handler), sHelp)
}

type cmdDesc struct {
	cmd  string
	help string
	cb   C.telemetry_cb
}

// TelemetryInit initializes telemetry callbacks for rings.
// Specify prefix for cmd path to avoid conflicts in the future.
// "/ring" is the good candidate for a prefix.
func TelemetryInit(prefix string) {
	desc := []cmdDesc{
		{
			cmd:  prefix + "/list",
			cb:   C.ring_list_cb,
			help: "Show list of rings. Takes no parameters.",
		}, {
			cmd:  prefix + "/info",
			cb:   C.ring_info_cb,
			help: "Show info on the ring. Param: ring name.",
		},
	}
	for _, d := range desc {
		if rc := telemetryRegisterCmd(d.cmd, d.cb, d.help); rc != 0 {
			panic(common.IntErr(int64(rc)))
		}
	}
}
