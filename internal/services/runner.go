package services

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/deathrashed/mac-cli/internal/commands"
)

type Result struct {
	Command string
	Output  string
	Err     string
}

type StreamEvent struct {
	Line string
	Done bool
	Err  error
}

type Runner struct {
	BaseDir string
}

func NewRunner(baseDir string) Runner {
	return Runner{BaseDir: baseDir}
}

func (r Runner) Preview(cmd commands.Command, args []string) string {
	return cmd.Preview(args)
}

func (r Runner) Run(ctx context.Context, cmd commands.Command, args []string) (Result, error) {
	if cmd.Native {
		return runNative(cmd, args)
	}
	return r.runLegacy(ctx, cmd, args)
}

func (r Runner) Stream(ctx context.Context, cmd commands.Command, args []string) <-chan StreamEvent {
	events := make(chan StreamEvent, 64)
	go func() {
		defer close(events)
		if cmd.Native {
			result, err := runNative(cmd, args)
			if result.Output != "" {
				for _, line := range strings.Split(strings.TrimRight(result.Output, "\n"), "\n") {
					events <- StreamEvent{Line: line}
				}
			}
			if result.Err != "" {
				for _, line := range strings.Split(strings.TrimRight(result.Err, "\n"), "\n") {
					events <- StreamEvent{Line: line}
				}
			}
			events <- StreamEvent{Done: true, Err: err}
			return
		}
		r.streamLegacy(ctx, cmd, args, events)
	}()
	return events
}

func (r Runner) runLegacy(ctx context.Context, command commands.Command, args []string) (Result, error) {
	legacy := filepath.Join(r.BaseDir, "mac")
	argv := append([]string{legacy, legacyName(command.Name)}, args...)
	c := exec.CommandContext(ctx, "bash", argv...)
	c.Dir = r.BaseDir
	c.Env = runnerEnv(r.BaseDir)
	c.Stdin = strings.NewReader(confirmInput(command))
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	return Result{Command: command.Preview(args), Output: stdout.String(), Err: stderr.String()}, err
}

func (r Runner) streamLegacy(ctx context.Context, command commands.Command, args []string, events chan<- StreamEvent) {
	legacy := filepath.Join(r.BaseDir, "mac")
	argv := append([]string{legacy, legacyName(command.Name)}, args...)
	c := exec.CommandContext(ctx, "bash", argv...)
	c.Dir = r.BaseDir
	c.Env = runnerEnv(r.BaseDir)
	c.Stdin = strings.NewReader(confirmInput(command))

	reader, writer := io.Pipe()
	c.Stdout = writer
	c.Stderr = writer

	scanDone := make(chan struct{})
	go func() {
		defer close(scanDone)
		scanner := bufio.NewScanner(reader)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)
		scanner.Split(scanOutputLine)
		for scanner.Scan() {
			events <- StreamEvent{Line: scanner.Text()}
		}
		if err := scanner.Err(); err != nil {
			events <- StreamEvent{Line: "output stream error: " + err.Error()}
		}
	}()

	err := c.Start()
	if err != nil {
		_ = writer.Close()
		<-scanDone
		events <- StreamEvent{Done: true, Err: err}
		return
	}
	err = c.Wait()
	_ = writer.Close()
	<-scanDone
	events <- StreamEvent{Done: true, Err: err}
}

func scanOutputLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i, b := range data {
		if b == '\n' || b == '\r' {
			return i + 1, data[:i], nil
		}
	}
	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func runNative(cmd commands.Command, args []string) (Result, error) {
	out := ""
	switch cmd.Name {
	case "list", "help", "usage":
		out = nativeHelp(commands.DefaultCatalog(), "")
	case "categories":
		out = nativeCategories(commands.DefaultCatalog())
	case "general", "search", "network", "dns", "ssh", "webdev", "performance", "terminal", "git", "web_utilities", "homebrew", "image", "magento":
		out = nativeHelp(commands.DefaultCatalog(), cmd.Name)
	case "ssh:list":
		out = nativeHelp(commands.DefaultCatalog(), "ssh")
	case "display":
		if err := exec.Command("pmset", "displaysleepnow").Run(); err != nil {
			return Result{}, err
		}
		out = "Display sleep requested.\n"
	case "convert:c-to-f":
		v, err := oneFloat(args, "celsius")
		if err != nil {
			return Result{}, err
		}
		out = fmt.Sprintf("%.2f°C is equal to %.2f°F\n", v, (v*9/5)+32)
	case "convert:f-to-c":
		v, err := oneFloat(args, "fahrenheit")
		if err != nil {
			return Result{}, err
		}
		out = fmt.Sprintf("%.2f°F is equal to %.2f°C\n", v, (v-32)*5/9)
	case "convert:to-lowercase":
		if len(args) == 0 {
			return Result{}, errors.New("text argument is required")
		}
		out = strings.ToLower(strings.Join(args, " ")) + "\n"
	case "convert:to-uppercase":
		if len(args) == 0 {
			return Result{}, errors.New("text argument is required")
		}
		out = strings.ToUpper(strings.Join(args, " ")) + "\n"
	case "text:urlencode":
		if len(args) == 0 {
			return Result{}, errors.New("text argument is required")
		}
		out = url.QueryEscape(strings.Join(args, " ")) + "\n"
	case "text:urldecode":
		if len(args) == 0 {
			return Result{}, errors.New("text argument is required")
		}
		decoded, err := url.QueryUnescape(strings.Join(args, " "))
		if err != nil {
			return Result{}, err
		}
		out = decoded + "\n"
	default:
		return Result{}, fmt.Errorf("no native handler for %s", cmd.Name)
	}
	return Result{Command: cmd.Preview(args), Output: out}, nil
}

func nativeCategories(catalog commands.Catalog) string {
	var b strings.Builder
	b.WriteString("MAC-CLI categories\n\n")
	for _, category := range catalog.Categories() {
		b.WriteString("  " + category + "\n")
	}
	return b.String()
}

func nativeHelp(catalog commands.Catalog, category string) string {
	var b strings.Builder
	if category == "" {
		b.WriteString("MAC-CLI commands\n\n")
		for _, cat := range catalog.Categories() {
			cmds := catalog.ByCategory(cat)
			if len(cmds) == 0 {
				continue
			}
			b.WriteString(cat + "\n")
			for _, cmd := range cmds {
				b.WriteString(fmt.Sprintf("  %-24s %s\n", cmd.Name, cmd.Description))
			}
			b.WriteString("\n")
		}
		return b.String()
	}
	title := strings.ReplaceAll(category, "_", " ")
	b.WriteString("MAC-CLI " + title + " commands\n\n")
	for _, cmd := range catalog.ByCategory(titleCaseCategory(category)) {
		b.WriteString(fmt.Sprintf("  %-24s %s\n", cmd.Name, cmd.Description))
	}
	return b.String()
}

func titleCaseCategory(category string) string {
	switch category {
	case "dns":
		return "DNS"
	case "ssh":
		return "SSH"
	case "web_utilities":
		return "Web_Utilities"
	case "webdev":
		return "Webdev"
	default:
		if category == "" {
			return category
		}
		return strings.ToUpper(category[:1]) + category[1:]
	}
}

func oneFloat(args []string, name string) (float64, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("%s argument is required", name)
	}
	return strconv.ParseFloat(args[0], 64)
}

func shellArgs(args []string) string {
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		quoted = append(quoted, strconv.Quote(arg))
	}
	return strings.Join(quoted, " ")
}

func runnerEnv(baseDir string) []string {
	return append(os.Environ(),
		"MAC_CLI_BASE_DIR="+baseDir,
		"echocommand=false",
		"CLICOLOR_FORCE=1",
		"FORCE_COLOR=1",
		"TERM=xterm-256color",
	)
}

func confirmInput(command commands.Command) string {
	if command.Dangerous {
		return "y\n"
	}
	return "n\n"
}

func legacyName(name string) string {
	aliases := map[string]string{
		"ports":             "network:ports",
		"speedtest":         "network:speedtest",
		"ip:local":          "network:ip:local",
		"ip:public":         "network:public-ip",
		"git:branches:date": "git:branches",
		"git:settings":      "git:config",
	}
	if alias, ok := aliases[name]; ok {
		return alias
	}
	return name
}
