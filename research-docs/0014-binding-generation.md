# Research: Binding generation for service contracts declared in a type language

**ID:** 0014
**Date:** 9 November 2026
**Status:** advised
**Run-ID:** run_250b7afe-dcac-4972-a8b2-4d3962226188
**Reviewed:** https://github.com/Bitspark/bittype/issues/194#issuecomment-6728826766
**Owner:** claude-bittype-gaps
**Issue:** https://github.com/Bitspark/bittype/issues/194
**Code baseline:** bittype 6132170e394f2261d6a24c8efd7839c93299fdcc; horos 3a51a06e6db1800a39f94782213609a2d8723d10 (Go core 0.10.2; published npm package 0.9.0); ontos 0.16.0; bitwire 0.7.0; bitruntime 0.8.0; platform repository `libraries/service` 24fedd4 (service-api 0.3.2, service-contracts 0.1.2, service-cli 0.6.2, contract-generator 0.1.2); platform repository `services/content` eaae493

## Question

A software platform has four network services. Today each service's contract is
written as a TypeScript interface. A generator in the platform repository reads
that interface and emits two things:

- an operation description, which a command line and web forms interpret at
  run time;
- a binding module, which converts native values to and from the message
  runtime's values and composes the adapters owned by contract libraries.

The owner has decided that service contracts should instead be declared once in
**horos**, a declarative type language over a shared value model, and turned
into code by **bittype**, a Go library that already generates native types from
horos declarations. In the owner's words (8 November 2026), on using horos and
bittype for the platform's contracts: "bittypes and horos were made for exactly
our use case, platform was its target from the beginning". The owner did not
decide who generates the bindings.

From a horos contract, a build must produce more than native types. We call the
rest **binding generation**:

1. conversion between native values and message values (exists today only for
   contracts sourced from TypeScript, in the platform repository's generator);
2. conversion between native values and JSON for the command line and web forms
   (exists today only as a hand-written run-time interpreter, plus hand-written
   per-type JSON converters);
3. operation data: paths, parameters, summaries;
4. shape checks of received values (do not exist today);
5. declared operation failures (today the error mapping is hand-written per
   service; the call protocol has no declared failure type).

Nothing produces these from a horos contract, and their owner is undecided.

**How should binding generation be factored and owned?** The goal is that every
meaning has exactly one owner. Three kinds of change must stay additive:

- adding a contract library, such as a new callable or stream contract, changes
  only that library and the consumers that compose it;
- adding a service changes only that service;
- adding a target language changes only that language's implementation.

**The tensions we see:**

- **The representation decision.** The native representation of a type is
  decided by bittype's representation policy. An adapter must agree with it
  exactly. Today both come from one TypeScript source; with horos they would not.
- **The message form.** Today the platform's runtime library decides how a record
  or union looks as a message value. A horos contract would let the horos type
  and the shared value model decide that instead. Whether it should is open.
- **Types-only target syntax.** bittype's native syntax covers type
  declarations only, not executable code.
- **Library contributions.** We must decide where a contract library declares
  its native type and its adapter, per target language.
- **Presentation and documentation.** JSON conversion has no agreed source, and
  the type language has no place for documentation.

Sections 1 to 3 give the context. Section 4 is the crux: what flows where and
which representations exist. Section 5 lists what is fixed. Section 6 lists the
options and our present leanings, none decided. Section 7 holds the five
questions.

## 1. System overview

### 1.1 The components

| Component | What it is | Language(s) |
| --- | --- | --- |
| **ontos** | The value model. A value is an **atom** (a byte string) or a **tuple** (an ordered sequence of values). On that floor it defines canonical data forms (integers, text, lists, maps, booleans, null) as labelled tuples, and one canonical binary encoding. | Go, Rust, TypeScript |
| **ontos/wire** | An optional extension of the value model: a **mixed value** may also contain live **capability leaves**, references that can send messages to an endpoint. A value without such leaves is **ground**. | Go, TypeScript |
| **horos** | A type language whose canonical form is itself ontos values. It has no text syntax of its own; bittype supplies a readable one. Its two parts are listed below this table. | Descriptors and conformance: Go, Rust, TypeScript. Declaration semantics: **Go only** |
| **bittype** | A Go library: readable syntax and analysis for horos declarations, and **native type generation** for TypeScript, Go, Rust and C++ through a public generation interface. Its charter excludes "application-value validators, serialization, Wire adapters, RPC/streaming implementations and service migrations". | Go; emits TypeScript, Go, Rust, C++ |
| **bitwire** | The message model's interfaces: the hydrated value type and the send capability interface. | Go, TypeScript |
| **bitruntime** | The message runtime. It **hydrates**: it converts mixed values to and from ground values for transport, exporting local endpoints and importing remote ones as proxies. It never consults types. | Go, TypeScript |
| **platform repository, `libraries/service`** | Four TypeScript packages (detail in section 3.8). | TypeScript |
| **platform services** | Content, Cubes, Pointers and Drafts: 29 operations in total. Each has a model package (domain types and validators) and an API package. | TypeScript |
| **bitlink** | A dormant repository whose September design owned "model-to-wire bindings, adapter support and generation". It contains Markdown only and has had no commit since 26 September 2026. | none |

horos's two parts:

- **Descriptors.** Ontos values that describe sets of values, with a
  **conformance** checker for *ground* values.
- **A declaration language.** Revision 4 has modules, imports pinned by content
  digest, generic type declarations, **kinds** (type-level predicates over types,
  for example "this type is data"), and a primitive type `wire<E>` for capability
  leaves that accept messages of type E. A descriptor may mention `wire`. Only a
  **ground descriptor**, one with no `wire`, can be checked against values.

The four packages in `libraries/service`:

- `service-api`: the adapter contract, structural adapters, and the operation
  description types;
- `service-contracts`: the **Function** and **Stream** contract library, with
  its call protocol and adapters;
- `service-cli`: the command line and web forms derived from an operation
  description;
- `contract-generator`: reads a TypeScript service interface and emits the
  description and the binding module.

### 1.2 Terms used throughout

| Term | Meaning in this document |
| --- | --- |
| **Native value / native type** | A value or type of the target programming language, for example a TypeScript `Uint8Array` or interface. |
| **Capability (leaf), `wire<E>`** | A live reference that can send messages of type E. Sending means only that the local runtime **admitted** the message (accepted it for delivery); it implies no reply or delivery. In horos the primitive is written `wire<E>`; libraries wrap it as `Wire<E>`. |
| **Hydrated value, message form** | The runtime's in-memory message value: ground atoms and tuples plus capability leaves (TypeScript type `HydratedValue`). The platform's TypeScript code calls conversion to this form `toWire`; its `WireSchema`, `WireErrors` and "Wire codecs" also refer to this form. **"Wire format"** in this document means the bytes on the network, which follow from the hydrated form. |
| **Descriptor, ground descriptor, conformance** | A horos type as an ontos value; a ground descriptor has no capability parts; `conforms(value, descriptor)` checks a ground value. |
| **Declaration identity** | A horos declaration is identified by (module digest, name). The module digest is the SHA-256 of the module's canonical bytes, and imports name it. |
| **Representation profile** | In bittype: a named policy deciding how horos forms become native types in one target language, for example "a closed record becomes a TypeScript interface; an optional field becomes `field?:`". Not to be confused with a **protocol tag** such as `bitspark/service/1`, which names a message protocol. |
| **Mapping library, recipe** | A bittype contribution supplying native representations. A **recipe** is its rule for one horos declaration. A library selects declarations by **exact declaration identity**, never by name or shape. |
| **Contract library** | A library defining a protocol type such as `Function<A, B>`, together with its behaviour and adapters. Function, Stream and Cell are contract libraries. |
| **Adapter** | A run-time object converting one native type to and from hydrated values, in both directions. |
| **Library-owned adapter** | An adapter supplied by the library that declares the type: the Function library's adapter for `Function<A, B>`, or a domain library's adapter for a type whose native form is custom. |
| **Structural adapter** | An adapter for an ordinary structural form (record, list, union, optional field, recursion), composed from its parts' adapters. |
| **Semantic codec** | In today's platform code, a converter for a named domain type (such as a hash), supplied at run time under a string key `package#TypeName`. Distinct from the **canonical encoding** (ontos's bytes). |
| **Binding** | The composition of adapters for a whole service, plus the call protocol, giving a client proxy and a provider dispatcher. Two unrelated uses also occur in quoted material: horos "binding" (name resolution), and bittype's `BindingDescription` (a rule for uses of a declaration). |
| **Operation** | A field of a service record whose type is a library Function. Its **path** is its field path, for example `value.put`. |
| **Presentation** | The command line and web forms, which read and print values as JSON. |
| **Semantic owner vs. lifetime owner** | "Owner" usually means the component that owns a meaning (rule AC-02). The platform runtime also has a **lifetime owner** (`ApiWireOwner`) that keeps exported endpoints alive and closes them. |
| **Qualified** | Validated by tests and checks for a stated scope. |

