package rof

import (
	"github.com/reloonfire/rof-language/helpers"
)

type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

func NewEnv() Environment {
	var e Environment
	e.Enclosing = nil
	e.Values = make(map[string]interface{}, 15)
	return e
}

func (e *Environment) Get(name Token) interface{} {
	if helpers.ContainsKey(e.Values, name.Lexeme) {
		return e.Values[name.Lexeme]
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	panic(&RuntimeError{name, "Undefined Variable '" + name.Lexeme + "'."})
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
	//fmt.Println("[DEBUG] Env -> ", e.Values)
}

func (e *Environment) Assign(name Token, value interface{}) {
	if helpers.ContainsKey(e.Values, name.Lexeme) {
		e.Values[name.Lexeme] = value
		return
	}

	if e.Enclosing != nil {
		e.Enclosing.ASsign(name, value)
		return
	}

	panic(&RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."})
}
