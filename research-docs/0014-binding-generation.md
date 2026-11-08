# Research: Binding generation for service contracts declared in a type language

**ID:** 0014
**Date:** 9 November 2026
**Status:** draft
**Owner:** claude-bittype-gaps
**Issue:** https://github.com/Bitspark/bittype/issues/194
**Code baseline:** bittype 6132170e394f2261d6a24c8efd7839c93299fdcc; horos 3a51a06e6db1800a39f94782213609a2d8723d10 (Go core 0.10.2, npm 0.9.0); ontos 0.16.0; bitwire 0.7.0; bitruntime 0.8.0; platform `libraries/service` 24fedd4 (service-api 0.3.2, service-contracts 0.1.2, service-cli 0.6.2, contract-generator 0.1.2); platform `services/content` eaae493

## Question

A software platform has four network services. Each service's contract is
written today as a TypeScript interface. A platform generator reads that
interface and produces everything a service and its clients need, apart from
the native types themselves:

- an operation description used by a command line and web forms;
- a binding module that converts native values to and from the runtime's message
  values and composes contract-library adapters.

The owner has decided that service contracts should instead be declared once in
**horos**, a declarative type language over a shared value model, and generated
by **bittype**, a Go library that already generates native types from horos
declarations. In the owner's words (8 November 2026): "bittypes and horos were
made for exactly our use case, platform was its target from the beginning".

Native type generation from horos exists. Everything else is what this document
calls **binding generation**:

- conversion between native values and the runtime's message values;
- conversion between native values and JSON for presentation;
- operation data;
- runtime shape checks.

Binding generation exists today only for contracts sourced from TypeScript, in
the platform's own generator. Nothing produces it from a horos contract, and its
owner is undecided.

**How should binding generation be factored and owned?** The goal is that every
meaning has exactly one owner. Two kinds of change must stay additive:

- adding a contract library (such as a new callable or stream contract) changes
  only that library and the consumers that compose it;
- adding a target language changes only that language's implementation.

The answer decides:

- whether bittype's scope extends beyond types;
- the public contract between generated bindings and the platform's runtime
  libraries;
- how presentation and declared failures fit.

Section 7 holds the five questions for the expert.

## 1. System overview

### 1.1 The components

| Component | What it is | Language(s) |
| --- | --- | --- |
| **ontos** | The value model. A value is an **atom** (a byte string) or a **tuple** (an ordered sequence of values). On top of that floor it defines canonical data forms (integers, text, lists, maps, booleans, null) as labelled tuples, and a canonical binary encoding. | Go, Rust, TypeScript |
| **ontos/wire** | An optional extension of the value model: a **mixed value** may also contain live **Wire** leaves, that is, send capabilities (references to remote or local endpoints). A value with no Wire leaf is **ground**. | Go, TypeScript |
| **horos** | A type language over ontos values. (1) **Descriptors**: ontos values that describe sets of ground values, with a **conformance** checker. (2) A **declaration language** (revision 4): modules, imports pinned by content digest, generic type declarations, general structural kinds, and a primitive `wire<E>` type for send capabilities. | Descriptors and conformance: Go, Rust, TypeScript. Declaration semantics: **Go only** |
| **bittype** | A Go library: readable syntax and analysis for horos declarations, and **native type generation** for TypeScript, Go, Rust and C++ through a public generation interface. Its charter currently excludes "application-value validators, serialization, Wire adapters, RPC/streaming implementations and service migrations". | Go (generator); outputs TypeScript, Go, Rust, C++ source |
| **bitwire, bitruntime** | The message runtime. bitruntime **hydrates**: it converts mixed values to and from ground frames, exporting local endpoints and importing remote ones as proxies. It does not consult types. | Go, TypeScript |
| **platform `libraries/service`** | Four TypeScript packages: `service-api` (adapter contract, structural adapter factories, operation description types), `service-contracts` (the **Function** and **Stream** contract library: call protocol and adapters), `service-cli` (command line and web forms derived from an operation description), `contract-generator` (reads a TypeScript service interface and emits the description and the binding module). | TypeScript |
| **platform services** | Content, Cubes, Pointers and Drafts: 29 operations in total. Each has a model package (domain types and validators) and an API package. | TypeScript |
| **bitlink** | A dormant repository whose September design owned "model-to-wire bindings, adapter support and generation". It contains Markdown only and has had no commit since 26 September 2026. | none |

### 1.2 Terms used throughout

| Term | Meaning here |
| --- | --- |
| **Native value / native type** | A value or type of the target programming language, for example a TypeScript `Uint8Array` or interface. |
| **Hydrated value** | The runtime's in-memory message value: ground atoms and tuples plus live Wire leaves (TypeScript type `HydratedValue`). |
| **Wire** | A send capability. `Wire<E>` in horos describes a capability that accepts messages of type E. Sending is local admission only: it implies no reply or delivery. |
| **Descriptor, ground descriptor** | A horos type as an ontos value. A ground descriptor describes Wire-free values and can be checked by `conforms(value, descriptor)`. |
| **Declaration identity** | A horos declaration is identified by (module digest, name). The module digest is the SHA-256 of the module's canonical bytes; imports name it. |
| **Representation profile** | In bittype, a named policy that decides how horos forms become native types in one target language, for example "a closed record becomes a TypeScript interface, an optional field becomes `field?:`". |
| **Mapping library** | A bittype contribution that supplies native representations. A library selects horos declarations by **exact declaration identity**, never by name or shape. |
| **Adapter** | A runtime object converting one native type to and from hydrated values, in both directions. |
| **Library-owned adapter** | An adapter supplied by the library that declares the type, for example the Function library's adapter for `Function<A, B>`, or a domain library's adapter for a type whose native form is custom. |
| **Structural adapter** | An adapter for an ordinary structural form (record, list, union, optional field, recursion) composed from its parts' adapters. |
| **Binding** | The composition of adapters for a whole service contract, plus the call protocol, producing a client proxy and a provider dispatcher. |
| **Operation** | A field of a service record whose type is a library Function. Its **path** is its field path, for example `value.put`. |
| **Presentation** | The command line and web forms, which read and print values as JSON. |

### 1.3 Scale and direction

- **Consumers today.** Four TypeScript services with 6, 6, 7 and 10 operations. Each service's hand-written binding hooks are 81 to 144 lines. Its model packages are a few hundred lines.
- **Target languages.** TypeScript first. Go consumers are expected later.
- **Contract libraries.** Function and Stream exist. Cell and others are expected.
- **Expected change.** Other contract libraries and services will be added by different owners, possibly by different agents in parallel. That is why additive composition is a hard requirement rather than a preference.

## 2. Architectural context

### 2.1 The conversion layers

```text
 native value (e.g. a TypeScript object of type ContentLink)
      │   ▲
      │   │   adapters: generated structure + library-owned adapters   ← the subject of this document
      ▼   │
 hydrated value (atoms, tuples, live Wire leaves)
      │   ▲
      │   │   bitruntime hydration: exports/imports capabilities (existing; type-blind)
      ▼   │
 ground frame (atoms and tuples only)
      │   ▲
      │   │   ontos codec and a carrier such as WebSocket (existing)
      ▼   │
 bytes
```

