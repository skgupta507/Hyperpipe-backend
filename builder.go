package main

import (
	"sync"

	"github.com/tidwall/gjson"
)

func TwoRowItemRenderer(t string, a gjson.Result) []Item {

	r := []Item{}

	var id string

	if t == "album" || t == "singles" {
		id = "menu.menuRenderer.items.#(menuNavigationItemRenderer.text.runs.0.text == Shuffle play).menuNavigationItemRenderer.navigationEndpoint.watchPlaylistEndpoint.playlistId"
	} else {
		id = "navigationEndpoint.browseEndpoint.browseId"
	}

	wg := sync.WaitGroup{}

	v := a.Get("#.musicTwoRowItemRenderer")

	v.ForEach(
		func(_, j gjson.Result) bool {

			wg.Add(1)

			go func() {
				defer wg.Done()

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

func CarouselShelfRenderer(i gjson.Result) map[string]interface{} {

	shelf := make(map[string]interface{})

	shelf["title"] = RunsText(i.Get("header.musicCarouselShelfBasicHeaderRenderer.title"))

	ct := i.Get("contents")

	if ct.Get("#(musicTwoRowItemRenderer)").Exists() {
		shelf["contents"] = TwoRowItemRenderer("", ct)
	} else if ct.Get("#(musicResponsiveListItemRenderer)").Exists() {
		shelf["contents"] = ResponsiveListItemRenderer(ct)
	}

	shelf["raw"] = ct.Value()

	return shelf
}

func ResponsiveListItemRenderer(s gjson.Result) []Item {

	r := []Item{}

	wg := sync.WaitGroup{}

	s.ForEach(
		func(_, v gjson.Result) bool {

			wg.Add(1)

			go func() {
				defer wg.Done()

				j := v.Get("musicResponsiveListItemRenderer")
				flex := j.Get("flexColumns.#.musicResponsiveListItemFlexColumnRenderer.text.runs.0")

				r = append(r, Item{
					Id:         j.Get("playlistItemData.videoId").String(),
					Title:      flex.Get("#(navigationEndpoint.watchEndpoint.videoId).text").String(),
					SubId:      flex.Get("#.navigationEndpoint.browseEndpoint").Get("#(browseId).browseId").String(),
					Sub:        flex.Get("#(navigationEndpoint.browseEndpoint.browseId).text").String(),
					Thumbnails: GetThumbnails(j.Get("thumbnail.musicThumbnailRenderer.thumbnail.thumbnails")),
				})

			}()

			wg.Wait()

			return true
		},
	)

	return r
}
