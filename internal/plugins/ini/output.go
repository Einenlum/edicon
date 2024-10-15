package ini

type OutputType int

const (
	FullOutput OutputType = iota
	KeyValuesOnlyOutput
)

func OutputIniFile(iniFile *IniFile, outputType OutputType) string {
	output := ""

	shouldBePrinted := func(line Line) bool {
		if outputType != KeyValuesOnlyOutput {
			return true
		}

		return line.ContentType != OtherType
	}

	for _, section := range *iniFile.Sections {
		for _, line := range *section.Lines {
			if shouldBePrinted(line) {
				output += line.ToString() + "\n"
			}
		}
	}

	return output
}
