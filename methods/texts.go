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

//CompareFiles - сравнение оригинального файла с оригинальным файлом, в котором заменили строки
func CompareFiles(firstFile, secondFile string, selIndex int) (err error) {
	first, err := os.Open(firstFile)
	if err != nil {
		return err
	}
	defer first.Close()

	second, err := os.Open(secondFile)
	if err != nil {
		return err
	}
	defer second.Close()

	firstScan := bufio.NewScanner(first)
	secondScan := bufio.NewScanner(second)

	var fCountTxts, fCountLocTxts int
	var sCountTxts, sCountLocTxts int

	firstScan.Scan()
	secondScan.Scan()

	tmpStrs := strings.Split(firstScan.Text(), " ")
	fCountTxts, _ = strconv.Atoi(tmpStrs[0])
	fCountLocTxts, _ = strconv.Atoi(tmpStrs[1])

	tmpStrs = strings.Split(firstScan.Text(), " ")
	sCountTxts, _ = strconv.Atoi(tmpStrs[0])
	sCountLocTxts, _ = strconv.Atoi(tmpStrs[1])

	if fCountTxts != sCountTxts {
		err = fmt.Errorf("Количество строк из первого файла не совпадает с количеством строк из второго файла")
		return err
	}

	if fCountLocTxts != sCountLocTxts {
		err = fmt.Errorf("Количество локализованных строк из первого файла не совпадает с количеством локализованных строк из второго файла")
		return err
	}

	var fStrs [][]string
	var sStrs [][]string
	var row   []string

	for i := 0; i < fCountTxts; i++ {
		for j := 0; j < fCountLocTxts; j++ {
			firstScan.Scan()
			row = append(row, firstScan.Text())
		}
		fStrs = append(fStrs, row)
		row = nil
		firstScan.Scan()
	}

	for i := 0; i < sCountTxts; i++ {
		for j := 0; j < sCountLocTxts; j++ {
			secondScan.Scan()
			row = append(row, secondScan.Text())
		}
		sStrs = append(sStrs, row)
		row = nil
		secondScan.Scan()
	}

	var nonTranslatedStrs []string

	for i := 0; i < fCountTxts; i++ {
		if fStrs[i][selIndex - 1] == sStrs[i][selIndex - 1] {
			nonTranslatedStrs = append(nonTranslatedStrs, fStrs[i][selIndex - 1])
		}
	}

	if len(nonTranslatedStrs) > 0 {
		newFilePath := firstFile[:strings.LastIndex(firstFile, ".txt")] + "_non-translated.txt"
		newFile, err := os.Create(newFilePath)

		if err != nil {
			return err
		}
		defer newFile.Close()

		for i := 0; i < len(nonTranslatedStrs); i++ {
			newFile.WriteString(nonTranslatedStrs[i])
			newFile.WriteString("\r\n")
		}
	} else {
		err = fmt.Errorf("На найдено непереведённых строк")
		return err
	}

	return
}

