package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var phpCmd = &cobra.Command{
	Use:   "php",
	Short: "PHP INI configuration",
	Long: `Something
Longer
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	InitCommonCommands(phpCmd)
}
