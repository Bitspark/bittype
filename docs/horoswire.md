# Optional Horos/Wire frontend and lifting

**Implemented in Go, 4 November 2026.** The optional frontend projects
[horos/wire revision 1](https://github.com/Bitspark/horos/blob/v0.7.0/docs/spec/wire-types.md).
Its canonical input is ground Ontos module data describing mixed types.
It never accepts live endpoints as source syntax or performs hydration.
Ground Horos retains its existing domain and public profiles.

## Packages and dependencies

| Package | Responsibility |
| --- | --- |
| [syntax/horoswire/go](../syntax/horoswire/go/nodes.go) | Named Go language marker, typed AST and declaration vocabulary |
| [check/horoswire/go](../check/horoswire/go/check.go) | Syntax admission, including unresolved applications |
| [resolve/horoswire/go](../resolve/horoswire/go/resolver.go) | Authoritative Horos resolution and export into shared semantic facts |
| [analyze/horoswire/go](../analyze/horoswire/go/analyze.go) | Consume an explicit compatible resolution receipt and reuse its model |
| [provider/horoswire/go](../provider/horoswire/go/provider.go) | Compose check, resolution and analysis capabilities |
| [transforms/horos-to-horoswire/go](../transforms/horos-to-horoswire/go/lift.go) | Explicit partial lift, canonical derivation evidence and AST origins |
| [internal/declsyntax](../internal/declsyntax) | Shared syntax framing; no external model dependencies |
| [internal/horosmodel](../internal/horosmodel) | Private shared projection, correspondence, resolution and fact-export mechanics |

The two concrete resolvers supply fixed model entry points, canonical prefixes
and semantic profiles to those private mechanics. There is no profile registry,
constructor switch or inferred selection in that helper. Horos determines type
positions and binding. The agnostic trees, semantic model, engine, Wire and
runtime acquire no new kinds, imports or registrations. Native language stages
remain independent of Horos/Ontos. Architecture tests enforce these boundaries.

The dependency in `go.mod` is public Horos 0.7.0 with Ontos 0.16.0. Consumers must
upgrade together; there is no local replacement, legacy dependency path or
compatibility shim. Neither Bittype frontend
depends on the live `ontos/wire` carrier. Its static input remains ordinary Ontos
data. Go is the implemented tooling language; this increment makes no claim of
Rust or TypeScript optional declaration APIs.

## Profiles, syntax and semantics

| Boundary | Ground Horos | Optional Horos/Wire |
| --- | --- | --- |
| Language ID | `horos` | `horoswire` |
| Display name | `Horos` | `Horos/Wire` |
| AST profile | `horos-declarations-1/ast-2` | `horos-wire-1/ast-1` |
| Resolution profile | `horos-declarations-1/resolution-1` | `horos-wire-1/resolution-1` |
| Analysis profile | `horos-declarations-1/analysis-2` | `horos-wire-1/analysis-1` |
| Semantic key owner | `horos` | `horos/wire` |
| Canonical module tag | `horos/decl/module` | `horos/wire/module` |

Both use the [declaration AST framing](horos.md#vocabulary-and-bindings): Program,
Module, Import, Declaration, TypeParameter, Parameter, Application, Descriptor,
Atom and Tuple. A Wire expression is an ordinary provider-owned Descriptor with
`constructor = "wire"` and a message-type operand. Data has `constructor = "data"`;
its closed ground descriptor remains literal Atom/Tuple syntax. Upper parameter
binding does not cross that boundary, even when a literal resembles a parameter.

`Resolve(root, dependencies, treeID)` admits canonical optional modules and
projects their complete import closure. `Project(program, treeID)` projects an
already checked optional program, re-admitting it under the AST profile's fixed
model limits. `Model(view)` captures the borrowed view once, checks syntax,
reconstructs the canonical documents and returns the optional Horos program.
It rejects incorrect document pins, type/literal positions and supplied binding
annotations; it never repairs caller syntax. `Program.Specialize` can then bind
an optional slot to a Wire expression. Ground `Model` and specialization cannot.

Registration remains explicit through `representation.NewLanguage`,
`resolution.NewResolver` and `analysis.NewAnalyzer`, followed by their engine
registration methods. Syntax registration installs no semantic capability or
transformation. Analysis requires the exact compatible resolution receipt and
reuses its immutable model. Ground receipts and configuration keys cannot stand
in for optional-profile evidence. Neither provider has a text parser/formatter.

Successful resolution establishes static declaration validity and bindings.
Provider-owned semantic terms retain Wire/Data operand locations; the generic
semantic model does not interpret them. Exhaustion produces incomplete evidence,
not a type rejection. Missing dependencies, invalid source and invalid request
configuration remain distinct outcomes. Nothing establishes capability authority,
liveness, remote protocol behavior or instance conformance.

## Explicit transformation

The `lift` package's `Derive(input, options)` takes a ground Horos AST. Options
select exact source declaration identities, the source entry module, a distinct
output tree address and Horos's `LiftOptions`. It delegates the ordered partial
mapping to Horos, then projects the derived optional program. The result contains
the target tree and immutable canonical `Derivation` evidence.

For `Page<T> = list-of(T)`, the output is a distinct optional `Page<T> =
list-of(T)`. Its slot can be specialized with `Wire(Data(Person))`. An authored
application `Page<Person>` retains an application to the derived Page and the
corresponding lifted argument. Parameter order, caller/callee scope, imported
targets and source identities remain explicit. There is no duplicate authored
Page definition or operation table.

The selected output entry must reach every selected root. A one-AST transformer
refuses independent roots outside that entry instead of silently dropping them.
Use separate transformations or Horos's multi-module derivation API for that
case. Generated module/import/declaration/binder/expression nodes carry direct
source addresses and copied inherited origins. Canonical pins do not include
these AST addresses. Closed embedding can prune derived dependencies while the
Horos derivation retains the original source pins and expression paths.

`NewTransformer(from, to, options)` checks the exact endpoint profiles and captures
the selection, including copies of roots and naming maps. The caller explicitly
registers the returned existing `ASTTransformer`. Registration adds one partial
language edge, not an edge per library contract. A failed lift returns no usable
tree, derivation or envelope; it never retries embedding or another route.

Selection entries and bytes are bounded before copying configuration or resolving
source. Provenance accumulation has its own bounded storage, and the complete
target AST must also pass the fixed carrier/profile limits before publication.

Horos's separate explicit `Embed` operation remains available for a concrete
ground specialization. The transformer does not call it as a fallback. For
example, lifting a closed application of `Other<T> = not(T)` still fails because
the retained generic target has no lift. The closed specialization can instead
be deliberately embedded through Horos. See the full
[lifting rules](type-representations.md#2-embedding-and-lifting-are-explicit-different-operations).

## Qualification and remaining scope

[Frontend fixtures](../integration/go/horoswire_test.go) exercise optional
pagination, generic Function/Stream/Cell and unfamiliar Mailbox declarations,
profile separation, literal opacity, syntax-only admission, located semantic
rejection, budget exhaustion and explicit resolution/model reuse. These names
exist in consumer fixtures, not in production dispatch or registration.
[Transformation fixtures](../transforms/horos-to-horoswire/go/lift_test.go)
exercise template/application preservation, distinct identities, copied origins,
explicit registration, captured configuration and refusal without fallback.
Existing ground frontend, architecture, native compiler and text tests remain
required. Compilation and formatting run on the repository's remote gates.

Native type generators and library-owned representation recipes remain future
work. Readable Horos text, mixed application-value validation, executable
adapters, hydration alignment, library protocols and package releases are not
implemented here. Their preserved designs remain in the
[concept ledger](wire-extension-preservation.md) and [runtime guide](runtime-values.md).
