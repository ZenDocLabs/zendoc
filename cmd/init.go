package cmd

import (
	"github.com/dterbah/zendoc/internal/doc"
	"github.com/spf13/cobra"

	"github.com/fatih/color"
)

var initZenDoc = &cobra.Command{
	Use:   "init",
	Short: "Init a ZenDoc project",
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.InitZenDoc()
		if err != nil {
			color.Red("error when init zendoc on your project. Error: %s", err)
		} else {
			color.Green("ZenDoc project initialized !")
		}
	},
}

func init() {
	rootCmd.AddCommand(initZenDoc)
}
