package mainmenu

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deathrashed/mac-cli/internal/commands"
	"github.com/deathrashed/mac-cli/internal/services"
	"github.com/deathrashed/mac-cli/pkg/utils"
	"github.com/deathrashed/mac-cli/ui/components"
	"github.com/deathrashed/mac-cli/ui/modals"
	"github.com/deathrashed/mac-cli/ui/styles"
)

type Mode int

const (
	modeBrowse Mode = iota
	modePalette
	modeHelp
	modeSettings
	modeArgs
	modePreview
	modeConfirm
	modeRunning
	modeResult
)

type streamMsg services.StreamEvent

type focusPane int

const (
	focusCategories focusPane = iota
	focusCommands
)

type Model struct {
	catalog    commands.Catalog
	runner     services.Runner
	theme      styles.Theme
	categories list.Model
	commands   list.Model
	palette    list.Model
	argInput   textinput.Model
	detail     table.Model
	outputView viewport.Model
	spinner    spinner.Model
	progress   progress.Model
	mode       Mode
	width      int
	height     int
	category   string
	focus      focusPane
	selected   commands.Command
	args       []string
	status     string
	errText    string
	stream     <-chan services.StreamEvent
	output     []string
	outputFile string
}

func New(catalog commands.Catalog, runner services.Runner) Model {
	theme := styles.New()
	categoryItems := make([]list.Item, 0)
	for _, cat := range catalog.Categories() {
		categoryItems = append(categoryItems, components.CategoryItem(cat))
	}
	allItems := commandItems(catalog.Commands)
	categoryDelegate := list.NewDefaultDelegate()
	categoryDelegate.ShowDescription = false
	categoryDelegate.SetSpacing(0)
	categoryDelegate.Styles.NormalTitle = categoryDelegate.Styles.NormalTitle.Padding(0, 0, 0, 1)
	categoryDelegate.Styles.SelectedTitle = categoryDelegate.Styles.SelectedTitle.Padding(0, 0, 0, 1)
	commandDelegate := list.NewDefaultDelegate()
	commandDelegate.ShowDescription = false
	commandDelegate.SetSpacing(0)
	commandDelegate.Styles.NormalTitle = commandDelegate.Styles.NormalTitle.Padding(0, 0, 0, 1)
	commandDelegate.Styles.SelectedTitle = commandDelegate.Styles.SelectedTitle.Padding(0, 0, 0, 1)
	paletteDelegate := list.NewDefaultDelegate()
	paletteDelegate.SetSpacing(1)
	paletteDelegate.Styles.NormalTitle = paletteDelegate.Styles.NormalTitle.Padding(0, 0, 0, 1)
	paletteDelegate.Styles.NormalDesc = paletteDelegate.Styles.NormalDesc.Padding(0, 0, 0, 1)
	paletteDelegate.Styles.SelectedTitle = paletteDelegate.Styles.SelectedTitle.Padding(0, 0, 0, 1)
	paletteDelegate.Styles.SelectedDesc = paletteDelegate.Styles.SelectedDesc.Padding(0, 0, 0, 1)

	cats := list.New(categoryItems, categoryDelegate, 28, 20)
	cats.Title = ""
	cats.SetShowTitle(false)
	cats.SetShowStatusBar(false)
	cats.SetShowPagination(false)
	cats.SetShowHelp(false)
	cats.SetFilteringEnabled(true)
	cmds := list.New(nil, commandDelegate, 52, 20)
	cmds.Title = ""
	cmds.SetShowTitle(false)
	cmds.SetShowStatusBar(false)
	cmds.SetShowPagination(false)
	cmds.SetShowHelp(false)
	cmds.SetFilteringEnabled(true)
	palette := list.New(allItems, paletteDelegate, 80, 20)
	palette.Title = "Command Palette"
	palette.SetFilteringEnabled(true)
	input := textinput.New()
	input.Placeholder = "Type the requested values here"
	input.Prompt = "Input: "
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	prog := progress.New(progress.WithDefaultGradient())
	detail := table.New(table.WithColumns([]table.Column{
		{Title: "Field", Width: 14},
		{Title: "Value", Width: 54},
	}))
	outputView := viewport.New(80, 24)
	m := Model{
		catalog:    catalog,
		runner:     runner,
		theme:      theme,
		categories: cats,
		commands:   cmds,
		palette:    palette,
		argInput:   input,
		detail:     detail,
		outputView: outputView,
		spinner:    spin,
		progress:   prog,
		mode:       modeBrowse,
		focus:      focusCategories,
		status:     "Choose a category",
	}
	if len(categoryItems) > 0 {
		m.category = categoryItems[0].(components.CategoryItem).Title()
		m.refreshCommands()
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resize()
	case streamMsg:
		if msg.Line != "" {
			m.output = append(m.output, cleanOutputLine(msg.Line))
			m.trimOutput()
			m.syncOutputView(true)
		}
		if msg.Done {
			m.stream = nil
			m.mode = modeResult
			if msg.Err != nil {
				m.errText = msg.Err.Error()
				m.status = "Command failed: " + msg.Err.Error()
			} else {
				m.errText = ""
				m.status = "Command finished"
			}
			m.syncOutputView(false)
			return m, nil
		}
		return m, m.waitForStream()
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	if m.mode == modeRunning {
		if key, ok := msg.(tea.KeyMsg); ok && (key.String() == "esc" || key.String() == "q") {
			return m, nil
		}
		var vpCmd tea.Cmd
		m.outputView, vpCmd = m.outputView.Update(msg)
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, tea.Batch(cmd, vpCmd)
	}

	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m.updateLists(msg)
	}

	switch m.mode {
	case modeHelp, modeSettings:
		if key.String() == "esc" || key.String() == "q" {
			m.mode = modeBrowse
			m.status = "Ready"
		}
		return m, nil
	case modeResult:
		switch key.String() {
		case "c":
			m.copySmartOutput()
			return m, nil
		case "a":
			m.copyAllOutput()
			return m, nil
		case "s":
			m.saveOutput()
			return m, nil
		}
		if key.String() == "esc" || key.String() == "q" {
			m.mode = modeBrowse
			m.status = "Ready"
			m.output = nil
			m.outputView.SetContent("")
			m.outputFile = ""
			return m, nil
		}
		var cmd tea.Cmd
		m.outputView, cmd = m.outputView.Update(msg)
		return m, cmd
	case modeArgs:
		switch key.String() {
		case "esc":
			m.mode = modeBrowse
			return m, nil
		case "enter":
			m.args = strings.Fields(m.argInput.Value())
			m.mode = modePreview
			return m, nil
		}
		var cmd tea.Cmd
		m.argInput, cmd = m.argInput.Update(msg)
		return m, cmd
	case modePreview:
		switch key.String() {
		case "esc":
			m.mode = modeBrowse
			return m, nil
		case "enter":
			if m.selected.Dangerous {
				m.mode = modeConfirm
				m.status = "Confirm " + m.selected.Name
				return m, nil
			}
			return m.startCommand()
		}
		return m, nil
	case modeConfirm:
		switch key.String() {
		case "esc", "n":
			m.mode = modeBrowse
			m.status = "Cancelled " + m.selected.Name
			return m, nil
		case "enter", "y":
			return m.startCommand()
		}
		return m, nil
	case modePalette:
		switch key.String() {
		case "esc", "ctrl+p":
			m.mode = modeBrowse
			return m, nil
		case "enter":
			if item, ok := m.palette.SelectedItem().(components.CommandItem); ok {
				m.chooseCommand(item.Command)
			}
			return m, nil
		}
		var cmd tea.Cmd
		m.palette, cmd = m.palette.Update(msg)
		return m, cmd
	}

	switch key.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "?":
		m.mode = modeHelp
	case "s":
		m.mode = modeSettings
	case "ctrl+p":
		m.mode = modePalette
	case "tab":
		m.toggleFocus()
	case "left", "h":
		m.focus = focusCategories
		m.status = "Choose a category"
	case "right", "l":
		m.focus = focusCommands
		m.status = "Choose a command in " + m.category
	case "enter":
		if m.focus == focusCategories {
			m.focus = focusCommands
			m.status = "Choose a command in " + m.category
		} else if item, ok := m.commands.SelectedItem().(components.CommandItem); ok {
			m.chooseCommand(item.Command)
		}
	default:
		return m.updateLists(msg)
	}
	return m, nil
}

