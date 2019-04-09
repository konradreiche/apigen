package main

import (
	"net"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/konradreiche/apigen/api"
	"github.com/konradreiche/apigen/server"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	statsdClient, err := testClient()
	if err != nil {
		panic(err)
	}

	a := api.NewAPI()
	a = api.NewLoggingMiddleware(a, logger)
	a = api.NewInstrumentingMiddleware(a, statsdClient, logger)
	server := server.NewServer(a)
	server.Serve()
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
