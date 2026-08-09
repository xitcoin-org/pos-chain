package keeper

import (
	"bytes"
	"context"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	ed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/xitcoin-org/pos-chain/x/validatoradmission/types"
)

type testStakingKeeper struct {
	getCalls    int
	validator   stakingtypes.Validator
	getErr      error
	jailed      bool
	jailAddress sdk.ConsAddress
}

func (s *testStakingKeeper) GetValidator(
	_ context.Context,
	_ sdk.ValAddress,
) (stakingtypes.Validator, error) {
	s.getCalls++
	if s.getErr != nil {
		return stakingtypes.Validator{}, s.getErr
	}
	return s.validator, nil
}

func (s *testStakingKeeper) Jail(
	_ context.Context,
	consAddress sdk.ConsAddress,
) error {
	s.jailed = true
	s.jailAddress = consAddress
	return nil
}

func TestApproveAndRevokeValidator(t *testing.T) {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	ctx := sdktestutil.DefaultContext(
		key,
		storetypes.NewTransientStoreKey("validator_admission_test"),
	)

	authority := sdk.AccAddress(bytes.Repeat([]byte{1}, 20)).String()
	wrongAuthority := sdk.AccAddress(bytes.Repeat([]byte{2}, 20)).String()
	validator := sdk.ValAddress(bytes.Repeat([]byte{3}, 20)).String()

	k := NewKeeper(key)
	k.SetAuthority(ctx, authority)

	staking := &testStakingKeeper{getErr: stakingtypes.ErrNoValidatorFound}
	server := NewMsgServer(k, staking)
	goCtx := sdk.WrapSDKContext(ctx)

	_, err := server.ApproveValidator(goCtx, &types.MsgApproveValidator{
		Authority:        wrongAuthority,
		ValidatorAddress: validator,
	})
	if err == nil {
		t.Fatal("approval from a non-authority was accepted")
	}
	if k.IsApprovedValidator(ctx, validator) {
		t.Fatal("non-authority changed the approval list")
	}

	_, err = server.ApproveValidator(goCtx, &types.MsgApproveValidator{
		Authority:        authority,
		ValidatorAddress: validator,
	})
	if err != nil {
		t.Fatalf("authority approval rejected: %v", err)
	}
	if !k.IsApprovedValidator(ctx, validator) {
		t.Fatal("approved validator missing from on-chain state")
	}

	_, err = server.RevokeValidator(goCtx, &types.MsgRevokeValidator{
		Authority:        authority,
		ValidatorAddress: validator,
	})
	if err != nil {
		t.Fatalf("authority revocation rejected: %v", err)
	}
	if k.IsApprovedValidator(ctx, validator) {
		t.Fatal("revoked validator remained approved")
	}
	if staking.getCalls != 1 {
		t.Fatalf("expected one validator lookup, got %d", staking.getCalls)
	}
}

func TestRevokeExistingValidatorSuspendsIt(t *testing.T) {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	ctx := sdktestutil.DefaultContext(
		key,
		storetypes.NewTransientStoreKey("validator_admission_revoke_test"),
	)

	authority := sdk.AccAddress(bytes.Repeat([]byte{11}, 20)).String()
	validatorAddress := sdk.ValAddress(bytes.Repeat([]byte{12}, 20)).String()

	consensusKey := &ed25519.PubKey{Key: bytes.Repeat([]byte{13}, 32)}
	consensusAny, err := codectypes.NewAnyWithValue(consensusKey)
	if err != nil {
		t.Fatalf("consensus key packaging failed: %v", err)
	}

	validator := stakingtypes.Validator{
		OperatorAddress: validatorAddress,
		ConsensusPubkey: consensusAny,
	}

	k := NewKeeper(key)
	k.SetAuthority(ctx, authority)
	k.SetApprovedValidator(ctx, validatorAddress, true)

	staking := &testStakingKeeper{validator: validator}
	server := NewMsgServer(k, staking)

	_, err = server.RevokeValidator(
		sdk.WrapSDKContext(ctx),
		&types.MsgRevokeValidator{
			Authority:        authority,
			ValidatorAddress: validatorAddress,
		},
	)
	if err != nil {
		t.Fatalf("revocation rejected: %v", err)
	}

	if k.IsApprovedValidator(ctx, validatorAddress) {
		t.Fatal("revoked validator remained approved")
	}
	if !staking.jailed {
		t.Fatal("existing revoked validator was not suspended")
	}

	expected := sdk.ConsAddress(consensusKey.Address())
	if !bytes.Equal(staking.jailAddress, expected) {
		t.Fatalf("wrong consensus address suspended: got %X want %X",
			staking.jailAddress, expected)
	}
}