**Presentation adds a parallel conversion,** native ↔ JSON, used by the command line and web forms.

### 2.2 Who may judge what (horos specification, quoted)

```text
| Judgment | Meaning and authority |
| --- | --- |
| Declaration syntax | Exact tags, arities, sorts, kind contracts and revision |
| Static semantics | Horos binding, primitive formation, kind membership, constructor contracts and materialization |
| Ground conformance | Existing WF/Conforms and diagnostics over ground Ontos data |
| Mixed structural admission | Ontos/wire finite positional structure and payload admission |
| Capability recognition | Runtime recognition of a particular sending face |
| Operational validity | Authority, scope, endpoint lifetime and transport state |
| Protocol compliance | Library-owned message, reply, state and cancellation laws |

None silently establishes the next judgment. Static Horos has no live
Wire/runtime dependency. Ontos/wire owns optional mixed values; runtime
composition owns recognition and hydration. No general package dispatches on
Function, Stream, Cell or a structural fingerprint to select behavior.
```

### 2.3 Standing architecture rules

The platform owner has made architecture part of correctness. These numbered rules apply to every design here:

| Rule | Requirement |
| --- | --- |
| AC-01 | Architecture is part of correctness. |
| AC-02 | Meanings remain with their semantic owners. Independently owned implementations supply their domain interpretations, policies and laws; shared mechanisms own their invariant guarantees. |
| AC-03 | Extensions compose additively. Within the agreed extension domain, adding a contribution changes only its owning implementation and consumer composition, not shared orchestration, its case catalogue or unrelated implementations. |
| AC-04 | Contributions use the same public mechanism. Builtins and the first implementation use the same public extension contract as independently supplied ones. |
| AC-05 | Composition preserves guarantees: identity, scope, ownership, authority, assumptions, evidence association, resource limits and failure meaning. |
| AC-06 | Limited coverage does not redefine the architecture. |
| AC-07 | Evidence supports only the claims it establishes. |
| AC-08 | General contract changes require a general justification. |
| AC-09 | Design interfaces from the domain, its owners and the guarantees consumers need. The existing API and the effort of changing it must never influence the design, not even as a tiebreaker. Reuse and migration are decided after the interface is chosen. |

### 2.4 Contract-owned bindings (platform decision, 29 October 2026, quoted)

> Function, Stream, Cell and future contracts are ordinary libraries above Wire.
> Their meaning, protocols and binding implementations belong to their libraries.
> None receives privileged treatment in the general generator, binding runtime,
> Wire, hydration or routing layers. Adding another contract must not require
> those layers to learn its name or semantics.
>
> Generated application bindings select and compose library-owned adapters.
> Dependencies point from those bindings to the selected contract libraries, and
> from the libraries to Wire/hydration. [...]
>
> The generator follows the consumer's imports, resolves declaration identities,
> binds generic arguments and emits ordinary static imports and factory calls.
> The association must survive aliases and nested generic applications; matching
> only a bare type-name string is insufficient. The runtime receives the composed
> adapters. It does not discover, load or register a contract by a runtime name.

Its responsibility table:

| Part | Owns | Does not own |
| --- | --- | --- |
| Wire, hydration and routing | Their existing value, native Wire, transport, export and resource laws | Function invocation, Stream demand/cancellation, Cell mutation or contract-specific type identity |
| General generation and binding machinery | Declaration/import resolution, generic slots and applications, uniform adapter composition and resource context | A catalogue of privileged library contracts, their protocols or semantic policies |
| Function library | Invocation/request-reply, its outcomes and callable lifetime, using supplied input/output adapters | The meaning or implementation of an arbitrary argument/result type |
| Stream library | Its declared ordering, demand, completion and cancellation behavior | Changes to the general runtime to recognize Stream |
| Generated consuming application | Static imports and composition of the particular bindings its API needs | Another Wire engine or a handwritten duplicate operation mapping |

The decision forbids eight things. Five are relevant here:

1. a contract-specific discriminator (`kind: 'function'`) in general machinery;
2. central switches, name checks, hardcoded imports, default registrations or built-in adapter tables, including a central registry that must be edited for each library;
4. a general-runtime traversal of library shapes that invents their structural identity;
5. moving special cases into another file, renaming them as a plugin, or generating the central switch;
6. making a new contract require edits to the general generator, runtime, Wire, hydration or their bootstrap configuration.

### 2.5 Earlier advice on this boundary

The same owners asked an earlier consultation (4 November 2026) about optional Wire types. Its final recommendation, which was applied:

> - Horos owns the optional Wire type and declaration semantics, not runtime recognition.
> - Bitruntime owns hydration, export identity, scope and operational admission.
> - Contract libraries own protocol definitions, behavioral laws, adapters and native recipes.
> - Bittype owns language-agnostic semantic tooling and explicit transformation orchestration, not another canonical type algebra.
>
> [...] Crucially, a native recipe does not implement a runtime adapter or prove protocol behavior.
> The proposed Adapter API should remain separate future work.

This consultation is that future work. An older consultation (28 September 2026), on a type language that horos has since replaced, recommended:

- **A JSON projection is a named projection of abstract values.** It is published independently of any network protocol, and a validation description is specific to that projection.
- **Mapping a declaration onto an existing native type needs conversion laws.** It "should use an explicit generation-side mapping and conversion functions, held to round-trip and validation laws", kept outside declaration identity.

## 3. Relevant code

### 3.1 horos declarations: the Function, Cell and Stream libraries

horos writes canonical modules as ontos values. `D(tag, ...)` is a labelled tuple, `T` a text atom, `I` an integer, `L` a list, `M` a map, `H(constructor, ...)` a primitive type expression and `Pin(module)` the module's digest. The specification's example libraries:

```text
WireLibrary = D("module", I(4), T("example.wire"), M{
  "k": Kinds
}, M{
  "Wire": Def(L(Param("M")), H("wire", P(0)))
})

FunctionLibrary = D("module", I(4), T("example.contracts.function"), M{
  "k": Kinds,
  "wire": D("import", T("example.wire"), Pin(WireLibrary))
}, M{
  "Unit": Def(L(), H("tuple")),
  "Function": Def(L(Param("A"), Param("B")),
    Use("wire", "Wire", H("tuple", P(0), Use("wire", "Wire", P(1)))))
})

CellLibrary = D("module", I(4), T("example.contracts.cell"), M{
  "k": Kinds,
  "fn": D("import", T("example.contracts.function"), Pin(FunctionLibrary))
}, M{
  "Cell": Def(L(Param("T")), H("record", A("closed"),
    H("field", T("get"), A("req"),
      Use("fn", "Function", Ref("fn", "Unit"), P(0))),
    H("field", T("set"), A("req"),
      Use("fn", "Function", P(0), Ref("fn", "Unit")))))
})
```

Read as types:

