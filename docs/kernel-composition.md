# Kernel-owned generation composition

**Owner decision, 6 November 2026; implementation pending.** The kernel composes
supplied source, target and mapping providers from an explicit request. Libraries
construct their own invocation state and own representation meaning. Consumers
select capabilities and requested output. No library supplies an assembled
pipeline; the kernel has no language or named-contract catalogue.

This replaces the whole-Composition factory in
[automatic composition](generation-composition.md), preserving the
[generation operations](generation-spi.md), admission, real binder identities,
association owners, definition scheduling and whole-invocation failure. It does
not repair [verification ACX-001](architecture-cases/0001-verification-composition.md),
add native constraint enforcement, or generate runtime bindings.

## Uniform layout

Implementations follow `<stage>/<dimension>/<scope>/<component>/go`.
Scope is one language, an ordered pair, or shared. The final go means implemented
in Go. Directory names consistently use golang, including pairs; language IDs
do not change. Only implemented capabilities get directories.

```text
contracts/generation/go/                    operations and construction SPI
contracts/transform/go/                     ASTTransformer and typed adapter
engine/go/                                 checked requests and registration
parse/syntax/<language>/provider/go/
parse/files/horos/provider/go/
check/syntax/<language>/provider/go/
check/syntax/shared/{ast,cst}/go/
format/syntax/<language>/provider/go/
format/files/horos/provider/go/
print/{syntax,files}/shared/provider/go/
resolve/files/<language>/provider/go/
analyze/files/<language>/provider/go/
transform/syntax/shared/generation/go/       registry and automatic composition
transform/syntax/shared/target/go/           symbolic native storage
transform/syntax/horos/source/go/            semantic reader
transform/syntax/horos/mapping/go/           shared Horos mapping policy
transform/syntax/<native>/target/go/         native AST construction
transform/syntax/horos-to-<native>/mapping/go/ supplied representations
```

Generic adapters stay in contracts/kernel. Remove aggregate provider/language
wrappers, the separate transforms and generation trees, pair-specific factories
and default compositions, without aliases. Vocabularies, trees, semantic models
and verification retain their responsibilities. Source readers may depend on
their semantic providers. Native targets never import source semantics or
mappings. Mappings use the source and target APIs they need. The shared kernel
imports only contracts and carriers, never concrete integrations.

## Construction SPI

F/T are language markers. S/N are provider-owned APIs, not a new universal type
algebra. The proposed surface extends existing concepts:

```go
type SourceSession[F any] interface {
    Source[F]
    Catalog[F]
    ResolveRoots(context.Context, []ast.NodeRef) ([]Declaration[F], error)
}
type SourceProvider[F any, S SourceSession[F]] interface {
    Contract() transform.Contract
    Open(context.Context, transform.Request[F]) (S, error)
}
type TargetProvider[F, T any, N Native[F, T]] interface {
    Open(context.Context, TargetRequest) (N, error)
}
type MappingProvider[F, T any, S SourceSession[F], N Native[F, T]] interface {
    Open(context.Context, MappingRequest[S, N]) (Contribution[F, T], error)
}
type Contribution[F, T any] struct {
    Bindings   []Binding[F, T]
    Structural *Mapping[F, T]
    Helpers    []Plan[F, T]
}
```

Catalog, selectors, mappings, bindings and plans have one public definition in
the generation contracts. Use, Publish, Define, Introduce/Apply and association
checks keep their meanings. Providers return only their own components; they
cannot select another provider, supply invocation roots, run the composer or
return a complete Composition. Plan factories remain declaration-local setup.

TargetRequest and MappingRequest are read-only kernel-created objects with an
invocation identity, output/target metadata and limits. MappingRequest exposes
the actual S/N instances opened for this invocation. It cannot be constructed
by attaching unrelated values to a claimed ticket. Zero, expired and foreign
construction requests fail. No public arbitrary object map or AST[any] is added.

Provider-specific settings are configured at registration: Go owns Package,
C++ owns Namespace. Registration captures descriptors and evidence requirements.
Providers must keep their configuration stable; copying a Go interface does not
freeze its object. Open creates fresh invocation state, with no hidden analysis
or IO. These are trusted in-process collaborators, not a sandbox.

## Registration and request

The kernel exposes generic free functions (the S/N constraints are as above):

```go
RegisterSource[F, S](registry, language, provider) (SourceRef, error)
RegisterTarget[F, T, N](registry, from, to, provider) (TargetRef, error)
RegisterMapping[F, T, S, N](registry, from, to, provider) (MappingRef, error)
```

