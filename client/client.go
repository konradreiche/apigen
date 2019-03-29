package client

import (
	"net/http"

	"github.com/konradreiche/apigen/api"
)

type Client struct {
	conn     *http.Client
	endpoint string
}

func NewClient(endpoint string) api.API {
	client := http.DefaultClient
	return &Client{
		conn:     client,
		endpoint: endpoint,
	}
}
