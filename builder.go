package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/tidwall/gjson"
)

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

func TwoRowItemRenderer(t string, a gjson.Result) []Item {

	r := []Item{}

	var id string

	if t == "album" {
		id = "menu.menuRenderer.items.#(menuNavigationItemRenderer" +
			".text.runs.0.text == Shuffle play).menuNavigationItemRenderer" +
			".navigationEndpoint.watchPlaylistEndpoint.playlistId"
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

				gid := j.Get(id).String()

				if len(gid) > 2 && gid[:2] == "VL" {
					gid = gid[2:]
				}

				r = append(r, Item{
					Id:    gid,
					Title: RunsText(j.Get("title")),
					Sub:   RunsText(j.Get("subtitle")),
					Thumbnails: GetThumbnails(j.Get("thumbnailRenderer" +
						".musicThumbnailRenderer.thumbnail.thumbnails")),
				})

			}()

			wg.Wait()

			return true
		},
	)

	return r
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
				flex := j.Get("flexColumns.#.musicResponsiveListItemFlexColumnRenderer" +
					".text.runs.0")

				r = append(r, Item{
					Id: j.Get("playlistItemData.videoId").String(),
					Title: flex.Get("#(navigationEndpoint.watchEndpoint" +
						".videoId).text").String(),
					SubId: flex.Get("#.navigationEndpoint.browseEndpoint" +
						".browseId").Get("0").String(),
					Sub: flex.Get("#(navigationEndpoint.browseEndpoint" +
						".browseId).text").String(),
					Thumbnails: GetThumbnails(j.Get("thumbnail.musicThumbnailRenderer" +
						".thumbnail.thumbnails")),
				})

			}()

			wg.Wait()

			return true
		},
	)

	return r
}

func ResponsiveListItemRendererCH(s gjson.Result) []Item {

	r := []Item{}

	wg := sync.WaitGroup{}

	s.ForEach(
		func(_, v gjson.Result) bool {

			wg.Add(1)

			go func() {
				defer wg.Done()

				j := v.Get("musicResponsiveListItemRenderer")
				flex := j.Get("flexColumns.#.musicResponsiveListItemFlexColumnRenderer" +
					".text.runs.0.text")

				r = append(r, Item{
					Id:    j.Get("navigationEndpoint.browseEndpoint.browseId").String(),
					Title: flex.Get("0").String() + " • " + flex.Get("1").String(),
					Thumbnails: GetThumbnails(j.Get("thumbnail.musicThumbnailRenderer" +
						".thumbnail.thumbnails")),
				})

			}()

			wg.Wait()

			return true
		},
	)

	return r
}

func NavigationButton(s gjson.Result) []Item {

	r := []Item{}

	wg := sync.WaitGroup{}

	s.ForEach(
		func(_, v gjson.Result) bool {

			wg.Add(1)

			go func() {
				defer wg.Done()

				j := v.Get("musicNavigationButtonRenderer")

				color := j.Get("solid.leftStripeColor").Uint() & 0xffffff

				r = append(r, Item{
					Id:    j.Get("clickCommand.browseEndpoint.params").String(),
					Title: RunsText(j.Get("buttonText")),
					Sub:   fmt.Sprintf("#%06x", color),
				})
			}()

			wg.Wait()

			return true
		},
	)

	return r
}

func MultiSelectMenuItemRenderer(j, ref gjson.Result) []Item {

	r := []Item{}

	wg := sync.WaitGroup{}

	j.ForEach(
		func(_, v gjson.Result) bool {

			wg.Add(1)

			go func() {
				defer wg.Done()

				steg := v.Get("formItemEntityKey").String()
				id := ref.Get("#(entityKey == " + steg + ")" +
					".payload.musicFormBooleanChoice.opaqueToken")

				r = append(r, Item{
					Id:    id.String(),
					Title: RunsText(v.Get("title")),
				})
			}()

			wg.Wait()

			return true
		},
	)

	return r
}
