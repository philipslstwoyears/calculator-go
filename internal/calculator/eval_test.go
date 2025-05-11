package calc

import (
	"github.com/joho/godotenv"
	"log"
	"testing"
)

func TestCalc(t *testing.T) {
	err := godotenv.Load("../../cmd/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "/",
			expression:     "1/2",
			expectedResult: 0.5,
		},
		{
			name:           "priority max",
			expression:     "-(-11-(1*20/2)-11/2*3)",
			expectedResult: 37.5,
		},
		{
			name:           "priority 50/50",
			expression:     "-11-1*20/2",
			expectedResult: -21,
		},
		{
			name:           "priority not too much",
			expression:     "-11-1*20/001",
			expectedResult: -31,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := Calc(testCase.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error: %v", testCase.expression, err)
			}
			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}
}
