package main

import (
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
)

func RunsText(j gjson.Result) string {

	var s []string

	a := j.Get("runs.#.text").Array()

	for i := 0; i < len(a); i++ {
		s = append(s, a[i].String())
	}

	return strings.Join(s, "")
}

func GetThumbnails(j gjson.Result) []Thumbnail {

	t := []Thumbnail{}

	wg := sync.WaitGroup{}

	j.ForEach(
		func(_, v gjson.Result) bool {

			wg.Add(1)

			go func() {
				defer wg.Done()

				u, err := ParseUrl(v.Get("url").String())
				if err != nil {
					log.Fatal(err)
				}

				t = append(t, Thumbnail{
					Url:    u,
					Width:  v.Get("width").Int(),
					Height: v.Get("height").Int(),
				})
			}()

			wg.Wait()

			return true
		},
	)

	return t
}

func GetNextSongs(n gjson.Result) []Item {

	r := []Item{}

	wg := sync.WaitGroup{}

	n.Get("#.playlistPanelRenderer").ForEach(
		func(_, j gjson.Result) bool {

			wg.Add(1)

			go func() {
				defer wg.Done()

				r = append(r, Item{
					Id:         j.Get("navigationEndpoint.watchEndpoint.videoId").String(),
					Title:      RunsText(j.Get("title")),
					Sub:        j.Get("longBylineText.runs.2.text").String(),
					Thumbnails: GetThumbnails(j.Get("thumbnail.thumbnails")),
				})
			}()

			wg.Wait()

			return true
		},
	)

	return r
}

func ParseHome(raw string) (string, error) {

	j := gjson.Parse(raw)

	c := j.Get("contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer.content.sectionListRenderer")

	m := c.Get("contents")
	n := c.Get("continuations.0.nextContinuationData.continuation").String()

	var body []map[string]interface{}

	m.ForEach(
		func(key, value gjson.Result) bool {

			var v map[string]interface{}

			icsr := value.Get("musicImmersiveCarouselShelfRenderer")
			csr := value.Get("musicCarouselShelfRenderer")

			switch {
			case icsr.Exists():
				v = CarouselShelfRenderer(icsr)
			case csr.Exists():
				v = CarouselShelfRenderer(csr)
			default:
				v = make(map[string]interface{})
				v["raw"] = value.Value()
			}

			body = append(body, v)

			return true
		},
	)

	data := Home{
		Contents: body,
		Continue: n,
	}

	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func ParseExplore(raw string) (string, error) {

	j := gjson.Parse(raw)

	c := j.Get("contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer.content.sectionListRenderer.contents")

	a := c.Get("#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == New albums & singles).musicCarouselShelfRenderer")
	t := c.Get("#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == Trending).musicCarouselShelfRenderer")

	val := Explore{
		Albums:   TwoRowItemRenderer("album", a.Get("contents")),
		Trending: ResponsiveListItemRenderer(t.Get("contents")),
	}

	res, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func ParseArtist(raw string) (string, error) {

	j := gjson.Parse(raw)

	h := j.Get("header.musicImmersiveHeaderRenderer")
	c := j.Get("contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer.content.sectionListRenderer.contents")

	s := c.Get("#(musicShelfRenderer.title.runs.0.text == Songs).musicShelfRenderer")
	a := c.Get("#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == Albums).musicCarouselShelfRenderer")
	m := c.Get("#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == Singles).musicCarouselShelfRenderer")
	u := c.Get("#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == Fans might also like).musicCarouselShelfRenderer")

	val := Artist{
		Title:            RunsText(h.Get("title")),
		Description:      RunsText(h.Get("description")),
		SubscriberCount:  RunsText(h.Get("subscriptionButton.subscribeButtonRenderer.subscriberCountText")),
		Thumbnails:       GetThumbnails(h.Get("thumbnail.musicThumbnailRenderer.thumbnail.thumbnails")),
		BrowsePlaylistId: h.Get("playButton.buttonRenderer.navigationEndpoint.watchEndpoint.playlistId").String(),
		PlaylistId:       s.Get("contents.0.musicResponsiveListItemRenderer.flexColumns.0.musicResponsiveListItemFlexColumnRenderer.text.runs.0.navigationEndpoint.watchEndpoint.playlistId").String(),
		Items: Items{
			Songs:   ResponsiveListItemRenderer(s.Get("contents")),
			Albums:  TwoRowItemRenderer("album", a.Get("contents")),
			Singles: TwoRowItemRenderer("singles", m.Get("contents")),
			Artists: TwoRowItemRenderer("artist", u.Get("contents")),
		},
	}

	res, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func ParseLyrics(raw string) (string, error) {

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

func ParseNext(raw string) (string, error) {

	j := gjson.Parse(raw)

	c := j.Get("contents.singleColumnMusicWatchNextResultsRenderer.tabbedRenderer.watchNextTabbedResultsRenderer.tabs")
	m := j.Get("playerOverlays.playerOverlayRenderer.browserMediaSession.browserMediaSessionRenderer")

	n := c.Get("#(tabRenderer.title == Up next).tabRenderer.content.musicQueueRenderer.content.playlistPanelRenderer")
	l := c.Get("#(tabRenderer.title == Lyrics).tabRenderer.endpoint.browseEndpoint.browseId")

	val := Next{
		LyricsId: l.String(),
		MediaSession: MediaSession{
			Album:      RunsText(m.Get("album")),
			Thumbnails: GetThumbnails(m.Get("thumbnailDetails.thumbnails")),
		},
		Songs: GetNextSongs(n.Get("contents")),
	}

	res, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
