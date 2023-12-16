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
	STRING // ab

	NULL
	VOID

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

	LEQ // <=
	GEQ // >=

	// Delims
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]

	Q_MARK // ?

	// Keywords
	FUNCTION
	RETURN

	LET
	CONST

	IF
	ELSE

	TRUE
	FALSE

	// Primitive Typing
	INT_T  // integer type
	STR_T  // string type
	FLT_T  // float type
	BOOL_T // bool type

	ARR_T // array type
	SET_T // set

	MAP_T // map type
	OBJ_T // object

	OPTIONAL_T // optional
	RSLT_T     // result type

	ANY_OBJ_T // any object
	ANY_T     // any non null value

	STRUCT // struct declaration

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

	"null": NULL,
	"void": VOID,

	"struct": STRUCT,
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
	':': COLON,

	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
	'[': LBRACKET,
	']': RBRACKET,
	'?': Q_MARK,
}

var builtin_types = map[string]TokenType{
	"int":    INT_T,
	"string": STR_T,
	"float":  FLT_T,
	"bool":   BOOL_T,

	"array": ARR_T,
	"set":   SET_T,

	"map":    MAP_T,
	"object": OBJ_T,

	"optional": OPTIONAL_T,
	"result":   RSLT_T,

	"any_object": ANY_OBJ_T,
	"any":        ANY_T,

	"null": NULL,
	"void": VOID,
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

func LookUpBuiltInType(ident string) TokenType {
	if tok, ok := builtin_types[ident]; ok {
		return tok
	}

	return IDENTIFIER
}
