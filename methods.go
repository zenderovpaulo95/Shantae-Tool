package methods

import (
	"bufio"
	"encoding/binary"
	"errors"
	"os"
)

type ListFiles struct {
	fileName         string
	isCompressed     byte
	fileOffset       uint64
	compressedSize   uint64
	uncompressedSize uint64
	index            int //Для правильной пересборки архивов
}

type ArcHeader struct {
	header              uint
	unknown1            uint16
	unknown2            uint16
	unknown3            uint16
	unknown4            uint16
	offset1             uint
	filesCount          int
	filenameOffset      uint
	unknownDataOffset   uint
	compressLogicOffset uint
	cSizeOffset         uint
	sizeOffset          uint
	fileOffset          uint
	baseOffset          uint
}

func sortListFiles(arcHead []ListFiles) []ListFiles {
	for i := 1; i < len(arcHead); i++ {
		for j := i; j > 0; j-- {
			if arcHead[j-1].fileOffset > arcHead[j].fileOffset {
				tmpArc := arcHead[j-1]
				arcHead[j-1] = arcHead[j]
				arcHead[j] = tmpArc
			}
		}
	}

	return arcHead
}

func readArcHeader(fileName string) (arcHead []ListFiles, err error) {
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

	buf := bufio.NewReader(file)
	tmpByte = make([]byte, 4)
	_, bufErr := buf.Read(tmpByte)

	if bufErr != nil {
		return nil, bufErr
	}

	head.header = uint(binary.LittleEndian.Uint32(tmpByte))

	if head.header != 0x18f32f12 {
		err = errors.New("Archive header: Unknown header!")
		return nil, err
	}

	tmpByte = make([]byte, 2)
	_, bufErr = buf.Read(tmpByte)
	head.unknown1 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 2)
	_, bufErr = buf.Read(tmpByte)
	head.unknown2 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 2)
	_, bufErr = buf.Read(tmpByte)
	head.unknown3 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 2)
	_, bufErr = buf.Read(tmpByte)
	head.unknown4 = binary.LittleEndian.Uint16(tmpByte)

	tmpByte = make([]byte, 4)
	_, bufErr = buf.Read(tmpByte)
	head.offset1 = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = buf.Read(tmpByte)
	head.filesCount = int(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = buf.Read(tmpByte)
	head.filenameOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = buf.Read(tmpByte)
	head.unknownDataOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = buf.Read(tmpByte)
	head.compressLogicOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = buf.Read(tmpByte)
	head.cSizeOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = buf.Read(tmpByte)
	head.sizeOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = buf.Read(tmpByte)
	head.fileOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, bufErr = buf.Read(tmpByte)
	head.baseOffset = uint(binary.LittleEndian.Uint32(tmpByte))

	arcHead = make([]ListFiles, head.filesCount)

	//Просчитываю индексы для будущих списков
	for i := 0; i < head.filesCount; i++ {
		arcHead[i].index = i + 1
	}

	//Считываю логику сжатых файлов (сжаты файлы или нет)
	_, err = file.Seek(int64(head.compressLogicOffset), 0)

	for i := 0; i < head.filesCount; i++ {
		arcHead[i].isCompressed, err = buf.ReadByte()
	}

	//Считываю смещения к файлам
	_, err = file.Seek(int64(head.fileOffset), 0)

	for i := 0; i < head.filesCount; i++ {
		tmpByte = make([]byte, 8)
		_, err = buf.Read(tmpByte)
		arcHead[i].fileOffset = binary.LittleEndian.Uint64(tmpByte)
	}

	//Считываю размеры сжатых файлов
	_, err = file.Seek(int64(head.cSizeOffset), 0)

	for i := 0; i < head.filesCount; i++ {
		tmpByte = make([]byte, 8)
		_, err = buf.Read(tmpByte)
		arcHead[i].compressedSize = binary.LittleEndian.Uint64(tmpByte)
	}

	//Считываю размеры файлов
	_, err = file.Seek(int64(head.sizeOffset), 0)

	for i := 0; i < head.filesCount; i++ {
		tmpByte = make([]byte, 8)
		_, err = buf.Read(tmpByte)
		arcHead[i].uncompressedSize = binary.LittleEndian.Uint64(tmpByte)
	}

	//Считываю названия файлов
	/*_, err = file.Seek(int64(head.sizeOffset), 0)

	for i := 0; i < head.filesCount; i++ {
		tmpByte = make([]byte, 8)
		_, err = buf.Read(tmpByte)
		arcHead[i].uncompressedSize = binary.LittleEndian.Uint64(tmpByte)
	}*/

	return
}
