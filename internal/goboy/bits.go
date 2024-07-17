package goboy

func SetBit(b byte, position byte, on bool) byte {
	if on {
		return b | (1 << position)
	} else {
		return b &^ (1 << position)
	}
}

func GetBit(b byte, position byte) bool {
	return (b & (1 << position)) != 0
}

func BoolToByte(val bool) byte {
	if val {
		return 1
	}

	return 0
}

func BoolToUint16(val bool) uint16 {
	if val {
		return 1
	}

	return 0
}
