<div align="center">

# 🔌 nanocom

### A Modern, Cross-Platform Serial Communication Terminal

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue?style=for-the-badge)]()

*A lightweight, feature-rich serial terminal with a beautiful minicom-inspired TUI, built entirely in Go*

[Features](#features) • [Why nanocom?](#why-nanocom) • [Getting Started](#getting-started) • [Usage](#usage) • [Architecture](#architecture) • [Contributing](#contributing)

</div>

---

## 🎯 Overview

**nanocom** is a modern serial communication program designed for developers, embedded engineers, and IoT enthusiasts who need a reliable, cross-platform tool for interacting with serial devices. Whether you're debugging an Arduino, configuring a router, or communicating with embedded systems, nanocom provides an intuitive terminal interface that just works.

Built with the elegant [Charm](https://charm.sh/) ecosystem, nanocom combines the familiarity of classic tools like minicom with modern Go performance and a beautiful, responsive TUI.

## Why nanocom?

- 💡 **Developer-first ergonomics** – Every interaction is mapped to predictable key bindings, live reconfiguration, and clear visual feedback.
- ⚙️ **Instant onboarding** – Single static binary, zero dependencies, and intuitive defaults to start shipping logs in under a minute.
- 🛡️ **Production-grade reliability** – Asynchronous I/O, precise buffering, and configurable session logging keep critical data intact.
- 🤝 **Friendly Contribution Surface** – Modular Go packages, documented architecture, and formatter-friendly code make pull requests painless.

## Features

### Developer Experience

- 🖥️ Cross-platform support across Linux, macOS, and Windows with native serial handling
- 🎨 Minicom-inspired TUI with adaptive colors, rounded borders, and subtle animations
- ⚡ Real-time asynchronous I/O so keystrokes, logs, and transfers never block each other

### Productivity Boosters

- 📁 Built-in X/Y/Zmodem flows powered by [trzsz-go](https://github.com/trzsz/trzsz-go) for reliable file transfer sessions (Under Development 🍃)
- 📝 Session logging with rolling files for copy/paste-friendly transcripts
- 🔧 Live-tunable serial settings (baud, data bits, parity, endings) without restarting

### Power Tools

- 🔍 Auto port discovery and smart defaults for USB-UART devices
- ⌨️ Vim-style navigation with both arrow and hjkl bindings
- 🚀 Zero external dependencies—ship a single binary to your team

### Additional Features

- **🔍 Auto Port Discovery** - Automatically detects available serial ports
- **⌨️ Vim-Style Navigation** - Use `j/k` or arrow keys for menu navigation
- **🎯 Command Mode** - Minicom-compatible `Ctrl-A` command prefix
- **💾 Persistent Settings** - Configure once, use everywhere
- **🚀 Zero Dependencies** - Single binary, no runtime dependencies

## Getting Started

1. **Install the CLI**

    ```bash
    go install github.com/shoaibashk/nanocom@latest
    ```

2. **(Optional) Build From Source**

    ```bash
    git clone https://github.com/shoaibashk/nanocom.git
    cd nanocom
    go build -o nanocom .
    ```

3. **Validate Your Setup**

    ```bash
    nanocom --help
    ```

4. **Connect**

    ```bash
    nanocom -p /dev/ttyUSB0 -b 115200   # Linux / macOS
    nanocom -p COM3 -b 115200           # Windows
    ```

### Developer Checklist

- [x] Go 1.24+ installed
- [x] Serial device connected (USB, UART, etc.)
- [x] User in dialout/serial group (Linux)
- [x] Terminal supports truecolor for best TUI experience

## Usage

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

## 📁 File Transfers (X/Y/Zmodem)

nanocom routes file transfers through the embedded [trzsz-go](https://github.com/trzsz/trzsz-go) filter, so you get a modern trz/tsz workflow while staying compatible with classic `rz`/`sz` tooling. When a transfer is staged, the status bar shows a spinner so you can tell the terminal is waiting on the remote endpoint.

### Sending files to the remote host

1. Press `Ctrl+A`, then `S`, and enter the local file or directory path.
2. After the dialog closes, look for the terminal hint that lists the command to run remotely.
3. Run the matching command on the device you are connected to:

    | Protocol | Remote command |
    |----------|----------------|
    | Zmodem   | `rz -y`        |
    | Ymodem   | `rb -y`        |
    | Xmodem   | `rx`           |

4. The spinner disappears and nanocom logs the result once trzsz reports completion or an error.

### Receiving files from the remote host

1. Press `Ctrl+A`, then `R`, and choose the directory where incoming files should be saved.
2. Ask the remote endpoint to start the transfer with the appropriate command:

    | Protocol | Remote command example |
    |----------|------------------------|
    | Zmodem   | `sz firmware.bin`      |
    | Ymodem   | `sb *.bin`             |
    | Xmodem   | `sx boot.img`          |

3. Files are written to the directory you selected, and nanocom appends a completion line to the terminal buffer.

> **Tip:** Inside either transfer dialog you can press `Ctrl+P` to cycle protocols if the connected device only speaks a specific modem flavor.

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

## Architecture

nanocom is built with a **clean, modular architecture** following Go best practices. The TUI layer is organized into focused packages, each with a single responsibility:

```text
nanocom/
├── main.go                 # Application entry point
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go             # Root command & flags
│   ├── list.go             # Port listing command
│   └── config.go           # Configuration command
└── internal/
    └── tui/                # Terminal UI (Bubble Tea)
        ├── app.go          # Entry point & tea.Model wrapper
        ├── state/          # View states & messages
        │   ├── viewstate.go    # ViewState enum
        │   └── messages.go     # All tea.Msg types
        ├── keys/           # Key bindings
        │   └── keys.go         # Key map definitions
        ├── serial/         # Serial communication
        │   ├── connection.go   # Config & connection logic
        │   ├── reader.go       # Async serial reader
        │   └── ports.go        # Port discovery
        ├── transfer/       # File transfers
        │   ├── state.go        # Transfer state & protocol
        │   ├── trzsz.go        # trzsz bridge
        │   └── utils.go        # Path utilities
        ├── logging/        # Session capture
        │   └── state.go        # Logging state
        ├── menu/           # Menu definitions
        │   └── menu.go         # Menu items & constants
        ├── model/          # Core model
        │   └── model.go        # Main Model struct
        ├── update/         # Event handling
        │   └── update.go       # Update function & handlers
        ├── views/          # UI rendering
        │   └── views.go        # All view render functions
        └── styles/         # Visual styling
            └── styles.go       # Colors & lipgloss styles
```

### Package Responsibilities

| Package | Purpose |
|---------|---------|
| `state` | View state enum and all message types |
| `keys` | Key bindings configuration |
| `serial` | Serial port connection, reading, and port discovery |
| `transfer` | File transfer state, protocols, and trzsz bridge |
| `logging` | Session capture to file |
| `menu` | Menu item definitions |
| `model` | Main application model and state |
| `update` | Message handling and business logic |
| `views` | All UI rendering functions |
| `styles` | Colors, themes, and lipgloss styles |

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

## Configuration

### Serial Port Settings

Configure these settings via the TUI (`Ctrl-A O`) or command flags:

| Setting | Options | Default |
|---------|---------|---------|
| **Baud Rate** | 300 - 921600 | 9600 |
| **Data Bits** | 5, 6, 7, 8 | 8 |
| **Stop Bits** | 1, 2 | 1 |
| **Parity** | None, Even, Odd | None |
| **Line Ending** | NL, CR, CRLF, None | CRLF |

## Contributing

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

## Release Process

- Follow [Semantic Versioning](https://semver.org/) for every git tag: `vMAJOR.MINOR.PATCH` (for example `v0.2.0`).
- Create annotated tags so release notes include author/date: `git tag -a v0.2.0 -m "Release v0.2.0"`.
- Push the tag (`git push origin v0.2.0`) to trigger the `release` workflow, which signs artifacts with GoReleaser.
- Let GitHub Actions finish and confirm the release shows a verified signature before announcing the build.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Built with ❤️ by [Shoaibashk](https://github.com/shoaibashk)** · ⭐ Star the repo if it streamlines your serial workflows!
