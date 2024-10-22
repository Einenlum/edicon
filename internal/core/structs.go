package core

type Configuration interface {
	OutputFile(outputType int) (string, error)

	WriteToFile(filepath string, outputType int) error
}

type Configurator interface {
	GetParameter(
		notationStyle NotationStyle,
		filePath string,
		key string,
	) (string, error)

	SetParameter(
		notationStyle NotationStyle,
		filePath string,
		key string,
		value string,
	) (*Configuration, error)
}
