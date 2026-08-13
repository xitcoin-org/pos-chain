package keeper

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/xitcoin-org/pos-chain/x/bridge/types"
)

func TestSubmitAttestationRejectsDefaultUnconfiguredBridge(t *testing.T) {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	ctx := sdktestutil.DefaultContext(key, storetypes.NewTransientStoreKey("bridge_msg_disabled_test")).WithBlockTime(time.Unix(1800000000, 0))
	k := NewKeeper(key)
	msg := validBridgeMessage(t, "cronos-testnet-xitcoin-testnet", nil)

	if _, err := NewMsgServer(k).SubmitAttestation(sdk.WrapSDKContext(ctx), &msg); err != ErrRouteDisabled {
		t.Fatalf("unconfigured bridge error = %v, want %v", err, ErrRouteDisabled)
	}
}

func TestSubmitAttestationRecordsOnlyApprovedAttestation(t *testing.T) {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	ctx := sdktestutil.DefaultContext(key, storetypes.NewTransientStoreKey("bridge_msg_admit_test")).WithBlockTime(time.Unix(1800000000, 0))
	k := NewKeeper(key)

	first, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	second, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	third, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	guardian, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	config := types.RouteConfig{
		RouteID:           "cronos-testnet-xitcoin-testnet",
		BridgeSigners:     []string{crypto.PubkeyToAddress(first.PublicKey).Hex(), crypto.PubkeyToAddress(second.PublicKey).Hex(), crypto.PubkeyToAddress(third.PublicKey).Hex()},
		Guardian:          crypto.PubkeyToAddress(guardian.PublicKey).Hex(),
		MaxTransferAmount: "10",
		DailyLimit:        "15",
		Enabled:           true,
	}
	if err := k.SetRouteConfig(ctx, config); err != nil {
		t.Fatal(err)
	}
	msg := validBridgeMessage(t, config.RouteID, nil)
	attestation := msg.Attestation()
	digest, err := types.SigningDigest(attestation)
	if err != nil {
		t.Fatal(err)
	}
	firstSignature, err := crypto.Sign(digest.Bytes(), first)
	if err != nil {
		t.Fatal(err)
	}
	secondSignature, err := crypto.Sign(digest.Bytes(), second)
	if err != nil {
		t.Fatal(err)
	}
	msg.Signatures = [][]byte{firstSignature, secondSignature}

	response, err := NewMsgServer(k).SubmitAttestation(sdk.WrapSDKContext(ctx), &msg)
	if err != nil {
		t.Fatalf("valid submission rejected: %v", err)
	}
	id, err := attestation.ID()
	if err != nil {
		t.Fatal(err)
	}
	if response.AttestationId != hex.EncodeToString(id[:]) {
		t.Fatalf("attestation ID = %s, want %s", response.AttestationId, hex.EncodeToString(id[:]))
	}
	if !k.IsProcessed(ctx, id) {
		t.Fatal("accepted attestation was not replay-protected")
	}
}

func validBridgeMessage(t *testing.T, routeID string, signatures [][]byte) types.MsgSubmitAttestation {
	t.Helper()
	return types.MsgSubmitAttestation{
		Submitter:     sdk.AccAddress(make([]byte, 20)).String(),
		RouteId:       routeID,
		Direction:     string(types.DirectionCronosToXitcoin),
		SourceChainId: "cronos-testnet",
		SourceRef:     "0x" + strings.Repeat("c", 64),
		Nonce:         1,
		Destination:   "xitcoin1testdestination",
		Amount:        "10",
		DeadlineUnix:  1800003600,
		Signatures:    signatures,
	}
}
