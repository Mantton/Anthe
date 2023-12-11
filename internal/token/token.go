package token

type TokenType byte

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = iota
	EOF               // EOF

	// identifiers + literals
	IDENTIFIER // add, foo,bar, x, y

	INTEGER // 1233
	FLOAT
	STRING

	// operator
	ASSIGN // =
	NOT    // !

	// arithmetic
	ADD // +
	SUB // -
	MUL // *
	QUO // /

	// boolean
	LSS // <
	GTR // >
	EQL // ==
	NEQ // !=

	// Delims
	COMMA     // ,
	SEMICOLON // ;
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]

	// Keywords
	FUNCTION
	RETURN

	LET
	CONST

	IF
	ELSE

	TRUE
	FALSE
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"return": RETURN,

	"let":   LET,
	"const": CONST,

	"if":   IF,
	"else": ELSE,

	"true":  TRUE,
	"false": FALSE,
}

var symbols = map[rune]TokenType{
	'=': ASSIGN,
	'!': NOT,

	// arithmetic
	'+': ADD,
	'-': SUB,
	'*': MUL,
	'/': QUO,

	// boolean
	'>': GTR,
	'<': LSS,

	// delim
	',': COMMA,
	';': SEMICOLON,
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
	'[': LBRACKET,
	']': RBRACKET,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}

func IsSymbol(ch rune) bool {
	_, ok := symbols[ch]
	return ok
}

func LookUpSymbol(ch rune) TokenType {
	return symbols[ch]
}
