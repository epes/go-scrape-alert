package gsa

import (
	"io/ioutil"
	"net/http"
	"time"
)

func scrape(scrapeURL string, filter func([]byte, string) (string, error)) (string, error) {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := c.Get(scrapeURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	msg, err := filter(body, scrapeURL)
	if err != nil {
		return "", err
	}

	return msg, nil
}
