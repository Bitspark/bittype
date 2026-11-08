package generationfixture

import (
	"context"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// The fixture's own analysis publishes an ordinary semantic graph. Its entity
// kinds, operators and fact relations are fixture vocabulary: the generation
// kernel must traverse the graph without reading any of them.
func key(name string) semantics.Key { return semantics.Key{Owner: "generation-fixture", Name: name} }

var (
	declarationKind = key("declaration")
	parameterKind   = key("parameter")
	moduleKind      = key("module")
	encoding        = key("fixture-encoding")
	// judged marks a term the fixture analysis judged admissible; refuted marks
	// one it judged inadmissible. A term with neither is outside its coverage.
	judged  = key("judged")
	refuted = key("refuted")
)

const analysisProfile = "generation-fixture-analysis-v1"

func analysisContract() analysis.Contract {
	return analysis.Contract{Semantic: semantic.Contract{
		Profile: analysisProfile,
		Obligations: []semantic.Obligation{
			{Key: key("validity"), Domain: key("entry-program")},
			{Key: key("declarations"), Domain: key("entry-program")},
		},
		SemanticsRequirements: []semantic.Requirement{
			{Obligation: key("validity"), Verdict: semantic.Holds},
			{Obligation: key("declarations"), Verdict: semantic.Undetermined},
		},
	}}
}

// Evidence is the requirement the fixture's source integration declares.
func Evidence() semantic.EvidenceRequirement {
	return semantic.EvidenceRequirement{Operation: analysis.Kind(), Profile: analysisProfile, Reports: analysisContract().Semantic.SemanticsRequirements, Semantics: true}
}

type analyzer struct {
	SourceLanguage
	data semantics.Data
}

func (a analyzer) Analyze(ctx context.Context, r analysis.Request[SourceLanguage]) (semantic.Result, error) {
	if err := ctx.Err(); err != nil {
		return semantic.Result{}, err
	}
	data := a.data
	return r.Publish(semantic.Publication{Facts: &data, Reports: []semantic.Report{
		{Obligation: key("validity"), Status: semantic.Complete, Verdict: semantic.Holds},
		{Obligation: key("declarations"), Status: semantic.Complete},
	}})
}

// graph builds the fixture's semantic records. Declaration sites are the
// declaration nodes of the captured fixture tree.
type graph struct {
	tree  ast.TreeID
	data  semantics.Data
	names map[string]semantics.EntityID
	order []string
}

func newGraph(tree ast.TreeID) *graph {
	g := &graph{tree: tree, names: map[string]semantics.EntityID{}}
	g.data.Entities = append(g.data.Entities, semantics.Entity{Kind: moduleKind, Name: "fixture"})
	// Term 1 is the shared fact value; the fixture judgment is the relation itself.
	g.data.Terms = append(g.data.Terms, semantics.Term{Kind: semantics.ProviderExpression, Operator: key("verdict"), Encoding: encoding})
	return g
}

const marker semantics.TermID = 1

func (g *graph) declare(name string, parameters ...string) semantics.EntityID {
	id := semantics.EntityID(len(g.data.Entities) + 1)
	site := ast.NodeRef{Tree: g.tree, Node: ast.NodeID(len(g.order) + 2)}
	g.data.Entities = append(g.data.Entities, semantics.Entity{Kind: declarationKind, Name: name, Owner: 1, Declarations: []ast.NodeRef{site}})
	for _, p := range parameters {
		pid := semantics.EntityID(len(g.data.Entities) + 1)
		g.data.Entities = append(g.data.Entities, semantics.Entity{Kind: parameterKind, Name: p, Owner: id})
		g.data.Entities[id-1].Parameters = append(g.data.Entities[id-1].Parameters, semantics.Parameter{Entity: pid, Sort: key("type")})
	}
	g.names[name] = id
	g.order = append(g.order, name)
	return id
}

func (g *graph) parameter(declaration string, i int) semantics.EntityID {
	return g.data.Entities[g.names[declaration]-1].Parameters[i].Entity
}

func (g *graph) term(t semantics.Term) semantics.TermID {
	g.data.Terms = append(g.data.Terms, t)
	id := semantics.TermID(len(g.data.Terms))
	g.data.Facts = append(g.data.Facts, semantics.Fact{Subject: semantics.Subject{Term: id}, Relation: judged, Term: marker})
	return id
}

func (g *graph) ref(entity semantics.EntityID) semantics.TermID {
	return g.term(semantics.Term{Kind: semantics.Reference, Target: entity})
}

func (g *graph) apply(callee string, arguments ...semantics.TermID) semantics.TermID {
	return g.term(semantics.Term{Kind: semantics.Application, Callee: g.ref(g.names[callee]), Arguments: arguments})
}

func (g *graph) expression(operator string, operands ...semantics.TermID) semantics.TermID {
	t := semantics.Term{Kind: semantics.ProviderExpression, Operator: key(operator), Encoding: encoding}
	for i, o := range operands {
		t.Operands = append(t.Operands, semantics.Operand{Locator: []uint32{uint32(i)}, Term: o})
	}
	return g.term(t)
}

func (g *graph) body(name string, term semantics.TermID) {
	g.data.Entities[g.names[name]-1].Body = term
}

// refute replaces a term's fixture judgment: refuted, or (with none) uncovered.
func (g *graph) refute(term semantics.TermID, uncovered bool) {
	facts := g.data.Facts[:0]
	for _, f := range g.data.Facts {
		if f.Subject.Term != term || f.Relation != judged {
			facts = append(facts, f)
		}
	}
	g.data.Facts = facts
	if !uncovered {
		g.data.Facts = append(g.data.Facts, semantics.Fact{Subject: semantics.Subject{Term: term}, Relation: refuted, Term: marker})
	}
}
