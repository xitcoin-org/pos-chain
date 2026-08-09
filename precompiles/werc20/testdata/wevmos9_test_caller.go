package testdata

import (
	contractutils "github.com/xitcoin-org/pos-chain/contracts/utils"
	evmtypes "github.com/xitcoin-org/pos-chain/x/vm/types"
)

func LoadWEVMOS9TestCaller() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("WEVMOS9TestCaller.json")
}
