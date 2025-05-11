package convert

import (
	"github.com/philipslstwoyears/calculator-go/internal/model/dto"
	"github.com/philipslstwoyears/calculator-go/proto"
)

func ExpressionToDTO(e *proto.Expression) *dto.Expression {
	return &dto.Expression{
		ID:         int(e.Id),
		Expression: e.Expression,
		UserID:     int(e.UserId),
		Status:     e.Status,
		Result:     float64(e.Result),
	}
}
