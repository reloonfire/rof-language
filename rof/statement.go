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

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (v If) Statement() Stmt {
	return v
}

type While struct {
	Condition Expr
	Body      Stmt
}

func (w While) Statement() Stmt {
	return w
}
