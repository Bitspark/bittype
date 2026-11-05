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
func (e *exporter[L]) entity(kind semantics.Key, name string, owner semantics.EntityID, site ast.NodeRef, standard bool) (semantics.EntityID, error) {
	if site.Node == 0 && !standard {
		return 0, fmt.Errorf("Horos entity has no current syntax site")
	}
	edges, bytes := 0, keySize(kind)+len(name)
	var sites []ast.NodeRef
	if site.Node != 0 {
		sites = []ast.NodeRef{site}
		edges = 1
		bytes += len(site.Tree)
	}
	if err := e.reserve(1, edges, bytes); err != nil {
		return 0, err
	}
	e.facts.Entities = append(e.facts.Entities, semantics.Entity{Kind: kind, Name: name, Owner: owner, Declarations: sites})
	return semantics.EntityID(len(e.facts.Entities)), nil
}
func (e *exporter[L]) association(site ast.NodeRef, term semantics.TermID) error {
	if site.Node == 0 {
		return nil
	}
	at := semantics.Occurrence{Site: site, Role: e.frontend.ExpressionRole()}
	if err := e.reserve(1, 0, occurrenceSize(at)); err != nil {
		return err
	}
	e.facts.Expressions = append(e.facts.Expressions, semantics.Expression{At: at, Term: term})
	return nil
}
func (e *exporter[L]) binding(site ast.NodeRef, target semantics.EntityID) error {
	if site.Node == 0 {
		return nil
	}
	at := semantics.Occurrence{Site: site, Role: e.frontend.TargetRole()}
	if err := e.reserve(1, 0, occurrenceSize(at)); err != nil {
		return err
	}
	e.facts.Bindings = append(e.facts.Bindings, semantics.Binding{At: at, Target: target})
	return nil
}
func (e *exporter[L]) dependency(owner, target semantics.EntityID, site ast.NodeRef) error {
	if owner == 0 || target == 0 {
		return fmt.Errorf("Horos reference has no declaration identity")
	}
	dependency := semantics.Dependency{From: owner, To: target, Relation: e.frontend.ReferenceDependency()}
	bytes := keySize(dependency.Relation)
	if site.Node != 0 {
		at := semantics.Occurrence{Site: site, Role: e.frontend.ExpressionRole()}
		dependency.At = &at
		bytes += occurrenceSize(at)
	}
	if err := e.reserve(1, 0, bytes); err != nil {
		return err
	}
	e.facts.Dependencies = append(e.facts.Dependencies, dependency)
	return nil
}
func (e *exporter[L]) allocateTerm(depth int) (semantics.TermID, error) {
	if depth > e.limits.MaxDepth {
		return 0, semantics.ErrLimit
	}
	if err := e.reserve(1, 0, 0); err != nil {
		return 0, err
	}
	e.facts.Terms = append(e.facts.Terms, semantics.Term{})
	return semantics.TermID(len(e.facts.Terms)), nil
}
func (e *exporter[L]) fillTerm(id semantics.TermID, term semantics.Term) error {
	edges := len(term.Arguments) + len(term.Operands)
	if term.Callee != 0 {
		edges++
	}
	for _, operand := range term.Operands {
		edges += len(operand.Locator)
	}
	if err := e.reserve(0, edges, keySize(term.Operator)+keySize(term.Encoding)+len(term.Data)); err != nil {
		return err
	}
	e.facts.Terms[id-1] = term
	return nil
}
func locator(path []int) []uint32 {
	result := make([]uint32, len(path))
	for i, part := range path {
		result[i] = uint32(part)
	}
	return result
}
func (e *exporter[L]) expression(value *decl.Expression, site ast.NodeRef, owner semantics.EntityID, depth int) (semantics.TermID, error) {
	id, err := e.allocateTerm(depth)
	if err != nil {
		return 0, err
	}
	term := semantics.Term{}
	switch value.Kind() {
	case decl.Parameter:
		parameter, _ := value.Parameter()
		term.Kind, term.Target = semantics.Reference, e.parameters[parameter]
		if err := e.binding(site, term.Target); err != nil {
			return 0, err
		}
	case decl.Reference:
		target, _ := value.Target()
		term.Kind, term.Target = semantics.Reference, e.declarations[target]
		if err := e.binding(site, term.Target); err != nil {
			return 0, err
		}
		if err := e.dependency(owner, term.Target, site); err != nil {
			return 0, err
		}
	case decl.Application:
		term.Kind = semantics.Application
		term.Callee, err = e.expression(value.Callee(), e.input.CalleeSite(site), owner, depth+1)
		if err != nil {
			return 0, err
		}
		for i, argument := range value.Arguments() {
			child, err := e.expression(argument, e.input.ArgumentSite(site, i), owner, depth+1)
			if err != nil {
				return 0, err
			}
			term.Arguments = append(term.Arguments, child)
		}
	case decl.Descriptor:
		term.Kind, term.Operator, term.Encoding = semantics.ProviderExpression, e.frontend.key("descriptor/"+value.Constructor()), semantics.Key{Owner: "ontos", Name: "codec-v1"}
		term.Data = codec.Encode(value.Value())
		for _, operand := range value.Operands() {
			path := operand.Path()
			child, err := e.expression(operand.Expression(), e.input.OperandSite(site, path), owner, depth+1)
			if err != nil {
				return 0, err
			}
			term.Operands = append(term.Operands, semantics.Operand{Locator: locator(path), Term: child})
		}
	default:
		return 0, fmt.Errorf("unsupported checked Horos expression")
	}
	if err := e.fillTerm(id, term); err != nil {
		return 0, err
	}
	if err := e.association(site, id); err != nil {
		return 0, err
	}
	subject := semantics.Subject{Term: id}
	if err := e.sortFact(subject, e.frontend.SortRelation(), value.Sort(), depth+1); err != nil {
		return 0, err
	}
	if err := e.memberships(subject, value.Memberships(), depth+1); err != nil {
		return 0, err
	}
	return id, nil
}

