# Mac CLI

 macOS command line tools for developers


---

### Sponsors

>[![Sponsor logo](images/sponsor-icon.png)](https://apps.apple.com/us/app/superplanner/id6443725564)
>
>Mac CLI is sponsored by 📒 [SuperPlanner](https://superplanner.app/), an innovative daily planner and task manager for iPhone, iPad and Mac.
>
>SuperPlanner combines the calendar with task management features to keep everything organized in one place.
>
>There is no login or user registration required. All data is stored locally and syncs between your devices using private and secure iCloud sync.
>
>[![Sponsor download badge](images/sponsor-download-badge.png)](https://apps.apple.com/us/app/superplanner/id6443725564)

---

### Introduction

Mac CLI is the ultimate tool for developers to manage their Mac. It provides a huge set of command line commands that automate the usage of your macOS system. When you run a function, the executed command is displayed, helping you memorize each of the utilities for future usage.

The tool is designed to be easily extendable with additional commands through the use of plugins. To view the currently available commands, you can navigate to the plugins folder and explore the different categories: [/mac-cli/plugins](https://github.com/guarinogabriel/mac-cli/tree/master/mac-cli/plugins)

_Contributions to add new plugins and keep improving the existing ones are welcome and very much appreciated!_

![image](images/demo.gif)

---

### Installation in 1 Simple Step - Including Configuration Wizard!

Via curl
> `sh -c "$(curl -fsSL https://raw.githubusercontent.com/guarinogabriel/mac-cli/master/mac-cli/tools/install)"`

Via wget
> `sh -c "$(wget -qO- https://raw.githubusercontent.com/guarinogabriel/mac-cli/master/mac-cli/tools/install)"`

---

### Features

Mac CLI offers a wide range of utilities categorized into plugins for easy management and extension. Here are some of the key categories and commands:

- **General Utilities**: Commands for system management, file operations, and more.
- **Search**: Tools for searching files, directories, content, and documentation (`search:file`, `search:spotlight`, `search:dash`).
- **Network**: Network diagnostics and information (`network:speedtest`, `network:ping`, `network:wifi-password`).
- **File**: File manipulation and management (`file:clean-names`, `file:backup`, `file:compress`).
- **Archive**: Archive creation and extraction (`archive:tar`, `archive:zstd`, `archive:extract`).
- **Media**: Media file processing (`media:clip`, `media:to-avif`, `media:upload-imgbb`).
- **System**: System resource management (`system:kill`, `system:cleanup-dsstore`, `system:usage`).
- **URL**: URL shortcut creation (`url:create`, `url:create-force`).
- **Regex**: Text extraction and manipulation using regular expressions (`regex:urls`, `regex:emails`).
- **Text**: Text processing and URL utilities (`text:urlencode`, `text:shortenurl`).
- **Productivity**: Task and time management (`productivity:todo`, `productivity:note`, `productivity:timer`).
- **Dev**: Developer tools for Git and Xcode (`dev:git-clear`, `dev:gitignore`, `dev:xcode-clear`).
- **Web**: Web-related tasks and downloads (`web:download-video`, `web:shorten-url`, `web:weather`).
- **Convert**: Conversion tools for temperature and formats (`convert:c-to-f`, `convert:json-to-yaml`).
- **Finance**: Financial calculations (`finance:cagr`, `finance:simple-interest`).

To see the full list of commands, run `mac list` or explore the plugins directory.

---

### Usage

Once installed, you can run Mac CLI commands using the `mac` prefix followed by the command name and any parameters:

```bash
mac system:usage
mac file:clean-names ./directory
mac web:download-video https://youtube.com/watch?v=example
```

Use `mac help` to view all available commands and categories.

---

### Extending Mac CLI

Mac CLI is built with a plugin architecture, making it easy to add new commands. To create a new plugin:

1. Add a new shell script in the `mac-cli/plugins/` directory with a unique name.
2. Follow the structure of existing plugins, using `case "$fn" in` to define commands.
3. Update the `COMMANDS` array in the main `mac` script to include your new commands.
4. Update the help content in `mac-cli/misc/help` to document your commands.

Feel free to contribute by submitting pull requests with new plugins or enhancements to existing ones.

---

### Community and Support

Join the Mac CLI community for discussions, support, and contributions:

- GitHub: [guarinogabriel/mac-cli](https://github.com/guarinogabriel/mac-cli)
- Issues: [Report bugs or suggest features](https://github.com/guarinogabriel/mac-cli/issues)

---

### Configuration

The configuration is done when you install Mac CLI for the first time though the installer configuration wizard.
After that, you can update your Mac CLI configuration by editing the following file: `/usr/local/bin/mac`

---

### Requirements

These are the requirements to be able to run all the commands (the dependencies/requirements are installed when you install Mac CLI for the first time):

* Homebrew
* Git
* Pipe Viewer (pv)

---

### Update

You can update Mac CLI to the latest version by running:
> `sh -c "$(curl -fsSL https://raw.githubusercontent.com/guarinogabriel/mac-cli/master/mac-cli/tools/update)"`

---

### Uninstallation

You can uninstall Mac CLI by running:
> `sh -c "$(curl -fsSL https://raw.githubusercontent.com/guarinogabriel/mac-cli/master/mac-cli/tools/uninstall)"`

---

### Help / Commands List

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac help`  | List all available commands in mac script  | |
| `mac list`  | List all available commands by category  | |
| `mac usage`  | Show usage information for commands  | |
| `mac categories`  | List all command categories  | |

### General Commands

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac update`  | Install macOS software updates, update installed Ruby gems, Homebrew, npm and their installed packages | |
| `mac lock`  | Lock  | |
| `mac restart`  | Restart macOS  | |
| `mac sleep`  | Sleep mode  | |
| `mac shutdown`  | Shutdown  | |
| `mac time`  | Show clock at top right position in Terminal/iTerm | |
| `mac screensaver`  | Start screensaver  | |
| `mac folders:list`  | List folders in current directory with their current size | |
| `mac folder:size`  | Calculate current folder size  | |
| `mac bluetooth:status`  | Get the bluetooth status  | |
| `mac bluetooth:enable`  | Enable bluetooth  | |
| `mac bluetooth:disable`  | Disable bluetooth  | |
| `mac wifi:status`  | Get the wifi status  | |
| `mac wifi:scan`  | Scan available wifi networks  | |
| `mac wifi:enable`  | Enable wifi  | |
| `mac wifi:disable`  | Disable wifi  | |
| `mac dock:add-space N`  | Add blank space to dock  | N = number of spaces |
| `mac eject-all`  | Eject all mounted volumes and disks  | |
| `mac battery`  | Get battery status  | |
| `mac info`  | Get macOS version information  | |
| `mac hidden:show`  | Show hidden files  | |
| `mac hidden:hide`  | Hide hidden files  | |
| `mac find:text X`  | Find exact phrase recursively inside directory | X = Text string |
| `mac find:biggest-files`  | Find biggest files inside directory  | |
| `mac find:biggest-directories`  | Find biggest directories inside directory | |
| `mac zip:extract X` | Extract Zip file to current folder | X = Zip file to extract |
| `mac gzip:compress X` | Compress current file using Gzip | X = File to compress |
| `mac gzip:extract X` | Extract Gzip file to current folder | X = Gzip file to extract |
| `mac tar:compress X`  | Compress X file/directory using tar with progress indicator | X = File or directory |
| `mac tar:extract X` | Extract tar file to current folder | X = Tar file to extract |

### Search Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac find:recent N`  | Find files modified in the last N minutes  |  N = number of minutes  |
| `mac search:replace X` | Search and replace string in file | X = File to perform the search and replace operation |
| `mac search:file X` | Search for files by name using fd or find | X = Pattern to search for |
| `mac search:dir X` | Search for directories by name using fd or find | X = Pattern to search for |
| `mac search:content X` | Search file contents for a pattern using rg or grep | X = Pattern to search for |
| `mac search:spotlight X` | Search using Spotlight for quick results | X = Query to search |
| `mac search:dash X` | Search for text in Dash documentation | X = Query to search |
| `mac search:devdocs X` | Search for text in DevDocs online documentation | X = Query to search |

### Network Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac speedtest`  | Internet connection speed test  | |
| `mac ports`  | List of used ports  | |
| `mac ip:local`  | Get local IP address  | |
| `mac ip:public`  | Get public IP address  | |
| `mac network:public-ip` | Get public IP address using curl | |
| `mac network:cert X` | Get SSL certificate information for a domain | X = Domain to check |
| `mac network:ping X` | Ping a host to check network connectivity | X = Host to ping |
| `mac network:traceroute X` | Trace the route to a host | X = Host to trace |
| `mac network:info` | Get detailed network information | |
| `mac network:wifi-password` | Get Wi-Fi password for the current network | |
| `mac network:dns-flush` | Flush DNS cache | |

### DNS Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac dns:list`  | List DNS server(s)  | |
| `mac dns:add`  | Add DNS server  | |
| `mac dns:remove`  | Remove DNS server  | |
| `mac dns:flush`  | Flush DNS cache  | |

### SSH Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac ssh:download-file X`  | Download file from remote server through SSH  |  X = Path of the remote file to download  |
| `mac ssh:download-folder X`  | Download entire folder from remote server through SSH  |  X = Path of the remote folder to download  |
| `mac ssh:sync:local X`  | Sync local folder with remote folder using rsync through SSH (download remote folder to local folder)  |  X = Path of the remote folder to sync to local folder  |
| `mac ssh:sync:remote X`  | Sync remote folder with local folder using rsync through SSH (upload local folder to remote folder)  |  X = Path of the remote folder to sync from local folder  |
| `mac ssh:upload X`  | Upload file to remote server through SSH  |  X = Path of the file to upload to the remote server  |
| `mac ssh:public-key`  | Copy SSH Public Key  |  |
| `mac ssh:list`  | List all the saved SSH credentials  |  |

### Performance and Maintenance Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac system`  | Show system information to review mac performance  | |
| `mac temp`  | Show temperature, fan and battery statistics  | |
| `mac memory`  | See memory usage sorted by memory consumption  | |
| `mac trash:empty`  | Empty trash | |
| `mac trash:size`  | Calculate trash size | |
| `mac system:kill X` | Kill a process by name or ID | X = Process name or ID |
| `mac system:ports` | List processes with open ports | |
| `mac system:processes` | List all running processes with details | |
| `mac system:cleanup-dsstore` | Remove .DS_Store files from current directory | |
| `mac system:unmount-all` | Unmount all external volumes | |
| `mac system:clean` | Remove common temporary files and caches | |
| `mac system:usage` | Show CPU, memory, disk usage stats | |

### Git Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac git:config`  | Display local Git configuration  | |
| `mac git:open`  | Open current repository on Github  | |
| `mac git:create:branch`  | Create branch based on current branch  | |
| `mac git:branches:date`  | Get last update date for all branches in current project  | |
| `mac git:undo-commit`  | Undo latest commit  | |
| `mac git:log`  | See latest commits IDs and titles for current branch  | |
| `mac git:branch`  | See all branches  | |
| `mac git:branch:rename`  | Rename Git branch | |
| `mac git:branch:remove-local`  | Remove local Git branch | |
| `mac git:branch:remove-remote`  | Remove local and remote Git branch | |
| `mac git:settings`  | Check Git settings  | |
| `mac git:add-removed`  | Add removed files to staged files  | |
| `mac git:size`  | Get size for current Git directory  | |
| `mac git:branch-fuzzy` | Fuzzy search and switch Git branch | |
| `mac git:add-fuzzy` | Fuzzy search and add files to staging | |
| `mac git:unpushed` | List unpushed commits | |
| `mac git:copy-branch` | Copy current branch name to clipboard | |
| `mac dev:git-clear` | Clear Git changes, resetting to the last commit | |
| `mac dev:git-save` | Save Git changes as a stash | |
| `mac dev:git-standup` | View Git standup report for recent activity | |
| `mac dev:git-switch` | Switch Git branch interactively if fzf is available | |
| `mac dev:git-sync` | Sync Git changes with remote repository | |
| `mac dev:gitignore X` | Generate a .gitignore file for a specific language or framework | X = Language or framework |

### Homebrew Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac brew`  | Get a list of installed Homebrew packages  | |
| `mac brew:update` | Update Homebrew and installed packages | |

### File Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac file:clean-names X` | Clean file names in directory (remove special characters) | X = Directory path |
| `mac file:rename-clean X` | Rename files with clean name (preview mode) | X = Directory path |
| `mac file:backup X` | Create a backup of a file | X = File to backup |
| `mac file:move-link X Y` | Move file to target and create symlink | X = Source file, Y = Target location |
| `mac file:compress X` | Compress file with tar.gz | X = File or directory to compress |
| `mac file:extract X` | Extract compressed file | X = File to extract |
| `mac file:dataurl X` | Convert file to data URL | X = File to convert |

### Archive Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac archive:tar X Y` | Create tar archive with progress | X = Source, Y = Output file (optional) |
| `mac archive:zstd X Y` | Create zstd compressed tar archive | X = Source, Y = Output file (optional) |
| `mac archive:gz-fast X Y` | Create fast gzip compressed tar archive | X = Source, Y = Output file (optional) |
| `mac archive:extract X` | Extract various archive formats | X = Archive file to extract |

### Media Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac image` | Optimize images in current directory | |
| `mac media:clip X Y Z W` | Clip media file (video/audio) | X = Input file, Y = Start time, Z = End time, W = Output file (optional) |
| `mac media:to-jxl X Y` | Convert image to JPEG XL | X = Input file, Y = Output file (optional) |
| `mac media:to-avif X Y` | Convert image to AVIF | X = Input file, Y = Output file (optional) |
| `mac media:svg-to-png X Y` | Convert SVG to PNG | X = Input file, Y = Output file (optional) |
| `mac media:upload-imgbb X` | Upload image to ImgBB | X = Image file to upload |

### System Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac volume`  | Get volume level  | |
| `mac volume:ismute`  | Check if volume is muted  | |
| `mac volume:mute`  | Mute volume  | |
| `mac volume:unmute`  | Unmute volume  | |
| `mac system:kill X` | Kill a process by name or ID | X = Process name or ID |
| `mac system:ports` | List processes with open ports | |
| `mac system:processes` | List all running processes with details | |
| `mac system:cleanup-dsstore` | Remove .DS_Store files from current directory | |
| `mac system:unmount-all` | Unmount all external volumes | |
| `mac system:clean` | Remove common temporary files and caches | |
| `mac system:usage` | Show CPU, memory, disk usage stats | |

### URL Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac url:create X` | Create URL shortcut on desktop | X = URL to create shortcut for |
| `mac url:create-force X` | Create URL shortcut on desktop (overwrite if exists) | X = URL to create shortcut for |

### Regex Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac regex:urls X` | Extract URLs from text or file | X = Text or file to process (optional) |
| `mac regex:emails X` | Extract email addresses from text or file | X = Text or file to process (optional) |
| `mac regex:ips X` | Extract IP addresses from text or file | X = Text or file to process (optional) |
| `mac regex:clean X Y` | Clean text with regex pattern | X = Pattern, Y = Text or file to clean |
| `mac regex:filename X` | Extract filename from path or URL | X = Path or URL to process |

### Text Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac text:urlencode X` | URL encode text | X = Text to encode |
| `mac text:urldecode X` | URL decode text | X = Text to decode |
| `mac text:shortenurl X` | Shorten URL using tinyurl | X = URL to shorten |

### Productivity Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac productivity:todo X` | Add or list tasks in a simple TODO list | X = Task to add (optional) |
| `mac productivity:note X` | Create or open a quick note | X = Note content or title (optional) |
| `mac productivity:timer X` | Set a timer for a specified duration | X = Minutes for timer |

### Developer Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac dev:git-clear` | Clear Git changes, resetting to the last commit | |
| `mac dev:git-save` | Save Git changes as a stash | |
| `mac dev:git-standup` | View Git standup report for recent activity | |
| `mac dev:git-switch` | Switch Git branch interactively if fzf is available | |
| `mac dev:git-sync` | Sync Git changes with remote repository | |
| `mac dev:gitignore X` | Generate a .gitignore file for a specific language or framework | X = Language or framework |
| `mac dev:format-json X` | Format JSON content using Python's json.tool | X = JSON file (optional) |
| `mac dev:xcode-clear` | Clear Xcode derived data to free up space | |
| `mac dev:xcode-recent` | Open the most recent Xcode project | |

### Web Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac web:download-video X` | Download a video from a URL using yt-dlp or youtube-dl | X = URL of video to download |
| `mac web:shorten-url X` | Shorten a URL using tinyurl | X = URL to shorten |
| `mac web:weather X` | Get current weather information | X = Location (optional) |

### Convert Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac convert:c-to-f X` | Convert temperature from Celsius to Fahrenheit | X = Temperature in Celsius |
| `mac convert:f-to-c X` | Convert temperature from Fahrenheit to Celsius | X = Temperature in Fahrenheit |
| `mac convert:to-lowercase X` | Convert text to lowercase | X = Text to convert (optional) |
| `mac convert:to-uppercase X` | Convert text to uppercase | X = Text to convert (optional) |
| `mac convert:json-to-yaml X` | Convert JSON to YAML (requires yq) | X = JSON file (optional) |
| `mac convert:yaml-to-json X` | Convert YAML to JSON (requires yq) | X = YAML file (optional) |

### Finance Utilities

| Command  | Description | Arguments |
| ------------- | ------------- | ------------- |
| `mac finance:cagr X Y Z` | Calculate Compound Annual Growth Rate (CAGR) | X = Beginning value, Y = Ending value, Z = Number of years |
| `mac finance:simple-interest X Y Z` | Calculate simple interest | X = Principal amount, Y = Rate of interest, Z = Time in years |
