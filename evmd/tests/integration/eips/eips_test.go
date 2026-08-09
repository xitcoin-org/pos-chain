package eips_test

import (
	"testing"

	evm "github.com/xitcoin-org/pos-chain"
	"github.com/xitcoin-org/pos-chain/evmd/tests/integration"
	"github.com/xitcoin-org/pos-chain/tests/integration/eips"
	testapp "github.com/xitcoin-org/pos-chain/testutil/app"
	//nolint:revive // dot imports are fine for Ginkgo
	//nolint:revive // dot imports are fine for Ginkgo
)

func TestEIPs(t *testing.T) {
	create := testapp.ToEvmAppCreator[evm.IntegrationNetworkApp](integration.CreateEvmd, "evm.IntegrationNetworkApp")
	eips.RunTests(t, create)
}