```text
Unit          = H("tuple")                        -- the empty tuple type
Function<A,B> = Wire<tuple(A, Wire<B>)>           -- send A with a reply capability for B
Cell<T>       = record closed { get: Function<Unit,T>, set: Function<T,Unit> }
Step<T>       = compound("item", T) | compound("end")
Stream<T>     = record closed { next: Function<Unit,Step<T>>, cancel: Function<Unit,Unit> }
```

The specification says the example Function "is illustrative and success-only. Its library owns fresh reply endpoints, matching, one-reply discipline, failures, completion and resource management." **A service contract is the same shape as Cell:** a closed record whose fields are applications of a library's Function type, with nested records for namespaces. No horos module for a Function with declared failures exists yet.

### 3.2 horos ground descriptors and the ontos data forms

The ground descriptor constructors, quoted from the horos specification, in abridged form:

```text
| any()                       | Every value |
| atom()                      | Any atom |
| atom-exact(Atom(bytes))     | Exactly those atom bytes |
| atom-len(int(lo),int(hi))   | Atom byte length between inclusive bounds |
| scalar(Atom(kind))          | Canonical `int`, `utf8-text` or `bool`, delegated to ontos/data |
| tuple(T0,...,Tn)            | Bare tuple with exactly this arity and matching children |
| seq(T)                      | Bare tuple of any arity, all children matching T |
| compound(Atom(label),T0,...)| Tuple beginning with that exact atom label, followed by exactly the matching children |
| list-of(T) / set-of(T)      | Canonical list / set whose elements match T |
| map-of(K,V)                 | Canonical map whose keys match K and values match V |
| record(Atom(mode),fields...)| A text-keyed canonical map; mode open or closed; each field required or optional |
| union(T0,...) / and(...) / not(T) / exact(value)
| rec(Atom(name),T) / ref(Atom(name))   guarded recursion
```

Text refinements are `text-len`, `key-grammar` (`^[a-z0-9][a-z0-9._-]{0,63}$`), `text-non-blank` and `lowercase-uuid`. There is no general pattern or prefix refinement.

The canonical ontos data forms (ontos data specification, quoted):

```text
| bytes   | a **bare** `Atom(bytes)` — opaque bytes *are* the floor's leaf |
| string  | `utf8-text(Atom(bytes))` — a **compound**                      |
| integer | `int(Atom(…))` — a **compound**                                 |

`Tuple(Atom("int"), Atom(<sign-magnitude bytes>))` — one leading sign byte followed
by the big-endian, minimal-length magnitude; zero is exactly the single byte 0x00.
`Tuple(Atom("utf8-text"), Atom(<utf-8 bytes>))` — the valid UTF-8 bytes verbatim.
```

The other forms:

- `bool` is `Tuple(Atom("bool"), Atom(0x00|0x01))`.
- `list` is `Tuple(Atom("list"), e…)`.
- `map` is `Tuple(Atom("map"), Tuple(k,v)…)`, with entries sorted by the codec bytes of each key.
- `null` is `Tuple(Atom("null"))`.

Canonical integers are unbounded.

### 3.3 horos's Go declarations API (build time)

```go
// horos core/go/declarations/model.go (signatures; bodies omitted)
func Resolve(root core.Value, dependencies []core.Value) (*Program, error)
func (p *Program) Declaration(id DeclarationID) (*Declaration, bool)
func (p *Program) Value(id DeclarationID) (TypeValue, error)       // checked value of a type declaration
func (p *Program) Proper(value core.Value) (TypeValue, error)      // admit a closed proper type expression
func (p *Program) Apply(callee TypeValue, arguments []TypeValue) (TypeValue, error)
func (p *Program) Descriptor(value TypeValue) (core.Value, error)  // materialize; may contain wire
func (p *Program) GroundDescriptor(value TypeValue) (core.Value, error) // only Wire-free (ground) types

type DeclarationID struct { Module Digest; Name string }
```

Facts that matter for binding generation:

- **`GroundDescriptor` fails for any type containing Wire.** A whole service record or a `Function<A,B>` fails. An operation's argument type A or result type B succeeds when it is Wire-free.
- **No API takes a checked type value apart.** To get A out of a field `Function<A,B>`, a consumer walks the checked expression tree: an application whose callee targets the Function declaration's identity, whose first argument is A. That tree is exposed by `Expression.Kind/Callee/Target/Arguments/Operands`. Doing this is ordinary semantic analysis; bittype already performs it for type generation.
- **A module has no slot for documentation or annotations.** The module is `D("module", I(4), T(name), imports, declarations)` and a declaration is `D("type", parameters, kind, expression)`, both with exactly checked arities. A record field is `Tuple(Atom("horos/field"), text(name), Atom("req"|"opt"), T)`. The specification says "Source spans and native naming choices are separate Bittype metadata". bittype's horos source syntax keeps comments only as trivia: they are not canonical data.

### 3.4 horos's TypeScript checker (run time)

```ts
// @bitspark/horos 0.9.0 (npm)
export function wf(t: Value): boolean;                    // descriptor well-formed
export function conforms(v: Value, t: Value): boolean;    // total, never throws
export type Diagnosis =
  | { verdict: "conforms" } | { verdict: "mismatch" }
  | { verdict: "not_well_formed" } | { verdict: "exhausted"; limit: Limit };
export function diagnose(v: Value, t: Value): Diagnosis;  // no mismatch location is reported
```

- **Values must be ground ontos values.** A hydrated tuple containing a Wire leaf never conforms, except to `any()`. Wire-free subtrees of a hydrated value are ordinary ground tuples, so they can be checked.
- **The Go checker, the TypeScript checker and the shared conformance corpus are unchanged since the TypeScript package was released.** Descriptors produced by the Go declarations layer should therefore be accepted by the TypeScript checker. No test exercises that combination yet.

### 3.5 ontos's JSON projection

ontos publishes one named JSON projection, `ontos-data-json/1`. It is **one-directional**, from domain JSON to ontos data:

- object becomes map;
- array becomes list;
- string becomes `utf8-text`;
- an integer literal becomes `int` (fractions and exponents are rejected);
- `true`/`false` become `bool`;
- `null` becomes `null`.

There is no reverse function, and no JSON form produces a set. The TypeScript package exports `project(raw: string, options): Value`.

### 3.6 bittype's generation interface

Mapping libraries **describe** native intent. A kernel **plans**: it walks the source's semantic graph itself and selects one rule per occurrence by exact selector. A target **builds** native syntax from the closed plan.

