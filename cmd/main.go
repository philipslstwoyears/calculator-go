package main

import (
	"github.com/joho/godotenv"
	"github.com/philipslstwoyears/calculator-go/internal/agent"
	calc "github.com/philipslstwoyears/calculator-go/internal/calculator"
	"github.com/philipslstwoyears/calculator-go/internal/dto"
	"github.com/philipslstwoyears/calculator-go/internal/server"
	"github.com/philipslstwoyears/calculator-go/internal/storage"
	"log"
)

func main() {
	err := godotenv.Load("cmd/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	data := storage.New()
	ch := make(chan dto.Expression)
	workers := calc.New(data, ch)
	agent := agent.New(data, ch)
	serv := server.New()

	if err := workers.Start(); err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := agent.RunServer(); err != nil {
			log.Fatal(err)
		}
	}()
	if err := serv.RunServer(); err != nil {
		log.Fatal(err)
	}
}
