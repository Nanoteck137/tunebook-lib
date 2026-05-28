package main

import (
	"log/slog"
	"os"

	"github.com/nanoteck137/tunebooklib/library"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use: "init",
}

var initLibraryCmd = &cobra.Command{
	Use: "library",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		err := library.InitializeLibrary(dir)
		if err != nil {
			slog.Error("failed to initialize library", "err", err)
			os.Exit(1)
		}
	},
}

var initAlbumCmd = &cobra.Command{
	Use: "album",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		err := library.InitializeAlbum(dir)
		if err != nil {
			slog.Error("failed to initialize album", "err", err)
			os.Exit(1)
		}
	},
}

var initArtistCmd = &cobra.Command{
	Use: "artist",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		artistName, _ := cmd.Flags().GetString("artist-name")
		coverUrl, _ := cmd.Flags().GetString("cover-url")

		err := library.InitializeArtist(dir, library.InitializeArtistParams{
			ArtistName: artistName,
			CoverUrl:   coverUrl,
		})
		if err != nil {
			slog.Error("failed to initialize artist", "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	initLibraryCmd.Flags().String("dir", ".", "directory to use")
	initLibraryCmd.MarkFlagDirname("dir")

	initAlbumCmd.Flags().String("dir", ".", "directory to use")
	initAlbumCmd.MarkFlagDirname("dir")

	initArtistCmd.Flags().String("dir", ".", "directory to use")
	initArtistCmd.Flags().String("artist-name", "", "set the artist name (when empty it uses the directory name)")
	initArtistCmd.Flags().String("cover-url", "", "url to image for downloading")
	initArtistCmd.MarkFlagDirname("dir")

	initCmd.AddCommand(initLibraryCmd, initAlbumCmd, initArtistCmd)

	rootCmd.AddCommand(initCmd)
}
