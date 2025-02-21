package methods

import (
	"encoding/binary"
	"errors"
	"os"
	"strings"
)

type ListFiles struct {
	FileName         string
	IsCompressed     bool
	FileNameOffset   uint32
	FileOffset       uint64
	CompressedSize   uint64
	UncompressedSize uint64
	Index            int //Для правильной пересборки архивов
}

type ArcHeader struct {
	Header              uint
	Unknown1            uint16
	Unknown2            uint16
	Unknown3            uint16
	Unknown4            uint16
	Offset1             uint
	FilesCount          int
	FilenameOffset      uint
	UnknownDataOffset   uint
	CompressLogicOffset uint
	CSizeOffset         uint
	SizeOffset          uint
	FileOffset          uint
	BaseOffset          uint
}

type ListVolFiles struct {
	CRC      uint
	Unknown1 uint //Смешение к названию файла
	Offset   uint
	Zero1    uint
	Zero2    uint
	Size     uint
	FileName string
	Index    int
}

type AnotherListFiles struct {
	CRC        uint
	Offset     uint64
	NameOffset uint
	Size       uint
	FileName   string
	Index      int
}

type VolAnotherHeader struct {
	Header     uint
	Unknown1   uint16
	Unknown2   uint16
	Unknown3   uint
	Unknown4   uint
	FilesCount uint
	InfoOff    uint
	DataOff    uint
	FileSize   uint
}

type VolHeader struct {
	Header     uint
	Unknown1   uint
	Unknown2   uint
	FileCount  int
	FileOffset uint
	BaseOffset uint
	FileSize   uint //Полный размер файла
}

func RemoveDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func SortListFiles(arcHead []ListFiles) []ListFiles {
	for i := 1; i < len(arcHead); i++ {
		for j := i; j > 0; j-- {
			if arcHead[j-1].FileOffset > arcHead[j].FileOffset {
				tmpArc := arcHead[j-1]
				arcHead[j-1] = arcHead[j]
				arcHead[j] = tmpArc
			}
		}
	}

	return arcHead
}

func PadSize(inSize int64) int64 {
	for inSize%16 != 0 {
		inSize++
	}

	return inSize
}

func SortVolListFiles(volHead []ListVolFiles) []ListVolFiles {
	for i := 1; i < len(volHead); i++ {
		for j := i; j > 0; j-- {
			if volHead[j-1].Offset > volHead[j].Offset {
				tmpVol := volHead[j-1]
				volHead[j-1] = volHead[j]
				volHead[j] = tmpVol
			}
		}
	}

	return volHead
}

func SortAnotherVolListFiles(AnotherVolHead []AnotherListFiles) []AnotherListFiles {
	for i := 1; i < len(AnotherVolHead); i++ {
		for j := i; j > 0; j-- {
			if AnotherVolHead[j-1].Offset > AnotherVolHead[j].Offset {
				tmpAnVol := AnotherVolHead[j-1]
				AnotherVolHead[j-1] = AnotherVolHead[j]
				AnotherVolHead[j] = tmpAnVol
			}
		}
	}

	return AnotherVolHead
}

func SortFonts(chars []Coordinates) []Coordinates {
	for i := 1; i < len(chars); i++ {
		for j := i; j > 0; j-- {
			//if chars[j-1].Char > chars[j].Char {
			if chars[j-1].Unknown7 > chars[j].Unknown7 {
				tmpChar := chars[j-1]
				chars[j-1] = chars[j]
				chars[j] = tmpChar
			}
		}
	}

	return chars
}

func SortUD(data []UnknownStruct) []UnknownStruct {
	for i := 1; i < len(data); i++ {
		for j := i; j > 0; j-- {
			if data[j-1].Unknown5 > data[j].Unknown5 {
				tmpData := data[j-1]
				data[j-1] = data[j]
				data[j] = tmpData
			}
		}
	}

	return data
}

func ReadFileHeader(fileName string) (volHead []ListVolFiles, volAnHead []AnotherListFiles, err error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()

	if statsErr != nil {
		return nil, nil, statsErr
	}

	var size int64 = stats.Size()

	if size <= 0 {
		err = errors.New("Vol stat: File either empty or it doesn't exists")
		return nil, nil, err
	}

	tmpByte := make([]byte, 4)
	_, err = file.Read(tmpByte)

	if err != nil {
		return nil, nil, err
	}

	checkHeader := binary.LittleEndian.Uint32(tmpByte)

	switch checkHeader {
	case 0x9ee6c000:
		volHead, err = ReadVolHeader(fileName)

		if err != nil {
			return nil, nil, err
		}

		volAnHead = nil

		break

	case 0xb53d32cb:
		volAnHead, err = ReadAnotherVolHeader(fileName)

		if err != nil {
			return nil, nil, err
		}

		volHead = nil

		break
	}

	return
}

