# bittype

A language-independent type model and native type generation for TypeScript,
Go, Rust and C++.

**Current status, 31 October 2026:** generic AST/CST carriers, typed pipeline
contracts, native syntax vocabularies, parsers, formatters, checks and a shared
printer are implemented in Go. Specialized providers extend the common
LanguageProvider identity interface. The canonical model remains a proposal;
the canonical AST, generators, registry and package releases remain deferred.

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
- [Generation design](docs/generation.md): planned AST-to-AST pipeline.
- [Charter](CHARTER.md), [contributor instructions](AGENTS.md) and
  [source layout](LAYOUT.md): ownership, delivery and package boundaries.

Historical research remains in Git history, not as competing active specifications.
