package main

import (
	"bufio"
	"fmt"
	calc "github.com/philipslstwoyears/calculator-go/internal/calculator"
	"log"
	"os"
)

func main() {
	expression, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	result, err := calc.Calc(expression)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
