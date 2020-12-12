package ast

func NewAst() AST {
	ret := new(RootNode)
	ret.root = ret
	ret.parent = ret
	ret.flags = NFlagRoot | NFlagBlock
	ret.current = ret
	ret.children = make([]Node, 0, 4)

	return ret
}

/***************************************************************************************/

type AstNode struct {
	root     AST
	parent   Node
	children []Node
	flags    uint64
}

func (n *AstNode) AddFlags(mask uint64)  { n.flags |= mask }
func (n *AstNode) SetFlags(mask uint64)  { n.flags = mask }
func (n *AstNode) SetParent(parent Node) { n.parent = parent }
func (n *AstNode) SetRoot(root AST)      { n.root = root }
func (n *AstNode) initChildren()         { n.children = make([]Node, 0) }
func (n *AstNode) GetChildren() []Node   { return n.children }
func (n *AstNode) GetFlags() uint64      { return n.flags }
func (n *AstNode) GetParent() Node       { return n.parent }
func (n *AstNode) GetRoot() AST          { return n.root }
func (n *AstNode) NumChildren() int      { return len(n.children) }

// Returns true if the node contains all provided flags
func (n *AstNode) HasFlags(flags ...NodeFlag) bool {
	for _, flg := range flags {
		if (n.flags & flg) == 0 {
			return false
		}
	}
	return true
}

func (parent *AstNode) init(child Node) {
	child.SetRoot(parent.GetRoot())
	child.SetParent(parent)
	child.SetFlags(NFlagNone)
	child.initChildren()

	parent.AddChild(child)
}

func (n *AstNode) AddChild(child Node) {
	n.AddFlags(NFlagBlock)
	n.children = append(n.children, child)
}

/***************************************************************************************/

type RootNode struct {
	XMLStatement
	current Node
}

func (r *RootNode) GetC() Node        { return r.current }
func (r *RootNode) SetC(current Node) { r.current = current }
func (r *RootNode) StartNode() Node {
	switch r.NumChildren() {
	case 0:
		return nil
	case 1:
		return r.GetChildren()[0]
	default:
		panic("Invalid root node!")
	}
}

/***************************************************************************************/

func astNodeSlice2Node(s []*AstNode) []Node {
	ret := make([]Node, len(s))
	for i, v := range s {
		ret[i] = v
	}
	return ret
}
