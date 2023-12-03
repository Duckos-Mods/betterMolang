package bettermolang

import (
	"fmt"
	"strings"
)

// Tokeniser token struct

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
	if token, ok := TOKEN_VLT_MAP[singleChar]; ok {
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
