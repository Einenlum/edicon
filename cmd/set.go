package cmd

import (
	"fmt"

	"github.com/einenlum/edicon/internal/core"
	"github.com/einenlum/edicon/internal/plugins"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a parameter",
	Long: `Something
Longer
`,
	Run: func(cmd *cobra.Command, args []string) {
		configurator, err := plugins.GetConfiguratorFromParentCmd(cmd.Parent())
		if err != nil {
			panic(err)
		}
		key, value, file := getSetCmdArguments(args)

		notationStyle := getNotationStyle(cmd)
		outputType := getOutputType(cmd)
		shouldOverwrite := shouldOverwrite(cmd)

		config, err := configurator.SetParameter(notationStyle, file, key, value)
		if err != nil {
			panic(err)
		}

		if shouldOverwrite {
			err = config.WriteToFile(file, outputType)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println(config.OutputFile(outputType))
		}
	},
}

func shouldOverwrite(cmd *cobra.Command) bool {
	overwrite, err := cmd.Flags().GetBool("write")
	if err != nil {
		panic(err)
	}

	return overwrite
}

func getOutputType(cmd *cobra.Command) core.OutputType {
	onlyValues, err := cmd.Flags().GetBool("only-values")
	if err != nil {
		panic(err)
	}

	if onlyValues {
		return core.MeaningFullOutput
	}

	return core.FullOutput
}

func getSetCmdArguments(args []string) (string, string, string) {
	if len(args) < 3 {
		panic("Not enough arguments")
	}

	key := args[0]
	value := args[1]
	file := args[2]

	return key, value, file
}

func init() {
	setCmd.Flags().BoolP("brackets", "b", false, "Use brackts notation \"key[foo.bar]\" instead of dot notation")
	setCmd.Flags().BoolP("write", "w", false, "Write the changes to the file")
	setCmd.Flags().BoolP("only-values", "o", false, "Only output the values (remove empty lines and comments)")
}
