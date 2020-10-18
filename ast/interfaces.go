package ast

type Node interface {
	Init(Node)

	Root() AST
	Parent() Node
	Flags() uint64
	Children() []Node

	SetRoot(AST)
	SetParent(Node)
	SetFlags(uint64)
	AddFlags(uint64)
	HasFlags(...NodeFlag) bool
	AddChild(Node)

	AddXMLStatement(string) *XMLStatement
	AddComment(string) *XMLComment
	AddTextNode(string) *XMLText
}

type AST interface {
	Node
	C() Node   // Current node
	SetC(Node) // Set current node
}
