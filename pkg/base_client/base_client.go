package base_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type BaseHTTPClient struct {
	Builder    *clientBuilder
	clientOnce sync.Once // make sure that BaseHTTPClient will be created only once
	client     *http.Client
}

type BaseClient interface {
	Get(endpoint string, headers ...http.Header) (*InternalResponse, error)
	Post(endpoint string, payload interface{}, headers ...http.Header) (*InternalResponse, error)
	Delete(url string, headers ...http.Header) (*InternalResponse, error)
}

func (c *BaseHTTPClient) Get(endpoint string, headers ...http.Header) (*InternalResponse, error) {
	return c.do(http.MethodGet, endpoint, getHeaders(headers...), nil)
}

func (c *BaseHTTPClient) Post(endpoint string, payload interface{}, headers ...http.Header) (*InternalResponse, error) {
	return c.do(http.MethodPost, endpoint, getHeaders(headers...), payload)
}

func (c *BaseHTTPClient) Delete(url string, headers ...http.Header) (*InternalResponse, error) {
	return c.do(http.MethodDelete, url, getHeaders(headers...), nil)
}

func (c *BaseHTTPClient) do(method, endpoint string, headers http.Header, body interface{}) (*InternalResponse, error) {
	// request setup
	buffer, err := c.marshalRequestBody(body)
	requestBody := bytes.NewReader(buffer)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("failed to create resp body: %w", err)
	}

	fullURL, err := url.JoinPath(c.Builder.baseUrl, endpoint)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("failed to create full URL: %w", err)
	}

	request, err := NewRequest(method, fullURL, requestBody)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("failed to create request: %w", request)
	}

	request.request.Header = c.getRequestHeaders(headers)
	var i int
	for {
		var responseCode int

		if request.body != nil {
			if _, err = request.body.Seek(0, 0); err != nil {
				return nil, fmt.Errorf("failed to seek body: %w", err)
			}
		}

		//	attempt to make request
		resp, err := c.getHttpClient().Do(request.request)
		checkOK, checkErr := c.Builder.checkForRetry(resp, err)

		if err != nil {
			log.Error(err)
			return nil, fmt.Errorf("failed to do request %s %s: %w", method, fullURL, err)
		}

		if !checkOK {
			return ConvertToInternal(resp)
		}

		switch checkErr {
		case nil:
			err = c.drainBody(resp.Body)
			if err != nil {
				return nil, err
			}
		default:
			return nil, err
		}

		remain := c.Builder.retryMax - 1
		i++
		if remain == 0 {
			break
		}

		wait := c.Builder.backoff(c.Builder.retryWaitMin, c.Builder.retryWaitMax, i)
		desc := fmt.Sprintf("%s %s", method, fullURL)
		if responseCode > 0 {
			desc = fmt.Sprintf("%s (status: %d)", desc, responseCode)
		}
		log.Printf("%s, retrying in %s, (%d left)\n", desc, wait, remain)
		time.Sleep(wait)
	}
	return nil, fmt.Errorf("%s %s giving up after %d attempts", method, fullURL, c.Builder.retryMax+1)

}

func (c *BaseHTTPClient) drainBody(body io.ReadCloser) error {
	defer body.Close()
	_, err := io.Copy(io.Discard, body)
	if err != nil {
		return err
	}
	return nil
}

func (c *BaseHTTPClient) getHttpClient() *http.Client {
	c.clientOnce.Do(func() {
		if c.Builder.client != nil {
			c.client = c.Builder.client
			return
		}
		c.client = &http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: c.Builder.responseTimeout,
				DialContext:           (&net.Dialer{Timeout: c.Builder.connectionTimeout}).DialContext,
			},
			Timeout: c.Builder.connectionTimeout + c.Builder.responseTimeout,
		}
	})
	return c.client
}

func (c *BaseHTTPClient) marshalRequestBody(body interface{}) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	return json.Marshal(body)
}