func (m Model) updateLists(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.focus == focusCategories {
		m.categories, cmd = m.categories.Update(msg)
		m.syncCategory()
		return m, cmd
	}
	if m.focus == focusCommands {
		m.commands, cmd = m.commands.Update(msg)
		if item, ok := m.commands.SelectedItem().(components.CommandItem); ok {
			m.updateDetail(item.Command)
		}
		return m, cmd
	}
	return m, nil
}

func (m *Model) chooseCommand(cmd commands.Command) {
	m.selected = cmd
	m.args = nil
	m.argInput.SetValue("")
	m.argInput.Placeholder = placeholderFor(cmd)
	m.updateDetail(cmd)
	if len(cmd.Args) > 0 {
		m.mode = modeArgs
		m.argInput.Focus()
		m.status = "Waiting for input: " + cmd.Name
		return
	}
	m.mode = modePreview
	m.status = "Preview " + cmd.Name
}

func (m *Model) syncCategory() {
	if item, ok := m.categories.SelectedItem().(components.CategoryItem); ok && item.Title() != m.category {
		m.category = item.Title()
		m.refreshCommands()
	}
}

func (m *Model) refreshCommands() {
	m.commands.SetItems(commandItems(m.catalog.ByCategory(m.category)))
	m.commands.Select(0)
	if item, ok := m.commands.SelectedItem().(components.CommandItem); ok {
		m.updateDetail(item.Command)
	}
}

