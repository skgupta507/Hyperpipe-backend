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
	Sub        string      `json:"subtitle,omitempty"`
	SubId      string      `json:"subId,omitempty"`
	Thumbnails []Thumbnail `json:"thumbnails,omitempty"`
}

type Items struct {
	Songs   []Item `json:"songs"`
	Albums  []Item `json:"albums"`
	Singles []Item `json:"singles"`
	Artists []Item `json:"recommendedArtists"`
}

type Explore struct {
	TrendingId string `json:"trendingId"`
	Albums     []Item `json:"albums_and_singles"`
	Trending   []Item `json:"trending"`
}

type Genres struct {
	Moods  []Item `json:"moods"`
	Genres []Item `json:"genres"`
}

type Genre struct {
	Title     string `json:"title"`
	Spotlight []Item `json:"spotlight"`
	Featured  []Item `json:"featured"`
	Community []Item `json:"community"`
}

type Artist struct {
	Title            string      `json:"title"`
	Description      string      `json:"description,omitempty"`
	BrowsePlaylistId string      `json:"browsePlaylistId,omitempty"`
	PlaylistId       string      `json:"playlistId,omitempty"`
	SubscriberCount  string      `json:"subscriberCount,omitempty"`
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
	Name    string `json:"clientName,omitempty"`
	Version string `json:"clientVersion,omitempty"`
	Hl      string `json:"hl",omitempty`
}

type Context struct {
	Client Client `json:"client",omitempty`
}

type PageType struct {
	PageType string `json:"pageType",omitempty`
}

type BrowseMusicConfig struct {
	MusicConfig PageType `json:"browseEndpointContextMusicConfig,omitempty"`
}

type WatchMusicConfig struct {
	Panel bool   `json:"hasPersistentPlaylistPanel,omitempty"`
	Type  string `json:"musicVideoType,omitempty"`
}

type BrowseData struct {
	Context     Context           `json:"context,omitempty"`
	MusicConfig BrowseMusicConfig `json:"browseEndpointContextMusicConfig,omitempty"`
	BrowseId    string            `json:"browseId,omitempty"`
	Params      string            `json:"params,omitempty"`
}

type NextData struct {
	Id          string           `json:"videoId,omitempty"`
	PlaylistId  string           `json:"playlistId",omitempty`
	Context     Context          `json:"context,omitempty"`
	Audio       bool             `json:"isAudioOnly,omitempty"`
	Tuner       string           `json:"tunerSettingValue,omitempty"`
	Panel       bool             `json:"enablePersistentPlaylistPanel,omitempty"`
	MusicConfig WatchMusicConfig `json:"watchEndpointMusicConfig,omitempty"`
}
