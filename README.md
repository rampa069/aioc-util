# AIOC-Util (Go Version)

Cross-platform utility for configuring AIOC (All-In-One-Cable) hardware settings, written in Go.

This is a port of the Python aioc-util to Go, providing native cross-platform binaries without requiring a Python runtime.

## Features

- Native binaries for Windows, Linux, and macOS
- No runtime dependencies (statically linked)
- Complete feature parity with Python version
- HID USB communication with AIOC devices

## Building

### Prerequisites

- Go 1.21 or later
- C compiler (gcc, clang, or MSVC)
- libhidapi development files (Linux only)

### Linux

```bash
# Install dependencies (Debian/Ubuntu)
sudo apt-get install libhidapi-dev

# Build
cd go
go mod download
go build -o aioc-util
```

### macOS

```bash
# Install dependencies
brew install hidapi

# Build
cd go
go mod download
go build -o aioc-util
```

### Windows

```bash
# No additional dependencies needed (uses Windows HID API)
cd go
go mod download
go build -o aioc-util.exe
```

### Cross-compilation

```bash
# For Windows from Linux/macOS
GOOS=windows GOARCH=amd64 go build -o aioc-util.exe

# For Linux from macOS/Windows
GOOS=linux GOARCH=amd64 go build -o aioc-util

# For macOS from Linux/Windows
GOOS=darwin GOARCH=amd64 go build -o aioc-util
```

## Usage

The Go version has identical command-line arguments to the Python version:

```bash
# Dump all registers
./aioc-util --dump

# Set PTT1 source
./aioc-util --ptt1 VPTT --store

# Configure foxhunt mode
./aioc-util --foxhunt-interval 60 --foxhunt-wpm 20 --foxhunt-message "DE N0CALL" --store

# Configure audio settings
./aioc-util --audio-rx-gain 4x --audio-tx-boost on --store

# List all PTT sources
./aioc-util --list-ptt-sources

# Get help
./aioc-util --help
```

## Installation

After building, you can copy the binary to a location in your PATH:

```bash
# Linux/macOS
sudo cp aioc-util /usr/local/bin/

# Or for user-only install
mkdir -p ~/.local/bin
cp aioc-util ~/.local/bin/
```

On Windows, add the directory containing `aioc-util.exe` to your PATH environment variable.

## Dependencies

The Go version uses:
- [go-hid](https://github.com/sstallion/go-hid) - Go bindings for the hidapi library

## License

Same as the Python version of aioc-util.
