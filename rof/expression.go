package rof

type Visitor interface {
	visitBinaryExpr(expr Binary) interface{}
	visitGroupingExpr(expr Grouping) interface{}
	visitLiteralExpr(expr Literal) interface{}
	visitUnaryExpr(expr Unary) interface{}
}

type Expr interface{}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) Accept(visitor Visitor) interface{} {
	return visitor.visitBinaryExpr(*b)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor Visitor) interface{} {
	return visitor.visitGroupingExpr(*g)
}

type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(visitor Visitor) interface{} {
	return visitor.visitLiteralExpr(*l)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(visitor Visitor) interface{} {
	return visitor.visitUnaryExpr(*u)
}
