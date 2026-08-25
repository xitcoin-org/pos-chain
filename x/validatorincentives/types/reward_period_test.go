package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func rewardTestXTC(value int64) sdkmath.Int {
	return sdkmath.NewInt(value).MulRaw(1_000_000_000_000_000_000)
}

func TestValidateFundedRewardPeriodAtTwoBillionBonded(t *testing.T) {
	quarterBlocks := DefaultBlocksPerYear / 4

	provision, err := ValidateFundedRewardPeriod(
		rewardTestXTC(2_000_000_000),
		rewardTestXTC(100_000_000),
		rewardTestXTC(160_000_000),
		800,
		quarterBlocks,
		DefaultBlocksPerYear,
	)
	require.NoError(t, err)
	require.True(t, provision.Equal(rewardTestXTC(40_000_000)))
}

func TestValidateFundedRewardPeriodAtMaximumSupply(t *testing.T) {
	quarterBlocks := DefaultBlocksPerYear / 4

	provision, err := ValidateFundedRewardPeriod(
		rewardTestXTC(5_250_000_000),
		rewardTestXTC(1_000_000_000),
		rewardTestXTC(1_050_000_000),
		2_000,
		quarterBlocks,
		DefaultBlocksPerYear,
	)
	require.NoError(t, err)
	require.True(t, provision.Equal(rewardTestXTC(262_500_000)))
}

func TestFundedRewardPeriodRejectsUnderfunding(t *testing.T) {
	quarterBlocks := DefaultBlocksPerYear / 4

	_, err := ValidateFundedRewardPeriod(
		rewardTestXTC(2_000_000_000),
		rewardTestXTC(100_000_000),
		rewardTestXTC(159_999_999),
		800,
		quarterBlocks,
		DefaultBlocksPerYear,
	)
	require.Error(t, err)

	_, err = ValidateFundedRewardPeriod(
		rewardTestXTC(2_000_000_000),
		rewardTestXTC(39_999_999),
		rewardTestXTC(160_000_000),
		800,
		quarterBlocks,
		DefaultBlocksPerYear,
	)
	require.Error(t, err)
}

func TestValidateRateTransition(t *testing.T) {
	require.NoError(t, ValidateRateTransition(800, 900))
	require.NoError(t, ValidateRateTransition(800, 400))
	require.Error(t, ValidateRateTransition(800, 901))
	require.Error(t, ValidateRateTransition(2_000, 2_001))
}

func TestRewardPeriodValidationBoundaries(t *testing.T) {
	_, err := ValidateFundedRewardPeriod(
		sdkmath.ZeroInt(),
		sdkmath.ZeroInt(),
		sdkmath.ZeroInt(),
		0,
		0,
		DefaultBlocksPerYear,
	)
	require.Error(t, err)

	_, err = ValidateFundedRewardPeriod(
		rewardTestXTC(1),
		rewardTestXTC(1),
		rewardTestXTC(1),
		800,
		DefaultBlocksPerYear+1,
		DefaultBlocksPerYear,
	)
	require.Error(t, err)
}
