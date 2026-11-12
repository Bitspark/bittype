package horosmodel

import (
	"math/big"

	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	codec "github.com/bitspark/ontos/codec/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

func (e *exporter[L]) fact(subject semantics.Subject, relation semantics.Key, value semantics.TermID) error {
	if err := e.reserve(1, 2, keySize(relation)); err != nil {
		return err
	}
	e.facts.Facts = append(e.facts.Facts, semantics.Fact{Subject: subject, Relation: relation, Term: value})
	return nil
}
func (e *exporter[L]) provider(name string, value core.Value, operands []semantics.Operand, depth int) (semantics.TermID, error) {
	id, err := e.allocateTerm(depth)
	if err != nil {
		return 0, err
	}
	term := semantics.Term{Kind: semantics.ProviderExpression, Operator: n.Key(name), Encoding: semantics.Key{Owner: "ontos", Name: "codec-v1"}, Data: codec.Encode(value), Operands: operands}
	if err := e.fillTerm(id, term); err != nil {
		return 0, err
	}
	return id, nil
}
func (e *exporter[L]) sort(value decl.Sort, depth int) (semantics.TermID, error) {
	if depth > e.limits.MaxDepth {
		return 0, semantics.ErrLimit
	}
	name := n.ProperSort().Name
	var operands []semantics.Operand
	if !value.IsProper() {
		name = "sort/arrow"
		for i, input := range value.Inputs() {
			child, err := e.sort(input, depth+1)
			if err != nil {
				return 0, err
			}
			operands = append(operands, semantics.Operand{Locator: []uint32{1, uint32(i + 1)}, Term: child})
		}
		result, _ := value.Result()
		child, err := e.sort(result, depth+1)
		if err != nil {
			return 0, err
		}
		operands = append(operands, semantics.Operand{Locator: []uint32{2}, Term: child})
	}
	return e.provider(name, value.Value(), operands, depth)
}
func (e *exporter[L]) sortFact(subject semantics.Subject, relation semantics.Key, value decl.Sort, depth int) error {
	term, err := e.sort(value, depth)
	if err != nil {
		return err
	}
	return e.fact(subject, relation, term)
}

// Each membership is one established source guarantee with its original rigid
// assumptions. Query results never mutate this graph or imply a principal kind.
func (e *exporter[L]) memberships(subject semantics.Subject, memberships []decl.Membership, depth int) error {
	for _, membership := range memberships {
		kind, err := e.kind(membership.Kind(), ast.NodeRef{}, 0, false, depth+1)
		if err != nil {
			return err
		}
		operands := []semantics.Operand{{Locator: []uint32{0}, Term: kind}}
		for i, assumption := range membership.Assumptions() {
			reference, err := e.allocateTerm(depth + 2)
			if err != nil {
				return err
			}
			if err := e.fillTerm(reference, semantics.Term{Kind: semantics.Reference, Target: e.parameters[assumption.Parameter()]}); err != nil {
				return err
			}
			contract, err := e.kind(assumption.Kind(), ast.NodeRef{}, 0, false, depth+2)
			if err != nil {
				return err
			}
			value, err := e.provider("membership/assumption", core.NewAtom(nil), []semantics.Operand{{Locator: []uint32{0}, Term: reference}, {Locator: []uint32{1}, Term: contract}}, depth+1)
			if err != nil {
				return err
			}
			operands = append(operands, semantics.Operand{Locator: []uint32{1, uint32(i)}, Term: value})
		}
		value, err := e.provider("membership/proven", core.NewAtom(nil), operands, depth)
		if err != nil {
			return err
		}
		if err := e.fact(subject, n.MembershipRelation(), value); err != nil {
			return err
		}
	}
	return nil
}
func (e *exporter[L]) identity(subject semantics.Subject, name string, value core.Value) error {
	term, err := e.provider("identity/"+name, value, nil, 1)
	if err != nil {
		return err
	}
	return e.fact(subject, n.IdentityRelation(), term)
}
func (e *exporter[L]) moduleIdentity(id semantics.EntityID, pin decl.Digest) error {
	// The vocabulary owns this operator, so identity does not assemble its name.
	term, err := e.provider(n.ModuleIdentity().Name, core.NewAtom(pin[:]), nil, 1)
	if err != nil {
		return err
	}
	return e.fact(semantics.Subject{Entity: id}, n.IdentityRelation(), term)
}
func (e *exporter[L]) declarationIdentity(id semantics.EntityID, value decl.DeclarationID) error {
	return e.identity(semantics.Subject{Entity: id}, "declaration", core.NewTuple(core.NewAtom(value.Module[:]), data.MustEncodeText(value.Name)))
}
func (e *exporter[L]) parameterIdentity(id semantics.EntityID, value decl.ParameterID) error {
	return e.identity(semantics.Subject{Entity: id}, "parameter", core.NewTuple(core.NewAtom(value.Declaration.Module[:]), data.MustEncodeText(value.Declaration.Name), data.EncodeInt(big.NewInt(int64(value.Index)))))
}
