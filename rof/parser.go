package rof

import "fmt"

type Parser struct {
	Tokens     []Token
	Statements []interface{}
	Current    int
	HadError   bool
}

func (p *Parser) Parse() []interface{} {
	for !p.isAtEnd() {
		p.Statements = append(p.Statements, p.statement())
	}

	return p.Statements
}

func (p *Parser) expression() interface{} {
	return p.equality()
}

func (p *Parser) statement() interface{} {
	if p.match(PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() interface{} {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ; after value.")
	return Print{Expr: value}
}

func (p *Parser) expressionStatement() interface{} {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ; after expression.")
	return Expression{Expr: value}
}

func (p *Parser) equality() interface{} {
	//fmt.Println("[DEBUG] Equality ->", p.peek())
	expr := p.comparison()
	for p.match(BANG, BANG_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

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

func (p *Parser) comparison() interface{} {
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

func (p *Parser) addition() interface{} {
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

func (p *Parser) multiplication() interface{} {
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

func (p *Parser) unary() interface{} {
	//fmt.Println("[DEBUG] Unary ->", p.peek())
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		//fmt.Println("[DEBUG] IS Unary")
		return Unary{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() interface{} {
	//fmt.Println("[DEBUG] Primary ->", p.peek())
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

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{Expression: expr}
	}

	p.error(p.peek(), "Expect expression.")
	return nil
}

func (p *Parser) consume(t TokenType, text string) Token {
	if p.check(t) {
		return p.advance()
	}

	p.error(p.peek(), text)

	return Token{}
}

func (p *Parser) error(token Token, text string) {
	if token.TokenType == EOF {
		p.report(token.Line, " at end", text)
	} else {
		p.report(token.Line, " at '"+token.Lexeme+"'", text)
	}
}

func (p *Parser) report(line int, where, message string) {
	fmt.Println("[ line ", line, "] Error", where, ": ", message)
	p.HadError = true
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
