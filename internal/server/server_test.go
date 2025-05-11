package server

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/philipslstwoyears/calculator-go/internal/mocks"
	"github.com/philipslstwoyears/calculator-go/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name           string
		cookie         *http.Cookie
		body           string
		mockBehavior   func(m *mocks.MockCalcServiceClient)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "Success",
			cookie: &http.Cookie{Name: "id", Value: "1"},
			body:   `{"expression": "5+5"}`,
			mockBehavior: func(m *mocks.MockCalcServiceClient) {
				m.EXPECT().Calc(gomock.Any(), &proto.Request{Expression: "5+5", UserId: 1}).Return(&proto.Id{Id: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"id": float64(1)}, // JSON декодирует int32 как float64
		},
		{
			name:           "NoCookie",
			cookie:         nil,
			body:           `{"expression": "5+5"}`,
			mockBehavior:   func(m *mocks.MockCalcServiceClient) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "http: named cookie not present"},
		},
		{
			name:           "InvalidCookie",
			cookie:         &http.Cookie{Name: "id", Value: "invalid"},
			body:           `{"expression": "5+5"}`,
			mockBehavior:   func(m *mocks.MockCalcServiceClient) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "strconv.Atoi: parsing \"invalid\": invalid syntax"},
		},
		{
			name:           "InvalidJSON",
			cookie:         &http.Cookie{Name: "id", Value: "1"},
			body:           `{"expression": "5+5"`, // Неверный JSON
			mockBehavior:   func(m *mocks.MockCalcServiceClient) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "unexpected EOF"},
		},
		{
			name:   "AgentError",
			cookie: &http.Cookie{Name: "id", Value: "1"},
			body:   `{"expression": "5+5"}`,
			mockBehavior: func(m *mocks.MockCalcServiceClient) {
				m.EXPECT().Calc(gomock.Any(), &proto.Request{Expression: "5+5", UserId: 1}).Return(nil, status.Error(codes.InvalidArgument, "invalid expression"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "rpc error: code = InvalidArgument desc = invalid expression"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAgent := mocks.NewMockCalcServiceClient(ctrl)
			test.mockBehavior(mockAgent)

			app := &Application{agent: mockAgent}
			req := httptest.NewRequest(http.MethodPost, "/calculate", strings.NewReader(test.body))
			if test.cookie != nil {
				req.AddCookie(test.cookie)
			}
			rr := httptest.NewRecorder()

			app.CalculateHandler(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)
			if test.expectedStatus == http.StatusOK {
				assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
			}

			var actualBody map[string]interface{}
			err := json.NewDecoder(rr.Body).Decode(&actualBody)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedBody, actualBody)
		})
	}
}

func TestExpressionsHandler(t *testing.T) {
	tests := []struct {
		name           string
		cookie         *http.Cookie
		mockBehavior   func(m *mocks.MockCalcServiceClient)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "NoCookie",
			cookie:         nil,
			mockBehavior:   func(m *mocks.MockCalcServiceClient) {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "http: named cookie not present"},
		},
		{
			name:           "InvalidCookie",
			cookie:         &http.Cookie{Name: "id", Value: "invalid"},
			mockBehavior:   func(m *mocks.MockCalcServiceClient) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "strconv.Atoi: parsing \"invalid\": invalid syntax"},
		},
		{
			name:   "AgentError",
			cookie: &http.Cookie{Name: "id", Value: "1"},
			mockBehavior: func(m *mocks.MockCalcServiceClient) {
				m.EXPECT().GetExpressions(gomock.Any(), &proto.Id{Id: 1}).Return(nil, status.Error(codes.Internal, "db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "rpc error: code = Internal desc = db error"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAgent := mocks.NewMockCalcServiceClient(ctrl)
			test.mockBehavior(mockAgent)

			app := &Application{agent: mockAgent}
			req := httptest.NewRequest(http.MethodGet, "/expressions", nil)
			if test.cookie != nil {
				req.AddCookie(test.cookie)
			}
			rr := httptest.NewRecorder()

			app.expressionsHandler(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)
			if test.expectedStatus == http.StatusOK {
				assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
			}

			var actualBody interface{}
			if test.expectedStatus == http.StatusOK {
				actualBody = make([]map[string]interface{}, 0)
			} else {
				actualBody = make(map[string]interface{})
			}
			err := json.NewDecoder(rr.Body).Decode(&actualBody)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedBody, actualBody)
		})
	}
}

