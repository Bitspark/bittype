# Research: optional Wire value and type extensions without losing the existing design

> Historical consultation and dated evaluation. See the [research map](README.md)
> for current contracts and subsequent decisions; original advice/submissions
> retain their inspected baseline.

**ID:** 0003
**Date:** 3 November 2026
**Status:** applied
**Graded:** https://github.com/Bitspark/bittype/issues/35#issuecomment-6611597047
**Run-ID:** run_31c1813d-dd34-46a5-9873-105a934581fb
**Reviewed:** https://github.com/Bitspark/bittype/issues/35#issuecomment-6611112491
**Owner:** codex-hydrated-type-design
**Issue:** https://github.com/Bitspark/bittype/issues/35
**Related issues:** https://github.com/Bitspark/ontos/issues/23; https://github.com/Bitspark/horos/issues/17

## Question

**Response requirement for the revised submission:** provide a complete, concise
answer to all five questions below in at most 1,500 words. Do not repeat the brief
or write an extended introductory report. Prioritize counterexamples, precise
amendments and qualification obligations, especially partial lifting, profile-aware
identity and library-owned native recipes. An earlier genuine response stopped
mid-sentence during section 3; a fresh capture returned identical bytes, so it
did not cover these questions completely. The technical brief and its four
reviews are unchanged; this revision only bounds the response to obtain a
complete assessment. Do not treat any prior response as available context.

How should optional `ontos/wire` and `horos/wire` extensions own the mixed
data/capability value domain and its types, while preserving the frozen ground
models, existing generic/lifting ideas and Bittype's language-agnostic tooling?

The owner explicitly selected this documentation direction, then requested an
expert consultation and preservation of the ideas already developed. Exact words:

> "OK, document this in ontos, horos and bittype: ontos/wire and horos/wire as optional extensions, replacing the Bittype-owned contract model proposal. No implementation yet."

> "Also make sure this doesn't destroy any concepts or ideas that we had, and consult an expert as well"

The owner also clarified that Deixis is proposed for the **value model**, not as
a replacement for Horos's type semantics or Bittype's AST. We seek a critical
assessment, including counterexamples and modifications to the proposed profile;
the question is not a request to ratify a rename. This delivery changes docs only.
Future implementation obligations must be distinguished from delivered behavior.

## Context

### Vocabulary and system overview

All components below are independent libraries/repositories in one software family.

| Component | Existing responsibility |
| --- | --- |
| Deixis | Finite trees with one opaque own value at each node and finite byte-keyed child maps. Selection/reconstruction do not interpret the own value. |
| Ontos | Ground values: immutable byte atoms and ordered tuples, structural identity, a frozen canonical byte codec, and data readings such as text/list/map. |
| Horos | Descriptions of sets of Ontos values; checking descriptor well-formedness and value conformance. Go declaration tooling additionally resolves modules and first-order generics and specializes them to ground descriptors. |
| Bittype | Go-implemented generic AST/CST tooling, language-owned syntax/checking, agnostic semantic facts, explicit resolver/analyzer/transformer registration, route planning and native text tooling. AST means abstract syntax tree; CST means concrete syntax tree with tokens/trivia. |
| Bitwire | Ground sending/endpoint/addressed contracts, a hydrated live-value contract and pure codecs for its data frames. A Wire grants sending only; an Endpoint additionally owns receiving and closure. |
| Bitruntime | Concrete queues/carriers/endpoints and hydration: context-dependent export of recognized sending capabilities and import of proxies. |
| Contract libraries | Ordinary Function, Cell, Stream and future contracts, their protocols, native type representations and adapters. They are not built-in kinds in any general runtime or generator. |

Ground means finite, complete and context-free; a live capability is not ground.
A descriptor is data that describes a type, not an instance of that type.
A declaration identity is an exact module-document digest plus declaration name,
not a claim that its expanded shape is unique. A profile fixes an admitted domain
and interpretation; a matching name alone does not establish compatibility.
Ontos L0 is its core value model, L1 its labelled-compound reading and L2 its data
embeddings. SPI means service-provider interface: a boundary implemented by
explicitly supplied providers, not a built-in catalogue of languages or contracts.

Ontos/Horos have Go, Rust and TypeScript implementations (Python supplies validation
coverage for their ground model). Bittype native AST/text tooling targets TypeScript,
Go, Rust and C++; actual canonical-to-native transformers are not implemented.
There is no measured performance problem to solve in this documentation task.

### Inspected sources and evidence boundaries

| Repository | Exact inspected revision |
| --- | --- |
| deixis | e0b0410d6f76d795992cc0a5792c507a19ef05ee |
| ontos | a3f1d081f2e82f1e0ec39688c4060b5606b16c63 |
| horos | 14befc691738d1d51568e0c53843c067817ab33c |
| bittype | 9c868646c09a449c6e682328840434126be1301f |
| bitwire | e2d2d641744a49f1d8098b6ea5f2724454d72b96 |
| bitruntime | 64dbe709f2f002781a5010b08a6650c7eeb129df |

These are source observations, not a deployment certification. Bitwire currently
bundles a checked consumer mirror of Ontos v0.9.0 (core/codec/data), source
5eed85de5697f53397407f9f82f3b310648827f7. Runtime consumers use that value family.
None of these Wire/runtime layers currently imports Horos to hydrate messages.
The proposal neither changes that mirror nor makes unrelated Go package types
nominally identical. Such migration would be separate qualified implementation.

