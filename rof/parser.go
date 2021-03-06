package rof

import (
	"fmt"
)

type Parser struct {
	Tokens     []Token
	Statements []Stmt
	Current    int
	HadError   bool
}

func (p *Parser) Parse() []Stmt {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Parse Error:", r.(error))
			p.HadError = true
			p.synchronize()
		}
	}()

	for !p.isAtEnd() {
		p.Statements = append(p.Statements, p.declaration())
	}

	return p.Statements
}

func (p *Parser) declaration() Stmt {
	if p.match(VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() Var {
	tokenName := p.consume(IDENTIFIER, "Expect variable name.")
	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return Var{Name: tokenName, Initializer: initializer}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()
	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()
		exprVar, ok := expr.(Variable)
		if ok {
			return Assign{exprVar.Name, value}
		}

		panic(&ParseError{equals, "Invalid assignment target."})
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		return Logical{expr, operator, right}
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		return Logical{expr, operator, right}
	}

	return expr
}

func (p *Parser) statement() Stmt {
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LEFT_BRACE) {
		return Block{p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer Stmt
	switch {
	case p.match(SEMICOLON):
		initializer = nil
	case p.match(VAR):
		initializer = p.varDeclaration()
	default:
		initializer = p.expressionStatement()
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	body := p.statement()

	if increment != nil {
		body = Block{[]Stmt{body, Expression{increment}}}
	}

	if condition == nil {
		condition = Literal{true}
	}

	body = While{condition, body}

	if initializer != nil {
		body = Block{[]Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' before while condition.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after while condition")
	body := p.statement()

	return While{condition, body}
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' before if condition.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition")

	then := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return If{condition, then, elseBranch}
}

func (p *Parser) block() []Stmt {
	s := []Stmt{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		s = append(s, p.declaration())
	}

	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return s
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ; after value.")
	return Print{value}
}

func (p *Parser) expressionStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ; after expression.")
	return Expression{value}
}

func (p *Parser) equality() Expr {
	//fmt.Println("[DEBUG] Equality ->", p.peek())
	expr := p.comparison()
	for p.match(BANG, BANG_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	//fmt.Println("[DEBUG] Comparison ->", p.peek())
	expr := p.addition()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.addition()
		//fmt.Println("[DEBUG] IS Comparison")
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) addition() Expr {
	//fmt.Println("[DEBUG] Addition ->", p.peek())
	expr := p.multiplication()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.multiplication()
		//fmt.Println("[DEBUG] IS Addition")
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) multiplication() Expr {
	//fmt.Println("[DEBUG] Multiplication ->", p.peek())
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		//fmt.Println("[DEBUG] IS Multiplication")
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	//fmt.Println("[DEBUG] Unary ->", p.peek())
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		//fmt.Println("[DEBUG] IS Unary")
		return Unary{Operator: operator, Right: right}
	}

	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(expr Expr) Expr {
	var args []Expr

	if !p.check(RIGHT_PAREN) {
		args = append(args, p.expression())
		for p.match(COMMA) {
			if len(args) > 255 {
				panic(ParseError{p.peek(), "Cannot have more than 255 arguments."})
			}
			args = append(args, p.expression())
		}
	}
	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")

	return Call{expr, paren, args}

}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return Literal{Value: false}
	}
	if p.match(TRUE) {
		return Literal{Value: true}
	}
	if p.match(NIL) {
		return Literal{Value: nil}
	}

	if p.match(NUMBER, STRING) {
		return Literal{p.previous().Literal}
	}

	if p.match(IDENTIFIER) {
		return Variable{p.previous()}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{expr}
	}

	panic(&ParseError{p.peek(), "Expect expression"})
}

// Helper

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == t
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.Current++
	}
	//fmt.Println("[DEBUG] Current cursor -> ", p.Current)
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) consume(t TokenType, text string) Token {
	if p.check(t) {
		return p.advance()
	}

	panic(&ParseError{p.peek(), text})
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		p.advance()
	}
}
