# Type-owned generator plugins

**Proposal, 5 November 2026. No implementation is included.** This extends the
[native representation design](type-representations.md) to builtin mappings and
explicit consumer overrides. The existing pagination transformers remain the
implemented baseline; the contracts and paths below describe their proposed
replacement. They are not current Go APIs.

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
their target backends. No new engine interface or composition-specific SPI is
needed beyond the proposed contracts below.

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

Everything in this section is proposed. All shared declarations belong to package
`generation` at `contracts/generation/go`. `F` and `T` are the existing source and
target language marker types. Existing `transform`, `representation`, `semantics`
and `ast` types retain their current meanings. Opaque handles have construction
capabilities below, so independent source and target implementations can create
them without importing private Bittype packages.

### Configuration

```go
type PluginID string
type BindingID string

// ID identifies this configured binding, including its revision/configuration.
// Plugin.ID identifies the implementation for diagnostics.
type Mapping[F, T any] struct {
    ID     BindingID
    Plugin Plugin[F, T]
}

type Binding[F, T any] struct {
    Source  Selector[F]
    Mapping Mapping[F, T]
}

type Override[F, T any] struct {
    Source   Selector[F]
    Replaces BindingID
    With     Mapping[F, T]
}

type Config[F, T any] struct {
    Source     Source[F]
    Target     Target[T]
    Bindings   []Binding[F, T]
    Overrides  []Override[F, T]
    Structural *Mapping[F, T]
    Limits     Limits
}
```

Bindings contain builtin and library mappings uniformly. Structural is a selected
source/target integration plugin, not a hidden implementation in the kernel.
The factory belongs to the proposed `generation/syntax/go` package:

```go
func NewTransformer[F, T any](
    from, to representation.Language,
    config Config[F, T],
) (transform.ASTTransformer, error)
```

It captures the effective configuration and adapts through the existing
`transform.NewTransformer` bridge. It registers nothing. Composition provides
all required source evidence through Source.Contract and the invocation's
receipts. Plugins never invoke analysis implicitly. A mapping that needs more
facts requires an explicitly enlarged source contract and corresponding receipts;
there is no new per-plugin analyzer registry in this initial proposal.

### Use, publication and definition

```go
type Plugin[F, T any] interface {
    ID() PluginID
    Use(context.Context, UseRequest[F, T]) (Type[T], error)
}

type Publisher[F, T any] interface {
    Publish(
        context.Context,
        PublishRequest[F, T],
    ) (Symbol[T], error)
}

type Definer[F, T any] interface {
    Define(
        context.Context,
        DefineRequest[F, T],
    ) (Decl[T], error)
}

type UseRequest[F, T any] struct {
    Input   Ref[F]
    Context Context[F, T]
}

type PublishRequest[F, T any] struct {
    Input   SourceDecl[F] // Always open.
    Context Context[F, T]
}

type DefineRequest[F, T any] struct {
    Spec    DefinitionSpec[F, T]
    Output  Owned[T]
    Context Context[F, T]
}
```

Use creates a proper type expression such as `[]Person` or `Page<Person>`.
Publish exposes a requested open declaration such as Page<T> and returns a local
owned symbol. An open generic declaration is not itself a proper type expression;
calling Use on it for side effects would conflate these operations.

Publisher and Definer are optional capabilities of the selected Plugin. Publish
must explicitly arrange an owned alias/wrapper when the normal representation is
inline or external, or report unsupported publication. It cannot silently omit a
requested API root. Every requested definition requires Definer. There is no
Mode enum: an inline type can need a generated helper, and an owned declaration
can reference an external type.

### Owned declarations and helpers

