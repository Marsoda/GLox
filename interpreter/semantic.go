package interpreter

import (
	"GLox/parser"
	"GLox/scanner"
	"fmt"
)

func (i *Interpreter) VisitBinaryExpr(expr *parser.Binary) interface{} {
	// (递归)计算左右子表达式的值
	lv, rv := i.evaluate(expr.Left), i.evaluate(expr.Right)
	switch expr.Operator.Type {
	case scanner.MINUS:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) - rv.(float64)
	case scanner.STAR:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) * rv.(float64)
	case scanner.SLASH:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) / rv.(float64)
	// 加法操作可以定义在数字和字符之上
	case scanner.PLUS:
		return doPlus(expr.Operator, lv, rv)
	case scanner.GREATER:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) > rv.(float64)
	case scanner.GREATER_EQUAL:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) >= rv.(float64)
	case scanner.LESS:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) < rv.(float64)
	case scanner.LESS_EQUAL:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) <= rv.(float64)
	// == 和 != 运算的结果是bool类型
	case scanner.BANG_EQUAL:
		return !isEqual(lv, rv)
	case scanner.EQUAL_EQUAL:
		return isEqual(lv, rv)
	}
	return nil
}

func (i *Interpreter) VisitGroupingExpr(expr *parser.Grouping) interface{} {
	// 计算中间部分的expression即可
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *parser.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitUnaryExpr(expr *parser.Unary) interface{} {
	// 先计算右侧表达式的值
	rv := i.evaluate(expr.Right)
	switch expr.Operator.Type {
	case scanner.MINUS:
		checkNumberOperands(expr.Operator, rv)
		return -(rv.(float64))
	case scanner.BANG:
		return !isTruth(rv)
	}

	return nil
}

func (i *Interpreter) VisitVariableExpr(expr *parser.Variable) interface{} {
	// return i.environment.lookup(expr.Name)
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) VisitAssignExpr(expr *parser.Assign) interface{} {
	// 计算Assign的语法树上的value节点
	value := i.evaluate(expr.Value)
	/*
		i.environment.assign(expr.Name, value)

		// 因为赋值也是一个表达式，所以这里返回所求的value
		return value
	*/
	if distance, ok := i.locals[expr]; ok {
		i.environment.assignAt(distance, expr.Name, value)
	} else {
		i.globals.assign(expr.Name, value)
	}

	return nil
}

func (i *Interpreter) VisitLogicExpr(expr *parser.Logic) interface{} {
	left := i.evaluate(expr.Left)
	if expr.Operator.Type == scanner.OR {
		if isTruth(left) {
			return left
		}
	} else {
		if !isTruth(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitCallExpr(call *parser.Call) interface{} {
	callee, ok := i.evaluate(call.Callee).(LoxCallable)
	if !ok {
		panic(NewRuntimeError(call.Paren, "Can only call functions and classes."))
	}

	var args []interface{}
	for _, arg := range call.Arguments {
		args = append(args, i.evaluate(arg))
	}

	// 判断实参和形参的个数是否相同
	if len(args) != callee.Arity() {
		panic(NewRuntimeError(call.Paren, fmt.Sprintf("Expect %d arguments buf got %d.", len(args), callee.Arity())))
	}

	return callee.Call(i, args)
}

func (i *Interpreter) VisitExprStmt(stmt *parser.ExprStmt) {
	i.evaluate(stmt.Expr)
}

func (i *Interpreter) VisitFuncDeclStmt(stmt *parser.FuncDeclStmt) {
	// 结束函数定义的区别在于，会创建一个保存了函数节点引用的新变量
	function := NewLoxFunction(stmt, i.environment)
	i.environment.define(stmt.Name, function)
}

func (i *Interpreter) VisitReturnStmt(stmt *parser.ReturnStmt) {
	var value interface{}
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}

	panic(NewReturn(value))
}

func (i *Interpreter) VisitPrintStmt(stmt *parser.PrintStmt) {
	value := i.evaluate(stmt.Expr)
	// 需要打印计算的值
	fmt.Printf("%v\n", value)
}

func (i *Interpreter) VisitVarDeclStmt(stmt *parser.VarDeclStmt) {
	var value interface{}
	if stmt.Initializer != nil {
		// 对变量的初始化语句求值
		value = i.evaluate(stmt.Initializer)
	}
	i.environment.define(stmt.Name, value)
}

func (i *Interpreter) VisitBlockStmt(stmt *parser.BlockStmt) {
	// 把当前作用域的env传入下一个block
	i.executeBlock(stmt, NewEnvironment(i.environment))
}

func (i *Interpreter) VisitIfStmt(stmt *parser.IfStmt) {
	if isTruth(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ThenBranch)
	}
}

func (i *Interpreter) VisitWhileStmt(stmt *parser.WhileStmt) {
	for isTruth(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
}

func (i *Interpreter) VisitClassDeclStmt(stmt *parser.ClassDeclStmt) {
	i.environment.define(stmt.Name, nil)
	class := NewLoxClass(stmt.Name.Lexeme)
	i.environment.assign(stmt.Name, class)
}
