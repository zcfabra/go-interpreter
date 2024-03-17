package ast

import "lang/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier // the variable name to bind data to
	Value Expression
}

type Identifier struct {
	Token token.Token // IDENT type token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// let x = 10; would be represented by an AST of:
// (Program
//		(LetStatment)
// )
/*
let x = 10; would be represented by an AST of:
(Program (statements)->
		(LetStatement (
			name -> (Identifier),
			value -> (Expression)
		)
	)
)

*/
