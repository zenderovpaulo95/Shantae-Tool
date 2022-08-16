package methods

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"os"
	"io"
)

func RepackAnotherVol(listFiles []AnotherListFiles, fileName string, resourceDir string) (err error) {
	file, err := os.Open(fileName)
	listFiles = SortAnotherVolListFiles(listFiles)

	if err != nil {
		return err
	}
	defer file.Close()

	tmpByte := make([]byte, listFiles[0].Offset)
	_, err = file.Read(tmpByte)

	commonFileSize := uint(listFiles[0].Offset)

	_, _ = file.Seek(16, 0)
	tmpBOff := make([]byte, 4)
	file.Read(tmpBOff)
	listFilesOff := binary.LittleEndian.Uint32(tmpBOff)

	newFile, err := os.Create(fileName + ".tmp")

	if err != nil {
		return err
	}
	defer newFile.Close()

	newFile.Write(tmpByte)
	offset := listFiles[0].Offset

	for i := 0; i < len(listFiles); i++ {
		_, _ = file.Seek(int64(listFiles[i].Offset), 0)
		listFiles[i].Offset = offset
		file, err = os.Open(resourceDir + "/" + listFiles[i].FileName)
		if err != nil {
			return err
		}
		defer file.Close()

		fi, _ := file.Stat()

		tmpByte = make([]byte, fi.Size())
		_, err = file.Read(tmpByte)

		listFiles[i].Size = uint(fi.Size())

		paddedSize := PadSize(fi.Size())
		paddedFile := make([]byte, paddedSize)

		copy(paddedFile, tmpByte)

		newFile.Write(paddedFile)

		offset += uint64(paddedSize)
		commonFileSize += uint(paddedSize)

		fmt.Printf("%d. %08x\t%d     %s\n", (i + 1), listFiles[i].Offset, listFiles[i].Size, listFiles[i].FileName)
	}

	newFile.Seek(24, 0)
	tmpByte = make([]byte, 4)
	binary.LittleEndian.PutUint32(tmpByte, uint32(commonFileSize))
	newFile.Write(tmpByte)

	for i := 0; i < len(listFiles); i++ {
		_, _ = newFile.Seek(int64(listFilesOff+(uint32(listFiles[i].Index)*24)+8), 0)
		tmpByte = make([]byte, 4)
		binary.LittleEndian.PutUint32(tmpByte, uint32(listFiles[i].Offset))
		newFile.Write(tmpByte)

		_, _ = newFile.Seek(int64(listFilesOff+(uint32(listFiles[i].Index)*24)+20), 0)
		tmpByte = make([]byte, 4)
		binary.LittleEndian.PutUint32(tmpByte, uint32(listFiles[i].Size))
		newFile.Write(tmpByte)
	}

	return
}

func RepackVol(listFiles []ListVolFiles, fileName string, resourceDir string) (err error) {
	file, err := os.Open(fileName)
	listFiles = SortVolListFiles(listFiles)

	if err != nil {
		return err
	}
	defer file.Close()

	tmpByte := make([]byte, listFiles[0].Offset)
	_, err = file.Read(tmpByte)

	commonFileSize := listFiles[0].Offset

	_, _ = file.Seek(16, 0)
	tmpBOff := make([]byte, 4)
	file.Read(tmpBOff)
	listFilesOff := binary.LittleEndian.Uint32(tmpBOff)

	newFile, err := os.Create(fileName + ".tmp")

	if err != nil {
		return err
	}
	defer newFile.Close()

	newFile.Write(tmpByte)
	offset := listFiles[0].Offset

	for i := 0; i < len(listFiles); i++ {
		_, _ = file.Seek(int64(listFiles[i].Offset), 0)
		listFiles[i].Offset = offset
		file, err = os.Open(resourceDir + "/" + listFiles[i].FileName)
		if err != nil {
			return err
		}
		defer file.Close()

		fi, _ := file.Stat()

		tmpByte = make([]byte, fi.Size())
		_, err = file.Read(tmpByte)

		listFiles[i].Size = uint(fi.Size())

		paddedSize := PadSize(fi.Size())
		paddedFile := make([]byte, paddedSize)

		copy(paddedFile, tmpByte)

		newFile.Write(paddedFile)

		offset += uint(paddedSize)
		commonFileSize += uint(paddedSize)

		fmt.Printf("%d. %08x\t%d     %s\n", (i + 1), listFiles[i].Offset, listFiles[i].Size, listFiles[i].FileName)
	}

	newFile.Seek(24, 0)
	tmpByte = make([]byte, 4)
	binary.LittleEndian.PutUint32(tmpByte, uint32(commonFileSize))
	newFile.Write(tmpByte)

	for i := 0; i < len(listFiles); i++ {
		_, _ = newFile.Seek(int64(listFilesOff+(uint32(listFiles[i].Index)*24)+8), 0)
		tmpByte = make([]byte, 4)
		binary.LittleEndian.PutUint32(tmpByte, uint32(listFiles[i].Offset))
		newFile.Write(tmpByte)

		_, _ = newFile.Seek(int64(listFilesOff+(uint32(listFiles[i].Index)*24)+20), 0)
		tmpByte = make([]byte, 4)
		binary.LittleEndian.PutUint32(tmpByte, uint32(listFiles[i].Size))
		newFile.Write(tmpByte)
	}

	return
}

