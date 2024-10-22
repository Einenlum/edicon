package ini

import (
	"errors"
	"fmt"

	"github.com/einenlum/edicon/internal/core"
)

type OutputType int

const (
	FullOutput OutputType = iota
	KeyValuesOnlyOutput
)

func toIniOutputType(otype core.OutputType) (OutputType, error) {
	switch otype {
	case core.FullOutput:
		return FullOutput, nil
	case core.MeaningFullOutput:
		return KeyValuesOnlyOutput, nil
	default:
		err := fmt.Sprintf("Invalid output type: %d", otype)

		return FullOutput, errors.New(err)
	}
}

func OutputConfigFile(iniFile *IniConfiguration, outputType OutputType) string {
	output := ""

	shouldBePrinted := func(line Line) bool {
		if outputType != KeyValuesOnlyOutput {
			return true
		}

		return line.ContentType != OtherType
	}

	for _, line := range iniFile.GlobalSection.Lines {
		if shouldBePrinted(*line) {
			output += line.ToString() + "\n"
		}
	}

	for _, section := range iniFile.Sections {
		for _, line := range section.Lines {
			if shouldBePrinted(*line) {
				output += line.ToString() + "\n"
			}
		}
	}

	return output
}

func GetParsedIniFile(filePath string) (IniConfiguration, error) {
	parsedLines, err := ParseIniFile(filePath)
	if err != nil {
		return IniConfiguration{}, err
	}

	globalSection, sections := getSections(parsedLines)
	return IniConfiguration{globalSection, sections, filePath}, nil
}

func GetSectionByName(sections []*Section, name string) *Section {
	for _, section := range sections {
		if section.Name == name {
			return section
		}
	}

	return nil
}

func getKeyLine(lines []*Line, key string) *Line {
	for _, line := range lines {
		if line.ContentType == KeyValueType && line.KeyValue.Key == key {
			return line
		}
	}

	return nil
}

func getKeyLineBySectionName(sections []*Section, sectionName string, key string) *Line {
	section := GetSectionByName(sections, sectionName)
	if section == nil {
		return nil
	}

	return getKeyLine(section.Lines, key)
}

func EditConfigFile(
	notationStyle core.NotationStyle,
	filePath string,
	key string,
	value string,
) (*IniConfiguration, error) {
	decomposedKey := core.DecomposeKey(notationStyle, key)

	iniFile, err := GetParsedIniFile(filePath)
	if err != nil {
		return &IniConfiguration{}, err
	}

	if len(decomposedKey) == 1 {
		keyLine := getKeyLine(iniFile.GlobalSection.Lines, decomposedKey[0])
		if keyLine == nil {
			return &IniConfiguration{}, errors.New("Key not found")
		}

		keyLine.SetValue(value)
	} else {
		keyLine := getKeyLineBySectionName(iniFile.Sections, decomposedKey[0], decomposedKey[1])
		if keyLine == nil {
			return &IniConfiguration{}, errors.New("Key not found")
		}

		keyLine.SetValue(value)
	}

	return &iniFile, nil
}

func GetParameterFromPath(notationStyle core.NotationStyle, filePath string, key string) (string, error) {
	iniFile, err := GetParsedIniFile(filePath)
	if err != nil {
		return "", err
	}

	decomposedKey := core.DecomposeKey(notationStyle, key)
	if len(decomposedKey) == 1 {
		keyLine := getKeyLine(iniFile.GlobalSection.Lines, decomposedKey[0])
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
