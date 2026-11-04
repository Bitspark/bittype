# Optional Horos/Wire frontend

**Implemented in Go, 4 November 2026.** The optional frontend projects
[horos/wire revision 1](https://github.com/Bitspark/horos/blob/v0.8.0/docs/spec/wire-types.md).
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
| [internal/declsyntax](../internal/declsyntax) | Shared syntax framing; no external model dependencies |
| [internal/horosmodel](../internal/horosmodel) | Private shared projection, correspondence, resolution and fact-export mechanics |

The two concrete resolvers supply fixed model entry points, canonical prefixes
and semantic profiles to those private mechanics. There is no profile registry,
constructor switch or inferred selection in that helper. Horos determines type
positions and binding. The agnostic trees, semantic model, engine, Wire and
runtime acquire no new kinds, imports or registrations. Native language stages
remain independent of Horos/Ontos. Architecture tests enforce these boundaries.

The dependency in `go.mod` is public Horos 0.8.0 with Ontos 0.16.0. Consumers must
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

These are the current implementation's separate profiles and frontend IDs, not
a claim that data-only constraints inherently require two declaration languages.
This cleanup preserves those public boundaries and does not introduce generic
bounds or a consolidated declaration grammar.

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
`check.NewChecker`, `resolution.NewResolver` and `analysis.NewAnalyzer`, followed
by their independent engine registration methods. Identity registration installs
no checker, semantic capability or transformation. Analysis requires the exact compatible resolution receipt and
reuses its immutable model. Ground receipts and configuration keys cannot stand
in for optional-profile evidence. Neither provider has a text parser/formatter.

Successful resolution establishes static declaration validity and bindings.
Provider-owned semantic terms retain Wire/Data operand locations; the generic
semantic model does not interpret them. Exhaustion produces incomplete evidence,
not a type rejection. Missing dependencies, invalid source and invalid request
configuration remain distinct outcomes. Nothing establishes capability authority,
liveness, remote protocol behavior or instance conformance.

## Direct optional declarations

Author `Page<T> = list-of(T)` directly in horos/wire when T must accept data or
capabilities. `Resolve` projects that canonical module; `Model` recovers its
checked optional program. `Program.Specialize` can bind the same Page to
`Data(Person)` or `Wire(Data(Person))`. The source generic and its identity remain
unchanged. Imported applications preserve exact targets, argument order, caller
scope and binder identity. No second Page or generated operation table is needed.

Horos's separate `Embed` operation specializes a ground declaration and wraps its
closed descriptor in Data. It creates no optional declaration and cannot widen
a ground generic. See [direct authoring and embedding](type-representations.md#2-direct-declarations-and-ground-embedding).

Template lifting and `transforms/horos-to-horoswire/go` were removed completely
on 4 November 2026, along with derivation-specific indexing and dependency
allowances. There is no compatibility transformer. The generic ASTTransformer
contract and engine remain; no production ASTTransformer implementation exists
yet. The next mappings are source-profile-to-native type generation.

## Qualification and remaining scope

[Frontend fixtures](../integration/go/horoswire_test.go) exercise optional
pagination, generic Function/Stream/Cell and unfamiliar Mailbox declarations,
profile separation, literal opacity, syntax-only admission, located semantic
rejection, budget exhaustion and explicit resolution/model reuse. These names
exist in consumer fixtures, not in production dispatch or registration.
The [installed-consumer fixture](../integration/go/testdata/installed/consumer_test.go)
authors ground and optional pagination directly, resolves the public package
model and specializes the same optional generic with both data and Wire arguments.
It also verifies that ground generic slots reject Wire arguments.

Existing ground frontend, architecture, native compiler and text tests remain
required. Compilation and formatting run on the repository's remote gates.

Native type generators and library-owned representation recipes remain future
work. Readable Horos text, mixed application-value validation, executable
adapters, hydration alignment, library protocols and package releases are not
implemented here. Their preserved designs remain in the
[concept ledger](wire-extension-preservation.md) and [runtime guide](runtime-values.md).
