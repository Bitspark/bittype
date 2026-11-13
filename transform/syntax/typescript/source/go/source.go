package source

import (
	"context"
	"encoding/json"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Contract requires the TypeScript analysis of the base profile: the one
// evidence requirement that publishes the semantic graph the kernel walks.
func Contract() transform.Contract {
	return transform.Contract{Evidence: []semantic.EvidenceRequirement{n.Evidence()}}
}

// primitiveOperator is the operator of the term of a primitive type. The
// compiler helper writes it, and the analysed vocabulary publishes no name
// for it (it publishes ObjectOperator), so this reader spells it.
func primitiveOperator() semantics.Key { return n.Key("primitive") }

// interpretation is the TypeScript reading of one captured graph for the
// kernel. It reads through the model of the analysed vocabulary
// (semantics/typescript/go) where that model has an entry: a declaration by
// its entity. The model reads the properties of a type and gives no handle
// for the term that carries them, so that term is read from the graph with
// the names and the metadata type the vocabulary publishes.
//
// The methods do not check the context: the kernel checks it before each step
// of its walk, and no method here does work of its own between two steps.
type interpretation struct {
	captured g.Captured[s.Language] // handles resolve against this capture only
	graph    semantics.Semantics[s.Language]
	model    n.Model
	language representation.Language
}

func open(c g.Captured[s.Language]) (*interpretation, error) {
	if !c.Valid() {
		return nil, g.ErrInvalid
	}
	request := c.Request()
	// Contract states one requirement, so the graph's receipt is evidence 0.
	result, ok := request.Evidence(0)
	if !ok {
		return nil, semantic.ErrEvidenceUnavailable
	}
	model, err := n.Open(result)
	if err != nil {
		return nil, err
	}
	return &interpretation{captured: c, graph: c.Graph(), model: model, language: request.Snapshot().Language()}, nil
}

// entity reads the entity of a declaration the kernel reached and the unit
// that owns it. named is false for an entity that is no interface and no type
// alias of a captured unit: a property, a type parameter, a type of a
// package, a builtin.
//
// COPY (AC-11, finding AG-8 of the two other pairs): resolving a declaration
// handle against the capture, reading its entity and its owner, and the
// unsupported judgment below are the same steps in the Go source reader of
// the same lane, with only the kinds that count as "named" different. They
// have no agnostic implementation.
func (i *interpretation) entity(d g.Declaration[s.Language]) (e, owner semantics.Entity, named bool, err error) {
	id, ok := d.In(i.captured)
	if !ok {
		return e, owner, false, g.ErrScope
	}
	if e, ok = i.graph.Entity(id); !ok {
		return e, owner, false, g.ErrScope
	}
	owner, _ = i.graph.Entity(e.Owner)
	named = (e.Kind == n.AliasKind() || e.Kind == n.InterfaceKind()) && owner.Kind != n.PackageKind() && len(owner.Declarations) > 0
	return e, owner, named, nil
}

// site is where a judgment about a use is reported: the declaring site of the
// declaration the use occurs in. The kernel knows the site of the occurrence
// itself and gives a source no way to ask for it.
func (i *interpretation) site(u g.Use[s.Language]) ast.NodeRef {
	if id, ok := u.Declaration().In(i.captured); ok {
		if e, ok := i.graph.Entity(id); ok && len(e.Declarations) > 0 {
			return e.Declarations[0]
		}
	}
	return ast.NodeRef{}
}

func unsupported(site ast.NodeRef, message string) g.Judgment {
	diagnostic := &text.Diagnostic{Code: text.Unsupported, Message: message}
	if site.Tree != "" {
		diagnostic.Node = &site
	}
	return g.Judgment{Verdict: g.Unsupported, Diagnostic: diagnostic}
}

func (i *interpretation) Declaration(_ context.Context, d g.Declaration[s.Language]) (g.Selector, bool, error) {
	e, owner, named, err := i.entity(d)
	if err != nil || !named {
		return g.Selector{}, false, err
	}
	key, err := Named(i.language, owner.Name, e.Name)
	return key, true, err
}

// definition is the term that carries the properties of the type that a
// declaration declares, or zero where the declared type is no object type.
func (i *interpretation) definition(e semantics.Entity) semantics.TermID {
	if terms := i.graph.Facts(semantics.Subject{Term: e.Body}, n.PropertiesRelation()); len(terms) == 1 {
		return terms[0]
	}
	return 0
}

// Definition states the term that defines a declaration. The TypeScript
// analysis gives an interface or an alias of an object type an entity whose
// body is a reference to that entity, and relates the object type to it by
// the fact "properties". The kernel reads no fact, so without this statement
// it meets such a declaration as a use of itself. ok is false for a
// declaration whose declared type is no object type: its body is the term
// the graph records.
func (i *interpretation) Definition(_ context.Context, d g.Declaration[s.Language]) (semantics.TermID, bool, error) {
	e, _, named, err := i.entity(d)
	if err != nil || !named {
		return 0, false, err
	}
	term := i.definition(e)
	return term, term != 0, nil
}

// primitive reads the TypeScript spelling of a primitive type's term.
func primitive(term semantics.Term) (string, bool) {
	var name string
	if term.Operator != primitiveOperator() || json.Unmarshal(term.Data, &name) != nil || name == "" {
		return "", false
	}
	return name, true
}

func (i *interpretation) Expression(_ context.Context, u g.Use[s.Language]) (g.Selector, error) {
	id, ok := u.In(i.captured)
	if !ok {
		return g.Selector{}, g.ErrScope
	}
	term, _ := i.graph.Term(id)
	if term.Operator == n.ObjectOperator() {
		return ObjectType(i.language)
	}
	if name, ok := primitive(term); ok {
		return Primitive(i.language, name)
	}
	return g.Selector{}, semantic.ErrInvalidRequest
}

// JudgeDeclaration admits an interface or a type alias of a captured unit
// that has no type parameter. What it declares is judged as a use, and what a
// pair makes of it is the pair's declaration contribution's decision.
func (i *interpretation) JudgeDeclaration(_ context.Context, d g.Declaration[s.Language]) (g.Judgment, error) {
	e, _, named, err := i.entity(d)
	if err != nil {
		return g.Judgment{}, err
	}
	var site ast.NodeRef
	if len(e.Declarations) > 0 {
		site = e.Declarations[0]
	}
	switch {
	case !named:
		return unsupported(site, "this TypeScript source reader reads an interface or a type alias of a captured file"), nil
	case len(e.Parameters) != 0:
		return unsupported(site, "this TypeScript source reader reads no generic declaration"), nil
	}
	return g.Judgment{Verdict: g.Admitted}, nil
}

// JudgeUse admits an object type and a primitive type. The successful
// TypeScript analysis established that the use is valid TypeScript. This adds
// what the reader reads.
func (i *interpretation) JudgeUse(_ context.Context, u g.Use[s.Language]) (g.Judgment, error) {
	id, ok := u.In(i.captured)
	if !ok {
		return g.Judgment{}, g.ErrScope
	}
	term, ok := i.graph.Term(id)
	if !ok {
		return g.Judgment{}, g.ErrScope
	}
	site := i.site(u)
	switch term.Kind {
	case semantics.Reference:
		return unsupported(site, "this TypeScript source reader reads no reference to a named type or a type parameter"), nil
	case semantics.Application:
		return unsupported(site, "this TypeScript source reader reads no instantiated generic type"), nil
	}
	if _, ok := primitive(term); ok || term.Operator == n.ObjectOperator() {
		return g.Judgment{Verdict: g.Admitted}, nil
	}
	display := "of the form " + term.Operator.Name
	if t, err := i.model.Type(id); err == nil {
		display = t.String()
	}
	return unsupported(site, "this TypeScript source reader reads object types and primitive types, and not the type "+display), nil
}
