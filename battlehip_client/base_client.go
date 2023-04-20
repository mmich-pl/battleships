package battlehip_client

import "net/http"

type Client struct {
	client *http.Client
}

type HTTPClient interface {
	Get(endpoint string, headers ...http.Header) (*http.Response, error)
	Post(endpoint string, payload interface{}, headers ...http.Header) (*http.Response, error)
}
