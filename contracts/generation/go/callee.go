package generation

import "github.com/Bitspark/bittype/trees/representation/go"

func (s *Session[F, T]) Acquire(d Declaration[F]) (Callee[F, T], error) {
	if err := s.check(); err != nil {
		return Callee[F, T]{}, err
	}
	if err := s.authenticate(References[F]{Declarations: []Declaration[F]{d}}); err != nil {
		return Callee[F, T]{}, err
	}
	p := s.declarations[d]
	if p == nil {
		return Callee[F, T]{}, s.fail(ErrUnsupported)
	}
	return s.acquire(p)
}

// AcquireHelper selects an explicitly supplied native-only plan. It creates no
// source declaration/binder and inherits no caller's source environment.
func (s *Session[F, T]) AcquireHelper(key string) (Callee[F, T], error) {
	if err := s.check(); err != nil {
		return Callee[F, T]{}, err
	}
	p := s.plans[key]
	if p == nil {
		return Callee[F, T]{}, s.fail(ErrUnsupported)
	}
	if p.value.Declaration.Valid() {
		return Callee[F, T]{}, s.fail(ErrInvalid)
	}
	c, err := s.acquire(p)
	if err != nil {
		return Callee[F, T]{}, err
	}
	if !c.Owned() {
		return Callee[F, T]{}, s.fail(ErrInvalid)
	}
	return c, nil
}
func (s *Session[F, T]) acquire(p *planState[F, T]) (Callee[F, T], error) {
	if err := s.step(); err != nil {
		return Callee[F, T]{}, err
	}
	if p.callee != nil {
		if p.callee.phase < ready {
			return Callee[F, T]{}, s.fail(ErrInvalid)
		}
		return Callee[F, T]{p.callee}, nil
	}
	if p.value.Acquire == nil {
		return Callee[F, T]{}, s.fail(ErrUnsupported)
	}
	if s.reservations >= s.limits.Reservations {
		return Callee[F, T]{}, s.fail(ErrLimit)
	}
	s.reservations++
	c := &calleeState[F, T]{run: s, plan: p, frame: FrameID{fresh()}, signature: SignatureID{fresh()}, phase: reserved, formals: make(map[Binder[F]]Product[F, T])}
	if p.value.Declaration.Valid() {
		o, err := s.open(p.value.Declaration)
		if err != nil {
			return Callee[F, T]{}, err
		}
		c.open = o
	}
	p.callee = c
	o, err := s.enter(p, c.frame, c.open.Environment, c)
	if err != nil {
		return Callee[F, T]{}, err
	}
	defer o.close()
	site := Site[F, T]{o}
	z, owned, err := p.value.Acquire(s.ctx, site)
	if err != nil {
		return Callee[F, T]{}, s.fail(err)
	}
	if err := o.check(); err != nil {
		return Callee[F, T]{}, err
	}
	if !z.Valid() {
		return Callee[F, T]{}, s.fail(ErrInvalid)
	}
	if err := s.native.CheckDestination(s.ctx, site, z, owned); err != nil {
		return Callee[F, T]{}, s.fail(err)
	}
	if err := o.check(); err != nil {
		return Callee[F, T]{}, err
	}
	c.destination, c.owned = z, owned
	c.phase = introducing
	if !representation.IsNil(p.value.Convention) {
		// The same ticket cannot issue after acquisition. Introduction has its
		// own operation lifetime but shares the retained canonical frame/layout.
		introOp, err := s.enter(p, c.frame, c.open.Environment, c)
		if err != nil {
			return Callee[F, T]{}, err
		}
		intro, err := p.value.Convention.Introduce(s.ctx, &introduceRequest[F, T]{operation: introOp, callee: c})
		introOp.close()
		if err != nil {
			return Callee[F, T]{}, s.fail(err)
		}
		if err := s.check(); err != nil {
			return Callee[F, T]{}, err
		}
		if !introOp.completed || !intro.Valid() || intro.state.ticket != introOp.ticket || intro.state.callee != c || !c.layout.Valid() {
			return Callee[F, T]{}, s.fail(ErrInvalid)
		}
		c.formals = intro.state.formals
	} else {
		if len(c.open.Formals) != 0 || p.value.Header == nil {
			return Callee[F, T]{}, s.fail(ErrUnsupported)
		}
		h, err := p.value.Header(s.ctx, site, z)
		if err != nil {
			return Callee[F, T]{}, s.fail(err)
		}
		if err := o.check(); err != nil {
			return Callee[F, T]{}, err
		}
		if !h.Valid() {
			return Callee[F, T]{}, s.fail(ErrInvalid)
		}
		if err := s.native.CheckLayout(s.ctx, site, z, h); err != nil {
			return Callee[F, T]{}, s.fail(err)
		}
		c.layout = h
	}
	if err := o.check(); err != nil {
		return Callee[F, T]{}, err
	}
	s.frames[c.frame] = c.formals
	c.phase = ready
	return Callee[F, T]{c}, nil
}

