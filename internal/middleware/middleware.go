package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type responseWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
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
