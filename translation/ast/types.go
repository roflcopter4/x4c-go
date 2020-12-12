package ast

type NodeFlag = uint64

const (
	NFlagNone NodeFlag = 0
	NFlagRoot NodeFlag = 1 << iota
	NFlagBlock
	NFlagText
)

/***************************************************************************************/

type XMLStatement struct {
	AstNode
	Name       string
	Attributes []*XMLAttribute
}

type XMLAttribute struct {
	Name string
	Val  *Expression
}

type Expression struct {
	Raw string
	XML []*XMLAttribute
}

func (n *AstNode) AddXMLStatement(name string) *XMLStatement {
	child := new(XMLStatement)
	n.init(child)

	child.Name = name
	child.Attributes = make([]*XMLAttribute, 0)
	return child
}

func (stmt *XMLStatement) AddAttribute(name string, expr *Expression) {
	stmt.Attributes = append(stmt.Attributes, &XMLAttribute{name, expr})
}

func NewExpression(raw string) *Expression {
	expr := new(Expression)
	expr.Raw = raw
	return expr
}

/***************************************************************************************/

type XMLComment struct {
	AstNode
	Text string
}

type XMLText struct {
	AstNode
	Text string
}

func (n *AstNode) AddComment(text string) *XMLComment {
	child := new(XMLComment)
	n.init(child)

	child.flags = NFlagText
	child.Text = text
	return child
}

func (n *AstNode) AddTextNode(text string) *XMLText {
	if text == "" {
		return nil
	}
	child := new(XMLText)
	n.init(child)

	child.flags = NFlagText
	child.Text = text
	return child
}

/***************************************************************************************/

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
