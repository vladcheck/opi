package parser

// NodeType is an enum that represents type of an AST Node (e.g. Expression, Program, Statement, etc.).
type NodeType = int

const (
	// NTProgram is source code of a single file as a whole.
	NTProgram NodeType = iota
)

// Node is a single element in the AST, holding zero or more child Nodes.
type Node struct {
	Type NodeType
	// Children holds the direct children of the Node.
	Children []Node
}

// NewNode creates a new AST Node.
func NewNode(t NodeType) Node {
	return Node{
		Type:     t,
		Children: make([]Node, 0),
	}
}
