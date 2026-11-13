package parser

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", 0)

// DebugPrintAST prints an AST.
func DebugPrintAST(root *Node) {
	printNode(root, "")
}

func printNode(node *Node, prefix string) {
	if node == nil {
		return
	}

	label := nodeTypeToString(node.NodeType)
	if node.Value != nil {
		label = fmt.Sprintf("%s(%v)", label, *node.Value)
	}

	logger.Printf("%s%s\n", prefix+" ", label)
	printNode(node.Left, prefix+"-")
	printNode(node.Right, prefix+"-")
}