### 1.3 Scale, people and direction

- **Consumers.** Four TypeScript services with 6, 6, 7 and 10 operations,
  running in Node.js and browsers. Each service's hand-written binding hooks are
  81 to 144 lines.
- **People.** One human owner decides for all the repositories involved:
  ontos, horos, bittype, the runtime and the platform. Each repository is
  maintained by AI coding-agent sessions, several of which may work in parallel
  on different repositories. That is why additive composition with clear
  ownership is a hard requirement rather than a preference.
- **Target languages for bindings.** TypeScript now; Go is expected (no date).
  bittype also emits Rust and C++ types; bindings for those are not planned, but
  the design must not preclude them.
- **Expected growth.** More services, and more contract libraries (Cell, and
  others), written by different sessions.
- **What may change.**
  - The native TypeScript API, the hydrated form and the wire format may all
    change: the services share one installation, and no external client
    depends on today's bytes.
  - Horos changes, for example a documentation slot or a text pattern
    refinement, may be recommended; horos's maintainers decide.

## 2. Architectural context

### 2.1 The conversion layers

```text
 native value (e.g. a TypeScript object of type ContentLink)
      │   ▲
      │   │   adapters: structure + library-owned adapters        ← the subject of this document
      ▼   │
 hydrated value (atoms, tuples, capability leaves)
      │   ▲
      │   │   bitruntime: exports/imports capabilities (existing; never consults types)
      ▼   │
 ground value for transport (atoms and tuples only)
      │   ▲
      │   │   ontos canonical encoding and a carrier such as WebSocket (existing)
      ▼   │
 bytes (the "wire format")
```

**Presentation adds a parallel conversion,** native ↔ JSON, for the command
line and web forms.

### 2.2 Who may judge what

The horos specification separates the following judgments, and **no judgment
implies the next one**:

- **horos judges static validity of declarations, and conformance of ground
  values to ground descriptors.**
- **The ontos/wire extension judges the shape of mixed values.**
- **The runtime judges capabilities and resources:**
  - it recognises a capability;
  - it enforces authority, endpoint lifetime and transport state.
- **Each contract library judges its protocol:** messages, replies, state and
  cancellation.

The specification adds: "No general package dispatches on Function, Stream,
Cell or a structural fingerprint to select behavior."

### 2.3 Standing architecture rules (AC = architectural correctness)

The owner has made architecture part of correctness. The numbered rules apply to
every design here:

| Rule | Requirement |
| --- | --- |
| AC-01 | Architecture is part of correctness: ownership, dependency and composition contracts are requirements, like behaviour. |
| AC-02 | Meanings remain with their semantic owners. Independently owned implementations supply their domain interpretations, policies and laws; shared mechanisms own their invariant guarantees. |
| AC-03 | Extensions compose additively. Within the agreed extension domain, adding a contribution changes only that contribution and the consumers that compose it. It does not change shared orchestration, a central list of cases, or unrelated implementations. |
| AC-04 | Built-in and first contributions use the same public extension mechanism as independent ones. |
| AC-05 | Composition preserves guarantees: identity, scope, ownership, authority, assumptions, resource limits and failure meaning. |
| AC-06 | A small first coverage does not redefine the intended extension boundary. |
| AC-07 | A test or a compiler's acceptance supports only the claim it establishes. |
| AC-08 | A change to a shared contract needs a general justification. |
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

The decision's responsibility table, quoted:

| Part | Owns | Does not own |
| --- | --- | --- |
| Wire, hydration and routing | Their existing value, native Wire, transport, export and resource laws | Function invocation, Stream demand/cancellation, Cell mutation or contract-specific type identity |
| General generation and binding machinery | Declaration/import resolution, generic slots and applications, uniform adapter composition and resource context | A catalogue of privileged library contracts, their protocols or semantic policies |
| Function library | Invocation/request-reply, its outcomes and callable lifetime, using supplied input/output adapters | The meaning or implementation of an arbitrary argument/result type |
| Generated consuming application | Static imports and composition of the particular bindings its API needs | Another Wire engine or a handwritten duplicate operation mapping |

"Generic slots" means type-parameter positions filled with argument adapters.
"Resource context" means the facilities for the lifetime of exported endpoints.

**The decision forbids eight things.** Five are relevant here, numbered as in
the source; items 3, 7 and 8 are omitted as unrelated.

1. a contract-specific discriminator (`kind: 'function'`) in general machinery;
2. central switches, name checks, hardcoded imports, default registrations or
   built-in adapter tables, including a central registry that must be edited
   for each library;
4. a general-runtime traversal of library shapes that invents their structural
   identity;
5. moving special cases into another file, renaming them as a plugin, or
   generating the central switch;
6. making a new contract require edits to the general generator, runtime, Wire,
   hydration or their bootstrap configuration.

### 2.5 Earlier advice on this boundary

An earlier consultation (4 November 2026), on optional capability types,
concluded as follows; its conclusions were applied:

> - Horos owns the optional Wire type and declaration semantics, not runtime recognition.
> - Bitruntime owns hydration, export identity, scope and operational admission.
> - Contract libraries own protocol definitions, behavioral laws, adapters and native recipes.
> - Bittype owns language-agnostic semantic tooling and explicit transformation orchestration, not another canonical type algebra.
>
> [...] Crucially, a native recipe does not implement a runtime adapter or prove protocol behavior.
> The proposed Adapter API should remain separate future work.

