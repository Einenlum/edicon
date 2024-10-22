package cmd

import (
	"github.com/spf13/cobra"
)

func InitCommonCommands(cmd *cobra.Command) {
	cmd.AddCommand(getCmd)
	cmd.AddCommand(setCmd)
}

func InitConfigCommands(cmd *cobra.Command) {
	cmd.AddCommand(phpCmd)
}
