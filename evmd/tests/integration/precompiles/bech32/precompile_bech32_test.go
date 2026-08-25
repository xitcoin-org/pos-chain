package bech32

import (
	"testing"

	"github.com/stretchr/testify/suite"

	evm "github.com/xitcoin-org/pos-chain"
	"github.com/xitcoin-org/pos-chain/evmd/tests/integration"
	"github.com/xitcoin-org/pos-chain/tests/integration/precompiles/bech32"
	testapp "github.com/xitcoin-org/pos-chain/testutil/app"
)

func TestBech32PrecompileTestSuite(t *testing.T) {
	create := testapp.ToEvmAppCreator[evm.Bech32PrecompileApp](integration.CreateEvmd, "evm.Bech32PrecompileApp")
	s := bech32.NewPrecompileTestSuite(create)
	suite.Run(t, s)
}
