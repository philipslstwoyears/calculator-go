package agent

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/philipslstwoyears/calculator-go/internal/dto"
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
		Status:     "Запрос находиться в обработке",
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
