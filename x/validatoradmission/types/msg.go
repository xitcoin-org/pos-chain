package types

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgApproveValidator{}
	_ sdk.Msg = &MsgRevokeValidator{}
)

func validateAuthorityAndValidator(authority, validatorAddress string) error {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		return errorsmod.Wrap(err, "invalid Xitcoin authority address")
	}

	if _, err := sdk.ValAddressFromBech32(validatorAddress); err != nil {
		return errorsmod.Wrap(err, "invalid Xitcoin validator address")
	}

	return nil
}

// ValidateBasic validates an approval request.
func (m *MsgApproveValidator) ValidateBasic() error {
	return validateAuthorityAndValidator(m.Authority, m.ValidatorAddress)
}

// GetSignBytes returns legacy sign bytes.
func (m MsgApproveValidator) GetSignBytes() []byte {
	return AminoCdc.MustMarshalJSON(&m)
}

// ValidateBasic validates a revocation request.
func (m *MsgRevokeValidator) ValidateBasic() error {
	return validateAuthorityAndValidator(m.Authority, m.ValidatorAddress)
}

// GetSignBytes returns legacy sign bytes.
func (m MsgRevokeValidator) GetSignBytes() []byte {
	return AminoCdc.MustMarshalJSON(&m)
}
