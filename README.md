# Mac CLI

Mac CLI is a macOS developer toolkit being modernised from shell plugins into a Charmbracelet-powered terminal UI. The existing shell commands remain available through a compatibility bridge while command groups are migrated to native Go modules.

## New TUI

Run the Bubble Tea interface:

```bash
mac
```

Build the cached TUI binary used by the `mac` launcher:

```bash
./tools/build-tui
```

If `bin/mac-tui` does not exist yet, bare `mac` falls back to `go run ./cmd/mac --tui`.

The interface includes:

- multi-column category and command navigation
- searchable Bubbles lists
- command palette with `ctrl+p`
- help modal with `?`
- settings screen with `s`
- command preview modal before execution
- spinner and progress feedback while commands run
- status bar with command metadata

## CLI Compatibility

Existing commands still work from the Go entrypoint:

```bash
go run ./cmd/mac convert:c-to-f 25
go run ./cmd/mac system:ports
go run ./cmd/mac --yes system:cleanup-dsstore
```

They also still work through the installed shell command:

```bash
mac categories
mac dns:list
mac system:ports
```

Dangerous commands ask for confirmation through Gum when `gum` is installed. Set `MAC_CLI_CONFIRM_DANGER=0` or pass `--yes` for non-interactive runs.

The legacy shell script is still present as `./mac`. The Go runner invokes it with:

```bash
MAC_CLI_BASE_DIR=/path/to/checkout bash ./mac <command> [args...]
```

## Architecture

```text
cmd/mac/              Go entrypoint and CLI flags
internal/commands/    command catalog and native module metadata
internal/services/    execution, Gum fallback, legacy shell bridge
internal/config/      environment-backed configuration
ui/mainmenu/          Bubble Tea MVU application
ui/submenu/           future category-specific views
ui/palette/           command-palette expansion point
ui/modals/            help, preview, result, settings modal rendering
ui/styles/            Lip Gloss theme and shared styles
ui/components/        Bubbles list/table item adapters
pkg/utils/            small reusable helpers
docs/MIGRATION.md     migration plan and roadmap
plugins/              legacy shell command implementations
```

## Keyboard

| Key | Action |
| --- | --- |
| `↑/↓`, `j/k` | Move selection |
| `←/→`, `h/l` | Switch category |
| `/` | Filter the focused list |
| `enter` | Preview selected command |
| `ctrl+p` | Open command palette |
| `?` | Open help |
| `s` | Open settings |
| `esc` | Close modal or cancel |
| `q` | Quit |

## Native Rewrite Examples

These commands have native Go implementations already:

- `convert:c-to-f`
- `convert:f-to-c`
- `convert:to-lowercase`
- `convert:to-uppercase`
- `text:urlencode`
- `text:urldecode`

All other cataloged commands are routed through the Bash compatibility bridge so existing functionality is preserved while migration continues.

## Development

Install dependencies:

```bash
go mod tidy
```

Run tests:

```bash
go test ./...
```

Validate legacy plugin syntax:

```bash
bash -n mac
for plugin in plugins/*; do bash -n "$plugin"; done
```

## Migration Roadmap

See [docs/MIGRATION.md](docs/MIGRATION.md) for the staged rewrite plan. The short version: keep shell behavior working, migrate low-risk parsing commands first, add tests around every migrated command, then move system-mutating commands into native Go with explicit Gum/Bubble Tea confirmations.
