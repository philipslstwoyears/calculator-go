package server

import (
	"encoding/json"
	"fmt"
	calc "github.com/philipslstwoyears/calculator-go/internal/calculator"
	"github.com/philipslstwoyears/calculator-go/internal/dto"
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

type Request struct {
	Expression string `json:"expression"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		data := &dto.ErrorResponse{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(data)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := calc.Calc(request.Expression)
	if err != nil {
		data := &dto.ErrorResponse{
			Error: err.Error(),
		}
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(data)
	} else {
		data1 := &dto.ResultResponse{
			Result: fmt.Sprintf("%.2f", result),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data1)
	}
}
func (a *Application) RunServer() error {
	http.HandleFunc("/", CalcHandler)
	log.Println("Listening on port", a.config.Addr)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
