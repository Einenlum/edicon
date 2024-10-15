/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getCmd.Flags().BoolP("brackets", "b", false, "Use brackts notation \"key[foo.bar]\" instead of dot notation")
}
