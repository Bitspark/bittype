package golang

import (
	"errors"
	"testing"

	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// packageGraphModel builds a model over a hand-written graph, which no bridge
// produced: the package "example.org/p" with the declaration T, a declaration
// without an owner, a declaration owned by a field, and two packages at the
// one path "example.org/twice". names are the data of the package-name facts
// of the first package; a nil entry is a fact whose term is a reference, which
// carries no data.
func packageGraphModel(t *testing.T, names ...[]byte) (Model, map[string]Type) {
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
	name := func(text []byte) semantics.TermID {
		term := semantics.Term{Kind: semantics.ProviderExpression, Operator: Key("package-name"), Encoding: Key("utf8"), Data: text}
		if text == nil {
			term = semantics.Term{Kind: semantics.Reference, Target: 2}
		}
		data.Terms = append(data.Terms, term)
		return semantics.TermID(len(data.Terms))
	}
	for _, text := range names {
		data.Facts = append(data.Facts, semantics.Fact{Subject: semantics.Subject{Entity: 1}, Relation: PackageNameRelation(), Term: name(text)})
	}
	// The field and both packages at the repeated path carry a name fact too,
	// so only the field's kind and the repetition keep them from being read.
	for id, text := range map[semantics.EntityID]string{4: "f", 6: "twice", 7: "twice"} {
		data.Facts = append(data.Facts, semantics.Fact{Subject: semantics.Subject{Entity: id}, Relation: PackageNameRelation(), Term: name([]byte(text))})
	}
	graph, err := semantics.New[s.Language](data, semantics.DefaultLimits(), func(ast.NodeRef) bool { return true })
	if err != nil {
		t.Fatal(err)
	}
	state := &model{graph: graph}
	return Model{state}, map[string]Type{"T": {state, 1}, "Orphan": {state, 2}, "Inner": {state, 3}}
}

// The readers of a package's name read the fact and nothing else: without
// exactly one well-formed fact they fail, and never fall back to the path.
func TestPackageNameReadersFailWithoutTheFact(t *testing.T) {
	model, types := packageGraphModel(t, []byte("declared"))
	if name, err := model.PackageName("example.org/p"); err != nil || name != "declared" {
		t.Fatal("package name", name, err)
	}
	if pkg, err := types["T"].Package(); err != nil || pkg != (DeclaringPackage{Path: "example.org/p", Name: "declared"}) {
		t.Fatal("declaring package", pkg, err)
	}
	for name, facts := range map[string][][]byte{
		"no fact":                     nil,
		"two facts":                   {[]byte("declared"), []byte("other")},
		"an empty name":               {{}},
		"a term that carries no data": {nil},
	} {
		model, types := packageGraphModel(t, facts...)
		if got, err := model.PackageName("example.org/p"); !errors.Is(err, ErrForeign) || got != "" {
			t.Fatal(name, "gave a package name:", got, err)
		}
		if got, err := types["T"].Package(); !errors.Is(err, ErrForeign) || got != (DeclaringPackage{}) {
			t.Fatal(name, "gave a declaring package:", got, err)
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
