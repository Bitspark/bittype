# Working here as an agent

Read [CHARTER.md](CHARTER.md), [docs/model.md](docs/model.md),
[docs/ast.md](docs/ast.md), [docs/syntax.md](docs/syntax.md),
[docs/generation.md](docs/generation.md) and
[LAYOUT.md](LAYOUT.md) before changing this tree.

## Current scope

The owner selected canonical definitions as input and TypeScript, Go, Rust and C++
types as outputs. The canonical model remains a proposal. The current native-text
slice implements parsers, AST-to-CST formatters, syntax checks and a generic
CST-to-text printer. Read the [concrete syntax design](docs/parsing-formatting.md).

- Canonical AST, actual language transformers, public builders and application
  CLI remain deferred. The syntax engine and ASTTransformer boundary are now
  implemented; read [docs/engine.md](docs/engine.md). Do not introduce runtime adapters, callable behavior,
  application-value validators, serialization or networking.
- Function, Stream, Cell and future behavioral contracts remain library-owned.
  Never add their names as core kinds, switches, imports, registrations,
  fingerprints or hidden defaults.
- Preserve generic binders/applications, declaration references and origins.
  Equal-looking syntax does not authorize substituting library contracts.
- Do not make TypeScript extraction the authoritative canonical input.
- Implement the toolchain in Go in component-owned packages in the root module.
  Formatters construct CST tokens/trivia; only the printer emits source text.
  Do not print source and reparse it as an intermediate AST/CST.
- AST, pipeline and language identity remain independent. CST depends on AST.
  Syntax depends on AST/CST carriers and identity, never on another language
  or parser/formatter. Native stages own their grammar and spelling.
- LanguageProvider supplies ID and Name for configured language providers.
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
  Settle canonical normalization/identity before fixing an encoding.

## Delivery

Use a task branch or worktree and a pull request. Squash a green change onto
`main`; keep the designated primary checkout on `main`.

For documentation, check relative Markdown links, JSON example syntax and
`git diff --check`, including staged additions. The documentation workflow runs
these checks on the exact committed revision on a GitHub runner. Portable
compilation and later language checks also run on the shared build service or
GitHub runners, using the maintained repository gate when one exists.
The Go workflow checks the exact commit with formatting, contract/dependency
tests and vet. Use that remote gate; do not run Go workloads locally. Maintain
independent consumer examples and negative type-check cases for API boundaries.

Keep historical commits, research evidence, license notices and unrelated work.
Do not add private checkout dependencies, local coordination state or credentials.
Repository names are lowercase; the company name is Bitspark.
