package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func testAuthorityAddress() string {
	return sdk.AccAddress(
		[]byte("validator-incentive"),
	).String()
}

func TestMsgUpdateParamsValidateBasic(t *testing.T) {
	params := DefaultParams()
	message := MsgUpdateParams{
		Authority:             testAuthorityAddress(),
		AnnualRateBasisPoints: params.AnnualRateBasisPoints,
		BlocksPerYear:         params.BlocksPerYear,
		RewardPeriodBlocks:    params.RewardPeriodBlocks,
	}

	require.NoError(t, message.ValidateBasic())

	message.Authority = ""
	require.Error(t, message.ValidateBasic())
}

func TestMsgActivateFundedPeriodValidateBasic(t *testing.T) {
	message := MsgActivateFundedPeriod{
		Authority:                   testAuthorityAddress(),
		CommittedAnnualBudgetAtomic: "1000000000000000000",
	}

	require.NoError(t, message.ValidateBasic())

	message.CommittedAnnualBudgetAtomic = "0"
	require.Error(t, message.ValidateBasic())

	message.CommittedAnnualBudgetAtomic = "-1"
	require.Error(t, message.ValidateBasic())

	message.CommittedAnnualBudgetAtomic = "1.5"
	require.Error(t, message.ValidateBasic())
}

func TestNilValidatorIncentiveMessages(t *testing.T) {
	var update *MsgUpdateParams
	require.Error(t, update.ValidateBasic())

	var activate *MsgActivateFundedPeriod
	require.Error(t, activate.ValidateBasic())
}
