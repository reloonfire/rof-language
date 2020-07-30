package rof

type Stmt interface {
	Statement() Stmt
}

type Expression struct {
	Expr
}

func (e Expression) Statement() Stmt {
	return e
}

type Print struct {
	Expr
}

func (p Print) Statement() Stmt {
	return p
}

type Var struct {
	Name        Token
	Initializer Expr
}

func (v Var) Statement() Stmt {
	return v
}

type Block struct {
	Statements []Stmt
}

func (v Block) Statement() Stmt {
	return v
}
