package dto

type ErrorResponse struct {
	Error string `json:"error"`
}

type Expression struct {
	Id         int     `json:"id"`
	Status     string  `json:"status"`
	Result     float64 `json:"result"`
	Expression string
}

type Request struct {
	Expression string `json:"expression"`
}
