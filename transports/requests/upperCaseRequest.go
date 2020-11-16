package requests

import (
	"context"
	"encoding/json"
	"net/http"
)

type UppercaseRequest struct {
	S string `json:"s"`
}

func DecodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request UppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}