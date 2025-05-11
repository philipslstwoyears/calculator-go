package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/philipslstwoyears/calculator-go/internal/model/convert"
	"github.com/philipslstwoyears/calculator-go/internal/model/dto"
	"github.com/philipslstwoyears/calculator-go/proto"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (a *Application) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	userId, err := strconv.Atoi(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	request := new(dto.Request)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	id, err := a.agent.Calc(r.Context(), &proto.Request{Expression: request.Expression, UserId: int32(userId)})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int32{"id": id.Id})
}
func (a *Application) expressionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("id")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	userId, err := strconv.Atoi(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	expressions, err := a.agent.GetExpressions(r.Context(), &proto.Id{Id: int32(userId)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	result := make([]*dto.Expression, len(expressions.GetExpressions()))
	for i, expression := range expressions.GetExpressions() {
		result[i] = convert.ExpressionToDTO(expression)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (a *Application) expressionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie("id")
	if err != nil {
		log.Printf("expressionHandler: no cookie, err=%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	userId, err := strconv.Atoi(cookie.Value)
	if err != nil {
		log.Printf("expressionHandler: invalid cookie, err=%v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("expressionHandler: invalid ID, err=%v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	expression, err := a.agent.GetExpression(r.Context(), &proto.Id{Id: int32(id)})
	if err != nil {
		log.Printf("expressionHandler: agent error, err=%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	if expression.GetUserId() != int32(userId) {
		log.Printf("expressionHandler: forbidden, userId=%d, expressionUserId=%d", userId, expression.GetUserId())
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: "It is not your expression"})
		return
	}

	log.Printf("expressionHandler: success, id=%d", id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(convert.ExpressionToDTO(expression))
}

func (a *Application) loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("loginHandler: start")
	w.Header().Set("Content-Type", "application/json")

	request := new(dto.User)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("loginHandler: invalid JSON, err=%v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	id, err := a.agent.Login(r.Context(), &proto.User{Login: request.Login, Password: request.Password})
	if err != nil {
		log.Printf("loginHandler: agent error, err=%v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "id",
		Value:    strconv.Itoa(int(id.Id)),
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	})
	log.Printf("loginHandler: success, id=%d", id.Id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int32{"id": id.Id})
}

func (a *Application) registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("registerHandler: start")
	w.Header().Set("Content-Type", "application/json") // Устанавливаем Content-Type в начале

	request := new(dto.User)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("registerHandler: invalid JSON, err=%v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	id, err := a.agent.Register(r.Context(), &proto.User{Login: request.Login, Password: request.Password})
	if err != nil {
		log.Printf("registerHandler: agent error, err=%v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Error: err.Error()})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "id",
		Value:    strconv.Itoa(int(id.Id)),
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	})
	log.Printf("registerHandler: success, id=%d", id.Id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int32{"id": id.Id})
}
