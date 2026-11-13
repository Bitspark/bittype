package golang

import (
	"errors"
	"testing"

	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// nameTerm is a package-name term in the form the bridge writes.
func nameTerm(text string) semantics.Term {
	return semantics.Term{Kind: semantics.ProviderExpression, Operator: semantics.Key{Owner: "go", Name: "package-name"}, Encoding: semantics.Key{Owner: "go", Name: "utf8"}, Data: []byte(text)}
}

// packageGraphModel builds a model over a hand-written graph, which no bridge
// produced: the package "example.org/p" with the declaration T, a declaration
// without an owner, a declaration owned by a field, and two packages at the
// one path "example.org/twice". Each of names becomes the term of one
// package-name fact of the first package.
func packageGraphModel(t *testing.T, names ...semantics.Term) (Model, map[string]Type) {
	t.Helper()
	data := semantics.Data{
		Entities: []semantics.Entity{
			{Kind: PackageKind(), Name: "example.org/p"},
			{Kind: DefinedKind(), Name: "T", Owner: 1, Body: 1},
			{Kind: DefinedKind(), Name: "Orphan", Body: 2},
			{Kind: Key("field"), Name: "f"},
			{Kind: DefinedKind(), Name: "Inner", Owner: 4, Body: 3},
			{Kind: PackageKind(), Name: "example.org/twice"},
			{Kind: PackageKind(), Name: "example.org/twice"},
		},
		Terms: []semantics.Term{
			{Kind: semantics.Reference, Target: 2},
			{Kind: semantics.Reference, Target: 3},
			{Kind: semantics.Reference, Target: 5},
		},
	}
	fact := func(subject semantics.EntityID, term semantics.Term) {
		data.Terms = append(data.Terms, term)
		data.Facts = append(data.Facts, semantics.Fact{Subject: semantics.Subject{Entity: subject}, Relation: PackageNameRelation(), Term: semantics.TermID(len(data.Terms))})
	}
	for _, term := range names {
		fact(1, term)
	}
	// The field and both packages at the repeated path carry a name fact too,
	// so only the field's kind and the repetition keep them from being read.
	fact(4, nameTerm("f"))
	fact(6, nameTerm("twice"))
	fact(7, nameTerm("twice"))
	graph, err := semantics.New[s.Language](data, semantics.DefaultLimits(), func(ast.NodeRef) bool { return true })
	if err != nil {
		t.Fatal(err)
	}
	state := &model{graph: graph}
	return Model{state}, map[string]Type{"T": {state, 1}, "Orphan": {state, 2}, "Inner": {state, 3}}
}

// The readers of a package's name read the fact and nothing else: without
// exactly one fact in the published form they fail, and never fall back to
// the path.
func TestPackageNameReadersFailWithoutTheFact(t *testing.T) {
	model, types := packageGraphModel(t, nameTerm("declared"))
	if name, err := model.PackageName("example.org/p"); err != nil || name != "declared" {
		t.Fatal("package name", name, err)
	}
	if pkg, err := types["T"].Package(); err != nil || pkg != (DeclaringPackage{Path: "example.org/p", Name: "declared"}) {
		t.Fatal("declaring package", pkg, err)
	}
	other := func(change func(*semantics.Term)) semantics.Term {
		term := nameTerm("declared")
		change(&term)
		return term
	}
	for name, facts := range map[string][]semantics.Term{
		"no fact":       nil,
		"two facts":     {nameTerm("declared"), nameTerm("other")},
		"an empty name": {nameTerm("")},
		// A fact whose term has another form is no name fact. Each row differs
		// from the published form in one thing.
		"a reference":                {{Kind: semantics.Reference, Target: 2}},
		"another operator":           {other(func(term *semantics.Term) { term.Operator = DisplayRelation() })},
		"another owner's operator":   {other(func(term *semantics.Term) { term.Operator.Owner = "other" })},
		"another encoding":           {other(func(term *semantics.Term) { term.Encoding = Key("json-1") })},
		"another owner's encoding":   {other(func(term *semantics.Term) { term.Encoding.Owner = "other" })},
		"an operand":                 {other(func(term *semantics.Term) { term.Operands = []semantics.Operand{{Locator: []uint32{0}, Term: 1}} })},
		"text that is no UTF-8":      {other(func(term *semantics.Term) { term.Data = []byte{0xff, 0xfe} })},
		"a truncated UTF-8 sequence": {other(func(term *semantics.Term) { term.Data = []byte("pag\xc3") })},
	} {
		model, types := packageGraphModel(t, facts...)
		if got, err := model.PackageName("example.org/p"); !errors.Is(err, ErrForeign) || got != "" {
			t.Error(name, "gave a package name:", got, err)
		}
		if got, err := types["T"].Package(); !errors.Is(err, ErrForeign) || got != (DeclaringPackage{}) {
			t.Error(name, "gave a declaring package:", got, err)
		}
	}
	// The readers compare the form of the term and not Go's rule for a package
	// name: text in the published form is returned as it is written.
	for _, text := range []string{"my-pages x", "example.org/p", "_", "pagés"} {
		model, types := packageGraphModel(t, nameTerm(text))
		if got, err := model.PackageName("example.org/p"); err != nil || got != text {
			t.Fatal("a name in the published form was not returned:", text, got, err)
		}
		if got, err := types["T"].Package(); err != nil || got.Name != text {
			t.Fatal("a name in the published form was not returned:", text, got, err)
		}
	}
	// A declaration that no package owns has no package, whatever its owner carries.
	for _, name := range []string{"Orphan", "Inner"} {
		if got, err := types[name].Package(); !errors.Is(err, ErrForeign) || got != (DeclaringPackage{}) {
			t.Fatal(name, "has a package:", got, err)
		}
	}
	// An entity that is not a package is not found by its name, and a path
	// that two packages have names neither.
	for _, path := range []string{"f", "example.org/twice", "example.org/missing"} {
		if got, err := model.PackageName(path); !errors.Is(err, ErrForeign) || got != "" {
			t.Fatal("a package name was read for", path, got, err)
		}
	}
	if _, err := (Model{}).PackageName("example.org/p"); !errors.Is(err, ErrForeign) {
		t.Fatal("an invalid model has a package", err)
	}
}
