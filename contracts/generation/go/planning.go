package generation

import (
	"context"
	"fmt"
	"slices"

	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// PlannedDeclaration is one target declaration of a closed program. Role and
// Provenance are the planner's own facts. Provenance lists the declaring sites
// of the source declarations this declaration represents. It is a list because
// the contract does not assume that one source declaration yields one target
// declaration. A helper represents none.
type PlannedDeclaration[T any] struct {
	ID          uint32
	Description TargetDeclaration[T]
	Role        DeclarationRole
	Provenance  []ast.NodeRef
}
type Program[T any] struct {
	declarations []PlannedDeclaration[T]
	exposed      []uint32
}

func (p Program[T]) Declarations() []PlannedDeclaration[T] {
	out := slices.Clone(p.declarations)
	for i := range out {
		out[i].Description = out[i].Description.Clone()
		out[i].Provenance = slices.Clone(out[i].Provenance)
	}
	return out
}
func (p Program[T]) Exposed() []uint32 { return slices.Clone(p.exposed) }

type Step struct {
	Binding, Source, Declaration, Form string
	Expression                         string
}

// GenerationPlan is issued only after all reachable rules and dependencies have
// been resolved. Run does not retain a source, reader or rule. The plan
// carries the selected representation plan of the same planning: the target
// program and every entry of that plan are projections of the same frozen
// decision records.
type GenerationPlan struct {
	steps          []Step
	run            func(context.Context) (Output, error)
	representation RepresentationPlan
}

func (p GenerationPlan) Steps() []Step { return slices.Clone(p.steps) }

// Representation is the selected representation plan.
func (p GenerationPlan) Representation() RepresentationPlan { return p.representation }
func (p GenerationPlan) Valid() bool                        { return p.run != nil }
func (p GenerationPlan) Run(ctx context.Context) (Output, error) {
	if p.run == nil || ctx == nil {
		return Output{}, ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return Output{}, err
	}
	return p.run(ctx)
}

// traversal is what the planner needs of one run besides selection: the walk
// over the captured graph, the source's interpretation and reader, the target
// execution, the program under construction and the budgets.
type traversal[F, T, R any] struct {
	ctx       context.Context
	language  representation.Language // the capture's source language
	walk      *walk[F]
	source    Interpretation[F]
	reader    R
	execution Execution[T]
	active    map[uint32]bool
	program   Program[T]
	steps     []Step
	ops       int
	bytes     int
}

func (p *traversal[F, T, R]) tick(depth int) error {
	if err := p.ctx.Err(); err != nil {
		return err
	}
	p.ops++
	if depth > p.execution.Limits.Depth || p.ops > p.execution.Limits.Operations {
		return ErrLimit
	}
	return nil
}

// root resolves a requested syntax site to the one entity that declares it and
// requires the source to classify that entity as a generatable declaration.
func (p *traversal[F, T, R]) root(site ast.NodeRef) (Declaration[F], error) {
	entities := p.walk.sites[site]
	if len(entities) != 1 {
		return Declaration[F]{}, ErrScope
	}
	d := Declaration[F]{p.walk, entities[0]}
	if _, ok, err := p.source.Declaration(p.ctx, d); err != nil {
		return Declaration[F]{}, err
	} else if !ok {
		return Declaration[F]{}, ErrScope
	}
	return d, nil
}

// selector is the source's stable identity for a declaration the walk reached.
// term and at locate a failure: the use that reaches the declaration and the
// declaration that use occurs in, or for a generated declaration itself, no
// term and that declaration.
func (p *traversal[F, T, R]) selector(d Declaration[F], term semantics.TermID, at semantics.EntityID) (string, error) {
	key, ok, err := p.source.Declaration(p.ctx, d)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", p.walk.unsupported(term, at, "the source does not generate this declaration")
	}
	return p.qualified(key, term, at, "declaration selector")
}

// qualified is the key of a selector the source returned. A domain's selectors
// are in the capture's source language, so requiring the same language here
// makes every selection exact: a selector of another language or profile is
// outside the run and selects neither an exact contribution nor the policy.
func (p *traversal[F, T, R]) qualified(s Selector, term semantics.TermID, entity semantics.EntityID, what string) (string, error) {
	if !s.Valid() {
		return "", ErrScope
	}
	if !s.Language().Compatible(p.language) {
		return "", p.walk.located(ErrScope, term, entity, what+" outside the source language")
	}
	return s.Key(), nil
}

// admit establishes an occurrence and its whole subtree before any contribution
// represents it. Every occurrence in the subtree needs a structural reading in
// its declaration's scope and an Admitted source judgment, including actuals a
// target convention omits and operands no rule requests. Representation stays
// demand-driven; admission is not. Each occurrence is admitted once per run.
func (p *traversal[F, T, R]) admit(o *occurrence[F], depth int) error {
	if o.admitted {
		return o.admission
	}
	err := p.tick(depth)
	var view UseView[F]
	if err == nil {
		view, err = p.walk.structure(o)
	}
	if err == nil {
		var j Judgment
		if j, err = p.source.JudgeUse(p.ctx, Use[F]{o}); err == nil {
			err = verdict(j)
		}
	}
	for _, a := range view.Actuals {
		if err == nil {
			err = p.admit(a.Use.o, depth+1)
		}
	}
	for _, operand := range view.Operands {
		if err == nil {
			err = p.admit(operand.Use.o, depth+1)
		}
	}
	o.admitted, o.admission = true, err
	return err
}

