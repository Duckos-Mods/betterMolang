package bettermolang

import (
	"fmt"
	"strings"
)

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
	TokenType     int
	NextToken     []VLTokenNode
	consumeLength int
}

func (t *VLToken) verifyTokenType(scan *Scanner) (int, int) {
	if scan.isAtEnd() {
		return t.TokenType, t.consumeLength
	}

	for _, node := range t.NextToken {
		consLeng := t.consumeLength
		bulkMatchData := make([]byte, 0)
		var tempNode *VLTokenNode = &node
		// This will loop till break!
		for i := 0; i < 1; i = 0 {
			bulkMatchData = append(bulkMatchData, tempNode.SelfVal)
			if tempNode.NextVal == nil {
				break
			}
			consLeng++
			tempNode = tempNode.NextVal
		}
		if scan.bulkMatch(bulkMatchData) {
			return tempNode.SucessType, consLeng
		}
	}
	return t.TokenType, t.consumeLength
}

var (
	TOKEN_SINGLE_CHAR_MAP = map[byte]int{
		'(':  TOKEN_LEFT_PAREN,
		')':  TOKEN_RIGHT_PAREN,
		'{':  TOKEN_LEFT_BRACE,
		'}':  TOKEN_RIGHT_BRACE,
		'[':  TOKEN_LEFT_BRACKET,
		'.':  TOKEN_DOT,
		'\'': TOKEN_SINGLE_QUOTE,
		'?':  TOKEN_QUESTION,
		':':  TOKEN_COLON,
	}

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
			},
			consumeLength: 1,
		},
		'&': {
			TokenType: TOKEN_AND,
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
			TokenType: TOKEN_OR,
			NextToken: []VLTokenNode{
				{
					SucessType: TOKEN_OR_OR,
					SelfVal:    '|',
					NextVal:    nil,
				},
			},
			consumeLength: 1,
		},
		'#': {
			TokenType: TOKEN_NULL,
			NextToken: []VLTokenNode{
				{
					SucessType: TOKEN_COMMENT,
					SelfVal:    '#',
					NextVal: &VLTokenNode{
						SucessType: TOKEN_COMMENT,
						SelfVal:    '#',
						NextVal:    nil,
					},
				},
			},
			consumeLength: 1,
		},
	}

	TOKEN_SPECIAL_MAP = map[int]func(*Scanner){
		TOKEN_SINGLE_QUOTE: func(scan *Scanner) {
			consumeLength := 1
			for scan.peak() != '\'' && !scan.isAtEnd() {
				if scan.peak() == '\n' {
					scan.Line++
				}
				consumeLength++
				scan.advance()
			}
			if scan.isAtEnd() {
				scan.throw(fmt.Sprintf("Unterminated string on line %d", scan.Line))
			}
			scan.advance()
			consumeLength++
			scan.addToken(TOKEN_STRING, consumeLength)
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
	// Special tokens
	TOKEN_COMMENT

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
}

func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) advance() byte {
	s.Current++
	return s.Source[s.Current-1]
}

func (s *Scanner) peak() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.Source[s.Current]
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
		if !s.match(char) {
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

func (s *Scanner) addToken(tokenType int, consumeLength int) {
	s.Tokens = append(s.Tokens, Token{
		TokenType: tokenType,
		Value:     strings.TrimSpace(s.Source[s.Current-consumeLength : s.Current]),
		Line:      s.Line,
	})
}

func (s *Scanner) scanToken() {
	var singleChar byte = s.advance()
	switch singleChar {
	case '\n':
		s.Line++
	case ' ', '\r', '\t':
		// Ignore whitespace
		break
	}
	if token, ok := TOKEN_SINGLE_CHAR_MAP[singleChar]; ok {
		if function, ok := TOKEN_SPECIAL_MAP[token]; ok {
			function(s)
		} else {
			s.addToken(token, 1)
		}

		return
	} else if token, ok := TOKEN_VLT_MAP[singleChar]; ok {
		verified, consumeLength := token.verifyTokenType(s)
		if verified != TOKEN_COMMENT {
			s.addToken(verified, consumeLength)
		} else {
			s.removeComment()
		}
		return
	}
}

func (s *Scanner) removeComment() {
	// Just remove the comment
	for s.peak() != '\n' && !s.isAtEnd() {
		s.advance()
	}
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) throw(err string) {
	panic(err)
}

func (s *Scanner) ScanTokens(code string) []Token {
	s.Source = code
	for !s.isAtEnd() {
		s.scanToken()
	}
	return s.Tokens
}

func NewScanner() *Scanner {
	return &Scanner{
		Source:  "",
		Tokens:  make([]Token, 0),
		Current: 0,
		Line:    1,
	}
}