func RepackArchive(listFiles []ListFiles, fileName string, resourceDir string) (err error) {
	file, err := os.Open(fileName)

	if err != nil {
		return err
	}
	defer file.Close()

	listFiles = SortListFiles(listFiles)

	file.Seek(32, 0)
	tmpByte := make([]byte, 4)
	_, _ = file.Read(tmpByte)
	cSizeOff := uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, _ = file.Read(tmpByte)
	SizeOff := uint(binary.LittleEndian.Uint32(tmpByte))

	tmpByte = make([]byte, 4)
	_, _ = file.Read(tmpByte)
	FileOffs := uint(binary.LittleEndian.Uint32(tmpByte))

	file.Seek(0, 0)
	tmpByte = make([]byte, listFiles[0].FileOffset)
	_, _ = file.Read(tmpByte)

	newFile, err := os.Create(fileName + ".tmp")
	if err != nil {
		return err
	}
	defer newFile.Close()
	newFile.Write(tmpByte)
	offset := listFiles[0].FileOffset

	for i := 0; i < len(listFiles); i++ {
		newFile.Seek(int64(offset), 0)
		listFiles[i].FileOffset = offset
		tmpFile, err := os.Open(resourceDir + "/" + listFiles[i].FileName)
		if err != nil {
			panic(err)
			//return err
		}
		defer tmpFile.Close()

		fi, _ := tmpFile.Stat()
		listFiles[i].UncompressedSize = uint64(fi.Size())
		listFiles[i].CompressedSize = uint64(fi.Size())
		tmpByte = make([]byte, listFiles[i].UncompressedSize)
		tmpFile.Read(tmpByte)

		if listFiles[i].IsCompressed {
			r := bytes.NewReader(tmpByte)
			var b bytes.Buffer
			flateWriter, err := flate.NewWriter(&b, flate.BestCompression)

			if err != nil {
				panic(err)
			}

			io.Copy(flateWriter, r)

			listFiles[i].CompressedSize = uint64(b.Len())
			paddedSize := PadSize(int64(listFiles[i].CompressedSize))

			tmpByte = make([]byte, paddedSize)
			copy(tmpByte, b.Bytes())
			newFile.Write(tmpByte)

			flateWriter.Flush()

			offset += uint64(paddedSize)
		} else {
			paddedSize := PadSize(int64(listFiles[i].UncompressedSize))
			tmpB := make([]byte, paddedSize)
			copy(tmpB, tmpByte)
			newFile.Write(tmpB)
			offset += uint64(paddedSize)
		}

		newFile.Seek(int64(cSizeOff + uint(listFiles[i].Index * 8)), 0)
		tmpByte = make([]byte, 8)
		binary.LittleEndian.PutUint64(tmpByte, listFiles[i].CompressedSize)
		newFile.Write(tmpByte)

		newFile.Seek(int64(SizeOff + uint(listFiles[i].Index * 8)), 0)
		tmpByte = make([]byte, 8)
		binary.LittleEndian.PutUint64(tmpByte, listFiles[i].UncompressedSize)
		newFile.Write(tmpByte)

		newFile.Seek(int64(FileOffs + uint(listFiles[i].Index * 8)), 0)
		tmpByte = make([]byte, 8)
		binary.LittleEndian.PutUint64(tmpByte, listFiles[i].FileOffset)
		newFile.Write(tmpByte)

		fmt.Printf("%d. %016x\t%d\t%d\t%d     %s\n", (i + 1), listFiles[i].FileOffset, listFiles[i].UncompressedSize, listFiles[i].CompressedSize, listFiles[i].IsCompressed, listFiles[i].FileName)
	}
	return nil
}
