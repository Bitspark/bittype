# bittype

A language-independent type model and native type generation for TypeScript,
Go and Rust.

**Current status, 31 October 2026:** generic AST data types and a typed pipeline
contract are implemented in Go. The canonical language is still a proposal.
There are no language grammars, parsers, generators, formatters or published
type packages yet.

The intended pipeline is canonical input to canonical AST, target-language AST,
then source text. Generators and formatters are separate Go components.
TypeScript is an output language, alongside Go and Rust. The present slice
provides only the shared AST and stage declarations, without transformation or
formatting implementations.

Function, Stream, Cell and future behavioral contracts remain library-owned.
This milestone adds no adapters, callable implementations, serialization,
validation of application values or Wire dependencies.

## Read first

- [Proposed type model](docs/model.md): input format, semantics, target
  representations, limitations and implementation stages.
- [Generic AST](docs/ast.md): implemented data types, language separation,
  addressing, provenance and limits.
- [Generation design](docs/generation.md): target layout, stage contract,
  composition strategy, dependencies and remaining work.
- [Charter](CHARTER.md): current ownership and delivery boundary.
- [Contributor instructions](AGENTS.md): repository workflow and prohibited
  coupling.
- [Source layout](LAYOUT.md): component-first implementation paths.

The active documents describe this milestone. Earlier bootstrap plans and
research remain in Git history; they are not a second current specification.