### Relevant code: the value floor and its existing projection

Deixis's accepted model is `Node[T] = T × FinMap[Bytes, Node[T]]`. Each node has
both parts, even when its child map is empty. There are no structural node kinds.
The tree is finite/well-founded, but that restriction does not inspect inside T.
Construction/navigation require only T; optional equivalence and codec operations
require extra supplied laws. Live handles may be carried without any equality or
portable encoding claim. `Data` in historical Deixis `DataTree` means a byte-reader
capability, not an Ontos ground value; those terms must not be conflated.

Existing Go structural interface, deixis `core/go/tree.go:22-28`:

```go
type Key = []byte
type TreePath = [][]byte
type Child[T any] struct { Key Key; Tree DeixisNode[T] }
type DeixisNode[T any] interface {
    Own() T
    Children() []Child[T]
    At(TreePath) (DeixisNode[T], bool)
    Decompose() (T, []Child[T])
}
```

Deixis `core/go/tonative.go:54-65` says its structural conversion copies keys and
retains own values unchanged: it never clones/invokes them, fetches, or binds.
Reconstruction preserves capability handles and their aliasing. It is not a
domain adapter or hydration operation.

Ontos L0 is the closed algebra `Value = Atom(Bytes) | Tuple(Value*)`; no live
resource, unknown or generic parameter is an L0 value. Atom/Tuple identity and
`ontos-codec-v1` bytes are frozen. The existing optional projection
`ontos-over-deixis-v1` uses the unchanged positional key spelling
`kappa : natural number -> Bytes`, from `deixis-pos-v1`. Keys must be exact,
dense positions 0 through n-1; a recognizer never repairs sparse/noncanonical keys.

Normative projection equations, ontos `projection/deixis/PROJECTION.md:43-57`:

```text
P(Atom(b))          = Node(Some(b), {})
P(Tuple(v0,...,v(n-1))) = Node(None, { kappa(i): P(vi) for each position i })

R(Node(Some(b), children)) = Atom(b) only if children is empty
R(Node(None, children))    = Tuple(R(child0), ...)
    only if keys are exactly the dense canonical positions and all children recognize
```

Actual Go aliases, `projection/deixis/go/projection.go:39-42`:

```go
type Payload = deixis.Option[[]byte]
type Node = deixis.Node[Payload]
```

`None` here is a supplied payload alternative, not a missing own value, missing
field, missing node or null. Empty atom and empty tuple are different childless
nodes. The profile supplies equality over the whole Option, including that distinction.
Its laws include `R(P(v)) = v`, exact recognition of P's image and identity transport.
The byte law is recovery followed by the Ontos codec; it does **not** equate the
Deixis codec's bytes with Ontos codec bytes. Three source implementations and a
shared authored corpus exist. Their documented evidence ceiling is not an
independent external implementation proof; this proposal must not inflate it.

### Relevant code: current hydrated values and stateful conversion

Bitwire `wire/go/wire.go:15-17` uses the mirrored Ontos type:

```go
type Wire interface { Send(core.Value) error }
```

Bitwire `wire/go/hydrated.go:3-25` defines these separate live contracts:

```go
type HydratedValue = any // only documented, runtime-admitted cases are legal
type ReceivedContext = any // established by composition, never asserted by message data
type HydratedTuple interface { Items() []HydratedValue }
type HydratedWire interface { Send(HydratedValue) error }
type HydratedEndpoint interface {
    HydratedWire
    Wire() HydratedWire
    Receive(func(HydratedValue, ReceivedContext)) (func(), error)
    Closed() <-chan struct{}
    Termination() Termination
    Close() error
}
```

The actual admissible live domain is atoms, finite immutable tuples and
runtime-recognized send capabilities. An arbitrary implementation of the Go
HydratedWire interface is not automatically recognized/exportable. Endpoint owners
are projected to send-only faces, including local delivery.

Bitruntime `hydrated/go/hydrated.go:58-82` constructs tuples by capturing children,
recognizing atoms/tuples and its own tuple/face/proxy types, rejecting others, and
returning an ordinary Ontos tuple when every child is ground. `Scope.send` at
line551 stages a whole value; its `toData` at644 replaces capability occurrences
with scoped export descriptors. Bitwire's pure codec packs those descriptors into
an Ontos frame. Scope `fromData` at703 and `Deliver` at742 reconstruct proxies and
admit a message after frame/route/scope checks. Ordinary reference-looking data
remains data: conversion does not scan for magic tuples and grant authority.

The transport-oriented structural domain is distinct from live values:

```text
HydratedData = DataAtom | DataTuple(HydratedData*) | DataReference(Path, Scope, ID)
```

Here Path is a sequence of exact byte segments; Scope and ID are data tokens
identifying a hydration scope and export, not live Go Scope objects or authority
proofs. In the inspected protocol each token is exactly 16 octets.

The runtime owns export identity, limits, admission, proxy recognition and cleanup.
Send success is local admission, not remote result or protocol conformance. Closing
a scope ends its exports/references but does not close borrowed domain endpoints.
Raw Deixis or Ontos codecs perform none of those operations.

### Relevant code: current Horos generic semantics

