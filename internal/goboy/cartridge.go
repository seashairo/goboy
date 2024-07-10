package goboy

import (
	"encoding/binary"
	"fmt"
	"os"
)

var CARTRIDGE_TYPE_MAP = map[byte]string{
	0x00: "ROM ONLY",
	0x01: "MBC1",
	0x02: "MBC1+RAM",
	0x03: "MBC1+RAM+BATTERY",
	0x05: "MBC2",
	0x06: "MBC2+BATTERY",
	0x08: "ROM+RAM",
	0x09: "ROM+RAM+BATTERY",
	0x0B: "MMM01",
	0x0C: "MMM01+RAM",
	0x0D: "MMM01+RAM+BATTERY",
	0x0F: "MBC3+TIMER+BATTERY",
	0x10: "MBC3+TIMER+RAM+BATTERY",
	0x11: "MBC3",
	0x12: "MBC3+RAM",
	0x13: "MBC3+RAM+BATTERY",
	0x19: "MBC5",
	0x1A: "MBC5+RAM",
	0x1B: "MBC5+RAM+BATTERY",
	0x1C: "MBC5+RUMBLE",
	0x1D: "MBC5+RUMBLE+RAM",
	0x1E: "MBC5+RUMBLE+RAM+BATTERY",
	0x20: "MBC6",
	0x22: "MBC7+SENSOR+RUMBLE+RAM+BATTERY",
	0xFC: "POCKET CAMERA",
	0xFD: "BANDAI TAMA5",
	0xFE: "HuC3",
	0xFF: "HuC1+RAM+BATTERY",
}

var ROM_SIZE_MAP = map[byte]string{
	0x00: "32 KiB",
	0x01: "64 KiB",
	0x02: "128 KiB",
	0x03: "256 KiB",
	0x04: "512 KiB",
	0x05: "1 MiB",
	0x06: "2 MiB",
	0x07: "4 MiB",
	0x08: "8 MiB",
}

// RAM is split into 8 KiB banks
var RAM_SIZE_MAP = map[byte]string{
	0x00: "0 KiB",
	0x01: "2 KiB",
	0x02: "8 KiB",
	0x03: "32 KiB",
	0x04: "128 KiB",
	0x05: "64 KiB",
}

var DESTINATION_CODE_MAP = map[byte]string{
	0x00: "Japan",
	0x01: "Overseas",
}

