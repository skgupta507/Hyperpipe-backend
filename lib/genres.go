package lib

import (
	"encoding/json"

	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

type Genres struct {
	Moods  []Item `json:"moods"`
	Genres []Item `json:"genres"`
}

func parseGenres(raw string) (string, error) {

	j := gjson.Parse(raw)

	c := j.Get("contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer" +
		".content.sectionListRenderer.contents.#.gridRenderer")

	m := c.Get("#(header.gridHeaderRenderer.title.runs.0.text == Moods & moments)")
	g := c.Get("#(header.gridHeaderRenderer.title.runs.0.text == Genres)")

	val := Genres{
		Moods:  NavigationButton(m.Get("items")),
		Genres: NavigationButton(g.Get("items")),
	}

	res, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetGenres() (string, int) {

	context := utils.TypeBrowse("", "FEmusic_moods_and_genres", "", "")

	raw, status := utils.FetchBrowse(context)

	res, err := parseGenres(raw)
	if err != nil {
		return utils.ErrMsg(err), 500
	}

	return res, status
}
