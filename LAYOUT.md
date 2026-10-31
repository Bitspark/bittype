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
engine/go/                explicit registration, planning and execution
syntax/<language>/go/     native abstract/concrete vocabulary and identity
parse/<language>/go/      source text to CST
check/<language>/go/      native syntax checking
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

The implemented language directories are `typescript`, `golang`, `rust` and
`cpp`. The final `go` directory denotes implementation language, not target:
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
| Engine | Checked representations and shared contracts; no concrete language, provider or mapping |
| Native syntax | Carriers and identity; no implementation stage or other language |
| Native checker | Own syntax, shared checking mechanics and own lexical helpers |
| Native parser/formatter | Own syntax/checker and appropriate shared mechanics; neither imports the other |
| Native provider | Own syntax and stages plus the shared printer; composition only |

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
union of native directories across all five stages and require every language
to have an implemented package in each stage. Generic carrier check directories
are excluded explicitly. An orphan parser and an incomplete syntax package both
fail. The checks also classify every production package and enforce its directed
dependencies, including private helpers and the absence of external dependencies.

The [shared conformance suite](integration/go/conformance_test.go) must cover
every discovered language. It checks stable identity, matching AST/CST checker
views, formatter snapshot ownership and metadata, whole-input rejection, located
parse errors and empty output on failure. The same test-only fixture table feeds
source round trips and fuzzing. Tests never register a production default.
Independent compiler fixtures remain a separate required qualification.

Existing public consumer/negative type-check cases, unfamiliar-language tests,
route oracle, concurrency checks and compiler gates remain required. Directory
consistency does not prove grammar or mapping correctness.

## Future components

The [canonical model](docs/model.md) remains a proposal. Its implementation will
use the same `syntax`, `parse`, `check`, `format` and `provider` paths, with explicit
resolution under `resolve/canonical/go` when implemented. Directed mappings will
live under `transforms/<source>-to-<target>/go`; application composition and IO
will live under `cmd/bittype/go`. Their dependency rules and qualification will
be added with those implementations, not through empty packages now.

Adding a native language supplies the five stages, any private lexical helpers,
conformance fixtures and compiler qualification. It requires no change to the
generic carriers, engine or production registration catalogue: there is no such
catalogue. Applications choose their registrations. Adding Function, Stream,
Cell or another library declaration does not add a stage, language or core kind.

See the [charter](CHARTER.md), [engine API](docs/engine.md),
[native text design](docs/parsing-formatting.md) and
[future transformation plan](docs/generation.md).
