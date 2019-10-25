package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>
*/
import "C"
import "unsafe"

// ActionType is the rte_flow_action type.
type ActionType uint32

// Reload implements Action interface.
func (t ActionType) Reload() {}

// Pointer implements Action interface.
func (t ActionType) Pointer() unsafe.Pointer { return nil }

// Type implements Action interface.
func (t ActionType) Type() ActionType { return t }

// Action is the definition of a single action.
//
// A list of actions is terminated by a END action.
//
// For simple actions without a configuration object, conf remains
// NULL.
type Action interface {
	// Pointer returns a valid C pointer to underlying struct.
	Pointer() unsafe.Pointer

	// Reload is used to apply changes so that the underlying struct
	// reflects the up-to-date configuration.
	Reload()

	// Type returns implemented rte_flow_action_* struct.
	Type() ActionType
}
