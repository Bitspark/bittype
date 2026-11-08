# bittype/ACX-007: Bittype matrix dependencies keyed by display names

**Status:** Corrected within the stated scope, including unnamed arguments.
**Scope:** Source-derived provider availability in Bittype's documentation tool.
The first correction landed in [bittype PR126](https://github.com/Bitspark/bittype/pull/126)
at `a582c5bc9327b38b75b16dba59fc558067ed1171`.
The completed identity correction landed in [bittype PR128](https://github.com/Bitspark/bittype/pull/128)
at `580aca20548553e083d5fd6a7c65ec55a72528fb`.

## Rules involved

[AC-05](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-05) requires composition to preserve
contract identity and dependency guarantees. The discovery tool resolved real Go
types, then discarded their package ownership when deciding whether a constructor's
required provider existed. An implementation of an unrelated contract could satisfy
that dependency in the report.

This is an identity failure in shared orchestration, not a missing language case.
The requirement applies even when all existing package names happen to be distinct.
[AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07) qualifies the repair's evidence:
correct availability does not establish semantic completeness or execution coverage.

## Intended boundary

Each package owns its Go declarations. Discovery may format a short name for a
table, but constructor requirements must retain the original contract and its
concrete type arguments. Availability is the least fixed point of those exact
requirements; an unrelated implementation must not create an available cell.

## Actual implementation and consequence

At [Bittype 7acd3ba](https://github.com/Bitspark/bittype/blob/7acd3ba12cf260e31d581eb1ec317128e5d0ede9/tools/provider-matrix/go/discover.go),
`contractName` used `types.TypeString` with `Package.Name()` and replaced the
language marker's spelling with `L`. The `available` map then used that display
string for both discovered providers and constructor dependencies.

Two packages may both be named `print` and export `Provider[I, O]` while declaring
different methods. The concrete counterexample has:

- `contracts/print/go.Provider` requiring `Print`, with an existing Syntax printer.
- `contracts/elsewhere/go.Provider` requiring `Render`, with no implementation.
- A generic Files printer whose constructor requires the second provider.

Both requirements displayed as `print.Provider[cst.CST[L], string]`. The existing
`Print` implementation therefore discharged the unavailable `Render` dependency,
and the matrix advertised Files printing. Go's nominal ownership information was
available; presentation had been allowed to substitute for it.

## What worked and what it did not establish

Ordinary discovery, generic construction and direct transformer-direction examples
worked with the repository's existing names. Those examples did not exercise two
independent declaration owners with the same display spelling. Green examples
therefore did not establish preservation of contract identity through availability
composition.

## Required correction

Use actual Go type identity, including concrete arguments, for dependency
resolution. Keep labels exclusively for presentation. If different contracts would
occupy an indistinguishable display row, diagnose that ambiguity rather than silently
merging their implementations. Preserve the existing least-fixed-point construction
rules and direct transformer directions.

## Correction and qualification status

PR126 replaced short display keys with package-qualified type strings in the importable
[`tools/provider-matrix/report/go`](https://github.com/Bitspark/bittype/tree/a582c5bc9327b38b75b16dba59fc558067ed1171/tools/provider-matrix/report/go)
host. `TestDependenciesDoNotUseDisplayContractNames` requires the Files printer to
remain unavailable in the counterexample. A second regression,
`TestDiscoveryRejectsSameDisplayContractFromDifferentOwners`, requires an explicit
diagnostic for colliding display rows.

Both regressions passed in the full Go test suite at `f209457`, the head merged
by PR126, in [run 42254515087](https://github.com/Bitspark/bittype/actions/runs/42254515087).
The tooling tests also passed locally on Windows. The independent declaration
owners in that counterexample establish separation of those named declarations.
They did not establish preservation of every concrete argument's Go identity.

Independent review found that `types.TypeString` still omits package ownership of
private fields and methods inside unnamed types. For example,
`Provider[struct{value int}]` has the same string when `value` belongs to package
`a` or package `b`, although `types.Identical` distinguishes the instantiations.
The same loss affects row-collision checks. A Syntax printer returning the first
struct must not satisfy a Files printer's dependency on the second.

PR128 uses `go/types.Identical` for availability, operation deduplication, call
relationships and case-context association. Matrix row grouping substitutes only
the exact selected language marker with a shared Go type parameter before that
comparison. It retains other nominal origins, private field/method packages,
field tags, channel direction and function/interface structure. Named underlying
definitions are never rewritten. This abstraction is for grouping rows, not an
admitted program instantiation or a replacement Go identity implementation.

Display strings no longer decide semantic equivalence. Unequal operations survive
collection; a collision in a rendered row label or anchor produces an explicit
error instead of merging them. The copied independent-review discovery regression
failed on the previous implementation and passes with this repair. Additional
tests cover private/exported structures and interfaces, nested arguments, shared
private aliases, row abstraction, non-marker names and field tags, ambiguous case
contexts, call arguments and anchor collisions.

The tooling suite and vet passed locally. The full
[Go contracts/native workflow](https://github.com/Bitspark/bittype/actions/runs/42258794400)
and [documentation checks](https://github.com/Bitspark/bittype/actions/runs/42258794478)
passed at `b773f1a`, the exact tree merged by PR128. Independent review inspected the
substitution and association code and reran seven targeted regressions with no
remaining blocker in this bounded repair. This evidence does not prove every
possible Go type or establish whole-language capability.
This correction does not
mark the separate [ACX-001 verification-composition repair](0001-verification-composition.md)
complete. The owning [capability-report design](https://github.com/Bitspark/bittype/blob/main/docs/capability-reports.md)
keeps discovered availability, definition inventories and test observations distinct.

## Assessment lesson

Inspect the identity used at the decision point, not merely whether the earlier
parser or type checker resolved declarations correctly. A display string can be
adequate for a heading and invalid as a dependency key. An independently owned,
same-spelled contract exposes that distinction without introducing a language list
or a new reporting policy.
Package qualification alone is also insufficient: formatting a type can omit
identity-bearing details inside its arguments. Use the owning compiler's identity
operation for semantic decisions and treat display collision handling separately.

This is the maintained Bittype-owned record. The
[original platform record](https://github.com/Bitspark/platform/blob/main/docs/concepts/architecture-cases/0007-bittype-matrix-contract-identity.md)
serves as a cross-project reference. [Capability report contract](../capability-reports.md) ·
[Owning case index](README.md).
