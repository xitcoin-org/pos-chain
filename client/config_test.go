package client

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
)

func TestInitConfigReturnsNonNotExistError(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config")
	if err := os.WriteFile(configPath, []byte("not a directory"), 0o600); err != nil {
		t.Fatalf("Failed to create config path: %v", err)
	}
	cmd := &cobra.Command{}
	cmd.PersistentFlags().String(flags.FlagHome, "", "")
	if err := cmd.PersistentFlags().Set(flags.FlagHome, tempDir); err != nil {
		t.Fatalf("Could not set home flag [%T] %v", err, err)
	}

	if err := InitConfig(cmd); err == nil || os.IsNotExist(err) {
		t.Fatalf("Expected a non-not-exist error, got: [%T] %v", err, err)
	}
}
