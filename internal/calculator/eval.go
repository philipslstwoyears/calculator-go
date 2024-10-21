package calculator

import "fmt"

func Calc(expression string) (float64, error) {
	if len(expression) == 0 {
		return 0, fmt.Errorf("expression is empty")
	}
	return 0, nil
}
func convertToPolishNotation(expression string) (float64, error) {
	//tokens := make([]string, 0)
	//for _, ch := range expression {
	//	switch ch {
	//	case '+', '-', '*', '/', '(':
	//	}
	//}
	return 0, nil
}
