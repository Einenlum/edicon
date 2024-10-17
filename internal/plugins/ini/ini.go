package ini

import (
	"einenlum/edicon/internal/core"
	"errors"
)

type OutputType int

const (
	FullOutput OutputType = iota
	KeyValuesOnlyOutput
)

func OutputIniFile(iniFile *IniFile, outputType OutputType) string {
	output := ""

	shouldBePrinted := func(line Line) bool {
		if outputType != KeyValuesOnlyOutput {
			return true
		}

		return line.ContentType != OtherType
	}

	for _, section := range *iniFile.Sections {
		for _, line := range *section.Lines {
			if shouldBePrinted(line) {
				output += line.ToString() + "\n"
			}
		}
	}

	return output
}

func GetParsedIniFile(filePath string) (IniFile, error) {
	parsedLines, err := ParseIniFile(filePath)
	if err != nil {
		return IniFile{}, err
	}

	return IniFile{getSections(&parsedLines), filePath}, nil
}

func GetSectionByName(sections *[]Section, name string) *Section {
	for _, section := range *sections {
		if section.Name == name {
			return &section
		}
	}

	return nil
}

func getKeyLine(section *Section, key string) *Line {
	for _, line := range *section.Lines {
		if line.ContentType == KeyValueType && line.KeyValue.Key == key {
			return &line
		}
	}

	return nil
}

func getKeyLineBySectionName(sections *[]Section, sectionName string, key string) *Line {
	section := GetSectionByName(sections, sectionName)
	if section == nil {
		return nil
	}

	return getKeyLine(section, key)
}

func EditIniFile(notationStyle core.NotationStyle, filePath string, key string, value string) (IniFile, error) {
	decomposedKey := core.DecomposeKey(notationStyle, key)

	iniFile, err := GetParsedIniFile(filePath)
	if err != nil {
		return IniFile{}, err
	}

	if len(decomposedKey) == 1 {
		keyLine := getKeyLineBySectionName(iniFile.Sections, "PHP", decomposedKey[0])
		if keyLine == nil {
			return IniFile{}, errors.New("Key not found")
		}

		keyLine.SetValue(value)
	} else {
		keyLine := getKeyLineBySectionName(iniFile.Sections, decomposedKey[0], decomposedKey[1])
		if keyLine == nil {
			return IniFile{}, errors.New("Key not found")
		}

		keyLine.SetValue(value)
	}

	return iniFile, nil
}

func GetIniParameterFromPath(notationStyle core.NotationStyle, filePath string, key string) (string, error) {
	iniFile, err := GetParsedIniFile(filePath)
	if err != nil {
		return "", err
	}

	decomposedKey := core.DecomposeKey(notationStyle, key)
	if len(decomposedKey) == 1 {
		keyLine := getKeyLineBySectionName(iniFile.Sections, "PHP", decomposedKey[0])
		if keyLine == nil {
			return "", errors.New("Key not found")
		}

		return keyLine.KeyValue.Value, nil
	} else {
		keyLine := getKeyLineBySectionName(iniFile.Sections, decomposedKey[0], decomposedKey[1])
		if keyLine == nil {
			return "", errors.New("Key not found")
		}

		return keyLine.KeyValue.Value, nil
	}
}
