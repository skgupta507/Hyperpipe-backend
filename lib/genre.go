package lib

import (
	"encoding/json"

	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

type Genre struct {
	Title     string                 `json:"title"`
	Spotlight []Item                 `json:"spotlight"`
	Featured  []Item                 `json:"featured"`
	Community []Item                 `json:"community"`
	Shelf     map[string]interface{} `json:"shelf"`
}

func parseGenre(raw string) (string, error) {
	j := gjson.Parse(raw)

	c := j.Get(
		"contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer.content.sectionListRenderer.contents",
	)

	g := c.Get("#.gridRenderer")

	var val Genre

	s := g.Get("#(header.gridHeaderRenderer.title.runs.0.text == Spotlight)")
	f := g.Get("#(header.gridHeaderRenderer.title.runs.0.text == Featured playlists)")
	cp := g.Get("#(header.gridHeaderRenderer.title.runs.0.text == Community playlists)")

	val = Genre{
		Title:     RunsText(j.Get("header.musicHeaderRenderer.title")),
		Spotlight: TwoRowItemRenderer(s.Get("items"), false),
		Featured:  TwoRowItemRenderer(f.Get("items"), false),
		Community: TwoRowItemRenderer(cp.Get("items"), false),
		Shelf: ShelfRenderer(
			c.Get("#.musicCarouselShelfRenderer"),
		),
	}

	res, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetGenre(param string) (string, int) {

	context := utils.TypeBrowse("", "FEmusic_moods_and_genres_category", param, "")

	raw, status := utils.FetchBrowse(context)

	res, err := parseGenre(raw)
	if err != nil {
		return utils.ErrMsg(err), 500
	}

	return res, status
}
