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

/*
Parsing expressions:
Operators can be prefix (e.g -5, !some_bool) or infix (eg. 3 - 5)
We can also "boost" the precendence of an expression by grouping it w/ parentheses
E.g 10 * 5 + 5 vs 10 * (5 + 5)

Top down operator precedence (Pratt Parsing)
Parsing functions (e.g. `parseLetStatement`, `parseReturnStatement`) are
associated w/ token types, rather than grammar rules

Each token type can have two parsing fns associated w/ it depending on the
token's position -- this allows distinguishing btwn. infix and prefix operator


How the parser ACTUALLY works

*/

import (
	"fmt"
	"lang/ast"
	"lang/lexer"
	"lang/token"
	"strconv"
)

// Operator Precedence
const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParserFns map[token.TokenType]prefixParseFn
	infixParserFns  map[token.TokenType]infixParseFn
}

/*
The argument for the infix function is the 'lhs' of the infix
operator that's being parsed
*/
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Register fns for token types
func (p *Parser) registerPrefixFn(tt token.TokenType, fn prefixParseFn) {
	p.prefixParserFns[tt] = fn
}

func (p *Parser) registerInfixFn(tt token.TokenType, fn infixParseFn) {
	p.infixParserFns[tt] = fn
}

// Errors
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token type to be '%s', found '%s'", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
func (p *Parser) noPrefixParseError(tt token.TokenType) {
	msg := fmt.Sprintf("Expected a valid prefix for %s", tt)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	// The 'peek' token is actually the latest consumed token
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// While not ended
	for p.curToken.Type != token.EOF {
		// Parse statement + advance
		if stmt := p.ParseStatement(); stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) ParseStatement() ast.Statement {
	// Parse statements by type
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Grab the prefix parsing function, and apply it (if any are found)
	prefix := p.prefixParserFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	right := p.parseExpression(PREFIX)
	expr.Right = right
	return expr
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	// TODO: Parse expression

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: parse expressions rather than skipping them
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) expectPeek(tt token.TokenType) bool {
	// Expect peek is an 'assertion' function (very commmon in parsers)
	// Meant to enforce correctness of token order
	if p.peekTokenIs(tt) {
		p.nextToken()
		return true
	}
	p.peekError(tt)
	return false
}
func (p *Parser) peekTokenIs(tt token.TokenType) bool {
	return p.peekToken.Type == tt
}

func (p *Parser) curTokenIs(tt token.TokenType) bool {
	return p.curToken.Type == tt
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("%q Could Not Be Parsed As Int", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: value}
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{l: lexer, errors: []string{}}
	p.prefixParserFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefixFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	// Sets curToken to peekToken (which is nil at this point), sets peekToken = 0
	p.nextToken()
	// Sets curToken to peekToken (which is now 0), sets peekToken = 1
	p.nextToken()

	return p
}
