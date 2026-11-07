package source

import (
	"context"
	"fmt"
	analyze "github.com/Bitspark/bittype/analyze/files/horos/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	resolve "github.com/Bitspark/bittype/resolve/files/horos/provider/go"
	"github.com/Bitspark/bittype/semantics/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"slices"
)

func Contract() transform.Contract {
	return transform.Contract{Evidence: []semantic.EvidenceRequirement{{Operation: analysis.Kind(), Profile: analyze.Profile, Reports: analyze.Contract().Semantic.SemanticsRequirements, Semantics: true}}}
}

type Source struct {
	snapshot       semantic.Snapshot
	evidence       semantic.Result
	facts          semantics.Data
	declarations   g.Space[g.DeclarationRole[h.Language]]
	binders        g.Space[g.BinderRole[h.Language]]
	uses           g.Space[g.UseRole[h.Language]]
	environments   g.Space[g.EnvironmentRole[h.Language]]
	entities       map[uint64]semantics.EntityID
	decls          map[semantics.EntityID]g.Declaration[h.Language]
	bound          map[semantics.EntityID]g.Binder[h.Language]
	formalSubjects map[semantics.EntityID]g.Use[h.Language]
	sorts          map[semantics.Subject]semantics.TermID
	results        map[semantics.EntityID]semantics.TermID
	identities     map[semantics.EntityID]semantics.TermID
	sites          map[semantics.TermID]ast.NodeRef
	termOwners     map[semantics.TermID]semantics.EntityID
	order          []g.Declaration[h.Language]
	bySite         map[ast.NodeRef]g.Declaration[h.Language]
	roots          map[g.Declaration[h.Language]]bool
	rootModule     ast.NodeRef
}

