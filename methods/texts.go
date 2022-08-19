package methods

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
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

	if (Text.Header != 0xB9247E83) && (Text.Header != 0x54585453) {
		err = fmt.Errorf("Заголовок не text файла")
		return
	}

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
	file.Seek(int64(Text.TableTextOff), 0)

	var offset int64 = 0

	for i := 0; i < int(Text.CountTexts); i++ {
		Text.TextStrings[i].TextOffsets = make([]uint, Text.CountLocTexts)
		Text.TextStrings[i].Texts = make([]string, Text.CountLocTexts)

		tmpByte = make([]byte, 4)
		file.Read(tmpByte)
		Text.TextStrings[i].CRC = uint(binary.LittleEndian.Uint32(tmpByte))

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

			file.Seek(offset, 0)
		}
	}

	return
}

func ExtractText(text TextHeader, fileName string, outputDir string) (err error) {
	if (text.CountTexts > 0) && (text.CountLocTexts > 0) {
		fi, _ := os.Stat(fileName)
		newFileName := outputDir + "/" + fmt.Sprintf("%s", fi.Name())
		newFileName = newFileName[:strings.LastIndex(newFileName, ".text")] + ".txt"

		file, err := os.Create(newFileName)
		if err != nil {
			return err
		}
		defer file.Close()

		str := fmt.Sprintf("%d %d\r\n", text.CountTexts, text.CountLocTexts)
		file.WriteString(str)

		for i := 0; i < int(text.CountTexts); i++ {
			for j := 0; j < int(text.CountLocTexts); j++ {
				file.WriteString(text.TextStrings[i].Texts[j])
				file.WriteString("\r\n")
			}
			file.WriteString("\r\n")
		}
	} else {
		err = fmt.Errorf("В файле нет строк")
		return err
	}

	return
}

func ReplaceText(text TextHeader, fileName string, txtFileName string) (err error) {
	if (text.CountTexts > 0) && (text.CountLocTexts > 0) {
		txtFile, err := os.Open(txtFileName)

		if err != nil {
			return err
		}
		defer txtFile.Close()

		scanner := bufio.NewScanner(txtFile)

		scanner.Scan()
		newStr := strings.Split(scanner.Text(), " ")
		if len(newStr) == 2 {
			countTexts, _ := strconv.Atoi(newStr[0])
			countLocTexts, _ := strconv.Atoi(newStr[1])

			newStrings := []string{}

			if countTexts <= int(text.CountTexts) && countLocTexts <= int(text.CountLocTexts) {
				for i := 0; i < countTexts; i++ {
					for j := 0; j < countLocTexts; j++ {
						scanner.Scan()
						//fmt.Printf("Было: %s\n", text.TextStrings[i].Texts[j])
						text.TextStrings[i].Texts[j] = scanner.Text()
						newStrings = append(newStrings, text.TextStrings[i].Texts[j])
						//fmt.Printf("Стало: %s\n", text.TextStrings[i].Texts[j])
					}
					scanner.Scan()
					//fmt.Println()
				}
			} else {
				err = fmt.Errorf("Количество строк больше строк в оригинальном файле")
				return err
			}

			startOffset := text.TableTextOff + uint(text.CountTexts*(4+(4*uint(text.CountLocTexts))))
			//fmt.Printf("%x\n", startOffset)

			sort.Slice(newStrings, func(i, j int) bool {
				return newStrings[i] < newStrings[j]
			})

			newStrings = RemoveDuplicate(newStrings)
			newOffsets := make([]uint, len(newStrings))

			newOffsets[0] = startOffset

			for i := 1; i < len(newStrings); i++ {
				len := len([]byte(newStrings[i-1])) + 1
				newOffsets[i] = newOffsets[i-1] + uint(len)
			}
		}

		/*for scanner.Scan() {
			newStr := scanner.Text()

			if strings.Index(newStr, ". ") > 0 {

				newStr = newStr[strings.Index(newStr, ". ")+2:]
				fmt.Println(newStr)
			}
		}*/

	} else {
		err = fmt.Errorf("В файле нет строк")
		return err
	}

	return
}
