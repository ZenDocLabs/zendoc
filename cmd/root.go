package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zendoc",
	Short: "A CLI tool to generate and view documentation of Go project",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Errorf("error when executing cli %s", err)
		os.Exit(1)
	}
}
