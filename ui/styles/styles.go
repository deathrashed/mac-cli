package styles

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Doc       lipgloss.Style
	Header    lipgloss.Style
	Bar       lipgloss.Style
	Tab       lipgloss.Style
	ActiveTab lipgloss.Style
	PaneTitle lipgloss.Style
	Panel     lipgloss.Style
	Sidebar   lipgloss.Style
	Detail    lipgloss.Style
	Status    lipgloss.Style
	Nugget    lipgloss.Style
	Accent    lipgloss.Style
	Muted     lipgloss.Style
	Danger    lipgloss.Style
	Modal     lipgloss.Style
	Button    lipgloss.Style
	ActiveBtn lipgloss.Style
}

func New() Theme {
	subtle := lipgloss.Color("#5f6472")
	accent := lipgloss.Color("#9B7CFF")
	cyan := lipgloss.Color("#56D4DD")
	pink := lipgloss.Color("#FF6FAE")
	bg := lipgloss.Color("#1e1e1e")
	panel := lipgloss.Color("#1e1e1e")
	text := lipgloss.Color("#E9E7FF")

	tabBorder := lipgloss.RoundedBorder()
	return Theme{
		Doc: lipgloss.NewStyle().
			Background(bg).
			Foreground(text),
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(text).
			Background(lipgloss.Color("#4B38A7")).
			Padding(0, 1),
		Bar: lipgloss.NewStyle().
			Foreground(text).
			Background(lipgloss.Color("#2A2A2A")),
		Tab: lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderForeground(subtle).
			Foreground(lipgloss.Color("#BBB7D6")).
			Padding(0, 1),
		ActiveTab: lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderForeground(accent).
			Foreground(text).
			Background(lipgloss.Color("#2A2440")).
			Padding(0, 1),
		PaneTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#101014")).
			Background(cyan).
			Bold(true).
			Padding(0, 1),
		Panel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4A4F66")).
			Background(panel).
			Padding(1, 2),
		Sidebar: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accent).
			Background(panel).
			Padding(1, 1),
		Detail: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lipgloss.Color("#3C4055")).
			Padding(1, 0, 0, 0),
		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DAD7FF")).
			Background(lipgloss.Color("#2A2A2A")),
		Nugget: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#101014")).
			Background(cyan).
			Bold(true).
			Padding(0, 1),
		Accent: lipgloss.NewStyle().Foreground(cyan).Bold(true),
		Muted:  lipgloss.NewStyle().Foreground(subtle),
		Danger: lipgloss.NewStyle().Foreground(pink).Bold(true),
		Modal: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(pink).
			Background(bg).
			Foreground(text).
			Padding(1, 2),
		Button: lipgloss.NewStyle().
			Foreground(text).
			Background(lipgloss.Color("#50566A")).
			Padding(0, 2).
			MarginRight(1),
		ActiveBtn: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#11131A")).
			Background(cyan).
			Bold(true).
			Padding(0, 2).
			MarginRight(1),
	}
}
