# Working here as an agent

## Owner delivery rule, 4 November 2026

No backups, checksum/validation/check workflows, or provenance workflows are
implicit task or delivery prerequisites. Deletion means deletion within the
authorized scope; do not replace it with retained copies or archives. Do not claim
unperformed checks passed. This owner rule supersedes conflicting older delivery
instructions in this file. It does not remove the requested semantic-analysis
behavior of the library. The owner clarified that ordinary automated software
tests and typechecking remain required for affected changes; a build alone is
not a test result. Use the relevant existing checks and repository resource
rules. Documentation changes require local-link and whitespace checks under
the platform contribution instructions, not a new validation workflow.

## Architectural correctness: standing owner rule, 6 November 2026

Architecture is part of correctness from the first bounded implementation.
Passing tests and working examples do not discharge semantic ownership or
additive-composition requirements. Shared mechanisms own invariant guarantees;
language and mapping meanings enter through owner-controlled contributions.
Within the declared extension boundary, adding a case must change only its
implementation and consumer composition, not shared orchestration, closed case
types or unrelated implementations. Package isolation and a "plugin" label do
not establish this. An extension must preserve scope, identity, association,
authority and failure guarantees; arbitrary payloads or unchecked callbacks do
not suffice. Derive from domain concepts and composition before adding abstractions.

For interface design and redesign, derive the best public surface from domain
requirements, ownership and guarantees. Incumbent signatures, available code
and implementation or migration effort must have zero influence on the design.
Decide how to implement the selected contract afterward.
[AC-09](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-09).

Read the [general rule and its rationale](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md).
Reason about responsibilities and change propagation; do not replace this with
keyword heuristics, scoring or new lint gates. Limited coverage may be deferred;
a violated extension boundary must be identified as an architectural defect,
not described as merely future coverage or reported architecturally complete.

**Owner direction, 11 November 2026:** "if something is placed incorrectly, it is
incorrect and must change. Something that is missing can be healed additively,
i.e. there is nothing wrong requiring removal." Correct a misplaced or wrongly
named part first, as its own change, and add missing coverage afterwards. Defer
such a move only when its destination needs a huge effort, and record the
deferral with its reason and trigger as an open violation.
[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10).

