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
- Native type generation for TypeScript, Go, Rust and C++.
- Native AST vocabularies and an agnostic language identity/registration
  boundary, independent of particular language providers.
- The specification of canonical declaration encoding and its independent
  vectors, once identity and normalization rules have been settled.
- Small data types needed to represent presence and nullness in generated code.

The generic AST, pipeline and language identity interface are independent
packages inside the initial root Go module. None imports another foundation,
concrete language, generator, formatter or contract library. CST depends on AST. Syntax packages
depend on the AST/CST carriers and identity interface. A future central language registry
will accept that interface without knowing concrete providers; the application
root selects providers and composes stages. The shared infrastructure does not
know the names or semantics of particular library contracts.

## Current delivery boundary

The [canonical model](docs/model.md) remains a proposal. Generic AST/CST carriers,
typed views and pipeline contracts, four native type vocabularies, native parsers,
AST-to-CST formatters, syntax checks and the shared CST printer are implemented.
The [native text design](docs/parsing-formatting.md) defines their ownership and
qualification limits. Specialized providers all extend LanguageProvider.

The canonical AST, AST-to-AST generators, registry, public builders and pipeline
executor remain deferred. Native parsing does not make a target language the
canonical authority. The [generation design](docs/generation.md) remains separate.

This scope excludes application-value validators, serialization, Wire adapters,
RPC/streaming implementations, consumer migration and package publication.
Function, Stream, Cell and future behavioral contracts remain ordinary library
declarations with no privileged cases in shared infrastructure. There are no
production dependencies on other Bitspark implementation repositories.

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