Horos ground descriptors are Ontos tuples tagged `horos/<constructor>`.
They include exact values, scalars, records, typed lists/maps/sets, unions,
intersections, ground complements and guarded recursion. `WF(d)` checks a closed
descriptor; `Conforms(v,d)` checks a ground value. Ground `any` excludes live values.

The separately implemented Go declaration package adds module, import, type,
parameter and application descriptions. Ordered parameters are first-order type
slots: exact arity, no inference/defaults/partial application/higher-kinded slots.
The supplied module bundle is immutable and pinned. All transitive imports and
declarations are resolved, including unused imports; loading performs no network IO.

The actual input uses D(tag, operands...) = Tuple(A("horos/decl/" + tag),
operands...). A is a raw UTF-8 atom, T canonical Ontos text, I its integer
embedding, L its ordered list and M its canonical text-keyed map. H(tag, ...)
is the ground descriptor Tuple(A("horos/" + tag), ...). A compact legal example:

```text
D("module", I(1), T("example.boxes"), M{}, M{
  "Box": D("type", L(T("T")), H("list-of", D("param", I(0)))),
  "Texts": D("type", L(), D("apply", D("local", T("Box")),
    L(H("scalar", A("utf8-text")))))
})
```

This declares Box<T> and Texts = Box<Text>. An imported target instead uses
D("imported", T(alias), T(name)); its module import is
D("import", T(expectedName), Pin(module)), where Pin is exactly 32 SHA-256 bytes
over the expanded module's canonical Ontos encoding. These equations are value
notation, not a delivered text parser. The revision is a semantic boundary.

Actual identity types, horos `core/go/declarations/model.go:17-31`:

```go
type Digest [32]byte // hash of exact canonical module bytes, not semantic equivalence
type DeclarationID struct { Module Digest; Name string }
type ParameterID struct { Declaration DeclarationID; Index int }
```

The immutable Expression model at `model.go:125-159` retains descriptor operand
paths, resolved target declarations, caller-bound application arguments and exact
parameter identities. Operands mark type positions; labels, field names, scalar
selectors and literal contents are not substitution positions.

Specialization, `core/go/declarations/specialize.go:5-8,32-49`, includes:

```go
// Concrete, recursively closed descriptor arguments; original Program unchanged.
func (p *Program) Specialize(id DeclarationID, arguments []core.Value) (core.Value, error)

// Inside expansion (env is the calling declaration's bound arguments):
switch expression.kind {
case Parameter:
    return env[expression.parameter.Index], nil
case Application:
    actual := make([]core.Value, len(expression.arguments))
    for i, argument := range expression.arguments {
        value, err := expand(argument, env, depth+1)
        if err != nil { return nil, err }
        actual[i] = value
    }
    return expand(p.declarations[expression.target].expression, actual, depth+1)
}
```

The excerpt illustrates argument expansion in caller scope before switching to
the callee; the actual implementation also validates descriptor closure, arity,
depth/node/byte/module bounds and the final descriptor, and returns no partial
usable result. Missing context, invalid source/context/configuration and exhausted
requests are distinct diagnoses. Recursive descriptor `rec/ref` binders are not
generic parameters. Module imports and named declaration applications remain acyclic.

### Architectural context: what is being replaced, and what must survive

The present documentation assigns the following **unimplemented** language to
Bittype. With A = raw UTF-8 atom, T = canonical Ontos text, I = canonical integer,
L = canonical list, M = canonical text-keyed map and Pin = SHA-256 atom of exact
canonical module bytes, its proposed prefix is:

```text
C(tag, operands...) = Tuple(A("bittype/contract/" + tag), operands...)
E = C("data", G) | C("wire", E) | C("tuple", E...)
  | C("compound", A(label), E...) | C("list-of", E)
  | C("record", A("closed"), C("field", T(name), A(presence), E)...)
  | C("union", E, E, ...) | C("param", I(index))
  | C("apply", target, L(E...))
target = C("local", T(name)) | C("imported", T(alias), T(name))
C("module", I(1), T(name), importsMap, declarationsMap)
C("import", T(name), Pin(module))
C("type", L(T(parameterName), ...), E)
```

G is a complete closed Horos ground descriptor; presence is req/opt. Data(G) is
opaque to upper substitution. The proposed structural layouts preserve Ontos's
tuple/labelled-compound/list/text-keyed-map readings, with live children admitted
only above the ground boundary. Records are closed; union has at least two branches.
The proposal has no live map keys, sets, open-record extras, top/complement or
named recursion. A Wire's message type is a schema operand, not a structural
child of the live Wire instance. An absence of Wire syntax does not switch profiles.

The new owner direction replaces this Bittype-owned authority. It does **not**
revive an even earlier, also unimplemented idea to add a `horos/wire` constructor
directly to Horos ground declaration revision2. Current ground revision1 and
existing checkers must keep rejecting that syntax. The optional package name
must not become an accidental compatibility decoder for either old proposal.

### Candidate value and type extension

Candidate value profile, proposed only:

```text
MixedPayload = None | Some(BytesCase(bytes)) | Some(WireCase(opaque sending handle))
MixedValue is a recognized subset of Deixis.Node<MixedPayload>:
  Atom:  Some(BytesCase(bytes)), no children
  Wire:  Some(WireCase(handle)), no children
  Tuple: None, exactly dense canonical positional children, each a MixedValue
```

