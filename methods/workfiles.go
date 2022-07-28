package methods

import (
	"encoding/binary"
	"errors"
	"os"

)

func Unpack(listFiles []ListFiles, fileName string) (err error) {
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

		if bool(listFiles[i].IsCompressed) {
			fileSize = listFiles[i].CompressedSize
		}

		tmpFile := make([]byte, fileSize)

		_, err = file.Read(tmpFile)

		if bool(listFiles[i].IsCompressed) {
			//КАК ИЗВЛЕКАТЬ АРХИВЫ С DEFLATE ЗНАЧЕНИЕМ?!
		}
	}

	return nil
}
