package lib

import (
	"encoding/json"

	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

type Options struct {
	Default string `json:"default"`
	All     []Item `json:"all"`
}

type Charts struct {
	Options  Options `json:"options"`
	Artists  []Item  `json:"artists"`
	Trending []Item  `json:"trending"`
}

func parseCharts(raw string) (string, error) {

	j := gjson.Parse(raw)

	c := j.Get("contents.singleColumnBrowseResultsRenderer.tabs.0").Get("tabRenderer.content.sectionListRenderer.contents")

	o := c.Get(
		"0.musicShelfRenderer.subheaders.0.musicSideAlignedItemRenderer",
	).Get(
		"startItems.0.musicSortFilterButtonRenderer",
	)

	a := c.Get(
		"#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer",
	).Get(
		"title.runs.0.text == Top artists).musicCarouselShelfRenderer.contents",
	)

	t := c.Get(
		"#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer",
	).Get("title.runs.0.text == Trending).musicCarouselShelfRenderer.contents")

	opts := o.Get(
		"menu.musicMultiSelectMenuRenderer.options.#.musicMultiSelectMenuItemRenderer",
	)
	ref := j.Get("frameworkUpdates.entityBatchUpdate.mutations")

	val := Charts{
		Options: Options{
			Default: RunsText(o.Get("title")),
			All:     MultiSelectMenuItemRenderer(opts, ref),
		},
		Artists:  ResponsiveListItemRendererCH(a),
		Trending: ResponsiveListItemRenderer(t),
	}

	res, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetCharts(params, code string) (string, int) {

	context := utils.TypeBrowse("", "FEmusic_charts", params, code)

	raw, status := utils.FetchBrowse(context)

	res, err := parseCharts(raw)
	if err != nil {
		return utils.ErrMsg(err), 500
	}

	return res, status
}
