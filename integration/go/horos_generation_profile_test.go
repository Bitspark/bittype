package integration_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
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
)

// secondProfile and unregisteredProfile belong to this file. The horos mapping
// package knows neither name. Nothing is ever registered under
// unregisteredProfile.
const (
	secondProfile       = "integration-test/second-representation/1"
	unregisteredProfile = "integration-test/unregistered/1"
)

// profileLibraries is the builtin library set of one pair, built under the
// profile it is given.
type profileLibraries[L any] func(profile string) []g.MappingLibrary[h.Language, L, source.Reader]

// eachProfileTarget runs one check against all four builtin pairs.
func eachProfileTarget(t *testing.T, check func(t *testing.T, pair profilePair)) {
	t.Helper()
	for _, pair := range []profilePair{
		newProfilePair[ts.Language]("typescript", typescriptText(), tgtarget.Provider{}, typescriptLibraries),
		newProfilePair[gs.Language]("go", golangText(), ggtarget.Provider{Package: "pagination"}, golangLibraries),
		newProfilePair[rs.Language]("rust", rustText(), rgtarget.Provider{}, rustLibraries),
		newProfilePair[cs.Language]("cpp", cppText(), cgtarget.Provider{Namespace: "pagination"}, cppLibraries),
	} {
		t.Run(pair.name, func(t *testing.T) { check(t, pair) })
	}
}

// profilePair hides the target language of one pair behind three operations.
// compose registers the source, the target and the pair's libraries built
// under each of the given profiles, and reports every registration error.
// plan and generate request the pagination module under one profile.
type profilePair struct {
	name     string
	compose  func(t *testing.T, registered ...string) (*profileComposition, []error)
	plan     func(t *testing.T, q *profileComposition, requested string) (g.GenerationPlan, error)
	generate func(t *testing.T, q *profileComposition, requested string) (string, bool, error)
}
type profileComposition struct {
	registry  c.Registry
	request   c.Request
	libraries int // libraries offered to the registry, refused or not
}

func newProfilePair[L any](name string, p nativeProvider[L], target g.Target[L], libraries profileLibraries[L]) profilePair {
	return profilePair{
		name: name,
		compose: func(t *testing.T, registered ...string) (*profileComposition, []error) {
			t.Helper()
			_, from, input := generationInput(t, paginationModule())
			to := nativeChecker(t, p).Language()
			q := &profileComposition{}
			var err error
			// The fixture registers the source and the target. It registers no
			// library here, so that every library's own answer is seen below.
			if q.request, err = genfixture.Request[L](&q.registry, input, to, "generated-profile", target); err != nil {
				t.Fatal(err)
			}
			var refusals []error
			for _, profile := range registered {
				for _, library := range libraries(profile) {
					q.libraries++
					if _, err := c.RegisterMapping(&q.registry, from, to, library); err != nil {
						refusals = append(refusals, err)
					}
				}
			}
			return q, refusals
		},
		plan: func(t *testing.T, q *profileComposition, requested string) (g.GenerationPlan, error) {
			t.Helper()
			request := q.request
			request.Profile = requested
			return q.registry.PlanGeneration(context.Background(), request, nativeChecker(t, p))
		},
		generate: func(t *testing.T, q *profileComposition, requested string) (string, bool, error) {
			t.Helper()
			request := q.request
			request.Profile = requested
			out, err := q.registry.Generate(context.Background(), request, nativeChecker(t, p))
			if err != nil || !out.Snapshot().Valid() {
				return "", out.Snapshot().Valid(), err
			}
			return formatted(t, p, outputTree[L](t, out)), true, nil
		},
	}
}

// producedNothing reports whether a request produced nothing and failed as unsupported.
func producedNothing(text string, valid bool, err error) bool {
	return errors.Is(err, g.ErrUnsupported) && !valid && text == ""
}

// TestHorosMappingLibrariesAnswerUnderASecondProfile: the common horos mapping
// adapter holds no list of admitted profile names. The builtin libraries of
// every pair, built under a profile that only this file names, register, and a
// request for that profile generates through them. The kernel selects by the
// requested identity alone: a registry that holds only the second profile
// refuses the structural one.
func TestHorosMappingLibrariesAnswerUnderASecondProfile(t *testing.T) {
	eachProfileTarget(t, func(t *testing.T, pair profilePair) {
		structural, refusals := pair.compose(t, rules.RepresentationProfile)
		if len(refusals) != 0 {
			t.Fatal("the structural profile did not register:", refusals)
		}
		want, valid, err := pair.generate(t, structural, rules.RepresentationProfile)
		if err != nil || !valid || want == "" {
			t.Fatal("the structural profile did not generate:", err)
		}
		second, refusals := pair.compose(t, secondProfile)
		if len(refusals) != 0 {
			t.Fatal("a second profile was refused at registration:", refusals)
		}
		got, valid, err := pair.generate(t, second, secondProfile)
		if err != nil || !valid {
			t.Fatal("a request for the second profile did not generate:", err)
		}
		if got != want {
			t.Fatalf("the same libraries generated another text under the second profile:\n%s\nwant:\n%s", got, want)
		}
		if text, valid, err := pair.generate(t, second, rules.RepresentationProfile); !producedNothing(text, valid, err) {
			t.Fatal("a registry that holds only the second profile answered the structural one:", err)
		}
	})
}

