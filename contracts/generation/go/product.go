package generation

import "github.com/Bitspark/bittype/trees/representation/go"

func (s *Session[F, T]) accept(o *operation[F, T], want Want[F, T], token Token[F, T]) (Product[F, T], error) {
	if err := o.check(); err != nil {
		return Product[F, T]{}, err
	}
	if !token.Valid() || !want.Valid() || want.site.operation != o {
		return Product[F, T]{}, s.fail(ErrInvalid)
	}
	if representation.IsNil(o.plan.value.Associations) {
		return Product[F, T]{}, s.fail(ErrUnsupported)
	}
	if err := s.step(); err != nil {
		return Product[F, T]{}, err
	}
	captured, err := o.plan.value.Associations.Capture(s.ctx, want, token)
	if err != nil {
		return Product[F, T]{}, s.fail(err)
	}
	if err := o.check(); err != nil {
		return Product[F, T]{}, err
	}
	if captured.Token != token || captured.Binding != want.Binding() || !captured.Payload.Valid() {
		return Product[F, T]{}, s.fail(ErrScope)
	}
	if err := s.native.CheckProduct(s.ctx, want, captured.Payload); err != nil {
		return Product[F, T]{}, s.fail(err)
	}
	if err := o.check(); err != nil {
		return Product[F, T]{}, err
	}
	return Product[F, T]{&productState[F, T]{captured: captured, ticket: o.ticket}}, nil
}
func (s *Session[F, T]) returned(o *operation[F, T], want Want[F, T], token Token[F, T]) (Product[F, T], error) {
	p, err := s.accept(o, want, token)
	if err != nil {
		return Product[F, T]{}, err
	}
	o.completed = true
	return p, nil
}

// forward only changes the callback return envelope. It preserves the accepted
// issuer, signature, payload and support. No public retag/ForwardAs route exists.
func (s *Session[F, T]) forward(o *operation[F, T], u UseData[F], p Product[F, T]) (Product[F, T], error) {
	if err := o.check(); err != nil {
		return Product[F, T]{}, err
	}
	b := p.Binding()
	if !p.Valid() || b.Run != s.id || b.Config != s.configuration || b.Frame != o.frame || b.Subject != u.Subject || b.Environment != u.Environment {
		return Product[F, T]{}, s.fail(ErrScope)
	}
	o.completed = true
	return Product[F, T]{&productState[F, T]{captured: p.state.captured, ticket: o.ticket}}, nil
}

type useRequest[F, T any] struct {
	operation *operation[F, T]
	data      UseData[F]
}

func (r *useRequest[F, T]) Want() Want[F, T] {
	return Want[F, T]{Site[F, T]{r.operation}, r.data.Clone()}
}
func (r *useRequest[F, T]) Child(u Use[F]) (Product[F, T], error) {
	if err := r.operation.check(); err != nil {
		return Product[F, T]{}, err
	}
	return r.operation.run.use(r.operation.environment, r.operation.frame, u)
}
func (r *useRequest[F, T]) Bound() (Product[F, T], error) {
	o, s := r.operation, r.operation.run
	if err := o.check(); err != nil {
		return Product[F, T]{}, err
	}
	if !r.data.Bound.Valid() {
		return Product[F, T]{}, s.fail(ErrScope)
	}
	p, ok := s.frames[o.frame][r.data.Bound]
	if !ok {
		return Product[F, T]{}, s.fail(ErrScope)
	}
	return s.forward(o, r.data, p)
}
func (r *useRequest[F, T]) Return(t Token[F, T]) (Product[F, T], error) {
	return r.operation.run.returned(r.operation, r.Want(), t)
}
func (r *useRequest[F, T]) RequireOwned(d Declaration[F]) (Callee[F, T], error) {
	if err := r.operation.check(); err != nil {
		return Callee[F, T]{}, err
	}
	c, err := r.operation.run.Acquire(d)
	if err != nil {
		return Callee[F, T]{}, err
	}
	if !c.Owned() {
		return Callee[F, T]{}, r.operation.run.fail(ErrInvalid)
	}
	return c, nil
}
func (r *useRequest[F, T]) RequireHelper(key string) (Callee[F, T], error) {
	if err := r.operation.check(); err != nil {
		return Callee[F, T]{}, err
	}
	return r.operation.run.AcquireHelper(key)
}
func (r *useRequest[F, T]) Apply() (Product[F, T], error) {
	o, s := r.operation, r.operation.run
	if err := o.check(); err != nil {
		return Product[F, T]{}, err
	}
	if !r.data.Callee.Valid() {
		return Product[F, T]{}, s.fail(ErrInvalid)
	}
	c, err := s.Acquire(r.data.Callee)
	if err != nil {
		return Product[F, T]{}, err
	}
	if c.state.plan != o.plan || representation.IsNil(c.state.plan.value.Convention) {
		return Product[F, T]{}, s.fail(ErrUnsupported)
	}
	a, err := s.enter(c.state.plan, o.frame, o.environment, c.state)
	if err != nil {
		return Product[F, T]{}, err
	}
	defer a.close()
	out, err := c.state.plan.value.Convention.Apply(s.ctx, &applyRequest[F, T]{operation: a, data: r.data.Clone(), callee: c})
	if err != nil {
		return Product[F, T]{}, s.fail(err)
	}
	if err := s.check(); err != nil {
		return Product[F, T]{}, err
	}
	if !a.completed || !out.Valid() || out.state.ticket != a.ticket {
		return Product[F, T]{}, s.fail(ErrInvalid)
	}
	return s.forward(o, r.data, out)
}

type applyRequest[F, T any] struct {
	operation *operation[F, T]
	data      UseData[F]
	callee    Callee[F, T]
}

func (r *applyRequest[F, T]) Want() Want[F, T] {
	return Want[F, T]{Site[F, T]{r.operation}, r.data.Clone()}
}
func (r *applyRequest[F, T]) Callee() Callee[F, T] { return r.callee }
func (r *applyRequest[F, T]) Actual(b Binder[F]) (Use[F], error) {
	if err := r.operation.check(); err != nil {
		return Use[F]{}, err
	}
	for _, a := range r.data.Actuals {
		if a.Binder == b {
			return a.Use, nil
		}
	}
	return Use[F]{}, r.operation.run.fail(ErrScope)
}
func (r *applyRequest[F, T]) Represent(b Binder[F]) (Product[F, T], error) {
	u, err := r.Actual(b)
	if err != nil {
		return Product[F, T]{}, err
	}
	return r.operation.run.use(r.operation.environment, r.operation.frame, u)
}
func (r *applyRequest[F, T]) Return(t Token[F, T]) (Product[F, T], error) {
	return r.operation.run.returned(r.operation, r.Want(), t)
}
