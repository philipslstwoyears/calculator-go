package calc

import (
	"errors"
	"fmt"
	"github.com/philipslstwoyears/calculator-go/internal/stack"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func operate(num1 float64, num2 float64, operation string) (float64, error) {
	switch operation {
	case "+":
		additionMs := os.Getenv("TIME_ADDITION_MS")
		ms, err := strconv.Atoi(additionMs)
		if err != nil {
			return 0, err
		}
		time.Sleep(time.Duration(ms))
		return num1 + num2, nil
	case "-":
		subMs := os.Getenv("TIME_SUBTRACTION_MS")
		ms, err := strconv.Atoi(subMs)
		if err != nil {
			return 0, err
		}
		time.Sleep(time.Duration(ms))
		return num1 - num2, nil
	case "*":
		multMs := os.Getenv("TIME_MULTIPLICATIONS_MS")
		ms, err := strconv.Atoi(multMs)
		if err != nil {
			return 0, err
		}
		time.Sleep(time.Duration(ms))
		return num1 * num2, nil
	case "/":
		divMs := os.Getenv("TIME_DIVISIONS_MS")
		ms, err := strconv.Atoi(divMs)
		if err != nil {
			return 0, err
		}
		time.Sleep(time.Duration(ms))
		if num2 == 0 {
			return 0, errors.New("Division by zero")
		}
		return num1 / num2, nil
	}
	return -1, nil
}

func calcPolishNotation(polishNotation []string) (float64, error) {
	numStack := stack.New[float64]()
	for _, elem := range polishNotation {
		switch elem {
		case "+", "-", "*", "/":
			if numStack.Size() < 2 {
				return -1, fmt.Errorf("Wrong expression")
			}
			num2 := numStack.Pop()
			num1 := numStack.Pop()
			result, err := operate(num1, num2, elem)
			if err != nil {
				return -1, err
			}
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
	operationStack := stack.New[string]()
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
			return []string{}, fmt.Errorf("inccorect simbol: %s", string(character))
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
