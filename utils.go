package main

import (
	"fmt"
	"net/url"
	"time"
)

func ParseUrl(raw string) (string, error) {

	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	prehost := u.Host

	u.Host = "hyperpipe-proxy.onrender.com"
	q := u.Query()
	q.Set("host", prehost)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func ErrorMessage(err error) string {
	return fmt.Sprintf("{\"error\": \"%s\", \"message\": \"Please Report this error\"}", err)
}

func calc(url string) func() {
	start := time.Now()
	return func() {
		msg := fmt.Sprintf("Responded %v in %v\n", url, time.Since(start))
		fmt.Println(msg)
	}
}
