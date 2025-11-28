# nanocom

A cross-platform serial communication program with a minicom-like TUI (Terminal User Interface).

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles), and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

## Features

- 🖥️ Cross-platform (Linux, macOS, Windows)
- 📡 Serial port communication
- 🎨 Colorful minicom-like TUI interface
- ⌨️ Keyboard-driven interface with Ctrl-A command mode
- 🔧 Configurable serial port settings (baud rate, data bits, stop bits, parity)
- 📋 Port discovery and selection

## Installation

```bash
go install github.com/shoaibashk/nanocom@latest
```

Or build from source:

```bash
git clone https://github.com/shoaibashk/nanocom.git
cd nanocom
go build
```

## Usage

```bash
# Start with TUI interface
nanocom

# Start with specific port
nanocom -p /dev/ttyUSB0

# Start with specific baud rate
nanocom -b 115200

# Start with both options
nanocom -p /dev/ttyUSB0 -b 115200
```

## Key Bindings

The TUI uses a minicom-like command mode activated with `Ctrl-A`:

| Key | Action |
|-----|--------|
| `Ctrl-A Z` | Show Help Menu |
| `Ctrl-A P` | Communication Parameters |
| `Ctrl-A C` | Clear Screen |
| `Ctrl-A O` | Configuration Menu |
| `Ctrl-A Q` | Quit |
| `Ctrl-A X` | Exit |
| `Ctrl-A A` | Send Ctrl-A literally |

### Navigation

| Key | Action |
|-----|--------|
| `↑/k` | Move up |
| `↓/j` | Move down |
| `←/h` | Decrease value |
| `→/l` | Increase value |
| `Enter` | Select/Confirm |
| `Esc` | Back/Cancel |

## Subcommands

```bash
# List available serial ports
nanocom list
```

## Screenshots

The TUI interface resembles minicom with:
- Title bar at the top
- Terminal area for serial I/O
- Input line for typing commands
- Status bar showing port info, connection status, and time
- Modal menus for configuration

## License

MIT
