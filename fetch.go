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

func Fetch(path string, data []byte, c string) (string, int, error) {

	var query string

	if c != "" {
		query = "&ctoken=" + c + "&continuation=" + c + "&type=next"
	} else {
		query = ""
	}

	url := "https://music.youtube.com/youtubei/v1/" + path + "?key=AIzaSyC9XL3ZjWddXya6X74dJoCTL-WEYFDNX30" + query + "&prettyPrint=false"

	log.Println(url)

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

func FetchBrowse(id string, c string, browse BrowseData) (string, int) {

	data, err := json.Marshal(browse)
	if err != nil {
		return ErrorMessage(err), 500
	}

	raw, status, err := Fetch("browse", data, c)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return raw, status

}

func FetchHome(c string) (string, int) {

	/* WIP */

	id := "FEmusic_home"

	if c != "" {
		id = ""
	}

	log.Println(id)

	context := GetTypeBrowse("home", id)

	raw, status := FetchBrowse(id, c, context)

	/*res, err := ParseHome(raw)
	if err != nil {
		return ErrorMessage(err), 500
	}*/

	return raw, status
}

func FetchExplore() (string, int) {

	id := "FEmusic_explore"

	log.Println(id)

	context := GetTypeBrowse("", id)

	raw, status := FetchBrowse(id, "", context)

	res, err := ParseExplore(raw)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return res, status
}

func FetchMoods() (string, int) {

	/* WIP */

	id := "FEmusic_moods_and_genres"

	log.Println(id)

	context := GetTypeBrowse("", id)

	raw, status := FetchBrowse(id, "", context)

	return raw, status
}

func FetchArtist(id string) (string, int) {

	context := GetTypeBrowse("artist", id)

	raw, status := FetchBrowse(id, "", context)

	res, err := ParseArtist(raw)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return res, status
}

func FetchLyrics(id string) (string, int) {

	context := GetTypeBrowse("lyrics", id)

	raw, status := FetchBrowse(id, "", context)

	res, err := ParseLyrics(raw)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return res, status
}

func FetchPlaylist(id string) (string, int) {

	context := GetTypeBrowse("playlist", id)

	raw, status := FetchBrowse(id, "", context)

	return raw, status
}

func FetchAlbum(id string) (string, int) {

	context := GetTypeBrowse("album", id)

	raw, status := FetchBrowse(id, "", context)

	url := gjson.Parse(raw).Get("microformat.microformatDataRenderer.urlCanonical").String()

	data := "{\"canonicalUrl\": \"" +
		strings.Replace(url, "https://music.youtube.com", "", -1) +
		"\"}"

	return data, status
}

func FetchNext(id string) (string, int) {

	data, err := json.Marshal(GetTypeNext(id))
	if err != nil {
		return ErrorMessage(err), 500
	}

	raw, status, err := Fetch("next", data, "")
	if err != nil {
		return ErrorMessage(err), 500
	}

	res, err := ParseNext(raw)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return res, status
}
