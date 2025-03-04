package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverMiddleware(t *testing.T) {
	handler := RecoverMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Test panic")
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rec.Code)
	}

	if body := rec.Body.String(); body != "Internal Server Error\n" {
		t.Errorf("Unexpected response body: got %q, want %q", body, "Internal Server Error\n")
	}
}

func TestLoggerMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	handler := LoggerMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	logOutput := buf.String()
	if logOutput == "" {
		t.Errorf("Expected log output, but got empty string")
	}
	if !bytes.Contains([]byte(logOutput), []byte("Запрос: GET /test")) {
		t.Errorf("Expected log to contain request method and path, got: %q", logOutput)
	}
}
