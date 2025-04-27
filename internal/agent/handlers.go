package agent

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/philipslstwoyears/calculator-go/internal/dto"
	"github.com/philipslstwoyears/calculator-go/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

func (a *Application) CalcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	request := new(dto.Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		data := &dto.ErrorResponse{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(data)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expression := dto.Expression{
		Expression: request.Expression,
		Status:     "Выражение принято для вычисления",
	}
	id := a.storage.Add(expression)
	expression.Id = id
	a.input <- expression
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}
func (a *Application) ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	exp := a.storage.GetAll()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exp)
}
func (a *Application) ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	atoi, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	i, ok := a.storage.Get(atoi)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(i)
}

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
