// +build shared

package memzone

/*
#cgo LDFLAGS: -lrte_eal -lrte_kvargs
#cgo LDFLAGS: -Wl,--push-state,-Bdynamic -lpthread -lnuma -ldl -Wl,--pop-state
*/
import "C"
