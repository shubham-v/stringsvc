package responses

import (
	"context"
	"encoding/json"
	"net/http"
)

type UppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

func DecodeUppercaseResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response UppercaseResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

