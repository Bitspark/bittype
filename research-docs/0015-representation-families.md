# Research: Native representation families for a type language's values, designed together with their conversions

**ID:** 0015
**Date:** 9 November 2026
**Status:** reviewed
**Reviewed:** https://github.com/Bitspark/bittype/issues/203#issuecomment-6748762823
**Owner:** claude-bittype-gaps
**Issue:** https://github.com/Bitspark/bittype/issues/203
**Code baseline:** bittype 4fc859e2882efdc51b2991ebb7e43b5c4d6c515f; horos 3a51a06e6db1800a39f94782213609a2d8723d10; ontos 0d5f34e; platform repository `libraries/service` 24fedd4; platform repository `services/content` eaae493; services `cubes` bed6164, `drafts` e1811f5, `pointers` e2ab221

## Question

**horos** is a declarative type language whose types are themselves values of a
shared value model, **ontos**. **bittype** is a Go library that reads horos
declarations and generates native type declarations for TypeScript, Go, Rust
and C++. A software platform will declare its network services' contracts in
horos. From each contract a build produces native types, and **bindings**: code
that converts native values to and from ontos values.

An earlier consultation (research 0014) decided who owns which part of that
build (section 2.3):

- bittype selects the native representation of every type and publishes those
  choices as a typed **selected representation plan**;
- a **binding backend** per target language reads the plan and emits the
  conversion code, choosing no representation rule itself;
- message values are canonical ontos data conforming to the horos type.

Today bittype represents only text, booleans, lists, named closed records with
optional fields, and generics. The pilot service needs integers, bytes, maps,
unions and recursive types. The other services also need nullable fields,
tagged unions and a text-or-integer union.

**How should a representation family be defined so that it is complete the
first time?** For each target, a complete family fixes four things:

- its native type;
- what the plan records about each use of it;
- how its values convert, and the laws that conversion obeys;
- what it refuses, with binding support qualified separately from type support.

Families must compose (generics, nesting, recursion) and stay additive.

**What "complete the first time" means here.** The owner asked for this work to
be done properly, without technical debt that forces a refactor later
(9 November 2026). We count as a refactor any later change to one of these:

- the plan's shared envelope or its version semantics;
- shared planning code;
- a delivered family's native type, plan entries or conversion.

Adding a new family later, or supporting a form that is explicitly refused
today, is additive and not a refactor. So a family may refuse cases, provided
the refusal is explicit and support can be added without touching delivered
parts. Replacing `number` with `bigint` after release would be a refactor.

That is why the obvious order fails. If native types came first and conversions
later, the conversion would constrain which native type is admissible at all.
A TypeScript `number` cannot hold every horos integer, and a JavaScript `Map`
compares object keys by reference.

**The tensions**, each developed in section 6:

- **Integers.** horos integers are unbounded, and ergonomic native integers are not.
- **Unions.** horos unions are untagged and often overlap; native code wants tags.
- **Maps.** Keys are ordered by their canonical encoding; native maps use native equality and order.
- **Recursion.** It is anonymous in horos; native targets need names and sometimes indirection.
- **Mutation.** Native buffers and big integers are mutable; message values are not.
- **The plan's form.** It must carry all of this without becoming a universal schema.

Sections 1 to 3 give the context. Section 4 follows two values through every
layer. Section 5 lists what is fixed. Section 6 gives the options and our
present leanings, none decided. Section 7 holds five questions.

## 1. System overview

### 1.1 The components

| Component | What it is | Language(s) |
| --- | --- | --- |
| **ontos** | The value model. A value is an **atom** (a byte string) or a **tuple** (an ordered sequence of values). On that floor it defines canonical data forms (integers, text, booleans, lists, maps, sets, null, decimals) as labelled tuples, and one canonical binary encoding. | Go, Rust, TypeScript; no C++ implementation |
| **horos** | A type language whose types are ontos values called **descriptors**. It has a conformance checker for values and a declaration language with modules, generic declarations and imports pinned by content digest. bittype supplies its readable text syntax. | Descriptors and conformance: Go, Rust, TypeScript. Declarations: Go only |
| **bittype** | A Go library: horos syntax and analysis, and native type generation for TypeScript, Go, Rust and C++ through a public generation interface. Its native syntax covers type declarations only, not executable code. | Go; emits TypeScript, Go, Rust, C++ |
| **bitwire, bitruntime** | The message layer. Values in transit are ontos atoms and tuples, plus live references to remote endpoints. The platform's code calls these in-memory message values **hydrated** values, and conversion to them `toWire`. A horos type may contain such references, written `wire<E>`. This document concerns values without them. | Go, TypeScript |
| **platform repository, `libraries/service`** | The TypeScript service library: today's converters and call protocol, and the home of the future TypeScript binding backend. | TypeScript |
| **platform services** | Content, Cubes, Pointers and Drafts. Content is the pilot. | TypeScript |

### 1.2 Terms used throughout

