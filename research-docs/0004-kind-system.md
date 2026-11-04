# Research: a proper kind system for unified Horos and Bittype

**ID:** 0004
**Date:** 4 November 2026
**Status:** draft
**Owner:** codex-hydrated-type-design
**Issue:** https://github.com/Bitspark/bittype/issues/55

## Question

What coherent, implementable kind system should replace Horos's first-order
Type/Data parameter bounds, and how should Bittype expose and consume it without
acquiring Horos-specific rules in its agnostic kernel? The owner requested:

> "I want a proper kind system as well. Let's get advice"

Please give a complete answer to the five questions below in at most 5,000 words.
Prioritize a precise core calculus, counterexamples, canonical representation,
ownership and independently checkable acceptance criteria. We need a sound
design to implement, not a registry of named contracts or an endorsement of the
options below. Distinguish indispensable features from independently deferrable
extensions; a proper kind system need not implement every feature of a general
programming language.

Spend the response budget on rules, worked positive/negative derivations and a
compact native-policy table rather than repeating background. Complete all five
answers; avoid a long introductory report that crowds out the checking rules.

## Context

### System overview and vocabulary

These are independent public repositories, not packages of one compiler:

| Component | Responsibility and inspected state |
| --- | --- |
| Ontos | Immutable ground values: byte atoms and ordered tuples, with standard data embeddings and a frozen canonical byte codec. An optional Ontos/wire package represents mixed values with live sending capabilities. |
| Horos | Types describing Ontos data, ground well-formedness and value conformance; a unified declaration language additionally describes Wire capabilities, modules, generic declarations and applications. Go declaration semantics are authoritative. Ground checkers also exist in Rust and TypeScript. |
| Bittype | Go-implemented AST/CST infrastructure, readable Horos syntax, semantic-provider interfaces, an agnostic syntax engine, four native syntax frontends and types-only Horos-to-native transformations. |
| Bitwire/Bitruntime | Live capability recognition, transport and hydration (exported references and imported proxies). This consultation changes none of these operational contracts. |
| Ordinary libraries | Function, Stream, Cell and future behavioral contracts. They compose types and own protocols/native representations. General tools must never recognize these names or their structural fingerprints to select behavior. |

AST means abstract syntax tree. CST means concrete syntax tree, retaining source
tokens/trivia and an AST view. A provider implements a separately registered
processing stage. A semantic receipt binds established reports and an immutable
fact graph to the exact captured syntax, configuration and dependencies.

A **value** inhabits a **type**. A **type constructor** accepts type-level
arguments and produces a type-level result. A **kind** would classify these
type-level expressions and their applications. In existing Go code, `Kind` also
names an AST/expression variant tag; that use is unrelated to the proposed kind
system. A generic **binder** identifies a parameter independently of its spelling.
`Data` currently means a statically proved data-only type, not the value type
`any`, a runtime wrapper or a capability-authority claim.

The current canonical bound `type` admits proper type expressions describing
data or capabilities. In candidate signatures below, `Type` denotes the
proper-type classifier, not a top value type. Whether `Data` becomes a subkind
or remains a separate predicate/refinement is deliberately undecided.

Inspected baselines, both clean on main:

- Horos `e2896317a1f4b080cf0c4e8e46234d6ad680d141`, unified declaration revision 2,
  publicly distributed as Go module v0.9.0.
- Bittype `c143440bbf703a4ae02add267a5699c2bf7374c9`, depending on Horos v0.9.0 and
  Ontos v0.16.0. Its pagination source and exact main passed hosted Go tests,
  vet/race/fuzz, installed-consumer checks and four native compiler checks in
  GitHub run 41838812640. This proves the present profile, not higher kinds.

No local build, performance measurement or running service inspection is used
as evidence here. The workload is finite type/declaration documents; limits are
explicit. A neighboring task owns source-layout/multi-file design only and is
outside this consultation.

### Existing canonical model and examples

Canonical declaration documents are ground Ontos values. Even a description of
a Wire is data: it contains no endpoint. In the notation below, `A(s)` is a byte
atom containing UTF-8 s, `T(s)` and `I(n)` are text/integer embeddings, `L(...)` is
the standard list embedding, and `M{...}` is a canonical text-keyed map.
`H(c,...) = Tuple(A("horos/" + c),...)` and
`D(c,...) = Tuple(A("horos/decl/" + c),...)`. These are equations, not builders
executed by the compiler. `Pin(M)` is SHA-256 of exact canonical module bytes.

