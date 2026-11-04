package representation

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// ASTSet is a complete collection, not a Source and not semantic evidence.
type ASTSet struct {
	Entry ast.NodeRef
	Trees []AST
}

func (s ASTSet) EntryAST() AST {
	for _, tree := range s.Trees {
		info, err := tree.Inspect()
		if err == nil && info.ID == s.Entry.Tree && s.Entry.Node > 0 && int(s.Entry.Node) <= info.Nodes {
			return tree
		}
	}
	return AST{}
}
func (s ASTSet) LanguageID() language.ID { return s.EntryAST().LanguageID() }
func (s ASTSet) Snapshot() (ASTSet, error) {
	if len(s.Trees) == 0 || len(s.Trees) > 256 {
		return ASTSet{}, ErrInvalidAST
	}
	seen := map[ast.TreeID]bool{}
	catalogue := map[ast.TreeID]ast.NodeID{}
	out := ASTSet{Entry: s.Entry, Trees: make([]AST, len(s.Trees))}
	for i, tree := range s.Trees {
		info, err := tree.Inspect()
		if err != nil {
			return ASTSet{}, err
		}
		if seen[info.ID] || tree.id != s.Trees[0].id || tree.profile != s.Trees[0].profile || tree.payload.marker() != s.Trees[0].payload.marker() {
			return ASTSet{}, ErrLanguageMismatch
		}
		seen[info.ID] = true
		catalogue[info.ID] = ast.NodeID(info.Nodes)
		owned, err := tree.Snapshot()
		if err != nil {
			return ASTSet{}, err
		}
		out.Trees[i] = owned
	}
	for _, tree := range out.Trees {
		if err := tree.CheckReferences(catalogue); err != nil {
			return ASTSet{}, err
		}
	}
	if out.EntryAST().LanguageID() == "" {
		return ASTSet{}, ErrInvalidAST
	}
	return out, nil
}
func PackSet[L any](definition Language, input ast.Set[L]) (ASTSet, error) {
	if len(input.Trees) == 0 || len(input.Trees) > 256 {
		return ASTSet{}, ErrInvalidAST
	}
	out := ASTSet{Entry: input.Entry, Trees: make([]AST, len(input.Trees))}
	for i, tree := range input.Trees {
		packed, err := Pack(definition, tree)
		if err != nil {
			return ASTSet{}, err
		}
		out.Trees[i] = packed
	}
	return out.Snapshot()
}
