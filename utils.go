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

	u.Host = "pipedproxy-bom.kavin.rocks"
	q := u.Query()
	q.Set("host", prehost)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func ErrorMessage(err error) string {
	return fmt.Sprintf("{\"error\": \"%s\", \"message\": \"Please Report this error\"}", err)
}

func calc() func() {
	start := time.Now()
	return func() {
		msg := fmt.Sprintf("Responded in %v\n", time.Since(start))
		fmt.Println(msg)
	}
}
