package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	BaseDir       string
	RequireGum    bool
	ConfirmDanger bool
}

func Load(cwd string) Config {
	base := os.Getenv("MAC_CLI_BASE_DIR")
	if base == "" {
		base = cwd
	}
	abs, err := filepath.Abs(base)
	if err == nil {
		base = abs
	}
	return Config{
		BaseDir:       base,
		RequireGum:    os.Getenv("MAC_CLI_REQUIRE_GUM") == "1",
		ConfirmDanger: os.Getenv("MAC_CLI_CONFIRM_DANGER") != "0",
	}
}
