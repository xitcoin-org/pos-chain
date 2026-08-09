package contracts

import (
	contractutils "github.com/xitcoin-org/pos-chain/contracts/utils"
	evmtypes "github.com/xitcoin-org/pos-chain/x/vm/types"
)

func LoadSimpleERC20() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("account_abstraction/tokens/SimpleERC20.json")
}

func LoadSimpleEntryPoint() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("account_abstraction//entrypoint/SimpleEntryPoint.json")
}

func LoadSimpleSmartWallet() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("account_abstraction/smartwallet/SimpleSmartWallet.json")
}