var OLD_LICENSEE_CODE_MAP = map[byte]string{
	0x00: "None",
	0x01: "Nintendo",
	0x08: "Capcom",
	0x09: "HOT-B",
	0x0A: "Jaleco",
	0x0B: "Coconuts Japan",
	0x0C: "Elite Systems",
	0x13: "EA (Electronic Arts)",
	0x18: "Hudson Soft",
	0x19: "ITC Entertainment",
	0x1A: "Yanoman",
	0x1D: "Japan Clary",
	0x1F: "Virgin Games Ltd.3",
	0x24: "PCM Complete",
	0x25: "San-X",
	0x28: "Kemco",
	0x29: "SETA Corporation",
	0x30: "Infogrames5",
	0x31: "Nintendo",
	0x32: "Bandai",
	0x33: "Indicates that the New licensee code should be used instead.",
	0x34: "Konami",
	0x35: "HectorSoft",
	0x38: "Capcom",
	0x39: "Banpresto",
	0x3C: ".Entertainment i",
	0x3E: "Gremlin",
	0x41: "Ubi Soft1",
	0x42: "Atlus",
	0x44: "Malibu Interactive",
	0x46: "Angel",
	0x47: "Spectrum Holoby",
	0x49: "Irem",
	0x4A: "Virgin Games Ltd.3",
	0x4D: "Malibu Interactive",
	0x4F: "U.S. Gold",
	0x50: "Absolute",
	0x51: "Acclaim Entertainment",
	0x52: "Activision",
	0x53: "Sammy USA Corporation",
	0x54: "GameTek",
	0x55: "Park Place",
	0x56: "LJN",
	0x57: "Matchbox",
	0x59: "Milton Bradley Company",
	0x5A: "Mindscape",
	0x5B: "Romstar",
	0x5C: "Naxat Soft13",
	0x5D: "Tradewest",
	0x60: "Titus Interactive",
	0x61: "Virgin Games Ltd.3",
	0x67: "Ocean Software",
	0x69: "EA (Electronic Arts)",
	0x6E: "Elite Systems",
	0x6F: "Electro Brain",
	0x70: "Infogrames5",
	0x71: "Interplay Entertainment",
	0x72: "Broderbund",
	0x73: "Sculptured Software6",
	0x75: "The Sales Curve Limited7",
	0x78: "THQ",
	0x79: "Accolade",
	0x7A: "Triffix Entertainment",
	0x7C: "Microprose",
	0x7F: "Kemco",
	0x80: "Misawa Entertainment",
	0x83: "Lozc",
	0x86: "Tokuma Shoten",
	0x8B: "Bullet-Proof Software2",
	0x8C: "Vic Tokai",
	0x8E: "Ape",
	0x8F: "I’Max",
	0x91: "Chunsoft Co.8",
	0x92: "Video System",
	0x93: "Tsubaraya Productions",
	0x95: "Varie",
	0x96: "Yonezawa/S’Pal",
	0x97: "Kemco",
	0x99: "Arc",
	0x9A: "Nihon Bussan",
	0x9B: "Tecmo",
	0x9C: "Imagineer",
	0x9D: "Banpresto",
	0x9F: "Nova",
	0xA1: "Hori Electric",
	0xA2: "Bandai",
	0xA4: "Konami",
	0xA6: "Kawada",
	0xA7: "Takara",
	0xA9: "Technos Japan",
	0xAA: "Broderbund",
	0xAC: "Toei Animation",
	0xAD: "Toho",
	0xAF: "Namco",
	0xB0: "Acclaim Entertainment",
	0xB1: "ASCII Corporation or Nexsoft",
	0xB2: "Bandai",
	0xB4: "Square Enix",
	0xB6: "HAL Laboratory",
	0xB7: "SNK",
	0xB9: "Pony Canyon",
	0xBA: "Culture Brain",
	0xBB: "Sunsoft",
	0xBD: "Sony Imagesoft",
	0xBF: "Sammy Corporation",
	0xC0: "Taito",
	0xC2: "Kemco",
	0xC3: "Square",
	0xC4: "Tokuma Shoten",
	0xC5: "Data East",
	0xC6: "Tonkinhouse",
	0xC8: "Koei",
	0xC9: "UFL",
	0xCA: "Ultra",
	0xCB: "Vap",
	0xCC: "Use Corporation",
	0xCD: "Meldac",
	0xCE: "Pony Canyon",
	0xCF: "Angel",
	0xD0: "Taito",
	0xD1: "Sofel",
	0xD2: "Quest",
	0xD3: "Sigma Enterprises",
	0xD4: "ASK Kodansha Co.",
	0xD6: "Naxat Soft13",
	0xD7: "Copya System",
	0xD9: "Banpresto",
	0xDA: "Tomy",
	0xDB: "LJN",
	0xDD: "NCS",
	0xDE: "Human",
	0xDF: "Altron",
	0xE0: "Jaleco",
	0xE1: "Towa Chiki",
	0xE2: "Yutaka",
	0xE3: "Varie",
	0xE5: "Epcoh",
	0xE7: "Athena",
	0xE8: "Asmik Ace Entertainment",
	0xE9: "Natsume",
	0xEA: "King Records",
	0xEB: "Atlus",
	0xEC: "Epic/Sony Records",
	0xEE: "IGS",
	0xF0: "A Wave",
	0xF3: "Extreme Entertainment",
	0xFF: "LJN",
}