This consultation is that future work.

An older consultation (28 September 2026) concerned a type language that horos
has since replaced. Two of its recommendations still bear on this question:

- JSON should be a *named projection* of abstract values, published
  independently of any network protocol.
- Mapping a declaration onto an existing native type "should use an explicit
  generation-side mapping and conversion functions, held to round-trip and
  validation laws", kept outside declaration identity.

## 3. Relevant code

### 3.1 horos declarations: a service is a record of Functions

horos writes modules as ontos values. The specification's notation:

```text
A(s)           = Atom(UTF8(s))                         -- a raw atom
T(s)           = ontos utf8-text(s)                    -- a text value
I(n)           = ontos int(n)
L(x, ...)      = ontos list
M{ k: v, ... } = ontos map with text keys T(k)
H(c, x, ...)   = Tuple(A("horos/" + c), x, ...)        -- a primitive type expression
D(c, x, ...)   = Tuple(A("horos/decl/" + c), x, ...)   -- a declaration form
Pin(module)    = Atom(SHA256(canonical encoding of module))
Module         = D("module", I(4), T(name), M{ imports }, M{ declarations })
```

The specification's example libraries, with its helper abbreviations:

```text
Kinds                 = D("import", T("horos.kinds"), Pin(Prelude))   -- the standard kinds module
KType                 = D("kind-imported", T("k"), T("Type"))         -- the kind "any proper type"
P(i)                  = D("param", I(i))                              -- reference to type parameter i
Param(name)           = D("parameter", T(name), KType)
Def(parameters, body) = D("type", parameters, KType, body)            -- a declaration: parameters, result kind, body
Ref(alias, name)      = D("imported", T(alias), T(name))              -- a non-generic imported declaration
Use(alias, name, ...) = D("apply", Ref(alias, name), L(...))          -- an applied imported generic

WireLibrary = D("module", I(4), T("example.wire"), M{ "k": Kinds }, M{
  "Wire": Def(L(Param("M")), H("wire", P(0)))
})
FunctionLibrary = D("module", I(4), T("example.contracts.function"),
  M{ "k": Kinds, "wire": D("import", T("example.wire"), Pin(WireLibrary)) }, M{
  "Unit": Def(L(), H("tuple")),
  "Function": Def(L(Param("A"), Param("B")),
    Use("wire", "Wire", H("tuple", P(0), Use("wire", "Wire", P(1)))))
})
CellLibrary = D("module", I(4), T("example.contracts.cell"),
  M{ "k": Kinds, "fn": D("import", T("example.contracts.function"), Pin(FunctionLibrary)) }, M{
  "Cell": Def(L(Param("T")), H("record", A("closed"),
    H("field", T("get"), A("req"), Use("fn", "Function", Ref("fn", "Unit"), P(0))),
    H("field", T("set"), A("req"), Use("fn", "Function", P(0), Ref("fn", "Unit")))))
})
```

Read as types:

```text
Unit          = the empty tuple type
Function<A,B> = Wire<tuple(A, Wire<B>)>     -- send A together with a capability for the reply B
Cell<T>       = record closed { get: Function<Unit,T>, set: Function<T,Unit> }
Stream<T>     = record closed { next: Function<Unit, Step<T>>, cancel: Function<Unit,Unit> }   (module omitted)
```

**What the specification says about Function.** It "is illustrative and
success-only. Its library owns fresh reply endpoints, matching, one-reply
discipline, failures, completion and resource management."

**A service contract has the same shape as Cell:** a closed record whose fields
are applications of a library's Function type, with nested records for
namespaces such as `value.put`. No horos module for a Function with declared
failures exists yet.

### 3.2 horos ground descriptors and the ontos data forms

Ground descriptor constructors, abridged from the horos specification:

```text
any()  atom()  atom-exact(bytes)  atom-len(lo,hi)    -- atoms (byte strings)
scalar(int | utf8-text | bool)                       -- canonical ontos scalars
tuple(T0..Tn)  seq(T)  compound(label, T0..Tn)       -- fixed tuple, any-length tuple, labelled tuple
list-of(T)  set-of(T)  map-of(K,V)                   -- canonical collections
record(open|closed, field(name, req|opt, T)...)      -- a text-keyed canonical map
union(T..)  and(T..)  not(T)  exact(value)           -- untagged union, intersection, complement, singleton
rec(name, T)  ref(name)                              -- guarded recursion
```

Text refinements are `text-len`, `key-grammar` (`^[a-z0-9][a-z0-9._-]{0,63}$`),
`text-non-blank` and `lowercase-uuid`. There is no general pattern or prefix
refinement. **A horos union is untagged:** a value conforms when it conforms to
at least one alternative.

The canonical ontos data forms (ontos data specification):

```text
bytes   : a bare Atom(bytes)
integer : Tuple(Atom("int"), Atom(<one sign byte, 0x00 non-negative or 0x01 negative,
                                   then the big-endian minimal-length magnitude>))   -- unbounded
text    : Tuple(Atom("utf8-text"), Atom(<UTF-8 bytes, no normalization>))
bool    : Tuple(Atom("bool"), Atom(0x00 | 0x01))
list    : Tuple(Atom("list"), e...)
map     : Tuple(Atom("map"), Tuple(k, v)...)   -- entries sorted by the canonical encoded bytes of each key
null    : Tuple(Atom("null"))
```

### 3.3 horos's Go declarations API (build time)

```go
// horos core/go/declarations (signatures; bodies omitted)
func Resolve(root core.Value, dependencies []core.Value) (*Program, error)
func (p *Program) Declaration(id DeclarationID) (*Declaration, bool)
func (p *Program) Value(id DeclarationID) (TypeValue, error)        // the checked value of a type declaration (opaque)
func (p *Program) Proper(expr core.Value) (TypeValue, error)        // admit a closed type expression
func (p *Program) Apply(callee TypeValue, arguments []TypeValue) (TypeValue, error)
func (p *Program) Descriptor(value TypeValue) (core.Value, error)   // materialize; may contain wire
func (p *Program) GroundDescriptor(value TypeValue) (core.Value, error) // only for Wire-free (ground) types
type DeclarationID struct { Module Digest; Name string }
```

Facts that matter here:

- **`GroundDescriptor` fails for any type containing `wire`.** A whole service
  record or a `Function<A, B>` fails. An operation's Wire-free argument type A
  or result type B succeeds.
- **Getting A out of a field `Function<A, B>` means walking the checked
  expression tree.** `(*Declaration).Expression()` returns it. An application's
  `Callee().Target()` returns the callee's `DeclarationID`, and `Arguments()[0]`
  is A. bittype's analysis already walks this tree for type generation.
- **There is no slot for documentation or annotations.** A module is
  `D("module", I(4), T(name), imports, declarations)`. A declaration is
  `D("type", parameters, kind, body)`. A record field is
  `H("field", T(fieldName), A("req"|"opt"), FieldType)`. All of these have
  exactly checked arities. The specification says "Source spans and native
  naming choices are separate Bittype metadata". bittype's readable syntax
  treats comments as trivia, not canonical data.

### 3.4 horos's TypeScript checker (run time)

