# Repository layout

**Owner decision, 31 October 2026:** organize by processing stage and enforce a
uniform structure across languages. The stage-by-language matrix is deliberate:
the same responsibilities occupy the same paths, making omissions and drift
visible. One public root Go module contains independently owned packages,
including the [concrete verification composition](docs/type-system-architecture.md#proposed-package-roles).

**Owner decision, 6 November 2026:** the accepted
[kernel-owned composition and uniform layout](docs/kernel-composition.md)
uses `<stage>/<dimension>/<scope>/<component>/go`, with generic factories in
kernel/contracts and no per-pair assembly or presets. This hierarchy is implemented.

**Owner direction, 13 November 2026:** "So the root split to me would be
between model and ops. model must not contain any ops and must not import from
ops either. ops work on the model". The design
[models and operations](docs/models-and-operations.md) carries it out: models
under `model`, by sort and then by language, and operations under `ops`, where
the stage-by-language matrix and the path pattern above stay. **Not
implemented.** The hierarchy and the role table below describe `main` as it
is, and no path in this layout has changed. The
[order of work](docs/models-and-operations.md#10-the-order-of-work) says which
slice moves what. The names below the two roots are proposals of that page
and are not decided.
[bittype/ACX-015](docs/architecture-cases/0015-models-hold-operations.md)
records the violation that stays open until the slices are done.

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
  generation/go/          declarative mappings, immutable native programs and guarded planning
  semantic/go/            common snapshots, requests, reports and publication
  resolution/go/          resolver capability and dependency contract
  analysis/go/            analysis capability and explicit resolution input
  verification/go/        agnostic joint capture/receipt/site association
semantics/go/             typed immutable Semantics and provider-owned facts
semantics/horos/go/       horos vocabulary: roles, relations, entity kinds, term operators and encoding, profiles, limits setting, prelude pin, contract identities
semantics/golang/go/      Go vocabulary, contract identities and immutable native type readers
internal/golangsem/       private direct Go compiler bridge and fact export
semantics/typescript/go/  TypeScript vocabulary, graph names, contract identities, compiler configuration and immutable native type readers
internal/typescriptsem/   private captured-source compiler helper and fact export
internal/typescriptimports/ private to TypeScript Files formatting: the source paths its profile admits and excludes, relative import specifiers and import node construction
semantics/rust/go/        Rust vocabulary, contract identities, environment and immutable native type readers
internal/rustsem/         private captured-source compiler bridge
  helper/                private pinned rustc_driver executable (Rust/Cargo)
engine/go/                explicit registration, planning and execution
syntax/<language>/go/                          vocabulary and identity
parse/syntax/<language>/provider/go/           text to CST
check/syntax/<language>/provider/go/           native syntax admission
format/syntax/<language>/provider/go/          AST to CST
parse/files/horos/provider/go/                 Files to Source
format/files/{horos,typescript}/provider/go/   Semantics to Source
resolve/files/<language>/provider/go/          captured reference resolution
analyze/files/<language>/provider/go/          semantic analysis
analyze/files/{golang,typescript}/queries/go/  native query service over an opened model
print/{syntax,files}/shared/provider/go/       shared printers
check/syntax/shared/{ast,cst}/go/              carrier structure
transform/syntax/shared/generation/go/         registry and request composition
transform/syntax/shared/target/go/             native construction and scoped symbol storage
transform/syntax/horos/source/go/              Horos semantic reader
transform/syntax/horos/mapping/go/             representation profile identities and the mapping library function type
transform/syntax/horos/mapping/pattern/go/     horos source patterns
transform/syntax/horos/mapping/declarations/go/ horos declaration recipe
transform/syntax/horos/mapping/naming/go/      horos name to exported native identifier
transform/syntax/<native>/model/go/            native semantic forms and validation
transform/syntax/<native>/target/go/           AST construction from native forms
transform/syntax/typescript/symbol/go/        public type symbols for generation and Files
transform/syntax/golang/symbol/go/            predeclared type symbols for generation and Files
transform/syntax/rust/symbol/go/              standard-library item and primitive type symbols for generation and Files
transform/syntax/horos-to-<native>/mapping/<name>/go/ one selectable mapping library per package
internal/<language>/                           private lexical helpers
integration/go/          independent consumers and cross-package qualification
verification/            concrete static verification; public Logos contract
  files/source/horos/go/  explicit authoritative queries bound to captured Horos
  files/horos-to-go/go/   scoped representation contributions and named questions
  files/horos-to-go/mappings/<name>/go/ one selectable law per package: scalar, list, optional, record
  logos/go/              typed interpretation selection, receipt association, Logos checking
  integration/go/        TextBox/Page, independent tuple composition, theorem reuse
tools/provider-matrix/
  go/                    development command; explicit extractor/result composition
  report/go/             public discovery, association, rendering and write host
  capability/go/         public source/definition/structural extraction contract
  extract/grammar/go/    exact owning Rules-format extractor
  observation/go/        invocation import and source-context association
tools/horos-module/go/    ordinary Horos source to canonical Ontos module command
testdata/text/<language>/valid/
docs/
```

The native text language directories are `typescript`, `golang`, `rust` and
`cpp`; `horos`, `golang`, `typescript` and `rust` also implement resolve/analyze.
The final `go` directory denotes implementation language, not target:
`syntax/rust/go` is Go code describing Rust syntax. `ast` and `cst` in the check
stage are reserved generic-carrier components, not native language names.

**Clarified on 13 November 2026 ([models and operations](docs/models-and-operations.md#5-intent)).**
The hierarchy says of `transform/syntax/<native>/model/go`: "native semantic
forms and validation". That describes the package as it is. Validation is an
operation, and under the owner's direction a model contains none. In the
design the package becomes a declared table of a language's intended forms on
the model side, one checker in the kernel reads the table of every language,
and the typed constructors are a package of the operations side. Not
implemented.

Implementations use the same stage, dimension, scope, component and implementation-language order. A scope identifies one language, an ordered pair or shared machinery. Go packages use their role name. Syntax vocabularies retain their language names. Applications construct independent stage adapters through their contracts. A configured resolver belongs to its language's resolve stage. See [kernel-owned composition](docs/kernel-composition.md).

**Path form for several contributions, decided in
[#222](https://github.com/Bitspark/bittype/issues/222) on 11 November 2026.** A
component that holds several independently selectable contributions has one
package per contribution, one level below the component:
`<stage>/<dimension>/<scope>/<component>/<name>/go`. The Go package is named
`<name>`. No package lists the contributions of a component. A consumer
registers each one it selects. The role stays readable from the path, which
the architecture check uses. Two components have this form today:

- **Mapping libraries:** `transform/syntax/horos-to-<native>/mapping/<name>/go`
  with `text`, `boolean`, `list`, `record` and `declarations` for each of the
  four native languages, and `optional` for Go. Each package exports `Library`.
- **Verification laws:** `verification/files/horos-to-go/mappings/<name>/go`
  with `scalar`, `list`, `optional` and `record`. The unit is the law. `scalar`
  exports the two rules `Text()` and `Boolean()` of one law. The other three
  export `Rule`.

A native query service has its own package beside its language's analysis
provider: `analyze/files/golang/queries/go` and
`analyze/files/typescript/queries/go`. Rust and horos have none. The Go package
moved out of `analyze/files/golang/provider/go` with pull request
[#246](https://github.com/Bitspark/bittype/pull/246).

Stateless language stages expose a free operation and a zero-value `Provider`.
Both shared printer packages expose `Provider[L]` and `New[L]`; Files printing
also receives its Syntax printer collaborator. Their `provider.go` files own
construction and `print.go` files own emission. Files parsing/formatting retain
their explicit semantic collaborators. Uniform construction does not erase
these genuine dependency differences. Go, TypeScript and Rust supply bounded native
semantic providers through the same independent role constructors.

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
same dependency level. A package gets its import rule from the role that its
path gives. The roles are in
[`architecture_roles_test.go`](integration/go/architecture_roles_test.go) and
the explicit tables are in
[`architecture_test.go`](integration/go/architecture_test.go). The table below
summarises them. The exact lists are in the test.

**Clarified on 13 November 2026 ([models and operations](docs/models-and-operations.md#8-roles-per-root)).**
The first sentence of this section holds for the directory groups below a
root. The two roots of the owner's direction are an assertion about the
dependency level, in one direction: no package under `model` imports a
package under `ops`. In the design the role Foundation is divided between the
roots: the agnostic model packages on one side, the contract packages and the
operations halves of the carriers on the other. The role Vocabulary exists
under both roots: names and declared tables under `model`, and identities
that are values of a contract package's type under `ops`. A language's typed
readers, typed constructors and symbol admission get a role of their own,
which that page proposes to call public part. A pair's contribution
vocabulary names two languages, as a mapping does. Not implemented: the table
below gives the roles of `main` as they are.

| Role | Paths | Imports |
| --- | --- | --- |
| Foundation | `trees/*/go`, `files/go`, `source/go`, `semantics/go`, `contracts/*/go` | Each package has its own entry in the table `foundationRules` and imports exactly what the entry lists. It never imports a package outside that table. The AST imports nothing. The CST, `files/go`, `semantics/go` and `contracts/text/go` import only the AST. `trees/representation/go` and `contracts/semantic/go` do not import `source/go` |
| Kernel | `engine/go`, `transform/syntax/shared/*/go`, `check/syntax/shared/*/go`, `print/*/shared/provider/go` | The same table and the same rule. So no foundation or kernel package imports a vocabulary, a contribution or a language's private part |
| Shared private helper | `internal/tree`, `internal/grammar`, `internal/lex`, `internal/parseutil`, `internal/concrete` | The same table and the same rule. A helper of no single language cannot conceal a language |
| Vocabulary | `syntax/<language>/go`, `semantics/<language>/go`, `transform/syntax/<native>/model/go`, `transform/syntax/<native>/symbol/go` | Foundations and its own language's vocabulary. `semantics/<language>/go` and `symbol` may also import their language's lexical helper. Never a contribution. `semantics/<language>/go` never imports `internal/<language>sem`. TypeScript, Go and Rust have a `symbol` package today |
| Contribution: stage provider | `<stage>/<dimension>/<language>/provider/go` for `check/syntax`, `parse/syntax`, `format/syntax`, `resolve/files`, `analyze/files`, `parse/files` and `format/files` | Each of the seven has its own list: foundations, shared private helpers and what belongs to its own language. A checker gets no concrete tree, a formatter no lexer, and no Syntax stage a semantic contract. A parser and a formatter import their own checker under an exception, and neither list names the other. A resolver may import its language's semantic helper. A Files parser or formatter reaches its collaborators through their contracts |
| Contribution: query service | `analyze/files/<language>/queries/go` | The shared semantic contract, the representation, and its own syntax, vocabulary and semantic helper. Never its resolver. Go and TypeScript have this package today. A package that is not a composition imports a query service only under an exception that names the symbols it uses |
| Contribution: source reader | `transform/syntax/<language>/source/go` | Foundations and its own syntax and vocabulary. Only horos has one today |
| Contribution: native target | `transform/syntax/<native>/target/go` | The generation contract, the AST, the shared target construction, and its own syntax, model and lexical helper. It reads no source and calls no mapping: a source reader, a mapping and a common mapping part are refused |
| Contribution: mapping | `transform/syntax/<source>-to-<native>/mapping/<name>/go` | The generation contract, the representation, the syntax of both languages, the source reader and the common mapping parts of its source language, and the model of its native language. Never another mapping and never a native target |
| Contribution: verification law | `verification/files/horos-to-go/mappings/<name>/go` | `verification/files/horos-to-go/go`, `semantics/golang/go`, `semantics/go` and the external models `github.com/Bitspark/horos/core/go/declarations` and `github.com/bitspark/ontos/core/go`. Never another law |
| Contribution: capability extractor | `tools/provider-matrix/extract/<format>/go` | `tools/provider-matrix/capability/go` only |
| Common mapping part | `transform/syntax/<source>/mapping/go` and `transform/syntax/<source>/mapping/<part>/go` | The generation contract, the representation, and the syntax and source reader of its source language. The list names no other part, and the one-language rule below refuses every native package |
| Private part | `internal/<language>`, `internal/<language>sem` | A lexical helper imports `internal/grammar` only, and never a contract or a stage provider. A semantic helper imports foundations and its own syntax and vocabulary, and may reparse with its own parser under an exception. Both have their role only when `syntax/<language>/go` exists |
| Composition | `integration/go`, `verification/integration/go`, `tools/**` | Chooses and registers contributions. The two test compositions are not classified. A package under `tools/` needs no exception to import a stage provider or a query service. It is still classified, by an explicit entry or as an extractor |

**One language.** A package that names a language imports no package of
another language. A mapping names two. A package belongs to a language by its
path, or by the table `languageOwners` where its path gives no role. No entry
and no extension lifts this rule. The verification laws are outside it: the
rule spells out their pair, because the directory says `go` where the language
is `golang`.

**Explicit entries.** A role is written once for every language, so a new
package of an existing role needs no edit to the test. Two groups keep one
explicit entry per package, and a new package of either group is an edit:

- **The foundation and kernel set, `foundationRules`.** These are the first
  three rows of the table above. Their layering among each other is a fact per
  package, and one rule for "a foundation" would lose it. The set is closed.
  The paths in those rows describe its members and are not patterns: a new
  `contracts/<name>/go` is unclassified until it has an entry.
- **Packages of which there is only one, `oneOfAKindRules`.** These are the
  private horos parts `internal/declsyntax`, `internal/horosmodel` and
  `internal/horosfiles`, the path and import helper of TypeScript Files
  formatting `internal/typescriptimports`, the projection
  `verification/files/source/horos/go`, the horos-to-go verification pair
  `verification/files/horos-to-go/go`, the logos adapter
  `verification/logos/go`, the reporting-tool packages
  `tools/provider-matrix/go`, `report/go`, `capability/go` and
  `observation/go`, and the command `tools/horos-module/go`. The report host
  never imports a concrete extractor. The verification pair lists the Go query
  service and never the Go analysis provider. It imports the query service
  under an exception that names `Admission`.

**Role extensions.** Languages differ. The table `roleExtensions` names, per
package, the imports that one language needs and its role does not give. Ten
packages have one: the six horos stage providers, three TypeScript packages
and the horos declaration recipe. An extension must be used, must not repeat
the role and cannot lift a prohibition of the role.

**No further hand entries.** `TestProductionArchitecture` holds no entry of
its own. Eight were held there for Go until pull request
[#246](https://github.com/Bitspark/bittype/pull/246) converted them.
`semantics/golang/go`, `internal/golangsem`,
`resolve/files/golang/provider/go`, `analyze/files/golang/provider/go` and
`analyze/files/rust/provider/go` take their role. `internal/typescriptimports`,
`verification/files/horos-to-go/go` and `verification/logos/go` are entries of
`oneOfAKindRules`. Where a package has both a role by its path and a hand
entry, the entry replaces the role's allow-list and none of its prohibitions.
No package has both today.

A future `cmd/**` composition is not known to the test. Adding it is an edit.

The parse, check, format and print interfaces are independent; their configured
adapters use representations for owned heterogeneous inputs/outputs. The engine
registers them separately after identity/profile registration. Representation
depends on generic carrier checks and diagnostic records, never a stage-provider
contract. Source input is `parse.Input`; `text` contains diagnostics only. The shared implementation remains under `print/syntax/shared/provider/go`, separate
from `contracts/print/go`. Imports may use a local alias when both packages are
needed. No compatibility interfaces or input aliases remain under text.

The [reporting tools](tools/provider-matrix/README.md) are outside the runtime
dependency graph. The public host selects one module using this repository
layout; dependency export data does not select dependency annotations or docs.
The command generates compact matrices and static details together, while
execution observations belong to per-invocation artifacts. Independent extractors
use the same fixed structural contract. This is not a second language registry,
an open semantic payload model or a sandbox for hostile extractor code. The
[capability report design](docs/capability-reports.md) specifies those boundaries.

The structurally checked AST/CST and language descriptors share one home in
`trees/representation/go`. `contracts/transform/go` owns ASTTransformer and
ErrInvalidTransformer; representation owns tree and identity errors. Its narrow
CheckMarker operation checks actual Go type identity without exposing payloads
or providing unchecked construction. No aliases remain at the former paths or
under engine for these moved types.

Every private helper also has a dependency rule. A shared helper cannot conceal
a native-language import, and a language helper cannot depend on another
language. A lexical helper depends on no stage. A semantic helper may import
its own language's parser under an exception. This keeps the rule enforceable
through indirect imports.

## Enforced consistency

The [architecture checks](integration/go/architecture_test.go) discover the
union of language directories across syntax, Syntax stages and Files semantic stages
and require each language to have every stage in its tested capability profile.
The four native languages retain all five native text stages; Horos requires
the same Syntax stages plus Files resolve/analyze, with explicit qualification.
Generic carrier check directories
are excluded explicitly. An orphan parser and an incomplete syntax package both
fail. The checks also classify every production package and enforce its directed
dependencies, including private helpers and narrow external-model allowances
for the Horos frontend/text stages, source reader, module encoder and concrete
verification integrations. The exact package allowances live in the test. Generic carriers,
contracts, engine and printer must remain independent of Horos/ontos, and native
syntax/text stages must remain independent of the canonical frontend.

**Implemented. State of 12 November 2026.**
[Package boundaries](docs/package-boundaries.md) records the rules: one
contribution per package, identities in vocabulary packages, only a
composition imports a contribution, and import rules stated per role. The plan
is [#218](https://github.com/Bitspark/bittype/issues/218). Its last step, the
rules by role, merged with
[#244](https://github.com/Bitspark/bittype/pull/244). Go's contract identities
and query service moved afterwards, with pull request
[#246](https://github.com/Bitspark/bittype/pull/246). The horos names that
a writer and a separate reader each spelled moved into `semantics/horos/go`
with pull requests [#247](https://github.com/Bitspark/bittype/pull/247) and
[#248](https://github.com/Bitspark/bittype/pull/248), and the tag
prefixes of canonical horos values with
[#250](https://github.com/Bitspark/bittype/pull/250). The prefixes moved on
to `syntax/horos/go` with
[#285](https://github.com/Bitspark/bittype/pull/285), with the labels, the
field tag and the record words: names that the horos tree carries or is
lowered to. The horos syntax node
kinds have one declaration since
[#254](https://github.com/Bitspark/bittype/pull/254), and
`semantics/typescript/go` publishes the TypeScript graph names that Go reads
since [#255](https://github.com/Bitspark/bittype/pull/255). Every confirmed
instance of
[bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md) is
corrected. The check is the function
`checkArchitecture` in
[`architecture_check_test.go`](integration/go/architecture_check_test.go). It
does the following:

- **Rules by role.** A package gets its import rule from the role that its
  path gives, as the table under
  [dependency direction](#dependency-direction) summarises. At `c8d0344`, the merge
  of pull request #246, 84 of the 127 production packages are
  classified this way, ten of them with a role extension. The other 43 have an
  explicit entry: 31 in the foundation and kernel table and 12 one-of-a-kind
  packages. A production package with neither a role nor an entry fails as
  unclassified. A new language, mapping, common mapping part, verification
  law, extractor or language helper needs no edit to the rules. An import that
  its role does not
  give needs a role extension. A new language whose parser and formatter call
  its own checker still adds two entries to the exception table. A new stage
  or dimension gives no role and is an edit.
- **Every allowance is used.** An explicit entry and an extension are checked
  per package: an import that they allow and the package does not make fails.
  A role is checked per shape: each import it permits must be imported by at
  least one package of that role. The external model allowances are checked
  the same way.
- **Symbols of a provider or query service import.** The check reads which
  exported symbols a production package uses from each module package it
  imports. An import of a stage provider package
  (`<stage>/<dimension>/<language>/provider/go`) or of a query service package
  (`analyze/files/<language>/queries/go`) by a package that is not a
  composition needs an entry in the
  [exception table](integration/go/architecture_exceptions_test.go). The entry
  names the importer, the imported package, the exact symbols and the reason.
  A use of another symbol fails. An entry or a symbol that is not used fails.
  The table holds fourteen entries. All are accepted. In thirteen, a parser, a
  formatter and the horos resolver call their own checker, and the Rust and
  TypeScript semantic helpers call their own parser. In one, the horos-to-go
  verification pair takes `Admission` from the Go query service.
- **All files of a package.** The check also reads the files that the build
  constraints of the test run exclude: a file behind a build tag, a file of
  another operating system or architecture, and a cgo file when cgo is off.
  Test files stay out. An import of a module package that `go list` does not
  report fails, because nothing checks that package.
- **A private part belongs to its language.** A role grants
  `internal/<language>` and `internal/<language>sem` only where the imported
  package is that language's helper by its own role. A package cannot take
  another owner's private part by choosing a name.
- **Dot imports.** A dot import of a package of this module fails, because it
  hides the symbols used.

The compositions `integration/go` and `verification/integration/go` are not
classified. `tools/**` needs no exception to import a provider or a query
service.

The plan in #218 is complete. The profile gate that
[#238](https://github.com/Bitspark/bittype/pull/238) left in the adapter at
`transform/syntax/horos/mapping/go` is removed since 12 November 2026
([#216](https://github.com/Bitspark/bittype/issues/216),
[bittype/ACX-008](docs/architecture-cases/0008-structural-profile-gate.md)).
A mapping states its profile in its own description. That package keeps the
two profile identities and a function type that writes no profile and refuses
none. This remains open and belongs to other work:

- **Optional fields in Record.** Record still decides how an optional field is
  represented in Go, Rust and C++. What a mapping's `Expand` returns is
  [#196](https://github.com/Bitspark/bittype/issues/196), item 6.
- **Two candidates.** They are not confirmed and are listed in
  [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md#candidates).
- **The `Self` rule in the shared naming function.** Answered in principle on
  12 November 2026: a reserved word is the target's. It is not corrected.
  [#274](https://github.com/Bitspark/bittype/issues/274) delivers it with
  [#268](https://github.com/Bitspark/bittype/issues/268), as the
  [naming design](docs/naming.md#reserved-words) says.
- **Names that were not examined.**
  [#251](https://github.com/Bitspark/bittype/issues/251) examined the names
  that the native compiler helpers write and the horos node kinds that two
  packages declared, and both are corrected. What it left unexamined has no
  issue yet and is listed in the
  [limits of bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md#limits):
  what the TypeScript helper reads, its response status, the TypeScript
  generation form names, the Go and Rust vocabulary packages' own readers, and
  the admission grammars of the native checkers.

A green run shows that the forbidden imports are absent. It does not show that
a package holds one contribution, that a decision made inside a package is
right, or that each package of a role needs every import its role permits.
[What the check cannot establish](docs/package-boundaries.md#what-the-check-cannot-establish)
lists the limits.

The findings and their state are in
[bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md).

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
in `verification/logos/go`. General contracts retain agnostic association.
The implemented paths appear in the current hierarchy above. A wider verification
provider registry and additional language interpretations remain deferred.

The [agnostic semantic model](docs/semantics.md), [analysis contract](docs/analysis.md)
and first Horos analyzer are implemented. Additional analyzers use the same
optional analyze stage with an explicitly tested capability profile. Native
semantic analysis is not implied by native text tooling.

The [Horos declaration integration](docs/horos.md) implements
`syntax/horos/go`, `check/syntax/horos/provider/go`, `analyze/files/horos/provider/go` and
`resolve/files/horos/provider/go`, with private projection/correspondence in `internal/horosmodel`.
The horos semantic identities are published by `semantics/horos/go`.
Readable syntax now implements `parse/syntax/horos/provider/go` and
`format/syntax/horos/provider/go` under the same stage matrix. These replace the former proposed
`canonical` paths; there are no compatibility aliases.

Directed pagination mappings live under `transform/syntax/horos-to-<native>/mapping/<name>/go`, one package per mapping. A partial reverse implementation would use the reverse ordered pair. Application
composition and IO are proposed under `cmd/bittype/go`. Horos owns the declaration semantics and
encoding; these packages supply tooling projections and native representations.
New dependency rules and qualification accompany further implementation, with no
empty packages created now. All five current languages have qualified text stages;
new AST-only providers still require an explicit tested capability profile.

Adding a native language supplies its implemented stages, any private lexical helpers,
conformance fixtures and compiler qualification. It requires no change to the
generic carriers, engine or production registration catalogue: there is no such
catalogue. Applications choose their registrations. Adding Function, Stream,
Cell or another library declaration does not add a stage, language or core kind.

The [unified declaration model](docs/horos.md) uses one `horos` family under
syntax/check/resolve/analyze. Private syntax framing remains in
`internal/declsyntax`; model projection, correspondence and export remain in
`internal/horosmodel`. Both take the node kinds, attributes and roles from
`syntax/horos/go`, which declares them once. No `horoswire` stage or compatibility alias remains.
Bounds are explicit syntax and provider-owned semantic facts, not kernel kinds.

Model dependency allowances cover the concrete owners listed above, including
Horos source reading and verification. They do not extend to agnostic contracts.
Architecture checks preserve the agnostic and native-stage boundaries. Native
pagination mappings use `transform/syntax/horos-to-<native>/mapping/<name>/go`.
Template lifting and its exclusive dependency allowances remain removed.

The [native representation design](docs/type-representations.md) and
[type-owned generator plugin proposal](docs/type-generation-plugins.md) place
builtin and library mappings behind one supplied contract. Each type library owns
its mappings and type artifacts through the same SPI; repository placement is
separate from semantic ownership. The consumer selects
defaults and explicit overrides. Adding Cell or another ordinary contract adds no
Bittype stage, engine edge per contract or central registration entry. The first
[public SPI and guarded invocation](docs/generation-spi.md) are implemented in
`contracts/generation/go`. [Automatic composition](docs/generation-composition.md)
now drives all four native generators. The private pagination Dialects are deleted.

The contract package owns native descriptions, scoped planning and the sealed
plan/output. The registry and concrete integrations follow the stage-first layout:

```text
contracts/generation/go/                  descriptions, invariant planner and output
transform/syntax/shared/generation/go/      registry, selection and composition
transform/syntax/horos/source/go/           receipt-backed source reader
transform/syntax/horos/mapping/go/          profile identities and the mapping library function type
transform/syntax/horos/mapping/{pattern,declarations,naming}/go/  source patterns, declaration recipe, naming
transform/syntax/shared/target/go/          shared native construction and scoped symbol storage
transform/syntax/<native>/model/go/         intended native forms
transform/syntax/<native>/target/go/        native intent to AST construction
transform/syntax/horos-to-<native>/mapping/<name>/go/  independently selected mappings, one package each
```

Consumers register source, target and mapping libraries, then request source roots and an output identity. The kernel plans fixed native intent; Run constructs ASTs without invoking a mapping. The [construction SPI](docs/kernel-composition.md) gives the public interfaces and a Horos to TypeScript example.

The [naming design](docs/naming.md) is accepted and not implemented. Today the `naming` and `declarations` packages above spell every target's names. Under the design generation records a naming intent, and each language owns its realization: its reserved words, in Go the case of a name, and in its Files stage the collisions of a placed scope. The design names one new package, a Go-owned one for Go's check of a spelled name, which [#274](https://github.com/Bitspark/bittype/issues/274) adds. No path in this layout has changed.

The [contribution contract](docs/contribution-contract.md) is a design of 13 November 2026 from research 0017 and is not implemented. It changes what the mapping packages above register and return, in one migration behind a gate ([#196](https://github.com/Bitspark/bittype/issues/196), item 6). It names one new package role, the payload codec, which encodes, validates and freezes the value a rule returns. The role has a shared part and one codec package for each payload schema, and decision O37 on #196 divides them: a codec may import the ontos value packages and reaches horos only through the shared part. That role lies outside `contracts/generation/go`, which imports neither horos nor ontos, and outside the horos source reader. Its path and its import allowances are not decided: the migration edits this layout. The design also puts the binding facets of a mapping contribution into facet packages, and the contribution's identity with its generated payload types into a vocabulary package of the family for the target. This layout has no role for a facet package of a mapping contribution, and none for a vocabulary package of one family for one target: its role "Vocabulary" is per language. The nearest existing role is the verification law in the role table above, whose packages only tests and test compositions import. A new role is an edit of that table and of the architecture test, made with the prototype's report. The prototype of that design is done on a branch that never merges. It added roles and paths for these packages to its architecture test: a vocabulary of a contribution, a payload codec, a facet package, the vocabulary and the backend of a completeness product under a top-level directory `completeness/`, and two shared packages. They are a proposal, listed in [the contribution contract](docs/contribution-contract.md#the-roles-the-prototype-proposes). One point of the proposal is already taken back: it puts the shared part of the payload codec role under `transform/syntax/shared/`, which is a kernel path in the role table above, with imports of horos and ontos. Generic infrastructure stays independent of those packages, so that part does not live under a kernel path (decision O41 on #196); the migration names its place. The migration edits the role table above, and its review judges the proposal. No path in this layout has changed.

The [Syntax and Files design](docs/source.md) is implemented with
`Syntax : CST : AST` and `Files : Source : Semantics`. Both dimension names are
explicit in implementation paths. Horos composes the full Files text pipeline.
Go, TypeScript and Rust supply bounded Files resolution and analysis. TypeScript
also supplies [Files formatting with named type imports](docs/typescript-type-imports.md).
Applications construct a Files printer for any of the five languages from the shared printer providers.

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