```go
// contracts/generation/go/mapping.go
// Selector is an exact source-owned identity, qualified by language/profile/marker.
// Its key is never interpreted by the composer.
type Selector struct { language representation.Language; key string }
// horos keys: "declaration/<module pin>/<name>" and "builtin/<len>:<constructor>/<len>:<discriminator>"

// contracts/generation/go/description.go
type Pattern[F, R any] struct {
	Selector Selector
	Slots    []Slot
	Capture  func(context.Context, R, UseView[F]) ([]Capture[F], error) // names children (e.g. type arguments)
}
type BindingDescription[F, T, R any] struct {      // a rule for USES of a selected declaration
	ID     string
	Input  Pattern[F, R]
	Output Expression[T]                           // native intent, with captured slots
	Expand func(context.Context, Expansion[F, T, R]) (Expression[T], error)
}
type DeclarationRecipe[F, T, R any] struct {       // a rule for PUBLISHING a declaration natively
	ID       string
	Selector Selector                              // zero selector = the default structural policy
	Expand   func(context.Context, DeclarationRequest[F, T, R]) (DeclarationDescription[T], error)
}
type MappingDescription[F, T, R any] struct {
	Profile      string                            // the representation profile it implements
	Bindings     []BindingDescription[F, T, R]
	Declarations []DeclarationRecipe[F, T, R]
	Helpers      []HelperRecipe[T]
}
type MappingLibrary[F, T, R any] interface {
	Describe(from, to representation.Language) (MappingDescription[F, T, R], error)
}
```

**Native intent is an immutable expression algebra.** Its nodes are target-owned forms with scalar data, plus captured slots, generic parameters and references:

> Native expressions contain target-owned forms, immutable scalar metadata and
> child expressions. The shared interpreter understands substitution and
> references, not the meanings of these forms.

**The kernel's selection is exact.** It refuses ambiguity and recursion, and it never falls back from an exact use mapping to a structural one:

```go
// planning.go (declaration planning)
choices := p.recipes[key]
if len(choices) == 0 {
	// An exact use mapping cannot silently become a structural publication.
	if len(p.bindings[key]) > 0 {
		return 0, fmt.Errorf("%w: exact use has no declaration recipe", ErrUnsupported)
	}
	choices = p.recipes[""]   // the default structural policy
}
if len(choices) != 1 { return 0, ErrAmbiguous }
```

**A target turns the closed plan into native syntax:**

```go
type Target[T any] interface {
	CheckExpression(context.Context, Expression[T]) error
	CheckDeclaration(context.Context, NativeDeclaration[T]) error
	CheckProgram(context.Context, Program[T]) error
	Build(context.Context, BuildInput[T]) (Assembly, error)  // closed program → native AST
}
```

The output is re-checked by the target language's own semantic checker before it is released.

**A library maps one of its declarations to an existing native type by exact identity.** This is how the Cell library's TypeScript mapping is written in bittype's tests:

```go
key, _ := source.Named(from, pin, "Cell")    // exact selector: module pin + name
slot := g.TypeSlot("V")
pattern := g.Pattern[h.Language, source.Reader]{Selector: key, Slots: []g.Slot{slot},
	Capture: func(ctx context.Context, src source.Reader, use g.UseView[h.Language]) ([]g.Capture[h.Language], error) {
		for _, a := range use.Actuals {                        // select the actual for binder T
			if p, err := src.Parameter(a.Binder); err != nil { return nil, err
			} else if p.Name == "T" { return []g.Capture[h.Language]{{Slot: slot, Use: a.Use}}, nil }
		}
		return nil, g.ErrScope
	}}
binding := g.BindingDescription[h.Language, ts.Language, source.Reader]{ID: "cells/typescript/1",
	Input: pattern, Output: model.External("@bitspark/cell", "Cell", g.Mapped[ts.Language](slot))}
```

A later stage places the files and spells `import type { Cell as Cell_1 } from "@bitspark/cell"`. Cell's body is never opened.

**What the targets can express today.** The native syntax vocabularies describe **types only**. TypeScript and Go value declarations, function bodies, statements and executable expressions are explicitly outside them, and the TypeScript file stage handles type-only imports, not runtime imports. bittype's own design note says:

> The existing native AST/text profiles cover type declarations, not executable
> adapter bodies. Generating such bodies requires separately scoped target syntax,
> checking, formatting and compiler qualification.

**The current representation profile** (one table, four targets):

| horos form | TypeScript | Go | Rust | C++ |
| --- | --- | --- | --- | --- |
| UTF-8 text | string | string | String | std::string |
| Boolean | boolean | bool | bool | bool |
| list-of(T) | T[] | []T | Vec<T> | std::vector<T> |
| Named closed record | interface | struct | struct | struct |
| Optional field | field?: T | BittypeOptional[T] | Option<T> | std::optional<T> |
| Proper type parameter | native parameter | native parameter | native parameter | template parameter |
| Ordinary named application | F<T> | F[T] | F<T> | F<T> |

- **Outside the profile today:** integers, bytes, unions, tuples, maps, recursion, open records and Wire.
- **Naming.** Each target's mapping library decides names: Go capitalises field names and wraps optional fields in a generated `BittypeOptional[T]{Present, Value}`; TypeScript keeps names.
- **No source-to-native name table is emitted.** Only syntax-node origins are kept.

### 3.7 bittype's own sketch of the binding layer

From bittype's design notes (proposals, not implemented):

```go
type Adapter[T any] interface {
    ToValue(ctx BindingContext, value T) (wire.HydratedValue, error)
    FromValue(ctx BindingContext, value wire.HydratedValue) (T, error)
}
```

```text
Type generation:
  type description + native representation -> Go declaration AST -> source
Future binding generation:
  resolved descriptions + selected library bindings -> adapter source
Runtime conversion:
  Adapter[T] + hydrated instance + binding context -> native T (or error)
  Adapter[T] + native T + binding context -> hydrated instance (or error)
```

> Future binding generation can emit static imports and compose the matching
> library adapter factories, preserving the selected contract and representation
> identities. Recipes for types and for runtime behavior must be qualified
> together where behavior is required; successful compilation alone is
> insufficient.

### 3.8 The platform's current binding machinery (TypeScript, generated from TypeScript)

All four services build on this machinery, which has been on their main branches since 8 November 2026. It already follows the 29 October decision for TypeScript-sourced contracts. Its parts:

- a structural runtime;
- a contract library;
- a presentation library;
- a generator.

**The adapter contract and the structural runtime** (`@bitspark/service-api` 0.3.2):

```ts
export interface AdapterContext {        // a transaction over the endpoints created during one conversion
  readonly owner?: ApiWireOwner;         // lifetime owner of exported endpoints (default limit 256)
  mark(): readonly [number, number];
  undo(handler: () => unknown): void;
  find(value: object, key: unknown): AdapterEntry | undefined;   // reuse an export of the same native object
  add(value: object, key: unknown, entry: AdapterEntry): AdapterEntry;
  rollback(mark?: readonly [number, number]): void;
  commit(): void;
}
export interface Adapter<T = unknown> {
  readonly requiresOwner?: boolean;      // the value may contain live capabilities
  readonly acceptsUndefined?: boolean;
  readonly arity?: readonly [number, number];   // for argument tuples: [required, allowed]
  toWire(value: T, context?: AdapterContext): HydratedValue;
  fromWire(value: HydratedValue, context?: AdapterContext): T;
}
export type WireSchema =
  | { kind: 'string' | 'number' | 'boolean' | 'null' | 'bytes' | 'void' }
  | { kind: 'literal'; value: string | number | boolean | null }
  | { kind: 'array'; element: WireSchema }
  | { kind: 'tuple'; elements: readonly { optional: boolean; schema: WireSchema }[] }
  | { kind: 'object'; properties: Record<string, { optional: boolean; schema: WireSchema }>; additional?: WireSchema }
  | { kind: 'union'; variants: readonly WireSchema[] }
  | { kind: 'ref' | 'codec'; name: string }            // codec: a semantic codec supplied by name at run time
  | { kind: 'binding'; adapter: Adapter<unknown> };    // a library-owned adapter object
export function dataAdapter<T>(schema: WireSchema, options?: { types?: Record<string, WireSchema>; codecs?: Record<string, WireCodec> }): Adapter<T>;
export function tupleAdapter(elements): Adapter;   export function arrayAdapter(element): Adapter;
export function recordAdapter(properties, additional?): Adapter;   export function unionAdapter(variants): Adapter;
export function lazyAdapter<T>(build: () => Adapter<T>): Adapter<T>;          // recursion
export function genericAdapter<T>(names: readonly string[], build): (scope) => Adapter<T>;  // cached by argument adapter objects
```

