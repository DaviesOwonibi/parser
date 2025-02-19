package lexer

import (
	"bufio"
	"io"
	"unicode"
)

type Token int

const (
	EOF = iota
	ILLEGAL
	IDENT
	INT
	SEMI // ;
	COLON

	OPEN_BRACKET
	CLOSE_BRACKET
	COMMA
	DOT
	OPEN_BRACE
	CLOSE_BRACE
	FORWARDSLASH
	DOUBLEQUOTE
	SINGLEQUOTE

	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /
	MOD // %

	ASSIGN // =
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	INT:     "INT",
	SEMI:    ";",
	COLON:   ":",

	OPEN_BRACKET:  "OPEN_BRACKET",
	CLOSE_BRACKET: "CLOSE_BRACKET",
	COMMA:         "COMMA",
	DOT:           "DOT",
	OPEN_BRACE:    "OPEN BRACE",
	CLOSE_BRACE:   "CLOSE BRACE",
	FORWARDSLASH:  "FORWARDSLASH",
	DOUBLEQUOTE:   "DOUBLEQUOTE",
	SINGLEQUOTE:   "SINGLEQUOTE",

	// Infix ops
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
	MOD: "%",

	ASSIGN: "=",
}

func (t Token) String() string {
	return tokens[t]
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) Lex() (Position, Token, string) {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}

			// at this point there isn't much we can do, and the compiler
			// should just return the raw error to the user
			panic(err)
		}

		l.pos.column++

		switch r {
		case '\n':
			l.resetPosition()
		case ';':
			return l.pos, SEMI, ";"
		case ':':
			return l.pos, COLON, ":"
		case '+':
			return l.pos, ADD, "+"
		case '-':
			return l.pos, SUB, "-"
		case '*':
			return l.pos, MUL, "*"
		case '/':
			return l.pos, DIV, "/"
		case '%':
			return l.pos, MOD, "%"
		case '=':
			return l.pos, ASSIGN, "="
		case '(':
			return l.pos, OPEN_BRACKET, "("
		case ')':
			return l.pos, CLOSE_BRACKET, ")"
		case ',':
			return l.pos, COMMA, ","
		case '.':
			return l.pos, DOT, "."
		case '{':
			return l.pos, OPEN_BRACE, "{"
		case '}':
			return l.pos, CLOSE_BRACE, "}"

		case '\\':
			return l.pos, FORWARDSLASH, "\\"
		case '"':
			return l.pos, DOUBLEQUOTE, `"`
		case '\'':
			return l.pos, SINGLEQUOTE, "'"
		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if unicode.IsDigit(r) {
				// backup and let lexInt rescan the beginning of the int
				startPos := l.pos
				l.backup()
				lit := l.lexInt()
				return startPos, INT, lit
			} else if unicode.IsLetter(r) {
				// backup and let lexIdent rescan the beginning of the ident
				startPos := l.pos
				l.backup()
				lit := l.lexIdent()
				return startPos, IDENT, lit
			} else {
				return l.pos, ILLEGAL, string(r)
			}
		}
	}
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the int
				return lit
			}
		}

		l.pos.column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the integer
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit
			}
		}

		l.pos.column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the identifier
			l.backup()
			return lit
		}
	}
}
