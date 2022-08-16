package methods

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

func RepackArchive(listFiles []ListFiles, fileName string, outputDir string) (err error) {
	file, err := os.Open(fileName)

	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < len(listFiles); i++ {
		_, err = file.Seek(int64(listFiles[i].FileOffset), 0)

		if err != nil {
			return err
		}

		fileSize := listFiles[i].UncompressedSize

		if listFiles[i].IsCompressed {
			fileSize = listFiles[i].CompressedSize
		}

		tmpFile := make([]byte, fileSize)

		_, err = file.Read(tmpFile)

		if listFiles[i].IsCompressed {
			r := flate.NewReader(bytes.NewReader(tmpFile))

			enflated, err := ioutil.ReadAll(r)
			if err != nil {
				return err
			}

			tmpFile = make([]byte, len(enflated))
			copy(tmpFile, enflated)
			enflated = nil
		}

		_, err = os.Stat(outputDir + "/" + listFiles[i].FileName)

		if !os.IsNotExist(err) {
			os.Remove(outputDir + "/" + listFiles[i].FileName)
		}

		if strings.ContainsAny(listFiles[i].FileName, "/") {
			dirPath := listFiles[i].FileName[:strings.LastIndex(listFiles[i].FileName, "/")+1]

			fmt.Println(dirPath)

			_, err = os.Stat(outputDir + "/" + dirPath)

			if os.IsNotExist(err) {
				os.MkdirAll(outputDir+"/"+dirPath, os.ModePerm)
			}
		}

		outFile, err := os.Create(outputDir + "/" + listFiles[i].FileName)

		if err != nil {
			return err
		}
		defer outFile.Close()

		outFile.Write(tmpFile)

		fmt.Printf("%d. %016x\t%d     %s\n", (i + 1), listFiles[i].FileOffset, listFiles[i].UncompressedSize, listFiles[i].FileName)
	}

	return nil
}
