package flow

/*
#include <rte_config.h>
#include <rte_flow.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// Error is a verbose error structure definition.
//
// This object is normally allocated by applications and set by PMDs,
// the message points to a constant string which does not need to be
// freed by the application, however its pointer can be considered
// valid only as long as its associated DPDK port remains configured.
// Closing the underlying device or unloading the PMD invalidates it.
//
// Both cause and message may be NULL regardless of the error type.
type Error C.struct_rte_flow_error

func (e *Error) Error() string {
	return fmt.Sprintf("%v: %s", e.Unwrap(), C.GoString(e.message))
}

func (e *Error) Unwrap() error {
	return ErrorType(e._type)
}

// Cause returns object responsible for error.
func (e *Error) Cause() unsafe.Pointer {
	return e.cause
}

// ErrorType is a type of an error.
type ErrorType uint

func (e ErrorType) Error() string {
	if s, ok := errStr[e]; ok {
		return s
	}
	return ""
}

var (
	errStr = make(map[ErrorType]string)
)

func registerErr(c uint, str string) ErrorType {
	et := ErrorType(c)
	errStr[et] = str
	return et
}

// Error types.
var (
	ErrTypeNone         = registerErr(C.RTE_FLOW_ERROR_TYPE_NONE, "No error")
	ErrTypeUnspecified  = registerErr(C.RTE_FLOW_ERROR_TYPE_UNSPECIFIED, "Cause unspecified")
	ErrTypeHandle       = registerErr(C.RTE_FLOW_ERROR_TYPE_HANDLE, "Flow rule (handle)")
	ErrTypeAttrGroup    = registerErr(C.RTE_FLOW_ERROR_TYPE_ATTR_GROUP, "Group field")
	ErrTypeAttrPriority = registerErr(C.RTE_FLOW_ERROR_TYPE_ATTR_PRIORITY, "Priority field")
	ErrTypeAttrIngress  = registerErr(C.RTE_FLOW_ERROR_TYPE_ATTR_INGRESS, "Ingress field")
	ErrTypeAttrEgress   = registerErr(C.RTE_FLOW_ERROR_TYPE_ATTR_EGRESS, "Egress field")
	ErrTypeAttrTransfer = registerErr(C.RTE_FLOW_ERROR_TYPE_ATTR_TRANSFER, "Transfer field")
	ErrTypeAttr         = registerErr(C.RTE_FLOW_ERROR_TYPE_ATTR, "Attributes structure")
	ErrTypeItemNum      = registerErr(C.RTE_FLOW_ERROR_TYPE_ITEM_NUM, "Pattern length")
	ErrTypeItemSpec     = registerErr(C.RTE_FLOW_ERROR_TYPE_ITEM_SPEC, "Item specification spec")
	ErrTypeItemLast     = registerErr(C.RTE_FLOW_ERROR_TYPE_ITEM_LAST, "Item specification range")
	ErrTypeItemMask     = registerErr(C.RTE_FLOW_ERROR_TYPE_ITEM_MASK, "Item specification mask")
	ErrTypeItem         = registerErr(C.RTE_FLOW_ERROR_TYPE_ITEM, "Specific pattern item")
	ErrTypeActionNum    = registerErr(C.RTE_FLOW_ERROR_TYPE_ACTION_NUM, "Number of actions")
	ErrTypeActionConf   = registerErr(C.RTE_FLOW_ERROR_TYPE_ACTION_CONF, "Action configuration")
	ErrTypeAction       = registerErr(C.RTE_FLOW_ERROR_TYPE_ACTION, "Specific action")
)