type introduceRequest[F, T any] struct {
	operation *operation[F, T]
	callee    *calleeState[F, T]
	frozen    bool
}

func (r *introduceRequest[F, T]) Site() Site[F, T]            { return Site[F, T]{r.operation} }
func (r *introduceRequest[F, T]) Open() OpenData[F]           { return r.callee.open.Clone() }
func (r *introduceRequest[F, T]) Destination() Destination[T] { return r.callee.destination }
func (r *introduceRequest[F, T]) FreezeLayout(h Layout[T]) (FormalRequest[F, T], error) {
	o := r.operation
	if err := o.check(); err != nil {
		return nil, err
	}
	if r.frozen || !h.Valid() {
		return nil, o.run.fail(ErrInvalid)
	}
	if err := o.run.native.CheckLayout(o.run.ctx, r.Site(), r.callee.destination, h); err != nil {
		return nil, o.run.fail(err)
	}
	if err := o.check(); err != nil {
		return nil, err
	}
	r.frozen = true
	r.callee.layout = h
	return &formalRequest[F, T]{operation: o, callee: r.callee}, nil
}

type formalRequest[F, T any] struct {
	operation *operation[F, T]
	callee    *calleeState[F, T]
}

func (r *formalRequest[F, T]) WantFor(b Binder[F]) (Want[F, T], error) {
	if err := r.operation.check(); err != nil {
		return Want[F, T]{}, err
	}
	for _, f := range r.callee.open.Formals {
		if f.Binder == b {
			return Want[F, T]{site: Site[F, T]{r.operation}, use: UseData[F]{Use: f.Subject, Subject: f.Subject, Environment: r.callee.open.Environment, Bound: b}}, nil
		}
	}
	return Want[F, T]{}, r.operation.run.fail(ErrScope)
}
func (r *formalRequest[F, T]) Finish(bindings []FormalBinding[F, T]) (Introduction[F, T], error) {
	o, s := r.operation, r.operation.run
	if err := o.check(); err != nil {
		return Introduction[F, T]{}, err
	}
	if len(bindings) != len(r.callee.open.Formals) {
		return Introduction[F, T]{}, s.fail(ErrInvalid)
	}
	products := make(map[Binder[F]]Product[F, T], len(bindings))
	for _, b := range bindings {
		if _, exists := products[b.Binder]; exists {
			return Introduction[F, T]{}, s.fail(ErrInvalid)
		}
		want, err := r.WantFor(b.Binder)
		if err != nil {
			return Introduction[F, T]{}, err
		}
		p, err := s.accept(o, want, b.Token)
		if err != nil {
			return Introduction[F, T]{}, err
		}
		products[b.Binder] = p
	}
	if err := o.check(); err != nil {
		return Introduction[F, T]{}, err
	}
	o.completed = true
	return Introduction[F, T]{&introductionState[F, T]{callee: r.callee, ticket: o.ticket, formals: products}}, nil
}

func (s *Session[F, T]) Publish(root Declaration[F]) (Publication[F, T], error) {
	if err := s.authenticate(References[F]{Declarations: []Declaration[F]{root}}); err != nil {
		return Publication[F, T]{}, err
	}
	p := s.declarations[root]
	if p == nil || representation.IsNil(p.value.Publisher) {
		return Publication[F, T]{}, s.fail(ErrUnsupported)
	}
	for _, previous := range s.publications {
		if previous.Root() == root {
			return Publication[F, T]{}, s.fail(ErrInvalid)
		}
	}
	o, err := s.enter(p, s.frame, s.root, nil)
	if err != nil {
		return Publication[F, T]{}, err
	}
	defer o.close()
	out, err := p.value.Publisher.Publish(s.ctx, &publishRequest[F, T]{operation: o, root: root})
	if err != nil {
		return Publication[F, T]{}, s.fail(err)
	}
	if err := s.check(); err != nil {
		return Publication[F, T]{}, err
	}
	if !o.completed || !out.Valid() || out.state.ticket != o.ticket || out.Root() != root {
		return Publication[F, T]{}, s.fail(ErrInvalid)
	}
	s.publications = append(s.publications, out)
	return out, nil
}

