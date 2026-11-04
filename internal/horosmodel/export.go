package horosmodel

import (
	"context"
	"fmt"

	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	codec "github.com/bitspark/ontos/codec/go"
)

type exporter[L any] struct {
	frontend              Frontend[L]
	ctx                   context.Context
	input                 *Decoded[L]
	facts                 semantics.Data
	modules               map[decl.Digest]semantics.EntityID
	declarations          map[decl.DeclarationID]semantics.EntityID
	parameters            map[decl.ParameterID]semantics.EntityID
	limits                semantics.Limits
	records, edges, bytes int
}

func (e *exporter[L]) reserve(records, edges, bytes int) error {
	if err := e.ctx.Err(); err != nil {
		return err
	}
	if records < 0 || records > e.records || edges < 0 || edges > e.edges || bytes < 0 || bytes > e.bytes {
		return semantics.ErrLimit
	}
	e.records -= records
	e.edges -= edges
	e.bytes -= bytes
	return nil
}

func keySize(k semantics.Key) int                { return len(k.Owner) + len(k.Name) }
func occurrenceSize(at semantics.Occurrence) int { return len(at.Site.Tree) + keySize(at.Role) }

func (e *exporter[L]) entity(kind semantics.Key, name string, owner semantics.EntityID, site ast.NodeRef) (semantics.EntityID, error) {
	if site.Node == 0 {
		return 0, fmt.Errorf("Horos entity has no current syntax site")
	}
	if err := e.reserve(1, 1, keySize(kind)+len(name)+len(site.Tree)); err != nil {
		return 0, err
	}
	e.facts.Entities = append(e.facts.Entities, semantics.Entity{Kind: kind, Name: name, Owner: owner, Declarations: []ast.NodeRef{site}})
	return semantics.EntityID(len(e.facts.Entities)), nil
}

func (e *exporter[L]) association(site ast.NodeRef, term semantics.TermID) error {
	at := semantics.Occurrence{Site: site, Role: e.frontend.ExpressionRole()}
	if err := e.reserve(1, 0, occurrenceSize(at)); err != nil {
		return err
	}
	e.facts.Expressions = append(e.facts.Expressions, semantics.Expression{At: at, Term: term})
	return nil
}

func (e *exporter[L]) binding(site ast.NodeRef, target semantics.EntityID) error {
	at := semantics.Occurrence{Site: site, Role: e.frontend.TargetRole()}
	if err := e.reserve(1, 0, occurrenceSize(at)); err != nil {
		return err
	}
	e.facts.Bindings = append(e.facts.Bindings, semantics.Binding{At: at, Target: target})
	return nil
}

