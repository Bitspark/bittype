# Representation families

**Accepted design, 9 November 2026. Not implemented.** Decided by
[research 0015](../research-docs/0015-representation-families.md)
([#203](https://github.com/Bitspark/bittype/issues/203)). It details the selected
representation plan and the paired contributions of
[binding generation](binding-generation.md) for concrete value families.

A **representation family** is one kind of horos type together with its native
representation: integers, text, bytes, lists, maps and sets, unions, optional and
nullable fields, recursive types. Today the structural profile
`pagination-representation/1` represents only text, booleans, lists, named closed
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
  handling, and keeps its free generic parameters.
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

| Work | Owner | Tracking |
| --- | --- | --- |
| The plan envelope, payload schemas, version axes, generic requirements, six vertical cases | bittype | [#195](https://github.com/Bitspark/bittype/issues/195) |
| The representation contract, lowering facet, operational requirements, failure classes | bittype | [#196](https://github.com/Bitspark/bittype/issues/196) |
| Binding qualification, partial coverage, the information-loss criterion | bittype | [#197](https://github.com/Bitspark/bittype/issues/197) |
| Integers | bittype | [#163](https://github.com/Bitspark/bittype/issues/163) |
| Recursion | bittype | [#164](https://github.com/Bitspark/bittype/issues/164) |
| Text, maps, unions, nullable and the remaining families | bittype | [#165](https://github.com/Bitspark/bittype/issues/165) |
| Explicit specialization | bittype | [#166](https://github.com/Bitspark/bittype/issues/166) |
| The TypeScript backend: factory composition, recursive adapters, conversion state, the call context, the structural run-time package's hosting | platform `libraries/service` | [bitspark-service#35](https://github.com/Bitspark/bitspark-service/issues/35) |
| The expanded canonical size of a shared value | ontos | [ontos#35](https://github.com/Bitspark/ontos/issues/35) |
| The domain of `check().objects` | Content | [bitspark-content#52](https://github.com/Bitspark/bitspark-content/issues/52) |
