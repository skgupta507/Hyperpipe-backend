package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"net/url"
	"os"
	"strings"
	"time"
)

func RunsText(j gjson.Result) string {

	var s []string

	a := j.Get("runs.#.text").Array()

	for i := 0; i < len(a); i++ {
		s = append(s, a[i].String())
	}

	return strings.Join(s, "")
}

func ParseUrl(raw string) (string, error) {

	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	prehost := u.Host

	u.Scheme = "https"
	u.Host = string(os.Getenv("HYP_PROXY"))

	q := u.Query()
	q.Set("host", prehost)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func ErrorMessage(err error) string {
	data := url.QueryEscape(err.Error())
	fmt.Println(err)
	return fmt.Sprintf("{\"error\":\"%s\",\"message\":\"Got Error: %s\"}", data)
}

func calc() func() {
	start := time.Now()
	return func() {
		msg := fmt.Sprintf("Responded in %v\n", time.Since(start))
		fmt.Println(msg)
	}
}
