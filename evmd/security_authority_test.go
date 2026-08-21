package evmd

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKCALBAdministrativeAuthority(t *testing.T) {
	authority := mustKCALBAdministrativeAuthority()
	require.Equal(t, KCALBAdministrativeMultisigAddress, authority)

	address, err := sdk.AccAddressFromBech32(authority)
	require.NoError(t, err)
	require.Len(t, address, 20)
	require.NotEqual(t, mustGovernanceAuthority(), authority)
}
