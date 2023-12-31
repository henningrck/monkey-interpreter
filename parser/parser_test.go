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
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let something = y;", "something", "y"},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		assert.Len(t, program.Statements, 1)

		letStmt, ok := program.Statements[0].(*ast.LetStatement)
		assert.True(t, ok)
		assert.Equal(t, test.expectedIdentifier, letStmt.Name.Value)
		checkLiteral(t, letStmt.Value, test.expectedValue)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue any
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return 993322;", 993322},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		assert.Len(t, program.Statements, 1)

		retStmt, ok := program.Statements[0].(*ast.ReturnStatement)
		assert.True(t, ok)
		checkLiteral(t, retStmt.ReturnValue, test.expectedValue)
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

func TestFunctionLiteralExpressions(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	assert.Len(t, program.Statements, 1)

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	function, ok := expStmt.Expression.(*ast.FunctionLiteral)
	assert.True(t, ok)
	assert.Len(t, function.Parameters, 2)
	checkLiteral(t, function.Parameters[0], "x")
	checkLiteral(t, function.Parameters[1], "y")
	assert.Len(t, function.Body.Statements, 1)

	bodyExpStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	checkInfixExpression(t, bodyExpStmt.Expression, "x", "+", "y")
}

func TestFunctionLiteralParameters(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y) {};", expectedParams: []string{"x", "y"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		assert.Len(t, program.Statements, 1)

		expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		funcLit, ok := expStmt.Expression.(*ast.FunctionLiteral)
		assert.True(t, ok)
		assert.Len(t, funcLit.Parameters, len(test.expectedParams))

		for i, ident := range test.expectedParams {
			checkLiteral(t, funcLit.Parameters[i], ident)
		}
	}
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
		checkInfixExpression(t, expStmt.Expression, test.leftValue, test.operator, test.rightValue)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	assert.Len(t, program.Statements, 1)

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	ifExp, ok := expStmt.Expression.(*ast.IfExpression)
	assert.True(t, ok)
	checkInfixExpression(t, ifExp.Condition, "x", "<", "y")
	assert.Len(t, ifExp.Consequence.Statements, 1)
	assert.Nil(t, ifExp.Alternative)

	consequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	checkLiteral(t, consequence.Expression, "x")
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	assert.Len(t, program.Statements, 1)

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	ifExp, ok := expStmt.Expression.(*ast.IfExpression)
	assert.True(t, ok)
	checkInfixExpression(t, ifExp.Condition, "x", "<", "y")
	assert.Len(t, ifExp.Consequence.Statements, 1)
	assert.Len(t, ifExp.Alternative.Statements, 1)

	consequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	checkLiteral(t, consequence.Expression, "x")

	alternative, ok := ifExp.Alternative.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	checkLiteral(t, alternative.Expression, "y")
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	assert.Len(t, program.Statements, 1)

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	callExp, ok := expStmt.Expression.(*ast.CallExpression)
	assert.True(t, ok)
	checkLiteral(t, callExp.Function, "add")
	assert.Len(t, callExp.Arguments, 3)
	checkLiteral(t, callExp.Arguments[0], 1)
	checkInfixExpression(t, callExp.Arguments[1], 2, "*", 3)
	checkInfixExpression(t, callExp.Arguments[2], 4, "+", 5)
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

func checkInfixExpression(t *testing.T, exp ast.Expression, leftValue any, operator string, rightValue any) {
	infixExp, ok := exp.(*ast.InfixExpression)
	assert.True(t, ok)
	assert.Equal(t, operator, infixExp.Operator)

	checkLiteral(t, infixExp.Left, leftValue)
	checkLiteral(t, infixExp.Right, rightValue)
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
