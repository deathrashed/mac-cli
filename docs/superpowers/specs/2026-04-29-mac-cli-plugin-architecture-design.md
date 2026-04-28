# Mac CLI Plugin Architecture Phase 1 Design

## Status

Approved for implementation planning.

## Goals

- Replace hardcoded command, help, and completion registries with plugin-owned metadata.
- Keep normal `mac` command execution fast by avoiding JSON parsing at runtime.
- Preserve existing command names and plugin behavior during Phase 1.
- Make adding a new command require one implementation change and one metadata change.
- Add generated Zsh completion alongside refreshed Bash and Fish completion.
- Prepare the codebase for later plugin isolation and command handler refactors.

## Non-goals

- Rewriting every plugin command into namespaced functions in Phase 1.
- Changing the public command syntax from `mac git:open` to `mac git open` in Phase 1.
- Removing every command that currently uses `sudo`; Phase 1 should remove installer/runtime registry assumptions first and document risky command behavior for follow-up refactors.
- Requiring `jq`, Python, or Node during normal CLI execution or shell completion.

## Current State

The current `mac` entrypoint contains a hardcoded `COMMANDS` array and manually sources each plugin file. `mac-cli/misc/help` contains a separate hardcoded category and usage registry. Bash and Fish completion scrape the command array from `/usr/local/bin/mac`, which breaks custom install locations and duplicates registry logic. Plugins are flat shell files that dispatch through `case "$fn"` and depend on global variables such as `fn`, `firstParameter`, `allParameters`, colors, and `echocommand`.

The first audit also identified plugin-specific technical debt:

- `git` has duplicate `git:open` cases, an apparent typo command `brew:git:config`, and a mismatch between `git:branches` and documented `git:branches:date`.
- `dns` edits `/etc/resolv.conf`, which is not the preferred durable macOS DNS path.
- `network` uses `sudo lsof` for `ports` and `wget` for public IP lookup.
- `performance` uses `sudo rm -rf ~/.Trash/*` for user trash cleanup.
- `search` searches `/` with `sudo find` for recent files.
- `ssh` repeats prompt logic, has unquoted variables, includes typo-like `${firstparameter}` references, and contains bare string statements where `echo` was likely intended.
- Existing completions hardcode `/usr/local/bin/mac` and have no Zsh implementation.

## Phase 1 Architecture

Phase 1 uses JSON metadata as the human-editable source of truth and generated shell-native artifacts as the runtime source of truth.

The recommended layout is:

```text
mac
mac-cli/
  bin/
    generate-registry
  plugins/
    git/
      plugin.sh
      commands.json
    dns/
      plugin.sh
      commands.json
  generated/
    registry.sh
    help.sh
  completion/
    _mac-cli
    mac.bash
    mac.fish
```

A lower-churn transition layout may temporarily keep flat implementation files with sibling metadata files:

```text
mac-cli/plugins/git
mac-cli/plugins/git.commands.json
```

The implementation plan should choose one transition strategy and apply it consistently. Plugin directories are preferred for long-term maintainability because they allow plugin-local helpers, fixtures, tests, and documentation without crowding the top-level plugin folder.

## Plugin Interface

Each plugin owns its metadata and implementation.

Required files for the long-term plugin directory format:

- `plugin.sh`: shell implementation for the plugin's commands.
- `commands.json`: metadata for discovery, help, validation, and completion.

Optional plugin-local paths:

- `lib/`: helpers private to the plugin.
- `README.md`: plugin-specific developer or user notes.
- `tests/`: plugin-specific test fixtures or smoke tests.

A minimal `commands.json` looks like:

```json
{
  "plugin_name": "git",
  "category": "git",
  "description": "Git utilities",
  "commands": [
    {
      "name": "git:open",
      "description": "Open current repository in the browser",
      "usage": "mac git:open",
      "args": [],
      "completion": []
    }
  ]
}
```

Command names should keep the existing colon-separated public syntax in Phase 1. Future metadata can add aliases, flags, validators, and richer completion hints without changing the runtime registry model.

## Registry Generation

`mac-cli/bin/generate-registry` reads all plugin `commands.json` files and writes generated artifacts.

Generated artifacts:

- `mac-cli/generated/registry.sh`: shell-readable command and plugin registry for the core CLI.
- `mac-cli/generated/help.sh`: shell-readable help/category data or help-rendering functions.
- `mac-cli/completion/_mac-cli`: generated Zsh completion.
- `mac-cli/completion/mac.bash`: generated Bash completion.
- `mac-cli/completion/mac.fish`: generated Fish completion.

