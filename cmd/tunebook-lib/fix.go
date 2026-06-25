package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nanoteck137/tunebooklib/library"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use: "fix",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HELLO WORLD")
		dir, _ := cmd.Flags().GetString("dir")

		err := library.FixAlbumType(dir)
		if err != nil {
			slog.Error("failed to fix albums", "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	fixCmd.Flags().String("dir", ".", "directory containing albums to fix")
	fixCmd.MarkFlagDirname("dir")

	rootCmd.AddCommand(fixCmd)
}
