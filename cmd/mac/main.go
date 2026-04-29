package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/deathrashed/mac-cli/internal/commands"
	"github.com/deathrashed/mac-cli/internal/config"
	"github.com/deathrashed/mac-cli/internal/services"
	"github.com/deathrashed/mac-cli/ui/mainmenu"
)

func main() {
	var (
		tuiFlag = flag.Bool("tui", false, "launch the full-screen Bubble Tea interface")
		yesFlag = flag.Bool("yes", false, "skip Gum confirmation for dangerous commands")
	)
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		fatal(err)
	}
	cfg := config.Load(cwd)
	catalog := commands.DefaultCatalog()
	runner := services.NewRunner(cfg.BaseDir)

	if *tuiFlag || flag.NArg() == 0 {
		if _, err := tea.NewProgram(mainmenu.New(catalog, runner), tea.WithAltScreen()).Run(); err != nil {
			fatal(err)
		}
		return
	}

	name := flag.Arg(0)
	cmd, ok := catalog.Find(name)
	if !ok {
		fatal(fmt.Errorf("unknown command %q; run mac --tui or mac list", name))
	}
	args := flag.Args()[1:]
	if cmd.Dangerous && cfg.ConfirmDanger && !*yesFlag {
		confirmed, err := services.GumConfirm(context.Background(), "Run "+cmd.Preview(args)+"?")
		if err != nil {
			fatal(err)
		}
		if !confirmed {
			fmt.Fprintln(os.Stderr, "Cancelled.")
			os.Exit(1)
		}
	}
	result, err := runner.Run(context.Background(), cmd, args)
	if result.Output != "" {
		fmt.Print(result.Output)
	}
	if result.Err != "" {
		fmt.Fprint(os.Stderr, result.Err)
	}
	if err != nil {
		fatal(fmt.Errorf("%s failed: %w", strings.Join(append([]string{name}, args...), " "), err))
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "mac:", err)
	os.Exit(1)
}