func (m *Model) updateDetail(cmd commands.Command) {
	args := "-"
	if len(cmd.Args) > 0 {
		args = strings.Join(cmd.Args, ", ")
	}
	kind := "legacy shell bridge"
	if cmd.Native {
		kind = "native Go"
	}
	m.detail.SetRows([]table.Row{
		{"Command", cmd.Name},
		{"Category", cmd.Category},
		{"Arguments", args},
		{"Runner", kind},
		{"Preview", cmd.Preview(nil)},
	})
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}
	base := m.baseView()
	switch m.mode {
	case modePalette:
		return m.overlay(base, m.palette.View())
	case modeHelp:
		return m.overlay(base, modals.Help(m.theme))
	case modeSettings:
		return m.overlay(base, modals.Settings(m.runner.BaseDir, services.GumAvailable(), m.theme))
	case modeArgs:
		return m.overlay(base, m.argsView())
	case modePreview:
		return m.overlay(base, modals.Preview(m.selected.Name, m.runner.Preview(m.selected, m.args), m.selected.Dangerous, m.theme))
	case modeConfirm:
		return m.overlay(base, m.confirmView())
	case modeRunning:
		return m.outputScreen(true)
	case modeResult:
		return m.outputScreen(false)
	default:
		return base
	}
}

func (m Model) baseView() string {
	header := m.header()
	status := m.statusBar()
	sidebarWidth, contentWidth, height := m.layout()
	bodyHeight := max(8, m.height-lipgloss.Height(header)-lipgloss.Height(status))
	paneHeight := max(6, min(height, bodyHeight-2))
	sidebarStyle := m.theme.Sidebar.BorderForeground(lipgloss.Color("#6B5BFF"))
	panelStyle := m.theme.Panel.BorderForeground(lipgloss.Color("#4A4F66"))
	if m.focus == focusCategories {
		sidebarStyle = sidebarStyle.BorderForeground(lipgloss.Color("#56D4DD"))
	} else {
		panelStyle = panelStyle.BorderForeground(lipgloss.Color("#56D4DD"))
	}
	leftHeader := m.paneHeader("Categories", "Enter to open")
	left := sidebarStyle.Width(sidebarWidth).Height(paneHeight).Render(lipgloss.JoinVertical(lipgloss.Left, leftHeader, "", m.categories.View()))
	rightHeader := m.paneHeader(m.category+" Commands", fmt.Sprintf("%d available", len(m.catalog.ByCategory(m.category))))
	rightBody := lipgloss.JoinVertical(lipgloss.Left, rightHeader, "", m.commands.View(), "", m.detailView())
	right := panelStyle.Width(contentWidth).Height(paneHeight).Render(rightBody)
	body := lipgloss.Place(
		m.width,
		bodyHeight,
		lipgloss.Left,
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right),
		lipgloss.WithWhitespaceBackground(lipgloss.Color("#1e1e1e")),
	)
	page := lipgloss.JoinVertical(lipgloss.Left, header, body, status)
	return m.theme.Doc.Width(m.width).Height(m.height).Render(lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Left,
		lipgloss.Top,
		page,
		lipgloss.WithWhitespaceBackground(lipgloss.Color("#1e1e1e")),
	))
}

