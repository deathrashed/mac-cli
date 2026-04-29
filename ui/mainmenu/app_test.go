package mainmenu

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/deathrashed/mac-cli/internal/commands"
	"github.com/deathrashed/mac-cli/internal/services"
)

func TestBaseViewFitsApprox1000x750Terminal(t *testing.T) {
	m := New(commands.DefaultCatalog(), services.NewRunner("."))
	m.width = 125
	m.height = 41
	m.resize()

	view := m.baseView()
	if got := lipgloss.Width(view); got > m.width {
		t.Fatalf("view width = %d, want <= %d", got, m.width)
	}
	if got := lipgloss.Height(view); got > m.height {
		t.Fatalf("view height = %d, want <= %d", got, m.height)
	}
}

func TestLayoutScalesDownForNarrowTerminal(t *testing.T) {
	m := New(commands.DefaultCatalog(), services.NewRunner("."))
	m.width = 80
	m.height = 24
	sidebar, content, height := m.layout()
	if sidebar+content+7 > m.width {
		t.Fatalf("layout overflows: sidebar=%d content=%d width=%d", sidebar, content, m.width)
	}
	if height+6 > m.height {
		t.Fatalf("layout height overflows: panel=%d terminal=%d", height, m.height)
	}
}

func TestCleanOutputLineNormalizesLegacyEscapes(t *testing.T) {
	got := cleanOutputLine(`\033[0;32mVol:\033[0m 45%\n`)
	want := "\033[0;32mVol:\033[0m 45%"
	if got != want {
		t.Fatalf("cleanOutputLine() = %q, want %q", got, want)
	}
}

func TestStripANSIForCopyAndSave(t *testing.T) {
	got := stripANSI([]string{"\033[0;32mgreen\033[0m", "plain"})
	if got[0] != "green" || got[1] != "plain" {
		t.Fatalf("stripANSI() = %#v", got)
	}
}

func TestSanitizeFilename(t *testing.T) {
	if got := sanitizeFilename("network:dns/flush"); got != "network-dns-flush" {
		t.Fatalf("sanitizeFilename() = %q", got)
	}
}
