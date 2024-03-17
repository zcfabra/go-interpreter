package parser

/*
How recursive descent parsers work:
Assuming these node types (Program, Let Stmt., Expression, Identifier)

Top Level (Program):
Scan tokens until a token such as let, return, if is found -- they indicate
that a statment has begun
Call a fn that handles each of the various statement types

Next level down (Statements):
Let statement: requires a variable (identifier) and a value
	|> Drop down a level to extract the identifier --> Pop back up
	|> Parse expression (i.e. get computation to find value) --> Pop back up
	|> Create new node w/ the identifier and value

Next level down (Expressions):
Various expresion types (Unary, Binary, etc) are parsed here
*/

import (
	"lang/ast"
	"lang/lexer"
	"lang/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func (p *Parser) nextToken() {
	// The 'peek' token is actually the latest consumed token
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
func (p *Parser) ParseProgram() *ast.Program {
	return nil
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{l: lexer}
	p.nextToken()
	p.nextToken()

	return p
}
