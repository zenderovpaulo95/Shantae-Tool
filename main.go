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
		if len(args) == 3 && args[1] == "-l" {
			if _, err := os.Stat(args[2]); err == nil {
				list, err := methods.ReadArcHeader(args[2])

				if err != nil {
					panic(err)
				}

				for i := 0; i < len(list); i++ {
					fmt.Printf("%d. %016x\t%d     %s\n", (i + 1), list[i].FileOffset, list[i].UncompressedSize, list[i].FileName)
				}
			}
		}
		if ((len(args) == 3) || (len(args) == 4)) && args[1] == "-ea" {
			if _, err := os.Stat(args[2]); err == nil {
				list, err := methods.ReadArcHeader(args[2])

				if err != nil {
					panic(err)
				}
				outputFilePath := filepath.Dir(args[0]) + "Unpacked"

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

				err = methods.Unpack(list, args[2], outputFilePath)

				if err != nil {
					panic(err)
				}
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
		fmt.Printf("%s -l arc.data - get list of files in archive\n", args[0])
	}

}
