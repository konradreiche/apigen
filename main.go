package main

import (
	"github.com/konradreiche/apigen/api"
	"github.com/konradreiche/apigen/server"
)

func main() {
	api := api.NewAPI()
	server := server.NewServer(api)
	server.Serve()
}
