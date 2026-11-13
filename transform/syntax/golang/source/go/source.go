package source

import (
	"context"
	"go/types"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Contract requires the Go analysis of the base profile: the one evidence
// requirement that publishes the semantic graph the kernel walks.
func Contract() transform.Contract {
	return transform.Contract{Evidence: []semantic.EvidenceRequirement{n.Evidence()}}
}

// interpretation is the Go reading of one captured graph for the kernel. It
// classifies declarations, names forms and judges what this reader reads. It
// reads through the model of the analysed vocabulary (semantics/golang/go)
// where that model has an entry: a type by its term. The model has no entry by
// a declaration's entity, so the entity record is read from the graph.
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

// entity reads the entity of a declaration the kernel reached and the package
// that owns it. named is false for an entity that is no named type of a
// package: a field, a predeclared type, a type parameter, a package.
func (i *interpretation) entity(d g.Declaration[s.Language]) (e, owner semantics.Entity, named bool, err error) {
	id, ok := d.In(i.captured)
	if !ok {
		return e, owner, false, g.ErrScope
	}
	if e, ok = i.graph.Entity(id); !ok {
		return e, owner, false, g.ErrScope
	}
	owner, _ = i.graph.Entity(e.Owner)
	named = (e.Kind == n.DefinedKind() || e.Kind == n.AliasKind()) && owner.Kind == n.PackageKind()
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
	key, err := Defined(i.language, owner.Name, e.Name)
	return key, true, err
}

// Definition states the term that defines a declaration. The Go analysis
// gives a named type an entity whose body is a reference to that entity, and
// relates the type's definition to it by the fact "underlying". The kernel
// reads no fact, so without this statement it meets every Go declaration as a
// use of itself. ok is false for a declaration that this reader does not read.
func (i *interpretation) Definition(_ context.Context, d g.Declaration[s.Language]) (semantics.TermID, bool, error) {
	e, _, named, err := i.entity(d)
	if err != nil || !named || e.Kind != n.DefinedKind() {
		return 0, false, err
	}
	declared, err := i.model.Type(e.Body)
	if err != nil {
		return 0, false, nil
	}
	underlying := declared.Underlying()
	return underlying.ID(), underlying.Valid(), nil
}

func (i *interpretation) Expression(_ context.Context, u g.Use[s.Language]) (g.Selector, error) {
	id, ok := u.In(i.captured)
	if !ok {
		return g.Selector{}, g.ErrScope
	}
	t, err := i.model.Type(id)
	if err != nil {
		return g.Selector{}, semantic.ErrInvalidRequest
	}
	term, _ := i.graph.Term(id)
	switch {
	case term.Operator == n.StructOperator():
		return StructType(i.language)
	case t.IsBasic(types.String):
		return Basic(i.language, BasicString)
	case t.IsBasic(types.Bool):
		return Basic(i.language, BasicBool)
	}
	return g.Selector{}, semantic.ErrInvalidRequest
}

// JudgeDeclaration admits a defined type of a captured package that has no
// type parameter. What its underlying type may be is judged as a use, and
// what a pair makes of it is the pair's declaration contribution's decision.
func (i *interpretation) JudgeDeclaration(_ context.Context, d g.Declaration[s.Language]) (g.Judgment, error) {
	e, owner, named, err := i.entity(d)
	if err != nil {
		return g.Judgment{}, err
	}
	var site ast.NodeRef
	if len(e.Declarations) > 0 {
		site = e.Declarations[0]
	}
	switch {
	case !named:
		return unsupported(site, "this Go source reader reads a named type of a package"), nil
	case e.Kind != n.DefinedKind():
		return unsupported(site, "this Go source reader reads no alias declaration"), nil
	case len(owner.Declarations) == 0:
		return unsupported(site, "this Go source reader reads no type of a supplied package"), nil
	case len(e.Parameters) != 0:
		return unsupported(site, "this Go source reader reads no generic declaration"), nil
	}
	return g.Judgment{Verdict: g.Admitted}, nil
}

// JudgeUse admits a struct type, the basic types string and bool, and a use
// of a defined type. The successful Go analysis established that the use is
// valid Go. This adds what the reader reads.
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
		if target, _ := i.graph.Entity(term.Target); target.Kind != n.DefinedKind() {
			return unsupported(site, "this Go source reader reads a use of a defined type, and of no alias or type parameter"), nil
		}
		return g.Judgment{Verdict: g.Admitted}, nil
	case semantics.Application:
		return unsupported(site, "this Go source reader reads no instantiated generic type"), nil
	}
	t, err := i.model.Type(id)
	if err != nil {
		return g.Judgment{Verdict: g.Undetermined}, nil
	}
	if term.Operator == n.StructOperator() || t.IsBasic(types.String) || t.IsBasic(types.Bool) {
		return g.Judgment{Verdict: g.Admitted}, nil
	}
	return unsupported(site, "this Go source reader reads struct types, string and bool, and not the type "+t.String()), nil
}
