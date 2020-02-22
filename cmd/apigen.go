package main

import (
	"log"
	"os"

	"github.com/konradreiche/apigen/parser"
)

func main() {
	p, err := parser.NewParser(os.Getenv("GOFILE"))
	check(err)
	err = p.Parse()
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
