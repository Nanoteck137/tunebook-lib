package main

import (
	"os"

	tunebooklib "github.com/nanoteck137/tunebooklib"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     tunebooklib.AppName,
	Version: tunebooklib.Version,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	rootCmd.SetVersionTemplate(tunebooklib.VersionTemplate(tunebooklib.AppName))
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
