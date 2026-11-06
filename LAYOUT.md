# Repository layout

**Owner decision, 31 October 2026:** organize by processing stage and enforce a
uniform structure across languages. The stage-by-language matrix is deliberate:
the same responsibilities occupy the same paths, making omissions and drift
visible. One public root Go module contains independently owned packages,
including the [concrete verification composition](docs/type-system-architecture.md#proposed-package-roles).

## Current hierarchy

```text
go.mod
files/go/                supplied Files and captured file/tree/source links
source/go/               sealed Source, organization and retained semantic evidence
trees/
  ast/go/                 generic AST, nodes, addresses, origins and borrowed view
  cst/go/                 generic concrete tree and stored AST view
  representation/go/      structurally checked AST/CST snapshots and language/profile descriptors
contracts/
  language/go/            language.Provider identity
  pipeline/go/            typed Stage contract
  parse/go/               symmetric Syntax/Files parse contracts and adapters
  check/go/               AST syntax-admission contract and registration adapter
  format/go/              symmetric Syntax/Files format contracts and adapters
  print/go/               symmetric Syntax/Files print contracts and adapters
  text/go/                shared source and syntax diagnostics only
  transform/go/           ASTTransformer contract and typed registration bridge
  semantic/go/            common snapshots, requests, reports and publication
  resolution/go/          resolver capability and dependency contract
  analysis/go/            analysis capability and explicit resolution input
  verification/go/        agnostic joint capture/receipt/site association
semantics/go/             typed immutable Semantics and provider-owned facts
semantics/golang/go/      Go vocabulary and immutable native type readers
internal/golangsem/       private direct Go compiler bridge and fact export
engine/go/                explicit registration, planning and execution
syntax/<language>/go/              abstract/concrete vocabulary and identity
parse/syntax/<language>/go/        text to CST
check/syntax/<language>/go/        AST syntax admission
format/syntax/<language>/go/       AST to CST
parse/files/<language>/go/         Files to Source (Horos implemented)
resolve/files/<language>/go/       reference resolution (Horos and bounded Go)
analyze/files/<language>/go/       semantic analysis (Horos and bounded Go)
format/files/<language>/go/        Semantics to Source (Horos implemented)
print/syntax/go/                   shared CST to text
print/files/go/                    shared Source to Files
provider/<language>/go/
  provider.go                     identity and direct typed composition
  syntax.go                       Syntax delegates and constructors
  files.go                        Files delegates and constructors
check/syntax/ast/go/               generic abstract carrier checks
check/syntax/cst/go/               generic concrete carrier checks
internal/<language>/     private language-specific lexical helpers
internal/                shared private checking/parsing/tree mechanics
transforms/horos-to-<target>/go/  analyzed Horos to native type ASTs
internal/horosgen/        source-owned semantic reading and bounded target arenas
integration/go/          independent consumers and cross-package qualification
verification/            concrete static verification; public Logos contract
  files/source/horos/go/  explicit authoritative queries bound to captured Horos
  files/horos-to-go/go/   selected mapping, prepared claim and direct Logos check
  integration/go/        TextBox/Page, nested reuse and dependency qualification
tools/provider-matrix/go/ development-only source-derived documentation generation
tools/horos-module/go/    ordinary Horos source to canonical Ontos module command
testdata/text/<language>/valid/
docs/
```

The native text language directories are `typescript`, `golang`, `rust` and
`cpp`; `horos` and `golang` also implement resolve/analyze.
The final `go` directory denotes implementation language, not target:
`syntax/rust/go` is Go code describing Rust syntax. `ast` and `cst` in the check
stage are reserved generic-carrier components, not native language names.

**Owner cleanup, 5 November 2026:** implementation paths follow
`<role>/<dimension>/<language>/go`; shared implementations omit the language.
Both `syntax` and `files` are explicit dimension segments. Language vocabularies
remain at `syntax/<language>/go`; composition packages span both dimensions and
remain at `provider/<language>/go`. No empty directories represent missing roles.

Implementation Go packages are named by their role: `parse`, `check`, `format`,
`print`, `resolve`, `analyze` and `provider`, for every language and dimension.
Language vocabulary packages retain their language names. Consumers choose local
import aliases when importing a contract and implementation with the same name.
Those import aliases are not exported compatibility APIs.

Every composition package has `provider.go`, `syntax.go` and `files.go`.
`NewSyntaxParser`, `NewSyntaxChecker`, `NewSyntaxFormatter`, `NewSyntaxPrinter`
and `NewFilesPrinter` each take an explicit `representation.Language` and return
one configured adapter plus an error. Horos additionally has `NewFilesParser`,
`NewFilesFormatter`, `NewFilesResolver` and `NewFilesAnalyzer` with the same
descriptor-first signature. Constructors register nothing and never select a
profile implicitly. `Provider` remains a direct typed composition of the roles
that actually exist; it is not an engine registration contract.

Stateless language stages expose a free operation and a zero-value `Provider`.
Both shared printer packages expose `Provider[L]` and `New[L]`; Files printing
also receives its Syntax printer collaborator. Their `provider.go` files own
construction and `print.go` files own emission. Files parsing/formatting retain
their explicit semantic collaborators. Uniform construction does not erase
these genuine dependency differences. Go now supplies bounded native semantic
providers through the same independent role constructors.

The source and its package-level tests belong with each implementation.
Independent consumer fixtures, stage composition, language conformance and native
compiler qualification remain under integration. Private helpers are not an
alternative public API. Do not create a language-first packages tree, additional
Go modules, compatibility aliases or placeholder source directories.

The [Go admission proposal](docs/go-native-admission.md) is qualified by
`integration/go/go_native_admission_test.go` and isolated model/issuer/consumer
fixtures under `integration/go/testdata/go-admission/`. These are compiler inputs,
not production evidence APIs or generated code. The test creates its temporary
fixture module at execution time; no nested Go module is committed.

## Dependency direction

Directory groups are not Go packages or an assertion that all members share the
same dependency level. The package boundaries are:

| Component | Permitted responsibilities and dependencies |
| --- | --- |
| AST, Stage, language identity | Independent foundations; no native syntax or other foundation dependency |
| CST | Generic AST carrier |
| Parse contract/adapter | Syntax and Files inputs, Source publication, carriers, semantic evidence and identity |
| Check contract/adapter | AST views, representations and language identity; no concrete grammar or text-stage dependency |
| Format contract/adapter | AST/CST and Source carriers, semantic evidence and identity |
| Print contract/adapter | CST and Source/Files carriers, representations and identity |
| Text diagnostics | AST addresses only; no stage contracts |
| Carrier checks and printer | Generic tree structure and shared diagnostic contracts; no native grammar or formatting choices |
| Structurally checked representation | Carriers, identity and generic carrier checks; no engine or transformer dependency |
| Transformer contract/bridge | Typed requests, shared semantic evidence admission, checked representations and endpoint identity; no concrete mapping |
| Semantic facts | Generic AST addresses; graph integrity and owned storage, no language rules |
| Shared semantic contract/bridge | Semantic facts, Files context, representations, explicit checker and identity; no Source dependency |
| Resolution contract/bridge | Shared semantic contract and identity; binding coverage and typed invocation |
| Analysis contract/bridge | Shared semantic and resolution contracts; explicit dependency composition |
| Engine | Checked representations and shared contracts; no concrete language, provider or mapping |
| Native syntax | Carriers and identity; no implementation stage or other language |
| Native checker | Own syntax, shared checking mechanics and own lexical helpers |
| Native parser/formatter | Own syntax/checker and appropriate shared mechanics; neither imports the other |
| Native provider | Own syntax and stages plus the shared printer; composition only |
| Horos syntax | Generic AST/CST and identity; no semantic-model imports |
| Horos checker | Own syntax and generic structural/grammar checks; no semantic resolution |
| Horos resolver | Own syntax, private projection/correspondence helper, Horos Go model and Ontos values/codec |
| Horos analyzer | Own syntax, Horos resolver contract and shared semantic evidence; no direct external model imports |
| Horos parser/formatter | Own syntax/checker, private spelling and canonical-value helpers, Ontos values; no semantic resolution |
| Horos provider | Own text/checker/resolver/analyzer stages and shared printer; composition only |
| Horos module command | Existing Horos parser, private module lowering and Ontos encoder; no semantic resolution or prelude-specific logic |

The parse, check, format and print interfaces are independent; their configured
adapters use representations for owned heterogeneous inputs/outputs. The engine
registers them separately after identity/profile registration. Representation
depends on generic carrier checks and diagnostic records, never a stage-provider
contract. Source input is `parse.Input`; `text` contains diagnostics only. The shared implementation remains under `print/syntax/go`, separate
from `contracts/print/go`. Imports may use a local alias when both packages are
needed. No compatibility interfaces or input aliases remain under text.

The structurally checked AST/CST and language descriptors share one home in
`trees/representation/go`. `contracts/transform/go` owns ASTTransformer and
ErrInvalidTransformer; representation owns tree and identity errors. Its narrow
CheckMarker operation checks actual Go type identity without exposing payloads
or providing unchecked construction. No aliases remain at the former paths or
under engine for these moved types.

Every private helper also has a dependency rule. A shared helper cannot conceal
a native-language import, and a language helper cannot depend on another
language or a stage. This keeps the rule enforceable through indirect imports.

## Enforced consistency

The [architecture checks](integration/go/architecture_test.go) discover the
union of language directories across syntax, provider, Syntax stages and Files semantic stages
and require each language to have every stage in its tested capability profile.
The four native languages retain all five native text stages; Horos requires
the same Syntax stages plus Files resolve/analyze, with explicit qualification.
Generic carrier check directories
are excluded explicitly. An orphan parser and an incomplete syntax package both
fail. The checks also classify every production package and enforce its directed
dependencies, including private helpers and narrow external-model allowances
only for the Horos frontend and source-owned transformation reader. Generic carriers,
contracts, engine and printer must remain independent of Horos/ontos, and native
syntax/text stages must remain independent of the canonical frontend.

The [native conformance suite](integration/go/conformance_test.go) covers
every native language. It checks stable identity, matching AST/CST checker
views, formatter snapshot ownership and metadata, whole-input rejection, located
parse errors and empty output on failure. The same test-only fixture table feeds
source round trips and fuzzing. Tests never register a production default.
Independent compiler fixtures remain a separate required qualification.
The [semantic suite](integration/go/horos_test.go) separately checks Horos identity,
registration, checking, snapshot ownership and specialization. Horos text tests qualify
source retention, canonical pagination round trips and independently registered stages.

Existing public consumer/negative type-check cases, unfamiliar-language tests,
route oracle, concurrency checks and compiler gates remain required. Directory
consistency does not prove grammar or mapping correctness.

## Verification and future components

The [evaluated type-system reasoning architecture](docs/type-system-architecture.md#proposed-package-roles)
places `verification/files` collaborators above existing captures and explicit
providers. Its scoped composition and agnostic association are now implemented;
language theories and Logos proof machinery stay with their owners.
The [Logos contract map](docs/logos-verification.md#ownership-map-for-the-verification-spi)
narrows any eventual shared contract to artifact association and explicit
composition. A direct Logos checker call needs no second checker-provider registry.
The [joint request/result contract](docs/verification.md) defines that association
and its TextBox/Page walkthroughs. Logos-specific prepared requests/results stay
in the concrete integration; general contracts retain only agnostic association.
The implemented paths appear in the current hierarchy above. A wider verification
provider registry and additional language interpretations remain deferred.

The [agnostic semantic model](docs/semantics.md), [analysis contract](docs/analysis.md)
and first Horos analyzer are implemented. Additional analyzers use the same
optional analyze stage with an explicitly tested capability profile. Native
semantic analysis is not implied by native text tooling.

The [Horos declaration integration](docs/horos.md) implements
`syntax/horos/go`, `check/syntax/horos/go`, `analyze/files/horos/go`, `provider/horos/go` and
`resolve/files/horos/go`, with private projection/correspondence in `internal/horosmodel`.
Readable syntax now implements `parse/syntax/horos/go` and
`format/syntax/horos/go` under the same stage matrix. These replace the former proposed
`canonical` paths; there are no compatibility aliases.

Directed pagination mappings live under `transforms/horos-to-<native>/go`; later
partial reverse mappings would use `transforms/<native>-to-horos/go`. Application
composition and IO are proposed under `cmd/bittype/go`. Horos owns the declaration semantics and
encoding; these packages supply tooling projections and native representations.
New dependency rules and qualification accompany further implementation, with no
empty packages created now. All five current languages have qualified text stages;
new AST-only providers still require an explicit tested capability profile.

Adding a native language supplies the five stages, any private lexical helpers,
conformance fixtures and compiler qualification. It requires no change to the
generic carriers, engine or production registration catalogue: there is no such
catalogue. Applications choose their registrations. Adding Function, Stream,
Cell or another library declaration does not add a stage, language or core kind.

The [unified declaration model](docs/horos.md) uses one `horos` family under
syntax/check/resolve/analyze/provider. Private syntax framing remains in
`internal/declsyntax`; model projection, correspondence and export remain in
`internal/horosmodel`. No `horoswire` stage or compatibility alias remains.
Bounds are explicit syntax and provider-owned semantic facts, not kernel kinds.

Model dependency allowances cover only this frontend and its private bridge.
Architecture checks preserve the agnostic and native-stage boundaries. Native
pagination mappings use `transforms/horos-to-<native>/go`.
Template lifting and its exclusive dependency allowances remain removed.

The [native representation design](docs/type-representations.md) and
[type-owned generator plugin proposal](docs/type-generation-plugins.md) place
builtin and library mappings behind one supplied contract. Each type library owns
its mappings and type artifacts outside Bittype's module; the consumer selects
defaults and explicit overrides. Adding Cell or another ordinary contract adds no
Bittype stage, engine edge per contract or central registration entry. Plugin
contracts remain proposed; the private pagination Dialect remains implemented.

The proposed generation extension preserves stage-first organization and makes
its Syntax dimension explicit:

```text
contracts/generation/go/                  shared configuration/plugin contracts
generation/syntax/go/                     language-neutral composition
generation/syntax/source/horos/go/        source semantic adapter and selectors
generation/syntax/target/typescript/go/   TypeScript fragments/symbols/assembly
generation/syntax/target/golang/go/       Go fragments/symbols/assembly
generation/syntax/target/rust/go/         Rust fragments/symbols/assembly
generation/syntax/target/cpp/go/          C++ fragments/symbols/assembly
transforms/horos-to-<target>/go/          configuration-only transformer factories
```

The source/target subdirectories distinguish collaborators within Syntax
generation, not additional dimensions or registered providers. Their identical
package structure is specified in the proposal. Existing transformer path names
remain `horos-to-typescript`, `horos-to-go`, `horos-to-rust` and `horos-to-cpp`;
native vocabulary/implementation directories use `golang` for Go as elsewhere.

Builtin mappings and opinionated default assembly belong in a Horos-owned
integration module, with `generation/syntax/<language>/go` for each target. Other
type libraries use the same layout in their own integrations. These may be
separate Go integration modules inside their owning repositories; they are not
new modules inside Bittype or a required new repository. Bittype factories never
import them. This preserves the dependency direction: integration -> Bittype,
consumer -> integration + Bittype; Horos core remains independent of Bittype.

No proposed source directory is created by the documentation task. Syntax/Files
providers and source-derived capability matrices remain unchanged.

The [Syntax and Files design](docs/source.md) is implemented with
`Syntax : CST : AST` and `Files : Source : Semantics`. Both dimension names are
explicit in implementation paths. Horos composes the full Files text pipeline; Go adds bounded Files resolution
and analysis without Files parsing/formatting;
all five language composition packages construct Files printers.

Files depends on AST addresses. Source retains semantic capture/result handles,
CSTs and typed Semantics. Shared semantic capture depends on Files but never
Source, and tree representation must not import Source. The symmetric role
contracts depend on those carriers; concrete Files compositions depend on explicit
syntax/check/resolution/analysis contracts. `internal/horosfiles` is private
composition, not a substitute language registry. Generic infrastructure never
imports Horos or Ontos. No layout package or combined Source provider is retained.

The [native provider design](docs/native-semantics.md) specifies the compiler
boundary and explicit query sessions. Native stages do not depend on Horos or
Logos; joint verification remains the consuming composition.
