package main

type Thumbnail struct {
	Height int64  `json:"height"`
	Width  int64  `json:"width"`
	Url    string `json:"url"`
}

type Item struct {
	Id         string      `json:"id"`
	Title      string      `json:"title"`
	Sub        string      `json:"subtitle"`
	Thumbnails []Thumbnail `json:"thumbnails"`
}

type Items struct {
	Songs   []Item `json:"songs"`
	Albums  []Item `json:"albums"`
	Artists []Item `json:"recommendedArtists"`
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

type MusicConfig struct {
	MusicConfig PageType `json:"browseEndpointContextMusicConfig"`
}

type BrowseData struct {
	Context     Context     `json:"context"`
	MusicConfig MusicConfig `json:"browseEndpointContextMusicConfig"`
	BrowseId    string      `json:"browseId"`
}
