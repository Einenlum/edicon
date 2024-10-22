package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var iniCmd = &cobra.Command{
	Use:   "ini",
	Short: "INI configuration",
	Long: `Something
Longer
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	InitCommonCommands(iniCmd)
}
