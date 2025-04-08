package cmd

import "github.com/spf13/cobra"

var generateZenDoc = &cobra.Command{
	Use:   "generate",
	Short: "Generate doc for the current go project",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(generateZenDoc)
}
