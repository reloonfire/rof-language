package rof

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
	Env Environment
}

func NewInterpreter() Interpreter {
	var i Interpreter
	i.Env = NewEnv()
	return i
}

func (i Interpreter) Interpret(stmts []interface{}) {
	for _, stmt := range stmts {
		i.execute(stmt)
	}
}

func (i Interpreter) execute(stmt interface{}) {
	AcceptMethodCall(stmt, []reflect.Value{reflect.ValueOf(i)})
}

func (i Interpreter) visitBinaryExpr(expr Binary) interface{} {
	right := i.evaluate(expr.Right)
	left := i.evaluate(expr.Left)

	switch expr.Operator.TokenType {
	case GREATER:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case MINUS:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case SLASH:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case STAR:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	case PLUS:
		if reflect.TypeOf(left).String() == "float64" && reflect.TypeOf(right).String() == "float64" {
			return left.(float64) + right.(float64)
		} else if reflect.TypeOf(left).String() == "string" && reflect.TypeOf(right).String() == "string" {
			return left.(string) + right.(string)
		}
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	}

	// Non raggiungibile
	return nil

}

func (i Interpreter) visitGroupingExpr(expr Grouping) interface{} {
	return i.evaluate(expr.Expression)
}

func (i Interpreter) visitLiteralExpr(expr Literal) interface{} {
	return expr.Value
}

func (i Interpreter) visitUnaryExpr(expr Unary) interface{} {
	right := i.evaluate(expr)

	switch expr.Operator.TokenType {
	case BANG:
		return !i.isTruthy(right)
	case MINUS:
		i.checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	}

	// Non raggiungibile
	return nil

}

func (i Interpreter) visitVariableExpr(expr Variable) interface{} {
	return i.Env.Get(expr.Name)
}

func (i Interpreter) visitExprStmt(stmt Expression) interface{} {
	return i.evaluate(stmt.Expr)
}

func (i Interpreter) visitPrintStmt(stmt Print) interface{} {
	//fmt.Println("[DEBUG] Print Called ->", stmt)
	value := i.evaluate(stmt.Expr)
	fmt.Println(i.stringify(value))
	return nil
}

func (i Interpreter) visitVarStmt(stmt Var) interface{} {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.Env.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i Interpreter) evaluate(expr interface{}) interface{} {
	return AcceptMethodCall(expr, []reflect.Value{reflect.ValueOf(i)}).Interface()
}

func (i Interpreter) isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}
	if reflect.TypeOf(obj).String() == "bool" {
		return obj.(bool)
	}
	return true
}

func (i Interpreter) isEqual(obj1, obj2 interface{}) bool {
	return obj1 == obj2
}

func (i Interpreter) checkNumberOperand(operator Token, operand interface{}) {
	if reflect.ValueOf(operand).String() == "float64" {
		return
	}
	fmt.Println("[Line "+string(operator.Line)+"] RuntimeError: Operand must be a number, not", reflect.TypeOf(operand), ".")
}

func (i Interpreter) checkNumberOperands(operator Token, left, right interface{}) {
	if reflect.TypeOf(left).String() == "float64" && reflect.TypeOf(right).String() == "float64" {
		return
	}
	fmt.Println("[Line "+string(operator.Line)+"] RuntimeError: Operands must be a numbers, not ", reflect.TypeOf(left), " and ", reflect.TypeOf(right), ".")
}

func (i Interpreter) stringify(obj interface{}) string {
	if obj == nil {
		return "nil"
	}
	if reflect.TypeOf(obj).String() == "float64" {
		return fmt.Sprintf("%5.2f", obj.(float64))
	}
	if reflect.TypeOf(obj).String() == "bool" {
		return fmt.Sprint(obj.(bool))
	}

	return obj.(string)
}
