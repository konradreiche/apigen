package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/konradreiche/apigen/api"
)

type Server struct {
	api api.API
}

func NewServer(api api.API) *Server {
	return &Server{
		api: api,
	}
}

func (s *Server) Serve() {
	http.HandleFunc("/price", s.GetPriceHandleFunc)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Error struct {
	Message string `json:"error,omitempty"`
}

func Encode(ctx context.Context, w io.Writer, data interface{}, err error) {
	encoder := json.NewEncoder(w)
	if err != nil {
		encoder.Encode(&Error{Message: err.Error()})
	}
	encoder.Encode(data)
}
