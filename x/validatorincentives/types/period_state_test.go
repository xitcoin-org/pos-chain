package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPeriodState(t *testing.T) {
	params := DefaultParams()

	state, err := NewPeriodState(
		100,
		rewardTestXTC(2_000_000_000),
		rewardTestXTC(100_000_000),
		rewardTestXTC(160_000_000),
		params,
	)
	require.NoError(t, err)
	require.Equal(t, uint64(100)+params.RewardPeriodBlocks, state.EndBlock)
	require.Equal(t, "40000000000000000000000000", state.PeriodProvisionAtomic)
	require.NoError(t, state.Validate())

	remaining, err := state.RemainingProvision()
	require.NoError(t, err)
	require.True(t, remaining.Equal(rewardTestXTC(40_000_000)))
}

func TestPeriodStateValidation(t *testing.T) {
	valid, err := NewPeriodState(
		100,
		rewardTestXTC(20_000_000),
		rewardTestXTC(100_000_000),
		rewardTestXTC(1_600_000),
		DefaultParams(),
	)
	require.NoError(t, err)

	invalid := valid
	invalid.DistributedAtomic = "400000000000000000000001"
	require.Error(t, invalid.Validate())

	invalid = valid
	invalid.EligibleBondedAtomic = "not-an-integer"
	require.Error(t, invalid.Validate())

	invalid = valid
	invalid.EndBlock = invalid.StartBlock
	require.Error(t, invalid.Validate())
}
