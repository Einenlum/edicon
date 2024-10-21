package cmd

import (
	"einenlum/edicon/internal/core"
	"einenlum/edicon/internal/plugins/ini"
	"fmt"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a parameter",
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
		notationStyle := core.GetNotationStyle(useBrackets)

		iniFile, err := ini.EditConfigFile(notationStyle, file, key, value)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(ini.OutputConfigFile(iniFile, ini.FullOutput))
	},
}

func init() {
	setCmd.Flags().BoolP("brackets", "b", false, "Use brackts notation \"key[foo.bar]\" instead of dot notation")
}
