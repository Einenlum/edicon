package main

import (
	"einenlum/edicon/internal/plugins/ini"
	"fmt"
)

func main() {
	sections, err := ini.GetSectionsFromIniFile("data/php.ini")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	for _, section := range sections {
		// fmt.Printf("%+v\n", line)
		fmt.Println("[" + section.Name + "]")

		for _, line := range section.Lines {
			if line.ContentType == ini.KeyValueType {
				fmt.Println(line.StringContent)
			}
		}
	}
}
