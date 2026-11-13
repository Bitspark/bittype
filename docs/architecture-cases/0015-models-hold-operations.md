# bittype/ACX-015: Models hold operations and import operations

**Status.** Identified on 13 November 2026 at main
`058140171592882239d759663e0f7387a171c021`. Correction pending. The
correction is designed in [models and operations](../models-and-operations.md)
and delivered by its slices S1 to S8. For the four intended-forms packages
the move on `main` is deferred until the migration of
[#196](https://github.com/Bitspark/bittype/issues/196) merges, with the
reason and the trigger recorded [below](#the-deferral-on-main).

**Scope.** The packages of `main` that declare what exists in a language or
in a generic carrier: the carriers, the vocabulary packages, the grammar
tables and the lexical word lists. Read in source on 13 November 2026 for
research 0018 and again for the design page. Nothing was built or run.

## The rule

The owner directed on 13 November 2026:

> So the root split to me would be between model and ops. model must not contain any ops and must not import from ops either. ops work on the model

The direction has two rules. **Rule 1:** a model imports no operation.
**Rule 2:** a model contains no operation. The test that tells the two
apart is decision M1 of the
[grade of research 0018](https://github.com/Bitspark/bittype/issues/305#issuecomment-6855741567):
a responsibility that declares what is true or admissible is model, and one
that determines whether something is true, interprets it or produces
something from it is an operation.

## Rules involved

[AC-01](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-01)
makes the agreed ownership and dependency contracts part of correctness. The
direction is such a contract, and the code does not meet it.

[AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)
keeps a meaning with its owner. A language's grammar table is a declaration
of that language. Today an operation owns it: the table is an unexported
variable of the language's checker.

[AC-11](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-11)
forbids a copy of an agnostic mechanism for each language. The checks of the
intended-forms packages are such copies.

[AC-06](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-06)
keeps coverage apart from the boundary. Two conditions in Go's
intended-forms package narrow the model to what today's generation builds.

[AC-05](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-05),
in the owner's clarification of 13 November 2026, asks that a relation
between a declaration and its implementation is stated. Where a constraint
exists only as its evaluator, no page states that the declared model is an
upper bound of what is admitted.

[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
orders the correction and governs the deferral.

## Intended boundary

A model package holds plain data, names and declared tables, and no
function that evaluates them. It imports model packages only. An operation
takes the model as an explicit input and does not own the declaration it
interprets. A type whose values only a refusing function creates is the
admitted state and lives with that function on the operations side.
[Models and operations](../models-and-operations.md#2-the-test) states the
test, and its [table](../models-and-operations.md#every-package-of-main)
places every package.

## Actual implementation and consequence

**Packages that hold both (rule 2).** By the design page's table 25 of the
129 production packages of `main` hold a model-side part and an operation.
One example of each group:

| Group | Example | The model-side part | The operation beside it |
| --- | --- | --- | --- |
| Intended-forms packages (4) | `transform/syntax/golang/model/go/model.go` | Seven form names (`:18-26`), two primitive names (`:30-33`) | The issuer (`:12`), five constructors (`:37-51`), `CheckExpression` (`:52-93`), `CheckDeclaration` (`:94-112`) |
| Semantic vocabularies (4) | `semantics/golang/go` | Profiles and 21 key functions (`vocabulary.go:13-71`), the payload shape `FieldInfo` (`model.go:34`) | `Open` and 33 reader methods (`model.go:42-316`), three contracts (`contracts.go:15-40`), the encoding of a compiler run's inputs (`environment.go`) |
| Symbol packages (3) | `transform/syntax/golang/symbol/go/symbol.go` | The record of a symbol and its category | Admission by asking the universe scope of the toolchain (`:33-43`), the issuer (`:48`), `ExternalUse` and `ReadExternal` (`:56-81`) |
| Carriers (3) | `semantics/go/model.go` | The plain records and `Data` (`:20-125`) | The sealed graph and `New` (`:146`, `:230-405`) |
| Checkers with their grammar table (4) | `check/syntax/golang/provider/go/check.go` | The table of 23 entries (`:17-41`) | `Check` with two rules between nodes (`:45-83`) |
| Lexical helpers (4) | `internal/rust` | The reserved-word list (`names.go:11`), a table of character ranges (`cratename.go:48`) | The evaluators and the literal decoders |
| Grammar and horos framing (2) | `internal/grammar/rules.go` | The table types (`:15-33`) | The interpreter `Check` (`:80`) and the predicates |
| The common horos mapping package (1) | `transform/syntax/horos/mapping/go/mapping.go` | Two profile identities (`:28`, `:34`) | The function type `Library[T]`, which is also [bittype/ACX-010](0010-agnostic-mapping-adapter.md) |

**Imports from a model-side package into an operation (rule 1).** The
sixteen vocabulary packages make 24 imports of a contract package: five
syntax vocabularies import `contracts/language/go`, four semantic
vocabularies import `contracts/semantic/go`, `contracts/resolution/go` and
`contracts/analysis/go`, and the four intended-forms packages and the three
symbol packages import `contracts/generation/go`. The five into the language
contract become imports between models, because language identity is model.
The other 19 are imports of an operation. `source/go` imports
`contracts/semantic/go` for the two receipts of a publication
(`source/go/source.go:38-44`). Four packages that hold model-side parts
import a private lexical helper.

**An operation owns a declaration** (AC-02). The four native grammar tables
are unexported variables of their checkers (`check.go:17` for Go), and the
horos table is in `internal/declsyntax/check.go:67`. Go's table binds 17
predicate function values, so it is not data today. The Go compiler bridge
writes eleven keys as literals that the vocabulary does not publish
(`grep -oh 'n\.Key("[a-z0-9-]*")' internal/golangsem/*.go | sort -u | wc -l`).

**One mechanism, a copy for each language** (AC-11). The four intended-forms
packages hold the same issuer, the same refusal of a foreign form, the same
arity checks and the same record rule. A fifth copy, for horos, is on the
branch of the two other pairs. Decision
[O46 on #196](https://github.com/Bitspark/bittype/issues/196#issuecomment-6855097337)
lists them as AG-7.

**A model narrowed to what one operation supports** (AC-06). Go's
intended-forms package refuses a primitive name other than `string` and
`bool` as invalid (`model.go:63-64`), and a `defined` declaration over a
body that is no record as invalid (`:105-110`). The first is what today's
target lowers. The second is the declaration policy's choice, which the Go
target does not need: it lowers a `defined` declaration over any body
(`transform/syntax/golang/target/go/provider.go:96-106`).

**The migration adds to it.** On `origin/feat/contribution-contract-196` at
`af1dcde` each of the four intended-forms packages gains a third check,
`CheckMember` (`transform/syntax/golang/model/go/model.go:122` on that
branch), and 24 contribution vocabularies and six family packages import
`contracts/generation/go` for their identity types.

## What worked and what it did not establish

The architecture test is green, and it is right to be: every import above is
one that a role permits. The role Vocabulary lists "Foundations and its own
language's vocabulary" (`LAYOUT.md:158`), and the contract packages are
foundations. So a green run says nothing about this boundary
([AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07)).
The accepted definition of a vocabulary already said "The names two parties
must agree on and nothing else" (`docs/package-boundaries.md:98`), while two
accepted pages described the intended-forms package as holding "validation"
and "validity" (`LAYOUT.md:67`, `docs/generation-composition.md:17`). The
corrections of [bittype/ACX-009](0009-package-boundaries.md) stand for their
recorded scope: they moved identities into vocabulary packages and did not
ask what else those packages hold. Its limits name the vocabulary packages'
own readers and the admission grammars of the checkers as not examined.

## Required correction

The design page gives it in full. In short:

1. `CheckMember` leaves the four intended-forms packages on the migration's
   branch at once (slice S1).
2. Each intended-forms package becomes a declared table of forms on the
   model side. One checker and one constructor in the agnostic layer read
   it. The target keeps lowering and coverage. Go's record-body condition
   goes to the contribution that supplies the declaration (slice S2).
3. The grammar tables and their type go to the model side, with predicates
   named by identity (slice S3). The word lists leave the lexical helpers
   (slice S4).
4. The typed readers and the contract identities leave the semantic
   vocabularies, and the symbol packages split (slice S5).
5. The plain records of the carriers part from admission, and language
   identity becomes a model (slice S6).
6. The identity types of a contribution get a model-side owner (slice S7).
7. The scripted move under the two roots, the roles for each root and the
   import rule of rule 1, with a test that breaks it (slice S8).

Rule 2 gets no check. Whether a package under the model root holds an
operation is judged in review, by the owner's rule against keyword
heuristics and lint gates in place of reasoning.

## The deferral on `main`

[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10):
"Deferring the move itself is permitted only when its destination needs a
huge effort. Record such a deferral in the case with its reason and the
trigger that ends it; the violation stays open."

**What is deferred.** The move of the checks and the constructors out of
`transform/syntax/{cpp,golang,rust,typescript}/model/go` on `main`.

**The reason.** The move is not put off. It is made once, on the migration's
branch, as its own change and before the migration's gate is shown. It is
not made on `main` first, because every file that imports one of the four
packages is a file that the branch rewrites: over all Go files of `main`,
`git grep -l 'transform/syntax/<language>/model/go"' -- '*.go'` gives 9
files for Go, 9 for TypeScript, 7 for Rust and 6 for C++, and each is among
the files that `git diff --name-only origin/main...origin/feat/contribution-contract-196`
lists. A move on `main` would be made against code that the migration
deletes, and made a second time on the branch. Decision M8 of the grade
gives this reason. It differs from the order that the verifier proposed,
which was `main` first.

**The trigger.** The migration of #196, pull request
[#308](https://github.com/Bitspark/bittype/pull/308), merges into `main`.
From that merge on `main` holds the split. Until then `main` holds the
misplacement, and this case stays open for it.

**What the deferral does not cover.** Slices S3 to S6 touch no file of the
migration's branch and are not deferred.

## Correction and qualification status

This change documents the violation and the designed correction. No code
has changed. The design page holds nine proposals that are not decided,
among them the names below the two roots, so the paths it gives are
proposed paths. The classification of each package is the design page's
table, made by reading. Each slice updates this case with what it delivered,
its qualification and its limits.

## Assessment lesson

The role of a package was read from its path and checked by its imports,
and both were in order. What the package held was described by the pages
that listed it, and those pages followed the code: a package of names was
described as "forms and validation" once it validated. A definition that the
code contradicts is a finding. The pages recorded the contradiction twice
and treated it as a description.

[Models and operations](../models-and-operations.md) ·
[Package boundaries](../package-boundaries.md) · [Case index](README.md)
