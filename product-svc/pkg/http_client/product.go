package http_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type ProductClient struct {
	url    string
	client *http.Client
}

func NewProductClient(url string, timeout time.Duration) *ProductClient {
	return &ProductClient{
		url: url,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (r *ProductClient) call(suffix, method string, request any, response any) error {

	jsonData, err := json.Marshal(request)
	if err != nil {
		return errors.New("error marshalling HTTP request")
	}
	log.Println("URL: ", fmt.Sprintf("%s%s", r.url, suffix))

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", r.url, suffix), bytes.NewReader(jsonData))
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

func (r *ProductClient) GET(suffix string, request any, response any) error {
	return r.call(suffix, http.MethodGet, request, response)
}

func (r *ProductClient) POST(suffix string, request any, response any) error {
	return r.call(suffix, http.MethodPost, request, response)
}
