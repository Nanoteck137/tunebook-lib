package main

import (
	"log/slog"
	"os"

	"github.com/nanoteck137/tunebooklib/library"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use: "upgrade",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		zipFile, _ := cmd.Flags().GetString("zip")

		err := library.UpgradeAlbum(library.UpgradeAlbumParams{
			Dir:     dir,
			ZipFile: zipFile,
		})
		if err != nil {
			slog.Error("failed to upgrade album", "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	upgradeCmd.Flags().String("dir", ".", "directory containing album.toml and track files")
	upgradeCmd.MarkFlagDirname("dir")

	upgradeCmd.Flags().String("zip", "", "path to a zip file containing new tracks to extract")
	upgradeCmd.MarkFlagFilename("zip", "zip")

	rootCmd.AddCommand(upgradeCmd)
}
