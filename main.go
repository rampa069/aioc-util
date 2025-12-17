package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sstallion/go-hid"
)

type Config struct {
	Defaults             bool
	Reboot               bool
	Dump                 bool
	SwapPTT              bool
	AutoPTT1             bool
	PTT1                 string
	PTT2                 string
	ListPTTSources       bool
	SetUSBVID            int
	SetUSBPID            int
	OpenUSBVID           int
	OpenUSBPID           int
	VolUp                string
	VolDn                string
	VPTTLvlCtrl          int
	VPTTTimCtrl          int
	VCOSLvlCtrl          int
	VCOSTimCtrl          int
	Store                bool
	SetPTT1State         string
	SetPTT2State         string
	EnableHWCOS          bool
	EnableVCOS           bool
	FoxhuntVolume        int
	FoxhuntWPM           int
	FoxhuntInterval      int
	FoxhuntGetSettings   bool
	FoxhuntMessage       string
	FoxhuntGetMessage    bool
	AudioRXGain          string
	AudioTXBoost         string
	AudioGetSettings     bool
}

func parsePTTSource(val string) (PTTSource, error) {
	if val == "" {
		return 0, nil
	}
	parts := strings.Split(val, "|")
	var result PTTSource
	for _, p := range parts {
		p = strings.TrimSpace(p)
		switch p {
		case "NONE":
			result |= PTTSourceNONE
		case "CM108GPIO1":
			result |= PTTSourceCM108GPIO1
		case "CM108GPIO2":
			result |= PTTSourceCM108GPIO2
		case "CM108GPIO3":
			result |= PTTSourceCM108GPIO3
		case "CM108GPIO4":
			result |= PTTSourceCM108GPIO4
		case "SERIALDTR":
			result |= PTTSourceSERIALDTR
		case "SERIALRTS":
			result |= PTTSourceSERIALRTS
		case "SERIALDTRNRTS":
			result |= PTTSourceSERIALDTRNRTS
		case "SERIALNDTRRTS":
			result |= PTTSourceSERIALNDTRRTS
		case "VPTT":
			result |= PTTSourceVPTT
		default:
			return 0, fmt.Errorf("unknown PTT source: %s", p)
		}
	}
	return result, nil
}

func parseCM108ButtonSource(val string) (CM108ButtonSource, error) {
	if val == "" {
		return 0, nil
	}
	parts := strings.Split(val, "|")
	var result CM108ButtonSource
	for _, p := range parts {
		p = strings.TrimSpace(p)
		switch p {
		case "NONE":
			result |= CM108ButtonSourceNONE
		case "IN1":
			result |= CM108ButtonSourceIN1
		case "IN2":
			result |= CM108ButtonSourceIN2
		case "VCOS":
			result |= CM108ButtonSourceVCOS
		default:
			return 0, fmt.Errorf("unknown button source: %s", p)
		}
	}
	return result, nil
}

func pttSourceString(src PTTSource) string {
	if src == PTTSourceNONE {
		return "NONE"
	}
	var parts []string
	if src&PTTSourceCM108GPIO1 != 0 {
		parts = append(parts, "CM108GPIO1")
	}
	if src&PTTSourceCM108GPIO2 != 0 {
		parts = append(parts, "CM108GPIO2")
	}
	if src&PTTSourceCM108GPIO3 != 0 {
		parts = append(parts, "CM108GPIO3")
	}
	if src&PTTSourceCM108GPIO4 != 0 {
		parts = append(parts, "CM108GPIO4")
	}
	if src&PTTSourceSERIALDTR != 0 {
		parts = append(parts, "SERIALDTR")
	}
	if src&PTTSourceSERIALRTS != 0 {
		parts = append(parts, "SERIALRTS")
	}
	if src&PTTSourceSERIALDTRNRTS != 0 {
		parts = append(parts, "SERIALDTRNRTS")
	}
	if src&PTTSourceSERIALNDTRRTS != 0 {
		parts = append(parts, "SERIALNDTRRTS")
	}
	if src&PTTSourceVPTT != 0 {
		parts = append(parts, "VPTT")
	}
	if len(parts) == 0 {
		return fmt.Sprintf("0x%08x", src)
	}
	return strings.Join(parts, "|")
}

