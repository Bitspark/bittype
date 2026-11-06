# Horos native type mappings

These packages implement the [public generation SPI](../../../../docs/generation-spi.md)
inside Bittype's root Go module. Horos owns canonical type meaning; these packages
own the selected native representations. The shared composer, source reader,
native backends and thin transformer factories do not import them.

```text
generation/syntax/mapping/horos/
  go/                    structural declaration policy and associations
  typescript/go/         TypeScript mappings and exact bindings
  golang/go/             Go mappings and native optional helper
  rust/go/               Rust mappings and exact bindings
  cpp/go/                C++ mappings and exact bindings

generation/syntax/composition/
  horos-to-typescript/go/ optional explicit default assembly
  horos-to-go/go/
  horos-to-rust/go/
  horos-to-cpp/go/
```

The default compositions require explicit TreeID and Profile, plus Go Package
or C++ Namespace where applicable. They consume the supplied Horos semantic
capture and analysis evidence. They do not resolve dependencies, run hidden
analysis, produce proofs or implement runtime adapters.

The `pagination-representation/1` profile covers text, bool, lists, named closed
records, optional fields, proper generic parameters and named applications.
It retains an open Page<T> and reuses it for Page<Person> and nested Page uses.
Source kind constraints are checked but may be erased in native parameters.
`pagination-enforcement/1` remains explicitly unsupported. Neither profile
promises arbitrary native instantiations satisfy Horos contracts.

Use Config.Configure to add exact library bindings or guarded overrides to a
fresh Setup. An override must identify the existing mapping. Duplicate, missing
or stale overrides fail. All child mappings use the same effective configuration;
a selected mapping's failure is terminal, with no fallback. Native publication
policy and body representation must agree.

Independent libraries use the same SPI; these mapping helpers are optional.
The [composition fixtures](../../../../integration/go/horos_generation_composition_test.go)
replace the list mapping and add an ordinary Cell representation without changes
to the shared composer, runtime or Wire. The
[pagination fixtures](../../../../integration/go/horos_generation_pagination_test.go)
exercise all four targets, source admission, imported identities and native compilers.
The [installed consumer](../../../../integration/go/testdata/installed/horos_generation_test.go)
uses public Bittype packages without local replacements. Existing generated-Go
verification remains a separate [test](../../../../integration/go/horos_generation_verification_test.go),
subject to the [documented verification limitation](../../../../docs/verification.md#known-architectural-violation).

Owner placement correction, 6 November 2026: these implementations moved from
Horos into Bittype. The old Horos module and its module-only CI job are removed,
without compatibility packages. Semantic ownership and public extension boundaries
are distinct from repository placement; Horos has no generator dependency on Bittype.
