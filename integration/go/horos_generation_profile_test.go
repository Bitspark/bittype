package integration_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	"github.com/Bitspark/bittype/integration/go/testdata/readonlyprofile"
	cs "github.com/Bitspark/bittype/syntax/cpp/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	rs "github.com/Bitspark/bittype/syntax/rust/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	cgtarget "github.com/Bitspark/bittype/transform/syntax/cpp/target/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	rgtarget "github.com/Bitspark/bittype/transform/syntax/rust/target/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	tgtarget "github.com/Bitspark/bittype/transform/syntax/typescript/target/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// unregisteredProfile belongs to this file. Nothing is ever registered under
// it.
const unregisteredProfile = "integration-test/unregistered/1"

// builtinPair is one builtin horos-to-native pair with its target language
// hidden. name is the pair's directory, transform/syntax/horos-to-<name>.
// compose registers the source, the target and the pair's builtin libraries in
// a fresh registry for the given modules. Every builtin library states the
// structural representation profile itself, so compose takes no profile.
type builtinPair struct {
	name    string
	check   func(t *testing.T)
	compose func(t *testing.T, modules ...core.Value) *profileComposition
}

// profileComposition is one registry with a request for one root module. plan
// and generate send that request under the profile they are given.
type profileComposition struct {
	registry *c.Registry
	request  c.Request
	checker  check.Checker
	text     func(t *testing.T, out g.Output) string
}

func (q *profileComposition) under(profile string) c.Request {
	request := q.request
	request.Profile = profile
	return request
}
func (q *profileComposition) plan(profile string) (g.GenerationPlan, error) {
	return q.registry.PlanGeneration(context.Background(), q.under(profile), q.checker)
}
func (q *profileComposition) generate(profile string) (g.Output, error) {
	return q.registry.Generate(context.Background(), q.under(profile), q.checker)
}

// refuses reports whether a request under the profile fails with
// ErrUnsupported and yields neither a plan nor an output.
func (q *profileComposition) refuses(profile string) error {
	if p, err := q.plan(profile); !errors.Is(err, g.ErrUnsupported) || p.Valid() {
		return errors.Join(errors.New("the request planned or failed with another error"), err)
	}
	if out, err := q.generate(profile); !errors.Is(err, g.ErrUnsupported) || out.Snapshot().Valid() {
		return errors.Join(errors.New("the request generated or failed with another error"), err)
	}
	return nil
}

// composeProfiles registers the libraries for the root module modules[0]. The
// other modules are the ones it imports.
func composeProfiles[L any](t *testing.T, modules []core.Value, p nativeProvider[L], target g.Target[L], libraries ...g.MappingLibrary[h.Language, L, source.Reader]) *profileComposition {
	t.Helper()
	_, _, input := generationInput(t, modules[0], modules[1:]...)
	checker := nativeChecker(t, p)
	q := &profileComposition{registry: &c.Registry{}, checker: checker, text: func(t *testing.T, out g.Output) string {
		t.Helper()
		return formatted(t, p, outputTree[L](t, out))
	}}
	var err error
	if q.request, err = genfixture.Request(q.registry, input, checker.Language(), "generated-profile", target, libraries...); err != nil {
		t.Fatal(err)
	}
	return q
}

// directoryNames returns the sorted names of the directories in a directory of
// this repository.
func directoryNames(t *testing.T, path ...string) []string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(append([]string{"..", ".."}, path...)...))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	slices.Sort(names)
	return names
}

