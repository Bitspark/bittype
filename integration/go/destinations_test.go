package integration_test

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/format/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	tsfiles "github.com/Bitspark/bittype/format/files/typescript/provider/go"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	"github.com/Bitspark/bittype/semantics/go"
	gn "github.com/Bitspark/bittype/semantics/golang/go"
	tn "github.com/Bitspark/bittype/semantics/typescript/go"
	"github.com/Bitspark/bittype/source/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	rs "github.com/Bitspark/bittype/syntax/rust/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	gomodel "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	gosymbol "github.com/Bitspark/bittype/transform/syntax/golang/symbol/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/boolean/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/declarations/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/list/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/optional/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/record/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/text/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	patterns "github.com/Bitspark/bittype/transform/syntax/horos/mapping/pattern/go"
	hsource "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	rustmodel "github.com/Bitspark/bittype/transform/syntax/rust/model/go"
	rustsymbol "github.com/Bitspark/bittype/transform/syntax/rust/symbol/go"
	rgtarget "github.com/Bitspark/bittype/transform/syntax/rust/target/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	tssymbol "github.com/Bitspark/bittype/transform/syntax/typescript/symbol/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// These tests read what a sealed generation output records about the
// destination of each name use. docs/destinations.md is their specification.

// externalUsesBySite returns the external uses of an output by their site. A
// site that two uses name fails the test.
func externalUsesBySite(t *testing.T, out g.Output) map[ast.NodeRef]g.ExternalUse {
	t.Helper()
	uses := map[ast.NodeRef]g.ExternalUse{}
	for _, use := range out.ExternalUses() {
		if _, twice := uses[use.At.Site]; twice {
			t.Fatalf("two external uses at %v", use.At.Site)
		}
		uses[use.At.Site] = use
	}
	return uses
}

func site[L any](tree ast.AST[L], index int) ast.NodeRef {
	return ast.NodeRef{Tree: tree.ID, Node: ast.NodeID(index + 1)}
}

// rustPathNames returns the names of the segments of a path-type node and the
// node of its last segment.
func rustPathNames(tree rs.AST, node ast.Node[rs.Language]) ([]string, ast.NodeID) {
	var names []string
	var last ast.NodeID
	for _, child := range node.Children {
		if child.Role == rs.RoleSegment {
			names = append(names, nodeAttribute(tree.Nodes[child.Node-1], rs.AttributeName))
			last = child.Node
		}
	}
	return names, last
}

// For the Go pagination output, ExternalUses holds one Go-issued predeclared
// use for each string and bool node, and no other use.
func TestGoPaginationRecordsPredeclaredDestinations(t *testing.T) {
	e, _, input := generationInput(t, paginationModule())
	out := generatedOutput(t, e, input, golangText(), goRequest)
	tree := outputTree[gs.Language](t, out)
	uses := externalUsesBySite(t, out)
	recorded := map[string]int{}
	for i, node := range tree.Nodes {
		if node.Kind != gs.KindNamedType || len(node.References) != 0 {
			continue
		}
		name := nodeAttribute(node, gs.AttributeName)
		use, ok := uses[site(tree, i)]
		if !ok {
			t.Errorf("the named type %s at node %d records no destination", name, i+1)
			continue
		}
		symbol, err := gosymbol.ReadExternal(out.Snapshot().Language(), use)
		if err != nil {
			t.Errorf("the use at node %d is not read by the Go issuer: %v", i+1, err)
			continue
		}
		if predeclared, ok := symbol.Predeclared(); !ok || predeclared != name {
			t.Errorf("the named type %s at node %d records the destination %q", name, i+1, predeclared)
			continue
		}
		recorded[name]++
	}
	// Page.next and Person.name are text. Person.active and the helper's
	// Present are booleans.
	if want := map[string]int{"string": 2, "bool": 2}; recorded["string"] != want["string"] || recorded["bool"] != want["bool"] || len(recorded) != len(want) {
		t.Errorf("recorded predeclared destinations %v, want %v", recorded, want)
	}
	if len(uses) != recorded["string"]+recorded["bool"] {
		t.Errorf("the output holds %d external uses for %d predeclared nodes", len(uses), recorded["string"]+recorded["bool"])
	}
}