type publishRequest[F, T any] struct {
	operation *operation[F, T]
	root      Declaration[F]
}

func (r *publishRequest[F, T]) Site() Site[F, T]     { return Site[F, T]{r.operation} }
func (r *publishRequest[F, T]) Root() Declaration[F] { return r.root }
func (r *publishRequest[F, T]) Acquire() (Callee[F, T], error) {
	if err := r.operation.check(); err != nil {
		return Callee[F, T]{}, err
	}
	return r.operation.run.Acquire(r.root)
}
func (r *publishRequest[F, T]) Expose(c Callee[F, T]) (Publication[F, T], error) {
	o, s := r.operation, r.operation.run
	if err := o.check(); err != nil {
		return Publication[F, T]{}, err
	}
	if !c.Valid() || c.state.run != s || c.state.phase < ready || c.state.plan != o.plan || c.state.open.Declaration != r.root {
		return Publication[F, T]{}, s.fail(ErrScope)
	}
	if err := s.native.CheckExposure(s.ctx, r.Site(), r.root, c); err != nil {
		return Publication[F, T]{}, s.fail(err)
	}
	if err := o.check(); err != nil {
		return Publication[F, T]{}, err
	}
	o.completed = true
	return Publication[F, T]{&publicationState[F, T]{callee: c, root: r.root, ticket: o.ticket}}, nil
}

func (s *Session[F, T]) Define(c Callee[F, T]) (Definition[F, T], error) {
	if err := s.step(); err != nil {
		return Definition[F, T]{}, err
	}
	if !c.Valid() || c.state.run != s || !c.Owned() || c.state.phase != ready {
		return Definition[F, T]{}, s.fail(ErrInvalid)
	}
	p := c.state.plan
	if representation.IsNil(p.value.Definer) {
		return Definition[F, T]{}, s.fail(ErrUnsupported)
	}
	o, err := s.enter(p, c.Frame(), c.state.open.Environment, c.state)
	if err != nil {
		return Definition[F, T]{}, err
	}
	defer o.close()
	c.state.phase = defining
	out, err := p.value.Definer.Define(s.ctx, &defineRequest[F, T]{operation: o, callee: c})
	if err != nil {
		return Definition[F, T]{}, s.fail(err)
	}
	if err := s.check(); err != nil {
		return Definition[F, T]{}, err
	}
	if !o.completed || !out.Valid() || out.state.ticket != o.ticket || out.state.callee.state != c.state {
		return Definition[F, T]{}, s.fail(ErrInvalid)
	}
	c.state.phase = defined
	s.definitions = append(s.definitions, out)
	return out, nil
}

type defineRequest[F, T any] struct {
	operation *operation[F, T]
	callee    Callee[F, T]
}

func (r *defineRequest[F, T]) Site() Site[F, T]     { return Site[F, T]{r.operation} }
func (r *defineRequest[F, T]) Callee() Callee[F, T] { return r.callee }
func (r *defineRequest[F, T]) Open() OpenData[F]    { return r.callee.Open() }
func (r *defineRequest[F, T]) Bound(b Binder[F]) (Product[F, T], error) {
	if err := r.operation.check(); err != nil {
		return Product[F, T]{}, err
	}
	p, exists := r.callee.state.formals[b]
	if !exists {
		return Product[F, T]{}, r.operation.run.fail(ErrScope)
	}
	return p, nil
}
func (r *defineRequest[F, T]) Use(u Use[F]) (Product[F, T], error) {
	if err := r.operation.check(); err != nil {
		return Product[F, T]{}, err
	}
	if !r.callee.state.open.Declaration.Valid() {
		return Product[F, T]{}, r.operation.run.fail(ErrScope)
	}
	return r.operation.run.use(r.operation.environment, r.operation.frame, u)
}
func (r *defineRequest[F, T]) CompleteBody(body Fragment[T]) (Definition[F, T], error) {
	o, s := r.operation, r.operation.run
	if err := o.check(); err != nil {
		return Definition[F, T]{}, err
	}
	if !body.Valid() {
		return Definition[F, T]{}, s.fail(ErrInvalid)
	}
	if err := s.native.CheckBody(s.ctx, r.callee, body); err != nil {
		return Definition[F, T]{}, s.fail(err)
	}
	if err := o.check(); err != nil {
		return Definition[F, T]{}, err
	}
	o.completed = true
	return Definition[F, T]{&definitionState[F, T]{callee: r.callee, body: body, ticket: o.ticket}}, nil
}
