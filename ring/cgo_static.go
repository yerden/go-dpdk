// +build static

package ring

/*
#cgo LDFLAGS: -Wl,--push-state,-Bstatic -lrte_ring -lrte_eal -lrte_kvargs -Wl,--pop-state
#cgo LDFLAGS: -Wl,--push-state,-Bdynamic -lpthread -lnuma -ldl -Wl,--pop-state
*/
import "C"