// For the Rust pagination output, ExternalUses holds one Rust-issued use for
// each external path, on the segment that names the item, and one for each
// bool.
func TestRustPaginationRecordsExternalDestinations(t *testing.T) {
	e, _, input := generationInput(t, paginationModule())
	out := generatedOutput(t, e, input, rustText(), rustRequest)
	tree := outputTree[rs.Language](t, out)
	uses := externalUsesBySite(t, out)
	recorded := map[string]int{}
	read := func(i int, at ast.NodeID) (rustsymbol.Symbol, bool) {
		use, ok := uses[ast.NodeRef{Tree: tree.ID, Node: at}]
		if !ok {
			t.Errorf("the %s at node %d records no destination", tree.Nodes[i].Kind, i+1)
			return rustsymbol.Symbol{}, false
		}
		symbol, err := rustsymbol.ReadExternal(out.Snapshot().Language(), use)
		if err != nil {
			t.Errorf("the use at node %d is not read by the Rust issuer: %v", at, err)
			return rustsymbol.Symbol{}, false
		}
		return symbol, true
	}
	for i, node := range tree.Nodes {
		switch node.Kind {
		case rs.KindPathType:
			names, last := rustPathNames(tree, node)
			if nodeAttribute(node, rs.AttributeAbsolute) != "true" {
				continue
			}
			written := strings.Join(names, "::")
			symbol, ok := read(i, last)
			if !ok {
				recorded["missing "+written]++
				continue
			}
			crate, path, ok := symbol.Item()
			if !ok || crate != rustsymbol.Std() || strings.Join(append([]string{crate.Name()}, path...), "::") != written {
				t.Errorf("the path %s at node %d records the destination %v %v", written, i+1, crate, path)
				continue
			}
			recorded[written]++
		case rs.KindPrimitiveType:
			name := nodeAttribute(node, rs.AttributeName)
			symbol, ok := read(i, ast.NodeID(i+1))
			if !ok {
				recorded["missing "+name]++
				continue
			}
			if primitive, ok := symbol.Primitive(); !ok || primitive != name || symbol != rustsymbol.Bool() {
				t.Errorf("the primitive type %s at node %d records the destination %q", name, i+1, primitive)
				continue
			}
			recorded[name]++
		}
	}
	// Page.items and Other are lists, Page.next is an optional text,
	// Person.name is text and Person.active is a boolean.
	want := map[string]int{"std::vec::Vec": 2, "std::option::Option": 1, "std::string::String": 2, "bool": 1}
	if !mapsEqual(recorded, want) {
		t.Errorf("recorded external destinations %v, want %v", recorded, want)
	}
	if len(uses) != 6 {
		t.Errorf("the output holds %d external uses, want 6", len(uses))
	}
}

