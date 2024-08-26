package utils

import (
	"net/http"
)

// MakeGetRequest makes an HTTP GET request to the specified URL
func MakeGetRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
