package parser_test

import (
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

	ident, ok := expStmt.Expression.(*ast.Identifier)
	assert.True(t, ok)
	assert.Equal(t, "something", ident.Value)
	assert.Equal(t, "something", ident.TokenLiteral())
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

	lit, ok := expStmt.Expression.(*ast.IntegerLiteral)
	assert.True(t, ok)
	assert.Equal(t, int64(5), lit.Value)
	assert.Equal(t, "5", lit.TokenLiteral())
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	assert.Len(t, errors, 0)
}
