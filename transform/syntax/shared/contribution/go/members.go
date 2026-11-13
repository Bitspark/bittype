package contribution

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
)

// Members is the projector's part of a record contribution: the target's
// record from the recorded member names and the views of the dependencies. A
// dependency with a fragment is a delegated field, and its members are the
// fragment. Every other dependency is a field that the record lays out itself,
// named by the next recorded member. field and record are the target model's
// forms. It reads no source.
func Members[T any](p g.Projection[T], members []string, field func(name string, ty g.Expression[T]) g.Expression[T], record func(fields ...g.Expression[T]) g.Expression[T]) (g.Projected[T], error) {
	var laid []g.Expression[T]
	var table []g.Member
	own := 0
	for i := 0; i < p.Dependencies(); i++ {
		d := p.Dependency(i)
		if len(d.Fragment) > 0 {
			laid = append(laid, d.Fragment...)
			table = append(table, g.FragmentMembers(i))
			continue
		}
		if own >= len(members) {
			return g.Projected[T]{}, g.ErrInvalid
		}
		laid = append(laid, g.At(field(members[own], d.Type), d.Origin))
		table = append(table, g.Member{Name: members[own]})
		own++
	}
	if own != len(members) {
		return g.Projected[T]{}, g.ErrInvalid
	}
	return g.Projected[T]{Type: g.At(record(laid...), p.Origin()), Members: table}, nil
}

// Single is the target type of the one dependency of a contribution that
// names exactly one, such as the element of a list. It reads no source.
func Single[T any](p g.Projection[T]) (g.Expression[T], error) {
	if p.Dependencies() != 1 {
		return g.Expression[T]{}, g.ErrScope
	}
	return p.Dependency(0).Type, nil
}

// Constant completes a contribution that represents every selection of its
// domain by one constant target expression: its rule names no request and
// decides the given payload, and its projector returns the expression. The
// caller states the reference, the target, the profile, the domain and the
// schema, and the payload is a value of that schema's type. Nothing varies
// between two selections, so the expression is an invariant that the
// representation revision fixes (row 5).
func Constant[F, T, R any](c g.Contribution[F, T, R], payload any, expression g.Expression[T]) g.Contribution[F, T, R] {
	c.Rule = g.Rule[F, R]{Decide: func(context.Context, g.DecideStep[F, R]) (g.Decision, error) {
		return g.Decision{Payload: payload}, nil
	}}
	c.Project = func(context.Context, g.Projection[T]) (g.Projected[T], error) {
		return g.Projected[T]{Type: expression}, nil
	}
	return c
}