```ts
// @bitspark/horos 0.9.0, published (the repository's unpublished source says 0.10.0; the checker is unchanged)
// Value is an ontos ground value (Atom | Tuple)
export function wf(t: Value): boolean;                    // the descriptor is well-formed
export function conforms(v: Value, t: Value): boolean;    // total, never throws
export type Diagnosis = { verdict: "conforms" } | { verdict: "mismatch" }
  | { verdict: "not_well_formed" } | { verdict: "exhausted"; limit: Limit };  // Limit: which bound was hit
export function diagnose(v: Value, t: Value): Diagnosis;  // no mismatch location is reported
```

- **The input must be ground.** A hydrated tuple containing a capability leaf
  never conforms, except to `any()`. Its Wire-free subtrees are ordinary ground
  tuples and can be checked.
- **Go and TypeScript should agree.** The Go and TypeScript checkers and their
  shared conformance corpus are unchanged since this package was published, so
  descriptors produced by the Go declarations layer should be accepted. No test
  exercises that combination.

### 3.5 ontos's JSON projection

ontos publishes one named JSON projection, `ontos-data-json/1`. It is **one-directional**, from domain JSON to ontos data:

| JSON | ontos value |
| --- | --- |
| object | map |
| array | list |
| string | `utf8-text` |
| integer literal | `int` (fractions and exponents are rejected) |
| `true`, `false` | `bool` |
| `null` | `null` |

- There is no reverse function.
- No JSON form produces bytes or a set.

### 3.6 bittype's generation interface

**The roles.** Mapping libraries **describe** native intent. A kernel **plans**: it walks the source's semantic graph itself and selects exactly one rule per occurrence by exact identity. A target **builds** native syntax from the closed plan.

**The type parameters** used in these signatures: F is the source language (horos), T the target language (for example TypeScript), and R the source-specific reader that libraries use to inspect a horos form.

```go
// Selector is an exact source-owned identity (the source language plus a key the kernel never interprets).
// horos keys: "declaration/<module pin>/<name>" and "builtin/<constructor>/<discriminator>".
type Selector struct { /* language, key */ }

type Pattern[F, R any] struct {      // which uses a rule applies to
	Selector Selector
	Slots    []Slot                   // named placeholders, e.g. for type arguments
	Capture  func(context.Context, R, UseView[F]) ([]Capture[F], error) // fills slots with child uses
}
type BindingDescription[F, T, R any] struct {   // a rule for USES of a selected declaration
	ID     string
	Input  Pattern[F, R]
	Output Expression[T]              // native intent, a static template over the slots
}
type DeclarationRecipe[F, T, R any] struct {    // a rule for PUBLISHING a declaration natively
	ID       string
	Selector Selector                 // the zero selector supplies the default structural policy
	Expand   func(context.Context, DeclarationRequest[F, T, R]) (DeclarationDescription[T], error)
}
type MappingDescription[F, T, R any] struct {
	Profile      string               // the representation profile it implements
	Bindings     []BindingDescription[F, T, R]
	Declarations []DeclarationRecipe[F, T, R]
}
type MappingLibrary[F, T, R any] interface {
	Describe(from, to representation.Language) (MappingDescription[F, T, R], error)
}
```

**Native intent is an immutable expression tree.** Its nodes are target-owned
forms carrying scalar data, together with slots, generic parameters and
references. In bittype's words: "The shared interpreter understands substitution
and references, not the meanings of these forms."

**Selection is exact.** For a declaration, the planner:

