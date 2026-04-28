# AGENTS.md

## Project Overview

Mac CLI is a command-line tool for macOS developers, offering a wide range of utilities to automate system management, development workflows, and more. It is built with a plugin architecture to allow easy extension of commands. Key technologies include shell scripting (primarily `sh`), Git for version control, and macOS native tools.

The project structure includes:
- **Main Script**: `mac` - The entry point for all commands.
- **Plugins**: `mac-cli/plugins/` - Directory containing categorized shell scripts for various commands.
- **Help Content**: `mac-cli/misc/help` - Documentation for command usage.

## Setup Commands

- **Install Dependencies**: Mac CLI does not require external dependencies beyond macOS native tools for basic functionality. Some plugins may suggest installing additional tools like `fd`, `rg`, `yt-dlp`, etc., if not already present.
- **Clone Repository**: If working on development, clone the repository using:
  ```bash
  git clone https://github.com/guarinogabriel/mac-cli.git
  cd mac-cli
  ```
- **Run Installation Script**: For a fresh setup, use the installation script:
  ```bash
  sh -c "$(curl -fsSL https://raw.githubusercontent.com/guarinogabriel/mac-cli/master/mac-cli/tools/install)"
  ```

## Development Workflow

- **Adding New Commands**: To add a new command, create or update a plugin script in `mac-cli/plugins/`:
  1. Use the `case "$fn" in` structure to define command logic.
  2. Add the command to the `COMMANDS` array in the main `mac` script.
  3. Update help content in `mac-cli/misc/help`.
- **Testing Changes**: Test new commands by running them directly:
  ```bash
  mac <command>
  ```
- **Environment**: No specific environment variables are required beyond standard macOS shell configurations.

## Testing Instructions

- **Syntax Check**: Validate shell script syntax for plugins:
  ```bash
  sh -n mac-cli/plugins/*.sh
  ```
- **Manual Testing**: Run individual commands to ensure they function as expected. There are no automated test suites currently.
- **Dry Run**: Many destructive commands include an `echocommand` flag to preview actions without executing them:
  ```bash
  mac <command> --echo
  ```

## Code Style

- **Shell Scripting**: Use `#!/bin/sh` for compatibility. Follow POSIX-compliant shell scripting practices.
- **Naming Conventions**: Commands use a `category:command` format (e.g., `system:kill`, `file:backup`).
- **File Organization**: Plugins are organized by category in `mac-cli/plugins/` (e.g., `system.sh`, `file.sh`).
- **Formatting**: Keep scripts readable with consistent indentation (typically 4 spaces) and clear comments for sections and commands.

## Build and Deployment

- **No Build Process**: Mac CLI does not require a build step; it runs directly from shell scripts.
- **Deployment**: Updates are typically deployed by committing changes to the repository. No automated deployment pipeline is currently in place for end-users beyond the installation script.
- **Versioning**: Follow semantic versioning if updating the main tool. Update the version in documentation if applicable.

## Pull Request Guidelines

- **Title Format**: `[category] Brief description` (e.g., `[system] Add cleanup command for temporary files`).
- **Required Checks**: Ensure syntax validation with `sh -n` on modified scripts and manual testing of new commands.
- **Review Process**: Submit PRs to the main repository for review by maintainers.

## Additional Notes

- **Plugin Extensions**: When adding new plugins, ensure they align with existing categories or create a new category if necessary.
- **Common Gotchas**: Be cautious with destructive commands; always include confirmation prompts where data loss is possible.
- **Performance**: Keep commands lightweight to maintain quick execution on the command line.

