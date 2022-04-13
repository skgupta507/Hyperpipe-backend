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

func GetSongs(s gjson.Result) []Item {

	r := []Item{}

	wg := sync.WaitGroup{}

	s.ForEach(
		func(_, v gjson.Result) bool {

			wg.Add(1)

			go func() {
				defer wg.Done()

				j := v.Get("musicResponsiveListItemRenderer")
				f := j.Get("flexColumns.#.musicResponsiveListItemFlexColumnRenderer.text.runs.0")

				r = append(r, Item{
					Id:         j.Get("playlistItemData.videoId").String(),
					Title:      f.Get("#(navigationEndpoint.watchEndpoint.videoId).text").String(),
					Thumbnails: GetThumbnails(j.Get("thumbnail.musicThumbnailRenderer.thumbnail.thumbnails")),
				})

			}()

			wg.Wait()

			return true
		},
	)

	return r
}

func GetTwoRowItem(t string, a gjson.Result) []Item {

	r := []Item{}

	var id string

	if t == "album" {
		id = "menu.menuRenderer.items.#(menuNavigationItemRenderer.text.runs.0.text == Shuffle play).menuNavigationItemRenderer.navigationEndpoint.watchPlaylistEndpoint.playlistId"
	} else if t == "artist" {
		id = "navigationEndpoint.browseEndpoint.browseId"
	}

	wg := sync.WaitGroup{}

	a.ForEach(
		func(_, v gjson.Result) bool {

			wg.Add(1)

			go func() {
				defer wg.Done()

				j := v.Get("musicTwoRowItemRenderer")

				r = append(r, Item{
					Id:         j.Get(id).String(),
					Title:      RunsText(j.Get("title")),
					Sub:        RunsText(j.Get("subtitle")),
					Thumbnails: GetThumbnails(j.Get("thumbnailRenderer.musicThumbnailRenderer.thumbnail.thumbnails")),
				})

			}()

			wg.Wait()

			return true
		},
	)

	return r
}

func ParseArtist(raw string) (string, error) {

	j := gjson.Parse(raw)

	h := j.Get("header.musicImmersiveHeaderRenderer")
	c := j.Get("contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer.content.sectionListRenderer.contents")

	s := c.Get("#(musicShelfRenderer.title.runs.0.text == Songs).musicShelfRenderer")
	a := c.Get("#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == Albums).musicCarouselShelfRenderer")
	u := c.Get("#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer.title.runs.0.text == Fans might also like).musicCarouselShelfRenderer")

	val := Artist{
		Title:            RunsText(h.Get("title")),
		Description:      RunsText(h.Get("description")),
		SubscriberCount:  RunsText(h.Get("subscriptionButton.subscribeButtonRenderer.subscriberCountText")),
		Thumbnails:       GetThumbnails(h.Get("thumbnail.musicThumbnailRenderer.thumbnail.thumbnails")),
		BrowsePlaylistId: h.Get("playButton.buttonRenderer.navigationEndpoint.watchEndpoint.playlistId").String(),
		PlaylistId:       s.Get("contents.0.musicResponsiveListItemRenderer.flexColumns.0.musicResponsiveListItemFlexColumnRenderer.text.runs.0.navigationEndpoint.watchEndpoint.playlistId").String(),
		Items: Items{
			Songs:   GetSongs(s.Get("contents")),
			Albums:  GetTwoRowItem("album", a.Get("contents")),
			Artists: GetTwoRowItem("artist", u.Get("contents")),
		},
	}

	res, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
