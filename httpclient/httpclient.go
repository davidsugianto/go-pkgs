package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Headers    map[string]string
}

type Option func(*Client)

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.HTTPClient.Timeout = timeout
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(c *Client) {
		for k, v := range headers {
			c.Headers[k] = v
		}
	}
}

func New(baseURL string, opts ...Option) *Client {
	c := &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Headers: make(map[string]string),
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) makeRequest(ctx context.Context, method, endpoint string, body interface{}, contentType string) (*http.Response, error) {
	var bodyReader io.Reader

	switch v := body.(type) {
	case string:
		bodyReader = bytes.NewBufferString(v)
	case []byte:
		bodyReader = bytes.NewBuffer(v)
	case nil:
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonData)
		if contentType == "" {
			contentType = "application/json"
		}
	}

	fullURL := c.BaseURL + endpoint

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, err
	}

	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	if bodyReader != nil && contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(data))
	}
	return resp, nil
}

func (c *Client) Get(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	return c.makeRequest(ctx, http.MethodGet, endpoint, body, "application/json")
}

func (c *Client) Post(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	return c.makeRequest(ctx, http.MethodPost, endpoint, body, "application/json")
}

func (c *Client) PostRaw(ctx context.Context, endpoint string, rawBody string, contentType string) (*http.Response, error) {
	return c.makeRequest(ctx, http.MethodPost, endpoint, rawBody, contentType)
}

func (c *Client) Put(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	return c.makeRequest(ctx, http.MethodPut, endpoint, body, "application/json")
}

func (c *Client) PutRaw(ctx context.Context, endpoint string, rawBody string, contentType string) (*http.Response, error) {
	return c.makeRequest(ctx, http.MethodPut, endpoint, rawBody, contentType)
}

func (c *Client) Delete(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	return c.makeRequest(ctx, http.MethodDelete, endpoint, body, "application/json")
}
