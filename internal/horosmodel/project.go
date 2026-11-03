// Package horosmodel projects the checked Horos model and verifies its AST view.
// It is private to this frontend, never a dependency of generic infrastructure.
package horosmodel

import (
	"encoding/hex"
	"errors"
	"strconv"

	structure "github.com/Bitspark/bittype/check/ast/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	s "github.com/Bitspark/bittype/internal/declsyntax"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

var ErrInvalid = errors.New("invalid Horos AST")

type importID struct {
	module decl.Digest
	alias  string
}

type projector[L any] struct {
	tree         ast.AST[L]
	modules      map[decl.Digest]ast.NodeID
	declarations map[decl.DeclarationID]ast.NodeID
	parameters   map[decl.ParameterID]ast.NodeID
	imports      map[importID]ast.NodeID
	bytes        int
	entries      int
	err          error
}

func attributes(pairs ...string) []ast.Attribute {
	result := make([]ast.Attribute, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		result = append(result, ast.Attribute{Name: pairs[i], Text: pairs[i+1]})
	}
	return result
}

func (p *projector[L]) add(kind ast.Kind[L], attrs ...string) ast.NodeID {
	if p.err != nil {
		return 0
	}
	bytes := len(kind)
	for _, value := range attrs {
		bytes += len(value)
	}
	if len(p.tree.Nodes) >= structure.MaxNodes || !p.take(bytes, len(attrs)/2) {
		p.exhausted()
		return 0
	}
	p.tree.Nodes = append(p.tree.Nodes, ast.Node[L]{Kind: kind, Attributes: attributes(attrs...)})
	return ast.NodeID(len(p.tree.Nodes))
}

func (p *projector[L]) child(parent ast.NodeID, role string, child ast.NodeID) {
	if !p.take(len(role), 1) {
		return
	}
	p.tree.Nodes[parent-1].Children = append(p.tree.Nodes[parent-1].Children, ast.Child{Role: role, Node: child})
}

func (p *projector[L]) reference(node ast.NodeID, role string, target ast.NodeID) {
	if !p.take(len(role)+len(p.tree.ID), 1) {
		return
	}
	p.tree.Nodes[node-1].References = append(p.tree.Nodes[node-1].References,
		ast.Reference{Role: role, Target: ast.NodeRef{Tree: p.tree.ID, Node: target}})
}

func (p *projector[L]) exhausted() {
	if p.err == nil {
		p.err = &text.Diagnostic{Code: text.LimitExceeded, Message: "Horos AST projection exceeds resource limit"}
	}
}

func (p *projector[L]) take(bytes, entries int) bool {
	if p.err != nil {
		return false
	}
	if bytes > structure.MaxTextBytes-p.bytes || entries > structure.MaxEntries-p.entries {
		p.exhausted()
		return false
	}
	p.bytes += bytes
	p.entries += entries
	return true
}

func (p *projector[L]) depth(depth int) bool {
	if depth > structure.MaxDepth {
		p.exhausted()
	}
	return p.err == nil
}

// Project keeps module/import/declaration/binder identities in ordinary reference
// edges. The caller chooses the AST address; it is not a semantic identity.
func (f Frontend[L]) Project(program Program, treeID ast.TreeID) (ast.AST[L], error) {
	if program == nil || treeID == "" {
		return ast.AST[L]{}, ErrInvalid
	}
	root, ok := program.Module(program.Root())
	if !ok {
		return ast.AST[L]{}, ErrInvalid
	}
	// A Program may have been admitted with custom Horos limits. This fixed AST
	// profile must admit the same closure under the defaults used by Read.
	modules := program.Modules()
	if len(modules) > decl.DefaultLimits().MaxModules {
		return ast.AST[L]{}, &decl.Error{Code: decl.Exhausted, Subject: decl.SubjectRequest, Message: "module closure exceeds AST profile's Horos limits"}
	}
	dependencies := make([]core.Value, 0, len(modules)-1)
	for _, module := range modules {
		if module.Pin() != program.Root() {
			dependencies = append(dependencies, module.Value())
		}
	}
	return f.Resolve(root.Value(), dependencies, treeID)
}

// Resolve admits canonical values directly with the profile's Horos limits.
func (f Frontend[L]) Resolve(root core.Value, dependencies []core.Value, treeID ast.TreeID) (ast.AST[L], error) {
	program, err := f.ResolveModel(root, dependencies, decl.DefaultLimits())
	if err != nil {
		return ast.AST[L]{}, err
	}
	return projectChecked[L](program, treeID)
}

// Callers supply a fully checked Program. Projection has the fixed syntax
// profile's storage bounds independently of the resolver's selected limits.
func projectChecked[L any](program Program, treeID ast.TreeID) (ast.AST[L], error) {
	if treeID == "" {
		return ast.AST[L]{}, ErrInvalid
	}
	p := &projector[L]{
		tree:    ast.AST[L]{ID: treeID},
		modules: make(map[decl.Digest]ast.NodeID), declarations: make(map[decl.DeclarationID]ast.NodeID),
		parameters: make(map[decl.ParameterID]ast.NodeID), imports: make(map[importID]ast.NodeID),
	}
	if !p.take(len(treeID), 0) {
		return ast.AST[L]{}, p.err
	}
	p.tree.Root = p.add(s.KindProgram)
	modules := program.Modules()
	// Allocate identities before projecting any body, including forward uses.
	for _, module := range modules {
		m := p.add(s.KindModule, s.AttributeName, module.Name(), s.AttributePin, module.Pin().String())
		p.modules[module.Pin()] = m
		p.child(p.tree.Root, s.RoleModule, m)
		for _, imported := range module.Imports() {
			i := p.add(s.KindImport, s.AttributeAlias, imported.Alias(), s.AttributeName, imported.Name(), s.AttributePin, imported.Pin().String())
			p.imports[importID{module.Pin(), imported.Alias()}] = i
			p.child(m, s.RoleImport, i)
		}
		for _, declaration := range module.Declarations() {
			d := p.add(s.KindDeclaration, s.AttributeName, declaration.ID().Name)
			p.declarations[declaration.ID()] = d
			p.child(m, s.RoleDeclaration, d)
			for index, name := range declaration.Parameters() {
				parameter := p.add(s.KindTypeParameter, s.AttributeName, name, s.AttributeIndex, strconv.Itoa(index))
				p.parameters[decl.ParameterID{Declaration: declaration.ID(), Index: index}] = parameter
				p.child(d, s.RoleParameter, parameter)
			}
		}
		if p.err != nil {
			return ast.AST[L]{}, p.err
		}
	}
	p.reference(p.tree.Root, s.RoleRoot, p.modules[program.Root()])
	for _, module := range modules {
		for _, imported := range module.Imports() {
			p.reference(p.imports[importID{module.Pin(), imported.Alias()}], s.RoleTarget, p.modules[imported.Pin()])
		}
		for _, declaration := range module.Declarations() {
			body := p.expression(declaration.Expression(), module.Pin(), 3)
			p.child(p.declarations[declaration.ID()], s.RoleBody, body)
			if p.err != nil {
				return ast.AST[L]{}, p.err
			}
		}
	}
	if p.err != nil {
		return ast.AST[L]{}, p.err
	}
	if err := structure.Structure(p.tree); err != nil {
		return ast.AST[L]{}, err
	}
	return p.tree, nil
}

func (p *projector[L]) expression(expression *decl.Expression, module decl.Digest, depth int) ast.NodeID {
	if !p.depth(depth) {
		return 0
	}
	switch expression.Kind() {
	case decl.Parameter:
		binder, _ := expression.Parameter()
		n := p.add(s.KindParameter, s.AttributeIndex, strconv.Itoa(binder.Index))
		p.reference(n, s.RoleTarget, p.parameters[binder])
		return n
	case decl.Application:
		// Syntax remains authored; semantic target comes from the checked model.
		raw := expression.Value().(core.Tuple)
		target := raw.At(1).(core.Tuple)
		scope := "local"
		attrs := []string{s.AttributeScope, scope}
		if target.Len() == 3 {
			scope = "imported"
			alias, _ := data.ReadText(target.At(1))
			name, _ := data.ReadText(target.At(2))
			attrs = []string{s.AttributeScope, scope, s.AttributeAlias, alias, s.AttributeName, name}
		} else {
			name, _ := data.ReadText(target.At(1))
			attrs = append(attrs, s.AttributeName, name)
		}
		n := p.add(s.KindApplication, attrs...)
		resolved, _ := expression.Target()
		p.reference(n, s.RoleTarget, p.declarations[resolved])
		if scope == "imported" {
			alias, _ := data.ReadText(target.At(1))
			p.reference(n, s.RoleImport, p.imports[importID{module, alias}])
		}
		for _, argument := range expression.Arguments() {
			p.child(n, s.RoleArgument, p.expression(argument, module, depth+1))
			if p.err != nil {
				return 0
			}
		}
		return n
	default:
		n := p.add(s.KindDescriptor, s.AttributeConstructor, expression.Constructor())
		// Type positions are supplied by Horos; the frontend has no descriptor
		// constructor switch and never scans a literal for slot-shaped tuples.
		slots := &slot{}
		for _, operand := range expression.Operands() {
			position := slots
			for _, index := range operand.Path() {
				if position.children == nil {
					position.children = make(map[int]*slot)
				}
				if position.children[index] == nil {
					position.children[index] = &slot{}
				}
				position = position.children[index]
			}
			position.expression = operand.Expression()
		}
		raw := expression.Value().(core.Tuple)
		for i := 1; i < raw.Len(); i++ {
			p.child(n, s.RoleOperand, p.fragment(raw.At(i), slots.children[i], module, depth+1))
			if p.err != nil {
				return 0
			}
		}
		return n
	}
}

type slot struct {
	expression *decl.Expression
	children   map[int]*slot
}

func (p *projector[L]) fragment(value core.Value, position *slot, module decl.Digest, depth int) ast.NodeID {
	if !p.depth(depth) {
		return 0
	}
	if position != nil && position.expression != nil {
		return p.expression(position.expression, module, depth)
	}
	switch v := value.(type) {
	case core.Atom:
		// Check before copying the atom and expanding it into hexadecimal text.
		if v.Len() > (structure.MaxTextBytes-p.bytes-len(s.KindAtom)-len(s.AttributeHex))/2 {
			p.exhausted()
			return 0
		}
		return p.add(s.KindAtom, s.AttributeHex, hex.EncodeToString(v.Bytes()))
	default:
		tuple := v.(core.Tuple)
		n := p.add(s.KindTuple)
		for i := 0; i < tuple.Len(); i++ {
			var next *slot
			if position != nil {
				next = position.children[i]
			}
			p.child(n, s.RoleElement, p.fragment(tuple.At(i), next, module, depth+1))
			if p.err != nil {
				return 0
			}
		}
		return n
	}
}
