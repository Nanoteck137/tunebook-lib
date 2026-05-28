package library

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/maruel/natural"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/tunebooklib/timer"
	"github.com/nanoteck137/tunebooklib/utils"
)

func validateImage(p string) error {
	ext := filepath.Ext(p)
	if !utils.IsValidImageExt(ext) {
		return errors.New("image not valid file extention: " + ext)
	}

	return nil
}

func dedupStringArr(arr []string) []string {
	seen := map[string]bool{}
	res := make([]string, 0, len(arr))

	for _, value := range arr {
		value = anvil.String(value)
		if value == "" {
			continue
		}

		if !seen[value] {
			seen[value] = true
			res = append(res, value)
		}
	}

	return res
}

func transformTagsArr(tags []string) []string {
	for i, tag := range tags {
		tags[i] = utils.Slug(strings.TrimSpace(tag))
	}

	return dedupStringArr(tags)
}

func processArtistMetadata(artist *Artist) {
	artist.Name = anvil.String(artist.Name)
	artist.SearchName = utils.Slug(artist.SearchName)

	artist.Tags = transformTagsArr(artist.Tags)
}

func processAlbumMetadata(metadata *Album) {
	album := &metadata.Album

	album.Name = anvil.String(album.Name)

	if album.Year == 0 {
		album.Year = metadata.General.Year
	}

	album.Artists = dedupStringArr(album.Artists)

	album.Tags = append(album.Tags, metadata.General.Tags...)
	album.Tags = transformTagsArr(album.Tags)

	for i := range metadata.Tracks {
		t := &metadata.Tracks[i]

		if t.Year == 0 {
			t.Year = metadata.General.Year
		}

		t.Name = anvil.String(t.Name)

		t.Tags = append(t.Tags, metadata.General.Tags...)
		t.Tags = append(t.Tags, metadata.General.TrackTags...)
		t.Tags = transformTagsArr(t.Tags)

		t.Artists = dedupStringArr(t.Artists)
	}
}

func ReadLibraryConfig(dir string) (LibraryConfig, error) {
	var res LibraryConfig

	p := filepath.Join(dir, libraryMetadataFilename)
	res, err := utils.ReadToml[LibraryConfig](p)
	if err == nil {
		res.Path = dir
		return res, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		p, err := FindFile(dir, libraryMetadataFilename)
		if err != nil {
			return LibraryConfig{}, err
		}

		res, err := utils.ReadJson[LibraryConfig](p)
		if err != nil {
			return LibraryConfig{}, err
		}

		res.Path = filepath.Dir(p)

		return res, nil
	}

	return LibraryConfig{}, err
}

// TODO(patrik): Move to utils
func FindFile(dir, filename string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	for {
		p := filepath.Join(dir, filename)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("%s not found in any parent directory", filename)
		}

		dir = parent
	}
}

type fetchResult struct {
	artists []Artist
	albums  []Album
}

func checkIfExcluded(relPath string, excludeDirs []string) bool {
	for _, exclude := range excludeDirs {
		if relPath == exclude {
			return true
		}

		if strings.HasPrefix(relPath, exclude+string(filepath.Separator)) {
			return true
		}
	}

	return false
}

