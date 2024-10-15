package ini

type OutputType int

const (
	FullOutput OutputType = iota
	KeyValuesOnlyOutput
)

func OutputIniFile(iniFile *IniFile, outputType OutputType) string {
	output := ""

	for _, section := range *iniFile.Sections {
		for _, line := range section.Lines {
			output += line.ToString() + "\n"
		}
	}

	return output
}