```go
type DefinitionSubject[F any] struct { /* private tagged alternative */ }

func SourceDefinition[F any](open Ref[F]) DefinitionSubject[F]
func HelperDefinition[F any](localKey string) DefinitionSubject[F]

func (s DefinitionSubject[F]) Source() (Ref[F], bool)
func (s DefinitionSubject[F]) Helper() (string, bool)

// Explicit correspondence; TargetIndex is zero-based.
type BinderLink[F any] struct {
    Source      Ref[F]
    TargetIndex int
}

type DefinitionSpec[F, T any] struct {
    Subject DefinitionSubject[F]
    Header  DeclarationHeader[T]
    Binders []BinderLink[F]
}

type Owned[T any] struct {
    Symbol     Symbol[T]
    Parameters []Symbol[T]
    Scope      Scope[T]
}

type Context[F, T any] interface {
    Source() SourceSession[F]
    Syntax() Builder[T]

    Lower(context.Context, Ref[F]) (Type[T], error)

    Own(
        context.Context,
        DefinitionSpec[F, T],
    ) (Owned[T], error)

    External(
        context.Context,
        ExternalRef[T],
    ) (Symbol[T], error)
}
```

Own reserves the declaration and parameter symbols immediately, records a copied
specification, and queues the selected mapping's Define operation. It returns
without recursively executing Define. Repeated reservations must agree on the
header and binder correspondence. There are no captured definition closures.

A resource is owned by `(invocation, configured mapping, source declaration or
helper key)`. SourceDefinition always names the canonical open declaration, not
the first concrete application that encountered it. A builtin text wrapper can
use HelperDefinition("text-wrapper") without inventing a source declaration.
Helpers with the same local name in different mappings remain distinct.

The host enters Define with a fresh formal environment. BinderLink explicitly
pairs real source binders with reserved target parameter slots. A helper begins
without captured source binders; HelperDefinition requires an empty Binders list,
and its native parameters are explicit in the header and Output.Parameters. Nested
Lower preserves the lexical target environment while changing resource ownership
to the newly selected mapping.

Header.Parameters can reserve native slots, but BinderLink only associates a
source binder with a target slot. It does not specify a source-witness/native-type/
evidence tuple, a representation projection from a witness, or synthesis of extra
application arguments. Reserving more slots alone therefore does not implement
the companion's witness sketches. Such an encoding needs an explicit binder and
application contract shared by Use, Publish and Define; it remains deferred.
Do not invent source binders for generated native-only evidence parameters.

### Source references, binding and projection

```go
type Selector[F any] struct { /* private, comparable identity */ }
type Ref[F any] struct { /* private invocation and source token */ }

func NewSelector[F any](
    source representation.Language,
    identitySchema string,
    canonical []byte,
) (Selector[F], error)

func (s Selector[F]) SourceLanguage() representation.Language
func (s Selector[F]) IdentitySchema() string
func (s Selector[F]) Canonical() []byte // Copy.

// Created by the host for one invocation and supplied to Source.Open.
// The source adapter keeps it private.
type RefFactory[F any] interface {
    Intern(localToken string) (Ref[F], error)
    Token(Ref[F]) (string, error)
}

type SelectorClass uint8

const (
    BuiltinSelector SelectorClass = iota + 1
    DeclarationSelector
)

type HeadKind uint8

const (
    ParameterHead HeadKind = iota + 1
    BuiltinHead
    DeclarationHead
)

type Parameter[F any] struct {
    Binder   Ref[F]
    NameHint string
    Sort     semantics.Key
}

type Argument[F any] struct {
    Parameter Ref[F]
    Value     Ref[F]
}

type View[F any] struct {
    Input   Ref[F] // Effective reference after captured substitutions.
    Subject semantics.Subject
    Origins []ast.NodeRef

    Kind HeadKind

    Selector Selector[F] // BuiltinHead or DeclarationHead.
    Binder   Ref[F]      // ParameterHead: canonical real binder.
    Open     Ref[F]      // DeclarationHead: original open declaration.

    Arguments []Argument[F] // Effective declaration argument order.
}

type SourceDecl[F any] struct {
    Open       Ref[F]
    Selector   Selector[F]
    NameHint   string
    Parameters []Parameter[F]
    Body       Ref[F]
    Origins    []ast.NodeRef
}

type Source[F any] interface {
    Language() representation.Language
    Contract() transform.Contract

    ClassifySelector(Selector[F]) (SelectorClass, error)

    Open(
        context.Context,
        transform.Request[F],
        RefFactory[F],
    ) (SourceSession[F], error)
}

type SourceSession[F any] interface {
    // Initial profile: requested named declarations, in stable order.
    // Each reference is the canonical open declaration.
    Roots() []Ref[F]

    Inspect(Ref[F]) (View[F], error)

    // Accepts an open declaration or an admitted direct application.
    // Body carries the environment of the supplied reference.
    Declaration(Ref[F]) (SourceDecl[F], error)

    // Projects an actual child while preserving the parent's effective
    // substitution, lexical, and captured-assumption environment.
    Project(
        parent Ref[F],
        directChild semantics.Subject,
    ) (Ref[F], error)

    Model() semantics.Semantics[F]
}
```

