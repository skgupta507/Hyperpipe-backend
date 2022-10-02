package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

func Fetch(path string, data []byte) (string, int, error) {

	url := "https://music.youtube.com/youtubei/v1/" + path + "?alt=json&key=AIzaSyC9XL3ZjWddXya6X74dJoCTL-WEYFDNX30&prettyPrint=false"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", 500, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Origin", "https://music.youtube.com")
	req.Header.Set("x-origin", "https://music.youtube.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", 500, err
	}

	defer resp.Body.Close()

	log.Println(resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 500, err
	}

	return string(body), resp.StatusCode, nil
}

func FetchBrowse(browse BrowseData) (string, int) {

	data, err := json.Marshal(browse)
	if err != nil {
		return ErrorMessage(err), 500
	}

	raw, status, err := Fetch("browse", data)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return raw, status

}

func FetchExplore() (string, int) {

	id := "FEmusic_explore"

	context := GetTypeBrowse("", id, "", "")

	raw, status := FetchBrowse(context)

	res, err := ParseExplore(raw)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return res, status
}

func FetchGenres() (string, int) {

	id := "FEmusic_moods_and_genres"

	context := GetTypeBrowse("", id, "", "")

	raw, status := FetchBrowse(context)

	res, err := ParseGenres(raw)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return res, status
}

func FetchGenre(param string) (string, int) {

	id := "FEmusic_moods_and_genres_category"

	context := GetTypeBrowse("", id, param, "")

	raw, status := FetchBrowse(context)

	res, err := ParseGenre(raw)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return res, status
}

func FetchCharts(params, code string) (string, int) {

	id := "FEmusic_charts"

	context := GetTypeBrowse("", id, params, code)

	raw, status := FetchBrowse(context)

	res, err := ParseCharts(raw)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return res, status
}

func FetchArtist(id string) (string, int) {

	context := GetTypeBrowse("artist", id, "", "")

	raw, status := FetchBrowse(context)

	res, err := ParseArtist(raw)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return res, status
}

func FetchLyrics(id string) (string, int) {

	context := GetTypeBrowse("lyrics", id, "", "")

	raw, status := FetchBrowse(context)

	res, err := ParseLyrics(raw)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return res, status
}

func FetchAlbum(id string) (string, int) {

	context := GetTypeBrowse("album", id, "", "")

	raw, status := FetchBrowse(context)

	url := gjson.Parse(raw).Get("microformat.microformatDataRenderer.urlCanonical").String()

	data := "{\"canonicalUrl\": \"" +
		strings.Replace(url, "https://music.youtube.com", "", -1) +
		"\"}"

	return data, status
}

func FetchNext(id string) (string, int) {

	pldata, err := json.Marshal(GetTypeNext(id, ""))
	if err != nil {
		return ErrorMessage(err), 500
	}

	plraw, plstatus, err := Fetch("next", pldata)
	if err != nil || plstatus > 399 {
		return ErrorMessage(err), plstatus
	}

	pl := gjson.Parse(plraw).Get("contents." +
		"singleColumnMusicWatchNextResultsRenderer." +
		"tabbedRenderer.watchNextTabbedResultsRenderer.tabs.0.tabRenderer.content." +
		"musicQueueRenderer.content.playlistPanelRenderer.contents." +
		"#(automixPreviewVideoRenderer).automixPreviewVideoRenderer." +
		"content.automixPlaylistVideoRenderer.navigationEndpoint." +
		"watchPlaylistEndpoint.playlistId").String()

	data, err := json.Marshal(GetTypeNext(id, pl))
	if err != nil {
		return ErrorMessage(err), 500
	}

	raw, status, err := Fetch("next", data)

	res, err := ParseNext(raw)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return res, status
}
