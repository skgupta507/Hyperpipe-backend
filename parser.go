package main

import (
	"encoding/json"
	"sync"

	"github.com/tidwall/gjson"
)

func GetNextSongs(n gjson.Result) []Item {

	r := []Item{}

	wg := sync.WaitGroup{}

	n.Get("#.playlistPanelVideoRenderer").ForEach(
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

func ParseExplore(raw string) (string, error) {

	j := gjson.Parse(raw)

	m := j.Get("contents.singleColumnBrowseResultsRenderer.tabs.0" +
		".tabRenderer.content.sectionListRenderer.contents")

	c := m.Get("#.musicCarouselShelfRenderer")

	a := c.Get("#(header.musicCarouselShelfBasicHeaderRenderer" +
		".title.runs.0.text == New albums & singles)")
	t := c.Get("#(header.musicCarouselShelfBasicHeaderRenderer" +
		".title.runs.0.text == Trending)")

	charts := m.Get("#(gridRenderer).gridRenderer.items.#" +
		".musicNavigationButtonRenderer").Get("#(buttonText" +
		".runs.0.text == Charts)")

	val := Explore{
		TrendingId: t.Get("header.musicCarouselShelfBasicHeaderRenderer" +
			".title.runs.0.navigationEndpoint.browseEndpoint.browseId").String(),
		ChartsId: charts.Get("clickCommand.browseEndpoint.params").String(),
		Albums:   TwoRowItemRenderer("album", a.Get("contents")),
		Trending: ResponsiveListItemRenderer(t.Get("contents")),
	}

	res, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func ParseGenres(raw string) (string, error) {

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

func ParseGenre(raw string) (string, error) {

	j := gjson.Parse(raw)

	c := j.Get("contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer" +
		".content.sectionListRenderer.contents.#.gridRenderer")

	s := c.Get("#(header.gridHeaderRenderer.title.runs.0.text == Spotlight)")
	f := c.Get("#(header.gridHeaderRenderer.title.runs.0.text == Featured playlists)")
	cp := c.Get("#(header.gridHeaderRenderer.title.runs.0.text == Community playlists)")

	val := Genre{
		Title:     RunsText(j.Get("header.musicHeaderRenderer.title")),
		Spotlight: TwoRowItemRenderer("", s.Get("items")),
		Featured:  TwoRowItemRenderer("", f.Get("items")),
		Community: TwoRowItemRenderer("", cp.Get("items")),
	}

	res, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func ParseCharts(raw string) (string, error) {

	j := gjson.Parse(raw)

	c := j.Get("contents.singleColumnBrowseResultsRenderer" +
		".tabs.0.tabRenderer.content.sectionListRenderer.contents")

	o := c.Get("0.musicShelfRenderer.subheaders.0" +
		".musicSideAlignedItemRenderer.startItems.0" +
		".musicSortFilterButtonRenderer")

	a := c.Get("#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer" +
		".title.runs.0.text == Top artists)" +
		".musicCarouselShelfRenderer.contents")

	t := c.Get("#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer" +
		".title.runs.0.text == Trending)" +
		".musicCarouselShelfRenderer.contents")

	opts := o.Get("menu.musicMultiSelectMenuRenderer.options" +
		".#.musicMultiSelectMenuItemRenderer")
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

func ParseArtist(raw string) (string, error) {

	j := gjson.Parse(raw)

	h := j.Get("header.musicImmersiveHeaderRenderer")
	c := j.Get("contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer.content" +
		".sectionListRenderer.contents")

	s := c.Get("#(musicShelfRenderer.title.runs.0.text == Songs).musicShelfRenderer")
	a := c.Get("#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer" +
		".title.runs.0.text == Albums).musicCarouselShelfRenderer")
	m := c.Get("#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer" +
		".title.runs.0.text == Singles).musicCarouselShelfRenderer")
	u := c.Get("#(musicCarouselShelfRenderer.header.musicCarouselShelfBasicHeaderRenderer" +
		".title.runs.0.text == Fans might also like).musicCarouselShelfRenderer")

	val := Artist{
		Title:       RunsText(h.Get("title")),
		Description: RunsText(h.Get("description")),
		SubscriberCount: RunsText(h.Get("subscriptionButton" +
			".subscribeButtonRenderer.subscriberCountText")),
		Thumbnails: GetThumbnails(h.Get("thumbnail.musicThumbnailRenderer" +
			".thumbnail.thumbnails")),
		BrowsePlaylistId: h.Get("playButton.buttonRenderer.navigationEndpoint" +
			".watchEndpoint.playlistId").String(),
		PlaylistId: s.Get("contents.0.musicResponsiveListItemRenderer" +
			".flexColumns.0.musicResponsiveListItemFlexColumnRenderer.text" +
			".runs.0.navigationEndpoint.watchEndpoint.playlistId").String(),
		Items: Items{
			Songs:   ResponsiveListItemRenderer(s.Get("contents")),
			Albums:  TwoRowItemRenderer("album", a.Get("contents")),
			Singles: TwoRowItemRenderer("album", m.Get("contents")),
			Artists: TwoRowItemRenderer("", u.Get("contents")),
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

	c := j.Get("contents.singleColumnMusicWatchNextResultsRenderer." +
		"tabbedRenderer.watchNextTabbedResultsRenderer.tabs")
	m := j.Get("playerOverlays.playerOverlayRenderer.browserMediaSession." +
		"browserMediaSessionRenderer")

	n := c.Get("#(tabRenderer.title == Up next).tabRenderer.content" +
		".musicQueueRenderer.content.playlistPanelRenderer")
	l := c.Get("#(tabRenderer.title == Lyrics).tabRenderer.endpoint" +
		".browseEndpoint.browseId")

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
