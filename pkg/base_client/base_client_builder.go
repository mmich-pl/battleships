package base_client

import (
	"net/http"
	"time"
)

type ClientBuilder interface {
	SetHeaderFromMap(headers map[string]string) ClientBuilder
	AddHeader(headerName, headerValue string) ClientBuilder
	SetConnectionTimeout(timeout time.Duration) ClientBuilder
	SetResponseTime(timeout time.Duration) ClientBuilder
	SetHttpClient(c *http.Client) ClientBuilder
	SetBaseURL(URL string) ClientBuilder
	Build() *BaseHTTPClient
}

type clientBuilder struct {
	Headers           http.Header
	connectionTimeout time.Duration
	responseTimeout   time.Duration
	baseUrl           string
	client            *http.Client
}

func (c *clientBuilder) AddHeader(headerName, headerValue string) ClientBuilder {
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

func (c *clientBuilder) SetResponseTime(timeout time.Duration) ClientBuilder {
	c.responseTimeout = timeout * time.Second
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
