package generation_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	gen "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type sourceLanguage struct{}
type nativeLanguage struct{}
type declaration = gen.Declaration[sourceLanguage]
type binder = gen.Binder[sourceLanguage]
type use = gen.Use[sourceLanguage]
type environment = gen.Environment[sourceLanguage]
type fragment = gen.Fragment[nativeLanguage]
type token = gen.Token[sourceLanguage, nativeLanguage]
type binding = gen.Binding[sourceLanguage]
type product = gen.Product[sourceLanguage, nativeLanguage]
type callee = gen.Callee[sourceLanguage, nativeLanguage]
type site = gen.Site[sourceLanguage, nativeLanguage]
type want = gen.Want[sourceLanguage, nativeLanguage]
type plan = gen.Plan[sourceLanguage, nativeLanguage]
type config = gen.Config[sourceLanguage, nativeLanguage]
type session = gen.Session[sourceLanguage, nativeLanguage]
type openData = gen.OpenData[sourceLanguage]
type useData = gen.UseData[sourceLanguage]
type formalBinding = gen.FormalBinding[sourceLanguage, nativeLanguage]
type captured = gen.Captured[sourceLanguage, nativeLanguage]

var errRecord = errors.New("fixture record is absent or belongs to another owner")
var errAdmission = errors.New("fixture whole-use applicability is not established")

type fixtureIdentity string

func (i fixtureIdentity) ID() language.ID { return language.ID(i) }
func (i fixtureIdentity) Name() string    { return string(i) }

type fixtureChecker struct{ fixtureIdentity }

func (fixtureChecker) Check(ast.View[sourceLanguage]) error { return nil }

func ref[K any](space gen.Space[K], slot uint64) gen.Ref[K] {
	r, ok := space.Ref(slot)
	if !ok {
		panic("invalid fixture address")
	}
	return r
}

// The fixture's records are immutable once an invocation starts. Hooks change
// what a deliberately faulty collaborator returns, not an already issued record.
type sourceStore struct {
	snapshot                            semantic.Snapshot
	declarations                        gen.Space[gen.DeclarationRole[sourceLanguage]]
	binders                             gen.Space[gen.BinderRole[sourceLanguage]]
	uses                                gen.Space[gen.UseRole[sourceLanguage]]
	environments                        gen.Space[gen.EnvironmentRole[sourceLanguage]]
	openRecords                         map[declaration]openData
	useRecords                          map[use]useData
	binderRecords                       map[binder]bool
	environmentRecords                  map[environment]bool
	invalidUses                         map[use]bool
	root                                environment
	openHook                            func(openData) openData
	admitHook                           func(useData) useData
	checkHook                           func()
	checks, authentications, admissions int
}

func (s *sourceStore) CheckInput(_ context.Context, snapshot semantic.Snapshot, evidence []semantic.Result) error {
	s.checks++
	if !snapshot.Same(s.snapshot) || len(evidence) != 0 {
		return errRecord
	}
	if s.checkHook != nil {
		s.checkHook()
	}
	return nil
}
func (s *sourceStore) Owners() gen.SourceOwners[sourceLanguage] {
	return gen.SourceOwners[sourceLanguage]{Declarations: s.declarations.Owner(), Binders: s.binders.Owner(), Uses: s.uses.Owner(), Environments: s.environments.Owner()}
}
func (s *sourceStore) RootEnvironment() environment { return s.root }
func (s *sourceStore) Authenticate(_ context.Context, r gen.References[sourceLanguage]) error {
	s.authentications++
	for _, v := range r.Declarations {
		if _, ok := s.openRecords[v]; !ok {
			return errRecord
		}
	}
	for _, v := range r.Binders {
		if !s.binderRecords[v] {
			return errRecord
		}
	}
	for _, v := range r.Uses {
		if _, ok := s.useRecords[v]; !ok {
			return errRecord
		}
	}
	for _, v := range r.Environments {
		if !s.environmentRecords[v] {
			return errRecord
		}
	}
	return nil
}
func (s *sourceStore) Open(_ context.Context, d declaration) (openData, error) {
	o, ok := s.openRecords[d]
	if !ok {
		return openData{}, errRecord
	}
	o = o.Clone()
	if s.openHook != nil {
		o = s.openHook(o)
	}
	return o, nil
}
func (s *sourceStore) Admit(_ context.Context, e environment, u use) (useData, error) {
	s.admissions++
	d, ok := s.useRecords[u]
	if !ok || d.Environment != e {
		return useData{}, errRecord
	}
	if s.invalidUses[u] {
		return useData{}, errAdmission
	}
	// Whole-use admission includes actuals even if the selected mapping will
	// deliberately erase them and never call Represent.
	for _, actual := range d.Actuals {
		if s.invalidUses[actual.Use] {
			return useData{}, errAdmission
		}
	}
	d = d.Clone()
	if s.admitHook != nil {
		d = s.admitHook(d)
	}
	return d, nil
}

