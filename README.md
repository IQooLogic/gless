# GLess

[![CI](https://github.com/iqoologic/gless/workflows/CI/badge.svg)](https://github.com/iqoologic/gless/actions/workflows/ci.yml)
[![Release](https://github.com/iqoologic/gless/workflows/Release/badge.svg)](https://github.com/iqoologic/gless/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/iqoologic/gless)](https://goreportcard.com/report/github.com/iqoologic/gless)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A terminal pager like `less` but with full ANSI escape code support, perfect for viewing colorful log files.

## Features

- 🎨 **Full ANSI support** - Preserves colors and formatting (bold, italic, underline, etc.)
- 🚀 **Fast and efficient** - Handles large files with buffering
- ⌨️ **Familiar keybindings** - Similar to `less` with vim-style navigation
- 🔍 **Search functionality** - Find text with case-insensitive search
- 📊 **Line numbers** - Optional line number display
- 🎯 **Multiple color modes** - Supports 8/16/256-color and RGB ANSI codes

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/iqool/gless/releases).

### Build from Source

**Quick build:**
```bash
cd gless
go build -o gless.exe
```

Or install directly:
```bash
go install github.com/iqool/gless@latest
```

**Cross-platform builds using Makefile:**
```bash
# Build for all platforms
make build-all

# Build for specific platforms
make build-linux    # Linux amd64, arm64, arm
make build-macos    # macOS amd64, arm64 (Apple Silicon)
make build-windows  # Windows amd64, arm64

# See all available targets
make help
```

## Usage

View a file:
```bash
gless myfile.log
```

Read from stdin:
```bash
cat colorful.log | gless -
```

## Keyboard Shortcuts

### Navigation
- `↑`, `k` - Move up one line
- `↓`, `j` - Move down one line
- `PageUp`, `b` - Move up one page
- `PageDown`, `f`, `Space` - Move down one page
- `u` - Move up half page
- `d` - Move down half page
- `Home`, `g` - Go to first line
- `End`, `G` - Go to last line

### Search
- `/` - Enter search mode
- `n` - Next search result
- `N` - Previous search result

### Display
- `#` - Toggle line numbers

### Other
- `h`, `?` - Show help
- `q`, `Ctrl+C` - Quit

## Why GLess?

Standard `less` can display ANSI codes, but many log viewers and terminal pagers don't handle them well. GLess is specifically designed to:
- Parse and preserve ANSI escape sequences
- Display colorful application logs correctly
- Handle output from tools like `docker logs`, `kubectl logs`, etc.

## License

MIT
