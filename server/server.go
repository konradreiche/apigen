package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	router := mux.NewRouter()
	router.Use(metadataMiddleware)
	router.HandleFunc(api.GetPriceEndpoint, s.GetPriceHandleFunc).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

type Error struct {
	Message string `json:"error,omitempty"`
}

func Encode(ctx context.Context, w io.Writer, data interface{}, err error) {
	encoder := json.NewEncoder(w)
	if err != nil {
		fmt.Println(err)
		encoder.Encode(&Error{Message: err.Error()})
		return
	}
	encoder.Encode(data)
}
