package library

// NOTE(patrik): These are copied from tunebook (library/db.go), and are
// the final entries used when importing

type ArtistEntry struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	CoverArt string   `json:"coverArt"`
	Tags     []string `json:"tags"`

	Path string `json:"path"`
}

type AlbumEntry struct {
	Id                 string   `json:"id"`
	Name               string   `json:"name"`
	CoverArt           string   `json:"coverArt"`
	Year               int64    `json:"year"`
	ArtistId           string   `json:"artistId"`
	FeaturingArtistIds []string `json:"featuringArtistIds"`
	Tags               []string `json:"tags"`

	Path string `json:"path"`
}

type TrackEntry struct {
	Id                 string   `json:"id"`
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
