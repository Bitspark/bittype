package generation

import (
	"context"
	"slices"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type planState[F, T any] struct { value Plan[F, T]; issuer IssuerID; callee *calleeState[F, T] }

// Session is an explicitly driven, preselected invocation, not a selector or
// definition scheduler. It supports nested callbacks on one call stack; callers
// must serialize access to a Session. Distinct sessions may run concurrently if
// their supplied collaborators support it. Finish is the only output commit.
type Session[F, T any] struct {
	ctx context.Context
	id RunID
	configuration ConfigID
	snapshot semantic.Snapshot
	source Source[F]
	owners SourceOwners[F]
	native Native[F, T]
	target representation.Language
	selectPlan func(UseData[F]) (string, error)
	plans map[string]*planState[F, T]
	declarations map[Declaration[F]]*planState[F, T]
	opens map[Declaration[F]]OpenData[F]
	frames map[FrameID]map[Binder[F]]Product[F, T]
	root Environment[F]
	frame FrameID
	limits Limits
	operations, depth, references, reservations int
	publications []Publication[F, T]
	definitions []Definition[F, T]
	failed error
	finished bool
}

func Begin[F, T any](ctx context.Context, input Input, config Config[F, T], limits Limits) (*Session[F, T], error) {
	if ctx == nil || !limits.Valid() || representation.IsNil(config.Source) || representation.IsNil(config.Native) || config.Select == nil || len(config.Plans) == 0 || len(config.Plans) > limits.Plans { return nil, ErrInvalid }
	if err := ctx.Err(); err != nil { return nil, err }
	if err := representation.CheckMarker[F](input.Snapshot.Language()); err != nil { return nil, err }
	if err := representation.CheckMarker[T](config.Target); err != nil { return nil, err }
	// Requirements are receiver-owned; this is admission, never provider dispatch.
	evidence, err := semantic.AdmitEvidence(input.Snapshot, config.Evidence, input.Evidence)
	if err != nil { return nil, err }
	if err := config.Source.CheckInput(ctx, input.Snapshot, evidence); err != nil { return nil, err }
	s := &Session[F, T]{ctx: ctx, id: RunID{fresh()}, configuration: ConfigID{fresh()}, snapshot: input.Snapshot,
		source: config.Source, owners: config.Source.Owners(), native: config.Native, target: config.Target, selectPlan: config.Select,
		plans: make(map[string]*planState[F, T]), declarations: make(map[Declaration[F]]*planState[F, T]), opens: make(map[Declaration[F]]OpenData[F]),
		frames: make(map[FrameID]map[Binder[F]]Product[F, T]), root: config.Source.RootEnvironment(), frame: FrameID{fresh()}, limits: limits}
	if !s.owners.Valid() { return nil, ErrInvalid }
	if err := s.authenticate(References[F]{Environments: []Environment[F]{s.root}}); err != nil { return nil, err }
	for _, p := range config.Plans {
		if p.Key == "" || s.plans[p.Key] != nil { return nil, ErrInvalid }
		ps := &planState[F, T]{value: p, issuer: IssuerID{fresh()}}
		s.plans[p.Key] = ps
		if p.Declaration.Valid() {
			if err := s.authenticate(References[F]{Declarations: []Declaration[F]{p.Declaration}}); err != nil { return nil, err }
			if s.declarations[p.Declaration] != nil { return nil, ErrInvalid }
			s.declarations[p.Declaration] = ps
		}
	}
	if err := s.check(); err != nil { return nil, err }
	return s, nil
}

func (s *Session[F, T]) fail(err error) error {
	if err == nil { err = ErrInvalid }
	if s != nil && s.failed == nil { s.failed = err }
	if s == nil { return err }
	return s.failed
}
func (s *Session[F, T]) check() error {
	if s == nil { return ErrInvalid }
	if s.failed != nil { return s.failed }
	if s.finished { return ErrInvalid }
	if err := s.ctx.Err(); err != nil { return s.fail(err) }
	return nil
}
func (s *Session[F, T]) step() error {
	if err := s.check(); err != nil { return err }
	if s.operations >= s.limits.Operations { return s.fail(ErrLimit) }
	s.operations++
	return nil
}
func (s *Session[F, T]) enter(p *planState[F, T], frame FrameID, env Environment[F], c *calleeState[F, T]) (*operation[F, T], error) {
	if err := s.step(); err != nil { return nil, err }
	if s.depth >= s.limits.Depth { return nil, s.fail(ErrLimit) }
	s.depth++
	return &operation[F, T]{run: s, plan: p, ticket: Ticket{fresh()}, frame: frame, environment: env, callee: c, live: true}, nil
}
func (o *operation[F, T]) close() { o.live = false; o.run.depth-- }
func (o *operation[F, T]) check() error {
	if o == nil { return ErrInvalid }
	if err := o.run.check(); err != nil { return err }
	if !o.live || o.completed { return o.run.fail(ErrInvalid) }
	return nil
}
func (s *Session[F, T]) authenticate(r References[F]) error {
	if err := s.step(); err != nil { return err }
	count := len(r.Declarations) + len(r.Binders) + len(r.Uses) + len(r.Environments)
	if count > s.limits.References - s.references { return s.fail(ErrLimit) }
	s.references += count
	for _, v := range r.Declarations { if !s.owners.Declarations.Owns(v) { return s.fail(ErrScope) } }
	for _, v := range r.Binders { if !s.owners.Binders.Owns(v) { return s.fail(ErrScope) } }
	for _, v := range r.Uses { if !s.owners.Uses.Owns(v) { return s.fail(ErrScope) } }
	for _, v := range r.Environments { if !s.owners.Environments.Owns(v) { return s.fail(ErrScope) } }
	if err := s.source.Authenticate(s.ctx, r); err != nil { return s.fail(err) }
	return s.check()
}
func (s *Session[F, T]) open(d Declaration[F]) (OpenData[F], error) {
	if err := s.authenticate(References[F]{Declarations: []Declaration[F]{d}}); err != nil { return OpenData[F]{}, err }
	if o, ok := s.opens[d]; ok { return o.Clone(), nil }
	o, err := s.source.Open(s.ctx, d)
	if err != nil { return OpenData[F]{}, s.fail(err) }
	if err := s.check(); err != nil { return OpenData[F]{}, err }
	if o.Declaration != d || len(o.Formals) > s.limits.Formals { return OpenData[F]{}, s.fail(ErrInvalid) }
	o = o.Clone()
	r := References[F]{Environments: []Environment[F]{o.Environment}, Uses: []Use[F]{o.Body}}
	seen := make(map[Binder[F]]bool, len(o.Formals))
	for _, f := range o.Formals {
		if seen[f.Binder] { return OpenData[F]{}, s.fail(ErrInvalid) }; seen[f.Binder] = true
		r.Binders = append(r.Binders, f.Binder); r.Uses = append(r.Uses, f.Subject)
	}
	if err := s.authenticate(r); err != nil { return OpenData[F]{}, err }
	s.opens[d] = o
	return o.Clone(), nil
}
func (s *Session[F, T]) admit(env Environment[F], use Use[F]) (UseData[F], error) {
	if err := s.authenticate(References[F]{Environments: []Environment[F]{env}, Uses: []Use[F]{use}}); err != nil { return UseData[F]{}, err }
	u, err := s.source.Admit(s.ctx, env, use)
	if err != nil { return UseData[F]{}, s.fail(err) }
	if err := s.check(); err != nil { return UseData[F]{}, err }
	if u.Use != use || u.Environment != env || len(u.Actuals) > s.limits.Formals || u.Bound.Valid() && (u.Callee.Valid() || len(u.Actuals) != 0) { return UseData[F]{}, s.fail(ErrInvalid) }
	u = u.Clone()
	r := References[F]{Uses: []Use[F]{u.Subject}}
	if u.Callee.Valid() { r.Declarations = append(r.Declarations, u.Callee) }
	if u.Bound.Valid() { r.Binders = append(r.Binders, u.Bound) }
	seen := make(map[Binder[F]]bool, len(u.Actuals))
	for _, a := range u.Actuals {
		if seen[a.Binder] { return UseData[F]{}, s.fail(ErrInvalid) }; seen[a.Binder] = true
		r.Binders = append(r.Binders, a.Binder); r.Uses = append(r.Uses, a.Use)
	}
	if err := s.authenticate(r); err != nil { return UseData[F]{}, err }
	if u.Callee.Valid() {
		o, err := s.open(u.Callee); if err != nil { return UseData[F]{}, err }
		if len(o.Formals) != len(u.Actuals) { return UseData[F]{}, s.fail(ErrInvalid) }
		for _, f := range o.Formals { if !seen[f.Binder] { return UseData[F]{}, s.fail(ErrScope) } }
	} else if len(u.Actuals) != 0 { return UseData[F]{}, s.fail(ErrInvalid) }
	return u, nil
}

func (s *Session[F, T]) Use(u Use[F]) (Product[F, T], error) {
	if err := s.check(); err != nil { return Product[F, T]{}, err }
	return s.use(s.root, s.frame, u)
}
func (s *Session[F, T]) use(env Environment[F], frame FrameID, u Use[F]) (Product[F, T], error) {
	data, err := s.admit(env, u); if err != nil { return Product[F, T]{}, err }
	key, err := s.selectPlan(data.Clone()); if err != nil { return Product[F, T]{}, s.fail(err) }
	if err := s.check(); err != nil { return Product[F, T]{}, err }
	p := s.plans[key]
	if p == nil || representation.IsNil(p.value.User) { return Product[F, T]{}, s.fail(ErrUnsupported) }
	if p.value.Declaration.Valid() && p.value.Declaration != data.Callee { return Product[F, T]{}, s.fail(ErrScope) }
	o, err := s.enter(p, frame, env, nil); if err != nil { return Product[F, T]{}, err }; defer o.close()
	r := &useRequest[F, T]{operation: o, data: data}
	out, err := p.value.User.Use(s.ctx, r)
	if err != nil { return Product[F, T]{}, s.fail(err) }
	if err := s.check(); err != nil { return Product[F, T]{}, err }
	if !o.completed || !out.Valid() || out.state.ticket != o.ticket { return Product[F, T]{}, s.fail(ErrInvalid) }
	return out, nil
}

// Finish checks the obligation ledger and returns an owned complete AST set.
// All failures return the zero set. The target owns grammar/expansion checks;
// the host independently bounds the actual returned output and endpoint marker.
func (s *Session[F, T]) Finish() (representation.ASTSet, error) {
	if err := s.step(); err != nil { return representation.ASTSet{}, err }
	if s.depth != 0 { return representation.ASTSet{}, s.fail(ErrInvalid) }
	for _, p := range s.plans {
		if c := p.callee; c != nil && (c.phase < ready || c.owned && c.phase != defined) { return representation.ASTSet{}, s.fail(ErrIncomplete) }
	}
	output, err := s.native.Assemble(s.ctx, slices.Clone(s.publications), slices.Clone(s.definitions), s.limits)
	if err != nil { return representation.ASTSet{}, s.fail(err) }
	if err := s.check(); err != nil { return representation.ASTSet{}, err }
	if len(output.Trees) == 0 || len(output.Trees) > 256 { return representation.ASTSet{}, s.fail(ErrInvalid) }
	nodes, bytes := 0, 0
	for _, envelope := range output.Trees {
		info, err := envelope.Inspect(); if err != nil { return representation.ASTSet{}, s.fail(err) }
		if info.Nodes > s.limits.OutputNodes - nodes || info.TextBytes > s.limits.OutputBytes - bytes { return representation.ASTSet{}, s.fail(ErrLimit) }
		nodes += info.Nodes; bytes += info.TextBytes
		if _, err := representation.Unpack[T](s.target, envelope); err != nil { return representation.ASTSet{}, s.fail(err) }
	}
	owned, err := output.Snapshot(); if err != nil { return representation.ASTSet{}, s.fail(err) }
	if err := s.check(); err != nil { return representation.ASTSet{}, err }
	s.finished = true
	return owned, nil
}
