# Syntax engine and AST transformations

The engine registers explicit syntax capabilities and directed transformations. Engine.Generate composes a request from a source, target, selected mappings, receipts and declaration sites. The [generation contract and walkthrough](kernel-composition.md) show Horos to TypeScript. TransformTo independently plans routes over explicitly registered transformer edges.

The [Horos frontend](horos.md) registers its syntax checker independently
from language identity, [resolution](resolution.md) and [semantic analysis](analysis.md).
Horos text capabilities are independently registered in the same way. Go, TypeScript and Rust have [bounded native semantic providers](native-semantics.md). C++ supplies syntax stages.

## Public components

| Package | Responsibility |
| --- | --- |
| [trees/representation/go](../trees/representation/go/ast.go) | Structurally checked AST/CST envelopes and language/profile descriptors |
| [contracts/transform/go](../contracts/transform/go/transformer.go) | ASTTransformer and typed stage registration bridge |
| [engine/go](../engine/go/engine.go) | Registration, route planning/execution and text-stage composition |
| [contracts/language/go](../contracts/language/go/language.go) | Shared ID/Name identity contract |
| [contracts/check/go](../contracts/check/go/provider.go) | Independent AST syntax-admission contract |
| [contracts/parse/go](../contracts/parse/go/provider.go) | Source input and text-to-CST contract |
| [contracts/format/go](../contracts/format/go/provider.go) | AST-to-CST contract |
| [contracts/print/go](../contracts/print/go/provider.go) | CST-to-text contract |
| [contracts/text/go](../contracts/text/go/diagnostic.go) | Shared syntax and source diagnostics |
| [contracts/pipeline/go](../contracts/pipeline/go/stage.go) | Independent typed Stage function contract |
| Existing native providers | Language-specific parsing, checking, formatting and printing |

The engine, transformation boundary and their production dependencies import no
concrete language, actual transformer or behavioral contract. There is no
language list, canonical route preference, automatic discovery or init-time
registration. Function, Stream and Cell remain ordinary library references.

## Identity, profiles and storage

The existing ast.AST[L] and cst.CST[L] carriers remain unchanged. L separates
Go types. A dynamic registry needs an opaque envelope, not ast.AST[any].

`representation.NewLanguage[L](identity, profile)` accepts `language.Provider`.
It captures ID, Name, the actual Go marker type and the explicit profile string.
It installs no checker or other capability. Configured descriptors implement
language.Provider; their zero values are invalid. Empty metadata and typed nil
providers are rejected. Native language markers are usable identity providers.

The four stage contracts embed only language identity and their own typed method.
A parser implementation need not implement checking, formatting or printing.
An implementation may offer several methods, but the application registers each
selected capability explicitly. Resolution and analysis remain separate too.

A profile names the admitted syntax/representation revision chosen by the
application, such as native-types-v1. It is not a generated fingerprint,
toolchain detector or published canonical language revision. Matching requires
all three of language ID, profile and Go marker. Matching names or IDs alone
cannot cast an AST into another language. Descriptor compatibility is the
provider owner's promise about its checker/profile; the engine cannot prove
two arbitrary checker implementations have identical semantics.

Pack[L](descriptor, view) accepts a native AST or its existing CST view,
checks generic carrier structure and owns a deep snapshot. It does not invoke a
language checker or establish syntax admission. `PackCST` likewise checks both
carrier arenas, not native grammar. Call a configured `check.Checker` or
`Engine.Check` when syntax admission is required.
Unpack[L](descriptor, envelope) verifies the endpoint and returns fresh mutable
storage. AST.Snapshot also copies deeply. Copying an envelope value does not
expose its private arena. Zero envelopes are invalid.

Attribute, child, reference and origin slices all retain their contents and
addresses in a snapshot. The frame does not resolve references, rename nodes or
invent new declaration identities. A transformer that changes an arena or tree
identity owns the corresponding reference/provenance mapping.

The existing carrier limits apply before data is copied: node and depth limits,
aggregate string bytes and list-entry counts. Checker callbacks receive private
copies. Typed stages receive owned input, and their output is structurally checked and
snapshotted, including when the stage retains a reference to its output.

## Transformer contract

The public interface is:

```go
type ASTTransformer interface {
    FromID() language.ID
    ToID() language.ID
    Transform(context.Context, transform.Input) (representation.ASTSet, error)
}
```

FromID and ToID are stable language endpoint keys. Transform is a potentially
partial semantic operation. An unsupported construct or missing migration
decision returns an error; successful syntax checking alone does not prove a
mapping's semantic equivalence.

NewTransformer[F, T](from, to, contract, stage) adapts:

```go
func(context.Context, transform.Request[F]) (ast.Set[T], error)
```

The bridge uses representation.CheckMarker to check actual marker types at
construction, without access to payloads or unchecked constructors. It verifies
input identity and explicitly supplied semantic evidence, passes the stage an
owned typed request, and structurally packs its
result against the destination descriptor. The standalone bridge does not invoke
syntax checkers; engine execution checks every endpoint explicitly. A nil stage is an error. The adapter returns zero output whenever
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