Current exact grammar (Horos `docs/spec/declarations.md:49-120`):

```text
Module = D("module", I(2), T(name), importsMap, declarationsMap)
Import = D("import", T(expectedModuleName), Atom(32-byte pin))
Declaration = D("type", L(D("parameter", T(name), A(bound)), ...), body)
bound = "type" | "data"

E = Horos descriptor with E in designated child-type positions
  | H("wire", E)
  | D("param", I(zeroBasedSlot))
  | D("apply", target, L(E, ...))
target = D("local", T(name)) | D("imported", T(alias), T(name))
```

There is one ordered parameter scope per declaration. Applications supply every
argument; targets are named declarations, not arbitrary expressions. Bare
generic declarations, constructor parameters, partial application, inference,
defaults, value parameters and higher-kinded slots are rejected. All declarations
are exported. Imports contain an expected module name and immutable pin; a
complete reachable import closure is supplied by the caller, with no implicit IO.
General declaration dependency cycles are rejected. Ground descriptor `rec`/`ref`
support guarded data recursion separately; they are not generic binders.

Bittype's implemented readable syntax expresses this existing example:

```text
module example.pagination;
type Page<T> = { items: list<T>; next?: text; };
type Person = { name: text; active: bool; };
type PeoplePage = Page<Person>;
type NestedPage = Page<Page<Person>>;
type DataPage<D: Data> = Page<D>;
type Ignore<T> = text;
type Send<T> = wire<T>;
```

An omitted source bound means explicit canonical `type`, not an omitted generic
argument or a substitution of `any`. A closed record admits data/capability
fields; optional means absent or present, not nullable. `any` is the existing
type of all **data** values. `wire<M>` describes a sending capability accepting
messages of type M; it is never data, even when M is data. Tuple, compound, list,
closed record and union classify as data when their value-bearing children do.
Other existing constructors (sets/maps, open records, intersection, complement,
guarded recursion and bare sequences) retain data-only operands in revision 2.

These judgments must remain distinct:

```text
Page<Person>                   data
Page<wire<Person>>             admitted type, not data
DataPage<Person>               admitted, data
DataPage<wire<Person>>         rejected bound
Ignore<wire<Person>>           data (argument does not occur in result)
Send<Person>                   not data
Forward<D: Data> = DataPage<D>  valid open generic
Forward<T> = DataPage<T>        invalid open generic before specialization
```

The motivating missing capability is passing a constructor itself. These are
candidate source spellings, not currently implemented syntax:

```text
type Wrapped<F: Type -> Type, T: Type> = { value: F<T>; };
type WrappedPage = Wrapped<Page, Person>;
type Compose<F: Type -> Type, G: Type -> Type, T: Type> = F<G<T>>;
type Nested = Compose<Page, Page, Person>;
type DataWrapped<F: Data -> Data, D: Data> = { value: F<D>; };
```

The spellings do not select the eventual calculus, but these intended outcomes
state the precision we want to preserve when known constructor arguments bind:

| Candidate use | Required judgment or design obligation |
| --- | --- |
| `Wrapped<Page, Person>` | Admitted; the known application is data. An abstract F's broad bound must not permanently erase known-argument precision. |
| `Wrapped<Ignore, wire<Person>>` | Admitted and data; Ignore's unused argument does not contaminate its result. |
| `Wrapped<Send, Person>` | Admitted, not data. |
| `Compose<Page, Page, Person>` | Admitted and data. |
| `Wrapped<Person, Person>` | Reject: a proper type is not a unary constructor. |
| `DataWrapped<Page, Person>` | The one Page must be usable with a proved data-preserving guarantee; propose how the signature/evidence expresses that guarantee. |
| `DataWrapped<Send, Person>` | Reject: Send does not guarantee a data result. |
| `DataPage<F<D>>` under only F: Type -> Type and D: Data | Reject the generic body; a future known caller cannot repair an invalid open definition. |
| The same body under F: Data -> Data and D: Data | Admit using the constructor's explicit result guarantee. |

Distinguish a guarantee necessary to check an open generic body from stronger
facts available for a well-formed body's concrete application. If a sound
alternative intentionally sacrifices one desired precision property, explain
the tradeoff and its user-visible consequence rather than silently widening it.

