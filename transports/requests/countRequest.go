package requests

import (
	"context"
	"encoding/json"
	"net/http"
)

type CountRequest struct {
	S string `json:"s"`
}

func DecodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request CountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
