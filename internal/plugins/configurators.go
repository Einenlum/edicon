package plugins

import (
	"einenlum/edicon/internal/core"
	"einenlum/edicon/internal/plugins/ini"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func GetConfiguratorFromParentCmd(parentCmd *cobra.Command) (core.Configurator, error) {
	if parentCmd == nil {
		return nil, errors.New("Parent command is nil")
	}

	parentConfigName := parentCmd.Use

	configurator, err := getConfigurator(parentConfigName)
	if err != nil {
		return nil, errors.New("Configurator not found")
	}

	return configurator, nil
}

func getConfigurator(ctype string) (core.Configurator, error) {
	switch ctype {
	case "php":
		return ini.IniConfigurator{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("No configurator found for %s", ctype))
	}
}
