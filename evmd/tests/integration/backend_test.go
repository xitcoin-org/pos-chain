package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"

	evm "github.com/xitcoin-org/pos-chain"
	"github.com/xitcoin-org/pos-chain/tests/integration/rpc/backend"
	testapp "github.com/xitcoin-org/pos-chain/testutil/app"
)

func TestBackend(t *testing.T) {
	create := testapp.ToEvmAppCreator[evm.IntegrationNetworkApp](CreateEvmd, "evm.IntegrationNetworkApp")
	s := backend.NewTestSuite(create)
	suite.Run(t, s)
}
