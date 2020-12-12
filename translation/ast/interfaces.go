package ast

type NodeFlag = uint64

const (
	NFlagNone NodeFlag = 0
	NFlagRoot NodeFlag = 1 << iota
	NFlagBlock
	NFlagText
)

type Node interface {
	init(Node)
	initChildren()

	GetRoot() AST
	GetParent() Node
	GetFlags() uint64
	GetChildren() []Node

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
	GetC() Node      // Current node
	SetC(Node)       // Set current node
	StartNode() Node // Get the first real node
}