// Kind payloads and designated kind operands come from Horos. Recursive lexical
// variables remain finite syntax, so no recursive expansion enters the fact graph.
func (e *exporter[L]) kind(value *decl.Kind, site ast.NodeRef, owner semantics.EntityID, publish bool, depth int) (semantics.TermID, error) {
	if value == nil {
		return 0, fmt.Errorf("missing checked Horos kind")
	}
	id, err := e.allocateTerm(depth)
	if err != nil {
		return 0, err
	}
	term := semantics.Term{}
	switch value.Form() {
	case decl.KindParameterReference:
		parameter, _ := value.Parameter()
		term.Kind, term.Target = semantics.Reference, e.parameters[parameter]
		if publish {
			if err := e.binding(site, term.Target); err != nil {
				return 0, err
			}
		}
	case decl.KindReference:
		target, _ := value.Target()
		term.Kind, term.Target = semantics.Reference, e.declarations[target]
		if publish {
			if err := e.binding(site, term.Target); err != nil {
				return 0, err
			}
			if err := e.dependency(owner, term.Target, site); err != nil {
				return 0, err
			}
		}
	case decl.KindApplication:
		term.Kind = semantics.Application
		term.Callee, err = e.kind(value.Operands()[0].Kind(), e.input.CalleeSite(site), owner, publish, depth+1)
		if err != nil {
			return 0, err
		}
		for i, argument := range value.Inputs() {
			child, err := e.kind(argument, e.input.ArgumentSite(site, i), owner, publish, depth+1)
			if err != nil {
				return 0, err
			}
			term.Arguments = append(term.Arguments, child)
		}
	default:
		term.Kind, term.Operator, term.Encoding = semantics.ProviderExpression, e.frontend.key("kind/"+string(value.Form())), semantics.Key{Owner: "ontos", Name: "codec-v1"}
		term.Data = codec.Encode(value.Value())
		for _, operand := range value.Operands() {
			path := operand.Path()
			child, err := e.kind(operand.Kind(), e.input.OperandSite(site, path), owner, publish, depth+1)
			if err != nil {
				return 0, err
			}
			term.Operands = append(term.Operands, semantics.Operand{Locator: locator(path), Term: child})
		}
	}
	if err := e.fillTerm(id, term); err != nil {
		return 0, err
	}
	if publish {
		if err := e.association(site, id); err != nil {
			return 0, err
		}
		if err := e.sortFact(semantics.Subject{Term: id}, e.frontend.SortRelation(), value.Sort(), depth+1); err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (f Frontend[L]) project(ctx context.Context, input *Decoded[L], program Program, limits semantics.Limits) (semantics.Data, error) {
	e := exporter[L]{frontend: f, ctx: ctx, input: input, limits: limits, records: limits.MaxRecords, edges: limits.MaxEdges, bytes: limits.MaxBytes,
		modules: make(map[decl.Digest]semantics.EntityID), declarations: make(map[decl.DeclarationID]semantics.EntityID), parameters: make(map[decl.ParameterID]semantics.EntityID)}
	modules := program.Modules()
	// Allocate every identity before any expression, including source-less bundled
	// definitions. Only actual snapshot sites enter Declarations and occurrences.
	for _, module := range modules {
		standard := module.Pin() == decl.PreludePin()
		moduleID, err := e.entity(f.key("module"), module.Name(), 0, input.ModuleSite(module.Pin()), standard)
		if err != nil {
			return semantics.Data{}, err
		}
		e.modules[module.Pin()] = moduleID
		if err := e.moduleIdentity(moduleID, module.Pin()); err != nil {
			return semantics.Data{}, err
		}
		for _, declaration := range module.Declarations() {
			kind, parameterKind := "type-definition", "type-parameter"
			if declaration.Form() == decl.KindDefinition {
				kind, parameterKind = "kind-definition", "kind-parameter"
			}
			id, err := e.entity(f.key(kind), declaration.ID().Name, moduleID, input.DeclarationSite(declaration.ID()), standard)
			if err != nil {
				return semantics.Data{}, err
			}
			e.declarations[declaration.ID()] = id
			if err := e.declarationIdentity(id, declaration.ID()); err != nil {
				return semantics.Data{}, err
			}
			names := []string{}
			sorts := []decl.Sort{}
			if declaration.Form() == decl.KindDefinition {
				for _, parameter := range declaration.KindParameters() {
					names = append(names, parameter.Name)
					sorts = append(sorts, parameter.Sort)
				}
			} else {
				for _, parameter := range declaration.Parameters() {
					names = append(names, parameter.Name)
					sorts = append(sorts, parameter.Kind.Sort())
				}
			}
			for index, name := range names {
				parameter := decl.ParameterID{Declaration: declaration.ID(), Index: index}
				parameterID, err := e.entity(f.key(parameterKind), name, id, input.ParameterSite(parameter), standard)
				if err != nil {
					return semantics.Data{}, err
				}
				e.parameters[parameter] = parameterID
				coarse := f.key("proper")
				if !sorts[index].IsProper() {
					coarse = f.key("constructor")
				}
				if err := e.reserve(0, 1, keySize(coarse)); err != nil {
					return semantics.Data{}, err
				}
				e.facts.Entities[id-1].Parameters = append(e.facts.Entities[id-1].Parameters, semantics.Parameter{Entity: parameterID, Sort: coarse})
				if err := e.parameterIdentity(parameterID, parameter); err != nil {
					return semantics.Data{}, err
				}
				if err := e.sortFact(semantics.Subject{Entity: parameterID}, f.SortRelation(), sorts[index], 1); err != nil {
					return semantics.Data{}, err
				}
			}
		}
	}
	for _, module := range modules {
		for _, declaration := range module.Declarations() {
			id := e.declarations[declaration.ID()]
			subject := semantics.Subject{Entity: id}
			if err := e.sortFact(subject, f.SortRelation(), declaration.Sort(), 1); err != nil {
				return semantics.Data{}, err
			}
			var body semantics.TermID
			var err error
			if declaration.Form() == decl.KindDefinition {
				if err := e.sortFact(subject, f.ResultSortRelation(), declaration.Sort(), 1); err != nil {
					return semantics.Data{}, err
				}
				body, err = e.kind(declaration.KindBody(), input.BodySite(declaration.ID()), id, true, 1)
			} else {
				if err := e.sortFact(subject, f.ResultSortRelation(), declaration.Result().Sort(), 1); err != nil {
					return semantics.Data{}, err
				}
				for index, parameter := range declaration.Parameters() {
					pid := decl.ParameterID{Declaration: declaration.ID(), Index: index}
					kind, err := e.kind(parameter.Kind, input.ParameterContractSite(pid), id, true, 1)
					if err != nil {
						return semantics.Data{}, err
					}
					if err := e.fact(semantics.Subject{Entity: e.parameters[pid]}, f.DeclaredKindRelation(), kind); err != nil {
						return semantics.Data{}, err
					}
				}
				result, err := e.kind(declaration.Result(), input.ResultSite(declaration.ID()), id, true, 1)
				if err != nil {
					return semantics.Data{}, err
				}
				if err := e.fact(subject, f.ResultKindRelation(), result); err != nil {
					return semantics.Data{}, err
				}
				contract, err := e.kind(declaration.Contract(), ast.NodeRef{}, 0, false, 1)
				if err != nil {
					return semantics.Data{}, err
				}
				if err := e.fact(subject, f.DeclaredKindRelation(), contract); err != nil {
					return semantics.Data{}, err
				}
				if err := e.memberships(subject, declaration.Memberships(), 1); err != nil {
					return semantics.Data{}, err
				}
				body, err = e.expression(declaration.Expression(), input.BodySite(declaration.ID()), id, 1)
			}
			if err != nil {
				return semantics.Data{}, err
			}
			e.facts.Entities[id-1].Body = body
			if err := e.fact(subject, f.SymbolicBodyRelation(), body); err != nil {
				return semantics.Data{}, err
			}
		}
	}
	return e.facts, nil
}
