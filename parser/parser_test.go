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
	checkParserErrors(t, p)
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
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser had %d", len(errors))
	for _, msg := range errors {
		t.Errorf("Paser Error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatement(t *testing.T) {

	input := `
	return 10;
	return 90;
	return 10123123;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != 3 {
		t.Fatalf("Wrong number of statments. Expected 3 got %d", len(program.Statements))
	}
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement, got '%T'", stmt)
			continue
		}
		if tokLiteral := returnStmt.TokenLiteral(); tokLiteral != "return" {
			t.Errorf("Return statement token literal was not 'return', got '%s'", tokLiteral)
		}

	}
}
