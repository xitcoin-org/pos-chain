package types

import "testing"

func validRouteConfig() RouteConfig {
	return RouteConfig{
		RouteID: "cronos-testnet-xitcoin-testnet",
		BridgeSigners: []string{
			"0x0000000000000000000000000000000000000001",
			"0x0000000000000000000000000000000000000002",
			"0x0000000000000000000000000000000000000003",
		},
		Guardian:          "0x0000000000000000000000000000000000000004",
		MaxTransferAmount: "1000000000000000000",
		DailyLimit:        "5000000000000000000",
		Enabled:           false,
	}
}

func TestRouteConfigValidation(t *testing.T) {
	if err := validRouteConfig().Validate(); err != nil {
		t.Fatalf("valid route config rejected: %v", err)
	}
	c := validRouteConfig()
	c.Guardian = c.BridgeSigners[0]
	if err := c.Validate(); err == nil {
		t.Fatal("guardian shared with signer was accepted")
	}
	c = validRouteConfig()
	c.BridgeSigners = c.BridgeSigners[:2]
	if err := c.Validate(); err == nil {
		t.Fatal("two configured signers were accepted")
	}
	c = validRouteConfig()
	c.DailyLimit = "0"
	if err := c.Validate(); err == nil {
		t.Fatal("zero daily limit was accepted")
	}
}
