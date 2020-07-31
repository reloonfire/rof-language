package rof

import "fmt"

type LoxCallable func(Interpreter, []Expr) interface{}

type Callable interface {
	Arity() int
	Call(i Interpreter, args []Expr) interface{}
}

type NativeFunction struct {
	Callable
	NativeCall LoxCallable
	A          int
}

// Call is the operation that executes a builtin function
func (n NativeFunction) Call(i Interpreter, arguments []Expr) interface{} {
	return n.NativeCall(i, arguments)
}

// Arity returns the number of allowed parameters for the native function
func (n NativeFunction) Arity() int {
	return n.A
}

// String returns the name of the native function
func (n NativeFunction) String() string {
	return fmt.Sprintf("<native/%p>", n.NativeCall)
}
