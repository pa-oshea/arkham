package parser

import (
	"arkham/ast"
	"arkham/lexer"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initProgramTest(t *testing.T, input string) *ast.Program {
	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	return program
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		program := initProgramTest(t, tt.input)

		require.Len(t, program.Statements, 1, "program.Statements does not contain correct number of statements")

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	program := initProgramTest(t, input)

	require.Len(t, program.Statements, 3, "program.Statements")

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !assert.Truef(t, ok, "stmt not *ast.returnStatement. got=%T", stmt) {
			continue
		}

		assert.Equal(t, "return", returnStmt.TokenLiteral())

	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	program := initProgramTest(t, input)

	require.Len(t, program.Statements, 1, "program has not enough statements")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	ident, ok := stmt.Expression.(*ast.Identifier)
	require.Truef(t, ok, "exp nto *ast.Identifier. got=%T", stmt.Expression)

	assert.Equal(t, "foobar", ident.Value, "ident.Value")
	assert.Equal(t, "foobar", ident.TokenLiteral(), "ident.TokenLiteral")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	program := initProgramTest(t, input)

	require.Len(t, program.Statements, 1, "program has not enough statements")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	require.Truef(t, ok, "exp not *ast.IntegerLiteral. got=%T", stmt.Expression)

	assert.Equal(t, int64(5), literal.Value, "Literal value not correct")
	assert.Equal(t, "5", literal.TokenLiteral(), "literal.TokenLiteral() not correct")
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		program := initProgramTest(t, tt.input)

		require.Len(t, program.Statements, 1, "program.Statements has not enough statements")

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.Truef(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		require.Truef(t, ok, "stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		require.Equal(t, tt.operator, exp.Operator)

		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		program := initProgramTest(t, tt.input)

		require.Len(t, program.Statements, 1, "program.Statements does not contain the correct number of statements")

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.Truef(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		require.Truef(t, ok, "exp is not ast.InfixExpression. got=%T", stmt.Expression)

		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}

		require.Equal(t, tt.operator, exp.Operator)

		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}
	for _, tt := range tests {
		program := initProgramTest(t, tt.input)

		assert.Equal(t, tt.expected, program.String())
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		program := initProgramTest(t, tt.input)

		require.Len(t, program.Statements, 1, "program has not enough statements")

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement")

		boolean, ok := stmt.Expression.(*ast.Boolean)
		require.Truef(t, ok, "exp not *ast.Boolean. got=%T", stmt.Expression)

		assert.Equal(t, tt.expectedBoolean, boolean.Value)

	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	program := initProgramTest(t, input)

	require.Len(t, program.Statements, 1, "program.Body does not contain correct number of statements")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "program.Statements is not ast.ExpressionStatement. got=%T", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.Truef(t, ok, "stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	assert.Len(t, exp.Consequence.Statements, 1, "consequence is not 1 statements")

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	assert.Nil(t, exp.Alternative, "exp.Alternative.Statements was not nil")
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	program := initProgramTest(t, input)

	require.Len(t, program.Statements, 1, "program.Body does not contain enough statements")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	require.True(t, ok, "stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	require.Len(t, function.Parameters, 2, "Not enough function parameters. need 2")

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	require.Len(t, function.Body.Statements, 1, "function.Body.Statements incorrect statements")

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "Function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		program := initProgramTest(t, tt.input)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		assert.Len(t, function.Parameters, len(tt.expectedParams), "length parameters wrong")

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	program := initProgramTest(t, input)

	require.Len(t, program.Statements, 1, "program.Statements does not contain enough statements")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.CallExpression)
	require.Truef(t, ok, "stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	require.Len(t, exp.Arguments, 3, "wrong length of arguments.")
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("Type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
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

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatroExpressino. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
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
