package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type HttpClient struct {
	client *http.Client
}

func NewFetcher(proxy string) (*HttpClient, error) {

	transport := http.Transport{}

	if proxy != "" {

		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return nil, fmt.Errorf("Error parse proxy address: %s", err.Error())
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}

	return &HttpClient{
		client: &http.Client{
			Transport: &transport,
			Timeout:   time.Second * 15,
		},
	}, nil
}

func (f *HttpClient) Get(url string, headers map[string]string) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (f *HttpClient) Post(url string, headers map[string]string, payload interface{}) (*http.Response, error) {

	body, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(body)

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
