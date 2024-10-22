package core

type Configuration interface{}

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
