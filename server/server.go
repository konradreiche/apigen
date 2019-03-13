package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/konradreiche/apigen/api"
)

func Serve() {
	http.HandleFunc("/price", GetPriceHandleFunc)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetPriceHandleFunc(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	var req api.GetPriceRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		EncodeError(enc, err)
		return
	}
	resp, err := api.GetPrice(context.Background(), req)
	if err != nil {
		EncodeError(enc, err)
		return
	}
	enc.Encode(resp)
}

type Error struct {
	Message string `json:"error"`
}

func EncodeError(enc *json.Encoder, err error) {
	enc.Encode(&Error{Message: err.Error()})
}