var NEW_LICENSEE_CODE_MAP = map[string]string{
	"00": "None",
	"01": "Nintendo Research & Development 1",
	"08": "Capcom",
	"13": "EA (Electronic Arts)",
	"18": "Hudson Soft",
	"19": "B-AI",
	"20": "KSS",
	"22": "Planning Office WADA",
	"24": "PCM Complete",
	"25": "San-X",
	"28": "Kemco",
	"29": "SETA Corporation",
	"30": "Viacom",
	"31": "Nintendo",
	"32": "Bandai",
	"33": "Ocean Software/Acclaim Entertainment",
	"34": "Konami",
	"35": "HectorSoft",
	"37": "Taito",
	"38": "Hudson Soft",
	"39": "Banpresto",
	"41": "Ubi Soft1",
	"42": "Atlus",
	"44": "Malibu Interactive",
	"46": "Angel",
	"47": "Bullet-Proof Software2",
	"49": "Irem",
	"50": "Absolute",
	"51": "Acclaim Entertainment",
	"52": "Activision",
	"53": "Sammy USA Corporation",
	"54": "Konami",
	"55": "Hi Tech Expressions",
	"56": "LJN",
	"57": "Matchbox",
	"58": "Mattel",
	"59": "Milton Bradley Company",
	"60": "Titus Interactive",
	"61": "Virgin Games Ltd.3",
	"64": "Lucasfilm Games4",
	"67": "Ocean Software",
	"69": "EA (Electronic Arts)",
	"70": "Infogrames5",
	"71": "Interplay Entertainment",
	"72": "Broderbund",
	"73": "Sculptured Software6",
	"75": "The Sales Curve Limited7",
	"78": "THQ",
	"79": "Accolade",
	"80": "Misawa Entertainment",
	"83": "lozc",
	"86": "Tokuma Shoten",
	"87": "Tsukuda Original",
	"91": "Chunsoft Co.8",
	"92": "Video System",
	"93": "Ocean Software/Acclaim Entertainment",
	"95": "Varie",
	"96": "Yonezawa/s’pal",
	"97": "Kaneko",
	"99": "Pack-In-Video",
	"A4": "Konami (Yu-Gi-Oh!)",
	"BL": "MTO",
	"DK": "Kodansha",
}

// @see https://gbdev.io/pandocs/The_Cartridge_Header.html
type RomHeader struct {
	entry [4]byte    // 0x0100 - 0x0103
	logo  [0x30]byte // 0x0104 - 0x0133

	title           [16]byte // 0x0134 - 0x0143
	newLicenseeCode [2]byte  // 0x0144 – 0x0145
	isSgb           byte     // 0x0146
	cartridgeType   byte     // 0x0147
	romSize         byte     // 0x0148
	ramSize         byte     // 0x0149
	destinationCode byte     // 0x014A
	licenseeCode    byte     // 0x014B
	version         byte     // 0x014C
	headerChecksum  byte     // 0x014D
	globalChecksum  uint16   // 0x014E - 0x14F
}

func (header *RomHeader) GetFriendlyTitle() string {
	return string(header.title[:])
}

func (header *RomHeader) GetFriendlyLicenseeName() string {
	if header.licenseeCode == 0x33 {
		return mapToFriendlyName(NEW_LICENSEE_CODE_MAP, string(header.newLicenseeCode[:]))
	} else {
		return mapToFriendlyName(OLD_LICENSEE_CODE_MAP, header.licenseeCode)
	}
}

func (header *RomHeader) GetFriendlyType() string {
	return mapToFriendlyName(CARTRIDGE_TYPE_MAP, header.cartridgeType)
}

func (header *RomHeader) GetFriendlyDestinationName() string {
	return mapToFriendlyName(DESTINATION_CODE_MAP, header.destinationCode)
}

func (header *RomHeader) GetFriendlyRomSize() string {
	return mapToFriendlyName(ROM_SIZE_MAP, header.romSize)
}

func (header *RomHeader) GetFriendlyRamSize() string {
	return mapToFriendlyName(RAM_SIZE_MAP, header.ramSize)
}

