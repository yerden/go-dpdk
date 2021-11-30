package lpm

/*
#include <stdlib.h>
#include <rte_config.h>
#include <rte_lpm6.h>
*/
import "C"

import (
	"net"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// LPM6 is an RTE Longest Prefix Match lookup object.
type LPM6 C.struct_rte_lpm6

// Config6 is used to configure LPM6 object while creation.
type Config6 struct {
	MaxRules    uint32
	NumberTbl8s uint32
	Flags       int
}

// Create6 an LPM6 object.
//
// Specify name as an LPM6 object name, socket_id as a NUMA socket ID
// for LPM6 table memory allocation and config as a structure
// containing the configuration.
//
// Returns handle to LPM6 object on success, and errno value:
//   E_RTE_NO_CONFIG - function could not get pointer to rte_config structure
//   E_RTE_SECONDARY - function was called from a secondary process instance
//   EINVAL - invalid parameter passed to function
//   ENOSPC - the maximum number of memzones has already been allocated
//   EEXIST - a memzone with the same name already exists
//   ENOMEM - no appropriate memory area found in which to create memzone
func Create6(name string, socket int, cfg *Config6) (*LPM6, error) {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	config := &C.struct_rte_lpm6_config{
		max_rules:    C.uint32_t(cfg.MaxRules),
		number_tbl8s: C.uint32_t(cfg.NumberTbl8s),
		flags:        C.int(cfg.Flags),
	}

	if r := (*LPM6)(C.rte_lpm6_create(s, C.int(socket), config)); r != nil {
		return r, nil
	}

	return nil, common.RteErrno()
}

// Free an LPM6 object.
func (r *LPM6) Free() {
	C.rte_lpm6_free((*C.struct_rte_lpm6)(r))
}

// Add a rule to LPM6 object.
//
// ip/prefix is an IP address/subnet to add, nextHop is a value
// associated with added IP subnet. Panics if ip is not IPv6.
func (r *LPM6) Add(ipnet net.IPNet, nextHop uint32) error {
	b, prefix := cvtIPv6Net(ipnet)
	rc := C.rte_lpm6_add((*C.struct_rte_lpm6)(r), (*C.uint8_t)(&b[0]), C.uint8_t(prefix), C.uint32_t(nextHop))
	return common.IntToErr(rc)
}

// Delete a rule from LPM6 object. Panics if ip is not IPv6.
func (r *LPM6) Delete(ipnet net.IPNet) error {
	b, prefix := cvtIPv6Net(ipnet)
	rc := C.rte_lpm6_delete((*C.struct_rte_lpm6)(r), (*C.uint8_t)(&b[0]), C.uint8_t(prefix))
	return common.IntToErr(rc)
}

// Lookup an IP in LPM6 object. Panics if ip is not IPv6.
func (r *LPM6) Lookup(ip net.IP) (uint32, error) {
	var res uint32
	rc := C.rte_lpm6_lookup((*C.struct_rte_lpm6)(r), (*C.uint8_t)(&ip[0]), (*C.uint32_t)(&res))
	return res, common.IntToErr(rc)
}

// DeleteAll removes all rules from LPM6 object.
func (r *LPM6) DeleteAll() {
	C.rte_lpm6_delete_all((*C.struct_rte_lpm6)(r))
}

// IsRulePresent checks if a rule present in the LPM6 and returns
// nextHop if it is. Panics if ip is not IPv6.
func (r *LPM6) IsRulePresent(ipnet net.IPNet, nextHop *uint32) (bool, error) {
	b, prefix := cvtIPv6Net(ipnet)
	rc := C.rte_lpm6_is_rule_present((*C.struct_rte_lpm6)(r), (*C.uint8_t)(&b[0]), C.uint8_t(prefix), (*C.uint32_t)(nextHop))
	n, err := common.IntOrErr(rc)
	return n != 0, err
}
