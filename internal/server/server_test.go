package server

//
//import (
//	"bytes"
//	"github.com/philipslstwoyears/calculator-go/internal/middleware"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestGoodRequest(t *testing.T) {
//	testCases := []struct {
//		GetBody        string
//		expectedResult string
//	}{
//		{`{"expression":"1+1"}`, `{"result":"2.00"}`},
//		{`{"expression":"11*(2+2)"}`, `{"result":"44.00"}`},
//		{`{"expression":"102/3"}`, `{"result":"34.00"}`},
//		{`{"expression":"-(-11-(1*20/2)-11/2*3)"}`, `{"result":"37.50"}`},
//		{`{"expression":"1*001*2/4"}`, `{"result":"0.50"}`},
//		{`{"expression":"1-12-1*(-1)"}`, `{"result":"-10.00"}`},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.GetBody, func(t *testing.T) {
//			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tc.GetBody)))
//			request.Header.Set("Content-Type", "application/json")
//
//			r := httptest.NewRecorder()
//			handler := http.HandlerFunc(CalcHandler)
//			middleware.LoggerMiddleware(middleware.RecoverMiddleware(handler)).ServeHTTP(r, request)
//
//			if r.Code != http.StatusOK {
//				t.Errorf("Wrong status code, expected %d, got: %d", http.StatusOK, r.Code)
//				return
//			}
//
//			rBody, err := io.ReadAll(r.Body)
//			if err != nil {
//				t.Errorf("Could not read body %v", err)
//			}
//			if string(rBody) != tc.expectedResult+"\n" {
//				t.Errorf("Wrong data, expected '%s', but got '%s'", tc.expectedResult, string(rBody))
//			}
//		})
//	}
//}
