package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestAnnualProvision(t *testing.T) {
	xtc := func(value int64) sdkmath.Int {
		return sdkmath.NewInt(value).MulRaw(1_000_000_000_000_000_000)
	}

	tests := []struct {
		name        string
		bonded      sdkmath.Int
		treasury    sdkmath.Int
		cap         sdkmath.Int
		rateBPS     uint32
		expected    sdkmath.Int
		expectError bool
	}{
		{
			name:     "four initial validators at eight percent",
			bonded:   xtc(20_000_000),
			treasury: xtc(100_000_000),
			cap:      xtc(10_000_000),
			rateBPS:  800,
			expected: xtc(1_600_000),
		},
		{
			name:     "lower rate is configurable",
			bonded:   xtc(20_000_000),
			treasury: xtc(100_000_000),
			cap:      xtc(10_000_000),
			rateBPS:  400,
			expected: xtc(800_000),
		},
		{
			name:     "funded annual budget applies",
			bonded:   xtc(200_000_000),
			treasury: xtc(100_000_000),
			cap:      xtc(10_000_000),
			rateBPS:  800,
			expected: xtc(10_000_000),
		},
		{
			name:     "lower annual cap is configurable",
			bonded:   xtc(200_000_000),
			treasury: xtc(100_000_000),
			cap:      xtc(4_000_000),
			rateBPS:  800,
			expected: xtc(4_000_000),
		},
		{
			name:     "two billion bonded scales to eight percent when fully funded",
			bonded:   xtc(2_000_000_000),
			treasury: xtc(200_000_000),
			cap:      xtc(160_000_000),
			rateBPS:  800,
			expected: xtc(160_000_000),
		},
		{
			name:     "governance may raise a fully funded rate to twelve percent",
			bonded:   xtc(200_000_000),
			treasury: xtc(100_000_000),
			cap:      xtc(24_000_000),
			rateBPS:  1_200,
			expected: xtc(24_000_000),
		},
		{
			name:     "funded balance is a hard limit",
			bonded:   xtc(200_000_000),
			treasury: xtc(1_000_000),
			cap:      xtc(10_000_000),
			rateBPS:  800,
			expected: xtc(1_000_000),
		},
		{
			name:     "empty treasury stops provision",
			bonded:   xtc(20_000_000),
			treasury: sdkmath.ZeroInt(),
			cap:      xtc(10_000_000),
			rateBPS:  800,
			expected: sdkmath.ZeroInt(),
		},
		{
			name:        "rate above twenty percent rejected",
			bonded:      xtc(20_000_000),
			treasury:    xtc(100_000_000),
			cap:         xtc(10_000_000),
			rateBPS:     2_001,
			expectError: true,
		},
		{
			name:        "negative treasury rejected",
			bonded:      xtc(20_000_000),
			treasury:    sdkmath.NewInt(-1),
			cap:         xtc(10_000_000),
			rateBPS:     800,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := AnnualProvision(
				test.bonded,
				test.treasury,
				test.cap,
				test.rateBPS,
			)

			if test.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.True(t, actual.Equal(test.expected))
		})
	}
}

func TestCumulativeProvisionHasNoAnnualRoundingDrift(t *testing.T) {
	annual := sdkmath.NewInt(10_000_003)
	blocks := uint64(17)

	final, err := CumulativeProvision(annual, blocks, blocks)
	require.NoError(t, err)
	require.True(t, final.Equal(annual))

	sum := sdkmath.ZeroInt()
	for block := uint64(1); block <= blocks; block++ {
		provision, provisionErr := BlockProvision(annual, block, blocks)
		require.NoError(t, provisionErr)
		require.False(t, provision.IsNegative())
		sum = sum.Add(provision)
	}

	require.True(t, sum.Equal(annual))
}

func TestProvisionValidation(t *testing.T) {
	_, err := CumulativeProvision(sdkmath.OneInt(), 1, 0)
	require.Error(t, err)

	_, err = CumulativeProvision(sdkmath.OneInt(), 2, 1)
	require.Error(t, err)

	_, err = BlockProvision(sdkmath.OneInt(), 0, 1)
	require.Error(t, err)
}
