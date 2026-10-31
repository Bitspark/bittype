# Syntax engine and AST transformations

**Implemented, 31 October 2026.** This increment supplies the frame for registered
AST transformations. It implements no canonical AST and no TypeScript, Go, Rust
or C++ mapping. Applications explicitly register the existing native text
providers and future transformer implementations.

## Public components

| Package | Responsibility |
| --- | --- |
| [trees/representation/go](../trees/representation/go/ast.go) | Checked AST/CST envelopes and language/profile descriptors |
| [contracts/transform/go](../contracts/transform/go/transformer.go) | ASTTransformer and typed stage registration bridge |
| [engine/go](../engine/go/engine.go) | Registration, route planning/execution and text-stage composition |
| [contracts/language/go](../contracts/language/go/language.go) | Shared ID/Name identity contract |
| [contracts/pipeline/go](../contracts/pipeline/go/stage.go) | Independent typed Stage function contract |
| Existing native providers | Language-specific parsing, checking, formatting and printing |

The engine, transformation boundary and their production dependencies import no
concrete language, actual transformer or behavioral contract. There is no
language list, canonical route preference, automatic discovery or init-time
registration. Function, Stream and Cell remain ordinary library references.

## Identity, profiles and storage

The existing ast.AST[L] and cst.CST[L] carriers remain unchanged. L separates
Go types. A dynamic registry needs an opaque envelope, not ast.AST[any].

representation.NewLanguage[L](checker, profile) creates a Language descriptor from a
text.CheckerProvider[L]. It captures ID and Name, the actual Go marker type, the
explicit profile string and the checker. Configured descriptors implement
LanguageProvider; their zero values are invalid. Native marker providers still
have usable zero values. Empty metadata and typed nil providers are rejected.

A profile names the admitted syntax/representation revision chosen by the
application, such as native-types-v1. It is not a generated fingerprint,
toolchain detector or published canonical language revision. Matching requires
all three of language ID, profile and Go marker. Matching names or IDs alone
cannot cast an AST into another language. Descriptor compatibility is the
provider owner's promise about its checker/profile; the engine cannot prove
two arbitrary checker implementations have identical semantics.

Pack[L](descriptor, view) accepts a native AST or its existing CST view,
checks carrier structure and native grammar, and owns a deep snapshot.
Unpack[L](descriptor, envelope) verifies the endpoint and returns fresh mutable
storage. AST.Snapshot also copies deeply. Copying an envelope value does not
expose its private arena. Zero envelopes are invalid.

Attribute, child, reference and origin slices all retain their contents and
addresses in a snapshot. The frame does not resolve references, rename nodes or
invent new declaration identities. A transformer that changes an arena or tree
identity owns the corresponding reference/provenance mapping.

The existing carrier limits apply before data is copied: node and depth limits,
aggregate string bytes and list-entry counts. Checker callbacks receive private
copies. Typed stages receive owned input, and their output is checked and
snapshotted, including when the stage retains a reference to its output.

## Transformer contract

The public interface is:

```go
type ASTTransformer interface {
    FromID() language.ID
    ToID() language.ID
    Transform(representation.AST) (representation.AST, error)
}
```

FromID and ToID are stable language endpoint keys. Transform is a potentially
partial semantic operation. An unsupported construct or missing migration
decision returns an error; successful syntax checking alone does not prove a
mapping's semantic equivalence.

NewTransformer[F, T](from, to, stage) adapts:

```go
pipeline.Stage[ast.AST[F], ast.AST[T]]
```

The bridge uses representation.CheckMarker to check actual marker types at
construction, without access to payloads or unchecked constructors. It checks input, passes the
stage an owned typed AST, and packs its result against the destination
descriptor. A nil stage is an error. The adapter returns zero output whenever
the stage fails, even if the stage also supplied a partial AST.

The adapter additionally implements DescribedTransformer, which exposes
FromLanguage and ToLanguage descriptors. Registration checks these typed/profile
endpoints early. A separately authored minimal ASTTransformer remains supported;
execution checks its input and result against registered descriptors regardless.
An implementation cannot bypass output type/profile checking by returning a
valid envelope for a different language.

The engine does not recover panics, sandbox callbacks, fetch dependencies or
verify arbitrary semantic claims. Providers must follow their stage contracts,
keep metadata/checker configuration stable and return diagnostics for expected
unsupported inputs. Concurrent callers require concurrency-safe providers.

## Registration

Engine's zero value is usable; New is a convenience constructor. Both start empty.

- RegisterLanguage installs an AST-only descriptor.
- NewProvider[L](textProvider, profile) packages a full existing text.Provider[L].
  RegisterProvider installs its descriptor and text capabilities atomically.
- RegisterTransformer installs one directed edge between registered languages.
- Language(id) looks up a descriptor; Languages returns a sorted independent list.

Language keys are exact and case-sensitive. Duplicate keys or directed edges
are errors, including attempts to replace a language with another profile or
marker. Use a separate explicitly composed engine for an alternative mapping
profile. Choose a full provider at registration when text capabilities are
needed; registration does not silently upgrade or replace an AST-only entry.