func mapsEqual(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// destinationCheck is what one language states about its sealed outputs: the
// nodes that are name uses, and the reader of its issuer.
type destinationCheck[L any] struct {
	// uses returns the nodes of a tree that are name uses.
	uses func(t *testing.T, tree ast.AST[L]) []ast.NodeID
	// reference is the role of a reference to a declaration or a parameter.
	reference string
	// read fails for a use that the language's issuer did not issue.
	read func(representation.Language, g.ExternalUse) error
}

// checkDestinations asserts that every name use of a sealed output carries
// exactly one of a reference and an external use, that no other node carries
// an external use, and that the language's issuer reads every use. It returns
// the numbers of references and of external uses.
func checkDestinations[L any](t *testing.T, out g.Output, check destinationCheck[L]) (references, external int) {
	t.Helper()
	tree := outputTree[L](t, out)
	uses := externalUsesBySite(t, out)
	nameUses := map[ast.NodeID]bool{}
	for _, id := range check.uses(t, tree) {
		nameUses[id] = true
		node := tree.Nodes[id-1]
		at := ast.NodeRef{Tree: tree.ID, Node: id}
		_, recorded := uses[at]
		linked := len(node.References) == 1 && node.References[0].Role == check.reference
		switch {
		case len(node.References) > 1 || (len(node.References) == 1 && !linked):
			t.Errorf("the %s %q at node %d carries references that are no destination: %v", node.Kind, nodeAttribute(node, "name"), id, node.References)
		case linked && recorded:
			t.Errorf("the %s %q at node %d carries a reference and an external use", node.Kind, nodeAttribute(node, "name"), id)
		case !linked && !recorded:
			t.Errorf("the %s %q at node %d carries no destination", node.Kind, nodeAttribute(node, "name"), id)
		case linked:
			references++
		default:
			external++
		}
	}
	for at, use := range uses {
		if at.Tree != tree.ID || !nameUses[at.Node] {
			t.Errorf("an external use is recorded at %v, which is no name use", at)
		}
		if err := check.read(out.Snapshot().Language(), use); err != nil {
			t.Errorf("the external use at %v is not read by its language's issuer: %v", at, err)
		}
	}
	for i, node := range tree.Nodes {
		if len(node.References) != 0 && !nameUses[ast.NodeID(i+1)] {
			t.Errorf("the %s at node %d carries a reference and is no name use", node.Kind, i+1)
		}
	}
	return references, external
}

func typescriptDestinations() destinationCheck[ts.Language] {
	return destinationCheck[ts.Language]{
		uses: func(t *testing.T, tree ts.AST) (out []ast.NodeID) {
			for i, node := range tree.Nodes {
				if node.Kind == ts.KindTypeReference {
					out = append(out, ast.NodeID(i+1))
				}
			}
			return out
		},
		reference: ts.RoleTarget,
		read: func(language representation.Language, use g.ExternalUse) error {
			_, err := tssymbol.ReadExternal(language, use)
			return err
		},
	}
}

func golangDestinations() destinationCheck[gs.Language] {
	return destinationCheck[gs.Language]{
		uses: func(t *testing.T, tree gs.AST) (out []ast.NodeID) {
			for i, node := range tree.Nodes {
				switch node.Kind {
				case gs.KindNamedType:
					out = append(out, ast.NodeID(i+1))
				case gs.KindQualifiedType:
					t.Errorf("node %d is a qualified type, which the Go target does not write", i+1)
				}
			}
			return out
		},
		reference: gs.RoleTarget,
		read: func(language representation.Language, use g.ExternalUse) error {
			_, err := gosymbol.ReadExternal(language, use)
			return err
		},
	}
}

func rustDestinations() destinationCheck[rs.Language] {
	return destinationCheck[rs.Language]{
		uses: func(t *testing.T, tree rs.AST) (out []ast.NodeID) {
			for i, node := range tree.Nodes {
				switch node.Kind {
				case rs.KindPathType:
					// The last segment names the item. The segments before it
					// qualify the use.
					_, last := rustPathNames(tree, node)
					if last == 0 {
						t.Errorf("the path at node %d has no segment", i+1)
						continue
					}
					out = append(out, last)
				case rs.KindPrimitiveType:
					out = append(out, ast.NodeID(i+1))
				case rs.KindQualifiedPathType:
					t.Errorf("node %d is a qualified path, which the Rust target does not write", i+1)
				}
			}
			return out
		},
		reference: rs.RoleTarget,
		read: func(language representation.Language, use g.ExternalUse) error {
			_, err := rustsymbol.ReadExternal(language, use)
			return err
		},
	}
}

// For TypeScript, Go and Rust, every name use of a sealed output carries
// exactly one of a reference and an external use.
func TestEveryNameUseOfASealedOutputCarriesOneDestination(t *testing.T) {
	t.Run("typescript", func(t *testing.T) {
		e, _, input := generationInput(t, paginationModule())
		references, external := checkDestinations(t, generatedOutput(t, e, input, typescriptText(), tsRequest), typescriptDestinations())
		if references == 0 || external != 0 {
			t.Errorf("the pagination output has %d references and %d external uses", references, external)
		}
		// One output with a declaration named Cell and an external Cell.
		references, external = checkDestinations(t, generatedCellPage(t), typescriptDestinations())
		if references == 0 || external != 1 {
			t.Errorf("the Cell output has %d references and %d external uses", references, external)
		}
	})
	t.Run("go", func(t *testing.T) {
		e, _, input := generationInput(t, paginationModule())
		references, external := checkDestinations(t, generatedOutput(t, e, input, golangText(), goRequest), golangDestinations())
		if references == 0 || external != 4 {
			t.Errorf("the pagination output has %d references and %d external uses", references, external)
		}
	})
	t.Run("rust", func(t *testing.T) {
		e, _, input := generationInput(t, paginationModule())
		references, external := checkDestinations(t, generatedOutput(t, e, input, rustText(), rustRequest), rustDestinations())
		if references == 0 || external != 6 {
			t.Errorf("the pagination output has %d references and %d external uses", references, external)
		}
	})
}

// stringHelper is a fixture mapping library: it maps the horos boolean to a
// helper declaration that is named string, a defined Go type.
func stringHelper(key g.HelperKey) rules.Library[gs.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, gs.Language, hsource.Reader], error) {
		pattern, err := patterns.Head(from, "scalar", "bool")
		if err != nil {
			return g.MappingDescription[h.Language, gs.Language, hsource.Reader]{}, err
		}
		return g.MappingDescription[h.Language, gs.Language, hsource.Reader]{
			Profile:  rules.RepresentationProfile,
			Bindings: []g.BindingDescription[h.Language, gs.Language, hsource.Reader]{{ID: "fixture/string-helper/1", Input: pattern, Output: g.Helper[gs.Language](key)}},
			Helpers: []g.HelperRecipe[gs.Language]{{Key: key, Declaration: g.NativeDeclaration[gs.Language]{
				Name: "string", Form: gomodel.Form(gomodel.FormDefined),
				Body: gomodel.Record(gomodel.Field("Flag", false, gomodel.Primitive(gomodel.PrimitiveBool))),
			}}},
		}, nil
	}
}