- uses the recipe registered for its selector;
- falls back to the default structural recipe only if no use rule exists for
  that selector ("an exact use mapping cannot silently become a structural
  publication");
- refuses with an error when no recipe applies, when more than one applies, and
  on a recursive declaration.

**A target turns the closed plan into native syntax:**

```go
type Target[T any] interface {
	CheckExpression(context.Context, Expression[T]) error
	CheckDeclaration(context.Context, NativeDeclaration[T]) error
	CheckProgram(context.Context, Program[T]) error
	Build(context.Context, BuildInput[T]) (Assembly, error)   // closed program → native syntax tree
}
```

The result is re-checked by the target language's own semantic checker before
it is released.

**A library maps one of its declarations to an existing native type.** Here
bittype's tests map a library type `Cell<T>` to a TypeScript package type. The
test's horos fixture is a simplified `type Cell<T> = wire<T>`, not the record
of section 3.1. In the snippet:

- `source.Named(lang, pin, name)` builds the exact selector;
- `g.Mapped(slot)` refers to the native form of whatever fills the slot;
- `model.External(module, export, args...)` is a TypeScript reference to an
  exported type.

```go
key, _ := source.Named(from, pin, "Cell")       // exact identity: the module pin and the name
slot := g.TypeSlot("V")
pattern := g.Pattern[h.Language, source.Reader]{Selector: key, Slots: []g.Slot{slot},
	Capture: func(ctx context.Context, src source.Reader, use g.UseView[h.Language]) ([]g.Capture[h.Language], error) {
		for _, a := range use.Actuals {           // the actual argument bound to Cell's parameter T
			if p, err := src.Parameter(a.Binder); err != nil { return nil, err
			} else if p.Name == "T" { return []g.Capture[h.Language]{{Slot: slot, Use: a.Use}}, nil }
		}
		return nil, g.ErrScope
	}}
rule := g.BindingDescription[h.Language, ts.Language, source.Reader]{ID: "cells/typescript/1",
	Input: pattern, Output: model.External("@bitspark/cell", "Cell", g.Mapped[ts.Language](slot))}
```

A later placement stage chooses the files, writes
`import type { Cell } from "@bitspark/cell"`, and verifies the binding of each
use with the TypeScript compiler. Cell's body is never opened.

**What the targets can express today.** The native syntax vocabularies describe
**types only**:

- value declarations, function bodies, statements and executable expressions are
  explicitly outside them;
- the TypeScript placement stage handles type-only imports, not run-time
  imports.

bittype's own note:

> The existing native AST/text profiles cover type declarations, not executable
> adapter bodies. Generating such bodies requires separately scoped target syntax,
> checking, formatting and compiler qualification.

**The current representation profile:**

| horos form | TypeScript | Go | Rust | C++ |
| --- | --- | --- | --- | --- |
| UTF-8 text | string | string | String | std::string |
| Boolean | boolean | bool | bool | bool |
| list-of(T) | T[] | []T | Vec<T> | std::vector<T> |
| Named closed record | interface | struct | struct | struct |
| Optional field | field?: T | BittypeOptional[T] | Option<T> | std::optional<T> |
| Type parameter, named application | native generics | native generics | native generics | templates |

**Outside the profile today:** integers, bytes, unions, tuples, maps,
recursion, open records and capabilities. These extensions are tracked
separately. They are needed before any of the four services can be generated,
but their representation choices are not the subject of this consultation.

**Each target's mapping library decides naming:**

- Go capitalises field names and wraps optional fields in a generated
  `BittypeOptional[T]{Present, Value}`.
- TypeScript keeps names.

bittype's charter says it "emits no generation manifest": no table of these
choices is emitted.

### 3.7 bittype's own sketch of the binding layer (a proposal, not implemented)

```go
// BindingContext stands for explicitly supplied resource/lifetime facilities.
type Adapter[T any] interface {
    ToValue(ctx BindingContext, value T) (wire.HydratedValue, error)
    FromValue(ctx BindingContext, value wire.HydratedValue) (T, error)
}
```

```text
Type generation:          type description + native representation -> native declaration -> source
Future binding generation: resolved descriptions + selected library bindings -> adapter source
Runtime conversion:       Adapter[T] + hydrated instance + binding context <-> native T (or error)
```

> Future binding generation can emit static imports and compose the matching
> library adapter factories, preserving the selected contract and representation
> identities. Recipes for types and for runtime behavior must be qualified
> together where behavior is required; successful compilation alone is
> insufficient.

### 3.8 The platform's current binding machinery (TypeScript contracts)

All four services have used this machinery since 8 November 2026. It follows
the 29 October decision for TypeScript-sourced contracts.

**The adapter contract and the structural runtime** (`@bitspark/service-api`):

```ts
export interface AdapterContext {   // a transaction over the endpoints exported during one conversion
  readonly owner?: ApiWireOwner;    // the lifetime owner of those endpoints
  mark(): readonly [number, number];
  rollback(mark?: readonly [number, number]): void;   // close endpoints exported after the mark
  commit(): void;
  find(value: object, key: unknown): AdapterEntry | undefined;   // reuse an existing export of the same object
  add(value: object, key: unknown, entry: AdapterEntry): AdapterEntry;
}
export interface Adapter<T = unknown> {
  readonly requiresOwner?: boolean;              // the value may contain capabilities
  readonly acceptsUndefined?: boolean;
  readonly arity?: readonly [number, number];    // for argument tuples: [required, allowed]
  toWire(value: T, context?: AdapterContext): HydratedValue;
  fromWire(value: HydratedValue, context?: AdapterContext): T;
}
export interface WireCodec {                     // a semantic codec: returns a GROUND value, so no capabilities
  toWire(value: unknown): Value;
  fromWire(value: HydratedValue): unknown;
}
export type WireSchema =
  | { kind: 'string' | 'number' | 'boolean' | 'null' | 'bytes' | 'void' }
  | { kind: 'literal'; value: string | number | boolean | null }
  | { kind: 'array'; element: WireSchema }
  | { kind: 'tuple'; elements: readonly { optional: boolean; schema: WireSchema }[] }
  | { kind: 'object'; properties: Record<string, { optional: boolean; schema: WireSchema }>; additional?: WireSchema }
  | { kind: 'union'; variants: readonly WireSchema[] }
  | { kind: 'ref' | 'codec'; name: string }         // codec: a semantic codec looked up by its string key
  | { kind: 'binding'; adapter: Adapter<unknown> }; // a library-owned adapter object
export function dataAdapter<T>(schema: WireSchema, options?: { types?: Record<string, WireSchema>; codecs?: Record<string, WireCodec> }): Adapter<T>;
// structural factories, each a dataAdapter whose children are {kind:'binding', adapter}:
export function tupleAdapter(elements): Adapter;  export function arrayAdapter(element): Adapter;
export function recordAdapter(properties, additional?): Adapter;  export function unionAdapter(variants): Adapter;
// composition helpers:
export function lazyAdapter<T>(build: () => Adapter<T>): Adapter<T>;             // recursion: build on first use
export function scopedAdapter<T>(build: (scope) => Adapter<T>): (scope) => Adapter<T>;   // cache per generic scope
export function genericAdapter<T>(names: readonly string[], build): (scope) => Adapter<T>; // cache per argument adapters
export function adapterSlots(names: readonly string[], supplied): Readonly<Record<string, Adapter>>;  // check supplied slots
```

**`dataAdapter` is one generic interpreter, so it decides the hydrated form of
every structural value:**

| Native value | Hydrated form chosen by `dataAdapter` |
| --- | --- |
| string | an atom holding the UTF-8 bytes |
| number | an atom holding JavaScript `String(n)` of a finite number |
| boolean | an atom `true` or `false` |
| null, void | an empty atom |
| bytes (`Uint8Array`) | the raw atom |
| array | a tuple of the elements |
| object (record) | a tuple of `(keyAtom, presence)` pairs sorted by key; presence is `()` when absent, `(value)` when present |
| union | `(indexAtom, body)`. Encoding tries the variants in order and the first success wins; endpoints exported by failed attempts are rolled back. Decoding uses the index. |
| binding | whatever the library adapter returns |

The limits per conversion are 100,000 nodes, depth 256 and 16 MiB of bytes.

**The operation description** is static data, also read by presentation:

```ts
export type ApiSchema = /* the WireSchema kinds, except that codecs are named and a library binding is static data: */
  | { kind: 'binding'; factory: { module: string; export: string }; arguments: readonly ApiSchema[] } /* | ... */;
export interface ApiOperation { path: readonly string[]; summary: string; parameters: readonly ApiParameter[]; result: ApiSchema }
export interface ApiParameter { name: string; optional: boolean; flatten: boolean /* present an options object as flags */; schema: ApiSchema }
export interface ApiDescription { name: string; summary: string; operations: readonly ApiOperation[]; types: Record<string, ApiSchema> }
```

**The Function contract library** (`@bitspark/service-contracts`). Its native
type carries its adapter association as a doc-comment annotation:

```ts
/** @wireAdapter {"module":"@bitspark/service-contracts/function-wire","export":"functionAdapter"} */
export interface Function<Args extends readonly unknown[], Result> {
  (...args: Args): Promise<Result>;
}
export function functionAdapter<Args extends readonly unknown[], Result>(
  input: Adapter<Args>, output: Adapter<Result>, options?): Adapter<Function<Args, Result>>;
export function serviceBinding<Api>(description: ApiDescription, operations: readonly OperationAdapters[], options?: BindingOptions): ServiceBinding<Api>;
export interface ServiceBinding<Api> {
  fromWire(capability: HydratedWire, options?): Client<Api>;   // a client proxy with the Api's methods
  toWire(service: Api, options?): Endpoint;                    // export a provider
}
export interface OperationAdapters { path: readonly string[]; input: Adapter<readonly unknown[]>; output: Adapter }
export interface BindingOptions {                              // service-supplied hooks
  errors?: WireErrors;                                         // the service's error mapping
  captureArguments?: (path: readonly string[], args: readonly unknown[]) => readonly unknown[];
  validateResult?: (path: readonly string[], args: readonly unknown[], result: unknown) => void;
}
export interface WireErrors {
  request(error: unknown): Value;            // provider: the request could not be decoded or captured
  encode(error: unknown): Value;             // provider: the implementation threw
  unknown(reason: string, message: string): Value;
  readAnswer<R>(answer: HydratedValue, decode: (value: HydratedValue) => R): R;   // caller
  transport(code: 'refused' | 'unknown', message: string, cause?: unknown): unknown;
}
```

**The call protocol:**

- **A callable value** travels as `(A("bitspark/function/1"), invokeCapability, releaseCapability)`. `functionAdapter.toWire` exports an endpoint under the lifetime owner, reuses the export when the same native function appears twice, and returns an imported proxy unchanged.
- **A request** is `(A(protocolTag), (pathAtoms...), argumentsTuple, replyCapability)`, with protocol tag `bitspark/service/1`.
- **An answer** is `(A("ok"), value)` or an error answer.
- **Default error answers:**
  - an undecodable request becomes `("protocol", message)`;
  - a provider exception becomes `("unknown", "The function failed after invocation")`;
  - `refused` is produced only locally, when the runtime does not admit a send.

  A service may replace these with its own `WireErrors`. `Function` has no declared failure parameter.

**The presentation library** (`@bitspark/service-cli`) derives command-line
commands and web forms from an `ApiDescription`:

```ts
export interface ApiCodec {                                  // per semantic codec key, supplied by the service
  fromJson(value: unknown): unknown;                         // conversion + semantic validation
  toJson(value: unknown): unknown;
  readonly kind?: 'text' | 'flag' | 'json' | 'bytes';        // how the value is entered
}
export function deriveApiProjection<Api>(description: ApiDescription, options?: {
  codecs?: Record<string, ApiCodec>;
  error?: (error: unknown) => { family?: string; code: string; message: string; details?: unknown } | undefined;
}): ServiceProjection<Api>;                                  // the command and form definitions
```

**Its JSON conversion is a second hand-written interpreter, over `ApiSchema`.**
Its rules differ from `dataAdapter`'s:

- JSON unions are untagged, and the first variant that matches wins;
- nested bytes become `{"$bytes": base64}`;
- `void` becomes `null`.

**A library binding cannot be presented.** If any reachable schema contains
`binding`, the whole derivation fails ("This library binding has no finite
presentation adapter"). Each operation's `summary` appears in command-line help
and web forms. A field's description is just its parameter name.

**The generator** (`@bitspark/contract-generator`) uses the TypeScript
compiler.

- **Operations.** For each member of the service interface it produces an
  `ApiOperation`; the summary is the member's doc comment.
- **Type positions,** resolved in this order:
  1. a bound generic parameter;
  2. a **semantic codec**, when the type's `package#TypeName` key appears in a
     service-supplied codec table. Template-literal types such as
     `` `sha256:${string}` `` must be codecs. This selection is by string key,
     and the same keys supply the run-time Wire codecs and JSON codecs;
  3. a **library binding**: a declaration carrying exactly one
     `@wireAdapter {"module", "export"}` annotation, found by declaration
     identity through aliases;
  4. a generic application;
  5. structure.
- **Refusals.** An unannotated callable type is refused.
- **The root factory.** The consumer's build configuration names it
  (`serviceBinding`); the generator has none built in.
- **Output.** The emitter renders the description as a module, as in this
  abridged real output:
  - `binding` becomes a factory import and call;
  - `ref` becomes a cached type entry;
  - each structural kind becomes a factory call.

```js
// Generated from declaration-owned bindings. Do not edit.
import { dataAdapter as __data, tupleAdapter as __tuple, recordAdapter as __record, unionAdapter as __union,
         scopedAdapter as __scoped, adapterSlots as __slots } from '@bitspark/service-api/wire';
import { serviceBinding as __factory0 } from "@bitspark/service-contracts/service";    // root, named by the consumer
import { functionAdapter as __factory1 } from "@bitspark/service-contracts/function-wire";  // from @wireAdapter
const __description = { /* the ApiDescription as JSON */ };
export function createExampleBinding(__options = {}, __supplied = {}) {
  const __scope = __slots([], __supplied);          // the root's generic slots (none here)
  const __types = Object.create(null);
  __types["T0"] = __scoped(__scope => __factory1(__tuple([{ optional: false, adapter: __data({"kind":"number"}) }]), __data({"kind":"string"})));
  __types["T3"] = __scoped(__scope => __union([__data({"kind":"null"}), __record({["text"]: { optional: false, adapter: __data({"kind":"string"}) }})]));
  return __factory0(__description, [
    { path: ["retain"], input: __tuple([{ optional: false, adapter: __types["T0"](__scope) }]), output: __data({"kind":"void"}) },
    { path: ["notes","get"], input: __tuple([{ optional: false, adapter: __data({"kind":"codec","name":"@example/model#NoteId"}, { codecs: __options.codecs ?? {} }) }]),
      output: __types["T3"](__scope) }], __options);
}
```

**Two properties of this machinery matter most here:**

- **Agreement by construction.** The generator reads the native TypeScript type
  and derives the adapter from that same type, so the native representation and
  the adapter agree by construction.
- **The runtime library decides the message form.** `dataAdapter` decides how a
  record or a union looks as a hydrated value. Neither the type language nor the
  value model does.

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
export type ContentHash = `sha256:${string}`;          // validator contentHash(): /^sha256:[0-9a-f]{64}$/
export type ContentRepresentation = 'raw' | 'map/1' | 'list/1';
export interface ContentLink { readonly representation: ContentRepresentation; readonly hash: ContentHash }  // validator contentLink()
export type ContentValue = Uint8Array | ContentMap | readonly ContentValue[];   // recursive; untagged natively
export interface ContentMap { readonly [key: string]: ContentValue }
```

**Its hand-written Wire codecs.** Each gives a semantic type its own hydrated
form:

```text
ContentHash : A(hash)                                   -- a UTF-8 atom; decoding re-validates
ContentLink : (A(representation), A(hash))
ContentValue: TAGGED on the wire although UNTAGGED natively; a shared sub-value is converted once (DAG sharing)
              ("bytes", atom) | ("list", (value...)) | ("map", ((keyAtom, value)...))
              decoding enforces depth/step/byte limits, rejects cycles, unknown tags and duplicate keys,
              then deep-copies and freezes the result
```

**The binding hooks the service still writes by hand** (abridged; `field(x, name)`
reads exactly one own data property):

```ts
const binding = createContentServiceBinding({          // the generated module (section 3.8)
  codecs: {                                            // Wire codecs for the three semantic types, by string key
    '@bitspark/content-model#ContentHash':  { toWire: v => contentHashToWire(contentHash(v)), fromWire: contentHashFromWire },
    '@bitspark/content-model#ContentLink':  { toWire: v => contentLinkToWire(contentLink(v)), fromWire: contentLinkFromWire },
    '@bitspark/content-model#ContentValue': { toWire: v => contentValueToWire(snapshotContentValue(v)), fromWire: contentValueFromWire },
  },
  captureArguments(path, args) {                       // caller: before encoding; provider: after decoding
    switch (path.join('.')) {
      case 'put': {                                    // copy the caller's bytes before any await
        const bytes = field(args[0], 'bytes');
        if (!(bytes instanceof Uint8Array)) throw new ContentModelError('Expected bytes');
        return [{ bytes: new Uint8Array(bytes) }];
      }
      case 'value.put': return [{ value: snapshotContentValue(field(args[0], 'value')) }];  // deep capture with limits
      default: return args;
    }
  },
  validateResult(path, _args, result) {                // `number` natively; the contract means a non-negative integer:
    if (path.join('.') === 'check') decimalAtom((result as { objects: number }).objects);  // throws otherwise
  },
  errors: { request: /* protocol or domain answer */, encode: errorToWire, unknown: /* … */, readAnswer, transport: /* … */ },
});
```

**Content's errors are declared only in code:**

- domain codes: `not-found | permission-denied | authentication-required | authorization-unavailable`;
- model codes: `invalid-data | validation-limit`, with details `{limit, maximum}`.

How they travel:

- **Shape on the wire:** `("domain", family, code, message, ((key, ("text", t) | ("integer", digits))...))`.
- **Limits:** at most 16 details, keys of at most 64 bytes, text of at most 256 bytes.
- **Caller side:** the caller rebuilds the same error classes.
- **Everything else** becomes "unknown outcome", never "definitely not executed".

**JSON for presentation is a third hand-written conversion.** For
`ContentValue` it is `{type:'bytes', base64} | {type:'map', entries} | {type:'list', items}`.

**Hand-written code, in lines** (non-blank, non-comment):

| Category | Lines |
| --- | --- |
| Interface and model types | 42 |
| Validators and deep capture | 134 |
| Wire codecs and their helpers | 168 |
| Error mapping and classes | 207 |
| JSON codecs | 79 |
| Codec table | 17 |
| Binding hooks | 81 |

**The other three services need more.** Cubes, Pointers and Drafts (23
operations) use:

- unions discriminated by a `kind` literal, including one variant with two
  literals;
- `string | number`;
- intersections with unions;
- nullable references;
- integers typed as `number`;
- the value model's own ground values, as arbitrary data;
- recursive JSON-like data;
- types imported from other services.

None uses streams or generics at its root.

### 3.10 The dormant bitlink design (prior art)

bitlink's September design for model-to-wire bindings kept types with the type
language and gave interfaces, bindings, clients and validators to bitlink. It
never started, because it waited on a type-language design that horos later
replaced. The principles that still bear on this question:

- "Declare once."
- "No runtime dependency on the generator."
- A candidate transparency law for generated bindings,
  `stub(bind(f)) ≈ f` and `bind(stub(a)) ≈ a`:
  - `bind` exports a native function as a capability;
  - `stub` imports a capability as a native proxy;
  - `≈` means observably equivalent.

## 4. Data and control flow

### 4.1 One call today

```text
caller                                                   provider
──────                                                   ────────
client.value.put({ value })
  check the argument count (input adapter's arity)
  captureArguments (service hook)
  input adapter.toWire → hydrated arguments
    (structure by dataAdapter, semantic codecs by key,
     library adapters; callables exported into the
     call's transaction)
  request = (tag, path, arguments, replyCapability)
  bitruntime: export capabilities, send ─────────────▶ bitruntime delivers
                                                         check tag and path; decode arguments;
                                                         captureArguments; re-encode and decode
                                                         invoke the implementation
                                                         output adapter.toWire; validateResult
  ◀───────────────────────────────── ("ok", value) or an error answer (the service's WireErrors)
  readAnswer → output adapter.fromWire → validateResult → result, or a rebuilt error
```

### 4.2 Three representations of one declared type

For a type declared in horos, three representations exist:

| Representation | Who could decide it under a horos contract | Who decides it today |
| --- | --- | --- |
| Native type and value | bittype's representation profile, and a library's exact recipe for its own declarations | the TypeScript interface itself |
| Hydrated (message) form | the horos descriptor and ontos's canonical data forms, **or** the binding runtime, as today. **Open** (question 4). | `service-api`'s `dataAdapter`, or a service's semantic codec |
| JSON (presentation) | open (question 3) | `service-cli`'s interpreter, plus service JSON codecs |

**The candidates differ concretely.** Take Content's `ContentLink`, declared in
horos as `record closed { representation: union(exact "raw", exact "map/1",
exact "list/1"), hash: ContentHash }`, with `ContentHash` declared as text:

```text
conforming to the horos descriptor (canonical ontos data; text(s) = Tuple(A("utf8-text"), Atom(bytes of s))):
  Tuple(A("map"),
        Tuple(text("hash"),           text("sha256:…")),
        Tuple(text("representation"), text("raw")))           -- entries in canonical key order
today's generic dataAdapter, if ContentLink were structural (the literal union is index-tagged):
  Tuple( Tuple(A("hash"),           Tuple(A("sha256:…"))),
         Tuple(A("representation"), Tuple(Tuple(A("0"), A("raw")))) )
today's hand-written Content codec:
  Tuple(A("raw"), A("sha256:…"))
```

Two more differences:

- **Integers.** Today a count is an atom holding decimal text. In ontos it is
  `Tuple(A("int"), Atom(<sign-magnitude bytes>))`, unbounded.
- **Unions.** Today a union carries its variant index. A horos union is
  untagged: alternatives are told apart by content, such as compound labels or
  exact literals.

**If the hydrated form follows the horos descriptor:**

- the wire format changes;
- a structural adapter takes its two sides from two different owners: the native
  side from the representation profile, the hydrated side from horos and ontos;
- received values can be checked with horos `conforms`.

**If it stays the runtime's choice:** horos descriptors cannot check messages,
and the horos contract no longer describes them.

### 4.3 What a build from a horos contract must produce

**Inputs**, all pinned by digest:

- the service module;
- the model modules it imports, for example a `content-model` module declaring
  `ContentHash`, `ContentLink` and `ContentValue`;
- the contract library modules (Function).

**Products:**

1. **Native types.** bittype's mechanism exists, but the profile extensions
   listed in section 3.6 are needed, and so are exact mappings for `Function`,
   `ContentHash` and `ContentValue`. Neither exists yet.
2. **Per operation:** the path, the argument and result types, and the declared
   failure type.
3. **Per argument and result type, an adapter for each direction.** Structural
   parts follow from the native representation and the chosen hydrated form.
   Library parts (`Function`, `ContentHash`, `ContentValue`) come from the
   declaring libraries' factories, selected by declaration identity.
4. **A ground descriptor per Wire-free argument and result,** produced by horos
   Go `GroundDescriptor`, if received values are to be checked with
   `conforms`.
5. **Presentation data:** parameter names and entry hints, summaries, and a
   JSON conversion.

### 4.4 Where Content's hand-written parts might go (a hypothesis)

| Hand-written today | Possible home under a horos contract |
| --- | --- |
| The TypeScript service interface | the horos service module; TypeScript types generated |
| The codec table keyed by `package#Type` strings | library adapter factories selected by declaration identity |
| `ContentLink` codec | structural (needs union support in the profile) |
| `ContentHash` codec and its pattern | text in a model module; the model library supplies semantic validation, since horos has no pattern refinement |
| `ContentValue` codec (tagged on the wire, untagged natively, DAG sharing, limits) | a library-owned custom representation and adapter |
| `captureArguments` (copy bytes, deep-snapshot before any await) | perhaps inherent, if conversion runs synchronously when the call starts |
| `validateResult` (`check` returns an integer) | gone, given an integer type in the declaration |
| Error mapping (207 lines) | declared failure types plus the Function library's protocol |
| JSON codecs | open (question 3) |

## 5. Constraints

- **Fixed rules.** The architecture rules (section 2.3) and the contract-owned
  binding decision (section 2.4) are fixed.
- **horos.**
  - It owns type semantics. It has no generator and no runtime binding.
  - Its declaration semantics exist only in Go; ground conformance exists in Go,
    Rust and TypeScript.
  - Changes to horos may be recommended, but horos's maintainers decide them.
- **bittype.**
  - It is a Go library and must not depend on the runtime packages (bitwire,
    bitruntime). Emitting code that imports them is not a dependency.
  - Its generated output is re-checked by each target language's checker.
- **Selection.** No component may hold a list of contract libraries, and no
  component may select behaviour by type name or by shape.
- **What may change.** The native TypeScript API, the hydrated form and the
  wire format may change. Whether the hydrated form should follow the horos
  descriptor is open; it is not a constraint.
- **Process.** No release checksum or validation workflow is a prerequisite;
  ordinary tests and type checks are. Security hardening is deferred during
  development.

## 6. What we have considered

### 6.1 Who emits, and what is emitted

Two choices are open and partly independent: **who emits** binding products,
and **what is emitted**.

**What is emitted:**

- **Code.** Generated adapter code per target language.
- **Descriptions.** Data that a hand-written run-time interpreter per target
  language reads. The description carries:
  - the hydrated form;
  - the native representation choices;
  - static imports of library adapter factories.

  Today's platform machinery is mostly this second kind: `dataAdapter`
  interprets a schema, and only the composition is generated code.
  - **For descriptions:** the generator needs no executable target syntax.
  - **Against descriptions:** the interpreter is an executable rendering of the
    representation profile, a second place where profile semantics live, unless
    the profile itself is defined as data both sides read.

**Who emits:**

- **A. bittype.** Binding products become a second product of the same kernel,
  representation profiles and exact selection. A library contribution then
  supplies, beside its native type recipe, a reference to its adapter factory
  for that target.
  - **For:**
    - One owner decides both the native representation and its conversion.
    - The kernel already resolves imports and exact identities, binds generic
      arguments and records external symbols for imports.
  - **Against:**
    - Emitting *code* would need executable syntax per target, which is a large,
      new and unqualified surface.
    - bittype's charter would change.
- **B. A generator owned on the run-time side.** It consumes bittype's analysis
  plus a **representation manifest** that bittype would publish: per
  declaration, its native names, field names, optional encoding, union
  representation and external symbols. The platform repository's existing
  generator is a candidate.
  - **For:**
    - The generator lives next to the run-time libraries whose adapter contract
      and resource semantics it targets.
    - It already composes library factories by declaration identity.
    - One generator reading the manifest could serve several target languages,
      if it has each language's run-time knowledge.
  - **Against:** The manifest must be complete and authoritative. Otherwise one
    representation decision gets two implementations, bittype's and the
    generator's, which can drift.
- **C. Runtime interpretation of horos modules** (rejected for now).
  - It needs a second implementation of horos declaration semantics in
    TypeScript.
  - It needs a run-time table from declaration identity to factory, which
    section 2.4 forbids or comes close to forbidding.

### 6.2 Narrower decisions, with our present leaning (none decided)

1. **The hydrated form.** Our leaning is canonical, descriptor-conforming data,
   because the horos declaration is the contract (section 4.2).
2. **Shape checks of received values.** Wire-free arguments and results would
   be checked with horos `conforms` against build-emitted ground descriptors
   before conversion; adapters would then add only semantic validation, such as
   the hash pattern.
   - **Limits:** no mismatch location is reported, and capability-containing
     values can be checked only in their Wire-free parts.
3. **Library contributions.** Per declaration and target, a library supplies a
   native type recipe (bittype mapping libraries already do) and an adapter
   factory reference parameterised by argument adapters (today's
   `@wireAdapter`). Where both live is undecided:
   - a per-library mapping library in Go;
   - metadata beside the library's horos module;
   - something else.
4. **Declared failures.** Candidate: `Function<A, B, E>`, with reply
   `Outcome<B, E> = compound("ok", B) | compound("error", E)`. Protocol
   outcomes (refused, unknown, protocol, limit) would stay the Function
   library's own. The Function library's maintainers must agree.
5. **Documentation.** Three options:
   - horos adds a documentation slot, which changes canonical bytes and
     therefore declaration identity;
   - bittype's readable syntax turns doc comments into metadata emitted beside
     the canonical module;
   - a separate documentation module keyed by declaration identity.
6. **JSON.** Three options:
   - derive it from the canonical value form through a total, two-way named
     projection; ontos's is one-way and has no bytes or sets;
   - derive it from the native representation, as `service-cli` does now;
   - library-owned JSON adapters for custom representations, structural ones
     derived.

   The two existing interpreters already disagree on unions, `void` and arity.
7. **Capture before asynchronous work.** If conversion to the hydrated form
   happens synchronously when a call starts, the conversion *is* the capture,
   and a separate hook is unnecessary for data. Callables and streams are
   converted by their libraries.

### 6.3 Two arguments we want tested

- **For a single emitter (A, emitting code or descriptions).** One
  representation decision must have one owner. Today the platform achieves this
  by deriving types and adapters from one TypeScript source. With a horos
  contract the native types come from bittype. Any other emitter must re-derive
  bittype's choices and keep them in step as profiles and targets grow.
- **For a run-time-side emitter (B, with a manifest).** bittype is a type tool.
  The adapter contract, resource transactions and protocols belong to the
  run-time libraries, so the generator that composes them belongs beside them.
  An authoritative manifest from bittype removes the duplication. It may also
  be what keeps bittype's scope small and stable.

AC-09 rules out using the platform's existing generator, or the cost of moving
work, as a reason for either.

## 7. Questions for the expert

1. **Who emits, and what.**
   - Who should emit binding products, and in what form? The choices are
     bittype, or a generator on the run-time side reading a representation
     manifest that bittype publishes (section 6.1); and adapter code, or
     descriptions interpreted by a per-language run-time library.
   - Which principles decide each choice, given one owner per representation
     decision and additive change for libraries, services and target languages?
   - Please weigh the two arguments in section 6.3.
   - Which one or two lessons from comparable systems transfer most directly
     (for example protobuf code-generation plugins, Smithy, serde, Avro's
     specific and generic records)?
2. **The library contribution contract.**
   - How should a contract or domain library contribute its native type and its
     adapter factory for one of its declarations and one target language, so
     that composition selects them by exact identity, with argument adapters for
     generic parameters?
   - Should that contribution live in bittype's mapping library, beside the
     library's horos module, or elsewhere?
   - Which laws should every contribution be tested against (for example round
     trip, conformance of produced values, capability ownership and rollback)?
   - As a worked case: how should the Function library separate an operation's
     declared domain failures from its own protocol outcomes?
3. **Presentation and documentation.**
   - The command line and web forms need native↔JSON conversion and
     documentation (summaries, entry hints). What principle should decide what
     the JSON form derives from, so that libraries do not write one adapter per
     projection?
   - Where should documentation live when the type language's canonical form
     has no place for it?
4. **The message form and the pilot.**
   - What should decide the hydrated form of a structural value: the
     descriptor-conforming canonical data, or the binding runtime?
   - Given that, how should shape checking be divided between conformance and
     adapters?
   - What should a first pilot on Content qualify, and which falsification tests
     would show that adding a second contract library, a second service, or Go
     as a second target language changes only its own owner?
5. **What are we not asking?** What would you want to know, or change, before
   committing to any of these designs, including patterns from systems we have
   not named? Which risk in this layering is most likely to be underestimated?
