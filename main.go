package main

import (
	"einenlum/edicon/internal/plugins/php"
	"fmt"
)

func main() {
	php_ini_content, err := php.GetPhpIniFileContent("/etc/php/php.ini")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("PHP ini file content: ", php_ini_content)
}
