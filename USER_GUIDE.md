# Hyperion Installation and User Guide

## Overview

hyperion is a powerful command-line utility for visualizing directory structures with extensive customization options. It provides a clear view of your filesystem with options to filter out specific files or folders, limit depth, show statistics, and more.

## Installation

### Requirements

- Go 1.18 or higher

### Using Go Install

The easiest way to install hyperion is using Go's installation tools:

```bash
go install https://github.com/Anouar-A-Alaoui/hyperion@latest
```

### From Source

```bash
git clone https://https://github.com/Anouar-A-Alaoui/hyperion.git
cd hyperion
make install
```

## Basic Usage

To display the directory tree of the current directory:

```bash
hyperion
```

By default, this will only show directories. To include files:

```bash
hyperion --show-files
```

## Command Line Options

hyperion has many command-line options to customize its behavior:

### Path Options

| Flag               | Type    | Default           | Description                      |
|--------------------|---------|-----------------|------------------------------------|
| `--path`           | string  | `"."`           | Root directory to scan             |
| `--max-depth`      | int     | `-1`            | Maximum depth (-1 for unlimited)   |

### Filtering Options

| Flag               | Type      | Default            | Description                       |
|--------------------|-----------|--------------------|-----------------------------------|
| `--show-files`     | bool      | `false`            | Show files in output              |
| `--exclude-folders`| string[]  | `["node_modules"]` | Folders to exclude                |
| `--exclude-files`  | string[]  | `[]`               | File extensions to exclude        |
| `--exclude-names`  | string[]  | `[]`               | File names to exclude exactly     |

### Visual Style Options

| Flag               | Type      | Default           | Description                        |
|--------------------|-----------|-------------------|------------------------------------|
| `--unicode`        | bool      | `true`            | Use Unicode box-drawing characters |
| `--color`          | bool      | `true`            | Use colors in output               |
| `--bg-color`       | bool      | `false`           | Use background color for items     |
| `--compact`        | bool      | `false`           | Enable compact tree layout         |

### Statistics Options

| Flag               | Type      | Default           | Description                       |
|--------------------|-----------|-------------------|-----------------------------------|
| `--show-stats`     | bool      | `false`           | Show total counts and sizes       |
| `--stat-table`     | bool      | `false`           | Show table of largest files       |
| `--stats-count`    | int       | `10`              | Number of files in stats table    |
| `--chart`          | bool      | `false`           | Show file size distribution chart |

### Help

| Flag               | Type      | Default           | Description                       |
|--------------------|-----------|-------------------|-----------------------------------|
| `--help`           | bool      | `false`           | Show usage and examples           |

## Examples

### Basic Directory Listing

Show only directories:

```bash
hyperion
```

Show directories and files:

```bash
hyperion --show-files
```

### Filtering Examples

Exclude multiple folder types:

```bash
hyperion --exclude-folders "node_modules,bin,obj"
```

Show files but exclude certain extensions:

```bash
hyperion --show-files --exclude-files ".exe,.dll,.obj"
```

Exclude specific filenames:

```bash
hyperion --show-files --exclude-names "config.json,README.md"
```

Limit directory depth:

```bash
hyperion --max-depth 2
```

### Visual Styles

Use ASCII characters instead of Unicode:

```bash
hyperion --unicode=false
```

Disable colors:

```bash
hyperion --color=false
```

Use background colors for better visibility:

```bash
hyperion --bg-color
```

Use compact layout:

```bash
hyperion --compact
```

### Statistics

Show basic statistics:

```bash
hyperion --show-stats
```

Show stats with a table of largest files:

```bash
hyperion --show-files --stat-table
```

Show top 20 largest files:

```bash
hyperion --show-files --stat-table --stats-count 20
```

Show all statistics with chart:

```bash
hyperion --show-files --show-stats --stat-table --chart
```

## Tips & Tricks

- Use `--compact` for large directories to make the output more condensed
- Combine `--stat-table` with `--chart` to get a complete overview of your disk usage
- Use `--color` with terminals that support ANSI colors for better readability
- Use `--exclude-folders` to skip large vendor directories like `node_modules`, `vendor`, etc.
- Redirect output to a file with `hyperion > tree.txt` to save the tree structure

## Troubleshooting

### Unicode Characters Display as Boxes or Question Marks

If Unicode characters don't display correctly, use ASCII mode:

```bash
hyperion --unicode=false
```

### Performance Issues with Large Directories

- Exclude large directories: `--exclude-folders "node_modules,vendor,dist"`
- Limit traversal depth: `--max-depth 3`
- Disable file display: remove `--show-files`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
