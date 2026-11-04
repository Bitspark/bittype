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
			kind := ast.Kind[L](s.KindTypeDefinition)
			if declaration.Form() == decl.KindDefinition { kind = s.KindKindDefinition }
			d := p.add(kind, s.AttributeName, declaration.ID().Name)
			p.declarations[declaration.ID()] = d
			p.child(m, s.RoleDeclaration, d)
			if declaration.Form() == decl.KindDefinition {
				for index, slot := range declaration.KindParameters() {
					parameter := p.add(s.KindKindParameter, s.AttributeName, slot.Name, s.AttributeIndex, strconv.Itoa(index))
					p.parameters[decl.ParameterID{Declaration: declaration.ID(), Index: index}] = parameter
					p.child(d, s.RoleParameter, parameter)
				}
			} else {
				for index, slot := range declaration.Parameters() {
					parameter := p.add(s.KindTypeParameter, s.AttributeName, slot.Name, s.AttributeIndex, strconv.Itoa(index))
					p.parameters[decl.ParameterID{Declaration: declaration.ID(), Index: index}] = parameter
					p.child(d, s.RoleParameter, parameter)
				}
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
			d := p.declarations[declaration.ID()]
			if declaration.Form() == decl.KindDefinition {
				for index, slot := range declaration.KindParameters() {
					parameter := p.parameters[decl.ParameterID{Declaration: declaration.ID(), Index: index}]
					p.child(parameter, s.RoleSort, p.sort(slot.Sort, 4))
				}
				p.child(d, s.RoleResult, p.sort(declaration.Sort(), 3))
				p.child(d, s.RoleBody, p.kind(declaration.KindBody(), module.Pin(), 3))
			} else {
				for index, slot := range declaration.Parameters() {
					parameter := p.parameters[decl.ParameterID{Declaration: declaration.ID(), Index: index}]
					p.child(parameter, s.RoleContract, p.kind(slot.Kind, module.Pin(), 4))
				}
				p.child(d, s.RoleResult, p.kind(declaration.Result(), module.Pin(), 3))
				p.child(d, s.RoleBody, p.expression(declaration.Expression(), module.Pin(), 3))
			}
			if p.err != nil { return ast.AST[L]{}, p.err }
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
	case decl.Reference:
		target, _ := expression.Target()
		return p.namedReference(s.KindReference, expression.Value(), target, module)
	case decl.Application:
		n := p.add(s.KindApplication)
		p.child(n, s.RoleCallee, p.expression(expression.Callee(), module, depth+1))
		for _, argument := range expression.Arguments() {
			p.child(n, s.RoleArgument, p.expression(argument, module, depth+1))
			if p.err != nil { return 0 }
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

func (p *projector[L]) namedReference(kind ast.Kind[L], value core.Value, target decl.DeclarationID, module decl.Digest) ast.NodeID {
	raw := value.(core.Tuple)
	attrs := []string{s.AttributeScope, "local"}
	if raw.Len() == 3 {
		alias, _ := data.ReadText(raw.At(1))
		name, _ := data.ReadText(raw.At(2))
		attrs = []string{s.AttributeScope, "imported", s.AttributeAlias, alias, s.AttributeName, name}
	} else {
		name, _ := data.ReadText(raw.At(1))
		attrs = append(attrs, s.AttributeName, name)
	}
	n := p.add(kind, attrs...)
	p.reference(n, s.RoleTarget, p.declarations[target])
	if raw.Len() == 3 {
		alias, _ := data.ReadText(raw.At(1))
		p.reference(n, s.RoleImport, p.imports[importID{module, alias}])
	}
	return n
}

func (p *projector[L]) sort(value decl.Sort, depth int) ast.NodeID {
	if !p.depth(depth) { return 0 }
	if value.IsProper() { return p.add(s.KindProperSort) }
	n := p.add(s.KindArrowSort)
	for _, input := range value.Inputs() { p.child(n, s.RoleInput, p.sort(input, depth+1)) }
	result, _ := value.Result()
	p.child(n, s.RoleResult, p.sort(result, depth+1))
	return n
}

func (p *projector[L]) kind(value *decl.Kind, module decl.Digest, depth int) ast.NodeID {
	if !p.depth(depth) { return 0 }
	var n ast.NodeID
	switch value.Form() {
	case decl.KindTop, decl.KindBottom:
		kind := ast.Kind[L](s.KindTop)
		if value.Form() == decl.KindBottom { kind = s.KindBottom }
		n = p.add(kind)
		p.child(n, s.RoleSort, p.sort(value.Sort(), depth+1))
	case decl.KindParameterReference:
		binder, _ := value.Parameter()
		n = p.add(s.KindKindParameterReference, s.AttributeIndex, strconv.Itoa(binder.Index))
		p.reference(n, s.RoleTarget, p.parameters[binder])
	case decl.KindReference:
		target, _ := value.Target()
		n = p.namedReference(s.KindKindReference, value.Value(), target, module)
	case decl.KindApplication:
		n = p.add(s.KindKindApplication)
		p.child(n, s.RoleCallee, p.kind(value.Operands()[0].Kind(), module, depth+1))
		for _, argument := range value.Inputs() { p.child(n, s.RoleArgument, p.kind(argument, module, depth+1)) }
	case decl.KindUnion, decl.KindIntersection:
		kind := ast.Kind[L](s.KindUnion)
		if value.Form() == decl.KindIntersection { kind = s.KindIntersection }
		n = p.add(kind)
		p.child(n, s.RoleSort, p.sort(value.Sort(), depth+1))
		for _, operand := range value.Inputs() { p.child(n, s.RoleOperand, p.kind(operand, module, depth+1)) }
	case decl.KindArrow:
		n = p.add(s.KindArrow)
		for _, input := range value.Inputs() { p.child(n, s.RoleInput, p.kind(input, module, depth+1)) }
		p.child(n, s.RoleResult, p.kind(value.Result(), module, depth+1))
	case decl.KindNode:
		n = p.add(s.KindNodePattern, s.AttributeConstructor, value.NodeForm())
		for _, selector := range value.Selectors() { p.child(n, s.RoleSelector, p.add(s.KindSelector, s.AttributeName, selector.Key, s.AttributeValue, selector.Value)) }
		p.child(n, s.RolePattern, p.sequence(value.Sequence(), module, depth+1))
	case decl.KindRecursive:
		n = p.add(s.KindRecursiveKind, s.AttributeIndex, strconv.Itoa(value.Entry()))
		for i, equation := range value.Equations() {
			wrapper := p.add(s.KindEquation, s.AttributeIndex, strconv.Itoa(i))
			p.child(n, s.RoleEquation, wrapper)
			p.child(wrapper, s.RoleBody, p.kind(equation, module, depth+2))
		}
	case decl.KindVariable:
		distance, slot, _ := value.Recursive()
		n = p.add(s.KindRecursiveVariable, s.AttributeDepth, strconv.Itoa(distance), s.AttributeIndex, strconv.Itoa(slot))
	}
	return n
}

func (p *projector[L]) sequence(value *decl.Sequence, module decl.Digest, depth int) ast.NodeID {
	if !p.depth(depth) { return 0 }
	kind := ast.Kind[L](s.KindSequenceEpsilon)
	switch value.Form() {
	case decl.SequenceItem: kind = s.KindSequenceItem
	case decl.SequenceConcat: kind = s.KindSequenceConcat
	case decl.SequenceAlternative: kind = s.KindSequenceAlternative
	case decl.SequenceStar: kind = s.KindSequenceStar
	}
	n := p.add(kind)
	if value.Element() != nil { p.child(n, s.RoleElement, p.kind(value.Element(), module, depth+1)) }
	for _, child := range value.Children() { p.child(n, s.RoleElement, p.sequence(child, module, depth+1)) }
	return n
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
