# Working here as an agent

## Owner delivery rule, 4 November 2026

No backups, checksum/validation/check workflows, or provenance workflows are
implicit task or delivery prerequisites. Deletion means deletion within the
authorized scope; do not replace it with retained copies or archives. Do not claim
unperformed checks passed. This owner rule supersedes conflicting older delivery
instructions in this file. It does not remove the requested semantic-analysis
behavior of the library.

Read [CHARTER.md](CHARTER.md), [docs/model.md](docs/model.md),
[docs/ast.md](docs/ast.md), [docs/syntax.md](docs/syntax.md),
[docs/generation.md](docs/generation.md) and
[LAYOUT.md](LAYOUT.md) before changing this tree.

## Current scope

The owner selected canonical definitions as input and TypeScript, Go, Rust and C++
types as outputs. Horos owns data descriptors and canonical
module/generic/application semantics with Type/Data bounds; bittype owns syntax,
AST/CST tooling and transformations. The Go declaration model and Horos text/semantic frontend are implemented;
read [docs/horos.md](docs/horos.md). Do not revive
the superseded independent bittype value-type algebra. The current native-text
slice implements parsers, AST-to-CST formatters, syntax checks and a generic
CST-to-text printer. Read the [concrete syntax design](docs/parsing-formatting.md).

The [agnostic semantic model](docs/semantics.md), [analysis SPI](docs/analysis.md)
and [resolver SPI](docs/resolution.md) are implemented, with Horos providers.
Keep semantic rules provider-owned. Shared snapshots, requests/results and reports
live in contracts/semantic/go; no aliases remain at the old analysis paths.
Resolution and analysis are independently registered. Analysis consumes an explicit
compatible receipt; never invoke a hidden resolver or bypass snapshot identity.
Keep the [provider matrix](docs/providers.md) synchronized with capability tests.
Syntax admission must accept unresolved names; analysis publishes request-bound,
immutable facts with explicit completion, verdict and model availability.
Only Horos implements production semantics; native analysis remains deferred.

The owner's 4 November [Syntax and Files symmetry](docs/source.md) replaces the
separate layout proposal: `Syntax : CST : AST` and `Files : Source : Semantics`.
Use parse/format/print roles at both levels; do not revive a layout.Model,
layout.Provider or combined Source provider. The typed carriers, generalized role interfaces, Horos Files compositions and
multi-unit output are implemented. Preserve exact capture/receipt association and
keep native Files semantics unavailable until their providers exist. No aliases.

The owner's 4 November [unified declaration decision](docs/unified-declarations.md)
supersedes separate Horos and Horos/Wire languages. One Horos revision-2 grammar
has explicit Type/Data bounds; `Page<T>` is authored once and
`DataPage<D: Data> = Page<D>` is an ordinary constrained alias. Remove the old
profiles, packages, wrapper syntax and readers without compatibility aliases.
Ground Ontos/Horos judgments remain data-only; optional Ontos/wire still owns
mixed *values*. Hydration stays in Bitwire/Bitruntime. Preserve the
[concept inventory](docs/wire-extension-preservation.md) and library ownership.
The one [Horos frontend](docs/horos.md) projects explicit bounds through existing
provider-owned semantic terms. The agnostic kernel interprets none of them.
Template lifting remains removed; no second declaration or derived Page exists.

The [native representation design](docs/type-representations.md) remains proposed.
Libraries own optional recipes, selected by exact declaration identity in consuming
build configuration. Concrete transformers consume them uniformly; no
contract-specific engine registrations, named-library imports in general
transformers or raw-source recipes. Missing/conflicting required representations
fail rather than selecting a weaker fallback.

- Readable Horos syntax and native pagination transformers are implemented;
  read [Horos syntax](docs/horos-syntax.md) and [semantic generation](docs/semantic-generation.md).
  Public builders and application CLI remain deferred. The syntax engine and ASTTransformer boundary are now
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
  syntax/parse/check/format/provider language matrix, with explicitly qualified
  optional resolve/analyze stages. Do not add compatibility
  aliases or a language-first packages tree. Keep the root Go module.
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
- Only the Horos frontend (including its text stages), its private model helpers
  and source-specific transformation readers consume external Horos/Ontos model dependencies.
  Preserve the narrow package import allowances; generic infrastructure and
  native stages must remain independent of these semantic packages.

## Delivery

Use a task branch or worktree and a pull request. Squash the task change onto
`main`; keep the designated primary checkout on `main`.

Do not add documentation checks, CI waiting, independent review or verification
records as delivery prerequisites under the owner rule above. Existing workflow
definitions remain unless their change is in scope. Portable workloads belong on the shared build service.

Keep historical commits, research evidence, license notices and unrelated work.
Do not add private checkout dependencies, local coordination state or credentials.
Repository names are lowercase; the company name is Bitspark.
