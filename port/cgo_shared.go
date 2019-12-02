// +build shared

package port

/*
#cgo LDFLAGS: -lrte_port -lrte_ethdev -lrte_mempool -lrte_eal -lrte_kvargs
#cgo LDFLAGS: -Wl,--push-state,-Bdynamic -lpthread -lpcap -lnuma -ldl -Wl,--pop-state
*/
import "C"
