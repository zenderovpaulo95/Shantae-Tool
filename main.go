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
			outputFilePath := filepath.Dir(args[0])

			if len(args == 4) {
				_, err := os.Stat(args[3])

				if os.IsNotExists(err) {
					panic(err)
				}

				outputFilePath = args[3]
			}
		}
			
		}
		if ((len(args) == 3) || (len(args) == 4)) && args[1] == "-ra" {
			//Do something later...
		}
	} else {
		fmt.Println("How to use my tool.")
		fmt.Printf("%s -ea arc.file - extract files from archive. Default extraction path near tool's path.\n", args[0])
		fmt.Printf("%s -ea arc.file \"path/to/extracted/files\" - extract files from archive into extracted folder\n", args[0])
		fmt.Printf("%s -ra arc.file - repack files into archive. Default resource folder is a tool's path.\n", args[0])
		fmt.Printf("%s -ra arc.file \"path/to/extracted/files\" - repack files from resource folder into archive\n", args[0])
		fmt.Printf("%s -l arc.file - get list of files in archive\n", args[0])
	}

}