The slot cases are supplied by the optional profile, not built into Deixis.
There is a total ground embedding and an exact partial ground recovery on its
wire-free admitted subset. No live-capability equality, hashing, portability,
authority or runtime recognition follows from structural admission. Wire-free
arbitrary Deixis nodes need not be ground Ontos values. The denotation need not
force a particular in-memory layout on every existing language/runtime.
The profile must capture immutable atom bytes and complete immutable child
structure while retaining the opaque live handles themselves. Deixis freezes
structure, not arbitrary T: mutating a borrowed BytesCase buffer after admission
would violate a claimed ground-value preservation law. A capability's ongoing
state/lifetime is deliberately not a mutable atom or structural child.

`ontos/wire` is proposed to own that profile without depending on bitruntime or
creating a bitwire import cycle. One possible design parameterizes the opaque
capability carrier and lets runtime composition supply recognition; the concrete
native API and separation from the Wire interaction interface are not decided.

`horos/wire` would own optional type/declaration semantics over that domain,
including Wire<E> and explicit ground embedding, and reuse Horos declaration and
structural machinery without widening any existing API or old identity. Bittype
would project/analyze those definitions and execute transformations; it would not
own another canonical type algebra. The exact new envelope/profile identity is
to be specified deliberately, not a search-and-replace of tag strings.

### Preserve the explicit, ordered partial lift

Existing Bittype docs distinguish **closed embedding**, which specializes a
resolved ground expression to G and produces Data(G), from **template lifting**,
which derives an independently identified mixed-domain declaration with preserved
parameters. Neither is selected implicitly by engine route availability.

The current ordered template rules are a preservation requirement, subject to
explicit counterexamples and improvements from this review:

| Input, in order | Result |
| --- | --- |
| Parameter slot | Corresponding derived binder, with source provenance. |
| Resolved application | Application of the derived exact target to translated arguments, preserving order/arity/caller scope; no shape-based rebinding. |
| Tuple / labelled compound / typed list / closed record | Translate children, retaining labels, field names, optionality and order. Apply even when children are closed. |
| Union with two or more branches | Translate branches in order. |
| Singleton union | Translate its branch, retain union provenance. |
| Empty union | Embed the empty ground type. |
| Any other independently specializable, closed, well-formed subtree | Data(G), retaining source dependencies/provenance. |
| Otherwise | Located unsupported-lift failure, no top substitution, erased constraint or usable partial output. |

Closed here excludes free recursion binders as well as generic parameters and
requires application resolution/specialization; a token scan is insufficient.
All original imports are validated before selecting the output closure. Derived
dependencies follow preserved applications and stop at Data boundaries, so a
ground-only dependency need not be lifted. Translate declarations once by resolved
identity/configuration, group by source module, and pin output dependencies before
their importers. Retain source pins, selected roots, lifting version, output names,
origins and distinct derived identities. Reject cycles and naming conflicts.

Important counterexamples:

- Data(any) remains ground; Data(not(G)) complements only ground Ontos values.
- `not(T)`, `set-of(T)` and recursion containing free T have no template lift in
  the initial profile; their closed specializations may instead be embedded.
- `exact(value resembling param/apply/wire)` is literal data, not a type position.
- Given `Other<T> = not(T)`, `Root = not(Other<Person>)` can embed its entire closed
  body without lifting Other. `Alias = Other<Person>` preserves the top-level
  application and therefore requires the unsupported Other lift. Explicit closed
  embedding of Alias is a separate requested operation, not a fallback.

Pagination remains one authored template:

```text
Ground Page<T> = closed record { items: list-of(T), next?: scalar(utf8-text) }
Derived Page<T> = closed record { items: list-of(T), next?: Data(scalar(utf8-text)) }
```

Ground Page<Person> remains valid. The derived Page may take Wire<Message> or an
ordinary Cell type. No user authors another Page or an expanded specialized record.
Open generics/applications remain available to native generation; specialization
to a ground descriptor does not replace those identities or become a prerequisite.

The qualification obligation is denotational: for an admitted template D and
closed ground arguments A, the lifted D bound to explicit Data embeddings of A
must describe precisely the embedded instances of the specialized original,
preserving their ground layout/membership. This is not equality of source ASTs,
module pins, or Deixis/Ontos codec bytes. Source and derived identities remain distinct.
This template lifting is separate from Deixis's lifting of a supplied slot equality.

### Library-owned contracts, native generation and runtime adapters

Existing illustrative library definitions, not installed protocols:

```text
Unit = Data(H("tuple"))
Function<A,B> = Wire<tuple(A, Wire<B>)>
Cell<T> = closed record { get: Function<Unit,T>, set: Function<T,Unit> }
Step<T> = compound("item", T) | compound("end")
Stream<T> = closed record {
  next: Function<Unit,Step<T>>
  cancel: Function<Unit,Unit>
}
```

H(tag, operands...) denotes the existing `horos/<tag>` descriptor tuple. Thus
H("tuple") is the descriptor accepting the empty tuple; it is not the empty
tuple value passed where Data requires a descriptor.

Function's shape is success-only; fresh reply endpoints, matching, failure,
completion and cleanup are its library's responsibility. Cell owns mutation.
The illustrative pull Stream owns demand, order, concurrent-next behavior and
terminal/cancel races. Shapes do not prove any of these laws. Unit is an empty
tuple, not automatically native void. Future protocols/qualification remain needed.

