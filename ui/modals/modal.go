package modals

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/deathrashed/mac-cli/ui/styles"
)

func Center(width, height int, content string, theme styles.Theme) string {
	box := theme.Modal.Width(min(width-8, 88)).Render(content)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}

func Help(theme styles.Theme) string {
	lines := []string{
		theme.Accent.Render("Keyboard"),
		"",
		"↑/↓ or j/k     Move selection",
		"←/→ or h/l     Switch category focus",
		"enter          Preview selected command",
		"ctrl+p         Command palette",
		"/              Filter current list",
		"?              Help modal",
		"s              Settings",
		"esc            Close modal/back",
		"q              Quit",
	}
	return strings.Join(lines, "\n")
}

func Preview(command, preview string, dangerous bool, theme styles.Theme) string {
	title := theme.Accent.Render("Command Preview")
	if dangerous {
		title += "  " + theme.Danger.Render("requires confirmation")
	}
	return strings.Join([]string{
		title,
		"",
		preview,
		"",
		theme.Muted.Render("enter: run  esc: cancel"),
	}, "\n")
}

func Settings(baseDir string, gum bool, theme styles.Theme) string {
	gumState := "available"
	if !gum {
		gumState = "not found"
	}
	return strings.Join([]string{
		theme.Accent.Render("Settings"),
		"",
		"Base directory: " + baseDir,
		"Gum fallback:   " + gumState,
		"Theme:          Riley dark / pastel Charmbracelet",
		"",
		theme.Muted.Render("esc: back"),
	}, "\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
