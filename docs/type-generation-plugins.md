# Type-owned generator plugins

**Declarative generation implemented, 8 November 2026.** Libraries own
native representation choices, including builtin mappings. The kernel selects
and composes descriptions; the target constructs the closed native program.
All four native generators use the [public SPI](generation-spi.md). The
procedural Use/Publish/Define and Introduce/Apply interfaces, products and tickets
have been removed. Their required responsibilities are expressed by mapping
rules, requested roots, declaration recipes and one generic convention.

The [paired consultation synthesis](generic-elaboration.md) retains the generic
coherence and native-admission requirements. Its earlier procedural surface is
not an alternative API. [ACX-005](architecture-cases/0005-generation-interface-design.md)
records why implementation reuse could not justify retaining that interface.

The later [Logos preservation proposal](logos-preservation.md) explains how
mapping-owned representation laws could compose into symbolic preservation
proofs under parameter assumptions, with native admission proved separately.
The support carrier and a plugin's declared capability alone do not
establish those claims; the proof must concern the actual selected output.

The [type-system modeling proposal](type-system-model.md) is now the design
prerequisite for that preservation-facing SPI. It separates language judgments
and rules from stored semantic facts and mapping laws. Its models remain
library-owned; selecting a mapping or constructing a product does not install
native semantics or prove a constraint-preservation claim.

