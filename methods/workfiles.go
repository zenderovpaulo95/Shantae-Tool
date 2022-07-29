package methods

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io/ioutil"
	"os"
)

func Unpack(listFiles []ListFiles, fileName string, outputDir string) (err error) {
	file, err := os.Open(fileName)

	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < len(listFiles); i++ {
		_, err = file.Seek(int64(listFiles[i].FileOffset), 0)

		if err != nil {
			return nil
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