Source here means the input-language collaborator, not the Files carrier
`source.Source[F]`. View is an operational view over the existing semantic graph,
not another canonical type algebra. Its three heads distinguish parameter uses,
builtin dispatch and declaration dispatch; they do not enumerate type libraries.

A selector is source/profile-qualified and uses a source-owned identity schema.
Horos's declaration key corresponds to `(decl.Digest, declaration name)`;
parameter identity adds the ordered slot. ClassifySelector interprets that
encoding without requiring the declaration to be in the current invocation.
The generic selection code compares identities without decoding Horos values.

Inspect applies captured substitutions while retaining declaration boundaries.
Declaration(view.Open) returns an open body; Declaration(use.Input) returns a
body in that application's environment. Project admits a real direct child and
preserves the parent's lexical/substitution/assumption context. It is not a
constructor for arbitrary graph addresses. Crossing a declaration boundary uses
Declaration. Horos-specific readers decode designated descriptor operands and
call Project; literal data in exact descriptors is never traversed as types.

Two environments remain separate: source binder -> effective source argument
belongs to SourceSession; source binder -> target parameter type belongs to the
lowering context. Matching parameter names or slice positions is insufficient.
RefFactory interns the complete source context and rejects references from other
invocations. Source adapters keep their factory private.

Roots means requested named API declarations, not every resolved dependency.
Source.Open admits source evidence and provides views; it does not eagerly reject
all recursive or unsupported bodies before representation selection. Required
negative or absent source admission is still a failure.

### Native fragments and symbols

```go
type Scope[T any] struct { /* private session identity */ }
type Symbol[T any] struct { /* private session identity */ }

type Fragment[T any] struct { /* immutable native syntax template */ }
type Type[T any] struct { /* fragment classified as a type use */ }
type Decl[T any] struct { /* fragment classified as a declaration */ }

func (t Type[T]) Fragment() Fragment[T]
func (d Decl[T]) Fragment() Fragment[T]
```

```go
type NameSlot[T any] struct {
    Attribute string
    Symbol    Symbol[T]
}

type Child[T any] struct {
    Role string
    Node Fragment[T]
}

type Link[T any] struct {
    Role   string
    Target Symbol[T] // Local resolved-link target only.
}

type NodeSpec[T any] struct {
    Kind       ast.Kind[T]
    Attributes []ast.Attribute
    Names      []NameSlot[T]
    Children   []Child[T]
    References []Link[T]

    Defines Symbol[T] // Zero, or the symbol anchored by this node.
    Origin  ast.Origin
}

type UseForm[T any] string // Target-owned construction vocabulary.

type SymbolUse[T any] struct {
    Form   UseForm[T]
    Symbol Symbol[T]
    Origin ast.Origin
}

type FragmentView[T any] struct {
    Node *NodeSpec[T]
    Use  *SymbolUse[T]
}

func (f Fragment[T]) View() FragmentView[T] // Copies; exactly one alternative.

type Builder[T any] interface {
    Node(NodeSpec[T]) (Fragment[T], error)
    Use(SymbolUse[T]) (Fragment[T], error)
    Type(Fragment[T]) (Type[T], error)
    Declaration(Fragment[T]) (Decl[T], error)
}
```

These fragments are immutable target syntax templates. Builder copies input
slices and enforces session/structural ownership. Type and Declaration mark the
intended native syntactic category; they do not add a language-independent
catalogue of valid native kinds. Native category/grammar admission remains with
the target adapter and existing checker.

A target package supplies typed convenience functions over this surface:

```go
gotarget.String(builder, origin)
gotarget.Slice(builder, element, origin)
gotarget.Named(builder, symbol, arguments, origin)
gotarget.DefinedType(builder, owned, body, origin)
```