type associationStore struct {
	space   gen.Space[gen.AssociationRole[sourceLanguage, nativeLanguage]]
	records map[token]captured
	count   int
}

func newAssociations() *associationStore {
	return &associationStore{space: gen.NewSpace[gen.AssociationRole[sourceLanguage, nativeLanguage]](), records: make(map[token]captured)}
}
func (a *associationStore) issue(w want, p fragment, alter func(*captured)) token {
	t := ref(a.space, uint64(len(a.records)+1))
	c := captured{Binding: w.Binding(), Token: t, Payload: p}
	if alter != nil {
		alter(&c)
	}
	a.records[t] = c
	return t
}
func (a *associationStore) Capture(_ context.Context, _ want, t token) (captured, error) {
	a.count++
	c, ok := a.records[t]
	if !ok {
		return captured{}, errRecord
	}
	return c, nil
}

type nativeStore struct {
	destinations         gen.Space[gen.DestinationRole[nativeLanguage]]
	layouts              gen.Space[gen.LayoutRole[nativeLanguage]]
	fragments            gen.Space[gen.FragmentRole[nativeLanguage]]
	destinationRecords   map[gen.Destination[nativeLanguage]]bool
	layoutRecords        map[gen.Layout[nativeLanguage]]gen.Destination[nativeLanguage]
	fragmentRecords      map[fragment]gen.FrameID
	output               representation.ASTSet
	assemblyHook         func() error
	assemblies, products int
}

func (n *nativeStore) destination(owned bool) gen.Destination[nativeLanguage] {
	d := ref(n.destinations, uint64(len(n.destinationRecords)+1))
	n.destinationRecords[d] = owned
	return d
}
func (n *nativeStore) layout(d gen.Destination[nativeLanguage]) gen.Layout[nativeLanguage] {
	l := ref(n.layouts, uint64(len(n.layoutRecords)+1))
	n.layoutRecords[l] = d
	return l
}
func (n *nativeStore) fragment(frame gen.FrameID) fragment {
	f := ref(n.fragments, uint64(len(n.fragmentRecords)+1))
	n.fragmentRecords[f] = frame
	return f
}
func (n *nativeStore) CheckDestination(_ context.Context, _ site, d gen.Destination[nativeLanguage], owned bool) error {
	actual, ok := n.destinationRecords[d]
	if !ok || actual != owned {
		return errRecord
	}
	return nil
}
func (n *nativeStore) CheckLayout(_ context.Context, _ site, d gen.Destination[nativeLanguage], l gen.Layout[nativeLanguage]) error {
	if saved, ok := n.layoutRecords[l]; !ok || saved != d {
		return errRecord
	}
	return nil
}
func (n *nativeStore) CheckProduct(_ context.Context, w want, p fragment) error {
	n.products++
	frame, ok := n.fragmentRecords[p]
	if !ok {
		return errRecord
	}
	if frame.Valid() && frame != w.Binding().Frame {
		return gen.ErrScope
	}
	return nil
}
func (n *nativeStore) CheckBody(_ context.Context, c callee, p fragment) error {
	frame, ok := n.fragmentRecords[p]
	if !ok {
		return errRecord
	}
	if frame.Valid() && frame != c.Frame() {
		return gen.ErrScope
	}
	return nil
}
func (n *nativeStore) CheckExposure(_ context.Context, _ site, _ declaration, c callee) error {
	_, ok := n.destinationRecords[c.Destination()]
	if !ok {
		return errRecord
	}
	return nil
}
func (n *nativeStore) Assemble(_ context.Context, _ []gen.Publication[sourceLanguage, nativeLanguage], _ []gen.Definition[sourceLanguage, nativeLanguage], _ gen.Limits) (representation.ASTSet, error) {
	n.assemblies++
	if n.assemblyHook != nil {
		if err := n.assemblyHook(); err != nil {
			return representation.ASTSet{}, err
		}
	}
	return n.output, nil
}

