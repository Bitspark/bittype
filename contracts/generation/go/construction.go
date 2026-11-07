package generation

import (
	"context"
	"sync"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// SourceSession resolves requested sites against the exact admitted inventory.
type SourceSession[F any] interface {
	Source[F]
	Catalog[F]
	ResolveRoots(context.Context, []ast.NodeRef) ([]Declaration[F], error)
}

// Construction providers own their state and meaning, never pipeline assembly.
// S and N are explicit provider-owned APIs. They are not kernel type algebras.
// bittype:dimension Syntax
type SourceProvider[F any, S SourceSession[F]] interface {
	Contract() transform.Contract
	Open(context.Context, transform.Request[F]) (S, error)
}

// bittype:dimension Syntax
type TargetProvider[F, T any, N Native[F, T]] interface {
	Open(context.Context, TargetRequest) (N, error)
}

// bittype:dimension Syntax
type MappingProvider[F, T any, S SourceSession[F], N Native[F, T]] interface {
	Open(context.Context, MappingRequest[S, N]) (Contribution[F, T], error)
}

type Contribution[F, T any] struct {
	Bindings   []RepresentationBinding[F, T]
	Structural *Mapping[F, T]
	Helpers    []Plan[F, T]
}
type TargetInfo struct {
	From, To representation.Language
	Output   ast.TreeID
	Limits   Limits
}

// Construction guards single-provider dispatch. The kernel owns selection and
// sequencing. A fresh construction has separate authority; it cannot attach
// unrelated source/native values to handles opened by another construction.
// Like Session, Construction must not be copied after use.
type Construction struct {
	mu           sync.Mutex
	ctx          context.Context
	info         TargetInfo
	closed, busy bool
	next         constructionStage
	calls        int
}

func NewConstruction(ctx context.Context, info TargetInfo) (*Construction, error) {
	if ctx == nil || info.From.Validate() != nil || info.To.Validate() != nil || info.Output == "" || !info.Limits.Valid() {
		return nil, ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return &Construction{ctx: ctx, info: info}, nil
}
func (c *Construction) Close() {
	if c != nil {
		c.mu.Lock()
		c.closed = true
		c.mu.Unlock()
	}
}

type constructionStage uint8

const (
	constructSource constructionStage = iota
	constructTarget
	constructMapping
)

func (c *Construction) begin(role constructionStage) (*constructionTicket, error) {
	if c == nil {
		return nil, ErrScope
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ctx == nil || c.closed || c.busy {
		return nil, ErrScope
	}
	if err := c.ctx.Err(); err != nil {
		c.closed = true
		return nil, err
	}
	if c.calls >= c.info.Limits.Operations {
		c.closed = true
		return nil, ErrLimit
	}
	if role != c.next {
		c.closed = true
		return nil, ErrScope
	}
	c.calls++
	c.busy = true
	if role < constructMapping {
		c.next++
	}
	return &constructionTicket{owner: c, active: true}, nil
}

type constructionTicket struct {
	owner  *Construction
	active bool
}

func (t *constructionTicket) check() error {
	if t == nil || t.owner == nil {
		return ErrScope
	}
	c := t.owner
	c.mu.Lock()
	defer c.mu.Unlock()
	if !t.active || c.closed {
		return ErrScope
	}
	return c.ctx.Err()
}

// The request is borrowed only during the call into provider code. Revoke the
// borrow before interpreting the result or publishing an opened handle.
func callProvider[V any](t *constructionTicket, open func() (V, error)) (V, error) {
	defer func() { t.owner.mu.Lock(); t.active = false; t.owner.mu.Unlock() }()
	return open()
}
func (t *constructionTicket) finish(ok bool) {
	c := t.owner
	c.mu.Lock()
	defer c.mu.Unlock()
	t.active = false
	c.busy = false
	if !ok {
		c.closed = true
	}
}

type TargetRequest struct{ ticket *constructionTicket }

func (r TargetRequest) Info() (TargetInfo, error) {
	if err := r.ticket.check(); err != nil {
		return TargetInfo{}, err
	}
	return r.ticket.owner.info, nil
}

type MappingRequest[S, N any] struct {
	ticket *constructionTicket
	source S
	target N
}

func (r MappingRequest[S, N]) Opened() (S, N, error) {
	if err := r.ticket.check(); err != nil {
		var s S
		var n N
		return s, n, err
	}
	return r.source, r.target, nil
}
func (r MappingRequest[S, N]) Info() (TargetInfo, error) { return (TargetRequest{r.ticket}).Info() }

// Opened handles have no public value constructor. Only successful provider
// dispatch can create them, and cross-construction mixing is rejected.
type OpenedSource[F any, S SourceSession[F]] struct {
	owner *Construction
	value S
}
type OpenedTarget[F, T any, N Native[F, T]] struct {
	owner *Construction
	value N
}

func (s OpenedSource[F, S]) Value(c *Construction) (S, error) {
	if err := opened(c, s.owner); err != nil {
		var zero S
		return zero, err
	}
	return s.value, nil
}
func (n OpenedTarget[F, T, N]) Value(c *Construction) (N, error) {
	if err := opened(c, n.owner); err != nil {
		var zero N
		return zero, err
	}
	return n.value, nil
}
func opened(c, owner *Construction) error {
	if c == nil || owner != c {
		return ErrScope
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return ErrScope
	}
	return c.ctx.Err()
}
func OpenSource[F any, S SourceSession[F]](c *Construction, p SourceProvider[F, S], r transform.Request[F]) (out OpenedSource[F, S], err error) {
	t, err := c.begin(constructSource)
	if err != nil {
		return out, err
	}
	ok := false
	defer func() { t.finish(ok) }()
	if representation.IsNil(p) || !r.Snapshot().Language().Compatible(c.info.From) {
		return out, ErrScope
	}
	value, err := p.Open(c.ctx, r)
	if err != nil {
		return out, err
	}
	if representation.IsNil(value) {
		return out, ErrInvalid
	}
	var evidence []semantic.Result
	for i := 0; ; i++ {
		e, exists := r.Evidence(i)
		if !exists {
			break
		}
		evidence = append(evidence, e)
	}
	if err = value.CheckInput(c.ctx, r.Snapshot(), evidence); err != nil {
		return out, err
	}
	if err = t.check(); err != nil {
		return out, err
	}
	ok = true
	return OpenedSource[F, S]{c, value}, nil
}
func OpenTarget[F, T any, N Native[F, T]](c *Construction, p TargetProvider[F, T, N]) (out OpenedTarget[F, T, N], err error) {
	t, err := c.begin(constructTarget)
	if err != nil {
		return out, err
	}
	ok := false
	defer func() { t.finish(ok) }()
	if representation.IsNil(p) || representation.CheckMarker[F](c.info.From) != nil || representation.CheckMarker[T](c.info.To) != nil {
		return out, ErrScope
	}
	value, err := callProvider(t, func() (N, error) { return p.Open(c.ctx, TargetRequest{t}) })
	if err != nil {
		return out, err
	}
	if representation.IsNil(value) {
		return out, ErrInvalid
	}
	if err = opened(c, c); err != nil {
		return out, err
	}
	ok = true
	return OpenedTarget[F, T, N]{c, value}, nil
}
func OpenMapping[F, T any, S SourceSession[F], N Native[F, T]](c *Construction, p MappingProvider[F, T, S, N], s OpenedSource[F, S], n OpenedTarget[F, T, N]) (out Contribution[F, T], err error) {
	if err = opened(c, s.owner); err != nil {
		c.Close()
		return out, err
	}
	if err = opened(c, n.owner); err != nil {
		c.Close()
		return out, err
	}
	t, err := c.begin(constructMapping)
	if err != nil {
		return out, err
	}
	ok := false
	defer func() { t.finish(ok) }()
	if representation.IsNil(p) {
		return out, ErrInvalid
	}
	out, err = callProvider(t, func() (Contribution[F, T], error) { return p.Open(c.ctx, MappingRequest[S, N]{t, s.value, n.value}) })
	if err != nil {
		return Contribution[F, T]{}, err
	}
	if err = opened(c, c); err != nil {
		return Contribution[F, T]{}, err
	}
	ok = true
	return out, nil
}