The generator may use Python 3, Node, or another build-time tool, but generated runtime files must not require JSON parsing.

### Collision handling

The generator must fail clearly when two plugins define the same command name. A collision is an error, not a warning, because silently choosing one implementation would make dispatch and help output unpredictable.

A collision error should include:

- The duplicated command name.
- The first plugin metadata file that defined it.
- The second plugin metadata file that attempted to define it.
- A non-zero exit status so install/update can stop before writing inconsistent generated files.

The generator should write artifacts atomically where practical: generate into temporary files first, validate them, then move them into place.

## Generated Runtime Registry

`registry.sh` should provide enough shell-native data for the core CLI to validate commands and load implementations without parsing JSON.

Minimum required values:

```sh
MAC_COMMANDS="list help usage categories git:open dns:flush"
MAC_PLUGINS="git dns general network"
MAC_PLUGIN_FILES="/path/to/plugins/git/plugin.sh /path/to/plugins/dns/plugin.sh"
```

It may also provide helper functions such as:

```sh
mac_command_exists() { ...; }
mac_command_description() { ...; }
mac_command_usage() { ...; }
```

The implementation plan should prefer simple generated shell functions when they reduce fragile string parsing in the core CLI.

## Core CLI Runtime Flow

At runtime, `mac` should:

1. Resolve `MAC_CLI_ROOT` relative to the `mac` executable.
2. Source `mac-cli/generated/registry.sh`.
3. Source `mac-cli/generated/help.sh`.
4. Normalize the requested command into `fn`, preserving current argument variables for compatibility.
5. Handle built-in commands such as `help`, `usage`, `list`, and `categories` through generated help data.
6. Validate command existence through the generated registry.
7. Dynamically source plugin implementation files.
8. Let existing `case "$fn"` plugin blocks dispatch commands during Phase 1.

This flow removes hardcoded plugin source lists and hardcoded command arrays while preserving current behavior.

## Completion Strategy

Completion files should be generated from the same metadata used for the runtime registry.

### Zsh

The generated Zsh file should be named `_mac-cli` for the repository artifact and installed so Zsh can complete the `mac` command. It should follow standard Zsh completion patterns:

- `#compdef mac`
- `_arguments`
- `zstyle`-friendly command descriptions
- generated command-description pairs

### Bash

The generated Bash completion should complete commands from generated metadata rather than scraping `/usr/local/bin/mac`.

### Fish

The generated Fish completion should emit command completions and descriptions from generated metadata.

## Install and Update Integration

The installer should default to the single-directory layout approved for this fork:

```text
~/.mac-cli/
  mac
  mac-cli/
```

The installer should add `~/.mac-cli` or the selected custom install root to `PATH` through an idempotent shell profile block. It should also configure completion paths for Zsh where practical.

Install and update must run `mac-cli/bin/generate-registry` after files are copied into place. If registry generation fails, install/update should report the error clearly and avoid presenting the install as successful.

## Maintenance Rules

- Metadata is the canonical source for command names, descriptions, categories, usage, and completion hints.
- Generated files should not be edited manually.
- Adding a command requires updating the plugin implementation and that plugin's `commands.json`.
- Completion behavior must come from generated artifacts, not from parsing shell source.
- Plugin-local helper variables and functions should be namespaced by plugin during Phase 2 refactors.
- High-risk commands should prefer dry-run modes, confirmation prompts, narrower defaults, or non-destructive behavior.

## Phase 1 Success Criteria

- `mac` no longer has a manually maintained `COMMANDS` array.
- `mac` no longer manually lists every plugin source file.
- `misc/help` no longer maintains a separate hardcoded command registry.
- Bash, Fish, and Zsh completion are generated from plugin metadata.
- Completion no longer references `/usr/local/bin/mac`.
- Normal command execution and completion do not parse JSON.
- Duplicate command names fail registry generation with a clear error.
- Existing public command names continue to work.

## Deferred Phase 2 Work

Phase 2 should convert plugins from `case "$fn"` dispatch to namespaced handler functions, reduce global state, normalize argument parsing, and modernize risky or deprecated command implementations. Priority candidates are DNS configuration, trash cleanup, root filesystem search, SSH prompting, duplicate Git commands, and shell quoting throughout plugins.