func newBuiltinPair[L any](name string, p nativeProvider[L], target g.Target[L], libraries func() []g.MappingLibrary[h.Language, L, source.Reader], constructors func() map[string]any, bundles func() map[string]genfixture.Bundle[L]) builtinPair {
	return builtinPair{
		name: name,
		compose: func(t *testing.T, modules ...core.Value) *profileComposition {
			t.Helper()
			return composeProfiles(t, modules, p, target, libraries()...)
		},
		check: func(t *testing.T) {
			from, err := representation.NewLanguage[h.Language](h.Language{}, h.Profile)
			if err != nil {
				t.Fatal(err)
			}
			to := nativeChecker(t, p).Language()
			// The constructors named in libraries_<pair>_test.go are those of
			// exactly the mapping packages that the pair has, by directory name.
			var named []string
			for name := range constructors() {
				named = append(named, name)
			}
			slices.Sort(named)
			// The contributions named in bundles_<pair>_test.go are those of
			// exactly the mapping packages that the pair has, by directory
			// name, and each has a vocabulary package and a codec package of
			// the same name. Each states the structural profile and this
			// target, and its codec is of the schema that it states.
			var contributions []string
			for name := range bundles() {
				contributions = append(contributions, name)
			}
			slices.Sort(contributions)
			for _, part := range []string{"mapping", "vocabulary", "codec"} {
				if packages := directoryNames(t, "transform", "syntax", "horos-to-"+name, part); len(packages) == 0 || !slices.Equal(contributions, packages) {
					t.Fatalf("the pair has the %s packages %v, and the test names the contributions %v", part, packages, contributions)
				}
			}
			for _, name := range contributions {
				bundle := bundles()[name]
				contribution, err := bundle.Contribution(from)
				if err != nil {
					t.Errorf("%s: %v", name, err)
					continue
				}
				codec, err := bundle.Codec()
				if err != nil {
					t.Errorf("%s: %v", name, err)
					continue
				}
				if contribution.Profile != rules.RepresentationProfile || contribution.Target != string(to.ID()) || !contribution.Reference.Valid() || contribution.Rule.Decide == nil {
					t.Errorf("%s states the profile %q, the target %q and the reference %v", name, contribution.Profile, contribution.Target, contribution.Reference)
				}
				if codec.Schema() != contribution.Schema {
					t.Errorf("%s states the schema %v and its codec is of %v", name, contribution.Schema, codec.Schema())
				}
			}
			// Until the mapping contract is removed, the libraries of the old
			// shape are checked as before, for the packages that still have one.
			for _, name := range named {
				if !slices.Contains(contributions, name) {
					t.Fatalf("the test names the constructor of %s, and the pair has no such mapping package", name)
				}
			}
			if len(libraries()) != len(named) {
				t.Errorf("%d libraries are registered by the tests, and the pair has %d mapping packages", len(libraries()), len(named))
			}
			// The evidence: build each library through its exported
			// constructor, describe it, and compare the profile it states.
			key := reflect.ValueOf(g.NewHelperKey("integration-test/helper"))
			for _, name := range named {
				constructor := reflect.ValueOf(constructors()[name])
				signature := constructor.Type()
				if signature.Kind() != reflect.Func || signature.NumOut() != 1 {
					t.Errorf("%s: Library is no constructor of one library", name)
					continue
				}
				// This test can build a helper key and nothing else. A
				// constructor with any other parameter, also a variadic one,
				// fails here until someone has looked at what it lets its
				// caller choose.
				var arguments []reflect.Value
				for i := 0; i < signature.NumIn(); i++ {
					if signature.IsVariadic() || signature.In(i) != key.Type() {
						t.Errorf("%s: Library takes a %s. Decide what its caller can choose with it before this test builds one", name, signature.In(i))
						arguments = nil
						break
					}
					arguments = append(arguments, key)
				}
				if len(arguments) != signature.NumIn() {
					continue
				}
				library, ok := constructor.Call(arguments)[0].Interface().(g.MappingLibrary[h.Language, L, source.Reader])
				if !ok {
					t.Errorf("%s: Library returns no mapping library of this pair", name)
					continue
				}
				description, err := library.Describe(from, to)
				if err != nil {
					t.Errorf("%s: %v", name, err)
					continue
				}
				if description.Profile != rules.RepresentationProfile {
					t.Errorf("%s states the profile %q", name, description.Profile)
				}
				if len(description.Bindings)+len(description.Declarations)+len(description.Helpers) == 0 {
					t.Errorf("%s describes nothing", name)
				}
			}
		},
	}
}

