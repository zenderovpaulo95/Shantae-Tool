package methods

import (
	"encoding/binary"
	"errors"
	"os"
)

type UnknownStruct struct {
	Unknown1 uint16
	Unknown2 uint16
	Unknown3 int16
	Unknown4 int16
	Unknown5 uint
}

type Coordinates struct {
	Char     uint
	X        uint16
	Y        uint16
	Unknown3 int16
	Unknown4 int16
	Unknown5 int16
	Unknown6 uint16
	Unknown7 uint
}

type KernPair struct {
	FirstChar  uint16
	SecondChar uint16
	Amount     int16
}

type FontHeader struct {
	Header          uint
	Unknown1        uint16
	CharsCount      uint16
	KernPairsCount  uint16
	Unknown2        byte
	Unknown3        byte
	Offset1         uint
	Offset2         uint
	CharsOffset     uint
	KernPairsOffset uint
	Unknown4        uint16
	Unknown5        int16
	Unknown6        uint16
	Unknown7        uint16
	UnknownOffset1  uint
	UnknownOffset2  uint

	KernPairs   []KernPair
	Chars       []Coordinates
	UnknownData []UnknownStruct
}

func ReadHeader(fileName string) (Font FontHeader, err error) {
	file, err := os.Open(fileName)

	if err != nil {
		return
	}
	defer file.Close()

	tmpByte := make([]byte, 4)
	_, err = file.Read(tmpByte)
	Font.Header = uint(binary.LittleEndian.Uint32(tmpByte))

	if Font.Header != 0xBAF8A21A {
		err = errors.New("Font: Wrong header.")
		return
	}

	tmpByte = make([]byte, 2)
	_, err = file.Read(tmpByte)
	Font.Unknown1 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 2)
	_, err = file.Read(tmpByte)
	Font.CharsCount = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 2)
	_, err = file.Read(tmpByte)
	Font.KernPairsCount = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 1)
	_, err = file.Read(tmpByte)
	Font.Unknown2 = tmpByte[0]

	tmpByte = make([]byte, 1)
	_, err = file.Read(tmpByte)
	Font.Unknown3 = tmpByte[0]

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	Font.Offset1 = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	Font.Offset2 = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	Font.CharsOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	Font.KernPairsOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 2)
	_, err = file.Read(tmpByte)
	Font.Unknown4 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 2)
	_, err = file.Read(tmpByte)
	Font.Unknown5 = int16(binary.LittleEndian.Uint16(tmpByte))

	tmpByte = make([]byte, 2)
	_, err = file.Read(tmpByte)
	Font.Unknown6 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 2)
	_, err = file.Read(tmpByte)
	Font.Unknown7 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	Font.UnknownOffset1 = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	Font.UnknownOffset2 = uint(binary.LittleEndian.Uint32(tmpByte))

	Font.KernPairs = make([]KernPair, Font.KernPairsCount)
	Font.Chars = make([]Coordinates, Font.CharsCount)
	Font.UnknownData = make([]UnknownStruct, Font.CharsCount)

	_, err = file.Seek(int64(Font.KernPairsOffset), 0)

	if err != nil {
		return
	}

	for i := 0; i < int(Font.KernPairsCount); i++ {
		tmpByte = make([]byte, 2)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.KernPairs[i].FirstChar = binary.LittleEndian.Uint16(tmpByte)

		tmpByte = make([]byte, 2)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.KernPairs[i].SecondChar = binary.LittleEndian.Uint16(tmpByte)

		tmpByte = make([]byte, 2)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.KernPairs[i].Amount = int16(binary.LittleEndian.Uint16(tmpByte))
	}

	_, err = file.Seek(int64(Font.CharsOffset), 0)

	if err != nil {
		return
	}

	for i := 0; i < int(Font.CharsCount); i++ {
		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.Chars[i].Char = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 2)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.Chars[i].X = binary.LittleEndian.Uint16(tmpByte)

		tmpByte = make([]byte, 2)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.Chars[i].Y = binary.LittleEndian.Uint16(tmpByte)

		tmpByte = make([]byte, 2)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.Chars[i].Unknown3 = int16(binary.LittleEndian.Uint16(tmpByte))

		tmpByte = make([]byte, 2)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.Chars[i].Unknown4 = int16(binary.LittleEndian.Uint16(tmpByte))

		tmpByte = make([]byte, 2)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.Chars[i].Unknown5 = int16(binary.LittleEndian.Uint16(tmpByte))

		tmpByte = make([]byte, 2)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.Chars[i].Unknown6 = binary.LittleEndian.Uint16(tmpByte)

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.Chars[i].Unknown7 = uint(binary.LittleEndian.Uint32(tmpByte))
	}

	_, err = file.Seek(int64(Font.UnknownOffset2), 0)

	for i := 0; i < int(Font.CharsCount); i++ {
		tmpByte = make([]byte, 2)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.UnknownData[i].Unknown1 = binary.LittleEndian.Uint16(tmpByte)

		tmpByte = make([]byte, 2)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.UnknownData[i].Unknown2 = binary.LittleEndian.Uint16(tmpByte)

		tmpByte = make([]byte, 2)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.UnknownData[i].Unknown3 = int16(binary.LittleEndian.Uint16(tmpByte))

		tmpByte = make([]byte, 2)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.UnknownData[i].Unknown4 = int16(binary.LittleEndian.Uint16(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)

		if err != nil {
			return
		}

		Font.UnknownData[i].Unknown5 = uint(binary.LittleEndian.Uint32(tmpByte))
	}

	return
}