func (m Model) header() string {
	sidebarWidth, contentWidth, _ := m.layout()
	totalWidth := min(m.width, sidebarWidth+contentWidth+5)
	left := " MAC-CLI "
	right := m.category
	if m.mode == modeRunning || m.mode == modeResult {
		right = m.selected.Name
	}
	gap := strings.Repeat(" ", max(1, totalWidth-lipgloss.Width(left)-lipgloss.Width(right)-2))
	return m.theme.Header.Width(totalWidth).Render(left + gap + right)
}

func (m Model) statusBar() string {
	left := m.theme.Nugget.Render(m.status)
	focus := "categories"
	if m.focus == focusCommands {
		focus = "commands"
	}
	right := m.theme.Status.Render(fmt.Sprintf(" %s · %d commands · tab switch · ctrl+p palette · ? help · q quit ", focus, len(m.catalog.Commands)))
	gap := strings.Repeat(" ", max(0, m.width-lipgloss.Width(left)-lipgloss.Width(right)))
	return m.theme.Status.Width(m.width).Render(left + gap + right)
}

func (m Model) overlay(_ string, content string) string {
	return modals.Center(m.width, m.height, content, m.theme)
}

func (m Model) argsView() string {
	hint := "Run " + m.selected.Name
	if len(m.selected.Args) > 0 {
		hint += "\nNeeded: " + strings.Join(m.selected.Args, ", ")
	}
	example := argumentHelp(m.selected)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.theme.Accent.Render(hint),
		"",
		example,
		"",
		m.argInput.View(),
		"",
		m.theme.Muted.Render("Enter previews the exact command. Esc cancels."),
	)
}

func (m Model) outputScreen(running bool) string {
	header := m.header()
	_, contentWidth, panelHeight := m.layout()
	width := min(m.width-4, contentWidth+8)
	height := max(10, m.height-5)
	title := "Output"
	if running {
		title = m.spinner.View() + " Running " + m.selected.Name
	}
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.paneHeader(title, m.selected.Preview(m.args)),
		"",
		m.outputView.View(),
	)
	panel := m.theme.Panel.Width(width).Height(min(panelHeight+5, height)).BorderForeground(lipgloss.Color("#FF6FAE")).Render(content)
	footerLeft := "Command finished"
	if running {
		footerLeft = "Streaming output"
	}
	left := m.theme.Nugget.Render(footerLeft)
	actions := " ↑/↓ scroll · pgup/pgdn page · c copy · s save · esc back "
	if running {
		actions = " ↑/↓ scroll · pgup/pgdn page · live output "
	}
	right := m.theme.Status.Render(actions)
	gap := strings.Repeat(" ", max(0, m.width-lipgloss.Width(left)-lipgloss.Width(right)))
	footer := m.theme.Status.Width(m.width).Render(left + gap + right)
	page := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		lipgloss.Place(m.width, m.height-2, lipgloss.Center, lipgloss.Center, panel, lipgloss.WithWhitespaceBackground(lipgloss.Color("#1e1e1e"))),
		footer,
	)
	return m.theme.Doc.Width(m.width).Height(m.height).Render(page)
}

