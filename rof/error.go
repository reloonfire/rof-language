package rof

import "fmt"

type RuntimeError struct {
	Token   Token
	Message string
}

func (re *RuntimeError) Error() string {
	return fmt.Sprintf("line #%d at '%v': '%s'", re.Token.Line, re.Token.Lexeme, re.Message)
}

type ParseError RuntimeError

func (pe *ParseError) Error() string {
	if pe.Token.TokenType == EOF {
		return fmt.Sprintf("line #%d at end: %s", pe.Token.Line, pe.Message)
	} else {
		return fmt.Sprintf("line #%d at '%v': %s", pe.Token.Line, pe.Token.Lexeme, pe.Message)
	}
}
