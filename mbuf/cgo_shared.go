// +build shared

package mbuf

/*
#cgo CFLAGS: -msse4.1 -mssse3
#cgo LDFLAGS: -lrte_mbuf -lrte_mempool -lrte_ring -lrte_eal -lrte_kvargs
*/
import "C"
