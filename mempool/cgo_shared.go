// +build shared

package mempool

/*
#cgo CFLAGS: -msse4.1 -mssse3
#cgo LDFLAGS: -lrte_mbuf -lrte_mempool -lrte_eal
*/
import "C"
