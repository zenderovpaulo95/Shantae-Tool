package methods

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func UnpackAnotherVol(listFiles []AnotherListFiles, fileName string, outputDir string) (err error) {
	file, err := os.Open(fileName)

	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < len(listFiles); i++ {
		_, err = file.Seek(int64(listFiles[i].Offset), 0)

		if err != nil {
			return err
		}

		tmpByte := make([]byte, listFiles[i].Size)
		_, err = file.Read(tmpByte)

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

		outFile.Write(tmpByte)

		fmt.Printf("%d. %08x\t%d     %s\n", (i + 1), listFiles[i].Offset, listFiles[i].Size, listFiles[i].FileName)
	}

	return
}

func UnpackVol(listFiles []ListVolFiles, fileName string, outputDir string) (err error) {
	file, err := os.Open(fileName)

	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < len(listFiles); i++ {
		_, err = file.Seek(int64(listFiles[i].Offset), 0)

		if err != nil {
			return err
		}

		tmpByte := make([]byte, listFiles[i].Size)
		_, err = file.Read(tmpByte)

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

		outFile.Write(tmpByte)

		fmt.Printf("%d. %08x\t%d     %s\n", (i + 1), listFiles[i].Offset, listFiles[i].Size, listFiles[i].FileName)
	}

	return
}

func UnpackArchive(listFiles []ListFiles, fileName string, outputDir string) (err error) {
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
