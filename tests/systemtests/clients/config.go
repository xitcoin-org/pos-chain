package clients

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

const (
	ChainID    = "local-4221"
	EVMChainID = 4221

	JsonRPCUrl0 = "http://127.0.0.1:8545"
	JsonRPCUrl1 = "http://127.0.0.1:8555"
	JsonRPCUrl2 = "http://127.0.0.1:8565"
	JsonRPCUrl3 = "http://127.0.0.1:8575"

	NodeRPCUrl0 = "http://127.0.0.1:26657"
	NodeRPCUrl1 = "http://127.0.0.1:26658"
	NodeRPCUrl2 = "http://127.0.0.1:26659"
	NodeRPCUrl3 = "http://127.0.0.1:26660"
)

type Config struct {
	ChainID     string
	EVMChainID  *big.Int
	PrivKeys    []string
	JsonRPCUrls []string
	NodeRPCUrls []string
}

// NewConfig creates a new Config instance.
func NewConfig() (*Config, error) {
	// Keep system tests aligned with local_node.sh without committing reusable
	// credentials. These deterministic identities are local-test-only.
	privKeys := make([]string, 4)
	for i, name := range []string{"dev0", "dev1", "dev2", "dev3"} {
		digest := sha256.Sum256([]byte("xitcoin-local-test-only:" + name))
		privKeys[i] = hex.EncodeToString(digest[:])
	}

	// jsonrpc urls of testnet nodes
	jsonRPCUrls := []string{JsonRPCUrl0, JsonRPCUrl1, JsonRPCUrl2, JsonRPCUrl3}

	// rpc urls of test nodes
	nodeRPCUrls := []string{NodeRPCUrl0, NodeRPCUrl1, NodeRPCUrl2, NodeRPCUrl3}

	return &Config{
		ChainID:     ChainID,
		EVMChainID:  big.NewInt(EVMChainID),
		PrivKeys:    privKeys,
		JsonRPCUrls: jsonRPCUrls,
		NodeRPCUrls: nodeRPCUrls,
	}, nil
}
