package rof

type Expr interface {
	Expression() Expr
}

type Assign struct {
	Name  Token
	Value Expr
}

func (b Assign) Expression() Expr {
	return b
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
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

type Literal struct {
	Value interface{}
}

func (l Literal) Expression() Expr {
	return l
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u Unary) Expression() Expr {
	return u
}

type Variable struct {
	Name Token
}

func (v Variable) Expression() Expr {
	return v
}

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (l Logical) Expression() Expr {
	return l
}

type Call struct {
	Callee Expr
	Paren  Token
	Args   []Expr
}

func (c Call) Expression() Expr {
	return c
}
