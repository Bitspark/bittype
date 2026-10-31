# bittype

A language-independent type model and native type generation for TypeScript,
Go, Rust and C++.

**Current status, 31 October 2026:** generic AST/CST carriers, typed pipeline
contracts, native syntax vocabularies, parsers, formatters, checks and a shared
printer are implemented in Go. The syntax engine registers languages and
ASTTransformer implementations, plans routes and composes existing text stages.
Specialized text providers extend LanguageProvider. The canonical model remains
a proposal; canonical AST, actual language mappings and package releases remain
deferred. The engine starts empty and has no built-in languages or mappings.

The [stage-first layout](LAYOUT.md) groups shared tree representations under
trees and public interfaces under contracts. Every native language implements
the same syntax/parse/check/format/provider matrix, enforced by architecture and
shared conformance tests. The engine contains no built-in language catalogue.

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

## Read first

- [Proposed type model](docs/model.md): canonical input proposal and semantics.
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
