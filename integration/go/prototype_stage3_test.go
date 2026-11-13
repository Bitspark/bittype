package integration_test

// PROTOTYPE (bittype#196), stage 3: identities and refusals. Rows 10, 15, 17,
// 21 and 25 of gate lines 2 and 3.

import (
	"context"
	"errors"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// refusalValues are the values of the closed set of refusal classes.
var refusalValues = map[string]error{
	"unsupported": g.ErrUnsupported, "refused by design": g.ErrRefusedByDesign, "requirement unsatisfied": g.ErrRequirementUnsatisfied,
	"requirement unresolved": g.ErrRequirementUnresolved, "missing facet": g.ErrMissingFacet, "ambiguous contribution": g.ErrAmbiguous,
	"missing contribution": g.ErrMissingContribution, "limit": g.ErrLimit,
}

// refusalClasses names every refusal value that an error matches.
func refusalClasses(err error) []string {
	var matched []string
	for name, value := range refusalValues {
		if errors.Is(err, value) {
			matched = append(matched, name)
		}
	}
	return matched
}

func expectRefusal(t *testing.T, name string, err error, class string) *g.Refusal {
	t.Helper()
	var refusal *g.Refusal
	if matched := refusalClasses(err); !errors.As(err, &refusal) || len(matched) != 1 || matched[0] != class {
		t.Errorf("%s: got %v, which matches %v; want exactly %q", name, err, matched, class)
		return &g.Refusal{}
	}
	t.Logf("GREEN %s: %v", name, err)
	return refusal
}

// Row 15: "Two contributions with one local key under two publishers
// register side by side, and two with one publisher and one key fail."
//
// Shows: the identity that duplicates compare is the publisher and the local
// key, across targets, and a registration carries it.
// Does not show: that a publisher's name is its owner's. A name does not
// authenticate; the composition lock is deferred (row 15a).
func TestPrototypeStage3PublisherAndLocalKey(t *testing.T) {
	from, _ := protoSource(t, "module example.identity;\n\ntype Name = text;\n")
	goLanguage, tsLanguage := nativeChecker(t, golangText()).Language(), nativeChecker(t, typescriptText()).Language()
	var registry c.Registry
	text := golangBundle(t, from).find(t, "horos-to-golang/text")
	if err := c.RegisterContribution(&registry, from, goLanguage, text); err != nil {
		t.Fatal(err)
	}
	other := text
	other.Reference.Identity.Publisher = "example.org/another-publisher"
	if err := c.RegisterContribution(&registry, from, goLanguage, other); err != nil {
		t.Errorf("one local key under two publishers does not register: %v", err)
	} else {
		t.Log("GREEN one local key under two publishers registers side by side")
	}
	if err := c.RegisterContribution(&registry, from, goLanguage, text); err == nil {
		t.Error("one publisher registered one key twice")
	} else {
		t.Logf("GREEN one publisher with one key twice fails: %v", err)
	}
	// The key is unique for its publisher across targets (O4): the target is
	// a stated field and no part of the key.
	sameKey := typescriptBundle(t, from).find(t, "horos-to-typescript/text")
	sameKey.Reference.Identity = text.Reference.Identity
	if err := c.RegisterContribution(&registry, from, tsLanguage, sameKey); err == nil {
		t.Error("one publisher registered one key for two targets")
	} else {
		t.Logf("GREEN one publisher with one key for a second target fails: %v", err)
	}
	// A contribution is registered for the target it states.
	if err := c.RegisterContribution(&registry, from, tsLanguage, golangBundle(t, from).find(t, "horos-to-golang/list")); err == nil {
		t.Error("a contribution that states Go registered for TypeScript")
	}
	misstated := typescriptBundle(t, from).find(t, "horos-to-typescript/list")
	misstated.Target = golangTarget
	if err := c.RegisterContribution(&registry, from, tsLanguage, misstated); !errors.Is(err, g.ErrScope) {
		t.Errorf("a contribution that states another target than it is registered for: %v", err)
	} else {
		t.Logf("GREEN the stated target is checked at registration: %v", err)
	}
}

// Row 17, O20: "The contributions of one family for two targets carry one
// meaning identity and two representation revisions, and a reference that
// lacks either part is refused."
//
// Shows: for the five families that the prototype moved for both targets, the
// two contributions name one meaning and each its own representation
// revision, and a reference without a meaning revision or without a
// representation revision cannot be registered.
// Does not show: who owns a family's meaning (O19), or that the two targets
// mean the same. The equality is of two spelled identities (finding F2).
func TestPrototypeStage3OneMeaningTwoRepresentationRevisions(t *testing.T) {
	from, _ := protoSource(t, "module example.identity;\n\ntype Name = text;\n")
	goBundle, tsBundle := golangBundle(t, from), typescriptBundle(t, from)
	for _, name := range []string{"text", "list", "record", "declarations", "optional"} {
		a, b := goBundle.find(t, "horos-to-golang/"+name).Reference, tsBundle.find(t, "horos-to-typescript/"+name).Reference
		if a.Meaning != b.Meaning || a.Representation == b.Representation || a.Identity == b.Identity {
			t.Errorf("%s: %s and %s", name, a, b)
		} else {
			t.Logf("GREEN %s: one meaning %s/%s, revisions %s and %s", name, a.Meaning.Identity, a.Meaning.Revision, a.Representation, b.Representation)
		}
	}
	for part, strip := range map[string]func(*g.ContributionReference){
		"the meaning revision":        func(r *g.ContributionReference) { r.Meaning.Revision = "" },
		"the meaning identity":        func(r *g.ContributionReference) { r.Meaning.Identity = "" },
		"the representation revision": func(r *g.ContributionReference) { r.Representation = "" },
	} {
		var registry c.Registry
		lacking := goBundle.find(t, "horos-to-golang/text")
		strip(&lacking.Reference)
		if err := c.RegisterContribution(&registry, from, nativeChecker(t, golangText()).Language(), lacking); !errors.Is(err, g.ErrInvalid) {
			t.Errorf("a reference that lacks %s registered: %v", part, err)
		} else {
			t.Logf("GREEN a reference that lacks %s is refused", part)
		}
	}
}

// Rows 10 and 21, O17, O26: the refusal values.
//
// Shows: missing and ambiguous are two values. A family that is not built is
// "missing contribution". A specialization is refused as unsupported with the
// application named, while another application of the same declaration
// plans. Each refusal of a rule of the moved packages matches exactly one
// value, and an unexpected error of a rule matches none.
// Does not show: the refusals of the ten mapping packages that are not moved.
// On the mapping contract a rule's refusal is a text diagnostic, and it
// matches none of the values.
func TestPrototypeStage3RefusalValues(t *testing.T) {
	const module = "module example.refusals;\n\ntype Page<T> = { items: list<T>; };\ntype Person = { name: text; };\ntype TextPage = Page<text>;\ntype People = Page<Person>;\n"
	from, input := protoSource(t, module)
	b := golangBundle(t, from)
	// Missing: no contribution answers the list form.
	_, err := protoPlanGo(t, input, b.without("horos-to-golang/list"), protoRequest{})
	missing := expectRefusal(t, "no contribution for a form", err, "missing contribution")
	if missing.Subject == "" || len(missing.Path) == 0 || missing.Profile != rules.RepresentationProfile || missing.Target != golangTarget {
		t.Errorf("the refusal carries subject %q, path %v, profile %q, target %q", missing.Subject, missing.Path, missing.Profile, missing.Target)
	}
	// Ambiguous: two contributions answer it. Registration order is no preference.
	second := b.find(t, "horos-to-golang/list")
	second.Reference.Identity.Publisher = "example.org/another-publisher"
	_, err = protoPlanGo(t, input, b.with(second), protoRequest{})
	expectRefusal(t, "two contributions for a form", err, "ambiguous contribution")
	// A family that is not built: TypeScript booleans are not moved.
	_, unbuilt := protoSource(t, "module example.unbuilt;\n\ntype Flag = bool;\n")
	_, err = protoPlanTS(t, unbuilt, typescriptBundle(t, from), protoRequest{})
	expectRefusal(t, "a family that is not built", err, "missing contribution")
	// Row 21: a specialization. The contribution's domain is the application
	// Page<text>. The first implementation refuses it and names the application.
	_, pin := protoSchemaModule(t, module)
	page, err := source.Named(from, pin, "Page")
	if err != nil {
		t.Fatal(err)
	}
	text, err := source.Builtin(from, "scalar", "utf8-text")
	if err != nil {
		t.Fatal(err)
	}
	specialization := b.find(t, "horos-to-golang/text")
	specialization.Reference.Identity = g.ContributionIdentity{Publisher: fixturePublisher, Key: "text-page"}
	specialization.Domain = g.ForSpecialization(page, text)
	_, err = protoPlanGo(t, input, b.with(specialization), protoRequest{})
	refused := expectRefusal(t, "a specialization of Page<text>", err, "unsupported")
	if !strings.Contains(refused.Application, "/Page<") || !strings.Contains(refused.Application, "utf8-text") || refused.Contribution != specialization.Reference {
		t.Errorf("the refusal names the application %q and the contribution %s", refused.Application, refused.Contribution)
	}
	if _, err = protoPlanGo(t, input, b.with(specialization), protoRequest{roots: []string{"People"}}); err != nil {
		t.Errorf("another application of the specialized declaration does not plan: %v", err)
	}
	// The refusals of the rules of the moved packages.
	for _, r := range []struct{ name, module, class, target string }{
		{"Go record: a nested record", "module example.nested;\n\ntype Outer = { inner: { name: text; }; };\n", "refused by design", "golang"},
		{"TypeScript record: a nested record", "module example.nested;\n\ntype Outer = { inner: { name: text; }; };\n", "refused by design", "typescript"},
		{"Go record: a field name the naming rule does not support", "module example.naming;\n\ntype Outer = { a__b: text; };\n", "unsupported", "golang"},
		{"Go declarations: a declaration name the naming rule does not support", "module example.naming;\n\ntype a__b = text;\n", "unsupported", "golang"},
		{"TypeScript declarations: a declaration name the naming rule does not support", "module example.naming;\n\ntype a__b = text;\n", "unsupported", "typescript"},
	} {
		from, input := protoSource(t, r.module)
		var err error
		if r.target == "golang" {
			_, err = protoPlanGo(t, input, golangBundle(t, from), protoRequest{})
		} else {
			_, err = protoPlanTS(t, input, typescriptBundle(t, from), protoRequest{})
		}
		refusal := expectRefusal(t, r.name, err, r.class)
		if refusal.Owner == "" || refusal.Reason == "" || !refusal.Contribution.Valid() || refusal.Subject == "" {
			t.Errorf("%s: the refusal lacks its owner, reason, contribution or subject: %+v", r.name, refusal)
		}
	}
	// An unexpected error of a rule is an implementation failure and no
	// refusal, even where the error it returns is a value of the contract. A
	// rule cannot use a class that is the planner's.
	for name, failure := range map[string]error{"a rule that returns the unsupported value itself": g.ErrUnsupported, "a rule that refuses with a class of the planner": g.Refuse(g.RefusalMissingFacet, fixturePublisher, "not mine to say")} {
		failing := b.find(t, "horos-to-golang/text")
		failing.Rule.Decide = func(context.Context, horosDecide) (g.Decision, error) { return g.Decision{}, failure }
		_, err = protoPlanGo(t, input, b.replace(failing), protoRequest{})
		var implementation *g.RuleFailure
		if matched := refusalClasses(err); !errors.As(err, &implementation) || len(matched) != 0 {
			t.Errorf("%s: got %v, which matches %v", name, err, matched)
		} else {
			t.Logf("GREEN %s is an implementation failure and matches no refusal value: %v", name, err)
		}
	}
}

// lenient is the Go target without its check that one output has no two
// declarations of one name. That check is a realization concern; the test
// below looks at the plan before it.
type lenient struct{ ggtarget.Provider }

func (lenient) CheckProgram(context.Context, g.Program[gs.Language]) error { return nil }

// Row 25: "Two contributions that spell a helper alike get two symbols".
//
// Shows: a declaration that a contribution owns is keyed by the contribution
// reference and a local name. Two contributions that spell their helper alike
// get two logical symbols, and one contribution that uses its helper twice
// gets one. Red under a kernel that shares owned declarations by spelling.
// Does not show: what a Files stage does with two declarations of one name in
// one scope. The Go target of today refuses such a program; the test plans
// with that one check switched off.
// Does not show: a shared helper that a contribution asks for by a request
// (row 25, second half). The request kind exists and is refused as
// unsupported: whose key it uses is not specified.
func TestPrototypeStage3OwnedDeclarationsAreKeyedByOwner(t *testing.T) {
	from, input := protoSource(t, "module example.owned;\n\ntype Box = { first?: text; second?: text; boxed: list<text>; };\n")
	b := golangBundle(t, from)
	// A fixture list contribution that owns a helper spelled like the helper
	// of the optional contribution.
	boxed := b.find(t, "horos-to-golang/list")
	boxed.Reference = fixtureReference("boxed-list", "golang/1")
	boxed.Owned = []g.OwnedDescription[gs.Language]{{Name: "box", Declaration: g.NativeDeclaration[gs.Language]{
		Name: "BittypeOptional", Form: model.Form(model.FormDefined), Parameters: []g.NativeParameter{{Name: "T"}},
		Body: model.Record(model.Field("Items", false, model.Collection(g.Parameter[gs.Language](0)))),
	}}}
	boxed.Project = func(_ context.Context, p g.Projection[gs.Language]) (g.Projected[gs.Language], error) {
		wrapped, err := p.Owned("box", p.Dependency(0).Type)
		return g.Projected[gs.Language]{Type: wrapped}, err
	}
	check := func() error {
		plan, err := protoPlanWith[gs.Language](t, input, nativeChecker(t, golangText()), lenient{ggtarget.Provider{Package: "generated"}}, b.replace(boxed), protoRequest{})
		if err != nil {
			return err
		}
		var owned []string
		for _, s := range plan.Plan.Symbols() {
			if s.Role == "owned" && s.Intent == "BittypeOptional" {
				owned = append(owned, s.ID)
			}
		}
		if len(owned) != 2 || owned[0] == owned[1] {
			return errors.New("the helpers spelled BittypeOptional have the symbols " + strings.Join(owned, ", "))
		}
		t.Logf("two symbols: %s", strings.Join(owned, " and "))
		return nil
	}
	redGreen(t, "two contributions that spell a helper alike get two symbols", check, "owned-shared-by-spelling")
	// The request for a shared helper is stated and refused.
	asking := b.find(t, "horos-to-golang/text")
	asking.Rule.Name = func(context.Context, horosName) ([]horosRequest, error) {
		return []horosRequest{g.SharedHelper[h.Language]("example.org/shared/optional")}, nil
	}
	_, err := protoPlanGo(t, input, b.replace(asking), protoRequest{})
	if refusal := expectRefusal(t, "a request for a shared helper", err, "unsupported"); refusal.Reason != g.ReasonSharedHelper {
		t.Errorf("the reason is %q", refusal.Reason)
	}
}

var _ = representation.Language{}
var _ ts.Language
