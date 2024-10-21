package ini

import (
	"einenlum/edicon/internal/core"
	"einenlum/edicon/internal/io"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	// "github.com/google/go-cmp/cmp"
	// "github.com/sergi/go-diff/diffmatchpatch"
	"github.com/pmezard/go-difflib/difflib"
)

func TestGetParsedIniFile(t *testing.T) {
	filePath := "../../../data/php.ini"

	iniFile, err := GetParsedIniFile(filePath)
	if err != nil {
		t.Fatal("Could not read the file", filePath, err.Error())
	}

	type TestElement struct {
		SectionName       string
		expectedLines     int
		expectedKeyValues int
	}

	t.Run("it parses sections", func(t *testing.T) {
		if len(iniFile.Sections) != 4 {
			t.Fatal("Expected 4 sections, got", len(iniFile.Sections))
		}

		sectionNames := []string{}
		for _, section := range iniFile.Sections {
			sectionNames = append(sectionNames, section.Name)
		}

		expectedSectionNames := []string{"PHP", "CLI Server", "Date", "mail function"}
		if !reflect.DeepEqual(sectionNames, expectedSectionNames) {
			t.Fatal("Expected sections: ", expectedSectionNames, "Real sections:", sectionNames)
		}
	})

	dataProvider := []TestElement{
		{"PHP", 19, 6},
		{"CLI Server", 3, 1},
		{"Date", 1, 0},
		{"mail function", 5, 2},
	}

	for _, element := range dataProvider {
		section := GetSectionByName(iniFile.Sections, element.SectionName)
		if section == nil {
			t.Fatal("Could not find section", element.SectionName)
		}

		t.Run("it parses "+element.SectionName+" lines", func(t *testing.T) {
			if len(section.Lines) != element.expectedLines {
				t.Fatal(fmt.Sprintf("Expected %d lines, got %d", element.expectedLines, len(section.Lines)))
			}
		})

		t.Run("it parses "+element.SectionName+" key values", func(t *testing.T) {
			keyValues := []*Line{}
			for _, line := range section.Lines {
				if line.ContentType == KeyValueType {
					keyValues = append(keyValues, line)
				}
			}

			if len(keyValues) != element.expectedKeyValues {
				t.Fatal(fmt.Sprintf("Expected %d key value lines, got %d", element.expectedKeyValues, len(keyValues)))
			}
		})
	}
}

func TestPrintIniFile(t *testing.T) {
	iniFilePath := "../../../data/php.ini"

	iniFile, err := GetParsedIniFile(iniFilePath)
	if err != nil {
		t.Fatal("Could not read the file", iniFilePath, err.Error())
	}

	t.Run("it prints the full ini file", func(t *testing.T) {
		fullIniFileContent, err := os.ReadFile(iniFilePath)
		if err != nil {
			t.Fatal("Could not read the file", iniFilePath, err.Error())
		}

		output := OutputConfigFile(&iniFile, FullOutput)
		diffOutput, minusLines, plusLines := getDiff(cleanContent(fullIniFileContent), removeEmptyTrailingLines(output))

		if len(minusLines) != 0 || len(plusLines) != 0 {
			t.Fatal("The diff should be empty. Diff: ", diffOutput)
		}
	})

	t.Run("it prints the key values only", func(t *testing.T) {
		keyValuesFileContent, err := os.ReadFile("../../../data/php_key_values_only.ini")
		if err != nil {
			t.Error("Could not read the file", iniFilePath, err.Error())
		}

		output := OutputConfigFile(&iniFile, KeyValuesOnlyOutput)

		diffOutput, minusLines, plusLines := getDiff(cleanContent(keyValuesFileContent), removeEmptyTrailingLines(output))

		if len(minusLines) != 0 || len(plusLines) != 0 {
			t.Fatal("The diff should be empty. Diff: ", diffOutput)
		}
	})
}

func TestGetParameter(t *testing.T) {
	iniFilePath := "../../../data/php.ini"

	missingCases := []string{"PHP.not_a_real_key", "not_a_real_key", "Foobar.baz"}
	for _, key := range missingCases {
		t.Run("it gets missing parameter "+key, func(t *testing.T) {
			value, err := GetParameterFromPath(core.DotNotation, iniFilePath, key)
			if err == nil {
				t.Error("Got " + value + " instead")
			}
		})
	}

	validCases := map[string]string{
		"PHP.engine":              "On",
		"PHP.precision":           "14",
		"PHP.disable_classes":     "",
		"PHP.error_reporting":     "E_ALL & ~E_DEPRECATED & ~E_STRICT",
		"PHP.default_mimetype":    "\"text/html\"",
		"PHP.zend_extension":      "opcache",
		"mail function.SMTP":      "localhost",
		"mail function.smtp_port": "25",
	}

	for key, expectedValue := range validCases {
		t.Run("it gets existing parameter "+key, func(t *testing.T) {
			value, err := GetParameterFromPath(core.DotNotation, iniFilePath, key)
			if err != nil {
				t.Error(err)
			}

			if expectedValue != value {
				t.Error(fmt.Sprintf("Expected %s got %s", expectedValue, value))
			}
		})
	}

	t.Run("it gets existing parameter CLI Server[cli_server.color]", func(t *testing.T) {
		value, err := GetParameterFromPath(core.BracketsNotation, iniFilePath, "CLI Server[cli_server.color]")
		if err != nil {
			t.Error(err)
		}

		if "On" != value {
			t.Error(fmt.Sprintf("Expected %s got %s", "On", value))
		}
	})
}

