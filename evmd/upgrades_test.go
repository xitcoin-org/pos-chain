package evmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpgradeNamesAreStableAndUnique(t *testing.T) {
	names := []string{
		UpgradeName,
		GovernanceSafeguardsUpgradeName,
		ValidatorIncentivesDailyV2UpgradeName,
	}

	require.Equal(t, "xitcoin-validator-incentives-daily-v2", ValidatorIncentivesDailyV2UpgradeName)
	require.NotEmpty(t, UpgradeName)
	require.NotEmpty(t, GovernanceSafeguardsUpgradeName)
	require.Len(t, names, 3)
	require.NotEqual(t, names[0], names[1])
	require.NotEqual(t, names[0], names[2])
	require.NotEqual(t, names[1], names[2])
}
