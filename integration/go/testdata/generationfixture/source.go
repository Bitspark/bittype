// Package generationfixture is an independent, deliberately bounded consumer of
// the public generation SPI. Its source table is not a Horos implementation.
package generationfixture

import (
	"context"
	"fmt"

	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type SourceLanguage struct{}

func (SourceLanguage) ID() language.ID { return "generation-fixture" }
func (SourceLanguage) Name() string    { return "Independent generation fixture" }

type sourceChecker struct{ SourceLanguage }

func (sourceChecker) Check(v ast.View[SourceLanguage]) error {
	if len(v.ASTView().Nodes) != 1 || v.ASTView().Nodes[0].Kind != "program" {
		return g.ErrInvalid
	}
	return nil
}

type sourceUse struct {
	data g.UseData[SourceLanguage]
	plan string
}

// Source keeps immutable source-owned identities, environments and admission
// records. Selection is an explicit per-use configuration table.
type Source struct {
	declarations       g.Space[g.DeclarationRole[SourceLanguage]]
	binders            g.Space[g.BinderRole[SourceLanguage]]
	uses               g.Space[g.UseRole[SourceLanguage]]
	environments       g.Space[g.EnvironmentRole[SourceLanguage]]
	opens              map[g.Declaration[SourceLanguage]]g.OpenData[SourceLanguage]
	useRecords         map[g.Use[SourceLanguage]]sourceUse
	binderRecords      map[g.Binder[SourceLanguage]]bool
	environmentRecords map[g.Environment[SourceLanguage]]bool
	root               g.Environment[SourceLanguage]
	snapshot           semantic.Snapshot
}

func newSource(ctx context.Context) (*Source, error) {
	s := &Source{
		declarations:       g.NewSpace[g.DeclarationRole[SourceLanguage]](),
		binders:            g.NewSpace[g.BinderRole[SourceLanguage]](),
		uses:               g.NewSpace[g.UseRole[SourceLanguage]](),
		environments:       g.NewSpace[g.EnvironmentRole[SourceLanguage]](),
		opens:              make(map[g.Declaration[SourceLanguage]]g.OpenData[SourceLanguage]),
		useRecords:         make(map[g.Use[SourceLanguage]]sourceUse),
		binderRecords:      make(map[g.Binder[SourceLanguage]]bool),
		environmentRecords: make(map[g.Environment[SourceLanguage]]bool),
	}
	s.root = s.environment()
	language, err := representation.NewLanguage[SourceLanguage](SourceLanguage{}, "fixture-v1")
	if err != nil {
		return nil, err
	}
	checker, err := check.NewChecker[SourceLanguage](language, sourceChecker{})
	if err != nil {
		return nil, err
	}
	tree := ast.AST[SourceLanguage]{ID: "source-fixture", Root: 1, Nodes: []ast.Node[SourceLanguage]{{Kind: "program"}}}
	packed, err := representation.Pack(language, tree)
	if err != nil {
		return nil, err
	}
	s.snapshot, err = semantic.Capture(ctx, checker, semantic.Input{Entry: ast.NodeRef{Tree: tree.ID, Node: 1}, Trees: []representation.AST{packed}})
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Source) environment() g.Environment[SourceLanguage] {
	r, _ := s.environments.Ref(uint64(len(s.environmentRecords) + 1))
	s.environmentRecords[r] = true
	return r
}

func (s *Source) declaration(generic bool) g.Declaration[SourceLanguage] {
	d, _ := s.declarations.Ref(uint64(len(s.opens) + 1))
	o := g.OpenData[SourceLanguage]{Declaration: d, Environment: s.environment()}
	if generic {
		b, _ := s.binders.Ref(uint64(len(s.binderRecords) + 1))
		s.binderRecords[b] = true
		u := s.use(o.Environment, "bound", g.Declaration[SourceLanguage]{}, b)
		o.Formals = []g.Formal[SourceLanguage]{{Binder: b, Subject: u}}
	}
	s.opens[d] = o
	return d
}

func (s *Source) use(env g.Environment[SourceLanguage], plan string, callee g.Declaration[SourceLanguage], bound g.Binder[SourceLanguage], actuals ...g.Actual[SourceLanguage]) g.Use[SourceLanguage] {
	u, _ := s.uses.Ref(uint64(len(s.useRecords) + 1))
	s.useRecords[u] = sourceUse{data: g.UseData[SourceLanguage]{Use: u, Subject: u, Environment: env, Callee: callee, Bound: bound, Actuals: append([]g.Actual[SourceLanguage](nil), actuals...)}, plan: plan}
	return u
}

func (s *Source) call(env g.Environment[SourceLanguage], plan string, decl g.Declaration[SourceLanguage], arg g.Use[SourceLanguage]) g.Use[SourceLanguage] {
	var actuals []g.Actual[SourceLanguage]
	if arg.Valid() {
		actuals = []g.Actual[SourceLanguage]{{Binder: s.opens[decl].Formals[0].Binder, Use: arg}}
	}
	return s.use(env, plan, decl, g.Binder[SourceLanguage]{}, actuals...)
}

// A distinct syntax occurrence can have the same source-certified effective
// subject. Native composition must compare that subject, not occurrence IDs.
func (s *Source) occurrence(original g.Use[SourceLanguage]) g.Use[SourceLanguage] {
	record := s.useRecords[original]
	u, _ := s.uses.Ref(uint64(len(s.useRecords) + 1))
	record.data = record.data.Clone()
	record.data.Use = u
	s.useRecords[u] = record
	return u
}

func (s *Source) body(d g.Declaration[SourceLanguage], u g.Use[SourceLanguage]) {
	o := s.opens[d]
	o.Body = u
	s.opens[d] = o
}

func (s *Source) CheckInput(ctx context.Context, snapshot semantic.Snapshot, evidence []semantic.Result) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !snapshot.Same(s.snapshot) || len(evidence) != 0 {
		return g.ErrScope
	}
	return nil
}

