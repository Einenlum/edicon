package cmd

import (
	"github.com/einenlum/edicon/internal/core"
	"github.com/spf13/cobra"
)

func getNotationStyle(cmd *cobra.Command) core.NotationStyle {
	useBrackets, err := cmd.Flags().GetBool("brackets")
	if err != nil {
		panic(err)
	}
	return core.GetNotationStyle(useBrackets)
}

func InitCommonCommands(cmd *cobra.Command) {
	cmd.AddCommand(getCmd)
	cmd.AddCommand(setCmd)
}

func InitConfigCommands(cmd *cobra.Command) {
	cmd.AddCommand(iniCmd)
	cmd.AddCommand(phpCmd)
}