func builtinPairs() []builtinPair {
	return []builtinPair{
		newBuiltinPair[ts.Language]("typescript", typescriptText(), tgtarget.Provider{}, typescriptLibraries, typescriptConstructors, typescriptBundleTable),
		newBuiltinPair[gs.Language]("golang", golangText(), ggtarget.Provider{Package: "pagination"}, golangLibraries, golangConstructors, golangBundleTable),
		newBuiltinPair[rs.Language]("rust", rustText(), rgtarget.Provider{}, rustLibraries, rustConstructors, rustBundleTable),
		newBuiltinPair[cs.Language]("cpp", cppText(), cgtarget.Provider{Namespace: "pagination"}, cppLibraries, cppConstructors, cppBundleTable),
	}
}

// TestBuiltinMappingsStateTheStructuralProfile: a profile is a promise about
// the native type, so the builtin mapping that makes the promise states it.
//
// The evidence is the comparison: every library that the exported constructor
// Library of a builtin mapping package returns describes itself under the
// structural representation profile. The pairs and the packages are read from
// the directories under transform/syntax, so a new pair or a new mapping
// package fails this test until it is named here.
//
// Two things around the comparison are tripwires and no proof. The test fails
// for a constructor parameter that is not a helper key, because it cannot
// build one. And the calls in the libraries_*_test.go files stop compiling
// when a constructor gains a required parameter. Neither sees a second
// exported constructor in an existing mapping package, a package variable or
// a setter. Those are for a reader of the package to see.
func TestBuiltinMappingsStateTheStructuralProfile(t *testing.T) {
	var directories, pairs []string
	for _, name := range directoryNames(t, "transform", "syntax") {
		if pair, found := strings.CutPrefix(name, "horos-to-"); found {
			directories = append(directories, pair)
		}
	}
	for _, pair := range builtinPairs() {
		pairs = append(pairs, pair.name)
	}
	slices.Sort(pairs)
	if len(directories) == 0 || !slices.Equal(pairs, directories) {
		t.Fatalf("transform/syntax has the horos pairs %v, and this test covers %v", directories, pairs)
	}
	for _, pair := range builtinPairs() {
		t.Run(pair.name, pair.check)
	}
}

// TestEnforcementProfileIsRefused: the enforcement profile is declared and not
// implemented (GEN-02). No builtin mapping states it, so the registry finds no
// mapping for a request that names it. The request fails with ErrUnsupported
// and yields no plan and no output, from the same registry that generates
// under the structural profile.
func TestEnforcementProfileIsRefused(t *testing.T) {
	for _, pair := range builtinPairs() {
		t.Run(pair.name, func(t *testing.T) {
			q := pair.compose(t, paginationModule())
			out, err := q.generate(rules.RepresentationProfile)
			if err != nil || q.text(t, out) == "" {
				t.Fatal("the structural profile did not generate:", err)
			}
			if err := q.refuses(rules.EnforcementProfile); err != nil {
				t.Fatal("the enforcement profile:", err)
			}
		})
	}
}

