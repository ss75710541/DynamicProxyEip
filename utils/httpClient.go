package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func HttpRequestExec(method string, url string, body io.Reader, contentType string, token string) ([]byte, error) {
	req, err := HttpRequestCreate(method, url, body, contentType, token)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("HttpRequestCreate error: %s", err))
	}

	respBody, err := HttpRequestDo(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("HttpRequestDo error: %s", err))
	}
	return respBody, nil
}

func HttpRequestCreate(method string, url string, body io.Reader, contentType string, token string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	return req, nil
}

func HttpRequestDo(req *http.Request) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

