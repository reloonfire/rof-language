package rof

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
	Env *Environment
}

func NewInterpreter() Interpreter {
	var i Interpreter
	i.Env = NewEnv(nil)
	return i
}

func (i Interpreter) Interpret(stmts []Stmt) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	for _, stmt := range stmts {
		i.execute(stmt)
	}

	return nil
}

func (i Interpreter) evaluate(expr Expr) interface{} {
	switch t := expr.(type) {
	case Binary:
		return i.BinaryExpr(t)
	case Grouping:
		return i.GroupingExpr(t)
	case Literal:
		return i.LiteralExpr(t)
	case Unary:
		return i.UnaryExpr(t)
	case Variable:
		return i.VariableExpr(t)
	case Assign:
		return i.AssignExpr(t)
	case Logical:
		return i.LogicalExpr(t)
	default:
		fmt.Println("[ERROR] Type -> ", reflect.TypeOf(t))
		return nil
	}
}

func (i Interpreter) execute(stmt Stmt) {
	switch t := stmt.(type) {
	case Expression:
		i.ExprStmt(t)
	case Print:
		i.PrintStmt(t)
	case Var:
		i.VarStmt(t)
	case Block:
		i.BlockStmt(t)
	case If:
		i.IfStmt(t)
	case While:
		i.WhileStmt(t)
	default:
		fmt.Println("[ERROR] Type -> ", reflect.TypeOf(t))
	}
}

func (i Interpreter) BinaryExpr(expr Binary) interface{} {
	right := i.evaluate(expr.Right)
	left := i.evaluate(expr.Left)
	//fmt.Println("[DEBUG] BinaryExpr Called -> ", expr, "\n\n	RIGHT -> ", right, "\n	LEFT -> ", left, "\n	Operator -> ", expr.Operator)

	switch expr.Operator.TokenType {
	case GREATER:
		l, r := i.checkNumberOperands(expr.Operator, left, right)
		return l > r
	case GREATER_EQUAL:
		l, r := i.checkNumberOperands(expr.Operator, left, right)
		return l >= r
	case LESS:
		l, r := i.checkNumberOperands(expr.Operator, left, right)
		return l < r
	case LESS_EQUAL:
		l, r := i.checkNumberOperands(expr.Operator, left, right)
		return l <= r
	case MINUS:
		l, r := i.checkNumberOperands(expr.Operator, left, right)
		return l - r
	case SLASH:
		l, r := i.checkNumberOperands(expr.Operator, left, right)
		return l / r
	case STAR:
		l, r := i.checkNumberOperands(expr.Operator, left, right)
		return l * r
	case PLUS:
		switch l := left.(type) {
		case float64:
			return l + right.(float64)
		case string:
			return l + fmt.Sprint(right)
		}

	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	}

	// Non raggiungibile
	return nil

}

func (i Interpreter) GroupingExpr(expr Grouping) interface{} {
	return i.evaluate(expr.Expr)
}

func (i Interpreter) LiteralExpr(expr Literal) interface{} {
	return expr.Value
}

func (i Interpreter) UnaryExpr(expr Unary) interface{} {
	right := i.evaluate(expr)

	switch expr.Operator.TokenType {
	case BANG:
		return !i.isTruthy(right)
	case MINUS:
		r := i.checkNumberOperand(expr.Operator, right)
		return -r
	}

	// Non raggiungibile
	return nil

}

func (i Interpreter) VariableExpr(expr Variable) interface{} {
	return i.Env.Get(expr.Name)
}

func (i Interpreter) AssignExpr(expr Assign) interface{} {
	value := i.evaluate(expr.Value)

	i.Env.Assign(expr.Name, value)
	return value
}

func (i Interpreter) LogicalExpr(expr Logical) interface{} {
	left := i.evaluate(expr.Left)

	if expr.Operator.TokenType == OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}
	return i.evaluate(expr.Right)
}

func (i Interpreter) ExprStmt(stmt Expression) {
	i.evaluate(stmt.Expr)
}

func (i Interpreter) PrintStmt(stmt Print) {
	//fmt.Println("[DEBUG] Print Called ->", stmt)
	value := i.evaluate(stmt.Expr)
	fmt.Println(Stringify(value))
}

func (i Interpreter) VarStmt(stmt Var) {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	//fmt.Println("[DEBUG] Create var ", stmt.Name.Lexeme, " -> ", value)
	i.Env.Define(stmt.Name.Lexeme, value)
}

func (i Interpreter) BlockStmt(stmt Block) {
	previous := i.Env

	i.Env = NewEnv(previous)
	defer func() { i.Env = previous }()
	for _, s := range stmt.Statements {
		i.execute(s)
	}
}

func (i Interpreter) IfStmt(stmt If) {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
}

func (i Interpreter) WhileStmt(stmt While) {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
}

// Helper

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

func (i Interpreter) checkNumberOperand(operator Token, operand interface{}) float64 {
	var n float64
	ok := false

	switch t := operand.(type) {
	case float64:
		ok = true
		n = t
	}

	if !ok {
		panic(&RuntimeError{operator, "Operand must be number"})
	}

	return n
}

func (i Interpreter) checkNumberOperands(operator Token, left, right interface{}) (float64, float64) {
	ok1, ok2 := false, false
	var n1, n2 float64
	switch t := left.(type) {
	case float64:
		ok1 = true
		n1 = t
	}
	switch t := right.(type) {
	case float64:
		ok2 = true
		n2 = t
	}
	if !ok1 || !ok2 {
		panic(&RuntimeError{operator, "Operands must be numbers"})
	}

	return n1, n2
}

func Stringify(obj interface{}) string {
	if obj == nil {
		return "nil"
	}

	//fmt.Println("[DEBUG] Print -> ", obj)

	switch t := obj.(type) {
	case string:
		return t
	case float64:
		return fmt.Sprintf("%v", t)
	case bool:
		return fmt.Sprintf("%v", t)
	default:
		return "nil"
	}

}
