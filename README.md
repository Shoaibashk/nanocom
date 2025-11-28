<div align="center">

# 🔌 nanocom

### A Modern, Cross-Platform Serial Communication Terminal

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue?style=for-the-badge)]()

*A lightweight, feature-rich serial terminal with a beautiful minicom-inspired TUI, built entirely in Go*

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Architecture](#-architecture) • [Contributing](#-contributing)

</div>

---

## 🎯 Overview

**nanocom** is a modern serial communication program designed for developers, embedded engineers, and IoT enthusiasts who need a reliable, cross-platform tool for interacting with serial devices. Whether you're debugging an Arduino, configuring a router, or communicating with embedded systems, nanocom provides an intuitive terminal interface that just works.

Built with the elegant [Charm](https://charm.sh/) ecosystem, nanocom combines the familiarity of classic tools like minicom with modern Go performance and a beautiful, responsive TUI.

## ✨ Features

<table>
<tr>
<td width="50%">

### 🖥️ Cross-Platform
Works seamlessly on Linux, macOS, and Windows with native serial port support

### 🎨 Beautiful TUI
Minicom-inspired interface with adaptive colors, rounded borders, and smooth interactions

### ⚡ Real-Time Communication
Asynchronous serial I/O with buffered reads and responsive input handling

</td>
<td width="50%">

### 📁 File Transfer
Built-in support for Zmodem, Xmodem, and Ymodem protocols

### 📝 Session Logging
Capture serial output to file for debugging and documentation

### 🔧 Fully Configurable
Adjust baud rate, data bits, stop bits, parity, and line endings on the fly

</td>
</tr>
</table>

### Additional Features

- **🔍 Auto Port Discovery** - Automatically detects available serial ports
- **⌨️ Vim-Style Navigation** - Use `j/k` or arrow keys for menu navigation
- **🎯 Command Mode** - Minicom-compatible `Ctrl-A` command prefix
- **💾 Persistent Settings** - Configure once, use everywhere
- **🚀 Zero Dependencies** - Single binary, no runtime dependencies

## 📦 Installation

### Using Go Install (Recommended)

```bash
go install github.com/shoaibashk/nanocom@latest
```

### Build from Source

```bash
git clone https://github.com/shoaibashk/nanocom.git
cd nanocom
go build -o nanocom .
```

### Verify Installation

```bash
nanocom --help
```

## 🚀 Usage

### Quick Start

```bash
# Launch the TUI interface
nanocom

# Connect to a specific port
nanocom -p /dev/ttyUSB0          # Linux/macOS
nanocom -p COM3                   # Windows

# Set baud rate
nanocom -b 115200

# Combine options
nanocom -p /dev/ttyUSB0 -b 115200
```

### List Available Ports

```bash
nanocom list
```

## ⌨️ Key Bindings

nanocom uses a minicom-compatible command mode. Press `Ctrl-A` followed by a command key:

### Command Mode (`Ctrl-A` + Key)

| Key | Action | Description |
|:---:|--------|-------------|
| `Z` | **Help Menu** | Display available commands |
| `O` | **Serial Settings** | Configure port parameters |
| `S` | **Send File** | Initiate file transfer |
| `R` | **Receive File** | Receive incoming file |
| `L` | **Toggle Logging** | Start/stop session capture |
| `C` | **Clear Screen** | Clear terminal buffer |
| `Q` | **Quit** | Exit without reset |
| `X` | **Exit & Reset** | Exit and reset terminal |

### Navigation

| Key | Action |
|:---:|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `←` / `h` | Decrease value |
| `→` / `l` | Increase value |
| `Enter` | Select / Confirm |
| `Esc` | Back / Cancel |

## 🏗️ Architecture

nanocom is built with a clean, modular architecture following Go best practices:

```
nanocom/
├── main.go                 # Application entry point
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go            # Root command & flags
│   ├── list.go            # Port listing command
│   └── config.go          # Configuration command
└── internal/
    └── tui/               # Terminal UI (Bubble Tea)
        ├── model.go       # TUI state & logic
        ├── styles.go      # Lipgloss styling
        └── transfer.go    # File transfer protocols
```

### Tech Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **CLI Framework** | [Cobra](https://github.com/spf13/cobra) | Command-line interface & flag parsing |
| **TUI Framework** | [Bubble Tea](https://github.com/charmbracelet/bubbletea) | Terminal UI with Elm architecture |
| **UI Components** | [Bubbles](https://github.com/charmbracelet/bubbles) | Pre-built TUI components |
| **Styling** | [Lip Gloss](https://github.com/charmbracelet/lipgloss) | Declarative terminal styling |
| **Serial I/O** | [go.bug.st/serial](https://github.com/bugst/go-serial) | Cross-platform serial port access |

### Design Patterns

- **Elm Architecture (TEA)** - Unidirectional data flow for predictable state management
- **Goroutine-based I/O** - Non-blocking serial reads with channel communication
- **Adaptive Theming** - Automatic light/dark mode based on terminal settings

## 🔧 Configuration

### Serial Port Settings

Configure these settings via the TUI (`Ctrl-A O`) or command flags:

| Setting | Options | Default |
|---------|---------|---------|
| **Baud Rate** | 300 - 921600 | 9600 |
| **Data Bits** | 5, 6, 7, 8 | 8 |
| **Stop Bits** | 1, 2 | 1 |
| **Parity** | None, Even, Odd | None |
| **Line Ending** | NL, CR, CRLF, None | CRLF |

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Setup

```bash
# Clone the repository
git clone https://github.com/shoaibashk/nanocom.git
cd nanocom

# Install dependencies
go mod download

# Run in development
go run .

# Run tests
go test ./...
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Built with ❤️ by [Shoaibashk](https://github.com/shoaibashk)**

⭐ Star this repo if you find it useful!

</div>
