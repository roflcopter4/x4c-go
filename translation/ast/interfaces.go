package ast

type Node interface {
	init(Node)
	initChildren()

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
	AddConditionStatement(*Expression, int) *ConditionStatement
}

type AST interface {
	Node
	C() Node         // Current node
	SetC(Node)       // Set current node
	StartNode() Node // Get the first real node
}
