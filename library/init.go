package library

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/nanoteck137/tunebooklib/probe"
	"github.com/nanoteck137/tunebooklib/utils"
	"github.com/pelletier/go-toml/v2"
)

// TODO(patrik): Move?
func downloadImage(url, dir, name string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", fmt.Errorf("failed to parse media type: %w", err)
	}

	// TODO(patrik): Add more types?
	ext := ""
	switch mediaType {
	case "image/png":
		ext = ".png"
	case "image/jpeg":
		ext = ".jpeg"
	default:
		return "", fmt.Errorf("unsupported media type: %s", mediaType)
	}

	p := path.Join(dir, name+ext)
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to copy http body to file: %w", err)
	}

	return p, nil
}

var dateRegex = regexp.MustCompile(`^([12]\d\d\d)`)

func parseArtist(s string) []string {
	if s == "" {
		return []string{}
	}

	splits := strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == ';'
	})

	artists := make([]string, 0, len(splits))
	for _, s := range splits {
		a := strings.TrimSpace(s)

		if a != "" {
			artists = append(artists, a)
		}
	}

	return artists
}

type trackInfo struct {
	Name   string
	Artist string
	Number int
	Year   int
}

func getTrackInfo(p string) (trackInfo, error) {
	probe, err := probe.ProbeMedia(context.Background(), p)
	if err != nil {
		return trackInfo{}, err
	}

	var res trackInfo

	res.Name, _ = probe.Tags.GetString("title")
	res.Artist, _ = probe.Tags.GetString("artist")

	if tag, err := probe.Tags.GetString("date"); err == nil {
		match := dateRegex.FindStringSubmatch(tag)
		if len(match) > 0 {
			res.Year, _ = strconv.Atoi(match[1])
		}
	}

	if tag, err := probe.Tags.GetInt("track"); err == nil {
		res.Number = int(tag)
	}

	return res, nil
}

func createTrackFromFile(dir, filename string) AlbumTrack {
	p := filepath.Join(dir, filename)
	info, _ := getTrackInfo(p)

	if info.Name == "" {
		info.Name = strings.TrimSuffix(filename, filepath.Ext(filename))
	}

	if info.Number == 0 {
		info.Number = utils.ExtractNumber(filename)
	}

	artists := parseArtist(info.Artist)

	return AlbumTrack{
		Id:      utils.CreateTrackId(),
		File:    filename,
		Name:    info.Name,
		Number:  int64(info.Number),
		Year:    0,
		Tags:    []string{},
		Artists: artists,
	}
}

type InitializeAlbumParams struct {
	ZipFile string
}

func InitializeAlbum(dir string, params InitializeAlbumParams) error {
	if params.ZipFile != "" {
		fmt.Printf("Extracting tracks from %s...\n", params.ZipFile)
		err := extractZip(params.ZipFile, dir)
		if err != nil {
			return fmt.Errorf("init album: extract zip: %w", err)
		}
	}

	metadata := Album{}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("init album: read dir: %w", err)
	}

	var tracks []string
	var images []string

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()

		// Skip files starting wtih .
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}

		p := filepath.Join(dir, name)

		ext := filepath.Ext(p)

		if utils.IsValidTrackExt(ext) {
			tracks = append(tracks, p)
		}

		if utils.IsValidImageExt(ext) {
			images = append(images, p)
		}
	}

	if len(tracks) <= 0 {
		fmt.Println("No tracks found... Quitting")
		return nil
	}

	p := tracks[0]

	probe, err := probe.ProbeMedia(context.Background(), p)
	if err != nil {
		return fmt.Errorf("init album: probe track %s: %w", p, err)
	}

	metadata.Album.Type = AlbumTypeAlbum

	isSingle := len(tracks) == 1

	if isSingle {
		metadata.Album.Type = AlbumTypeSingle
	}

	metadata.Album.Id = utils.CreateAlbumId()

	if len(images) > 0 {
		// TODO(patrik): Better selection?
		metadata.General.Cover = images[0]
	}

	if !isSingle {
		metadata.Album.Name, _ = probe.Tags.GetString("album")
	} else {
		// NOTE(patrik): If we only have one track then we make the
		// album name the same as the track name
		metadata.Album.Name, _ = probe.Tags.GetString("title")
	}

	if !isSingle {
		if tag, err := probe.Tags.GetString("album_artist"); err == nil {
			metadata.Album.Artists = parseArtist(tag)
		} else {
			if tag, err := probe.Tags.GetString("artist"); err == nil {
				metadata.Album.Artists = parseArtist(tag)
			}
		}
	} else {
		if tag, err := probe.Tags.GetString("artist"); err == nil {
			metadata.Album.Artists = parseArtist(tag)
		}
	}

	if tag, err := probe.Tags.GetString("date"); err == nil {
		match := dateRegex.FindStringSubmatch(tag)
		if len(match) > 0 {
			metadata.General.Year, _ = strconv.ParseInt(match[1], 10, 64)
		}
	}

	for _, p := range tracks {
		filename := filepath.Base(p)
		fmt.Printf("Found track: %s\n", filename)

		metadata.Tracks = append(metadata.Tracks, createTrackFromFile(dir, filename))
	}

	data, err := toml.Marshal(&metadata)
	if err != nil {
		return fmt.Errorf("init album: marshal: %w", err)
	}

	// TODO(patrik): Move album.toml to constant
	out := filepath.Join(dir, "album.toml")
	err = os.WriteFile(out, data, 0644)
	if err != nil {
		return fmt.Errorf("init album: write file: %w", err)
	}

	return nil
}

type InitializeArtistParams struct {
	ArtistName string
	CoverUrl   string
}

func InitializeArtist(dir string, params InitializeArtistParams) error {
	// TODO(patrik): Add check for artist.toml already exists

	params.ArtistName = strings.TrimSpace(params.ArtistName)

	if params.ArtistName == "" {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("init artist: absolute path: %w", err)
		}

		params.ArtistName = filepath.Base(absDir)
	}

	cover := ""

	if params.CoverUrl != "" {
		// TODO(patrik): I want a better setup to download images,
		// where we also validate the images using magick
		p, err := downloadImage(params.CoverUrl, dir, "cover")
		if err != nil {
			return fmt.Errorf("init artist: download cover: %w", err)
		}

		cover = p
	}

	metadata := Artist{
		Id:         utils.CreateArtistId(),
		SearchName: utils.Slug(params.ArtistName),
		Name:       params.ArtistName,
		Cover:      cover,
		Tags:       []string{},
	}

	d, err := toml.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("init artist: marshal: %w", err)
	}

	// TODO(patrik): Move to constant
	p := path.Join(dir, artistFilename)
	err = os.WriteFile(p, d, 0644)
	if err != nil {
		return fmt.Errorf("init artist: write file: %w", err)
	}

	return nil
}

func InitializeLibrary(dir string) error {
	metadata := LibraryConfig{}

	d, err := toml.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("init library: marshal: %w", err)
	}

	// TODO(patrik): Move to constant
	p := path.Join(dir, "library.toml")
	err = os.WriteFile(p, d, 0644)
	if err != nil {
		return fmt.Errorf("init library: write file: %w", err)
	}

	return nil
}