// goPredeclaredEntity says whether an entity of a Go graph is the predeclared
// type that a symbol names: an entity of the builtin kind with the symbol's
// name, without an owner and without a declaration site.
func goPredeclaredEntity(symbol gosymbol.Symbol, entity semantics.Entity) bool {
	name, ok := symbol.Predeclared()
	return ok && entity.Kind == gn.BuiltinKind() && entity.Name == name && entity.Owner == 0 && len(entity.Declarations) == 0
}

// A fixture mapping library supplies a Go helper named string, and a record
// has a text field. The recorded destination at the field's type is the
// predeclared symbol, and the compiler binds that site to the helper, an
// entity of the kind defined. The two differ. With the builtin boolean
// mapping, which declares no such helper, the same site binds to the
// predeclared entity.
func TestGoDeclarationNamedStringCapturesAPredeclaredUse(t *testing.T) {
	module := genModule("example.capture", genMap(), genEntry("Person", genType(genPerson())))
	type libraries = []g.MappingLibrary[h.Language, gs.Language, hsource.Reader]
	optionalKey, helperKey := g.NewHelperKey("horos/optional"), g.NewHelperKey("fixture/string")
	for name, fixture := range map[string]struct {
		libraries libraries
		captured  bool
	}{
		"a helper named string": {libraries{text.Library(), stringHelper(helperKey), list.Library(), record.Library(optionalKey), declarations.Library(), optional.Library(optionalKey)}, true},
		"no such helper":        {libraries{text.Library(), boolean.Library(), list.Library(), record.Library(optionalKey), declarations.Library(), optional.Library(optionalKey)}, false},
	} {
		t.Run(name, func(t *testing.T) {
			e, _, input := generationInput(t, module)
			out := generatedOutput(t, e, input, golangText(), func(reg *c.Registry, input transform.Input, to representation.Language) (c.Request, error) {
				return genfixture.Request(reg, input, to, "generated-go", ggtarget.Provider{Package: "capture"}, fixture.libraries...)
			})
			tree := outputTree[gs.Language](t, out)
			if text := formatted(t, golangText(), tree); fixture.captured != strings.Contains(text, "type string struct") {
				t.Fatalf("the fixture did not generate the expected declarations:\n%s", text)
			}
			// The type of the field Name: a named type string without a reference.
			var use ast.NodeRef
			for _, node := range tree.Nodes {
				if node.Kind != gs.KindField || nodeAttribute(node, gs.AttributeName) != "Name" {
					continue
				}
				for _, child := range node.Children {
					if child.Role == gs.RoleType {
						use = ast.NodeRef{Tree: tree.ID, Node: child.Node}
					}
				}
			}
			if use.Node == 0 || tree.Nodes[use.Node-1].Kind != gs.KindNamedType || nodeAttribute(tree.Nodes[use.Node-1], gs.AttributeName) != "string" || len(tree.Nodes[use.Node-1].References) != 0 {
				t.Fatalf("the field Name has no unlinked type string: %v", use)
			}
			// What generation recorded.
			recorded, ok := externalUsesBySite(t, out)[use]
			if !ok {
				t.Fatal("no destination is recorded at the type of the field Name")
			}
			symbol, err := gosymbol.ReadExternal(out.Snapshot().Language(), recorded)
			if err != nil {
				t.Fatal(err)
			}
			if name, ok := symbol.Predeclared(); !ok || name != "string" {
				t.Fatalf("the recorded destination is %q", name)
			}
			// What the compiler established for the same tree.
			_, result := goResults(t, out.Snapshot())
			model, err := gn.Open(result)
			if err != nil {
				t.Fatalf("the generated output has no Go semantics: %v; %+v", err, result.Diagnostics())
			}
			data := model.Graph().Data()
			target, ok := bindingAt(data, use, gn.NameRole())
			if !ok {
				t.Fatal("the compiler published no name binding at the type of the field Name")
			}
			entity := data.Entities[target-1]
			if fixture.captured {
				if entity.Kind != gn.DefinedKind() || len(entity.Declarations) != 1 {
					t.Fatalf("the use binds to %+v, want the helper, an entity of the kind defined", entity)
				}
				if goPredeclaredEntity(symbol, entity) {
					t.Fatal("the captured use is taken for its recorded destination")
				}
				return
			}
			if !goPredeclaredEntity(symbol, entity) {
				t.Fatalf("the use binds to %+v, which is not the recorded predeclared destination", entity)
			}
		})
	}
}

