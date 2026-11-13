package generation

import (
	"slices"

	"github.com/Bitspark/bittype/trees/representation/go"
)

// FreezeDescription captures mutable container fields and checks the composition
// grammar. Owner callbacks must capture stable configuration; they are trusted
// implementations, and are never retained by a finished plan.
func FreezeDescription[F, T, R any](from representation.Language, d MappingDescription[F, T, R], limits Limits) (MappingDescription[F, T, R], error) {
	fail := func(err error) (MappingDescription[F, T, R], error) { return MappingDescription[F, T, R]{}, err }
	if d.Profile == "" || !limits.Valid() {
		return fail(ErrInvalid)
	}
	if len(d.Bindings)+len(d.Declarations)+len(d.Helpers) > limits.Plans {
		return fail(ErrLimit)
	}
	d.Bindings = slices.Clone(d.Bindings)
	d.Declarations = slices.Clone(d.Declarations)
	d.Helpers = slices.Clone(d.Helpers)
	ids := map[string]bool{}
	slots := map[Slot]bool{}
	nodes, bytes := 0, 0
	var check func(Expression[T], map[Slot]bool, int) error
	check = func(e Expression[T], allowed map[Slot]bool, depth int) error {
		nodes++
		bytes += len(e.form.name)
		for _, s := range e.data {
			bytes += len(s)
		}
		if nodes > limits.OutputNodes || bytes > limits.OutputBytes || depth > limits.Depth {
			return ErrLimit
		}
		switch e.kind {
		case FormExpression:
			if !e.form.Valid() {
				return ErrScope
			}
		case CapturedExpression:
			if !allowed[e.slot] {
				return ErrScope
			}
		case HelperExpression:
			if !e.helper.Valid() {
				return ErrScope
			}
		case ArgumentExpression, ParameterExpression:
			if e.index < 0 {
				return ErrScope
			}
		default:
			return ErrInvalid
		}
		for _, a := range e.args {
			if err := check(a, allowed, depth+1); err != nil {
				return err
			}
		}
		return nil
	}
	for i, b := range d.Bindings {
		if b.ID == "" || ids[b.ID] || !b.Input.Selector.Valid() || !from.Compatible(b.Input.Selector.Language()) {
			return fail(ErrInvalid)
		}
		ids[b.ID] = true
		if len(b.Input.Slots) > limits.References {
			return fail(ErrLimit)
		}
		b.Input.Slots = slices.Clone(b.Input.Slots)
		allowed := map[Slot]bool{}
		for _, slot := range b.Input.Slots {
			if !slot.Valid() || slots[slot] {
				return fail(ErrScope)
			}
			slots[slot] = true
			allowed[slot] = true
		}
		if len(allowed) > 0 && b.Input.Capture == nil {
			return fail(ErrInvalid)
		}
		if b.Expand != nil {
			if b.Output.kind != 0 {
				return fail(ErrInvalid)
			}
		} else {
			if err := check(b.Output, allowed, 0); err != nil {
				return fail(err)
			}
		}
		d.Bindings[i] = b
	}
	for _, r := range d.Declarations {
		if r.ID == "" || ids[r.ID] || r.Expand == nil || r.Selector.Valid() && !from.Compatible(r.Selector.Language()) {
			return fail(ErrInvalid)
		}
		ids[r.ID] = true
	}
	for i, h := range d.Helpers {
		if !h.Key.Valid() {
			return fail(ErrInvalid)
		}
		h.Declaration = h.Declaration.Clone()
		if err := check(h.Declaration.Body, nil, 0); err != nil {
			return fail(err)
		}
		d.Helpers[i] = h
	}
	return d, nil
}
