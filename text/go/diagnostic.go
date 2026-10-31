package text

import (
	"fmt"
	"github.com/Bitspark/bittype/ast/go"
)

type Code string

const (
	SyntaxError Code = "syntax"
	Unsupported Code = "unsupported"
	InvalidAST Code = "invalid-ast"
	InvalidCST Code = "invalid-cst"
	LimitExceeded Code = "limit"
)

// Diagnostic is returned through error. Span offsets are half-open UTF-8 bytes;
// Node is optional when no abstract node has been built.
type Diagnostic struct {
	Code Code
	Message string
	Span *ast.Span
	Node *ast.NodeRef
}

func (d *Diagnostic) Error() string {
	if d.Span != nil {
		return fmt.Sprintf("%s:%d:%d: %s: %s", d.Span.Source, d.Span.Start, d.Span.End, d.Code, d.Message)
	}
	if d.Node != nil {
		return fmt.Sprintf("%s#%d: %s: %s", d.Node.Tree, d.Node.Node, d.Code, d.Message)
	}
	return string(d.Code) + ": " + d.Message
}

func AtNode[L any](tree ast.AST[L], id ast.NodeID, message string) error {
	d := &Diagnostic{Code: InvalidAST, Message: message, Node: &ast.NodeRef{Tree: tree.ID, Node: id}}
	if id > 0 && int(id) <= len(tree.Nodes) && len(tree.Nodes[id-1].Origin.Spans) > 0 {
		span := tree.Nodes[id-1].Origin.Spans[0]
		d.Span = &span
	}
	return d
}

