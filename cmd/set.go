package cmd

import (
	"einenlum/edicon/internal/notation"
	"einenlum/edicon/internal/plugins/ini"
	"fmt"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set an INI parameter",
	Long: `Something
Longer
`,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]
		file := args[2]

		useBrackets, err := cmd.Flags().GetBool("brackets")
		if err != nil {
			fmt.Println(err)
		}
		notationStyle := notation.GetNotationStyle(useBrackets)

		iniFile, err := ini.EditIniFile(notationStyle, file, key, value)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(ini.OutputIniFile(iniFile, ini.FullOutput))
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().BoolP("brackets", "b", false, "Use brackts notation \"key[foo.bar]\" instead of dot notation")
}
