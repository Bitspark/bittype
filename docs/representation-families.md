# Representation families

**Accepted design, 9 November 2026. Not implemented.** Decided by
[research 0015](../research-docs/0015-representation-families.md)
([#203](https://github.com/Bitspark/bittype/issues/203)). It details the selected
representation plan and the paired contributions of
[binding generation](binding-generation.md) for concrete value families.

A **representation family** is one kind of horos type together with its native
representation: integers, text, bytes, lists, maps and sets, unions, optional and
nullable fields, recursive types. Today the structural profile
`structural-representation/1` represents only text, booleans, lists, named closed
records and optional fields, and only as types.

## Delivered complete, the first time

A family is delivered with its native type, its plan entries, its value
conversion and its qualification together. It is never a types-only mapping that
a later binding step reopens. The conversion constrains which native types are
admissible at all: a TypeScript `number` cannot hold every horos integer.

**The promise.** Adding a family through the published identity, dependency,
invocation, plan-graph and target-extension contracts changes only that family
and the consumers that compose it. Existing representation identities keep their
meaning. Moving a consumer to a different native representation is a separate
migration of that consumer's contract. An explicit refusal that can later be
lifted additively is complete; replacing a delivered native type (for example
`number` by `bigint`) is not.

## The representation contract

A family contributes, per target, one contract with seven parts:

| Part | What it defines |
| --- | --- |
| Applicability | The source identities and forms it accepts, its generic assumptions, its targets, its explicit refusals |
| Native representation | Native type, dependencies, construction, access, discrimination, admissible native states |
| Value meaning | Promised canonical domain, domain predicates, native equivalence, any permitted normalization |
| Selected-use payload | A family-owned, versioned schema for the decisions recorded at each use |
| Conversion implementation | Run-time factory reference, child adapter slots, runtime compatibility, conversion laws; optionally a lowering facet (below) |
| Operational requirements | Budget use, mutation and ownership, recursion support, sharing behaviour, failures |
| Qualification | Separate evidence for type generation and binding support, tied to the representation and implementation versions |

**One instance, two projections.** The planner selects one immutable
representation instance per use. The native type and the plan entry are both
projections of that instance. No later pass reconstructs a decision from the
source, so a native declaration and its conversion cannot disagree.

**Clarified on 13 November 2026
([research 0017](../research-docs/0017-contribution-contract.md), rows 6 and
18).** Two words of this paragraph are read as
[the contribution contract](contribution-contract.md#what-a-rule-does-the-protocol)
states. "Cannot disagree" means that no sanctioned second decision point
exists: no interface lets a later pass select again or rebuild a decision. It
does not mean that a projector or a factory cannot be wrong, which
qualification tests. "Instance" has two senses on this page. Here it is the
selection at one use, which the plan records as a use entry. Under
[Conversion state](#conversion-state-sharing-and-budgets) the "selected
representation instance" that keys the memo is the applied representation
that several uses share, and not the record of one use
([identities](contribution-contract.md#identities)).

**Factories, and a lowering facet where needed.** By default a plan entry names
the family's run-time factory and gives it the family's payload as arguments. A
family may also supply a target-owned **lowering facet** that turns its recorded
decisions into typed target expressions or helpers, such as generated
accessors or constructors for Go. The binding backend owns imports, names,
placement and assembly. It never switches over families, and source-code strings
never stand in for typed target expressions. A shared interpreter is a valid
factory implementation only if its handlers are registered contributions
executing selected decisions. Generated source is never modified as an extension
mechanism.

**Run-time authority.** The representation owner owns its factories' semantic
contract and qualification: bittype's structural profile for the structural
families, a library for its exact mappings. Hosting a structural profile's
TypeScript run-time package is decided with platform `libraries/service`, which
owns the TypeScript backend; repository placement does not decide semantic
authority. The backend never decides what a structural integer, nullable field
or map means.

## The plan

| Shared envelope | Family-owned payload |
| --- | --- |
| Version axes; target and its configuration | Integer representation and admitted range |
| Checked source references and source semantic identity | Record field layout and presence behaviour |
| Root products and source-use identities | Union dispatch and construction |
| Scoped binders and actual-argument environments | Map representation and key requirements |
| Selected contribution and profile identities | Refinements and domain-predicate references |
| Typed dependency slots and recursive references | Family-specific construction and access |
| Native symbols, imports, factory references, runtime requirements | Other representation decisions |
| Diagnostic provenance | Family-specific diagnostic detail |

- **A graph.** Generic parameters and recursive adapter references are explicit
  links in the envelope. Nothing the orchestration must resolve hides inside a
  payload.
- **Payload schemas are horos types.** A payload is an ontos value. Its family
  declares the payload type in horos; bittype generates the Go and TypeScript
  types, and horos conformance validates a payload. No separate schema language
  is introduced.
- **Five version axes:** the envelope protocol; the family payload schema; the
  representation's semantic identity; generated-code/runtime compatibility; source
  semantics. An unknown required schema, semantic version or runtime requirement
  fails explicitly. Only metadata designated nonsemantic may be ignored.
- **Generic requirements from the first version.** Each generic use carries its
  scoped binder identities, native type arguments, adapter arguments and
  requirements (below).
- **Unstable until six vertical cases pass.** The envelope keeps an explicitly
  unstable version until the completeness test backend passes:
  1. recursive `ContentValue` with a shared DAG, a cycle and an expanded-size failure;
  2. `Page<Outcome>` with optional nullable fields and an exact discriminator;
  3. the exact external generic `Cell<T>`;
  4. a Go case with typed construction and a generic adapter argument;
  5. a generic map whose key requirement is rejected;
  6. an independently registered family that needs no kernel or backend change.
- **Clarified on 13 November 2026
  ([research 0017](../research-docs/0017-contribution-contract.md), row 30).**
  What stays unstable until then is the envelope's version: its layout and its
  encoding. The meanings that the plan shares with the contribution contract
  are fixed with that contract: the identities, the categories and scope of a
  request, the generic conventions and requirements, the meaning of a refusal,
  what a facet is attached to, and the relation between a native symbol and
  its realization
  ([what is final and what is provisional](contribution-contract.md#what-is-final-and-what-is-provisional)).

## Domains

Each selected representation distinguishes four sets:

| Set | Meaning |
| --- | --- |
| `H_T` | Every canonical value conforming to the horos descriptor |
| `V_T ⊆ H_T` | The application's promised domain, including identified domain predicates |
| `A_T,p ⊆ V_T` | The canonical values this representation promises to cover |
| `N_T,p` | The admissible native values |

A lossless binding has `A_T,p = V_T`; an ordinary structural one with no domain
predicates has `A_T,p = V_T = H_T`. `A_T,p ⊊ V_T` is **declared partial
coverage**, reported as such and never as full support. Resource exhaustion lies
outside these sets: a valid value stays valid when a call lacks the resources to
convert it.

- **The structural profile is lossless.** Narrowing appears only through a declared
  domain: a domain library's exact mapping, or an explicitly named bounded profile.
- **Domain predicates live in the domain library**, with a stable identity and
  version referenced by every target's binding. logos may later supply checked
  evidence; it is not a prerequisite. A domain violation and a representation
  restriction are different diagnostics.
- **A wider native type is a carrier, not a validator.** Encoding validates the
  native requirements, the domain predicate and ground conformance; decoding
  validates the canonical input and constructs an admissible native value.

## Families

**Integers** ([#163](https://github.com/Bitspark/bittype/issues/163)).
TypeScript `bigint`, Go `*big.Int`. In Go `nil` is inadmissible, encoding never
modifies its input, decoding returns an integer the caller owns, and copies use
`Int.Set`. Rust and C++ name an arbitrary-precision dependency before claiming
coverage. Rust's ontos stops at `i128` and C++ has no ontos implementation; both
are explicit qualification boundaries. A bounded native type (`number`, `int64`)
states its whole admissible native domain, including fractions, infinities and
`-0`. Non-negative alone does not justify `number`: Content's count needs a
declared bound such as `0 ≤ n ≤ 2^53 − 1`, or stays `bigint`
([bitspark-content#52](https://github.com/Bitspark/bitspark-content/issues/52)).

**Text and refinements** ([#165](https://github.com/Bitspark/bittype/issues/165)).
Adapters validate text: JavaScript strings may hold unpaired surrogates, Go
strings arbitrary bytes. Refinements use horos's meaning, such as code points,
never UTF-16 or byte lengths.

**Bytes.** Encoding into ontos atoms copies. Each family defines who owns decoded
storage and what repeated occurrences may share.

**Maps and sets** ([#165](https://github.com/Bitspark/bittype/issues/165)).

- Native maps only where native key equality is canonical equality. Otherwise use
  a canonical-map abstraction that captures a stable key representation, or an
  entry sequence with uniqueness rules.
- Go `comparable` is not evidence: a `*big.Int` compares as a pointer.
- Encoders detect duplicates after key conversion and sort by the codec order of
  the encoded keys. Decoders reject disorder.
- The TypeScript text-keyed default is `ReadonlyMap<string, V>`. Insertion order is
  not part of the value. An object form would make own properties, special
  property names and extra properties part of the family contract.

**Absence and null.** An optional nullable field has three states. Go uses
`Optional[Nullable[T]]` with defined treatment of unused payload fields; a pointer
form needs proof that its `nil` cannot collide with a valid value (`**big.Int`).
TypeScript uses `field?: T | null`, and the adapter separates an absent own
property from a present `undefined` at run time.

**Unions** ([#165](https://github.com/Bitspark/bittype/issues/165),
[#197](https://github.com/Bitspark/bittype/issues/197)).

- **The criterion is information preservation, not overlap.** The round-trip law
  implies `E(n1) = E(n2) ⇒ n1 ≈ n2`. An overlap without branch provenance binds:
  for `union(list-of(int), list-of(text))`, `D(empty) = []` and `E([]) = empty`
  satisfy both laws. A representation fails only where two native values that the
  equivalence distinguishes encode to the same canonical value.
- **A library representation cannot recover erased information.** It can share a
  native form, admit a normal form, define the equivalence so the distinction is
  irrelevant, or change the source contract so the distinction is canonical data
  (a real tag).
- **Source disjointness is not enough.** `union(list-of(int), seq(int))` is
  canonically disjoint, but with both children as `bigint[]` the native empty array
  does not say which canonical value to produce. Child representations must
  compose without losing canonical distinctions.
- **The first union family is conservative:** distinct canonical forms, distinct
  compound labels, a required discriminator field with disjoint exact values, null
  plus a type proved to exclude null, and compatible native tests and constructors.
  An overlap-supporting family may follow under coherence obligations.
- **The judgment is divided:**
  - horos semantics and bittype's source integration establish canonical facts
    (membership, disjointness, overlap witnesses, discriminators);
  - domain libraries add their predicates;
  - target families establish native discrimination and compatibility;
  - the planner links these facts;
  - the backend emits the selected strategy.

  A discriminator is not a membership proof.
- **The analysis is sound and deliberately incomplete.** Its outcomes are
  `proved`, `overlap`, `unknown` and `exhausted`, and neither of the last two
  establishes anything. Normalizing an intersection keeps declaration identity, and
  intersecting closed records is not merging their fields.

**Recursion** ([#164](https://github.com/Bitspark/bittype/issues/164)).

- **Naming** runs on the source integration's checked binder identities. A `rec`
  that forms a declaration's whole body takes the declaration's name. A lifted
  inner binder gets a deterministic name, with collision and reserved-word
  handling, and keeps its free generic parameters. The
  [naming design](naming.md#reconciliation-with-accepted-documents) of
  13 November 2026 gives a reading of that handling, which it marks as its
  own derivation and notes for #164: generation states a deterministic naming
  intent and checks the names inside one declaration, the language's Files
  stage owns collisions between top-level declarations, and the target owns
  its reserved words. The access of a lifted binder is stated by its
  contribution when that family is built.
- **Indirection** comes from the selected native layout graph. The planner may
  supply strongly connected components, the groups of mutually recursive
  declarations; each target defines its layout edges and where cycles break.
- **Ownership is part of the representation.** Rust and C++ recursive ownership
  (`Box` versus `Rc`, `unique_ptr`) is not frozen until those families' sharing
  promise is decided.
- **Recursive adapters** are built through explicit plan references
  (placeholders), separately from native value indirection.

**The remaining families** ([#165](https://github.com/Bitspark/bittype/issues/165)).

- Open records keep their extra fields or are refused.
- `any()` preserves arbitrary atoms and tuples without reinterpreting familiar
  labels.
- Closed records define what happens to extra native properties.
- Go defines whether nil and empty collections are equivalent.
- `seq` and list may share a standalone native form but must compose without
  losing their distinction.
- A validated canonical-value wrapper is a deliberate selected mapping for hard
  forms. It is never a fallback.

## Generic requirements

A native type parameter does not carry the adapter, equality, discriminator or
layout evidence needed to bind values of that type:

| Generic form | Required obligation |
| --- | --- |
| `Page<T>` | `T`'s selected adapter travels with its native type |
| `map-of(K, V)` | Compatible key semantics, or a uniform canonical-map representation |
| `union(T, text)` | Discrimination or overlap coherence for the admitted actuals |
| Nullable `T` | The null representation cannot collide with a valid `T` |
| An inner recursive binder using `T` | The free parameter and its adapter are kept when lifting |
| An exact external generic | Its library-owned factory arguments and requirements |

A generic family proves its requirements for every permitted actual, declares
requirements that actuals must satisfy, or defers qualification to checked
applications. Qualifying `U<int>` does not qualify `U<text>`. Explicit
specialization is represented as such
([#166](https://github.com/Bitspark/bittype/issues/166)), never as ordinary
substitution into an unchanged generic declaration.

## Conversion state, sharing and budgets

**Conversion state.** Active traversal detects cycles, and completed conversions
give sharing. The memo is keyed by the selected representation instance plus
input identity, including the descriptor and generic environment. An object-only
cache is unsound when one object appears under two mappings.

**Three sharing promises, distinguished:**

1. Abstract value equality ignores sharing.
2. Conversion may preserve physical DAG sharing. This is an operational guarantee,
   conditioned on the sharing the input exposes.
3. The sender's aliasing cannot cross the canonical tree encoding: `[x, x]` and
   `[x, y]` with equal `x` and `y` encode identically.

Hash-consing on decode introduces aliases the sender did not have, so it is a
visible policy.

**Budgets.** A DAG of O(k) distinct nodes can encode to Θ(2^k) occurrences, so
distinct-object and depth limits do not bound the expanded message. One call
context carries separate measures:

| Measure | Purpose |
| --- | --- |
| Conversion and validation work | Bounds computation |
| Distinct objects and allocated bytes | Bounds materialization |
| Path depth | Bounds nesting |
| Expanded canonical size | Bounds downstream processing and output |

Expanded size and height are computed by memo without expansion
([ontos#35](https://github.com/Bitspark/ontos/issues/35)). ontos owns the encoded
size, families account for their work, the service library owns the call context,
and Content owns its limits.

**Capture** stays with platform `libraries/service`. Families declare what they
copy, retain and require to stay stable while reading. Copying atom bytes protects
those bytes; it does not snapshot a mutable composite graph.

## Failures and conformance

- **Through `diagnose`.** Value facets use horos `diagnose`, never the boolean
  `conforms`. The boolean interface reports an exhausted check as false, and
  `not` can turn that into true. `diagnose` preserves exhaustion under `not` and
  in unions (`core/ts/src/diagnose.ts:114-137` in horos).
- **Distinct failures.** Unsupported use, missing facet, incompatible plan or
  runtime, invalid native value, structural mismatch, domain violation, declared
  coverage restriction, cycle and exhaustion stay distinct. Shared code never
  reads a conversion failure as a failed membership test.

## Qualification

Each family and target is qualified for bindings separately from types
([#197](https://github.com/Bitspark/bittype/issues/197)):

- compiled types in the target configuration;
- both laws over the declared domains;
- boundary values and malformed input, with classified failures;
- substitution, aliases and recursion;
- source conformance checked independently of successful native construction;
- cross-target canonical agreement;
- additivity: an unrelated contribution leaves existing selections unchanged.

Deliberately broken rules show that the tests detect a failure. They do not prove
the laws, so each family keeps its law-level argument beside its tests.

## Delivery

The order of work, the first slice and the owner of each vertical case are in [#206](https://github.com/Bitspark/bittype/issues/206). [#165](https://github.com/Bitspark/bittype/issues/165) is the index of the family issues.

**Changes before the additive work (11 November 2026).** Two existing parts change once, ahead of the families: the structural profile's identity and the closed profile gate ([#216](https://github.com/Bitspark/bittype/issues/216), [bittype/ACX-008](architecture-cases/0008-structural-profile-gate.md)), and what an existing mapping returns (item 6 of [#196](https://github.com/Bitspark/bittype/issues/196)). After them a family adds a contribution and changes no existing contract. A third change moves code and changes no contract: one package per mapping ([#222](https://github.com/Bitspark/bittype/issues/222), [bittype/ACX-009](architecture-cases/0009-package-boundaries.md), [package boundaries](package-boundaries.md)). The four pair packages are split on `main` since 11 November 2026: each mapping is `transform/syntax/horos-to-<native>/mapping/<name>/go` and exports `Library`. The common package `transform/syntax/horos/mapping/go` was divided the same day; its adapter stays at that path. Item 6 of #196 is unchanged.

The first of the two is done on 12 November 2026 (#216). The structural profile's identity is `structural-representation/1`, and the declared, unimplemented enforcement profile is `structural-enforcement/1`. A mapping states its profile in its own description, and its constructor takes none. No shared code lists, writes or refuses a profile, and a request selects the libraries that state the profile it names. No mapping states the enforcement profile. A bounded profile gets its own name beside these, as a contribution that states it, and needs no change to the common package. It can use the parts that state no profile: the function type, the common patterns, the naming function and the common declaration recipe. It reuses no builtin mapping, because each builtin mapping states the structural profile.

**Clarified on 13 November 2026 ([research 0017](../research-docs/0017-contribution-contract.md), rows 2 and 33).** "Change once" is a bounded commitment: one migration to the demonstrated extension contract for the agreed extension domain, before independently owned consumers rely on it. [The contribution contract](contribution-contract.md#the-extension-domain) writes that domain down, and it gives the order in which item 6 of #196 is delivered: a design, a prototype, and then the migration behind a gate. In the table below, "generics" among the delivered families names a division of responsibility between the kernel, the contribution that supplies a declaration and the selected representation contribution. It is no family with a contribution of its own ([what a contribution is](contribution-contract.md#what-a-contribution-is)).

| Work | Owner | Tracking |
| --- | --- | --- |
| The structural profile's identity, before the first plan entry records it. Done on 12 November 2026: `structural-representation/1` | bittype | [#216](https://github.com/Bitspark/bittype/issues/216) |
| One package per mapping contribution, before a family is written into a mapping package. Done on `main` on 11 November 2026 for the four pair packages and the common package | bittype | [#222](https://github.com/Bitspark/bittype/issues/222) |
| The plan envelope, payload schemas, version axes, generic requirements, six vertical cases | bittype | [#195](https://github.com/Bitspark/bittype/issues/195) |
| The representation contract, lowering facet, operational requirements, failure classes | bittype | [#196](https://github.com/Bitspark/bittype/issues/196) |
| Binding qualification, partial coverage, the information-loss criterion | bittype | [#197](https://github.com/Bitspark/bittype/issues/197) |
| Integers | bittype | [#163](https://github.com/Bitspark/bittype/issues/163) |
| Recursion | bittype | [#164](https://github.com/Bitspark/bittype/issues/164) |
| The delivered families (text, booleans, lists, closed records, optional fields, generics), completed with payloads, conversions and qualification | bittype | [#207](https://github.com/Bitspark/bittype/issues/207) |
| Bytes and atoms | bittype | [#208](https://github.com/Bitspark/bittype/issues/208) |
| Maps and sets | bittype | [#209](https://github.com/Bitspark/bittype/issues/209) |
| Unions | bittype | [#210](https://github.com/Bitspark/bittype/issues/210) |
| Nullable values and optional nullable fields | bittype | [#211](https://github.com/Bitspark/bittype/issues/211) |
| Exact forms and text refinements | bittype | [#212](https://github.com/Bitspark/bittype/issues/212) |
| Tuples, bare sequences and labelled compounds | bittype | [#213](https://github.com/Bitspark/bittype/issues/213) |
| Open records, `any()`, intersections and complements | bittype | [#214](https://github.com/Bitspark/bittype/issues/214) |
| Explicit specialization | bittype | [#166](https://github.com/Bitspark/bittype/issues/166) |
| The three contracts: service root, capture and lifetime, outcomes versus invocation uncertainty | platform `libraries/service` | [bitspark-service#37](https://github.com/Bitspark/bitspark-service/issues/37) |
| Hosting of the structural profile's TypeScript run-time package | platform `libraries/service` | [bitspark-service#38](https://github.com/Bitspark/bitspark-service/issues/38) |
| `Function<A, B, E>` and `Service<Ops>` declared in horos, with their paired contributions | platform `libraries/service` | [bitspark-service#39](https://github.com/Bitspark/bitspark-service/issues/39) |
| The TypeScript backend: factory composition, recursive adapters, conversion state, the call context | platform `libraries/service` | [bitspark-service#40](https://github.com/Bitspark/bitspark-service/issues/40) |
| The expanded canonical size of a shared value | ontos | [ontos#35](https://github.com/Bitspark/ontos/issues/35), delivered in ontos v0.17.0 on 10 November 2026. bittype pins v0.16.0 until [#164](https://github.com/Bitspark/bittype/issues/164) needs the measure |
| The domain of `check().objects` | Content | [bitspark-content#52](https://github.com/Bitspark/bitspark-content/issues/52) |
| Content's contract in horos | Content | [bitspark-content#53](https://github.com/Bitspark/bitspark-content/issues/53) |
| Content's exact mappings and the pilot's qualification | Content | [bitspark-content#54](https://github.com/Bitspark/bitspark-content/issues/54) |