func ReadAnotherVolHeader(fileName string) (volAnHead []AnotherListFiles, err error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var head VolAnotherHeader
	tmpByte := make([]byte, 4)
	_, err = file.Read(tmpByte)
	head.Header = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 2)
	_, err = file.Read(tmpByte)
	head.Unknown1 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 2)
	_, err = file.Read(tmpByte)
	head.Unknown2 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	head.Unknown3 = uint(binary.LittleEndian.Uint32(tmpByte))

	if head.Unknown3 == 0x1C {
		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		head.FilesCount = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		head.InfoOff = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		head.DataOff = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		head.FileSize = uint(binary.LittleEndian.Uint32(tmpByte))
	} else {
		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		head.Unknown4 = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		head.DataOff = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		head.FilesCount = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		head.InfoOff = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		head.FileSize = uint(binary.LittleEndian.Uint32(tmpByte))
	}

	_, err = file.Seek(int64(head.InfoOff), 0)

	if err != nil {
		return nil, err
	}

	volAnHead = make([]AnotherListFiles, head.FilesCount)

	for i := 0; i < int(head.FilesCount); i++ {
		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		volAnHead[i].CRC = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		volAnHead[i].NameOffset = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 8)
		_, err = file.Read(tmpByte)
		volAnHead[i].Offset = binary.LittleEndian.Uint64(tmpByte)

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		volAnHead[i].Size = uint(binary.LittleEndian.Uint32(tmpByte))

		volAnHead[i].Index = i
	}

	for i := 0; i < int(head.FilesCount); i++ {
		_, err = file.Seek(int64(volAnHead[i].NameOffset), 0)

		tmpByte = make([]byte, 1)
		tmpByte[0] = 0xFF

		var len int = 0

		for ; tmpByte[0] != 0; len++ {
			tmpByte = make([]byte, 1)
			_, err = file.Read(tmpByte)
		}

		_, err = file.Seek(int64(volAnHead[i].NameOffset), 0)
		tmpByte = make([]byte, len-1)
		_, err = file.Read(tmpByte)
		volAnHead[i].FileName = string(tmpByte)
		volAnHead[i].FileName = strings.Replace(volAnHead[i].FileName, "_", "/", -1)
	}

	return
}

