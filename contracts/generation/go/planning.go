package generation

import (
	"context"
	"fmt"
	"slices"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type PlannedDeclaration[T any] struct {
	ID          uint32
	Description NativeDeclaration[T]
	Exposed     bool
}
type Program[T any] struct {
	declarations []PlannedDeclaration[T]
	exposed      []uint32
}

func (p Program[T]) Declarations() []PlannedDeclaration[T] {
	out := slices.Clone(p.declarations)
	for i := range out {
		out[i].Description = out[i].Description.Clone()
	}
	return out
}
func (p Program[T]) Exposed() []uint32 { return slices.Clone(p.exposed) }

type Step struct {
	Binding, Source, Declaration, Form string
	Expression                         string
}

// GenerationPlan is issued only after all reachable rules and dependencies have
// been resolved. Run does not retain a source reader or mapping expansion.
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

type SelectedDescription[F, T any, S SourceSession[F]] struct {
	ID          string
	Description MappingDescription[F, T, S]
}
type PlanningInput[F, T any, S SourceSession[F]] struct {
	Source    S
	Snapshot  semantic.Snapshot
	Evidence  []semantic.Result
	Roots     []Declaration[F]
	Mappings  []SelectedDescription[F, T, S]
	Execution Execution[T]
}
type planner[F, T any, S SourceSession[F]] struct {
	ctx          context.Context
	source       S
	execution    Execution[T]
	bindings     map[string][]selectedBinding[F, T, S]
	recipes      map[string][]selectedRecipe[F, T, S]
	helpers      map[HelperKey]HelperRecipe[T]
	declarations map[Declaration[F]]uint32
	helperIDs    map[HelperKey]uint32
	opens        map[Declaration[F]]OpenData[F]
	conventions  map[uint32]GenericConvention[T]
	active       map[uint32]bool
	program      Program[T]
	steps        []Step
	ops          int
	bytes        int
}
type selectedBinding[F, T any, S SourceSession[F]] struct {
	name    string
	binding BindingDescription[F, T, S]
}
type selectedRecipe[F, T any, S SourceSession[F]] struct {
	name   string
	recipe DeclarationRecipe[F, T, S]
}

func Plan[F, T any, S SourceSession[F]](ctx context.Context, in PlanningInput[F, T, S]) (GenerationPlan, error) {
	fail := func(e error) (GenerationPlan, error) { return GenerationPlan{}, e }
	if ctx == nil || representation.IsNil(in.Source) || !in.Source.Owners().Valid() || !in.Snapshot.Valid() || representation.CheckMarker[F](in.Snapshot.Language()) != nil || !in.Execution.Limits.Valid() || representation.IsNil(in.Execution.Target) || representation.CheckMarker[T](in.Execution.Language) != nil || in.Execution.Output == "" || in.Execution.Checker.Validate() != nil || !in.Execution.Language.Compatible(in.Execution.Checker.Language()) || len(in.Roots) == 0 {
		return fail(ErrInvalid)
	}
	if err := ctx.Err(); err != nil {
		return fail(err)
	}
	if err := in.Source.CheckInput(ctx, in.Snapshot, in.Evidence); err != nil {
		return fail(err)
	}
	p := &planner[F, T, S]{ctx: ctx, source: in.Source, execution: in.Execution, bindings: map[string][]selectedBinding[F, T, S]{}, recipes: map[string][]selectedRecipe[F, T, S]{}, helpers: map[HelperKey]HelperRecipe[T]{}, declarations: map[Declaration[F]]uint32{}, helperIDs: map[HelperKey]uint32{}, opens: map[Declaration[F]]OpenData[F]{}, conventions: map[uint32]GenericConvention[T]{}, active: map[uint32]bool{}}
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
		d, err := FreezeDescription(in.Snapshot.Language(), part.Description, in.Execution.Limits)
		if err != nil {
			return fail(err)
		}
		part.Description = d
		entries += len(d.Bindings) + len(d.Declarations) + len(d.Helpers)
		if entries > in.Execution.Limits.Plans {
			return fail(ErrLimit)
		}
		for _, b := range part.Description.Bindings {
			p.bindings[b.Input.Selector.Key()] = append(p.bindings[b.Input.Selector.Key()], selectedBinding[F, T, S]{part.ID + "/" + b.ID, b})
		}
		for _, r := range part.Description.Declarations {
			p.recipes[r.Selector.Key()] = append(p.recipes[r.Selector.Key()], selectedRecipe[F, T, S]{part.ID + "/" + r.ID, r})
		}
		for _, h := range part.Description.Helpers {
			if _, ok := p.helpers[h.Key]; ok {
				return fail(ErrAmbiguous)
			}
			p.helpers[h.Key] = h
		}
	}
	seen := map[Declaration[F]]bool{}
	for _, root := range in.Roots {
		if seen[root] {
			return fail(ErrScope)
		}
		seen[root] = true
		id, err := p.declare(root, 0)
		if err != nil {
			return fail(err)
		}
		p.program.exposed = append(p.program.exposed, id)
		p.program.declarations[id-1].Exposed = true
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
	return GenerationPlan{steps: slices.Clone(p.steps), run: func(ctx context.Context) (Output, error) {
		built, err := exec.Target.Build(ctx, BuildInput[T]{program, exec.Language, exec.Output, exec.Limits})
		if err != nil {
			return Output{}, err
		}
		if err = ctx.Err(); err != nil {
			return Output{}, err
		}
		out, err := seal(ctx, exec.Checker, built, len(program.exposed), exec.Limits)
		if err != nil {
			return Output{}, err
		}
		if out.Snapshot().Entry().Tree != exec.Output {
			return Output{}, ErrScope
		}
		return out, nil
	}}, nil
}
func (p *planner[F, T, S]) tick(depth int) error {
	if err := p.ctx.Err(); err != nil {
		return err
	}
	p.ops++
	if depth > p.execution.Limits.Depth || p.ops > p.execution.Limits.Operations {
		return ErrLimit
	}
	return nil
}
func (p *planner[F, T, S]) declarationKey(d Declaration[F]) (string, error) {
	if err := p.source.Authenticate(p.ctx, References[F]{Declarations: []Declaration[F]{d}}); err != nil {
		return "", err
	}
	key, ok, err := p.source.DeclarationSelector(d)
	if err != nil {
		return "", err
	}
	if !ok || !key.Valid() {
		return "", ErrScope
	}
	return key.Key(), nil
}
func (p *planner[F, T, S]) reserve() (uint32, error) {
	if len(p.program.declarations) >= p.execution.Limits.Reservations {
		return 0, ErrLimit
	}
	id := uint32(len(p.program.declarations) + 1)
	p.program.declarations = append(p.program.declarations, PlannedDeclaration[T]{ID: id})
	p.active[id] = true
	return id, nil
}
func (p *planner[F, T, S]) declare(d Declaration[F], depth int) (uint32, error) {
	if err := p.tick(depth); err != nil {
		return 0, err
	}
	if id := p.declarations[d]; id != 0 {
		if p.active[id] {
			return 0, fmt.Errorf("%w: recursive native declaration", ErrUnsupported)
		}
		return id, nil
	}
	key, err := p.declarationKey(d)
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
	open, err := p.source.Open(p.ctx, d)
	if err != nil {
		return 0, err
	}
	open = open.Clone()
	if open.Declaration != d || !open.Body.Valid() || !open.Environment.Valid() {
		return 0, ErrScope
	}
	if len(open.Formals) > p.execution.Limits.Formals {
		return 0, ErrLimit
	}
	if err = p.source.Authenticate(p.ctx, References[F]{Uses: []Use[F]{open.Body}, Environments: []Environment[F]{open.Environment}}); err != nil {
		return 0, err
	}
	id, err := p.reserve()
	if err != nil {
		return 0, err
	}
	p.declarations[d] = id
	p.opens[d] = open
	state := &expansionState[F, T, S]{active: true, ctx: p.ctx, source: p.source, open: open, captures: map[Slot]Use[F]{}, limits: p.execution.Limits}
	desc, err := choice.recipe.Expand(p.ctx, DeclarationRequest[F, T, S]{Expansion[F, T, S]{state}})
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
	if len(desc.Generic.Bindings) != len(open.Formals) || len(desc.Generic.Arguments) != len(desc.Generic.Parameters) {
		return 0, ErrScope
	}
	if len(desc.Generic.Parameters) > p.execution.Limits.Formals {
		return 0, ErrLimit
	}
	env := map[Binder[F]]Expression[T]{}
	for i, f := range open.Formals {
		if !f.Binder.Valid() || env[f.Binder].kind != 0 {
			return 0, ErrScope
		}
		if err = p.source.Authenticate(p.ctx, References[F]{Binders: []Binder[F]{f.Binder}, Uses: []Use[F]{f.Subject}}); err != nil {
			return 0, err
		}
		value, err := p.template(desc.Generic.Bindings[i], id, len(desc.Generic.Parameters), nil, depth+1)
		if err != nil {
			return 0, err
		}
		env[f.Binder] = value
	}
	// Validate all application expressions now, before a use is encountered.
	for _, arg := range desc.Generic.Arguments {
		if err = p.convention(arg, len(open.Formals), depth+1); err != nil {
			return 0, err
		}
	}
	p.program.declarations[id-1].Description = nativeDeclaration(desc)
	p.conventions[id] = desc.Generic
	desc.Body, err = p.resolve(desc.Body, state.captures, open.Environment, env, id, depth+1)
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
func (p *planner[F, T, S]) use(use Use[F], environment Environment[F], env map[Binder[F]]Expression[T], scope uint32, depth int) (Expression[T], error) {
	if err := p.tick(depth); err != nil {
		return Expression[T]{}, err
	}
	data, err := p.source.Admit(p.ctx, environment, use)
	if err != nil {
		return Expression[T]{}, err
	}
	if data.Use != use || data.Environment != environment {
		return Expression[T]{}, ErrScope
	}
	if data.Bound.Valid() {
		value, ok := env[data.Bound]
		if !ok {
			return Expression[T]{}, ErrScope
		}
		return value, nil
	}
	// Whole-source admission applies even to arguments a representation omits.
	for _, a := range data.Actuals {
		if err = p.source.Authenticate(p.ctx, References[F]{Binders: []Binder[F]{a.Binder}, Uses: []Use[F]{a.Use}}); err != nil {
			return Expression[T]{}, err
		}
		if _, err = p.source.Admit(p.ctx, environment, a.Use); err != nil {
			return Expression[T]{}, err
		}
	}
	selector, err := p.source.UseSelector(data)
	if data.Callee.Valid() {
		var ok bool
		selector, ok, err = p.source.DeclarationSelector(data.Callee)
		if err == nil && !ok {
			err = ErrScope
		}
	}
	if err != nil {
		return Expression[T]{}, err
	}
	choices := p.bindings[selector.Key()]
	if len(choices) > 1 {
		return Expression[T]{}, fmt.Errorf("%w: %s", ErrAmbiguous, selector.Key())
	}
	if len(choices) == 1 {
		b := choices[0]
		captures := map[Slot]Use[F]{}
		if b.binding.Input.Capture != nil {
			values, err := b.binding.Input.Capture(p.ctx, p.source, data.Clone())
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
				if err = p.source.CheckChild(p.ctx, environment, use, v.Use); err != nil {
					return Expression[T]{}, err
				}
				captures[v.Slot] = v.Use
			}
		}
		expression := b.binding.Output
		if b.binding.Expand != nil {
			state := &expansionState[F, T, S]{active: true, ctx: p.ctx, source: p.source, use: data.Clone(), captures: captures, limits: p.execution.Limits}
			expression, err = b.binding.Expand(p.ctx, Expansion[F, T, S]{state})
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
		out, err := p.resolve(expression, captures, environment, env, scope, depth+1)
		if err != nil {
			return Expression[T]{}, err
		}
		p.steps = append(p.steps, Step{Binding: b.name, Source: selector.Key(), Form: out.form.Name(), Expression: describe(out)})
		return out, nil
	}
	if !data.Callee.Valid() {
		return Expression[T]{}, fmt.Errorf("%w: %s", ErrUnsupported, selector.Key())
	}
	id, err := p.declare(data.Callee, depth+1)
	if err != nil {
		return Expression[T]{}, err
	}
	convention := p.conventions[id]
	open := p.opens[data.Callee]
	if len(data.Actuals) != len(open.Formals) {
		return Expression[T]{}, ErrScope
	}
	actuals := make([]Expression[T], len(open.Formals))
	byBinder := map[Binder[F]]Use[F]{}
	for _, a := range data.Actuals {
		if byBinder[a.Binder].Valid() {
			return Expression[T]{}, ErrScope
		}
		byBinder[a.Binder] = a.Use
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
	for i, f := range open.Formals {
		a := byBinder[f.Binder]
		if !a.Valid() {
			return Expression[T]{}, ErrScope
		}
		if needed[i] {
			actuals[i], err = p.use(a, environment, env, scope, depth+1)
			if err != nil {
				return Expression[T]{}, err
			}
		}
	}

	args := make([]Expression[T], len(convention.Arguments))
	for i, arg := range convention.Arguments {
		args[i], err = p.template(arg, 0, 0, actuals, depth+1)
		if err != nil {
			return Expression[T]{}, err
		}
	}
	return Expression[T]{kind: DeclarationExpression, declaration: id, args: args}, nil
}
func (p *planner[F, T, S]) resolve(e Expression[T], captures map[Slot]Use[F], environment Environment[F], env map[Binder[F]]Expression[T], scope uint32, depth int) (Expression[T], error) {
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
		return p.use(use, environment, env, scope, depth+1)
	case NativeExpression, HelperExpression:
		args := make([]Expression[T], len(e.args))
		for i, a := range e.args {
			v, err := p.resolve(a, captures, environment, env, scope, depth+1)
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
		if scope == 0 || e.index < 0 || e.index >= len(p.program.declarations[scope-1].Description.Parameters) {
			return Expression[T]{}, ErrScope
		}
		e.declaration = scope
		return e, nil
	default:
		return Expression[T]{}, ErrScope
	}
}
func (p *planner[F, T, S]) template(e Expression[T], scope uint32, count int, actuals []Expression[T], depth int) (Expression[T], error) {
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
func (p *planner[F, T, S]) convention(e Expression[T], count, depth int) error {
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
func (p *planner[F, T, S]) helperUse(e Expression[T], depth int) (Expression[T], error) {
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
		id, err = p.reserve()
		if err != nil {
			return Expression[T]{}, err
		}
		p.helperIDs[e.helper] = id
		desc := recipe.Declaration.Clone()
		if len(desc.Parameters) > p.execution.Limits.Formals {
			return Expression[T]{}, ErrLimit
		}
		desc.Body, err = p.template(desc.Body, id, len(desc.Parameters), nil, depth+1)
		if err != nil {
			return Expression[T]{}, err
		}
		if err = p.execution.Target.CheckDeclaration(p.ctx, desc); err != nil {
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
func nativeDeclaration[T any](d DeclarationDescription[T]) NativeDeclaration[T] {
	out := NativeDeclaration[T]{Name: d.Name, Form: d.Form, Origin: d.Origin, Body: d.Body}
	for _, p := range d.Generic.Parameters {
		out.Parameters = append(out.Parameters, NativeParameter{Name: p.Name, Origin: p.Origin})
	}
	return out
}