| Term | Meaning in this document |
| --- | --- |
| **Native value / native type** | A value or type of a target language, for example a TypeScript `bigint` or `interface`. |
| **Canonical value** | An ontos value in the data form that a horos type describes, for example `Tuple(Atom("int"), Atom(0x00 0x05))` for the integer 5. |
| **Descriptor, conformance** | A horos type written as an ontos value; `conforms(value, descriptor)` checks a value against it. **Ground** means containing no endpoint references. |
| **Source, target** | In generation, the **source** is the horos side (declarations and their checked semantics) and the **target** is the native language. |
| **Use, actual, binder** | A **use** is one occurrence of a type expression inside a declaration. An **actual** is an argument supplied when a generic declaration is applied. A **binder** is a declared type parameter. |
| **The planner** | bittype's family-agnostic code that walks the source graph and selects representations (section 3.3). bittype's documents also call this shared machinery the **kernel**. |
| **Source integration** | bittype's horos-specific reader that interprets horos facts for the planner and judges whether each use is supported. It is bittype code. horos itself owns the meaning of membership. |
| **Representation family** | One kind of horos type together with how it is represented natively: integers, bytes, lists, maps and so on. |
| **Representation profile** | A named policy deciding how horos forms become native types in one target, for example `pagination-representation/1`. |
| **Mapping library, contribution** | A bittype contribution that supplies native representations. Its **use rules** (Go type `BindingDescription`, unrelated to value bindings) cover uses of a form, and its **recipes** cover declarations. A library selects by **exact identity**, never by name or shape. |
| **Structural profile vs exact mapping** | A **structural profile** maps horos's built-in forms (records, lists, text) for one target; `pagination-representation/1` is the current one. An **exact mapping** maps one particular declaration, identified by module digest and name, to a library-owned native type. Both use the same public mechanism. |
| **Selected representation plan** | Every representation decision bittype made for one generation, in typed form, published for binding backends. It does not exist yet. |
| **Facet** | One part of a contribution for one target. The **native facet** produces the native type. The **value facet** converts values; 0014's documents also call it the **adapter facet**. An **adapter** is the run-time implementation a value facet refers to, usually built by a **run-time factory**. |
| **Binding backend** | The per-target component that reads the plan and emits conversion code. |
| **Encode E, decode D** | The value facet's two directions: `E(native) → canonical` and `D(canonical) → native`. |
| **Value domain; admissible** | The set of values a representation can hold. A native value or canonical value is **admissible** if it lies in the domain the representation declares. A representation **narrows** if some valid horos value has no native form. |
| **Message capture** | The converter taking ownership of caller data, so that later mutation of the caller's objects cannot change a message already produced. (Distinct from bittype's generation-time *slot capture* in section 3.3.) |
| **Domain library** | A library owned by an application domain, such as Content's model package, which may add predicates (for example "non-negative") that horos cannot express. |
| **Binding support** | A family is binding-supported for a target when its conversion laws are tested, separately from being able to generate its type. |

### 1.3 Scale and direction

- **Consumers.** Four TypeScript services with 6, 6, 7 and 10 operations,
  running on Node.js 24 and in browsers. bittype's TypeScript checker targets
  ES2020, which has `bigint`.
- **Families.** We estimate ten to twelve families across the four services
  (sections 3.6 to 3.8).
- **Value sizes.** Content's default limits are 10,000 objects, 100,000 steps,
  64 MiB of bytes and depth 256. Today's call-wide conversion budget is 100,000
  nodes, depth 256 and 16 MiB. Conversion cost matters, but no service is
  throughput-bound today.
- **People.** One human owner decides for all the repositories involved. Each
  repository is maintained by AI coding-agent sessions, several of which work in
  parallel. Additive composition with clear ownership is therefore a hard
  requirement.
- **Targets.** TypeScript bindings now; Go bindings expected. bittype also emits
  Rust and C++ types. Bindings for those are not planned, but the design must
  not preclude them.

## 2. Architectural context

### 2.1 The conversion layers

```text
 native value (e.g. a TypeScript object of type ContentValue)
      │   ▲
      │   │   value facets, composed by the binding backend     ← the subject of this document
      ▼   │
 canonical ontos value, conforming to the horos type
      │   ▲
      │   │   bitruntime and ontos's canonical encoding (existing)
      ▼   │
 bytes on the network
```

### 2.2 Standing architecture rules (AC = architectural correctness)

The owner has made architecture part of correctness. The rules this document
relies on:

| Rule | Requirement |
| --- | --- |
| AC-02 | Meanings remain with their semantic owners. Independently owned implementations supply their domain interpretations, policies and laws; shared mechanisms own their invariant guarantees. |
| AC-03 | Extensions compose additively. Adding a contribution changes only that contribution and the consumers that compose it, not shared orchestration, a central list of cases, or unrelated implementations. |
| AC-04 | Built-in and first contributions use the same public extension mechanism as independent ones. |
| AC-05 | Composition preserves guarantees: identity, scope, ownership, authority, assumptions, resource limits and failure meaning. |
| AC-07 | A test or a compiler's acceptance supports only the claim it establishes. |
| AC-09 | Design interfaces from the domain, its owners and the guarantees consumers need. The existing API and the effort of changing it must never influence the design, not even as a tiebreaker. |

### 2.3 Research 0014's accepted decision

This decision is fixed, unless the expert finds that it forces a later
refactor. Such findings are welcome under question 5. It is recorded in
bittype's `docs/binding-generation.md`, quoted here.

**Responsibilities:**

| Owner | Responsibility |
| --- | --- |
| horos | Declaration semantics, checked identities and applications, descriptors, ground conformance |
| bittype planning | Exact contribution selection, substitutions, references, and the selected representation plan |
| Each representation implementation (a target's structural profile, or a library's exact mapping) | Native layout, presence rules, numeric representation, construction and access |
| Binding backend, one per target | Adapter composition, runtime imports, generated glue, assembly of the requested binding products |
| Contract or domain library | Its native recipe, its adapter implementation, its protocol or domain meaning, and the laws below |

**The plan.** For each use and generic instance it carries:

> - the checked source identity and the actual generic arguments;
> - the selected contribution, target and representation profile;
> - native symbols, distinguishing run-time imports from type-only imports;
> - field access and construction;
> - the rules for absence, presence, null and native `undefined`;
> - tuple layout and calling convention;
> - union discrimination and construction;
> - numeric range and exactness;
> - recursive references;
> - the selected adapter factories, their argument slots and their resource requirements;
> - the selected root contract [...]

> **Form.** These are typed references and target-owned expressions, the same
> principle as bittype's native intent [...]. They are not a universal schema
> and not another type algebra.

> **Completeness test.** The backend can implement each selected representation
> from the plan, the source semantics and the selected contributions alone,
> without guessing a rule.

Three facts shape the design:

- **No plan exists yet.** bittype's native program keeps no per-use record
  (section 3.3).
- **The plan must be serialized.** The TypeScript backend runs in another
  process, in Node.js. The plan carries an explicit format version, which
  bittype owns; a backend refuses an unknown version.
- **The payload is plain.** bittype's native expressions carry only lists of
  strings and child expressions, with no typed metadata.

