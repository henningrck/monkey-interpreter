package ast_test

import (
	"testing"

	"github.com/henningrck/monkey-interpreter/ast"
	"github.com/henningrck/monkey-interpreter/token"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	assert.Equal(t, "let myVar = anotherVar;", program.String())
}