These supply Go instantiations once, without handwritten assembly. The kernel
retains actual F/T/S/N Go type identity and exact language/profile descriptors.
An intentional API adaptation is explicit; display names or an unrelated value's
extra methods do not establish a declared construction contract. References name
particular registry-owned registrations. Zero, foreign and unknown handles fail.
Registration order cannot select a representation. Request snapshots freeze the
selected records; all provider callbacks run outside registry locks.

```go
type BindingRef struct { Provider MappingRef; ID string }
type Replacement struct {
    Selector Selector
    Expected BindingRef
    With BindingRef
}
type Selection struct {
    Source SourceRef
    Target TargetRef
    Mappings []MappingRef
    Replace []Replacement
}
type Request struct {
    Input transform.Input
    To representation.Language
    Roots []ast.NodeRef
    Output ast.TreeID
    Selection Selection
    Limits Limits
}
```

Roots are captured declaration sites, resolved and authenticated by the source
integration, not unqualified names. Empty roots are invalid, not export-all.
Output identities cannot collide with input trees. Each exact selector has one
effective binding. A replacement names both the expected and replacement
binding qualified by provider, and must retain the selector. Unknown, repeated,
stale or unused replacements fail. At most one structural policy is selected.
Helper collisions fail rather than sharing authority. Provider/mapping failure
is terminal, with no fallback.

## Kernel lifecycle

1. Snapshot request configuration and validate context, limits, handles, language
   profiles, Go markers/API types, roots and output identity before provider work.
2. Admit supplied evidence through the existing typed transform boundary against
   the source's captured requirements. Missing evidence never runs analysis.
3. Open source and target once. Authenticate the exact source snapshot/receipts.
   Open each selected mapping once against these same instances, with bounded
   work and invocation-local issuance.
4. Collect contributions, validate replacements, resolve/authenticate roots and
   build the internal composition. No provider returns that composition.
5. Publish roots, acquire definitions, introduce formals, compose child uses and
   applications, and drain definitions through the existing guarded SPI.
6. Finish a complete AST set and apply the engine's output syntax admission.
   Errors/cancellation invalidate the invocation. Construction tickets expire
   when their Open call returns; generation tickets retain their existing rules.

Engine.Generate is the direct generation entry. A generic adapter constructs a
request-bound ASTTransformer, without per-pair factories or global registration
per request. Existing TransformTo route planning stays separate. Generation
does not invent receipts for subsequent legs. Formatting/printing remain explicit.

## Horos to TypeScript

```horos
module example;
type Person = { name: text; };
type Page<T> = { items: list<T>; next?: text; };
type People = Page<Person>;
```

Register the Horos reader, TypeScript target and selected text, bool, list,
record, binder and named-declaration mappings. A library may expose a collection
of providers; it does not wire a whole pipeline. The representation profile is
explicit provider configuration; the existing enforcement profile stays
unsupported. Request the captured Person/Page/People sites, output identity,
limits and those registrations.

The reader exposes resolved declarations, original bounds and actual binders.
The named mapping introduces T; the record mapping requests its fields; the list
mapping requests the element and receives T's scoped product. People applies
Page's retained convention to Person. Formatting produces:

```typescript
export interface Person { name: string; }
export interface Page<T> { items: T[]; next?: string; }
export type People = Page<Person>;
```

An independent ReadonlyArray mapping replaces list through explicit selection.
Cell uses the same SPI and may stop expansion for its exact declaration. Adding
either changes only its provider and consumer configuration, not shared code,
source/target implementations or unrelated mappings. Builtins have no privileged
path. Representation generation does not establish native constraint enforcement.

## Qualification

Keep four-target pagination/compiler and installed-consumer tests. Exercise
foreign handles; incompatible markers/profiles/API contracts; missing/foreign
evidence; unknown/duplicate roots; output collisions; conflicting/stale
replacements; helper collisions; provider failure; cancellation; expired tickets;
concurrent invocations; and bounded work. Qualify a materially different
independent mapping, nesting and terminal failure. Test the kernel with unrelated
language markers as well as Horos. A green example does not replace ownership
and change-propagation review under AC-01 through AC-08.

Migrate consumers, directed imports, compiler paths and current documentation
links. Update source-derived reporting with actual layout/contracts; preserve
public-method identity, physical source links and separate delegate evidence.
No handwritten capability catalogue, fake route or blanket discovery exemption.
Historical research/source links remain historical.