> **Interpreters.** A shared structural interpreter is allowed only if it
> executes the plan's decisions rather than choosing representation rules
> itself.

> **Selection.** [...] A native mapping without an adapter facet still supports
> types, but a request for bindings that needs a missing facet fails
> explicitly. There is never a fallback to structural conversion [...].

**Laws** (condensed from the decision). For a type T and representation profile
p, let `V_T` be the admissible abstract values, including library-owned
predicates, and `N_T,p` the admissible native values. Every contribution must
satisfy:

```text
D(E(n)) ≈ n    for every n in N_T,p   (the equivalence states whether copying, class identity, branch tags and sharing matter)
E(D(v)) = v    for every v in V_T     (abstract-value equality; any normalization is specified)
```

It must also establish, quoted:

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

> **Coverage** is never weakened to "what this implementation happens to accept".

> **Capture in practice.** Both ontos implementations copy bytes when an atom is
> built [...]. Conversion into ontos atoms therefore discharges capture for byte
> buffers. Adapters that keep native references, views or capabilities must
> discharge it themselves.

> **Binding support is qualified separately from type support** [...]. A
> representation must keep enough information for the laws to be possible:
>
> - Two native branches can encode to the same canonical value, as with an overlapping union.
> - Two abstract alternatives can collapse into one native form.
>
> Neither can bind unless a library-owned representation supplies the missing information.

**Validation**, quoted:

| Judgment | Owner |
| --- | --- |
| Ground structural conformance | horos, against descriptors emitted at build time (prefer `diagnose`: exhaustion is not a mismatch) |
| Native representation validity and construction | the structural or custom adapter |
| Additional domain predicates | the domain library |

> **Resource accounting.** Validation and conversion share one budget: check
> maximal ground regions, and never give each child a fresh budget.

> **"Declared once" means structure.** Domain laws that horos cannot express stay
> with the model library. For example, `ContentHash` is text plus a
> content-model predicate, and a non-negative count is an integer plus a
> predicate [...].

### 2.4 What is still open around this work

- **Message capture.** The service library owns a pending contract on capture
  and lifetime. It is not needed here, and recommendations should not
  presuppose its answer.
- **The run-time factories.** The service library owns the TypeScript backend.
  It has not yet been asked who should implement the run-time factories of the
  structural families (section 6.1).

## 3. Relevant code

### 3.1 ontos values and canonical data forms

A value is `Atom(bytes)` or `Tuple(value, ...)`. Nothing else exists, and
`Atom("")` differs from `Tuple()`. Data forms are labelled tuples:

| Form | Canonical value | Notes |
| --- | --- | --- |
| Integer | `Tuple(Atom("int"), Atom(<sign-magnitude bytes>))` | Sign byte `0x00` non-negative or `0x01` negative, then the magnitude big-endian with leading zero bytes stripped. Zero is the single byte `0x00`. Unbounded. |
| Text | `Tuple(Atom("utf8-text"), Atom(<UTF-8 bytes>))` | Valid UTF-8, no Unicode normalization. Producers reject unpaired UTF-16 surrogates, which JavaScript strings may contain. |
| Boolean | `Tuple(Atom("bool"), Atom(0x00 or 0x01))` | |
| List | `Tuple(Atom("list"), element, ...)` | |
| Map | `Tuple(Atom("map"), Tuple(k0, v0), Tuple(k1, v1), ...)` | Keys may be any value, "sorted by key, strictly ascending, under the codec total order". Duplicates are forbidden ("Never last-wins"). |
| Set | `Tuple(Atom("set"), e0, e1, ...)` | Sorted by codec order, no duplicates. |
| Null | `Tuple(Atom("null"))` | "`null()` means present-but-no-value, NOT absence." |
| Decimal | an exact `mantissa × 10^exponent` built on integers | horos does not recognize it (section 3.2). |
| Bytes | a bare `Atom(bytes)` | Not a labelled form. |

**The codec order.** The canonical encoding is
`encode(Atom(b)) = 0x00 ‖ uvarint(len(b)) ‖ b` and
`encode(Tuple(v0 … vn)) = 0x01 ‖ uvarint(n) ‖ encode(v0) ‖ …`, where `uvarint`
is unsigned LEB128 in shortest form, which is little-endian. Map keys and set
elements are ordered by lexicographic comparison of these whole encodings.
Below 128 bytes the length byte makes shorter keys sort first, so `"b"`
precedes `"aa"`. From 128 bytes on, the little-endian length breaks that rule:
a 256-byte key (length bytes `80 02`) sorts before a 129-byte key (`81 01`).
**Only comparison of the actual encodings is correct.** Native string or
numeric order never is.

**Where order is checked.** The codec decodes any atoms and tuples. Map order and
duplicate keys are checked by the ontos data recognizer, or by horos
conformance of a `map-of` type. A binding must call one of them.

**Immutability.** Values are immutable. The Go and TypeScript constructors copy
atom bytes:

```go
// ontos core/go/value.go:76-80
func NewAtom(bytes []byte) Atom {
	cp := make([]byte, len(bytes))
	copy(cp, bytes)
	return Atom{bytes: cp}
}
```

```ts
// ontos core/ts/src/index.ts:56-58 (class Atom, condensed)
constructor(bytes) { this.#bytes = copyBytes(bytes); Object.freeze(this); }
```

Reading an atom's bytes also returns a copy. A TypeScript `Tuple` copies its
item array and freezes it; children are shared by reference.

**Trees, with sharing allowed.** "Well-founded (acyclic) … a finite tree, not a
graph … (Sharing … is unobservable and never makes a value cyclic.)" "Equal
subtrees MAY be shared … A value MAY be held as a DAG internally … `μ` (cyclic)
representations are forbidden." Encoding and equality do not memoize, so a
heavily shared value costs its full expanded size.

**Integer libraries.** TypeScript encodes from `bigint`, Go from `*big.Int`.
The Rust implementation materializes integers as `i128` and reports
`IntOutOfRange` beyond it. **There is no C++ ontos implementation.**

