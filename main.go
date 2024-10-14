package main

import (
	"einenlum/edicon/internal/plugins/ini"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: edicon <file-path> set <key> <value>")

		return
	}

	argsWithoutProg := os.Args[1:]
	filePath := argsWithoutProg[0]
	// action := argsWithoutProg[1]
	key := argsWithoutProg[2]
	value := argsWithoutProg[3]

	iniFile, err := ini.EditIniFile(filePath, key, value)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println(ini.OutputIniFile(iniFile, ini.FullOutput))
}
