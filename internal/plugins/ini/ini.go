package ini

import (
	"errors"
)

func GetParsedIniFile(filePath string) (IniFile, error) {
	parsedLines, err := ParseIniFile(filePath)
	if err != nil {
		return IniFile{[]Section{}, filePath}, err
	}

	return IniFile{getSections(parsedLines), filePath}, nil
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
	for _, line := range section.Lines {
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

func EditIniFile(filePath string, key string, value string) (*IniFile, error) {
	decomposedKey := DecomposeKeyWithDotNotation(key)

	iniFile, err := GetParsedIniFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(decomposedKey) == 1 {
		keyLine := getKeyLineBySectionName(&iniFile.Sections, "PHP", decomposedKey[0])
		if keyLine == nil {
			return nil, errors.New("Key not found")
		}

		keyLine.KeyValue.Value = value
	} else {
		keyLine := getKeyLineBySectionName(&iniFile.Sections, decomposedKey[0], decomposedKey[1])
		if keyLine == nil {
			return nil, errors.New("Key not found")
		}

		keyLine.KeyValue.Value = value
	}

	return &iniFile, nil
}

func GetIniParameterFromPath(filePath string, key string) (string, error) {
	iniFile, err := GetParsedIniFile(filePath)
	if err != nil {
		return "", err
	}

	decomposedKey := DecomposeKeyWithDotNotation(key)
	if len(decomposedKey) == 1 {
		keyLine := getKeyLineBySectionName(&iniFile.Sections, "PHP", decomposedKey[0])
		if keyLine == nil {
			return "", errors.New("Key not found")
		}

		return keyLine.KeyValue.Value, nil
	} else {
		keyLine := getKeyLineBySectionName(&iniFile.Sections, decomposedKey[0], decomposedKey[1])
		if keyLine == nil {
			return "", errors.New("Key not found")
		}

		return keyLine.KeyValue.Value, nil
	}
}
