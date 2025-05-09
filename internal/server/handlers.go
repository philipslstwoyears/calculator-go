package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/philipslstwoyears/calculator-go/internal/dto"
	"github.com/philipslstwoyears/calculator-go/proto"
	"net/http"
	"strconv"
	"time"
)

func (a *Application) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("id")
	if err != nil {
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId, err := strconv.Atoi(cookie.Value)
	if err != nil {
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	request := new(dto.Request)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := a.agent.Calc(r.Context(), &proto.Request{Expression: request.Expression, UserId: int32(userId)})
	if err != nil {
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int32{"id": id.Id})
}

func (a *Application) expressionsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("id")
	if err != nil {
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userId, err := strconv.Atoi(cookie.Value)
	if err != nil {
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	expressions, err := a.agent.GetExpressions(r.Context(), &proto.Id{Id: int32(userId)})
	if err != nil {
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expressions.GetExpressions())
}

func (a *Application) expressionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expression, err := a.agent.GetExpression(r.Context(), &proto.Id{Id: int32(id)})
	if err != nil {
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expression)

}

func (a *Application) loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	request := new(dto.User)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := a.agent.Login(r.Context(), &proto.User{Login: request.Login, Password: request.Password})
	if err != nil {
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "id", Value: strconv.Itoa(int(id.Id)), Expires: time.Now().Add(24 * time.Hour), Path: "/", HttpOnly: true})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int32{"id": id.Id})
}

func (a *Application) registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	request := new(dto.User)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := a.agent.Register(r.Context(), &proto.User{Login: request.Login, Password: request.Password})
	if err != nil {
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "id", Value: strconv.Itoa(int(id.Id)), Expires: time.Now().Add(24 * time.Hour), Path: "/", HttpOnly: true})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int32{"id": id.Id})
}
