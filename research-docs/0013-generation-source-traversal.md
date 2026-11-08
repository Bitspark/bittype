# Research: Kernel-owned semantic traversal and narrow source interpretation in type generation

**ID:** 0013
**Date:** 8 November 2026
**Status:** draft
**Owner:** claude-bittype-gaps
**Issue:** https://github.com/Bitspark/bittype/issues/158
**Code baseline:** bittype 57187bcb0621c8da3ac1ef12e12887d4d079cd94; horos Go core 0.10.2

## Question

A type generator turns declarations in a source language into type declarations in
native programming languages. Its generation kernel is language-agnostic. Today the
kernel asks each *source integration* (the plug-in for one source language) for two
different things through one interface:

- **Navigation** of an already established semantic graph: which declaration a use
  names, a declaration's parameters and body, and which argument of an application
  belongs to which parameter.
- **Interpretation** of source-specific meaning: stable declaration identity, the
  meaning of language-specific expression forms, and whether a use is admissible
  (well-kinded, within bounds) under the declaration's assumptions.

The navigation facts are already recorded in a shared, language-agnostic semantic
graph that every source integration consumes. So every new source integration has to
re-implement the same graph walk. The project's architecture rules treat this as a
defect: a shared mechanism should own the invariant mechanics, and a language owner
should contribute only its meanings.

