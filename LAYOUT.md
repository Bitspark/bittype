# Repository layout

**Owner decision, 31 October 2026:** organize by processing stage and enforce a
uniform structure across languages. The stage-by-language matrix is deliberate:
the same responsibilities occupy the same paths, making omissions and drift
visible. There is one root Go module, with independently owned packages.

## Current hierarchy

```text
go.mod
trees/
  ast/go/                 generic AST, nodes, addresses, origins and borrowed view
  cst/go/                 generic concrete tree and stored AST view
  representation/go/      checked AST/CST snapshots and language/profile descriptors
contracts/
  language/go/            LanguageProvider identity
  pipeline/go/            typed Stage contract
  text/go/                input, diagnostics and specialized provider contracts
  transform/go/           ASTTransformer contract and typed registration bridge
  semantic/go/            common snapshots, requests, reports and publication
  resolution/go/          resolver capability and dependency contract
  analysis/go/            analysis capability and explicit resolution input
semantics/go/             immutable provider-owned semantic facts and graph checks
engine/go/                explicit registration, planning and execution
syntax/<language>/go/     native abstract/concrete vocabulary and identity
resolve/horos/go/         Horos resolution, fact export and canonical-value helpers
parse/<language>/go/      source text to CST
check/<language>/go/      language syntax checking
analyze/horos/go/         consume resolution evidence and reuse its model
format/<language>/go/     AST view to CST
provider/<language>/go/   explicit composition of the language's stages
check/ast/go/             generic abstract carrier checks
check/cst/go/             generic concrete carrier checks
print/go/                 generic CST to source text
internal/<language>/     private language-specific lexical helpers
internal/                shared private checking/parsing/tree mechanics
integration/go/          independent consumers and cross-package qualification
testdata/text/<language>/valid/
docs/
```

The native text language directories are `typescript`, `golang`, `rust` and
`cpp`; `horos` implements the AST-only syntax/resolve/check/analyze/provider profile.
The final `go` directory denotes implementation language, not target:
`syntax/rust/go` is Go code describing Rust syntax. `ast` and `cst` in the check
stage are reserved generic-carrier components, not native language names.

The source and its package-level tests belong with each implementation.
Independent consumer fixtures, stage composition, language conformance and native
compiler qualification remain under integration. Private helpers are not an
alternative public API. Do not create a language-first packages tree, additional
Go modules, compatibility aliases or placeholder source directories.

## Dependency direction

Directory groups are not Go packages or an assertion that all members share the
same dependency level. The package boundaries are:

| Component | Permitted responsibilities and dependencies |
| --- | --- |
| AST, Stage, language identity | Independent foundations; no native syntax or other foundation dependency |
| CST | Generic AST carrier |
| Text contracts | AST/CST carriers and language identity |
| Carrier checks and printer | Generic tree structure and shared diagnostic contracts; no native grammar or formatting choices |
| Checked representation | Carriers, text checker/identity contracts and generic carrier checks; no engine or transformer dependency |
| Transformer contract/bridge | Typed stages, checked representations and endpoint identity; no concrete mapping |
| Semantic facts | Generic AST addresses; graph integrity and owned storage, no language rules |
| Shared semantic contract/bridge | Semantic facts, checked representations and identity; capture/publication, no concrete provider |
| Resolution contract/bridge | Shared semantic contract and identity; binding coverage and typed invocation |
| Analysis contract/bridge | Shared semantic and resolution contracts; explicit dependency composition |
| Engine | Checked representations and shared contracts; no concrete language, provider or mapping |
| Native syntax | Carriers and identity; no implementation stage or other language |
| Native checker | Own syntax, shared checking mechanics and own lexical helpers |
| Native parser/formatter | Own syntax/checker and appropriate shared mechanics; neither imports the other |
| Native provider | Own syntax and stages plus the shared printer; composition only |
| Horos syntax | Generic AST and identity; no semantic-model imports |
| Horos checker | Own syntax and generic structural/grammar checks; no semantic resolution |
| Horos resolver | Own syntax, private projection/correspondence helper, Horos Go model and Ontos values/codec |
| Horos analyzer | Own syntax, Horos resolver contract and shared semantic evidence; no direct external model imports |
| Horos provider | Own syntax/checker/resolver/analyzer and shared contracts; composition, no text stages |

The checked AST/CST and language descriptors share one home in
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
union of language directories across syntax/parse/check/format/provider/resolve/analyze
and require each language to have every stage in its tested capability profile.
The four native languages retain all five native text stages; Horos requires
syntax/resolve/check/analyze/provider and must not acquire unqualified text stages.
Generic carrier check directories
are excluded explicitly. An orphan parser and an incomplete syntax package both
fail. The checks also classify every production package and enforce its directed
dependencies, including private helpers and narrow external-model allowances
only for the Horos frontend. Generic carriers,
contracts, engine and printer must remain independent of Horos/ontos, and native
syntax/text stages must remain independent of the canonical frontend.

The [native conformance suite](integration/go/conformance_test.go) covers
every native language. It checks stable identity, matching AST/CST checker
views, formatter snapshot ownership and metadata, whole-input rejection, located
parse errors and empty output on failure. The same test-only fixture table feeds
source round trips and fuzzing. Tests never register a production default.
Independent compiler fixtures remain a separate required qualification.
The [AST-only suite](integration/go/horos_test.go) separately checks Horos identity,
registration, checking, snapshot ownership, specialization and absent text capabilities.

Existing public consumer/negative type-check cases, unfamiliar-language tests,
route oracle, concurrency checks and compiler gates remain required. Directory
consistency does not prove grammar or mapping correctness.

## Future components

The [agnostic semantic model](docs/semantics.md), [analysis contract](docs/analysis.md)
and first Horos analyzer are implemented. Additional analyzers use the same
optional analyze stage with an explicitly tested capability profile. Native
semantic analysis is not implied by native text tooling.

The [Horos declaration integration](docs/horos.md) implements
`syntax/horos/go`, `check/horos/go`, `analyze/horos/go`, `provider/horos/go` and
`resolve/horos/go`, with private projection/correspondence in `internal/horos`.
Readable syntax, when specified, adds `parse/horos/go` and
`format/horos/go` under the same stage matrix. These replace the former proposed
`canonical` paths; there are no compatibility aliases.

Directed mappings will live under `transforms/horos-to-<native>/go`; later partial
reverse mappings under `transforms/<native>-to-horos/go`. Application composition
and IO remain under `cmd/bittype/go`. Horos owns the declaration semantics and
encoding; these packages supply tooling projections and native representations.
Dependency rules and qualification will accompany implementation, with no empty
packages created now. The complete language-stage gate will accompany a complete
readable frontend; an initial AST-only provider must have an explicit separately
tested capability profile rather than pretending missing text stages exist.

Adding a native language supplies the five stages, any private lexical helpers,
conformance fixtures and compiler qualification. It requires no change to the
generic carriers, engine or production registration catalogue: there is no such
catalogue. Applications choose their registrations. Adding Function, Stream,
Cell or another library declaration does not add a stage, language or core kind.

See the [charter](CHARTER.md), [engine API](docs/engine.md),
[native text design](docs/parsing-formatting.md) and
[future transformation plan](docs/generation.md).

The [provider matrix](docs/providers.md) records implemented capabilities.
[Resolution composition](docs/resolution.md) defines snapshot identity and explicit model reuse.
