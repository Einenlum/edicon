package ini

type OutputType int

const (
	FullOutput OutputType = iota
	KeyValuesOnlyOutput
)

func OutputIniFile(iniFile *IniFile, outputType OutputType) string {
	output := ""

	for _, section := range iniFile.Sections {
		output += "[" + section.Name + "]\n"

		for _, line := range section.Lines {
			if outputType != KeyValuesOnlyOutput {
				output += line.StringContent + "\n"
			} else {
				if line.ContentType == KeyValueType {
					output += line.StringContent + "\n"
				}
			}
		}
	}

	return output
}
