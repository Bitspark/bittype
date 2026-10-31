# AST transformation and pipeline design

**Direction: 31 October 2026.** The toolchain is implemented in Go. Generic AST/CST
carriers, typed stages, LanguageProvider and four native text pipelines are
implemented. The [native text design](parsing-formatting.md) owns parser,
formatter, printer and checking details. Canonical AST, generators, registry,
public builders and application execution remain deferred. The owner's later
31 October proposal unifies generation and native-to-canonical migration under
ASTTransformer and adds a syntax engine with explicit registration and path
selection. The revised target layout below is a proposal, not an implementation.

The canonical model remains a [proposal](model.md). Its syntax and semantic
decisions must be settled before a language implementation claims conformance.

## Target pipeline

```text
canonical source
    -> parse, resolve and check once
    -> AST[Canonical]
         -> AST[TypeScript] -> CST[TypeScript] -> TypeScript source
         -> AST[Go]         -> CST[Go]         -> Go source
         -> AST[Rust]       -> CST[Rust]       -> Rust source
         -> AST[C++]        -> CST[C++]        -> C++ source
```

Canonical-to-native is one direction of AST transformation. Each such transformer
transforms a checked canonical AST into a target-language AST.
It chooses native type representations, names, imports, generic syntax and
recursive indirection according to an explicit target configuration.
It does not emit source strings.

Each formatter consumes a target AST view and constructs a CST. It owns escaping,
lexical spelling, layout and precedence. The generic printer emits supplied
leaves. The formatter does not resolve
canonical declarations or decide how a canonical nullable field should be
represented. Those choices must already be explicit in the target AST.

Canonical parsing, resolution and semantic checking belong to the future
canonical frontend. Native parsers are separate text-to-CST stages. File IO, dependency retrieval, configuration and output paths belong
to the application. A formatter is usable with a target AST built by another
producer, without a canonical input or generator dependency.

The output scope is native types and interfaces, including generic declarations
and applications. Executable bodies, value validators, serialization, callable
implementations, adapters and transport behavior are outside this milestone.

## Go layout

There is one root Go module, `github.com/Bitspark/bittype`. Components are
separate packages; separate modules are unnecessary for this initial foundation.
The final path component denotes the implementation language, not the target.
For example, `syntax/rust/go` is Go code describing Rust syntax.

Only the entries marked **present** exist in the current implementation.
Other paths are the target layout, not placeholder packages.