Literal positions are opaque. For example `H("exact", D("param", I(0)))`
describes one literal tuple; it contains no free generic variable. Record labels,
presence markers, compound labels and scalar/refinement operands are not all
type positions. Kind checking must not discover slots by inspecting arbitrary
Ontos subtrees. Descriptors preserve exact operands, ordering and metadata.

### Relevant Horos implementation

Horos `core/go/internal/declmodel/model.go:20-44,138-159` (actual fields):

```go
type DeclarationID struct {
    Module Digest // Digest is [32]byte: canonical module bytes, not type equality
    Name   string
}
type ParameterID struct {
    Declaration DeclarationID
    Index       int
}
type Bound string
const (
    TypeBound Bound = "type"
    DataBound Bound = "data"
)
type TypeParameter struct {
    Name  string
    Bound Bound
}
type Operand struct {
    path       []int
    expression *Expression
}
type Expression struct {
    path        []int
    kind        Kind // descriptor, parameter, application
    value       core.Value // an immutable Ontos value
    constructor string
    operands    []Operand
    parameter   ParameterID
    target      DeclarationID
    arguments   []*Expression
}
```

Public accessors copy mutable storage; the checked Program preserves open
generics and imported identity. The exported `declarations` package exposes
this model under the one fixed Horos grammar.

Horos `core/go/internal/declmodel/resolve.go:286-303`:

```go
if args, ok := r.program.grammar.form(value, "param", 1); ok {
    slot, ok := index(args[0])
    if !ok || slot >= len(owner.parameters) {
        return nil, nil, failure(Invalid, "unbound or out-of-range type parameter")
    }
    return &Expression{kind: Parameter, value: value,
        parameter: ParameterID{Declaration: owner.id, Index: slot}},
        r.program.grammar.Any(), nil
}
if args, ok := r.program.grammar.form(value, "apply", 2); ok {
    target, err := r.target(args[0], owner.id)
    if err != nil { return nil, nil, err }
    arguments, err := data.ReadList(args[1])
    if err != nil || len(arguments) != len(r.program.declarations[target].parameters) {
        return nil, nil, failure(Invalid, "application argument list or arity mismatch")
    }
    result := &Expression{kind: Application, value: value, target: target}
    r.edges[owner.id] = append(r.edges[owner.id], target)
    // The remaining branch reads each argument in the caller's scope,
    // records its path, checks its recursion closure and appends it.
}
```

`target` accepts only local/imported declaration forms and verifies the exact
resolved declaration. The shown `Any()` is an internal syntax-check skeleton
after checking bindings, not an inferred argument or a widening of value meaning.
The reader uses Horos's descriptor operand positions; afterward it rejects
declaration cycles and checks every generic body under its declared bounds.

The current data analysis is more precise than a flat Type/Data flag on a
declaration. Horos `core/go/internal/declmodel/bounds.go:3-51`:

```go
// A condition is a conjunction of caller-owned Type slots. Impossible denotes
// a capability-bearing type; an empty conjunction proves Data.
type dataCondition struct {
    impossible bool
    requires   map[ParameterID]bool
}
func (c dataCondition) proven() bool {
    return !c.impossible && len(c.requires) == 0
}

// Excerpt from declaration(), after a bounded-work step and cache lookup:
env := make([]dataCondition, len(declaration.parameters))
for i, parameter := range declaration.parameters {
    if parameter.Bound == TypeBound {
        env[i].requires = map[ParameterID]bool{
            {Declaration: declaration.id, Index: i}: true,
        }
    }
}
result, err := c.domain(declaration.expression, declaration, env, depth)
if err == nil { c.summaries[declaration.id] = result }
```

An application first checks every argument, including unused ones, against the
callee's bounds. It then substitutes the actual argument conditions only into
the callee's used-slot summary. This proves Page's conditional data result and
Ignore's unconditional data result without expanding shared declaration bodies
repeatedly. It cannot presently summarize an unknown constructor parameter F.

Horos `core/go/declarations/grammar.go:101-115` defines the source-owned rule:

```go
func descriptorRule(constructor string, args []core.Value) model.DomainRule {
    switch constructor {
    case "wire":
        return model.DomainRule{Capability: true}
    case "tuple", "compound", "list-of", "union":
        return model.DomainRule{}
    case "record":
        if len(args) > 0 {
            if mode, ok := args[0].(core.Atom); ok && string(mode.Bytes()) == "closed" {
                return model.DomainRule{}
            }
        }
    }
    return model.DomainRule{RequiresData: true}
}
```

