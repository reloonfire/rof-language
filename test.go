package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/reloonfire/rof-language/rof"
)

func main() {
	sc := new(rof.Scanner)
	//printer := new(rof.ASTPrinter)
	parser := new(rof.Parser)
	interpreter := new(rof.Interpreter)

	s, _ := ioutil.ReadFile("test.rof")
	sc.Source = string(s)
	// Scanner
	tokens := sc.Scan()
	if sc.HadError {
		return
	}
	// Parser
	parser.Tokens = tokens
	expr := parser.Parse()
	b, _ := json.Marshal(expr)
	fmt.Println("[DEBUG] Parsed Tokens -> ", string(b))
	if parser.HadError {
		return
	}
	// Interpreter
	interpreter.Interpret(expr)
}