func (m Model) confirmView() string {
	lines := []string{
		m.theme.Danger.Render("Confirmation required"),
		"",
		"This command can change system state or remove data.",
		"",
		m.theme.Accent.Render(m.selected.Preview(m.args)),
		"",
		m.theme.Muted.Render("Press y or Enter to run. Press n or Esc to cancel."),
	}
	return strings.Join(lines, "\n")
}

func (m Model) startCommand() (tea.Model, tea.Cmd) {
	m.mode = modeRunning
	m.status = "Running " + m.selected.Name
	m.outputFile = ""
	m.output = []string{m.theme.Muted.Render("$") + " " + m.theme.Accent.Render(m.selected.Preview(m.args)), ""}
	m.syncOutputView(true)
	m.stream = m.runner.Stream(context.Background(), m.selected, m.args)
	return m, tea.Batch(m.spinner.Tick, m.waitForStream())
}

func (m *Model) resize() {
	sidebarWidth, contentWidth, panelHeight := m.layout()
	listHeight := utils.Clamp(panelHeight-12, 6, max(6, panelHeight-5))
	detailWidth := max(24, contentWidth-8)
	m.commands.SetSize(max(24, contentWidth-6), listHeight)
	m.categories.SetSize(max(16, sidebarWidth-2), max(6, panelHeight-5))
	m.palette.SetSize(utils.Clamp(m.width-12, 40, 96), utils.Clamp(m.height-8, 10, 32))
	m.argInput.Width = utils.Clamp(m.width-18, 24, 80)
	m.progress.Width = utils.Clamp(m.width-18, 24, 80)
	m.outputView.Width = max(20, min(m.width-12, contentWidth+2))
	m.outputView.Height = max(6, m.height-13)
	m.detail.SetColumns([]table.Column{
		{Title: "Field", Width: 12},
		{Title: "Value", Width: detailWidth - 14},
	})
	m.detail.SetWidth(detailWidth)
	m.detail.SetHeight(utils.Clamp(panelHeight-listHeight-5, 4, 8))
}

func (m Model) waitForStream() tea.Cmd {
	ch := m.stream
	return func() tea.Msg {
		if ch == nil {
			return streamMsg{Done: true}
		}
		event, ok := <-ch
		if !ok {
			return streamMsg{Done: true}
		}
		return streamMsg(event)
	}
}

func (m *Model) trimOutput() {
	const maxLines = 2000
	if len(m.output) > maxLines {
		m.output = append([]string(nil), m.output[len(m.output)-maxLines:]...)
	}
}

func (m *Model) syncOutputView(follow bool) {
	m.outputView.SetContent(strings.Join(m.output, "\n"))
	if follow || m.outputView.AtBottom() {
		m.outputView.GotoBottom()
	}
}

func cleanOutputLine(line string) string {
	replacer := strings.NewReplacer(
		"\\033", "\033",
		"\\n", "",
		"\\t", "  ",
	)
	return replacer.Replace(line)
}

func (m *Model) copyOutput() {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(strings.Join(stripANSI(m.output), "\n"))
	if err := cmd.Run(); err != nil {
		m.status = "Copy failed: " + err.Error()
		return
	}
	m.status = "Copied output"
}

