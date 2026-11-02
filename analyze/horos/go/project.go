package analyze

import (
	"context"
	"fmt"

	bridge "github.com/Bitspark/bittype/internal/horos"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	codec "github.com/bitspark/ontos/codec/go"
)

type exporter struct {
	ctx context.Context
	input *bridge.Decoded
	facts semantics.Data
	modules map[decl.Digest]semantics.EntityID
	declarations map[decl.DeclarationID]semantics.EntityID
	parameters map[decl.ParameterID]semantics.EntityID
	limits semantics.Limits
	records, edges, bytes int
}

func (e *exporter) reserve(records, edges, bytes int) error {
	if err := e.ctx.Err(); err != nil { return err }
	if records < 0 || records > e.records || edges < 0 || edges > e.edges || bytes < 0 || bytes > e.bytes { return semantics.ErrLimit }
	e.records -= records
	e.edges -= edges
	e.bytes -= bytes
	return nil
}

func keySize(k semantics.Key) int { return len(k.Owner) + len(k.Name) }
func occurrenceSize(at semantics.Occurrence) int { return len(at.Site.Tree) + keySize(at.Role) }

func (e *exporter) entity(kind semantics.Key, name string, owner semantics.EntityID, site ast.NodeRef) (semantics.EntityID, error) {
	if site.Node == 0 { return 0, fmt.Errorf("Horos entity has no current syntax site") }
	if err := e.reserve(1, 1, keySize(kind)+len(name)+len(site.Tree)); err != nil { return 0, err }
	e.facts.Entities = append(e.facts.Entities, semantics.Entity{Kind: kind, Name: name, Owner: owner, Declarations: []ast.NodeRef{site}})
	return semantics.EntityID(len(e.facts.Entities)), nil
}

func (e *exporter) association(site ast.NodeRef, term semantics.TermID) error {
	at := semantics.Occurrence{Site: site, Role: ExpressionRole()}
	if err := e.reserve(1, 0, occurrenceSize(at)); err != nil { return err }
	e.facts.Expressions = append(e.facts.Expressions, semantics.Expression{At: at, Term: term})
	return nil
}

func (e *exporter) binding(site ast.NodeRef, target semantics.EntityID) error {
	at := semantics.Occurrence{Site: site, Role: TargetRole()}
	if err := e.reserve(1, 0, occurrenceSize(at)); err != nil { return err }
	e.facts.Bindings = append(e.facts.Bindings, semantics.Binding{At: at, Target: target})
	return nil
}

func (e *exporter) expression(expression *decl.Expression, site ast.NodeRef, owner semantics.EntityID, depth int) (semantics.TermID, error) {
	if depth > e.limits.MaxDepth { return 0, semantics.ErrLimit }
	if site.Node == 0 { return 0, fmt.Errorf("Horos expression has no current syntax site") }
	if err := e.reserve(1, 0, 0); err != nil { return 0, err }
	e.facts.Terms = append(e.facts.Terms, semantics.Term{})
	id := semantics.TermID(len(e.facts.Terms))
	term := semantics.Term{}
	switch expression.Kind() {
	case decl.Parameter:
		parameter, _ := expression.Parameter()
		term.Kind, term.Target = semantics.Reference, e.parameters[parameter]
		if err := e.binding(site, term.Target); err != nil { return 0, err }
	case decl.Application:
		target, _ := expression.Target()
		term.Kind, term.Target = semantics.Application, e.declarations[target]
		arguments := expression.Arguments()
		if err := e.reserve(0, len(arguments), 0); err != nil { return 0, err }
		term.Arguments = make([]semantics.TermID, len(arguments))
		for i, argument := range arguments {
			child, err := e.expression(argument, e.input.ArgumentSite(site, i), owner, depth+1)
			if err != nil { return 0, err }
			term.Arguments[i] = child
		}
		if err := e.binding(site, term.Target); err != nil { return 0, err }
		at := semantics.Occurrence{Site: site, Role: ExpressionRole()}
		if err := e.reserve(1, 0, keySize(ApplicationDependency())+occurrenceSize(at)); err != nil { return 0, err }
		e.facts.Dependencies = append(e.facts.Dependencies, semantics.Dependency{From: owner, To: term.Target, Relation: ApplicationDependency(), At: &at})
	case decl.Descriptor:
		term.Kind, term.Operator, term.Encoding = semantics.ProviderExpression, key("descriptor/"+expression.Constructor()), semantics.Key{Owner: "ontos", Name: "codec-v1"}
		// The value was bounded during Horos resolution. Each retained payload is
		// charged before publication; temporary encoding is bounded by input size.
		encoded := codec.Encode(expression.Value())
		operands := expression.Operands()
		if err := e.reserve(0, len(operands), len(encoded)+keySize(term.Operator)+keySize(term.Encoding)); err != nil { return 0, err }
		term.Data = encoded
		term.Operands = make([]semantics.Operand, len(operands))
		for i, operand := range operands {
			path := operand.Path()
			if err := e.reserve(0, len(path), 0); err != nil { return 0, err }
			locator := make([]uint32, len(path))
			for j, part := range path { locator[j] = uint32(part) }
			child, err := e.expression(operand.Expression(), e.input.OperandSite(site, path), owner, depth+1)
			if err != nil { return 0, err }
			term.Operands[i] = semantics.Operand{Locator: locator, Term: child}
		}
	default:
		return 0, fmt.Errorf("unsupported checked Horos expression kind")
	}
	e.facts.Terms[id-1] = term
	if err := e.association(site, id); err != nil { return 0, err }
	return id, nil
}

func project(ctx context.Context, input *bridge.Decoded, program *decl.Program, limits semantics.Limits) (semantics.Data, error) {
	e := exporter{ctx: ctx, input: input, limits: limits, records: limits.MaxRecords, edges: limits.MaxEdges, bytes: limits.MaxBytes,
		modules: make(map[decl.Digest]semantics.EntityID), declarations: make(map[decl.DeclarationID]semantics.EntityID), parameters: make(map[decl.ParameterID]semantics.EntityID)}
	modules := program.Modules()
	// Allocate identity records before translating any expression. Forward uses
	// and equal display names never select a target by spelling.
	for _, module := range modules {
		moduleID, err := e.entity(key("module"), module.Name(), 0, input.ModuleSite(module.Pin()))
		if err != nil { return semantics.Data{}, err }
		e.modules[module.Pin()] = moduleID
		for _, declaration := range module.Declarations() {
			id, err := e.entity(key("declaration"), declaration.ID().Name, moduleID, input.DeclarationSite(declaration.ID()))
			if err != nil { return semantics.Data{}, err }
			e.declarations[declaration.ID()] = id
			for index, name := range declaration.Parameters() {
				parameter := decl.ParameterID{Declaration: declaration.ID(), Index: index}
				parameterID, err := e.entity(key("parameter"), name, id, input.ParameterSite(parameter))
				if err != nil { return semantics.Data{}, err }
				e.parameters[parameter] = parameterID
				sort := key("type")
				if err := e.reserve(0, 1, keySize(sort)); err != nil { return semantics.Data{}, err }
				e.facts.Entities[id-1].Parameters = append(e.facts.Entities[id-1].Parameters, semantics.Parameter{Entity: parameterID, Sort: sort})
			}
		}
	}
	for _, module := range modules {
		for _, declaration := range module.Declarations() {
			id := e.declarations[declaration.ID()]
			body, err := e.expression(declaration.Expression(), input.BodySite(declaration.ID()), id, 1)
			if err != nil { return semantics.Data{}, err }
			e.facts.Entities[id-1].Body = body
		}
	}
	return e.facts, nil
}
