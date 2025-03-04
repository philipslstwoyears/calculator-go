package storage

import (
	"testing"

	"github.com/philipslstwoyears/calculator-go/internal/dto"
)

func TestStorage_Add(t *testing.T) {
	s := New()
	expr := dto.Expression{Expression: "2+2", Result: 4.0}
	id := s.Add(expr)

	s.mu.RLock()
	defer s.mu.RUnlock()
	storedExpr, exists := s.data[id]
	if !exists {
		t.Fatalf("Expression with ID %d not found", id)
	}
	if storedExpr.Expression != "2+2" || storedExpr.Result != 4.0 {
		t.Errorf("Stored expression mismatch: got %+v, want %+v", storedExpr, expr)
	}
}

func TestStorage_Update(t *testing.T) {
	s := New()
	expr := dto.Expression{Expression: "3+3", Result: 6.0}
	id := s.Add(expr)

	updatedExpr := dto.Expression{Id: id, Expression: "3*3", Result: 9.0}
	s.Update(updatedExpr)

	s.mu.RLock()
	defer s.mu.RUnlock()
	storedExpr, exists := s.data[id]
	if !exists {
		t.Fatalf("Updated expression with ID %d not found", id)
	}
	if storedExpr.Expression != "3*3" || storedExpr.Result != 9.0 {
		t.Errorf("Updated expression mismatch: got %+v, want %+v", storedExpr, updatedExpr)
	}
}

func TestStorage_Get(t *testing.T) {
	s := New()
	expr := dto.Expression{Expression: "5+5", Result: 10.0}
	id := s.Add(expr)

	retrievedExpr, exists := s.Get(id)
	if !exists {
		t.Fatalf("Expected expression with ID %d, but it was not found", id)
	}
	if retrievedExpr.Expression != "5+5" || retrievedExpr.Result != 10.0 {
		t.Errorf("Retrieved expression mismatch: got %+v, want %+v", retrievedExpr, expr)
	}
}

func TestStorage_GetAll(t *testing.T) {
	s := New()
	expr1 := dto.Expression{Expression: "1+1", Result: 2.0}
	expr2 := dto.Expression{Expression: "2+2", Result: 4.0}

	s.Add(expr1)
	s.Add(expr2)

	allExpressions := s.GetAll()
	if len(allExpressions) != 2 {
		t.Fatalf("Expected 2 expressions, got %d", len(allExpressions))
	}

	expected := []dto.Expression{expr1, expr2}
	for i, expr := range allExpressions {
		if expr.Expression != expected[i].Expression || expr.Result != expected[i].Result {
			t.Errorf("Mismatch at index %d: got %+v, want %+v", i, expr, expected[i])
		}
	}
}
