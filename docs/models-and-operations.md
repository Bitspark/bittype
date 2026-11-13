# Models and operations

**Design, 13 November 2026, from research 0018. The decisions M1 to M18 are the research owner's, `claude-contract-chain-issues`: M1 to M9 in the [grade](https://github.com/Bitspark/bittype/issues/305#issuecomment-6855741567) on #305, M10 to M18 in the [decisions of the same day](https://github.com/Bitspark/bittype/issues/305#issuecomment-6856947056) on nine choices that the first version of this page left to that owner, and [M19](https://github.com/Bitspark/bittype/issues/305#issuecomment-6857366489), which adjusts M11. Nothing on this page is implemented: the code on `main` has the layout that [`LAYOUT.md`](../LAYOUT.md) describes.**

The owner directed on 13 November 2026 that the repository splits at its root
into models and operations. This page says how the line is drawn, where every
package of `main` goes, which roles exist under each root and in what order the
work is done. A reader should be able to place any package and split any mixed
package from this page alone.

An expert consultation was made on the direction. Its question, the advice,
two verification reports against the code and the graded table of 39
recommendations are in
[research 0018](../research-docs/0018-models-and-operations.md#applied)
([advice](../research-docs/0018-models-and-operations.advice.md),
[verification of rows R1 to R20](../research-docs/0018-models-and-operations.verify-1.md),
[verification of rows R21 to R31](../research-docs/0018-models-and-operations.verify-2.md)).
This page does not repeat them.

**How to read the sources.** Each statement names its source.

- **The direction:** the owner's messages 7 and 8 of 13 November 2026, quoted
  in [section 1](#1-the-direction-and-its-standing).
- **An owner rule:** a numbered rule AC-01 to AC-11 of the platform's
  [architectural correctness](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#numbered-rules),
  or a dated rule in [`AGENTS.md`](../AGENTS.md) or [`LAYOUT.md`](../LAYOUT.md).
- **An accepted page:** a page of this repository on `main`, quoted with its
  file and line at `0581401`.
- **M1 to M9:** the nine decisions of the
  [grade](https://github.com/Bitspark/bittype/issues/305#issuecomment-6855741567).
  bittype has no architecture seat, so the research owner graded and said so
  there.
- **Row R1 to R31:** a row of the table "Applied" of research 0018. Its Action
  is the decision. Where it differs from the advice, the Action governs.
- **Derived here:** a consequence that this page draws from the sources
  above. It is no decision. An objection to it is weighed.
- **M10 to M18:** the nine
  [decisions of 13 November 2026](https://github.com/Bitspark/bittype/issues/305#issuecomment-6856947056)
  on the choices that no source above settles. The first version of this
  page marked each with its alternatives, and the research owner decided
  them. Each decision is stated where it applies, with its reason and with
  what was not taken. The same comment corrects the counts of M7.
- **M19:** the [decision of the same day](https://github.com/Bitspark/bittype/issues/305#issuecomment-6857366489)
  that adjusts M11, made after this page had checked M11 against the code.
  [Section 11](#11-decisions-and-deferrals) lists all nineteen.

Code is quoted with its file and line on `main` at `0581401`. The existing
code is given as a fact and never as a reason for a choice
([AC-09](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-09)).
Counts come with the command that gives them. The change that adds this
page inserts dated clarifications into `LAYOUT.md`, `AGENTS.md`,
`CHARTER.md`, `docs/package-boundaries.md` and
`docs/generation-composition.md`. A line number of those five files is
therefore that of `0581401`, before the
insertions.

## 1. The direction and its standing

The owner wrote on 13 November 2026, in two messages of a longer exchange that
the [research document](../research-docs/0018-models-and-operations.md#the-framing-the-owners-ten-messages)
quotes in full.

Message 7, a layout and a question:

```text
why don't we organize it like this:
model/...
ops/...
model never imports from ops
model:
model/{syntax,semantics,...}/go/...
ops:
ops/{transform,...}/go/...
and what is currently not declarative about our models? like, don't we already have an ast / cst model? why do we then need code for the model for the different languages even?
```

Message 8, the direction:

```text
So the root split to me would be between model and ops.
model must not contain any ops and must not import from ops either.
ops work on the model
```

| What | Standing |
| --- | --- |
| The repository root splits into `model` and `ops` | Direction. Fixed |
| **Rule 1 (no import):** a model imports no operation | Direction. Fixed |
| **Rule 2 (no containment):** a model contains no operation | Direction. Fixed |
| Operations work on the model | Direction. Fixed |
| The names below the two roots in message 7, `model/{syntax,semantics,...}/go/...` and `ops/{transform,...}/go/...` | A sketch. [Section 7](#7-the-layout) gives the order below a root (M6) and the names (M10, M11, M14, M19) |
| Messages 9 and 10, "you are talking about the correctness of a model, aren't you?" and "or consistency?" | Questions of the owner. [Section 3](#3-consistency-and-correctness) answers them with M2 |

The direction changes two older rules of the same owner at the root: "organize
by processing stage and enforce a uniform structure across languages"
(`LAYOUT.md:3-4`, 31 October 2026) and the path pattern
`<stage>/<dimension>/<scope>/<component>/go` (`LAYOUT.md:11`, 6 November
2026). Both stay below the operations root (M6, row R22).

The expert approved the direction "without qualification" and added the
sentence that the grade takes as its centre: "I would not approve a migration
that merely moves existing packages under two new directory names."
(advice, closing recommendation; row R31). So this page is about
responsibilities first and about paths second.

## 2. The test

**The test (M1, row R1).** One question is asked of one responsibility:

> Does this responsibility declare what is true or admissible, or does it determine whether something is true, interpret it, or produce something from it?

A responsibility that declares is model. One that determines, interprets or
produces is an operation. "A model may declare a constraint. It must not
execute the mechanism that evaluates that constraint." (advice, section 1.1.)
The test is applied to a responsibility and not to a package: a package that
holds both sorts is split.

**Two sentences that M1 adds,** because the bare test leaves real code open:

1. **Sealed types.** A type whose values only a refusing function creates is
   the admitted state. It lives with that function on the operations side.
   Model-side types are plain data. Go allows nothing else: a type with
   unexported fields and its only constructor cannot be in two packages.
2. **A constraint that exists only as its evaluator** is named by the model
   and stated by the evaluator. The declared model is then an upper bound of
   what is admitted: everything an operation admits satisfies the declaration,
   and a declaration alone does not show that a value is admitted.
   [AC-05](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-05)
   asks for exactly this statement: "State the required relationship, such as
   a conservative upper bound rather than equality".

**What a model package may hold and what it may not** (advice, section 1.2,
as rows R1 and R2 adapt it).

| A model package may hold | A model package holds none of these |
| --- | --- |
| Type declarations that are plain data | A parser, a checker or an interpreter |
| Names and vocabulary | A reader that reconstructs meaning from an encoded artifact |
| Declared constraints and tables, as data | A planner |
| Identity and revision of the model itself | A constructor or a builder: it is construction, and a refusing one seals a type (sentence 1) |
| Relations between declared concepts | A validator that executes a declaration |
| | An identity that is a value of a contract package's type (row R2). It goes to a vocabulary on the operations side ([section 8](#8-roles-per-root)) |

Derived here, from the test: a Go function that takes no input and returns one
fixed value is the form of a declared name, because Go has no constant of a
struct type. `SliceOperator()` (`semantics/golang/go/vocabulary.go:70`) and the
methods `ID()` and `Name()` of a language marker
(`syntax/golang/go/nodes.go:16-20`) declare and determine nothing. They are
model. The same holds for a method that returns the receiver or one of its
fields to satisfy a view interface (`trees/ast/go/view.go:11`,
`trees/cst/go/cst.go:30`): it declares a relation between two carriers.

**Five kinds of code that the bare test leaves open** (row R1; verification
report 1, section 2.1 and row R1). Each has its side here. The reasons are
derived from M1 and the rows named.

| Kind of code | Example on `main` | Side | Reason |
| --- | --- | --- | --- |
| An issuer: it creates an authority and neither declares nor evaluates | `var forms = g.NewFormIssuer[s.Language]()`, `transform/syntax/golang/model/go/model.go:12` | Operations | Sentence 1: the issuer is the state of the refusing function. A model's identity is declared and is no process-local pointer (row R19, [section 6](#6-declared-models)) |
| A constructor that refuses nothing | `Native`, `contracts/generation/go/expression.go:76-78`; the five constructors at `model.go:37-51` | Operations | Row R7a: "Constructors that refuse nothing are construction, operations side" |
| A constraint that exists only as its evaluator | `func(v string) bool { return v == "." \|\| identifier(v) }`, `check/syntax/golang/provider/go/check.go:19`; a raw Rust identifier is none of five words, `internal/rust/names.go:16-18` | The name on the model side, the evaluator on the operations side | Sentence 2 |
| A value typed by a contract package that a vocabulary publishes | `ResolutionContract() resolution.Contract`, `semantics/golang/go/contracts.go:15`; `Evidence() semantic.EvidenceRequirement`, `vocabulary.go:96` | Operations, in a vocabulary of the operations side | Row R2. Contract packages are operations (row R8), and rule 1 forbids a model to import one |
| A plain query on a sealed value | `Facts`, `semantics/go/model.go:507`; `DirectParameterBindings`, `:524` | Operations, with the sealed type | Sentence 1 places the sealed type, and Go keeps a method with its type. A predicate method on a plain record (`Key.Valid`, `:26-28`) leaves the record's package as a function of the operation that needs it |

**The test applied to real packages.** Each row was read again at the source
for this page. Line counts are `wc -l` over the non-test files.

| Package | Model | Operations |
| --- | --- | --- |
| Go's intended-forms package, `transform/syntax/golang/model/go` (112 lines) | Seven form names (`model.go:18-26`) and two primitive names (`:30-33`). With [section 5](#5-intent): the form table | The issuer (`:12`), `Form` and `Name` (`:35-36`), five constructors (`:37-51`), `CheckExpression` (`:52-93`), `CheckDeclaration` (`:94-112`) |
| Go's semantic vocabulary, `semantics/golang/go` (855 lines) | Two profile strings and two further names (`vocabulary.go:13-16`), 21 key functions (`:22-71`), the payload shape `FieldInfo` (`model.go:34`), the plain records of a compiler run's inputs (`environment.go:25-56`) | `Model`, `Type`, `Field`, `Open` and 33 reader methods (`model.go:15-316`), the encoding of inputs into settings and back (`environment.go:57-400`), `Configuration`, `MatchesConfiguration`, `Requirements` and `Evidence` (`vocabulary.go:74-99`), the three contracts (`contracts.go:15-40`) |
| Go's syntax vocabulary, `syntax/golang/go` (107 lines) | All of it: a marker type, two methods that return a constant, node kinds, roles, attribute names, concrete kinds | Nothing |
| Go's grammar table, in `check/syntax/golang/provider/go/check.go` | Two kind lists (`:15-16`) and the table of 23 entries (`:17-41`), once its 17 predicate function values are named by identity | `identifier` (`:13`), `Check` with two rules between nodes (`:45-83`), the provider (`:85-89`) |
| The tree carriers, `trees/ast/go` (130 lines) and `trees/cst/go` (30 lines) | All of it: "These are data types, not checked trees. Structural and grammatical checks live in separate packages; the carrier performs no validation or allocation." (`trees/ast/go/ast.go:7-8`) | Nothing |
| The graph carrier, `semantics/go/model.go` (560 lines) | The plain records and `Data` (`:20-125`) | The limits (`:127-140`), the sealed `Semantics[L]` (`:146`), `New` (`:230-405`), the accessors and queries (`:445-560`) |
| The envelope package, `trees/representation/go` (490 lines) | Nothing | All of it. An envelope exists only through `Pack`, which runs the structure check (`ast.go:55-68`). The descriptor `Language` has unexported fields and a refusing constructor (`language.go:15-30`). `Info` (`info.go:9-15`) is the result of one operation |
| A contract package, `contracts/transform/go` | Nothing | All of it: a request, a contract, interfaces, admission |
| The Rust lexical helper, `internal/rust` (310 lines) | The reserved-word list (`names.go:11`) and the table of character ranges of a crate name (`cratename.go:48`) | `Identifier`, `Declaration`, `LifetimeParameter`, `LifetimeReference`, `PathSegment` (`names.go:15-29`), `CrateName` (`cratename.go:19`), the literal decoders (`literal.go`) |

The commands behind the counts, run at the repository root:
`grep -o 'g\.Attr([a-z]*, [^)]*' check/syntax/golang/provider/go/check.go | wc -l`
gives 17, `grep -c 'semantics.Key *{ *return Key(' semantics/golang/go/vocabulary.go`
gives 21, and `grep -c '^func (m Model) \|^func (t Type) ' semantics/golang/go/model.go`
gives 33.

**What follows for a contract package (row R8).** There is no third root. "A
contract describes the boundary through which an operation is used." (advice,
section 1.5.) A value with a meaning of its own outside the operation is
model. A value that only expresses an operation's invocation, execution or
result is the operation's. Of the thirteen contract packages one is model:
`contracts/language/go`, the two declarations of a language's identity (row
R3). A receipt that one operation issues and another consumes stays with its
producer (row R3).

## 3. Consistency and correctness

The owner asked whether the code that the test could not assign is about "the
correctness of a model", or about its "consistency" (messages 9 and 10). M2
answers with two words and four levels. No page stated them before.

- **Consistency:** a value agrees with its declared model. It is decided from
  the value and the declaration alone.
- **Correctness:** a value agrees with something outside its own declared
  structure: a compiler's judgment, or another model. It needs evidence.

| Level | The question | Nature | What establishes it today |
| --- | --- | --- | --- |
| 1. Structural consistency | Is the generic data structure sound: identities in range, nothing dangling, no cycle? | Consistency | `check.Structure` for a tree (`check/syntax/shared/ast/go/structure.go:66`), called by `Pack`. `semantics.New` for a graph (`semantics/go/model.go:230`). `Files.Capture` for a file set (`files/go/files.go:55`) |
| 2. Model consistency | Does the value conform to the language's declared model: node kinds, child roles, attributes, rules between nodes? | Consistency | The language's checker: the grammar interpreter over the table (`internal/grammar/rules.go:80`) and the rules in code beside it. For intent, the checks of the intended-forms package |
| 3. Language validity | Does the language's own implementation accept the program: names resolve, the program type-checks? | Correctness | Resolution and analysis through a compiler bridge. The result is carried as a receipt |
| 4. Correctness against another model | Does this Go type represent that horos type? | Correctness | Verification |

All four are decided by operations. A model states the declaration that
levels 1 and 2 are checked against. **No guarantee is inferred from another**
(M2, row R5). Accepted pages already say so for each step: "success
establishes neither reference resolution nor semantic validity"
(`docs/parsing-formatting.md:121`); "`Semantics[L].Valid()` establishes a
constructed graph, not a language-valid program" (`docs/semantics.md:206`);
generation "establishes neither native kind enforcement nor a proof that the
chosen representation is correct" (`docs/generation.md:94-95`).

**Artifact, evidence, admitted artifact (M2, row R6).** The advice asks for a
checked state apart from the artifact for every sort of model. No new concept
is added for it. The repository has the three, under these words:

| Carrier | The artifact, unchecked | The check | The admitted state | What the evidence records |
| --- | --- | --- | --- | --- |
| Syntax tree | `ast.AST[L]`, `cst.CST[L]` | Level 1: `check.Structure`. Level 2: the language's checker | Level 1: the envelope, "an owned, structurally checked, opaque snapshot" (`trees/representation/go/ast.go:23`). Level 2: the snapshot `semantic.Snapshot` (`contracts/semantic/go/operation.go:28`), for a set of trees | The envelope: language identifier, profile string, Go marker. The snapshot: language descriptor, entry, envelopes, Files context. Neither records which checker ran |
| Semantic graph | `semantics.Data` | `semantics.New`: integrity and bounds only | `Semantics[L]` and the envelope `representation.Semantics` | The receipt `semantic.Result`: the snapshot, the operation, its dependencies, the contract with its profile, the requested keys, the configuration, the limits, the reports and the diagnostics (`contracts/semantic/go/request.go:17-25`, `:98-104`). Admission is `AdmitEvidence` (`evidence.go:44-62`) |
| Files | `files.Files`, `files.Context` | `Files.Capture`, `Context.Capture`, `Cover` | No type of its own: `Capture` returns `Files` again (`files/go/files.go:55`). The sealed states are the snapshot's private context and the Source, "private sealed storage" (`docs/source.md:380`) | For the Source: the receipt, "Matching addresses or bytes cannot substitute for exact receipt association." (`docs/source.md:326`) |

Accepted sentences fix the pattern: a syntax envelope is "An owned,
structurally checked AST snapshot; Capture separately requires syntax
admission" (`docs/analysis.md:28`), and the snapshot "is not itself a
successful semantic judgment" (`docs/abstraction-levels.md:189`).

**Four gaps close** (row R6). Each is an addition and has its slice in
[section 10](#10-the-order-of-work).

1. **No native language declares its syntax profile.** horos does
   (`syntax/horos/go/nodes.go:26`). For the four native languages a
   composition passes the literal `native-types-v1`:
   `git grep -c 'native-types-v1' -- '*.go'` gives 11 occurrences in 9 files,
   all of them tests or test data.
2. **A parsed or formatted envelope does not record that grammar admission
   ran.** The Go parser calls its checker (`parse/syntax/golang/provider/go/parse.go:89`)
   and returns an envelope that records level 1 only.
3. **No check of fact shapes exists.** Nothing checks that a `go/struct` term
   has one operand for each entry of its payload. The reader refuses on read
   (`semantics/golang/go/model.go:289`). Such a check is separate from graph
   admission, because the kernel does not decode provider payloads: "the
   kernel does not decode them or infer semantic equality from equal bytes"
   (`docs/semantics.md:119-121`).
4. **Captured Files have no type of their own.**

## 4. One language's models

**A language has several related models (M3, row R9).** They answer
different questions and keep separate names.

| Model | What it says | Declared today in | Who writes a value of it | Who reads it |
| --- | --- | --- | --- | --- |
| Syntax | What can be written: node kinds, child roles, attributes, what each node kind may hold, which strings are words of the language | `syntax/<language>/go`, the grammar table inside the language's checker, the word lists inside `internal/<language>` | A parser, a formatter and a target, each through its stage | Every stage of the language, a mapping that names the language, the documentation tool |
| Analysed semantics | What analysis established about one actual program: entity kinds, term operators, relations, roles, profiles | `semantics/<language>/go` | The compiler bridge of the language | Typed readers, query services, verification laws, a source reader, a Files formatter |
| Intent | What a transformation intends to construct before a program exists: forms, their data and their arguments | `transform/syntax/<language>/model/go` | A mapping, through the language's constructors | The planner, which does not read a form's meaning, and the target, which lowers it |
| Symbols | What a name in generated output refers to when the referent is not itself generated | `transform/syntax/<language>/symbol/go` | A target | A Files formatter |

Not every language has all four. C++ has syntax and intent and no analysed
semantics and no symbols. horos has syntax and analysed semantics on `main`
and no intent. The command of [section 7](#7-the-layout) lists the packages.

**Where intent stands beside the five levels.**
[Abstraction levels](abstraction-levels.md) names five levels: encoding,
text, concrete syntax, abstract syntax and declaration semantics. Intent is
none of them (row R9). It is a statement about a target program for which
none of the five exists yet. It reaches the abstract-syntax level only through
a target's lowering, and the level of declaration semantics only through
analysis of the realized program. The accepted sentences: "Native intent is
not an analyzed `Semantics[T]`: target analysis establishes facts about the
actual constructed/organized program." (`docs/generation.md:77-79`), and "The
planned native program is intent, not a fabricated semantic-analysis result."
(`docs/generation-planner.md:278`). The correspondence table of that page has
the row already: "Selected mappings and scoped substitutions fix native
intent. This is not an analyzed target `Semantics[L]`."
(`docs/abstraction-levels.md:193`). So intent is a fourth sort of model beside
the levels and no sixth level.

**The stated relation and the test for each form (M3, row R10).** The three
sets of names stay apart. Their relation is stated, and one test for each
form holds it.

- The relation from a form to the nodes it becomes stays code, and it is the
  target's: "Target implementations are trusted interpreters of their own
  native forms" (`docs/generation-planner.md:535-536`).
- "Declared" means two things and no table: the relation is written on the
  language's page, and one test for each form builds the form, checks the
  tree, has it analysed and compares what analysis reports with what the
  form intended (row R10). The tests are an addition of slice S9.
- A relation as data waits for a reader. No package would read one today
  ([derive before extending](https://github.com/Bitspark/platform/blob/main/docs/concepts/semantic-ownership-and-generalization.md#derive-before-extending)).
- One description from which all three are derived is not built (row R12).
  Syntax has notions without a type-level counterpart, such as a file and an
  import, and one Rust form becomes a path of several nodes.

**When two models share a name (M3, row R11).** Only by the platform's
criteria, which are the binding wording: "Compare identity, lifetime,
authority, allowed operations, guarantees and failure behavior. A consumer
using the proposed common contract must retain the same guarantees when one
valid specialization is substituted for another. Differences must be
represented explicitly rather than erased."
(platform `docs/concepts/semantic-ownership-and-generalization.md:251-254`.)
For intent and analysed semantics the finding is the row "Independent
identities or lifetimes", and its response is "Keep separate concepts and
define their relationship." (`:261`). An intended reference names a
destination, and an analysed one names the entity a compiler bound. Intent
lives within one generation request, and a graph is bound to the receipt of
one captured program.

A shared identity for single notions, such as Go's slice, defined type and
alias, is **deferred** (row R13a). No reader exists. The trigger is a state
of the repository that anyone can check: a non-test file imports both a
language's intent vocabulary and its semantic vocabulary. For Go today

```sh
git grep -l 'transform/syntax/golang/model/go"' -- '*.go' | xargs grep -l 'bittype/semantics/golang/go"'
```

names one file, `integration/go/destinations_test.go`, a test. Two points of
row R13 belong here. A type parameter is nothing of Go's to share on the
intent side: "The kernel owns capture identity, substitution and scope."
(`docs/generation-planner.md:435`). And Go's analysed kinds `alias` and
`defined` are wider than the two declaration forms of the same names
(`semantics/golang/go/vocabulary.go:40-42`).

**The case where an analysed model could not be read as a source.** Decision
[O45 on #196](https://github.com/Bitspark/bittype/issues/196#issuecomment-6854327944)
(research owner, 13 November 2026) records it. The first pair with Go as its
source was refused with "evidence cycle" before any field became a use. The
kernel starts a declaration at its entity's body and follows operands. The Go
and the TypeScript analysis publish as the body a reference to the entity
itself, and the definition sits behind a fact of the language. O45: "The
definition, which the shared graph has a place for, sits behind a
language-owned fact. That is a misplaced responsibility in two analyses, not a
gap in the contract." The correction is
[#310](https://github.com/Bitspark/bittype/issues/310): both analyses publish
the definition as the body, under a new profile version. The advice of
research 0018 does not mention O45. That silence is no confirmation, and the
decision stands as the research owner's (the grade, "What the advice did not
answer"). What the case shows for this page: an analysed model is read by
more than the readers it was written for, so its declared shape, and not only
its names, is part of the model ([section 6](#6-declared-models)).

## 5. Intent

What the code calls `model`, the package
`transform/syntax/<language>/model/go`, is called the intended-forms package
here. M4 (rows R7i, R14, R15, R24, R26) decides what it becomes.

**Decided.**

| Part | What it becomes | Source |
| --- | --- | --- |
| The form names and the shape of each form | One declared table per language, as data, on the model side, with a declared identity | M4; row R15, parts 1 and 2; row R19 |
| The checks of well-formedness | **One** checker in the agnostic layer, which reads a language's table. It lives in the kernel package that holds the single-scope check, `transform/syntax/shared/target/go` (`check.go:36`) | M4; row R24. [AC-11](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-11): "Per-language copies of that mechanism are forbidden." |
| The constructors | **One** constructor in the agnostic layer. A language keeps typed constructors over it, in a package that mappings may import | M4; rows R15 and R24 |
| Lowering and coverage | The target keeps both. It passes its language's table to the checker and adds its own coverage answers | M4; row R15, part 5. The accepted contract gives the checks to the target: "The planner invokes its checks against intended native forms and the completed program." (`docs/generation-spi.md:332-333`) |
| The condition "a Go defined type has a record body" (`model.go:105-110`) | It leaves the package. It is neither Go's rule nor the target's. Its owner is the contribution that supplies the declaration, and nothing refuses it afterwards | M4; row R14; M17. "alias or defined declaration \| The contribution that supplies the declaration: the default declaration policy or an exact mapping" (`docs/contribution-contract.md:203`) |
| The intent carrier | Only the native kind of an expression, a form with its data and its arguments, is intent. The planner's other kinds stay one unit with the planner. No value type of intent is on the model side | Row R15, part 3, as M12 reads it. "Parts that share a private invariant belong together. The planner and the plan types in `contracts/generation/go` are one unit, because only the planner may issue a plan." (`docs/package-boundaries.md:145-147`) |

**Why one checker.** The four packages on `main` are near copies: the same
issuer, the same refusal of a foreign form, the same arity checks and the
same record rule (`transform/syntax/{cpp,golang,rust,typescript}/model/go/model.go`).
The pairs' branch adds a fifth for horos. Decision
[O46 on #196](https://github.com/Bitspark/bittype/issues/196#issuecomment-6855097337)
lists them as AG-7, "five near copies", and leaves them to this evaluation.
The mechanism needs no knowledge of one language, so AC-11 gives it one owner,
and each language supplies its table.

**What the Go target keeps as coverage.** Two answers of `CheckExpression`
are statements about what today's generation builds and no rules of a form:
which primitive names exist (`model.go:63`) and the refusal of every external
type (`:87-88`). A model is not narrowed to what one operation supports (row
R26; [AC-06](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-06)).
Rust keeps its coverage answer for external types in the target already
(`transform/syntax/rust/model/go/model.go:57-60`). A third answer, the
refusal of an unknown form name (`:89-90`), is an answer against the declared
table and belongs to the one checker (M17 below).

**One change of category, made on purpose (row R24, M17).** Failure meaning
is one of the guarantees that composition preserves (AC-05), so a changed
category is recorded and never a side effect. A primitive name outside the
two that the target lowers is invalid today (`model.go:63-64`). As the
target's coverage answer it is unsupported. Slice S2 records the change in
its pull request, with a test that names the old and the new answer.

**M17: two refusals after the split.**

- *(a) The record body of Go's `defined`.* A `defined` declaration over a
  body that is no record is invalid today (`model.go:107-108`). After the
  split **nothing refuses it**. Go's table states no body rule, the Go target
  lowers a `defined` declaration over any body
  (`transform/syntax/golang/target/go/provider.go:96-106`), and
  the choice between an alias and a defined declaration is the declaration
  contribution's. So this is no change of category: the refusal ends. Row
  R24's "from invalid to unsupported" applies to the primitive set only, and
  row R14 is the binding wording for this condition. *Reason:* the accepted
  page gives the choice to the contribution
  (`docs/contribution-contract.md:203`, quoted above), and a refusal in the
  target would put a policy there that is not the target's. *Not taken:* the
  Go target refuses the case as unsupported.
- *(b) An unknown form name* stays unsupported, as today (`model.go:89-90`),
  and the answer is the one checker's. *Reason:* a later revision of a
  vocabulary can declare the form, which is the test of unsupported
  coverage: "a refusal is unsupported coverage when a later slice can admit
  the same input" (`docs/source.md:509-510`). And a category does not change
  without a case (AC-05). *Not taken:* invalid, as an input "malformed under
  the published vocabulary" (`docs/source.md:512-513`).

**M12: where the intent carrier is cut.** The model side holds the
form-table type and no type that carries a value. `Form[T]`, `FormIssuer[T]`
and `Expression[T]` stay whole in the generation contract package, with the
planner.

- *Reason.* `Form[T]` is created only by its issuer, which refuses an empty
  name (`contracts/generation/go/expression.go:18-23`), so sentence 1 of the
  test puts it on the operations side. `Expression[T]` has unexported fields,
  and one of its six kinds is issued by the planner alone
  (`expression.go:52-63`): "only the planner may issue a plan"
  (`docs/package-boundaries.md:145-147`). A native node's arguments can be
  any of the six kinds, so a tree of plain native nodes cannot hold a
  mapping's output. What a native node is made of is stated by the table
  type. The agnostic constructor is the function that creates an
  `Expression[T]`, and Go keeps it in that type's package.
- *What this narrows.* The advice says: "The generic expression machinery is
  a shared carrier." (advice, section 2.4), and its layout puts shared
  carriers under the model root. M4 took from it that only the native kind
  is intent. M12 narrows M4: the advice's generic intent expression carrier
  is **not on the model side**. It is the admitted state of the planner's
  operation. Row R15, part 3, of the table "Applied" is read so: "only the
  native kind of an expression (form, data, arguments) is intent" says what
  intent is, and places no type.
- *Not taken:* a plain record of one native node on the model side, generic
  in the type of its arguments, which the sealed expression would hold for
  its native kind. It stays the recorded alternative. It can be added
  without moving anything, and its trigger is the first package outside the
  generation contract that must hold native intent without a planner
  ([section 11](#11-decisions-and-deferrals)).
- *Not taken:* `Expression[T]` on the model side with exported fields. It
  would cost the guarantee that only the planner issues a declaration
  reference.

**M13: the form table is a sketch that slice S2 makes final.** "A table
type is settled by building the one checker against five tables, not on a
page" (M13). So the type names, the field names and the profile string
below are a sketch and are not decided. What the
table must be able to state is derived from the checks on `main`: the number
of data strings with a rule for each (not empty, one of a set), the number of
arguments, the form of every argument, a datum that differs between the
arguments, and for a declaration form an optional form of its body
(verification report 1, section 2.4). The record rule is a relation between
siblings. It is declared, because all five copies state it identically, which
is the concrete need that row R18 asks for before the table type grows.

```go
// Package intent is the agnostic model of intended forms. Plain data only.
package intent

// Vocabulary is one language's declared table of intended forms.
type Vocabulary struct {
	Profile      string        // identity and revision of this vocabulary
	Forms        []Form        // the forms of an expression
	Declarations []Declaration // the forms of a declaration
}

// Form declares one expression form.
type Form struct {
	Name      string
	Data      []Datum // one entry for each data string, in order
	MoreData  *Datum  // if set: any number of further data strings under this rule
	Arguments Arguments
}

// Datum declares one data string. Name says what the string states and is
// never compared.
type Datum struct {
	Name     string
	Nonempty bool
	OneOf    []string // if not empty: the string is one of these
}

// Arguments declares the arguments of a form. Max -1 means no upper bound.
type Arguments struct {
	Min, Max int
	Form     string // if not empty: every argument is a native node of this form
	Unique   string // if not empty: the datum of this name differs between the arguments
}

// Declaration declares one declaration form.
type Declaration struct {
	Name string
	Body string // if not empty: the body is a native node of this form
}
```

Go's table for all seven forms, beside the condition that states each shape
today in `transform/syntax/golang/model/go/model.go`:

```go
var Forms = intent.Vocabulary{
	Profile: "go/intended-forms-1",
	Forms: []intent.Form{
		{Name: FormPrimitive, Data: []intent.Datum{{Name: "name", Nonempty: true}}},
		{Name: FormSlice, Arguments: intent.Arguments{Min: 1, Max: 1}},
		{Name: FormField,
			Data: []intent.Datum{
				{Name: "name", Nonempty: true},
				{Name: "optional", OneOf: []string{"true", "false"}}},
			Arguments: intent.Arguments{Min: 1, Max: 1}},
		{Name: FormRecord, Arguments: intent.Arguments{Max: -1, Form: FormField, Unique: "name"}},
		{Name: FormExternal,
			Data:      []intent.Datum{{Name: "module"}, {Name: "name"}},
			Arguments: intent.Arguments{Max: -1}},
	},
	Declarations: []intent.Declaration{
		{Name: FormAlias},
		{Name: FormDefined},
	},
}
```

| Form | The condition today | What the table states | What it does not state |
| --- | --- | --- | --- |
| `primitive` | One datum, no argument, and the datum is `string` or `bool` (`:63`) | One datum that is not empty, no argument. "Not empty" is the sketch's own rule: it keeps an empty name invalid, as the condition of today does. Slice S2 confirms it or drops it (M13) | Which names the target lowers. That is coverage (the change of category above). That the name is a predeclared type of Go is a constraint that exists only as its evaluator, the universe scope of the building toolchain (`transform/syntax/golang/symbol/go/symbol.go:34`) |
| `slice` | No datum, one argument (`:67`) | The same | |
| `field` | Two data, the first not empty, the second `true` or `false`, one argument (`:71`) | The same | That the name is a Go identifier. The Go checker refuses another name when the output is sealed (`check/syntax/golang/provider/go/check.go:27`). What `optional` means for Go: the target reads the name only (`transform/syntax/golang/target/go/provider.go:122-123`) |
| `record` | No datum, every argument is a `field`, no name twice (`:75-86`) | The same | |
| `external` | Refused as unsupported before its shape is read (`:87-88`). Its constructor writes two data and any arguments (`:49-51`) | Two data and any arguments, as the constructor states them. No rule for either datum | Whether Go builds it. That is coverage |
| `alias` | A declaration form (`:102`) | A declaration form without a rule for its body | |
| `defined` | A declaration form whose body is a `record` (`:105-110`) | A declaration form without a rule for its body (row R14, M17) | The declaration policy's choice between `alias` and `defined` (`transform/syntax/horos/mapping/declarations/go/declarations.go:52-55`) |

Three things hold for every form and belong to the one checker and not to a
table: a form of another vocabulary is refused as `ErrScope` (`:56-59`), a
declaration has a name that is not empty (`:102`), and a parameter or a
reference is the planner's and no form of a language (row R17). Derived here:
for TypeScript's `interface` and for the `struct` of Rust and C++ the body
rule is the language's own, so their tables state `Body`, and Go's does not.
Slice S2 confirms it for each.

**The declared model is an upper bound here too** (sentence 2 of the test).
The table does not state where a form may occur: a `slice` of a `field`
passes the checks of `model.go` today and fails only when the output is
sealed. It does not state what a form becomes as nodes. The page of each
language says so with the relation of [section 4](#4-one-languages-models).

**Where the typed constructors live.** A mapping may not import a target:
"Mapping packages depend on native models, never target AST constructors."
(`AGENTS.md:190`). So the typed constructors of a language, `Primitive`,
`Collection`, `Field`, `Record` and `External`, with the language's issuer,
are a package of their own on the operations side that mappings and the
target import. [Section 7](#7-the-layout) names its path (M14) and
[section 8](#8-roles-per-root) its role (M15).

## 6. Declared models

The owner asked: "what is currently not declarative about our models?" M5
(rows R16 to R20) answers what is declared as data, with what, in which order
and under which identity.

**What is data and what stays code, by sort** (rows R17 and R18).

| Sort | Declared as data | Stays code, and whose |
| --- | --- | --- |
| Syntax | Node kinds, roles, attribute names (constants today). For each node kind its attributes, required or not, its roles with their cardinalities and allowed kinds, their order, and whether it refers (the table type, `internal/grammar/rules.go:15-33`). Word lists and character tables | The rules between nodes, the checker's: two for Go (`check.go:51-81`). The evaluation of a predicate. Parsers and formatters as a whole: concrete spelling, precedence and layout |
| Analysed semantics | Keys and profiles (functions and constants today). With slice S9 the shape of each fact: its operator, its encoding, the fields of its payload, its operands | The compiler bridge, which extracts. The typed readers, which interpret |
| Intent | Form names (constants today). With slice S2 the form table of [section 5](#5-intent) | Construction, the one checker, lowering and coverage |
| Symbols | The categories and the plain record of a symbol | The admission of a name, which for Go asks the universe scope of the toolchain. Issuing |

"Complex relationships between nodes" stay code "unless a concrete need
justifies extending the declarative language" (advice, section 3.3; row
R18). So the declared grammar is an upper bound of what a checker admits, as
sentence 2 of the test says, and each language's page states it.

**Which mechanism declares what (row R16).** No new schema language is
built, and the node patterns of horos kinds are not used for it. Two
mechanisms exist.

| What is declared | Mechanism | Why this one |
| --- | --- | --- |
| A grammar table, a form table, a word list | A Go table: an exported value of a plain table type in a model package, which shared code interprets | The grammar table type exists and has one interpreter. "Generic carriers, contracts, engine, printer and native syntax stages have no external Horos/Ontos imports." (`CHARTER.md:129-131`), so a native language's model is not declared in horos |
| A value-shaped payload of a contribution | A horos type, from which bittype generates the Go and TypeScript types | Accepted: "Payload schemas are horos types. A payload is an ontos value. Its family declares the payload type in horos; bittype generates the Go and TypeScript types, and horos conformance validates a payload. No separate schema language is introduced." (`docs/representation-families.md:98-101`) |

Which of the two declares the payload of an analysed fact is not decided. It
is decided with slice S9, which adds the fact schemas
([section 11](#11-decisions-and-deferrals)).

**The order (M5, row R17).** It follows
[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
and not the advice's "syntax, then facts, then intent". What is misplaced
moves first: the form table, where per-language copies of one mechanism
exist, and the grammar tables, which an operation owns today (row R4).
What is only missing is added afterwards: the schemas of analysed facts.
Parameter and reference structures belong to the shared carrier and to no
language's schema.

**Identity (M5, row R19).** A model's identity and revision is its
**profile**: a string that its owner declares. A content digest is
provenance and never the identity. Three accepted sentences say so:
"Profiles are caller-owned compatibility promises, not hashes or
automatically negotiated versions." (`trees/representation/go/language.go:13-14`);
a profile "is not a generated fingerprint, toolchain detector or published
canonical language revision" (`docs/engine.md:47-48`); "A digest of an
implementation is provenance. It replaces none of the five."
(`docs/contribution-contract.md:1034`). The identity depends on no package
path, so no move on this page changes one.

**When a profile changes** is decided already, on #278 and #276, and
[native semantics](native-semantics.md#rust-vocabulary) records it:

- A removed or changed member needs a new version: "A reader written for the
  earlier vocabulary cannot read a graph without that member, so the promise
  of the earlier version does not hold and the version changed."
  (`docs/native-semantics.md:841-843`).
- An addition keeps the version of the model's profile: "It adds an entity
  kind and a binding and removes nothing, so by the line drawn on #278 the
  Rust profile keeps its version" (`:853-854`).
- The profile of the writer changes with every added fact: "The helper's
  profile name changes when facts are added." (`:877`).

**What is missing and added** (row R19):

- A declared profile for each native syntax vocabulary (gap 1 of
  [section 3](#3-consistency-and-correctness)). Slice S9.
- An identity and a revision for each form vocabulary. Today a form's
  identity is "an allocation-unique owner"
  (`contracts/generation/go/expression.go:103`), a pointer of one process.
  Slice S2, as the field `Profile` of the sketch above.

**What keeps a declaration honest (row R20).** Three checks exist: the
generic checker interprets the table, and the ten parser and formatter
providers call their own checker
(`git grep -l 'check/syntax/[a-z]*/provider/go"' -- 'parse/**/*.go' 'format/**/*.go'`
gives 10 files). Added with the fact schemas in slice S9: a check of fact
shapes, tests of the readers against the schemas, and generated names for
the two bridges that are not written in Go. AC-05 asks that the mechanism be
general, "without requiring a new case-specific check for each
implementation change", and one generator from one description is that.

**M18: how a grammar table names a predicate.** An entry of a grammar table
holds a function value today (`Attribute.Accept`,
`internal/grammar/rules.go:17`). Go's table binds 17 of them, so it is not
data, and row R7c makes the naming of predicates by identity the
precondition of its move.

An entry names its predicate by an owner-qualified name, as a semantic key
names a fact (`semantics.Key`, `semantics/go/model.go:24`). The names that
every language can use (`bool`, `nonempty`, `string`, `decimal`) are
declared in `model/syntax/go` (M11). A language's own, such as Go's
`identifier`, are declared in its lexical part,
`model/syntax/<language>/lexical/go` (M19). A set of allowed
strings is data in the entry and no predicate. The checker that calls the
interpreter supplies the evaluators by name, and a name without an evaluator
fails the check and is never skipped. Where the evaluator is the building
toolchain, as for Go's identifier (`check.go:13`), the name is all that the
model holds, and the language's page says so (sentence 2 of the test).

- *Reason:* a closed enumeration would be a shared edit for each new
  predicate of one language ([AC-03](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-03)).
- *Not taken:* a closed enumeration of predicates in the package of the
  table type.
- *Not taken:* the tables stay with their checkers and keep their function
  values. Row R4 refuses it: an operation does not own the declaration it
  interprets.

## 7. The layout

**Decided (M6, rows R21 and R22).**

| Point | Decision | Source |
| --- | --- | --- |
| The roots | Two: `model` and `ops`. No third root | The direction; row R8 |
| Below `model` | The sort of model first, then the language | M6. The owner's sketch is `model/{syntax,semantics,...}/go/...`. The standing rule is "Do not create a language-first packages tree" (`LAYOUT.md:134`, `AGENTS.md:236-237`). The advice's tree puts the language first and its own section 4.3 says the opposite: "organized around the kinds of models being declared, with language identity as a further dimension" |
| Below `ops` | The stage matrix as it is: `<stage>/<dimension>/<scope>/<component>/go` below the root | M6; row R22. The owner's sketch is `ops/{transform,...}/go/...` |
| Modules | One. "One public root Go module contains independently owned packages" (`LAYOUT.md:6`) | M6. Its cost: the module graph cannot carry rule 1. The import test and Go's rule for `internal` carry it ([section 9](#9-enforcement)) |
| `internal/` | It goes under `ops`. It stays at the module root until no future model package imports a helper | M6; row R21 |
| Integration and tools | Under `ops`: `ops/integration`, `ops/tools`, and `ops/verification` | Row R21. With the helpers at `ops/internal`, Go lets only code below `ops/` import them. `integration/go` imports `internal/horosfiles` and `internal/tree`, `tools/horos-module/go` imports `internal/horosmodel` and `internal/lex`, and a test of the grammar extractor imports `internal/grammar` |
| The trailing `go` | It stays: "The final `go` directory denotes implementation language, not target" (`LAYOUT.md:94`) | Row R21 |

**`internal/` and Go's own rule.** Go lets a package below a directory named
`internal` be imported only from the tree rooted at that directory's parent.
With the helpers at `ops/internal`, the compiler itself refuses a model
package that imports one. Five packages that hold, or were read as holding,
a model-side part import a helper today: `trees/representation/go` imports
`internal/tree`, `semantics/rust/go` and `transform/syntax/rust/symbol/go`
import `internal/rust`, and `semantics/typescript/go` and
`transform/syntax/typescript/symbol/go` import `internal/typescript`. In each
the importing code is on the operations side by the table below: packing and
unpacking, the validation of a compiler run's inputs
(`semantics/rust/go/environment.go:156`, `:160`;
`semantics/typescript/go/environment.go:37`) and the admission of a symbol
(`transform/syntax/rust/symbol/go/symbol.go:56`, `:58`, `:73`;
`transform/syntax/typescript/symbol/go/symbol.go:24`). So every caller of a
lexical evaluator ends on the operations side, and no helper is needed on
both sides (row R7h). The model root gets no `internal` directory of its
own: a package below `model/internal` could not be imported from `ops/`.

**What is no package.** Derived here: `docs/`, `research-docs/`, `testdata/`,
`go.mod` and the files at the repository root describe or serve both roots
and stay where they are.

### Names below the roots

**M10: the shape of a path below `model`.** Two shapes, as `semantics/go`
and `semantics/<language>/go` have them today. `model/<sort>/go` is the
agnostic part of a sort: the types in which every language's declaration of
that sort is written. `model/<sort>/<scope>/go` is one scope's declaration
of that sort, where a scope is one language or an ordered pair. Where one
scope declares several independent things of one sort, each is
`model/<sort>/<scope>/<name>/go`, the path form that #222 decided for
several contributions (`LAYOUT.md:100-104`). The two tree carriers keep
their own directory, `model/trees/ast/go` and `model/trees/cst/go`.

- *Reason:* it is M6 with the shapes the vocabularies have today. A scope
  named `shared` is the stage matrix's device for a component level, which a
  sort does not have.
- *Not taken:* one shape with an explicit scope,
  `model/<sort>/<scope>/<part>/go`, where the scope is a language, a pair or
  `shared` (`LAYOUT.md:98`).

**M11: ten sorts.** A sort is a kind of model, not a kind of table.

| Sort | Directory | Holds | The word comes from |
| --- | --- | --- | --- |
| Language identity | `model/language` | `ID` and `Provider` | `contracts/language/go` today |
| Trees | `model/trees` | The generic abstract and concrete tree | `trees/` today |
| Syntax | `model/syntax` | In `model/syntax/go`: the grammar-table type and the predicate names that every language can use. In `model/syntax/<language>/go`: the language's node kinds, roles and attributes, its grammar table and its syntax profile. In `model/syntax/<language>/lexical/go`: its word lists, character tables and literal patterns, and the names of its own predicates (M19) | `syntax/` today |
| Semantics | `model/semantics` | The plain records of the graph, and a language's keys, profiles and payload shapes | `semantics/` today |
| Files | `model/files` | The plain records of a file set | `files/` today |
| Source | `model/source` | The plain records of a Source: unit, group, relation | `source/` today |
| Intent | `model/intent` | The form-table type, and each language's form table | "native intent" (`docs/generation.md:77`) |
| Symbols | `model/symbols` | A language's plain symbol record and its categories | `symbol` today |
| Families | `model/families` | The profile identities and the family vocabularies of a source language | "family" (`docs/representation-families.md`). On the migration's branch |
| Contributions | `model/contributions` | The identity types of a contribution, and each pair's contribution vocabularies | "contribution" (`docs/contribution-contract.md`). On the migration's branch |

- *Reason:* the derivation rule: map onto the concepts that exist before
  adding one. A grammar and a lexical rule are syntax. The owner's sketch
  names the sorts `syntax` and `semantics`.
- *Not taken:* `grammar` and `lexical` as sorts of their own, with
  `model/grammar/<language>/go` and `model/lexical/<language>/go`.
- *Not taken:* `forms`, or keeping `model`, for `intent`. The singular
  `symbol`.
- *Adjusted by M19.* M11 first put a language's kinds, its grammar table and
  its lexical tables into the one package `model/syntax/<language>/go`, and
  allowed a part its own package only where an import rule forbids the one.
  The check of that against the code is the reason for M19.

**M19: a language's lexical tables are a part below its syntax model.**

- The lexical tables of a language (reserved words, character ranges,
  literal patterns) and the names of the predicates over them are
  `model/syntax/<language>/lexical/go`, by the shape of M10 for a scope that
  declares several things. It imports no node kinds.
- Node kinds and the grammar table stay one package,
  `model/syntax/<language>/go`. It may import its language's lexical part to
  name a predicate. The lexical part never imports it. "So structure knows
  spelling and spelling does not know structure, and the import test can
  carry that." (M19.)
- The lexical helper on the operations side imports the lexical part only,
  beside what it imports today.
- The sorts stay ten. `lexical` is a part of a language's syntax and no sort
  of model. The table type and the shared predicate names stay in
  `model/syntax/go` (M11, M18).
- *Reason:* two accepted sentences that the one package works against. The
  role of a lexical helper is "A lexical helper imports `internal/grammar`
  only, and never a contract or a stage provider." (`LAYOUT.md:167`), and
  the helper "knows spelling, not structure"
  (`integration/go/architecture_roles_test.go:245-246`). With one package the
  helper would import the node kinds with the word lists, and nothing could
  keep it from naming structure. And "An importer that uses a slice of the
  package" is a sign of wrong grouping (`docs/package-boundaries.md:157-159`):
  the checker uses kinds and table, the helper only the word lists, a
  mapping only the kinds. Word lists and node kinds do not depend on each
  other. The table depends on both: it is keyed by kinds and names
  predicates.
- *Not taken:* the one package of M11 as first decided.

**M11 and M19 checked against the code.** Read at `0581401`.

| Language | The grammar table today | The predicates of its own that the table names, and where their evaluator is today | Word lists, character tables and patterns today | Does any of it need a node kind, a role or an attribute? |
| --- | --- | --- | --- | --- |
| C++ | `check/syntax/cpp/provider/go/check.go:29` | `Identifier` and `Builtin` of `internal/cpp`. One anonymous function: not empty and without a line break (`check.go:31`) | The reserved words (`internal/cpp/names.go:11`), the words of a builtin type (`internal/cpp/types.go:7`) | No |
| Go | `check/syntax/golang/provider/go/check.go:17` | `identifier`, ten times (`check.go:13`). One anonymous function: `.` or an identifier (`:19`). Both evaluators are in the checker | None. Go's words are the building toolchain's (`go/token`) | No |
| horos | `internal/declsyntax/check.go:67` | `identifier`, `moduleName`, `hexadecimal` and `pin`, all in `internal/declsyntax` (`check.go:17-60`) | The reserved words (`internal/horos/spelling.go:23-29`), the bound on the digits of an integer (`:90`) | No |
| Rust | `check/syntax/rust/provider/go/check.go:33` | `Identifier`, `Declaration`, `LifetimeParameter`, `LifetimeReference`, `PathSegment` and `ABI` of `internal/rust` | The reserved words (`internal/rust/names.go:11`), the character ranges of a crate name (`internal/rust/cratename.go:48`) | No |
| TypeScript | `check/syntax/typescript/provider/go/check.go:18` | `Identifier` and `TypeName` of `internal/typescript` | The reserved words (`internal/typescript/names.go:11`), the predefined type words (`names.go:24`), the two patterns of a number (`internal/typescript/literal.go:13-14`), as strings | No |

The tables also name the shared predicates (`bool`, `decimal`, `nonempty`,
`positive`, `signed`, `string`), which are declared in `model/syntax/go`,
and sets of allowed strings, which are data in the entry (M18).

What follows:

- **Nothing that the lexical part would hold needs a node kind.** Every
  predicate has the type `func(string) bool`
  (`internal/grammar/rules.go:17`) and sees one string. The four lexical
  helpers import `internal/grammar` and nothing else of this module, and
  `internal/horos` imports nothing. So the direction of M19 holds in the
  code: a table needs predicate names, and no word list and no predicate
  needs the table or the kinds.
- **The imports.** `model/syntax/<language>/lexical/go` imports at most
  `model/syntax/go`, for the type of a predicate name. It imports no
  language's kinds. `model/syntax/<language>/go` imports
  `model/language/go`, `model/trees/ast/go` and `model/trees/cst/go`, as
  `syntax/<language>/go` does today, and two more: `model/syntax/go` for
  the table type, and its own lexical part for predicate names.
  `model/syntax/go` imports `model/trees/ast/go`, because a table is keyed
  by `ast.Kind[L]` (`internal/grammar/rules.go:33`). Every one is a model
  package, and there is no cycle.
- **Nothing that needs only the kinds gets an import it may not have.** The
  import test reads the direct imports of a package. A parser, a formatter,
  a mapping or a target that imports the syntax package for its kinds makes
  no new import. No foundation or kernel package imports a syntax
  vocabulary at all. The prohibitions of the stage roles are "A checker
  gets no concrete tree, a formatter no lexer, and no Syntax stage a
  semantic contract." (`LAYOUT.md:159`). A grammar table is none of the
  three.
- **Go's lexical part holds predicate names and nothing else.** Derived
  here. Go has no word list of its own: its keywords are those of the
  toolchain that builds bittype (`token.Lookup`, `check.go:13`). M19 places
  the names of a language's predicates in its lexical part, and a predicate
  name is all that the model states of a rule that exists only as its
  evaluator (sentence 2 of the test). So `model/syntax/golang/lexical/go`
  declares the names of Go's two predicates and no table. The same shape
  for every language keeps the import rule one rule. It is no placeholder
  directory: it holds declarations that the table imports. Which release of
  the toolchain the rule follows is recorded nowhere, and that point stays
  open ([section 11](#11-decisions-and-deferrals)).
- **One horos rule has two evaluators today.** `identifier` in
  `internal/declsyntax/check.go:17-27` and `Identifier` in
  `internal/horos/spelling.go:11-21` state the same rule, and neither
  package imports the other. With one name in the lexical part, slice S4
  meets it. It is recorded here and is no part of this page's decisions.
- **How a table is written.** The entry constructors `Attr`, `One`, `Opt`,
  `Many` and `Attributes` (`internal/grammar/rules.go:35-40`, `:72-78`)
  refuse nothing and are construction, so they stay on the operations side
  (row R7a). A table in a model package is a composite literal of the table
  type.
- **One reader by path.** The grammar extractor of the documentation tool
  names the owner of the table type as the string
  `github.com/Bitspark/bittype/internal/grammar`
  (`tools/provider-matrix/extract/grammar/go/extract.go:16`). Slice S3
  changes it with the type's package.

**M14: the paths of the operations halves.** The general rule: a package or
a half that this page does not name otherwise keeps its path below its
root. That is all that happens to a whole operations package (row R22).
Four halves have a name of their own.

| Half | Path | Why there | Not taken |
| --- | --- | --- | --- |
| The contract identities, evidence requirements and settings of a language's analysis | `ops/semantics/<language>/go` | The general rule: it is what remains at today's path once the names have their model package | `ops/analyze/files/<language>/contracts/go`, a path under one stage for identities of two |
| The typed readers of a language's analysed graph | `ops/analyze/files/<language>/reader/go` | They open "an admitted analysis receipt" (`semantics/golang/go/model.go:41`), so they belong to the analyze stage of the Files dimension, beside the provider and the query service | Joining the query service, which a package that is no composition imports only under an exception |
| The typed constructors of a language's intent, with its issuer | `ops/transform/syntax/<language>/intent/go` | It is today's `transform/syntax/<language>/model/go` without the table and the checks. "The component `model` below `ops` would say the opposite of the root" (M14) | Keeping `model` as the component. Putting the constructors into the target, which a mapping may not import (`AGENTS.md:190`) |
| The admission and the issuing of a language's symbols | `ops/transform/syntax/<language>/symbol/go` | The general rule | |

- *Reason:* the general rule makes the two halves of a split package
  findable from each other and keeps the scripted move mechanical.

### Every place named

The advice's tree leaves six places without a name (row R21). The
migration's branch adds four more that this page has to place (row R23).

| What | Path | Side | Standing | Slice |
| --- | --- | --- | --- | --- |
| The Source records `Unit`, `GroupID`, `Group` and `Relation` (`source/go/source.go:20-37`) | `model/source/go` | Model | Row R7d places them. The name is M11 | S6 |
| Language identity, `ID` and `Provider` (`contracts/language/go/language.go:11-25`) | `model/language/go` | Model | Row R3 places it. The name is M11 | S6 |
| The intent carrier | `ops/contracts/generation/go`, where it is | Operations | M12 | None |
| The form-table type | `model/intent/go` | Model | M4 places it. The name is M11 | S2 |
| The grammar-table type (`internal/grammar/rules.go:15-33`) and the shared predicate names | `model/syntax/go` | Model | Row R7c places it. The place is M11 and M18 | S3 |
| A language's grammar table | `model/syntax/<language>/go`, with the node kinds | Model | Row R7c places it. The place is M11 and M19 | S3 |
| Lexical tables: a language's word lists, character tables and literal patterns, and the names of its predicates | `model/syntax/<language>/lexical/go` | Model | Row R7h places them. The place is M19 | S4 |
| The operations halves of split vocabulary packages | The four paths of M14 | Operations | M14 | S2, S5 |
| The identity types of a contribution (`ContributionIdentity`, `Meaning`, `ContributionReference`, `SchemaReference`, `EvidenceKind`; on the migration's branch `contracts/generation/go/contribution.go:22-46`, `:138`) | `model/contributions/go` | Model | Rows R7e and R23 place them. The name is M11 | S7 |
| A pair's contribution vocabulary (on the branch `transform/syntax/<source>-to-<target>/vocabulary/<name>/go`, 24 packages) | `model/contributions/<source>-to-<target>/<name>/go` | Model | Row R7g places it. The name is M11 | S7 |
| A family of one language (on the branch `transform/syntax/horos/family/<name>/go`, 6 packages), and the two profile identities of `transform/syntax/horos/mapping/go` | `model/families/<language>/<name>/go`, and `model/families/<language>/go` for the profile identities | Model | Row R23 asks for the place. The name is M11 | S7 |
| The shared part of the payload codec role (on the branch `transform/syntax/horos/payload/go`), and each codec | Operations. The path below `ops` is the migration's to name: "that part does not live under a kernel path (decision O41 on #196); the migration names its place" (`LAYOUT.md:476`) | Operations | Row R7g; decision O41 | S7, S8 |

The counts on the branch come from
`git ls-tree -r --name-only origin/feat/contribution-contract-196 -- transform | grep "/vocabulary/.*/go/vocabulary.go" | wc -l`
(24) and from the package list of the branch (six `family` packages), read
at `af1dcde`.

### Every package of `main`

The command for the production packages, at the repository root:

```sh
git ls-tree -r --name-only origin/main | grep '\.go$' | grep -v '_test\.go$' | grep -Ev '(^|/)testdata/' | xargs -n1 dirname | sort -u
```

It gives 129 at `0581401`. "Whole" means that the package takes the prefix
of its root and nothing else changes. The paths in the column "After" use
the names of M10, M11, M14 and M19. The side of each part is decided by the row
named, and where no row is named it is derived here from the test.

**Carriers (6).**

| Today | After | Source |
| --- | --- | --- |
| `trees/ast/go` | `model/trees/ast/go`, whole | Row R7d |
| `trees/cst/go` | `model/trees/cst/go`, whole | Row R7d |
| `trees/representation/go` | `ops/trees/representation/go`, whole. No part of it is plain data that has a meaning outside its operation | Rows R7d and R8; [section 2](#2-the-test) |
| `files/go` | **Splits.** The records `ID`, `Path`, `File`, `Files`, `SyntaxLink`, `Context` to `model/files/go`. The limits, `ValidPath`, `Capture`, `Cover`, `CheckFiles`, the errors, and the admitted type that slice S6 adds, to `ops/files/go` | Row R7d; gap 4 |
| `semantics/go` | **Splits.** The records and `Data` (`model.go:20-125`) to `model/semantics/go`. The limits, the sealed `Semantics[L]`, `New`, the accessors and the queries to `ops/semantics/go` | Row R7d |
| `source/go` | **Splits.** `Unit`, `GroupID`, `Group`, `Relation` to `model/source/go`. `Publication`, `Source`, `Policy`, `Publisher`, `Envelope` and the rest to `ops/source/go` | Row R7d. The accepted page declares the Source as "private sealed storage" with two receipts (`docs/source.md:373-380`) |

**Vocabularies (16).**

| Today | After | Source |
| --- | --- | --- |
| `syntax/{cpp,golang,horos,rust,typescript}/go` | `model/syntax/<language>/go`, whole, once language identity is at `model/language/go`. Each receives its language's grammar table (slice S3) | Rows R3 and R7d |
| `semantics/horos/go` | **Splits.** The profiles, the keys, the descriptor prefix, the term encoding, the limits setting key and the prelude pin to `model/semantics/horos/go`. `ResolutionContract`, `AnalysisContract`, `AnalysisContractWithResolutionConfiguration` (`vocabulary.go:75-91`) and the reader `DescriptorConstructor` (`:55`) to `ops/semantics/horos/go` | Row R2 |
| `semantics/{golang,rust,typescript}/go` | **Splits in three.** The profiles, the key functions, the payload shapes and the plain records of a compiler run's inputs to `model/semantics/<language>/go`. The contracts, the evidence requirement, the requirements, the configuration and the encoding of inputs into settings to `ops/semantics/<language>/go`. `Model`, `Open` and the typed readers to `ops/analyze/files/<language>/reader/go` | Rows R2 and R7b |
| `transform/syntax/{cpp,golang,rust,typescript}/model/go` | **Splits.** The form names, the primitive names and the form table to `model/intent/<language>/go`. The issuer and the typed constructors to `ops/transform/syntax/<language>/intent/go`. The checks become the one checker in `ops/transform/syntax/shared/target/go`. The coverage answers go to `ops/transform/syntax/<language>/target/go`. Go's record-body condition goes to the declaration contribution | M4; [section 5](#5-intent) |
| `transform/syntax/{golang,rust,typescript}/symbol/go` | **Splits.** The plain record of a symbol, its categories and the names it declares to `model/symbols/<language>/go`. The refusing constructors, the issuer, `ExternalUse` and `ReadExternal` to `ops/transform/syntax/<language>/symbol/go` | Row R7a; M8, slice S5 |

**Contracts (13).**

| Today | After | Source |
| --- | --- | --- |
| `contracts/language/go` | `model/language/go`, whole | Row R3 |
| `contracts/{analysis,check,format,parse,pipeline,print,resolution,semantic,text,transform,verification}/go` | `ops/contracts/<name>/go`, whole | Rows R7e and R8 |
| `contracts/generation/go` | `ops/contracts/generation/go`. On `main` it is whole (M12). On the migration's branch the identity types of a contribution leave it for `model/contributions/go` | Rows R7e and R23; M12 |

**Kernel (7).** All whole, to `ops/<same path>` (row R7f): `engine/go`,
`transform/syntax/shared/generation/go`, `transform/syntax/shared/target/go`,
`check/syntax/shared/ast/go`, `check/syntax/shared/cst/go`,
`print/syntax/shared/provider/go`, `print/files/shared/provider/go`. The
kernel package `transform/syntax/shared/target/go` receives the one form
checker in slice S2.

**Stage providers and query services (28).**

| Today | After | Source |
| --- | --- | --- |
| `check/syntax/{cpp,golang,rust,typescript}/provider/go` | **Splits.** The grammar table to `model/syntax/<language>/go`. The call of the interpreter, the rules in code and the provider stay, at `ops/<same path>` | Row R7c; M11 |
| `check/syntax/horos/provider/go` | `ops/<same path>`, whole. Its table is in `internal/declsyntax` | Row R22 |
| `parse/syntax/{cpp,golang,horos,rust,typescript}/provider/go`, `parse/files/horos/provider/go` | `ops/<same path>`, whole | Row R22 |
| `format/syntax/{cpp,golang,horos,rust,typescript}/provider/go`, `format/files/{horos,typescript}/provider/go` | `ops/<same path>`, whole | Row R22 |
| `resolve/files/{golang,horos,rust,typescript}/provider/go` | `ops/<same path>`, whole | Row R22 |
| `analyze/files/{golang,horos,rust,typescript}/provider/go`, `analyze/files/{golang,typescript}/queries/go` | `ops/<same path>`, whole | Row R22 |

**Generation by language and by pair (30).**

| Today | After | Source |
| --- | --- | --- |
| `transform/syntax/{cpp,golang,rust,typescript}/target/go` | `ops/<same path>`, whole. Each receives its coverage answers and the call of the one checker | M4 |
| `transform/syntax/horos/source/go` | `ops/<same path>`, whole | Row R22 |
| `transform/syntax/horos/mapping/{declarations,naming,pattern}/go` | `ops/<same path>`, whole | Row R22 |
| `transform/syntax/horos/mapping/go` | **Splits.** The two profile identities (`mapping.go:28`, `:34`) to `model/families/horos/go`. The function type `Library[T]` stays on the operations side, where [bittype/ACX-010](architecture-cases/0010-agnostic-mapping-adapter.md) moves it to the agnostic layer | Row R7g |
| `transform/syntax/horos-to-{cpp,rust,typescript}/mapping/{boolean,declarations,list,record,text}/go` (15) | `ops/<same path>`, whole | Row R22 |
| `transform/syntax/horos-to-golang/mapping/{boolean,declarations,list,optional,record,text}/go` (6) | `ops/<same path>`, whole | Row R22 |

**Private helpers (16).**

| Today | After | Source |
| --- | --- | --- |
| `internal/{concrete,golangsem,horosfiles,horosmodel,lex,parseutil,rustsem,tree,typescriptimports,typescriptsem}` | `ops/internal/<name>`, whole | Row R21 |
| `internal/grammar` | **Splits.** The table types `Attribute`, `Role`, `Rule`, `Rules` (`rules.go:15-33`) to `model/syntax/go`, with the names of the shared predicates. The interpreter `Check` (`:80`), the entry constructors and the predicates stay, at `ops/internal/grammar` | Row R7c; M11, M18 |
| `internal/{cpp,horos,rust,typescript}` | **Splits.** The word lists, the character tables and the literal patterns to `model/syntax/<language>/lexical/go`, with the names of the language's predicates. The evaluators and the literal decoders stay, at `ops/internal/<language>` | Row R7h; M19 |
| `internal/declsyntax` | **Splits.** The horos grammar table (`check.go:67`) to `model/syntax/horos/go`. `Check` and the rules in code stay, at `ops/internal/declsyntax` | Row R7c; M11 |

**Verification (7).** All whole, to `ops/<same path>` (row R22):
`verification/files/horos-to-go/go`,
`verification/files/horos-to-go/mappings/{list,optional,record,scalar}/go`,
`verification/files/source/horos/go`, `verification/logos/go`.

**Tools (6).** All whole, to `ops/<same path>` (row R21):
`tools/horos-module/go`,
`tools/provider-matrix/{go,capability/go,extract/grammar/go,observation/go,report/go}`.

**Packages that hold only tests (3),** outside the 129: `integration/go`,
`verification/integration/go` and
`tools/provider-matrix/extract/structfixture/go` go to `ops/<same path>`
(row R21).

**The sum.** By this table **8 packages of `main` are wholly model, 25 split
and 96 are wholly operations.**

| Group | Wholly model | Split | Wholly operations |
| --- | --- | --- | --- |
| Carriers (6) | 2 | 3 | 1 |
| Vocabularies (16) | 5 | 11 | 0 |
| Contracts (13) | 1 | 0 | 12 |
| Kernel (7) | 0 | 0 | 7 |
| Stage providers and query services (28) | 0 | 4 | 24 |
| Generation by language and by pair (30) | 0 | 1 | 29 |
| Private helpers (16) | 0 | 6 | 10 |
| Verification (7) | 0 | 0 | 7 |
| Tools (6) | 0 | 0 | 6 |
| **129** | **8** | **25** | **96** |

M11 and M19 change the home of the grammar tables and the word lists and
the side of no package, so these numbers are the same before and after
them. M7 first
gave 11, 27 and 91, from the placement table of verification report 2. The
[decisions of 13 November 2026](https://github.com/Bitspark/bittype/issues/305#issuecomment-6856947056)
correct M7 to the numbers above. The differences, each read again at the
source:

- `files/go` and `semantics/go` were counted as wholly model. Row R7d
  splits them.
- `semantics/horos/go` was counted as wholly model. It imports three
  contract packages and declares two contracts (`vocabulary.go:75-91`), so
  row R2 splits it like the three others.
- `trees/representation/go` was counted as a split. Every type of it is
  sealed or the result of one of its operations.
- `contracts/semantic/go`, `contracts/resolution/go` and
  `contracts/analysis/go` were counted as splits. Rows R2 and R7e keep them
  whole on the operations side, and the identities typed by them leave the
  vocabularies instead.
- `contracts/generation/go` was counted as a split. On `main` it is whole
  (M12). It splits on the migration's branch.

M7's statement stands with these numbers: every one of the 129 gets a new
import path, and no move is a prefix for all of them (row R31).

**The packages below `model` after the slices,** for the packages of `main`.
Derived here from the table above.

| Sort | Packages | Count |
| --- | --- | --- |
| `language` | `model/language/go` | 1 |
| `trees` | `model/trees/ast/go`, `model/trees/cst/go` | 2 |
| `syntax` | `model/syntax/go`; `model/syntax/<language>/go` for five languages; `model/syntax/<language>/lexical/go` for five languages | 11 |
| `semantics` | `model/semantics/go`; `model/semantics/<language>/go` for Go, horos, Rust and TypeScript | 5 |
| `files` | `model/files/go` | 1 |
| `source` | `model/source/go` | 1 |
| `intent` | `model/intent/go`; `model/intent/<language>/go` for C++, Go, Rust and TypeScript | 5 |
| `symbols` | `model/symbols/<language>/go` for Go, Rust and TypeScript | 3 |
| `families` | `model/families/horos/go` | 1 |
| `contributions` | None on `main` | 0 |
| | | **30** |

Eight of the 30 are today's wholly model packages at a new path, and 22 are
new packages that the splits create. The migration's branch adds
`model/contributions/go`, 24 contribution vocabularies and 6 family
vocabularies, and the pairs' branch a form table for horos.

## 8. Roles per root

A package gets its import rule from its role, and its role from its path
(`integration/go/architecture_roles_test.go:157`). The roots give a first
cut. The roles stay, and they do not all stay whole (M8, row R27).

| Role today | Under `model` | Under `ops` |
| --- | --- | --- |
| Foundation | The agnostic model packages, `model/<sort>/go`, among them `model/syntax/go`, and the two tree carriers. A closed table with one entry for each. A member imports only members of this table | The contract packages `ops/contracts/*/go` and the operations halves of the carriers: `ops/trees/representation/go`, `ops/files/go`, `ops/source/go`, `ops/semantics/go`. A closed table with one entry for each. A member imports members of both foundation tables and the kernel set, as today |
| Kernel | None | Whole (row R7f) |
| Shared private helper | None | `ops/internal/{concrete,grammar,lex,parseutil,tree}` |
| Vocabulary | `model/<sort>/<scope>/go`: the names of one language or one pair, and declared tables over them (M16). It imports the model-side foundations and the model packages of its own scope. It imports nothing under `ops` | `ops/semantics/<language>/go`: the identities of a language that are values of a contract package's type. It imports foundations and its own language's vocabulary, as today |
| Public part (new, M15) | None | A language's typed readers, typed constructors and symbol admission |
| Contribution of every kind, and common mapping part | None | Whole |
| Private part | None | `ops/internal/<language>`, `ops/internal/<language>sem`. A lexical helper imports its own language's lexical part, `model/syntax/<language>/lexical/go`, for the word lists it evaluates, and no other vocabulary (M19) |
| Composition | None | `ops/integration/go`, `ops/verification/integration/go`, `ops/tools/**` |

**The foundation table is divided (row R27).** `foundationRules` holds 31
entries today: six carriers, thirteen contracts, seven kernel packages and
five helpers of no single language
(`integration/go/architecture_test.go:162-203`). Two entries name a contract
package from a carrier: `trees/representation/go` lists
`contracts/language/go` (`:166`), and `source/go` lists
`contracts/semantic/go` (`:168`). After the splits of slice S6 the first is
an import from the operations side into the model, which is the permitted
direction, and the second lies inside the operations side. The model-side
table then lists model packages only.

**The vocabulary role exists on both sides (rows R2 and R27).** The accepted
definition lists "profile and contract identities" as vocabulary
(`docs/package-boundaries.md:98-99`). In code a contract identity is a value
of a contract package's type, contract packages are operations, and a model
imports none. So a profile string is a name on the model side, and a
contract built from it is an identity on the operations side.

**M15: the role of a language's typed readers, typed constructors and
symbol admission.** Row R27 asks that the reading half of a split
vocabulary package gets a named role, and no role of today fits it or
symbol admission. It is one role, **public part**, the public sibling of the
private part. The accepted term is: "**Private part.** Code that serves one
owner and is never selected by a consumer, such as a lexical helper of one
language." (`docs/package-boundaries.md:103-104`).

- *Definition.* A public part is code that serves one language, is never
  selected by a consumer, and is imported by several roles. It cannot live
  under `internal/`, because a mapping or a law that an independent library
  supplies imports it: "Private helpers are not an alternative public API."
  (`LAYOUT.md:133-134`).
- *What it may import.* The foundations of both roots, its own language's
  vocabularies under both roots and its own language's lexical helper. It
  imports no contribution and no semantic helper.
- *Who may import it.* A contribution that names its language (a stage
  provider, a query service, a source reader, a target, a mapping, a
  verification law), its language's private parts and a composition. No
  model package, no foundation or kernel package and no package of another
  language imports it.
- *Its members,* by the paths of M14:

  | Shape | Members on `main` after the splits |
  | --- | --- |
  | `ops/analyze/files/<language>/reader/go`, the typed readers | Go, Rust, TypeScript |
  | `ops/transform/syntax/<language>/intent/go`, the typed constructors with the issuer | C++, Go, Rust, TypeScript. horos on the pairs' branch |
  | `ops/transform/syntax/<language>/symbol/go`, symbol admission and issuing | Go, Rust, TypeScript |

- *Reason:* one role with one import rule, not three. The three share what
  matters for the rule: owned by one language, selected by no consumer,
  imported by several roles.
- *Not taken:* three roles. Leaving the three in vocabulary packages,
  against "and nothing else". Joining each to a contribution of its
  language: a mapping may not import a target, and an importer of a query
  service needs an exception that names its symbols.

**M16: a declared table in a vocabulary package.** A grammar table, a form
table and a word list are more than names: they declare what is admissible.
The packages `model/<sort>/<scope>/go` keep the role vocabulary, and the
definition gets a dated clarification: a vocabulary under `model` may hold,
beside names, a table over those names that declares what is admissible,
and it holds no function that evaluates the table. By M11 a language's
grammar table is part of its syntax vocabulary. By M19 its word lists are
a vocabulary of their own below it, the lexical part, which imports no
vocabulary of its language.

- *Reason:* no new role where the existing one fits (the derivation rule).
  The import rule is the same for names and for tables.
- *Not taken:* a role of its own for a package that holds a table.

### The sentences that change

Each sentence below stays readable as it was accepted, and a dated
clarification stands beside it. The clarifications are made in the change
that adds this page. They state the target and say that the code is not
there. The column "After the slices" is what each sentence says once the
[order of work](#10-the-order-of-work) is done. It uses the names of M10, M11, M14 and M19.

| Place | The sentence today | After the slices | Why | Source |
| --- | --- | --- | --- | --- |
| `LAYOUT.md:3-4`, owner decision of 31 October 2026 | "organize by processing stage and enforce a uniform structure across languages." | The root splits into `model` and `ops`. Operations are organized by processing stage with a uniform structure across languages, below `ops`. Models are organized by sort and then by language, below `model` | The direction changes the rule at the root and keeps it below `ops` | The direction; M6 |
| `LAYOUT.md:67`, in the hierarchy | "native semantic forms and validation" | `model/intent/<native>/go` holds the declared table of a language's intended forms. Their validation is the one checker of the kernel | A model contains no operation. The definition of a vocabulary already said "and nothing else" | M4; rows R2 and R7i |
| `LAYOUT.md:145-146` | "Directory groups are not Go packages or an assertion that all members share the same dependency level." | Directory groups below a root are not Go packages or an assertion that all members share the same dependency level. The two roots are such an assertion in one direction: no package under `model` imports a package under `ops` | A root that means "a model imports no operation" is exactly such an assertion | Row R27 |
| `LAYOUT.md:155`, role table, Foundation | "`trees/*/go`, `files/go`, `source/go`, `semantics/go`, `contracts/*/go`" | Two rows. On the model side: the agnostic model packages and the tree carriers. On the operations side: `ops/contracts/*/go` and the operations halves of the carriers | The root cuts the role in two | Row R27 |
| `LAYOUT.md:158`, role table, Vocabulary | "Foundations and its own language's vocabulary. `semantics/<language>/go` and `symbol` may also import their language's lexical helper." | On the model side: the model-side foundations and the model packages of its own scope, and nothing under `ops`. A syntax vocabulary also imports the agnostic syntax model, `model/syntax/go`, and its own language's lexical part. The lexical part imports at most `model/syntax/go` and never its language's syntax vocabulary. On the operations side: as today. No vocabulary under `model` imports a lexical helper | Rule 1; every caller of a lexical evaluator is on the operations side; the table type is model; structure knows spelling and not the reverse | Rows R2, R7h and R27; M11, M19 |
| `LAYOUT.md:167`, role table, Private part | "A lexical helper imports `internal/grammar` only, and never a contract or a stage provider." | A lexical helper imports the shared grammar helper and its own language's lexical part only, and never a contract or a stage provider. The "only" keeps its meaning: the lexical part holds word lists and predicate names and no node kind, so the helper still knows spelling and not structure | The word lists are data of the language's syntax model, and the evaluators read them | Row R7h; M19 |
| `LAYOUT.md:170-171`, one language | "A package that names a language imports no package of another language. A mapping names two." | The same, and a pair's contribution vocabulary under `model/contributions` names two as a mapping does | The model root gets packages that name a pair | Row R23 |
| `AGENTS.md:65-66` | "Identities live in vocabulary packages." | Identities live in vocabulary packages: a name on the model side, and an identity that is a value of a contract package's type on the operations side | The vocabulary role exists on both sides | Rows R2 and R27 |
| `AGENTS.md:179-181` | "Native semantic forms belong under transform/syntax/&lt;native&gt;/model/go, their AST interpretation under transform/syntax/&lt;native&gt;/target/go." | A language's intended forms are a declared table under `model/intent/<native>/go`. Their typed constructors are under `ops/transform/syntax/<native>/intent/go`, and their AST interpretation is under `ops/transform/syntax/<native>/target/go` | M4 | M4; M14 |
| `AGENTS.md:233-234` | "Follow the stage-first hierarchy in LAYOUT.md. Shared trees live under trees, shared contracts under contracts" | Follow the two roots of `LAYOUT.md`: models under `model`, by sort and then by language, and operations under `ops`, in the stage-first hierarchy. The generic trees live under `model/trees`, the shared contracts under `ops/contracts` | The direction | The direction; M6 |
| `AGENTS.md:251-252` | "Syntax depends on AST/CST carriers and identity, never on another language or parser/formatter. Native stages own their grammar and spelling." | A syntax vocabulary depends on the tree carriers, on language identity and on the agnostic syntax model, never on another language or on a parser or formatter. A language's syntax model declares its grammar table and its word lists. The native stages own the rules in code, the evaluation of a predicate and the spelling of concrete text | An operation does not own the declaration it interprets, and the table and the word lists are part of the syntax model | Rows R4, R7c and R7h; M11 |
| `CHARTER.md:80-81` | "Syntax packages depend on the AST/CST carriers and identity interface." | The same, and on two more model packages: the agnostic syntax model `model/syntax/go`, in whose table type a language's grammar table is written, and the language's own lexical part, whose predicate names the table uses | The table type is model, and the table lives in the syntax package | Row R7c; M11, M19 |
| `docs/package-boundaries.md:98-100`, Terms | "**Vocabulary.** The names two parties must agree on and nothing else: syntax kinds, semantic keys, profile and contract identities, intended native forms. A vocabulary package imports no contribution." | The same, and: the role exists under both roots. A vocabulary under `model` holds names and declared tables over them and imports nothing under `ops`. A contract identity is a value of a contract package's type and lives in a vocabulary under `ops` | Contradiction 2 of the evaluation: the definition, the direction and row R8 cannot all hold for one package | Rows R2 and R27; M16 |
| `docs/package-boundaries.md:225-226` | "`Validate` only requires two absolute paths. It reads no file and starts no process, so a vocabulary can hold it." | `Validate` evaluates a constraint, so it is an operation. The criterion "reads no file and starts no process" no longer admits a function to a vocabulary under `model` | Rule 2 with the test of M1 | M1; contradiction 1 of the evaluation |
| `docs/generation-composition.md:17`, ownership table | "Native model \| Native form identity and validity, without AST syntax" | Two rows. Native intent vocabulary: the declared forms of a language and their shapes, without AST syntax. Kernel form checker: the validity of a form against its language's declared table | M4 | M4; rows R2 and R7i |

Two sentences are kept and hold below both roots: "Do not create a
language-first packages tree, additional Go modules, compatibility aliases or
placeholder source directories." (`LAYOUT.md:134-135`), and "Mapping packages
depend on native models, never target AST constructors." (`AGENTS.md:190`).
The first also means that no directory under `model` is created before it
holds a package: no `model/relations` exists, because nothing on `main` or on
the branches is a declared relation in the sense of
[section 4](#4-one-languages-models).

## 9. Enforcement

**Decided (M9, row R30).** The advice lists five rules for the import test.

| Rule of the advice | Carried by | How |
| --- | --- | --- |
| 1. No import from `ops` into `model` | The import test | One prefix rule, beside the existing "a foundation or kernel package imports only foundation and kernel packages" (`integration/go/architecture_check_test.go:230-233`). It has a test that breaks it: a seeded import from a model package into an operations package fails |
| 2. No operation under `model` | Review. No check | Below |
| 3. No dependency across languages except the declared ones | The import test | It exists: the scope of a role and the table `languageOwners`. The language segment moves in the path, and the rule reads it at its new place |
| 4. No private part of an operation imported by a model | The compiler, and rule 1 | A subset of rule 1. Once the helpers are at `ops/internal`, Go refuses the import |
| 5. No cycles between model packages | The compiler | Go refuses an import cycle. The layering among the agnostic model packages is carried by the entries of the model-side foundation table |

**Why rule 2 gets no check.** Whether a package under `model` holds an
operation is a judgment about a responsibility. The owner's rule: "Reason
about responsibilities and change propagation; do not replace this with
keyword heuristics, scoring or new lint gates." (`AGENTS.md:36-37`). The
platform page says what a dependency check is and is not: "Such a check
enforces a boundary that reasoning has already decided. It does not
contradict the paragraph on lint gates above, and it does not replace that
reasoning: it cannot see a decision made inline, it cannot judge whether a
unit holds one contribution, and a green run shows only that the forbidden
edges are absent." (platform `docs/concepts/architectural-correctness.md:580-584`.)

A check that tried would also be wrong. `CheckExpression` in Go's
intended-forms package imports nothing forbidden, so no import rule sees it.
A rule on the presence of functions would refuse a model that is right:
`semantics/horos/go/vocabulary.go` has 31 function declarations
(`grep -c '^func ' semantics/horos/go/vocabulary.go`), and most of them
return one fixed name. So the reviewer of every change under `model` applies
the test of [section 2](#2-the-test). The advice says the same of its own
list: the test "cannot decide whether a package has been conceptually
classified correctly", and "A green architecture test establishes that
forbidden dependencies are absent. It does not prove that the package belongs
in the correct root." (advice, section 6, step 7.)

**How it fits the existing test.** The role function and the tables stay.
Slice S8 changes what they read:

- `packageRole` reads the root as the first segment and the roles of
  [section 8](#8-roles-per-root) below it.
- `foundationRules` (31 entries) becomes two closed tables, one for each
  root. `oneOfAKindRules` (12), `languageOwners` (5), `roleExtensions` (12)
  and the exception table (14) keep their meaning, with the new paths.
- `TestLanguageStageMatrix` lists `syntax` as the first of its native stages
  (`integration/go/architecture_test.go:15`). That column is under `model`
  after the move, so the matrix check reads both roots (row R22).
- The lexical part of a language has an allow-list of `model/syntax/go`
  and nothing else. So an import of its language's syntax vocabulary, and
  with it of a node kind, fails. The lexical helper's allow-list names the
  lexical part and not the syntax vocabulary (M19).
- Every added rule gets a synthetic test that breaks it, as the 37 tests of
  `architecture_synthetic_test.go` do for the present rules.

## 10. The order of work

**Decided (M7, M8; rows R27 to R29 and R31).** No move is a prefix for every
package: 25 packages split, so a move that only renames directories would
break rule 1 or rule 2 wherever a mixed package lands (row R31). The order:

1. This page: the test, the layout below the roots, the roles for each root.
2. The check that the migration added to four packages that will be models
   leaves them at once.
3. On `main`, each as its own change, the splits that touch no file of the
   migration's branch.
4. The split of the intended-forms packages on the migration's branch, as
   its own change, before the migration's gate is shown.
5. The scripted move of everything that is still at an old path, before the
   migration merges. The same script runs on the branches.
6. The gate of the migration, on the lasting boundaries, and its merge.
7. The Files formatters for Go and Rust are written at final paths.

| Slice | What it delivers | Where | Waits for | Validation, from the table "Applied" |
| --- | --- | --- | --- | --- |
| S0 | This page, with the clarifications in `LAYOUT.md`, `AGENTS.md`, `CHARTER.md`, `docs/package-boundaries.md` and `docs/generation-composition.md`, and the case bittype/ACX-015 | `main` | Nothing | "the page merged after an independent review that checks each quoted rule against its source" |
| S1 | `CheckMember` leaves the four intended-forms packages, and the target states it | The migration's branch, pull request #308 | Nothing | "a test that fails while a check sits in a model package, red on the branch before the move" |
| S2 | The split of the intended-forms packages ([section 5](#5-intent)) | The migration's branch, as its own pull request | S0 | "the five intended-forms packages hold no function; one checker test per table, red by mutant; outputs equal by digest" |
| S3 | The grammar-table type and the shared predicate names in `model/syntax/go`. Each language's grammar table in its syntax vocabulary, with predicates named by identity | `main` | S0. For each language also its lexical part of S4, whose predicate names the table uses | "the tables are data in a model package; the checker test red by mutant of a table entry" |
| S4 | The lexical part of each language, `model/syntax/<language>/lexical/go`: its word lists, character tables and literal patterns, and the names of its predicates, apart from their evaluators | `main` | S0 | "the evaluators import the tables; a test red by a changed table entry" |
| S5 | The semantic vocabularies: the typed readers and the contract identities leave, and the symbol packages split | `main` | S0 and #310 | "no model-side package imports a contract package; the import rule red by a seeded import" |
| S6 | The carriers: plain records apart from admission, language identity, and a type for captured Files | `main` | S0 | "the import rule red by a seeded import from a model carrier into an operation" |
| S7 | The identity types of a contribution get a model-side owner, and the pairs' vocabularies and the families follow | The migration's branch | S0 | "no vocabulary or family package imports the generation contract; the plan digests unchanged" |
| S8 | The scripted root move, the roles for each root, the import rule of the two roots | `main` and the branches | S2 to S7 | "the import rule of the two roots, red by a seeded import from model into ops; every test name passes as before the move" |
| S9 | The additions: a declared profile for each native syntax vocabulary, schemas of analysed facts with a shape check and generated names for the two bridges that are not written in Go, one test for each form | `main` | S8 | "per addition: its own red and green run" |

An issue is filed for a slice when it can start (the grade).

**What is free of the migration's branch.** The branch
`origin/feat/contribution-contract-196` at `af1dcde` changes production
files of 37 of the 129 packages of `main`, and leaves 92 alone:

```sh
git diff --name-only origin/main...origin/feat/contribution-contract-196 | grep '\.go$' | grep -v '_test\.go$' | grep -Ev '(^|/)testdata/' | xargs -n1 dirname | sort -u
```

intersected with the package list of [section 7](#every-package-of-main).
Slices S3 to S6 work in the 92.

**Why the split of the intended-forms packages is made on the branch.**
Every file that imports one of the four packages is a file that the branch
changes. Over all Go files of `main`,
`git grep -l 'transform/syntax/<language>/model/go"' -- '*.go'` gives 9
files for Go, 9 for TypeScript, 7 for Rust and 6 for C++, and each of them
is in the list of files that the branch changes. (The grade first gave 11,
10, 8 and 7. The decisions of 13 November 2026 correct them to these
numbers.) So the split on `main`
would be made against code that the migration deletes, and it would be made
again on the branch (the grade, M8). On the branch it is still its own
change, and it still comes before the gate.

`main` keeps the misplacement until the migration merges. That is a
deferral of a move, which
[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
permits only with a record: "Record such a deferral in the case with its
reason and the trigger that ends it; the violation stays open."
[bittype/ACX-015](architecture-cases/0015-models-hold-operations.md) records
it, with the merge of the migration as the trigger. This order differs from
the verifier's, who would make the split on `main` first (row R29).

**What the gate's line 7 has to do with it.** The migration "does not merge
while one line is red" (`docs/contribution-contract.md:2886`). Line 7 of
its gate is "A materially different contribution is added, and no shared
package changes", and its list of shared packages names "the native model
and target of each language" (`docs/contribution-contract.md:2966`). The
split redraws exactly those packages. Shown before the split, line 7 would
be evidence for a boundary that is about to change, and
[AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07)
lets evidence support only the claim it establishes. So the gate is shown
after S2 and after S8, on the boundaries that last. Line 7 is no freeze
against a correction: it measures the diff of adding a contribution.

**"Each split creates its new half at its final path" (M7).** Between the
first split and slice S8 both layouts exist. Derived here, these rules hold
from the first split on:

1. **A half that leaves into a new package is created at its final path,**
   under `model/` or `ops/`. The half that stays keeps its old path until
   S8. Nothing moves twice: the script of S8 moves only what is still at an
   old path.
   - A new package is created at its final path. By M19 the lexical part of
     a language is a new package and is created at
     `model/syntax/<language>/lexical/go` by slice S4. The table type is a
     new package and is created at `model/syntax/go` by slice S3.
   - Only a part that joins an existing package joins it where that
     package is. That is the grammar table: by M11 it joins
     `syntax/<language>/go`, which is wholly model and takes its prefix in
     S8 like every whole package. A checker imports that package already,
     so the table reaches it through an import line that the script
     rewrites in any case.
   - The five syntax vocabularies are not moved earlier than S8. Of the
     files of `main` that import one of them, the migration's branch
     changes 11 of 21 for C++, 17 of 47 for Go, 44 of 83 for horos, 12 of
     39 for Rust and 17 of 54 for TypeScript
     (`git grep -l '"github.com/Bitspark/bittype/syntax/<language>/go"' origin/main -- '*.go'`,
     intersected with the files the branch changes).
2. **A package under `model/` obeys rule 1 from the day it exists.** It
   imports packages under `model/`, and of the packages at old paths only
   those that [section 7](#every-package-of-main) marks as wholly model. It
   imports no package that the table places under `ops`, and no package that
   is still to be split.
3. **A package under `ops/` has the role of [section 8](#8-roles-per-root)
   from the day it exists.** It imports what that role permits, at old or
   new paths.
4. **The import test learns each new package with its split.** The change
   that creates a package adds its rule and a seeded import that the rule
   refuses. Until S8 the rule for a package under `model/` is an entry that
   lists its imports. S8 replaces those entries by the one prefix rule.
5. **`internal/` stays at the module root until S8.** A package under
   `ops/` that is created earlier imports a helper at its old path. Slices
   S4, S5 and S6 remove the imports of a helper from every future model
   package, and then the script moves `internal/` under `ops`.
6. **No alias stands in for an old path.** "Do not add compatibility aliases"
   (`AGENTS.md:236-237`).

**What the script of S8 is.** The table of [section 7](#every-package-of-main)
as a list from old path to new, for every package and half that is still at
an old path. It rewrites the import strings of the module
(`git grep -o '"github.com/Bitspark/bittype/[^"]*"' origin/main -- '*.go' | wc -l`
gives 2,137 in 367 files), the paths in the workflow files, the role
function and its tables, and the path prefixes that the documentation tool
reads. No identity of a model or of a contribution is a package path
([section 6](#6-declared-models)), so the move changes none. The numbers
say how large the step is. They are no reason for any choice on this page
(AC-09).

## 11. Decisions and deferrals

**The decisions.** Each is the research owner's, made on 13 November 2026
and recorded on #305: M1 to M9 in the
[grade](https://github.com/Bitspark/bittype/issues/305#issuecomment-6855741567),
M10 to M18 in the [decisions on this page's nine choices](https://github.com/Bitspark/bittype/issues/305#issuecomment-6856947056),
and [M19](https://github.com/Bitspark/bittype/issues/305#issuecomment-6857366489).

| Decision | What it decides | Rows | Section |
| --- | --- | --- | --- |
| M1 | The test, with the two added sentences | R1 to R4, R7a to R7h, R8 | [2](#2-the-test) |
| M2 | Consistency and correctness, four levels, no new evidence concept, four gaps | R5, R6 | [3](#3-consistency-and-correctness) |
| M3 | Three related models of one language, with a stated relation and a test for each form | R9 to R13 | [4](#4-one-languages-models) |
| M4 | The intended-forms packages: a declared table, one checker, one constructor | R7i, R14, R15, R24, R26 | [5](#5-intent) |
| M5 | Declared models with the mechanisms that exist, in the order of AC-10, identified by a profile | R16 to R20 | [6](#6-declared-models) |
| M6 | The layout: two roots, sort first below `model`, the stage matrix below `ops`, one module | R21, R22 | [7](#7-the-layout) |
| M7 | No prefix move: each split creates its new half at its final path. Of 129 packages 8 are wholly model, 25 split and 96 are wholly operations, as M10 to M18 correct the count | R28, R31 | [7](#every-package-of-main), [10](#10-the-order-of-work) |
| M8 | The order, with the work in flight | R27, R29 | [8](#8-roles-per-root), [10](#10-the-order-of-work) |
| M9 | Enforcement: one prefix rule by imports, rule 2 in review | R30 | [9](#9-enforcement) |
| M10 | The shape of a path below `model`: `model/<sort>/go`, `model/<sort>/<scope>/go`, `model/<sort>/<scope>/<name>/go` | R21 | [7](#names-below-the-roots) |
| M11 | Ten sorts. `grammar` and `lexical` are no sorts. A language's grammar table lives with its node kinds in `model/syntax/<language>/go`, the table type and the shared predicate names in `model/syntax/go`. M19 adjusts one point: M11 first put the lexical tables into the same package | R21, R7c, R7h | [7](#names-below-the-roots) |
| M12 | The model side holds the form-table type and no value type of intent. This narrows M4 and the reading of row R15, part 3 | R15 | [5](#5-intent) |
| M13 | The form-table type stands as a sketch, and slice S2 makes it final | R15, R24 | [5](#5-intent) |
| M14 | The paths of the operations halves, with the general rule that a half keeps its path below its root | R21, R27 | [7](#names-below-the-roots) |
| M15 | One role, public part, for typed readers, typed constructors and symbol admission | R27 | [8](#8-roles-per-root) |
| M16 | A declared table lives in a vocabulary package, and the definition is clarified | R27 | [8](#8-roles-per-root) |
| M17 | Nothing refuses Go's record-body condition after the split, and row R24's change of category covers the primitive set only. An unknown form name stays unsupported | R14, R24 | [5](#5-intent) |
| M18 | A grammar table names a predicate by an owner-qualified name | R7c | [6](#6-declared-models) |
| M19 | A language's lexical tables and the names of its predicates are a part below its syntax model, `model/syntax/<language>/lexical/go`. It imports no node kinds, the syntax vocabulary may import it, and the lexical helper imports it only. It adjusts M11 | R7h | [7](#names-below-the-roots) |

Two rows belong to no single decision. Row R23 places the identities of a
contribution ([section 7](#every-place-named), slice S7). Row R25 takes one
principle from the advice's four outcomes and no more: a valid value that an
operation cannot process is not invalid. Each operation keeps its accepted
classes, which have a category for a limit and another meaning of
"rejected". The principle is applied where it bites: in the one change of category of
[section 5](#5-intent), and in the end of the refusal of Go's record body
(M17).

**Deferred, with triggers.** A trigger is a state of a repository that
anyone can check.

| What is deferred | Why | The trigger that ends it | Source |
| --- | --- | --- | --- |
| A shared identity for single notions of intent and analysed semantics, such as Go's slice | No reader exists | A non-test file imports both a language's intent vocabulary and its semantic vocabulary. The command is in [section 4](#4-one-languages-models) | Row R13a |
| The move of the checks out of the four intended-forms packages on `main` | Every importer is a file that the migration's branch rewrites, so the move is made once, on the branch | The migration of #196 merges. The violation stays open until then, in [bittype/ACX-015](architecture-cases/0015-models-hold-operations.md) | M8; AC-10 |
| A relation from a form to what analysis reports, as data | No package would read it | The first package that must read such a relation without running the test of the form | Row R10 |
| A plain record of a native intent node on the model side | Nothing outside the generation contract would read it | The first package outside the generation contract that must hold native intent without a planner | M12 |

**What the consultation did not answer.** The grade and the table "Applied"
list these. Each is answered here or open.

| Point | State |
| --- | --- |
| Decision O45, which the advice does not mention | It stands as the research owner's. Silence is no confirmation ([section 4](#4-one-languages-models)) |
| A failure of a comparable system | None was named. Nothing on this page rests on one |
| When the root move itself happens | M8: slice S8, before the migration merges |
| What the split of a sealed Go type is | M1, sentence 1 |
| The six places that the advice's tree does not name | [Section 7](#every-place-named), with M10, M11, M12, M14 and M19 |
| Whether what stays code, the rules between nodes and the predicates, is still model | M1, sentence 2: the model names it, the evaluator states it, and the declared model is an upper bound |
| The role of an operations-side package that is neither contribution nor private part | M15: public part |
| What `model/relations` is, against `model/contributions` | No directory for relations is created. M11 names the sorts `families` and `contributions` |
| Whether the payload of an analysed fact is declared as a Go table or as a horos type | **Open.** Decided with slice S9 |
| How a declared model relates to a rule that the building toolchain owns, such as Go's keywords and its universe scope | **Open.** The model names the rule (sentence 2). Go's lexical part holds the names of its predicates and no word list ([section 7](#names-below-the-roots)). No profile records the toolchain's release today. Decided with slice S9, which adds the profiles of the native syntax vocabularies |
| Whether adding an intended form is an edit of a shared package in the sense of the gate's line 7 | **Open.** It belongs to the gate of [the contribution contract](contribution-contract.md#the-gate) and is decided there |
| What a family's meaning is for a source that is not horos | **Open**, by decision O46 on #196, item 6. Slice S7 places the family packages and does not answer it |
