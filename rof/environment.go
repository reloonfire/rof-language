package rof

import (
	"github.com/reloonfire/rof-language/helpers"
)

type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

func NewEnv(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
		Values:    make(map[string]interface{}),
	}
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
		e.Enclosing.Assign(name, value)
		return
	}

	panic(&RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."})
}