func fetch(dir string, excludeDirs []string, reporter *Reporter) (*fetchResult, error) {
	res := &fetchResult{
		artists: []Artist{},
		albums:  []Album{},
	}

	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if d == nil {
			return nil
		}

		if d.IsDir() {
			relPath, err := filepath.Rel(dir, p)
			if err != nil {
				return err
			}

			if checkIfExcluded(relPath, excludeDirs) {
				return filepath.SkipDir
			}

			return nil
		}

		name := d.Name()

		if strings.HasPrefix(name, ".") {
			return nil
		}

		switch name {
		case artistMetadataFilename:
			artist, err := utils.ReadToml[Artist](p)
			if err != nil {
				reporter.AddWarning(p, fmt.Errorf("failed to read artist: %w", err))
				return nil
			}

			p, err := filepath.Rel(dir, filepath.Dir(p))
			if err != nil {
				return err
			}

			artist.Path = p
			res.artists = append(res.artists, artist)
		case albumMetadataFilename:
			album, err := utils.ReadToml[Album](p)
			if err != nil {
				reporter.AddWarning(p, fmt.Errorf("failed to read album: %w", err))
				return nil
			}

			p, err := filepath.Rel(dir, filepath.Dir(p))
			if err != nil {
				return err
			}

			album.Path = p
			res.albums = append(res.albums, album)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func validateArtistMetadata(artist *Artist, reporter *Reporter) bool {
	file := filepath.Join(artist.Path, artistMetadataFilename)

	valid := true
	if artist.Id == "" {
		reporter.AddError(file, errors.New("id: missing id"))
		valid = false
	}

	if artist.Name == "" {
		reporter.AddError(file, errors.New("name: missing name"))
		valid = false
	}

	if artist.SearchName == "" {
		reporter.AddError(file, errors.New("searchName: missing search name"))
		valid = false
	}

	if artist.Cover != "" {
		p := filepath.Join(artist.Path, artist.Cover)

		err := validateImage(p)
		if err != nil {
			reporter.AddError(file, fmt.Errorf("cover: invalid cover art: %w", err))
			valid = false
		}
	}

	if len(artist.Tags) <= 0 {
		reporter.AddWarning(file, errors.New("tags: missing tags"))
	}

	return valid
}

func validateAlbumMetadata(album *Album, reporter *Reporter) bool {
	file := filepath.Join(album.Path, albumMetadataFilename)

	valid := true

	if album.Album.Id == "" {
		reporter.AddError(file, errors.New("album.id: missing id"))
		valid = false
	}

	if album.Album.Name == "" {
		reporter.AddError(file, errors.New("album.name: missing name"))
		valid = false
	}

	if album.General.Cover != "" {
		p := filepath.Join(album.Path, album.General.Cover)
		err := validateImage(p)
		if err != nil {
			reporter.AddError(file, fmt.Errorf("album.cover: invalid cover art: %w", err))
			valid = false
		}
	}

	if album.Album.Year == 0 {
		reporter.AddWarning(file, errors.New("album.year: year not set"))
	}

	if len(album.Album.Tags) == 0 {
		reporter.AddWarning(file, errors.New("album.tags: tags not set"))
	}

	return valid
}

func validateTrackMetadata(prefix, file string, track *AlbumTrack, reporter *Reporter) bool {
	valid := true

	if track.Id == "" {
		reporter.AddError(file, errors.New(prefix+".id"+": missing id"))
		valid = false
	}

	if track.File == "" {
		reporter.AddError(file, errors.New(prefix+".file"+": missing file"))
		valid = false
	}

	if track.Name == "" {
		reporter.AddError(file, errors.New(prefix+".name"+": missing name"))
		valid = false
	}

	// TODO(patrik): Should we have this check?
	if track.Number == 0 {
		reporter.AddWarning(file, errors.New(prefix+".number"+": missing number"))
	}

	if track.Number < 0 {
		reporter.AddWarning(file, errors.New(prefix+".number"+": should be positive"))
	}

	if track.Year == 0 {
		reporter.AddWarning(file, errors.New(prefix+".year"+": missing year"))
	}

	if track.Year < 0 {
		reporter.AddWarning(file, errors.New(prefix+".year"+": should be positive"))
	}

	if len(track.Tags) <= 0 {
		reporter.AddWarning(file, errors.New(prefix+".tags"+": missing tags"))
	}

	return valid
}

func writeEntry(w io.Writer, entry any) error {
	d, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = w.Write(d)
	if err != nil {
		return err
	}

	_, err = w.Write(([]byte)("\n"))
	if err != nil {
		return err
	}

	return nil
}

type UpdateLibraryOptions struct {
	OnlyArtists bool
}

type UpdateResult struct {
	Reports     map[string][]Report
	NumErrors   int
	NumWarnings int

	FetchingDuration   time.Duration
	ValidationDuration time.Duration
	WritingDuration    time.Duration
	TotalDuration      time.Duration
}

type Library struct {
	Path string

	Reporter Reporter

	Artists []ArtistEntry
	Albums  []AlbumEntry
	Tracks  []TrackEntry
}

func (lib *Library) WriteToDisk() error {
	openLib := func(name string) (*os.File, error) {
		p := filepath.Join(lib.Path, name)
		libFile, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, err
		}

		return libFile, nil
	}

	libArtists, err := openLib("artists")
	if err != nil {
		return fmt.Errorf("open artists library: %w", err)
	}
	defer libArtists.Close()

	libAlbums, err := openLib("albums")
	if err != nil {
		return fmt.Errorf("open albums library: %w", err)
	}
	defer libAlbums.Close()

	libTracks, err := openLib("tracks")
	if err != nil {
		return fmt.Errorf("open tracks library: %w", err)
	}
	defer libTracks.Close()

	for _, artist := range lib.Artists {
		err = writeEntry(libArtists, artist)
		if err != nil {
			return fmt.Errorf("write artist entry: %w", err)
		}
	}

	for _, album := range lib.Albums {
		err = writeEntry(libAlbums, album)
		if err != nil {
			return fmt.Errorf("write album entry: %w", err)
		}
	}

	for _, track := range lib.Tracks {
		err = writeEntry(libTracks, track)
		if err != nil {
			return fmt.Errorf("write track entry: %w", err)
		}
	}

	return nil
}

func ProcessMusicLibrary(dir string) (*Library, error) {
	libraryConfig, err := ReadLibraryConfig(dir)
	if err != nil {
		return nil, err
	}

	lib := &Library{
		Path: libraryConfig.Path,
		Reporter: Reporter{
			// TODO(patrik): Rename errors to reports
			Errors: map[string][]Report{},
		},
		Artists: []ArtistEntry{},
		Albums:  []AlbumEntry{},
		Tracks:  []TrackEntry{},
	}

	fetchingTimer := timer.Simple{}
	processingTimer := timer.Simple{}

	fetchingTimer.Start()

	fmt.Println("fetching files...")

	fetched, err := fetch(lib.Path, libraryConfig.ExcludedDirs, &lib.Reporter)
	if err != nil {
		return nil, fmt.Errorf("dir walk: %w", err)
	}

	fetchingTimer.Stop()

	fmt.Println("fetch done", fetchingTimer.Duration())
	fmt.Printf("Found %d artists\n", len(fetched.artists))
	fmt.Printf("Found %d albums\n", len(fetched.albums))
	fmt.Println()

	processingTimer.Start()

	artistMap := map[string]string{}

	for _, artist := range fetched.artists {
		processArtistMetadata(&artist)
		valid := validateArtistMetadata(&artist, &lib.Reporter)

		if valid {
			id, exists := artistMap[artist.SearchName]
			if exists {
				lib.Reporter.AddError(artist.Path, fmt.Errorf("duplicate artist: '%s' already exists with id '%s'", artist.Name, id))
				continue
			}

			artistMap[artist.SearchName] = artist.Id

			lib.Artists = append(lib.Artists, ArtistEntry{
				Id:       artist.Id,
				Name:     artist.Name,
				CoverArt: artist.Cover,
				Tags:     artist.Tags,
				Path:     artist.Path,
			})
		}
	}

	checkForArtist := func(name string) (string, bool) {
		id, exists := artistMap[utils.Slug(name)]
		if exists {
			return id, true
		}

		return "", false
	}

	type ResolvedArtists struct {
		ArtistId           string
		FeaturingArtistIds []string
	}

	resolveArtists := func(file string, prefix string, artists []string) (ResolvedArtists, bool) {
		if len(artists) <= 0 {
			lib.Reporter.AddError(file, errors.New(prefix+": no artists"))
			return ResolvedArtists{}, false
		}

		valid := true

		artistId := ""
		featuringArtistIds := []string{}

		if id, ok := checkForArtist(artists[0]); ok {
			artistId = id
		} else {
			lib.Reporter.AddError(file, fmt.Errorf("%s: missing artist: '%s'", prefix, artists[0]))
			valid = false
		}

		for _, artist := range artists[1:] {
			if id, ok := checkForArtist(artist); ok {
				featuringArtistIds = append(featuringArtistIds, id)
			} else {
				lib.Reporter.AddError(file, fmt.Errorf("%s: missing artist: '%s'", prefix, artist))
				valid = false
			}
		}

		return ResolvedArtists{
			ArtistId:           artistId,
			FeaturingArtistIds: featuringArtistIds,
		}, valid
	}

	for _, album := range fetched.albums {
		file := filepath.Join(album.Path, albumMetadataFilename)

		processAlbumMetadata(&album)
		valid := validateAlbumMetadata(&album, &lib.Reporter)

		artists, ok := resolveArtists(file, "album.artists", album.Album.Artists)
		if !ok {
			valid = false
		}

		if valid {
			lib.Albums = append(lib.Albums, AlbumEntry{
				Id:                 album.Album.Id,
				Name:               album.Album.Name,
				CoverArt:           album.General.Cover,
				Year:               album.Album.Year,
				ArtistId:           artists.ArtistId,
				FeaturingArtistIds: artists.FeaturingArtistIds,
				Tags:               album.Album.Tags,
				Path:               album.Path,
			})
		}

		for i, track := range album.Tracks {
			prefix := fmt.Sprintf("album.tracks[%d]", i)
			trackValid := validateTrackMetadata(prefix, file, &track, &lib.Reporter)

			artists, ok := resolveArtists(file, prefix+".artists", track.Artists)
			if !ok {
				trackValid = false
			}

			if valid && trackValid {
				lib.Tracks = append(lib.Tracks, TrackEntry{
					Id:                 track.Id,
					TrackFile:          track.File,
					Name:               track.Name,
					Number:             track.Number,
					Year:               track.Year,
					Tags:               track.Tags,
					AlbumId:            album.Album.Id,
					ArtistId:           artists.ArtistId,
					FeaturingArtistIds: artists.FeaturingArtistIds,
					Path:               album.Path,
				})
			}
		}
	}

	processingTimer.Stop()

	keys := slices.Collect(maps.Keys(lib.Reporter.Errors))
	sort.SliceStable(keys, func(i, j int) bool {
		return natural.Less(keys[i], keys[j])
	})

	for _, file := range keys {
		reports := lib.Reporter.Errors[file]

		color.Set(color.FgBlue)
		fmt.Fprintln(os.Stderr, file)

		for _, report := range reports {
			if report.IsWarning {
				color.Set(color.FgYellow)
				fmt.Fprintf(os.Stderr, " - warn:  ")
			} else {
				color.Set(color.FgRed)
				fmt.Fprintf(os.Stderr, " - error: ")
			}

			fmt.Fprintf(os.Stderr, "%s\n", report.Err.Error())
		}

		color.Unset()

		fmt.Fprintln(os.Stderr)
	}

	color.Set(color.FgGreen)

	fmt.Println("Report:")
	fmt.Printf(" Total:    %v\n", (lib.Reporter.NumErrors + lib.Reporter.NumWarnings))

	color.Set(color.FgRed)
	fmt.Printf(" Errors:   ")
	color.Set(color.FgGreen)
	fmt.Printf("%v\n", lib.Reporter.NumErrors)

	color.Set(color.FgYellow)
	fmt.Printf(" Warnings: ")
	color.Set(color.FgGreen)
	fmt.Printf("%v\n", lib.Reporter.NumWarnings)
	fmt.Println()

	color.Set(color.FgMagenta)

	total := fetchingTimer.Duration() + processingTimer.Duration()

	fmt.Println("Time:")
	fmt.Printf(" Total: %v\n", total)
	fmt.Printf(" Fetching: %v\n", fetchingTimer.Duration())
	fmt.Printf(" Processing: %v\n", processingTimer.Duration())

	color.Unset()

	return lib, nil
}