func cm108ButtonSourceString(src CM108ButtonSource) string {
	if src == CM108ButtonSourceNONE {
		return "NONE"
	}
	var parts []string
	if src&CM108ButtonSourceIN1 != 0 {
		parts = append(parts, "IN1")
	}
	if src&CM108ButtonSourceIN2 != 0 {
		parts = append(parts, "IN2")
	}
	if src&CM108ButtonSourceVCOS != 0 {
		parts = append(parts, "VCOS")
	}
	if len(parts) == 0 {
		return fmt.Sprintf("0x%08x", src)
	}
	return strings.Join(parts, "|")
}

func parseHexOrDec(s string) (int, error) {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		val, err := strconv.ParseInt(s[2:], 16, 64)
		return int(val), err
	}
	val, err := strconv.ParseInt(s, 0, 64)
	return int(val), err
}

func main() {
	if err := hid.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize HID library: %v\n", err)
		os.Exit(1)
	}
	defer hid.Exit()

	config := Config{
		VPTTLvlCtrl:   -1,
		VPTTTimCtrl:   -1,
		VCOSLvlCtrl:   -1,
		VCOSTimCtrl:   -1,
		FoxhuntVolume: -1,
		FoxhuntWPM:    -1,
		FoxhuntInterval: -1,
		SetUSBVID:     -1,
		SetUSBPID:     -1,
		OpenUSBVID:    -1,
		OpenUSBPID:    -1,
	}

	flag.BoolVar(&config.Defaults, "defaults", false, "Load hardware defaults")
	flag.BoolVar(&config.Reboot, "reboot", false, "Reboot the device")
	flag.BoolVar(&config.Dump, "dump", false, "Dump all known registers")
	flag.BoolVar(&config.SwapPTT, "swap-ptt", false, "Swap PTT1/PTT2 sources")
	flag.BoolVar(&config.AutoPTT1, "auto-ptt1", false, "Set AutoPTT on PTT1")
	flag.StringVar(&config.PTT1, "ptt1", "", "Set arbitrary PTT1 source (e.g. \"CM108GPIO1|SERIALDTR\")")
	flag.StringVar(&config.PTT2, "ptt2", "", "Set arbitrary PTT2 source (e.g. \"CM108GPIO2|VPTT\")")
	flag.BoolVar(&config.ListPTTSources, "list-ptt-sources", false, "List all possible PTT sources")

	var setUSB string
	flag.StringVar(&setUSB, "set-usb", "", "Set USB VID and PID (format: VID,PID in hex or decimal)")

	var openUSB string
	flag.StringVar(&openUSB, "open-usb", "", "USB VID and PID to use when opening (format: VID,PID)")

	flag.StringVar(&config.VolUp, "vol-up", "", "Set Volume Up button source")
	flag.StringVar(&config.VolDn, "vol-dn", "", "Set Volume Down button source")

	var vpttLvlCtrl string
	flag.StringVar(&vpttLvlCtrl, "vptt-lvlctrl", "", "Set VPTT_LVLCTRL register (hex or decimal)")

	var vpttTimCtrl string
	flag.StringVar(&vpttTimCtrl, "vptt-timctrl", "", "Set VPTT_TIMCTRL register (hex or decimal)")

	var vcosLvlCtrl string
	flag.StringVar(&vcosLvlCtrl, "vcos-lvlctrl", "", "Set VCOS_LVLCTRL register (hex or decimal)")

	var vcosTimCtrl string
	flag.StringVar(&vcosTimCtrl, "vcos-timctrl", "", "Set VCOS_TIMCTRL register (hex or decimal)")

	flag.BoolVar(&config.Store, "store", false, "Store settings into flash")
	flag.StringVar(&config.SetPTT1State, "set-ptt1-state", "", "Set PTT1 state via raw HID write: 'on' or 'off'")
	flag.StringVar(&config.SetPTT2State, "set-ptt2-state", "", "Set PTT2 state via raw HID write: 'on' or 'off'")
	flag.BoolVar(&config.EnableHWCOS, "enable-hwcos", false, "Enable hardware COS (needs an AIOC that supports it)")
	flag.BoolVar(&config.EnableVCOS, "enable-vcos", false, "Enable virtual COS (default behavior)")

	var foxhuntVolume string
	flag.StringVar(&foxhuntVolume, "foxhunt-volume", "", "Set foxhunt volume (0-65535)")

	var foxhuntWPM string
	flag.StringVar(&foxhuntWPM, "foxhunt-wpm", "", "Set foxhunt words per minute (0-255)")

	var foxhuntInterval string
	flag.StringVar(&foxhuntInterval, "foxhunt-interval", "", "Set foxhunt interval in seconds (0-255, 0 disables foxhunt mode)")

	flag.BoolVar(&config.FoxhuntGetSettings, "foxhunt-get-settings", false, "Read and display current foxhunt control settings")
	flag.StringVar(&config.FoxhuntMessage, "foxhunt-message", "", "Set foxhunt message (up to 16 characters)")
	flag.BoolVar(&config.FoxhuntGetMessage, "foxhunt-get-message", false, "Read and display current foxhunt message")
	flag.StringVar(&config.AudioRXGain, "audio-rx-gain", "", "Set audio RX gain: 1x, 2x, 4x, 8x, or 16x")
	flag.StringVar(&config.AudioTXBoost, "audio-tx-boost", "", "Set audio TX boost: off or on")
	flag.BoolVar(&config.AudioGetSettings, "audio-get-settings", false, "Read and display current audio settings")

	flag.Parse()

	// Parse hex/decimal values
	if setUSB != "" {
		parts := strings.Split(setUSB, ",")
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "Invalid --set-usb format. Use: VID,PID\n")
			os.Exit(1)
		}
		vid, err := parseHexOrDec(parts[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid VID: %v\n", err)
			os.Exit(1)
		}
		pid, err := parseHexOrDec(parts[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid PID: %v\n", err)
			os.Exit(1)
		}
		config.SetUSBVID = vid
		config.SetUSBPID = pid
	}

	if openUSB != "" {
		parts := strings.Split(openUSB, ",")
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "Invalid --open-usb format. Use: VID,PID\n")
			os.Exit(1)
		}
		vid, err := parseHexOrDec(parts[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid VID: %v\n", err)
			os.Exit(1)
		}
		pid, err := parseHexOrDec(parts[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid PID: %v\n", err)
			os.Exit(1)
		}
		config.OpenUSBVID = vid
		config.OpenUSBPID = pid
	}

	if vpttLvlCtrl != "" {
		val, err := parseHexOrDec(vpttLvlCtrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid --vptt-lvlctrl value: %v\n", err)
			os.Exit(1)
		}
		config.VPTTLvlCtrl = val
	}

	if vpttTimCtrl != "" {
		val, err := parseHexOrDec(vpttTimCtrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid --vptt-timctrl value: %v\n", err)
			os.Exit(1)
		}
		config.VPTTTimCtrl = val
	}

	if vcosLvlCtrl != "" {
		val, err := parseHexOrDec(vcosLvlCtrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid --vcos-lvlctrl value: %v\n", err)
			os.Exit(1)
		}
		config.VCOSLvlCtrl = val
	}

	if vcosTimCtrl != "" {
		val, err := parseHexOrDec(vcosTimCtrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid --vcos-timctrl value: %v\n", err)
			os.Exit(1)
		}
		config.VCOSTimCtrl = val
	}

	if foxhuntVolume != "" {
		val, err := parseHexOrDec(foxhuntVolume)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid --foxhunt-volume value: %v\n", err)
			os.Exit(1)
		}
		config.FoxhuntVolume = val
	}

	if foxhuntWPM != "" {
		val, err := parseHexOrDec(foxhuntWPM)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid --foxhunt-wpm value: %v\n", err)
			os.Exit(1)
		}
		config.FoxhuntWPM = val
	}

	if foxhuntInterval != "" {
		val, err := parseHexOrDec(foxhuntInterval)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid --foxhunt-interval value: %v\n", err)
			os.Exit(1)
		}
		config.FoxhuntInterval = val
	}

	// Show help if no args
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Handle list-ptt-sources
	if config.ListPTTSources {
		fmt.Println("CM108GPIO1 (0x00000001)")
		fmt.Println("CM108GPIO2 (0x00000002)")
		fmt.Println("CM108GPIO3 (0x00000004)")
		fmt.Println("CM108GPIO4 (0x00000008)")
		fmt.Println("SERIALDTR (0x00000100)")
		fmt.Println("SERIALRTS (0x00000200)")
		fmt.Println("SERIALDTRNRTS (0x00000400)")
		fmt.Println("SERIALNDTRRTS (0x00000800)")
		fmt.Println("VPTT (0x00001000)")
		os.Exit(0)
	}

	// Open device
	vid := uint16(AIOCVendorID)
	pid := uint16(AIOCProductID)
	if config.OpenUSBVID != -1 {
		vid = uint16(config.OpenUSBVID)
	}
	if config.OpenUSBPID != -1 {
		pid = uint16(config.OpenUSBPID)
	}

	aioc, err := Open(vid, pid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open AIOC device (VID: 0x%04x, PID: 0x%04x): %v\n", vid, pid, err)
		os.Exit(1)
	}
	defer aioc.Close()

	// Execute commands
	if config.Defaults {
		fmt.Println("Loading Defaults...")
		if err := aioc.SendCommand(CmdDEFAULTS); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load defaults: %v\n", err)
			os.Exit(1)
		}
	}

	if config.Dump {
		mfr, _ := aioc.GetManufacturer()
		prod, _ := aioc.GetProduct()
		serial, _ := aioc.GetSerialNumber()

		fmt.Printf("Manufacturer: %s\n", mfr)
		fmt.Printf("Product: %s\n", prod)
		fmt.Printf("Serial No: %s\n", serial)

		magic, _ := aioc.Read(RegMAGIC)
		magicBytes := []byte{
			byte(magic),
			byte(magic >> 8),
			byte(magic >> 16),
			byte(magic >> 24),
		}
		fmt.Printf("Magic: %s\n", magicBytes)

		ptt1Source, _ := aioc.Read(RegAIOCIOMUX0)
		ptt2Source, _ := aioc.Read(RegAIOCIOMUX1)
		fmt.Printf("Current PTT1 Source: %s\n", pttSourceString(PTTSource(ptt1Source)))
		fmt.Printf("Current PTT2 Source: %s\n", pttSourceString(PTTSource(ptt2Source)))

		btn1Source, _ := aioc.Read(RegCM108IOMUX0)
		btn2Source, _ := aioc.Read(RegCM108IOMUX1)
		btn3Source, _ := aioc.Read(RegCM108IOMUX2)
		btn4Source, _ := aioc.Read(RegCM108IOMUX3)
		fmt.Printf("Current CM108 Button 1 (VolUP) Source: %s\n", cm108ButtonSourceString(CM108ButtonSource(btn1Source)))
		fmt.Printf("Current CM108 Button 2 (VolDN) Source: %s\n", cm108ButtonSourceString(CM108ButtonSource(btn2Source)))
		fmt.Printf("Current CM108 Button 3 (PlbMute) Source: %s\n", cm108ButtonSourceString(CM108ButtonSource(btn3Source)))
		fmt.Printf("Current CM108 Button 4 (RecMute) Source: %s\n", cm108ButtonSourceString(CM108ButtonSource(btn4Source)))

		if err := aioc.DumpRegisters(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to dump registers: %v\n", err)
			os.Exit(1)
		}
	}

	if config.SwapPTT {
		ptt1Source, _ := aioc.Read(RegAIOCIOMUX0)
		ptt2Source, _ := aioc.Read(RegAIOCIOMUX1)

		fmt.Printf("Setting PTT1 Source to %s\n", pttSourceString(PTTSource(ptt2Source)))
		aioc.Write(RegAIOCIOMUX0, ptt2Source)
		fmt.Printf("Setting PTT2 Source to %s\n", pttSourceString(PTTSource(ptt1Source)))
		aioc.Write(RegAIOCIOMUX1, ptt1Source)

		newPTT1, _ := aioc.Read(RegAIOCIOMUX0)
		newPTT2, _ := aioc.Read(RegAIOCIOMUX1)
		fmt.Printf("Now PTT1 Source: %s\n", pttSourceString(PTTSource(newPTT1)))
		fmt.Printf("Now PTT2 Source: %s\n", pttSourceString(PTTSource(newPTT2)))
	}

	if config.AutoPTT1 {
		fmt.Printf("Setting PTT1 Source to %s\n", pttSourceString(PTTSourceVPTT))
		aioc.Write(RegAIOCIOMUX0, uint32(PTTSourceVPTT))

		newPTT1, _ := aioc.Read(RegAIOCIOMUX0)
		newPTT2, _ := aioc.Read(RegAIOCIOMUX1)
		fmt.Printf("Now PTT1 Source: %s\n", pttSourceString(PTTSource(newPTT1)))
		fmt.Printf("Now PTT2 Source: %s\n", pttSourceString(PTTSource(newPTT2)))
	}

	if config.PTT1 != "" || config.PTT2 != "" {
		if config.PTT1 != "" {
			val, err := parsePTTSource(config.PTT1)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse PTT1 source: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Setting PTT1 Source to %s\n", pttSourceString(val))
			aioc.Write(RegAIOCIOMUX0, uint32(val))
		}
		if config.PTT2 != "" {
			val, err := parsePTTSource(config.PTT2)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse PTT2 source: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Setting PTT2 Source to %s\n", pttSourceString(val))
			aioc.Write(RegAIOCIOMUX1, uint32(val))
		}

		newPTT1, _ := aioc.Read(RegAIOCIOMUX0)
		newPTT2, _ := aioc.Read(RegAIOCIOMUX1)
		fmt.Printf("Now PTT1 Source: %s\n", pttSourceString(PTTSource(newPTT1)))
		fmt.Printf("Now PTT2 Source: %s\n", pttSourceString(PTTSource(newPTT2)))
	}

	if config.SetUSBVID != -1 && config.SetUSBPID != -1 {
		value := uint32((config.SetUSBPID << 16) | config.SetUSBVID)
		aioc.Write(RegUSBID, value)
		newVal, _ := aioc.Read(RegUSBID)
		fmt.Printf("Now USBID: %08x\n", newVal)
	}

	if config.VolUp != "" || config.VolDn != "" {
		if config.VolUp != "" {
			su, err := parseCM108ButtonSource(config.VolUp)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse VolUp source: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Setting VolUP button source to %s\n", cm108ButtonSourceString(su))
			aioc.Write(RegCM108IOMUX0, uint32(su))
		}
		if config.VolDn != "" {
			sd, err := parseCM108ButtonSource(config.VolDn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse VolDn source: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Setting VolDN button source to %s\n", cm108ButtonSourceString(sd))
			aioc.Write(RegCM108IOMUX1, uint32(sd))
		}

		newVolUp, _ := aioc.Read(RegCM108IOMUX0)
		newVolDn, _ := aioc.Read(RegCM108IOMUX1)
		fmt.Printf("Now VolUP button source: %s\n", cm108ButtonSourceString(CM108ButtonSource(newVolUp)))
		fmt.Printf("Now VolDN button source: %s\n", cm108ButtonSourceString(CM108ButtonSource(newVolDn)))
	}

	if config.VPTTLvlCtrl != -1 {
		fmt.Printf("Setting VPTT_LVLCTRL to 0x%x\n", config.VPTTLvlCtrl)
		aioc.Write(RegVPTTLVLCTRL, uint32(config.VPTTLvlCtrl))
		newVal, _ := aioc.Read(RegVPTTLVLCTRL)
		fmt.Printf("Now VPTT_LVLCTRL: %08x\n", newVal)
	}

	if config.VPTTTimCtrl != -1 {
		fmt.Printf("Setting VPTT_TIMCTRL to 0x%x\n", config.VPTTTimCtrl)
		aioc.Write(RegVPTTTIMCTRL, uint32(config.VPTTTimCtrl))
		newVal, _ := aioc.Read(RegVPTTTIMCTRL)
		fmt.Printf("Now VPTT_TIMCTRL: %08x\n", newVal)
	}

	if config.VCOSLvlCtrl != -1 {
		fmt.Printf("Setting VCOS_LVLCTRL to 0x%x\n", config.VCOSLvlCtrl)
		aioc.Write(RegVCOSLVLCTRL, uint32(config.VCOSLvlCtrl))
		newVal, _ := aioc.Read(RegVCOSLVLCTRL)
		fmt.Printf("Now VCOS_LVLCTRL: %08x\n", newVal)
	}

	if config.VCOSTimCtrl != -1 {
		fmt.Printf("Setting VCOS_TIMCTRL to 0x%x\n", config.VCOSTimCtrl)
		aioc.Write(RegVCOSTIMCTRL, uint32(config.VCOSTimCtrl))
		newVal, _ := aioc.Read(RegVCOSTIMCTRL)
		fmt.Printf("Now VCOS_TIMCTRL: %08x\n", newVal)
	}

	if config.EnableHWCOS {
		fmt.Println("Enabling hardware COS (if your aioc supports it)...")
		aioc.Write(RegCM108IOMUX0, uint32(CM108ButtonSourceNONE))
		aioc.Write(RegCM108IOMUX1, uint32(CM108ButtonSourceIN2))

		newVal0, _ := aioc.Read(RegCM108IOMUX0)
		newVal1, _ := aioc.Read(RegCM108IOMUX1)
		fmt.Printf("Now CM108_IOMUX0: %s\n", cm108ButtonSourceString(CM108ButtonSource(newVal0)))
		fmt.Printf("Now CM108_IOMUX1: %s\n", cm108ButtonSourceString(CM108ButtonSource(newVal1)))
	}

	if config.EnableVCOS {
		fmt.Println("Enabling virtual COS...")
		aioc.Write(RegCM108IOMUX0, uint32(CM108ButtonSourceIN2))
		aioc.Write(RegCM108IOMUX1, uint32(CM108ButtonSourceVCOS))

		newVal0, _ := aioc.Read(RegCM108IOMUX0)
		newVal1, _ := aioc.Read(RegCM108IOMUX1)
		fmt.Printf("Now CM108_IOMUX0: %s\n", cm108ButtonSourceString(CM108ButtonSource(newVal0)))
		fmt.Printf("Now CM108_IOMUX1: %s\n", cm108ButtonSourceString(CM108ButtonSource(newVal1)))
	}

	if config.FoxhuntGetSettings {
		currentFoxhunt, _ := aioc.Read(RegFOXHUNTCTRL)
		currentVolume := (currentFoxhunt >> 16) & 0xFFFF
		currentWPM := (currentFoxhunt >> 8) & 0xFF
		currentInterval := currentFoxhunt & 0xFF

		fmt.Println("Current foxhunt settings:")
		fmt.Printf("  Volume: %d\n", currentVolume)
		fmt.Printf("  WPM: %d\n", currentWPM)
		fmt.Printf("  Interval: %d seconds\n", currentInterval)
		fmt.Printf("  Raw register: %08x\n", currentFoxhunt)
	}

	if config.FoxhuntGetMessage {
		msgRegisters := []Register{RegFOXHUNTMSG0, RegFOXHUNTMSG1, RegFOXHUNTMSG2, RegFOXHUNTMSG3}
		messageBytes := make([]byte, 0, 16)

		fmt.Println("Current foxhunt message registers:")
		for i, reg := range msgRegisters {
			val, _ := aioc.Read(reg)
			regBytes := []byte{
				byte(val),
				byte(val >> 8),
				byte(val >> 16),
				byte(val >> 24),
			}
			messageBytes = append(messageBytes, regBytes...)
			fmt.Printf("  MSG%d: %08x ('%s')\n", i, val, string(regBytes))
		}

		// Find null terminator
		nullIdx := -1
		for i, b := range messageBytes {
			if b == 0 {
				nullIdx = i
				break
			}
		}

		messageStr := ""
		if nullIdx != -1 {
			messageStr = string(messageBytes[:nullIdx])
		} else {
			messageStr = string(messageBytes)
		}

		fmt.Printf("Current foxhunt message: '%s'\n", messageStr)
	}

	if config.FoxhuntVolume != -1 || config.FoxhuntWPM != -1 || config.FoxhuntInterval != -1 {
		currentFoxhunt, _ := aioc.Read(RegFOXHUNTCTRL)
		currentVolume := int((currentFoxhunt >> 16) & 0xFFFF)
		currentWPM := int((currentFoxhunt >> 8) & 0xFF)
		currentInterval := int(currentFoxhunt & 0xFF)

		newVolume := currentVolume
		if config.FoxhuntVolume != -1 {
			newVolume = config.FoxhuntVolume
		}
		newWPM := currentWPM
		if config.FoxhuntWPM != -1 {
			newWPM = config.FoxhuntWPM
		}
		newInterval := currentInterval
		if config.FoxhuntInterval != -1 {
			newInterval = config.FoxhuntInterval
		}

		newFoxhunt := uint32((newVolume << 16) | (newWPM << 8) | newInterval)
		fmt.Printf("Setting FOXHUNT_CTRL: volume=%d, wpm=%d, interval=%d\n", newVolume, newWPM, newInterval)
		aioc.Write(RegFOXHUNTCTRL, newFoxhunt)

		updatedVal, _ := aioc.Read(RegFOXHUNTCTRL)
		fmt.Printf("Now FOXHUNT_CTRL: %08x\n", updatedVal)
	}

	if config.FoxhuntMessage != "" {
		messageBytes := []byte(config.FoxhuntMessage)
		if len(messageBytes) > 16 {
			messageBytes = messageBytes[:16]
		}
		// Pad with nulls
		for len(messageBytes) < 16 {
			messageBytes = append(messageBytes, 0)
		}

		msgRegisters := []Register{RegFOXHUNTMSG0, RegFOXHUNTMSG1, RegFOXHUNTMSG2, RegFOXHUNTMSG3}

		fmt.Printf("Setting foxhunt message: '%s'\n", config.FoxhuntMessage)
		for i := 0; i < 4; i++ {
			byteOffset := i * 4
			val := uint32(messageBytes[byteOffset]) |
				(uint32(messageBytes[byteOffset+1]) << 8) |
				(uint32(messageBytes[byteOffset+2]) << 16) |
				(uint32(messageBytes[byteOffset+3]) << 24)
			aioc.Write(msgRegisters[i], val)
			fmt.Printf("  MSG%d: %08x ('%s')\n", i, val, string(messageBytes[byteOffset:byteOffset+4]))
		}
	}

	if config.AudioGetSettings {
		currentRX, _ := aioc.Read(RegAUDIORX)
		currentTX, _ := aioc.Read(RegAUDIOTX)

		rxGainName := "unknown"
		switch RXGain(currentRX) {
		case RXGain1X:
			rxGainName = "1x"
		case RXGain2X:
			rxGainName = "2x"
		case RXGain4X:
			rxGainName = "4x"
		case RXGain8X:
			rxGainName = "8x"
		case RXGain16X:
			rxGainName = "16x"
		}

		txBoostName := "unknown"
		switch TXBoost(currentTX) {
		case TXBoostOFF:
			txBoostName = "off"
		case TXBoostON:
			txBoostName = "on"
		}

		fmt.Println("Current audio settings:")
		fmt.Printf("  RX Gain: %s\n", rxGainName)
		fmt.Printf("  TX Boost: %s\n", txBoostName)
		fmt.Printf("  Raw AUDIO_RX: %08x\n", currentRX)
		fmt.Printf("  Raw AUDIO_TX: %08x\n", currentTX)
	}

	if config.AudioRXGain != "" {
		gainMap := map[string]RXGain{
			"1x":  RXGain1X,
			"2x":  RXGain2X,
			"4x":  RXGain4X,
			"8x":  RXGain8X,
			"16x": RXGain16X,
		}

		gain, ok := gainMap[config.AudioRXGain]
		if !ok {
			fmt.Fprintf(os.Stderr, "Invalid audio RX gain: %s\n", config.AudioRXGain)
			os.Exit(1)
		}

		fmt.Printf("Setting Audio RX gain to %s\n", config.AudioRXGain)
		aioc.Write(RegAUDIORX, uint32(gain))
		newVal, _ := aioc.Read(RegAUDIORX)
		fmt.Printf("Now AUDIO_RX: %08x\n", newVal)
	}

	if config.AudioTXBoost != "" {
		boostMap := map[string]TXBoost{
			"off": TXBoostOFF,
			"on":  TXBoostON,
		}

		boost, ok := boostMap[config.AudioTXBoost]
		if !ok {
			fmt.Fprintf(os.Stderr, "Invalid audio TX boost: %s\n", config.AudioTXBoost)
			os.Exit(1)
		}

		fmt.Printf("Setting Audio TX boost to %s\n", config.AudioTXBoost)
		aioc.Write(RegAUDIOTX, uint32(boost))
		newVal, _ := aioc.Read(RegAUDIOTX)
		fmt.Printf("Now AUDIO_TX: %08x\n", newVal)
	}

	if config.Store {
		fmt.Println("Storing...")
		aioc.SendCommand(CmdSTORE)
	}

	if config.SetPTT1State != "" {
		on := config.SetPTT1State == "on"
		if err := aioc.SetPTTState(PTTChannel1, on); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to set PTT1 state: %v\n", err)
			os.Exit(1)
		}
	}

	if config.SetPTT2State != "" {
		on := config.SetPTT2State == "on"
		if err := aioc.SetPTTState(PTTChannel2, on); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to set PTT2 state: %v\n", err)
			os.Exit(1)
		}
	}

	if config.Reboot {
		fmt.Println("Rebooting device...")
		aioc.SendCommand(CmdREBOOT)
	}
}
