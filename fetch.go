package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func Fetch(path string, data []byte) (string, int, error) {

	url := "https://music.youtube.com/youtubei/v1/" + path + "?key=AIzaSyC9XL3ZjWddXya6X74dJoCTL-WEYFDNX30&prettyPrint=false"

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

func FetchArtist(id string) (string, int) {

	data, err := json.Marshal(GetTypeBrowse("artist", id))
	if err != nil {
		return ErrorMessage(err), 500
	}

	raw, status, err := Fetch("browse", data)
	if err != nil {
		return ErrorMessage(err), 500
	} else if status > 300 {
		return raw, status
	}

	res, err := ParseArtist(raw)
	if err != nil {
		return ErrorMessage(err), 500
	}

	return res, status
}

func FetchNext(id string) (string, int) {

	data, err := json.Marshal(GetTypeNext(id))
	if err != nil {
		return ErrorMessage(err), 500
	}

	raw, status, err := Fetch("next", data)
	if err != nil {
		return ErrorMessage(err), 500
	} else if status > 300 {
		return raw, status
	}

	res, err := ParseNext(raw)

	return res, status
}

/*func main() {

	res, status := FetchNext("KWLGyeg74es")

	log.Println(status)
	log.Println(res)
}*/
