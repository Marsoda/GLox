package parser

import "LoxGo/scanner"

type Stmt interface {
	Accept(visitor StmtVisitor)
}

type ExprStmt struct {
	Expr Expr
}

func NewExprStmt(expr Expr) *ExprStmt {
	return &ExprStmt{expr}
}

func (e *ExprStmt) Accept(visitor StmtVisitor) {
	visitor.VisitExprStmt(e)
}

type PrintStmt struct {
	Expr Expr
}

func NewPrintStmt(expr Expr) *PrintStmt {
	return &PrintStmt{expr}
}

func (p *PrintStmt) Accept(visitor StmtVisitor) {
	visitor.VisitPrintStmt(p)
}

type VarDeclStmt struct {
	Name        *scanner.Token
	Initializer Expr
}

func NewVarDeclStmt(name *scanner.Token, initializer Expr) *VarDeclStmt {
	return &VarDeclStmt{Name: name, Initializer: initializer}
}

func (v *VarDeclStmt) Accept(visitor StmtVisitor) {
	visitor.VisitVarDeclStmt(v)
}