Native recipes belong to exact library declarations and are explicitly selected
by the consuming build. Their information: source language/profile + module pin +
declaration identity; target language/syntax/representation profile; external
generic type or type-only target AST; ordered parameter mapping; locked native
imports/dependencies; supported arguments and representation limits. No raw source,
name/shape matching, last-writer-wins, silent weaker fallback or automatic runtime
installation. Structural generation is also an explicit profile choice.

This supports ordinary native generics such as Page<T> with a selected library's
Cell<T> in TypeScript/Go/Rust/C++. Native recipes do not implement adapters or prove
behavior. A future runtime adapter sketch remains a separate unimplemented API:

```go
type Adapter[T any] interface {
    ToValue(ctx BindingContext, value T) (wire.HydratedValue, error)
    FromValue(ctx BindingContext, value wire.HydratedValue) (T, error)
}
```

BindingContext here is proposed resource/lifetime composition, not the implemented
ReceivedContext. PageAdapter(CellAdapter(PersonAdapter)) is conceptual composition:
conversions may be lazy and bidirectional; wrapping a Cell is not reading a
snapshot. A dynamically received description does not create a statically usable
named Go type or select behavior by its shape. Existing Bitruntime Cell fixtures
are authored consumer evidence with their own Cell-of-Wire convention, not proof
that the above structural library declarations or universal adapters are implemented.

### Agnostic tooling constraints

Bittype already has a generic `AST[L]` arena of language-owned kinds, attributes,
containment children, non-containment references and origins. It does not depend
on Deixis for that carrier; this proposal does not replace it.

The static pipeline remains:

```text
AST[canonical profile] -> AST[native language] -> CST[native language] -> source
```

Current native AST/text profiles cover types and interfaces only. Executable
adapter-source generation needs separately qualified target syntax, checking,
formatting and native compiler coverage; the Adapter sketch does not imply that
support. A raw-source escape node is not an acceptable substitute.

Syntax admission permits unresolved names. Explicit resolution captures pinned
input/context snapshots; analysis consumes a compatible receipt and publishes
immutable request-bound facts with separate completion/verdict/model availability.
The receipt is immutable resolution evidence tied to those captured inputs; the
three result dimensions distinguish execution completion, acceptance/rejection,
and availability of semantic facts.
The kernel does not interpret Wire, Function, Cell or Stream. Concrete language
providers own semantic rules; native generation consumes explicit profiles/recipes.

Transformers expose FromID(), ToID() and Transform; configuration is immutable and
consumer-supplied. The engine rejects ambiguous shortest paths and does not retry a
weaker transformation. It rejects self-edges and permits one registered profile
per language ID, so a Horos-to-extension route needs distinct endpoint IDs. Different
Go marker types are an isolation choice, not an engine requirement; actual marker
and profile compatibility are checked. No new engine operation or per-library
registration follows from introducing the optional model.

### Explicit amendments and alternatives

The proposal changes ownership/scope, not merely spelling:

- Ontos CONTRIBUTING currently calls non-L0-expressible additions a fork and
  admits only its existing projection as an external-dependency exception.
  Its historical layered-architecture ADR says every upper layer projects totally
  to L0. A new dated amendment must distinguish the optional mixed-domain sibling
  from conservative L1/L2 data readings; preserve that ADR body and frozen specs.
- Horos current guidance says Bittype owns Wire and forbids Horos Wire altogether.
  Amend that scope only for the optional new profile. Existing ground APIs,
  declarations, supported package surfaces and meaning stay unchanged.
- Bittype's current proposed `bittype/contract/*` vocabulary must cease being the
  authoritative proposal. Preserve all design ideas with explicit new homes and
  links; do not leave two competing specifications or add compatibility aliases.

Alternatives considered fairly:

1. Keep the current Bittype-owned contract language: clear data boundary but gives
   the generator repository canonical type semantics. The proposed vocabulary also
   duplicates some structural work; moving ownership alone does not remove that
   duplication, so semantic and implementation reuse need separate justification.
2. Put Wire in Ontos L0 or Horos's existing ground checker: simpler-looking unification
   but breaks the frozen/data-only contract; outside the authorized direction.
3. Optional Ontos/Horos extensions over Deixis: selected direction, but requires
   honest profile/domain/identity boundaries and an acyclic dependency plan.
4. Keep the current runtime carrier as an optimized realization of a specified
   Deixis profile: avoids requiring a tree-allocation rewrite merely for denotation.
   This realization choice is compatible with option3, not an alternative owner.
5. General capability-extension machinery rather than a Wire-named package:
   possible factoring, but must be justified by current requirements, not speculative
   new abstractions or a centralized behavioral-contract registry.
   This is a possible factoring within the selected ownership, not a mutually
   exclusive alternative to option3.

### Prior research reuse and completion boundary

Bittype research0001 concerned a now-superseded independent value/type language.
Its lessons separate semantic identity, native representation and wire encoding;
its old normalization/identity recommendations are not current Horos policy.
Research0002 already qualified the agnostic semantic SPI and Horos provider.
We retain its snapshot/identity/completion and provider-ownership constraints,
not reopen that implemented architecture. Neither answers this optional-extension
ownership question. Ontos's older convergence panel rejects live capabilities in
the frozen foundation; it did not assess a separate optional sibling profile.
Deixis's mixed Handler|Ontos.Value survey is precedent for an opaque slot, not a
ratified profile, implemented runtime or ownership assignment.

