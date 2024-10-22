package plugins

import (
	"einenlum/edicon/internal/core"
	"einenlum/edicon/internal/plugins/ini"
	"errors"
	"fmt"
)

type ConfigurationType int

const (
	Php = iota
)

func GetConfigurator(ctype ConfigurationType) (core.Configurator, error) {
	switch ctype {
	case Php:
		return ini.IniConfigurator{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("No configurator found for %d", ctype))
	}
}
