package evmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBridgeInitialRouteUsesAdministrativeAuthority(t *testing.T) {
	app, _ := setup(true, 0, "xitcoin-bridge-authority-test", 101089)

	require.Equal(t, TestnetAdministrativeMultisigAddress, app.BridgeKeeper.Authority())
	require.NotEqual(t, mustGovernanceAuthority(), app.BridgeKeeper.Authority())
}
