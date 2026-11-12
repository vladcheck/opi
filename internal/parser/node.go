package parser

import "fmt"

// NodeType is an enum that represents type of an AST Node (e.g. Expression, Program, Statement, etc.).
type NodeType = int

const (
	// NtProgram is source code of a single file as a whole.
	NtProgram NodeType = iota
	// NtExpression represents an expression - an rvalue that evaluates to a value (e.g. a literal or a function).
	NtExpression
	// NtNumber is a number literal.
	NtNumber
	// NtString is a string literal.
	NtString
	// NtBool is a bool literal (true or false).
	NtBool
)

// Node is a single element in the AST, holding zero or more child Nodes.
type Node struct {
	NodeType NodeType
	Left     *Node
	Right    *Node
	Value    *string
}

// NewNode creates a new AST Node.
func NewNode(nt NodeType) *Node {
	return &Node{
		NodeType: nt,
		Left:     nil,
		Right:    nil,
		Value:    nil,
	}
}

// NewNodeWithValue creates a new AST Node with a Value - used for literals and identifiers.
func NewNodeWithValue(nt NodeType, value *string) *Node {
	return &Node{
		NodeType: nt,
		Left:     nil,
		Right:    nil,
		Value:    value,
	}
}

func (n *Node) String() string {
	nt := nodeTypeToString(n.NodeType)

	if n.Value != nil {
		return fmt.Sprintf("%s(%v)", nt, *n.Value)
	}

	return nt
}

func nodeTypeToString(nt NodeType) string {
	m := map[NodeType]string{
		NtProgram:    "Program",
		NtExpression: "Expr",
		NtNumber:     "Number",
		NtString:     "String",
		NtBool:       "Bool",
	}

	s, ok := m[nt]
	if !ok {
		return "UNMATCHED (edit nodeTypeToString())"
	}

	return s
}
