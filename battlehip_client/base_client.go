package battlehip_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
)

type Client struct {
	builder    *clientBuilder
	clientOnce sync.Once // make sure that Client will be created only once
	client     *http.Client
}

type HTTPClient interface {
	Get(endpoint string, headers ...http.Header) (*Response, error)
	Post(endpoint string, payload interface{}, headers ...http.Header) (*Response, error)
}

func (c *Client) Get(endpoint string, headers ...http.Header) (*Response, error) {
	return c.do(http.MethodGet, endpoint, getHeaders(headers...), nil)
}

func (c *Client) Post(endpoint string, payload interface{}, headers ...http.Header) (*Response, error) {
	return c.do(http.MethodPost, endpoint, getHeaders(headers...), payload)
}

func (c *Client) do(method, endpoint string, headers http.Header, body interface{}) (*Response, error) {
	fullHeaders := c.getRequestHeaders(headers)

	requestBody, err := c.marshalRequestBody(body)
	if err != nil {
		return nil, fmt.Errorf("failed to create resp body: %w", err)
	}

	fullURL, err := url.JoinPath(c.builder.baseUrl, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create full URL: %w", err)
	}

	request, err := http.NewRequest(method, fullURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", request)

	}

	request.Header = fullHeaders

	resp, err := c.getHttpClient().Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       respBody,
	}
	return &response, nil
}

func (c *Client) getHttpClient() *http.Client {
	c.clientOnce.Do(func() {
		if c.builder.client != nil {
			c.client = c.builder.client
			return
		}
		c.client = &http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: c.builder.responseTimeout,
				DialContext:           (&net.Dialer{Timeout: c.builder.connectionTimeout}).DialContext,
			},
			Timeout: c.builder.connectionTimeout + c.builder.responseTimeout,
		}
	})
	return c.client
}

func (c *Client) marshalRequestBody(body interface{}) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	return json.Marshal(body)
}
