package goboy

func BytesToUint16(hi byte, lo byte) uint16 {
	return (uint16(hi) << 8) | uint16(lo)
}

func Uint16ToBytes(value uint16) (byte, byte) {
	hi := byte(value >> 8)
	lo := byte(value & 0x00FF)

	return hi, lo
}

type MemoryBusser interface {
	readByte(address uint16) byte
	writeByte(address uint16, value byte)
}