The [evaluated architecture](type-system-architecture.md) places these theories
in language-owned Logos semantic packages. Generation returns a sealed output; explicit verification composition binds the actual output, selected
mappings and independently chosen statement to its checking result. Mapping
ownership supplies the law, not its truth. The scoped
[TextBox/Page integration](verification.md#implemented-go-surface) uses the landed
Logos typed checker; generation plans do not themselves certify those claims.

The baseline design is scoped to native types and interfaces. The subsequent
[Go admission proposal](go-native-admission.md) requires empty private marker
methods on certificate types: an explicit proposed
[scope extension](type-system-architecture.md#go-marker-method-scope-extension),
owned by the selected Go convention and emitted through native AST/CST tooling.
Production support for those method bodies is not implemented. Neither scope
includes calls, stream operations, runtime adapters, validators or transport code.
The expert consultation and its evaluation are recorded in
[research 0007](../research-docs/0007-type-owned-generator-plugins.md).

Read this together with [native kind-constraint preservation](native-kind-constraints.md).
This document owns mapping selection, generation operations and resource ownership;
the companion owns the guarantees and limits of native constraint encodings.
The [composition contract below](#constraint-preservation-in-plugin-composition)
connects implemented generation to still-proposed native enforcement encodings.

## Responsibility and ownership

A native representation has two axes: the source type and the target language.
For example, the Function library can supply a mapping to TypeScript's callable
syntax; the TypeScript syntax package owns that syntax independently of Function.
Another selected mapping can represent the same declaration as an external
handle. Neither choice changes the canonical Function declaration.

| Component | Responsibility |
| --- | --- |
| Source adapter | Interpret admitted source semantics, preserve binders, and supply exact mapping selectors. |
| Type-owned plugin | Describe the native representation of its builtin or named declaration and the child mapping demands. The kernel composes the descriptions. |
| Target-language tooling | The generation target constructs native declarations and bound type uses. The [Files formatter](source.md#from-generated-declarations-to-files) realizes imports and exports under the consumer's placement. Neither needs a canonical type catalogue. |
| Shared generation composition | Select configured mappings and coordinate invocation and ownership. It interprets neither source descriptors nor native types. |
| Consumer | Supply dependencies and select defaults, library plugins, explicit overrides and representation profiles. |
| Engine and generation registry | Register source, target and mapping libraries separately; plan requested roots and construct the selected native program. These registrations create no ASTTransformer edge. |

Builtins are ordinary supplied mappings. Bittype's Horos integration packages provide
text, boolean, list and record mappings; it is outside the language-neutral
generation kernel. Function, Stream, Cell and future libraries supply their own
optional mappings through the same contract. An ordinary record declaration can
also use the explicitly selected structural declaration policy without a custom
plugin.

A new representation of Cell changes its library integration and the consumer's
composition. It does not change the shared generator, ASTTransformer, engine,
Wire or runtime. A genuinely missing native grammatical form is different: its
language-owned AST/checker/formatter needs support before any plugin can emit it.

This boundary permits source-grammar dispatch in the Horos adapter and
native-grammar dispatch in a Go or TypeScript backend. A Go backend can know how
to construct a function type or channel type. The Function or Stream plugin owns
the decision to use that syntax for its canonical declaration. Moving that
decision into a general backend switch would violate the same ownership rule as
putting it in the kernel.

### Agnostic kernel and recursive composition

The implemented generator uses one language-agnostic composition kernel. Every native
representation rule is supplied through plugins, including builtin mappings and
the structural declaration policy. A language-specific generator becomes a
configured composition, with this division of work:

```text
agnostic generation kernel
  + source semantic adapter
  + selected builtin and library mappings
  + target syntax backend
  = one configured generation request
```

An individual mapping is selected for a canonical constructor or declaration and
a target language/profile. For example, the Horos list mapping for TypeScript
describes an array type; its Go counterpart describes a slice type. Each uses
a captured element slot under the same configuration. The shared
kernel coordinates selection, binding contexts and owned definitions, then asks
the target backend to assemble the output. It contains neither a canonical type
catalogue nor native representation cases.

An ordinary composed declaration does not require a custom mapping of its own.
Page can use the structural declaration policy, record mapping and list mapping.
Cell can likewise compose its fields using the selected Function mappings. A
library supplies a custom mapping when it wants a different native surface, such
as an existing Cell type. All these choices use the same plugin boundary.

Select a declaration's mapping before expanding its body. The selected plugin
chooses which source type operands it needs and requests their representations
through authenticated captures. An external Cell representation can use its argument without
expanding Cell's implementation shape; a structural Cell representation lowers
its fields. The kernel does not eagerly flatten every declaration into primitives
or require a complete child traversal before dispatch. Source adapters preserve
the actual type-operand and binding boundaries.

Here "kind-specific generator" means a mapping for a source constructor or named
type. Horos constraint kinds such as Data still govern applicability; their
overlapping memberships do not select a unique plugin. The same kernel is used
for TypeScript, Go, Rust and C++; their native grammar and scope rules remain in
their target integrations. The generic convention below supplies both the
formal and application templates.

### Canonical declarations, native types and later adapters

A library abstraction participating in cross-language mapping has a canonical
Horos declaration, such as Function<A,B> or Cell<T>. Its library owns the
behavioral contract as well. The native form is a selected representation of
that declaration. Other languages consume the same declaration through their
structural mappings or their own selected native/library representations.

Native syntax itself remains independent of Horos: a language AST can describe
a function or channel without importing a Horos library. The library plugin owns
the association between an exact canonical declaration and that native form.
There is no requirement for a Horos declaration corresponding to every native
AST node, or a one-to-one correspondence between canonical and native types.
Support in one language does not establish a faithful representation in another;
unsupported mappings remain explicit.

Later, a matching library-owned runtime adapter would connect a value conforming
to the canonical contract with the selected native representation:

```text
contract value of C<T> <-> selected native value of N<T>
```

For capability-containing values, decoding and hydration precede adaptation as
needed. A Function adapter could expose received contract endpoints as a native
callable facade; a reverse adapter could export a native callable through those
endpoints. Each supported direction must uphold the library's behavior and the
selected representation profile. Adapter composition supplies the interpretations
of generic arguments and fields, so a Page adapter can reuse its element adapter.

This is the future [runtime value boundary](runtime-values.md), not an adapter API
or implementation added by this proposal. Native type generation supplies no
invocation, cancellation or conversion behavior. Function, Stream and Cell remain
unknown to the general generator kernel, Wire and hydration machinery.

## Exact selection and explicit overrides

The source adapter produces opaque selectors qualified by the source language
and profile. Horos selectors distinguish exact builtin identities, including a
scalar's literal discriminator, from exact named declarations. A declaration key
contains the Horos module identity and declaration name. A graph-local EntityID,
syntax NodeRef, display name or structural fingerprint is not such an identity.

Kinds constrain applicability. They do not select mappings: a type may satisfy
several overlapping kinds, including Data. There is no most-specific-kind search.

Choosing direct bounds or witnesses is part of the explicitly selected
representation profile and mapping configuration. It is not another kind-based
dispatch pass, and failure to enforce a bound cannot select a weaker encoding.

The consumer registers mapping libraries; planning selects compatible offers
for the requested endpoints/profile and demanded occurrences. Registration order
is never a preference. More than one applicable demanded offer is ambiguous.
An exact use binding takes precedence over a structural declaration recipe; a
failure of that selected binding is final, with no structural fallback.

An explicit replacement names both actual registered offers. These are the
implemented composition contracts (the full registry owns MappingRef issuance):

```go
type BindingRef struct { Provider MappingRef; ID string }
type Replacement struct {
    Selector generation.Selector
    Expected, With BindingRef
}
func (r *Registry) Replace(Replacement) error
```

The application supplies the replacement; the registry authenticates its handles
and the planner checks it against applicable offers. Both bindings retain the
same exact selector. Replacements cannot form an order-dependent chain.

A binding's ID is local to its mapping registration. It is not a semantic proof,
a globally unique implementation revision or a configuration digest. There is no
PluginID API. The registry freezes the supplied descriptions; a plan retains its
selected configuration independently of subsequent registrations. Recursive
child demands use that same selection. See the [implemented selection contract](kernel-composition.md#registration-and-request).

Exact identity controls selection, not source equality. If the source language
makes an alias interchangeable with its target, choosing incompatible native
nominal types may violate the selected representation profile. The source/library
integration must establish eligibility; name or shape resemblance cannot do so.

For example, a consumer can choose these mappings:

```text
Builtin utf8-text       -> Go string
Builtin bool            -> Go bool
Builtin list-of         -> Go slice
Builtin closed record   -> Go struct body
Ordinary declaration    -> explicitly selected structural declaration policy
Exact Function identity -> Function library's Go callable representation

Explicit overrides:
  utf8-text             -> selected external Text type
  Function identity     -> selected invocation-handle type
```

The last two choices are type representations. The first does not add conversion
code and the second does not implement invocation. Imports refer to dependency
artifacts supplied by the consumer; generation performs no dependency download.

Useful override cases extend beyond otherwise inexpressible native types:

| Source contract | Alternative representations and reason |
| --- | --- |
| Text or an ordinary Identifier declaration | Native string, a nominal application type, or an existing dependency type to fit the consumer's API. |
| Generic list | A native collection or a library collection to select ownership, allocation or API conventions. |
| Function | Native callable syntax or a library invocation handle, depending on the intended native surface. |
| Stream | A selected iterator/channel-shaped type or a library stream handle; cancellation, buffering and termination behavior remain separate library responsibilities. |
| Cell | Structurally generated fields or an existing native Cell type; its mutation behavior is outside type generation. |
| A refined number | A native bounded number only for a compatible domain, or a type that retains the broader numeric range. |

These are design examples, not claims that these plugins currently exist. A
different representation can require a different later value adapter; selecting
the type plugin alone does not make existing native values interchangeable.

Overrides affect the artifact being generated. They do not rewrite the fields of
an already compiled external Page or change an external library's constraints.
External symbol configuration must state its compatible representation profile
and supplied dependency artifact. Known incompatibilities fail; matching native
names do not establish compatibility. The tool does not infer a behavioral or
representation contract from an external declaration's shape.

## Expert recommendation and chosen proposal

The historical [expert evaluation](../research-docs/0007-type-owned-generator-plugins.md)
endorsed type-owned mappings, supplied builtin mappings, consumer choice and
coherent generic families. Its operation-based sketch is evidence of that
consultation, not the current SPI. The subsequent owner decision to inspect
semantic choices before construction led to declarative descriptions and a
source-bound GenerationPlan. AC-09 requires that redesign to follow requirements
without incumbent API influence.
[Design correction](architecture-cases/0005-generation-interface-design.md)

Exact replacements and structural policy remain separately configured.
Type use, requested publication and owned definitions remain distinct
responsibilities, now expressed as rules, roots and recipes. The shared
constructor owns symbol scope and relocation; each native target owns grammar.

For consumers, the result is one place to choose native representations. A new
Cell mapping changes its library integration and consumer composition. Existing
structural Cell generation can instead reuse the selected Function and field
mappings without a Cell-specific plugin.

## Examples across languages

These are possible selected representations, not delivered plugins or a claim
that every row preserves every source contract. Library namespaces are
illustrative. Target profiles must define ownership, error/lifetime constraints
and retained source obligations; different languages need not choose the same
surface.

| Source type/family | TypeScript | Go | Rust | C++ |
| --- | --- | --- | --- | --- |
| Text | `string` | `string` | `String` | `std::string` |
| Generic list | `Array<T>` | `[]T` | `Vec<T>` | `std::vector<T>` |
| Function<A,B> | `(a: A) => B` | `func(A) B` | `function::Function<A, B>` | `std::function<B(A)>` |
| External Cell<T> | `cell.Cell<T>` | `cell.Cell[T]` | `cell::Cell<T>` | `cell::Cell<T>` |
| External Stream<T> | `stream.Stream<T>` | `stream.Stream[T]` | `stream::Stream<T>` | `stream::Stream<T>` |

Callable syntax and containers are existing native forms; the representation
choices are our proposal. See the [TypeScript function-type documentation](https://www.typescriptlang.org/docs/handbook/2/functions.html#function-type-expressions),
[Go specification](https://go.dev/ref/spec#Types), [Rust Vec documentation](https://doc.rust-lang.org/std/vec/struct.Vec.html),
and [C++ function-wrapper specification](https://eel.is/c++draft/func.wrap.func).
A nominal Text wrapper alone does not enforce UTF-8 in a language whose underlying
string can hold other values. That value-level guarantee depends on admitted
source values or separately supplied construction/validation behavior. Static
kind membership and native generic bounds do not inspect string contents. Neither
this plugin proposal nor the constraint proposal generates that runtime behavior.

For Function, the TypeScript library could offer a synchronous callable mapping,
an asynchronous `(a: A) => Promise<B>` mapping, and a handle mapping such as
`Invocation<A, B>`. The consumer selects the intended contract/profile; the
kernel does not infer sync/async/error behavior or generate implementation code.
Go's native function syntax likewise does not implement a remote call. Rust can
use a library handle rather than forcing every Function into a function pointer,
whose admitted callable values differ from closures with captured state.

For Stream, replacing a handle with TypeScript `AsyncIterable<T>`, Go `<-chan T>`,
a Rust iterator abstraction, or a C++ library range requires an explicit compatible
stream profile. A channel alone does not supply the canonical protocol's errors,
end/cancel rules or resource lifetime. These are possible library integrations,
not interchangeable spellings selected automatically by the generator.

For Cell, structural generation might produce fields corresponding to get and
set, each using the chosen Function mapping. Alternatively, its exact declaration
can select the library-owned external Cell<T> above. Native getters/setters,
atomic storage and synchronized mutation are behavior supplied separately.

### A builtin override that generates a helper

The Go text mapping can use a generated wrapper rather than a plain native string:

```go
// Proposed output, not implementation added by this task.
type Text string
```

In this proposed mapping, the library registers a HelperRecipe under its own
HelperKey and uses Helper(key) as the representation of text. The recipe's body
selects the native string form directly; mapping Horos text again would
reselect the override and create a cycle. The planner resolves the helper once.
The ordinary list rule still uses a mapped element slot, so list-of(text)
becomes []Text. No synthetic Horos declaration is needed.
[Exact helper contract](generation-spi.md#source-and-representation-environments)

The same distinction applies to a generic presence helper:

```go
type Presence[T any] struct {
    Present bool
    Value   T
}
```

The implementation flag uses native bool directly. A boolean field in the
canonical source still uses the configured Horos boolean mapping. The selected
plugin owns its mechanical native representation; canonical child types always
compose through declared captures or scoped Mapped demands.

### Pagination keeps its generic structure

The [worked generic-generation walkthrough](generic-elaboration.md#pagination-and-scope)
connects this composition to a single generic convention, scopes, all four target
representations and the [implementation sequence](generic-elaboration.md#implementation-sequence).

Consider this explanatory Horos source, with library imports and Person's body
omitted:

```text
type Page<T> = { items: list<T>; };
type EditablePeoplePage = Page<Cell<Person>>;
```

Generation separates the open declaration from its application:

1. Publish Page<T> once. The structural declaration policy reserves its name and
   binder; the record mapping lowers items through the selected list mapping,
   which retains the real T binder.
2. Lower Page<Cell<Person>> as an application. The selected Cell mapping composes
   Person's representation and supplies the native Cell<Person> type as Page's
   argument. It does not specialize or regenerate Page's record body.

With native array/slice mappings and an external Cell mapping, possible output
shapes are:

```typescript
// Imports and dependent declarations omitted from this shape illustration.
export interface Page<T> {
  items: T[];
}
export type EditablePeoplePage = Page<cell.Cell<Person>>;
```

```go
// Imports and dependent declarations omitted from this shape illustration.
type Page[T any] struct {
    Items []T
}
type EditablePeoplePage = Page[cell.Cell[Person]]
```

Rust and C++ can retain the same Page<cell::Cell<Person>> application with their
selected list and Cell mappings. Page has no Cell-specific case, and the
structural policy needs no Page-specific plugin. The shapes illustrate an
admitted-source-use profile; an enforcing profile must also preserve the
applicable source obligations through its selected native encoding.

An explicit override for every list could instead select Collection<T>: the
TypeScript field becomes Collection<T>, Go can use collection.Collection[T], and
Rust/C++ can use collection::Collection<T>. Child lowering propagates the selected
family through the same open Page. The [generic-coherence rule](#generics-and-representation-promises)
and [constraint-preservation contract](#constraint-preservation-in-plugin-composition)
govern that composition, including external artifact compatibility. A Data bound
remains a source-owned eligibility constraint, not a mapping selector or another
language.

## Go contracts

The current exact declarations are in the [public SPI](generation-spi.md).
Their ownership is:

| Contract | Supplied by | Invoked or consumed by |
| --- | --- | --- |
| MappingLibrary.Describe | Mapping library | Registry at registration; returns copied descriptions, not source-dependent AST code |
| Pattern selector and Capture | Source/mapping owner | Planner; captures only authenticated children of the selected use |
| BindingDescription.Output or Expand | Mapping owner | Planner; fixes an expression in the target-owned semantic vocabulary |
| DeclarationRecipe | Mapping owner | Planner for an owned declaration; returns name, form, generic convention and body |
| HelperRecipe | Owning library | Planner on demand by opaque helper identity; independent parameter scope |
| Target checks and Build | Native target integration | Planner admits native forms, then Run constructs ASTs from the closed program |
| Interpretation and reader | Source integration | Planner walks the admitted graph itself; the source classifies declarations, names descriptors and judges uses; mappings read payloads through the reader |

Implemented boundary excerpts (full fields and methods are in the linked SPI):

```go
type MappingLibrary[F, T, R any] interface {
    Describe(from, to representation.Language) (MappingDescription[F,T,R], error)
}
type GenericConvention[T any] struct {
    Parameters []ParameterDescription
    Bindings   []Expression[T]
    Arguments  []Expression[T]
}
type Target[T any] interface {
    CheckExpression(context.Context, Expression[T]) error
    CheckDeclaration(context.Context, NativeDeclaration[T]) error
    CheckProgram(context.Context, Program[T]) error
    Build(context.Context, BuildInput[T]) (Assembly, error)
}
```

### Source interpretation and semantic input

The source retains real binder identities, effective arguments and captured
assumptions. The native convention separately maps those binders to target
expressions and supplies complete native arguments. Source admission happens
before native projection, including omitted arguments. Literal data remains
opaque. No mapping invokes an implicit analyzer or resolver.

Common structural traversal is kernel-owned: the kernel walks the established
Semantics, and the source contributes classification, descriptor selectors and
judgments. [Three changes](generation-planner.md#three-spi-changes) ·
[GEN-01](generation-gaps.md#gen-01)

### Ownership and target construction

The plan contains immutable native expressions, declaration identities and
bound parameters. Source expansions finish before the target constructor runs.
Owned definitions are scoped by source identity; native helpers use opaque
owner keys without fabricated source identities or inherited binders.
External forms designate public symbols without claiming local ownership.

The constructor allocates fresh symbols and containment nodes. It records
exposures and external uses in the sealed generation.Output. Files supplies
local names and import/export syntax under the consumer's placement.
Construction failure returns no usable output. Stable target configuration
and faithful native-form interpretation remain trusted contracts, not consequences
of syntax checks. Cross-language verification stays separately supplied.

## Generics and representation promises

Open generic declarations remain open. Source binder identities and their ordered
arguments survive mapping; a plugin does not identify parameters by spelling.
For `Outer<T> = Page<T>`, Outer.T is the actual argument, not Page's formal binder.
Literal values inside `exact(...)` remain data, even when their bytes resemble a
type expression. Source adapters expose only the designated type operands.

The initial proposal selects exact builtin heads and exact declarations.
Concrete-application overrides require a separate specialization design. Exact
selection alone is insufficient: a plugin could inspect an argument and secretly
change representation even though the selected key remained the same. The profile
therefore requires one coherent representation family over its supported
substitutions: mapping after substitution must agree with substituting mapped
arguments into the mapped generic declaration. Administrative aliases and renamed
binders are permitted; unprovided conversions or changed storage/ownership
families are not. This is a trusted-plugin obligation, not a property Go interfaces
can prove. Consider
this explanatory notation:

```text
Page<T> = record { items: list-of(T) }
Bytes = Page<Byte>
```

Emitting a generic Page once with a slice does not make its field become a
ByteBuffer when Page is later applied to Byte. A local `list-of(Byte)` override
would disagree with that already emitted generic declaration. Supporting it
requires coherent specialization or a representation family, not a new dispatch
predicate. Such application selectors and automatic monomorphization are deferred.
Byte here is an illustrative library type, not a newly supported builtin.
Higher-order/computed constructors, argument-sensitive helper specialization and
anonymous type-entry publication also remain deferred. Open eligibility must be
established under binder assumptions over the profile's explicitly represented
substitutions; success on a few concrete applications does not establish it.

Representation profiles state what the native type retains and what remains a
source-level obligation. Representing Horos-admitted uses does not necessarily
enforce every kind constraint on arbitrary native callers. Unbounded integers
cannot silently become bounded integers, absence cannot collapse into a zero
value, and recursive type indirection cannot be guessed by the kernel. Unsupported
requested representations fail explicitly.

The companion [native kind-constraint proposal](native-kind-constraints.md)
describes direct bounds and source-witness/evidence encodings for enforcing
profiles. Its obligations also apply to selected plugins. The first public SPI
provides frozen layouts and whole-result associations; witness-dependent native
encodings, cross-package evidence and their language-specific syntax still require
qualified concrete implementations. Binding
selection is independent of that future encoding work; it never treats a plugin
or native representation as proof of kind membership.

## Constraint preservation in plugin composition

The [constraint proposal](native-kind-constraints.md) defines what an enforcing
encoding must retain. This plugin design defines which mapping is selected and
how it produces an owned output. A configured binding, a successful Horos judgment
and a target encoding are distinct requirements; none substitutes for the others.

A parameter's coarse sort distinguishes proper types from constructors. It is
not its kind contract. Source-owned readers must retain the actual parameter
and result contracts, membership subjects and assumptions from the supplied
semantic graph. Additional judgments require an implemented provider/profile
capability and corresponding receipts, not merely changed transformer
requirements. Neither the shared host nor a plugin invokes analysis on demand.

| Generation operation | Obligation when an enforcing profile is selected |
| --- | --- |
| Use mapping | Preserve admission of the effective source arguments, including ignored arguments; compose their selected representations and required evidence without widening the source/representation association |
| Public root request | Expose the requested open declaration under its original parameter assumptions and result contract; select one supported public parameter/header encoding rather than specializing from the first concrete use |
| Owned declaration recipe | Fill the reserved header consistently with Use and Publish; emit the required bounds and supported helpers while retaining their source assumptions and representation associations |

Witness/evidence helpers use the same ownership and deferred-definition mechanism
as other generated declarations. Reservation gives a helper a symbol; it does not
establish membership. Their source meaning remains with Horos and the source-owned
integration, and their syntax with the selected target mapping. Helper ownership
must keep different source obligations and representation configurations distinct;
matching display names cannot merge them. Unrestricted cross-package issuance
remains unqualified; explicit dependency contracts and consumer-owned composition
remain part of the broader [paired synthesis](generic-elaboration.md). The current SPI supplies declarative native helpers and automatic planning, not those native evidence encodings.

The generic-coherence requirement includes bounds, witness arguments, evidence
arguments and result representations, as well as payload storage. A future witness
profile must define how all of those substitute together. The single generic convention fixes both formal and actual substitutions. The complete native support/evidence encoding remains unimplemented.
The companion's language encodings remain target-specific candidates, with the
TypeScript and Go issuance flaws corrected after research 0009. An unsupported encoding
fails; extra header slots are not an implicit solution.

### Overrides and external representations

An override does not change the source fact `S : K`, but can change the selected
association `S -> N`. Source membership may remain usable; generated evidence or
a native bound tied to the old N must not be carried over as though the new
representation had already been justified. A changed binding registration distinguishes
configuration; it is not proof that the replacement preserves the requested
guarantee. Derive compatible output under the new configuration or reject it.

For example, replacing every list with Collection<T> can compose through an open
Page<T> only when the configured family supports the same substitutions and
constraint encoding. Selecting a different collection does not establish that
an external Page already compiled against another family is compatible. Native
artifacts must supply the required representation and constraint contract;
matching symbol names or source identities alone is insufficient.

A selected wrapper must also respect any alias-interchangeability promise of the
profile. Rust's non-enforcing type-alias bounds are a concrete reason to choose
another representation or fail; emitting a nominal wrapper does not automatically
preserve native interchangeability with its payload. Representation-only profiles
must state any conversion boundary, and this types-only design supplies no
conversions. See the [Rust constraint limits](native-kind-constraints.md#7-rust).

The current pagination representation profile still permits constraint erasure;
the current enforcement profile still returns Unsupported. Proposed plugin
composition alone implements neither an enforcing encoding nor a runtime value
checker. Selecting an enforcing promise cannot fall back to the representation
profile after a plugin, native bound or evidence requirement fails.

## Composition through the existing pipeline

1. The consumer registers a source, native target, mapping libraries and explicit
   replacements. Registration copies descriptions; no private assembly exists
   for a particular source/target pair.
2. A request supplies the source snapshot and semantic evidence, public roots,
   native language, generation profile, output identity and limits.
3. The kernel selects exact rules before structural expansion, authenticates
   captures and binds the complete generic convention.
4. Planning resolves reachable children and helpers and checks the closed native
   program. Run constructs that program without source or mapping callbacks.
5. The kernel admits the syntax and seals Output. Files organizes declarations
   from that output, checks bindings and publishes Source. Printing emits text.

[Exact registration and request signatures](kernel-composition.md) ·
[Page<Cell> reading and writing](generation-planner.md#pagination-on-read-and-write)

### File layout

The uniform stage/dimension/scope/component layout is implemented:

```text
contracts/generation/go/                         declarations and invariant planner
transform/syntax/shared/generation/go/           registry and request composition
transform/syntax/shared/target/go/               native construction mechanics
transform/syntax/<native>/model/go/              native semantic forms
transform/syntax/<native>/target/go/             form-to-AST construction
transform/syntax/horos/source/go/                source interpretation
transform/syntax/horos/mapping/go/               common source declaration policy
transform/syntax/horos-to-<native>/mapping/go/    representation descriptions
```

A real Cell library implements MappingLibrary in its own integration package.
Consumers register it; neither shared orchestration nor native targets import it.
The test fixture and the materially different reversed Pair rule use this public
boundary. [Layout and dependency rules](../LAYOUT.md) ·
[Independent composition](../integration/go/generation_planning_test.go)

## Forbidden shortcuts

- A central Function, Stream or Cell switch, name lookup, import, default catalogue
  or structural fingerprint, including a generated version of such a catalogue.
- Fixed builtin generation outside the supplied mapping mechanism.
- Implicit plugin discovery, init-time registration or dynamic loading by a
  canonical type's name.
- Order-dependent precedence, retrying defaults after override failure, or
  silently selecting a weaker representation profile.
- Raw source strings, print-and-reparse recipes, executable bodies or runtime
  bindings in type generation.
- Treating native syntax acceptance as proof of canonical equivalence or runtime
  behavior.

ASTTransformer remains the separate syntax-route contract. Engine.PlanGeneration
and Generate use the declarative generation registry and checked Output handoff.
The matrices report actual MappingLibrary implementations; they do not infer
which requests or configurations succeed.
