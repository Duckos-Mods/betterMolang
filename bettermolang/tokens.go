package bettermolang

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
	if t.NextToken == nil {
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
			// If the next node is nil then we break
			// That tells us that we have reached the end of the token
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

type Token struct {
	TokenType int
	Value     string
	Line      int
}

func (t *Token) ToString() string {
	switch t.TokenType {
	case TOKEN_BANG:
		{
			return "!"
		}
	case TOKEN_BANG_EQUAL:
		{
			return "!="
		}
	case TOKEN_EQUAL:
		{
			return "="
		}
	case TOKEN_EQUAL_EQUAL:
		{
			return "=="
		}
	case TOKEN_GREATER:
		{
			return ">"
		}
	case TOKEN_GREATER_EQUAL:
		{
			return ">="
		}
	case TOKEN_LESS:
		{
			return "<"
		}
	case TOKEN_LESS_EQUAL:
		{
			return "<="
		}
	case TOKEN_PLUS_EQUAL:
		{
			return "+="
		}
	case TOKEN_PLUS:
		{
			return "+"
		}
	case TOKEN_MINUS_EQUAL:
		{
			return "-="
		}
	case TOKEN_ENTITY_SELECTOR:
		{
			return "->"
		} // this isnt correct i need to handle this better
	case TOKEN_MINUS:
		{
			return "-"
		}
	case TOKEN_STAR_EQUAL:
		{
			return "*="
		}
	case TOKEN_STAR:
		{
			return "*"
		}
	case TOKEN_SLASH_EQUAL:
		{
			return "/="
		}
	case TOKEN_SLASH:
		{
			return "/"
		}
	case TOKEN_PLUS_PLUS:
		{
			return "++"
		}
	case TOKEN_MINUS_MINUS:
		{
			return "--"
		}
	case TOKEN_AND:
		{
			return "&"
		}
	case TOKEN_AND_AND:
		{
			return "&&"
		}
	case TOKEN_OR:
		{
			return "|"
		}
	case TOKEN_OR_OR:
		{
			return "||"
		}
	case TOKEN_QUESTION:
		{
			return "?"
		}
	case TOKEN_NULL_COALESCE:
		{
			return "??"
		}
	case TOKEN_COMMENT:
		{
			return "//"
		}
	case TOKEN_MULTI_LINE_COMMENT:
		{
			return "/*"
		}
	case TOKEN_MACRO:
		{
			return "macro"
		}
	case TOKEN_IDENTIFIER:
		{
			return "identifier"
		}
	case TOKEN_STRING:
		{
			return "string"
		}
	case TOKEN_NUMBER:
		{
			return "number"
		}
	case TOKEN_AND_:
		{
			return "and"
		}
	case TOKEN_BREAK:
		{
			return "break"
		}
	case TOKEN_ARRAY:
		{
			return "array"
		}
	case TOKEN_ELSE:
		{
			return "else"
		}
	case TOKEN_IF:
		{
			return "if"
		}
	case TOKEN_FALSE:
		{
			return "false"
		}
	case TOKEN_FOR:
		{
			return "for"
		}
	case TOKEN_FUNCTION:
		{
			return "function"
		}
	case TOKEN_NULL:
		{
			return "null"
		}
	case TOKEN_OR_:
		{
			return "or"
		}
	case TOKEN_RETURN:
		{
			return "return"
		}
	case TOKEN_TRUE:
		{
			return "true"
		}
	case TOKEN_VAR:
		{
			return "var"
		}
	case TOKEN_STRUCT:
		{
			return "struct"
		}
	case TOKEN_THIS:
		{
			return "this"
		}
	case TOKEN_FOREACH:
		{
			return "for_each"
		}
	case TOKEN_CONTINUE:
		{
			return "continue"
		}
	case TOKEN_QUERY:
		{
			return "query"
		}
	case TOKEN_MATH:
		{
			return "math"
		}
	case TOKEN_GEOMETRY:
		{
			return "geometry"
		}
	case TOKEN_TEXTURE:
		{
			return "texture"
		}
	case TOKEN_MATERIAL:
		{
			return "material"
		}
	case EOF:
		{
			return "EOF"
		}
	case TOKEN_LEFT_PAREN:
		{
			return "("
		}
	case TOKEN_RIGHT_PAREN:
		{
			return ")"
		}
	case TOKEN_LEFT_BRACE:
		{
			return "{"
		}
	case TOKEN_RIGHT_BRACE:
		{
			return "}"
		}
	case TOKEN_LEFT_BRACKET:
		{
			return "["
		}
	case TOKEN_RIGHT_BRACKET:
		{
			return "]"
		}
	case TOKEN_COMMA:
		{
			return ","
		}
	case TOKEN_DOT:
		{
			return "."
		}
	case TOKEN_SEMICOLON:
		{
			return ";"
		}
	case TOKEN_SINGLE_QUOTE:
		{
			return "'"
		}
	case TOKEN_COLON:
		{
			return ":"
		}
	default:
		{
			return "Unknown"
		}
	}
}
