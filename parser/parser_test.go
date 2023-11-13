package parser_test

import (
	"fmt"
	"testing"

	"github.com/henningrck/monkey-interpreter/ast"
	"github.com/henningrck/monkey-interpreter/lexer"
	"github.com/henningrck/monkey-interpreter/parser"
	"github.com/stretchr/testify/assert"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
	let y = 10;
	let something = 838383;`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	assert.Len(t, program.Statements, 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"something"},
	}

	for i, test := range tests {
		stmt := program.Statements[i]
		letStmt, ok := stmt.(*ast.LetStatement)
		assert.True(t, ok)
		assert.Equal(t, "let", stmt.TokenLiteral())
		assert.Equal(t, test.expectedIdentifier, letStmt.Name.Value)
		assert.Equal(t, test.expectedIdentifier, letStmt.Name.TokenLiteral())
	}
}

func TestReturnStatements(t *testing.T) {
	input := `return 5;
	return 10;
	return 993322;`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	assert.Len(t, program.Statements, 3)

	for _, stmt := range program.Statements {
		_, ok := stmt.(*ast.ReturnStatement)
		assert.True(t, ok)
		assert.Equal(t, "return", stmt.TokenLiteral())
	}
}

func TestIdentifierExpressions(t *testing.T) {
	input := `something;`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	assert.Len(t, program.Statements, 1)

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	assert.Equal(t, "something", expStmt.TokenLiteral())

	checkLiteral(t, expStmt.Expression, "something")
}

func TestIntegerLiteralExpressions(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	assert.Len(t, program.Statements, 1)

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	assert.Equal(t, "5", expStmt.TokenLiteral())
	checkLiteral(t, expStmt.Expression, 5)
}

func TestBooleanLiteralExpressions(t *testing.T) {
	input := `true;`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	assert.Len(t, program.Statements, 1)

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	assert.Equal(t, "true", expStmt.TokenLiteral())

	checkLiteral(t, expStmt.Expression, true)
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		assert.Len(t, program.Statements, 1)

		expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		prefixExp, ok := expStmt.Expression.(*ast.PrefixExpression)
		assert.True(t, ok)
		assert.Equal(t, test.operator, prefixExp.Operator)
		checkLiteral(t, prefixExp.Right, test.value)
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
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

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		assert.Len(t, program.Statements, 1)

		expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		infixExp, ok := expStmt.Expression.(*ast.InfixExpression)
		assert.True(t, ok)
		assert.Equal(t, test.operator, infixExp.Operator)

		checkLiteral(t, infixExp.Left, test.leftValue)
		checkLiteral(t, infixExp.Right, test.rightValue)
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
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		assert.Equal(t, test.expected, program.String())
	}
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	assert.Len(t, errors, 0)
}

func checkLiteral(t *testing.T, exp ast.Expression, expected any) {
	switch value := expected.(type) {
	case string:
		checkIdentifier(t, exp, value)
	case int:
		checkIntegerLiteral(t, exp, int64(value))
	case int64:
		checkIntegerLiteral(t, exp, value)
	case bool:
		checkBooleanLiteral(t, exp, value)
	default:
		t.Errorf("type of exp not handled, got %T", exp)
		t.Fail()
	}
}

func checkIdentifier(t *testing.T, exp ast.Expression, value string) {
	ident, ok := exp.(*ast.Identifier)
	assert.True(t, ok)
	assert.Equal(t, value, ident.Value)
	assert.Equal(t, value, ident.TokenLiteral())
}

func checkIntegerLiteral(t *testing.T, exp ast.Expression, value int64) {
	lit, ok := exp.(*ast.IntegerLiteral)
	assert.True(t, ok)
	assert.Equal(t, value, lit.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), lit.TokenLiteral())
}

func checkBooleanLiteral(t *testing.T, exp ast.Expression, value bool) {
	lit, ok := exp.(*ast.BooleanLiteral)
	assert.True(t, ok)
	assert.Equal(t, value, lit.Value)
	assert.Equal(t, fmt.Sprintf("%t", value), lit.TokenLiteral())
}
