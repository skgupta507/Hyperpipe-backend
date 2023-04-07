package lib

import (
	"encoding/json"
	"sync"

	"codeberg.org/Hyperpipe/hyperpipe-backend/utils"
	"github.com/tidwall/gjson"
)

type MediaSession struct {
	Thumbnails []Thumbnail `json:"thumbnails"`
	Album      string      `json:"album"`
}

type Next struct {
	Err error `json:"error"`
	LyricsId     string       `json:"lyricsId"`
	Songs        []Item       `json:"songs"`
	MediaSession MediaSession `json:"mediaSession"`
}

func parseNextSongs(n gjson.Result) []Item {

	np := n.Get("#.playlistPanelVideoRenderer")

	r := make([]Item, n.Get("#").Int())

	wg := sync.WaitGroup{}

	np.ForEach(
		func(n, v gjson.Result) bool {

			wg.Add(1)

			go func(i int64, j gjson.Result) {
				defer wg.Done()

				r[i] = Item{
					Id:         j.Get("navigationEndpoint.watchEndpoint.videoId").String(),
					Title:      RunsText(j.Get("title")),
					Sub:        j.Get("longBylineText.runs.2.text").String(),
					Thumbnails: GetThumbnails(j.Get("thumbnail.thumbnails")),
				}
			}(n.Int(), v)

			return true
		},
	)

	wg.Wait()

	return r
}

func parseNext(raw string) Next {

	j := gjson.Parse(raw)

	c := j.Get(
		"contents.singleColumnMusicWatchNextResultsRenderer.tabbedRenderer.watchNextTabbedResultsRenderer.tabs",
	)
	m := j.Get(
		"playerOverlays.playerOverlayRenderer.browserMediaSession.browserMediaSessionRenderer",
	)

	n := c.Get(
		"#(tabRenderer.title == Up next).tabRenderer.content",
	).Get("musicQueueRenderer.content.playlistPanelRenderer")

	l := c.Get(
		"#(tabRenderer.title == Lyrics).tabRenderer.endpoint.browseEndpoint.browseId",
	)

	return Next{
		LyricsId: l.String(),
		MediaSession: MediaSession{
			Album:      RunsText(m.Get("album")),
			Thumbnails: GetThumbnails(m.Get("thumbnailDetails.thumbnails")),
		},
		Songs: parseNextSongs(n.Get("contents")),
	}
}

func GetNext(id, queue string) (Next, int) {

	pldata, err := json.Marshal(utils.TypeNext(id, ""))
	if err != nil {
		return Next{Err: err}, 500
	}

	plraw, plstatus, err := utils.Fetch("next", pldata)
	if err != nil || plstatus > 399 {
		return Next{Err: err}, plstatus
	}

	if queue == "avoid" {
		return parseNext(plraw), plstatus
	}

	pl := gjson.Parse(plraw).Get(
		"contents.singleColumnMusicWatchNextResultsRenderer." +
			"tabbedRenderer.watchNextTabbedResultsRenderer.tabs.0.tabRenderer.content." +
			"musicQueueRenderer.content.playlistPanelRenderer.contents." +
			"#(automixPreviewVideoRenderer).automixPreviewVideoRenderer." +
			"content.automixPlaylistVideoRenderer.navigationEndpoint." +
			"watchPlaylistEndpoint.playlistId").String()

	data, err := json.Marshal(utils.TypeNext(id, pl))
	if err != nil {
		return Next{Err: err}, 500
	}

	raw, status, err := utils.Fetch("next", data)
	if err != nil {
		return Next{Err: err}, 500
	}

	return parseNext(raw), status
}
