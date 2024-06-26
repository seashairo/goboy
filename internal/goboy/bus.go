package goboy

type Bus struct {
}

func NewBus() Bus {
	return Bus{}
}

func (bus *Bus) read(address uint16) uint8 {
	return 0
}

func (bus *Bus) write(address uint16, value uint8) {

}
