package newast

type Node interface {
	GetParent() Node
	GetChildren() []Node
	GetNChildren() int

	SetParent(Node)
	AddChild(Node)
}