func (e *exporter[L]) expression(expression *decl.Expression, site ast.NodeRef, owner semantics.EntityID, depth int, publish bool) (semantics.TermID, error) {
	if depth > e.limits.MaxDepth {
		return 0, semantics.ErrLimit
	}
	if publish && site.Node == 0 {
		return 0, fmt.Errorf("Horos expression has no current syntax site")
	}
	if err := e.reserve(1, 0, 0); err != nil {
		return 0, err
	}
	e.facts.Terms = append(e.facts.Terms, semantics.Term{})
	id := semantics.TermID(len(e.facts.Terms))
	term := semantics.Term{}
	switch expression.Kind() {
	case decl.Parameter:
		parameter, _ := expression.Parameter()
		term.Kind, term.Target = semantics.Reference, e.parameters[parameter]
		if publish {
			if err := e.binding(site, term.Target); err != nil {
				return 0, err
			}
		}
	case decl.Reference:
		target, _ := expression.Target()
		term.Kind, term.Target = semantics.Reference, e.declarations[target]
		if publish {
			if err := e.binding(site, term.Target); err != nil {
				return 0, err
			}
			if err := e.dependency(site, owner, term.Target); err != nil {
				return 0, err
			}
		}
	case decl.Application:
		term.Kind = semantics.Application
		callee, err := e.expression(expression.Callee(), e.input.CalleeSite(site), owner, depth+1, publish)
		if err != nil {
			return 0, err
		}
		term.Callee = callee
		arguments := expression.Arguments()
		if err := e.reserve(0, 1+len(arguments), 0); err != nil {
			return 0, err
		}
		term.Arguments = make([]semantics.TermID, len(arguments))
		for i, argument := range arguments {
			child, err := e.expression(argument, e.input.ArgumentSite(site, i), owner, depth+1, publish)
			if err != nil {
				return 0, err
			}
			term.Arguments[i] = child
		}
	case decl.Descriptor:
		term.Kind, term.Operator, term.Encoding = semantics.ProviderExpression, e.frontend.key("descriptor/"+expression.Constructor()), semantics.Key{Owner: "ontos", Name: "codec-v1"}
		// The value was bounded during Horos resolution. Each retained payload is
		// charged before publication; temporary encoding is bounded by input size.
		encoded := codec.Encode(expression.Value())
		operands := expression.Operands()
		if err := e.reserve(0, len(operands), len(encoded)+keySize(term.Operator)+keySize(term.Encoding)); err != nil {
			return 0, err
		}
		term.Data = encoded
		term.Operands = make([]semantics.Operand, len(operands))
		for i, operand := range operands {
			path := operand.Path()
			if err := e.reserve(0, len(path), 0); err != nil {
				return 0, err
			}
			locator := make([]uint32, len(path))
			for j, part := range path {
				locator[j] = uint32(part)
			}
			child, err := e.expression(operand.Expression(), e.input.OperandSite(site, path), owner, depth+1, publish)
			if err != nil {
				return 0, err
			}
			term.Operands[i] = semantics.Operand{Locator: locator, Term: child}
		}
	default:
		return 0, fmt.Errorf("unsupported checked Horos expression kind")
	}
	e.facts.Terms[id-1] = term
	if publish {
		if err := e.association(site, id); err != nil {
			return 0, err
		}
		if err := e.typeFacts(semantics.Subject{Term: id}, expression.Interface(), expression.KindShape(), expression.DataSummary(), expression.DataStatus(), depth); err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (f Frontend[L]) project(ctx context.Context, input *Decoded[L], program Program, limits semantics.Limits) (semantics.Data, error) {
	e := exporter[L]{frontend: f, ctx: ctx, input: input, limits: limits, records: limits.MaxRecords, edges: limits.MaxEdges, bytes: limits.MaxBytes,
		modules: make(map[decl.Digest]semantics.EntityID), declarations: make(map[decl.DeclarationID]semantics.EntityID), parameters: make(map[decl.ParameterID]semantics.EntityID)}
	modules := program.Modules()
	// Allocate identity records before translating any expression. Forward uses
	// and equal display names never select a target by spelling.
	for _, module := range modules {
		moduleID, err := e.entity(f.key("module"), module.Name(), 0, input.ModuleSite(module.Pin()))
		if err != nil {
			return semantics.Data{}, err
		}
		e.modules[module.Pin()] = moduleID
		for _, declaration := range module.Declarations() {
			id, err := e.entity(f.key("declaration"), declaration.ID().Name, moduleID, input.DeclarationSite(declaration.ID()))
			if err != nil {
				return semantics.Data{}, err
			}
			e.declarations[declaration.ID()] = id
			for index, slot := range declaration.Parameters() {
				parameter := decl.ParameterID{Declaration: declaration.ID(), Index: index}
				parameterID, err := e.entity(f.key("parameter"), slot.Name, id, input.ParameterSite(parameter))
				if err != nil {
					return semantics.Data{}, err
				}
				e.parameters[parameter] = parameterID
				sort := f.key("type")
				if slot.Interface.IsConstructor() {
					sort = f.key("constructor")
				}
				if err := e.reserve(0, 1, keySize(sort)); err != nil {
					return semantics.Data{}, err
				}
				e.facts.Entities[id-1].Parameters = append(e.facts.Entities[id-1].Parameters, semantics.Parameter{Entity: parameterID, Sort: sort})
			}
		}
	}
	for _, module := range modules {
		if err := e.moduleIdentity(module.Pin()); err != nil {
			return semantics.Data{}, err
		}
		for _, declaration := range module.Declarations() {
			id := e.declarations[declaration.ID()]
			body, err := e.expression(declaration.Expression(), input.BodySite(declaration.ID()), id, 1, true)
			if err != nil {
				return semantics.Data{}, err
			}
			e.facts.Entities[id-1].Body = body
			if err := e.declarationFacts(declaration); err != nil {
				return semantics.Data{}, err
			}
		}
	}
	return e.facts, nil
}