func ReadVolHeader(fileName string) (volHead []ListVolFiles, err error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()

	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()

	if size <= 0 {
		err = errors.New("Vol stat: File either empty or it doesn't exists")
		return nil, err
	}

	var head VolHeader
	var tmpByte []byte

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	head.Header = uint(binary.LittleEndian.Uint32(tmpByte))

	if head.Header != 0x9ee6c000 {
		err = errors.New("Vol header: unknown header")
		return nil, err
	}

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	head.Unknown1 = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	head.Unknown2 = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	head.FileCount = int(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	head.FileOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	head.BaseOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, err = file.Read(tmpByte)
	head.FileSize = uint(binary.LittleEndian.Uint32(tmpByte))

	volHead = make([]ListVolFiles, head.FileCount)

	_, err = file.Seek(int64(head.FileOffset), 0)

	//Считываю сначала списки по 24 байта
	for i := 0; i < head.FileCount; i++ {
		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		volHead[i].CRC = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		volHead[i].Unknown1 = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		volHead[i].Offset = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		volHead[i].Zero1 = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		volHead[i].Zero2 = uint(binary.LittleEndian.Uint32(tmpByte))

		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		volHead[i].Size = uint(binary.LittleEndian.Uint32(tmpByte))

		volHead[i].Index = i
	}

	off := head.FileOffset + uint(24*head.FileCount)
	//А потом таблицу с названием файлов
	for i := 0; i < head.FileCount; i++ {
		_, err = file.Seek(int64(off), 0)
		len := 0

		tmpByte = make([]byte, 1)
		tmpByte[0] = 1

		for tmpByte[0] != 0 {
			_, err = file.Read(tmpByte)

			len++
		}

		_, err = file.Seek(int64(off), 0)

		tmpByte = make([]byte, len-1)
		_, err = file.Read(tmpByte)
		volHead[i].FileName = string(tmpByte)
		volHead[i].FileName = strings.Replace(volHead[i].FileName, "_", "/", -1)
		off += uint(len)
	}

	return
}

func ReadArcHeader(fileName string) (arcHead []ListFiles, err error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()

	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()

	if size <= 0 {
		err = errors.New("Stat: File either empty or it doesn't exists!")
		return nil, err
	}

	var head ArcHeader
	var tmpByte []byte

	tmpByte = make([]byte, 4)
	_, bufErr := file.Read(tmpByte)

	if bufErr != nil {
		return nil, bufErr
	}

	head.Header = uint(binary.LittleEndian.Uint32(tmpByte))

	if head.Header != 0x18f32f12 {
		err = errors.New("Archive header: Unknown header!")
		return nil, err
	}

	tmpByte = make([]byte, 2)
	_, bufErr = file.Read(tmpByte)
	head.Unknown1 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 2)
	_, bufErr = file.Read(tmpByte)
	head.Unknown2 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 2)
	_, bufErr = file.Read(tmpByte)
	head.Unknown3 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 2)
	_, bufErr = file.Read(tmpByte)
	head.Unknown4 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 4)
	_, bufErr = file.Read(tmpByte)
	head.Offset1 = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = file.Read(tmpByte)
	head.FilesCount = int(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = file.Read(tmpByte)
	head.FilenameOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = file.Read(tmpByte)
	head.UnknownDataOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = file.Read(tmpByte)
	head.CompressLogicOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = file.Read(tmpByte)
	head.CSizeOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = file.Read(tmpByte)
	head.SizeOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = file.Read(tmpByte)
	head.FileOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = file.Read(tmpByte)
	head.BaseOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	arcHead = make([]ListFiles, head.FilesCount)

	//Просчитываю индексы для будущих списков
	for i := 0; i < head.FilesCount; i++ {
		arcHead[i].Index = i
	}

	//Считываю логику сжатых файлов (сжаты файлы или нет)
	_, err = file.Seek(int64(head.CompressLogicOffset), 0)

	for i := 0; i < head.FilesCount; i++ {
		tmpByte = make([]byte, 1)
		_, bufErr = file.Read(tmpByte)

		arcHead[i].IsCompressed = true
		if int(tmpByte[0]) != 1 {
			arcHead[i].IsCompressed = false
		}
	}

	//Считываю смещения к файлам
	_, err = file.Seek(int64(head.FileOffset), 0)

	for i := 0; i < head.FilesCount; i++ {
		tmpByte = make([]byte, 8)
		_, err = file.Read(tmpByte)
		arcHead[i].FileOffset = binary.LittleEndian.Uint64(tmpByte)
	}

	//Считываю размеры сжатых файлов
	_, err = file.Seek(int64(head.CSizeOffset), 0)

	for i := 0; i < head.FilesCount; i++ {
		tmpByte = make([]byte, 8)
		_, err = file.Read(tmpByte)
		arcHead[i].CompressedSize = binary.LittleEndian.Uint64(tmpByte)
	}

	//Считываю размеры файлов
	_, err = file.Seek(int64(head.SizeOffset), 0)

	for i := 0; i < head.FilesCount; i++ {
		tmpByte = make([]byte, 8)
		_, err = file.Read(tmpByte)
		arcHead[i].UncompressedSize = binary.LittleEndian.Uint64(tmpByte)
	}

	//Считываю названия файлов
	_, err = file.Seek(int64(head.FilenameOffset), 0)
	for i := 0; i < head.FilesCount; i++ {
		tmpByte = make([]byte, 4)
		_, err = file.Read(tmpByte)
		arcHead[i].FileNameOffset = binary.LittleEndian.Uint32(tmpByte)
	}

	len := 0

	for i := 0; i < head.FilesCount; i++ {
		_, err = file.Seek(int64(arcHead[i].FileNameOffset), 0)
		len = 0

		tmpByte = make([]byte, 1)
		tmpByte[0] = 1

		for tmpByte[0] != 0 {
			_, bufErr = file.Read(tmpByte)

			len++
		}

		_, err = file.Seek(int64(arcHead[i].FileNameOffset), 0)

		tmpByte = make([]byte, len-1)
		_, err = file.Read(tmpByte)
		arcHead[i].FileName = string(tmpByte)
		arcHead[i].FileName = strings.Replace(arcHead[i].FileName, "_", "/", -1)
	}

	return
}
