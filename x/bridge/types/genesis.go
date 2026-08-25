package types

import (
	"errors"
	"math/big"
)

// GenesisState deliberately contains no enabled default route. A route may be
// present only when it passes the same validation used at runtime.
type GenesisState struct {
	RouteConfig       *RouteConfig `json:"route_config,omitempty"`
	Paused            bool         `json:"paused"`
	OutstandingAmount string       `json:"outstanding_amount,omitempty"`
	OutboundNonce     uint64       `json:"outbound_nonce,omitempty"`
}

func DefaultGenesisState() GenesisState { return GenesisState{} }

func (g GenesisState) Validate() error {
	if g.RouteConfig == nil {
		if g.Paused || g.OutstandingAmount != "" || g.OutboundNonce != 0 {
			return errors.New("bridge state requires a route configuration")
		}
		return nil
	}
	if err := g.RouteConfig.Validate(); err != nil {
		return err
	}
	if g.OutstandingAmount == "" {
		return nil
	}
	outstanding, ok := new(big.Int).SetString(g.OutstandingAmount, 10)
	if !ok || outstanding.Sign() < 0 {
		return errors.New("invalid bridge outstanding amount")
	}
	limit, _ := new(big.Int).SetString(g.RouteConfig.MaxOutstandingAmount, 10)
	if outstanding.Cmp(limit) > 0 {
		return errors.New("bridge outstanding amount exceeds route limit")
	}
	return nil
}
