package server

import (
	"testing"

	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/stretchr/testify/require"
)

func TestStartCmdDisablesInterBlockCacheByDefault(t *testing.T) {
	cmd := StartCmd(StartOptions{
		DefaultNodeHome: t.TempDir(),
	})

	flag := cmd.Flags().Lookup(sdkserver.FlagInterBlockCache)

	require.NotNil(t, flag)
	require.Equal(t, "false", flag.DefValue)
}