type useFunc func(context.Context, gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error)

func (f useFunc) Use(ctx context.Context, r gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) {
	return f(ctx, r)
}

type convention struct {
	introduce func(context.Context, gen.IntroduceRequest[sourceLanguage, nativeLanguage]) (gen.Introduction[sourceLanguage, nativeLanguage], error)
	apply     func(context.Context, gen.ApplyRequest[sourceLanguage, nativeLanguage]) (product, error)
}

func (c convention) Introduce(ctx context.Context, r gen.IntroduceRequest[sourceLanguage, nativeLanguage]) (gen.Introduction[sourceLanguage, nativeLanguage], error) {
	return c.introduce(ctx, r)
}
func (c convention) Apply(ctx context.Context, r gen.ApplyRequest[sourceLanguage, nativeLanguage]) (product, error) {
	return c.apply(ctx, r)
}

type defineFunc func(context.Context, gen.DefineRequest[sourceLanguage, nativeLanguage]) (gen.Definition[sourceLanguage, nativeLanguage], error)

func (f defineFunc) Define(ctx context.Context, r gen.DefineRequest[sourceLanguage, nativeLanguage]) (gen.Definition[sourceLanguage, nativeLanguage], error) {
	return f(ctx, r)
}

type publishFunc func(context.Context, gen.PublishRequest[sourceLanguage, nativeLanguage]) (gen.Publication[sourceLanguage, nativeLanguage], error)

func (f publishFunc) Publish(ctx context.Context, r gen.PublishRequest[sourceLanguage, nativeLanguage]) (gen.Publication[sourceLanguage, nativeLanguage], error) {
	return f(ctx, r)
}

