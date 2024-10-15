package cmd

import (
	"einenlum/edicon/internal/plugins/ini"
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an INI parameter",
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
		notationStyle := getNotationStyle(useBrackets)

		value, err := ini.GetIniParameterFromPath(notationStyle, file, key)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(value)
	},
}

func getNotationStyle(useBrackets bool) ini.NotationStyle {
	if useBrackets {
		return ini.BracketsNotation
	}

	return ini.DotNotation
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolP("brackets", "b", false, "Use brackts notation \"key[foo.bar]\" instead of dot notation")
}
