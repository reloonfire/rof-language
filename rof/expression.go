package rof

type VisitorExpression interface {
	visitBinaryExpr(expr Binary) interface{}
	visitGroupingExpr(expr Grouping) interface{}
	visitLiteralExpr(expr Literal) interface{}
	visitUnaryExpr(expr Unary) interface{}
	visitVariableExpr(expr Variable) interface{}
}

type Expr interface {
	Expression() Expr
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b Binary) Accept(visitor VisitorExpression) interface{} {
	return visitor.visitBinaryExpr(b)
}

func (b Binary) Expression() Expr {
	return b
}

type Grouping struct {
	Expr
}

func (b Grouping) Expression() Expr {
	return b
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

func (l Literal) Expression() Expr {
	return l
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u Unary) Accept(visitor VisitorExpression) interface{} {
	return visitor.visitUnaryExpr(u)
}

func (u Unary) Expression() Expr {
	return u
}

type Variable struct {
	Name Token
}

func (v Variable) Accept(visitor VisitorExpression) interface{} {
	return visitor.visitVariableExpr(v)
}

func (v Variable) Expression() Expr {
	return v
}
