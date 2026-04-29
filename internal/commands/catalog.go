package commands

import (
	"sort"
	"strings"
)

type Command struct {
	Name        string
	Category    string
	Description string
	Args        []string
	Dangerous   bool
	Native      bool
}

type Catalog struct {
	Commands []Command
}

func DefaultCatalog() Catalog {
	return Catalog{Commands: []Command{
		{Name: "list", Category: "Help", Native: true, Description: "List all commands by category."},
		{Name: "help", Category: "Help", Native: true, Description: "Show command help."},
		{Name: "usage", Category: "Help", Native: true, Description: "Show command usage."},
		{Name: "categories", Category: "Help", Native: true, Description: "List command categories."},
		{Name: "general", Category: "Help", Native: true, Description: "Show General command help."},
		{Name: "search", Category: "Help", Native: true, Description: "Show Search command help."},
		{Name: "network", Category: "Help", Native: true, Description: "Show Network command help."},
		{Name: "dns", Category: "Help", Native: true, Description: "Show DNS command help."},
		{Name: "ssh", Category: "Help", Native: true, Description: "Show SSH command help."},
		{Name: "webdev", Category: "Help", Native: true, Description: "Show Webdev command help."},
		{Name: "performance", Category: "Help", Native: true, Description: "Show Performance command help."},
		{Name: "terminal", Category: "Help", Native: true, Description: "Show Terminal command help."},
		{Name: "git", Category: "Help", Native: true, Description: "Show Git command help."},
		{Name: "web_utilities", Category: "Help", Native: true, Description: "Show Web Utilities command help."},
		{Name: "homebrew", Category: "Help", Native: true, Description: "Show Homebrew command help."},
		{Name: "image", Category: "Help", Native: true, Description: "Show Image command help."},
		{Name: "magento", Category: "Help", Native: true, Description: "Show Magento command help."},
		{Name: "update", Category: "General", Description: "Update macOS and installed package managers.", Dangerous: true},
		{Name: "upgrade", Category: "General", Description: "Upgrade Mac CLI.", Dangerous: true},
		{Name: "uninstall", Category: "General", Description: "Uninstall Mac CLI.", Dangerous: true},
		{Name: "lock", Category: "General", Description: "Lock the screen immediately."},
		{Name: "restart", Category: "General", Description: "Restart macOS.", Dangerous: true},
		{Name: "sleep", Category: "General", Description: "Put macOS to sleep."},
		{Name: "shutdown", Category: "General", Description: "Shutdown macOS.", Dangerous: true},
		{Name: "uptime", Category: "General", Description: "Show how long the system has been running."},
		{Name: "volume", Category: "General", Description: "Get or set system volume.", Args: []string{"level"}},
		{Name: "volume:ismute", Category: "General", Description: "Check whether audio output is muted."},
		{Name: "volume:mute", Category: "General", Description: "Mute audio output."},
		{Name: "volume:unmute", Category: "General", Description: "Unmute audio output."},
		{Name: "hidden:show", Category: "General", Description: "Show hidden files in Finder."},
		{Name: "hidden:hide", Category: "General", Description: "Hide hidden files in Finder."},
		{Name: "screensaver", Category: "General", Description: "Start the screensaver."},
		{Name: "display", Category: "General", Description: "Sleep the display immediately.", Native: true},
		{Name: "folders:list", Category: "General", Description: "List folders in the current directory with sizes."},
		{Name: "folder:size", Category: "General", Description: "Calculate current folder size."},
		{Name: "dock:add-space", Category: "General", Description: "Add a blank spacer to the Dock."},
		{Name: "bluetooth:status", Category: "General", Description: "Get Bluetooth status."},
		{Name: "bluetooth:enable", Category: "General", Description: "Enable Bluetooth."},
		{Name: "bluetooth:disable", Category: "General", Description: "Disable Bluetooth."},
		{Name: "wifi:status", Category: "General", Description: "Get Wi-Fi status."},
		{Name: "wifi:scan", Category: "General", Description: "Scan available Wi-Fi networks."},
		{Name: "wifi:enable", Category: "General", Description: "Enable Wi-Fi."},
		{Name: "wifi:disable", Category: "General", Description: "Disable Wi-Fi."},
		{Name: "eject-all", Category: "General", Description: "Eject mounted disks.", Dangerous: true},
		{Name: "battery", Category: "General", Description: "Show battery status."},
		{Name: "info", Category: "General", Description: "Show macOS version information."},
		{Name: "find:text", Category: "Search", Description: "Find text in files.", Args: []string{"pattern"}},
		{Name: "find:biggest-files", Category: "Search", Description: "Find largest files."},
		{Name: "find:biggest-directories", Category: "Search", Description: "Find largest directories."},
		{Name: "find:recent", Category: "Search", Description: "Find recently modified files."},
		{Name: "search:file", Category: "Search", Description: "Search for files by name using fd or find.", Args: []string{"pattern"}},
		{Name: "search:dir", Category: "Search", Description: "Search for directories by name.", Args: []string{"pattern"}},
		{Name: "search:content", Category: "Search", Description: "Search file contents using rg or grep.", Args: []string{"pattern"}},
		{Name: "search:recent", Category: "Search", Description: "Search recently modified files."},
		{Name: "search:spotlight", Category: "Search", Description: "Search with Spotlight.", Args: []string{"query"}},
		{Name: "search:dash", Category: "Search", Description: "Search Dash documentation.", Args: []string{"query"}},
		{Name: "search:devdocs", Category: "Search", Description: "Search DevDocs online.", Args: []string{"query"}},
		{Name: "search:replace", Category: "Search", Description: "Search and replace text.", Args: []string{"search", "replace"}},
		{Name: "network:speedtest", Category: "Network", Description: "Run a network speed test."},
		{Name: "speedtest", Category: "Network", Description: "Run a network speed test."},
		{Name: "network:ports", Category: "Network", Description: "List open ports."},
		{Name: "ports", Category: "Network", Description: "List open ports."},
		{Name: "network:ip:local", Category: "Network", Description: "Show local IP address."},
		{Name: "ip:local", Category: "Network", Description: "Show local IP address."},
		{Name: "network:public-ip", Category: "Network", Description: "Show public IP address."},
		{Name: "ip:public", Category: "Network", Description: "Show public IP address."},
		{Name: "network:cert", Category: "Network", Description: "Inspect an SSL certificate.", Args: []string{"domain"}},
		{Name: "network:ping", Category: "Network", Description: "Ping a host.", Args: []string{"host"}},
		{Name: "network:traceroute", Category: "Network", Description: "Trace route to a host.", Args: []string{"host"}},
		{Name: "network:info", Category: "Network", Description: "Show detailed network information."},
		{Name: "network:wifi-password", Category: "Network", Description: "Get current Wi-Fi password."},
		{Name: "network:dns-flush", Category: "Network", Description: "Flush DNS cache.", Dangerous: true},
		{Name: "dns:list", Category: "DNS", Description: "List current DNS servers."},
		{Name: "dns:add", Category: "DNS", Description: "Add a DNS server.", Args: []string{"server"}, Dangerous: true},
		{Name: "dns:remove", Category: "DNS", Description: "Remove a DNS server.", Args: []string{"server"}, Dangerous: true},
		{Name: "dns:flush", Category: "DNS", Description: "Flush DNS cache.", Dangerous: true},
		{Name: "hosts:edit", Category: "DNS", Description: "Edit /etc/hosts.", Dangerous: true},
		{Name: "ssh:download-file", Category: "SSH", Description: "Download a remote file over SSH.", Args: []string{"remote_path"}},
		{Name: "ssh:download-folder", Category: "SSH", Description: "Download a remote folder over SSH.", Args: []string{"remote_path"}},
		{Name: "ssh:sync:local", Category: "SSH", Description: "Sync remote folder to local folder.", Args: []string{"remote_path"}},
		{Name: "ssh:sync:remote", Category: "SSH", Description: "Sync local folder to remote folder.", Args: []string{"local_path"}},
		{Name: "ssh:upload", Category: "SSH", Description: "Upload a file over SSH.", Args: []string{"file"}},
		{Name: "ssh:public-key", Category: "SSH", Description: "Print SSH public key."},
		{Name: "ssh:list", Category: "SSH", Description: "List SSH utilities.", Native: true},
		{Name: "memory", Category: "Performance", Description: "Show memory usage."},
		{Name: "trash:empty", Category: "System", Description: "Empty the Trash.", Dangerous: true},
		{Name: "trash:size", Category: "System", Description: "Calculate Trash size."},
		{Name: "git:config", Category: "Git", Description: "Show Git configuration."},
		{Name: "git:open", Category: "Git", Description: "Open Git repository in browser."},
		{Name: "git:create:branch", Category: "Git", Description: "Create a branch.", Args: []string{"branch"}},
		{Name: "git:branches:date", Category: "Git", Description: "Show branch dates."},
		{Name: "git:undo-commit", Category: "Git", Description: "Undo last commit.", Dangerous: true},
		{Name: "git:log", Category: "Git", Description: "Show commit log."},
		{Name: "git:branch", Category: "Git", Description: "Show Git branches."},
		{Name: "git:branch:rename", Category: "Git", Description: "Rename a branch.", Args: []string{"old", "new"}},
		{Name: "git:branch:remove-local", Category: "Git", Description: "Remove a local branch.", Args: []string{"branch"}, Dangerous: true},
		{Name: "git:branch:remove-remote", Category: "Git", Description: "Remove a remote branch.", Args: []string{"branch"}, Dangerous: true},
		{Name: "git:settings", Category: "Git", Description: "Show Git settings."},
		{Name: "git:add-removed", Category: "Git", Description: "Stage removed files."},
		{Name: "git:size", Category: "Git", Description: "Get repository size."},
		{Name: "git:branch-fuzzy", Category: "Git", Description: "Select branch with fuzzy search."},
		{Name: "git:add-fuzzy", Category: "Git", Description: "Add files with fuzzy search."},
		{Name: "git:unpushed", Category: "Git", Description: "List unpushed commits."},
		{Name: "git:copy-branch", Category: "Git", Description: "Copy current branch name."},
		{Name: "brew:update", Category: "Homebrew", Description: "Update Homebrew and packages.", Dangerous: true},
		{Name: "file:clean-names", Category: "File", Description: "Clean filenames in a directory.", Args: []string{"dir", "--dry-run"}, Dangerous: true},
		{Name: "file:rename-clean", Category: "File", Description: "Rename one file to a clean name.", Args: []string{"file"}, Dangerous: true},
		{Name: "file:backup", Category: "File", Description: "Backup a file or directory.", Args: []string{"target"}},
		{Name: "file:move-link", Category: "File", Description: "Move a file and create a symlink.", Args: []string{"source", "dest"}, Dangerous: true},
		{Name: "file:compress", Category: "File", Description: "Compress a file.", Args: []string{"file"}},
		{Name: "file:extract", Category: "File", Description: "Extract an archive.", Args: []string{"file"}},
		{Name: "file:dataurl", Category: "File", Description: "Create a data URL.", Args: []string{"file"}},
		{Name: "archive:tar", Category: "Archive", Description: "Create a tar.gz archive.", Args: []string{"target", "output"}},
		{Name: "archive:zstd", Category: "Archive", Description: "Create a zstd-compressed tar archive.", Args: []string{"target", "output"}},
		{Name: "archive:gz-fast", Category: "Archive", Description: "Create a fast gzip archive.", Args: []string{"target", "output"}},
		{Name: "archive:extract", Category: "Archive", Description: "Extract an archive.", Args: []string{"file"}},
		{Name: "media:clip", Category: "Media", Description: "Clip media with ffmpeg.", Args: []string{"input", "start", "duration"}},
		{Name: "media:to-jxl", Category: "Media", Description: "Convert image to JXL.", Args: []string{"file"}},
		{Name: "media:to-avif", Category: "Media", Description: "Convert image to AVIF.", Args: []string{"file"}},
		{Name: "media:svg-to-png", Category: "Media", Description: "Convert SVG to PNG.", Args: []string{"file"}},
		{Name: "media:upload-imgbb", Category: "Media", Description: "Upload image to ImgBB.", Args: []string{"file"}},
		{Name: "system:kill", Category: "System", Description: "Kill processes matching a name or PID.", Args: []string{"pattern"}, Dangerous: true},
		{Name: "system:ports", Category: "System", Description: "List listening TCP ports."},
		{Name: "system:processes", Category: "System", Description: "List top processes by memory usage."},
		{Name: "system:cleanup-dsstore", Category: "System", Description: "Delete .DS_Store files recursively.", Dangerous: true},
		{Name: "system:unmount-all", Category: "System", Description: "Unmount external volumes.", Dangerous: true},
		{Name: "system:clean", Category: "System", Description: "Clean system caches.", Dangerous: true},
		{Name: "system:usage", Category: "System", Description: "Show CPU and memory usage."},
		{Name: "url:create", Category: "URL", Description: "Create a .webloc shortcut.", Args: []string{"url", "title"}},
		{Name: "url:create-force", Category: "URL", Description: "Create or replace a .webloc shortcut.", Args: []string{"url"}, Dangerous: true},
		{Name: "regex:urls", Category: "Regex", Description: "Extract URLs.", Args: []string{"file"}},
		{Name: "regex:emails", Category: "Regex", Description: "Extract email addresses.", Args: []string{"file"}},
		{Name: "regex:ips", Category: "Regex", Description: "Extract IP addresses.", Args: []string{"file"}},
		{Name: "regex:clean", Category: "Regex", Description: "Clean text with regex helpers.", Args: []string{"text"}},
		{Name: "regex:filename", Category: "Regex", Description: "Create a safe filename.", Args: []string{"text"}},
		{Name: "text:urlencode", Category: "Text", Description: "URL-encode text.", Args: []string{"text"}, Native: true},
		{Name: "text:urldecode", Category: "Text", Description: "URL-decode text.", Args: []string{"text"}, Native: true},
		{Name: "text:shortenurl", Category: "Text", Description: "Shorten a URL.", Args: []string{"url"}},
		{Name: "productivity:todo", Category: "Productivity", Description: "Append a todo.", Args: []string{"text"}},
		{Name: "productivity:note", Category: "Productivity", Description: "Append a note.", Args: []string{"text"}},
		{Name: "productivity:timer", Category: "Productivity", Description: "Start a timer.", Args: []string{"minutes"}},
		{Name: "dev:git-clear", Category: "Dev", Description: "Reset and clean Git changes.", Dangerous: true},
		{Name: "dev:git-save", Category: "Dev", Description: "Stash Git changes."},
		{Name: "dev:git-standup", Category: "Dev", Description: "Show recent Git activity."},
		{Name: "dev:git-switch", Category: "Dev", Description: "Switch Git branches interactively."},
		{Name: "dev:git-sync", Category: "Dev", Description: "Pull and push Git changes."},
		{Name: "dev:gitignore", Category: "Dev", Description: "Fetch a gitignore template.", Args: []string{"template"}},
		{Name: "dev:format-json", Category: "Dev", Description: "Format JSON.", Args: []string{"file"}},
		{Name: "dev:xcode-clear", Category: "Dev", Description: "Clear Xcode derived data.", Dangerous: true},
		{Name: "dev:xcode-recent", Category: "Dev", Description: "Open recent Xcode projects."},
		{Name: "web:download-video", Category: "Web", Description: "Download a video.", Args: []string{"url"}},
		{Name: "web:shorten-url", Category: "Web", Description: "Shorten a URL.", Args: []string{"url"}},
		{Name: "web:weather", Category: "Web", Description: "Show weather.", Args: []string{"location"}},
		{Name: "convert:c-to-f", Category: "Convert", Description: "Convert Celsius to Fahrenheit.", Args: []string{"celsius"}, Native: true},
		{Name: "convert:f-to-c", Category: "Convert", Description: "Convert Fahrenheit to Celsius.", Args: []string{"fahrenheit"}, Native: true},
		{Name: "convert:to-lowercase", Category: "Convert", Description: "Convert text to lowercase.", Args: []string{"text"}, Native: true},
		{Name: "convert:to-uppercase", Category: "Convert", Description: "Convert text to uppercase.", Args: []string{"text"}, Native: true},
		{Name: "convert:json-to-yaml", Category: "Convert", Description: "Convert JSON to YAML.", Args: []string{"file"}},
		{Name: "convert:yaml-to-json", Category: "Convert", Description: "Convert YAML to JSON.", Args: []string{"file"}},
		{Name: "finance:cagr", Category: "Finance", Description: "Calculate CAGR.", Args: []string{"beginning", "ending", "years"}},
		{Name: "finance:simple-interest", Category: "Finance", Description: "Calculate simple interest.", Args: []string{"principal", "rate", "years"}},
	}}
}

func (c Catalog) Categories() []string {
	seen := map[string]bool{}
	var out []string
	for _, cmd := range c.Commands {
		if !seen[cmd.Category] {
			seen[cmd.Category] = true
			out = append(out, cmd.Category)
		}
	}
	sort.Strings(out)
	return out
}

func (c Catalog) ByCategory(category string) []Command {
	var out []Command
	for _, cmd := range c.Commands {
		if strings.EqualFold(cmd.Category, category) {
			out = append(out, cmd)
		}
	}
	return out
}

func (c Catalog) Find(name string) (Command, bool) {
	for _, cmd := range c.Commands {
		if cmd.Name == name {
			return cmd, true
		}
	}
	return Command{}, false
}

func (cmd Command) Preview(args []string) string {
	parts := append([]string{"mac", cmd.Name}, args...)
	return strings.Join(parts, " ")
}
