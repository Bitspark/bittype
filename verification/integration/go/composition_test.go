package integration_test

import (
	"context"
	"errors"
	"go/types"
	"slices"
	"strings"
	"testing"

	association "github.com/Bitspark/bittype/contracts/verification/go"
	"github.com/Bitspark/bittype/semantics/go"
	v "github.com/Bitspark/bittype/verification/files/horos-to-go/go"
	"github.com/Bitspark/bittype/verification/files/horos-to-go/mappings/list/go"
	"github.com/Bitspark/bittype/verification/files/horos-to-go/mappings/record/go"
	"github.com/Bitspark/bittype/verification/files/horos-to-go/mappings/scalar/go"
	h "github.com/Bitspark/bittype/verification/files/source/horos/go"
	"github.com/Bitspark/bittype/verification/integration/go/testdata/pair"
	joint "github.com/Bitspark/bittype/verification/logos/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	"github.com/Bitspark/logos-contract/contract/go/typed"
	core "github.com/bitspark/ontos/core/go"
)

const pairGo = "type Pair[R,L any] struct { First L; Second R }; "

func TestConsumerPairComposesUnderListAndRecord(t *testing.T) {
	source := sourceModel(t, strings.Replace(sourcePage, "list<T>", "list<tuple<T,bool>>", 1))
	code := pairGo + strings.Replace(nativePage, "Items []N", "Items []Pair[bool,N]", 1)
	target := nativeModel(t, code)
	helper := native(t, target, "Pair")
	rule := pair.Rule(helper)
	if _, err := bindPage(t, source, target); !errors.Is(err, v.ErrCoverage) {
		t.Fatal("missing pair selection", err)
	}
	p, err := bindPage(t, source, target, rule)
	if err != nil {
		t.Fatal(err)
	}
	claim, err := p.family.Theorem()
	if err != nil {
		t.Fatal(err)
	}
	limits := joint.DefaultLimits()
	limits.Checker.Work = 20000000
	allowed := append(allowedLaws(), pair.Law)
	checker := composition(t, typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker, Dependencies: allowed}, limits, p.authorities...)
	theorem := checkClaim(t, checker, claim, p.family.Candidate())
	accepted(t, theorem)
	// The tuple law receives the actual Horos operand and the selected Go helper as origins.
	fields, err := declaration(t, source, "Page").Fields()
	if err != nil {
		t.Fatal(err)
	}
	index := slices.IndexFunc(fields, func(f h.Field) bool { return f.Name == "items" })
	if index < 0 {
		t.Fatal("missing items field")
	}
	items := fields[index]
	tupleSite := source.OperandSite(items.Site, items.Type.Operands()[0].Path())
	request := theorem.Request()
	sites := []association.Site{
		{Participant: request.Association().Source(), Node: tupleSite},
		{Participant: request.Association().Target(), Node: helper.Site()},
	}
	found := false
	statement := request.Statement()
	for _, origin := range request.Origins() {
		if origin.Location.Phase != typed.DeclarationPhase || !slices.Contains(statement.Declarations[origin.Location.Index].Dependencies(), pair.Law) {
			continue
		}
		found = slices.Contains(origin.Sites, sites[0]) && slices.Contains(origin.Sites, sites[1])
	}
	if !found {
		t.Fatal("tuple law lost its source operand or selected helper origin")
	}
	for _, depth := range []int{1, 2} {
		application, err := p.family.Application(theorem, p.argument, depth)
		if err != nil {
			t.Fatal(err)
		}
		result := checkClaim(t, checker, application, p.argument.Candidate(depth))
		accepted(t, result)
		closure, _ := result.Outcome().Accepted()
		for _, law := range []string{record.Law, list.Law, pair.Law} {
			if !slices.Contains(closure.Dependencies(), law) {
				t.Fatalf("missing %s in %v", law, closure.Dependencies())
			}
		}
		refused := checkClaim(t, receiver(t, p.authorities...), application, p.argument.Candidate(depth))
		if refused.Outcome().Code() != typed.UnadmittedDependency {
			t.Fatal("import bypassed receiver policy", refused.Outcome().Code())
		}
	}
	// Same artifacts and identical keys, but an independently selected interpretation.
	fresh, err := bindPage(t, source, target, rule)
	if err != nil {
		t.Fatal(err)
	}
	if r, err := receiver(t, fresh.authorities...).Prepare(context.Background(), claim); !errors.Is(err, joint.ErrInvalid) || r.Valid() {
		t.Fatal("accepted foreign question", err)
	}
	if _, err = fresh.family.Application(theorem, fresh.argument, 1); !errors.Is(err, joint.ErrInvalid) {
		t.Fatal("rebound theorem", err)
	}
	other := nativeModel(t, code)
	if _, err = bindPage(t, source, target, pair.Rule(native(t, other, "Pair"))); !errors.Is(err, v.ErrCoverage) {
		t.Fatal("foreign helper selected by spelling", err)
	}
	for _, changed := range []string{strings.Replace(code, "First L", "First R", 1), strings.Replace(code, "Pair[bool,N]", "Pair[N,bool]", 1)} {
		g := nativeModel(t, changed)
		if _, err := bindPage(t, source, g, pair.Rule(native(t, g, "Pair"))); !errors.Is(err, v.ErrCoverage) {
			t.Fatal("changed representation", err)
		}
	}
	duplicate := rule
	duplicate.Key = semantics.Key{Owner: "other", Name: "tuple-pair"}
	if _, err = bindPage(t, source, target, rule, duplicate); !errors.Is(err, v.ErrAmbiguous) {
		t.Fatal("ambiguous selection", err)
	}
	if _, err = v.Select(rule, rule); !errors.Is(err, joint.ErrInvalid) {
		t.Fatal("duplicate owner", err)
	}
}

