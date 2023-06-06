package base_client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type HTTPClientConfig struct {
	BaseUrl           string
	ConnectionTimeout time.Duration
	ResponseTimeout   time.Duration
	Headers           http.Header

	ProxyAddress string

	RetryWaitMin  int
	RetryWaitMax  int
	RetryMax      int
	CheckForRetry CheckForRetry
	Backoff       Backoff
}

type BaseHTTPClient struct {
	Config     HTTPClientConfig
	clientOnce sync.Once // make sure that BaseHTTPClient will be created only once
	client     *http.Client
}

type BaseClient interface {
	Get(endpoint string, headers ...http.Header) (*InternalResponse, error)
	Post(endpoint string, payload interface{}, headers ...http.Header) (*InternalResponse, error)
	Delete(url string, headers ...http.Header) (*InternalResponse, error)
}

func New(cfg HTTPClientConfig) *BaseHTTPClient {
	return &BaseHTTPClient{Config: cfg}
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

func (c *BaseHTTPClient) AddHeader(headerName, headerValue string) {
	if old := c.Config.Headers.Get(headerName); old != "" {
		c.Config.Headers.Del(headerName)
	}

	c.Config.Headers.Add(headerName, headerValue)
}

func (c *BaseHTTPClient) do(method, endpoint string, headers http.Header, body interface{}) (*InternalResponse, error) {
	buffer, err := c.marshalRequestBody(body)
	requestBody := bytes.NewReader(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create resp body: %w", err)
	}

	fullURL, err := url.JoinPath(c.Config.BaseUrl, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create full URL: %w", err)
	}

	request, err := NewRequest(method, fullURL, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", request)
	}

	request.request.Header = c.getRequestHeaders(headers)
	for i := 0; ; i++ {
		var responseCode int

		if request.body != nil {
			if _, err = request.body.Seek(0, 0); err != nil {
				return nil, fmt.Errorf("failed to seek body: %w", err)
			}
		}

		//	attempt to make request
		resp, err := c.getHttpClient().Do(request.request)
		log.Error(err)
		checkOK, checkErr := c.Config.CheckForRetry(resp, err)

		if err != nil {
			return nil, fmt.Errorf("failed to do request %s %s: %w", method, fullURL, err)
		}

		if !checkOK {
			return ConvertToInternal(resp)
		}

		switch checkErr {
		case nil:
			err = c.drainBody(resp.Body)
			if err != nil {
				log.Error(err)
				return nil, err
			}
		default:
			return nil, err
		}

		remain := c.Config.RetryMax - i
		i++
		if remain == 0 {
			break
		}

		wait := c.Config.Backoff(c.Config.RetryWaitMin, c.Config.RetryMax, i)
		desc := fmt.Sprintf("%s %s", method, fullURL)
		if responseCode > 0 {
			desc = fmt.Sprintf("%s (status: %d)", desc, responseCode)
		}
		log.Printf("%s, retrying in %s, (%d left)\n", desc, wait, remain)
		time.Sleep(wait)
	}
	return nil, fmt.Errorf("%s %s giving up after %d attempts", method, fullURL, c.Config.RetryMax+1)

}

func (c *BaseHTTPClient) drainBody(body io.ReadCloser) error {
	defer body.Close()
	_, err := io.Copy(io.Discard, body)
	if err != nil {
		return err
	}
	return nil
}

func createTLSConfig() *tls.Config {
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
		log.Error("No certs appended, using system certs only")
	}

	config := &tls.Config{
		RootCAs: rootCAs,
	}
	return config
}

func (c *BaseHTTPClient) getHttpClient() *http.Client {
	var transportConfig http.Transport

	if c.Config.ProxyAddress != "" {
		config := createTLSConfig()
		proxyUrl, _ := url.Parse(c.Config.ProxyAddress)
		transportConfig.Proxy = http.ProxyURL(proxyUrl)
		transportConfig.TLSClientConfig = config
	}

	c.clientOnce.Do(func() {
		c.client = &http.Client{
			Transport: &transportConfig,
			Timeout:   c.Config.ConnectionTimeout + c.Config.ConnectionTimeout,
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
