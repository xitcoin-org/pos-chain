package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"

	evm "github.com/xitcoin-org/pos-chain"
	"github.com/xitcoin-org/pos-chain/tests/integration/x/feemarket"
	testapp "github.com/xitcoin-org/pos-chain/testutil/app"
)

func TestFeeMarketKeeperTestSuite(t *testing.T) {
	create := testapp.ToEvmAppCreator[evm.IntegrationNetworkApp](CreateEvmd, "evm.IntegrationNetworkApp")
	s := feemarket.NewTestKeeperTestSuite(create)
	suite.Run(t, s)
}
