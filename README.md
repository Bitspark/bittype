# bittype

Syntax infrastructure and native type generation for TypeScript, Go, Rust and
C++, using horos as the canonical data type model.

**Current status, 2 November 2026:** generic AST/CST carriers, typed pipeline
contracts, native syntax vocabularies, parsers, formatters, checks and a shared
printer are implemented in Go. The syntax engine registers languages and
ASTTransformer implementations, plans routes and composes existing text stages.
Specialized providers extend LanguageProvider. Horos's Go declaration model now
supplies module resolution, generic binding and ground specialization. Bittype
implements its [AST-only frontend](docs/horos.md) and semantic analyzer, preserving
those bindings in the [shared semantic graph](docs/semantics.md). The engine's
optional analysis capability reports validity, coverage and model availability.
Readable Horos syntax, native mappings and package releases remain deferred.
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
application-value validation or Wire dependencies.

Ontos and Horos stay data-only. The [contract-type proposal](docs/hydrated-types.md)
places Wire and mixed data/capability types in a separate Bittype language
provider above Horos. Existing bitwire/bitruntime hydration remains above Horos;
the agnostic kernel acquires no Wire semantics. That provider is not implemented.

## Read first

- [Horos canonical input](docs/model.md): ownership, declarations and native representation plan.
- [Contract types above Horos](docs/hydrated-types.md): exact proposed Wire input,
  explicit data boundary, library-owned Function/Cell/Stream composition and future provider;
  documentation only.
- [Horos AST](docs/horos.md): model projection, binding references, checking and registration.
- [Semantic facts](docs/semantics.md): immutable entities, terms, occurrences and dependencies.
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
