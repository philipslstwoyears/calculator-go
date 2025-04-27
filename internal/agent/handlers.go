package agent

import (
	"context"
	"github.com/philipslstwoyears/calculator-go/internal/dto"
	"github.com/philipslstwoyears/calculator-go/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Application) Calc(ctx context.Context, r *proto.Request) (*proto.Id, error) {
	expression := dto.Expression{
		Expression: r.Expression,
		Status:     "Выражение принято для вычисления",
	}
	id := a.storage.Add(expression)
	expression.Id = id
	a.input <- expression
	return &proto.Id{
		Id: int32(id),
	}, nil
}

func (a *Application) GetExpressions(context.Context, *proto.Empty) (*proto.Expressions, error) {
	exp := a.storage.GetAll()
	result := make([]*proto.Expression, len(exp))
	for i, expression := range exp {
		result[i] = &proto.Expression{
			Expression: expression.Expression,
			Status:     expression.Status,
			Id:         int32(expression.Id),
			Result:     float32(expression.Result),
		}
	}
	return &proto.Expressions{
		Expressions: result,
	}, nil
}

func (a *Application) GetExpression(ctx context.Context, id *proto.Id) (*proto.Expression, error) {
	i, ok := a.storage.Get(int(id.GetId()))
	if !ok {
		return nil, status.Error(codes.NotFound, "expression not found")
	}
	expression := &proto.Expression{
		Expression: i.Expression,
		Status:     i.Status,
		Result:     float32(i.Result),
		Id:         int32(i.Id),
	}
	return expression, nil
}
