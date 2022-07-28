package methods

import (
	"encoding/binary"
	"errors"
	"os"
)

type ListFiles struct {
	FileName         string
	IsCompressed     byte
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
		arcHead[i].Index = i + 1
	}

	//Считываю логику сжатых файлов (сжаты файлы или нет)
	_, err = file.Seek(int64(head.CompressLogicOffset), 0)

	for i := 0; i < head.FilesCount; i++ {
		tmpByte = make([]byte, 1)
		_, bufErr = file.Read(tmpByte)

		arcHead[i].IsCompressed = tmpByte[0]
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

		tmpByte = make([]byte, len - 1)
		_, err = file.Read(tmpByte)
		arcHead[i].FileName = string(tmpByte)
	}

	return
}
