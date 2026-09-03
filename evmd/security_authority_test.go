package evmd

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestTestnetAdministrativeAuthority(t *testing.T) {
	authority := mustTestnetAdministrativeAuthority()
	require.Equal(t, "xtc1vza8zsgvrfwmve084ytd8xqdkkm7u9e5csctc2", authority)
	require.Equal(t, TestnetAdministrativeMultisigAddress, authority)

	address, err := sdk.AccAddressFromBech32(authority)
	require.NoError(t, err)
	require.Len(t, address, 20)
	require.NotEqual(t, mustGovernanceAuthority(), authority)
}
