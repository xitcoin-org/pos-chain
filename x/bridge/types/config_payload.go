package types

import (
	"crypto/sha256"
	"encoding/json"
)

// RouteConfigPayloadHash produces the exact payload commitment used by a
// signed update_route_config control action.
func RouteConfigPayloadHash(config RouteConfig) ([32]byte, error) {
	if err := config.Validate(); err != nil {
		return [32]byte{}, err
	}
	encoded, err := json.Marshal(config)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(encoded), nil
}
