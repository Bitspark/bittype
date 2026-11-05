# Type-owned generator plugins

**Proposal, 5 November 2026. No implementation is included.** This extends the
[native representation design](type-representations.md) to builtin mappings and
explicit consumer overrides. The existing pagination transformers remain the
implemented baseline; the contracts and paths below describe their proposed
replacement. They are not current Go APIs.

The [paired consultation synthesis](generic-elaboration.md) refines this design.
It replaces the earlier binder-to-slot concept with scoped representation
products and paired formal/application conventions. The old full Go SPI sketch
is withdrawn as a candidate interface; it survives only in frozen research
inputs. The complete replacement SPI awaits a selected support domain and the
construction/lifetime decisions listed below. No alternative compatibility API
is proposed.

The later [Logos preservation proposal](logos-preservation.md) explains how
mapping-owned representation laws could compose into checked preservation
claims. The support carrier and a plugin's declared capability alone do not
establish those claims; the proof must concern the actual selected output.

The design is scoped to native types and interfaces. It generates no calls,
stream operations, runtime adapters, validators or transport code. The expert
consultation and its evaluation are recorded in
[research 0007](../research-docs/0007-type-owned-generator-plugins.md).

Read this together with [native kind-constraint preservation](native-kind-constraints.md).
This document owns mapping selection, generation operations and resource ownership;
the companion owns the guarantees and limits of native constraint encodings.
The [composition contract below](#constraint-preservation-in-plugin-composition)
connects them without claiming that either proposal is implemented.

## Responsibility and ownership

A native representation has two axes: the source type and the target language.
For example, the Function library can supply a mapping to TypeScript's callable
syntax; the TypeScript syntax package owns that syntax independently of Function.
Another selected mapping can represent the same declaration as an external
handle. Neither choice changes the canonical Function declaration.

| Component | Responsibility |
| --- | --- |
| Source adapter | Interpret admitted source semantics, preserve binders, and supply exact mapping selectors. |
| Type-owned plugin | Choose and construct the native representation of its builtin or named declaration. Compose selected mappings of child types. |
| Target-language tooling | Supply native AST construction, symbols, imports and output-unit assembly. It has no canonical type catalogue. |
| Shared generation composition | Select configured mappings and coordinate invocation and ownership. It interprets neither source descriptors nor native types. |
| Consumer | Supply dependencies and select defaults, library plugins, explicit overrides and representation profiles. |
| Existing engine | Register the resulting ASTTransformer as one directed language edge. It never registers individual type plugins. |

Builtins are ordinary supplied mappings. Horos-owned integration code provides
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

The proposed generator is one language-agnostic composition kernel. Every native
representation rule is supplied through plugins, including builtin mappings and
the structural declaration policy. A language-specific generator becomes a
configured composition, with this division of work:

```text
agnostic generation kernel
  + source semantic adapter
  + selected builtin and library mappings
  + target syntax backend
  = one configured ASTTransformer
```

An individual mapping is selected for a canonical constructor or declaration and
a target language/profile. For example, the Horos list mapping for TypeScript
constructs an array type; its Go counterpart constructs a slice type. Each asks
Context.Lower for its element under the same effective configuration. The shared
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
through Lower. An external Cell representation can use its argument without
expanding Cell's implementation shape; a structural Cell representation lowers
its fields. The kernel does not eagerly flatten every declaration into primitives
or require a complete child traversal before dispatch. Source adapters preserve
the actual type-operand and binding boundaries.

Here "kind-specific generator" means a mapping for a source constructor or named
type. Horos constraint kinds such as Data still govern applicability; their
overlapping memberships do not select a unique plugin. The same kernel is used
for TypeScript, Go, Rust and C++; their native grammar and scope rules remain in
their target backends. The existing engine interface remains; the generic elaboration seam below
supplies the shared composition convention.

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

The consumer assembles exact bindings from supplied defaults and selected
libraries, then supplies a separate override set:

1. Reject duplicate selectors in the assembled bindings and in the override set.
2. Each override must name an existing exact binding and its expected BindingID.
   A missing binding or mismatched ID is an error. Overrides do not form an
   order-dependent replacement chain.
3. An exact declaration binding takes precedence over the separately configured
   structural declaration policy. Adding such a binding is how a consumer selects
   a custom representation for an otherwise structural declaration; it does not
   replace a virtual entry in the exact table.
4. At use time, select the exact binding if present, otherwise the structural
   policy for an eligible ordinary declaration, otherwise fail.
5. A selected plugin's failure is final. Unsupported input or missing evidence
   never causes a retry with a replaced mapping or the structural policy.

BindingID identifies a particular configured mapping, including its revision and
options. PluginID identifies its implementation for diagnostics; it is too weak a
replacement guard. Composition supplies IDs explicitly, without fingerprinting
Go closures. Changing a mapping's configuration requires a new BindingID.
Mapping IDs are unique in the effective configuration; one structural
mapping can serve many declarations, whose identities keep their resources apart.
Configuration is captured per transformer and invocation state is separate.
Recursive child mapping uses the same effective configuration.

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

The expert endorses type-owned mappings, supplied builtin mappings, consumer
choice and the unchanged ASTTransformer boundary. The adopted refinements are:

- Exact bindings plus guarded overrides; structural declaration handling stays
  a separate policy. An override names the configured BindingID it replaces.
- Separate type use, API publication and deferred definition. Inline, external
  and owned effects can coexist, so no Mode enum is proposed.
- Reserve symbols immediately and queue definitions. A mandatory public Plan
  method is unnecessary for the cases in this proposal.
- Require a coherent generic representation family. Exact selectors alone do
  not prevent an argument-sensitive plugin from breaking generic substitution.

The documented choice accepts those recommendations. It adds an explicit stable
ordering accessor for opaque resource keys, keeps the existing Syntax path
convention, and makes native syntax classification target-owned. The shared
fragment layer supplies ownership and relocation, not a second native grammar.
The proposal remains experimental until implementation is separately authorized.

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

Its Use operation reserves HelperDefinition("text-wrapper") and returns a named
type fragment. Own queues Define once. Define emits the backing native `string`
directly; calling Lower on Horos text again would reselect the same override and
cause a cycle. The ordinary list plugin still lowers its element through the
selected mappings, so list-of(text) becomes `[]Text`.

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
compose through Lower.

### Pagination keeps its generic structure

The [worked generic-generation walkthrough](generic-elaboration.md#pagination-and-scope)
connects this composition to Introduce/Apply, scoped products, all four target
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

The earlier complete interface sketch did not express witness-dependent payload
projections, extra actual arguments or result associations. It is replaced by the
[paired convention/product proposal](generic-elaboration.md#one-generic-convention-for-formals-and-applications),
not retained as an alternative SPI. The frozen inputs for research 0007/0008
retain what those experts reviewed. No package implements either sketch today.

The following are conceptual operations and carriers, not final Go signatures:

| Surface | Required meaning |
| --- | --- |
| Configured mapping | Exact selector, configured BindingID, implementation identity and explicit consumer options |
| Source adapter | Admitted source-use coverage, canonical open declarations, effective applications and context-preserving type-child projection |
| Representation product | Native payload plus integration-owned immutable association support, scoped to the source use and target formal environment |
| Generic convention | One immutable open-declaration convention with paired Introduce and Apply operations |
| Introduce | Bind real source formals to products and construct native formals/constraint clauses from reserved symbols |
| Apply | Construct the complete native use and its result association under the caller's actual arguments/assumptions |
| Use | Produce the selected representation product; a generic use follows its convention |
| Publish | Expose the requested open API declaration, using the same convention |
| Own | Reserve local declaration/native parameter symbols and queue an immutable definition request |
| Define | Fill the reservation in its introduced formal environment; preserve the authoritative native formals and constraints |
| Represent / payload projection | One admitted association path; requesting only the payload does not bypass admission |
| External reference | Native artifact/export identity accompanied by its supplied generic/representation contract |

A third Go type parameter for the integration's support carrier is a candidate
mechanic, not a required native source/payload/proof tuple. A carrier must be an
integration-issued immutable token or obey equivalent safe access rules; generic
Go typing or copying a struct cannot guarantee that. The kernel cannot bless an
arbitrary payload/support pair by attaching a requested source reference.

Generic conventions, publication and definition remain separately demanded
capabilities. Inline, external and owned effects can coexist. Apply must have an
explicit route to its callee without forcing every application to Own a local
definition. A requested root cannot silently disappear. Nongeneric definitions
need explicit zero-formal or nongeneric handling in the eventual API.

### Source interpretation and semantic input

Keep two environments distinct: source binder to effective source argument, and
source binder to scoped native product. Use real binder identities, not parameter
names or positions alone. Source references retain lexical, substitution and
captured-assumption context. Projection only enters actual designated type
operands; literal data in exact descriptions remains opaque. Crossing a named
declaration uses its explicit boundary, rather than flattening it before selection.

The source adapter admits coverage already established by the supplied semantic
receipts. Coarse proper/constructor sorts are not kind contracts. Additional
judgments require an implemented provider/profile capability and corresponding
receipts; configuration changes do not manufacture new facts. Generation invokes
no hidden analyzer or resolver. Association products do not enter source Semantics.

### Ownership and target construction

A source definition belongs to its configured mapping and canonical open
declaration for one invocation. Helper keys are mapping-local. The first concrete
application cannot choose a generic signature or capture its caller's symbols.
Repeated reservations must agree on one authoritative native layout. Native
helpers have explicit parameters and no inherited source-binder environment.

Reusable fragments remain immutable native syntax templates with symbolic names
and references. Materialized uses receive fresh containment nodes; owned symbols
and parameters have their own anchors and valid scopes. The target owns grammar,
name/import allocation and assembly, not canonical mapping decisions. An external
reference is not a fabricated resolved AST NodeRef.

Public issuers/accessors for source refs, admitted uses, products, signatures,
environments and symbols must let independent integrations operate without
private Bittype imports. The complete construction API, host-mediated Apply,
callee acquisition and typed support-demand/result operations remain to be
specified together. The expert's partial Go snippets are not a ready-to-compile
contract. See the [remaining API decisions](generic-elaboration.md#go-spi-status-and-next-scope).

The generation host retains per-invocation state, bounded work/depth/imports/
definitions/expanded output, and stable allocation ordering. Product caching must
include the target formal environment; cross-invocation caching remains deferred.
Configured plugins are trusted Go code, not a sandbox. Error or cancellation
discards the entire output. The existing ASTTransformer bridge returns an ast.Set;
it does not publish file placement, a native semantic receipt or an export manifest.

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
profiles. Its obligations also apply to selected plugins. The proposed plugin SPI
remains incomplete: witness-dependent parameter expansion, cross-package evidence
and the native syntax such encodings need still require concrete contracts. Binding
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
| Use | Preserve admission of the effective source arguments, including ignored arguments; compose their selected representations and required evidence without widening the source/representation association |
| Publish | Expose the requested open declaration under its original parameter assumptions and result contract; select one supported public parameter/header encoding rather than specializing from the first concrete use |
| Define | Fill the reserved header consistently with Use and Publish; emit the required bounds and supported helpers while retaining their source assumptions and representation associations |

Witness/evidence helpers use the same ownership and deferred-definition mechanism
as other generated declarations. Reservation gives a helper a symbol; it does not
establish membership. Their source meaning remains with Horos and the source-owned
integration, and their syntax with the selected target mapping. Helper ownership
must keep different source obligations and representation configurations distinct;
matching display names cannot merge them. Unrestricted cross-package issuance
remains unqualified; explicit dependency contracts and consumer-owned composition
are proposed in the [paired synthesis](generic-elaboration.md).

The generic-coherence requirement includes bounds, witness arguments, evidence
arguments and result representations, as well as payload storage. A future witness
profile must define how all of those substitute together. The paired
convention/product
direction replaces BinderLink; the complete support API remains undecided.
The companion's language encodings remain target-specific candidates, with the
TypeScript and Go issuance flaws corrected after research 0009. An unsupported encoding
fails; extra header slots are not an implicit solution.

### Overrides and external representations

An override does not change the source fact `S : K`, but can change the selected
association `S -> N`. Source membership may remain usable; generated evidence or
a native bound tied to the old N must not be carried over as though the new
representation had already been justified. A changed BindingID distinguishes
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

The proposed consumer constructs one transformer for each desired source/target pair:

1. Supply the source and target language/profile descriptors, source adapter,
   target backend, builtin bundle, selected library plugins and overrides.
2. Supply the already captured source collection and explicit semantic receipts.
   The existing typed transformer bridge admits them. Generation invokes no
   resolver or analyzer implicitly.
3. Select each requested root's mapping and invoke Publish before expanding its
   body. Own reserves symbols/parameters from the shared convention and queues
   immutable definition requests.
   It does not run Define recursively; publication returns an owned symbol.
4. Compose child mappings under the same configuration. A list mapping requests
   its element's representation; a structural Page mapping composes its field
   representations. An exact library mapping is selected before expanding that
   library declaration's body.
5. Drain the definition queue, introducing formal products and preserved native
   constraints, then running each selected Define once in that context.
   Definitions can reserve further helpers. After the queue closes,
   assemble target AST units with distinct containment nodes for each occurrence
   and symbolic references to owned declarations or explicit native imports.
   The result is the existing `ast.Set[T]` through the ASTTransformer bridge.
6. Use the existing target Syntax formatter and shared printer. Source-file
   placement and actual IO remain consumer responsibilities; generation does not
   fabricate the native semantic-provider receipt required by a Files formatter.

There is one engine edge per language pair. Alternative mapping configurations
use separate configured engine instances under the current registry policy.
Adding a plugin neither adds an engine edge nor registers another language.
There is no mandatory public Plan method: the reservation queue provides the
needed planning discipline. Any error discards the entire invocation output.
Inline expansion cycles fail. The initial profile may keep generated recursion
unsupported, but must inspect demanded representation dependencies after selection;
an external mapping may avoid expanding a recursive source body. Reserved symbols
do not establish native recursive-type legality.

For pagination, selecting a custom text representation affects the generated
Next field and other text occurrences consistently. Selecting a generic list
representation affects Items while preserving its T binder. Selecting a Cell
representation composes at the Cell application in `Page<Cell<Person>>`; no
Page-specific Cell handling is needed.

## Proposed file layout and dependencies

These are proposed source paths only; the final file split depends on completion
of the public SPI. This documentation task creates no Go
package, implementation directory, generator, provider or runtime binding.

```text
bittype/
  contracts/generation/go/
    selector.go                 exact identity and copying constructors
    source.go                   source references, views and sessions
    plugin.go                   bindings, overrides, Use and Publish
    ownership.go                definition subjects and authoritative headers
    product.go                  scoped representation products
    convention.go               paired formal/application conventions
    fragment.go                 immutable native templates and symbolic slots
    target.go                   target sessions, storage and limits
  generation/syntax/go/
    bindings.go                 exact selection and guarded replacement
    lower.go                    scoped representation and terminal errors
    conventions.go              signature lifecycle and formal environments
    definitions.go              reservation and deferred definition queue
    transform.go                existing ASTTransformer bridge
  generation/syntax/source/horos/go/
    source.go                   source view and captured environments
    identity.go                 builtin/declaration selector construction
    read.go                     Horos operands, fields and evidence queries
  generation/syntax/target/typescript/go/
  generation/syntax/target/golang/go/
  generation/syntax/target/rust/go/
  generation/syntax/target/cpp/go/
    builder.go                  native construction helpers
    symbols.go                  native scopes, imports and reservations
    assemble.go                 symbolic expansion and AST-set assembly
  transforms/horos-to-typescript/go/transform.go
  transforms/horos-to-go/go/transform.go
  transforms/horos-to-rust/go/transform.go
  transforms/horos-to-cpp/go/transform.go
                                thin configuration-only factories
```

Every target directory has the same three-file responsibility split. Syntax is
explicit as in the existing stage hierarchy. Source/target are collaborators of
Syntax generation, not new dimensions. Go uses `golang` in native implementation
paths and the existing `horos-to-go` transformer edge path remains unchanged.
Proposed Go package names are `generation` for the shared contract/composer,
`source` for source adapters and `target` for native adapters; consumers use local
import aliases when needed. Neither adapter is an engine-registered provider.

Outside Bittype, a type owner can add an independent Go integration module within
its repository. The location is illustrative; no new repository is required:

```text
type-library/
  tooling/bittype/go.mod
  tooling/bittype/generation/syntax/typescript/go/
  tooling/bittype/generation/syntax/golang/go/
  tooling/bittype/generation/syntax/rust/go/
  tooling/bittype/generation/syntax/cpp/go/
    mapping.go                  owned builtin/declaration mapping
    bindings.go                 supplied exact bindings

horos-owned integration module/
  generation/syntax/<language>/go/
    scalar.go list.go record.go structural.go bindings.go
  transforms/horos-to-<target>/go/
    defaults.go                 optional opinionated composition

consumer/
  generation/
    configuration.go            select bindings, overrides and artifacts
```

Only supported target packages are created during a future implementation;
missing support is explicit, with no empty placeholders. The Horos integration
owns builtin mappings and structural declaration policies. Horos core does not
import Bittype. Bittype factories do not import the integration module, so there
is no Bittype -> integration -> Bittype module dependency cycle.

Dependencies follow the ownership boundaries:

```text
consumer -> generation composer + source adapter + target adapter
consumer -> selected Horos/library integration modules
integration mapping -> generation contracts + source readers + target helpers
source adapter -> Horos model/front-end facts + shared contracts
target adapter -> own native syntax + shared contracts/carriers
shared generation -> contracts/carriers only
engine -> existing shared transformer contract only
```

Consumer construction supplies the source adapter, target backend, selected
builtin/library mappings, structural policy, explicit guarded overrides,
artifact contracts and the compatible association/convention protocol. For
example, replacing the configured text-to-string binding must name its expected
BindingID and supply the new wrapper mapping and compatible association rules.
It must fail on a missing/revised expected binding or an incompatible imported
fixed dependency. The exact Go configuration signature remains part of SPI
finalization; the earlier two-marker Config example is not an alternative API.

After construction, the consumer registers one configured ASTTransformer with
the existing engine. Other targets use the same ownership and operation roles
with their own integrations and native backends.

When implementation is authorized, replace the private fixed Dialect and its
builtin native mapping switch with this composition. Retire the superseded route
without compatibility aliases, duplicated builtin modes or hidden defaults. The
current task changes documentation only; the existing implementation is intact.

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

The existing ASTTransformer remains the public engine operation. This proposal
adds no engine capability, runtime registry or provider-matrix cell. The matrices
continue to describe implemented source, not proposed packages.
