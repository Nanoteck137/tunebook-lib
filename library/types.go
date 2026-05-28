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
