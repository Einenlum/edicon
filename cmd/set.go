package cmd

import (
	"fmt"

	"github.com/einenlum/edicon/internal/core"
	"github.com/einenlum/edicon/internal/plugins"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a parameter",
	Long: `Something
Longer
`,
	Run: func(cmd *cobra.Command, args []string) {
		configurator, err := plugins.GetConfiguratorFromParentCmd(cmd.Parent())
		if err != nil {
			panic(err)
		}
		key, value, file := getSetCmdArguments(args)

		useBrackets, err := cmd.Flags().GetBool("brackets")
		if err != nil {
			panic(err)
		}
		notationStyle := core.GetNotationStyle(useBrackets)

		config, err := configurator.SetParameter(notationStyle, file, key, value)
		if err != nil {
			panic(err)
		}

		fmt.Println(config.OutputFile(core.FullOutput))
	},
}

func getSetCmdArguments(args []string) (string, string, string) {
	if len(args) < 3 {
		panic("Not enough arguments")
	}

	key := args[0]
	value := args[1]
	file := args[2]

	return key, value, file
}

func init() {
	setCmd.Flags().BoolP("brackets", "b", false, "Use brackts notation \"key[foo.bar]\" instead of dot notation")
}
