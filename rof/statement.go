package rof

type Stmt interface {
	Statement() Stmt
}

type VisitorStatement interface {
	visitExprStmt(expr Expression) interface{}
	visitPrintStmt(expr Print) interface{}
	visitVarStmt(expr Var) interface{}
}

type Expression struct {
	Expr
}

func (e Expression) Accept(visitor VisitorStatement) interface{} {
	return visitor.visitExprStmt(e)
}

func (e Expression) Statement() Stmt {
	return e
}

type Print struct {
	Expr
}

func (p Print) Accept(visitor VisitorStatement) interface{} {
	return visitor.visitPrintStmt(p)
}

func (p Print) Statement() Stmt {
	return p
}

type Var struct {
	Name        Token
	Initializer Expr
}

func (p Var) Accept(visitor VisitorStatement) interface{} {
	return visitor.visitVarStmt(p)
}

func (v Var) Statement() Stmt {
	return v
}
