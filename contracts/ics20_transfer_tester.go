package contracts

import (
	contractutils "github.com/xitcoin-org/pos-chain/contracts/utils"
	evmtypes "github.com/xitcoin-org/pos-chain/x/vm/types"
)

func LoadICS20TransferTester() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("solidity/ICS20TransferTester.json")
}
