// +build static

package mbuf

/*
#cgo CFLAGS: -msse4.1 -mssse3
#cgo LDFLAGS: -Wl,--push-state,-Bstatic -lrte_mbuf -lrte_mempool -lrte_ring -lrte_eal -lrte_kvargs -Wl,--pop-state
#cgo LDFLAGS: -Wl,--push-state,-Bdynamic -lpthread -lnuma -ldl -Wl,--pop-state
*/
import "C"
