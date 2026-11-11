# Package boundaries and forbidden dependencies

**Proposed, 11 November 2026. Nothing in this document is implemented.** The
owner asked on 11 November 2026 to check for packages "that contain multiple
implementations / functionalities that either already are separate and neither
conceptually nor practically depend on each other, or which still depend on
each other but shouldn't, because it violates our kernel principles (e.g., a
function calling another function directly instead of invoking the kernel for
invoking the function's semantics, the kernel then deciding which registered
implementation to call)", and for a way to enforce the result "through
forbidden dependencies that we could check for and enforce".

This document records the difference between the two faults, what to look for,
the findings of a scan of `main` at `3bf45c3`, the rules to enforce and the
order of work.

[Plan, #218](https://github.com/Bitspark/bittype/issues/218) ·
[Open findings, bittype/ACX-009](architecture-cases/0009-package-boundaries.md) ·
[Layout](../LAYOUT.md) ·
[Existing architecture checks](../integration/go/architecture_test.go) ·
[Kernel-owned composition](kernel-composition.md) ·
[Shared guidance](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#package-boundaries-and-forbidden-dependencies)

## Terms

- **Kernel.** The shared mechanisms that select and compose: the engine, the
  generation registry and planner, and the contracts they are written against.
- **Contribution.** One thing a consumer registers or selects on its own: a
  stage provider, a source, a target, one mapping library, one verification
  law.
- **Vocabulary.** The names two parties must agree on and nothing else: syntax
  kinds, semantic keys, profile and contract identities, intended native forms.
  A vocabulary package imports no contribution.
- **Composition.** Code that chooses contributions and registers them: the
  integration tests, the tools and a future command.
- **Private part.** Code that serves one owner and is never selected by a
  consumer, such as a lexical helper of one language.

## Two different faults

| | Wrong grouping | Wrong coupling |
| --- | --- | --- |
| What it is | Parts that do not need each other share one package | A contribution names another contribution that the consumer is meant to choose |
| What is wrong | The package boundary. No dependency between the parts is wrong | The dependency. The caller makes the kernel's choice |
| Consequence | An importer of one part receives all of them, and the import check cannot tell the parts apart | The kernel's selection, scope and failure guarantees do not apply to that call. Replacing the callee means editing the caller |
| Correction | A move. Nothing is rewired | A demand that the kernel resolves. This can need a kernel capability |

The first fault hides the second. Inside one package Go lets every declaration
reach every other, and the import check works per package. Coupling between
two contributions in one package is therefore invisible to it. Splitting comes
first, which is the order
[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
requires: move what is misplaced, then add what is missing.

## The question that separates them

Not every direct call is wrong. Ask of each call from one part to another:

> Whose answer does the caller need: its own, or the composition's?

- **Its own.** The caller relies on a private invariant, a shared mechanism or
  a vocabulary. A direct call is right. If the callee sits in a contribution
  package, the question is only whether it is in the right package.
- **The composition's.** The caller needs whatever the consumer registered for
  that meaning. It must state a demand and let the kernel resolve it. A direct
  call fixes the consumer's choice inside the caller.

Two examples from `main` at `3bf45c3`:

- In a Record mapping, `r.Mapped(f.Use)` asks the kernel for the native type
  of a field. Whichever mapping the consumer selected answers. This is the
  correct form.
- Two lines below, the same Record mapping decides how an optional field is
  represented. Rust builds `std::option::Option`, C++ builds `std::optional`,
  and Go wraps the type in a helper through a package variable. A consumer who
  wants another representation of optional fields must replace Record.

Parts that share a private invariant belong together. The planner and the plan
types in `contracts/generation/go` are one unit, because only the planner may
issue a plan.

## What to look for

Signs of wrong grouping:

- **Several contributions in one package.** A consumer can register each one
  separately, and they share a package only because they were written together.
- **A function that lists the package's contributions.** Adding a contribution
  edits the list.
- **An importer that uses a slice of the package.** It imports a provider
  package and uses one identity, or one service that shares no declaration
  with the provider.
- **A name kept with an implementation.** A profile, a contract, a relation
  key or a pin is published by a provider package and not by a vocabulary.
- **A policy next to a shared mechanism.** A shared package also holds one
  case's rule.

Signs of wrong coupling:

- **A contribution imports a contribution.** Each such import needs a reason
  that survives the question above.
- **A callback that builds where it could ask.** It constructs the native form
  of something another contribution owns.
- **A package variable shared by two contributions.** One reaches the other
  without declaring it.
- **A shared function that carries one case's meaning.** Every caller inherits
  a rule that belongs to one of them.

## Findings on `main` at `3bf45c3`

The scan covered bittype only. Each finding was read in source after the scan
reported it.

| # | Where | Kind | Status | Tracked |
| --- | --- | --- | --- | --- |
| 1 | `transform/syntax/horos-to-{typescript,golang,rust,cpp}/mapping/go` each hold Text, Bool, List, Record and Declarations. Go also holds Optional | Grouping | Confirmed | [#222](https://github.com/Bitspark/bittype/issues/222) |
| 2 | `Libraries(profile)` in each of those packages lists the package's own mappings | Grouping | Confirmed | [#222](https://github.com/Bitspark/bittype/issues/222) |
| 3 | `verification/files/horos-to-go/mappings/go` holds the Text, Boolean, List and Record laws | Grouping | Confirmed | [#222](https://github.com/Bitspark/bittype/issues/222) |
| 4 | `transform/syntax/horos/mapping/go` holds the profile gate, the pattern constructors, the declaration recipe and the naming function | Grouping | Confirmed | [#222](https://github.com/Bitspark/bittype/issues/222); the gate in [#216](https://github.com/Bitspark/bittype/issues/216) |
| 5 | `analyze/files/{golang,typescript}/provider/go` each hold the analysis provider and a native query service that share no declaration | Grouping | Confirmed | [#221](https://github.com/Bitspark/bittype/issues/221) |
| 6 | Eight imports of a provider package use nothing but an identity. horos has no vocabulary package for its semantic names | Identity kept with an implementation | Confirmed | [#220](https://github.com/Bitspark/bittype/issues/220) |
| 7 | Record decides optional fields itself in Go, Rust and C++ | Coupling | Confirmed; the contract change is designed | [#196](https://github.com/Bitspark/bittype/issues/196), item 6 |
| 8 | The existing check passes unused allowances and cannot see which symbols an import uses | Enforcement gap | Confirmed | [#219](https://github.com/Bitspark/bittype/issues/219) |
| 9 | Thirteen direct calls into a same-language provider: `Check` from parse, format and the horos resolver; `Parse` from the Rust and TypeScript semantic helpers | Coupling candidate | Read as legitimate: the caller needs its own grammar's answer. `LAYOUT.md` permits them | Recorded as named exceptions in [#219](https://github.com/Bitspark/bittype/issues/219) |
| 10 | `Exported` appends `_` to `Self` for all four targets. `Self` is a keyword only in Rust | Coupling candidate | Not confirmed | Decided in [#222](https://github.com/Bitspark/bittype/issues/222) |
| 11 | `CheckSingleScope` states one placement's rule in the shared target construction package | Grouping candidate | Not confirmed | With the design point on [#161](https://github.com/Bitspark/bittype/issues/161) |

Finding 6 in detail. 22 imports go from a production package into a stage
provider package. Eight use only an identity (`Contract`, `Profile`, a relation
key, `PreludePin` or the `Compiler` configuration type), one uses the
co-located query service of finding 5, and thirteen are the calls of finding 9.
Go, TypeScript and Rust have `semantics/<language>/go` for their vocabulary and
still publish `Contract()` from the resolver package. horos has no such package:
most of its names are defined in `internal/horosmodel` and re-exported from the
resolver.

## Rules to enforce

A package's role follows from its path, which the [layout](../LAYOUT.md) fixes.

| Role | Paths |
| --- | --- |
| Foundation | `trees/*/go`, `files/go`, `source/go`, `semantics/go`, `contracts/*/go` |
| Kernel | `engine/go`, `transform/syntax/shared/*/go`, `check/syntax/shared/*/go`, `print/*/shared/provider/go` |
| Vocabulary | `syntax/<language>/go`, `semantics/<language>/go`, `transform/syntax/<native>/model/go`, `transform/syntax/typescript/symbol/go` |
| Contribution | `<stage>/<dimension>/<language>/provider/go`, `transform/syntax/<scope>/{source,target}/go`, each mapping and each verification law |
| Private part | `internal/*`, each bound to one owner |
| Composition | `integration/go`, `verification/integration/go`, `tools/**`, a future `cmd/**` |

1. A foundation or kernel package imports no vocabulary and no contribution.
2. A vocabulary package imports foundations and its own language's vocabulary,
   never a contribution.
3. A contribution imports foundations, kernel contracts, the vocabularies of
   its own scope and its own private parts.
4. Only a composition imports a contribution.
5. Every exception to rule 3 or 4 names the importer, the imported package, the
   exact exported symbols and the reason. A use outside the named symbols
   fails. An unused exception fails, and so does an unused symbol in one.
6. A production package whose path gives no role fails.

Rule 4 is the one that stops direct calls. If nothing but a composition can
import a contribution, no contribution can call another one.

Rules 1 to 3 are largely what the existing test already enforces, written once
per role where the test now writes them once per package. A new package of an
existing role then needs no edit to the test.

### What the check cannot establish

- **A decision made inline.** Record building `std::optional` imports nothing
  forbidden. That needs the contract change of #196 and a reviewer who asks the
  question above.
- **Whether a package holds one contribution.** That is a judgment. The scan
  lists candidates, and it over-reports: a setting key that a provider reads
  through a helper shows as a separate group.
- **Correct ownership.** A green run shows that the forbidden edges are absent
  ([AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07)).

The [agent instructions](../AGENTS.md) say not to replace reasoning about
responsibilities with keyword heuristics or new lint gates. These rules do not
replace it. They enforce boundaries that the reasoning has already decided, in
the test that already enforces directed imports.

## How the check is built

- It stays in `integration/go/architecture_test.go`. `go list -deps -json`
  already runs there and gives each package's imports and name.
- Symbol use comes from `go/parser` over each package's non-test files:
  resolve the import names per file, then collect the selectors on them. The
  standard library is enough, and `go.mod` gains no requirement.
- The check is a function over package data. Each rule has a test that feeds
  it a package set breaking exactly that rule and expects the failure.
- A dot import of a module package fails, because it hides symbol use.

## Order of work

1. **[#219](https://github.com/Bitspark/bittype/issues/219): name the symbols
   and fail unused allowances.** Moves nothing. It makes today's 22 provider
   imports explicit, nine of them as open exceptions that later steps remove.
2. **[#220](https://github.com/Bitspark/bittype/issues/220): identities move to
   the vocabularies.** Removes eight exceptions.
3. **[#221](https://github.com/Bitspark/bittype/issues/221): the query services
   get their own packages.** Removes one exception. After PR #204, which edits
   the same files.
4. **[#222](https://github.com/Bitspark/bittype/issues/222): one mapping
   contribution per package.** Before
   [#207](https://github.com/Bitspark/bittype/issues/207) writes a family into
   a mapping package, and together with #216.
5. **[#218](https://github.com/Bitspark/bittype/issues/218): rules by role.**
   The per-package list is replaced and deleted.

Steps 1 to 3 can be done in any order after step 1. Each step converts the
rules of the packages it touches to role rules.

No other platform repository imports a bittype package, and bittype has no
release tag (checked on 11 November 2026). The moves cost nothing outside this
repository today. That changes when `libraries/service` starts to consume the
generation packages.

## Repeating the scan

Two questions, answered from source without type checking:

1. Per package, which groups of top-level declarations never reference one
   another? Two or more groups that each export something are grouping
   candidates.
2. Per import of a module package, which exported symbols does the importer
   use? An importer that uses one identity from a provider package is a
   vocabulary candidate.

The tool's source is in [#218](https://github.com/Bitspark/bittype/issues/218).
It is a discovery aid and not a check. Read every candidate in source before
recording it.
