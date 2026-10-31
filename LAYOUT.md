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

The implemented source packages are `ast/go` and `pipeline/go`, with public
contract tests in `integration/go`. The [generation design](docs/generation.md)
records the full target layout for syntax, frontend, generation, formatting
and application composition. Those later packages do not yet exist.
Do not create empty source directories or present planned packages as delivered.

The generic AST and pipeline remain distinct from language grammars.
The [charter](CHARTER.md) sets ownership; this layout does not merge those
responsibilities or require a package or repository per directory.