- RegisterLanguage installs identity/profile/marker metadata only.
- RegisterSyntaxParser accepts `parse.SyntaxParser`, constructed with
  `parse.NewSyntaxParser[L](definition, parse.Provider[parse.Input, cst.CST[L]])`.
- RegisterChecker accepts `check.Checker`, constructed with
  `check.NewChecker[L](definition, check.Provider[L])`.
- RegisterSyntaxFormatter accepts `format.SyntaxFormatter`, constructed with
  `format.NewSyntaxFormatter[L](definition, format.Provider[ast.View[L], cst.CST[L]])`.
- RegisterSyntaxPrinter accepts `print.SyntaxPrinter`, constructed with
  `print.NewSyntaxPrinter[L](definition, print.Provider[cst.CST[L], string])`.
- RegisterFilesParser/Formatter/Printer independently attach the typed Files adapters
  through their heterogeneous registration bridges; see [Source](source.md).
- RegisterTransformer installs one directed edge between registered languages.
- RegisterResolver installs an optional resolution capability; Capture and Resolve
  produce an explicit input snapshot and resolution result.
- RegisterAnalyzer installs an optional semantic analysis capability for an
  existing compatible language; Analyze consumes a snapshot and any declared
  resolution result and returns its reports
  and optional shared model. See the [analysis contract](analysis.md).
- Language(id) looks up a descriptor; Languages returns a sorted independent list.

Language keys are exact and case-sensitive. Duplicate keys or directed edges
are errors, including attempts to replace a language with another profile or
marker. Use a separate explicitly composed engine for an alternative mapping
profile. Each capability attaches to an already registered compatible descriptor,
in any order, including after earlier operations. Unknown languages, duplicate
registrations of the same capability and ID/profile/marker mismatches fail without
changing the registry. Registering one capability never infers another.

The configured adapters are heterogeneous bridges around typed providers, like
resolution.Resolver and analysis.Analyzer. Their ID/Name/Language methods expose
captured metadata; their stage methods use owned boundary snapshots. Provider
identity must match the descriptor, and the actual Go marker is checked before
binding. `text.Provider`, `engine.Provider`, `engine.NewProvider` and
`Engine.RegisterProvider` are removed without compatibility aliases. The text
contract package contains diagnostics only.

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

Execute(ctx, input, plan) revalidates the route in the executing engine. Plans are
route values, not permissions or opaque bindings to the engine that planned
them. Registered profiles and actual input/output types are checked again.
TransformTo(ctx, input, targetID) is planning followed by execution.

Execution first requires registered checkers for every language in the route,
including intermediate endpoints, before invoking any provider or edge. A missing
checker yields ErrMissingCapability. Planning itself only checks connectivity.
Execution then checks the initial AST, takes an owned snapshot, invokes each
selected edge once, validates its destination and snapshots its result before
continuing. A zero-edge identity plan also validates and snapshots. Changed
transformer endpoint metadata after registration is an error.

A failed edge returns zero output and a StepError containing its zero-based edge
index and From/To IDs. Unwrap preserves the original error, including native
located diagnostics. There is no alternate-route retry after failure, no usable
partial result and no silent weakening of a representation. A route establishes
connectivity; it does not establish that every possible source AST is supported.

## Native text composition

Each configured stage adapter wraps an independently selected typed provider.
The engine offers:

```go
ParseSyntax(context.Context, language.ID, parse.Input) (representation.CST, error)
Check(representation.AST) error
FormatSyntax(context.Context, representation.AST) (representation.CST, error)
PrintSyntax(context.Context, representation.CST) (string, error)
RenderSyntax(context.Context, representation.AST) (string, error)
```

Parse requires only a parser; Check only a checker; Format only a formatter;
Print only a printer. Render captures the formatter and printer together before
invoking either, then formats and prints. It fails before formatting if either is
missing. These operations do not implicitly invoke other registered capabilities.
Native implementations may check their own grammar internally as part of their
stage contract. Missing selected capabilities return ErrMissingCapability;
unregistered languages return ErrUnknownLanguage.

The representation.CST envelope exposes ASTView() representation.AST without
another conversion provider. representation.PackCST and representation.UnpackCST
bridge typed concrete trees, validating
and copying abstract/concrete arenas, leaf pointers, links and metadata.
A parser or formatter error discards its returned CST. A printer error discards
its returned text. Providers remain responsible for their grammar, concrete
spelling and formatter preservation contracts; structural checks do not prove
that arbitrary leaf text means the stored abstract view.

For example, using the TypeScript marker for identity and its native
implementation for two independently selected capabilities:

