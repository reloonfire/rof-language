package rof

/*
import (
	"fmt"
	"reflect"
)

type ASTPrinter struct {
}

func AcceptMethodCall(expr interface{}, input []reflect.Value) reflect.Value {
	var ptr reflect.Value
	value := reflect.ValueOf(expr)
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem() // acquire value referenced by pointer
	} else {
		ptr = reflect.New(reflect.TypeOf(expr)) // create new pointer
		temp := ptr.Elem()                      // create variable to value of pointer
		temp.Set(value)                         // set value of variable to our passed in value
	}

	var finalMethod reflect.Value
	method := value.MethodByName("Accept")
	if method.IsValid() {
		finalMethod = method
	}
	// check for method on pointer
	method = ptr.MethodByName("Accept")
	if method.IsValid() {
		finalMethod = method
	}

	if finalMethod.IsValid() {
		return finalMethod.Call(input)[0]
	}
	return reflect.Value{}
}

func (a *ASTPrinter) Print(expr interface{}) string {
	return AcceptMethodCall(expr, []reflect.Value{reflect.ValueOf(a)}).Interface().(string)
}

func (a *ASTPrinter) visitBinaryExpr(expr Binary) interface{} {

	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *ASTPrinter) visitGroupingExpr(expr Grouping) interface{} {
	return a.parenthesize("group", expr.Expression)
}

func (a *ASTPrinter) visitLiteralExpr(expr Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return expr.Value.(interface{})
}

func (a *ASTPrinter) visitUnaryExpr(expr Unary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *ASTPrinter) parenthesize(name string, exprs ...interface{}) interface{} {
	var text string

	text += "(" + name
	for _, expr := range exprs {
		value := AcceptMethodCall(expr, []reflect.Value{reflect.ValueOf(a)}).Interface()
		text += " " + ConvertToString(value)
	}
	text += ")"

	return text
}

func ConvertToString(value interface{}) string {
	switch value.(type) {
	case string:
		return value.(string)
	case float64:
		return fmt.Sprintf("%5.2f", value.(float64))
	default:
		return ""
	}
}
*/
