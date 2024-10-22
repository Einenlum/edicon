package ini

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/einenlum/edicon/internal/core"
	"github.com/einenlum/edicon/internal/io"

	"github.com/pmezard/go-difflib/difflib"
)

const (
	PHP_FILE_PATH                 = "../../../data/ini/php.ini"
	INI_FILE_PATH                 = "../../../data/ini/ini.ini"
	PHP_KEY_VALUES_ONLY_FILE_PATH = "../../../data/ini/php_key_values_only.ini"
	INI_KEY_VALUES_ONLY_FILE_PATH = "../../../data/ini/ini_key_values_only.ini"
)

func testParseSections(t *testing.T, filepath string, expectedSectionNames []string) {
	config, err := GetParsedIniFile(filepath)
	if err != nil {
		t.Fatal("Could not read the file", filepath, err.Error())
	}

	if len(config.Sections) != len(expectedSectionNames) {
		t.Fatal(
			fmt.Sprintf("Expected %d sections, got %d", len(expectedSectionNames), len(config.Sections)),
			config.Sections,
		)
	}

	sectionNames := []string{}
	for _, section := range config.Sections {
		sectionNames = append(sectionNames, section.Name)
	}

	if !reflect.DeepEqual(sectionNames, expectedSectionNames) {
		t.Fatal("Expected sections: ", expectedSectionNames, "Real sections:", sectionNames)
	}
}

func testParseSection(
	t *testing.T,
	filepath string,
	sectionName string,
	expectedLines int,
	expectedKeyValues int,
) {
	config, err := GetParsedIniFile(filepath)
	if err != nil {
		t.Fatal("Could not read the file", filepath, err.Error())
	}

	section := GetSectionByName(config.Sections, sectionName)
	if section == nil {
		t.Fatal("Could not find section", sectionName)
	}

	if len(section.Lines) != expectedLines {
		t.Fatal(fmt.Sprintf("Expected %d lines, got %d", expectedLines, len(section.Lines)))
	}

	keyValues := []*Line{}
	for _, line := range section.Lines {
		if line.ContentType == KeyValueType {
			keyValues = append(keyValues, line)
		}
	}

	if len(keyValues) != expectedKeyValues {
		t.Fatal(fmt.Sprintf("Expected %d key value lines, got %d", expectedKeyValues, len(keyValues)))
	}
}

func testPrintFullFile(
	t *testing.T,
	originalFilepath string,
	expectedFilepath string,
	outputType OutputType,
) {
	config, err := GetParsedIniFile(originalFilepath)
	if err != nil {
		t.Fatal("Could not read the file", originalFilepath, err.Error())
	}

	expectedContent, err := os.ReadFile(expectedFilepath)
	if err != nil {
		t.Fatal("Could not read the file", expectedFilepath, err.Error())
	}

	output := OutputConfigFile(&config, outputType)
	diffOutput, minusLines, plusLines := getDiff(cleanContent(expectedContent), removeEmptyTrailingLines(output))

	if len(minusLines) != 0 || len(plusLines) != 0 {
		t.Fatal("The diff should be empty. Diff: ", diffOutput)
	}
}

func testGetMissingParameter(
	t *testing.T,
	notationStyle core.NotationStyle,
	filepath string,
	missingKey string,
) {
	value, err := GetParameterFromPath(notationStyle, filepath, missingKey)
	if err == nil {
		t.Error("Should be missing. Got " + value + " instead")
	}
}

func testGetExistingParameter(
	t *testing.T,
	notationStyle core.NotationStyle,
	filepath string,
	key string,
	expectedValue string,
) {
	value, err := GetParameterFromPath(notationStyle, filepath, key)
	if err != nil {
		t.Error(err)
	}

	if expectedValue != value {
		t.Error(fmt.Sprintf("Expected %s got %s", expectedValue, value))
	}
}

func testEditExistingParameter(
	t *testing.T,
	notationStyle core.NotationStyle,
	filepath string,
	fullKey string,
	sectionName string,
	keyName string,
	newValue string,
	removedLine string,
	addedLine string,
) {
	fixturesContent, err := io.GetFileContents(filepath)
	if err != nil {
		t.Fatal(err)
	}

	config, err := EditConfigFile(notationStyle, filepath, fullKey, newValue)
	if err != nil {
		t.Fatal(err)
	}

	keyLine := getKeyLineBySectionName(config.Sections, sectionName, keyName)
	if keyLine == nil {
		t.Fatal(fmt.Sprintf("Could not find key %s in section %s", keyName, sectionName))
	}

	if newValue != keyLine.KeyValue.Value {
		t.Fatal(fmt.Sprintf("Expected %s got %s", newValue, keyLine.KeyValue.Value))
	}

	output := OutputConfigFile(config, FullOutput)
	diffOutput, minusLines, plusLines := getDiff(
		cleanContent([]byte(fixturesContent)),
		removeEmptyTrailingLines(output),
	)

	if len(minusLines) != 1 || len(plusLines) != 1 {
		t.Fatal("Expected one line to be added and one line to be removed. Diff: ", diffOutput)
	}

	if minusLines[0] != strings.TrimRight(removedLine, "\n") {
		t.Fatal(fmt.Sprintf("Expected removed line to be \"%s\" got \"%s\"", removedLine, minusLines[0]))
	}

	if plusLines[0] != addedLine {
		t.Fatal(fmt.Sprintf("Expected added line to be \"%s\" got \"%s\"", addedLine, plusLines[0]))
	}
}