Self-edges are rejected: a same-language request is a checked snapshot, not
implicit normalization. Graph cycles across different languages are allowed,
but selected paths cannot repeat languages. There are no removal/replacement
operations or hidden fallback registrations.

Registration is synchronized and additive. Operations capture their registered
dependencies together, then release the lock before invoking user code.
Callbacks can perform registry reads without deadlocking. The Engine must not
be copied after use. This synchronization does not make a stateful plugin safe
for concurrent use by itself.

## Plans and execution

Plan(from, to) finds the unique shortest directed path. A breadth-first search
counts shortest routes, saturating counts at two; it does not enumerate all
routes. Missing connectivity and multiple shortest routes produce distinct
ErrNoPath and ErrAmbiguousPath errors. Registration order never breaks a tie.

PlanVia(ids...) accepts an explicit sequence including both endpoints. It checks
that every language and edge exists and that no language is repeated. It can
select a longer route. Plan.Languages returns a copy, and FromID/ToID expose its
endpoints. The zero plan is invalid.

Execute(input, plan) revalidates the route in the executing engine. Plans are
route values, not permissions or opaque bindings to the engine that planned
them. Registered profiles and actual input/output types are checked again.
TransformTo(input, targetID) is planning followed by execution.

Execution validates the initial AST, takes an owned snapshot, invokes each
selected edge once, validates its destination and snapshots its result before
continuing. A zero-edge identity plan also validates and snapshots. Changed
transformer endpoint metadata after registration is an error.

A failed edge returns zero output and a StepError containing its zero-based edge
index and From/To IDs. Unwrap preserves the original error, including native
located diagnostics. There is no alternate-route retry after failure, no usable
partial result and no silent weakening of a representation. A route establishes
connectivity; it does not establish that every possible source AST is supported.

## Native text composition

NewProvider adapts existing implementations; it creates no new native grammar.
The engine offers:

```go
Parse(languageID, text.Input) (representation.CST, error)
Check(representation.AST) error
Format(representation.AST) (representation.CST, error)
Print(representation.CST) (string, error)
Render(representation.AST) (string, error)
```

Render is Format followed by Print. These methods select the registered provider
by the envelope's language ID. An AST-only registration returns
ErrMissingCapability for text operations. Unregistered languages return
ErrUnknownLanguage.

The representation.CST envelope exposes ASTView() representation.AST without
another conversion provider. representation.PackCST and representation.UnpackCST
bridge typed concrete trees, validating
and copying abstract/concrete arenas, leaf pointers, links and metadata.
A parser or formatter error discards its returned CST. A printer error discards
its returned text. Providers remain responsible for their grammar, concrete
spelling and formatter preservation contracts; structural checks do not prove
that arbitrary leaf text means the stored abstract view.

For example, with explicit imports of engine, text, the TypeScript syntax
package and its provider:

```go
provider, err := engine.NewProvider[typescript.Language](
    typescriptprovider.Provider{}, "native-types-v1",
)
if err != nil { return err }
syntax := engine.New()
if err := syntax.RegisterProvider(provider); err != nil { return err }

parsed, err := syntax.Parse(provider.ID(), text.Input{
    SourceID: "input", TreeID: "input", Text: "type Page<T> = T[];",
})
if err != nil { return err }
original, err := syntax.Print(parsed)
if err != nil { return err }
formatted, err := syntax.Render(parsed.ASTView())
if err != nil { return err }
// original retains source bytes; formatted has the provider's selected layout.
_ = original
_ = formatted
```

After registering a destination provider and actual mapping, callers can pass
parsed.ASTView() to TransformTo, then Render the resulting AST. No actual mapping
or canonical registration is installed by this increment.

## Qualification and remaining work

Independent test languages exercise the public boundary without core changes.
The graph oracle enumerates all 4,096 directed graphs on four languages and
compares all endpoint pairs against independently enumerated simple paths.
Additional tests cover shared-tail ambiguity, cycles, explicit routes, missing
edges, identity operations, endpoint/profile mismatches, failure propagation,
no fallback, snapshot metadata, retained callback storage and concurrent
registration/reentrant reads.

Public consumer fixtures compile against actual production sources and reject
wrong typed stages/views, missing transformer endpoints and forged envelopes.
Production dependency tests prohibit concrete languages and mapping packages
throughout the frame's import closure. All four existing native providers also
run through the engine, including parse/print retention and formatting.

The maintained GitHub gate runs these tests, vet and a race-detector check for
concurrent registry operations, alongside existing native compiler and fuzz
checks. Documentation checks run on the same exact committed revision.
Workloads run remotely.

The [stage-first hierarchy](../LAYOUT.md) keeps checked tree representations
independent of transformer contracts and the engine. Shared conformance tests
and discovered stage/dependency checks enforce this boundary. The moved APIs
have no aliases at the former package paths.

The [canonical model](model.md), its frontend and actual language transformers
remain the next semantic increment. Public builders, a generic Stage composition
helper and an application CLI are also deferred; none is needed to register and
execute a correctly implemented ASTTransformer through this frame.
