package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/konradreiche/apigen/api"
	"github.com/konradreiche/apigen/client"
	"github.com/sirupsen/logrus"
)

func TestServer(t *testing.T) {
	server := NewServer(api.NewLoggingMiddleware(api.NewAPI(), logrus.New()))
	go func() {
		server.Serve()
	}()
	time.Sleep(1 * time.Second)

	client := client.NewClient("http://localhost:8080")
	fmt.Println(client.GetPrice(context.Background(), api.GetPriceRequest{
		AssetBase:  "USD",
		AssetQuote: "BTC",
	}))
}