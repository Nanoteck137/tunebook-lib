package main

import (
	"log/slog"
	"os"

	"github.com/nanoteck137/tunebooklib/library"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		lib, err := library.ProcessMusicLibrary(dir)
		if err != nil {
			slog.Error("failed to fetch library", "err", err)
			os.Exit(1)
		}

		err = lib.WriteToDisk()
		if err != nil {
			slog.Error("failed to write the processed library to disk", "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	updateCmd.Flags().StringP("dir", "d", ".", "The directory to update")
	updateCmd.MarkFlagDirname("dir")

	rootCmd.AddCommand(updateCmd)
}
