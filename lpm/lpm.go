/*
Package lpm wraps RTE LPM library.

Please refer to DPDK Programmer's Guide for reference and caveats.
*/
package lpm

/*
#include <rte_config.h>
#include <rte_lpm.h>
#include <rte_lpm6.h>
*/
import "C"

import (
	"net/netip"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// LPM is an RTE Longest Prefix Match lookup object.
type LPM C.struct_rte_lpm

// Config is used to configure LPM object while creation.
type Config struct {
	MaxRules    uint32
	NumberTbl8s uint32
	Flags       int
}

// FindExisting finds an existing LPM object and return a pointer to
// it.
//
// Specify pointer to *LPM or *LPM6 in ptr to find respective object
// by its memzone name.
func FindExisting(name string, ptr interface{}) error {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	var found bool

	switch r := ptr.(type) {
	case **LPM:
		*r = (*LPM)(C.rte_lpm_find_existing(s))
		found = *r != nil
	case **LPM6:
		*r = (*LPM6)(C.rte_lpm6_find_existing(s))
		found = *r != nil
	default:
		panic("incompatible value")
	}

	if !found {
		return common.RteErrno()
	}

	return nil
}

// Create an LPM object.
//
// Specify name as an LPM object name, socket_id as a NUMA socket ID
// for LPM table memory allocation and config as a structure
// containing the configuration.
//
// Returns handle to LPM object on success, and errno value:
//
//	E_RTE_NO_CONFIG - function could not get pointer to rte_config structure
//	E_RTE_SECONDARY - function was called from a secondary process instance
//	EINVAL - invalid parameter passed to function
//	ENOSPC - the maximum number of memzones has already been allocated
//	EEXIST - a memzone with the same name already exists
//	ENOMEM - no appropriate memory area found in which to create memzone
func Create(name string, socket int, cfg *Config) (*LPM, error) {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	config := &C.struct_rte_lpm_config{
		max_rules:    C.uint32_t(cfg.MaxRules),
		number_tbl8s: C.uint32_t(cfg.NumberTbl8s),
		flags:        C.int(cfg.Flags),
	}

	if r := (*LPM)(C.rte_lpm_create(s, C.int(socket), config)); r != nil {
		return r, nil
	}

	return nil, common.RteErrno()
}

// Free an LPM object.
func (r *LPM) Free() {
	C.rte_lpm_free((*C.struct_rte_lpm)(r))
}

// Add a rule to LPM object. Panics if ipnet is not IPv4 subnet.
//
// ip/prefix is an IP address/subnet to add, nextHop is a value
// associated with added IP subnet.
func (r *LPM) Add(ipnet netip.Prefix, nextHop uint32) error {
	ip, prefix := cvtIPv4Net(ipnet)
	rc := C.rte_lpm_add((*C.struct_rte_lpm)(r), C.uint32_t(ip), C.uint8_t(prefix), C.uint32_t(nextHop))
	return common.IntToErr(rc)
}

// Delete a rule from LPM object. Panics if ipnet is not IPv4 subnet.
func (r *LPM) Delete(ipnet netip.Prefix) error {
	ip, prefix := cvtIPv4Net(ipnet)
	rc := C.rte_lpm_delete((*C.struct_rte_lpm)(r), C.uint32_t(ip), C.uint8_t(prefix))
	return common.IntToErr(rc)
}

// Lookup an IP in LPM object. Panics if ip is not IPv4 subnet.
func (r *LPM) Lookup(ip netip.Addr) (uint32, error) {
	b := cvtIPv4(ip)
	var res uint32
	rc := C.rte_lpm_lookup((*C.struct_rte_lpm)(r), C.uint32_t(b), (*C.uint32_t)(&res))
	return res, common.IntToErr(rc)
}

// DeleteAll removes all rules from LPM object.
func (r *LPM) DeleteAll() {
	C.rte_lpm_delete_all((*C.struct_rte_lpm)(r))
}

// IsRulePresent checks if a rule present in the LPM and returns
// nextHop if it is. Panics if ipnet is not IPv4 subnet.
func (r *LPM) IsRulePresent(ipnet netip.Prefix, nextHop *uint32) (bool, error) {
	ip, prefix := cvtIPv4Net(ipnet)
	rc := C.rte_lpm_is_rule_present((*C.struct_rte_lpm)(r), C.uint32_t(ip), C.uint8_t(prefix), (*C.uint32_t)(nextHop))
	n, err := common.IntOrErr(rc)
	return n != 0, err
}
