package library

// TODO(patrik): Move?
const (
	libraryMetadataFilename = "library.toml"
	artistMetadataFilename  = "artist.toml"
	albumMetadataFilename   = "album.toml"
)

type LibraryMetadata struct {
	ExcludedDirs []string `json:"excludedDirs" toml:"excludedDirs"`

	Path string `json:"-" toml:"-"`
}

type AlbumMetadataGeneral struct {
	Cover     string   `json:"cover" toml:"cover"`
	Tags      []string `json:"tags" toml:"tags"`
	TrackTags []string `json:"trackTags" toml:"trackTags"`
	Year      int64    `json:"year" toml:"year"`
}

type AlbumMetadataAlbum struct {
	Id      string   `json:"id" toml:"id"`
	Name    string   `json:"name" toml:"name"`
	Year    int64    `json:"year" toml:"year"`
	Tags    []string `json:"tags" toml:"tags"`
	Artists []string `json:"artists" toml:"artists"`
}

type AlbumMetadataTrack struct {
	Id      string   `json:"id" toml:"id"`
	File    string   `json:"file" toml:"file"`
	Name    string   `json:"name" toml:"name"`
	Number  int64    `json:"number" toml:"number"`
	Year    int64    `json:"year" toml:"year"`
	Tags    []string `json:"tags" toml:"tags"`
	Artists []string `json:"artists" toml:"artists"`
}

// TODO(patrik): Rename
type AlbumMetadata struct {
	General AlbumMetadataGeneral `json:"general" toml:"general"`
	Album   AlbumMetadataAlbum   `json:"album" toml:"album"`
	Tracks  []AlbumMetadataTrack `json:"tracks" toml:"tracks"`

	Path string `json:"-" toml:"-"`
}

type ArtistMetadata struct {
	Id string `json:"id" toml:"id"`

	SearchName string   `json:"searchName" toml:"searchName"`
	Name       string   `json:"name" toml:"name"`
	Cover      string   `json:"cover" toml:"cover"`
	Tags       []string `json:"tags" toml:"tags"`

	Path string `json:"-" toml:"-"`
}

// NOTE(patrik): These are copied from tunebook (library/types.go), and are
// the final entries used when importing

type ArtistEntry struct {
	Id string `json:"id"`

	Name     string   `json:"name"`
	CoverArt string   `json:"coverArt"`
	Tags     []string `json:"tags"`

	Path string `json:"path"`
}

type AlbumEntry struct {
	Id string `json:"id"`

	Name               string   `json:"name"`
	CoverArt           string   `json:"coverArt"`
	Year               int64    `json:"year"`
	ArtistId           string   `json:"artistId"`
	FeaturingArtistIds []string `json:"featuringArtistIds"`
	Tags               []string `json:"tags"`

	Path string `json:"path"`
}

type TrackEntry struct {
	Id string `json:"id"`

	Name               string   `json:"name"`
	TrackFile          string   `json:"trackFile"`
	Number             int64    `json:"number"`
	Year               int64    `json:"year"`
	Tags               []string `json:"tags"`
	AlbumId            string   `json:"albumId"`
	ArtistId           string   `json:"artistId"`
	FeaturingArtistIds []string `json:"featuringArtistIds"`

	Path string `json:"path"`
}
