package server

import (
	"github.com/gorilla/mux"
	"github.com/philipslstwoyears/calculator-go/internal/middleware"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}
func (a *Application) RunServer() error {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/calculate", CalculateHandler)
	r.HandleFunc("/api/v1/expressions", ExpressionsHandler)
	r.HandleFunc("/api/v1/expressions/{id}", ExpressionHandler)
	r.Use(middleware.LoggerMiddleware, middleware.RecoverMiddleware)
	log.Println("Listening on port", a.config.Addr)
	return http.ListenAndServe(":"+a.config.Addr, r)
}
