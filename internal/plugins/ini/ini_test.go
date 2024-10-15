package ini

import (
	"einenlum/edicon/internal/notation"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetParsedIniFile(t *testing.T) {
	filePath := "../../../data/php.ini"

	iniFile, err := GetParsedIniFile(filePath)
	if err != nil {
		t.Error("Could not read the file", filePath, err.Error())
	}

	type TestElement struct {
		SectionName       string
		expectedLines     int
		expectedKeyValues int
	}

	t.Run("it parses sections", func(t *testing.T) {
		if len(*iniFile.Sections) != 4 {
			t.Error("Expected 4 sections, got", len(*iniFile.Sections))
		}

		sectionNames := []string{}
		for _, section := range *iniFile.Sections {
			sectionNames = append(sectionNames, section.Name)
		}

		expectedSectionNames := []string{"PHP", "CLI Server", "Date", "mail function"}
		if !reflect.DeepEqual(sectionNames, expectedSectionNames) {
			t.Error("Expected sections: ", expectedSectionNames, "Real sections:", sectionNames)
		}
	})

	dataProvider := []TestElement{
		{"PHP", 18, 6},
		{"CLI Server", 2, 1},
		{"Date", 1, 0},
		{"mail function", 4, 2},
	}

	for _, element := range dataProvider {
		section := GetSectionByName(iniFile.Sections, element.SectionName)

		t.Run("it parses "+element.SectionName+" lines", func(t *testing.T) {
			if len(section.Lines) != element.expectedLines {
				t.Error(fmt.Sprintf("Expected %d lines, got %d", element.expectedLines, len(section.Lines)))
			}
		})

		t.Run("it parses "+element.SectionName+" key values", func(t *testing.T) {
			keyValues := []Line{}
			for _, line := range section.Lines {
				if line.ContentType == KeyValueType {
					keyValues = append(keyValues, line)
				}
			}

			if len(keyValues) != element.expectedKeyValues {
				t.Error(fmt.Sprintf("Expected %d key value lines, got %d", element.expectedKeyValues, len(keyValues)))
			}
		})
	}
}

func TestPrintIniFile(t *testing.T) {
	iniFilePath := "../../../data/php.ini"

	iniFile, err := GetParsedIniFile(iniFilePath)
	if err != nil {
		t.Error("Could not read the file", iniFilePath, err.Error())
	}

	t.Run("it prints the full ini file", func(t *testing.T) {
		fullIniFileContent, err := os.ReadFile(iniFilePath)
		if err != nil {
			t.Error("Could not read the file", iniFilePath, err.Error())
		}

		output := OutputIniFile(&iniFile, FullOutput)
		diff := cmp.Diff(cleanContent(fullIniFileContent), removeEmptyTrailingLines(output))

		if diff != "" {
			t.Errorf("Mismatch (-expected +actual):\n%s", diff)
		}
	})

	t.Run("it prints the key values only", func(t *testing.T) {
		keyValuesFileContent, err := os.ReadFile("../../../data/php_key_values_only.ini")
		if err != nil {
			t.Error("Could not read the file", iniFilePath, err.Error())
		}

		output := OutputIniFile(&iniFile, KeyValuesOnlyOutput)

		diff := cmp.Diff(cleanContent(keyValuesFileContent), removeEmptyTrailingLines(output))

		if diff != "" {
			t.Errorf("Mismatch (-expected +actual):\n%s", diff)
		}
	})
}

func TestGetParameter(t *testing.T) {
	iniFilePath := "../../../data/php.ini"

	missingCases := []string{"PHP.not_a_real_key", "not_a_real_key", "Foobar.baz"}
	for _, key := range missingCases {
		t.Run("it tries to get missing parameter "+key, func(t *testing.T) {
			value, err := GetIniParameterFromPath(notation.DotNotation, iniFilePath, key)
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
		t.Run("it tries to get existing parameter "+key, func(t *testing.T) {
			value, err := GetIniParameterFromPath(notation.DotNotation, iniFilePath, key)
			if err != nil {
				t.Error(err)
			}

			if expectedValue != value {
				t.Error(fmt.Sprintf("Expected %s got %s", expectedValue, value))
			}
		})
	}

	t.Run("it tries to get existing parameter CLI Server[cli_server.color]", func(t *testing.T) {
		value, err := GetIniParameterFromPath(notation.BracketsNotation, iniFilePath, "CLI Server[cli_server.color]")
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
		"PHP.engine":              {"PHP", "engine", "Off"},
		"PHP.precision":           {"PHP", "precision", "140"},
		"PHP.disable_classes":     {"PHP", "disable_classes", "myclass"},
		"PHP.error_reporting":     {"PHP", "error_reporting", "E_ALL"},
		"PHP.default_mimetype":    {"PHP", "default_mimetype", "\"text/plain\""},
		"PHP.zend_extension":      {"PHP", "zend_extension", "opcache.so"},
		"mail function.SMTP":      {"mail function", "SMTP", "smtp.gmail.com"},
		"mail function.smtp_port": {"mail function", "smtp_port", "587"},
	}

	for key, values := range cases {
		sectionName := values[0]
		keyName := values[1]
		newValue := values[2]

		t.Run("it tries to edit existing parameter "+key, func(t *testing.T) {
			iniFile, err := EditIniFile(notation.DotNotation, iniFilePath, key, newValue)
			if err != nil {
				t.Error(err)
			}

			keyLine := getKeyLineBySectionName(iniFile.Sections, sectionName, keyName)
			if keyLine == nil {
				t.Error(fmt.Sprintf("Could not find key %s in section %s", keyName, sectionName))
			}

			if newValue != keyLine.KeyValue.Value {
				t.Error(fmt.Sprintf("Expected %s got %s", newValue, keyLine.KeyValue.Value))
			}
		})
	}

	t.Run("it tries to get existing parameter CLI Server[cli_server.color]", func(t *testing.T) {
		value, err := GetIniParameterFromPath(notation.BracketsNotation, iniFilePath, "CLI Server[cli_server.color]")
		if err != nil {
			t.Error(err)
		}

		if "On" != value {
			t.Error(fmt.Sprintf("Expected %s got %s", "On", value))
		}
	})
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
