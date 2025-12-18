# AIOC-Util

Command-line utility for configuring AIOC (All-In-One-Cable) hardware settings. Available in both Python and Go versions.

This utility is based on code from Hrafnkell Eiríksson TF3HR, G1LRO, and Simon Küppers/skuep. The Go port provides native cross-platform binaries without requiring a Python runtime.

**Note:** This code is only tested against AIOC firmware v1.3+. Firmware v1.2 does not seem to work. Please upgrade your AIOC if needed.

## Features

- View and modify AIOC internal registers
- Configure PTT (Push-to-Talk) sources
- Set up foxhunt mode for radio direction finding
- Configure audio settings (RX gain, TX boost)
- Support for custom USB VID/PID
- Available as Python script or native Go binary

## Installation

### Go Version (Recommended)

Download pre-compiled binaries from the [Releases](https://github.com/rampa069/aioc-util/releases) page for:
- Linux (amd64)
- macOS (amd64, arm64)
- Windows (amd64)

#### Linux Setup

```bash
# Download and install binary
wget https://github.com/rampa069/aioc-util/releases/latest/download/aioc-util-linux-amd64
chmod +x aioc-util-linux-amd64
sudo mv aioc-util-linux-amd64 /usr/local/bin/aioc-util

# Install udev rule for non-root access
sudo cp udev/91-aioc.rules /etc/udev/rules.d/
sudo udevadm control --reload
sudo udevadm trigger
```

Unplug and replug your AIOC USB device after installing the udev rule.

#### macOS Setup

```bash
# Download and install binary
wget https://github.com/rampa069/aioc-util/releases/latest/download/aioc-util-darwin-arm64  # for Apple Silicon
# or
wget https://github.com/rampa069/aioc-util/releases/latest/download/aioc-util-darwin-amd64  # for Intel
chmod +x aioc-util-darwin-*
sudo mv aioc-util-darwin-* /usr/local/bin/aioc-util
```

#### Windows Setup

1. Download `aioc-util-windows-amd64.exe` from [Releases](https://github.com/rampa069/aioc-util/releases)
2. Rename to `aioc-util.exe`
3. Place in a directory in your PATH

### Python Version

#### Requirements
- Python 3
- `hid` Python package
- A hid shared library for your platform

#### Linux

```bash
git clone https://github.com/rampa069/aioc-util.git
cd aioc-util
python3 -m venv venv
source venv/bin/activate
pip install hid

# Install udev rule
sudo cp udev/91-aioc.rules /etc/udev/rules.d/
sudo udevadm control --reload
sudo udevadm trigger

# Install hidapi libraries if needed
sudo apt install libhidapi-hidraw0 libhidapi-libusb0
```

#### Windows

```bash
git clone https://github.com/rampa069/aioc-util.git
cd aioc-util
python3 -m venv venv
.\venv\Scripts\activate
pip install hid
```

On Windows, download `hidapi.dll` from the [hidapi releases](https://github.com/libusb/hidapi/releases) and place it in the project root directory.

## Usage

### Basic Commands

```bash
# List all available options
aioc-util --help

# Dump all registers and current configuration
aioc-util --dump

# List all possible PTT sources
aioc-util --list-ptt-sources
```

### PTT Configuration

```bash
# Set PTT1 to virtual PTT
aioc-util --ptt1 VPTT --store

# Set PTT1 to multiple sources (combined with |)
aioc-util --ptt1 "CM108GPIO1|SERIALDTR" --store

# Set PTT2 to CM108 GPIO
aioc-util --ptt2 CM108GPIO2 --store

# Swap PTT1 and PTT2 sources
aioc-util --swap-ptt --store

# Key/unkey radio manually
aioc-util --set-ptt1-state on   # Key the radio
aioc-util --set-ptt1-state off  # Unkey the radio
```

### VPTT/VCOS Configuration

```bash
# Configure virtual PTT and COS control registers
aioc-util --vptt-lvlctrl 0x80 --vptt-timctrl 10 --vcos-lvlctrl 0xff --vcos-timctrl 20 --store

# Enable hardware COS (if your AIOC supports it)
aioc-util --enable-hwcos --store

# Enable virtual COS (default behavior)
aioc-util --enable-vcos --store
```

### Audio Settings

```bash
# Set RX gain (1x, 2x, 4x, 8x, or 16x)
aioc-util --audio-rx-gain 4x --store

# Enable TX boost
aioc-util --audio-tx-boost on --store

# View current audio settings
aioc-util --audio-get-settings
```

### Foxhunt Mode

The AIOC firmware v1.4+ includes a foxhunt mode for radio direction finding activities. The AIOC only needs USB power (e.g., from a power bank) in this mode.

```bash
# Check current foxhunt settings
aioc-util --foxhunt-get-settings --foxhunt-get-message

# Set up a basic foxhunt beacon
aioc-util --foxhunt-message "DE TF0FOX" --foxhunt-wpm 20 --foxhunt-interval 60 --foxhunt-volume 32000 --store

# Configure just the transmission speed
aioc-util --foxhunt-wpm 15 --store

# Set a new message without changing other settings
aioc-util --foxhunt-message "VVV DE F0XX" --store

# Disable foxhunt mode
aioc-util --foxhunt-interval 0 --store
```

**Foxhunt Parameters:**
- `--foxhunt-volume`: Audio output level (0-65535)
- `--foxhunt-wpm`: Morse code speed in words per minute (0-255)
- `--foxhunt-interval`: Time in seconds between transmissions (0 disables foxhunt mode)
- `--foxhunt-message`: Up to 16 character text message

### Custom USB VID/PID

```bash
# Open AIOC with custom VID/PID
aioc-util --open-usb 0x1234,0x5678 --dump

# Set custom VID/PID (to emulate CM108 for example)
aioc-util --set-usb 0x0d8c,0x000c --store
```

### Other Commands

```bash
# Load hardware defaults
aioc-util --defaults --store

# Reboot the device
aioc-util --reboot

# Store current settings to flash
aioc-util --store
```

## Application Examples

Before using these configurations, reset to defaults:

```bash
aioc-util --defaults --store
```

### APRSDroid

Enable virtual PTT so you don't have to rely on VOX:

```bash
aioc-util --ptt1 VPTT --store
```

### AllStarLink 3

Set up an AllStarLink node with AIOC:

```bash
# Set VCOS timing control
aioc-util --vcos-timctrl 1500 --store
```

ASL3 supports AIOC on its default USB VID/PID values. Edit `/etc/asterisk/res_usbradio.conf` and uncomment the AIOC USB VID/PID line.

Alternatively, change VID/PID to emulate CM108:

```bash
aioc-util --set-usb 0x0d8c,0x000c --store
```

## Finding USB VID/PID

**Linux:**
```bash
lsusb
```
Look for your device's VID:PID pair.

**Windows (PowerShell):**
```powershell
Get-PnpDevice -PresentOnly | Where-Object { $_.InstanceId -like "USB\VID*" } | Select-Object Name, InstanceId
```

## Building from Source

### Go Version

#### Prerequisites
- Go 1.21 or later
- C compiler (gcc, clang, or MSVC)
- libhidapi development files (Linux only)

#### Linux
```bash
sudo apt-get install libudev-dev libusb-1.0-0-dev
cd go
go mod download
go build -o aioc-util .
```

#### macOS
```bash
brew install hidapi
cd go
go mod download
go build -o aioc-util .
```

#### Windows
```bash
cd go
go mod download
go build -o aioc-util.exe .
```

## Credits

- Original Python version: Hrafnkell Eiríksson TF3HR
- Based on code from: G1LRO and Simon Küppers/skuep
- Go port: EA5IUE

## License

Same as the original aioc-util project.
