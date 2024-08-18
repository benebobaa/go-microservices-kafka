package http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Message string `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

type PaymentClient struct {
	url    string
	client *http.Client
}

func NewPaymentClient(url string, timeout time.Duration) *PaymentClient {
	return &PaymentClient{
		url: url,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (r *PaymentClient) call(ctx context.Context, suffix, method string, request any, response any) error {

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshalling request: %v", err)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", r.url, suffix), bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	res, err := r.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return fmt.Errorf("request cancelled: %v", ctx.Err())
		default:
			return fmt.Errorf("error sending HTTP request: %v", err)
		}
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(body, response)

	if err != nil {
		return fmt.Errorf("error unmarshalling response: %v", err)
	}

	if res.StatusCode == 429 {
		return fmt.Errorf("too many requests: %d", res.StatusCode)
	} else if res.StatusCode >= 500 {
		return fmt.Errorf("server error: %d", res.StatusCode)
	}

	return nil
}

func (r *PaymentClient) GET(ctx context.Context, suffix string, request any, response any) error {
	return r.call(ctx, suffix, http.MethodGet, request, response)
}

func (r *PaymentClient) POST(ctx context.Context, suffix string, request any, response any) error {
	return r.call(ctx, suffix, http.MethodPost, request, response)
}

func (r *PaymentClient) PATCH(ctx context.Context, suffix string, request any, response any) error {
	return r.call(ctx, suffix, http.MethodPatch, request, response)
}
