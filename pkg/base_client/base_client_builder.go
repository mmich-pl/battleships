package base_client

import (
	"net/http"
	"time"
)

type ClientBuilder interface {
	SetHeaderFromMap(headers map[string]string) ClientBuilder
	AddHeader(headerName, headerValue string) ClientBuilder
	SetConnectionTimeout(timeout time.Duration) ClientBuilder
	SetResponseTimeout(timeout time.Duration) ClientBuilder
	SetHttpClient(c *http.Client) ClientBuilder
	SetBaseURL(URL string) ClientBuilder
	SetRetryWaitMaxTime(duration int) ClientBuilder
	SetRetryWaitMinTime(duration int) ClientBuilder
	SetRetryMaxAttempts(attempts int) ClientBuilder
	SetRetryCheck(check CheckForRetry) ClientBuilder
	SetBackoff(backoff Backoff) ClientBuilder
	Build() *BaseHTTPClient
}

type clientBuilder struct {
	Headers           http.Header
	connectionTimeout time.Duration
	responseTimeout   time.Duration
	baseUrl           string
	client            *http.Client
	retryWaitMin      int
	retryWaitMax      int
	retryMax          int

	checkForRetry CheckForRetry
	backoff       Backoff
}

func (c *clientBuilder) AddHeader(headerName, headerValue string) ClientBuilder {
	if old := c.Headers.Get(headerName); old != "" {
		c.Headers.Del(headerName)
	}

	c.Headers.Add(headerName, headerValue)
	return c
}

func (c *clientBuilder) SetHeaderFromMap(headers map[string]string) ClientBuilder {
	headersList := make(http.Header)
	for key, val := range headers {
		headersList.Set(key, val)
	}
	c.Headers = headersList
	return c
}

func (c *clientBuilder) SetBaseURL(URL string) ClientBuilder {
	c.baseUrl = URL
	return c
}

func (c *clientBuilder) SetHttpClient(client *http.Client) ClientBuilder {
	c.client = client
	return c
}

func (c *clientBuilder) SetConnectionTimeout(timeout time.Duration) ClientBuilder {
	c.connectionTimeout = timeout * time.Second
	return c
}

func (c *clientBuilder) SetResponseTimeout(timeout time.Duration) ClientBuilder {
	c.responseTimeout = timeout * time.Second
	return c
}

func (c *clientBuilder) SetRetryWaitMaxTime(duration int) ClientBuilder {
	c.retryWaitMax = duration
	return c
}

func (c *clientBuilder) SetRetryWaitMinTime(duration int) ClientBuilder {
	c.retryWaitMin = duration
	return c
}

func (c *clientBuilder) SetRetryMaxAttempts(attempts int) ClientBuilder {
	c.retryMax = attempts
	return c
}

func (c *clientBuilder) SetRetryCheck(check CheckForRetry) ClientBuilder {
	c.checkForRetry = check
	return c
}

func (c *clientBuilder) SetBackoff(backoff Backoff) ClientBuilder {
	c.backoff = backoff
	return c
}

func (c *clientBuilder) Build() *BaseHTTPClient {
	client := BaseHTTPClient{
		Builder: c,
	}
	return &client
}

func NewBuilder() ClientBuilder {
	return &clientBuilder{}
}
