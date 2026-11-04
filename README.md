# bittype

Syntax infrastructure and native type generation for TypeScript, Go, Rust and
C++, using horos as the canonical type model.

**Current status, 4 November 2026:** generic AST/CST carriers, typed pipeline
contracts, native syntax vocabularies, parsers, formatters, checks and a shared
printer are implemented in Go. The syntax engine registers languages and
ASTTransformer implementations, plans routes and composes existing text stages.
Specialized providers extend language.Provider. The engine registers parsing,
checking, formatting and printing independently; identity registration installs
no capability. See the [engine composition API](docs/engine.md#registration). Horos's Go declaration model now
supplies module resolution, generic binding, Type/Data bounds and specialization. Bittype
implements its [AST-only frontend](docs/horos.md), resolver and semantic analyzer.
The engine captures syntax once and invokes each capability explicitly. Horos
analysis reuses the resolver's [shared semantic graph](docs/semantics.md), with
separate reports for validity, coverage and model availability.
Readable Horos syntax, native mappings and package releases remain deferred.
The concrete Horos frontend uses public Horos 0.9.0 and Ontos 0.16.0. Consumers
must upgrade together; no older dependency path or compatibility shim remains.
See the [dependency alignment and consumer checks](docs/dependencies.md).
The engine starts empty and has no built-in languages or mappings.

The [stage-first layout](LAYOUT.md) groups shared tree representations under
trees and public interfaces under contracts. Every native language implements
the same syntax/parse/check/format/provider matrix, enforced by architecture and
shared conformance tests. The engine contains no built-in language catalogue.
Horos has a separately qualified syntax/resolve/check/analyze/provider capability profile;
it does not advertise text stages that are not implemented.

The native pipeline separates source parsing, concrete layout and printing:

    text -> Parser -> CST[L]
    AST view[L] -> Formatter -> CST[L] -> Printer -> text

AST and CST expose the same borrowed AST-view interface. Parsers retain source
tokens and trivia; formatters construct a new concrete layout directly from
abstract syntax. The generic printer emits supplied leaves without knowing any
native grammar. Future generators will consume canonical definitions and produce
target ASTs, followed by this formatting/printing pipeline.

Function, Stream, Cell and future behavioral contracts remain library-owned.
This milestone adds no runtime adapters, callable implementations, serialization,
application-value validation or dependencies on bitwire/bitruntime.

Horos has one [declaration language](docs/model.md) with explicit Type/Data
bounds. `Page<T>` accepts data or capability types; `DataPage<D: Data> = Page<D>`
restricts arguments without a second language or copied template. The old
Horos/Wire frontend and profiles are removed without aliases. Ground Ontos values,
Horos ground checks and `any` remain data-only. Ontos/wire owns mixed values;
existing hydration remains in Bitwire/Bitruntime. The agnostic kernel does not
interpret these concrete type semantics.

## Read first

- [Horos canonical input](docs/model.md): ownership, declarations and native representation plan.
- [Type and value domains](docs/hydrated-types.md): ownership, authoritative canonical
  input and explicit data boundary. The [frontend guide](docs/horos.md)
  documents implemented APIs and qualification.
- [Concept preservation](docs/wire-extension-preservation.md): retained ideas,
  superseded proposals and expert-review evidence for the ownership change.
- [Native representations](docs/type-representations.md):
  direct generic pagination and library-owned type recipes configured for
  TypeScript, Go, Rust and C++; native
  recipes and generators remain proposed.
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
- [Charter](CHARTER.md), [contributor instructions](AGENTS.md) and
  [source layout](LAYOUT.md): ownership, delivery and package boundaries.

Historical research remains in Git history, not as competing active specifications.