These are Horos primitive constructors, not library protocol cases. Capability
means never data; RequiresData checks child obligations. Literal payloads are
not children for this purpose. The ground conformance checker is not widened.

Current `Program.Specialize(id DeclarationID, arguments []core.Value)`
(`core/go/internal/declmodel/specialize.go:8-113`) accepts exactly the declaration's
number of closed descriptor arguments. It checks well-formedness/Data bounds,
recursively replaces parameters, expands named applications in their definition
scope and reconstructs descriptor operands. It charges the actual output tree,
including repeated uses, and checks the closed result. A failure returns nil.
`DataDescriptor` additionally requires the result to be data before returning a
ground descriptor suitable for existing WF/Conforms. Open declarations remain
unchanged. Constructors as arguments will require a deliberately different
closure/application model, not passing arbitrary unvalidated Ontos values here.

### Bittype syntax, semantic graph and transformations

Bittype's stage-first frontend is `syntax/horos`, `parse/horos`, `check/horos`,
`resolve/horos`, `analyze/horos`, `format/horos`, `provider/horos` (each has a Go
package). Syntax admission is separate from semantic validity. The parser
explicitly rejects applying a parameter (`parse/horos/go/parse.go:331-335`).
Application AST nodes currently carry a local/imported name and argument
children; they have no callee-expression child. The syntax package imports no
Horos semantic implementation. Private bridge code reconstructs canonical
values, invokes Horos resolution, checks supplied binding annotations and
projects authoritative results without mutating the source AST.

The generic fact carrier in `semantics/go/model.go:21-67` is:

```go
type EntityID uint32 // model-local, one-based; not canonical identity
type TermID uint32
type Key struct{ Owner, Name string } // provider-owned vocabulary
type Entity struct {
    Kind         Key
    Name         string
    Owner        EntityID
    Declarations []ast.NodeRef // provenance into exact captured ASTs
    Parameters   []Parameter
    Body         TermID
}
type Parameter struct { Entity EntityID; Sort Key }
type TermKind uint8
const (
    Reference TermKind = iota + 1
    Application
    ProviderExpression
)
type Term struct {
    Kind      TermKind
    Target    EntityID
    Arguments []TermID
    Operator  Key
    Encoding  Key
    Data      []byte
    Operands  []Operand
}
type Operand struct { Locator []uint32; Term TermID }
```

`Entity.Kind` classifies an entity with a provider-owned key; `Parameter.Sort`
classifies what sort of argument its slot accepts, also with an opaque key.
`Term.Kind` selects one graph variant. None of these fields currently holds a
structural arrow kind or a kind scheme; equal spellings do not imply a common
language meaning.

`Reference` names an entity. `Application` names an entity plus effective
positional arguments. `ProviderExpression` carries a provider operator, encoding,
opaque bytes and explicitly located child terms. One variant is active at a
time. The immutable graph validates edges, ownership, term-containment acyclicity and budgets,
never language lookup, type equality or subtype rules. Entity-reference and
dependency graphs may be cyclic; those cycles have language-owned meaning.
However, its current
application check assumes every application is fully applied to a known entity
(`semantics/go/model.go:247-259`):

```go
case Reference, Application:
    if !validEntity(term.Target) || term.Operator != (Key{}) ||
        term.Encoding != (Key{}) || len(term.Data) != 0 || len(term.Operands) != 0 {
        return bad("invalid reference/application fields")
    }
    if term.Kind == Reference && len(term.Arguments) != 0 {
        return bad("reference has arguments")
    }
    if term.Kind == Application &&
        len(term.Arguments) != len(data.Entities[term.Target-1].Parameters) {
        return bad("application arity does not match effective parameters")
    }
    // Every argument must also refer to a present term.
```

The graph could structurally target a parameter entity with its own parameters,
but it has no general callee TermID for applying the result of an expression.
We want a principled agnostic representation, not fabricated declarations or
opaque payloads that hide ordinary application structure from all consumers.

Bounds already remain provider-owned. `internal/horosmodel/export.go:71-82`
constructs the following term, associated with a parameter's syntax occurrence:

```go
term := semantics.Term{Kind: semantics.ProviderExpression,
    Operator: e.frontend.key("bound/" + string(bound)),
    Encoding: semantics.Key{Owner: "ontos", Name: "codec-v1"},
    Data: codec.Encode(core.NewAtom([]byte(bound)))}
```

