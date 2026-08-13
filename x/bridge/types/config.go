package types

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// RouteConfig is the testnet route's immutable-format configuration.
// Changes must later be authorized through a verified 2-of-3 message.
type RouteConfig struct {
	RouteID           string   `json:"route_id"`
	BridgeSigners     []string `json:"bridge_signers"`
	Guardian          string   `json:"guardian"`
	MaxTransferAmount string   `json:"max_transfer_amount"`
	DailyLimit        string   `json:"daily_limit"`
	Enabled           bool     `json:"enabled"`
}

func (c RouteConfig) Validate() error {
	if !validRouteID(c.RouteID) {
		return errors.New("invalid route ID")
	}
	if _, err := configuredSigners(c.BridgeSigners); err != nil {
		return err
	}
	if !common.IsHexAddress(strings.TrimSpace(c.Guardian)) {
		return errors.New("invalid bridge guardian")
	}
	guardian := common.HexToAddress(c.Guardian)
	for _, signer := range c.BridgeSigners {
		if guardian == common.HexToAddress(signer) {
			return errors.New("bridge guardian must use a distinct address")
		}
	}
	if !positiveAmount(c.MaxTransferAmount) {
		return errors.New("max transfer amount must be positive")
	}
	if !positiveAmount(c.DailyLimit) {
		return errors.New("daily limit must be positive")
	}
	return nil
}

func positiveAmount(value string) bool {
	amount, ok := new(big.Int).SetString(value, 10)
	return ok && amount.Sign() > 0
}
