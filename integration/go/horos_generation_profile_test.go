package integration_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
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
// compose registers the source, the target and the pair's builtin
// contributions in a fresh registry for the given modules. Every builtin
// contribution states the structural representation profile itself, so compose
// takes no profile.
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

// missingProfile reports whether err is the kernel's refusal of a profile that
// no registered contribution states: a missing contribution for that profile,
// and no other refusal value. In particular it is not unsupported.
func missingProfile(err error, profile string) error {
	var refusal *g.Refusal
	if !errors.Is(err, g.ErrMissingContribution) || !errors.As(err, &refusal) {
		return errors.Join(errors.New("the request did not fail as a missing contribution"), err)
	}
	if refusal.Class != g.RefusalMissingContribution || refusal.Owner != g.KernelOwner || refusal.Profile != profile || !strings.HasPrefix(refusal.Reason, g.ReasonNoProfile) {
		return errors.Join(errors.New("the refusal does not name the kernel and the profile that nobody states"), err)
	}
	for _, other := range []error{g.ErrUnsupported, g.ErrAmbiguous, g.ErrRefusedByDesign, g.ErrMissingFacet, g.ErrLimit, g.ErrScope, g.ErrInvalid} {
		if errors.Is(err, other) {
			return errors.Join(errors.New("the refusal also matches "+other.Error()), err)
		}
	}
	return nil
}

// refuses reports whether a request under the profile fails as a missing
// contribution and yields neither a plan nor an output.
func (q *profileComposition) refuses(profile string) error {
	p, err := q.plan(profile)
	if err = missingProfile(err, profile); err != nil || p.Valid() {
		return errors.Join(errors.New("the request planned or failed with another error"), err)
	}
	out, err := q.generate(profile)
	if err = missingProfile(err, profile); err != nil || out.Snapshot().Valid() {
		return errors.Join(errors.New("the request generated or failed with another error"), err)
	}
	return nil
}

// composeProfiles registers the bundles for the root module modules[0]. The
// other modules are the ones it imports.
func composeProfiles[L any](t *testing.T, modules []core.Value, p nativeProvider[L], target g.Target[L], bundles ...genfixture.Bundle[L]) *profileComposition {
	t.Helper()
	_, _, input := generationInput(t, modules[0], modules[1:]...)
	checker := nativeChecker(t, p)
	q := &profileComposition{registry: &c.Registry{}, checker: checker, text: func(t *testing.T, out g.Output) string {
		t.Helper()
		return formatted(t, p, outputTree[L](t, out))
	}}
	var err error
	if q.request, err = genfixture.Compose(q.registry, input, checker.Language(), "generated-profile", target, bundles...); err != nil {
		t.Fatal(err)
	}
	return q
}

