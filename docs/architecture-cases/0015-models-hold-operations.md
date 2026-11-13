# bittype/ACX-015: Models hold operations and import operations

**Status.** Identified on 13 November 2026 at main
`058140171592882239d759663e0f7387a171c021`. Correction pending. The
correction is designed in [models and operations](../models-and-operations.md)
and delivered by its slices S1 to S8. For the four intended-forms packages
the correction is made on the branch of the migration of
[#196](https://github.com/Bitspark/bittype/issues/196), and `main` shows the
misplacement until that branch merges. The reason and the two conditions are
recorded [below](#where-the-correction-of-the-intended-forms-packages-is-made).
Since 13 November 2026 the platform's rules page names this case as the
example of [AC-12](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-12).

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

The direction is bittype's. The platform's rule for every project is
[AC-12](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-12),
"A description belongs to its subject, not to an operation that uses it",
accepted on the platform's `main` on 13 November 2026 from platform research
0004. It covers a part of this case and not all of it, as the next section
says.

## Rules involved

[AC-01](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-01)
makes the agreed ownership and dependency contracts part of correctness. The
direction is such a contract, and the code does not meet it.

[AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)
keeps a meaning with its owner. A language's grammar table is a declaration
of that language. Today an operation owns it: the table is an unexported
variable of the language's checker. AC-02 alone does not reach this: the
checker is the language's own contribution, so the meaning is with its
semantic owner. AC-12 is the rule that makes it assessable. Its
relationship to AC-02 reads: "This rule makes it assessable where one
owner supplies both the description and the operation that uses it first,
as a language's own checker that holds that language's grammar does.
AC-02's violation does not reach that case."

[AC-12](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-12)
gives a description to a unit of its subject: "No operation owns it, and
nothing it states refers to an operation or to an operation's request,
result or state." Three things in this case are what it covers:

- **The grammar tables are private data of their checkers.** The rule's
  first violation is "A description that is private data of an operation".
  A second operation cannot read a table without the checker: the
  documentation tool reads the checker's source for it, "without importing
  or executing inspected packages"
  (`tools/provider-matrix/extract/grammar/go/extract.go:2`).
- **The Go bridge writes eleven keys that the vocabulary does not
  publish.** They are names of Go's analysed facts, stated only inside the
  operation that writes them. A reader gets them by "its own copy where it
  could depend on the description".
- **Descriptions stand under a stage's name.** The intended forms of a
  language and the record of a TypeScript symbol are under
  `transform/syntax/<language>/`. "Its place and name, meaning the path or
  published name by which dependents refer to it, state no operation,
  because every dependent would then name that operation." The TypeScript
  Files formatter names the stage `transform` to name a symbol.

The rule also names what this case files under AC-06 and AC-02 below: "A
limit of one operation's coverage (AC-06) or a policy of one operation
(AC-02) stays with that operation."

Two things in this case are **not** AC-12. They are the owner's direction
for bittype, and the rule says so: "bittype has the direction from its
owner (13 November 2026) that a unit of descriptions holds no operation
and imports none; that direction is not extended to other projects here."

- **Containment, rule 2.** That a package of names also holds checks,
  typed readers or admission, where the names themselves are published.
  AC-12 leaves open whether operations share a description's unit.
- **The import rule on units, rule 1.** That a package with a model-side
  part imports a package of an operation. AC-12 asks only that nothing a
  description states refers to an operation.

[AC-11](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-11)
forbids a copy of an agnostic mechanism for each language. The checks of the
intended-forms packages are such copies.

[AC-06](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-06)
keeps coverage apart from the boundary. Two conditions in Go's
intended-forms package narrow the model to what today's generation builds:
the set of primitive names and the refusal of every external type. A third
condition, that a defined type has a record body, is a policy of the
declaration contribution in the wrong package, which is AC-02.

[AC-05](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-05),
in the owner's clarification of 13 November 2026, asks that a relation
between a declaration and its implementation is stated. Where a constraint
exists only as its evaluator, no page states that the declared model is an
upper bound of what is admitted.

[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
orders the correction: the misplaced part moves first, as its own change.

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

**Packages that hold both (rule 2).** By the design page's table 23 of the
129 production packages of `main` hold a model-side part and an operation.
One example of each group:

| Group | Example | The model-side part | The operation beside it |
| --- | --- | --- | --- |
| Intended-forms packages (4) | `transform/syntax/golang/model/go/model.go` | Seven form names (`:18-26`), two primitive names (`:30-33`) | The issuer (`:12`), five constructors (`:37-51`), `CheckExpression` (`:52-93`), `CheckDeclaration` (`:94-112`) |
| Semantic vocabularies (4) | `semantics/golang/go` | Profiles and 21 key functions (`vocabulary.go:13-71`), the payload shape `FieldInfo` (`model.go:34`) | `Open` and 33 reader methods (`model.go:42-316`), three contracts (`contracts.go:15-40`), the encoding of a compiler run's inputs (`environment.go`) |
| The TypeScript symbol package (1) | `transform/syntax/typescript/symbol/go/symbol.go` | The record of an exported type, `Type{Module, Name}` (`:15`) | The issuer (`:17`), `ExternalUse` and `ReadExternal` (`:20-43`). The symbol packages of Go and Rust hold sealed types only and are wholly operations |
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
Two are made by the symbol packages of Go and Rust, which are wholly
operations. The other 17 are imports of an operation by a package that holds
a model-side part. `source/go` imports `contracts/semantic/go` for the two
receipts of a publication (`source/go/source.go:38-44`). Three packages that
hold model-side parts import a private lexical helper.

**An operation owns a declaration** (AC-02, made assessable by AC-12). The four native grammar tables
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
`bool` as invalid (`model.go:63-64`), and every external type as unsupported
(`:87-88`). Both say what today's target lowers. The first stands in the
three other packages too.

**A policy in the wrong package** (AC-02). The same package refuses a
`defined` declaration over a body that is no record as invalid
(`:105-110`). That is the declaration policy's choice, which the Go target
does not need: it lowers a `defined` declaration over any body
(`transform/syntax/golang/target/go/provider.go:96-106`).

**The migration adds to it.** On `origin/feat/contribution-contract-196` at
`356c220`, 24 contribution vocabularies and six family packages import
`contracts/generation/go` for their identity types, and the 24 vocabularies
import the common mapping package for the publisher name. At `af1dcde` each
of the four intended-forms packages also held a third check, `CheckMember`.
Since commit `8c96c78` of the branch it is a method of each target
(`transform/syntax/golang/target/go/provider.go:138` on the branch).

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
   branch at once (slice S1). Delivered there with commit `8c96c78`.
2. The identity types of a contribution get a model-side owner, which also
   holds the publisher name. Their six methods become functions of the
   generation contract first. The contribution vocabularies, the families and
   the profile identities follow them (slice S7, before S2).
3. Each intended-forms package becomes a declared table of forms on the
   model side. One checker and one constructor in the agnostic layer read
   it. The target keeps lowering and coverage. Go's record-body condition is
   dropped, and no code is written for it anywhere: the declaration policy
   makes that choice already (slice S2).
4. The grammar tables go to the model side, each into its language's
   syntax vocabulary, with predicates named by identity. Their type goes to
   the agnostic syntax model (slices S3a and S3b). The word lists leave the lexical
   helpers for a lexical part below each language's syntax model, which
   imports no node kinds (slice S4).
5. The names, profiles and payload shapes of the semantic vocabularies get
   model packages, the Go bridge writes its eleven keys through published
   names, and the TypeScript symbol record gets its model package (slice
   S5, the first part). The typed readers leave the semantic vocabularies,
   and symbol admission gets its path, in the second part of S5, recorded
   [below](#what-waited-for-a-decision).
6. The records of the carriers part from admission (slice S6).
7. On the pairs' branch, when it merges the migration: the vocabularies and
   families of the two further pairs, their two declarations of the
   publisher name and the horos intended-forms package, as the pairs' part
   of slices S7 and S2.
8. The scripted move under the two roots, with language identity, the roles
   for each root and the import rule of rule 1, with a test that breaks it
   (slice S8). It waits for all of the above. Its table also moves the
   symbol packages of Go and Rust to their path.

Rule 2 gets no check. Whether a package under the model root holds an
operation is judged in review, by the owner's rule against keyword
heuristics and lint gates in place of reasoning.

## Where the correction of the intended-forms packages is made

**What it concerns.** The move of the checks and the constructors out of
`transform/syntax/{cpp,golang,rust,typescript}/model/go`.

**It is no deferral.**
[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
permits the deferral of a move "only when its destination needs a huge
effort". No such effort is claimed here, and the move is not put off. It is
made once, on the migration's branch, as its own change and before the
migration's gate is shown. This section records a decision about where the
correction is made (decisions M8 and M23 on #305).

**The reason.** Every file that imports one of the four packages is a file
that the branch rewrites: over all Go files of `main`,
`git grep -l 'transform/syntax/<language>/model/go"' -- '*.go'` gives 9
files for Go, 9 for TypeScript, 7 for Rust and 6 for C++, and each is among
the files that `git diff --name-only origin/main...origin/feat/contribution-contract-196`
lists. A move on `main` would be made against code that the migration
deletes, and made a second time on the branch. It differs from the order
that the verifier proposed, which was `main` first.

**Until when `main` shows the misplacement.** One of two things ends it, and
this case stays open until then:

1. The migration of #196, pull request
   [#308](https://github.com/Bitspark/bittype/pull/308), merges into `main`.
   From that merge on `main` holds the split.
2. Pull request #308 closes unmerged. Slice S2 is then made on `main`.

**What it does not concern.** Slices S3 to S6 split no package that the
migration's branch changes. They change files of the branch that import a
moved name, and the branch merges `main` after each split.

## What waited for a decision

**What it concerns.** Two parts of slice S5: taking the typed readers out
of the semantic vocabularies of Go, Rust and TypeScript, and the path of
symbol admission.

**It was no deferral.** No effort was claimed for the wait, so the deferral
clause of AC-10 did not cover it (decisions M27 and M29 on #305).

**The reason.** Each of the two is imported by several operations, and the
design first placed each under one of them: the readers under the stage
`analyze`, symbol admission under the stage `transform`. A package created
at such a path might have had to move again.

**The decision.** Platform research 0004 is evaluated, and the research
owner decided the two paths as
[M30](https://github.com/Bitspark/bittype/issues/305#issuecomment-6860238626): below
`ops`, a unit that one stage uses is in the stage matrix, and a unit that
several stages use is at the path of the sort and scope it works on. The
typed readers go to `ops/semantics/<language>/reader/go`, and symbol
admission to `ops/symbols/<language>/go`. [M32](https://github.com/Bitspark/bittype/issues/305#issuecomment-6861647384)
says how that is read: the substitution governs, and not a count of today's
users, so the readers and the symbol packages of all three languages take
that placement. AC-12 does not decide these
paths: the units are operations. What decides them is the substitution
that the rules page asks of a place.

**What stays open until the second part of S5 is done.** The typed readers
share a package with the contract identities, `semantics/<language>/go`,
after the names have left it. This case stays open for that. One point of
M30 is open, as M32 states it, and is decided at the start of that part:
which parts of the operations side of a symbol package work on symbols,
which are generation's use of symbols, and what the first kind refuses with
when it may not name the generation contract.

## Did we have the rule

The owner asked whether the rule existed and was violated. Three facts,
each read from the history of the repository that holds it.

1. **bittype's own definition existed and its packages already broke it.**
   "The names two parties must agree on and nothing else" entered
   `docs/package-boundaries.md` with commit `d23542a` on 11 November 2026.
   On that day Go's intended-forms package had held `CheckExpression` since
   commit `c054f80` of 8 November 2026, and Go's semantic vocabulary had
   held `Open` and its typed readers since its first commit, `a775892` of
   6 November 2026.
2. **AC-06 existed when the narrowing conditions were written.** The
   platform's rules page has AC-06, "Limited coverage does not redefine the
   architecture", since commit `7e588ae` of 6 November 2026. The primitive
   set and the refusal of every external type were written into Go's
   intended-forms package with commit `c054f80`, two days later
   (`transform/syntax/golang/model/go/model.go:42`, `:67` at that commit).
3. **The method page held the substitution tests when the grammar tables
   were placed.** The platform's
   [tests for placing a contract](https://github.com/Bitspark/platform/blob/main/docs/concepts/semantic-ownership-and-generalization.md#tests-for-placing-a-contract)
   stand since commit `5192506` of 6 October 2026, among them "Application
   substitution": "If bitwork is replaced by an unrelated hosted
   application, which promises must bithost still provide?". They were
   written for layers of services. The four native grammar tables became
   variables of their checkers with commit `09d74e8` of 31 October 2026
   (`check/golang/go/check.go:16` at that commit).

So the definition, the rule on coverage and the test existed, and none was
applied to the place of a description inside this library. The grade of
platform research 0004 says it in one sentence: "we had the test and most
of the rule, and did not apply them inside a library". What did not exist
was a rule that states the ownership of a description where the language
supplies both the description and its first user. That is AC-12.

## Correction and qualification status

This change documents the violation and the designed correction. No code
has changed. The names below the two roots are decided (M10, M11, M14,
M19, M28, M30 and M32 on #305). The form-table type is a sketch that slice S2 makes final (M13).
The classification of each package is the design page's
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
