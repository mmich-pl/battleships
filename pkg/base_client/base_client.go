package base_client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
)

type BaseHTTPClient struct {
	Builder    *clientBuilder
	clientOnce sync.Once // make sure that BaseHTTPClient will be created only once
	client     *http.Client
}

type BaseClient interface {
	Get(endpoint string, headers ...http.Header) (*Response, error)
	Post(endpoint string, payload interface{}, headers ...http.Header) (*Response, error)
	Delete(url string, headers ...http.Header) (*Response, error)
}

func (c *BaseHTTPClient) Get(endpoint string, headers ...http.Header) (*Response, error) {
	return c.do(http.MethodGet, endpoint, getHeaders(headers...), nil)
}

func (c *BaseHTTPClient) Post(endpoint string, payload interface{}, headers ...http.Header) (*Response, error) {
	return c.do(http.MethodPost, endpoint, getHeaders(headers...), payload)
}

func (c *BaseHTTPClient) Delete(url string, headers ...http.Header) (*Response, error) {
	return c.do(http.MethodDelete, url, getHeaders(headers...), nil)
}

func (c *BaseHTTPClient) do(method, endpoint string, headers http.Header, body interface{}) (*Response, error) {
	fullHeaders := c.getRequestHeaders(headers)

	requestBody, err := c.marshalRequestBody(body)
	if err != nil {
		return nil, fmt.Errorf("failed to create resp body: %w", err)
	}

	fullURL, err := url.JoinPath(c.Builder.baseUrl, endpoint)
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

func CreateTLSConfig() *tls.Config {
	insecure := flag.Bool("insecure-ssl", false, "Accept/Ignore all server SSL certificates")
	flag.Parse()

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	cert := os.Getenv("CERTIFICATE_PATH")
	certs, err := os.ReadFile(cert)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", cert, err)
	}

	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}

	config := &tls.Config{
		InsecureSkipVerify: *insecure,
		RootCAs:            rootCAs,
	}
	return config
}
