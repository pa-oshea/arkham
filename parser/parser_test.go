package parser

import (
	"arkham/ast"
	"arkham/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.NotNil(t, program, "ParseProgram return nil")
	require.Len(t, program.Statements, 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatment(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	require.Len(t, program.Statements, 3, "program.Statments")

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatment)

		if !assert.Truef(t, ok, "stmt not *ast.returnStatment. got=%T", stmt) {
			continue
		}

		assert.Equal(t, "return", returnStmt.TokenLiteral())

	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	require.Lenf(t, program.Statements, 1, "program has not enough statments. got=%d", len(program.Statements))

	stmt, ok := program.Statements[0].(*ast.ExpressionStatment)
	require.Truef(t, ok, "program.Statements[0] is not ast.ExpressionStatment. got=%T", program.Statements[0])

	ident, ok := stmt.Expression.(*ast.Identifier)
	require.Truef(t, ok, "exp nto *ast.Identifier. got=%T", stmt.Expression)

	assert.Equal(t, "foobar", ident.Value, "ident.Value")
	assert.Equal(t, "foobar", ident.TokenLiteral(), "ident.TokenLiteral")
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