func ReplaceFromLocalizedFile(originalFile string, translatedFile string, orIndex int) (err error) {
	original, err := os.Open(originalFile)
	if err != nil {
		return err
	}
	defer original.Close()

	translate, err := os.Open(translatedFile)
	if err != nil {
		return err
	}
	defer translate.Close()

	var orCountTexts, orCountLocTexts int
	var trCountTexts, trCountLocTexts int

	orScanner := bufio.NewScanner(original)
	trScanner := bufio.NewScanner(translate)

	orScanner.Scan()
	orFirstString := orScanner.Text()
	tmpStrs := strings.Split(orFirstString, " ")
	orCountTexts, _ = strconv.Atoi(tmpStrs[0])
	orCountLocTexts, _ = strconv.Atoi(tmpStrs[1])

	trScanner.Scan()
	tmpStrs = strings.Split(trScanner.Text(), " ")
	trCountTexts, _ = strconv.Atoi(tmpStrs[0])
	trCountLocTexts, _ = strconv.Atoi(tmpStrs[1])

	if orCountLocTexts < orIndex + 1 || trCountLocTexts < orIndex + 1 {
		err = fmt.Errorf("Недостаточно локализованных строк для сравнения")
		return err
	}

	var orStrings [][]string
	var trStrings [][]string
	var row []string

	for i := 0; i < orCountTexts; i++ {
		for j := 0; j < orCountLocTexts; j++ {
			orScanner.Scan()
			row = append(row, orScanner.Text())
		}
		orScanner.Scan()
		orStrings = append(orStrings, row)
		row = nil
	}

	for i := 0; i < trCountTexts; i++ {
		for j := 0; j < trCountLocTexts; j++ {
			trScanner.Scan()
			row = append(row, trScanner.Text())
		}
		trScanner.Scan()
		trStrings = append(trStrings, row)
		row = nil
	}

	//Английская версия обычно идёт первой. Значит, будем пробовать
	//искать значения у соседней строки и надеяться, что всё заменится без проблем

	for i := 0; i < orCountTexts; i++ {
		for j := 0; j < trCountTexts; j++ {
			if orStrings[i][orIndex] == trStrings[j][orIndex] {
				orStrings[i][orIndex - 1] = trStrings[j][orIndex - 1]
			}
		}
	}

	newFileName := originalFile[:strings.LastIndex(originalFile, ".txt")] + "_translated.txt"
	fmt.Println(newFileName)
	newFile, err := os.Create(newFileName)

	if err != nil {
		return err
	}
	defer newFile.Close()

	newFile.WriteString(orFirstString)
	newFile.WriteString("\r\n")

	for i := 0; i < orCountTexts; i++ {
		for j := 0; j < orCountLocTexts; j++ {
			newFile.WriteString(orStrings[i][j])
			newFile.WriteString("\r\n")
		}
		newFile.WriteString("\r\n")
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
				if strings.Contains(text.TextStrings[i].Texts[j], "\n") {
					text.TextStrings[i].Texts[j] = strings.Replace(text.TextStrings[i].Texts[j], "\n", "\\n", -1)
				}
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
						if strings.Contains(text.TextStrings[i].Texts[j], "\\n") {
							text.TextStrings[i].Texts[j] = strings.Replace(text.TextStrings[i].Texts[j], "\\n", "\n", -1)
						}
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
			for i := 0; i < len(newStrings); i++ {
				fmt.Println(newStrings[i])
			}
			newOffsets := make([]uint, len(newStrings))

			newOffsets[0] = startOffset

			for i := 1; i < len(newStrings); i++ {
				len := len([]byte(newStrings[i-1])) + 1
				newOffsets[i] = newOffsets[i-1] + uint(len)
			}

			for i := 0; i < len(newStrings); i++ {
				for j := 0; j < int(text.CountTexts); j++ {
					for k := 0; k < int(text.CountLocTexts); k++ {
						if text.TextStrings[j].Texts[k] == newStrings[i] {
							text.TextStrings[j].TextOffsets[k] = newOffsets[i]
						}
					}
				}
			}

			file, err := os.Open(fileName)
			if err != nil {
				return err
			}
			defer file.Close()

			tmpByte := make([]byte, startOffset)
			file.Read(tmpByte)
			newFile, err := os.Create(fileName + ".tmp")

			if err != nil {
				return err
			}
			defer newFile.Close()

			newFile.Write(tmpByte)

			for i := 0; i < len(newStrings); i++ {
				ch := 0
				tmpStr := fmt.Sprintf("%s%c", newStrings[i], ch)
				tmpByte = []byte(tmpStr)
				newFile.Write(tmpByte)
			}

			newFile.Seek(int64(text.TableTextOff), 0)

			for i := 0; i < int(text.CountTexts); i++ {
				tmpByte = make([]byte, 4)
				binary.LittleEndian.PutUint32(tmpByte, uint32(text.TextStrings[i].CRC))
				newFile.Write(tmpByte)

				for j := 0; j < int(text.CountLocTexts); j++ {
					tmpByte = make([]byte, 4)
					binary.LittleEndian.PutUint32(tmpByte, uint32(text.TextStrings[i].TextOffsets[j]))
					newFile.Write(tmpByte)
				}
			}
		}
	} else {
		err = fmt.Errorf("В файле нет строк")
		return err
	}

	return
}
