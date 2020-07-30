package rof

import (
	"fmt"

	"github.com/reloonfire/rof-language/helpers"
)

type Environment struct {
	Values map[string]interface{}
}

func NewEnv() Environment {
	var e Environment
	e.Values = make(map[string]interface{}, 15)
	return e
}

func (e *Environment) Get(name Token) interface{} {
	if helpers.ContainsKey(e.Values, name.Lexeme) {
		return e.Values[name.Lexeme]
	}

	fmt.Println("Undefined Variable '", name.Lexeme, "'.")
	return nil
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
	//fmt.Println("[DEBUG] Env -> ", e.Values)
}
