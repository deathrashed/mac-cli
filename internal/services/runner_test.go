package services

import (
	"context"
	"strings"
	"testing"

	"github.com/deathrashed/mac-cli/internal/commands"
)

func TestNativeTemperatureConversion(t *testing.T) {
	cmd, _ := commands.DefaultCatalog().Find("convert:c-to-f")
	got, err := NewRunner(".").Run(context.Background(), cmd, []string{"25"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got.Output, "77.00") {
		t.Fatalf("unexpected output: %q", got.Output)
	}
}

func TestNativeURLEncode(t *testing.T) {
	cmd, _ := commands.DefaultCatalog().Find("text:urlencode")
	got, err := NewRunner(".").Run(context.Background(), cmd, []string{"hello world"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(got.Output) != "hello+world" {
		t.Fatalf("unexpected output: %q", got.Output)
	}
}

func TestPreviewShowsUserFacingCommand(t *testing.T) {
	cmd, _ := commands.DefaultCatalog().Find("system:ports")
	preview := NewRunner("/tmp/mac-cli").Preview(cmd, nil)
	if preview != "mac system:ports" {
		t.Fatalf("preview = %q", preview)
	}
}

func TestLegacyNameMapsAdvertisedAliases(t *testing.T) {
	cases := map[string]string{
		"ports":             "network:ports",
		"speedtest":         "network:speedtest",
		"ip:local":          "network:ip:local",
		"ip:public":         "network:public-ip",
		"git:branches:date": "git:branches",
		"git:settings":      "git:config",
		"memory":            "memory",
	}
	for input, want := range cases {
		if got := legacyName(input); got != want {
			t.Fatalf("legacyName(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestConfirmInputAnswersYesOnlyForDangerousCommands(t *testing.T) {
	if got := confirmInput(commands.Command{Name: "trash:empty", Dangerous: true}); got != "y\n" {
		t.Fatalf("dangerous confirm input = %q", got)
	}
	if got := confirmInput(commands.Command{Name: "ports"}); got != "n\n" {
		t.Fatalf("safe confirm input = %q", got)
	}
}

func TestScanOutputLineSplitsCarriageReturns(t *testing.T) {
	advance, token, err := scanOutputLine([]byte("progress 50%\rnext"), false)
	if err != nil {
		t.Fatal(err)
	}
	if advance != len("progress 50%\r") {
		t.Fatalf("advance = %d", advance)
	}
	if string(token) != "progress 50%" {
		t.Fatalf("token = %q", token)
	}
}
