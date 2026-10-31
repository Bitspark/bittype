# Repository layout

Use component-first, language-second source paths:

```text
<component>/<lang>/...
cmd/<command>/<lang>/...
```

Language-neutral specifications and examples belong under `docs/` and, once
implemented, a shared fixture directory. The current model proposal is
[docs/model.md](docs/model.md).

Use `ts`, `go` and `rs` for this milestone's three target languages. The
language a generator is implemented in is independent of the languages it emits.

Source, native tests and package metadata belong with the owning implementation.
Root workspace manifests and CI configuration may stay at their normal locations.
Do not introduce repository-root language directories or a language-first
`packages/` tree.

Choose components from their responsibilities, such as the model frontend,
native type rendering or presence data types. Do not create empty source
directories or imply that any implementation already exists.

A generation kernel, if introduced, remains distinct from the declaration model.
The [charter](CHARTER.md) sets ownership; this layout does not merge those
responsibilities or require a package or repository per directory.