func (s *Source) Owners() g.SourceOwners[SourceLanguage] {
	return g.SourceOwners[SourceLanguage]{Declarations: s.declarations.Owner(), Binders: s.binders.Owner(), Uses: s.uses.Owner(), Environments: s.environments.Owner()}
}

func (s *Source) RootEnvironment() g.Environment[SourceLanguage] { return s.root }

func (s *Source) Authenticate(ctx context.Context, r g.References[SourceLanguage]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for _, d := range r.Declarations {
		if _, ok := s.opens[d]; !ok {
			return g.ErrScope
		}
	}
	for _, b := range r.Binders {
		if !s.binderRecords[b] {
			return g.ErrScope
		}
	}
	for _, u := range r.Uses {
		if _, ok := s.useRecords[u]; !ok {
			return g.ErrScope
		}
	}
	for _, e := range r.Environments {
		if !s.environmentRecords[e] {
			return g.ErrScope
		}
	}
	return nil
}

func (s *Source) Open(ctx context.Context, d g.Declaration[SourceLanguage]) (g.OpenData[SourceLanguage], error) {
	if err := ctx.Err(); err != nil {
		return g.OpenData[SourceLanguage]{}, err
	}
	o, ok := s.opens[d]
	if !ok {
		return g.OpenData[SourceLanguage]{}, g.ErrScope
	}
	return o.Clone(), nil
}

func (s *Source) Admit(ctx context.Context, env g.Environment[SourceLanguage], u g.Use[SourceLanguage]) (g.UseData[SourceLanguage], error) {
	if err := ctx.Err(); err != nil {
		return g.UseData[SourceLanguage]{}, err
	}
	r, ok := s.useRecords[u]
	if !ok || r.data.Environment != env {
		return g.UseData[SourceLanguage]{}, g.ErrScope
	}
	// Whole applications are admitted before a mapper chooses which operands to
	// represent. This table contains only explicitly constructed well-formed uses.
	for _, a := range r.data.Actuals {
		actual, exists := s.useRecords[a.Use]
		if !exists || actual.data.Environment != env {
			return g.UseData[SourceLanguage]{}, g.ErrScope
		}
	}
	return r.data.Clone(), nil
}

func (s *Source) selectPlan(u g.UseData[SourceLanguage]) (string, error) {
	r, ok := s.useRecords[u.Use]
	if !ok {
		return "", fmt.Errorf("unknown fixture use: %w", g.ErrScope)
	}
	return r.plan, nil
}
