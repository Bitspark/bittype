# Charter

**Current direction: 3 November 2026.** The owner selected horos for canonical
data type semantics and bittype for syntax/tooling and native type generation,
starting with types only. The [model boundary](docs/model.md) replaces the earlier
independent bittype type-language proposal. Horos declaration revision 1 and its
bittype AST projection are implemented in Go source; no package release is implied.

## Ownership

bittype owns:

- Readable source syntax and AST/CST projections of horos declarations.
- Frontend parsing, resolution tooling and located diagnostics implementing
  horos's declaration/binding rules, without a second semantic specification.
- Native type generation for TypeScript, Go, Rust and C++.
- Native AST vocabularies and an agnostic language identity/registration
  boundary, independent of particular language providers.
- Native representation profiles and independent generation/round-trip fixtures.
- The [agnostic semantic model](docs/semantics.md) and implemented
  [analysis-provider boundary](docs/analysis.md); language providers retain their
  semantic rules. The separate [resolver contract](docs/resolution.md) uses the
  same semantic infrastructure. Horos supplies the first production resolver and
  analyzer; analysis consumes an explicit receipt and can reuse its sealed model.
- Small data types needed to represent presence and nullness in generated code.
- Proposed [tooling for optional horos/wire](docs/hydrated-types.md), adapting
  its canonical declarations and lifting rules to ASTs and shared semantic facts.
  Bittype does not own that model or its mixed value domain.

Horos owns data descriptors, conformance and module/generic/application
semantics for data declarations, their encoding and semantic vectors. Ontos owns
ground values, structural identity and encoding. Horos does not depend on bittype;
the current Horos-specific frontend and future qualified providers/transformers
consume its model. Generic foundations and native text stages remain independent.

**Owner boundary correction, 3 November 2026:** ground Ontos and Horos remain
data-only. Optional ontos/wire owns a mixed value profile; optional horos/wire
owns its type model. These proposals replace the Bittype-owned contract model,
without reviving the earlier ground Horos revision-2 proposal. Hydration stays
in bitwire/bitruntime; Function, Stream and Cell remain ordinary libraries.
The [preservation map](docs/wire-extension-preservation.md) records retained
concepts, explicit supersessions and future qualification obligations.

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

The [Horos declaration integration](docs/horos.md) implements checked canonical
module values to AST, syntax admission, an optional semantic analyzer and an
explicitly registered AST-only provider.
Generic AST/CST carriers,
typed views and pipeline contracts, four native type vocabularies, native parsers,
AST-to-CST formatters, syntax checks and the shared CST printer are implemented.
The [native text design](docs/parsing-formatting.md) defines their ownership and
qualification limits. Specialized providers all extend language.Provider.

The ASTTransformer boundary and syntax engine now provide checked heterogeneous
registration, route planning/execution and native text composition. The canonical
readable source syntax, native AST-to-AST mappings, public builders and application CLI remain
deferred. Native parsing does not make a target language the
canonical authority. The [generation design](docs/generation.md) remains separate.

This scope excludes application-value validators, serialization, Wire adapters,
RPC/streaming implementations, consumer migration and package publication.
Function, Stream, Cell and future behavioral contracts remain ordinary library
declarations with no privileged cases in shared infrastructure. Only the Horos
frontend depends on the Horos Go declaration model and Ontos values. Architecture
checks enforce that narrow boundary; generic carriers, contracts, engine, printer
and native text stages have no external Horos/Ontos imports. The engine depends
on the generic analysis contract, whose graph does not encode language rules.

The optional Horos Wire frontend and its narrow model dependency allowance
are future work. Their proposal does not alter these current import checks,
instantiate a runtime or extend the implemented Horos language profile.

## Versioning and evidence

A published language revision fixes its admitted syntax and meaning.
The earlier `bittype: 1` proposal is superseded. Horos's declaration revision is
implemented in its Go source; released 0.6.0 packages still contain ground checking.

Canonical declaration identity, normalized value shape, generated-code
compatibility and undeclared behavior remain distinct. A digest does not prove
behavior or authorize a binding substitution. Specify normalization and identity
content before implementing an encoding or publishing digests. Horos's specification
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
generation. Historical evidence does not add
callables, adapters or migration gates to this milestone.

The optional [Horos/Wire frontend and explicit ground-to-optional lift](docs/horoswire.md)
are implemented through the existing stage and transformer interfaces. Canonical
semantics stay in Horos; no native generator or runtime binding is introduced.
