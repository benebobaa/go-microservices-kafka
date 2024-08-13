package http_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type UserClient struct {
	url    string
	client *http.Client
}

func NewUserClient(url string, timeout time.Duration) *UserClient {
	return &UserClient{
		url: url,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (r *UserClient) Call(method string, request any, response any) error {

	jsonData, err := json.Marshal(request)
	if err != nil {
		return errors.New("error marshalling HTTP request")
	}

	req, err := http.NewRequest(method, r.url, bytes.NewReader(jsonData))
	if err != nil {
		return errors.New("error creating HTTP request")
	}

	res, err := r.client.Do(req)
	if err != nil {
		return errors.New("error executing HTTP request")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New("error reading HTTP response")
	}

	err = json.Unmarshal(body, response)

	if err != nil {
		return errors.New("error unmarshalling HTTP response body")
	}

	return nil
}
