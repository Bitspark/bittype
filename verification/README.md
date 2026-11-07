# Static verification integration

These packages compose actual Bittype captures, authoritative Horos queries,
Go type checking and the public Logos typed checker. They belong to Bittype's
single root Go module.

- `files/source/horos/go` admits the captured Horos receipt and runs authoritative
  queries. Native facts and admission come from the [Go semantic providers](../docs/native-semantics.md#first-go-profile).
- `files/horos-to-go/go` owns the scoped subject and law protocol, the conditional
  family interpretation and the admission counterexample.
- `files/horos-to-go/mappings/go` supplies independently selectable text, bool,
  list, optional and closed-record rules.
- `logos/go` freezes typed questions, admits captures and origins, applies receiver
  limits and calls the public Logos checker.
- `integration/go` checks TextBox and Page, nested theorem reuse, and a separately
  owned tuple mapping composed with the existing list and record rules.

[Verification API](../docs/verification.md#implemented-go-surface) ·
[Architecture case and qualification](../docs/architecture-cases/0001-verification-composition.md#correction-and-qualification).

Dependency preparation, formatting, compilation, tests, vet and race checks run
through Bittype's existing GitHub workflow, including its anonymous installed
consumer. Do not run these workloads on the owner workstation.
