package rof

type Statement interface{}

type VisitorStatement interface {
	visitExprStmt(expr Expression) interface{}
	visitPrintStmt(expr Print) interface{}
	visitVarStmt(expr Var) interface{}
}

type Expression struct {
	Expr Expr
}

func (e Expression) Accept(visitor VisitorStatement) interface{} {
	return visitor.visitExprStmt(e)
}

type Print struct {
	Expr Expr
}

func (p Print) Accept(visitor VisitorStatement) interface{} {
	return visitor.visitPrintStmt(p)
}

type Var struct {
	Name        Token
	Initializer Expr
}

func (p Var) Accept(visitor VisitorStatement) interface{} {
	return visitor.visitVarStmt(p)
}