func mapToFriendlyName[K comparable](m map[K]string, k K) string {
	value, ok := m[k]

	if ok {
		return value
	} else {
		return "Unknown"
	}
}

type Cartridge struct {
	filename string
	romData  []byte
	header   RomHeader
}

func LoadCartridge(path string) Cartridge {
	fmt.Printf("Loading cartridge from %s\n", path)

	romData, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	header := RomHeader{}
	copy(header.entry[:], romData[0x100:0x104])
	copy(header.logo[:], romData[0x104:0x134])
	copy(header.title[:], romData[0x0134:0x0144])
	copy(header.newLicenseeCode[:], romData[0x0144:0x0146])
	header.isSgb = romData[0x0146]
	header.cartridgeType = romData[0x0147]
	header.romSize = romData[0x0148]
	header.ramSize = romData[0x0149]
	header.destinationCode = romData[0x014A]
	header.licenseeCode = romData[0x014B]
	header.version = romData[0x014C]
	header.headerChecksum = romData[0x14D]
	header.globalChecksum = binary.BigEndian.Uint16(
		[]byte{romData[0x14E], romData[0x14F]},
	)

	cartridge := Cartridge{
		filename: path,
		romData:  romData,
		header:   header,
	}

	cartridge.debugPrint()

	return cartridge
}

func (cartridge *Cartridge) debugPrint() {
	h := cartridge.header

	fmt.Println("Loaded cartridge:")
	fmt.Printf("\tTitle: %s\n", h.GetFriendlyTitle())
	fmt.Printf("\tLicensee: %s (0x%2.2x, %s)\n", h.GetFriendlyLicenseeName(), h.licenseeCode, string(h.newLicenseeCode[:]))
	fmt.Printf("\tCartridge Type: %s (0x%2.2x)\n", h.GetFriendlyType(), h.cartridgeType)
	fmt.Printf("\tROM Size: %s (0x%2.2x)\n", h.GetFriendlyRomSize(), h.romSize)
	fmt.Printf("\tRAM Size: %s (0x%2.2x)\n", h.GetFriendlyRamSize(), h.ramSize)
	fmt.Printf("\tDestination: %s (0x%2.2x)\n", h.GetFriendlyDestinationName(), h.destinationCode)
	fmt.Printf("\tVersion: 0x%2.2x\n", h.version)

	headerChecksum := cartridge.calculateHeaderChecksum()
	globalChecksum := cartridge.calculateGlobalChecksum()

	var headerChecksumString = "PASS"
	if h.headerChecksum != headerChecksum {
		headerChecksumString = fmt.Sprintf("FAIL - got 0x%2.2x", headerChecksum)
	}

	fmt.Printf("\tHeader Checksum: 0x%2.2x (%s)\n", h.headerChecksum, headerChecksumString)

	var globalChecksumString = "PASS"
	if h.globalChecksum != globalChecksum {
		globalChecksumString = fmt.Sprintf("FAIL - got 0x%2.2x", globalChecksum)
	}

	fmt.Printf("\tGlobal Checksum: 0x%2.2x (%s)\n", h.globalChecksum, globalChecksumString)
}

func (cartridge *Cartridge) calculateHeaderChecksum() byte {
	headerChecksum := byte(0)
	for address := 0x0134; address <= 0x014C; address++ {
		headerChecksum = headerChecksum - cartridge.romData[address] - 1
	}
	return headerChecksum
}

func (cartridge *Cartridge) calculateGlobalChecksum() uint16 {
	globalChecksum := uint16(0)
	for address := 0x0000; address < 0x014E; address++ {
		globalChecksum = globalChecksum + uint16(cartridge.romData[address])
	}

	for address := 0x0150; address < len(cartridge.romData); address++ {
		globalChecksum = globalChecksum + uint16(cartridge.romData[address])
	}

	return globalChecksum
}

func (cartridge *Cartridge) read(address uint16) byte {
	return cartridge.romData[address]
}

func (cartridge *Cartridge) write(address uint16, value byte) {
	panic("Can't write to ROM")
}
