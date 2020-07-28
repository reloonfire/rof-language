package rof

import (
	"github.com/reloonfire/rof-language/helpers"
)

// Scanner - Scanner look into the source looking for tokens
type Scanner struct {
	Source   string
	Tokens   []Token
	Start    int
	Current  int
	Line     int
	HadError bool
}

// Scan - Scan through source looking for tokens
func (s *Scanner) Scan() []Token {
	for !s.IsEnd() {
		s.Start = s.Current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, Token{TokenType: EOF, Lexeme: "EOF", Literal: nil, Line: s.Line})
	return s.Tokens
}

func (s *Scanner) advance() string {
	s.Current++
	return string(s.Source[s.Current-1])
}

// AddToken - Add Token
func (s *Scanner) addToken(t TokenType, literal interface{}) {
	lexeme := s.Source[s.Start:s.Current]
	s.Tokens = append(s.Tokens, Token{TokenType: t, Literal: literal, Line: s.Line, Lexeme: lexeme})
}

// ScanToken - Scan Token
func (s *Scanner) scanToken() {
	c := s.advance()
	//fmt.Println("[DEBUG] Analysing char '", c, "' at Pos [", s.Line, "]")
	switch c {
	case "(":
		s.addToken(LEFT_PAREN, nil)
		break
	case ")":
		s.addToken(RIGHT_PAREN, nil)
		break
	case "{":
		s.addToken(LEFT_BRACE, nil)
		break
	case "}":
		s.addToken(RIGHT_BRACE, nil)
		break
	case ",":
		s.addToken(COMMA, nil)
		break
	case ".":
		s.addToken(DOT, nil)
		break
	case "-":
		s.addToken(MINUS, nil)
		break
	case "+":
		s.addToken(PLUS, nil)
		break
	case ";":
		s.addToken(SEMICOLON, nil)
		break
	case "*":
		s.addToken(STAR, nil)
		break
	case "!":
		if s.match("=") {
			s.addToken(BANG_EQUAL, nil)
		} else {
			s.addToken(BANG, nil)
		}
		break
	case "=":
		if s.match("=") {
			s.addToken(EQUAL_EQUAL, nil)
		} else {
			s.addToken(EQUAL, nil)
		}
		break
	case "<":
		if s.match("=") {
			s.addToken(LESS_EQUAL, nil)
		} else {
			s.addToken(LESS, nil)
		}
		break
	case ">":
		if s.match("=") {
			s.addToken(GREATER_EQUAL, nil)
		} else {
			s.addToken(GREATER, nil)
		}
		break
	case "/":
		if s.match("/") {
			// A comment goes until the end of the line.
			for s.peek() != "\n" && !s.IsEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, nil)
		}
		break
	case " ":
	case "\r":
	case "\t":
		// Ignore whitespace.
		break
	case "\n":
		s.Line++
		break
	case "\"":
		s.string()
		break
	default:
		helpers.ReportError(s.Line, "Unexpected character.")
	}
}

func (s *Scanner) string() {
	for s.peek() != "\"" && !s.IsEnd() {
		if s.peek() == "\n" {
			s.Line++
		}
		s.advance()
	}

	// Unterminated string.
	if s.IsEnd() {
		helpers.ReportError(s.Line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := string(s.Source[s.Start+1 : s.Current-1])
	s.addToken(STRING, value)
}

func (s *Scanner) peek() string {
	if s.IsEnd() {
		return "\000"
	}
	return string(s.Source[s.Current])
}

func (s *Scanner) match(what string) bool {
	if s.IsEnd() {
		return false
	}
	if string(s.Source[s.Current]) != what {
		return false
	}
	s.Current++
	return true
}

// IsEnd - Check if current cursor is at the ending of source
func (s *Scanner) IsEnd() bool {
	return s.Current >= len(s.Source)
}
