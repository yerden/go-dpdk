package pipeline

/*
#include "pipeline_loop.h"
*/
import "C"
import "unsafe"

// RunLoop runs pipeline continuously, flushing it every 'flush'
// iterations. RunLoop returns if value referenced by 'stop' is set to
// non-zero value.
//
// flush must be power of 2.
func (pl *Pipeline) RunLoop(flush uint32, stop *uint8) {
	params := &C.struct_lcore_arg{}
	params.p = (*C.struct_rte_pipeline)(pl)
	params.flush = C.uint32_t(flush)
	params.stop = (*C.uint8_t)(stop)

	C.run_pipeline_loop(unsafe.Pointer(params))
}
