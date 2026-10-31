# bittype

A language-independent type model and native type generation for TypeScript,
Go and Rust.

**Current status, 31 October 2026:** design proposal only. There is no parser,
generator or published type package yet. Start with
[the proposed model](docs/model.md).

The current milestone defines value types, generic parameters, imports and
applications, then generates native type declarations from that canonical input.
TypeScript is an output language, alongside Go and Rust.

Function, Stream, Cell and future behavioral contracts remain library-owned.
This milestone adds no adapters, callable implementations, serialization,
validation of application values or Wire dependencies.

## Read first

- [Proposed type model](docs/model.md): input format, semantics, target
  representations, limitations and implementation stages.
- [Charter](CHARTER.md): current ownership and delivery boundary.
- [Contributor instructions](AGENTS.md): repository workflow and prohibited
  coupling.
- [Source layout](LAYOUT.md): component-first paths for future implementations.

The active documents describe this milestone. Earlier bootstrap plans and
research remain in Git history; they are not a second current specification.
