# Research: Native representation families for a type language's values, designed together with their conversions

**ID:** 0015
**Date:** 9 November 2026
**Status:** draft
**Owner:** claude-bittype-gaps
**Issue:** https://github.com/Bitspark/bittype/issues/203
**Code baseline:** bittype 4fc859e2882efdc51b2991ebb7e43b5c4d6c515f; horos 3a51a06e6db1800a39f94782213609a2d8723d10; ontos 0d5f34e; platform repository `libraries/service` 24fedd4; platform repository `services/content` eaae493; services `cubes` bed6164, `drafts` e1811f5, `pointers` e2ab221

## Question

**horos** is a declarative type language whose types are themselves values of a
shared value model, **ontos**. **bittype** is a Go library that reads horos
declarations and generates native type declarations for TypeScript, Go, Rust
and C++. A software platform will declare its network services' contracts in
horos. A build then produces two things from each contract: native types, and
**bindings**, which convert native values to and from ontos values, plus the
operation data that the command line and web forms need.

An earlier consultation (research 0014, 9 November 2026) settled who owns which
part of that build. Its decision is quoted in section 2.3. In short:

- bittype selects the native representation of every type and publishes those
  choices as a typed **selected representation plan**;
- a **binding backend** per target language reads the plan and emits the
  conversion code, without choosing any representation rule itself;
- message values are canonical ontos data that conform to the horos type.

Today bittype can represent only a small set of horos types: text, booleans,
lists, named closed records with optional fields, and generic declarations. The
pilot service, Content, needs more:

- **integers**: a count of stored objects;
- **bytes**: stored content;
- **maps**: a content value can be a map;
- **unions with tags**: a content value is bytes, a list or a map;
- **recursive types**: a list or map of further content values.

The other services need some of the same families (section 3.7).

The owner asked for this to be done properly, so that nothing has to be
refactored later. That rules out the obvious order. One could add native types
for these families first and the plan and conversions afterwards. But the
conversion for a family constrains which native type is admissible at all. A
TypeScript `number` cannot hold every horos integer. A JavaScript `Map` compares
object keys by reference, not by value. And the plan has to record decisions
that the type declaration alone does not show. Each family would have to be
reopened.

**How should a representation family be defined, so that it is complete the
first time?** By "complete" we mean that it fixes, for each target:

- its native type;
- what the plan records about each use of it;
- how its values convert to and from ontos values, and the laws that
  conversion obeys;
- what is refused, and how binding support is qualified, separately from type
  support.

Families must compose (generics, nesting, recursion) and stay additive: adding a
family must not change shared machinery or other families.

**The tensions we see:**

- **Unbounded integers.** horos integers have no size limit and horos has no
  range refinement, while ergonomic native integer types are bounded.
- **Untagged unions.** A horos union is plain set membership. Alternatives
  often overlap, and nothing reports which alternative a value belongs to.
  Native languages want a tag.
- **Maps.** Map keys may be any value and are ordered by their canonical
  encoding. Native maps use native equality and native order.
- **Recursion.** horos recursion is anonymous and structural, written inside one
  type expression. Every target needs a named declaration for it, and some need
  an explicit indirection.
- **Mutable native values.** Native byte buffers are mutable, while message
  values are immutable.
- **The plan's form.** The plan must carry these decisions without becoming a
  universal schema language or a second type algebra (a constraint from 0014).

Sections 1 to 3 give the context. Section 4 follows one value through every
layer. Section 5 lists what is fixed. Section 6 gives our present leanings,
none decided. Section 7 holds the five questions.

## 1. System overview

### 1.1 The components