All parameter `Sort` keys are currently `horos/type` (lines 183-197), with Data
recorded separately. The engine does not know Type, Data, Wire or library names.
This is not an existing general kind calculus disguised by a different name.

The explicit production flow is:

```text
source text -> parser -> CST with AST view -> syntax check/capture
  -> registered resolver -> Horos checked Program + semantic facts
  -> registered analyzer consuming exact resolver receipt
  -> transformer with same snapshot and explicit established analysis
  -> native AST -> formatter -> native CST -> generic printer -> source text
```

The analyzer does not hide a second resolver. Facts include resolved module,
declaration and parameter identities; descriptor operands/payloads preserve
labels, presence and closedness. An engine transform route recaptures each
intermediate AST and does not carry stale evidence forward or invoke analysis
implicitly. The kernel validates supplied receipts/profiles, not Horos semantics.

`contracts/transform/go/transformer.go:17-35`:

```go
type Input struct {
    Source   semantic.Snapshot
    Evidence []semantic.Result
}
type Contract struct { Evidence []semantic.EvidenceRequirement }
type ASTTransformer interface {
    FromID() language.ID
    ToID() language.ID
    Transform(context.Context, Input) (representation.AST, error)
}
```

Four qualified mappings consume the semantic model through a private
source-specific reader, `internal/horosgen`, with target-owned AST construction
in `transforms/horos-to-{typescript,go,rust,cpp}/go`. The existing pagination
profile emits an open native generic Page, Person and aliases retaining generic
applications. It does not specialize a copied Page for each caller. Text maps
to string/String/std::string, lists to native slices/arrays/vectors, records to
interfaces/structs, and optional fields to native presence or an explicit Go
presence struct. Bounds are checked on canonical input; native carriers do not
prove all Horos value invariants or enforce Data on arbitrary native arguments.

No mapping currently accepts Wire or higher-kinded source. Native language
capabilities and our AST/mapping coverage are separate: C++ offers
template-template parameters, but the current AST/checker/mapping supports only
ordinary type parameters, without their nested template parameter lists.
The current native mapping also rejects unused generic parameters
(`internal/horosgen/generate.go:206-209`): Ignore is a valid canonical example,
but is not presently an emitted native generic. Nested records must be named
declarations, and the first descriptor-expression subset is text, bool and list
plus references/applications. The broader Horos model is not the native profile.
Go/TypeScript do not directly provide an unapplied type-constructor parameter
corresponding to F in F<T>; Rust associated type constructors require an explicit
representation. We must decide which canonical programs can be represented
directly, require an explicit source specialization/representation policy, or
must fail. No runtime adapter, invented universal native any, raw-source escape
or silent monomorphization is authorized.

### Required use cases and pressure on the design

The earlier Wrapped, Compose and DataWrapped examples motivate these pressures;
their spellings remain candidates rather than a selected kind calculus.

1. Type constructors should be usable as arguments, with arity/kind errors
   caught at the canonical boundary. Applying a constructor parameter must
   preserve the caller's binding and the identity of supplied declarations.
2. The same Page declaration must retain its data-preserving use and its mixed
   use. A second Page, template lifting or hand-authored expanded definitions
   would undo the existing clean replacement.
3. Abstract F needs an honest result contract. From F: Type -> Type and D: Data
   alone, F<D> is not established Data: F could be Send. A useful design must
   express an actual data-preserving guarantee without a name-based whitelist.
4. Simply assigning Page: Type -> Type with Data <: Type is insufficient to
   derive Page: Data -> Data through ordinary arrow subkinding. It narrows the
   input but leaves the result Type. Meanwhile Ignore has stronger unconditional
   data behavior. Do not silently discard current used-parameter precision.
5. If arrow subkinding is chosen, constructor argument acceptance is contravariant
   and result acceptance covariant. A constructor accepting only Data cannot be
   passed where arbitrary Type arguments may be supplied. This is not a claim
   about covariance of values stored in native containers.
6. Decide whether partial application, constructor aliases, type-level lambdas
   and rank-1 kind polymorphism are necessary now or separable extensions. Define
   the boundary explicitly. Avoid an accidental Type:Type universe or unbounded
   type-level evaluation simply to make arrows work.
   A candidate constructor-valued result is `Return<F: Type -> Type> = F`, with
   `Return<Page>` itself a constructor. If this is in the first calculus, test
   its kind and further application; if deferred, test explicit rejection.
