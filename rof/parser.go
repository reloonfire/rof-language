package rof

import "fmt"

type Parser struct {
	Tokens   []Token
	Current  int
	HadError bool
}

func (p *Parser) expression() interface{} {
	return p.equality()
}

func (p *Parser) equality() interface{} {
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
	expr := p.addition()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.addition()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) addition() interface{} {
	expr := p.multiplication()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.multiplication()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) multiplication() interface{} {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() interface{} {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return Unary{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() interface{} {
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
	return nil
}

func (p *Parser) consume(t TokenType, text string) Token {
	if p.check(t) {
		p.advance()
	}

	p.error(p.peek(), text)

	return Token{}
}

func (p *Parser) error(token Token, text string) {
	if token.TokenType == EOF {
		p.report(token.Line, " at end", text)
	} else {
		p.report(token.Line, " ad '"+token.Lexeme+"'", text)
	}
}

func (p *Parser) report(line int, where, message string) {
	fmt.Println("[line ", line, "] Error", where, ": ", message)
	p.HadError = true
}