```text
go.mod                                      present
ast/go/
  ast.go, node.go, origin.go, view.go        present: carrier and borrowed view
  walk.go                                   planned: generic traversal
pipeline/go/
  stage.go                                  present: typed contract
  compose.go                                planned: typed composition
language/go/
  language.go                               present: identity interface
transform/go/
  ast.go                                    proposed: checked heterogeneous envelope
  transformer.go                            proposed: ASTTransformer contract
  adapter.go                                proposed: typed stage registration bridge
engine/go/
  engine.go                                 proposed: syntax engine API
  registry.go                               proposed: providers and directed edges
  path.go                                   proposed: route planning and ambiguity checks
  transform.go                              proposed: TransformTo execution
  text.go                                   proposed: registered parse/format/print composition
syntax/
  canonical/go/{nodes,concrete}.go           planned: canonical vocabulary
  typescript/go/nodes.go                     present
  golang/go/nodes.go                         present
  rust/go/nodes.go                           present
  cpp/go/nodes.go                            present
  {typescript,golang,rust,cpp}/go/
    concrete.go                             present: native concrete kinds
    builders.go                             planned
parse/canonical/go/parse.go                  planned: canonical text to CST
resolve/canonical/go/resolve.go              planned: references and generic binders
check/canonical/go/check.go                  planned: canonical semantics
format/canonical/go/format.go                planned: canonical AST to CST
provider/canonical/go/provider.go            planned: explicit capability composition
transform/
  canonical-to-typescript/go/                planned: transform, types, names, imports
  canonical-to-golang/go/                    planned
  canonical-to-rust/go/                      planned
  canonical-to-cpp/go/                       planned
  typescript-to-canonical/go/                proposed: partial migration mapping
  golang-to-canonical/go/                    proposed
  rust-to-canonical/go/                      proposed
  cpp-to-canonical/go/                       proposed
cst/go/cst.go                              present
text/go/{input,diagnostic,provider}.go       present
check/go/, check/cst/go/                     present
check/{typescript,golang,rust,cpp}/go/        present
parse/{typescript,golang,rust,cpp}/go/        present
format/{typescript,golang,rust,cpp}/go/       present: AST view to CST
print/go/print.go                           present: CST to text
provider/{typescript,golang,rust,cpp}/go/     present
cmd/bittype/go/
  main.go, compose.go, files.go              planned
integration/go/
  contracts_test.go, dependencies_test.go,
  syntax_test.go                            present
  *_text_test.go, formatter_fixtures_test.go present: native text/AST qualification
  provider_contracts_test.go                 present
  transform_test.go, engine_test.go          proposed
  testdata/*.go                             present: API type-check fixtures
  testdata/syntax/*.go                       present: target/extension fixtures
testdata/
  text/{typescript,golang,rust,cpp}/valid/   present: native source fixtures
  canonical/, transforms/                   planned: semantic/adoption fixtures
docs/
  model.md, ast.md, syntax.md, generation.md present
  parsing-formatting.md                     present
```

The [target vocabularies](syntax.md) cover a native type-declaration subset.
They are not full programming-language parsers. A concrete unsupported type
case should identify a missing grammatical form; it must not lead to a raw
code-string escape hatch. Their supported shapes have separate native checks. Public builders remain
deferred.

## Language identity and registration

Each target package supplies a zero-value `Language` implementing
`language.LanguageProvider`, with stable `ID() language.ID` and `Name() string`
methods. It also supplies aliases over `ast.AST[Language]`,
`ast.Node[Language]` and `ast.Kind[Language]`. The carrier now exposes ASTView; the stage API
remains unchanged and does not import the language identity interface.

The owner clarified on 31 October that a central agnostic language AST registry
will exist. The earlier no-language-registry direction is superseded. Its
implementation is deferred; the present interface is the common contract that
registered providers will implement. The registry must accept providers
through that interface, without a concrete language import, switch or builtin
list. Application composition owns the selected provider instances.

Identity methods do not discover packages, initialize resources or register
anything. A shared ID is not proof of equal Go marker types, syntax versions or
grammar compatibility. The identity interface neither erases `AST[L]` nor
implements construction, checking or transformations. See [the syntax
contract](syntax.md) for the IDs, provider contract and extension example.

## Typed stage contract

The [pipeline package](../pipeline/go/stage.go) defines:

```go
type Stage[I, O any] func(I) (O, error)
```

A generator will satisfy a stage from one AST type to another; a formatter
satisfies a stage from a target AST view to CST, and a printer from CST to text. The planned relationships
are shown below; target markers exist, but the Canonical marker does not yet:

```text
Stage[AST[Canonical], AST[TypeScript]]
Stage[ast.View[TypeScript], CST[TypeScript]]
Stage[CST[TypeScript], string]
```

`Stage` imports nothing and does not know that ASTs or languages exist.
The function type remains sufficient for statically composed producers and
consumers. Dynamic registration adds the separate ASTTransformer boundary below;
it does not replace Stage. Concrete stage factories can close over explicit
mapping and target configuration.

A stage treats input and reachable storage as read-only. It owns mutable
output storage it creates. On error, callers must discard the output.
The zero/nil stage means no configured stage and must not be invoked.
These behavioral requirements are documented contracts, not enforcement
provided by the function type.

The following signature describes future composition, **not a present API**:

```go
func Then[A, B, C any](first Stage[A, B], next Stage[B, C]) Stage[A, C]
```

