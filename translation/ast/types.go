package ast

/***************************************************************************************/

type Expression struct {
	Raw string
	XML *XMLStatement
}

func NewExpression(raw string) *Expression {
	expr := new(Expression)
	expr.Raw = raw
	return expr
}

/***************************************************************************************/

type ConditionStatement struct {
	AstNode
	Type int
	Expr *Expression
}

const (
	ConditionIf = 1 + iota
	ConditionElseif
	ConditionElse
	ConditionWhile
)

func (n *AstNode) AddConditionStatement(expr *Expression, ctype int) *ConditionStatement {
	cond := new(ConditionStatement)
	n.init(cond)
	cond.Type = ctype

	switch ctype {
	case ConditionIf, ConditionElseif, ConditionWhile:
		cond.Expr = expr
	case ConditionElse:
		cond.Expr = nil
	default:
		panic("invalid")
	}

	return cond
}

func (stmt *ConditionStatement) GetIdent() string {
	switch stmt.Type {
	case ConditionIf:
		return "if"
	case ConditionElseif:
		return "elseif"
	case ConditionElse:
		return "else"
	case ConditionWhile:
		return "while"
	default:
		panic("invalid")
	}
}
