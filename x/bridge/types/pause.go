package types

import (
	"crypto/sha256"
	"errors"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const GuardianPauseDomain = "xitcoin-bridge-testnet-guardian-pause-v1"

// GuardianPauseAction can request a suspension only. It cannot change limits,
// resume a route, transfer assets, mint, or access reserves.
type GuardianPauseAction struct {
	RouteID     string `json:"route_id"`
	Nonce       uint64 `json:"nonce"`
	ExpiresUnix int64  `json:"expires_unix"`
}

func (a GuardianPauseAction) Validate() error {
	if !validRouteID(a.RouteID) {
		return errors.New("invalid route ID")
	}
	if a.Nonce == 0 {
		return errors.New("guardian pause nonce must be positive")
	}
	if a.ExpiresUnix <= 0 {
		return errors.New("guardian pause expiry must be positive")
	}
	return nil
}

func (a GuardianPauseAction) ID() ([32]byte, error) {
	if err := a.Validate(); err != nil {
		return [32]byte{}, err
	}
	payload := strings.Join([]string{GuardianPauseDomain, a.RouteID, strconv.FormatUint(a.Nonce, 10), strconv.FormatInt(a.ExpiresUnix, 10)}, "\x00")
	return sha256.Sum256([]byte(payload)), nil
}

func GuardianPauseDigest(action GuardianPauseAction) (common.Hash, error) {
	id, err := action.ID()
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash([]byte(GuardianPauseDomain), id[:]), nil
}

// VerifyGuardianPause verifies a single signature from the dedicated guardian.
func VerifyGuardianPause(action GuardianPauseAction, guardian string, signature []byte) (common.Address, error) {
	if !common.IsHexAddress(strings.TrimSpace(guardian)) {
		return common.Address{}, errors.New("invalid bridge guardian")
	}
	digest, err := GuardianPauseDigest(action)
	if err != nil {
		return common.Address{}, err
	}
	recovered, err := recoverSigner(digest, signature)
	if err != nil {
		return common.Address{}, err
	}
	if recovered != common.HexToAddress(guardian) {
		return common.Address{}, errors.New("pause approval is not from the configured guardian")
	}
	return recovered, nil
}
