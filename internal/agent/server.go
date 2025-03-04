package agent

import (
	"github.com/gorilla/mux"
	"github.com/philipslstwoyears/calculator-go/internal/dto"
	"github.com/philipslstwoyears/calculator-go/internal/middleware"
	"github.com/philipslstwoyears/calculator-go/internal/storage"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT_AGENT")
	if config.Addr == "" {
		config.Addr = "8081"
	}
	return config
}

type Application struct {
	config  *Config
	storage *storage.Storage
	input   chan dto.Expression
}

func New(s *storage.Storage, input chan dto.Expression) *Application {
	return &Application{
		config:  ConfigFromEnv(),
		storage: s,
		input:   input,
	}
}
func (a *Application) RunServer() error {
	r := mux.NewRouter()
	r.HandleFunc("/internal/calculate", a.CalcHandler)
	r.HandleFunc("/internal/expressions", a.ExpressionsHandler)
	r.HandleFunc("/internal/expressions/{id}", a.ExpressionHandler)
	r.Use(middleware.LoggerMiddleware, middleware.RecoverMiddleware)
	log.Println("Listening on port", a.config.Addr)
	return http.ListenAndServe(":"+a.config.Addr, r)
}