**JSON** is outside this document. The existing JSON mapping is one-directional,
and a lossless two-way one has been proposed to ontos's maintainers.

### 3.2 horos descriptors for these families

A descriptor is `Tuple(Atom("horos/<constructor>"), argument, ...)`.
"Descriptors add no equality: value and descriptor identity remain exact ontos
structural identity."

| Descriptor | Admits |
| --- | --- |
| `any()` | Every value |
| `atom()`, `atom-exact(Atom)`, `atom-len(int lo, int hi)`, `atom-in(set of Atom)` | A bare atom; with exact bytes, a byte length within bounds, or membership in a finite set |
| `scalar(Atom("int" or "utf8-text" or "bool"))` | The canonical integer, text or boolean form. Integer recognition has "no host integer width bound" |
| `text-len`, `key-grammar()`, `text-non-blank()`, `lowercase-uuid()` | Text with a code-point length range, the fixed pattern `^[a-z0-9][a-z0-9._-]{0,63}$`, a non-blank character, or a lowercase UUID spelling |
| `tuple(T0, ..., Tn)` | Any tuple of exactly this arity whose children conform |
| `seq(T)` | Any tuple of any arity whose children all conform to T |
| `compound(Atom(label), T0, ...)` | A tuple whose first child is the exact label atom, then the given children |
| `list-of(T)`, `set-of(T)`, `map-of(K, V)` | A canonical list, set or map whose elements, keys and values conform |
| `record(Atom("open" or "closed"), field, ...)` | A canonical map with text keys (below) |
| `union(T, ...)`, `and(T, ...)`, `not(T)` | At least one, all, or not |
| `exact(value)` | Exactly that value |
| `rec(Atom(name), T)`, `ref(Atom(name))` | Recursion (below) |

There is **no general text pattern**, so `^sha256:[0-9a-f]{64}$` cannot be
expressed in horos. There is **no floating-point or decimal type**: "`decimal`
and `null` are not recognized scalar tags in v1."

**Records and absence.** A field is
`Tuple(Atom("horos/field"), text(name), Atom("req" or "opt"), T)`. Quoted:

> Mode is mandatory: `open` permits additional text-keyed fields with
> unconstrained values; `closed` rejects additional fields. Every required field
> must be present and conform. An optional field may be absent; if present, it
> must conform. An absent map entry and an entry containing null are different
> values.

"The singleton null type is expressible as `exact(null())`." A nullable field is
therefore `union(T, exact(null()))`, independent of whether it is optional.

**Unions are untagged.** "At least one alternative matches; empty union matches
nothing." The checker is a plain disjunction:

```ts
// horos core/ts/src/core.ts:518-523 (union arm of the conformance check)
for (const ct of args) {
  if (conformsIn(v, ct, env, pending, depth + 1)) return true;
}
```

No interface reports which alternative matched. Tags exist only where an author
writes them, for example `union(compound("item", T), compound("end"))`, or as an
exact-text field in each alternative.

**Alternatives often overlap.** The checker's documentation: "A value may
conform to many descriptors at once (overlapping classification, not a new
identity)". Overlaps that follow from the rules above:

- the empty list conforms to `list-of(T)` for every T, so
  `union(list-of(int), list-of(text))` is ambiguous on `[]`; likewise for sets
  and maps;
- two closed records whose fields are all optional both admit the empty map;
- two open records admit a map carrying both records' fields;
- `tuple` and `seq` admit any tuple, including labelled forms, so a canonical
  integer conforms to `tuple(atom(), atom())`;
- an open record with no fields and `map-of(scalar(utf8-text), any())` admit
  the same values.

**Integers have no refinement, by decision.** `scalar(int)` is the only integer
descriptor. A comment in horos's conformance vectors records the ruling:
"magnitude is NOT a form bound (int-range §5-principled-out to logos — CAA
ruling 2026-07-24)". logos is a separate project for checked logical claims.
horos's kinds, its predicates over types such as "this type is data", have
"no ... numeric-bound predicates". A finite set of integers can only be written
as a union of `exact` values.

**Recursion is anonymous and must be guarded.** Quoted:

> Every path from a `rec` to its matching `ref` must descend into a strict
> child of the value being checked. Tuple, compound, collection and record value
> positions provide descent; `union`, `and` and `not` do not. […] Inner bindings
> shadow matching outer names. The values themselves remain finite acyclic
> trees.

At the declaration layer: "Named type and kind dependency cycles are rejected."
**A recursive type is therefore written with `rec` and `ref` inside one type
expression**, never as two named declarations that refer to each other. Mutual
recursion uses nested `rec` binders.

**Declarations.** A horos module declares named, possibly generic types. A
declaration is identified by (module digest, name), where the digest is the
SHA-256 of the module's canonical bytes. horos's specification writes
declarations in a compact notation, which section 4 also uses:

- `H(c, x, ...)` is the descriptor `Tuple(Atom("horos/" + c), x, ...)`;
- `A(s)` is an atom, `T(s)` a canonical text value, `I(n)` a canonical integer
  and `L(x, ...)` a canonical list;
- `D(c, x, ...)` is `Tuple(Atom("horos/decl/" + c), x, ...)`, a
  declaration-level form: `D("param", I(0))` is the declaration's first type
  parameter, `D("local", T(name))` refers to a declaration in the same module,
  and `D("apply", E, L(E, ...))` applies a generic declaration to arguments.

A generic page of items, as written in horos's specification:

```text
Page = H("record", A("closed"),
  H("field", T("items"), A("req"), H("list-of", D("param", I(0)))),
  H("field", T("next"),  A("opt"), H("scalar", A("utf8-text"))))
```

**Limits.** Declarations: depth 256, 100,000 nodes, 16 MiB of canonical bytes,
256 modules. The checkers stop at 512 recursion frames. The boolean interface
reports an exhausted check as false, which inside `not` can become a true.
`diagnose` distinguishes `conforms`, `mismatch`, `not_well_formed` and
`exhausted`.