// withReadonlyFamily registers the read-only fixture family in the registry of
// a TypeScript composition: its contributions, and the codecs of the payload
// schemas that they state.
func withReadonlyFamily(t *testing.T, q *profileComposition) *profileComposition {
	t.Helper()
	from, to := q.request.Input.Source.Language(), q.checker.Language()
	family, err := readonlyprofile.Family(from)
	if err != nil {
		t.Fatal(err)
	}
	// The assignment compiles only while the fixture gives its contributions
	// as the type that the builtin mapping packages return and not as a type
	// of its own.
	var contributions []g.Contribution[h.Language, ts.Language, source.Reader] = family.Contributions
	for _, contribution := range contributions {
		if contribution.Profile != readonlyprofile.Profile {
			t.Fatalf("%s states the profile %q", contribution.Reference, contribution.Profile)
		}
		if err := c.RegisterContribution(q.registry, from, to, contribution); err != nil {
			t.Fatal(err)
		}
	}
	for _, codec := range family.Codecs {
		if err := q.registry.RegisterCodec(codec); err != nil {
			t.Fatal(err)
		}
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

func newBuiltinPair[L any](name string, p nativeProvider[L], target g.Target[L], set func() []genfixture.Bundle[L], bundles func() map[string]genfixture.Bundle[L]) builtinPair {
	return builtinPair{
		name: name,
		compose: func(t *testing.T, modules ...core.Value) *profileComposition {
			t.Helper()
			return composeProfiles(t, modules, p, target, set()...)
		},
		check: func(t *testing.T) {
			from, err := representation.NewLanguage[h.Language](h.Language{}, h.Profile)
			if err != nil {
				t.Fatal(err)
			}
			to := nativeChecker(t, p).Language()
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
			// The set that the tests register is the whole table.
			if len(set()) != len(contributions) {
				t.Errorf("%d contributions are registered by the tests, and the pair has %d mapping packages", len(set()), len(contributions))
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
		},
	}
}

func builtinPairs() []builtinPair {
	return []builtinPair{
		newBuiltinPair[ts.Language]("typescript", typescriptText(), tgtarget.Provider{}, typescriptBundles, typescriptBundleTable),
		newBuiltinPair[gs.Language]("golang", golangText(), ggtarget.Provider{Package: "pagination"}, golangBundles, golangBundleTable),
		newBuiltinPair[rs.Language]("rust", rustText(), rgtarget.Provider{}, rustBundles, rustBundleTable),
		newBuiltinPair[cs.Language]("cpp", cppText(), cgtarget.Provider{Namespace: "pagination"}, cppBundles, cppBundleTable),
	}
}

// TestBuiltinMappingsStateTheStructuralProfile: a profile is a promise about
// the native type, so the builtin mapping that makes the promise states it.
//
// The evidence is the comparison: every contribution that the exported
// constructor Contribution of a builtin mapping package returns states the
// structural representation profile and the target of its pair, and the codec
// beside it is of the schema that it states. The pairs and the packages are
// read from the directories under transform/syntax, so a new pair or a new
// mapping package fails this test until it is named here.
//
// One thing around the comparison is a tripwire and no proof. A bundle holds
// a constructor that takes the source language and nothing else, so the
// tables in the bundles_*_test.go files stop compiling when a constructor
// gains a parameter. That does not see a second exported constructor in an
// existing mapping package, a package variable or a setter. Those are for a
// reader of the package to see.
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
// implemented (GEN-02). No builtin contribution states it, so the registry has
// nothing to select for a request that names it. The request fails as a
// missing contribution and yields no plan and no output, from the same
// registry that generates under the structural profile.
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
// is registered fails as a missing contribution, whatever else the registry
// holds.
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

// TestRequestSelectsTheContributionsOfItsProfile: a second profile enters as
// contributions of its own. The read-only fixture family lives in its own
// package, states its own profile and maps a list to another native form than
// the builtin family. Its contributions are values of the contribution type,
// as the builtin ones are, and its declaration contribution uses the common
// declaration policy. One registry holds both families. Each request receives
// the contributions that state its profile and nothing of the other family.
// The horos mapping package and the registry know nothing of the fixture.
func TestRequestSelectsTheContributionsOfItsProfile(t *testing.T) {
	if readonlyprofile.Profile == rules.RepresentationProfile || readonlyprofile.Profile == rules.EnforcementProfile {
		t.Fatal("the fixture contribution states a profile of the builtin family:", readonlyprofile.Profile)
	}
	both := func(t *testing.T, modules ...core.Value) *profileComposition {
		t.Helper()
		return withReadonlyFamily(t, composeProfiles(t, modules, typescriptText(), tgtarget.Provider{}, typescriptBundles()...))
	}
	bindings := func(t *testing.T, q *profileComposition, profile string) map[string]bool {
		t.Helper()
		plan, err := q.plan(profile)
		if err != nil {
			t.Fatal(profile, err)
		}
		// A step names its binding as the reference of the contribution that
		// was selected for it.
		selected := map[string]bool{}
		for _, step := range plan.Steps() {
			selected[step.Binding] = true
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
		// The inspectable plan names the selected contributions.
		builtinList, err := typescriptBundleTable()["list"].Contribution(q.request.Input.Source.Language())
		if err != nil {
			t.Fatal(err)
		}
		structuralList, readonlyList := builtinList.Reference.String(), readonlyprofile.ListReference.String()
		if selected := bindings(t, q, rules.RepresentationProfile); !selected[structuralList] || selected[readonlyList] {
			t.Error("the structural plan selected", selected)
		}
		if selected := bindings(t, q, readonlyprofile.Profile); !selected[readonlyList] || selected[structuralList] {
			t.Error("the read-only plan selected", selected)
		}
		if err := q.refuses(rules.EnforcementProfile); err != nil {
			t.Error("the enforcement profile, with two families registered:", err)
		}
	})

	t.Run("a family is not completed from another profile", func(t *testing.T) {
		// The fixture family maps no record. The builtin record contribution
		// sits in the same registry and states another profile, so it does
		// not answer: the record's domain has no contribution.
		q := both(t, paginationModule())
		if text := generate(t, q, rules.RepresentationProfile); !strings.Contains(text, "interface Page<T>") {
			t.Fatal("the structural request lost its records:\n" + text)
		}
		out, err := q.generate(readonlyprofile.Profile)
		var diagnostic *text.Diagnostic
		if !errors.As(err, &diagnostic) || diagnostic.Code != text.Unsupported || out.Snapshot().Valid() {
			t.Fatal("the read-only request was completed by the structural family:", err)
		}
		var refusal *g.Refusal
		if !errors.Is(err, g.ErrMissingContribution) || !errors.As(err, &refusal) || refusal.Owner != g.KernelOwner || refusal.Profile != readonlyprofile.Profile {
			t.Fatal("the read-only request did not fail as a missing contribution of its own profile:", err)
		}
	})

	t.Run("the fixture alone", func(t *testing.T) {
		q := withReadonlyFamily(t, composeProfiles[ts.Language](t, listsModules(t), typescriptText(), tgtarget.Provider{}))
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

// republishedPublisher is the publisher of the contributions that
// TestRegistryDoesNotJudgeAStatedProfile republishes.
const republishedPublisher = "example.test/TestRegistryDoesNotJudgeAStatedProfile"

// republished is a bundle of its own author: it takes a builtin contribution,
// publishes it under an identity of its own and states a profile of its own
// choice. The codec is the builtin one, for the schema that the contribution
// keeps.
func republished(bundle genfixture.Bundle[ts.Language], profile string) genfixture.Bundle[ts.Language] {
	return genfixture.Bundle[ts.Language]{Codec: bundle.Codec, Contribution: func(from representation.Language) (g.Contribution[h.Language, ts.Language, source.Reader], error) {
		contribution, err := bundle.Contribution(from)
		contribution.Reference.Identity.Publisher = republishedPublisher
		contribution.Profile = profile
		return contribution, err
	}}
}

// TestRegistryDoesNotJudgeAStatedProfile records where the refusal of the
// enforcement profile ends. The registry selects by the profile that a
// contribution states and does not judge whether the contribution keeps that
// promise. Contributions that republish the builtin, constraint-erasing
// rules under the enforcement identity are selected by a request for it
// and yield representation output. Their author states that profile and
// answers for it. bittype states it nowhere, and no constructor of a builtin
// mapping takes it (bittype/ACX-008, limits).
func TestRegistryDoesNotJudgeAStatedProfile(t *testing.T) {
	var bundles []genfixture.Bundle[ts.Language]
	for _, bundle := range typescriptBundles() {
		bundles = append(bundles, republished(bundle, rules.EnforcementProfile))
	}
	q := composeProfiles(t, []core.Value{paginationModule()}, typescriptText(), tgtarget.Provider{}, bundles...)
	out, err := q.generate(rules.EnforcementProfile)
	if err != nil {
		t.Fatal("the registry refused a stated profile:", err)
	}
	if text := q.text(t, out); !strings.Contains(text, "interface Page<T>") {
		t.Fatal(text)
	}
	// The plan names the republished contributions and none of bittype's.
	plan, err := q.plan(rules.EnforcementProfile)
	if err != nil {
		t.Fatal(err)
	}
	for _, step := range plan.Steps() {
		if !strings.HasPrefix(step.Binding, republishedPublisher+"/") {
			t.Errorf("a step names the contribution %s", step.Binding)
		}
	}
	// The builtin contributions themselves were not registered, so nothing
	// states the structural profile in this registry.
	if err := q.refuses(rules.RepresentationProfile); err != nil {
		t.Fatal(err)
	}
}

// TestRegistrationNeedsAStatedProfile: a contribution that states no profile,
// and a contribution that states nothing at all, cannot be registered. Both
// fail with ErrInvalid.
func TestRegistrationNeedsAStatedProfile(t *testing.T) {
	from, err := representation.NewLanguage[h.Language](h.Language{}, h.Profile)
	if err != nil {
		t.Fatal(err)
	}
	to := nativeChecker(t, typescriptText()).Language()
	var registry c.Registry
	unstated, err := typescriptBundles()[0].Contribution(from)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.RegisterContribution(&registry, from, to, unstated); err != nil {
		t.Fatal("the contribution with its profile:", err)
	}
	unstated.Reference.Identity.Key += "/unstated"
	unstated.Profile = ""
	if err := c.RegisterContribution(&registry, from, to, unstated); !errors.Is(err, g.ErrInvalid) {
		t.Error("a contribution that states no profile:", err)
	}
	if err := c.RegisterContribution(&registry, from, to, g.Contribution[h.Language, ts.Language, source.Reader]{}); !errors.Is(err, g.ErrInvalid) {
		t.Error("a contribution that states nothing:", err)
	}
}
