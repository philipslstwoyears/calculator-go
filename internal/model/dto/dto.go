package dto

type ErrorResponse struct {
	Error string `json:"error"`
}

type Expression struct {
	UserID     int     `json:"user_id"`
	ID         int     `json:"id"`
	Status     string  `json:"status"`
	Result     float64 `json:"result"`
	Expression string  `json:"expression"`
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Id       int    `json:"id"`
}

type Request struct {
	Expression string `json:"expression"`
}
