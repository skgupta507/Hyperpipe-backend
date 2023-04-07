package lib

import (
	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

type Items struct {
	Songs   []Item `json:"songs"`
	Albums  []Item `json:"albums"`
	Singles []Item `json:"singles"`
	Artists []Item `json:"recommendedArtists"`
}

type MoreItem struct {
	Params string `json:"params"`
	Click  string `json:"click"`
	Visit  string `json:"visit"`
}

type ArtistMore struct {
	Album   MoreItem `json:"albums"`
	Singles MoreItem `json:"singles"`
}

type ArtistNext struct {
	Title string `json:"title"`
	Items []Item `json:"items"`
}

type Artist struct {
	Title            string      `json:"title"`
	Description      string      `json:"description,omitempty"`
	BrowsePlaylistId string      `json:"browsePlaylistId,omitempty"`
	PlaylistId       string      `json:"playlistId,omitempty"`
	SubscriberCount  string      `json:"subscriberCount,omitempty"`
	Thumbnails       []Thumbnail `json:"thumbnails"`
	Items            Items       `json:"items"`
	More             ArtistMore  `json:"more"`
}

func parseMoreButton(raw, v gjson.Result) MoreItem {
	nav := raw.Get("header.musicCarouselShelfBasicHeaderRenderer.moreContentButton.buttonRenderer.navigationEndpoint")

	return MoreItem{
		Params: nav.Get("browseEndpoint.params").String(),
		Click:  nav.Get("clickTrackingParams").String(),
		Visit:  v.String(),
	}
}

func parseArtist(raw string) Artist {

	j := gjson.Parse(raw)

	h := j.Get("header.musicImmersiveHeaderRenderer")
	c := j.Get(
		"contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer.content.sectionListRenderer.contents",
	)

	s := c.Get("#(musicShelfRenderer.title.runs.0.text == Songs).musicShelfRenderer")

	a := c.Get(
		"#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == Albums).musicCarouselShelfRenderer",
	)

	m := c.Get(
		"#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == Singles).musicCarouselShelfRenderer",
	)

	u := c.Get(
		"#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == Fans might also like).musicCarouselShelfRenderer",
	)

	return Artist{
		Title:       RunsText(h.Get("title")),
		Description: RunsText(h.Get("description")),
		SubscriberCount: RunsText(
			h.Get("subscriptionButton.subscribeButtonRenderer.subscriberCountText"),
		),
		Thumbnails: GetThumbnails(
			h.Get("thumbnail.musicThumbnailRenderer.thumbnail.thumbnails"),
		),
		BrowsePlaylistId: h.Get(
			"playButton.buttonRenderer.navigationEndpoint.watchEndpoint.playlistId",
		).String(),
		PlaylistId: s.Get(
			"contents.0.musicResponsiveListItemRenderer.flexColumns.0.musicResponsiveListItemFlexColumnRenderer",
		).Get("text.runs.0.navigationEndpoint.watchEndpoint.playlistId").String(),
		Items: Items{
			Songs:   ResponsiveListItemRenderer(s.Get("contents")),
			Albums:  TwoRowItemRenderer(a.Get("contents"), true),
			Singles: TwoRowItemRenderer(m.Get("contents"), true),
			Artists: TwoRowItemRenderer(u.Get("contents"), false),
		},
		More: ArtistMore{
			Album:   parseMoreButton(a, j.Get("responseContext.visitorData")),
			Singles: parseMoreButton(m, j.Get("responseContext.visitorData")),
		},
	}
}

func parseArtistNext(raw string) ArtistNext {

	j := gjson.Parse(raw)

	return ArtistNext{
		Title: RunsText(j.Get("header.musicHeaderRenderer.title")),
		Items: TwoRowItemRenderer(
			j.Get(
				"contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer.content.sectionListRenderer.contents.0.gridRenderer.items",
			),
			true,
		),
	}
}

func GetArtist(id string) (Artist, int) {

	context := utils.TypeBrowsePage(id, "artist")

	raw, status := utils.FetchBrowse(context)

	res := parseArtist(raw)

	return res, status
}

func GetArtistNext(id, params, ct, v string) (ArtistNext, int) {

	if id == "" || params == "" || ct == "" || v == "" {
		return ArtistNext{}, 404
	}
	
	context := utils.TypeBrowse(id, params, []string{ct, v})

	raw, status := utils.FetchBrowse(context)

	res := parseArtistNext(raw)

	return res, status
}
