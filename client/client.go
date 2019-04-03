package client

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

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

func EncodeVars(req interface{}, url string) string {
	v := reflect.ValueOf(req)
	t := reflect.TypeOf(req)
	for i := 0; i < t.NumField(); i++ {
		url = strings.Replace(url, fmt.Sprintf("{%s}", t.Field(i).Tag.Get("json")), v.Field(i).String(), 1)
	}
	return url
}
