package ante

import (
	"testing"

	evm "github.com/xitcoin-org/pos-chain"
	"github.com/stretchr/testify/suite"

	"github.com/xitcoin-org/pos-chain/evmd/tests/integration"
	"github.com/xitcoin-org/pos-chain/tests/integration/ante"
	testapp "github.com/xitcoin-org/pos-chain/testutil/app"
)

func TestEvmUnitAnteTestSuite(t *testing.T) {
	create := testapp.ToEvmAppCreator[evm.AnteIntegrationApp](integration.CreateEvmd, "evm.AnteIntegrationApp")
	suite.Run(t, ante.NewEvmUnitAnteTestSuite(create))
}
