package ast

const (
	NFlagNone NodeFlag = 0
	NFlagRoot NodeFlag = 1 << iota
	NFlagBlock
	NFlagText
)

type (
	NodeFlag = uint64
)

type XMLAttribute struct {
	Name string
	Val  *Expression
}

type XMLStatement struct {
	AstNode
	Name       string
	Attributes []XMLAttribute
}

type Expression struct {
	Raw string
}

/***************************************************************************************/

func (n *XMLStatement) Init(parent Node) {
	n.root = parent.Root()
	n.parent = parent
	n.flags = NFlagNone
	n.children = make([]Node, 0)

	parent.AddChild(n)
}

func (n *XMLStatement) AddChild(child Node) {
	n.AddFlags(NFlagBlock)
	n.children = append(n.children, child)
}

func (n *AstNode) AddXMLStatement(name string) *XMLStatement {
	child := new(XMLStatement)
	child.Init(n)
	child.Name = name
	child.Attributes = make([]XMLAttribute, 0)
	return child
}

func (stmt *XMLStatement) AddAttribute(name string, expr *Expression) {
	stmt.Attributes = append(stmt.Attributes, XMLAttribute{name, expr})
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

func (n *XMLComment) Init(parent Node) {
	n.root = parent.Root()
	n.parent = parent
	n.flags = NFlagText
	parent.AddChild(n)
}

func (n *AstNode) AddComment(text string) *XMLComment {
	child := new(XMLComment)
	child.Init(n)
	child.Text = text
	return child
}

type XMLText struct {
	AstNode
	Text string
}

func (n *XMLText) Init(parent Node) {
	n.root = parent.Root()
	n.parent = parent
	n.flags = NFlagText
	parent.AddChild(n)
}

func (n *AstNode) AddTextNode(text string) *XMLText {
	if text == "" {
		return nil
	}
	child := new(XMLText)
	child.Init(n)
	child.Text = text
	return child
}
