package battlehip_client

import (
	"net/http"
	"time"
)

type ClientBuilder interface {
	SetHeaders(headers http.Header) ClientBuilder
	SetConnectionTimeout(timeout time.Duration) ClientBuilder
	SetResponseTime(timeout time.Duration) ClientBuilder
	SetHttpClient(c *http.Client) ClientBuilder

	Build() Client
}

type clientBuilder struct {
	headers           http.Header
	connectionTimeout time.Duration
	responseTimeout   time.Duration
	baseUrl           string
	client            *http.Client
}

func (c *clientBuilder) SetHttpClient(client *http.Client) ClientBuilder {
	c.client = client
	return c
}

func (c *clientBuilder) SetHeaders(headers http.Header) ClientBuilder {
	c.headers = headers
	return c
}

func (c *clientBuilder) SetConnectionTimeout(timeout time.Duration) ClientBuilder {
	c.connectionTimeout = timeout
	return c
}

func (c *clientBuilder) SetResponseTime(timeout time.Duration) ClientBuilder {
	c.responseTimeout = timeout
	return c
}

func (c *clientBuilder) Build() Client {
	client := Client{
		builder: c,
	}
	return client
}

func NewBuilder() ClientBuilder {
	return &clientBuilder{}
}
