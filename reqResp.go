package main

import (
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func makeRequest(method string, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		log.Errorf("Error with create request: %v\n", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+ycToken.IAM)
	return req, nil

}

func makeGetRequest(url string) (req *http.Request, err error) {
	return makeRequest("GET", url, nil)
}

func makePostRequest(url string) (req *http.Request, err error) {
	return makeRequest("POST", url, nil)
}

func doRequest(req *http.Request) (resp *http.Response, err error) {
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Errorf("Failed request: %v\n", err)
		return nil, err
	}
	return resp, nil
}
