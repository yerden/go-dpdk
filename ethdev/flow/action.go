package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>
*/
import "C"
import (
	"github.com/yerden/go-dpdk/common"
)

// ActionType is the rte_flow_action type.
type ActionType uint32

// Action is the definition of a single action.
//
// A list of actions is terminated by a END action.
//
// For simple actions without a configuration object, conf remains
// NULL.
type Action interface {
	common.Transformer

	// ActionType returns implemented rte_flow_action_* struct.
	ActionType() ActionType
}
