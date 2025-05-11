package storage

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/philipslstwoyears/calculator-go/internal/dto"
	"github.com/philipslstwoyears/calculator-go/internal/middleware"
	"github.com/philipslstwoyears/calculator-go/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGoodRequest(t *testing.T) {
	// Создаем gomock контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для Storage
	mockStorage := mocks.NewMockStorage(ctrl)

	// Тестовые случаи
	testCases := []struct {
		body       string
		expression string
		result     float64
	}{
		{`{"expression":"1+1"}`, "1+1", 2.00},
		{`{"expression":"11*(2+2)"}`, "11*(2+2)", 44.00},
		{`{"expression":"102/3"}`, "102/3", 34.00},
		{`{"expression":"-(-11-(1*20/2)-11/2*3)"}`, "-(-11-(1*20/2)-11/2*3)", 37.50},
		{`{"expression":"1*001*2/4"}`, "1*001*2/4", 0.50},
		{`{"expression":"1-12-1*(-1)"}`, "1-12-1*(-1)", -10.00},
	}

	for _, tc := range testCases {
		t.Run(tc.body, func(t *testing.T) {
			// Настраиваем мок для AddExpression
			mockStorage.EXPECT().
				AddExpression(dto.Expression{
					Expression: tc.expression,
					UserID:     123,
					Result:     tc.result,
				}).
				Return(1, nil)

			// Создаем HTTP-запрос
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tc.body)))
			request.Header.Set("Content-Type", "application/json")
			request.AddCookie(&http.Cookie{Name: "id", Value: "123"})

			// Создаем Application с мок-хранилищем
			a := &Application{storage: mockStorage}
			r := httptest.NewRecorder()
			handler := http.HandlerFunc(a.CalculateHandler)
			middleware.LoggerMiddleware(middleware.RecoverMiddleware(handler)).ServeHTTP(r, request)

			// Проверяем статус ответа
			if r.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, r.Code)
			}

			// Проверяем тело ответа
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Could not read body: %v", err)
			}
			// Если CalculateHandler возвращает {"result":"X.XX"}
			expected := `{"result":"` + fmt.Sprintf("%.2f", tc.result) + `"}` + "\n"
			if string(body) != expected {
				t.Errorf("Expected body '%s', got '%s'", expected, string(body))
			}
			// Если CalculateHandler возвращает {"id":X}, замените на:
			// expected := `{"id":1}` + "\n"
		})
	}
}
