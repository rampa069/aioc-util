package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/sstallion/go-hid"
)

const (
	AIOCVendorID  = 0x1209
	AIOCProductID = 0x7388
)

// Register addresses
type Register uint8

const (
	RegMAGIC        Register = 0x00
	RegUSBID        Register = 0x08
	RegAIOCIOMUX0   Register = 0x24
	RegAIOCIOMUX1   Register = 0x25
	RegCM108IOMUX0  Register = 0x44
	RegCM108IOMUX1  Register = 0x45
	RegCM108IOMUX2  Register = 0x46
	RegCM108IOMUX3  Register = 0x47
	RegSERIALCTRL   Register = 0x60
	RegSERIALIOMUX0 Register = 0x64
	RegSERIALIOMUX1 Register = 0x65
	RegSERIALIOMUX2 Register = 0x66
	RegSERIALIOMUX3 Register = 0x67
	RegAUDIORX      Register = 0x72
	RegAUDIOTX      Register = 0x78
	RegVPTTLVLCTRL  Register = 0x82
	RegVPTTTIMCTRL  Register = 0x84
	RegVCOSLVLCTRL  Register = 0x92
	RegVCOSTIMCTRL  Register = 0x94
	RegFOXHUNTCTRL  Register = 0xA0
	RegFOXHUNTMSG0  Register = 0xA2
	RegFOXHUNTMSG1  Register = 0xA3
	RegFOXHUNTMSG2  Register = 0xA4
	RegFOXHUNTMSG3  Register = 0xA5
)

// Command flags
type Command uint8

const (
	CmdNONE        Command = 0x00
	CmdWRITESTROBE Command = 0x01
	CmdDEFAULTS    Command = 0x10
	CmdREBOOT      Command = 0x20
	CmdRECALL      Command = 0x40
	CmdSTORE       Command = 0x80
)

// PTT Source flags
type PTTSource uint32

const (
	PTTSourceNONE          PTTSource = 0x00000000
	PTTSourceCM108GPIO1    PTTSource = 0x00000001
	PTTSourceCM108GPIO2    PTTSource = 0x00000002
	PTTSourceCM108GPIO3    PTTSource = 0x00000004
	PTTSourceCM108GPIO4    PTTSource = 0x00000008
	PTTSourceSERIALDTR     PTTSource = 0x00000100
	PTTSourceSERIALRTS     PTTSource = 0x00000200
	PTTSourceSERIALDTRNRTS PTTSource = 0x00000400
	PTTSourceSERIALNDTRRTS PTTSource = 0x00000800
	PTTSourceVPTT          PTTSource = 0x00001000
)

// CM108 Button Source flags
type CM108ButtonSource uint32

const (
	CM108ButtonSourceNONE CM108ButtonSource = 0x00000000
	CM108ButtonSourceIN1  CM108ButtonSource = 0x00010000
	CM108ButtonSourceIN2  CM108ButtonSource = 0x00020000
	CM108ButtonSourceVCOS CM108ButtonSource = 0x01000000
)

// RX Gain values
type RXGain uint32

const (
	RXGain1X  RXGain = 0x00000000
	RXGain2X  RXGain = 0x00000001
	RXGain4X  RXGain = 0x00000002
	RXGain8X  RXGain = 0x00000003
	RXGain16X RXGain = 0x00000004
)

// TX Boost values
type TXBoost uint32

const (
	TXBoostOFF TXBoost = 0x00000000
	TXBoostON  TXBoost = 0x00000100
)

// PTT Channel values
const (
	PTTChannel1 = 3
	PTTChannel2 = 4
)

// AIOCDevice represents an AIOC HID device
type AIOCDevice struct {
	device *hid.Device
}

// Open opens an AIOC device by VID/PID
func Open(vid, pid uint16) (*AIOCDevice, error) {
	device, err := hid.OpenFirst(vid, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to open device: %w", err)
	}

	aioc := &AIOCDevice{device: device}

	// Verify magic
	magic, err := aioc.Read(RegMAGIC)
	if err != nil {
		device.Close()
		return nil, fmt.Errorf("failed to read magic: %w", err)
	}

	magicBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(magicBytes, magic)
	if !bytes.Equal(magicBytes, []byte("AIOC")) {
		device.Close()
		return nil, fmt.Errorf("unexpected magic: %s", magicBytes)
	}

	return aioc, nil
}

// Close closes the device
func (a *AIOCDevice) Close() error {
	return a.device.Close()
}