### 3.3 bittype's generation interface

A **mapping library** describes its rules once, at registration. The type
parameters `F`, `T` and `R` are the source language, the target language and
the source integration's reader.

```go
// contracts/generation/go/description.go (excerpts)
type MappingLibrary[F, T, R any] interface {
	Describe(from, to representation.Language) (MappingDescription[F, T, R], error)
}
type MappingDescription[F, T, R any] struct {
	Profile      string // the representation profile; libraries are selected by exact equality
	Bindings     []BindingDescription[F, T, R] // use rules
	Declarations []DeclarationRecipe[F, T, R]  // recipes; a zero selector is the structural default
	Helpers      []HelperRecipe[T]             // shared generated helper declarations
}
// A use rule: selected by exact identity, it captures child uses into named
// slots ("slot capture") and yields a native expression.
type BindingDescription[F, T, R any] struct {
	ID     string
	Input  Pattern[F, R]                                                    // selector plus slot capture
	Output Expression[T]                                                    // a fixed native expression, or
	Expand func(context.Context, Expansion[F, T, R]) (Expression[T], error) // a computed one
}
```

A **selector** is "an exact source-owned identity" that the planner compares but
never interprets. For a horos built-in form it is a key such as
`builtin/6:scalar/9:utf8-text`, with length-prefixed segments. For a declaration
it is `declaration/<module digest>/<name>`.

Native output is an **expression** over constructors that the target owns:

```go
// contracts/generation/go/expression.go (excerpts)
func Native[T any](form Form[T], data []string, args ...Expression[T]) Expression[T] // a target constructor with string data
func Mapped[T any](slot Slot) Expression[T]      // the native form of a captured child use
func Parameter[T any](index int) Expression[T]   // a declaration's own type parameter
func Argument[T any](index int) Expression[T]    // an actual generic argument
```

**Planning.** The planner walks the horos declaration graph itself:

- For each use it selects exactly one use rule by selector. More than one gives
  `ErrAmbiguous`.
- If none matches, it reports "no mapping for `<key>`" as unsupported at the
  use's source site.
- A declaration with no exact use rule is declared through a recipe: the exact
  recipe if there is one, otherwise the structural default.

The result is a closed native program. It keeps no per-use source identity and
no actual source arguments. Its only per-use record is debug strings. Recursion
is refused twice:
- the planner reports "recursive native declaration"
  (`contracts/generation/go/planning.go:322-325`);
- native construction refuses any declaration cycle
  (`transform/syntax/shared/target/go/compile.go:88-93`).

### 3.4 The implemented representation profile

All four targets implement the structural profile `pagination-representation/1`:

| horos form | TypeScript | Go | Rust | C++ |
| --- | --- | --- | --- | --- |
| UTF-8 text | `string` | `string` | `::std::string::String` | `::std::string` |
| Boolean | `boolean` | `bool` | `bool` | `bool` |
| `list-of(T)` | `T[]` | `[]T` | `::std::vec::Vec<T>` | `::std::vector<T>` |
| Named closed record | `interface` | defined `struct` | `struct` | `struct` |
| Optional field | `field?: T` | `BittypeOptional[T]` helper | `Option<T>` | `std::optional<T>` |
| Type parameter, application | native generics | native generics | native generics | templates |

TypeScript output is compiled with `strict` and `exactOptionalPropertyTypes`, so
`field?: T` does not admit an explicit `undefined`. Go's helper is a generated
struct, `BittypeOptional[T]{Present bool; Value T}`. Absence never becomes `null`
or a zero value. Inline records are refused: "nested records require a named
declaration in this profile". bittype's documentation states the boundary:

> Canonical integers remain Unsupported. The existence of native integer syntax
> does not choose an arbitrary-width representation or dependency. Open records,
> Wire, tuples, sets/maps, unions/intersections, complements, exact/refinement
> forms and ground recursion remain outside this profile.

("Wire" there is the endpoint-reference type `wire<E>`.)

**Exact mappings already work for types.** A `Cell<T>` contract library maps its
declaration, selected by exact identity, to the external TypeScript type `Cell`
exported by `@bitspark/cell`. It captures the actual for Cell's binder named `T`,
never a position. A test then generates `Page<Cell<string>>` with the import
added by the file-placement stage.

### 3.5 What the native target models can express

Each target has a small native model with the same builders: primitive,
collection, field, record and external type. **Every model admits only two
primitives**:

```go
// transform/syntax/<target>/model/go/model.go:42, identical in all four
if len(d) != 1 || len(a) != 0 || (d[0] != "string" && d[0] != "bool") {
```

| Family | TypeScript | Go | Rust | C++ |
| --- | --- | --- | --- | --- |
| Integers, bytes, maps, sets, unions | none | none | none | none |
| Indirection for recursion | not needed | none | none | none |
| External library types | yes | **no** (the model refuses them) | absolute paths | absolute paths |

The syntax layers underneath can express much more: TypeScript union, literal
and `bigint` types; Go pointer, map and array types; Rust enums and integer
primitives; C++ pointers and integer types. Only the models and mapping
libraries are missing.

### 3.6 Today's generic conversion, which this work replaces

The service library converts structural TypeScript values with one run-time
interpreter (`@bitspark/service-api` 0.3.2). Its message forms are not canonical
ontos data:

- strings, numbers and booleans are bare atoms;
- records and arrays are unlabelled tuples, each field wrapped in a tuple of
  zero or one element;
- numbers are any finite JavaScript number, written as `String(n)`;
- unions are `(alternative index, body)`, and the encoder keeps the first
  alternative that succeeds (first-success matching);
- shared subvalues are converted once per occurrence, so sharing is lost;
- one budget covers a whole call: 100,000 nodes, depth 256, 16 MiB.

Research 0014 replaced these forms with canonical data conforming to the horos
type. Only the budget principle carries over.

### 3.7 The pilot service: Content

Content stores byte strings and structured values by hash. Its contract is
today a TypeScript interface:

```ts
// services/content packages/api/src/content-service-api.ts
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
```

