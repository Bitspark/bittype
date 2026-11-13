# bittype/ACX-015: Models hold operations and import operations

**Status.** Identified on 13 November 2026 at main
`058140171592882239d759663e0f7387a171c021`. Correction in progress: slice S3a
is delivered on `main` with pull request [#318](https://github.com/Bitspark/bittype/pull/318), slice S4
with pull request [#320](https://github.com/Bitspark/bittype/pull/320), and
the rest is pending. The
correction is designed in [models and operations](../models-and-operations.md)
and delivered by its slices S1 to S8. For the four intended-forms packages
the correction is made on the branch of the migration of
[#196](https://github.com/Bitspark/bittype/issues/196), and `main` shows the
misplacement until that branch merges. The reason and the two conditions are
recorded [below](#where-the-correction-of-the-intended-forms-packages-is-made).

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
4. The type of a grammar table goes to the agnostic syntax model as plain
   data, with predicates named by identity, and each language gets a
   lexical part below its syntax model, which imports no node kinds and
   holds the names of its predicates (slice S3a, delivered). The word lists
   leave the lexical helpers for the lexical parts (slice S4, delivered). The grammar
   tables go to the model side, each into its language's syntax vocabulary
   (slice S3b). Decision [M31](https://github.com/Bitspark/bittype/issues/305#issuecomment-6860518608)
   cut the three slices so.
5. The names, profiles and payload shapes of the semantic vocabularies get
   model packages, the Go bridge writes its eleven keys through published
   names, and the TypeScript symbol record gets its model package (slice
   S5, the part that runs now). Taking the typed readers out of the semantic
   vocabularies, and the path of symbol admission, are the part of S5 that
   waits, recorded [below](#what-waits-for-a-decision).
6. The records of the carriers part from admission (slice S6).
7. On the pairs' branch, when it merges the migration: the vocabularies and
   families of the two further pairs, their two declarations of the
   publisher name and the horos intended-forms package, as the pairs' part
   of slices S7 and S2.
8. The scripted move under the two roots, with language identity, the roles
   for each root and the import rule of rule 1, with a test that breaks it
   (slice S8). It waits for all of the above.

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

## What waits for a decision

**What it concerns.** Two parts of slice S5: taking the typed readers out
of the semantic vocabularies of Go, Rust and TypeScript, and the path of
symbol admission.

**It is no deferral.** No effort is claimed for the wait, so the deferral
clause of AC-10 does not cover it (decisions M27 and M29 on #305).

**The reason.** Each of the two is imported by several operations, and the
design first placed each under one of them: the readers under the stage
`analyze`, symbol admission under the stage `transform`. A rule that
places a responsibility by the knowledge it needs would decide both
paths, and it is not accepted yet. A package created now might have to
move again.

**The condition.** Platform research 0004, the consultation on that rule,
is evaluated, and the research owner decides the two paths. The scripted
move, slice S8, waits for it.

**What stays open until then.** The typed readers share a package with
the contract identities, `semantics/<language>/go`, after the names have
left it. This case stays open for that.

## Correction and qualification status

The change that added this case documented the violation and the designed
correction, and changed no code. The names below the two roots are decided (M10, M11, M14,
M19 and M28 on #305). Two paths of M14 wait for a decision (M27, M29). The form-table type is a sketch that slice S2 makes final (M13).
The classification of each package is the design page's
table, made by reading. Each slice updates this case with what it delivered,
its qualification and its limits.

### Slice S3a, delivered with pull request [#318](https://github.com/Bitspark/bittype/pull/318)

**What is corrected.**

- The table types `Attribute`, `Role`, `Rule` and `Rules` are no longer in
  `internal/grammar`. They are plain data in the agnostic syntax model,
  [`model/syntax/go`](../../model/syntax/go/syntax.go), the first package at
  a final path under `model`. An attribute entry holds no function value: it
  names its predicate as `Predicate{Owner, Name}` and holds a set of
  admitted texts as data. The six shared names are declared there.
- Each of the five languages has a lexical part,
  `model/syntax/<language>/lexical/go`, with the names of its own
  predicates, 17 in all. It imports the agnostic syntax model and nothing
  else.
- The five tables are composite literals of the model's type. None binds a
  function value any more. The entry constructors of `internal/grammar`
  are deleted.
- The interpreter `internal/grammar.Check` takes the evaluators of the
  caller's language as a map from name to function and binds the shared
  names itself. A table that names a predicate without an evaluator fails
  every check before a tree is read, as unsupported.
- The horos identifier rule has one evaluator, `internal/horos.Identifier`.
  The copy in `internal/declsyntax` is deleted.

**What stays open after it.** The five tables are still unexported
variables of the four native checkers and of `internal/declsyntax`, so an
operation still owns the declaration it interprets (slice S3b). The word
lists, character tables and patterns are still inside the lexical helpers,
and Rust's keyword list still stands a second time in
`semantics/rust/go` (slice S4). Rust's primitive type words and
TypeScript's keyword-type words still stand in the checker and again in
the parser (slice S3b). No page states yet that a declared table is an
upper bound of what its checker admits; that sentence is written with
slice S3b, in a limited form for horos, whose checker takes the references
out of a tree before the interpreter sees it. Everything else of this case
is unchanged.

**Qualification.** Each of the five tables was written out, before and
after, as its node kinds, attributes, roles with their bounds and kinds,
and for each attribute which of 630 probe texts it admits. The two lists
are equal for every language. The outcome of every test by name and the
generated and formatted outputs of the integration tests are unchanged.
The import test has an exact entry for `model/syntax/go` and a role for
the lexical part with an allow-list and a prohibition on every syntax
vocabulary, and each was seen to refuse a seeded import. The pull request
holds the runs.

**Limits.** The probe texts are a sample, and equal answers on them do not
prove equal predicates. The tables are data now, and they are not yet
reachable from another package. That a package under `model` holds no
operation is judged in review and by no check. The import rule for
everything under `model` is slice S8's: until then each model package has
its own exact list.

### Slice S4, delivered with pull request [#320](https://github.com/Bitspark/bittype/pull/320)

**What is corrected.**

- The four lexical helpers hold no list any more. The lexical part of each
  language declares them as data:

  | Language | Declared in `model/syntax/<language>/lexical/go` | Where it stood |
  | --- | --- | --- |
  | C++ | 92 reserved words, the 14 words of a builtin type | `internal/cpp/names.go`, `internal/cpp/types.go` |
  | horos | 41 reserved words, the bound of 4096 digits of an integer | `internal/horos/spelling.go` |
  | Rust | 52 reserved words, the 14 names of calling conventions, the 840 ranges of the characters of a crate name | `internal/rust/names.go`, `literal.go`, `cratename.go` |
  | TypeScript | 46 reserved words, the 11 predefined type words, the two patterns of a number as strings | `internal/typescript/names.go`, `literal.go` |
  | Go | Nothing: two names and no list | |

- Each evaluator in a lexical helper reads its list from the lexical part.
  The question "is this text one of the words" has one implementation,
  `Words` of the shared grammar helper
  ([AC-11](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-11)).
- Rust's keyword list stood a second time in `semantics/rust/go`. That
  copy is deleted, and its evaluator asks the lexical helper.
- A list is published by a function without input that returns a fresh
  value, so no importer can change it for another. The bound and the two
  patterns are constants.
- [The page of the syntax](../syntax.md#go) states the relation that
  [AC-05](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-05)
  asks for where the building toolchain decides: Go's keywords, and the
  letters and digits of an identifier in C++, Go, Rust and TypeScript.

**What stays open after it.** The five grammar tables are still variables
of an operation, and two word sets still stand in a checker and again in
a parser (slice S3b). No page states yet that a declared table is an upper
bound of what its checker admits (slice S3b). Everything else of this case
is unchanged.

**Qualification.** Every list was generated from the helper's source and
read back: equal word for word and in order, and the deleted Rust copy was
equal to the kept list as a set. Every evaluator of the four helpers gave
the same answer before and after on the probe texts, the crate-name rule
on every code point (144435 admitted), and the Rust environment the same
refusal for every probe text as a crate name. The five grammar tables,
the outcome of every test by name and the generated and formatted outputs
are unchanged. A changed entry of a list turns a test red for each of the
four languages. The pull request holds the runs.

**Limits.** A word that an evaluator compares inline is not declared: the
five words that a raw Rust identifier may not be, the path keywords of a
Rust path segment, the escape letters of a horos string, the words of a
C++ builtin type that stand alone. The design page counts such a rule as
one that exists only as its evaluator. The tests of the helpers state
expected answers for selected words of each list and ask of every declared
word what its evaluator says. Only Rust's keywords are compared in full
with a second, independent list, that of the test of its environment. The
table of crate-name characters is compared with the pinned compiler only
where that compiler is provisioned.

## Assessment lesson

The role of a package was read from its path and checked by its imports,
and both were in order. What the package held was described by the pages
that listed it, and those pages followed the code: a package of names was
described as "forms and validation" once it validated. A definition that the
code contradicts is a finding. The pages recorded the contradiction twice
and treated it as a description.

[Models and operations](../models-and-operations.md) ·
[Package boundaries](../package-boundaries.md) · [Case index](README.md)
