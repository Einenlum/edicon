package ini

import (
	"strings"

	"github.com/einenlum/edicon/internal/io"
)

func parseLineString(lineNumber int, lineString string) Line {
	var spacePrefix string

	trimmedLineString := strings.TrimSpace(lineString)
	if trimmedLineString == lineString {
		spacePrefix = ""
	} else {
		spacePrefix = lineString[:len(lineString)-len(trimmedLineString)]
	}

	// Check if line is empty
	if len(trimmedLineString) == 0 {
		return Line{
			lineNumber,
			lineString,
			spacePrefix,
			Original,
			OtherType,
			nil,
			nil,
		}
	}

	// Check if line is a comment
	if strings.HasPrefix(trimmedLineString, ";") {
		return Line{
			lineNumber,
			lineString,
			spacePrefix,
			Original,
			OtherType,
			nil,
			nil,
		}
	}

	// Check if line is a section
	if strings.HasPrefix(trimmedLineString, "[") && strings.HasSuffix(trimmedLineString, "]") {
		sectionName := strings.Trim(trimmedLineString, "[]")

		sectionLine := SectionLine{sectionName}

		return Line{
			lineNumber,
			lineString,
			spacePrefix,
			Original,
			SectionLineType,
			nil,
			&sectionLine,
		}
	}

	// Check if line is a key value pair
	if strings.Contains(trimmedLineString, "=") {
		keyValue := strings.Split(trimmedLineString, "=")

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		return Line{
			lineNumber,
			lineString,
			spacePrefix,
			Original,
			KeyValueType,
			&KeyValue{key, value, false},
			nil,
		}
	}

	return Line{
		lineNumber,
		lineString,
		spacePrefix,
		Original,
		OtherType,
		nil,
		nil,
	}
}

func getLineFromLineString(lineNumber int, line string) Line {
	return Line{
		LineNumber:    lineNumber,
		StringContent: line,
	}
}

func ParseIniFile(file string) ([]*Line, error) {
	fileContent, err := io.GetFileContents(file)
	if err != nil {
		return []*Line{}, err
	}

	parsedLines := []*Line{}
	for idx, line := range strings.Split(fileContent, "\n") {
		lineNumber := idx + 1

		parsedLine := parseLineString(lineNumber, line)
		parsedLines = append(parsedLines, &parsedLine)
	}

	return parsedLines, nil
}

func getSections(parsedLines []*Line) []*Section {
	sections := []*Section{}
	currentSection := &Section{GlobalSectionName, []*Line{}}

	for _, line := range parsedLines {
		if line.ContentType == SectionLineType {
			currentSection = &Section{line.SectionLine.SectionName, []*Line{}}
			sections = append(sections, currentSection)
		}
		currentSection.Lines = append(currentSection.Lines, line)
	}

	return sections
}

func getGlobalSection(sections []*Section) *Section {
	return GetSectionByName(sections, GlobalSectionName)
}
