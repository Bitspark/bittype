package generation

import (
	"context"
	"fmt"
	"slices"

	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// PlannedDeclaration is one native declaration of a closed program. Role and
// Provenance are the planner's own facts. Provenance lists the declaring sites
// of the source declarations this declaration represents. It is a list because
// the contract does not assume that one source declaration yields one native
// declaration. A helper represents none.
type PlannedDeclaration[T any] struct {
	ID          uint32
	Description NativeDeclaration[T]
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
// been resolved. Run does not retain a source, reader or mapping expansion.
type GenerationPlan struct {
	steps []Step
	run   func(context.Context) (Output, error)
}

func (p GenerationPlan) Steps() []Step { return slices.Clone(p.steps) }
func (p GenerationPlan) Valid() bool   { return p.run != nil }
func (p GenerationPlan) Run(ctx context.Context) (Output, error) {
	if p.run == nil || ctx == nil {
		return Output{}, ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return Output{}, err
	}
	return p.run(ctx)
}

type SelectedDescription[F, T, R any] struct {
	ID          string
	Description MappingDescription[F, T, R]
}

// PlanningInput is one generation request: the source language and input, the
// source integration that interprets it, the requested root declaration sites
// and the selected mapping descriptions. Plan checks the descriptions, admits
// the input against the source's contract and opens the source over exactly
// that capture. Handles resolve only against that capture; interpreting it is
// the provider's responsibility.
type PlanningInput[F, T, R any] struct {
	From      representation.Language
	Input     transform.Input
	Source    SourceProvider[F, R]
	Roots     []ast.NodeRef
	Mappings  []SelectedDescription[F, T, R]
	Execution Execution[T]
}

// frame is the declaration being planned: its entity, planned ID and the native
// meaning of each of its real binders.
type frame[T any] struct {
	declaration semantics.EntityID
	scope       uint32
	env         map[semantics.EntityID]Expression[T]
}

type planner[F, T, R any] struct {
	ctx          context.Context
	language     representation.Language // the capture's source language
	walk         *walk[F]
	source       Interpretation[F]
	reader       R
	execution    Execution[T]
	bindings     map[string][]selectedBinding[F, T, R]
	recipes      map[string][]selectedRecipe[F, T, R]
	helpers      map[HelperKey]HelperRecipe[T]
	declarations map[semantics.EntityID]uint32
	helperIDs    map[HelperKey]uint32
	formals      map[uint32][]Binder[F]
	conventions  map[uint32]GenericConvention[T]
	active       map[uint32]bool
	program      Program[T]
	steps        []Step
	ops          int
	bytes        int
}
type selectedBinding[F, T, R any] struct {
	name    string
	binding BindingDescription[F, T, R]
}
type selectedRecipe[F, T, R any] struct {
	name   string
	recipe DeclarationRecipe[F, T, R]
}

func Plan[F, T, R any](ctx context.Context, in PlanningInput[F, T, R]) (GenerationPlan, error) {
	fail := func(e error) (GenerationPlan, error) { return GenerationPlan{}, e }
	if ctx == nil || representation.IsNil(in.Source) || representation.CheckMarker[F](in.From) != nil || !in.Execution.Limits.Valid() || representation.IsNil(in.Execution.Target) || representation.CheckMarker[T](in.Execution.Language) != nil || in.Execution.Output == "" || in.Execution.Checker.Validate() != nil || !in.Execution.Language.Compatible(in.Execution.Checker.Language()) || len(in.Roots) == 0 {
		return fail(ErrInvalid)
	}
	if err := ctx.Err(); err != nil {
		return fail(err)
	}
	p := &planner[F, T, R]{ctx: ctx, language: in.From, execution: in.Execution, bindings: map[string][]selectedBinding[F, T, R]{}, recipes: map[string][]selectedRecipe[F, T, R]{}, helpers: map[HelperKey]HelperRecipe[T]{}, declarations: map[semantics.EntityID]uint32{}, helperIDs: map[HelperKey]uint32{}, formals: map[uint32][]Binder[F]{}, conventions: map[uint32]GenericConvention[T]{}, active: map[uint32]bool{}}
	selected := map[string]bool{}
	entries := 0
	if len(in.Roots) > in.Execution.Limits.References || len(in.Mappings) > in.Execution.Limits.Plans {
		return fail(ErrLimit)
	}
	for _, part := range in.Mappings {
		if part.ID == "" || selected[part.ID] {
			return fail(ErrInvalid)
		}
		selected[part.ID] = true
		d, err := FreezeDescription(in.From, part.Description, in.Execution.Limits)
		if err != nil {
			return fail(err)
		}
		part.Description = d
		entries += len(d.Bindings) + len(d.Declarations) + len(d.Helpers)
		if entries > in.Execution.Limits.Plans {
			return fail(ErrLimit)
		}
		for _, b := range part.Description.Bindings {
			p.bindings[b.Input.Selector.Key()] = append(p.bindings[b.Input.Selector.Key()], selectedBinding[F, T, R]{part.ID + "/" + b.ID, b})
		}
		for _, r := range part.Description.Declarations {
			p.recipes[r.Selector.Key()] = append(p.recipes[r.Selector.Key()], selectedRecipe[F, T, R]{part.ID + "/" + r.ID, r})
		}
		for _, h := range part.Description.Helpers {
			if _, ok := p.helpers[h.Key]; ok {
				return fail(ErrAmbiguous)
			}
			p.helpers[h.Key] = h
		}
	}
	// Mapping descriptions are checked before the evidence, as the registry
	// checks them before it calls Plan. Then the kernel admits the evidence and
	// opens the source over exactly that capture; no caller pairs an opened
	// source with a capture. Admission requires the input in the From language.
	captured, err := admitInput[F](in.From, in.Source.Contract(), in.Input)
	if err != nil {
		return fail(err)
	}
	p.walk = newWalk(captured)
	source, err := in.Source.Open(ctx, captured)
	if err != nil {
		return fail(err)
	}
	if representation.IsNil(source.Interpretation) {
		return fail(ErrInvalid)
	}
	p.source, p.reader = source.Interpretation, source.Reader
	seen := map[ast.NodeRef]bool{}
	for _, site := range in.Roots {
		if seen[site] {
			return fail(ErrScope)
		}
		seen[site] = true
		root, err := p.root(site)
		if err != nil {
			return fail(err)
		}
		id, err := p.declare(root, 0)
		if err != nil {
			return fail(err)
		}
		if p.program.declarations[id-1].Role == RequestedDeclaration {
			return fail(ErrScope)
		}
		p.program.exposed = append(p.program.exposed, id)
		p.program.declarations[id-1].Role = RequestedDeclaration
	}
	if err := in.Execution.Target.CheckProgram(ctx, p.program); err != nil {
		return fail(err)
	}
	if err := ctx.Err(); err != nil {
		return fail(err)
	}
	// Capture only the closed target program and stable target configuration.
	program := p.program
	exec := in.Execution
	// What the plan states about its declarations, for the inventory of every run.
	owned := make([]planned, len(program.declarations))
	for i, d := range program.declarations {
		owned[i] = planned{role: d.Role, provenance: d.Provenance, values: d.Description.Values}
	}
	return GenerationPlan{steps: slices.Clone(p.steps), run: func(ctx context.Context) (Output, error) {
		built, err := exec.Target.Build(ctx, BuildInput[T]{program, exec.Language, exec.Output, exec.Limits})
		if err != nil {
			return Output{}, err
		}
		if err = ctx.Err(); err != nil {
			return Output{}, err
		}
		out, err := seal[T](ctx, exec.Checker, built, owned, program.exposed, exec.Limits)
		if err != nil {
			return Output{}, err
		}
		if out.Snapshot().Entry().Tree != exec.Output {
			return Output{}, ErrScope
		}
		return out, nil
	}}, nil
}
func (p *planner[F, T, R]) tick(depth int) error {
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
func (p *planner[F, T, R]) root(site ast.NodeRef) (Declaration[F], error) {
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
func (p *planner[F, T, R]) selector(d Declaration[F], term semantics.TermID, at semantics.EntityID) (string, error) {
	key, ok, err := p.source.Declaration(p.ctx, d)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", p.walk.unsupported(term, at, "the source does not generate this declaration")
	}
	return p.qualified(key, term, at, "declaration selector")
}

// qualified is the key of a selector the source returned. Mapping selectors
// are frozen in the capture's source language, so requiring the same language
// here makes every rule lookup exact: a selector of another language or profile
// is outside the run and chooses neither an exact nor a default rule.
func (p *planner[F, T, R]) qualified(s Selector, term semantics.TermID, entity semantics.EntityID, what string) (string, error) {
	if !s.Valid() {
		return "", ErrScope
	}
	if !s.Language().Compatible(p.language) {
		return "", p.walk.located(ErrScope, term, entity, what+" outside the source language")
	}
	return s.Key(), nil
}

// admit establishes an occurrence and its whole subtree before any mapping
// represents it. Every occurrence in the subtree needs a structural reading in
// its declaration's scope and an Admitted source judgment, including actuals a
// native convention omits and operands no mapping captures. Representation stays
// demand-driven; admission is not. Each occurrence is admitted once per run.
func (p *planner[F, T, R]) admit(o *occurrence[F], depth int) error {
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

func (p *planner[F, T, R]) reserve(role DeclarationRole, provenance []ast.NodeRef) (uint32, error) {
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
func (p *planner[F, T, R]) carry(values []IssuedValue) error {
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
func (p *planner[F, T, R]) declare(d Declaration[F], depth int) (uint32, error) {
	if err := p.tick(depth); err != nil {
		return 0, err
	}
	if id := p.declarations[d.entity]; id != 0 {
		if p.active[id] {
			return 0, fmt.Errorf("%w: recursive native declaration", ErrUnsupported)
		}
		return id, nil
	}
	key, err := p.selector(d, 0, d.entity)
	if err != nil {
		return 0, err
	}
	choices := p.recipes[key]
	if len(choices) == 0 {
		// An exact use mapping cannot silently become a structural publication.
		if len(p.bindings[key]) > 0 {
			return 0, fmt.Errorf("%w: exact use has no declaration recipe", ErrUnsupported)
		}
		choices = p.recipes[""]
	}
	if len(choices) == 0 {
		return 0, ErrUnsupported
	}
	if len(choices) != 1 {
		return 0, ErrAmbiguous
	}
	choice := choices[0]
	view, err := p.walk.declaration(d)
	if err != nil {
		return 0, err
	}
	if len(view.Formals) > p.execution.Limits.Formals {
		return 0, ErrLimit
	}
	j, err := p.source.JudgeDeclaration(p.ctx, d)
	if err == nil {
		err = verdict(j)
	}
	// A generated declaration's body is admitted even if its recipe does not
	// represent it, and before the recipe reads it.
	if err == nil {
		err = p.admit(view.Body.o, depth+1)
	}
	if err != nil {
		return 0, err
	}
	// A source declaration is a dependency until the request lists it as a root.
	id, err := p.reserve(DependencyDeclaration, p.walk.declared(d.entity))
	if err != nil {
		return 0, err
	}
	p.declarations[d.entity] = id
	p.formals[id] = view.Formals
	state := &expansionState[F, T, R]{active: true, ctx: p.ctx, reader: p.reader, declaration: &view, captures: map[Slot]Use[F]{}, limits: p.execution.Limits}
	desc, err := choice.recipe.Expand(p.ctx, DeclarationRequest[F, T, R]{Expansion[F, T, R]{state}})
	state.mu.Lock()
	state.active = false
	failure := state.failure
	state.mu.Unlock()
	if err != nil {
		return 0, err
	}
	if failure != nil {
		return 0, failure
	}
	desc = desc.Clone()
	if len(desc.Generic.Bindings) != len(view.Formals) || len(desc.Generic.Arguments) != len(desc.Generic.Parameters) {
		return 0, ErrScope
	}
	if err = p.carry(desc.Values); err != nil {
		return 0, err
	}
	if len(desc.Generic.Parameters) > p.execution.Limits.Formals {
		return 0, ErrLimit
	}
	f := frame[T]{declaration: d.entity, scope: id, env: map[semantics.EntityID]Expression[T]{}}
	for i, binder := range view.Formals {
		if _, dup := f.env[binder.entity]; dup {
			return 0, ErrScope
		}
		value, err := p.template(desc.Generic.Bindings[i], id, len(desc.Generic.Parameters), nil, depth+1)
		if err != nil {
			return 0, err
		}
		f.env[binder.entity] = value
	}
	// Validate all application expressions now, before a use is encountered.
	for _, arg := range desc.Generic.Arguments {
		if err = p.convention(arg, len(view.Formals), depth+1); err != nil {
			return 0, err
		}
	}
	p.program.declarations[id-1].Description = nativeDeclaration(desc)
	p.conventions[id] = desc.Generic
	desc.Body, err = p.resolve(desc.Body, state.captures, f, depth+1)
	if err != nil {
		return 0, err
	}
	if err = checkParameterUses(id, desc); err != nil {
		return 0, err
	}
	if err = p.execution.Target.CheckDeclaration(p.ctx, nativeDeclaration(desc)); err != nil {
		return 0, err
	}
	p.program.declarations[id-1].Description = nativeDeclaration(desc)
	p.active[id] = false
	p.steps = append(p.steps, Step{Binding: choice.name, Source: key, Declaration: desc.Name, Form: desc.Form.Name(), Expression: describe(desc.Body)})
	return id, nil
}

// use plans one occurrence: the kernel admits it with its whole subtree, then a
// binding or the callee's declaration represents it.
func (p *planner[F, T, R]) use(u Use[F], f frame[T], depth int) (Expression[T], error) {
	if err := p.tick(depth); err != nil {
		return Expression[T]{}, err
	}
	if !p.walk.owns(u) || u.o.declaration != f.declaration {
		return Expression[T]{}, ErrScope
	}
	if err := p.admit(u.o, depth); err != nil {
		return Expression[T]{}, err
	}
	view, err := p.walk.structure(u.o)
	if err != nil {
		return Expression[T]{}, err
	}
	var key string
	switch view.Kind {
	case BoundUse:
		value, ok := f.env[view.Bound.entity]
		if !ok {
			return Expression[T]{}, ErrScope
		}
		return value, nil
	case DeclarationUse:
		if key, err = p.selector(view.Callee, u.o.term, u.o.declaration); err != nil {
			return Expression[T]{}, err
		}
	case ExpressionUse:
		selector, err := p.source.Expression(p.ctx, u)
		if err != nil {
			return Expression[T]{}, err
		}
		if key, err = p.qualified(selector, u.o.term, u.o.declaration, "expression selector"); err != nil {
			return Expression[T]{}, err
		}
	default:
		return Expression[T]{}, ErrInvalid
	}
	choices := p.bindings[key]
	if len(choices) > 1 {
		return Expression[T]{}, fmt.Errorf("%w: %s", ErrAmbiguous, key)
	}
	if len(choices) == 1 {
		return p.bind(choices[0], key, view, f, depth)
	}
	if view.Kind != DeclarationUse {
		return Expression[T]{}, p.walk.unsupported(u.o.term, u.o.declaration, "no mapping for "+key)
	}
	id, err := p.declare(view.Callee, depth+1)
	if err != nil {
		return Expression[T]{}, err
	}
	convention := p.conventions[id]
	formals := p.formals[id]
	if len(view.Actuals) != len(formals) {
		return Expression[T]{}, ErrScope
	}
	needed := map[int]bool{}
	var neededArguments func(Expression[T])
	neededArguments = func(e Expression[T]) {
		if e.kind == ArgumentExpression {
			needed[e.index] = true
		}
		for _, a := range e.args {
			neededArguments(a)
		}
	}
	for _, arg := range convention.Arguments {
		neededArguments(arg)
	}
	actuals := make([]Expression[T], len(formals))
	for i, binder := range formals {
		a, ok := view.Actual(binder)
		if !ok {
			return Expression[T]{}, ErrScope
		}
		if needed[i] {
			if actuals[i], err = p.use(a, f, depth+1); err != nil {
				return Expression[T]{}, err
			}
		}
	}
	args := make([]Expression[T], len(convention.Arguments))
	for i, arg := range convention.Arguments {
		if args[i], err = p.template(arg, 0, 0, actuals, depth+1); err != nil {
			return Expression[T]{}, err
		}
	}
	return Expression[T]{kind: DeclarationExpression, declaration: id, args: args}, nil
}

// bind applies an exact binding: its captures must be children of the view.
func (p *planner[F, T, R]) bind(b selectedBinding[F, T, R], key string, view UseView[F], f frame[T], depth int) (Expression[T], error) {
	captures := map[Slot]Use[F]{}
	if b.binding.Input.Capture != nil {
		values, err := b.binding.Input.Capture(p.ctx, p.reader, view.Clone())
		if err != nil {
			return Expression[T]{}, err
		}
		if len(values) != len(b.binding.Input.Slots) {
			return Expression[T]{}, ErrScope
		}
		allowed := map[Slot]bool{}
		for _, slot := range b.binding.Input.Slots {
			allowed[slot] = true
		}
		for _, v := range values {
			if !allowed[v.Slot] || captures[v.Slot].Valid() {
				return Expression[T]{}, ErrScope
			}
			if !childOf(v.Use, view.Use) {
				return Expression[T]{}, errNotChild
			}
			captures[v.Slot] = v.Use
		}
	}
	expression := b.binding.Output
	if b.binding.Expand != nil {
		v := view.Clone()
		state := &expansionState[F, T, R]{active: true, ctx: p.ctx, reader: p.reader, use: &v, captures: captures, limits: p.execution.Limits}
		var err error
		expression, err = b.binding.Expand(p.ctx, Expansion[F, T, R]{state})
		state.mu.Lock()
		state.active = false
		failure := state.failure
		state.mu.Unlock()
		if err != nil {
			return Expression[T]{}, err
		}
		if failure != nil {
			return Expression[T]{}, failure
		}
	}
	out, err := p.resolve(expression, captures, f, depth+1)
	if err != nil {
		return Expression[T]{}, err
	}
	p.steps = append(p.steps, Step{Binding: b.name, Source: key, Form: out.form.Name(), Expression: describe(out)})
	return out, nil
}
func (p *planner[F, T, R]) resolve(e Expression[T], captures map[Slot]Use[F], f frame[T], depth int) (Expression[T], error) {
	if err := p.tick(depth); err != nil {
		return Expression[T]{}, err
	}
	size := len(e.form.name)
	for _, value := range e.data {
		size += len(value)
	}
	if size > p.execution.Limits.OutputBytes-p.bytes {
		return Expression[T]{}, ErrLimit
	}
	p.bytes += size
	switch e.kind {
	case CapturedExpression:
		use := captures[e.slot]
		if !use.Valid() {
			return Expression[T]{}, ErrScope
		}
		return p.use(use, f, depth+1)
	case NativeExpression, HelperExpression:
		args := make([]Expression[T], len(e.args))
		for i, a := range e.args {
			v, err := p.resolve(a, captures, f, depth+1)
			if err != nil {
				return Expression[T]{}, err
			}
			args[i] = v
		}
		e.args = args
		if e.kind == HelperExpression {
			return p.helperUse(e, depth+1)
		}
		if err := p.execution.Target.CheckExpression(p.ctx, e); err != nil {
			return Expression[T]{}, err
		}
		return e, nil
	case ParameterExpression:
		if f.scope == 0 || e.index < 0 || e.index >= len(p.program.declarations[f.scope-1].Description.Parameters) {
			return Expression[T]{}, ErrScope
		}
		e.declaration = f.scope
		return e, nil
	default:
		return Expression[T]{}, ErrScope
	}
}
func (p *planner[F, T, R]) template(e Expression[T], scope uint32, count int, actuals []Expression[T], depth int) (Expression[T], error) {
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
	case NativeExpression, HelperExpression:
		args := make([]Expression[T], len(e.args))
		for i, a := range e.args {
			v, err := p.template(a, scope, count, actuals, depth+1)
			if err != nil {
				return Expression[T]{}, err
			}
			args[i] = v
		}
		e.args = args
		if e.kind == HelperExpression {
			return p.helperUse(e, depth+1)
		}
		if err := p.execution.Target.CheckExpression(p.ctx, e); err != nil {
			return Expression[T]{}, err
		}
		return e, nil
	default:
		return Expression[T]{}, ErrScope
	}
}
func (p *planner[F, T, R]) convention(e Expression[T], count, depth int) error {
	if err := p.tick(depth); err != nil {
		return err
	}
	if e.kind == ArgumentExpression {
		if e.index < 0 || e.index >= count {
			return ErrScope
		}
		return nil
	}
	if e.kind != NativeExpression && e.kind != HelperExpression {
		return ErrScope
	}
	for _, a := range e.args {
		if err := p.convention(a, count, depth+1); err != nil {
			return err
		}
	}
	if e.kind == HelperExpression {
		helper, ok := p.helpers[e.helper]
		if !ok {
			return ErrUnsupported
		}
		if len(e.args) != len(helper.Declaration.Parameters) {
			return ErrScope
		}
		return nil
	}
	return p.execution.Target.CheckExpression(p.ctx, e)
}
func (p *planner[F, T, R]) helperUse(e Expression[T], depth int) (Expression[T], error) {
	recipe, ok := p.helpers[e.helper]
	if !ok {
		return Expression[T]{}, ErrUnsupported
	}
	id := p.helperIDs[e.helper]
	if id != 0 && p.active[id] {
		return Expression[T]{}, ErrUnsupported
	}
	if id == 0 {
		var err error
		id, err = p.reserve(HelperDeclaration, nil)
		if err != nil {
			return Expression[T]{}, err
		}
		p.helperIDs[e.helper] = id
		desc := recipe.Declaration.Clone()
		if len(desc.Parameters) > p.execution.Limits.Formals {
			return Expression[T]{}, ErrLimit
		}
		if err = p.carry(desc.Values); err != nil {
			return Expression[T]{}, err
		}
		desc.Body, err = p.template(desc.Body, id, len(desc.Parameters), nil, depth+1)
		if err != nil {
			return Expression[T]{}, err
		}
		// The target checks a copy. The plan keeps its own declaration.
		if err = p.execution.Target.CheckDeclaration(p.ctx, desc.Clone()); err != nil {
			return Expression[T]{}, err
		}
		p.program.declarations[id-1].Description = desc
		p.active[id] = false
		p.steps = append(p.steps, Step{Binding: "helper/" + e.helper.Name(), Declaration: desc.Name, Form: desc.Form.Name(), Expression: describe(desc.Body)})
	}
	if len(e.args) != len(p.program.declarations[id-1].Description.Parameters) {
		return Expression[T]{}, ErrScope
	}
	return Expression[T]{kind: DeclarationExpression, declaration: id, args: e.args}, nil
}
func describe[T any](e Expression[T]) string {
	switch e.kind {
	case ParameterExpression:
		return fmt.Sprintf("parameter(%d,%d)", e.declaration, e.index)
	case DeclarationExpression:
		return fmt.Sprintf("declaration(%d)%v", e.declaration, describeArgs(e.args))
	case NativeExpression:
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

// nativeDeclaration deliberately drops source binding/application templates.
// Construction can only observe the native declaration selected by planning.
// Each result owns its parameters and its values, so the declaration that a
// target checks shares no storage with the one the plan keeps.
func nativeDeclaration[T any](d DeclarationDescription[T]) NativeDeclaration[T] {
	out := NativeDeclaration[T]{Name: d.Name, Form: d.Form, Origin: d.Origin, Body: d.Body, Values: slices.Clone(d.Values)}
	for _, p := range d.Generic.Parameters {
		out.Parameters = append(out.Parameters, NativeParameter{Name: p.Name, Origin: p.Origin})
	}
	return out
}
