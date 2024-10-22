package ini

import (
	"github.com/einenlum/edicon/internal/core"
	"github.com/einenlum/edicon/internal/io"
)

type LineStatus int

const (
	Original LineStatus = iota
	Changed
)

type LineContentType int

const GlobalSectionName = "__GLOBAL__"

const (
	KeyValueType LineContentType = iota
	SectionLineType
	OtherType
)

type Section struct {
	Name  string
	Lines []*Line
}

type KeyValue struct {
	Key       string
	Value     string
	Commented bool
}

type SectionLine struct {
	SectionName string
}

type Line struct {
	LineNumber    int
	StringContent string
	SpacePrefix   string
	Status        LineStatus
	ContentType   LineContentType
	KeyValue      *KeyValue
	SectionLine   *SectionLine
}

type IniConfiguration struct {
	Sections []*Section
	FilePath string
}

func (line *Line) SetValue(value string) {
	line.Status = Changed

	if line.ContentType != KeyValueType {
		return
	}
	line.KeyValue.Value = value
}

func (line *Line) ToString() string {
	if line.Status == Original {
		return line.StringContent
	}

	var result string

	if line.ContentType == KeyValueType {
		result = line.SpacePrefix + line.KeyValue.Key + "=" + line.KeyValue.Value
		// prepend comment symbol if line is commented
		if line.KeyValue.Commented {
			result = ";" + result
		}

		return result
	}

	if line.ContentType == SectionLineType {
		return line.SpacePrefix + "[" + line.SectionLine.SectionName + "]"
	}

	return line.StringContent
}

func (config *IniConfiguration) OutputFile(outputType core.OutputType) (string, error) {
	iniOutputType, err := toIniOutputType(outputType)
	if err != nil {
		return "", err
	}

	return OutputConfigFile(config, iniOutputType), nil
}

func (config *IniConfiguration) WriteToFile(filepath string, outputType core.OutputType) error {
	output, err := config.OutputFile(outputType)
	if err != nil {
		return err
	}

	err = io.WriteFileContents(filepath, output)

	return err
}

type IniConfigurator struct{}

func (configurator IniConfigurator) GetParameter(
	notationStyle core.NotationStyle,
	filePath string,
	key string,
) (string, error) {
	return GetParameterFromPath(notationStyle, filePath, key)
}

func (configurator IniConfigurator) SetParameter(
	notationStyle core.NotationStyle,
	filePath string,
	key string,
	value string,
) (core.Configuration, error) {
	config, err := EditConfigFile(notationStyle, filePath, key, value)
	if err != nil {
		return nil, err
	}

	return config, nil
}
