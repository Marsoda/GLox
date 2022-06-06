package parser

import (
	"GLox/utils"
	"bytes"
)

// Printer ExprVisitor 子类之一，以特殊的形式打印出语法树上的节点
type Printer struct {
}

func (p *Printer) VisitBinaryExpr(expr *Binary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *Printer) VisitGroupingExpr(expr *Grouping) interface{} {
	return p.parenthesize("group", expr.Expression)
}

func (p *Printer) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	// 打印value的字面量
	return utils.ToString(expr.Value)
}

func (p *Printer) VisitUnaryExpr(expr *Unary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *Printer) VisitVariableExpr(expr *Variable) interface{} {
	// empty implementation

	return nil
}

func (p *Printer) VisitAssignExpr(expr *Assign) interface{} {
	// empty implementation

	return nil
}

func (p *Printer) VisitLogicExpr(expr *Logic) interface{} {
	// empty implementation

	return nil
}

func (p *Printer) VisitCallExpr(call *Call) interface{} {
	// empty implementation

	return nil
}

func (p *Printer) VisitGetExpr(expr *Get) interface{} {
	// empty implementation

	return nil
}

func (p *Printer) VisitSetExpr(expr *Set) interface{} {
	// empty implementation	

	return nil
}

func (p *Printer) VisitThisExpr(expr *This) interface{} {
	// empty implementation	

	return nil
}

func (p *Printer) VisitSuperExpr(expr *Super) interface{} {
	// empty implementation	

	return nil
}

func (p *Printer) parenthesize(name string, exprs ...Expr) string {
	var buffer bytes.Buffer
	buffer.WriteString("(" + name)
	for _, expr := range exprs {
		buffer.WriteString(" ")
		buffer.WriteString(expr.Accept(p).(string))
	}
	buffer.WriteString(")")

	return buffer.String()
}

func (p *Printer) Print(expr Expr) string {
	return expr.Accept(p).(string)
}
