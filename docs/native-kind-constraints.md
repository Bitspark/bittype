# Preserving Horos kind constraints in native types

**Proposal, 5 November 2026. Documentation only; no enforcement implementation.**
This document describes how TypeScript, Go, Rust and C++ generation could retain
Horos kind constraints without recognizing `Data`, `Type`, `Scalar` or any other
named kind specially. It extends the proposed [native representation design](type-representations.md)
and uses the existing [semantic transformation boundary](semantic-generation.md).

The companion [type-owned generator plugin design](type-generation-plugins.md)
owns exact binding selection, overrides and Use/Publish/Define. This document
owns the guarantees and limits of constraint encodings. Their
[composition contract](type-generation-plugins.md#constraint-preservation-in-plugin-composition)
keeps those responsibilities together without adding kind-based dispatch or
claiming that the proposed plugin SPI already implements witnesses.

The intended rule is: emit an encoding that preserves the requested constraint,
or return a located `Unsupported` diagnostic. Never silently erase a constraint
in an output advertised as enforcing it. The target language, representation
domain, supported kind fragment and public API shape must be explicit.

## Current implementation and proposed change

The [Horos frontend](horos.md) already supplies general kind contracts and
membership facts, including the assumptions under which they hold. The four
pagination mappings currently emit ordinary native generics under the
`pagination-representation/1` profile, which expressly permits constraint erasure.
The `pagination-enforcement/1` profile currently returns `Unsupported` with no
output. Neither the encodings nor the helper names below are implemented APIs.

The present private dialect operation receives a parameter name and origin,
without a constraint plan; see [the generator](../internal/horosgen/generate.go).
Consequently, for example, the TypeScript emitter produces a bare parameter.
The proposal changes that source-to-target lowering, not `ASTTransformer`, the
engine's registry or the shared semantic carrier.

The plugin proposal replaces the private fixed Dialect with source-owned
integration mappings and source/target collaborators. The existing private
generator is an implementation anchor, not a second planned enforcement route.

Implementing this proposal would replace constraint erasure in the enforcing
route. There must be no compatibility fallback to the unconstrained output.
This documentation change does not alter either existing profile's behavior.

## 1. What must survive generation

A kind classifies a type description or a type constructor. A native type
usually describes the representation of values. These are different levels.

For example, assuming an explicit import of the ordinary `horos.kinds` module
as `h`:

```text
type TextBox<T: h.Text> = { value: T; };
```

The current representation mapping produces the equivalent of:

```typescript
export interface TextBox<T> {
  value: T;
}
```

Native callers can therefore name `TextBox<boolean>`. Emitting
`T extends string` would retain a useful native restriction, but it is not a
general translation of a Horos kind. It does not establish the source meaning
of every native subtype of `string`.

A second example exposes the distinction:

```text
kind Leaf = h.Text | h.Bool;
```

`Leaf` admits a text description or a boolean description. A descriptor whose
root is `union`, containing those two descriptions, does not belong to `Leaf`.
The TypeScript constraint `T extends string | boolean` admits the native union
`string | boolean`. Translating a kind union into a union of native payload types
would therefore lose the source distinction.

The same problem occurs when different descriptor forms select identical native
representations, or when native aliases erase a distinction. A representation
recipe is not evidence that every source description with that representation
belongs to the same kind.

The required preservation relationship is:

```text
S = the source type expression, with its bound environment
K = the resolved source kind
N = the selected native representation of S

an admitted native application must retain an established S : K
and the selected association S -> N
```

For an open generic, the relationship retains its assumptions. A body proved to
preserve `K` when its argument satisfies `K` must not acquire unconditional
membership. Parameter obligations, including ignored parameters, and result
contracts both matter. A stored result contract is not permission to admit an
argument outside the original domain.

This is a static type-argument guarantee. It does not establish that arbitrary
native values conform to Horos refinements, closed-record rules or transport
contracts. Runtime conformance and hydration remain separate responsibilities.

## 2. Ownership and name-independent compilation

Horos owns formation, kind membership, inclusion and generic reasoning.
Bittype's Horos provider exposes the established facts through the existing
semantic model. Concrete transformations choose target encodings for those facts.

| Source information | Treatment |
| --- | --- |
| Named reference or kind application | Follow the resolved declaration and bindings; preserve recursive references instead of indefinitely expanding them |
| Top, bottom and sort | Retain the proper/constructor distinction; native `any`, `never` or an empty interface is not automatically an equivalent encoding |
| Kind union or intersection | Preserve alternatives or joint obligations on source subjects; do not substitute native value unions mechanically |
| Node pattern and selectors | Use the source form, permitted metadata selectors and designated type children |
| Child sequence pattern | Retain item, concatenation, alternative and finite repetition obligations |
| Guarded recursive group | Use Horos's established recursive judgment; ordinary target recursion is not automatically the same greatest fixed point |
| Constructor contract | Preserve input admission and result obligations; unsupported higher-order encodings fail explicitly |

All spellings in the prelude are ordinary declarations. Renaming `Data`, importing
an equivalent kind or defining an unrelated kind called `Data` does not select a
different algorithm. The Horos-specific lowering necessarily understands the
Horos kind algebra and primitive observation vocabulary. Those concepts do not
become cases in the agnostic engine, shared AST/CST carriers or native printers.

Kind-pattern metadata selectors are different from the plugin design's opaque
mapping Selector. The former describe source observations; the latter selects
an exact builtin or declaration binding. Membership is an eligibility obligation,
never a most-specific-kind mapping search. Direct bounds versus witnesses is an
explicit representation choice, not a retry after the selected mapping fails.

The target output encodes established relationships and supported generic
implications. It does not ship another Horos evaluator, invoke a runtime registry
or discover memberships through reflection.

If extra judgments are needed, request them through the Horos semantic-provider
composition and preserve their assumptions in the supplied evidence. A
transformer must not silently rerun resolution, query an unrelated Program or
manufacture a fact absent from the captured semantic input. `Unproved`, invalid
input and resource exhaustion do not authorize a permissive encoding.

## 3. Two explicit encoding choices

### Direct native bounds

Use ordinary target constraints when a defined lowering rule preserves the
required source distinction within the selected representation domain. Native
syntax alone does not establish that rule.

Possible mechanisms are TypeScript `extends`, Go interface type sets, Rust
traits and C++ concepts. A simple scalar restriction is a candidate for this
encoding. Mappings must account for aliases, native subtyping, unions and any
other supported representation that could satisfy the native bound.

If two supported source types become indistinguishable to the target constraint
but differ in membership, a direct bound cannot enforce that distinction. The
compiler must retain more information or reject that encoding. It must not
claim that every kind can be translated to a bound on an arbitrary native type.

### Source type witnesses and membership evidence

A source type witness is a generated compile-time type identifying a source
expression and its selected native representation. For example, `PersonType`
identifies the source `Person`, while `Person` is its native value representation.
A page's items still contain native person values, not witness objects.

Membership evidence binds a source subject, its kind obligation and the selected
representation. The subject includes the relevant application arguments and
binders; display names or equal-looking native layouts are insufficient.

An override may preserve the source fact `S : K` while changing `S -> N`.
Evidence and native bounds tied to the old representation cannot simply be
reused for the new one. The [plugin composition contract](type-generation-plugins.md#overrides-and-external-representations)
requires a compatible association under the selected configuration; BindingID
or PluginID alone establishes no semantic guarantee.

There are two distinct products:

- A concrete membership established for a source type.
- A generic evidence family whose prerequisites correspond to the assumptions
  of a Horos-established implication.

An evidence family must constrain every prerequisite. An unrestricted public
`Proof<S, N, K>` type that callers can instantiate for arbitrary arguments would
assert the conclusion instead of enforcing it. Empty marker interfaces and
user-writable boolean properties have the same problem.

The target encoding must also tie the evidence to the representation invariantly
where necessary. Evidence for one source/representation pair cannot be widened
into evidence for another. Merely retaining a phantom name does not establish
this property in every language.

Witness output changes the public generic API. A source parameter may become a
descriptor parameter, or several linked parameters, rather than one arbitrary
native type. The profile must state this change. Generated named applications
can hide the extra parameters without duplicating the generic definition.

The initial plugin SPI does not yet specify this binder expansion. Its BinderLink
pairs a real source binder with a reserved target slot; it does not define a
tuple of source witness, native representation and evidence, or a projection such
as Native<S>. A future encoding must specify the same header and argument
transformation for Use, Publish and Define, including substitution of assumptions
and result contracts. Extra native parameters are not new canonical binders.
The sketches below do not claim that reserving extra slots completes that design.

Adding a source type or kind requires ordinary source input and regeneration,
not a handwritten membership catalogue. An arbitrary independently authored
native type needs an explicit source association before it can participate in
this encoding. Native shape inspection is not a substitute.

## 4. Pagination through an ordinary recursive kind

The following source excerpt assumes `h` is explicitly imported. It uses a
custom kind to make the absence of named-kind cases visible:

```text
kind Tree = nu {
  X = h.Text | h.Bool | h.List<X> | h.ClosedRecord<X>;
} in X;

type Page<T> = {
  items: list<T>;
  next?: text;
};

type TreePage<T: Tree>: Tree = Page<T>;

type Person = {
  name: text;
  active: bool;
};

type People = TreePage<Person>;
```

The source reasoning establishes Person's membership from its field types. It
establishes TreePage's result membership under the assumption that its argument
belongs to Tree. Its `items` field is a list of that argument and its optional
`next` field is text. Optional fields participate in the kind judgment.

| Source application or change | Required treatment |
| --- | --- |
| `TreePage<Person>` | Retain Person's membership and its selected native representation |
| `TreePage<list<Person>>` | Compose established membership with a supported list representation |
| `TreePage<wire<text>>` | Reject at source admission: the Wire description is outside Tree |
| A union descriptor over text and bool used as the argument | Reject at source admission: Tree has no union-root alternative |
| Rename Tree or import it through another alias | Preserve the same semantics through resolved references |
| Replace the bound with ordinary `h.Data` | Use that definition's full kind graph; do not substitute Tree or a smaller JSON-like category |

These rows describe intended semantics, not results of executed examples.
The native helpers in the following sections are proposed shapes, not complete
implementations or a claim that every source form already has a native mapping.

## 5. TypeScript

Direct constraints can use `extends` where the representation rule supports it.
Conditional and mapped types may support more elaborate encodings, but a
structural payload type is not by itself a source witness.

A proposed witness-facing API is:

```typescript
export type TreePage<S extends TreeMember> = Page<Native<S>>;

export type People = TreePage<PersonType>;
```

`Native<PersonType>` would select `Person`. `TreeMember` would restrict admitted
source witnesses through generated support types. It is not a preinstalled
interface and must not mean an object with a writable `isTree` member.

This sketch leaves the exact enforcement construction unresolved. In particular:

- Private symbols and structural brands do not automatically make arbitrary
  type-only evidence impossible to manufacture or widen.
- A TypeScript union of possible witnesses is different from a Horos union node.
  Projecting it to a union of payload types can lose correlations. For example,
  `Page<string> | Page<boolean>` and `Page<string | boolean>` permit different
  collections. The encoding must preserve the distinction or reject that use.
- Conditional types can distribute over unions. Distribution must be deliberate
  at the whole application boundary, not an accidental broadening of a field.
- Mapping an invalid application to `never` does not reject the type expression
  itself. A profile promising formation rejection cannot advertise that as an
  equivalent result.
- Recursive aliases, compiler resource limits, `any`, bottom types and native
  assertion escape hatches need explicit treatment in the supported profile.

Start with a bounded generated witness domain and supported generic composition
rules whose guarantees can be stated precisely. A finite membership enumeration
must be described as finite; it is not an implementation of arbitrary open
recursive membership. Do not claim general TypeScript enforcement from the
signature sketch alone.

References: [generic constraints](https://www.typescriptlang.org/docs/handbook/2/generics.html),
[conditional types and distribution](https://www.typescriptlang.org/docs/handbook/2/conditional-types.html),
and [native type compatibility](https://www.typescriptlang.org/docs/handbook/type-compatibility.html).

## 6. Go

Direct bounds use interface type sets where their meaning fits the selected
representation. For example, an exact `string` term differs from `~string`, which
also admits user-defined types with that underlying type. That widening must be
justified rather than inserted for convenience.

Go interface constraints cannot recursively contain themselves, and do not
provide arbitrary recursive predicates over every field of an unknown struct.
The full ordinary Data definition cannot simply be rendered as one recursive
interface constraint on arbitrary native types.

The proposed witness encoding therefore makes the native representation and
membership parameter explicit:

```go
type TreePage[
    S, T any,
    P TreeEvidence[S, T],
] = Page[T]

type People = TreePage[PersonType, Person, PersonInTree]
```

A candidate generated-package mechanism is an exact type-set constraint over a
package-private evidence constructor:

```go
type evidence[S, T, K any] struct{}

type TreeEvidence[S, T any] interface {
    evidence[S, T, treeKind]
}

type PersonInTree = evidence[PersonType, Person, treeKind]
```

The example assumes the generated package also declares its source token and
kind token. `PersonInTree` is emitted only for the established association.
Different argument tuples denote different instantiations of the private named
type. Native callers can use exported evidence aliases, but there is no exported
unrestricted constructor for arbitrary tuples. No method or runtime evidence
field is required in Page.

Generic evidence aliases or families would preserve all their own prerequisite
constraints. A rule for list composition or Page's result must come from an
established Horos implication, not from a switch recognizing either declaration's
name. An unused evidence argument still has to satisfy its bound.

Initially this mechanism applies within one generated module closure/package.
Cross-package evidence composition remains unresolved: two private constructors
in different packages are not interchangeable. Exporting the unrestricted
constructor or replacing the exact constraint with an empty interface would
defeat the design. A later package composition contract must preserve the same
relationship without manual grants or a runtime registry.

Reference: [Go general interfaces and type-set restrictions](https://go.dev/ref/spec#General_interfaces).

## 7. Rust

Use trait bounds and associated representation types. Membership implementations
belong to a controlled generated domain and are emitted from source facts or
supported conditional implications. A candidate public shape is:

```rust
pub struct TreePage<S: TreeMember> {
    pub page: Page<S::Repr>,
}

pub type People = TreePage<PersonType>;
```

`TreeMember` would associate the selected `Repr` and require the generated
membership relationship. Appropriate sealing and ownership of the descriptor
types must prevent an unrelated native implementation from merely asserting that
relationship. Phantom type storage, if needed for a selected encoding, is part
of the representation choice and does not add runtime behavior.

The wrapper is an observable API choice: native users access the payload through
`page`. It must not be substituted for an alias without the profile stating it.
Ordinary Rust type-alias bounds are not enforced at use sites, so writing
`type TreePage<T: TreeMember> = Page<T>` would not solve the problem. A profile
must use an enforcing declaration/encoding or reject the requested alias shape.

If the profile requires the source alias to remain natively interchangeable with
Page, a nominal wrapper must not claim to satisfy that promise automatically.
A profile permitting a distinct wrapper must state the resulting conversion
boundary; this proposal generates no conversions. If the requested constraint
and interchangeability promises cannot both be retained, generation fails.

Kind alternatives cannot mechanically become overlapping blanket trait
implementations. Generate a coherent membership relation for the supported
descriptor domain. Account for orphan rules and cross-crate ownership. General
greatest-fixed-point kind recursion is not automatically supplied by Rust's
ordinary trait solver; use established source judgments and supported generic
rules, or report the unsupported case.

References: [traits and associated types](https://doc.rust-lang.org/reference/items/traits.html),
[implementation coherence](https://doc.rust-lang.org/reference/items/implementations.html#trait-implementation-coherence),
and [type-alias bounds](https://doc.rust-lang.org/rustc/lints/listing/warn-by-default.html#type-alias-bounds).

## 8. C++

Use concepts and requires-clauses under an explicit target standard. A proposed
witness-facing alias is:

```cpp
template<class S>
    requires tree_member<S>
using TreePage = Page<typename S::representation>;

using People = TreePage<PersonType>;
```

The concept is backed by generated membership relationships for the descriptor
domain. It does not infer Data membership from trivial copyability, class layout
or a library name. Native type aliases can erase identity; distinct descriptor
tokens preserve source distinctions when those distinctions are required.

Concept conjunction/disjunction can combine supported obligations. Membership
traits and partial specializations must avoid ambiguity and preserve generic
prerequisites. Unknown descriptors have no positive membership by default.
Unrestricted user specializations cannot be treated as Horos-established facts;
the supported domain and extension rules belong to the profile.

Compile-time expressions used by concepts are static support definitions, not
runtime adapters. C++ template recursion alone does not establish Horos's
coinductive membership rules. General recursive cases and compiler resource
limits need the same explicit supported-fragment treatment as other targets.

Reference: [C++ template constraints](https://eel.is/c++draft/temp.constr), including
[conjunction and disjunction](https://eel.is/c++draft/temp.constr.op).

## 9. Pipeline and target file responsibilities

The pipeline remains:

```text
Horos declarations and supplied dependencies
    -> explicit resolution and analysis
    -> semantic obligations, established facts and their assumptions
    -> source-specific constraint/representation plan
    -> target AST with selected bounds or witness/evidence encoding
    -> target CST
    -> native source text
```

No new central provider registry or core semantic algebra is required. Concrete
type-generation plugins and library recipes, when configured, must preserve the
same obligations. This document does not define a competing recipe-dispatch SPI.
The [type-owned generator plugin proposal](type-generation-plugins.md) specifies
that selection/composition contract and the
[operation obligations](type-generation-plugins.md#constraint-preservation-in-plugin-composition).
Its source adapter and type-owned integration mappings take over the lowering
currently located in `internal/horosgen` and concrete transformers. The proposed
ownership is:

| Owner or location | Responsibility |
| --- | --- |
| `internal/horosmodel/` and Horos Files providers | Expose required Horos-owned contracts, memberships and assumptions through existing semantic terms; add explicit semantic obligations if needed |
| `generation/syntax/source/horos/go/` | Read the supplied semantic evidence, preserve environments and expose source obligations without choosing native representations |
| Horos-owned and other type-owned integration mappings | Construct the selected native bounds, witnesses and evidence associations; compose child mappings and reject unsupported encodings |
| `generation/syntax/go/` | Select configured bindings and coordinate Use/Publish/Define and helper ownership; interpret no Horos kind or canonical type |
| `generation/syntax/target/<target>/go/` | Supply native construction, symbols and assembly; interpret no canonical membership rule |
| `transforms/horos-to-<target>/go/` | Configuration-only factories for the composed transformer, with no builtin or library mapping defaults |
| `syntax/<target>/go/` and its parse/check/format Syntax stages | Support the native grammatical constructs required by the encoding, without Horos semantics |
| Shared carriers, semantic kernel, engine and printers | Keep their existing language-agnostic contracts |

These are the [plugin proposal's target locations](type-generation-plugins.md#proposed-file-layout-and-dependencies),
not directories created by this documentation. Type-owned integration modules
depend on Bittype; Bittype factories do not import their default bundles. The
private fixed Dialect is replaced rather than retained beside this composition.

Membership evidence may become a generated helper only after its source meaning
and native association are established. Reserving a helper symbol does not prove
membership. All uses of a generic must agree with its one published header,
including any supported witness and evidence parameters; a first concrete use
must not specialize that header silently.

Requested public roots must remain published or fail. Other representation
dependencies are demanded by the selected mapping: an external representation
need not generate its source body's native fields, but it must still satisfy
required source admission and the selected constraint promise. This does not
authorize dropping a failed requested root or retrying a weaker mapping.

The present native vocabularies do not yet cover every construct shown here.
Likely additions include TypeScript conditional types and the selected nominal
support forms, Rust implementation items/qualified projections and any required
representation attributes, and C++ concepts, requires-expressions and selected
trait/specialization constructs. Go's existing generic/type-set vocabulary is a
starting point, not a claim that the evidence encoding is already delivered.
All additions follow the [stage-first layout](../LAYOUT.md); no raw-source node
or emit-and-reparse shortcut is permitted.

## 10. Initial scope and unresolved limits

Begin with proper-type constraints and the existing pagination representation
subset: named closed records, text, bool, lists and optional fields. Cover ordinary
user-defined kinds as well as prelude definitions through the same semantic path.
Retain open generic definitions and explicit named applications.

Using `h.Data` does not require every member of Data to have a native mapping.
The supported representation domain may be smaller. Unsupported representations
remain errors; they must not cause Data itself to be redefined as a smaller kind.

Before any encoding is described as enforcing constraints, its contract must
settle these concrete questions:

| Question | Required disposition |
| --- | --- |
| Does rejection occur at type formation, or only when attempting to use a value? | State the promised level; do not equate an uninhabitable result with a formation error |
| Can the target erase or widen the source/evidence/representation association? | Preserve the association for the supported domain or reject that encoding |
| How do new source types and generic families enter the generated domain? | Derive evidence from ordinary source facts and prerequisites; no manual semantic grants |
| How do separately generated packages/crates/modules compose evidence? | Explicit ownership and identity rules; initially unresolved |
| Can a recursive or higher-order obligation be represented faithfully? | Establish a supported rule or return `Unsupported`; retain `Unproved` and exhaustion distinctly |
| Does the encoding change native generic arity, nominality, layout or access syntax? | Make that part of the selected representation profile |

The TypeScript witness construction, cross-package evidence, arbitrary recursive
kind encodings and higher-order constructor enforcement remain design work.
The sketches are not a universal encoding of the full Horos kind calculus.
This proposal also does not expand the current payload mappings to integers,
all ground descriptors, Wire or library runtime implementations.

## 11. Forbidden shortcuts

- Dispatching on the names Data, Type, Scalar, Tree, Function, Stream or Cell.
- Moving kind interpretation into generic AST/CST, the syntax engine, Wire,
  Bitruntime, Bitnode or hydration.
- Using a library recipe, a native class name or structural resemblance as
  evidence of source membership or behavioral identity.
- Emitting unconstrained parameters while advertising constraint enforcement.
- Exporting arbitrary evidence constructors, or treating marker interfaces as
  automatically sufficient evidence in all four languages.
- Discarding a generic rule's assumptions, ignored-argument obligations or result
  contracts; treating missing evidence or `Unproved` as successful admission.
- Equating native assignability with Horos kind inclusion, or a native value
  union with a kind union over source descriptions.
- Hiding failed admission behind `never`, an empty representation, a permissive
  fallback, implicit specialization or removal of an unsupported public type.
- Reintroducing template lifting, handwritten operation mappings, runtime
  membership registries, serializers, value validators or protocol behavior.

Function, Stream and Cell remain ordinary library declarations with library-owned
native representations. A library representation may participate only when its
source contract and selected target encoding support the same static obligations.
No new runtime dependency follows from this design.
