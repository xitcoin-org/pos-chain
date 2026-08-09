package validatoradmission

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/xitcoin-org/pos-chain/x/validatoradmission/keeper"
)

func NewAdmissionAnteHandler(admissionKeeper keeper.Keeper, next sdk.AnteHandler) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		for _, message := range tx.GetMsgs() {
			switch msg := message.(type) {
			case *stakingtypes.MsgCreateValidator:
				if !admissionKeeper.IsApprovedValidator(ctx, msg.ValidatorAddress) {
					return ctx, errorsmod.Wrap(
						sdkerrors.ErrUnauthorized,
						"Validator Admission: validator is not approved by Xitcoin",
					)
				}

			case *slashingtypes.MsgUnjail:
				if !admissionKeeper.IsApprovedValidator(ctx, msg.ValidatorAddr) {
					return ctx, errorsmod.Wrap(
						sdkerrors.ErrUnauthorized,
						"Validator Admission: revoked validator cannot unjail",
					)
				}

			case *minttypes.MsgUpdateParams:
				return ctx, errorsmod.Wrap(
					sdkerrors.ErrUnauthorized,
					"Fixed monetary policy: Mint parameters cannot be changed",
				)

			case *govtypes.MsgSubmitProposal:
				for _, proposedMessage := range msg.Messages {
					if proposedMessage != nil &&
						proposedMessage.TypeUrl == "/cosmos.mint.v1beta1.MsgUpdateParams" {
						return ctx, errorsmod.Wrap(
							sdkerrors.ErrUnauthorized,
							"Fixed monetary policy: Mint parameter proposals are disabled",
						)
					}
				}
			}
		}

		return next(ctx, tx, simulate)
	}
}
