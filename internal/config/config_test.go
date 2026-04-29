package config

import (
	"os"
	"testing"
)

func TestLoadUsesEnvironmentBaseDir(t *testing.T) {
	t.Setenv("MAC_CLI_BASE_DIR", "/tmp/mac-cli-test")
	cfg := Load(".")
	if cfg.BaseDir != "/tmp/mac-cli-test" {
		t.Fatalf("expected env base dir, got %q", cfg.BaseDir)
	}
}

func TestLoadDangerConfirmationDefault(t *testing.T) {
	os.Unsetenv("MAC_CLI_CONFIRM_DANGER")
	cfg := Load(".")
	if !cfg.ConfirmDanger {
		t.Fatal("danger confirmation should default on")
	}
}
