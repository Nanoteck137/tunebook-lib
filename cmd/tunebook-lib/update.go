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
		noWarnings, _ := cmd.Flags().GetBool("no-warnings")

		lib, err := library.ProcessMusicLibrary(dir, library.UpdateLibraryOptions{
			SuppressWarnings: noWarnings,
		})
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

	updateCmd.Flags().Bool("no-warnings", false, "suppress displaying individual warnings")

	rootCmd.AddCommand(updateCmd)
}
