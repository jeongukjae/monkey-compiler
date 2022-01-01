package ast

import (
	"monkey/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProgramString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	require.Equal(t, "let myVar = anotherVar;", program.String(), "Wrong program.String()")
}