func (p *traversal[F, T, R]) reserve(role DeclarationRole, provenance []ast.NodeRef) (uint32, error) {
	if len(p.program.declarations) >= p.execution.Limits.Reservations {
		return 0, ErrLimit
	}
	id := uint32(len(p.program.declarations) + 1)
	p.program.declarations = append(p.program.declarations, PlannedDeclaration[T]{ID: id, Role: role, Provenance: provenance})
	p.active[id] = true
	return id, nil
}

// carry admits the values a description states about its declaration. Each
// must come from an issuer. Its content is charged to the output budget and is
// not read.
func (p *traversal[F, T, R]) carry(values []IssuedValue) error {
	for _, value := range values {
		if value.owner == nil {
			return ErrInvalid
		}
		if len(value.content)+1 > p.execution.Limits.OutputBytes-p.bytes {
			return ErrLimit
		}
		p.bytes += len(value.content) + 1
	}
	return nil
}

// template resolves an expression of a generic convention: a binding of a
// source binder, in the scope of its declaration, or an argument of an
// application, with the actuals of one use.
func (p *traversal[F, T, R]) template(e Expression[T], scope uint32, count int, actuals []Expression[T], depth int) (Expression[T], error) {
	if err := p.tick(depth); err != nil {
		return Expression[T]{}, err
	}
	switch e.kind {
	case ParameterExpression:
		if scope == 0 || e.index < 0 || e.index >= count {
			return Expression[T]{}, ErrScope
		}
		e.declaration = scope
		return e, nil
	case ArgumentExpression:
		if e.index < 0 || e.index >= len(actuals) {
			return Expression[T]{}, ErrScope
		}
		return actuals[e.index], nil
	case FormExpression:
		args := make([]Expression[T], len(e.args))
		for i, a := range e.args {
			v, err := p.template(a, scope, count, actuals, depth+1)
			if err != nil {
				return Expression[T]{}, err
			}
			args[i] = v
		}
		e.args = args
		if err := p.execution.Target.CheckExpression(p.ctx, e); err != nil {
			return Expression[T]{}, err
		}
		return e, nil
	default:
		return Expression[T]{}, ErrScope
	}
}
func (p *traversal[F, T, R]) convention(e Expression[T], count, depth int) error {
	if err := p.tick(depth); err != nil {
		return err
	}
	if e.kind == ArgumentExpression {
		if e.index < 0 || e.index >= count {
			return ErrScope
		}
		return nil
	}
	if e.kind != FormExpression {
		return ErrScope
	}
	for _, a := range e.args {
		if err := p.convention(a, count, depth+1); err != nil {
			return err
		}
	}
	return p.execution.Target.CheckExpression(p.ctx, e)
}
func describe[T any](e Expression[T]) string {
	switch e.kind {
	case ParameterExpression:
		return fmt.Sprintf("parameter(%d,%d)", e.declaration, e.index)
	case DeclarationExpression:
		return fmt.Sprintf("declaration(%d)%v", e.declaration, describeArgs(e.args))
	case FormExpression:
		return e.form.Name() + fmt.Sprint(e.data) + fmt.Sprint(describeArgs(e.args))
	default:
		return "unresolved"
	}
}
func describeArgs[T any](args []Expression[T]) []string {
	out := make([]string, len(args))
	for i, e := range args {
		out[i] = describe(e)
	}
	return out
}

func checkParameterUses[T any](id uint32, d DeclarationDescription[T]) error {
	used := map[int]bool{}
	var visit func(Expression[T])
	visit = func(e Expression[T]) {
		if e.kind == ParameterExpression && e.declaration == id {
			used[e.index] = true
		}
		for _, a := range e.args {
			visit(a)
		}
	}
	visit(d.Body)
	for i, p := range d.Generic.Parameters {
		if p.RequiredUse && !used[i] {
			return ErrUnsupported
		}
	}
	return nil
}

// targetDeclaration deliberately drops source binding/application templates.
// Construction can only observe the target declaration selected by planning.
// Each result owns its parameters and its values, so the declaration that a
// target checks shares no storage with the one the plan keeps.
func targetDeclaration[T any](d DeclarationDescription[T]) TargetDeclaration[T] {
	out := TargetDeclaration[T]{Name: d.Name, Form: d.Form, Origin: d.Origin, Body: d.Body, Values: slices.Clone(d.Values)}
	for _, p := range d.Generic.Parameters {
		out.Parameters = append(out.Parameters, TargetParameter{Name: p.Name, Origin: p.Origin})
	}
	return out
}
