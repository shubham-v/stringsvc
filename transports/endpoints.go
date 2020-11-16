package transports

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"stringsvc/services"
	requests "stringsvc/transports/requests"
	responses "stringsvc/transports/responses"
)

func MakeCountEndpoint(svc services.StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(requests.CountRequest)
		v := svc.Count(req.S)
		return responses.CountResponse{v}, nil
	}
}

func MakeUppercaseEndpoint(svc services.StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(requests.UppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return responses.UppercaseResponse{v, err.Error()}, nil
		}
		return responses.UppercaseResponse{v, ""}, nil
	}
}

