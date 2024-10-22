package core

// As crazy as it sounds I couldn't find for now a way to store
// lines in a global section with my current implementation.
// I use a very unplausable name for the global section to avoid
// conflicts with real sections.
const GLOBAL_SECTION_NAME = "EDICON_GLOBAL21162041383223945417"

type Configuration interface {
	OutputFile(outputType OutputType) (string, error)

	WriteToFile(filepath string, outputType OutputType) error
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
	) (Configuration, error)
}
