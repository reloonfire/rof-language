package rof

import (
	"fmt"
	"reflect"
)

type ASTPrinter struct {
}

func (a *ASTPrinter) Print(expr interface{}) string {
	fmt.Println("[DEBUG] Reflect Value Of ->", reflect.ValueOf(expr).Kind())
	return reflect.ValueOf(expr).MethodByName("Accept").Call([]reflect.Value{reflect.ValueOf(a)})[0].String()
}

func (a *ASTPrinter) visitBinaryExpr(expr Binary) string {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *ASTPrinter) visitGroupingExpr(expr Grouping) string {
	return a.parenthesize("group", expr.Expression)
}

func (a *ASTPrinter) visitLiteralExpr(expr Literal) string {
	if expr.Value == nil {
		return "nil"
	}
	return expr.Value.(string)
}

func (a *ASTPrinter) visitUnaryExpr(expr Unary) string {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *ASTPrinter) parenthesize(name string, exprs ...interface{}) string {
	var text string

	text += "(" + name
	for _, expr := range exprs {
		text += " " + reflect.ValueOf(expr).MethodByName("Accept").Call([]reflect.Value{reflect.ValueOf(a)})[0].String()
	}
	text += ")"
	return text
}
