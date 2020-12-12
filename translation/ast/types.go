package ast

const (
	ConditionIf = 1 + iota
	ConditionElseif
	ConditionElse
	ConditionWhile
)

type ConditionStatement struct {
	AstNode
	Type int
	Expr *Expression
}

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
