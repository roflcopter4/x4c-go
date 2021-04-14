package newast

/***************************************************************************************/

type XMLStatementNode struct {
	BaseNode
	Name       string
	Attributes []*XMLAttributeNode
}

type XMLAttributeNode struct {
	Name string
	Val  string
}

func NewXMLStatementNode(name string) *XMLStatementNode {
	child := new(XMLStatementNode)
	child.Name = name
	child.Attributes = make([]*XMLAttributeNode, 0)
	return child
}

func NewXMLAttributeNode(name string, expr string) *XMLAttributeNode {
	ret := new(XMLAttributeNode)
	*ret = XMLAttributeNode{Name: name, Val: expr}
	// ret.Name = name
	// ret.Val = expr
	return ret
}

func (stmt *XMLStatementNode) AddAttribute(name string, expr string) {
	stmt.Attributes = append(stmt.Attributes, NewXMLAttributeNode(name, expr))
}

/***************************************************************************************/

type XMLCommentNode struct {
	BaseNode
	Text    string
	isblock bool
}

type XMLTextNode struct {
	BaseNode
	Text string
}

func NewCommentNode(text string) *XMLCommentNode {
	ret := new(XMLCommentNode)
	ret.raw = text
	ret.Text = text
	return ret
}

func (n *BaseNode) AddTextNode(text string) *XMLTextNode {
	if text == "" {
		return nil
	}
	child := new(XMLTextNode)
	child.Text = text
	return child
}
