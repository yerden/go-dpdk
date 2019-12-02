// +build static

package eal

/*
#cgo LDFLAGS: -Wl,--push-state,-Bstatic -lrte_eal -lrte_kvargs -Wl,--pop-state
#cgo LDFLAGS: -Wl,--push-state,-Bdynamic -lpthread -lnuma -ldl -Wl,--pop-state
*/
import "C"
