# bittype/ACX-009: Package boundaries do not follow contributions

**Status:** Identified; correction pending.
**Scope:** The native and common horos mapping packages, the horos-to-Go
verification laws, the Go and TypeScript analysis provider packages, the
publication of contract identities and horos semantic names, and the import
check in
[`integration/go/architecture_test.go`](https://github.com/Bitspark/bittype/blob/3bf45c3/integration/go/architecture_test.go),
at `3bf45c3`. Identified on 11 November 2026 after the owner asked for this
check. Plan in [#218](https://github.com/Bitspark/bittype/issues/218); design in
[package boundaries](../package-boundaries.md).

## Rules involved

[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
states that a misplaced responsibility is incorrect and moves first. The
instances below are misplacements at the level of the package: a part sits in a
package whose other parts it does not need, or a name is published by an
implementation.

[AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07)
limits what evidence supports. `LAYOUT.md` says the architecture checks
"classify every production package and enforce its directed dependencies". That
is true per package. It does not establish that two contributions inside one
package are independent, that an allowed import uses only what it should, or
that an allowance is needed.

No instance is recorded here as a violation of
[AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)
or [AC-03](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-03).
An independent mapping library can still enter the kernel without editing any
of these packages. Two neighbouring faults are recorded elsewhere: the closed
profile gate is [bittype/ACX-008](0008-structural-profile-gate.md), and the
optional-field decision inside Record is the contract change of
[#196](https://github.com/Bitspark/bittype/issues/196), item 6.

## Intended boundary

One contribution, meaning one thing a consumer registers or selects on its own,
is one package. A name that two parties must agree on lives in a vocabulary
package that imports no implementation. Only a composition imports a
contribution. A contribution that needs another's result states a demand, and
the kernel resolves it with whatever the consumer registered.

## Actual implementation and consequence

**Contributions grouped in one package.**

- `transform/syntax/horos-to-{typescript,golang,rust,cpp}/mapping/go` each hold
  `Text`, `Bool`, `List`, `Record` and `Declarations`. The Go package also
  holds `Optional`.
- Each of the four exports `Libraries(profile)`, a list of its own mappings.
  [`docs/generation.md`](../generation.md) states "There are no per-pair
  assembly factories or presets."
- `verification/files/horos-to-go/mappings/go` holds the `Text`, `Boolean`,
  `List` and `Record` laws.
- `transform/syntax/horos/mapping/go` holds the profile gate with its `Library`
  adapter, the pattern constructors, the declaration recipe and the naming
  function.
- `analyze/files/golang/provider/go` and `analyze/files/typescript/provider/go`
  each hold the analysis provider and a native query service. The two share no
  declaration.

**Identities published by an implementation.** 22 imports go from a production
package into a stage provider package. Eight use only an identity:

| Importer | Provider package | Symbols |
| --- | --- | --- |
| `analyze/files/golang/provider/go` | `resolve/files/golang/provider/go` | `Contract` |
| `analyze/files/rust/provider/go` | `resolve/files/rust/provider/go` | `Contract` |
| `analyze/files/typescript/provider/go` | `resolve/files/typescript/provider/go` | `Contract`, `Compiler` |
| `analyze/files/horos/provider/go` | `resolve/files/horos/provider/go` | `Contract` |
| `internal/horosfiles` | `resolve/files/horos/provider/go` | `PreludePin` |
| `transform/syntax/horos/source/go` | `analyze/files/horos/provider/go` | `Contract`, `Profile` |
| `transform/syntax/horos/source/go` | `resolve/files/horos/provider/go` | `ExpressionRole`, `IdentityRelation`, `ResultSortRelation`, `SortRelation` |
| `verification/files/source/horos/go` | `resolve/files/horos/provider/go` | `Contract`, `Profile` |

A ninth, `verification/files/horos-to-go/go`, imports the Go analysis provider
only for the query service's `Admission`.

**The check.** The four native targets may import their checker, and none does.
The test passes. An allowed import may use any exported symbol.

Consequences:

- Coupling between two mappings of one pair is invisible to the import check.
  One exists: Go's Record reaches the Optional helper through the package
  variable `optionalKey`
  ([source](https://github.com/Bitspark/bittype/blob/3bf45c3/transform/syntax/horos-to-golang/mapping/go/mapping.go#L72-L88)).
- A representation family written into one of these packages has to move when
  the package is divided. [#207](https://github.com/Bitspark/bittype/issues/207)
  writes the first ones.
- A reader of horos semantic facts imports the resolver implementation to name
  a relation. A second resolver would have to be imported next to the first, or
  repeat the names.
- Every new mapping edits `Libraries`, and every new package edits the test's
  list.

## What worked and what it did not establish

The architecture test passes and its rules hold: the kernel, the source and the
targets import no concrete mapping package. Independent contributions enter the
kernel through the public contract, as the composition tests show. Those
results establish that the kernel is open and that the listed edges are absent.
They do not establish that packages follow contributions, because the test
reads imports per package and accepts any use of an allowed import.

## Required correction

1. Each builtin mapping and each verification law is its own package, and no
   package lists them ([#222](https://github.com/Bitspark/bittype/issues/222)).
2. The query services leave the analysis provider packages
   ([#221](https://github.com/Bitspark/bittype/issues/221)).
3. Contract identities and horos semantic names are published by vocabulary
   packages, and no package imports a provider for a name
   ([#220](https://github.com/Bitspark/bittype/issues/220)).
4. Every remaining import of a provider names its symbols and its reason, and
   an unused allowance fails
   ([#219](https://github.com/Bitspark/bittype/issues/219)).
5. The import rules are stated per role, so a new package of an existing role
   needs no edit to the test
   ([#218](https://github.com/Bitspark/bittype/issues/218)).
6. Every identity value and every generated declaration stays byte-identical.
   These are moves.

## Correction and qualification status

Nothing is implemented. All instances are unchanged at `3bf45c3`. The order is
the one in [package boundaries](../package-boundaries.md#order-of-work). The
mapping split is ordered before #207 under
[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10).

Three candidates are not confirmed and are not part of this case until they
are: the thirteen direct calls into a same-language checker or parser, the
`Self` rule in the shared naming function, and `CheckSingleScope` in the shared
target package. Their issues decide them.

Qualification needs, for each rule of the check, a test that feeds it a package
set breaking exactly that rule.

## Assessment lesson

The import check was read as proof that contributions are independent. It
proves that listed edges between packages are absent. Two questions were never
asked of it: what can reach what inside one package, and how much of an allowed
import is used. When a dependency check is green, ask what its unit is, and
whether the thing that must be independent is smaller than that unit.
