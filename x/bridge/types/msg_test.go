package types

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func validSubmitMessage() MsgSubmitAttestation {
	return MsgSubmitAttestation{
		Submitter:     sdk.AccAddress(make([]byte, 20)).String(),
		RouteId:       "cronos-testnet-xitcoin-testnet",
		Direction:     string(DirectionCronosToXitcoin),
		SourceChainId: "cronos-testnet",
		SourceRef:     "0x" + strings.Repeat("a", 64),
		Nonce:         1,
		Destination:   "xitcoin1testdestination",
		Amount:        "1",
		DeadlineUnix:  1800003600,
	}
}

func TestMsgSubmitAttestationValidateBasic(t *testing.T) {
	msg := validSubmitMessage()
	if err := msg.ValidateBasic(); err != nil {
		t.Fatalf("valid submission rejected: %v", err)
	}
	msg.Submitter = "not-an-address"
	if err := msg.ValidateBasic(); err == nil {
		t.Fatal("invalid submitter accepted")
	}
}

func TestMsgInitiateOutboundTransferValidateBasic(t *testing.T) {
	msg := MsgInitiateOutboundTransfer{
		Sender:      sdk.AccAddress(make([]byte, 20)).String(),
		RouteId:     "cronos-testnet-xitcoin-testnet",
		Destination: "0x0000000000000000000000000000000000000001",
		Amount:      "1",
	}
	if err := msg.ValidateBasic(); err != nil {
		t.Fatalf("valid outbound transfer rejected: %v", err)
	}
	msg.Sender = "not-an-address"
	if err := msg.ValidateBasic(); err == nil {
		t.Fatal("invalid outbound sender accepted")
	}
}
