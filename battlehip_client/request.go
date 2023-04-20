package battlehip_client

import (
	"encoding/json"
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
