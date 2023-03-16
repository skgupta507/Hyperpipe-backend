package lib

import (
	"encoding/json"

	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

type Explore struct {
	TrendingId string `json:"trendingId"`
	ChartsId   string `json:"chartsId"`
	Albums     []Item `json:"albums_and_singles"`
	Trending   []Item `json:"trending"`
}

func parseExplore(raw string) (string, error) {

	j := gjson.Parse(raw)

	m := j.Get(
		"contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer.content.sectionListRenderer.contents",
	)

	c := m.Get("#.musicCarouselShelfRenderer")

	a := c.Get(
		"#(header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == New albums & singles)",
	)
	t := c.Get("#(header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == Trending)")

	charts := m.Get(
		"#(gridRenderer).gridRenderer.items.#.musicNavigationButtonRenderer",
	).Get("#(buttonText.runs.0.text == Charts)")

	val := Explore{
		TrendingId: t.Get(
			"header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.navigationEndpoint.browseEndpoint.browseId",
		).String(),
		ChartsId: charts.Get("clickCommand.browseEndpoint.params").String(),
		Albums:   TwoRowItemRenderer(a.Get("contents"), true),
		Trending: ResponsiveListItemRenderer(t.Get("contents")),
	}

	res, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetExplore() (string, int) {

	context := utils.TypeBrowse("", "FEmusic_explore", "", "")

	raw, status := utils.FetchBrowse(context)

	res, err := parseExplore(raw)
	if err != nil {
		return utils.ErrMsg(err), 500
	}

	return res, status
}
