package newast

/***************************************************************************************/

type BaseNode struct {
	parent   Node
	children []Node
	raw      string
}

func NewTree() Node {
	ret := new(BaseNode)
	return ret
}

func (n *BaseNode) SetParent(parent Node) { n.parent = parent }
func (n *BaseNode) GetChildren() []Node   { return n.children }
func (n *BaseNode) GetParent() Node       { return n.parent }

func (n *BaseNode) GetNChildren() int {
	if n.children == nil {
		return 0
	}
	return len(n.children)
}

func (n *BaseNode) AddChild(child Node) {
	if n.children == nil {
		n.children = make([]Node, 0, 4)
	}
	n.children = append(n.children, child)
}

/***************************************************************************************/

type ExpressionNode struct {
	BaseNode
}
