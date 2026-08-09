package gov

import (
	"testing"

	"github.com/stretchr/testify/suite"

	evm "github.com/xitcoin-org/pos-chain"
	"github.com/xitcoin-org/pos-chain/evmd/tests/integration"
	"github.com/xitcoin-org/pos-chain/tests/integration/precompiles/gov"
	testapp "github.com/xitcoin-org/pos-chain/testutil/app"
)

func TestGovPrecompileTestSuite(t *testing.T) {
	create := testapp.ToEvmAppCreator[evm.GovPrecompileApp](integration.CreateEvmd, "evm.GovPrecompileApp")
	s := gov.NewPrecompileTestSuite(create)
	suite.Run(t, s)
}

func TestGovPrecompileIntegrationTestSuite(t *testing.T) {
	create := testapp.ToEvmAppCreator[evm.GovPrecompileApp](integration.CreateEvmd, "evm.GovPrecompileApp")
	gov.TestPrecompileIntegrationTestSuite(t, create)
}
