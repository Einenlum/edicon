package ini

type LineContentType int

const GlobalSectionName = "__GLOBAL__"

const (
	KeyValueType LineContentType = iota
	SectionLineType
	OtherType
)

type Section struct {
	Name  string
	Lines []Line
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
	ContentType   LineContentType
	KeyValue      *KeyValue
	SectionLine   *SectionLine
}

type IniFile struct {
	Sections *[]Section
	FilePath string
}

func (line *Line) ToString() string {
	var result string

	if line.ContentType == KeyValueType {
		result = line.KeyValue.Key + "=" + line.KeyValue.Value
		// prepend comment symbol if line is commented
		if line.KeyValue.Commented {
			result = ";" + result
		}

		return result
	}

	if line.ContentType == SectionLineType {
		return "[" + line.SectionLine.SectionName + "]"
	}

	return line.StringContent
}
