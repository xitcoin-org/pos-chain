package types

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ControlDomain keeps governance signatures separate from transfer approvals.
const ControlDomain = "xitcoin-bridge-testnet-control-v1"

const (
	ActionUpdateRouteConfig = "update_route_config"
	ActionPauseRoute        = "pause_route"
	ActionResumeRoute       = "resume_route"
	ActionCloseRoute        = "close_route"
)

// ControlAction describes a future route-governance action. PayloadHash binds
// an action to its exact, separately canonicalised payload. This primitive has
// no token, bank, reserve, minting, relayer, or transaction capability.
type ControlAction struct {
	RouteID       string `json:"route_id"`
	Action        string `json:"action"`
	PayloadHash   string `json:"payload_hash"`
	Nonce         uint64 `json:"nonce"`
	NotBeforeUnix int64  `json:"not_before_unix"`
	ExpiresUnix   int64  `json:"expires_unix"`
}

func (a ControlAction) Validate() error {
	if !validRouteID(a.RouteID) {
		return errors.New("invalid route ID")
	}
	switch a.Action {
	case ActionUpdateRouteConfig, ActionPauseRoute, ActionResumeRoute, ActionCloseRoute:
	default:
		return errors.New("invalid bridge control action")
	}
	if !common.IsHexHash(strings.TrimSpace(a.PayloadHash)) {
		return errors.New("invalid bridge control payload hash")
	}
	if a.Nonce == 0 {
		return errors.New("bridge control nonce must be positive")
	}
	if a.NotBeforeUnix <= 0 || a.ExpiresUnix <= a.NotBeforeUnix {
		return errors.New("invalid bridge control validity window")
	}
	return nil
}

// ID is deterministic and contains every signed governance field.
func (a ControlAction) ID() ([32]byte, error) {
	if err := a.Validate(); err != nil {
		return [32]byte{}, err
	}
	payload := strings.Join([]string{
		ControlDomain,
		a.RouteID,
		a.Action,
		strings.ToLower(strings.TrimSpace(a.PayloadHash)),
		strconv.FormatUint(a.Nonce, 10),
		strconv.FormatInt(a.NotBeforeUnix, 10),
		strconv.FormatInt(a.ExpiresUnix, 10),
	}, "\x00")
	return sha256.Sum256([]byte(payload)), nil
}

func ControlDigest(action ControlAction) (common.Hash, error) {
	id, err := action.ID()
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash([]byte(ControlDomain), id[:]), nil
}

// VerifyControlApprovals requires two distinct current bridge signers. It does
// not apply the action; execution will be added only after module integration.
func VerifyControlApprovals(action ControlAction, bridgeSigners []string, signatures [][]byte) ([]common.Address, error) {
	allowed, err := configuredSigners(bridgeSigners)
	if err != nil {
		return nil, err
	}
	if len(signatures) < RequiredApprovals || len(signatures) > MaxBridgeSigners {
		return nil, fmt.Errorf("signature count must be between %d and %d", RequiredApprovals, MaxBridgeSigners)
	}
	digest, err := ControlDigest(action)
	if err != nil {
		return nil, err
	}

	recovered := make([]common.Address, 0, len(signatures))
	seen := make(map[common.Address]struct{}, len(signatures))
	for _, signature := range signatures {
		address, err := recoverSigner(digest, signature)
		if err != nil {
			return nil, err
		}
		if _, ok := allowed[address]; !ok {
			return nil, errors.New("approval is not from a configured bridge signer")
		}
		if _, exists := seen[address]; exists {
			return nil, errors.New("duplicate signer approval")
		}
		seen[address] = struct{}{}
		recovered = append(recovered, address)
	}
	return recovered, nil
}
