# Mac CLI Charm TUI Migration

## Strategy

1. Keep the shell plugins working through the Bash compatibility bridge.
2. Move low-risk commands to native Go first: pure parsing, conversions, URL encoding, and read-only status commands.
3. Add Gum prompts for commands that mutate system state or user data.
4. Promote each category into `internal/commands/<category>` once its shell behavior has tests.
5. Replace shell command previews with structured execution plans as commands become native.

## First Native Modules

- `convert:c-to-f`
- `convert:f-to-c`
- `convert:to-lowercase`
- `convert:to-uppercase`
- `text:urlencode`
- `text:urldecode`

## Roadmap

- Phase 1: TUI shell bridge, catalog, preview modal, palette, settings, tests.
- Phase 2: Native search/file/text modules with dry-run support.
- Phase 3: Native Git and Homebrew wrappers with richer progress output.
- Phase 4: Native macOS services with privilege-aware prompts.
- Phase 5: Optional Wish SSH interface for remote Mac CLI sessions.