func New(ctx context.Context, r transform.Request[h.Language]) (*Source, error) {
	if ctx == nil {
		return nil, g.ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	input := r.Source()
	if r.Entry() != (ast.NodeRef{Tree: input.ID, Node: input.Root}) {
		return nil, semantic.ErrInvalidRequest
	}
	e, ok := r.Evidence(0)
	if !ok {
		return nil, semantic.ErrEvidenceUnavailable
	}
	model, ok := e.Semantics()
	if !ok {
		return nil, semantic.ErrEvidenceUnavailable
	}
	s := &Source{bySite: map[ast.NodeRef]g.Declaration[h.Language]{}, roots: map[g.Declaration[h.Language]]bool{}, snapshot: r.Snapshot(), evidence: e, facts: model.Data(), declarations: g.NewSpace[g.DeclarationRole[h.Language]](), binders: g.NewSpace[g.BinderRole[h.Language]](), uses: g.NewSpace[g.UseRole[h.Language]](), environments: g.NewSpace[g.EnvironmentRole[h.Language]](), entities: map[uint64]semantics.EntityID{}, decls: map[semantics.EntityID]g.Declaration[h.Language]{}, bound: map[semantics.EntityID]g.Binder[h.Language]{}, formalSubjects: map[semantics.EntityID]g.Use[h.Language]{}, sorts: map[semantics.Subject]semantics.TermID{}, results: map[semantics.EntityID]semantics.TermID{}, identities: map[semantics.EntityID]semantics.TermID{}, sites: map[semantics.TermID]ast.NodeRef{}, termOwners: map[semantics.TermID]semantics.EntityID{}}
	for _, ref := range input.Nodes[input.Root-1].References {
		if ref.Role == h.RoleRoot {
			s.rootModule = ref.Target
		}
	}
	if s.rootModule.Tree == "" {
		return nil, semantic.ErrInvalidRequest
	}
	for _, e := range s.facts.Expressions {
		if e.At.Role == resolve.ExpressionRole() {
			s.sites[e.Term] = e.At.Site
		}
	}
	for _, f := range s.facts.Facts {
		switch f.Relation {
		case resolve.SortRelation():
			if s.sorts[f.Subject] != 0 {
				return nil, semantic.ErrInvalidRequest
			}
			s.sorts[f.Subject] = f.Term
		case resolve.ResultSortRelation():
			if f.Subject.Entity == 0 || s.results[f.Subject.Entity] != 0 {
				return nil, semantic.ErrInvalidRequest
			}
			s.results[f.Subject.Entity] = f.Term
		case resolve.IdentityRelation():
			if f.Subject.Entity == 0 || s.identities[f.Subject.Entity] != 0 {
				return nil, semantic.ErrInvalidRequest
			}
			s.identities[f.Subject.Entity] = f.Term
		}
	}
	for i, e := range s.facts.Entities {
		id := semantics.EntityID(i + 1)
		if e.Kind != (semantics.Key{Owner: "horos", Name: "type-definition"}) {
			continue
		}
		if e.Owner == 0 || e.Body == 0 || len(e.Declarations) != 1 {
			return nil, semantic.ErrInvalidRequest
		}
		d, _ := s.declarations.Ref(uint64(id))
		s.entities[uint64(id)] = id
		s.decls[id] = d
		site := e.Declarations[0]
		if _, exists := s.bySite[site]; exists {
			return nil, semantic.ErrInvalidRequest
		}
		s.bySite[site] = d
		s.order = append(s.order, d)
		owner := s.facts.Entities[e.Owner-1]
		if len(owner.Declarations) != 1 {
			return nil, semantic.ErrInvalidRequest
		}
		if owner.Declarations[0] == s.rootModule {
			s.roots[d] = true
		}
		for _, p := range e.Parameters {
			b, _ := s.binders.Ref(uint64(p.Entity))
			s.bound[p.Entity] = b
			u, _ := s.uses.Ref(uint64(len(s.facts.Terms)) + uint64(p.Entity))
			s.formalSubjects[p.Entity] = u
		}
		// The graph owns lexical term sites. Callee references are not expanded.
		var walk func(semantics.TermID, int) error
		walk = func(t semantics.TermID, depth int) error {
			if err := ctx.Err(); err != nil {
				return err
			}
			if depth > 512 {
				return g.ErrLimit
			}
			if old := s.termOwners[t]; old != 0 {
				if old != id {
					return semantic.ErrInvalidRequest
				}
				return nil
			}
			s.termOwners[t] = id
			term := s.facts.Terms[t-1]
			for _, a := range term.Arguments {
				if err := walk(a, depth+1); err != nil {
					return err
				}
			}
			for _, a := range term.Operands {
				if err := walk(a.Term, depth+1); err != nil {
					return err
				}
			}
			return nil
		}
		if err := walk(e.Body, 0); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Source) Snapshot() semantic.Snapshot               { return s.snapshot }
func (s *Source) Declarations() []g.Declaration[h.Language] { return slices.Clone(s.order) }
func (s *Source) Owners() g.SourceOwners[h.Language] {
	return g.SourceOwners[h.Language]{Declarations: s.declarations.Owner(), Binders: s.binders.Owner(), Uses: s.uses.Owner(), Environments: s.environments.Owner()}
}
func (s *Source) RootEnvironment() g.Environment[h.Language] { r, _ := s.environments.Ref(1); return r }
func (s *Source) CheckInput(ctx context.Context, snapshot semantic.Snapshot, evidence []semantic.Result) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !s.snapshot.Same(snapshot) || len(evidence) != 1 || evidence[0] != s.evidence {
		return g.ErrScope
	}
	return nil
}
func (s *Source) Authenticate(ctx context.Context, r g.References[h.Language]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for _, d := range r.Declarations {
		i, ok := s.declarations.Index(d)
		if !ok || s.entities[i] == 0 {
			return g.ErrScope
		}
	}
	for _, b := range r.Binders {
		i, ok := s.binders.Index(b)
		if !ok || s.bound[semantics.EntityID(i)] != b {
			return g.ErrScope
		}
	}
	for _, u := range r.Uses {
		if _, _, err := s.useIndex(u); err != nil {
			return err
		}
	}
	for _, e := range r.Environments {
		i, ok := s.environments.Index(e)
		if !ok || (i != 1 && s.entities[i-1] == 0) {
			return g.ErrScope
		}
	}
	return nil
}
func (s *Source) Open(ctx context.Context, d g.Declaration[h.Language]) (g.OpenData[h.Language], error) {
	if err := s.Authenticate(ctx, g.References[h.Language]{Declarations: []g.Declaration[h.Language]{d}}); err != nil {
		return g.OpenData[h.Language]{}, err
	}
	i, _ := s.declarations.Index(d)
	e := s.facts.Entities[i-1]
	if err := s.proper(s.results[semantics.EntityID(i)], e.Body); err != nil {
		return g.OpenData[h.Language]{}, err
	}
	env, _ := s.environments.Ref(i + 1)
	o := g.OpenData[h.Language]{Declaration: d, Environment: env, Body: s.use(e.Body)}
	for _, p := range e.Parameters {
		if err := s.proper(s.sorts[semantics.Subject{Entity: p.Entity}], e.Body); err != nil {
			return g.OpenData[h.Language]{}, err
		}
		o.Formals = append(o.Formals, g.Formal[h.Language]{Binder: s.bound[p.Entity], Subject: s.formalSubjects[p.Entity]})
	}
	return o, nil
}
func (s *Source) Admit(ctx context.Context, env g.Environment[h.Language], use g.Use[h.Language]) (g.UseData[h.Language], error) {
	u := g.UseData[h.Language]{Use: use, Subject: use, Environment: env}
	if err := s.Authenticate(ctx, g.References[h.Language]{Uses: []g.Use[h.Language]{use}, Environments: []g.Environment[h.Language]{env}}); err != nil {
		return u, err
	}
	ti, formal, err := s.useIndex(use)
	if err != nil {
		return u, err
	}
	e, _ := s.environments.Index(env)
	if formal != 0 {
		if s.facts.Entities[formal-1].Owner != semantics.EntityID(e-1) {
			return u, g.ErrScope
		}
		u.Bound = s.bound[formal]
		return u, nil
	}
	id := semantics.TermID(ti)
	if s.termOwners[id] != semantics.EntityID(e-1) {
		return u, g.ErrScope
	}
	if err := s.proper(s.sorts[semantics.Subject{Term: id}], id); err != nil {
		return u, err
	}
	t := s.facts.Terms[id-1]
	if t.Kind == semantics.Reference || t.Kind == semantics.Application {
		target := t.Target
		if t.Kind == semantics.Application {
			callee := s.facts.Terms[t.Callee-1]
			if callee.Kind != semantics.Reference {
				return u, s.Unsupported(use, "computed constructor application")
			}
			target = callee.Target
		}
		if b := s.bound[target]; b.Valid() {
			if t.Kind != semantics.Reference || s.facts.Entities[target-1].Owner != semantics.EntityID(e-1) {
				return u, g.ErrScope
			}
			u.Bound = b
			u.Subject = s.formalSubjects[target]
			return u, nil
		}
		d := s.decls[target]
		if !d.Valid() {
			return u, s.Unsupported(use, "constructor-valued application")
		}
		u.Callee = d
		parameters := s.facts.Entities[target-1].Parameters
		if len(parameters) != len(t.Arguments) {
			return u, semantic.ErrInvalidRequest
		}
		for i, p := range parameters {
			u.Actuals = append(u.Actuals, g.Actual[h.Language]{Binder: s.bound[p.Entity], Use: s.use(t.Arguments[i])})
		}
	}
	// The exact successful Horos analysis admits the entire original occurrence,
	// including every actual, under this declaration's original assumptions.
	return u, nil
}
func (s *Source) use(id semantics.TermID) g.Use[h.Language] { r, _ := s.uses.Ref(uint64(id)); return r }
func (s *Source) useIndex(u g.Use[h.Language]) (uint64, semantics.EntityID, error) {
	i, ok := s.uses.Index(u)
	if !ok {
		return 0, 0, g.ErrScope
	}
	if i <= uint64(len(s.facts.Terms)) && s.termOwners[semantics.TermID(i)] != 0 {
		return i, 0, nil
	}
	b := semantics.EntityID(i - uint64(len(s.facts.Terms)))
	if s.formalSubjects[b] == u {
		return i, b, nil
	}
	return 0, 0, g.ErrScope
}
func (s *Source) Unsupported(use g.Use[h.Language], message string) error {
	i, _, _ := s.useIndex(use)
	site := s.sites[semantics.TermID(i)]
	return &text.Diagnostic{Code: text.Unsupported, Message: message, Node: &site}
}
func (s *Source) proper(sortID, term semantics.TermID) error {
	if sortID == 0 || int(sortID) > len(s.facts.Terms) {
		return semantic.ErrInvalidRequest
	}
	t := s.facts.Terms[sortID-1]
	if t.Kind != semantics.ProviderExpression || t.Operator != (semantics.Key{Owner: "horos", Name: "sort/proper"}) {
		return s.Unsupported(s.use(term), "constructor-valued types require an explicit representation")
	}
	return nil
}
func (s *Source) String() string { return fmt.Sprintf("Horos source (%d declarations)", len(s.order)) }