// This library supplies the correct law and deliberately applies it incorrectly.
type wrongOrder struct{ pair.ReversePremises }

func (wrongOrder) Apply(law typed.Expr, children []v.Evidence) typed.Expr {
	return (v.Congruence{}).Apply(law, children)
}
func TestIndependentLawCannotCertifyItsOwnCandidate(t *testing.T) {
	source := sourceModel(t, strings.Replace(sourcePage, "list<T>", "list<tuple<T,bool>>", 1))
	target := nativeModel(t, pairGo+strings.Replace(nativePage, "Items []N", "Items []Pair[bool,N]", 1))
	rule := pair.Rule(native(t, target, "Pair"))
	rule.Law = wrongOrder{}
	p, err := bindPage(t, source, target, rule)
	if err != nil {
		t.Fatal(err)
	}
	claim, err := p.family.Theorem()
	if err != nil {
		t.Fatal(err)
	}
	limits := joint.DefaultLimits()
	limits.Checker.Work = 20000000
	c := composition(t, typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker, Dependencies: append(allowedLaws(), pair.Law)}, limits, p.authorities...)
	bad := checkClaim(t, c, claim, p.family.Candidate())
	if bad.Outcome().Status() != typed.InvalidStatus {
		t.Fatal("incorrect proof", bad.Outcome().Status())
	}
}
func TestChildBelongsToThisProjection(t *testing.T) {
	s := sourceModel(t, sourcePage)
	g := nativeModel(t, nativePage)
	listRule := list.Rule()
	var saved v.Subject
	spy := v.Rule{Key: semantics.Key{Owner: "test", Name: "capture"}, Dependency: "test/capture", Law: v.Congruence{}, Project: func(p v.Subject) ([]v.Subject, bool, error) {
		view := p.Source()
		*view = decl.Expression{}
		children, ok, err := listRule.Project(p)
		if ok {
			saved = children[0]
		}
		return nil, false, err
	}}
	if _, err := bindPage(t, s, g, spy); err != nil {
		t.Fatal(err)
	}
	// Use a distinct capture for the parent and a child saved from the first.
	root := sourceModel(t, strings.Replace(sourcePage, "list<T>", "tuple<T,bool>", 1))
	replay := v.Rule{Key: semantics.Key{Owner: "test", Name: "replay"}, Dependency: "test/replay", Law: v.Congruence{}, Project: func(p v.Subject) ([]v.Subject, bool, error) {
		if p.Source().Constructor() != "tuple" {
			return nil, false, nil
		}
		return []v.Subject{saved}, true, nil
	}}
	if _, err := bindPage(t, root, g, replay); !errors.Is(err, joint.ErrInvalid) {
		t.Fatal("foreign child entered parent", err)
	}
}

