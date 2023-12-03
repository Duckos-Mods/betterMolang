package bettermolang

import "strings"

// Tokeniser regexes
// const (
// 	r_COMMENT    = `/[*][^*]*[*]+([^/*][^*]*[*]+)*/|//[^\n]*`
// 	r_STRING     = `[\"\'].*[\"\']`
// 	r_OPERATOR   = `\:|\?|=|-=|/=|\*=|\+=|\+|-|/|\*|>|<|>=|<=|==|!=|!|&&|`
// 	r_SEPARATOR  = `\(|\{|\[|\]|\}|\)|;|,`
// 	r_NUMBER     = `[1-9]\d*(\.\d+)?`
// 	r_IDENTIFIER = `[a-z]+[0-9_\-]*`
// 	r_MACRO      = `\#[a-z]+`
// )

// Create something like a linked list of tokens

type VLTokenNode struct {
	SucessType int
	SelfVal    byte
	NextVal    *VLTokenNode
}

// VL stands for Variable Length
type VLToken struct {
	TokenType int
	NextToken []VLTokenNode
}

func (t *VLToken) verifyTokenbType(scan *Scanner) int {
	if scan.isAtEnd() {
		return t.TokenType
	}
	for _, node := range t.NextToken {
		bulkMatchData := make([]byte, 0)
		var tempNode *VLTokenNode = &node
		// This will loop till break!
		for i := 0; i < 1; i = 0 {
			bulkMatchData = append(bulkMatchData, tempNode.SelfVal)
			if tempNode.NextVal == nil {
				break
			}
			tempNode = tempNode.NextVal
		}
		if scan.bulkMatch(bulkMatchData) {
			return tempNode.SucessType
		}
	}
	return t.TokenType
}

var (
	TOKEN_SINGLE_CHAR_MAP = map[byte]int{
		'(': TOKEN_LEFT_PAREN,
		')': TOKEN_RIGHT_PAREN,
		'{': TOKEN_LEFT_BRACE,
		'}': TOKEN_RIGHT_BRACE,
		'[': TOKEN_LEFT_BRACKET,
		'.': TOKEN_DOT,
		'"': TOKEN_DOUBLE_QUOTE,
		'?': TOKEN_QUESTION,
		':': TOKEN_COLON,
	}

	// We use the first byte of the token as the key
	TOKEN_VLT_MAP = map[byte]VLToken{
		'!': VLToken{
			TokenType: TOKEN_BANG,
			NextToken: []VLTokenNode{
				VLTokenNode{
					SucessType: TOKEN_BANG_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
			},
		},
		'=': VLToken{
			TokenType: TOKEN_EQUAL,
			NextToken: []VLTokenNode{
				VLTokenNode{
					SucessType: TOKEN_EQUAL_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
			},
		},
		'>': VLToken{
			TokenType: TOKEN_GREATER,
			NextToken: []VLTokenNode{
				VLTokenNode{
					SucessType: TOKEN_GREATER_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
			},
		},
		'<': VLToken{
			TokenType: TOKEN_LESS,
			NextToken: []VLTokenNode{
				VLTokenNode{
					SucessType: TOKEN_LESS_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
			},
		},
		'+': VLToken{
			TokenType: TOKEN_PLUS,
			NextToken: []VLTokenNode{
				VLTokenNode{
					SucessType: TOKEN_PLUS_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
				VLTokenNode{
					SucessType: TOKEN_PLUS_PLUS,
					SelfVal:    '+',
					NextVal:    nil,
				},
			},
		},
		'-': VLToken{
			TokenType: TOKEN_MINUS,
			NextToken: []VLTokenNode{
				VLTokenNode{
					SucessType: TOKEN_MINUS_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
				VLTokenNode{
					SucessType: TOKEN_MINUS_MINUS,
					SelfVal:    '-',
					NextVal:    nil,
				},
			},
		},
		'*': VLToken{
			TokenType: TOKEN_STAR,
			NextToken: []VLTokenNode{
				VLTokenNode{
					SucessType: TOKEN_STAR_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
			},
		},
		'/': VLToken{
			TokenType: TOKEN_SLASH,
			NextToken: []VLTokenNode{
				VLTokenNode{
					SucessType: TOKEN_SLASH_EQUAL,
					SelfVal:    '=',
					NextVal:    nil,
				},
			},
		},
		'&': VLToken{
			TokenType: TOKEN_AND,
			NextToken: []VLTokenNode{
				VLTokenNode{
					SucessType: TOKEN_AND_AND,
					SelfVal:    '&',
					NextVal:    nil,
				},
			},
		},
		'|': VLToken{
			TokenType: TOKEN_OR,
			NextToken: []VLTokenNode{
				VLTokenNode{
					SucessType: TOKEN_OR_OR,
					SelfVal:    '|',
					NextVal:    nil,
				},
			},
		},
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
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_QUESTION
	TOKEN_SEMICOLON
	TOKEN_DOUBLE_QUOTE
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

	EOF
)

// Tokeniser token struct
type Token struct {
	TokenType int
	Value     string
	Line      int
}

type Scanner struct {
	Source  string
	Tokens  []Token
	Current int
	Line    int
	Start   int
}

func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) advance() byte {
	s.Current++
	return s.Source[s.Current-1]
}

func (s *Scanner) peekNext() byte {
	if s.Current+1 >= len(s.Source) {
		return 0
	}
	return s.Source[s.Current+1]
}

func (s *Scanner) peakN(depth int) byte {
	if s.Current+depth >= len(s.Source) {
		return 0
	}
	return s.Source[s.Current+depth]
}

func (s *Scanner) bulkMatch(expected []byte) bool {
	if s.isAtEnd() {
		return false
	}
	// We do this so we can reset the index if it fails
	tempIndex := s.Current
	for _, char := range expected {
		if s.Source[s.Current] != char {
			return false
		}
		tempIndex++
	}
	s.Current = tempIndex
	return true
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.Source[s.Current] != expected {
		return false
	}
	s.Current++
	return true
}

func (s *Scanner) addToken(tokenType int) {
	s.Tokens = append(s.Tokens, Token{
		TokenType: tokenType,
		Value:     strings.TrimSpace(s.Source[s.Start:s.Current]),
		Line:      s.Line,
	})
}

func (s *Scanner) scanToken() {
	var singleChar byte = s.advance()
	if token, ok := TOKEN_SINGLE_CHAR_MAP[singleChar]; ok {
		s.addToken(token)
		return
	} else {
		panic("Unknown token")
	}
}

func (s *Scanner) ScanTokens(code string) []Token {
	s.Source = code
	for !s.isAtEnd() {
		s.scanToken()
	}
	return s.Tokens
}