func TestGetParsedIniFile(t *testing.T) {
	type TestElement struct {
		SectionName       string
		expectedLines     int
		expectedKeyValues int
	}

	t.Run("it parses php sections", func(t *testing.T) {
		testParseSections(t, PHP_FILE_PATH, []string{
			core.GLOBAL_SECTION_NAME,
			"PHP",
			"CLI Server",
			"Date",
			"mail function",
		})
	})

	t.Run("it parses ini sections", func(t *testing.T) {
		testParseSections(t, INI_FILE_PATH, []string{
			core.GLOBAL_SECTION_NAME,
			"user",
			"core",
			"alias",
			"push",
		})
	})

	dataProvider := []TestElement{
		{core.GLOBAL_SECTION_NAME, 2, 1},
		{"PHP", 19, 6},
		{"CLI Server", 3, 1},
		{"Date", 1, 0},
		{"mail function", 5, 2},
	}

	for _, element := range dataProvider {
		t.Run("it parses PHP: "+element.SectionName+" section", func(t *testing.T) {
			testParseSection(t, PHP_FILE_PATH, element.SectionName, element.expectedLines, element.expectedKeyValues)
		})
	}

	dataProvider = []TestElement{
		{core.GLOBAL_SECTION_NAME, 2, 1},
		{"user", 4, 2},
		{"core", 7, 4},
		{"alias", 6, 4},
		{"push", 5, 1},
	}

	for _, element := range dataProvider {
		t.Run("it parses INI: "+element.SectionName+" section", func(t *testing.T) {
			testParseSection(t, INI_FILE_PATH, element.SectionName, element.expectedLines, element.expectedKeyValues)
		})
	}
}

func TestPrintIniFile(t *testing.T) {
	t.Run("PHP: it prints the full ini file", func(t *testing.T) {
		testPrintFullFile(t, PHP_FILE_PATH, PHP_FILE_PATH, FullOutput)
	})

	t.Run("PHP: it prints the key values only", func(t *testing.T) {
		testPrintFullFile(t, PHP_FILE_PATH, PHP_KEY_VALUES_ONLY_FILE_PATH, KeyValuesOnlyOutput)
	})

	t.Run("INI: it prints the full ini file", func(t *testing.T) {
		testPrintFullFile(t, INI_FILE_PATH, INI_FILE_PATH, FullOutput)
	})

	t.Run("INI: it prints the key values only", func(t *testing.T) {
		testPrintFullFile(t, INI_FILE_PATH, INI_KEY_VALUES_ONLY_FILE_PATH, KeyValuesOnlyOutput)
	})
}

