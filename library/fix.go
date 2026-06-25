package library

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"io"

	"github.com/nanoteck137/tunebooklib/utils"
	"github.com/pelletier/go-toml/v2"
)

func FixAlbumType(dir string) error {
	return filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Name() != albumFilename {
			return nil
		}

		metadata, err := utils.ReadToml[Album](p)
		if err != nil {
			return fmt.Errorf("failed to read album: %w", err)
		}

		if len(metadata.Tracks) == 1 {
			metadata.Album.Type = AlbumTypeSingle
		} else {
			metadata.Album.Type = AlbumTypeAlbum
		}

		backup := p + ".bak"
		err = copyFile(p, backup)
		if err != nil {
			return fmt.Errorf("failed to backup album: %w", err)
		}

		data, err := toml.Marshal(&metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal album: %w", err)
		}

		err = os.WriteFile(p, data, 0644)
		if err != nil {
			return fmt.Errorf("failed to write album: %w", err)
		}

		fmt.Printf("Fixed: %s (type=%s, backup=%s)\n", p, metadata.Album.Type, backup)

		return nil
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Close()
}
