package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>
*/
import "C"
import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// ItemType represents rte_flow_item type.
type ItemType uint32

// ItemValue should be implemented to specify in Item.
type ItemValue interface {
	common.Transformer

	// ItemType returns implemented rte_flow_item_* struct.
	ItemType() ItemType

	// Mask returns pointer to rte_flow_item_*_mask variables. They
	// should not be changed by user so Mask returns pointer to C
	// struct.
	DefaultMask() unsafe.Pointer
}

// Item is the matching pattern item definition.
//
// A pattern is formed by stacking items starting from the lowest
// protocol layer to match. This stacking restriction does not apply
// to meta items which can be placed anywhere in the stack without
// affecting the meaning of the resulting pattern.
//
// Patterns are terminated by END items.
//
// The spec field should be a valid pointer to a structure of the
// related item type. It may remain unspecified (NULL) in many cases
// to request broad (nonspecific) matching. In such cases, last and
// mask must also be set to NULL.
//
// Optionally, last can point to a structure of the same type to
// define an inclusive range. This is mostly supported by integer and
// address fields, may cause errors otherwise. Fields that do not
// support ranges must be set to 0 or to the same value as the
// corresponding fields in spec.
//
// Only the fields defined to nonzero values in the default masks (see
// rte_flow_item_{name}_mask constants) are considered relevant by
// default.  This can be overridden by providing a mask structure of
// the same type with applicable bits set to one. It can also be used
// to partially filter out specific fields (e.g. as an alternate mean
// to match ranges of IP addresses).
//
// Mask is a simple bit-mask applied before interpreting the contents
// of spec and last, which may yield unexpected results if not used
// carefully. For example, if for an IPv4 address field, spec provides
// 10.1.2.3, last provides 10.3.4.5 and mask provides 255.255.0.0, the
// effective range becomes 10.1.0.0 to 10.3.255.255.
//
// Go: you may also specify ItemType as a Spec field if you don't want
// to specify any pattern for the item.
type Item struct {
	Spec, Last, Mask ItemValue
}
