package gsa

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type textbeltResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func alert(textbeltURL, number, message string) error {
	postBody := url.Values{}
	postBody.Set("message", message)
	postBody.Set("number", number)

	resp, err := http.PostForm(textbeltURL, postBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	s := textbeltResponse{}
	err = json.Unmarshal([]byte(body), &s)
	if err != nil {
		return err
	}

	if !s.Success {
		return errors.New(s.Message)
	}

	return nil
}
