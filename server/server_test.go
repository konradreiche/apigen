package server

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/konradreiche/apigen/api"
	"github.com/konradreiche/apigen/client"
	"github.com/sirupsen/logrus"
)

func TestServer(t *testing.T) {
	logger := logrus.New()
	statsdClient, err := testClient()
	if err != nil {
		panic(err)
	}

	a := api.NewAPI()
	a = api.NewLoggingMiddleware(a, logger)
	a = api.NewInstrumentingMiddleware(a, statsdClient, logger)
	server := NewServer(a)
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

func testClient() (*statsd.Client, error) {
	addr := "localhost:1201"
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	server, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	defer server.Close()

	client, err := statsd.New(addr)
	if err != nil {
		return nil, err
	}
	return client, nil
}
