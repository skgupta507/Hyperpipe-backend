package lib

import (
	"encoding/json"

	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

type Artist struct {
	Title            string      `json:"title"`
	Description      string      `json:"description,omitempty"`
	BrowsePlaylistId string      `json:"browsePlaylistId,omitempty"`
	PlaylistId       string      `json:"playlistId,omitempty"`
	SubscriberCount  string      `json:"subscriberCount,omitempty"`
	Thumbnails       []Thumbnail `json:"thumbnails"`
	Items            Items       `json:"items"`
}

func parseArtist(raw string) (string, error) {

	j := gjson.Parse(raw)

	h := j.Get("header.musicImmersiveHeaderRenderer")
	c := j.Get(
		"contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer.content.sectionListRenderer.contents",
	)

	s := c.Get("#(musicShelfRenderer.title.runs.0.text == Songs).musicShelfRenderer")

	a := c.Get(
		"#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer",
	).Get("title.runs.0.text == Albums).musicCarouselShelfRenderer")

	m := c.Get(
		"#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer",
	).Get("title.runs.0.text == Singles).musicCarouselShelfRenderer")

	u := c.Get(
		"#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer",
	).Get("title.runs.0.text == Fans might also like).musicCarouselShelfRenderer")

	val := Artist{
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
	}

	res, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetArtist(id string) (string, int) {

	context := utils.TypeBrowse("artist", id, "", "")

	raw, status := utils.FetchBrowse(context)

	res, err := parseArtist(raw)
	if err != nil {
		return utils.ErrMsg(err), 500
	}

	return res, status
}
