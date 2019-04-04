package main

import (
	"github.com/konradreiche/apigen/api"
	"github.com/konradreiche/apigen/server"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	a := api.NewAPI()
	a = api.NewLoggingMiddleware(a, logger)
	server := server.NewServer(a)
	server.Serve()
}