// A supplied TypeScript package exports Cell as an alias of another exported
// type. The recorded destination and the published binding both name the
// export Cell of that module, and neither names the type underneath.
func TestTypeScriptAliasExportIsTheDestination(t *testing.T) {
	const module = "@bitspark/cell"
	output := generatedCellPage(t)
	uses := output.ExternalUses()
	if len(uses) != 1 {
		t.Fatal("expected one external use", uses)
	}
	recorded, err := tssymbol.ReadExternal(output.Snapshot().Language(), uses[0])
	if err != nil {
		t.Fatal(err)
	}
	if recorded != (tssymbol.Type{Module: module, Name: "Cell"}) {
		t.Fatalf("the recorded destination is %+v", recorded)
	}
	tree := outputTree[ts.Language](t, output)
	placement := tsfiles.Placement{File: "types", Path: "types.ts"}
	for _, child := range tree.Nodes[tree.Root-1].Children {
		placement.Declarations = append(placement.Declarations, ast.NodeRef{Tree: tree.ID, Node: child.Node})
	}
	setting, err := (tsfiles.Organization{Entry: placement.File, Units: []tsfiles.Placement{placement}}).Setting()
	if err != nil {
		t.Fatal(err)
	}
	formatter, err := tsfiles.NewGenerated(placementStages(t, tn.PackageTypes{Module: module, Declaration: "export interface Inner<T> { value: T }\nexport type Cell<T> = Inner<T>;\n"}))
	if err != nil {
		t.Fatal(err)
	}
	var placed source.Source[ts.Language]
	placed, err = formatter.Format(context.Background(), format.GeneratedFilesInput{Output: output, Configuration: []semantic.Setting{setting}})
	if err != nil {
		t.Fatal("the placement of an aliased export was refused: ", err)
	}
	// The name bindings of the placed output that leave it: each names an
	// export of a package.
	data := placed.SemanticsView().Data()
	var published []tssymbol.Type
	for _, binding := range data.Bindings {
		entity := data.Entities[binding.Target-1]
		if binding.At.Role != tn.NameRole() || entity.Kind != tn.PackageTypeKind() {
			continue
		}
		if entity.Owner == 0 || data.Entities[entity.Owner-1].Kind != tn.PackageKind() {
			t.Fatalf("the package type %s has no package owner", entity.Name)
		}
		published = append(published, tssymbol.Type{Module: data.Entities[entity.Owner-1].Name, Name: entity.Name})
	}
	if !slices.Equal(published, []tssymbol.Type{recorded}) {
		t.Fatalf("the published bindings name %+v, and the recorded destination is %+v", published, recorded)
	}
}