```ts
// services/content packages/model/src (one declaration per file)
export type ContentHash = `sha256:${string}`;
export type ContentRepresentation = 'raw' | 'map/1' | 'list/1';
export interface ContentLink { readonly representation: ContentRepresentation; readonly hash: ContentHash; }
export type ContentValue = Uint8Array | ContentMap | readonly ContentValue[];
export interface ContentMap { readonly [key: string]: ContentValue; }
```

The families it uses:

| Family | Where |
| --- | --- |
| Integers | `check().objects`, a count. A hand-written validator requires a non-negative safe integer. |
| Bytes | `put`'s argument, `get`'s result, the leaves of `ContentValue` |
| Lists | `ContentHash[]`, lists of content values |
| Maps with text keys | `ContentMap` |
| A disjoint union | `ContentValue`, told apart at run time by `instanceof Uint8Array` and `Array.isArray` |
| Recursion | `ContentValue` and `ContentMap` refer to each other |
| A finite set of text values | `ContentRepresentation` |
| Text with a domain predicate | `ContentHash` must match `^sha256:[0-9a-f]{64}$` |

Content's semantics include **sharing**. Its snapshot function "preserve[s]
shared DAGs" and copies bytes. Its own hand-written converter says:

```ts
// services/content packages/api/src/content-value-wire.ts (comment)
/**
 * A captured ContentValue as tagged wire values: `("bytes", atom)`, `("list", (value...))`
 * or `("map", ((key, value)...))`. Each subvalue is converted once and shared wherever it
 * recurs, so a compact DAG stays compact and the runtime can bound its occurrences.
 */
```

Its decoder enforces depth, step and byte limits. It rejects cycles, duplicate
map keys and unknown tags. Its keys are bare UTF-8 atoms in JavaScript
enumeration order, so its form is not canonical ontos data either. A horos
contract would replace these hand-written converters with declared types.

### 3.8 The other three services

Type names are from each service's current main branch:

| Service | Shapes | Families beyond Content's |
| --- | --- | --- |
| Cubes | `CubeSelector = readonly (string \| number)[]`; a graph whose content is `ContentValue \| ContentLink`; annotations containing `ContentHash \| CubeAnnotation` | a text-or-integer union; unions of records |
| Drafts | `revision: number`; `ContentHash \| null`; outcomes `{ kind: 'created'; draft: Draft } \| { kind: 'published'; publication: KnownPublication } \| ...`; `pending?: {...}`; an intersection over a union | nullable fields; unions tagged by an exact-text field; optional fields; intersection |
| Pointers | `PointerJson = null \| boolean \| number \| string \| readonly PointerJson[] \| PointerMetadata`, inside a hand-written converter | recursive JSON-like data with **arbitrary finite numbers**, which horos cannot express (section 3.2) |

## 4. Two values through every layer

### 4.1 Content's recursive value

Content would declare `ContentValue` in horos with `rec` inside one expression
(section 3.2):

```text
ContentValue = H("rec", A("value"),
  H("union",
    H("atom"),
    H("list-of", H("ref", A("value"))),
    H("map-of", H("scalar", A("utf8-text")), H("ref", A("value")))))
```

This union is structurally disjoint: a bare atom, a tuple labelled `list` and a
tuple labelled `map` never coincide. The canonical value of
`{"a": bytes 01, "b": []}`:

```text
Tuple(Atom("map"),
  Tuple(Tuple(Atom("utf8-text"), Atom("a")), Atom(0x01)),
  Tuple(Tuple(Atom("utf8-text"), Atom("b")), Tuple(Atom("list"))))
```

A candidate TypeScript type:

```ts
export type ContentValue = Uint8Array | readonly ContentValue[] | { readonly [key: string]: ContentValue };
```

A `ReadonlyMap<string, ContentValue>` would serve equally well for the third
alternative. What conversion needs beyond the type, and so what the plan must
carry:

| Decision | For this value | Today decided by |
| --- | --- | --- |
| Recursion | The `rec` binder `value` becomes the native declaration `ContentValue`. An inner `rec` would need a fresh native name. | nobody (refused) |
| Union, decoding | atom → alternative 0, `list` tuple → 1, `map` tuple → 2 | nobody |
| Union, encoding | `instanceof Uint8Array` → 0, `Array.isArray` → 1, otherwise → 2. The order matters, since a `Uint8Array` is also an object. | nobody |
| Bytes | Encoding copies, which ontos's atom constructor already does. Decoding yields a fresh buffer. | ontos (copying) |
| Text keys | Unpaired surrogates in a JavaScript key must be refused | nobody |
| Map form | Index-signature object or `Map`. Objects have hazards such as the key `"__proto__"`. | nobody |
| Map order | Encoding sorts by the codec order of the encoded keys. Decoding must check order and duplicates through the ontos data recognizer or horos conformance. | nobody |
| Sharing | A buffer referenced twice: encode once and share? Preserve sharing on decoding? | Content's own converter |
| Cycles | A cyclic native object graph must be refused on encoding | Content's own converter |
| Limits | Depth, node and byte limits over the expanded tree or the shared graph | Content's own converter |

### 4.2 A generic outcome with an optional nullable field, and an integer

Suppose a horos version of Drafts' outcome, applied inside the generic `Page`:

```text
Outcome = H("union",
  H("record", A("closed"),
    H("field", T("kind"),  A("req"), H("exact", T("created"))),
    H("field", T("draft"), A("req"), D("local", T("Draft")))),
  H("record", A("closed"),
    H("field", T("kind"),    A("req"), H("exact", T("unknown"))),
    H("field", T("pending"), A("opt"), H("union", D("local", T("Pending")), H("exact", Tuple(A("null")))))))

Outcomes = D("apply", D("local", T("Page")), L(D("local", T("Outcome"))))
```

Here a family must decide:

- **Discrimination.** It can come from the `kind` field, whose exact text values
  differ between the alternatives. Recognizing that is a judgment about horos
  membership.
