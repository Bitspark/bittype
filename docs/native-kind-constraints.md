# Preserving Horos kind constraints in native types

**Proposal, 5 November 2026. Documentation only; no enforcement implementation.**
The [paired consultation synthesis](generic-elaboration.md) states the current
conceptual direction and admission promises. Consultations 0008/0009 exposed
flaws in the earlier TypeScript/Go evidence sketches; the target sections below
replace those claims. No candidate is compiler-qualified by this documentation.

The owner's [Logos preservation proposal](logos-preservation.md) develops an
explicit proof layer for these guarantees. It includes scoped facts and rules
for a symbolic pagination proof under `T : Data` and representation assumptions.
Native admission remains a separate theorem; the target-specific admission
candidates below do not bound that generic preservation theorem's domain.
Runtime value conformance remains separate from both.

The [type-system modeling proposal](type-system-model.md) specifies the common
structure for the language theories needed to derive those guarantees. Its
[five-language discussion](type-system-model.md#how-the-five-languages-fit)
uses the target-specific findings below as requirements, preserving their
different judgments, checking points and supported domains.

The [evaluated implementation design](type-system-architecture.md) keeps those
models language-owned and uses Logos's semantic-package/checking direction.
It separates conditional Page preservation from native admission, invalid proof
from false claim, and complete modeled coverage from a few compiler examples.
No new Bittype type-system foundation or target enforcement is implemented.

This document describes how TypeScript, Go, Rust and C++ generation could retain
Horos kind constraints without recognizing `Data`, `Type`, `Scalar` or any other
named kind specially. It extends the proposed [native representation design](type-representations.md)
and uses the existing [semantic transformation boundary](semantic-generation.md).

The companion [type-owned generator plugin design](type-generation-plugins.md)
owns declarative mapping selection, requested public roots and owned declaration
recipes. This document owns the guarantees and limits of constraint encodings. Their
[composition contract](type-generation-plugins.md#constraint-preservation-in-plugin-composition)
keeps those responsibilities together without adding kind-based dispatch or
claiming that the implemented mapping SPI already implements witnesses.
[GEN-02](generation-gaps.md#gen-02) records the missing construction contract,
current guarantee and completion criteria; source-reader correction is separate.

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

The [implemented composer](generation-composition.md) separates source
assumptions, representation conventions and native construction. The Horos-owned
representation policy currently requests bare native parameters, so TypeScript
emits no kind constraint.
The proposal requires source-to-target lowering and native-plan expressiveness
for the selected enforcing encoding. Today `NativeParameter` carries only a name
and origin. Its generic argument templates alone cannot express native bounds.
The exact enforcement contract is not finalized; any general extension must be
justified by that capability, without named-kind handling in shared orchestration.

The private fixed Dialect has been removed. Source-owned integration mappings
and source/target collaborators are the implemented generation path; native
enforcement must be supplied through that same path.

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

A proposed witness can identify a source expression, its selected association,
or a source expression plus a separately indexed representation family. The
choice must be explicit; a source-only identity with one fixed projection cannot
serve incompatible representations simultaneously. For example, PersonType can identify Person under a selected
association, while
Person is its native payload representation.
A page's items still contain native person values, not witness objects.

Membership evidence binds a source subject, its kind obligation and the selected
representation. The subject includes the relevant application arguments and
binders; display names or equal-looking native layouts are insufficient.

An override may preserve the source fact `S : K` while changing `S -> N`.
Evidence and native bounds tied to the old representation cannot simply be
reused for the new one. The [plugin composition contract](type-generation-plugins.md#overrides-and-external-representations)
requires a compatible association under the selected configuration; a registry-qualified
binding reference alone establishes no semantic guarantee.

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

The [single generic convention](generic-elaboration.md#one-generic-convention-for-formals-and-applications)
replaces the earlier binder-to-slot sketch. Its declaration and application
templates share a convention for native formals and complete actual arguments.
A stronger enforcing profile must additionally retain payload/result support. The complete support API and its target-specific authority remain
unfinalized. Additional native parameters create no additional source binders,
and a successfully reserved helper is not membership evidence.

Adding a source type or kind requires ordinary source input and regeneration,
not a handwritten membership catalogue. An arbitrary independently authored
native type needs an explicit source association before it can participate in
this encoding. Native shape inspection is not a substitute.

## 4. Pagination through an ordinary recursive kind

The [generation walkthrough](generic-elaboration.md#pagination-and-scope) follows
this example through parameter introduction, application and nested result
support. Its [implementation sequence](generic-elaboration.md#implementation-sequence)
retains the current representation-only guarantee until an explicit native
enforcing encoding is supported.

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

An upper bound on a separate proof argument does not establish issuance:

```typescript
// Rejected as a general enforcement scheme; support definitions are schematic.
type CheckedPage<S, N, P extends Evidence<S, N>> = Page<N>;
type Wrong = CheckedPage<WireTextType, number, never>;
```

Because `never` satisfies the bound, the resulting Page<number> may be inhabited
without evidence for the intended source tuple. Invariance does not repair this
bottom-type case. An outer `never` result is normalization, not formation
rejection; Page<never> can still contain an empty list. These are deductions from
the [documented bottom-type rules](https://www.typescriptlang.org/docs/handbook/release-notes/typescript-2-0.html),
not reported compiler runs.

Property-based structural evidence also permits mapped-bottom constructions.
Private-symbol spelling and readonly payload fields do not establish controlled
issuance or invariance. Witness unions must preserve whole-application choices:
Page<string> | Page<boolean> differs from Page<string | boolean>. See
[mapped types](https://www.typescriptlang.org/docs/handbook/2/mapped-types.html),
[strict function variance](https://www.typescriptlang.org/tsconfig/strictFunctionTypes.html)
and [conditional distribution](https://www.typescriptlang.org/docs/handbook/2/conditional-types.html).

Start with source-checked named applications. A stronger candidate is a finite,
non-augmentable selector of complete applications, with unions meaning choices
and bottom the empty selection. Payloads come from that configuration-specific
association, not an independently supplied native type. State strict settings,
consumer operations and excluded forging/any behavior. This is bounded
normalization; a requested open formation-enforcing API can remain Unsupported.

## 6. Go

**The former private singleton-evidence proposal is withdrawn.** An exported
constraint containing one exact private generic instantiation can permit that
instantiation to be inferred even when no alias for it was exported. A consumer
generic function assignment can infer the missing evidence argument, while an
alias erases it from the resulting payload. Private naming does not establish
selective issuance. This follows from the
[type-inference rules](https://go.dev/ref/spec#Type_inference) and is recorded as
specification/source reasoning, not an executed compiler example.

A candidate repair must restrict the admissible relation to actually issued
source/representation tuples; an unauthorized private instantiation must be
absent from that relation. The generic spelling of a finite intersection remains
unqualified. A private-method interface does not automatically solve type-only
issuance either.

The [compositional certificate proposal](go-native-admission.md) adds a second
restriction: evidence must have underlying type `func(S) N` as well as the
issuer's private method. Its Page certificate family requires valid child
evidence and supports nested pagination without a finite union. The source/native
indices prevent reusing a certificate for another pair; the function shape
excludes struct method-promotion attempts. Compiler fixtures cover the concrete
encoding, aliases, inference and forgery. These tests do not establish source
membership, representation laws or the complete native-model admission lemma.

Use distinct defined descriptor identities, not aliases of the same empty
struct. Exact payload constraints and underlying-type constraints (`string`
versus `~string`) admit different domains; neither is a general recursive record
predicate. Generic aliases can support admission followed by identity erasure.
See [type identity](https://go.dev/ref/spec#Type_identity) and
[Go 1.24 generic aliases](https://go.dev/doc/go1.24).

Source-checked named applications remain the implemented floor. The certificate
family is a proposal with a controlled issuing package; it is not an implemented
generation profile. Independent cross-package issuance remains deferred. A true
alias reexport retains the original issuer; a new issuer cannot adopt its private
method merely by copying its spelling.

## 7. Rust

Separate three mechanisms: ordinary alias bounds, demanded associated-type
projections, and nominal wrappers. Ordinary alias bounds do not supply a reliable
use-site gate under the baseline semantics. A projection through a constrained
implementation is a different candidate: selecting its output can demand the
argument obligation even if that output ignores the argument. See
[alias-bound limitations](https://doc.rust-lang.org/rustc/lints/listing/warn-by-default.html#type-alias-bounds)
and [associated types](https://doc.rust-lang.org/reference/items/associated-items.html).

Promise admission at specified well-formedness or demanded-normalization points,
not every unused alias mention. The compiler guide explicitly describes
[incomplete alias checking](https://rustc-dev-guide.rust-lang.org/analysis/well-formed.html).
A wrapper can retain association but changes identity/access and cannot silently
replace an alias-interchangeability promise.

A controlled-descriptor candidate must enforce Sized at the actual gate/carrier
or as a required supertrait; an alias's bound alone is insufficient. A dyn-compatible
membership trait can have a nameable unsized trait-object type satisfying it
without the intended concrete issuance; the effective Sized gate excludes that
route. Keep membership separate from configuration-selected representation. Sealing, orphan
rules, trait objects and overlapping blanket implementations constrain extension.
See [sealed traits](https://rust-lang.github.io/api-guidelines/future-proofing.html)
and [coherence](https://doc.rust-lang.org/reference/items/implementations.html#trait-implementation-coherence).

Controlled sized descriptors with generated facts are a candidate domain, not
an implementation of arbitrary coinductive Horos reasoning in Rust's solver.

## 8. C++20

Constrained alias templates can preserve an argument gate even for a constant
result, then erase separate source identity. Dependent constraints and requires
expressions have their own checking points; an unsatisfied requirement need not
make the whole translation unit fail. See
[constraints](https://timsong-cpp.github.io/cppwp/n4861/temp.constr) and
[alias templates](https://timsong-cpp.github.io/cppwp/n4861/temp.alias).

Concepts cannot be specialized, but supporting class/variable templates can.
Private placement alone is not a sound positive-membership authority, because
explicit specialization has special access rules. A finite candidate uses exact
identities of complete nongeneric generated descriptors under explicit definition
and customization contracts. See
[concepts](https://timsong-cpp.github.io/cppwp/n4861/temp.concept) and
[specializations](https://timsong-cpp.github.io/cppwp/n4861/temp.spec.general).

Required declarations must be visible consistently before use; constraint
satisfaction cannot change through late certificate registration. See
[atomic constraints](https://timsong-cpp.github.io/cppwp/n4861/temp.constr.atomic).
Generic descriptor families need a separately justified authority/extension model.
No native recursive template machinery automatically supplies Horos membership.

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
Its source adapter and type-owned integration mappings now own lowering;
`internal/horosgen` is deleted. Enforcement extends these implemented owners:

| Owner or location | Responsibility |
| --- | --- |
| `internal/horosmodel/` and Horos Files providers | Expose required Horos-owned contracts, memberships and assumptions through existing semantic terms; add explicit semantic obligations if needed |
| `transform/syntax/horos/source/go/` | Read the supplied semantic evidence, preserve environments and expose source obligations without choosing native representations |
| Horos-owned and other type-owned integration mappings | Construct the selected native bounds, witnesses and evidence associations; compose child mappings and reject unsupported encodings |
| `transform/syntax/shared/generation/go/` | Select descriptions, compose closed native intent and preserve requested roots/helper ownership; interpret no Horos kind or canonical type |
| `transform/syntax/<target>/target/go/` | Supply native construction, symbols and assembly; interpret no canonical membership rule |
| `transform/syntax/shared/generation/go/` | Admit independently registered source/target/mapping contributions and plan the requested roots; create no ASTTransformer edge. |
| `syntax/<target>/go/` and its parse/check/format Syntax stages | Support the native grammatical constructs required by the encoding, without Horos semantics |
| Shared carriers, semantic kernel, engine and printers | Keep their existing language-agnostic contracts |

These are the [plugin proposal's target locations](type-generation-plugins.md#file-layout),
implemented for representation, with enforcement still proposed. Type-owned integration modules
depend on Bittype; Bittype factories do not import their default bundles. The
private fixed Dialect is deleted.

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
Likely additions include TypeScript's selected nominal support forms,
Rust implementation items beyond type-only trait impls and any required
representation attributes, and C++ concepts, requires-expressions and selected
trait/specialization constructs. Go's existing generic/type-set vocabulary is a
starting point, not a claim that the evidence encoding is already delivered.
TypeScript conditional/mapped types and infer binders are available in the
[native syntax tooling](syntax.md#typescript); their evaluation and preservation
claims still require a separate semantic integration.
Rust qualified projections, associated-type equality arguments, generic
associated declarations and type-only trait impls are available
in the [native syntax tooling](syntax.md#rust); this supplies syntax support,
not the constraint-preservation analysis proposed here.
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
| How do separately generated packages/crates/modules compose evidence? | Consumer-owned composition and explicit dependency contracts are proposed; unrestricted cross-package issuance remains unqualified |
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
- Claiming formation rejection while returning `never` for a failed application;
  hiding failure through a permissive fallback, implicit specialization or removal
  of an unsupported requested type. An explicitly declared finite selector may
  instead give bottom the distinct meaning of an empty selection.
- Reintroducing template lifting, handwritten operation mappings, runtime
  membership registries, serializers, value validators or protocol behavior.

Function, Stream and Cell remain ordinary library declarations with library-owned
native representations. A library representation may participate only when its
source contract and selected target encoding support the same static obligations.
No new runtime dependency follows from this design.
