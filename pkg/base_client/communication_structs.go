package base_client

import (
	"encoding/json"
	"io"
	"net/http"
)

type InternalResponse struct {
	Status     string
	StatusCode int
	Headers    http.Header
	Body       []byte
}

func ConvertToInternal(resp *http.Response) (*InternalResponse, error) {
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	response := InternalResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       respBody,
	}
	return &response, nil

}

func (r *InternalResponse) UnmarshalJson(target interface{}) error {
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
