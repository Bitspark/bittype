# Repository layout

Use **component first, language second**, following the family
[layout policy](https://github.com/Bitspark/bitwire/blob/main/LAYOUT.md).
This rule applies to new source and to deliberate source-layout migrations.

```text
<repo>/<component>/<lang>/...
<repo>/cmd/<command>/<lang>/...
```

A component can have a semantic hierarchy, such as `core/node/go/` or
`providers/tree/memory/go/`. Choose that hierarchy before the language
directory. Do not organize implementations as repository-root `go/` or `ts/`,
`go/<component>/`, or a language-first `packages/` tree.

Use exactly two lowercase letters: `go`, `ts`, `py`, `rs`, `hs`, `cc`,
`jv`, `sw`; other assigned codes include `rb`, `kt`, `cs` and `sh`.
Bitwire already uses `hs` for Haskell. Service SDKs and their source manifests
also follow their own language registry; register a code there before using it.

Source, native tests and language-specific package metadata belong with the
language implementation. Shared specifications, schemas, vectors, documentation,
assets and deployment configuration stay language-neutral. Repository-root
workspace/build manifests and maintenance scripts may remain at their normal
tooling locations. A source directory does not automatically require its own
module, package publication or repository.

Create a language directory only when it contains a real implementation or
contract presentation. Keep module boundaries, dependencies and release rules
consistent with the repository's charter. When moving existing source, update
imports, manifests, generators, tests, CI and documentation together; never
rewrite an immutable published release or bypass a frozen-foundation policy.

## Adoption in this repository

Use component names for the distinct language, identity, presence and generation
kernel responsibilities in the [charter](CHARTER.md), with language directories
beneath each. Names and module boundaries are implementation decisions; the
layout does not merge the generation kernel into the declaration language.
No implementation is implied by these conventions.
