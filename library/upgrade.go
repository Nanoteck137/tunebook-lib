package library

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nanoteck137/tunebooklib/utils"
	"github.com/pelletier/go-toml/v2"
)

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		ext := filepath.Ext(f.Name)
		if !utils.IsValidTrackExt(ext) {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open file in zip: %w", err)
		}

		outPath := filepath.Join(dest, filepath.Base(f.Name))
		out, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			rc.Close()
			return fmt.Errorf("create output file: %w", err)
		}

		_, err = io.Copy(out, rc)
		rc.Close()
		out.Close()
		if err != nil {
			return fmt.Errorf("write file: %w", err)
		}

		fmt.Printf("Extracted: %s\n", filepath.Base(f.Name))
	}

	return nil
}

type UpgradeAlbumParams struct {
	Dir     string
	ZipFile string
}

func UpgradeAlbum(params UpgradeAlbumParams) error {
	if params.ZipFile != "" {
		fmt.Printf("Extracting tracks from %s...\n", params.ZipFile)
		err := extractZip(params.ZipFile, params.Dir)
		if err != nil {
			return fmt.Errorf("extract zip: %w", err)
		}
	}

	metadata, err := utils.ReadToml[Album](filepath.Join(params.Dir, albumFilename))
	if err != nil {
		return fmt.Errorf("read album: %w", err)
	}

	oldTracksByNumber := map[int64]*AlbumTrack{}
	for i := range metadata.Tracks {
		t := &metadata.Tracks[i]
		if t.Number != 0 {
			oldTracksByNumber[t.Number] = t
		}
	}

	entries, err := os.ReadDir(params.Dir)
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}

	var newTrackFiles []string
	for _, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		ext := filepath.Ext(e.Name())
		if utils.IsValidTrackExt(ext) {
			newTrackFiles = append(newTrackFiles, e.Name())
		}
	}

	sort.Strings(newTrackFiles)

	usedNumbers := map[int64]bool{}
	usedNames := map[string]bool{}
	var newTracks []AlbumTrack

	for _, filename := range newTrackFiles {
		stem := strings.TrimSuffix(filename, filepath.Ext(filename))
		num := int64(utils.ExtractNumber(filename))

		var matched *AlbumTrack

		if oldTrack, ok := oldTracksByNumber[num]; ok && num != 0 && !usedNumbers[num] {
			matched = oldTrack
			fmt.Printf("Matched by number %d: %s -> %s\n", num, oldTrack.File, filename)
		} else {
			for i := range metadata.Tracks {
				oldTrack := &metadata.Tracks[i]
				if usedNumbers[oldTrack.Number] || usedNames[oldTrack.File] {
					continue
				}
				oldStem := strings.TrimSuffix(oldTrack.File, filepath.Ext(oldTrack.File))
				if stem == oldStem {
					matched = oldTrack
					fmt.Printf("Matched by name: %s -> %s\n", oldTrack.File, filename)
					break
				}
			}
		}

		if matched != nil {
			t := *matched
			t.File = filename
			newTracks = append(newTracks, t)
			usedNumbers[t.Number] = true
			usedNames[matched.File] = true
		} else {
			p := filepath.Join(params.Dir, filename)
			info, _ := getTrackInfo(p)

			name := info.Name
			if name == "" {
				name = stem
			}

			artists := parseArtist(info.Artist)
			if len(artists) == 0 && len(metadata.Album.Artists) > 0 {
				artists = metadata.Album.Artists
			}

			year := int64(info.Year)
			if year == 0 {
				year = metadata.General.Year
			}

			t := AlbumTrack{
				Id:      utils.CreateTrackId(),
				File:    filename,
				Name:    name,
				Number:  int64(info.Number),
				Year:    year,
				Tags:    []string{},
				Artists: artists,
			}
			newTracks = append(newTracks, t)
			if t.Number != 0 {
				usedNumbers[t.Number] = true
			}
			fmt.Printf("New track: %s\n", filename)
		}
	}

	for _, oldTrack := range metadata.Tracks {
		if !usedNumbers[oldTrack.Number] && !usedNames[oldTrack.File] {
			fmt.Printf("Warning: old track %q (%s) has no matching new file\n", oldTrack.File, oldTrack.Name)
		}
	}

	metadata.Tracks = newTracks

	data, err := toml.Marshal(&metadata)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	out := filepath.Join(params.Dir, albumFilename)
	err = os.WriteFile(out, data, 0644)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
