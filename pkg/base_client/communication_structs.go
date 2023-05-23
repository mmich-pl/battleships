package base_client

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	Status     string
	StatusCode int
	Headers    http.Header
	Body       []byte
}

func (r *Response) UnmarshalJson(target interface{}) error {
	return json.Unmarshal(r.Body, target)
}

type Request struct {
	body    io.ReadSeeker
	request *http.Request
}

func NewRequest(method, url string, body io.ReadSeeker) (*Request, error) {
	var rcBody io.ReadCloser
	if body != nil {
		rcBody = io.NopCloser(body)
	}

	httpReq, err := http.NewRequest(method, url, rcBody)
	if err != nil {
		return nil, err
	}
	return &Request{body, httpReq}, nil
}
