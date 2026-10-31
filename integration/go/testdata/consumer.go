package consumer

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/pipeline/go"
)

// These named markers are owned by a consumer unknown to both foundations.
type inputLanguage struct{}
type outputLanguage struct{}

const declaration ast.Kind[inputLanguage] = "consumer:declaration"
const reference ast.Kind[inputLanguage] = "consumer:reference"

var input = ast.AST[inputLanguage]{
	ID:   "input",
	Root: 1,
	Nodes: []ast.Node[inputLanguage]{
		{
			Kind:       declaration,
			Attributes: []ast.Attribute{{Name: "name", Text: "Example"}},
			Children:   []ast.Child{{Role: "member", Node: 2}},
		},
		{
			Kind: reference,
			References: []ast.Reference{
				{Role: "declaration", Target: ast.NodeRef{Tree: "imported", Node: 1}},
			},
			Origin: ast.Origin{
				Spans: []ast.Span{{Source: "input.json", Start: 0, End: 4}},
			},
		},
	},
}

var output = ast.AST[outputLanguage]{
	ID:   "output",
	Root: 1,
	Nodes: []ast.Node[outputLanguage]{
		{
			Kind:   "consumer:output",
			Origin: ast.Origin{Nodes: []ast.NodeRef{{Tree: input.ID, Node: input.Root}}},
		},
	},
}

// Assignability witnesses only: these variables are not invoked.
// A function with the exact signature can satisfy the stage contract.
var generate func(ast.AST[inputLanguage]) (ast.AST[outputLanguage], error)
var generator pipeline.Stage[ast.AST[inputLanguage], ast.AST[outputLanguage]] = generate

var format func(ast.AST[outputLanguage]) (string, error)
var formatter pipeline.Stage[ast.AST[outputLanguage], string] = format

// The generic pipeline is also usable without ASTs or language markers.
var parse func([]byte) (int, error)
var unrelated pipeline.Stage[[]byte, int] = parse