func TestEditParameter(t *testing.T) {
	iniFilePath := "../../../data/php.ini"

	cases := map[string][]string{
		"PHP.engine": {
			"PHP",
			"engine",
			"Off",
			"engine = On",
			"engine=Off",
		},
		"PHP.precision": {
			"PHP",
			"precision",
			"140",
			"precision = 14",
			"precision=140",
		},
		"PHP.disable_classes": {
			"PHP",
			"disable_classes",
			"myclass",
			"disable_classes =",
			"disable_classes=myclass",
		},
		"PHP.error_reporting": {
			"PHP",
			"error_reporting",
			"E_ALL",
			"error_reporting = E_ALL & ~E_DEPRECATED & ~E_STRICT",
			"error_reporting=E_ALL",
		},
		"PHP.default_mimetype": {
			"PHP",
			"default_mimetype",
			"\"text/plain\"",
			"default_mimetype = \"text/html\"",
			"default_mimetype=\"text/plain\"",
		},
		"PHP.zend_extension": {
			"PHP",
			"zend_extension",
			"opcache.so",
			"zend_extension=opcache",
			"zend_extension=opcache.so",
		},
		"mail function.SMTP": {
			"mail function",
			"SMTP",
			"smtp.gmail.com",
			"SMTP = localhost",
			"SMTP=smtp.gmail.com",
		},
		"mail function.smtp_port": {
			"mail function",
			"smtp_port",
			"587",
			"smtp_port = 25",
			"smtp_port=587",
		},
	}

	fixturesIniFile, err := io.GetFileContents("../../../data/php.ini")
	if err != nil {
		t.Fatal(err)
	}

	for key, values := range cases {
		sectionName := values[0]
		keyName := values[1]
		newValue := values[2]
		removedLine := values[3]
		expectedLine := values[4]

		t.Run("it edits existing parameter "+key, func(t *testing.T) {
			iniFile, err := EditConfigFile(core.DotNotation, iniFilePath, key, newValue)
			if err != nil {
				t.Fatal(err)
			}

			keyLine := getKeyLineBySectionName(iniFile.Sections, sectionName, keyName)
			if keyLine == nil {
				t.Fatal(fmt.Sprintf("Could not find key %s in section %s", keyName, sectionName))
			}

			if newValue != keyLine.KeyValue.Value {
				t.Fatal(fmt.Sprintf("Expected %s got %s", newValue, keyLine.KeyValue.Value))
			}

			output := OutputConfigFile(iniFile, FullOutput)
			diffOutput, minusLines, plusLines := getDiff(
				cleanContent([]byte(fixturesIniFile)),
				removeEmptyTrailingLines(output),
			)

			if len(minusLines) != 1 || len(plusLines) != 1 {
				t.Fatal("Expected one line to be added and one line to be removed. Diff: ", diffOutput)
			}

			if minusLines[0] != strings.TrimRight(removedLine, "\n") {
				t.Fatal(fmt.Sprintf("Expected removed line to be \"%s\" got \"%s\"", removedLine, minusLines[0]))
			}

			if plusLines[0] != expectedLine {
				t.Fatal(fmt.Sprintf("Expected plus line to be \"%s\" got \"%s\"", expectedLine, plusLines[0]))
			}
		})
	}

	t.Run("it edits existing parameter CLI Server[cli_server.color]", func(t *testing.T) {
		iniFile, err := EditConfigFile(core.BracketsNotation, iniFilePath, "CLI Server[cli_server.color]", "black")
		if err != nil {
			t.Fatal(err)
		}

		keyLine := getKeyLineBySectionName(iniFile.Sections, "CLI Server", "cli_server.color")
		if keyLine == nil {
			t.Fatal("Could not find key cli_server.color in section CLI Server")
		}

		if "black" != keyLine.KeyValue.Value {
			t.Fatal(fmt.Sprintf("Expected \"black\" got %s", keyLine.KeyValue.Value))
		}
	})
}

func getLinesStartingWith(lines []string, prefix string) []string {
	filteredLines := []string{}

	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			filteredLines = append(filteredLines, line)
		}
	}

	return filteredLines
}

func cleanContent(output []byte) string {
	stringOutput := string(output)
	stringOutput = removeEmptyTrailingLines(stringOutput)

	return stringOutput
}

// I had to add this to avoid dealing with weird trailing empty lines
func removeEmptyTrailingLines(output string) string {
	lastNonEmptyLineNumber := 0

	for i := len(output) - 1; i >= 0; i-- {
		if output[i] != '\n' {
			lastNonEmptyLineNumber = i
			break
		}
	}

	return strings.TrimRight(output[:lastNonEmptyLineNumber+1], "\n")
}

func getDiff(expected string, actual string) (
	diffOutput string,
	minusLines []string,
	plusLines []string,
) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(expected),
		B:        difflib.SplitLines(actual),
		Context:  0,
		FromFile: "Expected",
		ToFile:   "Actual",
	}
	diffOutput, _ = difflib.GetUnifiedDiffString(diff)
	lines := difflib.SplitLines(diffOutput)
	if len(lines) < 3 {
		return diffOutput, minusLines, plusLines
	}
	lines = lines[3:]

	for _, line := range lines {
		if strings.HasPrefix(line, "-") {
			minusLines = append(
				minusLines,
				trimLine(strings.TrimLeft(line, "-")),
			)
		} else if strings.HasPrefix(line, "+") {
			plusLines = append(
				plusLines,
				trimLine(strings.TrimLeft(line, "+")),
			)
		}
	}

	return diffOutput, minusLines, plusLines
}

func trimLine(line string) string {
	trimmed := strings.TrimRight(line, "\n")
	trimmed = strings.TrimSpace(trimmed)

	return trimmed
}
