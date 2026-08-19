package types

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func commandAuthority() string {
	return sdk.AccAddress(bytes.Repeat([]byte{1}, 20)).String()
}

func TestUpdateParamsCommandValidation(t *testing.T) {
	command := UpdateParamsCommand{
		Authority: commandAuthority(),
		Params:    DefaultParams(),
	}
	require.NoError(t, command.ValidateBasic())

	command.Authority = ""
	require.ErrorContains(t, command.ValidateBasic(), "required")

	command.Authority = "not-an-address"
	require.ErrorContains(t, command.ValidateBasic(), "invalid authority")

	command.Authority = commandAuthority()
	command.Params.AnnualRateBasisPoints =
		MaxAnnualRateBasisPoints + 1
	require.ErrorContains(t, command.ValidateBasic(), "protocol ceiling")
}

func TestActivatePeriodCommandValidation(t *testing.T) {
	command := ActivatePeriodCommand{
		Authority:                   commandAuthority(),
		EligibleBondedAtomic:         "2000000000000000000000000000",
		TreasuryBalanceAtomic:        "100000000000000000000000000",
		CommittedAnnualBudgetAtomic: "160000000000000000000000000",
	}
	require.NoError(t, command.ValidateBasic())

	eligible, treasury, budget, err := command.Amounts()
	require.NoError(t, err)
	require.True(t, eligible.IsPositive())
	require.True(t, treasury.IsPositive())
	require.True(t, budget.IsPositive())
}

func TestActivatePeriodCommandRejectsInvalidAmounts(t *testing.T) {
	base := ActivatePeriodCommand{
		Authority:                   commandAuthority(),
		EligibleBondedAtomic:         "1",
		TreasuryBalanceAtomic:        "1",
		CommittedAnnualBudgetAtomic: "1",
	}

	tests := []struct {
		name   string
		mutate func(*ActivatePeriodCommand)
	}{
		{
			name: "invalid eligible amount",
			mutate: func(c *ActivatePeriodCommand) {
				c.EligibleBondedAtomic = "1.5"
			},
		},
		{
			name: "negative treasury",
			mutate: func(c *ActivatePeriodCommand) {
				c.TreasuryBalanceAtomic = "-1"
			},
		},
		{
			name: "zero budget",
			mutate: func(c *ActivatePeriodCommand) {
				c.CommittedAnnualBudgetAtomic = "0"
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := base
			test.mutate(&command)
			require.Error(t, command.ValidateBasic())
		})
	}
}
