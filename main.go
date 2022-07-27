/*
  Shantae Tool is fan localization tool, made by Sudakov Pavel
  Thanks to aluigi for script
*/
package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	fmt.Println("Shantae tool made by Sudakov Pavel")
	fmt.Println("Thanks to aluigi for script")
	fmt.Println("Version 1.0")

	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			fmt.Printf("%s\n", args[i])
		}
	} else {
		fmt.Println("How to use my tool.")
		fmt.Printf("%s -ea arc.file - extract files from archive. Default extraction path near tool's path.\n", args[0])
		fmt.Printf("%s -ea arc.file \"path/to/extracted/files\" - extract files from archive into extracted folder\n", args[0])
		fmt.Printf("%s -ra arc.file - repack files into archive. Default resource folder is a tool's path.\n", args[0])
		fmt.Printf("%s -ra arc.file \"path/to/extracted/files\" - repack files from resource folder into archive\n", args[0])
	}

}
