package bettermolang

import (
	"fmt"
	"strings"
)

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
			consLeng++
			if tempNode.NextVal == nil {
				break
			}
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
		']':  TOKEN_RIGHT_BRACKET,
		'.':  TOKEN_DOT,
		',':  TOKEN_COMMA,
		'\'': TOKEN_SINGLE_QUOTE,
		'?':  TOKEN_QUESTION,
		':':  TOKEN_COLON,
		';':  TOKEN_SEMICOLON,
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
		'#': {
			TokenType: TOKEN_MACRO,
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
		"struct": TOKEN_STRUCT,
		"var":    TOKEN_VAR,
		"if":     TOKEN_IF,
		"else":   TOKEN_ELSE,
		"for":    TOKEN_FOR,
		"return": TOKEN_RETURN,
		"break":  TOKEN_BREAK,
		"func":   TOKEN_FUNCTION,
		"nil":    TOKEN_NULL,
		"true":   TOKEN_TRUE,
		"false":  TOKEN_FALSE,
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
		return
	case ' ', '\r', '\t':
		// Ignore whitespace
		return
	}
	if token, ok := TOKEN_SINGLE_CHAR_MAP[singleChar]; ok {
		if token == TOKEN_NULL {
			s.throw(fmt.Sprintf("Unknown token on line %d", s.Line))
		}
		if function, ok := TOKEN_SPECIAL_MAP[token]; ok {
			function(s)
		} else {
			s.addToken(token, 1)
		}

		return
	} else if token, ok := TOKEN_VLT_MAP[singleChar]; ok {
		verified, consumeLength := token.verifyTokenType(s)
		if funct, ok := TOKEN_SPECIAL_MAP[verified]; ok {
			funct(s)
		} else {
			if verified == TOKEN_NULL {
				s.throw(fmt.Sprintf("Unknown token on line %d", s.Line))
			}
			s.addToken(verified, consumeLength)
			return
		}
	} else if s.isDigit(singleChar) {
		s.scanNumber()
		return
	} else if s.isAlpha(singleChar) {
		s.scanIdentifier()
		return
	} else {

		s.throw(fmt.Sprintf("Unknown token on line %d, Token Val of %s", s.Line, string(singleChar)))
	}
}

func (s *Scanner) scanTillDelim(delim string) {
	compArr := make([]byte, len(delim))
	// remove all data from the current position to the delim
	for i := 0; i < len(delim); i++ {
		compArr[i] = s.advance()
	}
	for !s.isAtEnd() {
		{
			// Shift everything in the compArr to the left by 1
			for i := 0; i < len(compArr)-1; i++ {
				compArr[i] = compArr[i+1]
			}
			compArr[len(compArr)-1] = s.advance()
			if string(compArr) == delim {
				return
			}
		}
	}
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) throw(err string) {
	panic(err)
}

func (s *Scanner) skip() {
	s.Current++ // Skip the current token
}

func (s *Scanner) scanNumber() {
	consumeLength := 1
	for s.isDigit(s.peak()) {
		consumeLength++
		s.advance()
	}
	if s.peak() == '.' && s.isDigit(s.peakN(1)) {
		consumeLength++
		s.advance()
		for s.isDigit(s.peak()) {
			consumeLength++
			s.advance()
		}
	}
	s.addToken(TOKEN_NUMBER, consumeLength)
}

func (s *Scanner) scanIdentifier() {
	consumeLength := 1
	for s.isAlphaNumeric(s.peak()) {
		consumeLength++
		s.advance()
	}
	// Check if the identifier is a keyword
	identifier := strings.TrimSpace(s.Source[s.Current-consumeLength : s.Current])
	if keyword, ok := TOKEN_KEYWORDS[identifier]; ok {
		s.addToken(keyword, consumeLength)
	} else {
		s.addToken(TOKEN_IDENTIFIER, consumeLength)
	}
}

func (s *Scanner) isAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
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
		Line:    0,
	}
}
