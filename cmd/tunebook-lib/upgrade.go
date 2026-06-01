package main

import (
	"log/slog"
	"os"

	"github.com/nanoteck137/tunebooklib/library"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use: "upgrade <SOURCE_DIR>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]

		target, _ := cmd.Flags().GetString("target")
		zipFile, _ := cmd.Flags().GetString("zip")

		err := library.UpgradeAlbum(library.UpgradeAlbumParams{
			SourceDir: dir,
			ZipFile:   zipFile,
			TargetDir: target,
		})
		if err != nil {
			slog.Error("failed to upgrade album", "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	upgradeCmd.Flags().String("target", ".", "target directory to write upgraded album.toml and scan for new tracks (defaults to --dir)")
	upgradeCmd.MarkFlagDirname("target")

	upgradeCmd.Flags().String("zip", "", "path to a zip file containing new tracks to extract")
	upgradeCmd.MarkFlagFilename("zip", "zip")

	rootCmd.AddCommand(upgradeCmd)
}
