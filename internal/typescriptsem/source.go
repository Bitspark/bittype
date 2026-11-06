// Package typescriptsem owns the private TypeScript compiler bridge.
package typescriptsem

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	parser "github.com/Bitspark/bittype/parse/syntax/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type Node struct { ID ast.NodeID; Kind s.Kind; Start int; End int; Attributes []ast.Attribute; Children []ast.Child; References []ast.Reference }
type Unit struct { Tree ast.TreeID; Path files.Path; Text string; Nodes []Node }
type Query struct { Source uint32; Target uint32 }
type input struct { Version int; Units []Unit; Query *Query }

func prepare(ctx context.Context,trees []s.AST,source *files.Context) ([]Unit,error) {
	if source==nil { return nil,nil }
	if len(trees)>16 { return nil,semantic.ErrLimit }
	var units []Unit; totalNodes,totalBytes:=0,0
	for _,tree:=range trees {
		if err:=ctx.Err(); err!=nil { return nil,err }
		var link files.SyntaxLink; var file files.File
		for _,candidate:=range source.Links { if candidate.Tree==tree.ID { link=candidate; break } }
		for _,candidate:=range source.Files.Entries { if candidate.ID==link.File { file=candidate; break } }
		if link.Tree=="" || file.ID=="" { return nil,semantic.ErrInvalidRequest }
		if !utf8.ValidString(file.Text) { return nil,semantic.ErrInvalidRequest }
		if !strings.HasSuffix(string(file.Path),".ts") || strings.HasSuffix(string(file.Path),".d.ts") { return nil,fmt.Errorf("%w: captured files must be .ts source units",semantic.ErrInvalidRequest) }
		totalNodes+=len(tree.Nodes); totalBytes+=len(file.Text); if totalNodes>4096 || totalBytes>1<<20 { return nil,semantic.ErrLimit }
		// Reparse the supplied original text to establish correspondence. No
		// formatter or printer participates, and caller origins are not trusted.
		parsed,err:=parser.Parse(ctx,parse.Input{Text:file.Text,SourceID:link.Source,TreeID:link.Tree}); if err!=nil { return nil,fmt.Errorf("%w: captured source/AST mismatch: %v",semantic.ErrInvalidRequest,err) }
		expected:=parsed.ASTView(); if expected.Root!=tree.Root || len(expected.Nodes)!=len(tree.Nodes) { return nil,semantic.ErrInvalidRequest }
		for i,node:=range expected.Nodes {
			actual:=tree.Nodes[i]
			if node.Kind!=actual.Kind || !reflect.DeepEqual(node.Attributes,actual.Attributes) || !reflect.DeepEqual(node.Children,actual.Children) { return nil,fmt.Errorf("%w: captured source/AST mismatch at %s:%d",semantic.ErrInvalidRequest,tree.ID,i+1) }
		}
		// TypeScript offsets use UTF-16; Bittype source spans use UTF-8 bytes.
		offsets:=make([]int,len(file.Text)+1); count:=0
		for offset,r:=range file.Text { offsets[offset]=count; count++; if r>0xffff { count++ }; offsets[offset+utf8.RuneLen(r)]=count }
		unit:=Unit{Tree:tree.ID,Path:file.Path,Text:file.Text}
		for i,node:=range expected.Nodes {
			if len(node.Origin.Spans)!=1 { return nil,semantic.ErrInvalidRequest }; span:=node.Origin.Spans[0]
			unit.Nodes=append(unit.Nodes,Node{ID:ast.NodeID(i+1),Kind:node.Kind,Start:offsets[span.Start],End:offsets[span.End],Attributes:node.Attributes,Children:node.Children,References:tree.Nodes[i].References})
		}
		units=append(units,unit)
	}
	return units,nil
}
