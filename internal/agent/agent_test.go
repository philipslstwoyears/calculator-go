package agent

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/philipslstwoyears/calculator-go/internal/mocks"
	"github.com/philipslstwoyears/calculator-go/internal/model/dto"
	"github.com/philipslstwoyears/calculator-go/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestAgent_Calc(t *testing.T) {
	tests := []struct {
		name          string
		input         *proto.Request
		expected      *proto.Id
		mockBehavior  func(r *mocks.MockStorage)
		expectedError error
	}{
		{
			name:     "Success",
			input:    &proto.Request{UserId: 1, Expression: "5+5"},
			expected: &proto.Id{Id: 1},
			mockBehavior: func(r *mocks.MockStorage) {
				r.EXPECT().AddExpression(dto.Expression{
					UserID:     1,
					Expression: "5+5",
					Status:     "Выражение принято для вычисления",
				}).Return(1, nil)
			},
			expectedError: nil,
		},
		{
			name:     "Fail",
			input:    &proto.Request{},
			expected: nil,
			mockBehavior: func(r *mocks.MockStorage) {
				r.EXPECT().AddExpression(dto.Expression{
					Status: "Выражение принято для вычисления",
				}).Return(0, errors.New("error"))
			},
			expectedError: status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			storage := mocks.NewMockStorage(c)
			test.mockBehavior(storage)
			ch := make(chan dto.Expression, 1)
			agent := New(storage, ch)
			id, err := agent.Calc(context.Background(), test.input)
			assert.Equal(t, test.expected, id)
			assert.Equal(t, test.expectedError, err)
			if err == nil {
				expression := <-ch
				assert.Equal(t, test.input.Expression, expression.Expression)
				assert.Equal(t, int(test.input.UserId), expression.UserID)
				assert.Equal(t, int(test.expected.Id), expression.ID)
			}
		})
	}
}
