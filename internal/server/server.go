package server

import (
	"github.com/gorilla/mux"
	"github.com/philipslstwoyears/calculator-go/internal/middleware"
	"github.com/philipslstwoyears/calculator-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	agent  proto.CalcServiceClient
}

func New() (*Application, error) {
	app := &Application{
		config: ConfigFromEnv(),
	}
	conn, err := grpc.NewClient("0.0.0.0:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	app.agent = proto.NewCalcServiceClient(conn)

	return app, nil
}
func (a *Application) RunServer() error {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/calculate", a.calculateHandler)
	r.HandleFunc("/api/v1/expressions", a.expressionsHandler)
	r.HandleFunc("/api/v1/expressions/{id}", a.expressionHandler)
	r.Use(middleware.LoggerMiddleware, middleware.RecoverMiddleware)
	log.Println("Listening on port", a.config.Addr)
	return http.ListenAndServe(":"+a.config.Addr, r)
}
