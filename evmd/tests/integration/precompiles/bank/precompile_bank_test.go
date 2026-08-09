package bank

import (
	"testing"

	"github.com/stretchr/testify/suite"

	evm "github.com/xitcoin-org/pos-chain"
	"github.com/xitcoin-org/pos-chain/evmd/tests/integration"
	"github.com/xitcoin-org/pos-chain/tests/integration/precompiles/bank"
	testapp "github.com/xitcoin-org/pos-chain/testutil/app"
)

func TestBankPrecompileTestSuite(t *testing.T) {
	create := testapp.ToEvmAppCreator[evm.BankPrecompileApp](integration.CreateEvmd, "evm.BankPrecompileApp")
	s := bank.NewPrecompileTestSuite(create)
	suite.Run(t, s)
}

func TestBankPrecompileIntegrationTestSuite(t *testing.T) {
	create := testapp.ToEvmAppCreator[evm.BankPrecompileApp](integration.CreateEvmd, "evm.BankPrecompileApp")
	bank.TestIntegrationSuite(t, create)
}