// TestUnregisteredProfileIsRefused: a request for a profile under which nothing
// is registered fails with ErrUnsupported, whatever else the registry holds.
func TestUnregisteredProfileIsRefused(t *testing.T) {
	empty := composeProfiles[ts.Language](t, []core.Value{paginationModule()}, typescriptText(), tgtarget.Provider{})
	for _, profile := range []string{unregisteredProfile, rules.RepresentationProfile, rules.EnforcementProfile, readonlyprofile.Profile} {
		if err := empty.refuses(profile); err != nil {
			t.Errorf("an empty registry, profile %q: %v", profile, err)
		}
	}
	for _, pair := range builtinPairs() {
		t.Run(pair.name, func(t *testing.T) {
			if err := pair.compose(t, paginationModule()).refuses(unregisteredProfile); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// listsModules returns a root module and the module it imports. Together they
// have lists, a generic declaration and an imported declaration, and no
// record. The builtin structural family and the read-only fixture family both
// map all of it.
func listsModules(t *testing.T) []core.Value {
	t.Helper()
	list := func(element core.Value) core.Value { return genTuple("horos/list-of", element) }
	shared := genModule("example.shared", genMap(), genEntry("Names", genType(list(genScalar("utf8-text")))))
	pin, err := decl.Pin(shared)
	if err != nil {
		t.Fatal(err)
	}
	imported := genTuple("horos/decl/import", data.MustEncodeText("example.shared"), core.NewAtom(pin[:]))
	root := genModule("example.lists", genMap(genEntry("shared", imported)),
		genEntry("Label", genType(genScalar("utf8-text"))),
		genEntry("Flags", genType(list(genScalar("bool")))),
		genEntry("Table", genType(list(genImported("shared", "Names")))),
		genEntry("Batch", genType(list(genParam(0)), "T")),
	)
	return []core.Value{root, shared}
}

// TestRequestSelectsTheContributionsOfItsProfile: a second profile enters as a
// contribution of its own. The read-only fixture library lives in its own
// package, states its own profile and maps a list to another native form than
// the builtin family. It is a mapping.Library function value, as the builtin
// mappings are, and it uses the common declaration recipe. One registry holds
// both families. Each request receives the contributions that state its
// profile and nothing of the other family. The horos mapping package and the
// registry know nothing of the fixture.
func TestRequestSelectsTheContributionsOfItsProfile(t *testing.T) {
	if readonlyprofile.Profile == rules.RepresentationProfile || readonlyprofile.Profile == rules.EnforcementProfile {
		t.Fatal("the fixture contribution states a profile of the builtin family:", readonlyprofile.Profile)
	}
	// The assignment compiles only while the fixture is given as the function
	// type of the horos mapping package and not as a type of its own.
	var fixture rules.Library[ts.Language] = readonlyprofile.Library()
	both := func(t *testing.T, modules ...core.Value) *profileComposition {
		t.Helper()
		return composeProfiles(t, modules, typescriptText(), tgtarget.Provider{}, append(typescriptLibraries(), fixture)...)
	}
	bindings := func(t *testing.T, q *profileComposition, profile string) map[string]bool {
		t.Helper()
		plan, err := q.plan(profile)
		if err != nil {
			t.Fatal(profile, err)
		}
		// A step names its binding as "<library>/<binding ID>". The library
		// part is the registry's own name for a registration.
		selected := map[string]bool{}
		for _, step := range plan.Steps() {
			if _, id, found := strings.Cut(step.Binding, "/"); found {
				selected[id] = true
			}
		}
		return selected
	}
	generate := func(t *testing.T, q *profileComposition, profile string) string {
		t.Helper()
		out, err := q.generate(profile)
		if err != nil {
			t.Fatal(profile, err)
		}
		return q.text(t, out)
	}

	t.Run("each request receives its own family", func(t *testing.T) {
		q := both(t, listsModules(t)...)
		structural, readonly := generate(t, q, rules.RepresentationProfile), generate(t, q, readonlyprofile.Profile)
		for _, want := range []string{"Label = string", "Flags = boolean[]", "_Names = string[]", "_Names[]", "Batch<T> = T[]"} {
			if !strings.Contains(structural, want) {
				t.Errorf("the structural request lacks %q:\n%s", want, structural)
			}
		}
		for _, want := range []string{"Label = string", "Flags = ReadonlyArray<boolean>", "_Names = ReadonlyArray<string>", "Table = ReadonlyArray<M", "Batch<T> = ReadonlyArray<T>"} {
			if !strings.Contains(readonly, want) {
				t.Errorf("the read-only request lacks %q:\n%s", want, readonly)
			}
		}
		if strings.Contains(structural, "ReadonlyArray") || strings.Contains(readonly, "[]") {
			t.Errorf("a request received the other family's list:\n%s\n%s", structural, readonly)
		}
		// The inspectable plan names the selected bindings.
		if selected := bindings(t, q, rules.RepresentationProfile); !selected["list/typescript/1"] || selected[readonlyprofile.ListBinding] {
			t.Error("the structural plan selected", selected)
		}
		if selected := bindings(t, q, readonlyprofile.Profile); !selected[readonlyprofile.ListBinding] || selected["list/typescript/1"] {
			t.Error("the read-only plan selected", selected)
		}
		if err := q.refuses(rules.EnforcementProfile); err != nil {
			t.Error("the enforcement profile, with two families registered:", err)
		}
	})

	t.Run("a family is not completed from another profile", func(t *testing.T) {
		// The fixture family maps no record. The builtin record library sits
		// in the same registry and states another profile, so it does not
		// answer.
		q := both(t, paginationModule())
		if text := generate(t, q, rules.RepresentationProfile); !strings.Contains(text, "interface Page<T>") {
			t.Fatal("the structural request lost its records:\n" + text)
		}
		out, err := q.generate(readonlyprofile.Profile)
		var diagnostic *text.Diagnostic
		if !errors.As(err, &diagnostic) || diagnostic.Code != text.Unsupported || out.Snapshot().Valid() {
			t.Fatal("the read-only request was completed by the structural family:", err)
		}
	})

	t.Run("the fixture alone", func(t *testing.T) {
		q := composeProfiles(t, listsModules(t), typescriptText(), tgtarget.Provider{}, fixture)
		if text := generate(t, q, readonlyprofile.Profile); !strings.Contains(text, "Table = ReadonlyArray<M") {
			t.Fatal(text)
		}
		for _, profile := range []string{rules.RepresentationProfile, rules.EnforcementProfile} {
			if err := q.refuses(profile); err != nil {
				t.Errorf("a registry that holds only the fixture, profile %q: %v", profile, err)
			}
		}
	})
}

// republished is a contribution of its own author: it takes another library's
// description and states a profile of its own choice.
type republished struct {
	library g.MappingLibrary[h.Language, ts.Language, source.Reader]
	profile string
}

func (p republished) Describe(from, to representation.Language) (g.MappingDescription[h.Language, ts.Language, source.Reader], error) {
	description, err := p.library.Describe(from, to)
	description.Profile = p.profile
	return description, err
}

// TestRegistryDoesNotJudgeAStatedProfile records where the refusal of the
// enforcement profile ends. The registry selects by the profile that a
// contribution states and does not judge whether the contribution keeps that
// promise. A contribution that republishes the builtin, constraint-erasing
// descriptions under the enforcement identity is selected by a request for it
// and yields representation output. Its author states that profile and
// answers for it. bittype states it nowhere, and no constructor of a builtin
// mapping takes it (bittype/ACX-008, limits).
func TestRegistryDoesNotJudgeAStatedProfile(t *testing.T) {
	var libraries []g.MappingLibrary[h.Language, ts.Language, source.Reader]
	for _, library := range typescriptLibraries() {
		libraries = append(libraries, republished{library: library, profile: rules.EnforcementProfile})
	}
	q := composeProfiles(t, []core.Value{paginationModule()}, typescriptText(), tgtarget.Provider{}, libraries...)
	out, err := q.generate(rules.EnforcementProfile)
	if err != nil {
		t.Fatal("the registry refused a stated profile:", err)
	}
	if text := q.text(t, out); !strings.Contains(text, "interface Page<T>") {
		t.Fatal(text)
	}
	// The builtin libraries themselves were not registered, so nothing states
	// the structural profile in this registry.
	if err := q.refuses(rules.RepresentationProfile); err != nil {
		t.Fatal(err)
	}
}

// TestRegistrationNeedsAStatedProfile: a library that states no profile, and a
// library without a description function, cannot be registered. Both fail
// with ErrInvalid: the generation contract refuses a description without a
// profile, and the function type refuses a nil function.
func TestRegistrationNeedsAStatedProfile(t *testing.T) {
	from, err := representation.NewLanguage[h.Language](h.Language{}, h.Profile)
	if err != nil {
		t.Fatal(err)
	}
	to := nativeChecker(t, typescriptText()).Language()
	var registry c.Registry
	unstated := rules.Library[ts.Language](func(from, to representation.Language) (g.MappingDescription[h.Language, ts.Language, source.Reader], error) {
		description, err := typescriptLibraries()[0].Describe(from, to)
		description.Profile = ""
		return description, err
	})
	if _, err := c.RegisterMapping(&registry, from, to, unstated); !errors.Is(err, g.ErrInvalid) {
		t.Error("a library that states no profile:", err)
	}
	if _, err := c.RegisterMapping(&registry, from, to, rules.Library[ts.Language](nil)); !errors.Is(err, g.ErrInvalid) {
		t.Error("a library without a description function:", err)
	}
}