Documentation completion must include a concept-preservation table, current versus
proposed status, exact supersession boundaries and future qualification obligations.
No implementation, dependency migration, runtime change, package publication or
deployment is authorized. Recommendations requiring those become explicit future
work with checkable triggers, rather than being represented as delivered now.

## Questions for the Expert

1. What semantic boundaries make the selected optional Ontos/Horos extensions
   coherent while retaining the frozen ground models? Identify where a proposed
   lifting/profile claim is false or requires stronger conditions.
2. How should the mixed Deixis value profile specify admission, ground recovery,
   handle aliasing and optional identity/encoding, without assuming live-capability
   equality or creating an Ontos/Wire/runtime dependency cycle?
3. How should horos/wire reuse declaration machinery and express Wire<E>, ground
   boundaries and partial lifting? Assess whether the ordered rules remain sound,
   propose any invariant-preserving amendments, and address the counterexamples
   above. Distinguish static type facts from runtime recognition and protocol evidence.
4. What explicit identity, profile, module and provenance choices preserve generics,
   library-owned target representations and the existing agnostic pipeline, without
   silently reinterpreting old declarations or creating competing semantic authorities?
5. What important idea, counterexample or simpler composition have we missed, and
   what concrete future qualification would falsify your recommended design?

## Applied

### Advice provenance and scope of evaluation

The complete response is [raw advice](0003-optional-wire-extensions.advice.md)
from attempt 2, run `run_31c1813d-dd34-46a5-9873-105a934581fb`, to document
`doc_04fefdb0-fc0a-45f6-a4dd-f6c5ef9e7fe4`. The frozen
[submitted bytes](0003-optional-wire-extensions.submitted.md) have SHA-256
`4d8a070c2984e0e6c6bdc3f06b0e91bb400954b88dc5f2f28323d40ea60dd517`.
The verified response hash is
`11d7ba8b0b218bc6e24aaec41a2353fb161bb4a417e0ff27a65c1d0a2a278a72`.
The service reported completed; `consult verify --output` returned exit 0 with
byte identity and no warnings. The response covers all five questions and has a
complete conclusion. The service labels its advisor model **6 Pro**; this is an
expert-model assessment, not a claim of human review or executed code validation.

Attempt 1's [input](0003-optional-wire-extensions.submitted.attempt-1.md) and
[partial advice](0003-optional-wire-extensions.advice.capture-1.md) are retained
unchanged. A supported recapture returned the same truncated answer, ending in
section 3.2. It is genuine evidence for its actual contents, not evidence for
missing recommendations. The revised submission adds only the bounded complete
response requirement; the reviewed technical context and questions are unchanged.
The [recovery record](0003-optional-wire-extensions.submission.json) retains both
attempts, all hashes and the first run's two capture result IDs.

Evaluation uses the six source pins in this brief, particularly Ontos
`a3f1d08`, Horos `14befc6` and Bittype `9c86864`; those remote main revisions were
reconfirmed before finalizing the documentation. Owners compared the advice to
actual ground projection, declaration/binder/specialization and registry code.
Accepted rows below deliver **documentation requirements and ownership** only.
They do not claim extension/runtime implementation or future conformance tests.

### Recommendation dispositions

Root owner `codex-hydrated-type-design` is responsible for every accepted row;
bounded Ontos/Horos reviewers supplied source verification and document edits.
The independent reviewer checks the exact candidate and delivered state.

