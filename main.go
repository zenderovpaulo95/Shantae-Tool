/*
  Shantae Tool is fan localization tool, made by Sudakov Pavel
  Thanks to aluigi for script
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"shantae/methods"
)

func main() {
	args := os.Args

	fmt.Println("Shantae tool made by Sudakov Pavel")
	fmt.Println("Thanks to aluigi for script")
	fmt.Println("Version 1.0")

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
					list, err := methods.ReadVolHeader(args[2])

					if err != nil {
						panic(err)
					}

					for i := 0; i < len(list); i++ {
						fmt.Printf("%d. %08x\t%d     %s\n", (i + 1), list[i].Offset, list[i].Size, list[i].FileName)
					}
				}
			}
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
					os.MkdirAll(outputFilePath, 0666)
				}

				fmt.Println("Unpacking...")

				if args[1] == "-ea" {
					list, err := methods.ReadArcHeader(args[2])
					err = methods.UnpackArchive(list, args[2], outputFilePath)

					if err != nil {
						panic(err)
					}
				} else {
					list, err := methods.ReadVolHeader(args[2])
					err = methods.UnpackVol(list, args[2], outputFilePath)
					if err != nil {
						panic(err)
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
				fmt.Printf("%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", font.Chars[i].Char, font.Chars[i].X, font.Chars[i].Y, font.Chars[i].Unknown3, font.Chars[i].Unknown4, font.Chars[i].Unknown5, font.Chars[i].Unknown6, font.Chars[i].Unknown7, font.Chars[i].Unknown8)
			}

			fmt.Println()

			for i := 0; i < int(font.CharsCount); i++ {
				fmt.Printf("%d\t%d\t%d\t%d\t%x\n", font.UnknownData[i].Unknown1, font.UnknownData[i].Unknown2, font.UnknownData[i].Unknown3, font.UnknownData[i].Unknown4, font.UnknownData[i].Unknown5)
			}
		}
		if ((len(args) == 3) || (len(args) == 4)) && args[1] == "-ra" {
			//Do something later...
		}
	} else {
		fmt.Println("How to use my tool.")
		fmt.Printf("%s -ea arc.data - extract files from archive. Default extraction path near tool's path.\n", args[0])
		fmt.Printf("%s -ea arc.data \"path/to/extracted/files\" - extract files from archive into extracted folder\n", args[0])
		fmt.Printf("%s -ra arc.data - repack files into archive. Default resource folder is a tool's path.\n", args[0])
		fmt.Printf("%s -ra arc.data \"path/to/extracted/files\" - repack files from resource folder into archive\n", args[0])
		fmt.Printf("%s -la arc.data - get list of files in archive\n", args[0])
		fmt.Printf("%s -lv arc.vol - get list of files in archive\n", args[0])
	}

}
