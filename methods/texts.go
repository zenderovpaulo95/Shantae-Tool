package methods

type texts struct {
	CRC         uint
	TextOffsets []uint
	Texts       []string
}

type textHeader struct {
	Header         uint
	Unknown1       uint
	Unknown2       byte
	Unknown3       byte
	CountLocTexts1 uint16
	Offset1        uint
	Offset2        uint
	CountTexts     uint
	TableTextOff   uint
	Offset3        uint

	Texts []texts
}