func TestExpressionHandler(t *testing.T) {
	tests := []struct {
		name           string
		cookie         *http.Cookie
		url            string
		mockBehavior   func(m *mocks.MockCalcServiceClient)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "Success",
			cookie: &http.Cookie{Name: "id", Value: "1"},
			url:    "/expression/1",
			mockBehavior: func(m *mocks.MockCalcServiceClient) {
				m.EXPECT().GetExpression(gomock.Any(), &proto.Id{Id: 1}).Return(&proto.Expression{
					Id:         1,
					Expression: "5+5",
					Status:     "completed",
					Result:     10,
					UserId:     1,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":         float64(1),
				"expression": "5+5",
				"status":     "completed",
				"result":     float64(10),
				"user_id":    float64(1),
			},
		},
		{
			name:           "NoCookie",
			cookie:         nil,
			url:            "/expression/1",
			mockBehavior:   func(m *mocks.MockCalcServiceClient) {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "http: named cookie not present"},
		},
		{
			name:           "InvalidCookie",
			cookie:         &http.Cookie{Name: "id", Value: "invalid"},
			url:            "/expression/1",
			mockBehavior:   func(m *mocks.MockCalcServiceClient) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "strconv.Atoi: parsing \"invalid\": invalid syntax"},
		},
		{
			name:           "InvalidID",
			cookie:         &http.Cookie{Name: "id", Value: "1"},
			url:            "/expression/invalid",
			mockBehavior:   func(m *mocks.MockCalcServiceClient) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "strconv.Atoi: parsing \"invalid\": invalid syntax"},
		},
		{
			name:   "AgentError",
			cookie: &http.Cookie{Name: "id", Value: "1"},
			url:    "/expression/1",
			mockBehavior: func(m *mocks.MockCalcServiceClient) {
				m.EXPECT().GetExpression(gomock.Any(), &proto.Id{Id: 1}).Return(nil, status.Error(codes.NotFound, "expression not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "rpc error: code = NotFound desc = expression not found"},
		},
		{
			name:   "Forbidden",
			cookie: &http.Cookie{Name: "id", Value: "1"},
			url:    "/expression/1",
			mockBehavior: func(m *mocks.MockCalcServiceClient) {
				m.EXPECT().GetExpression(gomock.Any(), &proto.Id{Id: 1}).Return(&proto.Expression{
					Id:         1,
					Expression: "5+5",
					Status:     "completed",
					Result:     10,
					UserId:     2, // Другой пользователь
				}, nil)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   map[string]interface{}{"error": "It is not your expression"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAgent := mocks.NewMockCalcServiceClient(ctrl)
			test.mockBehavior(mockAgent)

			app := &Application{agent: mockAgent}
			req := httptest.NewRequest(http.MethodGet, test.url, nil)
			if test.cookie != nil {
				req.AddCookie(test.cookie)
			}
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/expression/{id}", app.expressionHandler).Methods(http.MethodGet)
			router.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)
			if test.expectedStatus == http.StatusOK {
				assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
			}

			var actualBody map[string]interface{}
			err := json.NewDecoder(rr.Body).Decode(&actualBody)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedBody, actualBody)
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockBehavior   func(m *mocks.MockCalcServiceClient)
		expectedStatus int
		expectedBody   interface{}
		expectedCookie *http.Cookie
	}{
		{
			name: "Success",
			body: `{"login": "user1", "password": "pass123"}`,
			mockBehavior: func(m *mocks.MockCalcServiceClient) {
				m.EXPECT().Login(gomock.Any(), &proto.User{Login: "user1", Password: "pass123"}).Return(&proto.Id{Id: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"id": float64(1)},
			expectedCookie: &http.Cookie{
				Name:     "id",
				Value:    "1",
				Path:     "/",
				HttpOnly: true,
			},
		},
		{
			name:           "InvalidJSON",
			body:           `{"login": "user1", "password": "pass123"`, // Неверный JSON
			mockBehavior:   func(m *mocks.MockCalcServiceClient) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "unexpected EOF"},
			expectedCookie: nil,
		},
		{
			name: "AgentError",
			body: `{"login": "user1", "password": "pass123"}`,
			mockBehavior: func(m *mocks.MockCalcServiceClient) {
				m.EXPECT().Login(gomock.Any(), &proto.User{Login: "user1", Password: "pass123"}).Return(nil, status.Error(codes.NotFound, "user not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "rpc error: code = NotFound desc = user not found"},
			expectedCookie: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAgent := mocks.NewMockCalcServiceClient(ctrl)
			test.mockBehavior(mockAgent)

			app := &Application{agent: mockAgent}
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(test.body))
			rr := httptest.NewRecorder()

			app.loginHandler(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

			var actualBody map[string]interface{}
			err := json.NewDecoder(rr.Body).Decode(&actualBody)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedBody, actualBody)

			if test.expectedCookie != nil {
				cookies := rr.Result().Cookies()
				assert.Len(t, cookies, 1)
				cookie := cookies[0]
				assert.Equal(t, test.expectedCookie.Name, cookie.Name)
				assert.Equal(t, test.expectedCookie.Value, cookie.Value)
				assert.Equal(t, test.expectedCookie.Path, cookie.Path)
				assert.Equal(t, test.expectedCookie.HttpOnly, cookie.HttpOnly)
				assert.WithinDuration(t, time.Now().Add(24*time.Hour), cookie.Expires, time.Minute)
			} else {
				assert.Empty(t, rr.Result().Cookies())
			}
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockBehavior   func(m *mocks.MockCalcServiceClient)
		expectedStatus int
		expectedBody   interface{}
		expectedCookie *http.Cookie
	}{
		{
			name: "Success",
			body: `{"login": "user1", "password": "pass123"}`,
			mockBehavior: func(m *mocks.MockCalcServiceClient) {
				m.EXPECT().Register(gomock.Any(), &proto.User{Login: "user1", Password: "pass123"}).Return(&proto.Id{Id: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"id": float64(1)},
			expectedCookie: &http.Cookie{
				Name:     "id",
				Value:    "1",
				Path:     "/",
				HttpOnly: true,
			},
		},
		{
			name:           "InvalidJSON",
			body:           `{"login": "user1", "password": "pass123"`, // Неверный JSON
			mockBehavior:   func(m *mocks.MockCalcServiceClient) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "unexpected EOF"},
			expectedCookie: nil,
		},
		{
			name: "AgentError",
			body: `{"login": "user1", "password": "pass123"}`,
			mockBehavior: func(m *mocks.MockCalcServiceClient) {
				m.EXPECT().Register(gomock.Any(), &proto.User{Login: "user1", Password: "pass123"}).Return(nil, status.Error(codes.AlreadyExists, "user already exists"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "rpc error: code = AlreadyExists desc = user already exists"},
			expectedCookie: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAgent := mocks.NewMockCalcServiceClient(ctrl)
			test.mockBehavior(mockAgent)

			app := &Application{agent: mockAgent}
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(test.body))
			rr := httptest.NewRecorder()

			app.registerHandler(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

			var actualBody map[string]interface{}
			err := json.NewDecoder(rr.Body).Decode(&actualBody)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedBody, actualBody)

			if test.expectedCookie != nil {
				cookies := rr.Result().Cookies()
				assert.Len(t, cookies, 1)
				cookie := cookies[0]
				assert.Equal(t, test.expectedCookie.Name, cookie.Name)
				assert.Equal(t, test.expectedCookie.Value, cookie.Value)
				assert.Equal(t, test.expectedCookie.Path, cookie.Path)
				assert.Equal(t, test.expectedCookie.HttpOnly, cookie.HttpOnly)
				assert.WithinDuration(t, time.Now().Add(24*time.Hour), cookie.Expires, time.Minute)
			} else {
				assert.Empty(t, rr.Result().Cookies())
			}
		})
	}
}
