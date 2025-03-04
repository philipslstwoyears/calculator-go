package calc

import (
	"fmt"
	"github.com/philipslstwoyears/calculator-go/internal/dto"
	"github.com/philipslstwoyears/calculator-go/internal/storage"
	"os"
	"strconv"
)

type Worker struct {
	storage *storage.Storage
	input   chan dto.Expression
}

func New(storage *storage.Storage, input chan dto.Expression) *Worker {
	return &Worker{
		storage: storage,
		input:   input,
	}
}
func (w *Worker) Start() error {
	countWorkers := os.Getenv("COMPUTING_POWER")
	computerWorkers, err := strconv.Atoi(countWorkers)
	if err != nil {
		return err
	}
	for i := 0; i < computerWorkers; i++ {
		go w.worker()
	}
	return nil
}

func (w *Worker) worker() {
	for expression := range w.input {
		calc, err := Calc(expression.Expression)
		if err != nil {
			expression.Status = fmt.Sprintf("Ошибка: %v", err)
		} else {
			expression.Status = "Ok"
		}
		expression.Result = calc
		w.storage.Update(expression)
	}
}
