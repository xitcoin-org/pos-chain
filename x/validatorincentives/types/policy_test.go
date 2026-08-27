package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestLaunchReferenceCalculation(t *testing.T) {
	treasury := sdkmath.NewInt(20_000_000)
	stake := sdkmath.NewInt(20_000_000)

	annualized, err := AnnualizedRewardCapacity(
		treasury,
		DefaultTreasuryReleaseRateBasisPoints,
	)
	require.NoError(t, err)
	require.Equal(t, "2000000", annualized.String())

	apy, err := DerivedAPYBasisPoints(annualized, stake)
	require.NoError(t, err)
	require.Equal(t, "1000", apy.String())

	daily, err := PeriodProvision(
		annualized,
		DefaultCalculationPeriodBlocks,
		DefaultBlocksPerYear,
	)
	require.NoError(t, err)
	require.Equal(t, "5475", daily.String())
}

func TestFundingAndStakeChangeDerivedAPY(t *testing.T) {
	annualized, err := AnnualizedRewardCapacity(sdkmath.NewInt(30_000_000), 1_000)
	require.NoError(t, err)
	apy, err := DerivedAPYBasisPoints(annualized, sdkmath.NewInt(20_000_000))
	require.NoError(t, err)
	require.Equal(t, "1500", apy.String())

	apy, err = DerivedAPYBasisPoints(annualized, sdkmath.NewInt(40_000_000))
	require.NoError(t, err)
	require.Equal(t, "750", apy.String())
}

func TestPolicyRejectsInvalidValues(t *testing.T) {
	_, err := AnnualizedRewardCapacity(sdkmath.NewInt(-1), 1_000)
	require.Error(t, err)
	_, err = AnnualizedRewardCapacity(sdkmath.OneInt(), 10_001)
	require.Error(t, err)
	_, err = PeriodProvision(sdkmath.OneInt(), 0, 1)
	require.Error(t, err)
	_, err = DerivedAPYBasisPoints(sdkmath.OneInt(), sdkmath.ZeroInt())
	require.NoError(t, err)
}

func TestCumulativeProvisionHasNoRoundingDrift(t *testing.T) {
	provision := sdkmath.NewInt(101)
	previous := sdkmath.ZeroInt()
	sum := sdkmath.ZeroInt()
	for block := uint64(1); block <= 7; block++ {
		target, err := CumulativeProvision(provision, block, 7)
		require.NoError(t, err)
		sum = sum.Add(target.Sub(previous))
		previous = target
	}
	require.True(t, sum.Equal(provision))
}

func TestDefaultParamsAreDaily(t *testing.T) {
	params := DefaultParams()
	require.Equal(t, uint32(1_000), params.TreasuryReleaseRateBasisPoints)
	require.Equal(t, uint64(17_280), params.CalculationPeriodBlocks)
	require.NoError(t, params.Validate())

	params.TreasuryReleaseRateBasisPoints = 10_001
	require.Error(t, params.Validate())
}
