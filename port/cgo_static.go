// +build static

package port

/*
#cgo LDFLAGS: -Wl,--push-state,-Bstatic -lrte_port -lrte_ethdev -lrte_mempool -lrte_eal -lrte_kvargs -Wl,--pop-state
#cgo LDFLAGS: -Wl,--push-state,-Bdynamic -lpthread -lpcap -lnuma -ldl -Wl,--pop-state
*/
import "C"
