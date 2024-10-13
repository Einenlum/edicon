package io

import (
	"os"
)

func GetFileContents(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

func GenerateRandomTempFilePath() (string, error) {
	file, err := os.CreateTemp("", "tempfile-*.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	return file.Name(), nil
}

func WriteFileContents(filepath string, content string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
