package library

const (
	libraryFilename = "library.toml"
	artistFilename  = "artist.toml"
	albumFilename   = "album.toml"
	tagFilename     = "tag.toml"
)

type LibraryConfig struct {
	ExcludedDirs []string `toml:"excludedDirs"`

	Path string `toml:"-"`
}

type TagFile struct {
	Tags       []string            `toml:"tags"`
	Artist     []string            `toml:"artist"`
	Album      []string            `toml:"album"`
	Track      []string            `toml:"track"`
	Maps       map[string][]string `toml:"map"`
	ArtistMaps map[string][]string `toml:"mapArtist"`
	AlbumMaps  map[string][]string `toml:"mapAlbum"`
	TrackMaps  map[string][]string `toml:"mapTrack"`
}

type TagTarget string

const (
	TagTargetArtist TagTarget = "artist"
	TagTargetAlbum  TagTarget = "album"
	TagTargetTrack  TagTarget = "track"
)

type AlbumGeneral struct {
	Cover     string   `toml:"cover"`
	Tags      []string `toml:"tags"`
	TrackTags []string `toml:"trackTags"`
	Year      int64    `toml:"year"`
}

type AlbumAlbum struct {
	Id      string   `toml:"id"`
	Name    string   `toml:"name"`
	Year    int64    `toml:"year"`
	Tags    []string `toml:"tags"`
	Artists []string `toml:"artists"`
}

type AlbumTrack struct {
	Id      string   `toml:"id"`
	File    string   `toml:"file"`
	Name    string   `toml:"name"`
	Number  int64    `toml:"number"`
	Year    int64    `toml:"year"`
	Tags    []string `toml:"tags"`
	Artists []string `toml:"artists"`
}

type Album struct {
	General AlbumGeneral `toml:"general"`
	Album   AlbumAlbum   `toml:"album"`
	Tracks  []AlbumTrack `toml:"tracks"`

	Path string `toml:"-"`
}

type Artist struct {
	Id         string   `toml:"id"`
	SearchName string   `toml:"searchName"`
	Name       string   `toml:"name"`
	Cover      string   `toml:"cover"`
	Tags       []string `toml:"tags"`

	Path string `toml:"-"`
}