Every factory is a `dataAdapter` whose child positions are `{kind: 'binding', adapter}`, and `dataAdapter` is one generic interpreter of `WireSchema`. The interpreter therefore defines the hydrated encoding of every structural form:

| Native value | Hydrated encoding chosen by `dataAdapter` |
| --- | --- |
| string | an atom holding the UTF-8 bytes |
| number | an atom holding canonical decimal text (any finite number) |
| boolean | an atom `true` or `false` |
| null, void | an empty atom |
| bytes (`Uint8Array`) | the raw atom |
| array | a tuple of the elements |
| object (record) | a tuple of `(keyAtom, presence)` pairs sorted by key, where presence is `()` (absent) or `(value)` |
| union | `(indexAtom, body)`. Encoding tries the variants in order and the first that succeeds wins (rolling back any endpoints the failed attempts created); decoding uses the index. |
| binding | whatever the library adapter returns |

The limits per conversion are 100,000 nodes, depth 256 and 16 MiB.

```js
// the record case of dataAdapter, encoding direction (excerpt)
for (const name of Object.keys(source).sort()) {
  const key = text(name, state);
  const property = shape.properties[name], selected = property?.schema ?? shape.additional;
  if (!selected) fail('Unexpected object field: ' + name);
  if (source[name] === undefined && !acceptsUnit(selected)) {
    if (!property?.optional) fail('Undefined required field: ' + name);
    fields[name] = tuple([]);
  } else fields[name] = hydratedTuple([next(selected, source[name])]);
  entries.push(hydratedTuple([key, fields[name]]));
}
```

**The operation description** is static data, also used by presentation:

```ts
export type ApiSchema = /* the WireSchema kinds, except that a library binding is static data: */
  | { kind: 'binding'; factory: { module: string; export: string }; arguments: readonly ApiSchema[] }
  | { kind: 'parameter'; name: string } | { kind: 'apply'; name: string; arguments: Record<string, ApiSchema> } /* ... */;
export interface ApiOperation { path: readonly string[]; summary: string; parameters: readonly ApiParameter[]; result: ApiSchema }
export interface ApiParameter { name: string; optional: boolean; flatten: boolean /* presentation hint */; schema: ApiSchema }
export interface ApiDescription { name: string; summary: string; operations: readonly ApiOperation[]; types: Record<string, ApiSchema>; /* generics ... */ }
```

**The Function contract library** (`@bitspark/service-contracts` 0.1.2). Its native type carries the adapter association as a doc-comment annotation:

```ts
/** @wireAdapter {"module":"@bitspark/service-contracts/function-wire","export":"functionAdapter"} */
export interface Function<Args extends readonly unknown[], Result> {
  (...args: Args): Promise<Result>;
}
export function functionAdapter<Args extends readonly unknown[], Result>(
  input: Adapter<Args>, output: Adapter<Result>, options?): FunctionAdapter<Args, Result>;
export function serviceBinding<Api>(description: ApiDescription, operations: readonly OperationAdapters[], options?: BindingOptions): ServiceBinding<Api>;
export interface OperationAdapters { path: readonly string[]; input: Adapter<readonly unknown[]>; output: Adapter }
export interface BindingOptions {
  errors?: WireErrors;   // the service's error mapping (below)
  captureArguments?: (path: readonly string[], args: readonly unknown[]) => readonly unknown[];
  validateResult?: (path: readonly string[], args: readonly unknown[], result: unknown) => void;
}
export interface WireErrors {
  request(error: unknown): Value;      // provider: the request could not be decoded or captured
  encode(error: unknown): Value;       // provider: the implementation threw
  unknown(reason: string, message: string): Value;
  readAnswer<R>(answer: HydratedValue, decode: (value: HydratedValue) => R): R;   // caller
  transport(code: 'refused' | 'unknown', message: string, cause?: unknown): unknown;
}
```

**A callable value** travels as `(atom "bitspark/function/1", invokeWire, releaseWire)`. `functionAdapter.toWire` exports a fresh endpoint under the context's owner. It reuses the export when it sees the same native function twice, and returns a borrowed proxy unchanged.

**A call:**

- The request is `(atom profile, (pathAtom...), argumentsTuple, replyWire)`, with profile `bitspark/service/1` for services.
- The answer is `(atom "ok", value)` or an error answer.
- By default, a provider exception becomes `("unknown", message)` with no detail, and an undecodable request becomes `("protocol", message)`. `refused` is produced only locally, when the runtime refuses to admit a send.
- A service may replace this default with its own `WireErrors`. There is no declared-failure parameter on `Function`.

**The presentation library** (`@bitspark/service-cli` 0.6.2) derives command-line commands and web forms from an `ApiDescription`:

```ts
export interface ApiCodec {                // per semantic codec name, supplied by the service
  fromJson(value: unknown): unknown;       // synchronous conversion + semantic validation
  toJson(value: unknown): unknown;
  readonly kind?: 'text' | 'flag' | 'json' | 'bytes';   // how the field is entered: positional text, flag, JSON, bytes
}
export function deriveApiProjection<Api>(description: ApiDescription, options?: {
  codecs?: Record<string, ApiCodec>;
  error?: (error: unknown) => { family?: string; code: string; message: string; details?: unknown; status?: ProjectionStatus } | undefined;
}): ServiceProjection<Api>;
```

**The JSON conversion is a second hand-written interpreter, over `ApiSchema`.** Its rules differ from `dataAdapter`'s:

- unions are untagged in JSON, and the first variant that matches wins;
- nested bytes become `{"$bytes": base64}`;
- `void` becomes `null`.

**A library binding cannot be presented.** If any reachable schema contains a `binding`, the whole projection fails with "This library binding has no finite presentation adapter". Each operation's summary is printed in command-line help and web forms. A field's description is just its parameter name.

**The generator** (`@bitspark/contract-generator` 0.1.2) uses the TypeScript compiler. For each member of the service interface it produces an `ApiOperation`; the summary comes from the member's doc comment. A type position is resolved in this order:

