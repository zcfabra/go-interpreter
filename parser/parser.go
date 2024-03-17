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


How the parser ACTUALLY works:
** The goal is to have operators w/ higher precedence be lower in the AST
(i.e. they get evaluated first)

Example step through: 1 + 2 + 3;

1. Cur token is 1, peek token is +
parseExpression()
|> get prefix fn
|> parseIntegerLiteral()
|> infix for loop
|> pass in result of parseIntegerLiteral as leftExpr
2. Cur token is +, peek token is 2
|> get function parseInfixExpression()
|> get precedence of next token (2)
|> parse right expression
|> return InfixEXpression
3. Cur token is 2, peek token is +
|> get prefix fn
|> call parseIntegerLiteral()
|> Infix for loop skipped (curToken of 2 has lower precedence than peekToken +)
|> return prefix only (2 as an integer literal)
|> have now popped back to the outermost parseExpression frame -- the infix expression
made from 1 + 2 is now the lhs of an infix
4. Cur token is +, peek token is 3
|> call parseInfixExpression()
|> ...
|> once it pops back up, lhs = an infix w/ + as the operator, 3 (as an IntegerLiteral)
as the RHS and the previous InfixExpression (1 + 2) as the LHS

At this point, the AST is (+ (+ 1 2) 3)

In parseExpression when the precedence is evaluated, it represents the "right binding power"
of the current parseExpression invocation --> the amount of righwards tokens that can be binded
We saw that then the precendence was +, we could bind int literals to the right,
but when the precedence was an int and the next token was +, it could not be bound

Left binding power can be found through the peekPrecedence() call
the check the for loop after prefix parsing in parseExpression()
if all about weighing left binding power vs right binding power

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

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
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
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParserFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
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

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// Current token is the operator -- take in an already parsed expr as the
	// lhs
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.CurPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)
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

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) CurPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
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
	// Prefix fns
	p.prefixParserFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefixFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	// Infix fns
	p.infixParserFns = make(map[token.TokenType]infixParseFn)
	p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFn(token.EQ, p.parseInfixExpression)
	p.registerInfixFn(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixFn(token.LT, p.parseInfixExpression)
	p.registerInfixFn(token.GT, p.parseInfixExpression)
	// Sets curToken to peekToken (which is nil at this point), sets peekToken = 0
	p.nextToken()
	// Sets curToken to peekToken (which is now 0), sets peekToken = 1
	p.nextToken()

	return p
}
