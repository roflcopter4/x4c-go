package ast

type AstNode struct {
	root     AST
	parent   Node
	children []Node
	flags    uint64
}

type RootNode struct {
	XMLStatement
	current Node
}

/***************************************************************************************/

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

// Accessor method: return the root node
func (n *AstNode) Root() AST        { return n.root }
func (n *AstNode) SetRoot(root AST) { n.root = root }

// Accessor method: return the parent node
func (n *AstNode) Parent() Node          { return n.parent }
func (n *AstNode) SetParent(parent Node) { n.parent = parent }

// Accessor method: return the mask of flags
func (n *AstNode) Flags() uint64        { return n.flags }
func (n *AstNode) SetFlags(mask uint64) { n.flags = mask }
func (n *AstNode) AddFlags(mask uint64) { n.flags |= mask }

func (n *AstNode) Children() []Node { return n.children }
func (n *AstNode) initChildren()    { n.children = make([]Node, 0) }

// Returns true if the node contains all provided flags
func (n *AstNode) HasFlags(flags ...NodeFlag) bool {
	for _, flg := range flags {
		if (n.flags & flg) == 0 {
			return false
		}
	}
	return true
}

/***************************************************************************************/

func (parent *AstNode) init(child Node) {
	child.SetRoot(parent.Root())
	child.SetParent(parent)
	child.SetFlags(NFlagNone)
	child.initChildren()

	parent.AddChild(child)
}

func (n *AstNode) AddChild(child Node) {
	n.AddFlags(NFlagBlock)
	n.children = append(n.children, child)
}

func (root *RootNode) C() Node           { return root.current }
func (root *RootNode) SetC(current Node) { root.current = current }

/***************************************************************************************/

func astNodeSlice2Node(s []*AstNode) []Node {
	ret := make([]Node, len(s))
	for i, v := range s {
		ret[i] = v
	}
	return ret
}
