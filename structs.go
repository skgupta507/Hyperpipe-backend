package main

/* Parsers */
type Thumbnail struct {
	Height int64  `json:"height"`
	Width  int64  `json:"width"`
	Url    string `json:"url"`
}

type Item struct {
	Id         string      `json:"id"`
	Title      string      `json:"title"`
	Sub        string      `json:"subtitle"`
	SubId      string      `json:"subId"`
	Thumbnails []Thumbnail `json:"thumbnails"`
}

type Items struct {
	Songs   []Item `json:"songs"`
	Albums  []Item `json:"albums"`
	Singles []Item `json:"singles"`
	Artists []Item `json:"recommendedArtists"`
}

type Home struct {
	Contents []map[string]interface{} `json:"contents"`
	Continue string                   `json:"continue"`
}

type Explore struct {
	Albums   []Item `json:"albums_and_singles"`
	Trending []Item `json:"trending"`
}

type Artist struct {
	Title            string      `json:"title"`
	Description      string      `json:"description"`
	BrowsePlaylistId string      `json:"browsePlaylistId"`
	PlaylistId       string      `json:"playlistId"`
	SubscriberCount  string      `json:"subscriberCount"`
	Thumbnails       []Thumbnail `json:"thumbnails"`
	Items            Items       `json:"items"`
}

type MediaSession struct {
	Thumbnails []Thumbnail `json:"thumbnails"`
	Album      string      `json:"album"`
}

type Next struct {
	LyricsId     string       `json:"lyricsId"`
	Songs        []Item       `json:"songs"`
	MediaSession MediaSession `json:"mediaSession"`
}

type Lyrics struct {
	Text   string `json:"text"`
	Source string `json:"source"`
}

/* Structs and Types */
type Client struct {
	Name    string `json:"clientName"`
	Version string `json:"clientVersion"`
}

type Context struct {
	Client Client `json:"client"`
}

type PageType struct {
	PageType string `json:"pageType"`
}

type BrowseMusicConfig struct {
	MusicConfig PageType `json:"browseEndpointContextMusicConfig"`
}

type WatchMusicConfig struct {
	Panel bool   `json:"hasPersistentPlaylistPanel"`
	Type  string `json:"musicVideoType"`
}

type BrowseData struct {
	Context     Context           `json:"context"`
	MusicConfig BrowseMusicConfig `json:"browseEndpointContextMusicConfig"`
	BrowseId    string            `json:"browseId"`
}

type NextData struct {
	Id          string           `json:"videoId"`
	Context     Context          `json:"context"`
	Audio       bool             `json:"isAudioOnly"`
	Tuner       string           `json:"tunerSettingValue"`
	Panel       bool             `json:"enablePersistentPlaylistPanel"`
	MusicConfig WatchMusicConfig `json:"watchEndpointMusicConfig"`
}
