# Charter

**Current direction: 6 November 2026.** The owner selected horos for canonical
data type semantics and bittype for syntax/tooling and native type generation,
starting with types only. The [model boundary](docs/model.md) replaces the earlier
independent bittype type-language proposal. Horos declaration revision 4 and its
bittype AST projection are implemented in Go source; no package release is implied.

## Architectural correctness

The owner's [6 November standing rule](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md)
makes semantic ownership and additive composition part of correctness. A working
example or passing suite is evidence for its exercised behavior, not a waiver of
the architecture. The goal is to preserve independently owned meanings and their
ability to evolve and compose without changes to shared orchestration or
unrelated implementations within the declared extension boundary.

A first implementation may support a small domain, but it must use the same
extension contract available to later contributions. Scope limits and
architectural defects are different. This applies to concrete reusable
compositions as well as the lowest-level kernel. The [Horos/Go verification case](docs/architecture-cases/0001-verification-composition.md)
records the original closure and the qualification of an independent mapping
through the public contribution contract.

## Ownership

bittype owns:

- Readable source syntax and AST/CST projections of horos declarations.
- Frontend parsing, resolution tooling and located diagnostics implementing
  horos's declaration/binding rules, without a second semantic specification.
- Native type generation for TypeScript, Go, Rust and C++.
- The [public generation SPI](docs/generation-spi.md) and guarded declarative
  planning, with language and mapping interpretation supplied explicitly.
- Native AST vocabularies and an agnostic language identity/registration
  boundary, independent of particular language providers.
- Native representation profiles and independent generation/round-trip fixtures.
- The [agnostic semantic model](docs/semantics.md) and implemented
  [analysis-provider boundary](docs/analysis.md); language providers retain their
  semantic rules. The separate [resolver contract](docs/resolution.md) uses the
  same semantic infrastructure. Horos supplies the first production resolver and
  analyzer; analysis consumes an explicit receipt and can reuse its sealed model.
- Small data types needed to represent presence and nullness in generated code.
- Implemented [unified Horos tooling](docs/horos.md), adapting canonical
  general kind contracts to ASTs and shared semantic facts. Bittype does not own the
  canonical type model or the mixed value domain.

Horos owns data descriptors, conformance and module/generic/application
semantics with general structural kinds, their encoding and semantic judgments. Ontos owns
ground values, structural identity and encoding. Horos does not depend on bittype;
the current Horos-specific frontend and future qualified providers/transformers
consume its model. Generic foundations and native text stages remain independent.

**Owner boundary correction, 4 November 2026:** one Horos declaration language
contains data and capability types, with explicit generic bounds. Ground values
and conformance stay data-only; Ontos/wire owns mixed values. Separate Horos/Wire
profiles and Data schema wrappers are removed, without aliases. Hydration stays
in Bitwire/Bitruntime, and Function, Stream and Cell remain ordinary libraries.
The [preservation map](docs/wire-extension-preservation.md) records the retained
concepts and clean migration; no runtime behavior follows from type validity.

The [general kind model](docs/kinds.md) replaces special Type/Data claims and
dedicated Data proofs. Scalar, List, Record and recursive Data are ordinary kind
definitions. Bittype's kernel carries provider-owned facts without interpreting
them. Native representation profiles distinguish representing Horos-admitted
uses from enforcing source constraints on arbitrary native callers.

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
explicitly registered text and semantic provider. The [Horos source grammar](docs/horos-syntax.md)
implements parsing, formatting and printing on the same AST.
Generic AST/CST carriers,
typed views and pipeline contracts, four native type vocabularies, native parsers,
AST-to-CST formatters, syntax checks and the shared CST printer are implemented.
The [native text design](docs/parsing-formatting.md) defines their ownership and
qualification limits. Specialized providers all extend language.Provider.

The ASTTransformer boundary and syntax engine now provide checked heterogeneous
registration, route planning/execution and native text composition. Pagination AST-to-AST mappings are implemented; public builders and application CLI remain deferred. Native parsing does not make a target language the
canonical authority. The [generation design](docs/generation.md) remains separate.

The [type-owned generator plugin design](docs/type-generation-plugins.md)
places builtin and library native mappings behind one supplied contract.
The [public SPI](docs/generation-spi.md) uses declarative descriptions and
kernel-owned planning. A closed native program supplies target construction;
source readers, mapping callbacks and source-substitution templates are absent
from that construction input. All four generators use this contract, and the
superseded procedural interfaces are removed. No canonical builtin or
named-library catalogue is added to shared generation or the engine.

The [unfinished generation boundaries](docs/generation-gaps.md) remain open.
GEN-01 corrects the source reader's mix of common semantic navigation and
source-owned interpretation. GEN-02 concerns preserving constraints for native
consumers; the current representation profile permits erasure and the enforcing
profile is unsupported. Passing Page<Cell> tests does not close either gap.

This scope excludes application-value validators, serialization, Wire adapters,
RPC/streaming implementations and service migrations. The unified declaration
cut includes maintained frontend consumers and the required public Horos release.
Function, Stream, Cell and future behavioral contracts remain ordinary library
declarations with no privileged cases in shared infrastructure. Only the Horos
frontend depends on the Horos Go declaration model and Ontos values. Architecture
checks enforce that narrow boundary; generic carriers, contracts, engine, printer
and native text stages have no external Horos/Ontos imports. The engine depends
on the generic analysis contract, whose graph does not encode language rules.

The one Horos frontend and its private helpers have the only concrete model
import allowances. Removing the second frontend changes no agnostic SPI.

## Versioning and evidence

A published language revision fixes its admitted syntax and meaning.
The earlier `bittype: 1` proposal is superseded. Horos's declaration revision is
implemented in Go and versioned with the public Horos distribution. Rust and
TypeScript retain ground checking; Bittype supplies readable syntax and native
backends, with Horos mapping packages and optional explicit compositions in Bittype's root
module. Horos owns the canonical model; it has no generator dependency on Bittype.

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

The unified [Horos frontend](docs/horos.md) implements the existing stage
interfaces. Template lifting and separate profiles are removed without aliases.
Canonical semantics stay in Horos; no native generator or runtime binding is added.