1. a bound generic parameter;
2. a **semantic codec**, keyed `package#TypeName` from the nearest `package.json`. Template-literal types such as `` `sha256:${string}` `` must be codecs;
3. a **library binding**: a declaration carrying exactly one `@wireAdapter {"module": ..., "export": ...}` annotation, found by declaration identity through aliases;
4. a generic application;
5. structure.

An unannotated callable type is refused. It then emits a binding module like this real emitter output, abridged:

```js
// Generated from declaration-owned bindings. Do not edit.
import { dataAdapter as __data, tupleAdapter as __tuple, recordAdapter as __record, unionAdapter as __union,
         scopedAdapter as __scoped, adapterSlots as __slots } from '@bitspark/service-api/wire';
import { serviceBinding as __factory0 } from "@bitspark/service-contracts/service";
import { functionAdapter as __factory1 } from "@bitspark/service-contracts/function-wire";
const __description = { /* the ApiDescription as JSON */ };
export function createExampleBinding(__options = {}, __supplied = {}) {
  const __scope = __slots([], __supplied);
  const __types = Object.create(null);
  __types["T0"] = __scoped(__scope => __factory1(__tuple([{ optional: false, adapter: __data({"kind":"number"}) }]), __data({"kind":"string"})));
  __types["T3"] = __scoped(__scope => __union([__data({"kind":"null"}), __record({["text"]: { optional: false, adapter: __data({"kind":"string"}) }})]));
  return __factory0(__description, [
    { path: ["retain"], input: __tuple([{ optional: false, adapter: __types["T0"](__scope) }]), output: __data({"kind":"void"}) },
    { path: ["notes","get"], input: __tuple([{ optional: false, adapter: __data({"kind":"codec","name":"@example/model#NoteId"}, { codecs: __options.codecs ?? {} }) }]),
      output: __types["T3"](__scope) }], __options);
}
```

**Two things here are central to this question:**

- **Agreement by construction.** The generator reads the native TypeScript type and derives the adapter from that same type. The native representation and the adapter therefore agree by construction.
- **The runtime library decides the hydrated form.** `dataAdapter` chooses the encoding of a record or a union. Neither the type language nor the value model does.

### 3.9 The pilot service: Content

**The contract today** (TypeScript, six operations):

```ts
export interface ContentService {
  put(options: { bytes: Uint8Array }): Promise<{ hash: ContentHash }>;
  get(hash: ContentHash): Promise<Uint8Array>;
  list(): Promise<ContentHash[]>;
  check(): Promise<{ objects: number }>;
  readonly value: {
    put(options: { value: ContentValue }): Promise<ContentLink>;
    get(content: ContentLink): Promise<{ content: ContentLink; value: ContentValue }>;
  };
}
// the domain model package
export type ContentHash = `sha256:${string}`;          // validated: /^sha256:[0-9a-f]{64}$/
export type ContentRepresentation = 'raw' | 'map/1' | 'list/1';
export interface ContentLink { readonly representation: ContentRepresentation; readonly hash: ContentHash }
export type ContentValue = Uint8Array | ContentMap | readonly ContentValue[];   // recursive, untagged natively
export interface ContentMap { readonly [key: string]: ContentValue }
```

**Its hand-written Wire codecs.** They give each semantic type its own hydrated form:

```ts
contentHashToWire(hash)  = textAtom(hash)                          // UTF-8 atom "sha256:…"; decode re-validates
contentLinkToWire(link)  = tuple([textAtom(link.representation), textAtom(link.hash)])
// ContentValue: a TAGGED encoding of an UNTAGGED native union, with shared sub-values converted once (DAG sharing)
//   ("bytes", atom) | ("list", (value...)) | ("map", ((keyAtom, value)...))
// decoding enforces depth/steps/bytes limits, rejects cycles, unknown tags and duplicate keys,
// then deep-copies and freezes the result (snapshotContentValue)
```

**The binding hooks the service still writes by hand** (abridged):

```ts
const binding = createContentServiceBinding({          // the generated module (section 3.8)
  codecs: {                                            // Wire codecs for the three semantic types, keyed by string
    '@bitspark/content-model#ContentHash':  { toWire: v => contentHashToWire(contentHash(v)), fromWire: contentHashFromWire },
    '@bitspark/content-model#ContentLink':  { toWire: v => contentLinkToWire(contentLink(v)), fromWire: contentLinkFromWire },
    '@bitspark/content-model#ContentValue': { toWire: v => contentValueToWire(snapshotContentValue(v)), fromWire: contentValueFromWire },
  },
  captureArguments(path, args) {                       // runs on the caller before encoding, and on the provider after decoding
    switch (path.join('.')) {
      case 'put': {                                    // exactly one own data field `bytes`; copy the bytes before any await
        const bytes = field(args[0], 'bytes');
        if (!(bytes instanceof Uint8Array)) throw new ContentModelError('Expected bytes');
        return [{ bytes: new Uint8Array(bytes) }];
      }
      case 'value.put': return [{ value: snapshotContentValue(field(args[0], 'value')) }];  // deep capture, limits, no cycles
      default: return args;
    }
  },
  validateResult(path, _args, result) {                // `number` in TypeScript, but the contract means a non-negative integer
    if (path.join('.') === 'check') decimalAtom((result as { objects: number }).objects);
  },
  errors: { request: /* protocol or domain answer */, encode: errorToWire, unknown: /* … */, readAnswer, transport: /* … */ },
});
```

**Content's errors are declared only in code:**

- Domain error codes are `not-found | permission-denied | authentication-required | authorization-unavailable`.
- Model errors are `invalid-data | validation-limit`, with details `{limit, maximum}`.
- They travel as `("domain", family, code, message, ((key, ("text", t) | ("integer", digits))...))`, limited to 16 details, keys of 64 bytes and text of 256 bytes.
- The caller rebuilds the same error classes.
- Anything else becomes `unknown`, never "not run".

**JSON for presentation is a third hand-written conversion.** For `ContentValue` it is `{type: 'bytes', base64} | {type: 'map', entries} | {type: 'list', items}`, with its own limits. The codec table also says that `ContentHash` is entered as positional text.

**Hand-written code, in lines** (non-blank, non-comment):

| Category | Lines |
| --- | --- |
| Service interface and model types | 42 |
| Model validators and deep capture | 134 |
| Wire value codecs and their grammar helpers | 168 |
| Wire error mapping and error classes | 207 |
| JSON codecs | 79 |
| Codec table and JSON error projection | 17 |
| Binding hooks | 81 |

**The other three services need more.** Cubes, Pointers and Drafts (23 more operations) use:

- discriminated unions on a `kind` literal, including a variant with two literals (`'changed' | 'unchanged'`);
- `string | number`;
- intersections with unions;
- nullable references;
- integers typed as `number`;
- the value model's own `Value` type (arbitrary ground data);
- recursive JSON-like data;
- types imported from other services.

None uses streams or generics at its root.

### 3.10 The dormant bitlink design (prior art)

bitlink's September design for "model-to-wire binding and adapter generation" gave the type language and per-language type packages to bittype. It gave interfaces, bindings, clients and validators to itself. Its principles:

