package config

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	GethVersion = "1.16.3"
	EvmdHttpEndpoint = "http://localhost:8545"
	EvmdWsEndpoint   = "ws://localhost:8546"
	GethHttpEndpoint = "http://localhost:8547"
	GethWsEndpoint   = "ws://localhost:8548"
)

var (
	Dev0PrivateKey = localTestPrivateKey("dev0")
	Dev1PrivateKey = localTestPrivateKey("dev1")
	Dev2PrivateKey = localTestPrivateKey("dev2")
	Dev3PrivateKey = localTestPrivateKey("dev3")
)

func localTestPrivateKey(name string) string {
	digest := sha256.Sum256([]byte("xitcoin-local-test-only:" + name))
	return hex.EncodeToString(digest[:])
}

type Config struct {
	EvmdHttpEndpoint string `yaml:"evmd_http_endpoint"`
	EvmdWsEndpoint   string `yaml:"evmd_ws_endpoint"`
	GethHttpEndpoint string `yaml:"geth_http_endpoint"`
	GethWsEndpoint   string `yaml:"geth_ws_endpoint"`

	RichPrivKey string `yaml:"rich_privkey"`
	// Timeout is the timeout for the RPC (e.g. 5s, 1m)
	Timeout string `yaml:"timeout"`
}

func (c *Config) Validate() error {
	if c.EvmdHttpEndpoint == "" {
		return fmt.Errorf("rpc_endpoint must be set")
	}
	if c.EvmdWsEndpoint == "" {
		return fmt.Errorf("ws_endpoint must be set")
	}
	if c.GethHttpEndpoint == "" {
		return fmt.Errorf("geth_http_endpoint must be set")
	}
	if c.GethWsEndpoint == "" {
		return fmt.Errorf("geth_ws_endpoint must be set")
	}

	if c.RichPrivKey == "" {
		return fmt.Errorf("rich_privkey must be set")
	}
	if _, err := time.ParseDuration(c.Timeout); err != nil {
		return fmt.Errorf("invalid timeout: %v", err)
	}
	return nil
}

func MustLoadConfig() *Config {
	// Use environment variable if set, otherwise default to localhost
	evmdURL := os.Getenv("EVMD_URL")
	if evmdURL == "" {
		evmdURL = EvmdHttpEndpoint
	}

	gethURL := os.Getenv("GETH_URL")
	if gethURL == "" {
		gethURL = GethHttpEndpoint
	}

	// Handle WebSocket URLs - derive from HTTP URLs or use environment variables
	evmdWsURL := os.Getenv("EVMD_WS_URL")
	if evmdWsURL == "" {
		evmdWsURL = EvmdWsEndpoint
	}

	gethWsURL := os.Getenv("GETH_WS_URL")
	if gethWsURL == "" {
		gethWsURL = GethWsEndpoint
	}

	return &Config{
		EvmdHttpEndpoint: evmdURL,
		EvmdWsEndpoint:   evmdWsURL,
		GethHttpEndpoint: gethURL,
		GethWsEndpoint:   gethWsURL,
		RichPrivKey:      Dev0PrivateKey, // Default to dev0's private key
		Timeout:          "10s",
	}
}

// GetDev0PrivateKeyAndAddress returns dev0's private key and address for contract deployment
func GetDev0PrivateKeyAndAddress() (*ecdsa.PrivateKey, common.Address, error) {
	privateKey, err := crypto.HexToECDSA(Dev0PrivateKey)
	if err != nil {
		return nil, common.Address{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, fmt.Errorf("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, address, nil
}
