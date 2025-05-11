package agent

import (
	"context"
	"github.com/philipslstwoyears/calculator-go/internal/model/dto"
	"github.com/philipslstwoyears/calculator-go/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Application) Calc(ctx context.Context, r *proto.Request) (*proto.Id, error) {
	expression := dto.Expression{
		Expression: r.Expression,
		UserID:     int(r.UserId),
		Status:     "Выражение принято для вычисления",
	}
	id, err := a.Storage.AddExpression(expression)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	expression.ID = id
	a.input <- expression
	return &proto.Id{
		Id: int32(id),
	}, nil
}

func (a *Application) GetExpressions(ctx context.Context, id *proto.Id) (*proto.Expressions, error) {
	exp, err := a.Storage.GetExpressions(int(id.GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := make([]*proto.Expression, len(exp))
	for i, expression := range exp {
		result[i] = &proto.Expression{
			Expression: expression.Expression,
			Status:     expression.Status,
			Id:         int32(expression.ID),
			Result:     float32(expression.Result),
			UserId:     int32(expression.UserID),
		}
	}
	return &proto.Expressions{
		Expressions: result,
	}, nil
}

func (a *Application) GetExpression(ctx context.Context, id *proto.Id) (*proto.Expression, error) {
	i, ok := a.Storage.GetExpression(int(id.GetId()))
	if !ok {
		return nil, status.Error(codes.NotFound, "expression not found")
	}
	expression := &proto.Expression{
		Expression: i.Expression,
		Status:     i.Status,
		Result:     float32(i.Result),
		Id:         int32(i.ID),
		UserId:     int32(i.UserID),
	}
	return expression, nil
}
func (a *Application) Login(ctx context.Context, in *proto.User) (*proto.Id, error) {
	user, ok := a.Storage.GetUser(in.Login)
	if !ok {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if user.Password != in.Password {
		return nil, status.Error(codes.InvalidArgument, "wrong password")
	}
	return &proto.Id{
		Id: int32(user.Id),
	}, nil
}
func (a *Application) Register(ctx context.Context, in *proto.User) (*proto.Id, error) {
	_, ok := a.Storage.GetUser(in.Login)
	if ok {
		return nil, status.Error(codes.AlreadyExists, "user is already registered")
	}
	id, err := a.Storage.AddUser(dto.User{
		Login:    in.Login,
		Password: in.Password,
	})
	if err != nil {
		return nil, err
	}
	return &proto.Id{
		Id: int32(id),
	}, nil
}
