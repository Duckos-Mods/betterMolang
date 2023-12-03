package bettermolang

import "fmt"

var (

	// We use the first byte of the token as the key
	TOKEN_VLT_MAP = map[byte]VLToken{
		'!': {
			TokenType: TOKEN_BANG,
			NextToken: []VLTokenNode{
				{
					SucessType: TOKEN_BANG_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
			},
			consumeLength: 1,
		},
		'=': {
			TokenType: TOKEN_EQUAL,
			NextToken: []VLTokenNode{
				{
					SucessType: TOKEN_EQUAL_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
			},
			consumeLength: 1,
		},
		'>': {
			TokenType: TOKEN_GREATER,
			NextToken: []VLTokenNode{
				{
					SucessType: TOKEN_GREATER_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
			},
			consumeLength: 1,
		},
		'<': {
			TokenType: TOKEN_LESS,
			NextToken: []VLTokenNode{
				{
					SucessType: TOKEN_LESS_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
			},
			consumeLength: 1,
		},
		'+': {
			TokenType: TOKEN_PLUS,
			NextToken: []VLTokenNode{
				{
					SucessType: TOKEN_PLUS_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
				{
					SucessType: TOKEN_PLUS_PLUS,
					SelfVal:    '+',
					NextVal:    nil,
				},
			},
			consumeLength: 1,
		},
		'-': {
			TokenType: TOKEN_MINUS,
			NextToken: []VLTokenNode{
				{
					SucessType: TOKEN_MINUS_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
				{
					SucessType: TOKEN_MINUS_MINUS,
					SelfVal:    '-',
					NextVal:    nil,
				},
				{
					SucessType: TOKEN_ENTITY_SELECTOR,
					SelfVal:    '>',
					NextVal:    nil,
				},
			},
			consumeLength: 1,
		},
		'*': {
			TokenType: TOKEN_STAR,
			NextToken: []VLTokenNode{
				{
					SucessType: TOKEN_STAR_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
				{
					SucessType: TOKEN_MULTI_LINE_COMMENT,
					SelfVal:    '/',
					NextVal:    nil,
				},
			},
			consumeLength: 1,
		},
		'/': {
			TokenType: TOKEN_SLASH,
			NextToken: []VLTokenNode{
				{
					SucessType: TOKEN_SLASH_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
				{
					SucessType: TOKEN_MULTI_LINE_COMMENT,
					SelfVal:    '*',
					NextVal:    nil,
				},
				{
					SucessType: TOKEN_COMMENT,
					SelfVal:    '/',
					NextVal:    nil,
				},
			},
			consumeLength: 1,
		},
		'&': {
			TokenType: TOKEN_NULL,
			NextToken: []VLTokenNode{
				{
					SucessType: TOKEN_AND_AND,
					SelfVal:    '&',
					NextVal:    nil,
				},
			},
			consumeLength: 1,
		},
		'|': {
			TokenType: TOKEN_NULL,
			NextToken: []VLTokenNode{
				{
					SucessType: TOKEN_OR_OR,
					SelfVal:    '|',
					NextVal:    nil,
				},
			},
			consumeLength: 1,
		},
		'?': {
			TokenType: TOKEN_QUESTION,
			NextToken: []VLTokenNode{
				{
					SucessType: TOKEN_NULL_COALESCE,
					SelfVal:    '?',
					NextVal:    nil,
				},
			},
			consumeLength: 1,
		},
		';': {
			TokenType:     TOKEN_SEMICOLON,
			NextToken:     nil,
			consumeLength: 1,
		},
		':': {
			TokenType:     TOKEN_COLON,
			NextToken:     nil,
			consumeLength: 1,
		},
		'\'': {
			TokenType:     TOKEN_SINGLE_QUOTE,
			NextToken:     nil,
			consumeLength: 1,
		},
		',': {
			TokenType:     TOKEN_COMMA,
			NextToken:     nil,
			consumeLength: 1,
		},
		'(': {
			TokenType:     TOKEN_LEFT_PAREN,
			NextToken:     nil,
			consumeLength: 1,
		},
		')': {
			TokenType:     TOKEN_RIGHT_PAREN,
			NextToken:     nil,
			consumeLength: 1,
		},
		'{': {
			TokenType:     TOKEN_LEFT_BRACE,
			NextToken:     nil,
			consumeLength: 1,
		},
		'}': {
			TokenType:     TOKEN_RIGHT_BRACE,
			NextToken:     nil,
			consumeLength: 1,
		},
		'[': {
			TokenType:     TOKEN_LEFT_BRACKET,
			NextToken:     nil,
			consumeLength: 1,
		},
		']': {
			TokenType:     TOKEN_RIGHT_BRACKET,
			NextToken:     nil,
			consumeLength: 1,
		},
		'.': {
			TokenType:     TOKEN_DOT,
			NextToken:     nil,
			consumeLength: 1,
		},
	}

	TOKEN_SPECIAL_MAP = map[int]func(*Scanner){
		TOKEN_SINGLE_QUOTE: func(scan *Scanner) {
			strVal := make([]byte, 0)
			strVal = append(strVal, scan.peakN(-1)) // Add the single quote by reading the previous char which is unsafe but i dont care
			// because i am slow im going to build a token then push it back. Which is a waste of memory but i dont care
			var escapeChar byte = 0x0
			for scan.peak() != '\'' && !scan.isAtEnd() || escapeChar == '\\' {
				if escapeChar == '\\' {
					escapeChar = 0x0
					strVal = append(strVal, scan.peak())
					scan.skip()
					continue
				}
				if scan.peak() == '\n' {
					scan.Line++
				}
				if scan.peak() == '\\' {
					escapeChar = scan.peak()
					scan.skip()
				} else {
					strVal = append(strVal, scan.peak())
					scan.skip()
				}
			}
			if scan.isAtEnd() {
				scan.throw(fmt.Sprintf("Unterminated string on line %d", scan.Line))
			}
			strVal = append(strVal, scan.peak())
			scan.skip()
			scan.Tokens = append(scan.Tokens, Token{
				TokenType: TOKEN_STRING,
				Value:     string(strVal),
				Line:      scan.Line,
			})
		},
		TOKEN_MULTI_LINE_COMMENT: func(scan *Scanner) {
			scan.scanTillDelim("*/")
		},
		TOKEN_COMMENT: func(scan *Scanner) {
			scan.scanTillDelim("\n")
		},
	}
	TOKEN_KEYWORDS = map[string]int{
		"struct":   TOKEN_STRUCT,
		"var":      TOKEN_VAR,
		"if":       TOKEN_IF,
		"else":     TOKEN_ELSE,
		"for":      TOKEN_FOR,
		"return":   TOKEN_RETURN,
		"break":    TOKEN_BREAK,
		"func":     TOKEN_FUNCTION,
		"nil":      TOKEN_NULL,
		"true":     TOKEN_TRUE,
		"false":    TOKEN_FALSE,
		"this":     TOKEN_THIS,
		"for_each": TOKEN_FOREACH,
		"continue": TOKEN_CONTINUE,
		"q":        TOKEN_QUERY,
		"query":    TOKEN_QUERY,
		"math":     TOKEN_MATH,
		"geometry": TOKEN_GEOMETRY,
		"texture":  TOKEN_TEXTURE,
		"material": TOKEN_MATERIAL,
	}
)

// Tokeniser tokens
const (
	// Single-character tokens.
	TOKEN_LEFT_PAREN = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_LEFT_BRACKET
	TOKEN_RIGHT_BRACKET
	//TOKEN_LEFT_ARROW
	//TOKEN_RIGHT_ARROW
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_SEMICOLON
	TOKEN_SINGLE_QUOTE
	TOKEN_COLON

	// One or two character tokens.
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL
	TOKEN_PLUS_EQUAL
	TOKEN_PLUS
	TOKEN_MINUS_EQUAL
	TOKEN_ENTITY_SELECTOR
	TOKEN_MINUS
	TOKEN_STAR_EQUAL
	TOKEN_STAR
	TOKEN_SLASH_EQUAL
	TOKEN_SLASH
	TOKEN_PLUS_PLUS
	TOKEN_MINUS_MINUS
	TOKEN_AND
	TOKEN_AND_AND
	TOKEN_OR
	TOKEN_OR_OR
	TOKEN_QUESTION
	TOKEN_NULL_COALESCE
	// Special tokens
	TOKEN_COMMENT
	TOKEN_MULTI_LINE_COMMENT
	TOKEN_MACRO // Will handle this later

	// Literals.
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER

	// Keywords.
	TOKEN_AND_
	TOKEN_BREAK
	TOKEN_ARRAY
	TOKEN_ELSE
	TOKEN_IF
	TOKEN_FALSE
	TOKEN_FOR
	TOKEN_FUNCTION
	TOKEN_NULL
	TOKEN_OR_
	TOKEN_RETURN
	TOKEN_TRUE
	TOKEN_VAR
	TOKEN_STRUCT // Fuck who ever asked me to add these imma cry adding these fml fml fml
	TOKEN_THIS
	TOKEN_FOREACH
	TOKEN_CONTINUE
	TOKEN_QUERY
	TOKEN_MATH
	TOKEN_GEOMETRY
	TOKEN_TEXTURE
	TOKEN_MATERIAL

	EOF
)
