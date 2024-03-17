package ast

import (
	"lang/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name:  &Identifier{Token: token.Token{Type: token.IDENT, Literal: "hi"}, Value: "hi"},
				Value: &Identifier{Token: token.Token{Type: token.IDENT, Literal: "other_var"}, Value: "other_var"},
			},
		},
	}

	if programString := program.String(); programString != "let hi = other_var;" {
		t.Errorf("program.string() did not match expected value, got '%s'", programString)
	}

}
