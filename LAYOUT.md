# Repository layout

Use component-first, language-second source paths:

```text
<component>/<implementation-language>/...
cmd/<command>/<implementation-language>/...
```

Language-neutral specifications and examples belong under `docs/` and, once
implemented, a shared fixture directory. The current model proposal is
[docs/model.md](docs/model.md).

The toolchain is implemented in Go. The trailing `go` directory identifies
that implementation language. Target names belong to components: for example,
`syntax/rust/go` contains Go code for Rust syntax, not Rust source.
There is one root Go module, with independently owned component packages.

Source, native tests and package metadata belong with the owning implementation.
Root workspace manifests and CI configuration may stay at their normal locations.
Do not introduce repository-root language directories or a language-first
`packages/` tree.

Implemented packages include AST/CST carriers, typed pipeline and provider
contracts, four native syntax packages and their parse/check/format/provider
packages, plus the shared printer. Private helpers live under internal/.
The [native text design](docs/parsing-formatting.md) records their layout and
dependencies; the [generation design](docs/generation.md) separates them from
future canonical generation. The implemented frame uses transform/go for the
ASTTransformer contract and typed adapters, engine/go for the agnostic registry
and route execution. Future transform/<from>-to-<to>/go packages own each directed
mapping. Canonical AST, actual transformers, public builders and application CLI
remain deferred. The [engine reference](docs/engine.md) records the public API;
the generation document marks all planned paths separately.
Do not create empty source directories or present planned packages as delivered.

The generic AST, pipeline and language identity contract remain independent
of concrete syntax packages. Each target package owns its vocabulary and
implements language identity methods without registration side effects.
The [charter](CHARTER.md) sets ownership; this layout does not merge those
responsibilities or require a package or repository per directory.
