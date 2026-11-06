# Static verification integration

These packages compose actual Bittype captures, authoritative Horos queries,
Go type checking and the public Logos typed checker. They belong to Bittype's
single root Go module; only the owning Horos/Go composition imports the checker.

**Architectural correction pending:** [bittype/ACX-001](../docs/architecture-cases/0001-verification-composition.md)
records the closed mapping/claim contribution path and its required replacement.
The scoped proof results do not establish an independently extensible verifier.

- `files/source/horos/go`: admit the captured resolution receipt, reconstruct
  its owner-local Program and run explicit `Satisfies` queries.
- Native facts and admission come from the independently registered
  [Go resolution/analysis providers](../docs/native-semantics.md#first-go-profile).
  The former verification-only Go projection is removed.
- `files/horos-to-go/go`: select actual declarations/fields/helper, construct
  the fixed claim, prepare immutable requests and directly call Logos.
- `integration/go`: actual TextBox/Page inputs and nested reuse. The root
  architecture test enforces every production package's import boundary.

Read the [contract and API](../docs/verification.md#implemented-go-surface) and
[package ownership](../docs/type-system-architecture.md#proposed-package-roles).
This verification composition does not implement native admission enforcement,
generate methods or supply runtime adapters.

Dependency preparation, formatting, compilation, tests, vet and race checks run
through Bittype's existing GitHub workflow, including its anonymous installed
consumer. Do not run these workloads on the owner workstation.