**Owner request, 11 November 2026:** check for packages whose parts "neither
conceptually nor practically depend on each other", and for parts that call one
another directly "instead of invoking the kernel", and enforce the result
"through forbidden dependencies that we could check for and enforce". Read
[package boundaries](docs/package-boundaries.md) before adding a package or an
import between implementations. One contribution, meaning one thing a consumer
registers or selects on its own, is one package. Identities live in vocabulary
packages. Only a composition imports a contribution. Ask of every direct call
whose answer the caller needs: its own, or the composition's. The
[architecture test](integration/go/architecture_check_test.go) enforces part of
this since 11 November 2026. An import of a stage provider package outside a
composition needs an entry that names its symbols and its reason. Mapping
contributions, verification laws, language vocabularies and query services are
classified by their path. Every other package is still classified by a
per-package list. The plan is [#218](https://github.com/Bitspark/bittype/issues/218); the findings that remain open are
[bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md). The
import check enforces boundaries that this reasoning has decided. It does not
replace the reasoning, and a green run does not show that ownership is right.

The [Horos/Go verification path](docs/verification.md#library-owned-verification-contributions)
takes library-owned interpretations, mapping projections and laws. Preserve the
captured subject, receiver authority, dependency policy and Logos outcome when
adding a contribution.

Reference the [numbered architecture rules](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#numbered-rules)
when explaining architectural requirements or violations. Maintain a numbered
[case record](docs/architecture-cases/README.md)
for an identified violation and update it with the actual correction and its
qualification when fixed, linking the owning design. The current verification
case is [bittype/ACX-001](docs/architecture-cases/0001-verification-composition.md).
Maintain the record when a change affects its qualified extension boundary.

Read [CHARTER.md](CHARTER.md), [docs/model.md](docs/model.md),
[docs/ast.md](docs/ast.md), [docs/syntax.md](docs/syntax.md),
[docs/generation.md](docs/generation.md) and
[LAYOUT.md](LAYOUT.md) before changing this tree.

## Current scope

The owner selected canonical definitions as input and TypeScript, Go, Rust and C++
types as outputs. Horos owns data descriptors and canonical
module/generic/application semantics with general structural kinds; bittype owns syntax,
AST/CST tooling and transformations. The Go declaration model and Horos text/semantic frontend are implemented;
read [docs/horos.md](docs/horos.md). Do not revive
the superseded independent bittype value-type algebra. The current native-text
slice implements parsers, AST-to-CST formatters, syntax checks and a generic
CST-to-text printer. Read the [concrete syntax design](docs/parsing-formatting.md).

The [agnostic semantic model](docs/semantics.md), [analysis SPI](docs/analysis.md)
and [resolver SPI](docs/resolution.md) are implemented, with Horos and bounded Go/TypeScript/Rust providers.
Keep semantic rules provider-owned. Shared snapshots, requests/results and reports
live in contracts/semantic/go; no aliases remain at the old analysis paths.
Resolution and analysis are independently registered. Analysis consumes an explicit
compatible receipt; never invoke a hidden resolver or bypass snapshot identity.
The [provider matrix](docs/providers.md) is generated from source by
`tools/provider-matrix/go`; its workflow updates the two marked regions.
Do not hand-edit cells or add a capability catalogue or runtime registry.
Dimensions come from `bittype:dimension` comments on boundary type declarations;
`bittype:shared` explicitly marks shared contracts. Keep these source-owned and
do not restore package-name classification or add runtime metadata methods.
Syntax admission must accept unresolved names; analysis publishes request-bound,
immutable facts with explicit completion, verdict and model availability.
Horos, Go, TypeScript and Rust supply registered Files resolution/analysis providers. Read
[the native semantic-provider design](docs/native-semantics.md). The bounded Go
provider uses `go/types`; native queries are explicit and the verifier consumes
its immutable graph. TypeScript uses an explicitly configured, pinned Compiler API
helper over captured source with verified AST correspondence; analysis reuses the
resolver receipt. Rust uses an explicitly configured rustc_driver helper over captured source. C++ analysis remains unimplemented.

The owner's 4 November [Syntax and Files symmetry](docs/source.md) replaces the
separate layout proposal: `Syntax : CST : AST` and `Files : Source : Semantics`.
Use parse/format/print roles at both levels; do not revive a layout.Model,
layout.Provider or combined Source provider. The typed carriers, generalized role interfaces, Horos Files compositions and
multi-unit output are implemented. Preserve exact capture/receipt association and
keep native Files semantics unavailable until their providers exist. No aliases.

The owner's 4 November [unified declaration decision](docs/unified-declarations.md)
supersedes separate Horos and Horos/Wire languages. One Horos revision-4 grammar
has ordinary kind definitions and general contracts; `Page<T>` is authored once and
`DataPage<D: h.Data> = Page<D>` is an ordinary constrained alias using an explicit
prelude import. Read [general kinds](docs/kinds.md). Remove the old
profiles, packages, wrapper syntax and readers without compatibility aliases.
Ground Ontos/Horos judgments remain data-only; optional Ontos/wire still owns
mixed *values*. Hydration stays in Bitwire/Bitruntime. Preserve the
[concept inventory](docs/wire-extension-preservation.md) and library ownership.
The one [Horos frontend](docs/horos.md) projects explicit bounds through existing
provider-owned semantic terms. The agnostic kernel interprets none of them.
Template lifting remains removed; no second declaration or derived Page exists.
Do not restore Claim/Interface, DataSummary, a dedicated Data proof path, or
reserved Type/Data names. Scalar, List, Record, Data and user-defined kinds use
the same Horos machinery. Kind definitions are compile-time entities, not native
type exports. Keep semantic sort separate from declaration parameter count.

The [declarative generation SPI](docs/generation-spi.md) and invariant planner
live in contracts/generation/go. The [registry](docs/kernel-composition.md)
lives at transform/syntax/shared/generation/go. Register source, target and mapping
libraries independently. Libraries describe native types and scoped child demands;
the planner resolves them into a closed native program. Target construction
receives that program without a source reader or mapping callback. Do not restore
the procedural operations, association-product tickets, per-pair assembly,
per-request mapping tables or compatibility aliases.

Builtin and independent library mappings use the same contract. Exact overrides
identify both expected and replacement bindings; ambiguity and selected-rule
failure never trigger a weaker fallback. Native semantic forms belong under
transform/syntax/<native>/model/go, their AST interpretation under
transform/syntax/<native>/target/go. The horos mapping adapter lives in
transform/syntax/horos/mapping/go; its source patterns, declaration recipe and
naming rule are the packages pattern, declarations and naming below it. Each
mapping is one package under transform/syntax/horos-to-<native>/mapping/<name>/go
and exports Library.
The kernel, source and target packages cannot import concrete mapping packages.
Mapping packages depend on native models, never target AST constructors.

Preserve real binder identity and one convention for native formals and complete
applications. Admit every source argument even if native mapping omits it.
Native helpers use owner-qualified keys and their own parameters, never fabricated
source declarations or inherited caller binders. Scoped expansion requests expire
on return. A run uses fresh authenticated symbols and cannot select new mappings.
Source admission, native syntax and representation meaning remain separately
owned; no agnostic mechanism infers their semantics.

TypeScript generation carries exposures and external occurrences in the checked
generation.Output. Its [Files formatter](docs/typescript-type-imports.md) decides
imports and exports after placement and checks fresh native bindings. Formatting
already analyzed files retains its distinct input. The
[three SPI changes](docs/generation-planner.md#three-spi-changes) are implemented:
the kernel walks the admitted semantic graph, and sources supply interpretation,
judgments and a reader for their mapping libraries. Do not restore source-owned
navigation (Open/Admit/CheckChild or reader records).
Cross-language verification and native admission enforcement remain separate.
Use [GEN-02](docs/generation-gaps.md#gen-02) when describing unfinished
generation work, and [GEN-01](docs/generation-gaps.md#gen-01) for the delivered
traversal correction and its supported scope. Update those statuses and the
owning architecture case when a correction lands; working examples do not
close a gap.

- Readable Horos syntax and native pagination transformers are implemented;
  read [Horos syntax](docs/horos-syntax.md) and [semantic generation](docs/semantic-generation.md).
  Public builders and a general application CLI remain deferred. The
  `tools/horos-module/go` command lowers a single Horos source module to canonical
  Ontos bytes through the existing parser/private encoder, without resolution.
  It is a Horos-specific composition tool, never a prelude-aware core stage.
  The syntax engine and ASTTransformer boundary are now
  implemented; read [docs/engine.md](docs/engine.md). Do not introduce runtime
  adapters, callable behavior,
  application-value validators, serialization or networking.
- Function, Stream, Cell and future behavioral contracts remain library-owned.
  Never add their names as core kinds, switches, imports, registrations,
  fingerprints or hidden defaults.
- Preserve generic binders/applications, declaration references and origins.
  Equal-looking syntax does not authorize substituting library contracts.
- Do not make TypeScript extraction the authoritative canonical input.
- Follow the stage-first hierarchy in LAYOUT.md. Shared trees live under trees,
  shared contracts under contracts; native implementations use the same
  syntax and dimension-qualified parse/check/format/provider matrix, with explicitly qualified
  optional resolve/analyze stages. Do not add compatibility
  aliases or a language-first packages tree. Keep one public root Go module
  independent of private dependencies. The concrete verification composition
  consumes the public Logos contract, as documented in
  [the verification architecture](docs/type-system-architecture.md#proposed-package-roles).
- Implementations use stage/dimension/scope/component/go. Keep shared printers at print/{syntax,files}/shared/provider/go. Vocabularies use syntax/language/go. Construct stage adapters through their individual contracts. The [layout](LAYOUT.md) gives the implemented paths.
- Add every language to the shared conformance fixtures; AST-only providers
  require an explicit tested capability profile. Architecture checks
  discover stage directories and enforce complete coverage and directed imports.
  Keep package tests with their implementation and cross-package qualification
  under integration. Private helpers must obey the same ownership boundaries.
- Implement the toolchain in Go in component-owned packages in the root module.
  Formatters construct CST tokens/trivia; only the printer emits source text.
  Do not print source and reparse it as an intermediate AST/CST.
- AST, pipeline and language identity remain independent. CST depends on AST.
  Syntax depends on AST/CST carriers and identity, never on another language
  or parser/formatter. Native stages own their grammar and spelling.
- language.Provider supplies ID and Name for configured language providers.
  Name capability interfaces through their package: parse.Provider,
  check.Provider, format.Provider, print.Provider, resolution.Provider and
  analysis.Provider. Each stage contract is independent of text composition;
  parse.Input owns source input. Register parse/check/format/print capabilities
  independently through their configured adapters. Language registration captures
  identity/profile/marker only. Pack/PackCST check generic carrier structure;
  explicit checkers admit syntax for Check, semantic Capture and every transform
  endpoint. AST-only providers need no text capability. text contains diagnostics
  only; do not restore combined provider registrations or aliases. Avoid redundant
  package prefixes or Provider suffixes; retain no superseded interface aliases.
  Transformers identify two endpoints through FromID/ToID. Registration stays
  explicit and agnostic, with no known-language list,
  name switch or automatic registration. No compatibility identity alias.
- Registry profiles are explicit caller-owned compatibility promises. Check
  identity, profile and actual Go marker at dynamic boundaries; no AST[any]
  conversion, structural fingerprint or inferred contract identity.
- Keep route planning separate from semantic applicability. Reject ambiguous
  shortest routes and unsupported inputs; never silently retry a weaker mapping.
  Inputs and output envelopes own snapshots. Call plugins outside registry locks.
- ASTView is a borrowed read-only view. A formatter owns a deep AST snapshot in
  its output CST. Parsers retain all admitted source bytes; whole-input rejection,
  bounded resources and errors without partial output are required.
- The printer checks carrier structure only. It never interprets language kinds,
  chooses layout, escapes tokens, or resolves imports.
- Syntax checks are not complete compiler/type checking. Unsupported constructs
  fail explicitly; no raw-source escape node or skipped executable body.
- Native compiler qualification and source/AST round trips are separate evidence.
  Text cannot recreate external resolved identities or original provenance.
- Keep prior research and Nightseam as evidence, not competing specifications.
  Use Horos's declaration specification for normalization/identity and encoding;
  do not infer semantic equivalence or behavioral ownership from document pins.
- Only the Horos frontend (including its text stages), its private model helpers,
  source-specific transformation readers and the concrete Horos verification
  projections/composition consume external Horos/Ontos model dependencies.
  Only Logos-aware verification packages import the public Logos contract. The joint checker invokes it. Selected mapping libraries supply proof syntax.
  Preserve the narrow package import allowances; generic infrastructure and
  native stages must remain independent of these semantic packages.

## Delivery

Use a task branch or worktree and a pull request. Squash the task change onto
`main`; keep the designated primary checkout on `main`.

Do not add documentation checks, CI waiting, independent review or verification
records as delivery prerequisites under the owner rule above. Existing workflow
definitions remain unless their change is in scope.

Keep historical commits, research evidence, license notices and unrelated work.
Do not add private checkout dependencies, local coordination state or credentials.
Repository names are lowercase; the company name is Bitspark.