- **Three states for `pending`.** It can be absent, present as null, or present
  as a value. In TypeScript that is `pending?: Pending | null` under
  `exactOptionalPropertyTypes`. In Go, a presence wrapper around a nullable
  form. The plan must say which is which.
- **Composition.** Every rule must hold after `Page`'s parameter is substituted.
  Mapping `Page<Outcome>` must agree with substituting the mapped `Outcome` into
  the mapped `Page`.

Two contrasting cases:

- **An overlapping union.** `union(list-of(scalar(int)), list-of(scalar(utf8-text)))`
  admits `[]` in both alternatives. Native `bigint[] | string[]` cannot say
  which alternative an empty array came from, so decoding cannot pick a branch
  and `E(D(v))` has no defined branch.
- **A count.** Content's `{ objects: int }`. The horos type is unbounded. As a
  `bigint`, every value converts. As a `number`, encoding is exact only up to
  2^53 − 1, and decoding must refuse larger values. Content's own meaning,
  "non-negative", is a domain predicate (section 2.3) that horos cannot
  express.

## 5. Constraints

- **horos semantics are given.** Unions are untagged set membership. Integers are
  unbounded, with ranges deliberately left to logos. Named declaration cycles
  are rejected. horos has no floating-point or decimal type and no general text
  pattern. horos changes may be recommended; its maintainers decide.
- **Message values are canonical ontos data conforming to the horos type**
  (0014). There is no other message form.
- **No silent narrowing.** A family must not quietly lose information or refuse
  valid values it was not declared to refuse. A refusal is an explicit,
  classified failure.
- **One mechanism** (AC-04). The structural profile and exact mappings
  contribute through the same public mechanism. The planner learns no family by
  name, and shared code holds no catalogue of families.
- **The plan** is typed references and target-owned expressions, serialized and
  versioned (section 2.3).
- **bittype's native syntax covers types only**, and bittype owns no run-time
  code today.
- **Targets.** TypeScript first, then Go. Rust and C++ types must remain
  possible. C++ has no ontos implementation, and Rust's ontos integers stop at
  `i128`.
- **Compatibility is not a constraint.** Today's native API, message forms and
  bytes may all change: the services share one installation, and no external
  client depends on them.
- **Out of scope:** JSON presentation, endpoint references (`wire<E>`), and
  Pointers' floating-point numbers, which need a horos or domain decision first.

## 6. What we have considered

None of this is decided.

### 6.1 The unit of a family, and where conversion decisions live

A family would contribute, per target, one unit with three parts:
- a **native facet**: the use rules and recipes it already supplies;
- a **value facet**: a reference to the run-time conversion;
- the **plan entries** it records for each use.

Three ways to carry the conversion decisions:

| Option | How it works | Strengths | Costs |
| --- | --- | --- | --- |
| A. Factory arguments | Each plan entry names the family's run-time factory and gives it typed arguments the family defines: field names and absence rules for a record, alternatives and discrimination for a union. The backend composes static factory calls. | Each family owns its payload; the backend stays generic; adding a family adds a factory and a payload type. | Every family needs a payload schema readable in TypeScript, and a run-time factory per target. The plan's schema is spread across families. |
| B. Target expressions | The plan carries target-owned expressions for access, construction and discrimination, which the backend writes into generated code. | Little run-time code, and fast generated code. | Risks a small code language inside the plan. bittype's native syntax covers types only. |
| C. Shared interpreter | A run-time interpreter executes the plan's decisions. | One implementation per target. | It must stay an executor (0014), and its input is A's payloads again. |

Our present leaning is A, possibly with C as one implementation of A's
factories. We are unsure who implements the structural families' run-time
factories in TypeScript. bittype owns the structural profile's decisions but no
run-time code. The service library owns the backend. A separate small run-time
package per profile is a third possibility. The service library has not yet been
asked.

### 6.2 Integers

Three options:

| Option | Native types | Behaviour |
| --- | --- | --- |
| a. Lossless structural profile | TypeScript `bigint`, Go `*big.Int`, a Rust or C++ big-integer library | Every horos integer converts. |
| b. A declared narrowing profile | `number`, `int64`, `i128`, under a separately named profile such as `bounded-integers/1` | Decoding explicitly refuses values outside the native domain. A consumer chooses the profile knowingly. |
| c. Exact domain mappings | A domain library declares a type such as `Count` with its native type and its value domain | The domain's predicate decides the native type. Others stay lossless. |

bittype's current position, in its design notes: "Unbounded integers cannot
silently become bounded integers, absence cannot collapse into a zero value, and
recursive type indirection cannot be guessed by the kernel." Our present leaning
is a with c. Option b honours "no silent narrowing" too, since the narrowing is
declared, but it narrows every integer in a contract at once.

Open points:
- **Where a narrower value domain is declared**, since horos has no ranges: a
  logos claim, the domain library (0014 already gives domain predicates to the
  domain library), or a horos extension.
- **Whether `bigint` is acceptable** at the TypeScript surface.
- **Mutability.** Go's `big.Int` is mutable, which bears on message capture.

### 6.3 Bytes, text and refinements

- **Bytes.** `atom()` becomes `Uint8Array`, `[]byte`, `Vec<u8>` or
  `std::vector<std::uint8_t>`. Encoding into ontos atoms copies. Decoding must
  hand out a buffer the caller may own.
- **Wider native types.** A length range, an exact value or a finite set of byte
  strings has no native type form, and neither has a domain text predicate such
  as `ContentHash`'s pattern. Here the native type is wider than the horos type.
  Our leaning is to check horos refinements by conformance at the boundary, and
  domain predicates by the domain library (section 2.3). A finite set of text
  values, such as `ContentRepresentation`, can become a native literal union.
- **Text.** JavaScript strings may hold unpaired surrogates. Encoding must refuse
  them.

### 6.4 Tuples, sequences, labelled compounds and `any`

- **Tuples.** `tuple(...)` becomes a native tuple. Go has none, so it would be a
  struct with positional fields.
- **Sequences.** `seq(T)` resembles a list but is a bare tuple, so it needs its
  own native form or a refusal.
