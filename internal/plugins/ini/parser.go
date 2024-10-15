package ini

import (
	"einenlum/edicon/internal/io"
	"einenlum/edicon/internal/notation"
	"regexp"
	"strings"
)

func parseLineString(lineNumber int, lineString string) Line {
	// Check if line is empty
	if len(lineString) == 0 {
		return Line{lineNumber, lineString, OtherType, nil, nil}
	}

	// Check if line is a comment
	if strings.HasPrefix(lineString, ";") {
		return Line{lineNumber, lineString, OtherType, nil, nil}
	}

	// Check if line is a section
	if strings.HasPrefix(lineString, "[") && strings.HasSuffix(lineString, "]") {
		sectionName := strings.Trim(lineString, "[]")

		sectionLine := SectionLine{sectionName}

		return Line{lineNumber, lineString, SectionLineType, nil, &sectionLine}
	}

	// Check if line is a key value pair
	if strings.Contains(lineString, "=") {
		keyValue := strings.Split(lineString, "=")

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		return Line{lineNumber, lineString, KeyValueType, &KeyValue{key, value, false}, nil}
	}

	return Line{lineNumber, lineString, OtherType, nil, nil}
}

func getLineFromLineString(lineNumber int, line string) Line {
	return Line{
		LineNumber:    lineNumber,
		StringContent: line,
	}
}

func ParseIniFile(file string) (*[]Line, error) {
	fileContent, err := io.GetFileContents(file)
	if err != nil {
		return &[]Line{}, err
	}

	parsedLines := []Line{}
	for idx, line := range strings.Split(fileContent, "\n") {
		lineNumber := idx + 1

		parsedLine := parseLineString(lineNumber, line)
		parsedLines = append(parsedLines, parsedLine)
	}

	return &parsedLines, nil
}

func getSections(parsedLines *[]Line) []Section {
	sections := []Section{}
	currentSection := Section{GlobalSectionName, []Line{}}

	for _, line := range *parsedLines {
		if line.ContentType == SectionLineType {
			if len(currentSection.Lines) > 0 {
				sections = append(sections, currentSection)
			}

			currentSection = Section{line.SectionLine.SectionName, []Line{}}
		} else {
			currentSection.Lines = append(currentSection.Lines, line)
		}
	}

	sections = append(sections, currentSection)

	return sections
}

func DecomposeKey(notationStyle notation.NotationStyle, key string) []string {
	if notationStyle == notation.DotNotation {
		return DecomposeKeyWithDotNotation(key)
	}

	return DecomposeKeyWithBracketNotation(key)
}

func DecomposeKeyWithBracketNotation(key string) []string {
	re := regexp.MustCompile(`([^\[\]]+)|\[([^\[\]]+)\]`)
	matches := re.FindAllStringSubmatch(key, -1)

	var result []string
	for _, match := range matches {
		// Add the first part (outside the brackets) or any part inside the brackets
		if match[1] != "" {
			result = append(result, match[1])
		} else if match[2] != "" {
			result = append(result, match[2])
		}
	}

	return result
}

func DecomposeKeyWithDotNotation(key string) []string {
	if !strings.Contains(key, ".") {
		return []string{key}
	}

	keyParts := strings.Split(key, ".")

	return keyParts
}

func getGlobalSection(sections *[]Section) *Section {
	return GetSectionByName(sections, GlobalSectionName)
}
