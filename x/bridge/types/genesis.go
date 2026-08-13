package types

import "errors"

// GenesisState deliberately contains no enabled default route. A route may be
// present only when it passes the same validation used at runtime.
type GenesisState struct {
	RouteConfig *RouteConfig `json:"route_config,omitempty"`
	Paused      bool         `json:"paused"`
}

func DefaultGenesisState() GenesisState { return GenesisState{} }

func (g GenesisState) Validate() error {
	if g.RouteConfig == nil {
		if g.Paused {
			return errors.New("bridge cannot be paused without a route configuration")
		}
		return nil
	}
	return g.RouteConfig.Validate()
}