| # | Recommendation | Verdict | Action | Owner | Link | Validation |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | Optional Ontos/Horos ownership; frozen ground cores; dated scope amendments | HOLDS | Amend contribution guidance explicitly; keep Ontos ADR0003 historical body and frozen specs; retire Bittype semantic ownership. | Root | [ontos23](https://github.com/Bitspark/ontos/issues/23), [horos17](https://github.com/Bitspark/horos/issues/17), [bittype35](https://github.com/Bitspark/bittype/issues/35) | Cross-repo consistency review; historical-body comparison; changed paths show docs only. |
| 2 | Exact mixed payload/leaf/dense-child recognition | HOLDS | Specify closed recognized profile and negative sparse/noncanonical/payload-with-child cases. | Root | [ontos23](https://github.com/Bitspark/ontos/issues/23) | Source comparison with existing Deixis projection; explicit acceptance/refusal inventory. |
| 3 | Immutable byte/child capture and repeated-handle alias preservation | HOLDS | Require immutable observations, no mutable backing exposure and Tuple(w,w) preservation without invocation/rebinding. | Root | [ontos23](https://github.com/Bitspark/ontos/issues/23) | Review catches absent payload immutability requirement; final page must state bytes, keys, children and handles separately. |
| 4 | Recovery recognizes exactly embedding image | NEEDS ADAPTATION | Apply the image law to ground recovery, not mixed admission, which also accepts Wire leaves; distinguish wire-free arbitrary Deixis trees. | Root | [ontos23](https://github.com/Bitspark/ontos/issues/23) | Counterexample: a valid Wire leaf is mixed but outside ground embedding image; sparse wire-free node is outside both. |
| 5 | Ground identity/codec preserved; no generic mixed Eq/hash/portability | HOLDS | Recover ground then use Ontos codec; distinguish handle reference, capability designation and scoped export identity. | Root | [ontos23](https://github.com/Bitspark/ontos/issues/23) | Compare frozen projection/codec laws; no assertion of Deixis-byte equality or live portable encoding. |
| 6 | Opaque carrier, acyclic dependencies and optimized runtime realization | NEEDS ADAPTATION | Keep core independent of Deixis; optional projection/value extension may depend on it. No Bitwire/Bitruntime interface import in optional value model; static type tools need no live carrier. Treat first response's linear graph as semantic illustration only. | Root | [ontos23](https://github.com/Bitspark/ontos/issues/23), [bittype35](https://github.com/Bitspark/bittype/issues/35) | Existing module boundaries and current runtime mirror remain documented; no dependencies or allocation rewrite added. |
| 7 | Separate structure, runtime recognition, operational validity and protocol evidence | HOLDS | Document judgments, send-only attenuation, stale-handle limits and current endpoint/scope ownership. | Root | [horos17](https://github.com/Bitspark/horos/issues/17), [bittype35](https://github.com/Bitspark/bittype/issues/35) | Send-shaped foreign object and expired-handle counterexamples; existing runtime source trace. |
| 8 | Reuse declarations with separate optional profile/envelope; preserve ground APIs | HOLDS | Specify optional W grammar, explicit profile/revision, canonical pins, lexical binding and unchanged ground revision1. No compatibility decoder. | Root | [horos17](https://github.com/Bitspark/horos/issues/17) | Compare current D grammar/resolver and W examples; unknown tags/profiles remain refusal obligations. |
| 9 | Data(any), complement and exact retain ground meaning | HOLDS | Give ideal membership denotation; Data is closed and opaque; exhaustion is not non-membership. | Root | [horos17](https://github.com/Bitspark/horos/issues/17), [bittype35](https://github.com/Bitspark/bittype/issues/35) | Ground checker and literal operand boundaries; no mixed top/complement implied. |
| 10 | Ordered position-sensitive partial lift and Other/Root/Alias counterexample | HOLDS | Preserve every rule, branch order, singleton/empty union, full source validation and derived closure stopping at Data; no fallback. | Root | [horos17](https://github.com/Bitspark/horos/issues/17), [bittype35](https://github.com/Bitspark/bittype/issues/35) | Independent preservation review compares old/new tables and checks opposite Root/Alias outcomes. |
| 11 | Preserve exact binder/target identity during lifting | NEEDS ADAPTATION | Preserve exact source correspondence as provenance; derived declarations/binders receive distinct identities. Resolve arguments in caller scope before binding callee slots. | Root | [horos17](https://github.com/Bitspark/horos/issues/17), [bittype35](https://github.com/Bitspark/bittype/issues/35) | Current DeclarationID/ParameterID and specialization source; distinct derived namespace and shadowing cases. |
| 12 | Denotational ground preservation and substitution qualification | HOLDS | State exact embedded-membership/layout law and lexical coherence as future proof/fixture obligations, not established theorems. Preserve optional fields, nested apps, literals and recursion counterexamples. | Root | [horos17](https://github.com/Bitspark/horos/issues/17), [bittype35](https://github.com/Bitspark/bittype/issues/35) | Review formulas and explicit limits; empty union denotes no instances, not a bottom value. |
| 13 | Distinct immutable profile compatibility and document/declaration identity | NEEDS ADAPTATION | Qualify references/recipes with semantic profile separately; keep existing ground DeclarationID(module pin,name) unchanged. Profile fixes domain, interpretation, encoding and recovery conventions; derived receipts retain roots, pins, slot mappings and naming choices. | Root | [horos17](https://github.com/Bitspark/horos/issues/17), [bittype35](https://github.com/Bitspark/bittype/issues/35) | No source struct change; canonical names remain pin-bearing, native naming/origins metadata; no shape equivalence. |
| 14 | Library-owned native recipes and explicit structural alternative | HOLDS | Retain source/target profile identity, ordered operands, external/type-only AST choices, dependency locks, errors and four-language configuration. | Root | [bittype35](https://github.com/Bitspark/bittype/issues/35) | Preservation ledger and recipe table compared to base; no raw source, implicit lookup, ambiguous or last-writer fallback. |
| 15 | Separate future adapters, dynamic descriptions and behavior | HOLDS | Preserve lazy bidirectional adapters, explicit lifetimes/cleanup, dynamic/static distinction and type-only syntax limits; keep complete library examples with protocol limitations. | Root | [bittype35](https://github.com/Bitspark/bittype/issues/35), [horos17](https://github.com/Bitspark/horos/issues/17) | Compare prior runtime-values and library equations; no executable generation or native type fabrication claimed. |
| 16 | Existing AST/CST, semantic SPI and engine stay agnostic | HOLDS | Plan optional frontend stages and distinct engine endpoints using current interfaces, immutable receipts/configuration and explicit registration. | Root | [bittype35](https://github.com/Bitspark/bittype/issues/35) | Current registry self-edge/profile rules and import gates; only documentation paths change. |
| 17 | Resource budgets and distinct exhaustion outcomes | HOLDS | Require bounded capture/recovery, resolution/lifting/generation and runtime work with node/depth/byte/module limits where applicable; no usable partial output. | Root | [ontos23](https://github.com/Bitspark/ontos/issues/23), [horos17](https://github.com/Bitspark/horos/issues/17), [bittype35](https://github.com/Bitspark/bittype/issues/35) | Independent reviewer identified missing optional-operation bounds; final docs explicitly enumerate them and distinguish invalid/unsupported/exhausted. |
| 18 | Protocol laws and adversarial falsification cases | HOLDS | Record qualification gates for authority attenuation, stale scope, union emptiness, literal opacity, pinned recipes, adapter assumptions and correct endpoint/scope cleanup. | Root | [ontos23](https://github.com/Bitspark/ontos/issues/23), [horos17](https://github.com/Bitspark/horos/issues/17), [bittype35](https://github.com/Bitspark/bittype/issues/35) | All advice falsifiers map to documented future checks; none reported as executed by this docs task. |
| 19 | Avoid premature generalized capability framework | HOLDS | Keep the concrete Wire optional profile; require a second concrete extension with demonstrated shared semantics before proposing generalized machinery. | Root | [ontos23](https://github.com/Bitspark/ontos/issues/23), [bittype35](https://github.com/Bitspark/bittype/issues/35) | No generalized registry, placeholder package or new SPI introduced. |
| 20 | Implement profiles, generators, adapters and conformance suites | DEFERRED | Trigger: a separately authorized implementation change introduces the corresponding optional model, transform or runtime binding; that change must supply the documented laws/negative fixtures and target/runtime qualification before claiming support. | Future component owner | — | No implementation authorized by this request; current diff must contain only documentation/research evidence. |

### Delivery evidence

The documentation was reviewed independently by `codex-hydrated-type-review`
(session `01a19ddb-e273-7ef1-9192-c6be8a46705f`). Its initial review found missing
extension-operation resource bounds, ambiguous canonical/native target-name
wording and ambiguous payload-with-children wording. All were corrected before
the exact candidate review; no concept-loss or source implementation change was
found. The 20-row evaluation and frozen/advice hashes were checked independently.

| Repository | Reviewed candidate | Merged main and exact validation |
| --- | --- | --- |
| bittype | `84a21c8e636d1d59b7f01dc0ae83f4b9807b2e03`; [docs](https://github.com/Bitspark/bittype/actions/runs/41689863073), [Go/native](https://github.com/Bitspark/bittype/actions/runs/41689863165) passed. | [PR36](https://github.com/Bitspark/bittype/pull/36) merged `1dee5dbb45cd8cf1d632c9151bafca56876c86b5`; [main docs](https://github.com/Bitspark/bittype/actions/runs/41690741558), [main Go/native](https://github.com/Bitspark/bittype/actions/runs/41690741317) passed. |
| horos | `f79eced7c6769db5903faff644f1c210948edbcd`; [exact dispatch](https://github.com/Bitspark/horos/actions/runs/41689878790) and [PR CI](https://github.com/Bitspark/horos/actions/runs/41689882779) passed all six lanes. | [PR18](https://github.com/Bitspark/horos/pull/18) merged `56da2464962d879225b30dd991d0e9c8b5cb6cd5`; [main CI](https://github.com/Bitspark/horos/actions/runs/41690761746) passed all six lanes. |
| ontos | `87299ae50c3487d32869be84c3b511952170847f`; [exact dispatch](https://github.com/Bitspark/ontos/actions/runs/41691040503), [PR CI](https://github.com/Bitspark/ontos/actions/runs/41690976001) and [substrate/profile](https://github.com/Bitspark/ontos/actions/runs/41690975995) passed. The exact dispatch includes consumer smoke; both CI runs include fuzz qualification. | [PR24](https://github.com/Bitspark/ontos/pull/24) merged `dbd902087113b3adda8fe27676fcecc082e98db6`; its tree is identical to the independently reviewed and remotely validated candidate. |

The candidate and merged trees match exactly: Bittype
`7fd9a10a10d56ba115624c18ade3ee962ed21282`, Horos
`c08aaa0325f92e16473b1938eb728689c2e91e61`, Ontos
`7661d0af0a89adbe802b790950458cb95fba1e54`. Implementing issues35,17 and23 are
closed. Every accepted documentation recommendation is delivered and validated;
this consultation is applied within its agreed documentation-only scope.

The initial Ontos candidate correctly failed its maintained spelling lint and
public-profile documentation-index check. Only repository-name prose and two
formula fence labels were corrected; no naming rule or gate was weakened.
The required index was regenerated on a GitHub runner through the maintained
Atlas manifest function, not edited manually:
[export run41690725951](https://github.com/Bitspark/ontos/actions/runs/41690725951),
source `f64ea767118f97e35e02d78c2968b8a652b13c72`, generator
`Bitspark/atlas-gate@978f6720e995281ffb890ea4e92293eaaaf47978`.
The retrieved artifact hash matched
`1af7d77b9b54cd903b434c8411801216e3f99ce7b0c817b5d81592e407aa01d1`;
its 13-line diff adds only the new page and its 11 headings. The temporary export
workflow is not part of any task PR or main tree. Superseded failing-head runs
were cancelled after their failure evidence was retained; the corrected candidate
passed all maintained checks anew.

The delivered documentation retains explicit future qualification requirements.
No optional value/type model, AST frontend, generator, adapter, protocol, runtime,
dependency or release was implemented. Row20 remains deferred under its stated
trigger; it is not unfinished work in this documentation-only request.
