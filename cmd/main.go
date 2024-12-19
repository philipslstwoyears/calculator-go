package main

import (
	"github.com/philipslstwoyears/calculator-go/internal/server"
	"log"
)

func main() {
	app := server.New()
	err := app.RunServer()
	if err != nil {
		log.Fatal(err)
	}
}
