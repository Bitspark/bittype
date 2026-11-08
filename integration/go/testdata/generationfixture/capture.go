package generationfixture

import (
	"context"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// SourceLanguage is the independent fixture language. It is not horos.
type SourceLanguage struct{}

func (SourceLanguage) ID() language.ID { return "generation-fixture" }
func (SourceLanguage) Name() string    { return "Independent generation fixture" }

const programTree ast.TreeID = "source-fixture"

// program is the fixture source: plain declarations, a generic Page whose body
// names its parameter, applications nested two deep, a Cell owned by a
// separate mapping library, and negative and boundary cases. The requested
// roots decide which declarations are generated.
func program() *graph {
	p := newGraph(programTree)
	for _, d := range []struct {
		name       string
		parameters []string
	}{{"Person", nil}, {"Page", []string{"T"}}, {"Outer", []string{"T"}}, {"Nested", []string{"T"}}, {"People", nil}, {"Cell", []string{"V"}}, {"Cells", nil}} {
		p.declare(d.name, d.parameters...)
	}
	p.body("Person", p.expression("person-body"))
	p.body("Page", p.expression("page-body", p.ref(p.parameter("Page", 0))))
	p.body("Outer", p.apply("Page", p.ref(p.parameter("Outer", 0))))
	p.body("Nested", p.apply("Page", p.apply("Page", p.ref(p.parameter("Nested", 0)))))
	p.body("People", p.apply("Page", p.ref(p.names["Person"])))
	p.body("Cell", p.expression("opaque", p.ref(p.parameter("Cell", 0))))
	p.body("Cells", p.apply("Page", p.apply("Cell", p.ref(p.names["Person"]))))

	// A phantom parameter: the Ignore recipe omits its argument natively.
	p.declare("Ignore", "T")
	p.body("Ignore", p.expression("opaque", p.ref(p.parameter("Ignore", 0))))
	p.declare("Kept")
	p.body("Kept", p.apply("Ignore", p.ref(p.names["Person"])))
	p.declare("Leak")
	rejected := p.expression("person-body")
	p.refute(rejected, false)
	p.body("Leak", p.apply("Ignore", rejected))
	p.declare("Blind")
	uncovered := p.expression("person-body")
	p.refute(uncovered, true)
	p.body("Blind", p.apply("Ignore", uncovered))

	// A body naming another declaration's binder.
	p.declare("Stray")
	p.body("Stray", p.apply("Page", p.ref(p.parameter("Page", 0))))
	p.declare("Rec")
	p.body("Rec", p.apply("Page", p.ref(p.names["Rec"])))
	// A parameter-headed application and an unapplied generic declaration.
	p.declare("Later", "F")
	p.body("Later", p.term(semantics.Term{Kind: semantics.Application, Callee: p.ref(p.parameter("Later", 0)), Arguments: []semantics.TermID{p.ref(p.names["Person"])}}))
	p.declare("Loose")
	p.body("Loose", p.ref(p.names["Page"]))

	// One term shared by two bodies is two occurrences. Through Borrow it
	// names a binder that Borrow does not own; without binders it is valid twice.
	p.declare("Lend", "T")
	p.declare("Borrow", "T")
	lent := p.apply("Page", p.ref(p.parameter("Lend", 0)))
	p.body("Lend", lent)
	p.body("Borrow", lent)
	p.declare("SameA")
	p.declare("SameB")
	same := p.apply("Page", p.ref(p.names["Person"]))
	p.body("SameA", same)
	p.body("SameB", same)

	// Admission covers whole subtrees, whatever the representation omits:
	// omitted actuals, content nested in them, and operands no rule captures.
	p.declare("Smuggle")
	p.body("Smuggle", p.apply("Ignore", p.ref(p.parameter("Page", 0))))
	p.declare("SmuggleLoose")
	p.body("SmuggleLoose", p.apply("Ignore", p.ref(p.names["Page"])))
	p.declare("Deep")
	nested := p.expression("person-body")
	p.refute(nested, false)
	p.body("Deep", p.apply("Ignore", p.apply("Page", nested)))
	p.declare("Hidden")
	p.body("Hidden", p.expression("person-body", p.ref(p.parameter("Page", 0))))
	// A declaration the source refuses, and a judgment with no verdict.
	p.declare("Refused")
	p.body("Refused", p.expression("person-body"))
	p.mark(semantics.Subject{Entity: p.names["Refused"]}, refuted)
	p.declare("Void")
	void := p.expression("person-body")
	p.mark(semantics.Subject{Term: void}, zero)
	p.body("Void", void)
	// An entity with a site that the source does not generate.
	p.other("Kindish")
	return p
}

type programChecker struct{ SourceLanguage }

func (programChecker) Check(v ast.View[SourceLanguage]) error {
	nodes := v.ASTView().Nodes
	if len(nodes) < 2 || nodes[0].Kind != "program" {
		return g.ErrInvalid
	}
	return nil
}

// captured is one fixture program captured and analyzed: the snapshot, the
// analysis result that publishes its graph, and each declaration's syntax site,
// which is the generation root a caller requests.
type captured struct {
	language representation.Language
	snapshot semantic.Snapshot
	evidence semantic.Result
	sites    map[string]ast.NodeRef
}

func capture(ctx context.Context, p *graph) (captured, error) {
	language, err := representation.NewLanguage[SourceLanguage](SourceLanguage{}, "fixture-v1")
	if err != nil {
		return captured{}, err
	}
	checker, err := check.NewChecker[SourceLanguage](language, programChecker{})
	if err != nil {
		return captured{}, err
	}
	tree := ast.AST[SourceLanguage]{ID: programTree, Root: 1, Nodes: []ast.Node[SourceLanguage]{{Kind: "program"}}}
	c := captured{language: language, sites: map[string]ast.NodeRef{}}
	for i, name := range p.order {
		node := ast.NodeID(i + 2)
		tree.Nodes = append(tree.Nodes, ast.Node[SourceLanguage]{Kind: "declaration"})
		tree.Nodes[0].Children = append(tree.Nodes[0].Children, ast.Child{Role: "declaration", Node: node})
		c.sites[name] = ast.NodeRef{Tree: programTree, Node: node}
	}
	packed, err := representation.Pack(language, tree)
	if err != nil {
		return captured{}, err
	}
	if c.snapshot, err = semantic.Capture(ctx, checker, semantic.Input{Entry: ast.NodeRef{Tree: programTree, Node: 1}, Trees: []representation.AST{packed}}); err != nil {
		return captured{}, err
	}
	a, err := analysis.NewAnalyzer[SourceLanguage](language, analysisContract(), analyzer{data: p.data})
	if err != nil {
		return captured{}, err
	}
	if c.evidence, err = a.Analyze(ctx, analysis.Input{Snapshot: c.snapshot}); err != nil {
		return captured{}, err
	}
	if !c.evidence.HasSemantics() {
		return captured{}, semantic.ErrEvidenceUnavailable
	}
	return c, nil
}
