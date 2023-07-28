package utils

import (
	"io"
	"net/http"
	"strings"
)

func HttpGet(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func HttpPost(url, content string) ([]byte, error) {

	response, err := http.Post(url, "application/json", strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
