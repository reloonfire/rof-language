package rof

type VisitorExpression interface {
	visitBinaryExpr(expr Binary) interface{}
	visitGroupingExpr(expr Grouping) interface{}
	visitLiteralExpr(expr Literal) interface{}
	visitUnaryExpr(expr Unary) interface{}
	visitVariableExpr(expr Variable) interface{}
}

type Expr interface{}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b Binary) Accept(visitor VisitorExpression) interface{} {
	return visitor.visitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) Accept(visitor VisitorExpression) interface{} {
	return visitor.visitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func (l Literal) Accept(visitor VisitorExpression) interface{} {
	return visitor.visitLiteralExpr(l)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u Unary) Accept(visitor VisitorExpression) interface{} {
	return visitor.visitUnaryExpr(u)
}

type Variable struct {
	Name Token
}

func (v Variable) Accept(visitor VisitorExpression) interface{} {
	return visitor.visitVariableExpr(v)
}