// The Rust model's external form names its crate apart from the path. The
// model refuses a malformed form as invalid. The Rust target refuses a crate
// that no destination can name: a text that is no crate name as invalid, and
// the crate of a library as unsupported. The Go model's external form stays
// unsupported.
func TestExternalFormsNameADependencyThatADestinationCanName(t *testing.T) {
	ctx := context.Background()
	text := rustmodel.Primitive(rustmodel.PrimitiveString)
	rust := rgtarget.Provider{}
	for name, c := range map[string]struct {
		form   g.Expression[rs.Language]
		model  error
		target error
	}{
		"an item of the standard library": {rustmodel.External("std", []string{"option", "Option"}, text), nil, nil},
		"a nested item":                   {rustmodel.External("std", []string{"collections", "BTreeMap"}, text, text), nil, nil},
		"no path":                         {rustmodel.External("std", nil), g.ErrInvalid, g.ErrInvalid},
		"an empty segment":                {rustmodel.External("std", []string{"option", ""}), g.ErrInvalid, g.ErrInvalid},
		"no crate":                        {rustmodel.External("", []string{"Option"}), g.ErrInvalid, g.ErrInvalid},
		"the crate joined with the path":  {rustmodel.External("std::option", []string{"Option"}), nil, g.ErrInvalid},
		"a path joined in one segment":    {rustmodel.External("std", []string{"option::Option"}), nil, g.ErrInvalid},
		"a keyword as a segment":          {rustmodel.External("std", []string{"option", "fn"}), nil, g.ErrInvalid},
		"the crate of a library":          {rustmodel.External("serde_json", []string{"Value"}), nil, g.ErrUnsupported},
		"another crate of the sysroot":    {rustmodel.External("core", []string{"option", "Option"}), nil, g.ErrUnsupported},
	} {
		t.Run(name, func(t *testing.T) {
			if err := rustmodel.CheckExpression(ctx, c.form); !errors.Is(err, c.model) || (c.model == nil && err != nil) {
				t.Errorf("the model answers %v, want %v", err, c.model)
			}
			if err := rust.CheckExpression(ctx, c.form); !errors.Is(err, c.target) || (c.target == nil && err != nil) {
				t.Errorf("the target answers %v, want %v", err, c.target)
			}
		})
	}
	if err := (ggtarget.Provider{Package: "p"}).CheckExpression(ctx, gomodel.External("fmt", "Stringer")); !errors.Is(err, g.ErrUnsupported) {
		t.Errorf("the Go external form answers %v, want unsupported", err)
	}
}