type fixture struct {
	source                                                          *sourceStore
	native                                                          *nativeStore
	associations                                                    *associationStore
	target                                                          representation.Language
	input                                                           semantic.Input
	checker                                                         check.Checker
	generic, plain                                                  declaration
	first, second, spare                                            binder
	scalar, alternate, application, body, formalFirst, formalSecond use
	root, genericEnv                                                environment
	closed                                                          fragment
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	from, err := representation.NewLanguage[sourceLanguage](fixtureIdentity("generation-source"), "fixture")
	if err != nil {
		t.Fatal(err)
	}
	to, err := representation.NewLanguage[nativeLanguage](fixtureIdentity("generation-target"), "fixture")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[sourceLanguage](from, fixtureChecker{fixtureIdentity("generation-source")})
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := representation.Pack(from, ast.AST[sourceLanguage]{ID: "source", Root: 1, Nodes: []ast.Node[sourceLanguage]{{Kind: "module"}}})
	if err != nil {
		t.Fatal(err)
	}
	input := semantic.Input{Entry: ast.NodeRef{Tree: "source", Node: 1}, Trees: []representation.AST{envelope}}
	snapshot, err := semantic.Capture(context.Background(), checker, input)
	if err != nil {
		t.Fatal(err)
	}
	s := &sourceStore{snapshot: snapshot, declarations: gen.NewSpace[gen.DeclarationRole[sourceLanguage]](), binders: gen.NewSpace[gen.BinderRole[sourceLanguage]](), uses: gen.NewSpace[gen.UseRole[sourceLanguage]](), environments: gen.NewSpace[gen.EnvironmentRole[sourceLanguage]](), openRecords: make(map[declaration]openData), useRecords: make(map[use]useData), binderRecords: make(map[binder]bool), environmentRecords: make(map[environment]bool), invalidUses: make(map[use]bool)}
	n := &nativeStore{destinations: gen.NewSpace[gen.DestinationRole[nativeLanguage]](), layouts: gen.NewSpace[gen.LayoutRole[nativeLanguage]](), fragments: gen.NewSpace[gen.FragmentRole[nativeLanguage]](), destinationRecords: make(map[gen.Destination[nativeLanguage]]bool), layoutRecords: make(map[gen.Layout[nativeLanguage]]gen.Destination[nativeLanguage]), fragmentRecords: make(map[fragment]gen.FrameID)}
	f := &fixture{source: s, native: n, associations: newAssociations(), target: to, input: input, checker: checker, generic: ref(s.declarations, 1), plain: ref(s.declarations, 2), first: ref(s.binders, 1), second: ref(s.binders, 2), spare: ref(s.binders, 3), scalar: ref(s.uses, 1), alternate: ref(s.uses, 2), application: ref(s.uses, 3), body: ref(s.uses, 4), formalFirst: ref(s.uses, 5), formalSecond: ref(s.uses, 6), root: ref(s.environments, 1), genericEnv: ref(s.environments, 2)}
	s.root = f.root
	s.environmentRecords[f.root], s.environmentRecords[f.genericEnv] = true, true
	s.binderRecords[f.first], s.binderRecords[f.second], s.binderRecords[f.spare] = true, true, true
	s.useRecords[f.scalar] = useData{Use: f.scalar, Subject: f.scalar, Environment: f.root}
	s.useRecords[f.alternate] = useData{Use: f.alternate, Subject: f.alternate, Environment: f.root}
	s.useRecords[f.application] = useData{Use: f.application, Subject: f.application, Environment: f.root, Callee: f.generic, Actuals: []gen.Actual[sourceLanguage]{{Binder: f.first, Use: f.scalar}, {Binder: f.second, Use: f.alternate}}}
	s.useRecords[f.body] = useData{Use: f.body, Subject: f.body, Environment: f.genericEnv}
	s.useRecords[f.formalFirst] = useData{Use: f.formalFirst, Subject: f.formalFirst, Environment: f.genericEnv, Bound: f.first}
	s.useRecords[f.formalSecond] = useData{Use: f.formalSecond, Subject: f.formalSecond, Environment: f.genericEnv, Bound: f.second}
	s.openRecords[f.generic] = openData{Declaration: f.generic, Environment: f.genericEnv, Body: f.body, Formals: []gen.Formal[sourceLanguage]{{Binder: f.first, Subject: f.formalFirst}, {Binder: f.second, Subject: f.formalSecond}}}
	s.openRecords[f.plain] = openData{Declaration: f.plain, Environment: f.root, Body: f.scalar}
	f.closed = n.fragment(gen.FrameID{})
	n.output, err = representation.PackSet(to, ast.Set[nativeLanguage]{Entry: ast.NodeRef{Tree: "native", Node: 1}, Trees: []ast.AST[nativeLanguage]{{ID: "native", Root: 1, Nodes: []ast.Node[nativeLanguage]{{Kind: "module", Children: []ast.Child{{Role: "type", Node: 2}}}, {Kind: "type"}}}}})
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func (f *fixture) configuration(user useFunc) config {
	if user == nil {
		user = func(_ context.Context, r gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) {
			return r.Return(f.associations.issue(r.Want(), f.closed, nil))
		}
	}
	return config{Source: f.source, Native: f.native, Target: f.target, Plans: []plan{{Key: "direct", User: user, Associations: f.associations}}, Select: func(useData) (string, error) { return "direct", nil }}
}

func (f *fixture) genericPlan(owned bool, alter func([]formalBinding) []formalBinding) plan {
	p := plan{Key: "generic", Declaration: f.generic, Associations: f.associations}
	p.User = useFunc(func(_ context.Context, r gen.UseRequest[sourceLanguage, nativeLanguage]) (product, error) {
		return r.Apply()
	})
	p.Acquire = func(context.Context, site) (gen.Destination[nativeLanguage], bool, error) {
		return f.native.destination(owned), owned, nil
	}
	p.Convention = convention{
		introduce: func(_ context.Context, r gen.IntroduceRequest[sourceLanguage, nativeLanguage]) (gen.Introduction[sourceLanguage, nativeLanguage], error) {
			formals, err := r.FreezeLayout(f.native.layout(r.Destination()))
			if err != nil {
				return gen.Introduction[sourceLanguage, nativeLanguage]{}, err
			}
			bindings := make([]formalBinding, 0, len(r.Open().Formals))
			for _, formal := range r.Open().Formals {
				w, err := formals.WantFor(formal.Binder)
				if err != nil {
					return gen.Introduction[sourceLanguage, nativeLanguage]{}, err
				}
				bindings = append(bindings, formalBinding{Binder: formal.Binder, Token: f.associations.issue(w, f.native.fragment(r.Site().Frame()), nil)})
			}
			if alter != nil {
				bindings = alter(bindings)
			}
			return formals.Finish(bindings)
		},
		apply: func(_ context.Context, r gen.ApplyRequest[sourceLanguage, nativeLanguage]) (product, error) {
			return r.Return(f.associations.issue(r.Want(), f.closed, nil))
		},
	}
	p.Definer = defineFunc(func(_ context.Context, r gen.DefineRequest[sourceLanguage, nativeLanguage]) (gen.Definition[sourceLanguage, nativeLanguage], error) {
		return r.CompleteBody(f.closed)
	})
	return p
}

func (f *fixture) plainPlan(owned bool) plan {
	return plan{
		Key: "plain", Declaration: f.plain, Associations: f.associations,
		Acquire: func(context.Context, site) (gen.Destination[nativeLanguage], bool, error) {
			return f.native.destination(owned), owned, nil
		},
		Header: func(_ context.Context, _ site, d gen.Destination[nativeLanguage]) (gen.Layout[nativeLanguage], error) {
			return f.native.layout(d), nil
		},
		Publisher: publishFunc(func(_ context.Context, r gen.PublishRequest[sourceLanguage, nativeLanguage]) (gen.Publication[sourceLanguage, nativeLanguage], error) {
			c, err := r.Acquire()
			if err != nil {
				return gen.Publication[sourceLanguage, nativeLanguage]{}, err
			}
			return r.Expose(c)
		}),
		Definer: defineFunc(func(_ context.Context, r gen.DefineRequest[sourceLanguage, nativeLanguage]) (gen.Definition[sourceLanguage, nativeLanguage], error) {
			return r.CompleteBody(f.closed)
		}),
	}
}

func begin(t *testing.T, ctx context.Context, f *fixture, c config, limits gen.Limits) *session {
	t.Helper()
	s, err := gen.Begin[sourceLanguage, nativeLanguage](ctx, gen.Input{Snapshot: f.source.snapshot}, c, limits)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func failedFinish(t *testing.T, s *session, expected error) {
	t.Helper()
	out, err := s.Finish()
	if !errors.Is(err, expected) || len(out.Trees) != 0 || out.Entry != (ast.NodeRef{}) {
		t.Fatalf("failed invocation published output or lost its failure: output=%v error=%v, want %v", out, err, expected)
	}
}
