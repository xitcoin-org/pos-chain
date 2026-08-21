package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()

	require.Equal(t, uint32(800), params.AnnualRateBasisPoints)
	require.Equal(t, DefaultBlocksPerYear, params.BlocksPerYear)
	require.Equal(t, DefaultBlocksPerYear/4, params.RewardPeriodBlocks)
	require.NoError(t, params.Validate())
}

func TestParamsValidate(t *testing.T) {
	tests := []struct {
		name        string
		params      Params
		expectError bool
	}{
		{
			name:   "quarterly defaults",
			params: DefaultParams(),
		},
		{
			name: "monthly period",
			params: Params{
				AnnualRateBasisPoints: 800,
				BlocksPerYear:         DefaultBlocksPerYear,
				RewardPeriodBlocks:    DefaultBlocksPerYear / 12,
			},
		},
		{
			name: "rate above ceiling",
			params: Params{
				AnnualRateBasisPoints: 2_001,
				BlocksPerYear:         DefaultBlocksPerYear,
				RewardPeriodBlocks:    DefaultBlocksPerYear / 4,
			},
			expectError: true,
		},
		{
			name: "zero blocks per year",
			params: Params{
				AnnualRateBasisPoints: 800,
				RewardPeriodBlocks:    1,
			},
			expectError: true,
		},
		{
			name: "period does not divide year",
			params: Params{
				AnnualRateBasisPoints: 800,
				BlocksPerYear:         DefaultBlocksPerYear,
				RewardPeriodBlocks:    1_000_000,
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.params.Validate()
			if test.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestParamsValidateUpdate(t *testing.T) {
	previous := DefaultParams()

	next := previous
	next.AnnualRateBasisPoints = 900
	require.NoError(t, next.ValidateUpdate(previous))

	next.AnnualRateBasisPoints = 901
	require.Error(t, next.ValidateUpdate(previous))

	next.AnnualRateBasisPoints = 400
	require.NoError(t, next.ValidateUpdate(previous))
}
