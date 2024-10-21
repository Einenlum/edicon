package cmd

import (
	"einenlum/edicon/internal/core"
	"einenlum/edicon/internal/plugins/ini"
	"fmt"

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
		key := args[0]
		file := args[1]

		useBrackets, err := cmd.Flags().GetBool("brackets")
		if err != nil {
			fmt.Println(err)
		}
		notationStyle := core.GetNotationStyle(useBrackets)

		value, err := ini.GetParameterFromPath(notationStyle, file, key)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(value)
	},
}

func init() {
	getCmd.Flags().BoolP("brackets", "b", false, "Use brackts notation \"key[foo.bar]\" instead of dot notation")
}
