package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/konradreiche/apigen/api"
)

func (s *Server) GetPriceHandleFunc(w http.ResponseWriter, r *http.Request) {
	var req api.GetPriceRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		Encode(nil, w, nil, err)
	}
	resp, err := s.api.GetPrice(context.Background(), req)
	if err != nil {
		Encode(nil, w, nil, err)
	}
	Encode(nil, w, resp, nil)
}
