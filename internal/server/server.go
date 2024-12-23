package server

import (
	"encoding/json"
	"fmt"
	calc "github.com/philipslstwoyears/calculator-go/internal/calculator"
	"github.com/philipslstwoyears/calculator-go/internal/dto"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"
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
	http.Handle("/", LoggerMiddleware(RecoverMiddleware(http.HandlerFunc(CalcHandler))))
	log.Println("Listening on port", a.config.Addr)

	return http.ListenAndServe(":"+a.config.Addr, nil)
}
func RecoverMiddleware(next http.Handler) http.Handler { // ловим панику
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v\n%s", err, debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWrapper{ResponseWriter: w}
		// Обработчик для записи статуса
		w = rw
		next.ServeHTTP(rw, r)

		logEntry := fmt.Sprintf(
			"Запрос: %s %s %s [%s] - Статус: %d - Время: %s",
			r.Method,
			r.URL.Path,
			r.Proto,
			r.RemoteAddr,
			rw.statusCode,
			time.Since(start),
		)
		log.Println(logEntry)
	})
}

type responseWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