It would connect stages only when the intermediate type agrees and propagate
the first failure. There is no `Then` implementation or stub in this slice.
Scheduling, caching, cancellation and diagnostic aggregation are also deferred.

## Proposed ASTTransformer and syntax engine

Generation and canonicalization are names for directions of one operation:

```text
AST[Canonical] -> AST[TypeScript]    generation
AST[TypeScript] -> AST[Canonical]   canonicalization/adoption
AST[TypeScript] -> AST[Canonical] -> AST[Rust]   composed transformation
```

The proposed registry-facing interface is:

```go
type ASTTransformer interface {
    FromID() language.ID
    ToID() language.ID
    Transform(AST) (AST, error)
}

func (e *Engine) TransformTo(input AST, targetID language.ID) (AST, error)
```

These are design sketches, not implemented APIs. AST here is the planned
transform.AST envelope: a language identity plus an opaque, checked native AST
payload. It is not ast.AST[any]; different generic instantiations are not
interchangeable. Typed adapters retain Stage[ast.AST[From], ast.AST[To]] inside
each implementation and validate the actual payload type and endpoint identities
at the dynamic boundary. No generic AST kind, language switch or reflective
syntax conversion is needed. CST input can supply its existing AST view when
wrapped; transformation does not preserve concrete source spelling.
This registration adapter only bridges typed AST calls; it generates no
application, service or runtime bindings.

LanguageProvider still identifies a single language. A transformer relates two
languages, so it identifies its endpoints with FromID/ToID rather than pretending
to have one language ID. The engine registers language capabilities and directed
transformer edges explicitly. Application composition supplies every instance.
Canonical is an ordinary language provider; the engine has no built-in canonical
name or route. Adding a language or transformation changes its implementation
and application registration, not engine source.

Initial route policy:

- Reject duplicate language identities and duplicate directed edges. Alternate
  mapping profiles use separately composed engine instances initially.
- Plan a unique shortest directed route, without cycles; expose the selected
  route through a Plan operation. Missing paths and tied shortest routes are
  errors. Allow callers to provide an explicit route when desired.
- FromID/ToID establish candidate connectivity, not applicability to every AST.
  Each edge validates its admitted input, checks its output, and returns located
  diagnostics for unsupported syntax or missing semantic information.
- Execute the selected route once. Do not silently retry another route after a
  semantic rejection or substitute a weaker mapping. Return no usable partial
  result on failure. Input remains read-only and results own their storage.
- A same-language request validates and returns a deep snapshot. It does not
  implicitly normalize or migrate a schema revision.

Language IDs remain stable family keys. Registered providers and configured
transformer endpoints must agree on the supported syntax revision/profile;
matching strings alone is insufficient. An engine instance fixes that profile
per language. Version/profile negotiation is not hidden in graph traversal.

Canonicalization is a partial semantic mapping, not the inverse of formatting
or a promise to recover an original definition. For example, a native numeric
type alone may not identify the intended canonical numeric domain; a nullable
pointer alone may not distinguish optional presence from an explicit null value.
The mapper must obtain that choice from configured migration rules or report
that it cannot decide. It must not silently select a narrower domain. Ordinary
successful canonical-to-native mappings also retain their documented native
representation limits; the route planner does not prove semantic equivalence.

The adoption flow is native text -> native CST/AST view -> checked canonical AST
-> canonical CST -> canonical source. Reviewing and adopting that source makes
it authoritative for subsequent generation. Native source does not become a
permanent second source of truth. Generic binders, declaration references and
origin links survive the mapping; comments and formatting require separate
metadata if a migration needs them.

Implement and qualify canonical-to-native mappings first. Add reverse mappings
for explicitly supported migration cases afterwards, using the same interface.
Path tests must cover an unknown third-party language, duplicate registration,
missing/ambiguous routes, typed endpoint mismatches, unsupported input, metadata
and input ownership. Semantic fixtures must establish each edge's meaning; a
path that merely reaches the right language ID is not sufficient evidence.