func TestCounterexampleUsesSelectedNonBooleanWitness(t *testing.T) {
	s := sourceModel(t, "type Box<T: h.Bool> = {value:T;};")
	g := nativeModel(t, "type Box[T any] struct {Value T}")
	d := declaration(t, s, "Box")
	k, err := d.ParameterKind(0)
	if err != nil {
		t.Fatal(err)
	}
	text, err := s.Proper(core.NewTuple(core.NewAtom([]byte("horos/scalar")), core.NewAtom([]byte("utf8-text"))))
	if err != nil {
		t.Fatal(err)
	}
	nonmember, err := s.Satisfies(context.Background(), text, k)
	if err != nil {
		t.Fatal(err)
	}
	nativeText, err := g.Basic(types.String)
	if err != nil {
		t.Fatal(err)
	}
	box := native(t, g, "Box")
	admitted, err := nativeAdmission(t, context.Background(), g, box, nativeText)
	if err != nil {
		t.Fatal(err)
	}
	library := v.NewAdmissions(selectMappings(t, scalar.Text()))
	w, err := library.Counterexample(d, box, nonmember, admitted)
	if err != nil {
		t.Fatal(err)
	}
	accepted(t, checkClaim(t, receiver(t, library.Authority()), w.Claim, w.Candidate))
	other := v.NewAdmissions(selectMappings(t, scalar.Boolean()))
	if _, err := other.Counterexample(d, box, nonmember, admitted); !errors.Is(err, v.ErrCoverage) {
		t.Fatal(err)
	}
}

func TestMappingCannotEraseTheOriginalBinder(t *testing.T) {
	s := sourceModel(t, sourcePage)
	g := nativeModel(t, strings.Replace(nativePage, "Items []N", "Items []bool", 1))
	erasing := v.Rule{
		Key: semantics.Key{Owner: "test", Name: "erase"}, Dependency: "test/erase", Law: v.Congruence{},
		Project: func(p v.Subject) ([]v.Subject, bool, error) {
			_, parameter := p.Source().Parameter()
			return nil, parameter && p.Target().IsBasic(types.Bool), nil
		},
	}
	if _, err := bindPage(t, s, g, erasing); !errors.Is(err, v.ErrCoverage) {
		t.Fatal("mapping erased its captured parameter", err)
	}
}

// A library can project a named body's operands directly, without a record rule.
func TestNamedTupleUsesActualBodyAndNativeFieldOperands(t *testing.T) {
	s := sourceModel(t, "type Product<T> = tuple<T,bool>;")
	g := nativeModel(t, "type Product[N any] struct { First N; Second bool }")
	d := declaration(t, s, "Product")
	target := native(t, g, "Product")
	const dependency = "example/tuple-struct"
	rule := v.Rule{
		Key: semantics.Key{Owner: "example", Name: "tuple-struct"}, Dependency: dependency, Law: pair.ReversePremises{},
		Project: func(p v.Subject) ([]v.Subject, bool, error) {
			if p.Source().Constructor() != "tuple" || !p.Target().Same(target) || len(p.Source().Operands()) != 2 {
				return nil, false, nil
			}
			fields, ok := p.Target().Fields()
			if !ok || len(fields) != 2 || fields[0].Name != "First" || fields[1].Name != "Second" {
				return nil, false, nil
			}
			var children []v.Subject
			for i, field := range fields {
				if field.Tag != "" || field.Embedded {
					return nil, false, nil
				}
				child, err := p.Operand(i, field.Type)
				if err != nil {
					return nil, false, err
				}
				children = append(children, child)
			}
			return children, true, nil
		},
	}
	library := v.NewFamilies(selectMappings(t, rule, scalar.Boolean()))
	k := importedKind(t, s, "Data")
	bound, err := k.Arrow(s.Root())
	if err != nil {
		t.Fatal(err)
	}
	value, err := d.Value()
	if err != nil {
		t.Fatal(err)
	}
	guarantee, err := s.Satisfies(context.Background(), value, bound)
	if err != nil {
		t.Fatal(err)
	}
	family, err := library.Bind(d, target, nil, k, guarantee)
	if err != nil {
		t.Fatal(err)
	}
	claim, err := family.Theorem()
	if err != nil {
		t.Fatal(err)
	}
	limits := joint.DefaultLimits()
	limits.Checker.Work = 20000000
	policy := typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker, Dependencies: append(allowedLaws(), dependency)}
	checker := composition(t, policy, limits, library.Authorities()...)
	accepted(t, checkClaim(t, checker, claim, family.Candidate()))
}