- "Declare once."
- "A target overrides, never declares."
- "The generator never writes a file a human edits."
- "No runtime dependency on the generator. A rendering depends on the runtime it binds to and on the standard library of its language; whatever it needs of the declaration is embedded."
- "Deterministic output."
- "Reference handling, native mapping and validation have separate jobs": the reference mechanism exports and imports bindings, generated proxies and bindings map arguments and results to native signatures, and validators check data against its schema.
- A candidate transparency law for generated bindings: `stub(bind(f)) ≈ f` and `bind(stub(a)) ≈ a`.

Its generator seam returned rendered files from strings, with no syntax tree. It was blocked behind the September type-language design, which horos has since replaced, and it never started.

## 4. Data and control flow

### 4.1 One call today

```text
caller                                                   provider
──────                                                   ────────
client.value.put({ value })
  check the argument count (the input tuple adapter's arity)
  captureArguments (service hook): deep capture
  input adapter .toWire → hydrated arguments
    (structural parts by dataAdapter, semantic codecs by
     name, library adapters; callables are exported
     into the call's transaction)
  request = (profile, path, arguments, replyWire)
  bitruntime: hydrate, export, send ground frame ───────▶ bitruntime delivers the request
                                                         check profile and path; decode arguments
                                                         captureArguments again; re-encode and
                                                         decode (canonical arguments)
                                                         invoke the implementation
                                                         output adapter .toWire; validateResult
  ◀──────────────────────────────── ("ok", value) or an error answer (service WireErrors hooks)
  readAnswer → output adapter .fromWire → result or a rebuilt error
```

### 4.2 Three representations of one declared type

For a type declared in horos, three representations exist, and each has a different natural owner:

| Representation | Who decides it under a horos contract | Who decides it today (TypeScript source) |
| --- | --- | --- |
| Native type and value | bittype's representation profile, and a library's exact recipe for its own declarations | the TypeScript interface itself |
| Hydrated (message) value | the horos descriptor: the value must conform to it, so it is canonical ontos data | `service-api`'s `dataAdapter` encoding, or a service codec |
| JSON (presentation) | undecided | `service-cli`'s interpreter, plus service JSON codecs |

**The hydrated forms differ concretely.** Take Content's `ContentLink`, declared in horos as `record closed { representation: union(exact "raw", exact "map/1", exact "list/1"), hash: ContentHash }`:

```text
horos-conforming (canonical ontos data):
  Tuple(Atom("map"),
        Tuple(utf8-text("hash"),           utf8-text("sha256:…")),
        Tuple(utf8-text("representation"), utf8-text("raw")))
  where utf8-text(s) = Tuple(Atom("utf8-text"), Atom(<bytes of s>)), entries in key-codec order
today's generic dataAdapter encoding (if ContentLink were structural):
  Tuple( Tuple(Atom("hash"), Tuple(Atom("sha256:…"))), Tuple(Atom("representation"), Tuple(Atom("raw"))) )
today's hand-written Content codec:
  Tuple(Atom("raw"), Atom("sha256:…"))
```

Three more differences:

- **Integers.** Today a count is an atom holding decimal text. In ontos it is `Tuple(Atom("int"), <sign-magnitude bytes>)`.
- **Unions.** Today a union is tagged by variant index. A horos union is untagged: a value conforms when it matches one alternative. Alternatives are told apart by content, such as compound labels or exact literals.
- **The wire format changes.** A binding generated from a horos contract therefore changes the wire format, and the hydrated form of each structural value is fixed by horos and ontos rather than by the binding runtime.

**Two inputs determine a structural adapter:** the native representation (from the profile) and the canonical data form (from horos and ontos). With a horos contract, they no longer come from one TypeScript source.

### 4.3 What a build from a horos contract must produce

**Inputs:**

- the service module;
- the model modules it imports (for example a `content-model` module declaring `ContentHash`, `ContentLink` and `ContentValue`);
- the contract library modules (Function).

All are pinned by digest.

**Products:**

1. **Native types.** These exist: bittype, with exact library mappings for `Function`, `ContentHash` and `ContentValue`.
2. **Per operation:** its path, its argument and result types, and its declared failure type.
3. **Per argument and result type, an adapter for each direction.**
   - Structural parts follow from the native representation and the canonical form.
   - Library parts (`Function`, `ContentHash`, `ContentValue`) come from the declaring libraries' factories, selected by declaration identity.
4. **A ground descriptor per Wire-free argument and result.** horos's Go `GroundDescriptor` produces them; at run time `conforms` checks received values against them.
5. **Presentation data:** parameter names and entry hints, summaries, and a JSON conversion.

### 4.4 Where Content's hand-written parts might go (a hypothesis, not a decision)

| Hand-written today | Possible home under a horos contract |
| --- | --- |
| The TypeScript service interface | the horos service module; TypeScript types generated |
| Codec table keyed by `package#Type` strings | library adapter factories selected by declaration identity |
| `ContentLink` codec | structural: no hand-written code |
| `ContentHash` codec and `/^sha256:[0-9a-f]{64}$/` | declared as text in a model module; the model library supplies semantic validation (horos has no pattern refinement) |
| `ContentValue` codec (tagged on the wire, untagged natively, DAG sharing, limits) | a library-owned custom representation and adapter |
| `captureArguments` (copy bytes, deep-snapshot values before any await) | perhaps inherent if conversion happens synchronously when the call starts; today it duplicates conversion |
| `validateResult` (`check` returns an integer) | gone: an integer type in the declaration plus conformance |
| Error mapping (207 lines) | declared failure types plus the Function library's protocol |
| JSON codecs | undecided (section 6) |

## 5. Constraints

- **The architecture rules** (section 2.3) and the contract-owned binding decision (section 2.4) are fixed.
- **horos.**
  - It owns type semantics and has no generator and no runtime binding.
  - Its declaration semantics exist only in Go. Its ground conformance exists in Go, Rust and TypeScript.
  - Its module format has no documentation slot. Adding one is horos's decision, not this consultation's; this consultation may recommend it.
- **bittype.**
  - It is a Go library. It must not depend on the Wire runtime packages (bitwire, bitruntime).
  - Its generated type output is re-checked by each target language's checker.
- **No central catalogue.** No component may hold a list of contract libraries, and no component may select behavior by type name or shape.
- **Consumers.** The services and their clients are TypeScript, in Node and browsers. Go consumers are expected later. Bindings must eventually exist per target language.
- **Not prerequisites.** No release checksum or validation workflow is a prerequisite (owner rule, 4 November 2026), but ordinary tests and type checks are. Security hardening is deferred during development.
- **The wire format.**
  - The platform services' wire format may change: they share one installation, and no external client depends on the current byte layout.
  - Values described by a horos contract are expected to conform to that contract on the wire. Today's structural encoding does not (section 4.2).

## 6. What we have considered

### 6.1 Where binding generation could live

**A. A second product of bittype's generator.**

