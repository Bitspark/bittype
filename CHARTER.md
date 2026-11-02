# Charter

**Current direction: 2 November 2026.** The owner selected horos for canonical
data type semantics and bittype for syntax/tooling and native type generation,
starting with types only. The [model boundary](docs/model.md) replaces the earlier
independent bittype type-language proposal. The declaration format remains a
proposal, not a published language revision or an implementation.

## Ownership

bittype owns:

- Readable source syntax and AST/CST projections of horos declarations.
- Frontend parsing, resolution tooling and located diagnostics implementing
  horos's declaration/binding rules, without a second semantic specification.
- Native type generation for TypeScript, Go, Rust and C++.
- Native AST vocabularies and an agnostic language identity/registration
  boundary, independent of particular language providers.
- Native representation profiles and independent generation/round-trip fixtures.
- Small data types needed to represent presence and nullness in generated code.

Horos owns descriptors, conformance and the proposed module/generic/application
semantics, declaration encoding and semantic vectors. Ontos owns values,
structural identity and encoding. Horos does not depend on bittype; only the
Horos-specific frontend and relevant transformers in bittype will consume its
model. Generic foundations and native text stages remain independent.

The owner's 31 October structure decision selects the stage-first matrix in
[LAYOUT.md](LAYOUT.md). Shared trees, shared contracts and the engine have distinct
homes; all native languages follow the same stage structure. Conformance and
dependency checks enforce the uniformity without production language discovery.

The generic AST, pipeline and language identity interface are independent
packages inside the initial root Go module. None imports another foundation,
concrete language, generator, formatter or contract library. CST depends on AST. Syntax packages
depend on the AST/CST carriers and identity interface. The central syntax engine
accepts explicit registration without knowing concrete providers; the application
root selects providers and composes stages. The shared infrastructure does not
know the names or semantics of particular library contracts.

## Current delivery boundary

The [Horos declaration integration](docs/model.md) remains a proposal. Generic AST/CST carriers,
typed views and pipeline contracts, four native type vocabularies, native parsers,
AST-to-CST formatters, syntax checks and the shared CST printer are implemented.
The [native text design](docs/parsing-formatting.md) defines their ownership and
qualification limits. Specialized providers all extend LanguageProvider.

The ASTTransformer boundary and syntax engine now provide checked heterogeneous
registration, route planning/execution and native text composition. The canonical
AST, actual AST-to-AST mappings, public builders and application CLI remain
deferred. Native parsing does not make a target language the
canonical authority. The [generation design](docs/generation.md) remains separate.

This scope excludes application-value validators, serialization, Wire adapters,
RPC/streaming implementations, consumer migration and package publication.
Function, Stream, Cell and future behavioral contracts remain ordinary library
declarations with no privileged cases in shared infrastructure. There are currently
no production dependencies on other Bitspark implementation repositories. Future
Horos integration must narrowly extend the relevant dependency checks without
adding model imports to generic carriers, contracts, engine or printer.

## Versioning and evidence

A published language revision fixes its admitted syntax and meaning.
The earlier `bittype: 1` proposal is superseded. Horos's proposed declaration
revision does not make those documents accepted by its released checker.

Canonical declaration identity, normalized value shape, generated-code
compatibility and undeclared behavior remain distinct. A digest does not prove
behavior or authorize a binding substitution. Specify normalization and identity
content before implementing an encoding or publishing digests. Horos's proposal
specifies document pins separately from semantic equivalence; bittype consumes
that specification rather than inventing one.

A generation manifest records the input revision, target configuration and
generator version. Once an identity encoding exists, it records that version
too; no placeholder digest stands in for unresolved declarations.

Independent semantic examples must cover generic scope/application, recursion,
absence and null, numeric domains and invalid input. Each language's output
must compile against the same examples. Document native representation limits
rather than claiming that a compiler enforces every canonical value invariant.

## Historical scope

The September bootstrap described a broader generated-consumer chain.
The owner has narrowed the current work to Horos-owned definitions and types-only
generation. This update is documentation only. Historical evidence does not add
callables, adapters or migration gates to this milestone.
