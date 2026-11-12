package parser

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var logger = log.New(os.Stdout, "", 0)

// DebugPrintAST prints an AST.
func DebugPrintAST(root *Node) {
	queue := []*nodeWithDepth{newNodeWithDepth(root, 0)}

	for len(queue) > 0 {
		current := queue[0]

		if current.node.Left != nil {
			queue = append(queue, newNodeWithDepth(current.node.Left, current.Depth+1))
		}

		logger.Printf("%v\n", current)

		if current.node.Right != nil {
			queue = append(queue, newNodeWithDepth(current.node.Right, current.Depth+1))
		}

		// for GC kicking in immediately
		queue[0] = nil

		queue = queue[1:]
	}
}

type nodeWithDepth struct {
	node  *Node
	Depth int
}

func newNodeWithDepth(node *Node, depth int) *nodeWithDepth {
	return &nodeWithDepth{
		node:  node,
		Depth: depth,
	}
}

func (node *nodeWithDepth) String() string {
	prefix := strings.Repeat(" ", node.Depth)
	n := node.node
	nt := nodeTypeToString(n.NodeType)

	if n.Value != nil {
		return fmt.Sprintf("%s%s(%v)", prefix, nt, *n.Value)
	}

	return fmt.Sprintf("%s%s", prefix, nt)
}
