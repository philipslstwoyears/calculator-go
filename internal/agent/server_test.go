package agent

import (
	"bytes"
	"encoding/json"
	"github.com/philipslstwoyears/calculator-go/internal/dto"
	"github.com/philipslstwoyears/calculator-go/internal/middleware"
	"github.com/philipslstwoyears/calculator-go/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalcHandler(t *testing.T) {
	testCases := []struct {
		GetBody        string
		expectedResult string
	}{
		{`{"expression":"5+5"}`, `{"id":0}`},
	}

	for _, tc := range testCases {
		t.Run(tc.GetBody, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/internal/calculate", bytes.NewBuffer([]byte(tc.GetBody)))
			request.Header.Set("Content-Type", "application/json")

			r := httptest.NewRecorder()
			data := storage.New()
			ch := make(chan dto.Expression, 2)
			a := New(data, ch)
			handler := http.HandlerFunc(a.CalcHandler)
			middleware.LoggerMiddleware(middleware.RecoverMiddleware(handler)).ServeHTTP(r, request)

			if r.Code != http.StatusOK {
				t.Errorf("Wrong status code, expected %d, got: %d", http.StatusOK, r.Code)
				return
			}

			rBody, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Could not read body %v", err)
			}
			if string(rBody) != tc.expectedResult+"\n" {
				t.Errorf("Wrong data, expected '%s', but got '%s'", tc.expectedResult, string(rBody))
			}
			d := data.GetAll()
			if len(d) != 1 {
				t.Errorf("Wrong data, expected length %d, but got %d", 1, len(d))
			}
		})
	}
}

func TestExpressionsHandler(t *testing.T) {
	testCases := []struct {
		Name string
		data []dto.Expression
	}{
		{`Ok`, []dto.Expression{
			{
				Id:     0,
				Status: "Ok",
				Result: 10,
			},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/internal/expressions", nil)
			request.Header.Set("Content-Type", "application/json")

			r := httptest.NewRecorder()
			data := storage.New()
			for _, datum := range tc.data {
				data.Add(datum)
			}
			ch := make(chan dto.Expression, 2)
			a := New(data, ch)
			handler := http.HandlerFunc(a.ExpressionsHandler)
			middleware.LoggerMiddleware(middleware.RecoverMiddleware(handler)).ServeHTTP(r, request)

			if r.Code != http.StatusOK {
				t.Errorf("Wrong status code, expected %d, got: %d", http.StatusOK, r.Code)
				return
			}

			rBody, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Could not read body %v", err)
			}
			marshal, err := json.Marshal(tc.data)
			if err != nil {
				return
			}
			if string(rBody) != string(marshal)+"\n" {
				t.Errorf("Wrong data, expected '%s', but got '%s'", string(marshal), string(rBody))
			}
		})
	}
}