7. Separate canonical byte/declaration identity from kind/type equivalence and
   normalization. Preserve exact origins even if application reduction is needed.
   Equivalent expansions do not authorize substituting behavioral libraries.

### Constraints and failure cases

- Horos owns the calculus, canonical grammar, kinds/bounds, binding and
  specialization. Bittype owns readable spelling/projection and native mappings.
  The agnostic carrier may gain general graph forms, but no Horos vocabulary,
  subkinding algorithm, constructor table or implicit provider invocation.
- Ground Ontos, ground Horos WF/Conforms and `any` stay data-only. Static Wire
  remains available in the one Horos declaration language. Function/Stream/Cell
  remain ordinary declarations, with no core registrations or dispatch cases.
- Use a clean revisioned replacement of superseded declarations and tooling;
  no compatibility aliases/readers, copied type models or template-lifting paths.
  Preserve historical evidence and exact immutable released artifacts.
- Preserve lexical scope, same-spelled distinct binders, import pins/names and
  literal opacity. Validate all argument expressions, even when the result
  ignores them. A malformed unused argument cannot bypass kind checking.
- Reject wrong input kinds/arity, unbound constructor parameters, invalid result
  claims, cycles, stale semantic receipts and unsupported native forms without
  partial usable output. Resource exhaustion is incomplete checking, not proof
  that source is invalid or that a Data obligation holds.
- Bound syntax, stored kind structures, comparison/unification (if any),
  normalization, dependency graphs and expanded output. Existing Horos work
  budgets and cached declaration summaries prevent repeated source DAG expansion.
  Strong normalization alone is not a practical resource bound.
- Current Horos defaults are depth 256 (configurable within 1..512), 100,000
  nodes/work steps, 16 MiB of canonical bytes and 256 supplied modules
  (`core/go/internal/declmodel/model.go:112-126`). Raw admission, semantic phases
  and final expanded output have separately enforced budgets. Do not count shared
  representation nodes once and then encode their exponentially repeated uses
  without charging those occurrences.
- Source serialization/decoding and parsing remain free of network/file IO.
  Only caller-supplied immutable dependencies are visible. No runtime authority
  or protocol guarantee follows from type/kind validity.

Minimum independent qualification should include existing pagination unchanged;
Page and Ignore precision; known versus abstract data-preserving constructors;
same-spelled imported constructors; nested applications and any supported partial
applications; adverse arrow compatibility; supported constructor-valued results
or explicit rejection when deferred; kind/binder
scope; literal data resembling new tags; invalid unused arguments; bounded
normalization/output expansion; and clear native support/rejection for all four
languages. Syntax round trips, semantic judgments and native compiler acceptance
are separate evidence. Expected results must be authored independently of the
checker under test.

### What we have considered and prior research

Possible designs, not selected outcomes:

- **Arrow kinds plus a separate Data predicate/evidence.** Kinds classify
  application shape; a compositional proof system describes data preservation.
  This extends existing precision but needs explicit higher-order guarantees.
- **Type/Data base kinds with subkinding and arrows.** Compact surface syntax,
  but Page's conditional result needs principled checking/refinement, multiple
  signatures or quantified kinds; an arrow label alone is insufficient.
- **Bounded kind polymorphism.** It may capture data-preserving constructors
  uniformly, but requires scoping, instantiation, constraints and a finite
  decidable checking strategy. We want evidence that the extra machinery earns
  its place rather than assuming it is necessary.
- **A small explicitly kinded type-lambda calculus.** Uniform constructor values
  and application could remove ad hoc named-target limitations. It raises
  canonical binding, beta/eta equality, normalization and native realization
  questions. Full dependent types or unrestricted type computation are not
  implied by the owner's request.

The existing bound-only implementation is inadequate for the explicit new
constructor-parameter use case. Adding more enum names such as FunctionKind or
CellKind would violate ownership and would not solve higher-order application.
Pushing all new expressions into opaque semantic payloads would preserve a
surface ABI while hiding the structure generators need.

Reuse checks covered local research, Git branches/worktrees, open PRs, platform
number/topic reservations and Nightfall's research/consultation lists. No existing
consultation answers this kind-system question:

- Bittype research 0002, run `run_a409da7b-a708-41bc-bba3-bd0b011d33ff`, established
  provider-owned facts and exact semantic receipts. Its applied table row 24
  deferred higher kinds; full effective positional application was deliberate.