| Component | What it is | Language(s) |
| --- | --- | --- |
| **ontos** | The value model. A value is an **atom** (a byte string) or a **tuple** (an ordered sequence of values). On that floor it defines canonical data forms (integers, text, booleans, lists, maps, sets, null) as labelled tuples, and one canonical binary encoding. | Go, Rust, TypeScript; no C++ implementation |
| **horos** | A type language whose types are ontos values called **descriptors**. It has a conformance checker for values and a declaration language with modules, generic declarations and imports pinned by content digest. bittype supplies its readable text syntax. | Descriptors and conformance: Go, Rust, TypeScript. Declaration semantics: Go only |
| **bittype** | A Go library: horos syntax and analysis, and native type generation for TypeScript, Go, Rust and C++ through a public generation interface. Its native syntax covers type declarations only, not executable code. | Go; emits TypeScript, Go, Rust, C++ |
| **bitruntime** | The message runtime. It carries ontos values between endpoints and never consults types. | Go, TypeScript |
| **platform repository, `libraries/service`** | The TypeScript service library: today's adapters and call protocol, and the future TypeScript binding backend. | TypeScript |
| **platform services** | Content, Cubes, Pointers and Drafts. Content is the pilot for horos-declared contracts. | TypeScript |
| **logos** | A separate project for checked logical claims. horos deliberately leaves numeric ranges to it (section 3.2). | (not used in this document's code) |

### 1.2 Terms used throughout

| Term | Meaning in this document |
| --- | --- |
| **Native value / native type** | A value or type of a target language, for example a TypeScript `bigint` or `interface`. |
| **Canonical value** | An ontos value in the canonical data form that a horos type describes, for example `Tuple(Atom("int"), Atom(0x00 0x05))` for the integer 5. |
| **Descriptor, conformance** | A horos type written as an ontos value; `conforms(value, descriptor)` checks a value against it. |
| **Representation family** | One kind of horos type, together with how it is represented natively: integers, bytes, lists, maps and so on. |
| **Representation profile** | A named policy in bittype deciding how horos forms become native types in one target, for example `pagination-representation/1`. |
| **Mapping library, recipe** | A bittype contribution that supplies native representations. A **binding description** is its rule for uses of one horos form; a **recipe** is its rule for one horos declaration. Libraries select by exact identity, never by name or shape. |
| **Selected representation plan** | The plan of research 0014: every representation decision bittype made for one generation, in typed form, published for binding backends. Not implemented yet. |
| **Facet** | One part of a library's contribution for one target. The **native facet** produces the native type. The **value facet** converts values. 0014 also defines a root facet and an operation facet, which do not matter here. |
| **Binding backend** | The per-target component that reads the plan and emits conversion code. The TypeScript backend lives in `libraries/service`. |
| **Encode E, decode D** | The value facet's two directions: `E(native) → canonical`, `D(canonical) → native`. |
| **Native domain** | The set of horos values a native representation can hold. A representation **narrows** the domain if some valid horos value has no native form. |
| **Binding support** | A family is binding-supported for a target when its conversion laws are tested (section 2.3), separately from being able to generate its type. |
| **Qualified** | Validated by tests and checks for a stated scope. |

### 1.3 Scale, people and direction

- **Consumers.** Four TypeScript services with 6, 6, 7 and 10 operations,
  running in Node.js and browsers.
- **People.** One human owner decides for all the repositories involved. Each
  repository is maintained by AI coding-agent sessions, several of which work in
  parallel. Additive composition with clear ownership is therefore a hard
  requirement.
- **Targets.** TypeScript bindings now; Go bindings expected. bittype also emits
  Rust and C++ types. Bindings for those are not planned, but the design must
  not preclude them.
- **What may change.** The native TypeScript API, the message form and the
  bytes on the network may all change: the services share one installation and
  no external client depends on them. horos changes may be recommended, and
  horos's maintainers decide.

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

The owner has made architecture part of correctness. These numbered rules apply
to every design here:

| Rule | Requirement |
| --- | --- |
| AC-01 | Architecture is part of correctness: ownership, dependency and composition contracts are requirements, like behaviour. |
| AC-02 | Meanings remain with their semantic owners. Independently owned implementations supply their domain interpretations, policies and laws; shared mechanisms own their invariant guarantees. |
| AC-03 | Extensions compose additively. Adding a contribution changes only that contribution and the consumers that compose it, not shared orchestration, a central list of cases, or unrelated implementations. |
| AC-04 | Built-in and first contributions use the same public extension mechanism as independent ones. |
| AC-05 | Composition preserves guarantees: identity, scope, ownership, authority, assumptions, resource limits and failure meaning. |
| AC-06 | A small first coverage does not redefine the intended extension boundary. |
| AC-07 | A test or a compiler's acceptance supports only the claim it establishes. |
| AC-08 | A change to a shared contract needs a general justification. |
| AC-09 | Design interfaces from the domain, its owners and the guarantees consumers need. The existing API and the effort of changing it must never influence the design, not even as a tiebreaker. |

### 2.3 Research 0014's accepted decision

The decision is recorded in bittype's `docs/binding-generation.md`. These are
the parts this document builds on, quoted.

**Responsibilities:**

| Owner | Responsibility |
| --- | --- |
| horos | Declaration semantics, checked identities and applications, descriptors, ground conformance |
| bittype planning | Exact contribution selection, substitutions, references, and the selected representation plan |
| Each representation implementation (a target's structural profile, or a library's exact mapping) | Native layout, presence rules, numeric representation, construction and access |
| Binding backend, one per target | Adapter composition, runtime imports, generated glue, assembly of the requested binding products |
| Contract or domain library | Its native recipe, its adapter implementation, its protocol or domain meaning, and the laws below |

**The plan.** "For each use and generic instance, the plan carries":

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
>
> The serialized plan carries an explicit format version, which bittype owns. A
> backend consumes a pinned version and refuses an unknown one.

The plan is "typed references and target-owned expressions. They are not a
universal schema or a new type algebra." A completeness test must show that "a
test backend implements each supported representation from the plan and the
selected contributions alone, without guessing a rule."

**Interpreters.** "A shared structural interpreter is allowed only if it executes
the plan's decisions rather than choosing representation rules itself."

**No fallback.** "A native mapping without an adapter facet still supports types,
but a request for bindings that needs a missing facet fails explicitly. There is
never a fallback to structural conversion."

**Laws every value facet must satisfy** (from the 0014 advice):

| Law | Meaning |
| --- | --- |
| Round trip on native values | `D(E(n)) ≈ n` for every admissible native value `n`, with `≈` specified (copying, class identity, sharing) |
| Round trip on canonical values | `E(D(v)) = v` for every admissible canonical value `v`, with any normalization specified |
| Conformance | `E(n)` conforms to the horos descriptor |
| Coverage | the admissible native values and the admissible canonical values correspond as stated |
| Capture | later mutation of caller-owned data cannot change a captured message |
| Generic composition | substitution, aliases, nested applications and recursion preserve the laws |

**Binding support is qualified separately from type support** (0014, accepted):
"Information-losing representations, such as overlapping unions, are rejected
for bindings unless a library-owned representation supplies the missing
information."

### 2.4 Earlier advice that bears on the families

From the 0014 advice, quoted:

> An integer declaration does not by itself remove Content's result validator.
> The declared integer is signed and unbounded, while `check.objects` currently
> means a non-negative integer. Non-negativity still needs a refinement or
> domain-owned predicate. The native representation also needs range and
> exactness rules.

> Consider a horos union of integer values and the exact integer zero. Suppose a
> native tagged sum distinguishes `GeneralInteger(0)` from `ExactZero`. Both
> encode to the same canonical integer zero. Decoding cannot recover which native
> branch was used.

> Serde explicitly chooses declaration-order first-success matching for untagged
> enums. That is a specified policy, not a consequence of union membership.
> Horos's set-theoretic union semantics do not by themselves authorize importing
> that policy.

> Synchronous conversion does not by itself establish capture. A synchronous
> adapter can return an atom backed by the caller's mutable byte buffer. Later
> mutation can still change the message.

> Validation and conversion also need coordinated resource accounting. Giving
> every child a fresh whole-operation budget defeats the intended bound.

bittype's own design notes already say: "Unbounded integers cannot silently
become bounded integers, absence cannot collapse into a zero value, and
recursive type indirection cannot be guessed by the kernel." And on generics:
"mapping after substitution must agree with substituting mapped arguments into
the mapped generic declaration."

### 2.5 Pending contracts in `libraries/service`

The 0014 advice asked for four short contracts before any interface is
committed. bittype's, the native representation laws above, is written. The
service library owns the other three and has not started them:

- the service root: a declared `Service<Ops>` contract dispatching by path;
- capture and lifetime: what a call captures and how long exported endpoints
  live;
- domain outcomes versus invocation uncertainty: `Function<A, B, E>` with the
  reply protocol `ok | error | protocol | unknown`.

The capture contract in particular constrains value facets. This consultation
does not need its answer, but its recommendations should not presuppose one.

## 3. Relevant code

### 3.1 ontos values and canonical data forms

A value is `Atom(bytes)` or `Tuple(value, ...)`. Nothing else exists, and
`Atom("")` differs from `Tuple()`. Data forms are labelled tuples. Each family
below is quoted from the ontos data specification:

| Form | Canonical value | Notes |
| --- | --- | --- |
| Integer | `Tuple(Atom("int"), Atom(<sign-magnitude bytes>))` | Sign byte `0x00` non-negative or `0x01` negative, then the magnitude big-endian with leading zero bytes stripped. Zero is exactly the single byte `0x00`. Unbounded. |
| Text | `Tuple(Atom("utf8-text"), Atom(<UTF-8 bytes>))` | Valid UTF-8, no Unicode normalization. |
| Boolean | `Tuple(Atom("bool"), Atom(0x00 or 0x01))` | |
| List | `Tuple(Atom("list"), element, ...)` | |
| Map | `Tuple(Atom("map"), Tuple(k0, v0), Tuple(k1, v1), ...)` | Keys may be any value. Entries are "strictly ascending under the codec total order", with duplicates forbidden ("Never last-wins"). |
| Set | `Tuple(Atom("set"), e0, e1, ...)` | Sorted by codec order, no duplicates. |
| Null | `Tuple(Atom("null"))` | "`null()` means present-but-no-value, NOT absence." |
| Bytes | a bare `Atom(bytes)` | Not a labelled form. |

**The codec order.** The canonical encoding is
`encode(Atom(b)) = 0x00 ‖ uvarint(len(b)) ‖ b` and
`encode(Tuple(...)) = 0x01 ‖ uvarint(n) ‖ ...`. Map keys and set elements are
ordered by lexicographic comparison of these encodings. Because the length comes
first, shorter keys sort first: `"b"` precedes `"aa"`. Native string order and
native numeric order are therefore not canonical order.

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
// ontos core/ts/src/index.ts:56-58 (class Atom)
constructor(bytes: BytesLike) { this.#bytes = copyBytes(bytes); Object.freeze(this); }
```

A TypeScript `Tuple` shallow-copies its item array and freezes it, so children
are shared by reference.

**Trees, with sharing allowed.** "Well-founded (acyclic) … a finite tree, not a
graph … Sharing … is unobservable and never makes a value cyclic." "Equal
subtrees MAY be shared … A value MAY be held as a DAG internally … cyclic
representations are forbidden." Encoding and equality do not memoize, so a
heavily shared value is processed at its full expanded size.

**Integer libraries.** TypeScript encodes from `bigint`, and Go from
`*big.Int`. The Rust implementation materializes integers as `i128` and reports
`IntOutOfRange` beyond it. **There is no C++ ontos implementation.**

**JSON.** The existing named profile `ontos-data-json/1` is one-directional,
from JSON to ontos values. It cannot produce bytes, sets, tuples, labelled
compounds or non-text map keys. A lossless two-way profile has been proposed to
ontos's maintainers as a result of research 0014. Presentation is outside this
document.

### 3.2 horos descriptors for these families

A descriptor is `Tuple(Atom("horos/<constructor>"), argument, ...)`.
"Descriptors add no equality: value and descriptor identity remain exact ontos
structural identity." The table lists the constructors and what the
conformance checker tests for each:

| Descriptor | Admits |
| --- | --- |
| `any()` | Every value |
| `atom()`, `atom-exact(Atom)`, `atom-len(int lo, int hi)`, `atom-in(set of Atom)` | A bare atom; with exact bytes, a byte length within bounds, or membership in a finite set |
| `scalar(Atom("int" or "utf8-text" or "bool"))` | The canonical integer, text or boolean form. Integer recognition has "no host integer width bound" |
| `tuple(T0, ..., Tn)` | Any tuple of exactly this arity whose children conform |
| `seq(T)` | Any tuple of any arity whose children all conform to T |
| `compound(Atom(label), T0, ...)` | A tuple whose first child is the exact label atom, then the given children |
| `list-of(T)`, `set-of(T)`, `map-of(K, V)` | A canonical list, set or map whose elements, keys and values conform |
| `record(Atom("open" or "closed"), field, ...)` | A canonical map with text keys (below) |
| `union(T, ...)`, `and(T, ...)`, `not(T)` | At least one, all, or not |
| `exact(value)` | Exactly that value |
| `rec(Atom(name), T)`, `ref(Atom(name))` | Recursion (below) |

A field is `Tuple(Atom("horos/field"), text(name), Atom("req" or "opt"), T)`.

**Records and absence**, quoted:

> `open` permits additional text-keyed fields with unconstrained values;
> `closed` rejects additional fields. Every required field must be present and
> conform. An optional field may be absent; if present, it must conform. An
> absent map entry and an entry containing null are different values.

"The singleton null type is expressible as `exact(null())`." A nullable field is
therefore `union(T, exact(null()))`, independent of whether it is optional.

**Unions are untagged.** "At least one alternative matches; empty union matches
nothing." The checker is a plain disjunction:

```ts
// horos core/ts/src/core.ts:518-523 (union arm of conformsIn)
for (const ct of args) {
  if (conformsIn(v, ct, env, pending, depth + 1)) return true;
}
```

No interface reports which alternative matched. Tags exist only where an author
writes them, for example `union(compound("item", T), compound("end"))`, or as an
exact-text field in each alternative.

**Alternatives often overlap.** The checker states: "A value may conform to many
descriptors at once (overlapping classification, not a new identity)." Overlaps
that follow from the rules above:

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
descriptor. A horos conformance test records the ruling: "magnitude is NOT a
form bound (int-range principled-out to logos)". Kinds, horos's type-level
predicates, have "no … numeric-bound predicates". A finite set of integers can
only be written as a union of `exact` values. The only refinements are on text
and atoms: lengths, membership, exact values and a few text grammars.

**Recursion is anonymous and must be guarded.** Quoted:

> Every path from a `rec` to its matching `ref` must descend into a strict
> child of the value being checked. Tuple, compound, collection and record value
> positions provide descent; `union`, `and` and `not` do not. […] Inner bindings
> shadow matching outer names. The values themselves remain finite acyclic
> trees.

At the declaration layer, "Named type and kind dependency cycles are rejected."
**A recursive type is therefore written with `rec` and `ref` inside one type
expression**, never as two named declarations that refer to each other. Mutual
recursion uses nested `rec` binders.

**Declarations.** A horos module declares named, possibly generic types. A type
expression in a declaration is a descriptor, or a parameter reference, a
reference to a local or imported declaration, or an application of a generic
declaration to arguments. A declaration is identified by (module digest, name),
where the digest is the SHA-256 of the module's canonical bytes. This
declaration from horos's specification is a generic page:

```text
H("record", A("closed"),
  H("field", T("items"), A("req"), H("list-of", D("param", I(0)))),
  H("field", T("next"),  A("opt"), H("scalar", A("utf8-text"))))
```

`H(c, ...)` is the descriptor constructor `c`. `A`, `T` and `I` are an atom, a
text value and an integer. `D("param", I(0))` refers to the declaration's first
type parameter.

**Limits.** Declarations: depth 256, 100,000 nodes, 16 MiB of canonical bytes,
256 modules. The TypeScript checker stops at 512 recursion frames. Its boolean
interface reports an exhausted check as non-conforming. A separate `diagnose`
interface distinguishes `conforms`, `mismatch`, `not_well_formed` and
`exhausted`.

### 3.3 bittype's generation interface

A **mapping library** describes its rules once, at registration. `F`, `T` and `R`
are type parameters for the source language, the target language and the
source-specific reader.

```go
// contracts/generation/go/description.go (excerpts)
type MappingLibrary[F, T, R any] interface {
	Describe(from, to representation.Language) (MappingDescription[F, T, R], error)
}
type MappingDescription[F, T, R any] struct {
	Profile      string // the representation profile; libraries are selected by exact equality
	Bindings     []BindingDescription[F, T, R]
	Declarations []DeclarationRecipe[F, T, R]
	Helpers      []HelperRecipe[T]
}
// A rule for uses of one source form, selected by exact identity.
type BindingDescription[F, T, R any] struct {
	ID     string
	Input  Pattern[F, R] // a selector plus the child uses to capture into slots
	Output Expression[T] // a fixed native expression, or ...
	Expand func(context.Context, Expansion[F, T, R]) (Expression[T], error) // ... a computed one
}
type Pattern[F, R any] struct {
	Selector Selector
	Slots    []Slot
	Capture  func(context.Context, R, UseView[F]) ([]Capture[F], error)
}
// A rule for one declaration; a zero selector is the structural default.
type DeclarationRecipe[F, T, R any] struct {
	ID       string
	Selector Selector
	Expand   func(context.Context, DeclarationRequest[F, T, R]) (DeclarationDescription[T], error)
}
```

A **selector** is "an exact source-owned identity" that the planner compares
but never interprets. For horos built-in forms it is a key such as
`builtin/6:scalar/9:utf8-text`. For a declaration it is
`declaration/<module digest>/<name>`.

Native output is an **expression** over target-owned constructors:

```go
// contracts/generation/go/expression.go (excerpts)
// Form[T] is a target-owned constructor, readable only by its issuer;
// "not an AST kind or a globally registered type kind".
func Native[T any](form Form[T], data []string, args ...Expression[T]) Expression[T]
func Mapped[T any](slot Slot) Expression[T]    // the native form of a captured child use
func Parameter[T any](index int) Expression[T] // a declaration's own type parameter
func Argument[T any](index int) Expression[T]  // an actual generic argument
func Helper[T any](key HelperKey, args ...Expression[T]) Expression[T] // a shared generated helper
```

A native expression's payload is a list of strings plus child expressions. It
carries no numeric or other typed metadata.

**Planning.** The planner walks the horos declaration graph itself. For each use
it selects exactly one binding by selector. More than one is `ErrAmbiguous`.
None is "no mapping for `<key>`", reported as unsupported at the use's source
site. A use of a declaration with no exact binding is declared through a
recipe: the exact recipe if there is one, otherwise the structural default. The
result is a closed native program for the target:

```go
// contracts/generation/go/native.go and planning.go (excerpts)
type NativeDeclaration[T any] struct { Name string; Form Form[T]; Origin ast.NodeRef; Parameters []NativeParameter; Body Expression[T] }
type NativeParameter struct { Name string; Origin ast.NodeRef }
type Program[T any] struct { declarations []PlannedDeclaration[T]; exposed []uint32 }
type Step struct { Binding, Source, Declaration, Form string; Expression string } // debug strings only
```

`Program[T]` keeps no per-use source identity and no actual source arguments.
The only per-use record is `Step`, whose fields are debug strings. **Nothing
resembling the selected representation plan exists yet.** Recursion is refused
in two places: the planner reports "recursive native declaration", and native
construction refuses any declaration cycle.

### 3.4 The implemented representation profile

All four targets implement one profile, `pagination-representation/1`, through
mapping libraries with the same structure.

| horos form | TypeScript | Go | Rust | C++ |
| --- | --- | --- | --- | --- |
| UTF-8 text | `string` | `string` | `::std::string::String` | `::std::string` |
| Boolean | `boolean` | `bool` | `bool` | `bool` |
| `list-of(T)` | `T[]` | `[]T` | `::std::vec::Vec<T>` | `::std::vector<T>` |
| Named closed record | `interface` | defined `struct` | `struct` | `struct` |
| Optional field | `field?: T` | `BittypeOptional[T]` helper | `Option<T>` | `std::optional<T>` |
| Type parameter, application | native generics | native generics | native generics | templates |

TypeScript output is compiled with `strict` and `exactOptionalPropertyTypes`,
so `field?: T` does not admit an explicit `undefined`. Go's helper is a
generated struct, `BittypeOptional[T]{Present bool; Value T}`. Absence never
becomes `null` or a zero value.

The text mapping for TypeScript, which is typical, in outline:

```go
// transform/syntax/horos-to-typescript/mapping/go/mapping.go:15-25 (outline)
pattern, err := rules.Pattern(from, "scalar", "utf8-text")
// Bindings: {{ID: "utf8-text/typescript/1", Input: pattern, Output: model.Primitive("string")}}
```

The record recipe refuses inline records: "nested records require a named
declaration in this profile". The project's documentation states the boundary:

> Canonical integers remain Unsupported. The existence of native integer syntax
> does not choose an arbitrary-width representation or dependency. Open records,
> Wire, tuples, sets/maps, unions/intersections, complements, exact/refinement
> forms and ground recursion remain outside this profile.

**A library-owned representation already works for types.** A `Cell<T>` contract
library maps its declaration, by exact identity, to an external TypeScript type:

```go
// integration/go/horos_generation_composition_test.go:197-218 (abridged)
key, err := source.Named(from, pin, "Cell") // "declaration/<64-hex pin>/Cell"
slot := g.TypeSlot("V")
pattern := g.Pattern[h.Language, source.Reader]{Selector: key, Slots: []g.Slot{slot},
	Capture: func(ctx context.Context, src source.Reader, use g.UseView[h.Language]) ([]g.Capture[h.Language], error) {
		// select the actual for Cell's own binder T, never a position
		for _, a := range use.Actuals {
			if p, err := src.Parameter(a.Binder); err != nil {
				return nil, err
			} else if p.Name == "T" {
				return []g.Capture[h.Language]{{Slot: slot, Use: a.Use}}, nil
			}
		}
		return nil, g.ErrScope
	}}
// Binding: {ID: "cells/typescript/1", Input: pattern,
//           Output: model.External("@bitspark/cell", "Cell", g.Mapped[ts.Language](slot))}
```

`Page<Cell<Person>>` then generates `Page<Cell<Person>>` in TypeScript, with an
import of `Cell` from `@bitspark/cell` that the file-placement stage adds.

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
| Integers | none | none | none | none |
| Bytes | none (an external global `Uint8Array` is possible, untested) | none | none | none |
| Maps and sets | none | none | none | none |
| Unions | none | none | none | none |
| Indirection for recursion | not needed | none | none | none |
| External library types | yes | **no** (the model refuses `external`) | absolute paths | absolute paths |

The native syntax layers underneath can express much more: TypeScript union,
literal and `bigint` types; Go pointer, map and array types; Rust enums with
variants, `Box` paths and integer primitives; C++ pointers and integer types.
Only the models and the mapping libraries are missing.

### 3.6 The pilot service: Content

Content stores byte strings and structured values by hash. Its contract is
today a TypeScript interface (platform repository `services/content`, commit
`eaae493`):

```ts
// packages/api/src/content-service-api.ts
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
// packages/model/src (one declaration per file)
export type ContentHash = `sha256:${string}`;
export type ContentRepresentation = 'raw' | 'map/1' | 'list/1';
export interface ContentLink { readonly representation: ContentRepresentation; readonly hash: ContentHash; }
export type ContentValue = Uint8Array | ContentMap | readonly ContentValue[];
export interface ContentMap { readonly [key: string]: ContentValue; }
```

The families this contract uses:

- **Integers.** `check().objects` counts stored objects. TypeScript says
  `number`. A hand-written validator requires a non-negative safe integer in
  canonical decimal text.
- **Bytes.** `put`'s argument, `get`'s result, and the leaves of `ContentValue`.
- **Lists.** `ContentHash[]` and lists of content values.
- **Maps with text keys.** `ContentMap`.
- **A union, untagged in TypeScript.** `ContentValue` is told apart at run time by
  `instanceof Uint8Array` and `Array.isArray`.
- **Recursion.** `ContentValue` and `ContentMap` refer to each other.
- **A finite set of text values.** `ContentRepresentation`.
- **A text pattern.** `ContentHash` must match `^sha256:[0-9a-f]{64}$`.

Content's semantics include **sharing**. Its snapshot function "preserve[s]
shared DAGs" and copies bytes. Its own conversion of a content value says:

```ts
// packages/api/src/content-value-wire.ts (comment and core of the encoder)
/**
 * A captured ContentValue as tagged wire values: `("bytes", atom)`, `("list", (value...))`
 * or `("map", ((key, value)...))`. Each subvalue is converted once and shared wherever it
 * recurs, so a compact DAG stays compact and the runtime can bound its occurrences.
 */
const memo = new Map<object, Tuple>();
// ...
result = tuple([tags.map, tuple(Object.keys(map).map(key => tuple([textAtom(key), visit(map[key]!)])))]);
```

Its decoder enforces a depth limit, a step limit and a byte limit. It rejects
cycles ("Cyclic wire value"), non-atom map keys, duplicate map keys and unknown
tags. Then it snapshots the result. Two observations:

- The encoder takes map keys in JavaScript enumeration order, which is not a
  canonical order. Content's separate storage format sorts keys by UTF-8 bytes.
- `ContentHash`, `ContentLink` and `ContentValue` are hand-written "semantic
  codecs" today. A horos contract would replace them with declared types.

### 3.7 The other three services

The same families recur. Type names are from each service's current main
branch:

| Service | Shapes | Families |
| --- | --- | --- |
| Cubes | Validation reports with five counts; `CubeSelector = readonly (string \| number)[]`; a cube graph whose content is `ContentValue \| ContentLink`; annotations that contain `ContentHash \| CubeAnnotation` | integers, a text-or-integer union, recursion, unions of records |
| Drafts | `revision: number`; `ContentHash \| null`; outcomes such as `{ kind: 'created'; draft: Draft } \| { kind: 'published'; publication: KnownPublication }`; optional `pending` | integers, nullable fields, unions tagged by an exact text field, optional fields, an intersection over a union |
| Pointers | `PointerJson = null \| boolean \| number \| string \| readonly PointerJson[] \| PointerMetadata`, inside a hand-written codec; nullable unions; tagged outcome unions | recursive JSON-like data with arbitrary finite numbers, nullable fields, tagged unions |

### 3.8 Today's generic conversion, which this work replaces

The platform's service library converts structural TypeScript values with one
run-time interpreter (`@bitspark/service-api` 0.3.2, `dataAdapter`). Its choices
were made without a type language, and are listed so the expert can see what
not to inherit:

- **Numbers** are any finite JavaScript number, encoded as the text of
  `String(n)`, with `-0` written `"-0"`. There is no integer kind and no
  `bigint` support.
- **Booleans** are the atoms `"true"` and `"false"`; `null` is an empty atom.
- **Records** are entries sorted by JavaScript's default string sort, which is
  UTF-16 code-unit order and differs from UTF-8 byte order for some characters.
  An absent optional field and a field present as `undefined` encode
  differently.
- **Unions** are encoded as `(alternative index, body)`. The encoder tries each
  alternative in order and keeps the first that succeeds (first-success
  matching).
- **Shared subvalues** are converted once per occurrence, so sharing is lost.
- **Cycles** are rejected, and a whole call has one budget: 100,000 nodes,
  depth 256, 16 MiB.

None of these forms is canonical horos data. Research 0014 decided that message
values become canonical ontos data that conform to the horos type.

## 4. One value through every layer

`ContentValue` exercises most of the families at once. Suppose Content declares
it in horos. Recursion must be written with `rec` inside one expression
(section 3.2). In the descriptor notation of section 3.2, with `H` for a
descriptor constructor and `A` for an atom:

```text
ContentValue = H("rec", A("value"),
  H("union",
    H("atom"),
    H("list-of", H("ref", A("value"))),
    H("map-of", H("scalar", A("utf8-text")), H("ref", A("value")))))
```

This union happens to be structurally disjoint. A bare atom, a tuple labelled
`list` and a tuple labelled `map` never coincide.

**The canonical value** of the content value `{"a": bytes 01, "b": []}`:

```text
Tuple(Atom("map"),
  Tuple(Tuple(Atom("utf8-text"), Atom("a")), Atom(0x01)),
  Tuple(Tuple(Atom("utf8-text"), Atom("b")), Tuple(Atom("list"))))
```

**A candidate TypeScript type**:

```ts
export type ContentValue = Uint8Array | readonly ContentValue[] | { readonly [key: string]: ContentValue };
```

A `ReadonlyMap<string, ContentValue>` would serve equally well for the third
alternative.

**What conversion needs to know, and so what the plan must carry**, beyond the
type itself:

| Decision | For this value | Today decided by |
| --- | --- | --- |
| Recursion | The `rec` binder `value` becomes the native declaration `ContentValue`. An inner `rec` would need a fresh native name. | nobody (refused) |
| Union, decoding | atom → alternative 0, `list` tuple → 1, `map` tuple → 2 | nobody |
| Union, encoding | `instanceof Uint8Array` → 0, `Array.isArray` → 1, otherwise → 2. The order matters, since a `Uint8Array` is also an object. | nobody |
| Bytes | Copy on encoding, so that later mutation cannot change the message; produce a fresh buffer on decoding | nobody |
| Map form | Index-signature object or `Map`. Objects have hazards such as the key `"__proto__"`. | nobody |
| Map order | Encoding sorts by the codec order of the encoded keys, not by JavaScript order. Decoding rejects disorder and duplicates. | ontos enforces decoding |
| Sharing | A buffer referenced twice: encode once and share? Is sharing preserved on decoding? | Content's own codec |
| Cycles | A cyclic native object graph must be refused on encoding | Content's own codec |
| Limits | Depth, node and byte limits, counted over the expanded tree or the shared graph | Content's own codec |

**The integer in `check()`'s result** raises a separate choice. The horos type
is `scalar(int)`, which is unbounded. The native candidates are `bigint`, which
is lossless, and `number`, which is ergonomic but exact only up to 2^53 − 1.
Content's meaning is "a non-negative count", which horos cannot express.

## 5. Constraints

- **horos semantics are given.**
  - Unions are untagged set membership.
  - Integers are unbounded, and ranges are deliberately left to logos.
  - Named declaration cycles are rejected.

  horos changes may be recommended; its maintainers decide.
- **Message values are canonical ontos data conforming to the horos type**
  (research 0014). There is no other message form.
- **No silent narrowing.** A family must not quietly lose information or refuse
  valid values it was not declared to refuse. A refusal is an explicit,
  classified failure.
- **Same mechanism for everyone** (AC-04). The structural profile and library or
  domain mappings contribute through one public mechanism. The kernel learns no
  family by name. There is no catalogue of families in shared code.
- **The plan** consists of typed references and target-owned expressions, not a
  universal schema or a new type algebra. It is serialized and versioned, because
  the TypeScript backend runs in another process (Node.js).
- **bittype's native syntax covers types only**, and bittype does not own runtime
  adapters (0014). It owns the plan.
- **Targets.** TypeScript bindings come first and Go bindings next. Rust and C++
  types must remain possible. C++ has no ontos implementation, and Rust's ontos
  integers stop at `i128`.
- **Compatibility is not a constraint.** Today's native API, message forms and
  bytes may all change.
- **The capture and lifetime contract is pending** (section 2.5).

## 6. What we have considered

None of this is decided.

### 6.1 One contribution per family and target

A family contributes, per target, one unit with three parts:

- a **native facet**: the binding descriptions and recipes it already supplies;
- a **value facet**: a reference to the run-time conversion it requires;
- the **plan entries** it records for each use.

Our present leaning is that the plan entries are the **arguments of the value
facet**. For a record they would be the field names, native member names and
absence rules. For a union, the alternatives and their discrimination. For a
map, the key representation and the native map form. Their types would be owned
by the family that defines them. The plan would then consist of:
- a generic envelope: the selected contribution, its target, the source use, its
  children and the version;
- per entry, a family-owned payload.

The backend composes static calls of each family's run-time factory with these
arguments. Alternatives we considered:

- **The plan records target expressions for access and construction**, which the
  backend pastes into generated code. This is flexible, but it risks becoming a
  small code language inside the plan.
- **A shared structural interpreter** reads the plan at run time. 0014 allows it
  only if it executes the plan's decisions rather than choosing rules.

We are unsure who implements a structural family's run-time factory in
TypeScript. Candidates:
- the representation profile's owner, which today is bittype, but bittype owns no
  run-time code;
- the TypeScript backend's owner, the service library;
- a separate small run-time package per profile.

### 6.2 Integers

The options per target are:
- lossless types: TypeScript `bigint`, Go `*big.Int`, a Rust big-integer crate,
  a C++ library;
- bounded native types: `number`, `int64`, `i128`.

Our leaning:
- **The structural profile.** It uses only lossless types.
- **Bounded types.** A bounded type may appear only through an exact mapping
  owned by a library or domain that declares the narrower domain, for example a
  `Count` type that Content declares. Its decoding explicitly refuses values
  outside that domain.

Open points:
- **Where the domain is declared, since horos has no ranges.** It could be a
  logos claim, a domain library's own declaration, or a horos extension.
- **Ergonomics.** Whether `bigint` is acceptable at the TypeScript surface.
- **Mutability.** A Go `big.Int` is mutable, which touches capture.

### 6.3 Bytes and atoms

- **The type.** `atom()` becomes `Uint8Array`, `[]byte`, `Vec<u8>` or
  `std::vector<std::uint8_t>`.
- **Copying.** Values are copied on encoding and on decoding.
- **Refinements.** A length bound, an exact value or a finite set of byte
  strings has no native type form. Our leaning is to check refinements by horos
  conformance at the boundary, not to model them natively. A finite set of
  *text* values, such as `ContentRepresentation`, can become a native literal
  union.

### 6.4 Tuples, sequences, labelled compounds and `any`

- **Tuples.** `tuple(...)` becomes a native tuple. Go has none, so it would be a
  struct with positional fields.
- **Sequences.** `seq(T)` resembles a list but is a bare tuple, so it needs its
  own native form or a refusal.
- **Labelled compounds.** `compound(label, ...)` becomes a record-like native
  type whose label is fixed.
- **`any()`.** It becomes the ontos value type, as an external library type.

### 6.5 Maps and sets

- **Text keys.** Maps keyed by text, the common case, become a native map or an
  index-signature object.
- **Key equality.** Native maps compare keys by native equality. That is correct
  only for key types whose native equality is canonical equality: text,
  booleans, and `bigint` in TypeScript. Byte strings and records as keys are not
  (JavaScript compares objects by reference; Go cannot use `[]byte` as a key).
  Our leaning:
  - such keys need a library-owned canonical map keyed by encoded bytes;
  - otherwise they are refused for bindings.
- **Order.** Encoding always sorts by the codec order of the encoded keys.

### 6.6 Unions

Our leaning:
- **When to bind.** Bind a union only when its alternatives are provably
  distinguishable on both sides:
  - on the canonical side, by a structural test that picks exactly one
    alternative;
  - on the native side, by a test over the native representations.

  First-success matching is rejected.
- **Who judges.** Distinguishability on the canonical side is a judgment about
  horos membership. It belongs to horos's interpretation, as a source judgment,
  not to the kernel or a target.
- **Common tags we would recognize:**
  - different labels of `compound`;
  - different canonical forms (atom, list, map);
  - an exact-text field shared by all alternatives, as in Drafts' `kind`;
  - `union(T, exact(null()))` where T excludes null.
- **Overlapping unions** are refused for bindings unless a library-owned
  representation supplies the missing information (0014).

We are unsure how far a general disjointness judgment can go. It is easy for the
cases above and hard in general (`and`, `not`, `any`, open records).

### 6.7 Absence, null and open records

- **Optional fields** stay as they are (section 3.4).
- **Nullable fields** are `union(T, exact(null()))`, whose native form is
  `T | null`, `*T`, `Option<T>` or similar. They must stay distinct from absence,
  so an optional nullable field has three states.
- **Open records** carry extra fields of any value. A native form needs a place
  for them: a separate map member, or refusal. Our leaning is to refuse open
  records for bindings until a representation is chosen deliberately.

### 6.8 Recursion

- **Naming.** Each `rec` binder becomes a named native declaration. A binder that
  is a declaration's whole body takes the declaration's name. Inner binders need
  a naming rule, and the naming must avoid collisions.
- **Indirection** is chosen per recursive edge by the target's profile:
  - none in TypeScript;
  - Go: a pointer, or none through a slice;
  - Rust: `Box`, or none through `Vec`;
  - C++: `std::unique_ptr`, or `std::vector` where an incomplete element type is
    allowed.

  The plan records each recursive edge and its indirection.
- **Cycles and sharing.**
  - Native values may be cyclic or shared, and canonical values are trees.
    Encoding refuses cycles.
  - Sharing is unobservable in ontos, so the round-trip law `D(E(n)) ≈ n` would
    ignore sharing.
  - Content, however, relies on DAG compactness for resource bounds.

  Whether sharing is part of a family's contract or an implementation property is
  open. The same goes for whether limits count the expanded tree or the shared
  graph.

### 6.9 Qualification

Each family and target is qualified for bindings separately from types:
- by the laws of section 2.3, with each law shown failing when its rule is broken
  on purpose;
- by the 0014 completeness test, in which a test backend implements each family
  from the plan alone.

TypeScript comes first. Families arrive in the order Content needs them:
1. integers;
2. bytes;
3. text-keyed maps;
4. disjoint unions;
5. recursion.

### 6.10 Claims we want tested

1. "A structural profile should never narrow a domain. Only an exact library or
   domain mapping that declares the narrower domain may."
2. "Whether a union's alternatives are distinguishable is a source judgment
   (horos), not a target judgment."
3. "The plan is the selected contributions plus their contribution-owned
   arguments, and nothing else."
4. "Sharing is never part of a family's semantic contract."

## 7. Questions for the expert

1. **The unit of a family.** What should one representation family contribute,
   as a unit, so that its native type, its plan entries and its conversion are
   defined once and always agree? How would you shape the plan entries, such as
   absence rules, union discrimination, map key handling and recursive
   references? They must let a backend in another process carry the decisions
   out without re-deriving a rule, while the plan stays typed references and
   target-owned expressions rather than a universal schema.

2. **Where native and canonical domains differ.** Integers are the clearest case:
   horos integers are unbounded and horos has no ranges by design, while
   ergonomic native integers are bounded. Map and set keys are another: ontos
   orders keys by their canonical encoding, and native maps use native equality
   and order. How should a family decide between a faithful representation and
   an ergonomic narrower one? Who may declare a narrower domain, and how should
   conversion treat a valid value outside it?

3. **Unions.** horos unions are untagged and their alternatives often overlap.
   When can a union be bound to a native union, and which component should own
   the judgment that its alternatives can be told apart, on the canonical side
   and on the native side? How would you treat unions that are only partly
   distinguishable?

4. **Recursion and sharing.** horos recursion is anonymous and values are finite
   trees that may be shared internally. Native targets need named recursive
   declarations, sometimes explicit indirection, and their values can be cyclic
   or shared. How would you divide naming, indirection, cycle refusal, sharing
   and resource accounting between the source, the family and the backend?

5. **Open.** Looking at the whole design, which choices here are most likely to
   force a refactor later, and what would you do differently? Patterns from
   comparable systems are especially welcome: code generators for schema
   languages with unbounded integers, untagged unions or canonical map order,
   and what worked or failed for them.