## Composition and dependencies

The application root will explicitly compose selected concrete stages.
It resolves input once and passes the resulting canonical AST to independently
selected target generators. Each formatter receives the corresponding target
AST. Shared inputs remain read-only by contract; outputs have separate storage.
Whether targets run sequentially or concurrently is an application decision,
not a requirement imposed by the AST model.

Arrows in the table mean a direct package dependency is allowed:

| Component | Dependencies |
| --- | --- |
| Generic AST | Standard library only; currently no imports |
| Generic pipeline | Standard library only; currently no imports, including no AST dependency |
| Language identity interface | Standard library only; currently no imports |
| Proposed syntax engine and registry | Language/provider and transformation contracts; no concrete syntax imports |
| Proposed transformation boundary | Generic AST, language identity, typed stage contracts; no concrete syntax imports |
| Generic CST | Generic AST |
| Each syntax package | Generic AST/CST, language identity interface |
| Canonical frontend | Canonical syntax, generic AST |
| Each transformer | Its source and destination syntax/checkers, generic AST; pipeline and transformation contracts |
| Each formatter | Own syntax/checker, AST/CST carriers, private concrete construction |
| Shared printer | Generic CST/checks only; no concrete language imports |
| Application composition root | Explicitly selected language providers, engine, frontend, transformers and pipeline |
| Integration tests | The APIs under test and native compiler/checker tooling |

Standard-library utilities are permitted in every layer. There are no external
Go module dependencies in this slice and no dependencies on another Bitspark
implementation repository. The AST, pipeline and language identity packages do
not import each other or concrete syntax packages. Syntax packages do not import
one another, generators or formatters. A formatter does not import the canonical
grammar or its generator.

Go's parser uses the standard go/parser representation internally. All four
formatters construct CST directly in Go; no source-print-and-reparse shortcut
is used. Native compilers are remote test dependencies, not production dependencies.

The intended language registry is distinct from a behavioral contract catalogue.
No registry implementation, automatic discovery or hidden default registration
is introduced now. The generic AST, pipeline, language identity interface and
future registry must never acquire Function, Stream, Cell, Wire or runtime
imports, special kinds, semantic name switches or structural fingerprint matching.
Native function-signature syntax remains language-owned. Later library-owned
realizations require a separate scoped design; this slice supplies no binding
extension registry.

## Evidence and next increments

The present Go module targets Go 1.27.0 language features; CI uses the pinned
Go 1.27.3 toolchain. The [Go workflow](../.github/workflows/go.yml) checks the
exact committed revision on a GitHub runner: formatting, tests and vet.
[Documentation checks](../.github/workflows/docs.yml) also run remotely.
No compilation or formatting workload is run on the workstation.

Public API fixtures establish that an independent consumer can supply a
language identity implementation without changing the foundation, and that
incompatible language/stage types are rejected. Four authored target examples
include open/applied pagination and native reader signatures. Tests exercise
metadata through the common interface and enforce the production dependency
boundary. Native text syntax and compiler fixtures now provide separately scoped evidence.
Registry operation and canonical generation remain unimplemented.

Subsequent, separately authorized increments are:

1. Specify and implement the canonical AST and independent semantic fixtures.
2. Implement the typed transformation boundary and agnostic syntax engine registry.
3. Implement canonical frontend and canonical-to-native ASTTransformer providers.
4. Compose and qualify native generation, then add explicitly scoped reverse
   transformers for migration and qualify multi-step routes.

Native text stages and compiler fixtures are implemented separately; they do not
depend on these future increments.

Generator tests will compare target ASTs. Formatter tests will accept
independently authored target ASTs. End-to-end tests will compile and use the
same semantic examples, including open/applied pagination and presence/null
distinctions. Compiler success alone will not prove all canonical value
invariants.