- Research 0003, run `run_31c1813d-dd34-46a5-9873-105a934581fb`, established data,
  identity and runtime-authority boundaries. Separate Horos/Wire languages and
  template lifting in that historical brief have since been removed. They must
  not reappear merely to preserve a historical recommendation.
- Historical 0001 value-type-core favored first-order Value^n -> Value and
  deferred higher kinds. Its independent type algebra and structural identity
  proposal are superseded by Horos and exact declaration identity.

Bounded prior-implementation inspection found no reusable kind checker:
system2 at `66a1226ffb948d4f50b9f63380ba99d61c69e8f2` explicitly labels its
declaration boundary unimplemented (`docs/DECLARATIONS.md`). Nightseam at
`85a238979129a8f361a9c9c05ea38ecd0343a569` has named type/family parameters
and application-sensitive liveness analysis (`internal/model/family.go:174`,
`internal/model/expr.go:77`, `internal/analysis/liveness.go:145`). Its application
heads are named declarations, not arbitrary constructors. Its callable-name
special cases cannot be imported into our library-owned architecture. Its law
that specializing then generating should agree with generically generating then
binding is useful qualification inspiration, subject to actual target support.

Primary-source precedents consulted, with the useful content stated here:

- [Haskell 2010, kinds and type declarations](https://www.haskell.org/onlinereport/haskell2010/haskellch4.html#x10-630004.1.1)
  distinguishes proper types from constructor arrow kinds. Constructor parameters
  such as App f a do not require general kind polymorphism.
- [GHC kind polymorphism](https://ghc.gitlab.haskell.org/ghc/doc/users_guide/exts/poly_kinds.html)
  distinguishes higher-kinded parameters from forall-k polymorphism and explicitly
  separates kind signatures from saturation arity. Horos must choose saturation
  rules rather than infer them from arrows alone.
- [FreeST kinds](https://freest-lang.github.io/tutorial/kinds-and-types/) supplies
  a practical finite proper-kind hierarchy and contravariant-domain,
  covariant-result arrow subkinding. We do not import its session/channel semantics.
- [Refinement Kinds](https://arxiv.org/html/1908.00441v1) studies structural
  predicates over types and entailment-based subkinding; its reflection and solver
  machinery exceed our immediate need. [Higher-Order Subtyping with Type Intervals](https://arxiv.org/abs/2107.01883)
  illustrates the undecidability risk of a substantially richer bounded calculus,
  not a claim that simple finite arrow subkinding is undecidable.
- [TypeScript's higher-kinded parameter proposal](https://github.com/microsoft/TypeScript/issues/1213)
  is not an implemented general constructor-parameter feature. [Go's specification](https://go.dev/ref/spec#Type_parameter_declarations)
  requires instantiation of generic types when used. [Rust associated types](https://doc.rust-lang.org/reference/items/associated-items.html#associated-types)
  can encode a type family through generic associated types, with explicit trait
  and projection obligations. [C++ template parameters](https://eel.is/c++draft/temp.param)
  and [template argument matching](https://eel.is/c++draft/temp.arg.template)
  provide native template-template forms with C++-specific matching. None of
  these mechanisms establishes Horos Data conformance by native compilation alone.

## Questions for the Expert

1. **Core calculus:** What concrete kind grammar, judgments and application
   rules fit these use cases? Recommend a coherent first implementation and
   distinguish necessary constructor parameters, partial application, constructor
   aliases, type lambdas and kind polymorphism from separable future features.
2. **Data precision:** How should Data relate to Type while preserving Page,
   Ignore, Wire and open generic forwarding? Give the rules and counterexamples
   for abstract constructor result guarantees, arrow compatibility and any
   inference/checking boundary; compare the viable alternatives above.
3. **Canonical representation and checking:** What exact expression/kind forms,
   binder structure, equality/normalization policy and bounded algorithm should
   Horos implement? Address constructor references, application callees,
   specialization, literal opacity, cycle handling and revision migration.
4. **Tooling and generation:** What should Bittype's agnostic fact model expose
   versus leave provider-owned, and how should native generation handle the new
   kinds for TypeScript, Go, Rust and C++? Give a concrete incremental mapping
   policy with explicit unsupported cases and tests, preserving identity and
   avoiding implicit specialization or runtime adapters.
5. **Open assessment:** What important requirement, counterexample or simpler
   architectural alternative have we missed? What would you change in the
   framing, and which independent tests would most strongly falsify a flawed
   implementation before it becomes the canonical format?
