package main

import (
	"einenlum/edicon/internal/plugins/ini"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: edicon <file-path>")

		return
	}

	argsWithoutProg := os.Args[1:]
	filePath := argsWithoutProg[0]

	iniFile, err := ini.GetParsedIniFile(filePath)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	for _, section := range iniFile.Sections {
		// fmt.Printf("%+v\n", line)
		fmt.Println("[" + section.Name + "]")

		for _, line := range section.Lines {
			if line.ContentType == ini.KeyValueType {
				fmt.Println(line.StringContent)
			}
		}
	}
}