- **Shape.** The same kernel, representation profiles and exact selection produce the adapters as well as the types. A library mapping supplies, beside its native type recipe, a reference to its adapter factory for that target.
- **For:**
  - One owner decides both the native representation and its conversion.
  - The kernel already resolves imports and exact identities, binds generic arguments, and records external symbols for import placement.
  - A Go generator can later emit Go bindings through the same mechanism.
- **Against:**
  - bittype's target vocabularies describe types only. Executable adapter code needs value and statement syntax per target language, with checking and formatting: a large, new and unqualified surface. That cost disappears if the generated product is data (option C).
  - bittype's charter would change.

**B. A separate binding generator consuming bittype's analysis.**

- **Shape.** The concrete candidate is the platform's existing contract-generator, with its input switched from the TypeScript compiler to bittype's analysis of the horos contract.
- **For:** The generator stays next to the runtime libraries it targets, and it already composes library factories by declaration identity.
- **Against:**
  - It must reproduce every representation choice bittype made for the types: native names, Go's capitalised fields, the optional encoding, union representations, external symbols. bittype publishes none of these today ("It emits no generation manifest").
  - Unless such a manifest becomes the authoritative contract, one representation decision has two implementations.
  - Go bindings would need another generator.

**C. Generated data plus a generic runtime interpreter per language.**

- **Shape.** The build emits, per declaration, a description: the canonical form, the native representation choices, and static imports of library adapter factories. A hand-written runtime library per target language interprets it. Today's platform machinery mostly works this way: `dataAdapter` interprets a schema, and only the composition is generated code.
- **Two owners to choose.** Who emits the description: bittype, from the same decisions it uses for types? Who owns each language's interpreter: a runtime library?
- **For:**
  - The generator needs no executable target syntax.
  - Representation choices travel as data from their owner.
- **Against:** The interpreter is an executable rendering of the representation profile, so it is a second place where profile semantics live, unless the profile itself is defined as data both sides read.

**D. Runtime interpretation of horos modules (rejected for now).**

- **Why rejected:** It needs a second implementation of horos declaration semantics in TypeScript, and a run-time table from declaration identity to factory, which section 2.4 forbids or comes close to forbidding.

### 6.2 Narrower decisions, with our present leaning (none decided)

1. **Hydrated form: canonical ontos data that conforms to the horos descriptor.**
   - No encoding specific to the binding runtime. Leaning yes, since the contract is the horos declaration.
   - It changes today's wire format.
2. **Runtime shape checks: horos `conforms` against ground descriptors emitted at build time.**
   - Received Wire-free arguments and results are checked before conversion, so adapters only need semantic validation such as the hash pattern. Leaning yes.
   - Limits: the checker reports no mismatch location, and a Wire-containing value can only be checked in its Wire-free parts.
3. **Library contributions.** Per declaration and target, a library supplies a native type recipe (bittype mapping libraries already do this) and an adapter factory reference parameterised by its argument adapters (today's `@wireAdapter`).
   - Undecided: whether both live in one per-library contribution (a bittype mapping library in Go), or in metadata beside the library's horos module.
4. **Declared failures.**
   - `Function<A, B, E>`, whose reply is `Outcome<B, E> = compound("ok", B) | compound("error", E)`.
   - Protocol failures (refused, unknown, protocol, wire limit) stay the Function library's own outcomes.
   - Leaning yes. The Function library's owner must agree.
5. **Documentation.** horos modules have no documentation slot. Three options:
   - horos adds one, which changes canonical bytes and so declaration identity;
   - bittype's source syntax turns doc comments into metadata emitted beside the canonical module;
   - a separate documentation module keyed by declaration identity.
6. **JSON for presentation.** Three options:
   - derive it from the canonical form through a total named projection; ontos's `ontos-data-json/1` is one-directional and has no bytes or sets;
   - derive it from the native representation generically, as `service-cli` does now;
   - library-owned JSON adapters for custom representations, with structural ones derived.

   The two existing interpreters (`dataAdapter` and `service-cli`) already disagree on unions, `void` and arity.
7. **Capture before asynchronous work.**
   - If conversion to the hydrated form happens synchronously when the call starts, it is the capture, and a separate hook is unnecessary.
   - Leaning yes for data. Callables and streams are converted by their libraries.

### 6.3 The argument we want tested

**Our working argument for A, or for C with bittype emitting the description:** one representation decision must have one owner. Today the platform achieves this by deriving types and adapters from the same TypeScript source. With a horos contract, the native types come from bittype. A generator elsewhere would have to re-derive bittype's choices and keep them in step as profiles and targets grow.

**The counter-argument:**

- bittype is a type tool, and executable target syntax is a large new surface.
- The platform's runtime libraries define the adapter contract and its resource semantics, so the generator that composes them arguably belongs beside them.
- A representation manifest emitted by bittype could make B sound.

We have not decided, and AC-09 forbids letting the existing platform generator or the cost of moving it influence the answer.

## 7. Questions for the expert

1. **Factoring and ownership.**
   - How would you factor binding generation relative to the two existing owners: the representation profile that decides native types, and the contract libraries that own protocol adapters?
   - The aim: one owner for each representation decision, and additive change when a contract library or a target language is added.
   - Between generating adapter code (options A and B) and generating descriptions for a per-language runtime interpreter (option C), which principles decide?
   - What do comparable systems teach about where this layer lives and how it stays consistent with the generated types? For example: protobuf and gRPC code-generation plugins, Smithy, Cap'n Proto, TypeSpec emitters, serde derive, Avro specific versus generic records.
2. **The library contribution contract and validation.**
   - How should a contract or domain library declare, for one of its declarations and one target language, its native type and adapter factory, so that composition selects them by exact identity, with generic arguments passed as adapters? Examples: the Function library for `Function<A, B, E>`, or the Content model for `ContentHash` and the recursive, natively untagged `ContentValue`.
   - What laws should such a contribution satisfy and be tested against, such as round-trip, conformance of produced values, and capability ownership?
   - How would you divide validation between shape checks of received values (horos conformance against build-emitted descriptors) and the adapters themselves?
3. **Several projections of one declaration.** One declaration needs native↔hydrated conversion, native↔JSON conversion for a command line and web forms, and operation metadata (paths, entry hints, summaries).
   - How would you derive these without each library writing one adapter per projection?
   - Should JSON be a projection of the canonical value form, of the native form, or a family each library owns?
   - Where should documentation live when the type language's canonical form carries none?
   - Is fixing the hydrated form to the canonical, descriptor-conforming data (section 4.2) the right principle, given that it changes the wire format?
4. **Failures and a first pilot.**
   - How should declared operation failures be modelled relative to protocol and transport failures, so that the call library owns the protocol and each service owns only its domain errors?
   - What should a first pilot qualify so that it proves the architecture rather than one example? The pilot is Content: six operations, a recursive union, bytes and two custom representations. Which tests or falsification criteria would show that adding a second contract library, a second service, or Go as a second target language changes only its own owner?
5. **What are we not asking?** What would you want to know, or change, before committing to any of these designs? Which risk in this layering do you consider most likely to be underestimated?
