package integration_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
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
	core "github.com/bitspark/ontos/core/go"
)

// unregisteredProfile belongs to this file. Nothing is ever registered under
// it.
const unregisteredProfile = "integration-test/unregistered/1"

// builtinPair is one builtin horos-to-native pair with its target language
// hidden. compose registers the source, the target and the pair's builtin
// libraries in a fresh registry for the given module. Every builtin library
// states the structural representation profile itself, so compose takes no
// profile.
type builtinPair struct {
	name    string
	check   func(t *testing.T)
	compose func(t *testing.T, module core.Value) *profileComposition
}

// profileComposition is one registry with a request for one module. plan and generate
// send that request under the profile they are given.
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

func composeProfiles[L any](t *testing.T, module core.Value, p nativeProvider[L], target g.Target[L], libraries ...g.MappingLibrary[h.Language, L, source.Reader]) *profileComposition {
	t.Helper()
	_, _, input := generationInput(t, module)
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

func newBuiltinPair[L any](name string, p nativeProvider[L], target g.Target[L], libraries func() []g.MappingLibrary[h.Language, L, source.Reader], constructors func() []any) builtinPair {
	return builtinPair{
		name: name,
		compose: func(t *testing.T, module core.Value) *profileComposition {
			t.Helper()
			return composeProfiles(t, module, p, target, libraries()...)
		},
		check: func(t *testing.T) {
			from, err := representation.NewLanguage[h.Language](h.Language{}, h.Profile)
			if err != nil {
				t.Fatal(err)
			}
			to := nativeChecker(t, p).Language()
			// The lists of this package cover every mapping package of the pair.
			directories, err := os.ReadDir(filepath.Join("..", "..", "transform", "syntax", "horos-to-"+name, "mapping"))
			if err != nil {
				t.Fatal(err)
			}
			if len(directories) == 0 || len(libraries()) != len(directories) || len(constructors()) != len(directories) {
				t.Fatalf("%d mapping packages, %d libraries and %d constructors are listed", len(directories), len(libraries()), len(constructors()))
			}
			// No constructor takes a profile. A profile is a string identity,
			// so no constructor takes a string.
			for i, constructor := range constructors() {
				signature := reflect.TypeOf(constructor)
				if signature.Kind() != reflect.Func {
					t.Fatalf("constructor %d is no function", i)
				}
				for j := 0; j < signature.NumIn(); j++ {
					parameter := signature.In(j)
					if signature.IsVariadic() && j == signature.NumIn()-1 {
						parameter = parameter.Elem()
					}
					if parameter.Kind() == reflect.String {
						t.Errorf("constructor %d takes a string: its caller could choose the profile", i)
					}
				}
			}
			// Every library states the structural representation profile.
			for i, library := range libraries() {
				description, err := library.Describe(from, to)
				if err != nil {
					t.Errorf("library %d: %v", i, err)
					continue
				}
				if description.Profile != rules.RepresentationProfile {
					t.Errorf("library %d states the profile %q", i, description.Profile)
				}
				if len(description.Bindings)+len(description.Declarations)+len(description.Helpers) == 0 {
					t.Errorf("library %d describes nothing", i)
				}
			}
		},
	}
}

func builtinPairs() []builtinPair {
	return []builtinPair{
		newBuiltinPair[ts.Language]("typescript", typescriptText(), tgtarget.Provider{}, typescriptLibraries, typescriptConstructors),
		newBuiltinPair[gs.Language]("golang", golangText(), ggtarget.Provider{Package: "pagination"}, golangLibraries, golangConstructors),
		newBuiltinPair[rs.Language]("rust", rustText(), rgtarget.Provider{}, rustLibraries, rustConstructors),
		newBuiltinPair[cs.Language]("cpp", cppText(), cgtarget.Provider{Namespace: "pagination"}, cppLibraries, cppConstructors),
	}
}

// TestBuiltinMappingsStateTheStructuralProfile: a profile is a promise about
// the native type, so the builtin mapping that makes the promise states it.
// No exported constructor of a builtin mapping package lets its caller choose
// the profile, and every builtin library describes itself under the structural
// representation profile. The calls in the libraries_*_test.go files compile
// only while a constructor takes nothing but its helper key.
func TestBuiltinMappingsStateTheStructuralProfile(t *testing.T) {
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
	empty := composeProfiles[ts.Language](t, paginationModule(), typescriptText(), tgtarget.Provider{})
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

// listsModule has no record and no generic declaration, so both the builtin
// structural family and the read-only fixture family map all of it.
func listsModule() core.Value {
	list := func(element core.Value) core.Value { return genTuple("horos/list-of", element) }
	return genModule("example.lists", genMap(),
		genEntry("Label", genType(genScalar("utf8-text"))),
		genEntry("Names", genType(list(genScalar("utf8-text")))),
		genEntry("Flags", genType(list(genScalar("bool")))),
		genEntry("Table", genType(list(genLocal("Names")))),
	)
}

// TestRequestSelectsTheContributionsOfItsProfile: a second profile enters as a
// contribution of its own. The read-only fixture library lives in its own
// package, states its own profile and maps a list to another native form than
// the builtin family. One registry holds both families. Each request receives
// the contributions that state its profile and nothing of the other family.
// The horos mapping package and the registry know nothing of the fixture.
func TestRequestSelectsTheContributionsOfItsProfile(t *testing.T) {
	if readonlyprofile.Profile == rules.RepresentationProfile || readonlyprofile.Profile == rules.EnforcementProfile {
		t.Fatal("the fixture contribution states a profile of the builtin family:", readonlyprofile.Profile)
	}
	both := func(t *testing.T, module core.Value) *profileComposition {
		t.Helper()
		return composeProfiles(t, module, typescriptText(), tgtarget.Provider{}, append(typescriptLibraries(), readonlyprofile.Library())...)
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
		q := both(t, listsModule())
		structural, readonly := generate(t, q, rules.RepresentationProfile), generate(t, q, readonlyprofile.Profile)
		for _, want := range []string{"Names = string[]", "Flags = boolean[]", "Table = Names[]", "Label = string"} {
			if !strings.Contains(structural, want) {
				t.Errorf("the structural request lacks %q:\n%s", want, structural)
			}
		}
		for _, want := range []string{"Names = ReadonlyArray<string>", "Flags = ReadonlyArray<boolean>", "Table = ReadonlyArray<Names>", "Label = string"} {
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
		// The fixture family maps no record and no generic declaration. The
		// builtin record and declaration libraries sit in the same registry
		// and state another profile, so they do not answer.
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
		q := composeProfiles(t, listsModule(), typescriptText(), tgtarget.Provider{}, readonlyprofile.Library())
		if text := generate(t, q, readonlyprofile.Profile); !strings.Contains(text, "Table = ReadonlyArray<Names>") {
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
	q := composeProfiles(t, paginationModule(), typescriptText(), tgtarget.Provider{}, libraries...)
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
