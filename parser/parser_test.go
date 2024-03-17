package parser

import (
	"lang/ast"
	"lang/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {

	input := `
	let x = 5;
	let y = 10;
	let hi = 90;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if numStatements := len(program.Statements); numStatements != 3 {
		t.Fatalf("Unexpected number of statements: %d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"hi"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if literal := s.TokenLiteral(); literal != "let" {
		t.Errorf("Token Literal != 'let', was %s", literal)
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("S was not a Let Statement, found %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("Let Statement Name not '%s', found '%s'", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() !=  '%s', found '%s'", name, letStmt.TokenLiteral())
		return false
	}
	return true
}
