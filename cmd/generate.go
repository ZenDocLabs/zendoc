package cmd

import (
	"os"

	"github.com/dterbah/zendoc/internal"
	"github.com/dterbah/zendoc/internal/doc/generate"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var generateZenDoc = &cobra.Command{
	Use:   "generate [output]",
	Short: "Generate doc for the current go project",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			color.Red("Missing output format. Expected 'json' or 'web'")
			cmd.Usage()
			os.Exit(1)
		}

		format := args[0]
		if format != internal.JSON_EXPORT_TYPE && format != internal.WEB_EXPORT_TYPE {
			color.Red("Invalid output format. Must be '%s' or '%s'", internal.WEB_EXPORT_TYPE, internal.JSON_EXPORT_TYPE)
			cmd.Usage()
			os.Exit(1)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		outputFormat := args[0]
		err := generate.GenerateDoc(outputFormat)

		if err != nil {
			color.Red("error when generating doc %s", err)
			os.Exit(1)
		} else {
			color.Green("Documentation exported !")
		}
	},
}

func init() {
	rootCmd.AddCommand(generateZenDoc)
}
