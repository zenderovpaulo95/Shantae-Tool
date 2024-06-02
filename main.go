/*
  Shantae Tool is fan localization tool, made by Sudakov Pavel
  Thanks to aluigi for script
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"shantae/methods"
)

func main() {
	args := os.Args

	fmt.Println("Инструментарий для игры Shantae and Pirate's curse. Автор программы: Судаков Павел")
	fmt.Println("Особая благодарность aluigi за скрипт для разбора ресурсов")
	fmt.Println("Версия 1.0")

	fmt.Println("Исходный код на GitFlic: https://gitflic.ru/project/pashok6798/shantae-tool")

	if len(args) > 1 {
		if len(args) == 3 && ((args[1] == "-la") || (args[1] == "-lv")) {
			if _, err := os.Stat(args[2]); err == nil {
				if args[1] == "-la" {
					list, err := methods.ReadArcHeader(args[2])

					if err != nil {
						panic(err)
					}

					for i := 0; i < len(list); i++ {
						fmt.Printf("%d. %016x\t%d     %s\n", (i + 1), list[i].FileOffset, list[i].UncompressedSize, list[i].FileName)
					}
				} else {
					volList, anotherVolList, err := methods.ReadFileHeader(args[2])

					if err != nil {
						panic(err)
					}

					if volList != nil {
						for i := 0; i < len(volList); i++ {
							fmt.Printf("%d. %08x\t%d     %s\n", (i + 1), volList[i].Offset, volList[i].Size, volList[i].FileName)
						}
					} else {
						for i := 0; i < len(anotherVolList); i++ {
							fmt.Printf("%d. %08x\t%d     %s\n", (i + 1), anotherVolList[i].Offset, anotherVolList[i].Size, anotherVolList[i].FileName)
						}
					}
				}
			}
		}
		if ((len(args) == 4) || (len(args) == 5)) && (args[1] == "-compare") {
			orIndex := 1 //По умолчанию будет значение 1 (сравнение с первой строкой)

			_, err := os.Stat(args[2])
			if err != nil {
				panic(err)
			}

			_, err = os.Stat(args[3])
			if err != nil {
				panic(err)
			}

			if len(args) == 5 {
				orIndex, err = strconv.Atoi(args[4])
				if err != nil {
					orIndex = 1
				}
			}

			err = methods.CompareFiles(args[2], args[3], orIndex)

			if err != nil {
				panic(err)
			}

			fmt.Println("Найдены непереведённые строки.")
		}
		if ((len(args) == 4) || (len(args) == 5)) && (args[1] == "-replace") {
			orIndex := 1 //По умолчанию будет значение 1 (замена певой строки)

			_, err := os.Stat(args[2])

			if os.IsNotExist(err) {
				panic(err)
			}

			_, err = os.Stat(args[3])

			if err != nil {
				panic(err)
			}

			if len(args) == 5 {
				orIndex, err = strconv.Atoi(args[4])
				if err != nil {
					orIndex = 1
				}
			}

			err = methods.ReplaceFromLocalizedFile(args[2], args[3], orIndex)

			if err != nil {
				panic(err)
			}

			fmt.Println("В оригинальный файл успешно перенесены переведённые строки")
		}
		if ((len(args) == 3) || (len(args) == 4)) && ((args[1] == "-ea") || (args[1] == "-ev")) {
			if _, err := os.Stat(args[2]); err == nil {
				outputFilePath := filepath.Dir(args[0]) + "/Unpacked"

				if len(args) == 4 {
					_, err = os.Stat(args[3])

					if os.IsNotExist(err) {
						panic(err)
					}

					outputFilePath = args[3]
				}

				_, err = os.Stat(outputFilePath)

				if os.IsNotExist(err) {
					os.MkdirAll(outputFilePath, os.ModePerm)
				}

				fmt.Println("Идёт распаковка...")

				if args[1] == "-ea" {
					list, err := methods.ReadArcHeader(args[2])
					err = methods.UnpackArchive(list, args[2], outputFilePath)

					if err != nil {
						panic(err)
					}
				} else {
					volList, anotherVolList, err := methods.ReadFileHeader(args[2])

					if err != nil {
						panic(err)
					}

					if volList != nil {
						err = methods.UnpackVol(volList, args[2], outputFilePath)

						if err != nil {
							panic(err)
						}
					} else {
						err = methods.UnpackAnotherVol(anotherVolList, args[2], outputFilePath)

						if err != nil {
							panic(err)
						}
					}
				}
			}
		}
		if (len(args) == 3) && (args[1] == "-lf") {
			font, err := methods.ReadHeader(args[2])

			if err != nil {
				panic(err)
			}

			for i := 0; i < int(font.KernPairsCount); i++ {
				fmt.Printf("%d\t%d\t%d\n", font.KernPairs[i].FirstChar, font.KernPairs[i].SecondChar, font.KernPairs[i].Amount)
			}

			font.Chars = methods.SortFonts(font.Chars)

			for i := 0; i < int(font.CharsCount); i++ {
				fmt.Printf("%d.\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%x\n", (i + 1), font.Chars[i].Char, font.Chars[i].X, font.Chars[i].Y, font.Chars[i].Unknown3, font.Chars[i].Unknown4, font.Chars[i].Unknown5, font.Chars[i].Unknown6, font.Chars[i].Unknown7)
			}

			fmt.Println()

			font.UnknownData = methods.SortUD(font.UnknownData)

			for i := 0; i < int(font.CharsCount); i++ {
				fmt.Printf("%d.\t%d\t%d\t%d\t%d\t%x\n", (i + 1), font.UnknownData[i].Unknown1, font.UnknownData[i].Unknown2, font.UnknownData[i].Unknown3, font.UnknownData[i].Unknown4, font.UnknownData[i].Unknown5)
			}
		}
		if (len(args) == 3) && (args[1] == "-lt") {
			text, err := methods.ReadTextHeader(args[2])

			if err != nil {
				panic(err)
			}

			for i := 0; i < int(text.CountTexts); i++ {
				for j := 0; j < int(text.CountLocTexts); j++ {
					fmt.Printf("%d. %s\n", (j + 1), text.TextStrings[i].Texts[j])
				}

				fmt.Println()
			}
		}
		if ((len(args) == 3) || (len(args) == 4)) && (args[1] == "-et") {
			if _, err := os.Stat(args[2]); err == nil {
				outputFilePath := filepath.Dir(args[0]) + "/Unpacked"

				if len(args) == 4 {
					_, err = os.Stat(args[3])

					if os.IsNotExist(err) {
						panic(err)
					}

					outputFilePath = args[3]
				}

				_, err = os.Stat(outputFilePath)

				if os.IsNotExist(err) {
					os.MkdirAll(outputFilePath, os.ModePerm)
				}

				text, err := methods.ReadTextHeader(args[2])

				if err != nil {
					panic(err)
				}

				err = methods.ExtractText(text, args[2], outputFilePath)
				if err != nil {
					panic(err)
				}

				fmt.Println("Файл успешно извлечён.")
			}
		}
		if (len(args) == 4) && (args[1] == "-rt") {
			if _, err := os.Stat(args[2]); err == nil {
				txtFileName := args[3]

				text, err := methods.ReadTextHeader(args[2])

				if err != nil {
					panic(err)
				}

				err = methods.ReplaceText(text, args[2], txtFileName)
				if err != nil {
					panic(err)
				}

				fmt.Println("Файл успешно модифицирован.")
			}
		}
		if ((len(args) == 3) || (len(args) == 4)) && ((args[1] == "-ra") || (args[1] == "-rv")) {
			if _, err := os.Stat(args[2]); err == nil {
				outputFilePath := filepath.Dir(args[0]) + "/Unpacked"

				if len(args) == 4 {
					_, err = os.Stat(args[3])

					if os.IsNotExist(err) {
						panic(err)
					}

					outputFilePath = args[3]
				}

				_, err = os.Stat(outputFilePath)

				if os.IsNotExist(err) {
					os.MkdirAll(outputFilePath, os.ModePerm)
				}

				fmt.Println("Идёт запаковка...")

				if args[1] == "-ra" {
					list, err := methods.ReadArcHeader(args[2])
					err = methods.RepackArchive(list, args[2], outputFilePath)

					if err != nil {
						panic(err)
					}
				} else {
					volList, anotherVolList, err := methods.ReadFileHeader(args[2])

					if err != nil {
						panic(err)
					}

					if volList != nil {
						err = methods.RepackVol(volList, args[2], outputFilePath)

						if err != nil {
							panic(err)
						}
					} else {
						err = methods.RepackAnotherVol(anotherVolList, args[2], outputFilePath)

						if err != nil {
							panic(err)
						}
					}
				}
			}
		}
	} else {
		fmt.Println("Как пользоваться программой.")
		fmt.Printf("%s -ea arc.data - извлечь vol файлы из архива. По умолчанию файлы извлекутся рядом с программой в папку Unpacked.\n", args[0])
		fmt.Printf("%s -ea arc.data \"путь/к/папке/с извлечёнными ресурсами\" - извлечь vol файлы из архива в указанную папку.\n", args[0])
		fmt.Printf("%s -ra arc.data - перепаковать vol файлы в архив. По умолчанию папка Unpacked, находящаяся рядом с программой.\n", args[0])
		fmt.Printf("%s -ra arc.data \"путь/к/папке/с извлечёнными ресурсами\" - перепаковать vol файлы в архив из указанной папки с ресурсами.\n", args[0])
		fmt.Printf("%s -la arc.data - получить список файлов в архиве.\n", args[0])
		fmt.Printf("%s -lv file.vol - получить список файлов в vol файле.\n", args[0])
		fmt.Printf("%s -ev file.vol - распаковать vol файл.\n", args[0])
		fmt.Printf("%s -ev file.vol \"путь/к/папке/с извлечёнными ресурсами\" - распаковать vol файл в указанную папку.\n", args[0])
		fmt.Printf("%s -rv file.vol - перепаковать vol файл.\n", args[0])
		fmt.Printf("%s -rv file.vol \"путь/к/папке/с извлечёнными ресурсами\" - перепаковать vol файл из указанной папки\n", args[0])
		fmt.Printf("%s -et file.vol - распаковать text файл.\n", args[0])
		fmt.Printf("%s -et file.vol \"путь/к/папке/с извлечёнными ресурсами\" - распаковать text файл в указанную папку.\n", args[0])
		fmt.Printf("%s -rt file.vol loc_file.txt - перепаковать text файл.\n", args[0])
		fmt.Printf("%s -replace or_file.txt loc_file - Заменить оригинальные строки на локализованные.\n", args[0])
		fmt.Printf("%s -replace or_file.txt loc_file or_index - Заменить оригинальные строки на локализованные с указанием номера строки.\n", args[0])
		fmt.Printf("%s -compare original.txt replaced_original.txt - найти непереведённые строки в локализованном файле.\n", args[0])
		fmt.Printf("%s -compare original.txt replaced_original.txt or_index - найти непереведённые строки в локализованном файле с указанием номера строки.\n", args[0])

	}

}
