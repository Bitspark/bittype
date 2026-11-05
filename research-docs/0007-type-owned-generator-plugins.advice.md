# Advice for consult
**Research doc:** doc_4c3de827-b0ed-45cc-8870-6c6b892772f0.md
**Document:** doc_4c3de827-b0ed-45cc-8870-6c6b892772f0
**Advisor:** ChatGPT, model "6 Pro" as labelled on the answer page
**Session:** nsess_2961ad21-951d-4774-a3fc-9de617dcd960
**Run:** run_35a886b0-509f-4467-8535-9267d74ba42f
**Chat:** https://chatgpt.com/c/6aec7e89-e490-83eb-9bca-3b8e9fde5ffd
**Verify:** nightfall consult verify run_35a886b0-509f-4467-8535-9267d74ba42f --output <this file>

# Response to Research 0007

**I would adopt the proposal’s ownership and composition model, with four changes: simplify replacement semantics, remove `Mode`, distinguish declaration publication from type use, and make generic coherence an explicit plugin obligation.**

The existing `ASTTransformer` boundary is sufficient. The new mechanism can remain a library that constructs one configured transformer from a source adapter, a target adapter, and explicitly supplied mappings. Function, Stream, Cell, and builtin representations can all participate without adding their names to the generation kernel.

My review uses the supplied code excerpts and stated limitations in the [attached proposal](https://chatgpt.com/library?file_id=libfile_c907151601fc8191ae4c4dff5953d831&account_id=personal) as its baseline. The signatures below are proposed contracts.

The main decisions I recommend are:

| Area | Recommendation |
|---|---|
| Selection | Exact bindings, checked replacements, and a separately configured structural declaration policy |
| Plugin operations | `Use`, with optional `Publish` and `Define` operations |
| Ownership | Reserve symbols immediately; queue definitions; assemble concrete ASTs after the queue closes |
| Generics | Require one coherent representation family over the profile’s supported substitutions |
| Fragments | Immutable native syntax templates with symbolic names and references |
| Initial scope | Named declaration roots, proper-type parameters, direct applications, supplied external types, and explicit helpers |

## 1. Exact bindings and checked replacements are sufficient

The selector design is fundamentally right. A selector should identify an exact source dispatch head within a versioned source identity domain:

- A builtin scalar discriminator, such as Horos `utf8-text`.
- A builtin constructor, such as `list-of`.
- An ordinary declaration’s defining module digest and declaration name.

The generic package should compare these identities without decoding them. Horos owns the distinction between the scalar operator and its discriminator, the encoding of module pins, and the treatment of declaration aliases.

Kind membership belongs in eligibility checking. It must not participate in choosing the most specific mapping.

### Separate structural policy from the exact binding table

I would simplify the proposed “effective default” rule:

1. Assemble exact bindings supplied by the consumer and its dependencies.
2. Reject duplicate selectors.
3. Apply explicit replacements to existing exact bindings.
4. At invocation, select an exact binding if one exists.
5. Otherwise, select the configured structural policy for an eligible ordinary declaration.
6. Otherwise, return an error.

Adding an exact declaration binding therefore takes precedence over structural handling. It does not need to replace a virtual entry in the exact table.

This makes a missing replacement uniformly erroneous, whether its selector names a builtin or a declaration. It also lets most configuration errors be detected before a source snapshot exists.

A suitable configuration shape is:

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

I recommend making `Replaces` required. A plugin implementation ID is too weak a guard: two instances of the same implementation can have different options or external dependencies. A binding revision is a declaration by the composition owner of which configured mapping is being replaced. Arbitrary Go closures cannot be reliably fingerprinted to derive this automatically.

The **selection policy** is immutable per transformer. Invocation-specific references, eligibility judgments, reservations, and definitions are created after admitting the captured inputs.

A selected mapping’s failure is terminal. That includes unsupported input, missing evidence, invalid syntax, resource conflicts, and exceeded limits. None should trigger another candidate or the structural policy.

### Preserve declaration boundaries

Resolve a use to its declaration identity before considering body expansion. A genuine source alias can have its own mapping; an unrelated declaration with an identical body cannot inherit it.

One qualification should be added: **identity preservation does not by itself establish preservation of source equality or assignability.** If Horos considers an alias interchangeable with its target, mapping them to incompatible native nominal types may violate the chosen representation profile. The supplied excerpts do not establish Horos’s complete alias semantics, so this should be a source-owned eligibility obligation rather than an assumed consequence of exact selection.

## 2. Complete the source boundary around references and environments

The proposed `SourceSession` lacks several operations necessary to implement independent plugins:

- Distinguishing parameter uses from dispatchable heads.
- Recovering an application’s original open declaration.
- Obtaining effective argument bindings.
- Entering a declaration body with the correct substitutions.
- Projecting a descriptor operand without losing its environment.
- Identifying requested output roots separately from all dependency declarations.

`Model()` plus `Subject()` does not solve these problems. A plugin that reads an operand’s `TermID` still needs a legitimate way to obtain its environment-bearing `Ref`.

I would use the following source contract.

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

These methods need explicit semantics:

**`ClassifySelector` validates the identity domain and canonical encoding.** It does not require the named declaration to occur in this invocation.

**`Inspect` applies captured parameter substitutions while preserving declaration boundaries.** It must not unfold an alias before its mapping has been selected.

**`Declaration(view.Open)` returns the open body.** That is the input for emitting one generic definition. `Declaration(use.Input)` returns the body under that application’s effective substitutions, which is useful for intentional inline expansion.

**`Project` validates the source relationship.** It is not an arbitrary graph-address constructor. Crossing a declaration boundary uses `Declaration`.

**`RefFactory` provides an actual extension mechanism.** The source adapter interns a token representing the complete effective source context. `Token` rejects references from other invocations. Private fields alone would otherwise make independently implemented source adapters impossible.

Horos-specific readers should sit above this interface. For example, a list reader can decode the Horos operand locator and use `Project` to return its element. A record reader can return field names, optionality, and projected field references. No fabricated formal parameter is needed for a builtin descriptor operand.

Kind queries should likewise remain Horos-owned queries over captured evidence. Missing facts remain unknown.

Finally, `Open` must not eagerly reject every unsupported body or recursive dependency. It admits and interprets the captured source evidence; representation-specific traversal begins after selection. A custom mapping may never need a declaration’s body.

## 3. Remove `Mode`; retain three small operations with distinct purposes

`Inline`, `External`, and `Owned` describe overlapping effects.

An inline callable can require an owned helper. An owned wrapper can use an external type. A list representation can need an import. Making these mutually exclusive modes creates redundant consistency rules without simplifying ownership.

I would retain a small use interface and add optional interfaces for publication and deferred definitions:

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

### Publication is a real missing operation

An open `Page<T>` is a constructor to publish, not a proper type use to lower. Calling `Use(Page)` merely for side effects either invents free target binders or returns an unusable constructor expression. Go, for example, requires generic types to be instantiated when used as types. [The Go Programming Language](https://go.dev/ref/spec?utm_source=chatgpt.com)

For the first profile, `Publish` should return a **local owned symbol** representing the requested native declaration.

If a selected mapping normally returns inline or external syntax, publication must explicitly reserve an appropriate native alias or wrapper, or report unsupported publication. The kernel should not silently discard the requested root or guess how to expose it.

This also gives `Roots()` a clear meaning: declarations requested in the native API, rather than every declaration encountered while resolving dependencies.

### Generalize the definition subject

A deferred definition needs to describe either an open source declaration or a plugin-owned helper.

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

The ownership key is conceptually:

> Invocation + selected binding instance + open declaration identity or helper key.

It must not be based only on `PluginID`, a display name, or whichever concrete application happens to be visited first.

`Own` accepts copied metadata, **not a definition closure**. The host invokes the already-selected plugin’s `Define` method. This avoids both incomparable callback bodies and the temptation to capture `Page<Text>` while defining a supposedly open `Page<T>`.

Repeated requests for the same ownership key must agree on placement, declaration form, parameter signature, and binder correspondence. Conflicts are errors.

### Keep the two binding environments separate

There are two different mappings:

1. **Source binder → effective source argument**, owned by the source session.
2. **Source binder → target parameter type**, owned by the lowering context.

Before `Define`, the host creates a fresh definition context and installs the second map using `BinderLink` and the reserved target parameters. Helper definitions start without captured source binders.

A nested `Lower` retains the lexical parameter environment but switches the active ownership namespace to the newly selected mapping. Thus a record mapping’s presence helper belongs to that record mapping, even when reached while defining `Page<T>`.

Neither parameter names nor coincidentally matching slice positions establish binder identity.

## 4. Use a thin native fragment and symbol layer

The fragment mechanism should add ownership and symbolic relocation to native syntax. It should not become a new language-neutral type calculus.

The essential opaque handles are:

```go
type Scope[T any] struct { /* private session identity */ }
type Symbol[T any] struct { /* private session identity */ }

type Fragment[T any] struct { /* immutable native syntax template */ }
type Type[T any] struct { /* fragment classified as a type use */ }
type Decl[T any] struct { /* fragment classified as a declaration */ }

func (t Type[T]) Fragment() Fragment[T]
func (d Decl[T]) Fragment() Fragment[T]
```

A small construction surface can expose ordinary native nodes and deferred symbol occurrences:

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

Native packages provide convenient functions over that surface, such as:

```go
gotarget.String(builder, origin)
gotarget.Slice(builder, element, origin)
gotarget.Named(builder, symbol, arguments, origin)
gotarget.DefinedType(builder, owned, body, origin)
```

These are target-owned helpers. The generic kernel never decides that an arbitrary fragment can accept native generic arguments.

`SymbolUse` is a target-owned construction placeholder. The target assembler expands it into the appropriate native name, qualified path, or other structured reference syntax. Its vocabulary is not a generic semantic dispatch mechanism.

### The target lifecycle

The target adapter needs placement and native declaration information before reserving symbols:

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

`Storage[T]` is a host-supplied construction capability so independent adapters can mint valid opaque handles:

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

The contracts package must supply the underlying handle issuers and copying fragment constructors; the host supplies them with invocation limits. This prevents private fields from accidentally sealing the extension boundary.

Only the host receives the target lifecycle operations. Plugins obtain `Context.Syntax`, `Own`, and `External`.

For the initial profile, owned declarations can be restricted to configured native declaration scopes. Parameter constraints are constructed in `Define`, after parameter symbols exist, and checked against the reserved header. Arbitrary nested declaration scopes and fragment-local type or lifetime binders can remain deferred.

### Required fragment invariants

**Names must remain symbolic until allocation.** Fixing only reference edges is insufficient: a collision rename would leave stale identifier attributes.

**Each type occurrence gets fresh containment.** Reusing an immutable `Type[T]` template is allowed; reusing its eventual contained nodes is not.

**Definitions attach once.** Declaration and parameter symbols have unique primary anchors in the assembled output.

**Free symbols remain scope-checked.** A type referring to `Page<T>`’s parameter cannot be moved into an unrelated declaration.

**Target scope rules remain target-owned.** Go package declarations and import names have different scopes; imports are file-scoped. A generic “one name table per tree” policy would therefore be inadequate. [The Go Programming Language](https://go.dev/ref/spec?utm_source=chatgpt.com)

If fragment-local type or lifetime binders are later supported, materialization must clone those binders and redirect internal references together. Ordinary callable parameter-name syntax need not be deferred when it introduces no such symbolic dependencies.

### External symbols require special care

A locked external artifact is not automatically an AST declaration node.

When no declaration tree is supplied, the assembler emits imports/includes and qualified native syntax. It must not fabricate an `ast.NodeRef` for the external declaration.

Import aliases and qualification are resolved for each materialization unit, even if the external symbol handle itself is shared.

Also, the existing `ast.Set[T]` return value is not a build lockfile, native `Source[T]`, or semantic receipt. The caller retains responsibility for the configured external dependency environment.

## 5. Use deferred reservation rather than a mandatory public `Plan`

The necessary planning discipline can be implemented through the reservation queue:

1. **Validate configuration.** Check domains, duplicate bindings, replacement guards, and source/target endpoint compatibility.
2. **Admit the invocation.** Reuse `transform.NewTransformer` and its explicit evidence contract; open fresh source and target sessions.
3. **Publish requested declarations.** Select their mappings before body expansion and call `Publish`.
4. **Reserve on demand.** The first `Own` reserves the declaration and parameter symbols, records the immutable specification, and queues its definition.
5. **Return immediately.** `Own` does not recursively execute `Define`.
6. **Drain definitions.** Invoke the selected plugin’s `Define` once for each reservation, with its open formal environment. Definitions can request additional resources.
7. **Assemble after closure.** Allocate deterministic names and import aliases, expand symbolic uses, clone occurrences, and resolve local links.
8. **Return atomically.** Any error discards the complete invocation output.

Internally, resources can have `reserved`, `defining`, and `defined` states. Referencing a reserved symbol is distinct from recursively expanding a type expression.

An inline lowering cycle that never crosses an owned or external symbol boundary is an error. Native legality of a symbolic cycle is a separate matter.

A mandatory `Plan` method on every plugin would add another public phase without removing these requirements. It becomes worthwhile if a concrete feature needs whole-program information before representation construction. The cases in this document can use the smaller queue protocol.

## 6. Strengthen the generic coherence rule

The proposed deferral of application-specific selectors is correct, but insufficient.

A plugin selected by the open `list-of` key could inspect its argument and return `ByteBuffer` for `Byte`, while returning a generic collection otherwise. That recreates the same incoherence without introducing an application selector.

The profile therefore needs an explicit law:

\[
L_C(\tau[\sigma])
\;\approx_P\;
L_C(\tau)[L_C(\sigma)]
\]

Here:

- \(C\) is the immutable selected configuration.
- \(\sigma\) is a substitution within the profile’s represented domain.
- \(\approx_P\) is the representation relation promised by the profile.

Agreement can allow alpha-renaming and administrative aliases. It cannot silently permit a different storage family or an unprovided conversion.

Consequences follow directly:

- A mapping defines one representation family over its supported substitutions.
- Eligibility of an open family is checked under its binder assumptions and the profile’s stated domain.
- Success for several concrete applications does not establish open eligibility.
- Argument inspection must not secretly change native parameter structure, helper identity, ownership, or indirection.
- An owned generic definition is emitted from the open declaration, never from the first concrete application.
- Use-site syntax remains plugin-owned.
- Source kinds constrain admissibility; they do not select representation families.

This is a **trusted-plugin obligation**. A Go interface cannot prove parametricity of arbitrary plugin code. A `Plan` API would make some violations less convenient, but would not prove the law either.

The profile’s domain must also be precise. Supporting a generic proper-type parameter does not mean that every possible Horos type is now supported. Quantify the promise over the explicitly represented substitutions, and reject unsupported applications consistently.

### Recommended first-version boundary

| Support | Defer |
|---|---|
| Exact builtin and ordinary declaration selectors | Application-specific selectors |
| Proper-type parameters and admitted direct applications | Higher-order or computed constructors |
| Uniform inline, owned, and external representation families | Argument-sensitive specialization |
| Explicit source-to-target binder correspondence | Automatic monomorphization |
| Source-owned eligibility from captured evidence | Predicate or most-specific-kind dispatch |
| Helpers with explicit native parameters | Helpers specialized implicitly per concrete use |
| Named declaration publication | Anonymous type-entry publication until its target policy is specified |

## 7. The builtin text-wrapper case works without a source declaration

The selected text override requests a helper:

```go
DefinitionSpec[Horos, Go]{
    Subject: HelperDefinition[Horos]("text-wrapper"),
    Header: DeclarationHeader[Go]{
        Placement: Placement{
            Unit: "types",
        },
        Form:     gosyntax.KindTypeDeclaration,
        NameHint: "Text",
    },
}
```

The host reserves its symbol under the selected text binding and queues that binding’s `Define`.

`Use` returns a native named-type reference to the reservation. `Define` recognizes helper `"text-wrapper"` and constructs the target declaration:

```go
type Text string
```

The crucial detail is that `Define` constructs the native backing `string` directly. It does **not** call `Lower` on Horos text. Doing so would select the same override again and could produce a cycle or an invalid self-definition.

Using native primitive syntax inside the selected wrapper is ordinary implementation logic of that representation.

The list plugin continues to lower its source element through the effective table. Consequently:

```text
list-of(text)
```

produces a slice type referring to the helper:

```go
[]Text
```

There is one helper definition and fresh contained reference syntax at each occurrence.

`type Text string` does not enforce UTF-8: Go strings contain bytes. This wrapper therefore fits the admitted-source-use profile unless additional enforcement is separately provided. [The Go Programming Language](https://go.dev/ref/spec?utm_source=chatgpt.com)

The presence helper uses the same mechanism with an explicit native parameter:

```go
type Presence[T any] struct {
    Present bool
    Value   T
}
```

Its internal `Present` flag is native implementation metadata. Constructing that native `bool` need not reenter the Horos boolean mapping. Conversely, a source boolean field must pass through the selected source mapping.

That distinction prevents overrides from accidentally rewriting the implementation machinery of another mapping.

## 8. Recursion, compatibility, and evidence remain separate obligations

### Selection must precede representation-specific cycle rejection

A selected external representation may avoid a recursive source body entirely. Rejecting the raw source graph before selection would incorrectly rule that out.

This cannot rescue source input for which required Horos admission is absent or negative. Representation selection starts after the required source receipts have been admitted.

Reservation also does not establish native legality. Rust, for example, imposes restrictions on recursive types, including finite size and nominal structure in the recursion. Symbol availability alone proves neither. [The Rust Reference](https://doc.rust-lang.org/reference/types.html?utm_source=chatgpt.com)

The smallest initial profile can continue rejecting generated recursive definitions. It should do so on the demanded representation dependencies after overrides have determined which bodies need expansion.

### Deterministic selection does not guarantee compatible composition

A parent mapping can impose native requirements on its children. An external artifact can also embody fixed representation assumptions that an arbitrary child override does not preserve.

For the first profile, specify a conservative target admissibility contract: returned types must be usable in the supported native positions, and external integrations must declare their representation assumptions. A selected plugin should reject incompatible children or artifacts.

This need not introduce a general capability system. It does require acknowledging that grammatical ASTs and exact keys are insufficient to prove semantic compatibility.

### Evidence requirements must cover plugin eligibility

`Source.Contract()` must cover the facts required by the configured first-profile mappings.

If a new plugin needs additional evidence, composition must explicitly provide the corresponding source contract and receipts. A plugin must not run an analyzer, retrieve dependencies, or reinterpret missing facts as success.

A future per-mapping evidence declaration may be useful. It is unnecessary initially if the documented source contract already covers the supported plugins’ obligations.

### Bound expanded work and isolate invocations

Preserve the current output-node, depth, and tree-count boundaries, and add bounded definition/import/work accounting as needed.

Count materialized occurrences, expanded symbolic uses, and payload bytes. Counting only immutable templates would miss potentially large expansion.

Sequential execution within one invocation is sufficient. A configured transformer may still receive concurrent invocations, so plugin configuration must remain immutable and mutable state must live in fresh sessions.

Cancellation remains cooperative. Host operations can check cancellation and limits; these interfaces cannot preempt arbitrary Go code that never returns.

## 9. Package boundaries and composition

The proposed dependency direction is appropriate. I would refine the responsibilities as follows:

| Path | Responsibility |
|---|---|
| `contracts/generation/go/selector.go` | Selector identity and copying constructors |
| `contracts/generation/go/source.go` | References, source views, source/session contracts |
| `contracts/generation/go/plugin.go` | Bindings, replacements, use and publication requests |
| `contracts/generation/go/ownership.go` | Definition subjects, headers, binder links, reservations |
| `contracts/generation/go/fragment.go` | Immutable native templates, symbolic slots, handle construction |
| `contracts/generation/go/target.go` | Target/session and storage contracts |
| `generation/go/bindings.go` | Exact selection and checked replacements |
| `generation/go/lower.go` | Scoped lowering and terminal error propagation |
| `generation/go/definitions.go` | Ownership keys and deferred definition queue |
| `generation/go/transform.go` | Existing transformer adaptation and atomic completion |
| `generation/source/horos/go/` | Horos identities, substitutions, payload reading, evidence interpretation |
| `generation/target/<language>/go/` | Native builders, scopes, imports, declaration checks, AST assembly |
| `transforms/horos-to-<target>/go/` | Thin factories requiring supplied composition |
| Horos-owned integration module | Builtin mappings and optional opinionated composition |
| Library integration packages | Exact mappings for their ordinary declarations |
| Consumer | Binding selection, replacements, artifacts, and representation profile |

The structural declaration policy also interprets source declarations and chooses native declaration forms. Its implementation belongs in the relevant integration layer, even though the generic configuration has a slot for it.

For the illustrative composition:

| Selector or policy | Supplied mapping |
|---|---|
| Horos `utf8-text` | Native string mapping |
| Horos `bool` | Native boolean mapping |
| Horos `list-of` | Native collection mapping |
| Horos closed record descriptor | Native record-body mapping |
| Structural declaration policy | Publish a named declaration and lower its body |
| Exact Function declaration | Function library’s selected representation |
| Explicit Function replacement | Consumer-selected invocation-handle representation |

The integration or consumer package imports those implementations. Bittype’s generic generation package and engine do not import the integration module or a catalogue of library declarations.

One useful architectural restraint is to keep the new SPI experimental until these concrete cases fit it. With only Horos as a demonstrated source, a richer generic semantic model would add commitments without evidence that they are needed.

## 10. What the prior art supports—and what this proposal must establish itself

MLIR supports the distinction between type conversion and value materialization, and its documentation makes context-dependent caching concerns explicit. Its conversion mechanism also permits candidate retry and reverse registration precedence. Those selection rules should not be imported into this explicit-replacement design. [MLIR](https://mlir.llvm.org/docs/DialectConversion/?utm_source=chatgpt.com)

OpenAPI Generator demonstrates consumer-supplied references to existing model types, while noting language-dependent support. That is useful precedent for external representations, but not evidence of canonical source identity or semantic equivalence. [OpenAPI Generator](https://openapi-generator.tech/docs/customization/?utm_source=chatgpt.com)

Rust’s coherence rules illustrate why overlap needs a defined policy. Bittype can establish that policy through consumer-scoped exact selection and checked replacement; it does not need Rust’s global orphan restrictions. [The Rust Reference](https://doc.rust-lang.org/reference/items/implementations.html?utm_source=chatgpt.com)

**The proposal’s essential public model should be: exact selection, environment-bearing source references, proper type uses, explicit publication, and owned definition resources.** The reservation queue supplies the planning discipline. Native adapters supply syntax and scope rules. The representation profile supplies the semantic promise—including the uniform-family rule that exact dispatch alone cannot enforce.