SymbolUse is a target-owned construction placeholder, expanded into structured
native syntax during assembly. It is not a canonical type kind or runtime
binding. The generic kernel never appends generic arguments to arbitrary output;
the selected plugin and target helper construct the complete native type use.

Names as well as reference links stay symbolic until allocation. Every materialized
occurrence gets fresh containment nodes; reusable fragments do not share contained
AST nodes. Definitions and parameter symbols have one primary anchor. Free symbols
must remain in their native scope. For example, a Page<T> field fragment cannot be
moved into an unrelated declaration while retaining Page's T.

### Target session and construction capabilities

```go
type UnitKey string

type Placement struct {
    Unit  UnitKey
    Scope string // Target-configured declaration scope within/for the unit.
}

type ParameterHeader[T any] struct {
    Form     ast.Kind[T]
    NameHint string
}

type DeclarationHeader[T any] struct {
    Placement Placement
    Form      ast.Kind[T]
    NameHint  string
    Parameters []ParameterHeader[T]
}

type ExternalRef[T any] struct {
    Artifact string // Key into the caller-supplied locked artifact catalogue.
    Export   string // Target-owned export/member identity.
}

type ResourceKey struct { /* host-created ownership identity */ }

type Target[T any] interface {
    Language() representation.Language

    Open(
        context.Context,
        Storage[T],
    ) (TargetSession[T], error)
}

type TargetSession[T any] interface {
    EntryScope() Scope[T]

    Reserve(ResourceKey, DeclarationHeader[T]) (Owned[T], error)
    External(ExternalRef[T]) (Symbol[T], error)

    Builder(Scope[T]) Builder[T]
    Parameter(Scope[T], Symbol[T]) (Type[T], error)

    Define(Owned[T], Decl[T]) error

    Finish(
        context.Context,
        []Symbol[T], // Published roots.
    ) (ast.Set[T], error)
}
```

```go
type Limits struct {
    MaxNodes       int
    MaxDepth       int
    MaxUnits       int
    MaxWork        int
    MaxDefinitions int
    MaxImports     int
    MaxBytes       int
}

type Storage[T any] interface {
    Scope(parent Scope[T], stableKey string) (Scope[T], error)
    Symbol(scope Scope[T], stableKey string) (Symbol[T], error)
    Builder(scope Scope[T]) Builder[T]

    Limits() Limits
    ChargeOutput(nodes, bytes int) error
}
```

ResourceKey is comparable and supplied by the host, scoped to the active invocation.
The target also needs a stable ordering key for deterministic allocation:

```go
func (k ResourceKey) OrderKey() string
```

OrderKey is opaque to the target and excludes the invocation nonce. It distinguishes
the configured mapping and definition subject; it is not a structural fingerprint
or semantic equivalence claim. Source identity decoding stays in the source adapter.

Only the generation host receives TargetSession and Storage. Plugins receive the
restricted Context surface. The contract package supplies copying constructors
and handle issuers; Storage scopes and symbols are session-owned. Native adapters
own declaration/import scopes, spelling, placement, and artifact/export lookup.
They can capture those configurations in their constructors. No canonical mapping
is hidden in the target backend.

Placement identifies logical AST units and configured native declaration scopes.
The current ASTSet return boundary does not assign filesystem paths or publish a
Files Source. Parameter constraints are built in Define after reservation and
must agree with the native header. Arbitrary fragment-local type/lifetime binders
and arbitrary nested declaration scopes are deferred in the first profile.

ExternalRef identifies a supplied artifact/export, not an AST node. Without its
declaration tree, emit native import/include and qualified syntax rather than a
fabricated NodeRef. Import aliases are allocated per materialization unit under
native scope rules. The consumer retains the dependency environment; ASTSet is
neither a dependency lock nor a semantic receipt.

Limits cover actual materialized occurrences, expanded uses, definitions, imports,
units, payloads and work. They retain existing ceilings rather than counting only
reusable templates. Cancellation is cooperative; arbitrary Go plugins remain
trusted code. No cross-invocation cache or global mutable plugin state is proposed.

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
profiles. Its obligations also apply to selected plugins. This initial plugin
profile does not claim to settle witness-dependent parameter expansion,
cross-package evidence or all the native syntax such encodings need. Binding
selection is independent of that future encoding work; it never treats a plugin
or native representation as proof of kind membership.

