package cmd

import (
	"fmt"

	"github.com/einenlum/edicon/internal/core"
	"github.com/einenlum/edicon/internal/plugins"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a parameter",
	Long: `Something
Longer
`,
	Run: func(cmd *cobra.Command, args []string) {
		configurator, err := plugins.GetConfiguratorFromParentCmd(cmd.Parent())
		if err != nil {
			panic(err)
		}

		key, file := getGetCmdArguments(args)

		useBrackets, err := cmd.Flags().GetBool("brackets")
		if err != nil {
			fmt.Println(err)
		}
		notationStyle := core.GetNotationStyle(useBrackets)

		value, err := configurator.GetParameter(notationStyle, file, key)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(value)
	},
}

func getGetCmdArguments(args []string) (string, string) {
	if len(args) < 2 {
		panic("Not enough arguments")
	}

	key := args[0]
	file := args[1]

	return key, file
}

func init() {
	getCmd.Flags().BoolP("brackets", "b", false, "Use brackts notation \"key[foo.bar]\" instead of dot notation")
}
