package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/fengjx/go-kit-demo/greetsvc/endpoint"
)

func DecodeSumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.SumRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.ConcatRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
