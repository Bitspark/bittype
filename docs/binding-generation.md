# Binding generation for horos-declared contracts

**Accepted direction, 9 November 2026, from [research 0014](../research-docs/0014-binding-generation.md). Not implemented.** This page records who owns what when service contracts declared in horos are turned into running bindings, and what the interfaces between those owners must carry. Type generation is unchanged. Each piece of work is tracked by its owner (see [delivery](#delivery)).

## What binding generation is

From a horos contract, a build needs more than native types:

- conversion between native values and hydrated message values;
- conversion between native values and JSON for presentation (a command line, web forms);
- operation data: paths, parameters, summaries;
- shape checks of received values;
- declared operation failures.

That set is **binding generation**. [Runtime values](runtime-values.md#6-generation-and-runtime-remain-separate-composition-layers) describes it as a separate composition layer above type generation.

## The decision

> One semantic owner does not require one emitter. It requires an authoritative decision, an explicit interface carrying that decision, and tests that establish that its consumers preserve the promised meaning.

**bittype** remains the authority for selecting native representations. It publishes each selection as a typed **selected representation plan** ([#195](https://github.com/Bitspark/bittype/issues/195)).

**A binding backend per target** consumes three inputs:

- the plan;
- the source semantics;
- each library's **paired contribution** ([#196](https://github.com/Bitspark/bittype/issues/196)).

From these it emits static composition and typed glue.

**bittype does not emit runtime adapters or own their behaviour.** Its charter gains only the plan interface and the hooks that produce it.

| Owner | Responsibility |
| --- | --- |
| horos | Declaration semantics, checked identities and applications, descriptors, ground conformance |
| bittype planning | Exact contribution selection, substitutions, references, and the selected representation plan |
| Each representation implementation (a target's structural profile, or a library's exact mapping) | Native layout, presence rules, numeric representation, construction and access |
| Binding backend, one per target | Adapter composition, runtime imports, generated glue, assembly of the requested binding products |
| Contract or domain library | Its native recipe, its adapter implementation, its protocol or domain meaning, and the laws below |
| bitruntime | Capability recognition, export and import identity, lifetime, admission, transport |
| Consumer build | Selection and composition of the concrete contributions it uses |

**Placement.** Each target's backend sits with the owner of that target's runtime adapter contract:

- **TypeScript.** The backend belongs to platform's `libraries/service` (Bitspark/bitspark-service) and consumes a serialized plan. Whether that is its contract-generator or a new package is that repository's decision. platform agreed to this on 9 November 2026.
- **Go.** A Go backend belongs to the future Go runtime library.

Distribution follows from this interface. Neither an existing generator nor the dormant bitlink repository decides it.

## The selected representation plan

**A manifest of names is not enough.** For each use and generic instance, the plan carries:

- the checked source identity and the actual generic arguments;
- the selected contribution, target and representation profile;
- native symbols, distinguishing run-time imports from type-only imports;
- field access and construction;
- the rules for absence, presence, null and native `undefined`;
- tuple layout and calling convention;
- union discrimination and construction;
- numeric range and exactness;
- recursive references;
- the selected adapter factories, their argument slots and their resource requirements;
- the selected root contract, such as `Service<Ops>`, which the author declares at the root of the contract and which is selected like any other contribution.

The serialized plan carries an explicit format version, which bittype owns. A backend consumes a pinned version and refuses an unknown one.

**Form.** These are typed references and target-owned expressions, the same principle as bittype's native intent (the [generation SPI](generation-spi.md)). They are not a universal schema and not another type algebra.

**Serialization.** The plan starts as a typed build-time interface and is serialized only where a real process boundary needs it. A TypeScript backend in Node is such a boundary.

**Completeness test.** The backend can implement each selected representation from the plan, the source semantics and the selected contributions alone, without guessing a rule. If it must decide how Go optional values work, or how a custom union is represented, the plan is incomplete.

**What the backend emits:**

- static imports and factory calls;
- ground-descriptor constants;
- operation and presentation metadata;
- target-specific glue where needed.

**Interpreters.** A shared structural interpreter is allowed only if it executes the plan's decisions rather than choosing representation rules itself.

**The emitter's scope.** It is scoped separately, with its own syntax and compiler checks. bittype's native vocabularies stay types-only.

## Paired contributions

A library publishes, for one declaration and one supported target representation:

| Element | Meaning |
| --- | --- |
| Source selector | Exact module digest and declaration identity |
| Representation context | Target and supported profile |
| Native facet | Its bittype recipe or native intent |
| Value adapter facet | Run-time factory reference and typed argument-adapter slots, converting the value inside a message |
| Root (attachment) facet | Serving the value as one routed endpoint; Function's is `functionBinding`, its value facet `functionAdapter` |
| Requirements | Adapter-interface compatibility, resource facilities, supported scope |
| Optional product facets | For example an operation facet, through which a callable contract exposes its argument, result and failure positions |

**Selection.** A native facet and the facets a binding product needs are selected together. A native mapping without an adapter facet still supports types, but a request for bindings that needs a missing facet fails explicitly. There is never a fallback to structural conversion, in the same way that an exact use mapping never silently becomes a structural publication today.

**Ownership.** Contributions belong to their library's owner, outside canonical horos bytes. Native package paths and export names are not part of declaration identity.

**Composition:**

1. resolve through the checked graph;
2. preserve associations through aliases;
3. select by exact identity before any structural expansion;
4. bind slots by checked parameter identity;
5. reject duplicates, incompatible profiles, missing facets and unsupported applications.

**Supplying contributions.** A consumer supplies its contributions explicitly. Shared machinery has no preferred Function, Stream or Cell entries.

**Operations** are discovered through the same boundary:

- General traversal knows record fields.
- Only a contribution's operation facet says that a field is an operation.
- A mapped Cell is never flattened into operations unless a contribution supplies that reading.

## Laws every contribution must satisfy

For a type T and representation profile p, with:

- `V_T`, the admissible abstract values, including library-owned predicates;
- `N_T,p`, the admissible native values;
- E for encoding and D for decoding;

the contribution must satisfy:

```text
D(E(n)) ≈ n    for every n in N_T,p   (the equivalence states whether copying, class identity, branch tags and sharing matter)
E(D(v)) = v    for every v in V_T     (abstract-value equality; any normalization is specified)
```

It must also establish:

| Law | Establishes |
| --- | --- |
| Conformance | Produced values inhabit the declared descriptor where ground conformance applies |
| Semantic validity | Domain predicates survive both directions |
| Coverage | Every promised abstract value is representable, and restrictions are explicit |
| Capture | Once captured, later mutation of caller-owned data cannot change the message |
| Generic composition | Substitution, aliases, nested applications and recursion preserve the laws |
| Resource atomicity | Failed preparation releases newly acquired resources and leaves existing ones intact |
| Capability preservation | Authority, identity and reuse rules, scope and lifetime remain valid |
| Cross-target agreement | Different native representations agree on abstract values and protocol observations |

**Coverage** is never weakened to "what this implementation happens to accept". [Research 0015](../research-docs/0015-representation-families.md) refines `V_T` into four sets and makes narrower coverage a declared, reported restriction: see [representation families](representation-families.md#domains).

**Capture in practice.** Both ontos implementations copy bytes when an atom is built (Go `NewAtom`, the TypeScript `Atom` constructor). Conversion into ontos atoms therefore discharges capture for byte buffers. Adapters that keep native references, views or capabilities must discharge it themselves.

**Capability round trips** are behavioural equivalence under stated assumptions about connection state, lifetime, cancellation and observable failures. Rollback undoes only newly prepared exports. Conversion never invokes a protocol to decide how to encode a value.

**Binding support is qualified separately from type support** ([#197](https://github.com/Bitspark/bittype/issues/197)). A representation must keep enough information for the laws to be possible:

- Two native branches can encode to the same canonical value, as with an overlapping union.
- Two abstract alternatives can collapse into one native form.

**Corrected by research 0015.** The obstruction is information loss, not overlap. The round-trip law implies `E(n1) = E(n2) ⇒ n1 ≈ n2`, so a representation fails only where two native values that its equivalence distinguishes encode to the same canonical value. An overlap without branch provenance binds. A library-owned representation cannot recover erased information. It can share a native form, admit a normal form, define the equivalence, or change the source contract so the distinction becomes canonical data. See [unions](representation-families.md#families). "Every valid type has an automatic adapter" is not a promise this design makes.

## Message form and validation

**Message values follow the horos descriptor** as canonical ontos data. A custom native representation does not permit an unrelated abstract value on the wire; a compact form must be a declared protocol representation or a named projection.

**Validation is divided by judgment:**

| Judgment | Owner |
| --- | --- |
| Whether a value is ground or mixed | ontos/wire facilities |
| Ground structural conformance | horos, against descriptors emitted at build time (prefer `diagnose`: exhaustion is not a mismatch) |
| Native representation validity and construction | the structural or custom adapter |
| Additional domain predicates | the domain library |
| Capability recognition, authority, lifetime | the runtime |
| Protocol messages, state, replies, cancellation | the contract library |

**Mixed values** also need the structure around capability leaves checked. A mixed checker needs a specified owner consistent with horos and ontos/wire ([ontos#33](https://github.com/Bitspark/ontos/issues/33)). Until one is qualified, unsupported mixed constructs are rejected.

**Resource accounting.** Validation and conversion share one budget: check maximal ground regions, and never give each child a fresh budget.

**"Declared once" means structure.** Domain laws that horos cannot express stay with the model library. For example, `ContentHash` is text plus a content-model predicate, and a non-negative count is an integer plus a predicate ([horos#36](https://github.com/Bitspark/horos/issues/36) asks about integer refinements).

## Presentation and documentation

**JSON is a named projection P of abstract ground values:**

```text
toJSON   = P ∘ E
fromJSON = D ∘ P⁻¹
```

Libraries implement native↔abstract conversion once, and each format implements its projection once. `ontos-data-json/1` is one-directional and has no bytes or sets, so a lossless, reversible profile is proposed to ontos ([ontos#33](https://github.com/Bitspark/ontos/issues/33)). The projection is a pure ground conversion and never exports endpoints.

**Presentation conveniences** are library contributions over the abstract or projection interface: entry widgets and hints. They cannot establish semantic constraints.

**Unsupported operations** are reported one operation at a time. Presentation consumes the build's operation projection; it does not infer from native types.

**Documentation** lives in one sidecar outside canonical horos identity ([#198](https://github.com/Bitspark/bittype/issues/198)):

- It is keyed by declaration identity, member path and parameter identity, with use-site overlays.
- It defines precedence between defaults and overlays, and diagnoses stale references.
- Doc comments in bittype's readable syntax are one way to author it. Prose edits never change a module digest.

## Failures and the service root

**`Function<A, B, E>`** is designed by the Function library ([bitspark-service#35](https://github.com/Bitspark/bitspark-service/issues/35)):

- A completed answer is `Ok(B) | DomainError(E)`, and `E` must be ground. Function's declaration states that requirement.
- The full reply protocol is Function's: `ok(B) | error(E) | protocol(message) | unknown(reason, message)`. `protocol` means the peer refused before invoking. `unknown` means a failure after invocation, an unusable result, or an answer over the limit.
- Invocation failures travel on a separate channel: local non-admission, an unknown outcome after admission, an undecodable reply, ended observation. Function supplies no remote cancellation.
- A provider returns a domain error as a value. A thrown exception is `unknown`, and nothing is recognised by a `code` or `family` property. The native representation of `E` is the domain's own contribution.
- The protocol is incompatible with today's. libraries/service introduces it under new profiles with a coordinated cutover. Decided by libraries/service in [bitspark-service#35, comment 6729151207](https://github.com/Bitspark/bitspark-service/issues/35#issuecomment-6729151207).

**The service root is a declared root contract.** A service declared as a record of Function values is a different contract from a single capability that dispatches by path. The advice's default was to export the declared record, with a dispatcher allowed only as an explicit library-owned root protocol. libraries/service took that alternative ([bitspark-service#35, comment 6729151207](https://github.com/Bitspark/bitspark-service/issues/35#issuecomment-6729151207)), for four reasons:

- a routed location attaches exactly one endpoint, so a record would still need a bootstrap root protocol;
- per-operation exports multiply endpoints and lifetimes;
- path dispatch lets a provider add operations without changing what callers decode;
- authority is decided per request, not by holding an endpoint.

So `@bitspark/service-contracts` declares `Service<Ops>`, whose `Ops` is a record, with nested records as namespaces, of `Function<A, B, E>` leaves. The service author writes `Service<…>` at the root of the contract. The plan selects it like any other contribution, through its root facet. Each operation's argument, result and failure positions come from Function's operation facet, so outcomes stay Function's. Service rejects any other leaf at composition. The generator never infers a dispatcher, and a plain record of Functions stays a record value.

## Order of work and qualification

Before any interface is committed, four short contracts are written:

| Contract | Owner |
| --- | --- |
| Native representation laws (this page) | bittype |
| Service-root protocol | platform `libraries/service` |
| Capture and lifetime model | platform `libraries/service` |
| Domain outcomes versus invocation uncertainty | platform `libraries/service` |

They are then challenged with:

- Content, as the pilot;
- an independent second contract library with distinct behaviour, through aliases and nested generics;
- a small Go interoperability slice.

**Content qualifies all six operations,** together with:

- every binding product;
- its existing guarantees: hash validation, non-negative counts, snapshots, limits, rejection of cycles and duplicate keys, the intended treatment of sharing.

A `put`/`get` slice is not qualification of the whole.

**The falsification tests:**

- mutation after a call;
- fidelity for absence, null, unit, bytes, large integers and recursion;
- overlapping unions;
- Go-produced descriptors through the TypeScript checker ([horos#36](https://github.com/Bitspark/horos/issues/36));
- rollback with a retained export;
- a disconnect or corrupted reply;
- duplicate and late replies;
- a changed module digest with the same name and shape;
- presentation beside an unsupported capability operation.

**Additivity** is then tested with shared machinery frozen:

| Addition | What it may change |
| --- | --- |
| A second contract library | Only its declarations, contributions, adapters, tests and consumer imports |
| A second service | Only its contract, implementation, documentation and composition |
| Go bindings | Only the Go backend and the Go implementations each semantic owner supplies |

The unit of implementation is library × target.

## Delivery

| Work | Owner | Tracking |
| --- | --- | --- |
| Selected representation plan, serialized form, completeness test | bittype | [#195](https://github.com/Bitspark/bittype/issues/195) |
| Paired native and adapter facets, operation facet | bittype | [#196](https://github.com/Bitspark/bittype/issues/196) |
| Binding support reported separately from type support | bittype | [#197](https://github.com/Bitspark/bittype/issues/197) |
| Documentation sidecar | bittype | [#198](https://github.com/Bitspark/bittype/issues/198) |
| Service root, `Function<A, B, E>`, ownership and lifetime contract, per-operation presentation, TypeScript backend | platform `libraries/service` | [bitspark-service#35](https://github.com/Bitspark/bitspark-service/issues/35), with [bitspark-service#37](https://github.com/Bitspark/bitspark-service/issues/37), [bitspark-service#39](https://github.com/Bitspark/bitspark-service/issues/39) and [bitspark-service#40](https://github.com/Bitspark/bitspark-service/issues/40) |
| Lossless JSON projection; mixed-structure checking | ontos | [ontos#33](https://github.com/Bitspark/ontos/issues/33) |
| Go-descriptor through TypeScript-checker test; integer refinements (optional) | horos | [horos#36](https://github.com/Bitspark/horos/issues/36) |
| Representation families the services need, each delivered complete ([representation families](representation-families.md), research 0015) | bittype | [#163](https://github.com/Bitspark/bittype/issues/163), [#164](https://github.com/Bitspark/bittype/issues/164), [#165](https://github.com/Bitspark/bittype/issues/165) (the index of the family issues), [#166](https://github.com/Bitspark/bittype/issues/166); ordered in [#206](https://github.com/Bitspark/bittype/issues/206) |