**The decision we need help with:** how should the boundary between the generation
kernel and a source integration be drawn, so that the kernel follows declaration,
reference and application structure through the shared graph, while source-owned
contracts keep interpretation and applicability judgments? The answer must keep a set
of guarantees (listed under [Guarantees that must survive](#guarantees-that-must-survive)).
It must also say how mapping libraries, which today read the source integration's
records directly, see the kernel's structural view.

## Context

### System overview

**bittype** is a Go library. It parses, checks and prints source code (syntax), and
it generates native type declarations for TypeScript, Go, Rust and C++. Its input
language is **horos**, a small declaration language for data types. A horos module
declares named type definitions, which may be generic (have type parameters) and may
constrain them with *kinds*, which are classifications of types similar to type
classes or bounds. For example:

```horos
type Person = { name: h.Text; };
type Cell<V> = ...;            // owned by a separate library
type Page<T> = { items: List<T>; cursor: h.Text; };
type People = Page<Cell<Person>>;
```

Here `h` is an imported standard prelude module that defines kinds such as `h.Text`
and `h.Data`. A bound is written `type TextBox<T: h.Text> = { value: T; };`.

Several independently owned parts cooperate:

| Part | Owns |
| --- | --- |
| horos **analysis provider** | Name resolution, kind checking and the semantic facts of horos. It publishes them as one immutable semantic graph per captured input. |
| **Generation kernel** (`contracts/generation/go`) | Planning: selection of mappings, reservation of native declarations, substitution of native type arguments, scope checks, limits, and all-or-nothing output. It is meant to interpret no source or native rules. |
| **Source integration** (for horos: `transform/syntax/horos/source/go`) | Today: navigation of the graph *and* horos interpretation (the subject of this document). |
| **Mapping libraries** | One per source type family and target (for example horos `List` to TypeScript arrays, or a library that owns `Cell<V>` and maps it to an imported native `Cell<T>`). They describe native representations declaratively. |
| **Target** (one per native language) | Checks native expressions and declarations, and builds native syntax trees from a closed native program. |

Vocabulary used throughout:

- **Snapshot:** an immutable capture of input syntax trees. Capturing equal bytes
  twice yields two different snapshots; identity is by allocation, not content.
- **Result** or **evidence:** an immutable output of an analysis operation (resolution
  or analysis) bound to exactly one snapshot. It carries *reports* (named obligations,
  each with a completion status and a verdict), diagnostics and, optionally, a semantic
  graph.
- **Semantic graph** (`semantics.Semantics[L]`): the shared, language-agnostic store of
  established facts for one snapshot (described next). `L` is a compile-time marker
  type for the language.
- **Declaration**, **binder**, **use**, **environment:** a declaration is a named type
  definition. A binder is one of its real type parameters. A use is one type occurrence
  in a body. An environment is the set of binders and assumptions under which a use is
  interpreted; today it is the enclosing declaration.
- **Selector:** a stable, source-owned identity string qualified by language and
  profile, for example `declaration/<module pin>/Cell` or
  `builtin/6:scalar/9:utf8-text`. Mapping libraries register rules against selectors;
  the kernel compares selectors and never interprets them.
- **Native program:** the closed list of planned native declarations that a target
  turns into syntax trees. Planning resolves every choice before construction runs.
- **Profile:** a named configuration that selects which mapping libraries apply (for
  example `pagination-representation/1`).

### The shared semantic graph

The graph is the shared context that both the kernel and the source integration hold.
It stores entities (declarations, parameters, modules and so on), terms (the
expression trees of bodies), occurrence associations and provider-owned facts. Its
constructor validates only graph integrity, such as addresses, parameter ownership and
acyclic term containment. It does not check language rules. Excerpts from
`semantics/go/model.go`:

```go
// IDs are one-based addresses local to one Semantics[L], never persistent identities.
type EntityID uint32
type TermID uint32

// Key names vocabulary owned by a provider. It causes no registration or lookup.
type Key struct{ Owner, Name string }

type Entity struct {
	Kind         Key           // provider-owned, e.g. {horos, type-definition}
	Name         string
	Owner        EntityID      // enclosing entity (for a parameter: its declaration)
	Declarations []ast.NodeRef // syntax sites that declare it
	Parameters   []Parameter   // real binders, in declaration order
	Body         TermID
}

// Parameter retains a real binder and its coarse provider-owned sort.
type Parameter struct {
	Entity EntityID
	Sort   Key
}

type TermKind uint8

const (
	Reference TermKind = iota + 1
	Application
	ProviderExpression
)

// Term has exactly one variant. Application arguments are effective positional
// arguments, not a mixture of written, inferred, defaulted or packed arguments.
// Reference owns Target; Application owns Callee and Arguments. Applicability,
// arity and saturation belong to the provider, not graph admission.
// A provider expression's encoding defines its bytes and operand locators.
type Term struct {
	Kind      TermKind
	Target    EntityID   // Reference
	Callee    TermID     // Application
	Arguments []TermID   // Application
	Operator  Key        // ProviderExpression, e.g. {horos, descriptor/list-of}
	Encoding  Key        // ProviderExpression payload encoding
	Data      []byte     // ProviderExpression payload
	Operands  []Operand  // ProviderExpression children
}

type Operand struct {
	Locator []uint32 // provider-defined path inside the payload
	Term    TermID
}

// Subject selects exactly one model-local entity or term.
type Subject struct {
	Entity EntityID
	Term   TermID
}

// Fact associates a subject with a provider-owned fact term. A relation may have
// several values; its meaning and cardinality belong to the provider.
type Fact struct {
	Subject  Subject
	Relation Key
	Term     TermID
}
```

Graph admission enforces, among other things, that each parameter entity is listed by
exactly one entity and names that entity as its `Owner`, that a reference names an
existing entity, and that term *containment* (callee, arguments, operands) is acyclic.
Entity references may form cycles. Containment may form a DAG: one term may be shared.

The graph offers point queries: `Entity(id)`, `Term(id)`, `Facts(subject, relation)`,
occurrence lookups, and one structural pairing query:

```go
type ArgumentBinding struct {
	Parameter EntityID
	Argument  TermID
}

// DirectParameterBindings pairs the real parameters of a direct Reference
// callee with an equally sized argument list. Parameter-headed applications,
// other callees and arity mismatches have no such view. This structural query
// does not establish applicability.
// A direct zero-parameter application returns an empty successful view.
func (m Semantics[L]) DirectParameterBindings(id TermID) ([]ArgumentBinding, bool) {
	...
	term := m.state.Terms[id-1]
	if term.Kind != Application { return nil, false }
	callee := m.state.Terms[term.Callee-1]
	if callee.Kind != Reference { return nil, false }
	entity := m.state.Entities[callee.Target-1]
	if entity.Owner != 0 {               // callee is itself a parameter: no view
		for _, parameter := range m.state.Entities[entity.Owner-1].Parameters {
			if parameter.Entity == callee.Target { return nil, false }
		}
	}
	parameters := entity.Parameters
	if len(parameters) != len(term.Arguments) { return nil, false }
	result := make([]ArgumentBinding, len(parameters))
	for i, parameter := range parameters {
		result[i] = ArgumentBinding{Parameter: parameter.Entity, Argument: term.Arguments[i]}
	}
	return result, true
}
```

There are no reverse or indexed queries (for example "terms owned by a declaration" or
"entity declared at a syntax site"). A consumer builds those indexes itself.

### Evidence and capture identity

The graph reaches generation inside evidence. A source integration declares which
evidence it needs as an `EvidenceRequirement`:

```go
// EvidenceRequirement is a consumer-owned compatibility requirement, not a
// request to run a provider. Semantics additionally requires a published graph.
type EvidenceRequirement struct {
	Operation     semantics.Key   // e.g. {bittype, analysis}
	Profile       string          // provider profile
	Configuration []Setting
	Reports       []Requirement   // each: obligation key + required verdict
	Semantics     bool
}
```

`semantic.AdmitEvidence(snapshot, requirements, results)` matches requirement *i* to
result *i*. It rejects a result produced from a different snapshot (pointer identity),
a different operation, profile or configuration, a report that is not complete with
exactly the required verdict, and a missing graph when `Semantics` is true. A
**completion** is one of Complete, Incomplete or Unsupported. A **verdict** is one of
Holds, Rejected or Undetermined. "Complete/Undetermined" means a completed enumeration
that asserts no truth value. The horos generation source requires, among others,
`program-validity` Complete/Holds and `type-facts`, `applications`,
`occurrence-bindings` Complete/Undetermined.

A result returns its graph as an untyped envelope. A consumer obtains the typed graph
with `representation.UnpackSemantics[L](snapshot.Language(), envelope)`, which checks
the language identity, profile and marker type. A graph holds no pointer back to its
snapshot. The only binding between graph and snapshot is the result that published it.

### What horos records and what its analysis judges

The horos analysis provider emits five entity kinds into the graph:

| Entity kind | Owner | Parameters | Body |
| --- | --- | --- | --- |
| `{horos, module}` | none | none | none |
| `{horos, type-definition}` | its module | its type-parameter entities | the type expression |
| `{horos, kind-definition}` | its module | its kind-parameter entities | the kind expression |
| `{horos, type-parameter}` | its type-definition | none | none |
| `{horos, kind-parameter}` | its kind-definition | none | none |

Imported modules and the always-loaded prelude (which declares kinds such as `Text`
and `Data`, and no type-definitions) are ordinary module entities in the same graph.
An import creates no entity; a reference to an imported declaration simply targets
that declaration's entity. horos has no nested generic declarations: every
type-definition is owned by a module.

Type-definition bodies use the three term variants as follows:

| horos expression | Term |
| --- | --- |
| A parameter, e.g. `T` | `Reference` whose target is the type-parameter entity |
| A named type with no arguments, e.g. `Person` | `Reference` whose target is the type-definition (never a zero-argument application) |
| An application, e.g. `Page<Cell<Person>>` | `Application` whose callee is a term (normally a `Reference`) and whose arguments are terms |
| A built-in descriptor, e.g. `List<T>`, `{ name: h.Text; }`, a scalar | `ProviderExpression` with operator `{horos, descriptor/<constructor>}`, payload encoding `{ontos, codec-v1}`, and the nested type terms as operands |

Operand locators are paths into the encoded descriptor tuple, defined by horos. For
example, a list element is at `[1]`, and record field *i* has its type at `[i+1, 3]`.
Record field names and the optional flag are only in the payload.

Facts used by generation (all owner `horos`): `sort` (on every type term, declaration
and parameter; the value is a sort term `{horos, sort/proper}` or `{horos, sort/arrow}`,
an arrow being a type constructor such as `* -> *`), `result-sort` (on declarations),
`canonical-identity` (on modules, as a 32-byte content hash called the *module pin*;
also on declarations and parameters), `declared-kind` (a parameter's bound, as a kind
term) and `membership` (a proven kind membership, with its assumptions, on every type
term). There is no separate "bounds" relation: a bound is the `declared-kind` fact on
the parameter. A missing fact asserts nothing.

The analysis judges the whole module closure before it reports
`program-validity` Complete/Holds. Under each declaration's own parameter assumptions:

- every application exactly saturates its callee and has sort-correct arguments;
- every argument provably satisfies the callee's declared parameter kind;
- every descriptor operand satisfies its formation contract;
- the body satisfies the declared result kind.

When a judgment fails, no graph is published. A refuted bound yields
`program-validity` Complete/Rejected; an unprovable bound yields Incomplete/Undetermined.
The analysis does **not** require types to be *proper* (non-constructor). Generation adds
that as a profile restriction: constructor-valued types have no default native
representation. horos therefore accepts several forms that generation must reject, and
today it rejects them inconsistently:

| horos-valid form | Graph shape | Today's generation error |
| --- | --- | --- |
| Constructor-valued declaration or parameter | sort `sort/arrow` | Unsupported diagnostic (correct) |
| Curried application `Later<Alias><text>` | application whose callee is an application | Unsupported diagnostic (correct) |
| Parameter-headed application `F<A>` | application whose callee references a parameter | scope error (wrong kind of error) |
| Applying a zero-parameter constructor alias, `Alias<text>` | application with one argument whose callee has no parameters | invalid-request error (wrong kind of error) |

These misclassifications illustrate why "no structural view", "rejected judgment",
"uncovered fact" and "malformed input" need distinct outcomes.

### The current generation contract

The public contract (`contracts/generation/go/source.go` and `mapping.go`), complete:

```go
// Owner-qualified addresses. A Space allocates refs; only refs from the configured
// owner's space are accepted. A Ref is an address, not a proof of existence.
type Declaration[F any] = Ref[DeclarationRole[F]]
type Binder[F any]      = Ref[BinderRole[F]]
type Use[F any]         = Ref[UseRole[F]]
type Environment[F any] = Ref[EnvironmentRole[F]]

type Formal[F any] struct { Binder Binder[F]; Subject Use[F] }
type Actual[F any] struct { Binder Binder[F]; Use Use[F] }

// These are provider construction records. Only the configured source can admit them.
type OpenData[F any] struct {
	Declaration Declaration[F]
	Environment Environment[F]
	Body        Use[F]
	Formals     []Formal[F]
}

type UseData[F any] struct {
	Use         Use[F]          // the occurrence
	Subject     Use[F]          // the source-certified effective subject
	Environment Environment[F]
	Callee      Declaration[F]  // set for a use of a declaration
	Bound       Binder[F]       // set for a use of a bound formal
	Actuals     []Actual[F]     // callee binders paired with argument uses
}

type SourceOwners[F any] struct {
	Declarations Owner[DeclarationRole[F]]
	Binders      Owner[BinderRole[F]]
	Uses         Owner[UseRole[F]]
	Environments Owner[EnvironmentRole[F]]
}

// References requests authentication of existing, correctly typed source records.
type References[F any] struct {
	Declarations []Declaration[F]
	Binders      []Binder[F]
	Uses         []Use[F]
	Environments []Environment[F]
}

// Source consumes already supplied receipts/queries. It performs no hidden
// resolution or analysis. Its records and policy remain immutable during a run.
type Source[F any] interface {
	CheckInput(context.Context, semantic.Snapshot, []semantic.Result) error
	Owners() SourceOwners[F]
	RootEnvironment() Environment[F]
	Authenticate(context.Context, References[F]) error
	Open(context.Context, Declaration[F]) (OpenData[F], error)
	// Admit establishes the whole use's applicability under its original
	// assumptions using supplied evidence, including actuals the mapping may
	// never represent. Missing evidence must fail, never trigger hidden analysis.
	Admit(context.Context, Environment[F], Use[F]) (UseData[F], error)
	// CheckChild authenticates a direct semantic operand or generic actual of parent.
	CheckChild(context.Context, Environment[F], Use[F], Use[F]) error
}

// Catalog is source-owned interpretation; enumeration does not demand expansion.
type Catalog[F any] interface {
	Declarations() []Declaration[F]
	DeclarationSelector(Declaration[F]) (Selector, bool, error)
	UseSelector(UseData[F]) (Selector, error)
}

type SourceSession[F any] interface {
	Source[F]
	Catalog[F]
	ResolveRoots(context.Context, []ast.NodeRef) ([]Declaration[F], error)
}

type SourceProvider[F any, S SourceSession[F]] interface {
	Contract() transform.Contract   // the evidence requirements above
	Open(context.Context, transform.Request[F]) (S, error)
}
```

Mapping libraries are typed by the concrete session type `S`, so they can call the
source integration's own reading API in addition to the contract:

```go
// Pattern fixes selection and declares its captures. Capture interprets the
// selected source meaning only; the kernel checks each child relation.
type Pattern[F any, S SourceSession[F]] struct {
	Selector Selector
	Slots    []Slot
	Capture  func(context.Context, S, UseData[F]) ([]Capture[F], error)
}
type BindingDescription[F, T any, S SourceSession[F]] struct {
	ID     string
	Input  Pattern[F, S]
	Output Expression[T]                                             // fixed output, or
	Expand func(context.Context, Expansion[F, T, S]) (Expression[T], error) // computed output
}
type DeclarationRecipe[F, T any, S SourceSession[F]] struct {
	ID       string
	Selector Selector // zero selector: the structural declaration policy
	Expand   func(context.Context, DeclarationRequest[F, T, S]) (DeclarationDescription[T], error)
}
// GenericConvention: Bindings is indexed by source parameter (what each source
// binder means inside the native declaration); Arguments is indexed by native
// parameter (how a use fills each native parameter from the source actuals).
type GenericConvention[T any] struct {
	Parameters []ParameterDescription
	Bindings   []Expression[T]
	Arguments  []Expression[T]
}
type MappingLibrary[F, T any, S SourceSession[F]] interface {
	Describe(from, to representation.Language) (MappingDescription[F, T, S], error)
}
```

`Expression[T]` is a native expression tree with a few kernel-owned placeholders:
`Mapped(slot)` (substitute the mapped representation of a captured source use),
`Parameter(i)` (a native type parameter of the declaration being planned),
`Argument(i)` (the native representation of the *i*-th source actual at a use site),
and `Helper(key)` (a library-owned helper declaration).

During an expansion, a mapping library holds an `Expansion` value whose methods expire
when the callback returns:

```go
func (r Expansion[F, T, S]) Source() (S, error)        // the session, while active
func (r Expansion[F, T, S]) Use() (UseData[F], error)  // the admitted use being mapped
func (r Expansion[F, T, S]) Mapped(use Use[F]) (Expression[T], error) {
	...
	if s.open.Declaration.Valid() {
		if use != s.open.Body { err = ErrScope } // declaration recipes can map their body only
	} else {
		err = s.source.CheckChild(s.ctx, s.use.Environment, s.use.Use, use)
	}
	...
}
func (r DeclarationRequest[F, T, S]) Open() (OpenData[F], error)
```

### How the planner drives the source today

The planner (`contracts/generation/go/planning.go`) starts at requested root
declarations, declares each one, and walks uses depth-first. Condensed excerpts:

```go
func (p *planner[F, T, S]) declare(d Declaration[F], depth int) (uint32, error) {
	if id := p.declarations[d]; id != 0 {
		if p.active[id] { return 0, fmt.Errorf("%w: recursive native declaration", ErrUnsupported) }
		return id, nil
	}
	key, err := p.declarationKey(d) // Authenticate + DeclarationSelector
	choices := p.recipes[key]       // exact recipe, else the structural recipe
	...
	open, err := p.source.Open(p.ctx, d)
	if open.Declaration != d || !open.Body.Valid() || !open.Environment.Valid() { return 0, ErrScope }
	p.source.Authenticate(... open.Body, open.Environment ...)
	id, _ := p.reserve()
	desc, err := choice.recipe.Expand(ctx, DeclarationRequest{...open...})
	if len(desc.Generic.Bindings) != len(open.Formals) { return 0, ErrScope }
	env := map[Binder[F]]Expression[T]{}
	for i, f := range open.Formals {          // real binder -> its native meaning
		p.source.Authenticate(... f.Binder, f.Subject ...)
		env[f.Binder], _ = p.template(desc.Generic.Bindings[i], id, ...)
	}
	desc.Body, err = p.resolve(desc.Body, captures, open.Environment, env, id, depth+1)
	...
}

func (p *planner[F, T, S]) use(use Use[F], environment Environment[F], env map[Binder[F]]Expression[T], scope uint32, depth int) (Expression[T], error) {
	data, err := p.source.Admit(p.ctx, environment, use)
	if data.Use != use || data.Environment != environment { return Expression[T]{}, ErrScope }
	if data.Bound.Valid() {                   // a formal in scope
		value, ok := env[data.Bound]
		if !ok { return Expression[T]{}, ErrScope }
		return value, nil
	}
	// Whole-source admission applies even to arguments a representation omits.
	for _, a := range data.Actuals {
		p.source.Authenticate(... a.Binder, a.Use ...)
		if _, err = p.source.Admit(p.ctx, environment, a.Use); err != nil { return ..., err }
	}
	selector, err := p.source.UseSelector(data)
	if data.Callee.Valid() { selector, _, err = p.source.DeclarationSelector(data.Callee) }
	if choice, ok := single(p.bindings[selector.Key()]); ok {
		// an exact binding: capture children, each checked with CheckChild, then map
		...
	}
	if !data.Callee.Valid() { return Expression[T]{}, ErrUnsupported }
	id, err := p.declare(data.Callee, depth+1)       // generic or plain declaration
	open := p.opens[data.Callee]
	if len(data.Actuals) != len(open.Formals) { return Expression[T]{}, ErrScope }
	// pair actuals with the callee's real binders; map only those the native
	// convention needs, then fill native arguments from the convention
	...
	return Expression[T]{kind: DeclarationExpression, declaration: id, args: args}, nil
}
```

Construction (`Run`) later receives only the closed native program. It has no source
access and no mapping callbacks.

### The horos source integration

At open time (`New`), the horos integration copies the whole graph and pre-walks every
type-definition body to build a term-to-declaration ownership map. It also indexes the
facts it interprets: `sort` (on parameters and terms), `result-sort` (on declarations)
and `identity` (a module's content pin, a 32-byte hash).

```go
for i, e := range s.facts.Entities {
	id := semantics.EntityID(i + 1)
	if e.Kind != (semantics.Key{Owner: "horos", Name: "type-definition"}) { continue }
	if e.Owner == 0 || e.Body == 0 || len(e.Declarations) != 1 { return nil, semantic.ErrInvalidRequest }
	d, _ := s.declarations.Ref(uint64(id))
	...
	for _, p := range e.Parameters {
		b, _ := s.binders.Ref(uint64(p.Entity))
		s.bound[p.Entity] = b
		u, _ := s.uses.Ref(uint64(len(s.facts.Terms)) + uint64(p.Entity)) // a fabricated "formal subject" use
		s.formalSubjects[p.Entity] = u
	}
	// The graph owns lexical term sites. Callee references are not expanded.
	var walk func(semantics.TermID, int) error
	walk = func(t semantics.TermID, depth int) error {
		if old := s.termOwners[t]; old != 0 {
			if old != id { return semantic.ErrInvalidRequest } // a term shared by two declarations is rejected
			return nil
		}
		s.termOwners[t] = id
		term := s.facts.Terms[t-1]
		for _, a := range term.Arguments { walk(a, depth+1) }
		for _, a := range term.Operands  { walk(a.Term, depth+1) }
		return nil
	}
	walk(e.Body, 0)
}
```

`Open` and `Admit` then combine graph navigation with horos judgments:

```go
func (s *Source) Open(ctx context.Context, d g.Declaration[h.Language]) (g.OpenData[h.Language], error) {
	// authenticate d ...
	e := s.facts.Entities[i-1]
	if err := s.proper(s.results[semantics.EntityID(i)], e.Body); err != nil { ... } // horos judgment
	env, _ := s.environments.Ref(i + 1)
	o := g.OpenData[h.Language]{Declaration: d, Environment: env, Body: s.use(e.Body)} // navigation
	for _, p := range e.Parameters {                                                  // navigation
		if err := s.proper(s.sorts[semantics.Subject{Entity: p.Entity}], e.Body); err != nil { ... } // judgment
		o.Formals = append(o.Formals, g.Formal[h.Language]{Binder: s.bound[p.Entity], Subject: s.formalSubjects[p.Entity]})
	}
	return o, nil
}

func (s *Source) Admit(ctx context.Context, env g.Environment[h.Language], use g.Use[h.Language]) (g.UseData[h.Language], error) {
	u := g.UseData[h.Language]{Use: use, Subject: use, Environment: env}
	// authenticate use and env ...
	if formal != 0 {                                      // the fabricated formal-subject use
		if s.facts.Entities[formal-1].Owner != semantics.EntityID(e-1) { return u, g.ErrScope }
		u.Bound = s.bound[formal]
		return u, nil
	}
	id := semantics.TermID(ti)
	if s.termOwners[id] != semantics.EntityID(e-1) { return u, g.ErrScope }       // scope (navigation)
	if err := s.proper(s.sorts[semantics.Subject{Term: id}], id); err != nil { return u, err } // judgment
	t := s.facts.Terms[id-1]
	if t.Kind == semantics.Reference || t.Kind == semantics.Application {         // navigation
		target := t.Target
		if t.Kind == semantics.Application {
			callee := s.facts.Terms[t.Callee-1]
			if callee.Kind != semantics.Reference { return u, s.Unsupported(use, "computed constructor application") }
			target = callee.Target
		}
		if b := s.bound[target]; b.Valid() {               // a reference to a parameter
			if t.Kind != semantics.Reference || s.facts.Entities[target-1].Owner != semantics.EntityID(e-1) { return u, g.ErrScope }
			u.Bound = b
			u.Subject = s.formalSubjects[target]
			return u, nil
		}
		d := s.decls[target]
		if !d.Valid() { return u, s.Unsupported(use, "constructor-valued application") }
		u.Callee = d
		parameters := s.facts.Entities[target-1].Parameters
		if len(parameters) != len(t.Arguments) { return u, semantic.ErrInvalidRequest }
		for i, p := range parameters {                     // pairing actuals with real binders
			u.Actuals = append(u.Actuals, g.Actual[h.Language]{Binder: s.bound[p.Entity], Use: s.use(t.Arguments[i])})
		}
	}
	// The exact successful Horos analysis admits the entire original occurrence,
	// including every actual, under this declaration's original assumptions.
	return u, nil
}

// proper: the sort fact must be the horos provider expression {horos, sort/proper}.
// Constructor-valued types (sorts like * -> *) need an explicit representation.
func (s *Source) proper(sortID, term semantics.TermID) error {
	if sortID == 0 || int(sortID) > len(s.facts.Terms) { return semantic.ErrInvalidRequest }
	t := s.facts.Terms[sortID-1]
	if t.Kind != semantics.ProviderExpression || t.Operator != (semantics.Key{Owner: "horos", Name: "sort/proper"}) {
		return s.Unsupported(s.use(term), "constructor-valued types require an explicit representation")
	}
	return nil
}

func (s *Source) CheckChild(ctx context.Context, env g.Environment[h.Language], parent, child g.Use[h.Language]) error {
	p, err := s.Admit(ctx, env, parent)
	if _, err = s.Admit(ctx, env, child); err != nil { return err }
	for _, a := range p.Actuals { if a.Use == child { return nil } }       // a generic actual
	v, err := s.Read(parent)
	for _, o := range v.Operands { if o.Use == child { return nil } }      // a provider operand
	return g.ErrScope
}
```

Its catalogue functions are horos interpretation:

```go
func Named(language representation.Language, modulePin, name string) (g.Selector, error) {
	return g.NewSelector(language, "declaration/"+modulePin+"/"+name)
}
func Builtin(language representation.Language, constructor, discriminator string) (g.Selector, error) {
	return g.NewSelector(language, fmt.Sprintf("builtin/%d:%s/%d:%s", len(constructor), constructor, len(discriminator), discriminator))
}
func (s *Source) DeclarationSelector(d g.Declaration[h.Language]) (g.Selector, bool, error) {
	e := s.facts.Entities[...]
	pin, err := s.modulePin(e.Owner)          // reads the module's identity fact
	key, err := Named(s.snapshot.Language(), pin, e.Name)
	return key, true, err
}
func (s *Source) UseSelector(use g.UseData[h.Language]) (g.Selector, error) {
	if use.Bound.Valid() { return Bound(s.snapshot.Language()) }
	view, err := s.Read(use.Use)              // decodes {horos, descriptor/<constructor>} and its payload
	return Builtin(s.snapshot.Language(), view.Constructor, view.Discriminator)
}
```

Beyond the contract, the horos session type exposes a reading API that horos mapping
libraries use:

```go
type Reader interface {
	g.SourceSession[h.Language]
	Snapshot() semantic.Snapshot
	Declaration(g.Declaration[h.Language]) (Declaration, error) // name, module pin, root-module flag, origin site, body
	Parameter(g.Binder[h.Language]) (Parameter, error)         // name and declaration site
	IsBody(g.Environment[h.Language], g.Use[h.Language]) bool  // is this use the declaration's body?
	Read(g.Use[h.Language]) (Expression, error)                // constructor, discriminator, origin, operands
	Operand(g.Use[h.Language], ...uint32) (g.Use[h.Language], error) // operand at a payload locator
	Fields(g.Use[h.Language]) ([]Field, error)                 // record field name, optional flag, type use, origin
	Unsupported(g.Use[h.Language], string) error               // a diagnostic at the use's site
}
```

### How mapping libraries use the source today

The horos declaration recipe (shared by every horos-to-native target):

```go
Expand: func(ctx context.Context, r g.DeclarationRequest[h.Language, T, source.Reader]) (g.DeclarationDescription[T], error) {
	src, _ := r.Source()
	open, _ := r.Open()
	decl, err := src.Declaration(open.Declaration)       // horos identity
	name, err := Exported(decl.Name)
	if !decl.Root { name = "M" + decl.ModulePin[:16] + "_" + name } // imported module
	body, err := src.Read(open.Body)
	form := alias
	if body.Constructor == "record" { form = record }
	desc := g.DeclarationDescription[T]{Name: name, Form: form, Origin: decl.Origin}
	for i, f := range open.Formals {
		parameter, err := src.Parameter(f.Binder)
		pname, err := Exported(parameter.Name)
		desc.Generic.Parameters = append(desc.Generic.Parameters, g.ParameterDescription{Name: pname, Origin: parameter.Origin, RequiredUse: true})
		desc.Generic.Bindings = append(desc.Generic.Bindings, g.Parameter[T](i))
		desc.Generic.Arguments = append(desc.Generic.Arguments, g.Argument[T](i))
	}
	desc.Body, err = r.Mapped(open.Body)
	return desc, err
}
```

The record rule for TypeScript (Go, Rust and C++ are analogous):

```go
Expand: func(ctx context.Context, r g.Expansion[h.Language, s.Language, source.Reader]) (g.Expression[s.Language], error) {
	src, _ := r.Source()
	use, _ := r.Use()
	if !src.IsBody(use.Environment, use.Use) {
		return g.Expression[s.Language]{}, src.Unsupported(use.Use, "nested records require a named declaration in this profile")
	}
	view, _ := src.Read(use.Use)
	fields, _ := src.Fields(use.Use)
	var members []g.Expression[s.Language]
	for _, f := range fields {
		ty, err := r.Mapped(f.Use)          // kernel checks f.Use is a child via CheckChild
		members = append(members, g.At(model.Field(f.Name, f.Optional, ty), f.Origin))
	}
	return g.At(model.Record(members...), view.Origin), nil
}
```

The list element capture reads the operand at payload locator `[1]`:

```go
p.Capture = func(ctx context.Context, src source.Reader, u g.UseData[h.Language]) ([]g.Capture[h.Language], error) {
	operand, err := src.Operand(u.Use, 1)
	return []g.Capture[h.Language]{{Slot: slot, Use: operand}}, nil
}
```

A library that owns `Cell<V>` maps every use of that declaration to an imported native
type. It is keyed by the declaration selector and captures the actual through the
source's pairing:

```go
key, _ := source.Named(from, pin, "Cell")
pattern := g.Pattern[h.Language, source.Reader]{Selector: key, Slots: []g.Slot{slot},
	Capture: func(ctx context.Context, src source.Reader, use g.UseData[h.Language]) ([]g.Capture[h.Language], error) {
		if len(use.Actuals) != 1 { return nil, g.ErrScope }
		return []g.Capture[h.Language]{{Slot: slot, Use: use.Actuals[0].Use}}, nil
	}}
// Output: model.External("@bitspark/cell", "Cell", g.Mapped(slot))  =>  TypeScript Page<Cell<string>>
```

### The independent source

To show that the kernel is not horos-specific, the repository has an independent test
source (`integration/go/testdata/generationfixture`). It has its own toy language and
its own mapping library. It does **not** publish a semantic graph. It constructs
`OpenData` and `UseData` records by hand in tables and returns them from `Open` and
`Admit`:

```go
func (s *Source) Admit(ctx context.Context, env g.Environment[SourceLanguage], u g.Use[SourceLanguage]) (g.UseData[SourceLanguage], error) {
	r, ok := s.useRecords[u]
	if !ok || r.data.Environment != env { return g.UseData[SourceLanguage]{}, g.ErrScope }
	for _, a := range r.data.Actuals {
		actual, exists := s.useRecords[a.Use]
		if !exists || actual.data.Environment != env { return g.UseData[SourceLanguage]{}, g.ErrScope }
	}
	return r.data.Clone(), nil
}
```

Its `CheckChild` also accepts a declaration's formal subjects as children of the body,
and it builds a second occurrence with the same effective subject to show that the
kernel must not equate occurrences by handle. It therefore proves composability through
today's interface, but it re-implements the same navigation responsibility in a
different form. Under the target design it has to publish a semantic graph and supply
only its own interpretation.

### Existing qualification

All planning tests enter through the public registry. They cover the nested
Page/Cell output, a library-owned Cell, reversed native arguments, a native convention
that omits an argument, ambiguous mappings, failing mappings without fallback,
expired expansion requests, unsupported horos forms, invalid horos programs
(rejected because their evidence is unavailable), name collisions, and a depth
limit. Several guarantees have no negative test that reaches the planner: the
omitted-argument case fails earlier, at evidence admission, so the planner's
admission of an omitted actual has never been seen to fail. Foreign or forged
records, a wrong environment, a capture that returns a non-child, and recursion are
also untested. The replacement should be qualified against those cases directly.

### Guarantees that must survive

These come from earlier consultations and the delivered SPI. Each is a requirement on
the replacement, independent of the current method names.

| Guarantee | Meaning |
| --- | --- |
| Capture identity | Planning consumes exactly the admitted snapshot and evidence. A recapture of equal bytes, or evidence from another snapshot, is rejected. |
| Real binders | A source binder is the actual declared parameter, never a name or position. Native conventions may reorder, omit or add native parameters, so binder *i* is not native parameter *i*. |
| Environments and assumptions | A use is interpreted under its declaration's own assumptions (for example `D: h.Data`). A product built while defining a callee must not cross into a caller, and an assumption must not become an unconditional fact. |
| Foreign-reference rejection | Records, handles or terms from another run, another source or another graph cannot be passed off as belonging to this run. Possession of an address is not admission. |
| Whole-use admission | A use is admitted as a whole, including every actual, before any mapping represents it. An actual that the native convention omits (for example a phantom parameter) still needs admission. Otherwise a rule that emits constant text could bypass a source bound by never traversing its argument. |
| No hidden analysis | Generation never runs a resolver or analyzer. Missing evidence fails; it is never computed on demand. |
| Unknown is not false | An uncovered or absent fact is distinct from a failed judgment. Absence of a fact is not a negative answer. |
| Limits and whole-output failure | Work, depth, references, formals, reservations and output size are bounded. Any failure, cancellation or exhausted limit yields no usable output and never triggers a weaker fallback mapping. |
| Selector ownership | Stable selector identity belongs to the source language. The kernel compares selector identity and never spellings or graph-local IDs. |

### Why the current boundary is a defect

The project records this as architecture case ACX-006 and gap GEN-01. The relevant
architecture rules (paraphrased):

- **AC-02: meanings remain with their semantic owners.** Independently owned parts
  supply their interpretations; a shared mechanism owns its invariant guarantees.
- **AC-05: composition preserves guarantees.** Openness may not be bought with
  arbitrary payloads, unchecked claims, or callbacks that discard guarantees.
- **AC-06: limited coverage does not redefine the architecture.** Working examples
  through the first reader do not make the wrong responsibility split acceptable.
- **AC-09: design interfaces independently of the implemented surface and the effort
  to change it.** The current methods, handles and records have zero weight, including
  as a tiebreaker. A retained element needs the same independent justification as a
  new one.

The intended division for `People = Page<Cell<Person>>`:

| Decision | Intended owner |
| --- | --- |
| Follow established references to Page, Cell and Person | Kernel, over the captured graph |
| Follow `Page.T` to `Cell<Person>` and `Cell.V` to Person | Kernel, over established application bindings |
| Interpret horos constructors, kind judgments and pinned identity | horos integration |
| Select a native representation of Cell | Cell's mapping library |
| Admit source evidence and preserve capture/scope association | Shared contract guarantees, with source-owned judgments |

The completion criteria the correction must meet:

1. The kernel consumes declaration, reference and application structure through the
   shared semantic context. No source integration supplies a parallel traversal.
2. Narrow source-owned contracts keep constructor/operand interpretation, stable
   selector identity and applicability judgments. Unknown or uncovered facts stay
   distinct from failed judgments; there is no hidden analysis.
3. Capture identity, real binders, source environments, foreign-reference rejection,
   limits and whole-output failure stay enforced. Omitted native arguments still
   require source admission.
4. horos and an independent source exercise the same traversal, including nested
   `Page<Cell<Person>>` and invalid scope and evidence cases.
5. The superseded reader interfaces and duplicated traversal are removed.

### Constraints

- **No hidden analysis, no new semantic algebra.** The shared graph has "no scope
  graph, substitution map, universal type algebra or conditional proof system", and
  the kernel must not add one. Operand locators are provider-defined paths that the
  kernel does not decode.
- **Kernel interprets nothing language-specific.** No horos names, operators, entity
  kinds, fact relations or sort terms may appear in the kernel. A new source language
  is added by contributing its own interpretation and mapping libraries, without
  changing the kernel.
- **Mapping libraries stay declarative and source-typed.** A horos mapping library may
  use horos's own reading API (that is horos interpretation). Libraries describe native
  representations; the planner resolves them before construction, and construction
  never sees the source.
- **No compatibility aliases or retained superseded interfaces.** The replacement
  removes the old reader surface; there are no deprecated shims.
- **Single process, trusted code.** Mapping libraries and source integrations are
  trusted in-process Go code, but the contract still has to make scope and identity
  mistakes detectable rather than silently wrong.
- **Scale.** Inputs are bounded by explicit limits (default graph caps: 100,000 records,
  800,000 edges; default planning caps: 100,000 operations, depth 128). Correctness and
  clear ownership matter more than speed.
- **Out of scope here:** native enforcement of source constraints in generated APIs
  (a separate gap, GEN-02), recursive types (currently rejected), and file placement.
  The design should not make those harder.

### What we've considered

**Earlier consultations did not consider kernel-owned traversal.** The 2026 consultation
that produced `Open`/`Admit` framed the constraint as "without turning the agnostic host
into a semantic evaluator". Its expert wrote: "The source integration is responsible for
the meaning of the substitution, arity, saturation, applicability, and required
judgments. The driver does not infer any of those from the shape of `semantics.Term`."
An earlier consultation said a structural binding query "establishes correspondence,
not applicability". We read those as compatible with kernel traversal *if* the kernel
follows structural correspondence only and asks the source for every judgment. We
would like that reading confirmed or corrected.

**A candidate split** (ours, for critique; not a decision):

| Concern | Candidate owner | Notes |
| --- | --- | --- |
| Take the graph from the evidence entry whose requirement sets `Semantics: true` (exactly one), and unpack it as `Semantics[F]` | Kernel | Generic. Rules out a source supplying a different graph than the admitted one. |
| Resolve a requested root syntax site to the entity that declares it | Kernel (`Entity.Declarations`) | The source still classifies whether that entity is a generatable declaration. |
| Open a declaration: its body term and real binders (`Entity.Parameters`) | Kernel | |
| Walk a body: Reference to an own parameter is a bound formal; Reference to another entity's parameter is a scope error; Reference or Application to another entity is a use of that callee, paired with real binders by `DirectParameterBindings`; no structural view (computed callee, parameter-headed application, unpaired arity) is *unsupported*, not a failed judgment | Kernel | Arity, applicability and saturation *meaning* stay source-owned. The kernel only notices that it cannot build a complete native application from the structure. |
| Child relation: a captured child must be an operand or argument of the current term | Kernel | Replaces `CheckChild`. |
| Classify an entity as a generatable declaration and give its stable selector | Source | horos: kind `{horos, type-definition}`, selector from module pin and name. |
| Give a provider expression's selector (constructor and discriminator) | Source | horos: `{horos, descriptor/...}` plus payload discriminator. |
| Judge each declaration, formal and subject admissible under its declaration's assumptions, using only supplied evidence; distinguish holds / rejected / uncovered | Source | horos: the `sort/proper` checks on result, parameters and terms. |
| Read payloads and operand meanings for mapping libraries | Source (per-language reading API) | horos: record fields, list element at locator `[1]`. |

**Open points within that split:**

1. **The admission judgment's grain.** Per term, per whole use (term plus all actuals),
   or per declaration plus per term? The horos analyzer has already judged the entire
   body when it reports success, so horos's per-term re-check today mainly enforces a
   *profile* restriction (proper sorts only). For an omitted actual, should the kernel
   ask the source to judge the actual as one whole subject, or walk its subtree and ask
   for each term?
2. **Outcome shape.** A three-valued result (holds, rejected with a diagnostic,
   uncovered) versus distinct error kinds. Today absence of a sort fact surfaces as a
   generic invalid-request error, which conflates "uncovered" with "malformed".
3. **Handles.** Mapping libraries and the source's reading API could exchange
   graph-local term IDs that the kernel validates as structural children of the current
   subject. That makes foreign *content* impossible, because the output can only contain
   this graph's structure, but not foreign *handles* detectable. Alternatively the kernel
   could mint owner-qualified, run-scoped handles that the reading API must resolve
   through the kernel. Which identity model fits?
4. **Environment model.** horos has no nested generic declarations, so "the enclosing
   declaration's parameters" suffices. Should the kernel define scope as parameter
   ownership by the declaration being opened, or as an `Owner`-chain rule? And should a
   term shared between two declarations' bodies (the graph permits sharing) be
   accepted, with the environment taken from the walk path, or rejected as horos does
   today?
5. **Occurrence versus effective subject.** Today `UseData` carries both the occurrence
   and a source-certified effective subject. The planner never compares subjects. With
   graph terms, is a separate effective subject still needed, for example for a future
   memoization or recursion rule, and who would own the equivalence?
6. **The mapping library's view.** A library's `Capture` and `Expand` need the current
   subject, the declaration being opened, the formals, and for a declaration-keyed rule
   such as `Cell<V>` the actuals by binder. Should the kernel pass a structural view
   value, and should libraries identify actuals by binder (via the source's parameter
   API) or by position?

**Rejected directions:**

- Moving `Admit` wholesale into the kernel. It mixes navigation with evidence and
  assumption checks that the kernel must not perform.
- Keeping `Open`/`Admit` and adding a shared helper package that sources may call.
  The kernel would still depend on a source-supplied traversal, and nothing would stop
  divergent implementations.
- Unrestricted callbacks (for example, "the source walks and calls back into the
  kernel"). This inverts control but keeps the traversal source-owned and weakens the
  child-relation and scope checks.
- A package move of the current reader. It changes the location, not the
  responsibility.

## Questions for the Expert

1. **Where is the principled line?** Given a shared graph with exactly three term
   variants (reference, application, provider expression with provider-defined operand
   locators) and provider-owned facts, how would you divide the generation walk between
   a language-agnostic kernel and a source-owned interpretation contract? Are there
   patterns from comparable systems you would draw on, such as compilers with a shared
   intermediate representation plus language hooks, or multi-language code-intelligence
   graphs? We are particularly interested in how such systems keep "follow the
   structure" from quietly becoming "judge the structure".
2. **How should admission judgments be shaped?** How would you model the source's
   judgments, including their grain (declaration, whole use, term) and outcome (holds,
   rejected, uncovered), so that whole-use admission still covers actuals a native
   convention omits, unknown facts stay distinct from failures, and the kernel never
   interprets what a judgment means?
3. **What identity model for subjects?** What handle and identity design would you use
   for the subjects passed between the kernel, a source's reading API and mapping
   libraries, so that capture identity, foreign-reference rejection and expiry of scoped
   expansion requests are preserved without the source re-supplying navigation records?
4. **What should mapping libraries see?** How should source-typed mapping libraries
   receive the kernel's structural view (current subject, declaration, formals, actuals
   by real binder) alongside the source's own reading API, including declaration-keyed
   rules such as a library-owned `Cell<V>` that captures its actual?
5. **What else would you watch for?** Looking at this design as a whole, including later
   work on recursive types and on carrying source constraints into native output, what
   risks, missing concepts or simplifications do you see that we have not asked about?
