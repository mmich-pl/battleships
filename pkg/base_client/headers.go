package base_client

import "net/http"

func getHeaders(headers ...http.Header) http.Header {
	if len(headers) > 0 {
		return headers[0]
	}
	return http.Header{}
}

func (c *BaseHTTPClient) getRequestHeaders(requestHeaders http.Header) http.Header {
	result := make(http.Header)

	// Add headers from the HTTP client instance
	for header, value := range c.Builder.Headers {
		if len(value) > 0 {
			result.Set(header, value[0])
		}
	}

	// Add headers specific for current request
	for header, value := range requestHeaders {
		if len(value) > 0 {
			result.Set(header, value[0])
		}
	}

	return result
}