// Read reads a 32-bit value from a register
func (a *AIOCDevice) Read(address Register) (uint32, error) {
	// Prepare request: [0, NONE, address, 0x00000000]
	request := make([]byte, 8)
	request[0] = 0
	request[1] = uint8(CmdNONE)
	request[2] = uint8(address)
	binary.LittleEndian.PutUint32(request[3:7], 0)

	if _, err := a.device.SendFeatureReport(request); err != nil {
		return 0, fmt.Errorf("failed to send feature report: %w", err)
	}

	// Read response
	response := make([]byte, 7)
	n, err := a.device.GetFeatureReport(response)
	if err != nil {
		return 0, fmt.Errorf("failed to get feature report: %w", err)
	}
	if n < 7 {
		return 0, fmt.Errorf("short read: got %d bytes, expected 7", n)
	}

	// Extract value (last 4 bytes)
	value := binary.LittleEndian.Uint32(response[3:7])
	return value, nil
}

// Write writes a 32-bit value to a register
func (a *AIOCDevice) Write(address Register, value uint32) error {
	data := make([]byte, 8)
	data[0] = 0
	data[1] = uint8(CmdWRITESTROBE)
	data[2] = uint8(address)
	binary.LittleEndian.PutUint32(data[3:7], value)

	if _, err := a.device.SendFeatureReport(data); err != nil {
		return fmt.Errorf("failed to write register: %w", err)
	}
	return nil
}

// SendCommand sends a command to the device
func (a *AIOCDevice) SendCommand(cmd Command) error {
	data := make([]byte, 8)
	data[0] = 0
	data[1] = uint8(cmd)
	data[2] = 0x00
	binary.LittleEndian.PutUint32(data[3:7], 0)

	if _, err := a.device.SendFeatureReport(data); err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}
	return nil
}

// SetPTTState sets the PTT state via raw HID write
func (a *AIOCDevice) SetPTTState(channel int, on bool) error {
	state := uint8(0)
	if on {
		state = 1
	}
	ioMask := uint8(1 << (channel - 1))
	ioData := state << (channel - 1)

	data := []byte{0, 0, ioData, ioMask, 0}
	n, err := a.device.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write PTT state: %w", err)
	}
	if n != len(data) {
		return fmt.Errorf("incomplete write: wrote %d bytes, expected %d", n, len(data))
	}
	return nil
}

// GetManufacturer returns the manufacturer string
func (a *AIOCDevice) GetManufacturer() (string, error) {
	return a.device.GetMfrStr()
}

// GetProduct returns the product string
func (a *AIOCDevice) GetProduct() (string, error) {
	return a.device.GetProductStr()
}

// GetSerialNumber returns the serial number string
func (a *AIOCDevice) GetSerialNumber() (string, error) {
	return a.device.GetSerialNbr()
}

// DumpRegisters dumps all known registers
func (a *AIOCDevice) DumpRegisters() error {
	registers := map[string]Register{
		"MAGIC":         RegMAGIC,
		"USBID":         RegUSBID,
		"AIOC_IOMUX0":   RegAIOCIOMUX0,
		"AIOC_IOMUX1":   RegAIOCIOMUX1,
		"CM108_IOMUX0":  RegCM108IOMUX0,
		"CM108_IOMUX1":  RegCM108IOMUX1,
		"CM108_IOMUX2":  RegCM108IOMUX2,
		"CM108_IOMUX3":  RegCM108IOMUX3,
		"SERIAL_CTRL":   RegSERIALCTRL,
		"SERIAL_IOMUX0": RegSERIALIOMUX0,
		"SERIAL_IOMUX1": RegSERIALIOMUX1,
		"SERIAL_IOMUX2": RegSERIALIOMUX2,
		"SERIAL_IOMUX3": RegSERIALIOMUX3,
		"AUDIO_RX":      RegAUDIORX,
		"AUDIO_TX":      RegAUDIOTX,
		"VPTT_LVLCTRL":  RegVPTTLVLCTRL,
		"VPTT_TIMCTRL":  RegVPTTTIMCTRL,
		"VCOS_LVLCTRL":  RegVCOSLVLCTRL,
		"VCOS_TIMCTRL":  RegVCOSTIMCTRL,
		"FOXHUNT_CTRL":  RegFOXHUNTCTRL,
		"FOXHUNT_MSG0":  RegFOXHUNTMSG0,
		"FOXHUNT_MSG1":  RegFOXHUNTMSG1,
		"FOXHUNT_MSG2":  RegFOXHUNTMSG2,
		"FOXHUNT_MSG3":  RegFOXHUNTMSG3,
	}

	for name, reg := range registers {
		value, err := a.Read(reg)
		if err != nil {
			return err
		}
		fmt.Printf("Reg. %s: %08x\n", name, value)
	}
	return nil
}
