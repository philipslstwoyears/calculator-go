package calc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func operate(num1 float64, num2 float64, operation string) float64 {
	switch operation {
	case "+":
		return num1 + num2
	case "-":
		return num1 - num2
	case "*":
		return num1 * num2
	case "/":
		return num1 / num2
	}
	return -1
}

func calcPolishNotation(polishNotation []string) (float64, error) {
	numStack := New[float64]()
	for _, elem := range polishNotation {
		switch elem {
		case "+", "-", "*", "/":
			if numStack.Size() < 2 {
				return -1, fmt.Errorf("Wrong expression")
			}
			num2 := numStack.Pop()
			num1 := numStack.Pop()
			result := operate(num1, num2, elem)
			numStack.Push(result)
		case "~":
			if numStack.IsEmpty() {
				return -1, fmt.Errorf("Wrong expression")
			}
			numStack.Push(-numStack.Pop())
		default:
			num, err := strconv.ParseFloat(elem, 64)
			if err != nil {
				return -1, err
			}
			numStack.Push(num)
		}
	}
	return numStack.Pop(), nil
}

func convertToPolishNotation(expression string) ([]string, error) {
	polishNotation := []string{}
	operationStack := New[string]()
	num := ""
	countBrackets := 0

	for i, character := range expression {
		switch character {
		case '+', '-', '*', '/', '(':
			if num != "" {
				polishNotation = append(polishNotation, num)
				num = ""
			}
			if character == '-' && (i == 0 || (expression[i-1] != ')' && !unicode.IsDigit(rune(expression[i-1])))) {
				operationStack.Push("~")
				continue
			}
			if character == '(' {
				countBrackets++
			}
			for !operationStack.IsEmpty() {
				operationFromStack := operationStack.Peek()
				if (character == '+' || character == '-') && (operationFromStack != "(" && operationFromStack != ")") {
					polishNotation = append(polishNotation, operationStack.Pop())
				} else if (character == '*' || character == '/') && (operationFromStack == "*" || operationFromStack == "/") {
					polishNotation = append(polishNotation, operationStack.Pop())
				} else {
					break
				}
			}
			operationStack.Push(string(character))

		case ')':
			if num != "" {
				polishNotation = append(polishNotation, num)
				num = ""
			}
			countBrackets--
			if countBrackets < 0 {
				return []string{}, fmt.Errorf("Wrong expression: check brackets")
			}
			for operationFromStack := operationStack.Pop(); operationFromStack != "("; operationFromStack = operationStack.Pop() {
				polishNotation = append(polishNotation, operationFromStack)
			}

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			num += string(character)

		default:
			return []string{}, fmt.Errorf("Wrong expression: symbols should be nums or '*/+-().'")
		}

	}

	if num != "" {
		polishNotation = append(polishNotation, num)
	}

	for !operationStack.IsEmpty() {
		operationFromStack := operationStack.Pop()
		if operationFromStack == "(" {
			return []string{}, fmt.Errorf("Wrong expression: check brackets")
		}
		polishNotation = append(polishNotation, operationFromStack)
	}

	return polishNotation, nil
}

func Calc(expression string) (float64, error) {
	if expression == "" {
		return 0, fmt.Errorf("Expression is empty")
	}
	expression = strings.ReplaceAll(expression, " ", "")
	expression = strings.ReplaceAll(expression, "\n", "")
	polishNotation, err := convertToPolishNotation(expression)
	if err != nil {
		return -1, err
	}
	result, err := calcPolishNotation(polishNotation)
	if err != nil {
		return -1, err
	}
	return result, nil
}

// stack

type element interface {
	int | float64 | string
}

type Stack[T element] struct {
	items    []T
	topIndex int
}

func New[T element]() *Stack[T] {
	return &Stack[T]{
		items:    make([]T, 0),
		topIndex: -1,
	}
}

func (s *Stack[T]) IsEmpty() bool {
	return s.topIndex == -1
}

func (s *Stack[T]) Push(item T) {
	s.topIndex++
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
	popedElement := s.items[s.topIndex]
	s.topIndex--
	s.items = s.items[:s.topIndex+1]
	return popedElement
}

func (s *Stack[T]) Peek() T {
	return s.items[s.topIndex]
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}
