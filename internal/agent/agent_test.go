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

func TestAgent_GetExpressions(t *testing.T) {
	tests := []struct {
		name          string
		input         *proto.Id
		expected      *proto.Expressions
		mockBehavior  func(r *mocks.MockStorage)
		expectedError error
	}{
		{
			name:  "Success",
			input: &proto.Id{Id: 1},
			expected: &proto.Expressions{
				Expressions: []*proto.Expression{
					{
						Id:         1,
						Expression: "5+5",
						Status:     "completed",
						Result:     10.0,
						UserId:     1,
					},
				},
			},
			mockBehavior: func(r *mocks.MockStorage) {
				r.EXPECT().GetExpressions(1).Return([]dto.Expression{
					{
						ID:         1,
						Expression: "5+5",
						Status:     "completed",
						Result:     10.0,
						UserID:     1,
					},
				}, nil)
			},
			expectedError: nil,
		},
		{
			name:     "Fail",
			input:    &proto.Id{Id: 1},
			expected: nil,
			mockBehavior: func(r *mocks.MockStorage) {
				r.EXPECT().GetExpressions(1).Return(nil, errors.New("database error"))
			},
			expectedError: status.Error(codes.Internal, "database error"),
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
			result, err := agent.GetExpressions(context.Background(), test.input)
			assert.Equal(t, test.expected, result)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestAgent_GetExpression(t *testing.T) {
	tests := []struct {
		name          string
		input         *proto.Id
		expected      *proto.Expression
		mockBehavior  func(r *mocks.MockStorage)
		expectedError error
	}{
		{
			name:  "Success",
			input: &proto.Id{Id: 1},
			expected: &proto.Expression{
				Id:         1,
				Expression: "5+5",
				Status:     "completed",
				Result:     10.0,
				UserId:     1,
			},
			mockBehavior: func(r *mocks.MockStorage) {
				r.EXPECT().GetExpression(1).Return(dto.Expression{
					ID:         1,
					Expression: "5+5",
					Status:     "completed",
					Result:     10.0,
					UserID:     1,
				}, true)
			},
			expectedError: nil,
		},
		{
			name:     "NotFound",
			input:    &proto.Id{Id: 999},
			expected: nil,
			mockBehavior: func(r *mocks.MockStorage) {
				r.EXPECT().GetExpression(999).Return(dto.Expression{}, false)
			},
			expectedError: status.Error(codes.NotFound, "expression not found"),
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
			result, err := agent.GetExpression(context.Background(), test.input)
			assert.Equal(t, test.expected, result)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestAgent_Login(t *testing.T) {
	tests := []struct {
		name          string
		input         *proto.User
		expected      *proto.Id
		mockBehavior  func(r *mocks.MockStorage)
		expectedError error
	}{
		{
			name:  "Success",
			input: &proto.User{Login: "user1", Password: "pass123"},
			expected: &proto.Id{
				Id: 1,
			},
			mockBehavior: func(r *mocks.MockStorage) {
				r.EXPECT().GetUser("user1").Return(dto.User{
					Id:       1,
					Login:    "user1",
					Password: "pass123",
				}, true)
			},
			expectedError: nil,
		},
		{
			name:     "UserNotFound",
			input:    &proto.User{Login: "user1", Password: "pass123"},
			expected: nil,
			mockBehavior: func(r *mocks.MockStorage) {
				r.EXPECT().GetUser("user1").Return(dto.User{}, false)
			},
			expectedError: status.Error(codes.NotFound, "user not found"),
		},
		{
			name:     "WrongPassword",
			input:    &proto.User{Login: "user1", Password: "wrongpass"},
			expected: nil,
			mockBehavior: func(r *mocks.MockStorage) {
				r.EXPECT().GetUser("user1").Return(dto.User{
					Id:       1,
					Login:    "user1",
					Password: "pass123",
				}, true)
			},
			expectedError: status.Error(codes.InvalidArgument, "wrong password"),
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
			result, err := agent.Login(context.Background(), test.input)
			assert.Equal(t, test.expected, result)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestAgent_Register(t *testing.T) {
	tests := []struct {
		name          string
		input         *proto.User
		expected      *proto.Id
		mockBehavior  func(r *mocks.MockStorage)
		expectedError error
	}{
		{
			name:  "Success",
			input: &proto.User{Login: "user1", Password: "pass123"},
			expected: &proto.Id{
				Id: 1,
			},
			mockBehavior: func(r *mocks.MockStorage) {
				r.EXPECT().GetUser("user1").Return(dto.User{}, false)
				r.EXPECT().AddUser(dto.User{
					Login:    "user1",
					Password: "pass123",
				}).Return(1, nil)
			},
			expectedError: nil,
		},
		{
			name:     "UserAlreadyExists",
			input:    &proto.User{Login: "user1", Password: "pass123"},
			expected: nil,
			mockBehavior: func(r *mocks.MockStorage) {
				r.EXPECT().GetUser("user1").Return(dto.User{
					Id:       1,
					Login:    "user1",
					Password: "pass123",
				}, true)
			},
			expectedError: status.Error(codes.AlreadyExists, "user is already registered"),
		},
		{
			name:     "StorageError",
			input:    &proto.User{Login: "user1", Password: "pass123"},
			expected: nil,
			mockBehavior: func(r *mocks.MockStorage) {
				r.EXPECT().GetUser("user1").Return(dto.User{}, false)
				r.EXPECT().AddUser(dto.User{
					Login:    "user1",
					Password: "pass123",
				}).Return(0, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
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
			result, err := agent.Register(context.Background(), test.input)
			assert.Equal(t, test.expected, result)
			assert.Equal(t, test.expectedError, err)
		})
	}
}
