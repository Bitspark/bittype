package horosmodel

import (
	"fmt"
	"math/big"

	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	codec "github.com/bitspark/ontos/codec/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// These terms describe Horos-owned facts. The generic model validates their
// storage and edges, never their operators, claims or proof rules. Empty payloads
// are codec-v1 empty atoms; input/child locators are zero-based ordinal paths.
func (e *exporter[L]) providerTerm(operator string, payload core.Value, children []semantics.TermID) (semantics.TermID, error) {
	if payload == nil {
		payload = core.NewAtom(nil)
	}
	term := semantics.Term{Kind: semantics.ProviderExpression, Operator: e.frontend.key(operator), Encoding: semantics.Key{Owner: "ontos", Name: "codec-v1"}, Data: codec.Encode(payload)}
	if err := e.reserve(1, 2*len(children), keySize(term.Operator)+keySize(term.Encoding)+len(term.Data)); err != nil {
		return 0, err
	}
	for i, child := range children {
		term.Operands = append(term.Operands, semantics.Operand{Locator: []uint32{uint32(i)}, Term: child})
	}
	e.facts.Terms = append(e.facts.Terms, term)
	return semantics.TermID(len(e.facts.Terms)), nil
}

func (e *exporter[L]) fact(subject semantics.Subject, relation semantics.Key, term semantics.TermID) error {
	if err := e.reserve(1, 2, keySize(relation)); err != nil {
		return err
	}
	e.facts.Facts = append(e.facts.Facts, semantics.Fact{Subject: subject, Relation: relation, Term: term})
	return nil
}

func (e *exporter[L]) kindTerm(kind decl.Kind, depth int) (semantics.TermID, error) {
	if depth > e.limits.MaxDepth {
		return 0, semantics.ErrLimit
	}
	operator := "kind/proper"
	var inputs []semantics.TermID
	if kind.IsConstructor() {
		operator = "kind/constructor"
		for _, input := range kind.Inputs() {
			term, err := e.kindTerm(input, depth+1)
			if err != nil {
				return 0, err
			}
			inputs = append(inputs, term)
		}
	}
	return e.providerTerm(operator, nil, inputs)
}

func (e *exporter[L]) interfaceTerm(value decl.Interface, depth int) (semantics.TermID, error) {
	if depth > e.limits.MaxDepth {
		return 0, semantics.ErrLimit
	}
	operator := "interface/proper"
	var inputs []semantics.TermID
	if value.IsConstructor() {
		operator = "interface/constructor"
		for _, input := range value.Inputs() {
			term, err := e.interfaceTerm(input, depth+1)
			if err != nil {
				return 0, err
			}
			inputs = append(inputs, term)
		}
	}
	return e.providerTerm(operator, core.NewAtom([]byte(value.Result())), inputs)
}

func (e *exporter[L]) summaryTerm(summary *decl.DataSummary, depth int) (semantics.TermID, error) {
	if depth > e.limits.MaxDepth {
		return 0, semantics.ErrLimit
	}
	if summary == nil {
		return 0, fmt.Errorf("missing checked Horos Data summary")
	}
	var children []semantics.TermID
	switch summary.Kind() {
	case decl.SummaryTrue, decl.SummaryFalse:
	case decl.SummaryNeed:
		parameter, ok := summary.Parameter()
		target := e.parameters[parameter]
		if !ok || target == 0 {
			return 0, fmt.Errorf("Horos summary has no checked parameter identity")
		}
		if err := e.reserve(1, 0, 0); err != nil {
			return 0, err
		}
		e.facts.Terms = append(e.facts.Terms, semantics.Term{Kind: semantics.Reference, Target: target})
		children = append(children, semantics.TermID(len(e.facts.Terms)))
	case decl.SummaryAnd:
		for _, child := range summary.Children() {
			term, err := e.summaryTerm(child, depth+1)
			if err != nil {
				return 0, err
			}
			children = append(children, term)
		}
	case decl.SummaryProbe:
		probe := summary.Probe()
		if probe == nil {
			return 0, fmt.Errorf("Horos summary has no checked probe expression")
		}
		// A probe retains checked expression structure, not another invocation of
		// its own DataSummary accessor. No fabricated syntax provenance or duplicate
		// dependency/binding occurrence is attached to this reduction evidence.
		term, err := e.expression(probe, ast.NodeRef{}, 0, depth+1, false)
		if err != nil {
			return 0, err
		}
		children = append(children, term)
	default:
		return 0, fmt.Errorf("unsupported checked Horos Data summary")
	}
	return e.providerTerm("data-summary/"+string(summary.Kind()), nil, children)
}

func (e *exporter[L]) typeFacts(subject semantics.Subject, value decl.Interface, kind decl.Kind, summary *decl.DataSummary, status decl.DataStatus, depth int) error {
	shape, err := e.kindTerm(kind, depth)
	if err != nil {
		return err
	}
	if err := e.fact(subject, e.frontend.KindRelation(), shape); err != nil {
		return err
	}
	contract, err := e.interfaceTerm(value, depth)
	if err != nil {
		return err
	}
	if err := e.fact(subject, e.frontend.InterfaceRelation(), contract); err != nil {
		return err
	}
	if summary == nil && status == "" {
		return nil // Constructor expressions and parameter entities have no result proof.
	}
	switch status {
	case decl.DataProven, decl.DataFailed, decl.DataConditional:
	default:
		return fmt.Errorf("unsupported checked Horos Data status")
	}
	proof, err := e.summaryTerm(summary, depth)
	if err != nil {
		return err
	}
	if err := e.fact(subject, e.frontend.DataSummaryRelation(), proof); err != nil {
		return err
	}
	classification, err := e.providerTerm("data-status/"+string(status), core.NewAtom([]byte(status)), nil)
	if err != nil {
		return err
	}
	return e.fact(subject, e.frontend.DataStatusRelation(), classification)
}

func (e *exporter[L]) identity(subject semantics.Subject, operator string, payload core.Value) error {
	term, err := e.providerTerm("identity/"+operator, payload, nil)
	if err != nil {
		return err
	}
	return e.fact(subject, e.frontend.IdentityRelation(), term)
}

func (e *exporter[L]) moduleIdentity(pin decl.Digest) error {
	return e.identity(semantics.Subject{Entity: e.modules[pin]}, "module", core.NewAtom(pin[:]))
}

func (e *exporter[L]) declarationFacts(declaration *decl.Declaration) error {
	id := declaration.ID()
	subject := semantics.Subject{Entity: e.declarations[id]}
	if err := e.identity(subject, "declaration", core.NewTuple(core.NewAtom(id.Module[:]), data.MustEncodeText(id.Name))); err != nil {
		return err
	}
	if err := e.typeFacts(subject, declaration.Interface(), declaration.Kind(), declaration.DataSummary(), declaration.DataStatus(), 1); err != nil {
		return err
	}
	for index, parameter := range declaration.Parameters() {
		parameterID := decl.ParameterID{Declaration: id, Index: index}
		subject := semantics.Subject{Entity: e.parameters[parameterID]}
		if err := e.identity(subject, "parameter", core.NewTuple(core.NewAtom(id.Module[:]), data.MustEncodeText(id.Name), data.EncodeInt(big.NewInt(int64(index))))); err != nil {
			return err
		}
		if err := e.typeFacts(subject, parameter.Interface, parameter.Interface.Kind(), nil, "", 1); err != nil {
			return err
		}
	}
	return nil
}

func (e *exporter[L]) dependency(site ast.NodeRef, owner, target semantics.EntityID) error {
	at := semantics.Occurrence{Site: site, Role: e.frontend.ExpressionRole()}
	if err := e.reserve(1, 0, keySize(e.frontend.ReferenceDependency())+occurrenceSize(at)); err != nil {
		return err
	}
	e.facts.Dependencies = append(e.facts.Dependencies, semantics.Dependency{From: owner, To: target, Relation: e.frontend.ReferenceDependency(), At: &at})
	return nil
}