// TestRequestSelectsTheLibrariesOfItsProfile: one registry holds two profiles
// whose list mapping differs. The second profile's list library exists under
// that profile only and is built through the common adapter and the common
// patterns. Each request receives the libraries of its own profile.
func TestRequestSelectsTheLibrariesOfItsProfile(t *testing.T) {
	_, from, input := generationInput(t, paginationModule())
	checker := nativeChecker(t, typescriptText())
	to := checker.Language()
	var registry c.Registry
	request, err := genfixture.Request(&registry, input, to, "generated-profile", tgtarget.Provider{}, typescriptLibraries(rules.RepresentationProfile)...)
	if err != nil {
		t.Fatal(err)
	}
	readonly := rules.Library[ts.Language]{Profile: secondProfile, DescribeRules: listContribution("integration-test/readonly-list/1", "", "ReadonlyArray")}
	// Position 2 of the builtin set is its list library; the second profile
	// replaces it with its own.
	builtin := typescriptLibraries(secondProfile)
	for _, library := range append(append(builtin[:2:2], builtin[3:]...), readonly) {
		if _, err := c.RegisterMapping(&registry, from, to, library); err != nil {
			t.Fatal("a library of the second profile was refused:", err)
		}
	}
	text := func(profile string) string {
		t.Helper()
		q := request
		q.Profile = profile
		out, err := registry.Generate(context.Background(), q, checker)
		if err != nil {
			t.Fatal(profile, err)
		}
		return formatted(t, typescriptText(), outputTree[ts.Language](t, out))
	}
	if structural := text(rules.RepresentationProfile); !strings.Contains(structural, "items: T[]") || strings.Contains(structural, "ReadonlyArray") {
		t.Fatal("the structural request received another profile's list:\n" + structural)
	}
	if second := text(secondProfile); !strings.Contains(second, "items: ReadonlyArray<T>") || strings.Contains(second, "T[]") {
		t.Fatal("the second profile's request did not receive its own list:\n" + second)
	}
}

// TestEnforcementProfileIsRefused: the enforcement profile is declared and not
// implemented (GEN-02). No builtin library answers under it, so a consumer
// cannot register the representation libraries under that name. A request for
// it fails with ErrUnsupported and yields no plan and no output, also when the
// same registry holds the representation libraries.
func TestEnforcementProfileIsRefused(t *testing.T) {
	eachProfileTarget(t, func(t *testing.T, pair profilePair) {
		structural, refusals := pair.compose(t, rules.RepresentationProfile)
		if len(refusals) != 0 {
			t.Fatal("the structural profile did not register:", refusals)
		}
		want, valid, err := pair.generate(t, structural, rules.RepresentationProfile)
		if err != nil || !valid || want == "" {
			t.Fatal("the structural profile did not generate:", err)
		}
		// The registry that just generated representation output does not
		// answer an enforcement request with it.
		if p, err := pair.plan(t, structural, rules.EnforcementProfile); !errors.Is(err, g.ErrUnsupported) || p.Valid() {
			t.Error("an enforcement request planned over the representation libraries:", err)
		}
		if text, valid, err := pair.generate(t, structural, rules.EnforcementProfile); !producedNothing(text, valid, err) {
			t.Error("an enforcement request generated over the representation libraries:", err, text)
		}
		// A consumer that builds every builtin library under the enforcement
		// profile gets one refusal per library and still no output.
		enforcement, refusals := pair.compose(t, rules.EnforcementProfile)
		if enforcement.libraries == 0 || len(refusals) != enforcement.libraries {
			t.Errorf("%d of %d libraries refused the enforcement profile", len(refusals), enforcement.libraries)
		}
		for _, err := range refusals {
			if !errors.Is(err, g.ErrUnsupported) {
				t.Error("a library refused the enforcement profile with another error:", err)
			}
		}
		if p, err := pair.plan(t, enforcement, rules.EnforcementProfile); !errors.Is(err, g.ErrUnsupported) || p.Valid() {
			t.Error("an enforcement request planned:", err)
		}
		if text, valid, err := pair.generate(t, enforcement, rules.EnforcementProfile); !producedNothing(text, valid, err) {
			t.Error("an enforcement request generated:", err, text)
		}
	})
}

// TestUnregisteredProfileIsRefused: a request for a profile under which nothing
// is registered fails with ErrUnsupported, whatever else the registry holds.
func TestUnregisteredProfileIsRefused(t *testing.T) {
	eachProfileTarget(t, func(t *testing.T, pair profilePair) {
		for name, registered := range map[string][]string{"empty registry": nil, "another profile registered": {rules.RepresentationProfile}} {
			t.Run(name, func(t *testing.T) {
				q, refusals := pair.compose(t, registered...)
				if len(refusals) != 0 {
					t.Fatal("registration:", refusals)
				}
				if p, err := pair.plan(t, q, unregisteredProfile); !errors.Is(err, g.ErrUnsupported) || p.Valid() {
					t.Error("an unregistered profile planned:", err)
				}
				if text, valid, err := pair.generate(t, q, unregisteredProfile); !producedNothing(text, valid, err) {
					t.Error("an unregistered profile generated:", err, text)
				}
			})
		}
	})
}
