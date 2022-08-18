package methods

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type Texts struct {
	CRC         uint
	TextOffsets []uint
	Texts       []string
}

type TextHeader struct {
	Header        uint
	Unknown1      uint
	Unknown2      byte
	Unknown3      byte
	CountLocTexts uint16
	Offset1       uint
	Offset2       uint
	CountTexts    uint
	TableTextOff  uint
	Offset3       uint

	TextStrings []Texts
}

func ReadTextHeader(fileName string) (Text TextHeader, err error) {
	file, err := os.Open(fileName)

	if err != nil {
		return
	}
	defer file.Close()

	tmpByte := make([]byte, 4)
	file.Read(tmpByte)
	Text.Header = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	file.Read(tmpByte)
	Text.Unknown1 = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 1)
	file.Read(tmpByte)
	Text.Unknown2 = tmpByte[0]

	tmpByte = make([]byte, 1)
	file.Read(tmpByte)
	Text.Unknown3 = tmpByte[0]

	tmpByte = make([]byte, 2)
	file.Read(tmpByte)
	Text.CountLocTexts = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 4)
	file.Read(tmpByte)
	Text.Offset1 = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	file.Read(tmpByte)
	Text.Offset2 = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	file.Read(tmpByte)
	Text.CountTexts = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	file.Read(tmpByte)
	Text.TableTextOff = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	file.Read(tmpByte)
	Text.Offset3 = uint(binary.LittleEndian.Uint32(tmpByte))

	Text.TextStrings = make([]Texts, Text.CountTexts)

	for i := 0; i < int(Text.CountTexts); i++ {
		Text.TextStrings[i].TextOffsets = make([]uint, Text.CountLocTexts)
		Text.TextStrings[i].Texts = make([]string, Text.CountLocTexts)

		tmpByte = make([]byte, 4)
		file.Read(tmpByte)
		Text.TextStrings[i].CRC = uint(binary.LittleEndian.Uint32(tmpByte))

		var offset int64 = 0

		for j := 0; j < int(Text.CountLocTexts); j++ {
			tmpByte = make([]byte, 4)
			file.Read(tmpByte)
			Text.TextStrings[i].TextOffsets[j] = uint(binary.LittleEndian.Uint32(tmpByte))
			offset, _ = file.Seek(0, io.SeekCurrent)
			file.Seek(int64(Text.TextStrings[i].TextOffsets[j]), 0)

			tmpByte = make([]byte, 1)
			tmpByte[0] = 0xFF
			len := 0

			for ; tmpByte[0] != 0; len++ {
				tmpByte = make([]byte, 1)
				file.Read(tmpByte)
			}

			file.Seek(int64(Text.TextStrings[i].TextOffsets[j]), 0)
			tmpByte = make([]byte, len-1)
			file.Read(tmpByte)
			Text.TextStrings[i].Texts[j] = string(tmpByte)
			fmt.Println(Text.TextStrings[i].Texts[j])

			file.Seek(offset, 0)
		}
	}

	return
}
