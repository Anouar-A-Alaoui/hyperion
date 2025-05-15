# Hyperion

A customizable directory tree visualizer written in Go.

## Features

- Display directory structures with customizable options
- Exclude specific folders, file types, or exact file names
- Control the visualization depth
- Choose between Unicode or ASCII tree styles
- Enable colorized output for folders, files, and symlinks
- Show statistics about the scanned directory
- Display tables of the largest files
- Show visual charts of file size distribution
- Compact mode for more concise output

## Installation

### Using Go Install

```bash
go install github.com/Anouar-A-Alaoui/hyperion@latest
```

### Building from Source

```bash
git clone https://github.com/Anouar-A-Alaoui/hyperion.git
cd hyperion
go build -o hyperion .
```

## Usage

```
hyperion [flags]
```

### Command-Line Flags

| Flag                | Type      | Default            | Description                                         |
| ------------------- | --------- | ------------------ | --------------------------------------------------- |
| `--path`            | string    | `"."`              | Root directory to scan                              |
| `--exclude-folders` | string[]  | `["node_modules"]` | Folders to exclude from tree                        |
| `--show-files`      | bool      | `false`            | Whether to show files in output                     |
| `--exclude-files`   | string[]  | `[]`               | File extensions to exclude (e.g., `.exe`)           |
| `--exclude-names`   | string[]  | `[]`               | File names to exclude exactly (e.g., `config.json`) |
| `--max-depth`       | int       | `-1`               | Maximum depth to recurse (-1 for unlimited)         |
| `--unicode`         | bool      | `true`             | Use Unicode characters for pretty tree visuals      |
| `--color`           | bool      | `true`             | Use colors in output                                |
| `--bg-color`        | bool      | `false`            | Use background color for items                      |
| `--compact`         | bool      | `false`            | Enable compact tree layout                          |
| `--show-stats`      | bool      | `false`            | Show total files, dirs, size                        |
| `--stat-table`      | bool      | `false`            | Show a table of largest files and types             |
| `--stats-count`     | int       | `10`               | Number of top files to show in stats table          |
| `--chart`           | bool      | `false`            | Show a visual chart of file size distribution       |
| `--help`            | bool      | `false`            | Show usage and examples                             |
| `--about`           | bool      | `false`            | Show about                                          |
| `--version`         | bool      | `false`            | Show version                                        |

### Examples

```bash
# Basic usage (folders only)
hyperion

# Exclude folders and file types
hyperion --show-files --exclude-folders "bin,obj" --exclude-files ".exe,.dll"

# Show stats with Unicode and color
hyperion --show-files --unicode --color --show-stats

# Show top 15 largest files with chart
hyperion --show-files --stat-table --stats-count 15 --chart

# Compact view with background color
hyperion --show-files --compact --bg-color
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
