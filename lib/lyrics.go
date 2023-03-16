package lib

import (
	"encoding/json"

	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

type Lyrics struct {
	Text   string `json:"text"`
	Source string `json:"source"`
}

func parseLyrics(raw string) (string, error) {

	j := gjson.Parse(raw)

	d := j.Get("contents.sectionListRenderer.contents.0.musicDescriptionShelfRenderer")

	l := d.Get("description")
	s := d.Get("footer")

	val := Lyrics{
		Text:   RunsText(l),
		Source: RunsText(s),
	}

	res, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetLyrics(id string) (string, int) {

	context := utils.TypeBrowse("lyrics", id, "", "")

	raw, status := utils.FetchBrowse(context)

	res, err := parseLyrics(raw)
	if err != nil {
		return utils.ErrMsg(err), 500
	}

	return res, status
}
