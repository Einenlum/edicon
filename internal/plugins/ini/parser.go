package ini

import (
	"einenlum/edicon/internal/io"
	"strings"
)

func parseLineString(lineNumber int, lineString string) Line {
	// trim lineString
	lineString = strings.TrimSpace(lineString)

	// Check if line is empty
	if len(lineString) == 0 {
		return Line{
			lineNumber,
			lineString,
			Original,
			OtherType,
			nil,
			nil,
		}
	}

	// Check if line is a comment
	if strings.HasPrefix(lineString, ";") {
		return Line{
			lineNumber,
			lineString,
			Original,
			OtherType,
			nil,
			nil,
		}
	}

	// Check if line is a section
	if strings.HasPrefix(lineString, "[") && strings.HasSuffix(lineString, "]") {
		sectionName := strings.Trim(lineString, "[]")

		sectionLine := SectionLine{sectionName}

		return Line{
			lineNumber,
			lineString,
			Original,
			SectionLineType,
			nil,
			&sectionLine,
		}
	}

	// Check if line is a key value pair
	if strings.Contains(lineString, "=") {
		keyValue := strings.Split(lineString, "=")

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		return Line{
			lineNumber,
			lineString,
			Original,
			KeyValueType,
			&KeyValue{key, value, false},
			nil,
		}
	}

	return Line{
		lineNumber,
		lineString,
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
