# bittype

Syntax infrastructure and native type generation for TypeScript, Go, Rust and
C++, using horos as the canonical type model.

**Current status, 6 November 2026:** generic AST/CST carriers, typed pipeline
contracts, native syntax vocabularies, parsers, formatters, checks and a shared
printer are implemented in Go. The syntax engine registers languages and
ASTTransformer implementations, plans routes and composes existing text stages.
Specialized providers extend language.Provider. The engine registers parsing,
checking, formatting and printing independently; identity registration installs
no capability. See the [engine composition API](docs/engine.md#registration). Horos's Go declaration model now
supplies module resolution, generic binding, general structural kinds and specialization. Bittype
implements its [text and semantic frontend](docs/horos.md), resolver and semantic analyzer.
The engine captures syntax once and invokes each capability explicitly. Horos
analysis reuses the resolver's [shared semantic graph](docs/semantics.md), with
separate reports for validity, coverage and model availability.
The [readable Horos syntax](docs/horos-syntax.md) and its text stages are implemented.
The [Horos module command](tools/horos-module/README.md) lowers an ordinary source
module to canonical Ontos bytes for distribution, using the same parser and
module encoding as the frontend, without semantic resolution.
The first [semantic pagination mappings](docs/semantic-generation.md) generate
TypeScript, Go, Rust and C++ type ASTs from explicitly supplied Horos analysis.
The [Syntax and Files design](docs/source.md) specifies the symmetry
`Syntax : CST : AST` and `Files : Source : Semantics`. Typed carriers, symmetric parse/format/print interfaces, multi-unit transformations
and the Horos Files composition are implemented. Bounded Go, TypeScript and Rust Files
resolution/analysis providers are implemented; C++ semantics remain deferred.
Implementations use explicit Syntax/Files stage paths and role-named Go packages.
Each language composition package provides matching `NewSyntax...` and
`NewFiles...` constructors for its existing capabilities, taking an explicit
language/profile descriptor. See [naming, placement and construction](LAYOUT.md).
Package releases remain deferred.
The concrete Horos frontend uses the Horos Go module pinned in go.mod and Ontos 0.16.0. Consumers
must upgrade together; no older dependency path or compatibility shim remains.
See the [dependency alignment](docs/dependencies.md).
The engine starts empty and has no built-in languages or mappings.

The [joint verification composition](docs/verification.md#implemented-go-surface)
adds immutable source/target association and an optional Logos-backed Horos/Go
integration. Explicit Horos kind queries and registered bounded Go `go/types` semantic providers
supply TextBox counterexamples and conditional Page representation proofs,
including nested reuse. The owning composition imports the public Logos contract;
all packages use the root module and its anonymous installed-consumer gate.
Native admission enforcement and runtime bindings are separate.

The first [public generation SPI](docs/generation-spi.md) is implemented in Go
source: independent Use, Publish and Define capabilities, paired Introduce/Apply,
and a guarded invocation with explicitly supplied collaborators. The runner lives
with its contracts in `contracts/generation/go`. [Automatic composition](docs/generation-composition.md)
selects exact mappings and guarded overrides, publishes roots and schedules
definitions. All four generators use this SPI. Horos-owned defaults live in its
optional `tooling/bittype` module; the former private Dialects are deleted.

The [stage-first layout](LAYOUT.md) groups shared tree representations under
trees and public interfaces under contracts. Every native language implements
the same syntax/parse/check/format/provider matrix, enforced by architecture and
shared conformance tests. The engine contains no built-in language catalogue.
Horos implements the same text stages plus separately registered resolution and analysis.

The native pipeline separates source parsing, concrete layout and printing:

    text -> Parser -> CST[L]
    AST view[L] -> Formatter -> CST[L] -> Printer -> text

AST and CST expose the same borrowed AST-view interface. Parsers retain source
tokens and trivia; formatters construct a new concrete layout directly from
abstract syntax. The generic printer emits supplied leaves without knowing any
native grammar. Native pagination generators consume analyzed canonical definitions and produce
target ASTs, followed by this formatting/printing pipeline.

Function, Stream, Cell and future behavioral contracts remain library-owned.
This milestone adds no runtime adapters, callable implementations, serialization,
application-value validation or dependencies on bitwire/bitruntime.

Horos has one [declaration language](docs/model.md) with [general kinds](docs/kinds.md).
`Type`, `Data`, `Scalar`, `List<K>` and `Record<K>` are ordinary prelude definitions.
Horos 0.10.2 authors those definitions in `.horos` source and packages the canonical
module produced by the Bittype encoder; the Go-authored definition table is removed.
`Page<T>` accepts data or capability types; `DataPage<D: h.Data> = Page<D>`
restricts arguments without a second language or copied template. The old
Horos/Wire frontend and profiles are removed without aliases. Ground Ontos values,
Horos ground checks and `any` remain data-only. Ontos/wire owns mixed values;
existing hydration remains in Bitwire/Bitruntime. The agnostic kernel does not
interpret these concrete type semantics.

## Read first

- [Horos canonical input](docs/model.md): ownership, declarations and native representation plan.
- [General kinds](docs/kinds.md): structural families, ordinary Data, constructor
  contracts and the semantic/native generation boundary.
- [Type-system modeling](docs/type-system-model.md): proposed terms, judgments,
  binding and rules for Horos and native languages, their relationship to Logos,
  and concrete Horos/Go reasoning examples; documentation only.
- [Type-system reasoning architecture](docs/type-system-architecture.md): expert-reviewed
  ownership, stage-first package roles, exact TextBox/Page statements, Logos
  checking prerequisites, implementation boundary and module isolation.
- [Logos checking and the verification boundary](docs/logos-verification.md):
  source-traced production request, direct checking APIs, typed proposal/experiment,
  Bittype-owned association work, scoped integration and remaining dependencies.
- [Joint verification contract](docs/verification.md): implemented immutable
  artifact association, preparation/checking split, result/failure semantics and
  TextBox/Page walkthroughs, with the supported integration subset and limits.
- [Preserving native kind constraints](docs/native-kind-constraints.md): proposed
  direct bounds and source witnesses for all four targets, guarantees and
  unresolved limits; documentation only.
- [Logos preservation proofs](docs/logos-preservation.md): proposed source/native
  symbolic pagination proof with concrete scoped facts and structural rules;
  native admission is separate. Documentation only.
- [Go native admission](docs/go-native-admission.md): proposed compositional
  certificate encoding and independent compiler tests for pagination, nesting,
  aliases, inference and attempts to forge evidence; no production implementation.
- [Type and value domains](docs/hydrated-types.md): ownership, authoritative canonical
  input and explicit data boundary. The [frontend guide](docs/horos.md)
  documents implemented APIs and qualification.
- [Concept preservation](docs/wire-extension-preservation.md): retained ideas,
  superseded proposals and expert-review evidence for the ownership change.
- [Native representations](docs/type-representations.md):
  direct generic pagination and library-owned type recipes configured for
  TypeScript, Go, Rust and C++; exact library mappings and pagination compositions
  use the public SPI.
- [Public generation SPI](docs/generation-spi.md): implemented operations,
  association ownership, scope and callback lifetimes, and explicit invocation.
- [Type-owned generator plugins](docs/type-generation-plugins.md): implemented
  public contracts, automatic selectors, guarded overrides and the four migrated
  production mappings; native constraint enforcement remains proposed.
- [Generic generation walkthrough](docs/generic-elaboration.md#pagination-and-scope):
  pagination across four targets, paired parameter/application conventions,
  nested results, Cell composition and the remaining implementation stages.
- [Runtime values and hydration](docs/runtime-values.md): type descriptions versus
  native instances, existing Wire/runtime carriers and future typed/dynamic
  adapters; documentation only, with pinned implementation evidence.
- [Horos AST](docs/horos.md): model projection, binding references, checking and registration.
- [Semantic facts](docs/semantics.md): immutable entities, terms, occurrences and dependencies.
- [Resolution capability](docs/resolution.md): shared snapshots, explicit prerequisites
  and immutable model reuse.
- [Provider matrix](docs/providers.md): implemented capabilities by language.
- [Analysis capability](docs/analysis.md): request capture, explicit registration,
  completion/verdict reports and publication requirements.
- [Generic AST](docs/ast.md): carriers, typed views, addressing and provenance.
- [Target syntax](docs/syntax.md): native vocabulary and supported type subsets.
- [Parsing and formatting](docs/parsing-formatting.md): CST model, providers,
  ownership, round trips and qualification.
- [Syntax engine](docs/engine.md): registration, typed boundaries, path finding,
  execution and text-stage composition.
- [Generation design](docs/generation.md): AST-to-AST pipeline and future mappings.
- [Syntax and Files](docs/source.md): Source/Semantics carriers, symmetric
  parse/format/print roles, multi-unit transformations and source-base parsing.
- [Charter](CHARTER.md), [contributor instructions](AGENTS.md) and
  [source layout](LAYOUT.md): ownership, delivery and package boundaries.

Historical research remains in Git history, not as competing active specifications.

See [native semantic providers](docs/native-semantics.md)
for the bounded Go, TypeScript and Rust resolvers/analyzers, immutable native facts and
explicit compiler queries. C++ semantic checking and broader Rust coverage remain planned.
