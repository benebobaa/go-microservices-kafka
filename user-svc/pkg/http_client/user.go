package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Message string `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

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

func (r *UserClient) call(suffix, method string, request any, response any) error {

	jsonData, err := json.Marshal(request)
	if err != nil {
		return &ErrorResponse{
			Message: fmt.Sprintf("error marshalling request: %v", err),
		}
	}
	log.Println("URL: ", fmt.Sprintf("%s%s", r.url, suffix))

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", r.url, suffix), bytes.NewReader(jsonData))
	if err != nil {
		return &ErrorResponse{
			Message: fmt.Sprintf("error creating HTTP request: %v", err),
		}
	}

	res, err := r.client.Do(req)
	if err != nil {
		return &ErrorResponse{
			Message: fmt.Sprintf("error sending HTTP request: %v", err),
		}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &ErrorResponse{
			Message: fmt.Sprintf("error reading HTTP response body: %v", err),
		}
	}

	err = json.Unmarshal(body, response)

	if err != nil {
		return &ErrorResponse{
			Message: fmt.Sprintf("error unmarshalling response: %v", err),
		}
	}

	if res.StatusCode >= 400 {
		return &ErrorResponse{
			Message: fmt.Sprintf("error response: %s", string(body)),
		}
	}

	return nil
}

func (r *UserClient) GET(suffix string, request any, response any) error {
	return r.call(suffix, http.MethodGet, request, response)
}

func (r *UserClient) POST(suffix string, request any, response any) error {
	return r.call(suffix, http.MethodPost, request, response)
}