func TestGetParameter(t *testing.T) {
	phpMissingCases := []string{"PHP.not_a_real_key", "not_a_real_key", "Foobar.baz"}
	for _, key := range phpMissingCases {
		t.Run("PHP: it gets missing parameter "+key, func(t *testing.T) {
			testGetMissingParameter(t, core.DotNotation, PHP_FILE_PATH, key)
		})
	}

	iniMissingCases := []string{"user.not_a_real_key", "not_a_real_key", "Foobar.baz"}
	for _, key := range iniMissingCases {
		t.Run("INI: it gets missing parameter "+key, func(t *testing.T) {
			testGetMissingParameter(t, core.DotNotation, INI_FILE_PATH, key)
		})
	}

	validPhpDotCases := map[string]string{
		"orphan_key":              "value",
		"PHP.engine":              "On",
		"PHP.precision":           "14",
		"PHP.disable_classes":     "",
		"PHP.error_reporting":     "E_ALL & ~E_DEPRECATED & ~E_STRICT",
		"PHP.default_mimetype":    "\"text/html\"",
		"PHP.zend_extension":      "opcache",
		"mail function.SMTP":      "localhost",
		"mail function.smtp_port": "25",
	}

	for key, expectedValue := range validPhpDotCases {
		t.Run("PHP: it gets existing parameter "+key, func(t *testing.T) {
			testGetExistingParameter(t, core.DotNotation, PHP_FILE_PATH, key, expectedValue)
		})
	}

	validIniDotCases := map[string]string{
		"orphan_key":      "value",
		"user.name":       "User",
		"user.email":      "user@example.com",
		"core.editor":     "vim",
		"core.autocrlf":   "input",
		"core.fileMode":   "true",
		"core.ignoreCase": "false",
		"alias.co":        "checkout",
		"alias.br":        "branch",
		"alias.ci":        "commit",
		"alias.st":        "status",
		"push.default":    "simple",
	}

	for key, expectedValue := range validIniDotCases {
		t.Run("INI: it gets existing parameter "+key, func(t *testing.T) {
			testGetExistingParameter(t, core.DotNotation, INI_FILE_PATH, key, expectedValue)
		})
	}

	validPhpBracketsCases := map[string]string{
		"CLI Server[cli_server.color]": "On",
	}
	for key, expectedValue := range validPhpBracketsCases {
		t.Run("it gets existing parameter "+key, func(t *testing.T) {
			testGetExistingParameter(t, core.BracketsNotation, PHP_FILE_PATH, key, expectedValue)
		})
	}
}

func TestEditParameter(t *testing.T) {
	type EditTestElement struct {
		sectionName string
		keyName     string
		newValue    string
		removedLine string
		addedLine   string
	}

	phpDotCases := map[string]EditTestElement{
		"orphan_key": {
			core.GLOBAL_SECTION_NAME,
			"orphan_key",
			"foobar",
			"orphan_key = value",
			"orphan_key=foobar",
		},
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

	for fullKey, value := range phpDotCases {
		t.Run("PHP: it edits existing parameter "+fullKey, func(t *testing.T) {
			testEditExistingParameter(
				t,
				core.DotNotation,
				PHP_FILE_PATH,
				fullKey,
				value.sectionName,
				value.keyName,
				value.newValue,
				value.removedLine,
				value.addedLine,
			)
		})
	}

	phpBracketsCases := map[string]EditTestElement{
		"CLI Server[cli_server.color]": {"CLI Server", "cli_server.color", "black", "cli_server.color = On", "cli_server.color=black"},
	}

	for fullKey, value := range phpBracketsCases {
		t.Run("PHP: it edits existing parameter "+fullKey, func(t *testing.T) {
			testEditExistingParameter(
				t,
				core.BracketsNotation,
				PHP_FILE_PATH,
				fullKey,
				value.sectionName,
				value.keyName,
				value.newValue,
				value.removedLine,
				value.addedLine,
			)
		})
	}

	iniDotCases := map[string]EditTestElement{
		"orphan_key": {
			core.GLOBAL_SECTION_NAME,
			"orphan_key",
			"foobar",
			"orphan_key = value",
			"orphan_key=foobar",
		},
		"user.name": {
			"user",
			"name",
			"John",
			"name = User",
			"name=John",
		},
		"user.email": {
			"user",
			"email",
			"something@example.com",
			"email = user@example.com",
			"email=something@example.com",
		},
		"core.editor": {
			"core",
			"editor",
			"vscode",
			"editor = vim",
			"editor=vscode",
		},
		"core.autocrlf": {
			"core",
			"autocrlf",
			"true",
			"autocrlf = input",
			"autocrlf=true",
		},
		"core.fileMode": {
			"core",
			"fileMode",
			"false",
			"fileMode = true",
			"fileMode=false",
		},
		"core.ignoreCase": {
			"core",
			"ignoreCase",
			"true",
			"ignoreCase = false",
			"ignoreCase=true",
		},
		"alias.co": {
			"alias",
			"co",
			"command",
			"co = checkout",
			"co=command",
		},
		"alias.br": {
			"alias",
			"br",
			"foo",
			"br = branch",
			"br=foo",
		},
		"alias.ci": {
			"alias",
			"ci",
			"bar",
			"ci = commit",
			"ci=bar",
		},
		"alias.st": {
			"alias",
			"st",
			"baz",
			"st = status",
			"st=baz",
		},
		"push.default": {
			"push",
			"default",
			"upstream",
			"default = simple",
			"default=upstream",
		},
	}

	for fullKey, value := range iniDotCases {
		t.Run("INI: it edits existing parameter "+fullKey, func(t *testing.T) {
			testEditExistingParameter(
				t,
				core.DotNotation,
				INI_FILE_PATH,
				fullKey,
				value.sectionName,
				value.keyName,
				value.newValue,
				value.removedLine,
				value.addedLine,
			)
		})
	}
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
