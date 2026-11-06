# Optional static verification integration

This Go module composes actual Bittype captures, authoritative Horos queries,
Go type checking and the landed Logos typed checker. It is isolated because
Logos is a private dependency. The root Bittype module stays public and has no
dependency on this module.

- `files/source/horos/go`: admit the captured resolution receipt, reconstruct
  its owner-local Program and run explicit `Satisfies` queries.
- `files/target/golang/go`: translate captured Go ASTs directly to `go/ast`,
  check the supported type-only package with `go/types`, and query instantiation.
- `files/horos-to-go/go`: select actual declarations/fields/helper, construct
  the fixed claim, prepare immutable requests and directly call Logos.
- `integration/go`: actual TextBox/Page inputs, nested reuse and architecture.

Read the [contract and API](../docs/verification.md#implemented-go-surface) and
[module decision](../docs/type-system-architecture.md#proposed-package-roles).
This does not register a general Go semantic provider, implement native admission
enforcement, generate methods or supply runtime adapters.

All dependency preparation, formatting, compilation, tests and vet involving
this module run on the private Logos GitHub runner. Do not run them on the owner
workstation or expose private dependency material through public caches/artifacts.