## Constraint preservation in plugin composition

The [constraint proposal](native-kind-constraints.md) defines what an enforcing
encoding must retain. This plugin design defines which mapping is selected and
how it produces an owned output. A configured binding, a successful Horos judgment
and a target encoding are distinct requirements; none substitutes for the others.

Parameter.Sort distinguishes proper types from constructors. It is not the
parameter's kind contract. Source-owned readers must retain the actual parameter
and result contracts, membership subjects and assumptions from Model(). Any
additional required facts belong in Source.Contract and the explicitly supplied
receipts. Neither the shared host nor a selected plugin invokes analysis on demand.

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
matching display names cannot merge them. Cross-package evidence remains unresolved.

The generic-coherence requirement includes bounds, witness arguments, evidence
arguments and result representations, as well as payload storage. A future witness
profile must define how all of those substitute together. The simple BinderLink
surface above does not yet define that expansion or a dependent projection. The
Go three-parameter and TypeScript/Rust/C++ witness sketches in the companion are
therefore future encodings, not capabilities implied by these Go interfaces.
An unsupported encoding fails; extra header slots are not an implicit solution.

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

The consumer constructs one transformer for each desired source/target pair:

1. Supply the source and target language/profile descriptors, source adapter,
   target backend, builtin bundle, selected library plugins and overrides.
2. Supply the already captured source collection and explicit semantic receipts.
   The existing typed transformer bridge admits them. Generation invokes no
   resolver or analyzer implicitly.
3. Select each requested root's mapping and invoke Publish before expanding its
   body. Own reserves symbols/parameters and queues immutable definition specs.
   It does not run Define recursively; publication returns an owned symbol.
4. Compose child mappings under the same configuration. A list mapping requests
   its element's representation; a structural Page mapping composes its field
   representations. An exact library mapping is selected before expanding that
   library declaration's body.
5. Drain the definition queue, running each selected Define once in its formal
   context. Definitions can reserve further helpers. After the queue closes,
   assemble target AST units with distinct containment nodes for each occurrence
   and symbolic references to owned declarations or explicit native imports.
   The result is the existing `ast.Set[T]` through the ASTTransformer bridge.
6. Use the existing target Syntax formatter and shared printer. Source-file
   placement and actual IO remain consumer responsibilities; generation does not
   fabricate the native semantic evidence required by a Files formatter.

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

These are proposed source paths only. This documentation task creates no Go
package, implementation directory, generator, provider or runtime binding.

```text
bittype/
  contracts/generation/go/
    selector.go                 exact identity and copying constructors
    source.go                   source references, views and sessions
    plugin.go                   bindings, overrides, Use and Publish
    ownership.go                definition subjects, headers and binder links
    fragment.go                 immutable native templates and symbolic slots
    target.go                   target sessions, storage and limits
  generation/syntax/go/
    bindings.go                 exact selection and guarded replacement
    lower.go                    scoped lowering and terminal errors
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

A consumer's configuration can look like this. Horos and Go are shorthand for
existing language marker types; the named plugin instances, selectors, source
adapter and target backend are supplied by their owning packages. This is a
composition sketch using proposed APIs, not a runnable current API example:

```go
config := generation.Config[Horos, Go]{
    Source:     sourceAdapter,
    Target:     goTarget,
    Bindings:   slices.Concat(horosDefaults, libraryBindings),
    Structural: &structuralMapping,
    Limits:     limits,
    Overrides: []generation.Override[Horos, Go]{
        {
            Source:   textSelector,
            Replaces: "horos.text.go.string/v1",
            With: generation.Mapping[Horos, Go]{
                ID:     "application.text.go.wrapper/v1",
                Plugin: textWrapper,
            },
        },
    },
}
transformer, err := generate.NewTransformer[Horos, Go](from, to, config)
```

After handling construction errors, the consumer registers this transformer with
the existing engine. The same composition works for the other three targets with
their respective plugins/backend. A missing or changed expected text binding
makes construction fail, rather than quietly changing the override's meaning.

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
