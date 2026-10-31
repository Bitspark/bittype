# bittype

A language-independent type model and native type generation for TypeScript,
Go, Rust and C++.

**Current status, 31 October 2026:** generic AST data types, a typed pipeline
contract, and TypeScript, Go, Rust and C++ AST vocabularies are implemented in Go.
Each target supplies identity methods through a common language interface for
a future agnostic registry. The canonical AST and registry remain deferred.
There are no parsers, builders, grammar checkers, generators, formatters or
published type packages yet.

The intended pipeline is canonical input to canonical AST, target-language AST,
then source text. Generators and formatters are separate Go components.
TypeScript is an output language, alongside Go, Rust and C++. The present slice
provides target syntax declarations and language identity metadata without
transformation or formatting implementations.

Function, Stream, Cell and future behavioral contracts remain library-owned.
This milestone adds no adapters, callable implementations, serialization,
validation of application values or Wire dependencies.

## Read first

- [Proposed type model](docs/model.md): input format, semantics, target
  representations, limitations and implementation stages.
- [Generic AST](docs/ast.md): implemented data types, language separation,
  addressing, provenance and limits.
- [Target syntax](docs/syntax.md): the four native AST vocabularies, language
  identity interface, future registration and current limits.
- [Generation design](docs/generation.md): target layout, stage contract,
  composition strategy, dependencies and remaining work.
- [Charter](CHARTER.md): current ownership and delivery boundary.
- [Contributor instructions](AGENTS.md): repository workflow and prohibited
  coupling.
- [Source layout](LAYOUT.md): component-first implementation paths.

The active documents describe this milestone. Earlier bootstrap plans and
research remain in Git history; they are not a second current specification.
