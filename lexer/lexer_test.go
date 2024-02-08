package lexer

import (
	"arkham/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	let add = fn(x, y) {
		x + y;
	};

	let result = add(five, ten);
	!-/*5;
	5 < 10 > 5;
	
	if (5 < 10) {
		return true;
	} else {
		return false;
	}

	10 == 10;
	10 != 9;`

	tests := []struct {
		index           int
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{0, token.LET, "let"},
		{1, token.IDENT, "five"},
		{2, token.ASSIGN, "="},
		{3, token.INT, "5"},
		{4, token.SEMICOLON, ";"},
		{5, token.LET, "let"},
		{6, token.IDENT, "ten"},
		{7, token.ASSIGN, "="},
		{8, token.INT, "10"},
		{9, token.SEMICOLON, ";"},
		{10, token.LET, "let"},
		{11, token.IDENT, "add"},
		{12, token.ASSIGN, "="},
		{13, token.FUNCTION, "fn"},
		{14, token.LPAREN, "("},
		{15, token.IDENT, "x"},
		{16, token.COMMA, ","},
		{17, token.IDENT, "y"},
		{18, token.RPAREN, ")"},
		{19, token.LBRACE, "{"},
		{20, token.IDENT, "x"},
		{21, token.PLUS, "+"},
		{22, token.IDENT, "y"},
		{23, token.SEMICOLON, ";"},
		{24, token.RBRACE, "}"},
		{25, token.SEMICOLON, ";"},
		{26, token.LET, "let"},
		{27, token.IDENT, "result"},
		{28, token.ASSIGN, "="},
		{29, token.IDENT, "add"},
		{30, token.LPAREN, "("},
		{31, token.IDENT, "five"},
		{32, token.COMMA, ","},
		{33, token.IDENT, "ten"},
		{34, token.RPAREN, ")"},
		{35, token.SEMICOLON, ";"},
		{36, token.BANG, "!"},
		{37, token.MINUS, "-"},
		{38, token.SLASH, "/"},
		{39, token.ASTERISK, "*"},
		{40, token.INT, "5"},
		{41, token.SEMICOLON, ";"},
		{42, token.INT, "5"},
		{43, token.LT, "<"},
		{44, token.INT, "10"},
		{45, token.GT, ">"},
		{46, token.INT, "5"},
		{47, token.SEMICOLON, ";"},
		{48, token.IF, "if"},
		{49, token.LPAREN, "("},
		{50, token.INT, "5"},
		{51, token.LT, "<"},
		{52, token.INT, "10"},
		{53, token.RPAREN, ")"},
		{54, token.LBRACE, "{"},
		{55, token.RETURN, "return"},
		{56, token.TRUE, "true"},
		{57, token.SEMICOLON, ";"},
		{58, token.RBRACE, "}"},
		{59, token.ELSE, "else"},
		{60, token.LBRACE, "{"},
		{61, token.RETURN, "return"},
		{62, token.FALSE, "false"},
		{63, token.SEMICOLON, ";"},
		{64, token.RBRACE, "}"},
		{65, token.INT, "10"},
		{66, token.EQ, "=="},
		{67, token.INT, "10"},
		{68, token.SEMICOLON, ";"},
		{69, token.INT, "10"},
		{70, token.NOT_EQ, "!="},
		{71, token.INT, "9"},
		{72, token.SEMICOLON, ";"},
		{73, token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		require.Equalf(t, tt.expectedType, tok.Type, "Test[%d] tokentype wrong", i)
		require.Equalf(t, tt.expectedLiteral, tok.Literal, "Test[%d]", i)
	}
}

func TestIsLetter(t *testing.T) {
	assert.Equal(t, true, isLetter('a'))
	assert.Equal(t, true, isLetter('z'))
	assert.Equal(t, true, isLetter('c'))
	assert.Equal(t, true, isLetter('A'))
	assert.Equal(t, true, isLetter('Z'))
	assert.Equal(t, true, isLetter('R'))
	assert.Equal(t, true, isLetter('_'))
	assert.Equal(t, false, isLetter('5'))
}

func TestIsDigit(t *testing.T) {
	assert.Equal(t, true, isDigit('1'))
	assert.Equal(t, true, isDigit('5'))
}