- **Labelled compounds.** `compound(label, ...)` becomes a record-like native
  type whose label is fixed.
- **`any()`.** It becomes the ontos value type, as an external library type.

### 6.5 Maps and sets

- **Text keys.** These, the common case, become a native map or an
  index-signature object.
- **Key equality.** Native maps compare keys by native equality. That is correct
  only where native equality coincides with canonical equality: text, booleans,
  and `bigint` in TypeScript. Byte strings and records as keys do not
  (JavaScript compares objects by reference, and Go cannot key a map by
  `[]byte`). Our leaning: a library-owned canonical map keyed by encoded bytes,
  or refusal for bindings.
- **Order.** Encoding always sorts by the codec order of the encoded keys.

### 6.6 Unions

Our leaning:
- **When to bind.** Bind a union only when its alternatives are provably
  distinguishable on both sides:
  - on the canonical side, by a structural test that selects exactly one
    alternative;
  - on the native side, by a test over the native representations.

  First-success matching is rejected.
- **Easy cases:**
  - different `compound` labels;
  - different canonical forms (atom, list, map);
  - an exact-text field shared by all alternatives, as in section 4.2;
  - `union(T, exact(null()))` where T excludes null.
- **Overlapping unions** are refused for bindings unless a library-owned
  representation supplies the missing information (0014).

Who judges canonical distinguishability is open:
- the horos project, since it is a fact about membership;
- bittype's horos source integration, which already judges uses for the planner;
- each union family.

The judgment is easy for the cases above and hard in general (`and`, `not`,
`any`, open records, and Drafts' intersection over a union).

### 6.7 Absence, null and open records

- **Optional fields** stay as they are (section 3.4).
- **Nullable fields.** `union(T, exact(null()))` becomes `T | null`, `*T`,
  `Option<T>` or similar. It stays distinct from absence, so an optional nullable
  field has three states (section 4.2).
- **Open records** carry extra fields of any value. A native form needs a place
  for them, a separate map member, or refusal. Our leaning is to refuse open
  records for bindings until a representation is chosen deliberately.

### 6.8 Recursion and sharing

- **Naming.** Each `rec` binder becomes a named native declaration. A binder that
  forms a declaration's whole body takes the declaration's name. Inner binders
  need a naming rule that avoids collisions.
- **Indirection** is chosen per recursive edge by the target's profile, and the
  plan records each edge and its indirection:
  - none in TypeScript;
  - Go: a pointer, or none through a slice;
  - Rust: `Box`, or none through `Vec`;
  - C++: `std::unique_ptr`, or `std::vector` where an incomplete element type is
    allowed.
- **Cycles and sharing.** Native values may be cyclic or shared; canonical values
  are trees. Encoding refuses cycles. Sharing is unobservable in ontos, so
  `D(E(n)) ≈ n` could ignore it. But Content relies on DAG compactness for its
  resource bounds. Two questions are open:
  - whether sharing is part of a family's contract or only an implementation
    property;
  - whether limits count the expanded tree or the shared graph.

### 6.9 Qualification and order

Each family and target is qualified for bindings separately from types:
- by the laws of section 2.3, each shown failing when its rule is broken on
  purpose;
- by the completeness test, in which a test backend implements each family from
  the plan, the source semantics and the contributions alone.

TypeScript comes first. Families arrive in the order Content needs them:
1. integers;
2. bytes;
3. text-keyed maps;
4. disjoint unions;
5. recursion.

### 6.10 Hypotheses we want attacked

1. "A structural profile never narrows a value domain; only a declared profile or
   an exact domain mapping may."
2. "Whether a union's alternatives are distinguishable is a judgment about horos
   membership, owned on the source side, not by a target."
3. "Beyond the source semantics, the plan needs only the selected contributions
   and their contribution-owned arguments."
4. "Sharing is never part of a family's semantic contract." This is in tension
   with Content's reliance on compact shared graphs.

## 7. Questions for the expert

For each question, principles and patterns from comparable systems, including
what failed, are especially welcome. Such systems include code generators for
schema and serialization languages. Section 6.10 lists four hypotheses: where an
answer bears on one, please say where it holds and where it breaks.

1. **The unit of a family, and its owners.** What should one representation
   family contribute, per target, so that its native type, plan entries and
   conversion are defined once, always agree, and can be qualified for bindings
   separately from types? Our leaning (6.1) is that plan entries are typed
   arguments of a family-owned run-time factory. The alternatives are target
   expressions carried in the plan, or a shared interpreter. bittype owns the
   structural profile's decisions but no run-time code. Who should own the
   structural families' run-time factories?

2. **Where native and canonical domains differ.** A native type can hold less
   than the horos type: `number` for an unbounded integer, or a native map whose
   key equality is not canonical equality. It can also hold more: `Uint8Array`
   for a length-bounded atom, `string` for text with a domain predicate. Our
   leaning (6.2, 6.3) is that the structural profile never narrows, and that
   wider native types are checked at the boundary. How would you approach both
   directions, and where should a narrower value domain be declared, given that
   horos deliberately has no ranges?

3. **Unions.** horos unions are untagged set membership, and their alternatives
   often overlap (3.2). How would you decide when a union can be bound to a
   native union, and where should the judgment that its alternatives are
   distinguishable live? It is easy for distinct compound labels, distinct
   canonical forms, a shared exact-text field or `T | null`, and hard with
   `and`, `not`, `any` and open records. How far can such a judgment usefully go?

4. **Recursion and native identity.** horos recursion is anonymous, and
   canonical values are immutable finite trees whose internal sharing cannot be
   observed. Native targets need named declarations and sometimes indirection,
   and native values can be mutated, shared or cyclic, while Content relies on
   compact shared graphs for its resource bounds. How would you divide naming,
   indirection, cycle refusal, sharing and limit accounting between the source
   side, each family and the binding backend? Please do not presuppose the
   pending capture contract.

5. **Open.** What should we know that these questions do not ask? In particular,
   which choices here are most likely to force a refactor as more families and
   the Go backend arrive, and what would you do instead?
