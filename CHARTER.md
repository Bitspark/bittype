# Charter

**Current direction: 31 October 2026.** The owner selected bittype for canonical
type definitions and native type generation, starting with types only.
The detailed [model](docs/model.md) is a proposal, not a published language
revision or an implementation.

## Ownership

bittype owns:

- The language-independent value-type model, its source syntax, declarations,
  generic slots, applications, imports and resolution rules.
- Parsing and checking declarations, including scope and recursive type rules.
- Native type generation for TypeScript, Go and Rust.
- The specification of canonical declaration encoding and its independent
  vectors, once identity and normalization rules have been settled.
- Small data types needed to represent presence and nullness in generated code.

The generic AST and pipeline are independent packages inside the initial root
Go module. Neither imports a language grammar, generator, formatter or contract
library; the pipeline also does not import the AST package. The application
root will compose explicit stages and targets. The shared infrastructure does
not know the names or semantics of particular library contracts.

## Current delivery boundary

The proposed model is documented in docs/model.md. The present implementation
contains only language-parameterized AST data declarations and a typed pipeline
stage contract in Go. No language grammar, parser, generator, formatter,
composition executor or command is implemented.

The [generation design](docs/generation.md) specifies the target file layout,
composition and dependencies. Later stages will implement the canonical
frontend, three AST-to-AST type generators and separate target formatters.
All those components will be implemented in Go; emitted types target
TypeScript, Go and Rust.

This scope does not include application-value validators, descriptors for a
particular validator, serialization, Wire adapters, RPC, streaming, callable
implementations, consumer migration, package publication or deployment.

Function, Stream, Cell and future behavioral contracts belong to ordinary
libraries. Neither their protocols nor their language-specific realizations
become privileged cases in the general type model, generator or runtime.
Preserve resolved declaration identities and applications so later consumers
can select library-owned realizations without matching structural lookalikes.

bittype has no dependency on another Bitspark implementation repository.
References to theory or previous designs supply evidence, not runtime
dependencies. A declaration has the same meaning with a different storage,
serialization or communication mechanism.

## Versioning and evidence

A published language revision fixes its admitted syntax and meaning.
The proposed `bittype: 1` spelling does not itself publish revision 1.

Canonical declaration identity, normalized value shape, generated-code
compatibility and undeclared behavior remain distinct. A digest does not prove
behavior or authorize a binding substitution. Specify normalization and identity
content before assigning a canonical encoding version or publishing digests.

A generation manifest records the input revision, target configuration and
generator version. Once an identity encoding exists, it records that version
too; no placeholder digest stands in for unresolved declarations.

Independent semantic examples must cover generic scope/application, recursion,
absence and null, numeric domains and invalid input. Each language's output
must compile against the same examples. Document native representation limits
rather than claiming that a compiler enforces every canonical value invariant.

## Historical scope

The September bootstrap described a broader generated-consumer chain.
The owner has narrowed the current work to the value model and types-only
generation. Historical commits and research remain evidence; they do not add
callables, adapters or migration gates to this milestone.
