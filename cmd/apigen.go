package main

import (
	"log"
	"os"

	"github.com/konradreiche/apigen/parser"
)

func main() {
	p, err := parser.NewParser(os.Getenv("GOFILE"))
	if err != nil {
		log.Fatal(err)
	}
	err = p.Parse()
	if err != nil {
		log.Fatal(err)
	}
}