```go
definition, err := representation.NewLanguage[typescript.Language](
    typescript.Language{}, "native-types-v1",
)
if err != nil { return err }
parser, err := parse.NewSyntaxParser[typescript.Language](definition, tsparse.Provider{})
if err != nil { return err }
printerProvider, err := sharedprint.New[typescript.Language](typescript.Language{})
if err != nil { return err }
printer, err := print.NewSyntaxPrinter[typescript.Language](definition, printerProvider)
if err != nil { return err }
syntax := engine.New()
if err := syntax.RegisterLanguage(definition); err != nil { return err }
if err := syntax.RegisterSyntaxParser(parser); err != nil { return err }
if err := syntax.RegisterSyntaxPrinter(printer); err != nil { return err }

parsed, err := syntax.ParseSyntax(ctx, definition.ID(), parse.Input{
    SourceID: "input", TreeID: "input", Text: "type Page<T> = T[];",
})
if err != nil { return err }
original, err := syntax.PrintSyntax(ctx, parsed)
if err != nil { return err }
_ = original // retained source bytes; Check and Render are not registered
```

Here tsparse denotes `parse/syntax/typescript/provider/go`, sharedprint denotes
`print/syntax/shared/provider/go`, and parse/print are the role contract packages.
Register `format.NewSyntaxFormatter[typescript.Language](definition, tsformat.Provider{})`
to enable Render and `check.NewChecker[typescript.Language](definition, tscheck.Provider{})`
to enable Check, semantic Capture and checked
transformation execution. tsformat and tscheck denote the corresponding
`format/syntax/typescript/provider/go` and `check/syntax/typescript/provider/go` packages. `semantic.Capture(ctx, checker, input)` also accepts an
explicit configured checker outside an engine. Capture preserves budget,
cancellation, catalogue and snapshot checks before semantic providers run.

After registering source/destination checkers, destination Syntax stages and a
mapping, capture the input and supply required receipts in transform.Input.
TransformTo returns an ASTSet; FormatSyntax/PrintSyntax can emit each selected
output tree. No canonical registration is installed implicitly.

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
throughout the frame's import closure. Independent single-capability providers exercise partial registration, per-stage
identity/profile rejection, render preflight and route checker preflight. All four
existing native implementations also run through the engine, including parse/print retention and formatting.

The maintained GitHub gate runs these tests, vet and a race-detector check for
concurrent registry operations, alongside existing native compiler and fuzz
checks. Documentation checks run on the same exact committed revision.
Workloads run remotely.

The [stage-first hierarchy](../LAYOUT.md) keeps structurally checked tree representations
independent of transformer contracts and the engine. Shared conformance tests
and discovered stage/dependency checks enforce this boundary. The moved APIs
have no aliases at the former package paths.

The [Horos frontend](horos.md) and [declarative native generation](generation-spi.md)
are implemented. Their registration is separate from ASTTransformer routing;
production reverse canonicalizers remain deferred. The [agnostic semantic model](semantics.md) and separate
[resolution SPI](resolution.md) and [analysis SPI](analysis.md) are implemented.
Capture admits a shared syntax snapshot; Resolve and Analyze return reports and
optional shared facts through independent explicit registrations.
Public builders, a generic Stage composition helper and an
application CLI also remain deferred; none is needed to register and execute a
correctly implemented ASTTransformer through this frame.

## Explicit semantic transformation input

`transform.Input` contains a captured `semantic.Snapshot` and an ordered slice of
`semantic.Result`. A transformer contract declares evidence requirements: exact
operation, profile and configuration, completed report verdicts, and whether an
available model is required. The typed adapter admits them using shared semantic
dependency checks before invoking its callback. Equal source text captured again
is not the same snapshot. Requirements are owned copies and metadata is bounded.

`transform.Request[L]` exposes owned source syntax, the entry address, captured
input envelopes and admitted evidence. The generic bridge interprets no provider
terms. Concrete maps own semantic applicability and located unsupported errors.
Cancellation before/after the callback returns no output.

Intermediate route outputs are syntax-checked and captured afresh; source
receipts are discarded. A later semantic edge cannot consume evidence about an
earlier tree. Execute the legs explicitly when intermediate analysis is needed.
See [semantic generation](semantic-generation.md) for the first production maps.

The [Files design](source.md) is implemented alongside the Syntax operations.
Both dimensions use context-bearing parse/format/print interfaces. Source retains
its exact capture and analysis receipt. ASTTransformer output is a complete
`representation.ASTSet`; typed callbacks return `ast.Set[T]`. Every route leg
admits and captures all output units and preserves the explicit entry. Receipt
requirements at subsequent legs are never satisfied by carrying old evidence.
`EntryAST()` selects the explicit entry for consumers that deliberately want one
unit; consumers emitting a source set should traverse all `Trees`.


## Declarative generation

Semantic generation has a separate registry from ASTTransformer route planning.
The application registers source interpretation, a native target and mapping
libraries; the kernel selects descriptions for the source capture and requested
profile. Implemented engine signatures are:

```go
func (e *Engine) PlanGeneration(ctx context.Context, r generation.Request) (g.GenerationPlan, error)
func (e *Engine) Generate(ctx context.Context, r generation.Request) (g.Output, error)
```

Here generation is transform/syntax/shared/generation/go and g is
contracts/generation/go. PlanGeneration resolves semantic choices without AST
construction. The returned plan's Run constructs fresh syntax. Generate performs
both and returns the same sealed output. Registration does not create language
route edges. [Exact contracts and consumer example](kernel-composition.md)