func (m *Model) saveOutput() {
	name := sanitizeFilename(m.selected.Name)
	path := filepath.Join(os.TempDir(), fmt.Sprintf("mac-cli-%s-%s.log", name, time.Now().Format("20060102-150405")))
	if err := os.WriteFile(path, []byte(strings.Join(stripANSI(m.output), "\n")+"\n"), 0644); err != nil {
		m.status = "Save failed: " + err.Error()
		return
	}
	m.outputFile = path
	m.status = "Saved output to " + path
}

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

func stripANSI(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		out = append(out, ansiPattern.ReplaceAllString(line, ""))
	}
	return out
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(":", "-", "/", "-", " ", "-")
	return replacer.Replace(name)
}

func (m Model) visibleOutput(height int) []string {
	if len(m.output) == 0 {
		return []string{m.theme.Muted.Render("No output yet.")}
	}
	start := max(0, len(m.output)-height)
	lines := m.output[start:]
	out := make([]string, 0, len(lines)+1)
	for _, line := range lines {
		out = append(out, lipgloss.NewStyle().MaxWidth(max(20, m.width-12)).Render(line))
	}
	if m.errText != "" {
		out = append(out, m.theme.Danger.Render(m.errText))
	}
	return out
}

func (m Model) layout() (sidebarWidth int, contentWidth int, panelHeight int) {
	width := max(60, m.width)
	height := max(20, m.height)
	sidebarWidth = utils.Clamp(width/3, 26, 38)
	if width < 90 {
		sidebarWidth = utils.Clamp(width/3, 22, 28)
	}
	contentWidth = max(30, width-sidebarWidth-7)
	panelHeight = max(8, height-8)
	return sidebarWidth, contentWidth, panelHeight
}

func (m *Model) toggleFocus() {
	if m.focus == focusCategories {
		m.focus = focusCommands
		m.status = "Choose a command in " + m.category
		return
	}
	m.focus = focusCategories
	m.status = "Choose a category"
}

func (m Model) paneHeader(title, meta string) string {
	label := m.theme.PaneTitle.Render(title)
	metaText := m.theme.Muted.Render(" " + meta)
	return lipgloss.JoinHorizontal(lipgloss.Center, label, metaText)
}

func (m Model) detailView() string {
	cmd, ok := m.commands.SelectedItem().(components.CommandItem)
	if !ok {
		return m.theme.Detail.Render(m.theme.Muted.Render("Select a command to see details."))
	}
	args := "No arguments required"
	if len(cmd.Command.Args) > 0 {
		args = strings.Join(cmd.Command.Args, ", ")
	}
	runner := "Shell compatibility bridge"
	if cmd.Command.Native {
		runner = "Native Go module"
	}
	danger := "No"
	if cmd.Command.Dangerous {
		danger = m.theme.Danger.Render("Yes, confirmation required")
	}
	rows := []string{
		m.theme.Accent.Render("Details"),
		"",
		m.detailRow("Command", cmd.Command.Name),
		m.detailRow("Purpose", cmd.Command.Description),
		m.detailRow("Inputs", args),
		m.detailRow("Runner", runner),
		m.detailRow("Dangerous", danger),
	}
	return m.theme.Detail.Width(max(24, m.detail.Width())).Render(strings.Join(rows, "\n"))
}

func (m Model) detailRow(label, value string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.theme.Muted.Width(12).Render(label),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#E9E7FF")).Render(value),
	)
}

func placeholderFor(cmd commands.Command) string {
	if len(cmd.Args) == 0 {
		return "No input needed"
	}
	return "Example: " + strings.Join(cmd.Args, " ")
}

func argumentHelp(cmd commands.Command) string {
	if len(cmd.Args) == 0 {
		return "This command does not need extra input."
	}
	if len(cmd.Args) == 1 {
		return "Type the " + cmd.Args[0] + " value for this command."
	}
	return "Type the values in this order: " + strings.Join(cmd.Args, " ")
}

func commandItems(cmds []commands.Command) []list.Item {
	items := make([]list.Item, 0, len(cmds))
	for _, cmd := range cmds {
		items = append(items, components.CommandItem{Command: cmd})
	}
	return items
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
