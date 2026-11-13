# Models and operations

**Design, 13 November 2026, from research 0018. The decisions M1 to M9 are the research owner's, `claude-contract-chain-issues`, graded on [#305](https://github.com/Bitspark/bittype/issues/305#issuecomment-6855741567). The nine points marked PROPOSAL are this page's and are not decided. Nothing on this page is implemented: the code on `main` has the layout that [`LAYOUT.md`](../LAYOUT.md) describes.**

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
  file and line at `0daee73`.
- **M1 to M9:** the nine decisions of the
  [grade](https://github.com/Bitspark/bittype/issues/305#issuecomment-6855741567).
  bittype has no architecture seat, so the research owner graded and said so
  there.
- **Row R1 to R31:** a row of the table "Applied" of research 0018. Its Action
  is the decision. Where it differs from the advice, the Action governs.
- **Derived here:** a consequence that this page draws from the sources
  above. It is no decision. An objection to it is weighed.
- **PROPOSAL P1 to P9:** a choice that none of the sources settles. Each
  gives the alternatives and what each costs. The research owner decides
  them, and they then get a decision number.
  [Section 11](#11-decisions-and-deferrals) lists all nine.

Code is quoted with its file and line on `main` at `0daee73`. The existing
code is given as a fact and never as a reason for a choice
([AC-09](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-09)).
Counts come with the command that gives them.

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
| The names below the two roots in message 7, `model/{syntax,semantics,...}/go/...` and `ops/{transform,...}/go/...` | A sketch. [Section 7](#7-the-layout) decides the order below a root with M6 and proposes the names |
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
  native forms" (`docs/generation-planner.md:529-530`).
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
| The checks of well-formedness | **One** checker in the agnostic layer, which reads a language's table. It lives in the kernel package that holds the single-scope check, `transform/syntax/shared/target/go` (`check.go:11`) | M4; row R24. [AC-11](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-11): "Per-language copies of that mechanism are forbidden." |
| The constructors | **One** constructor in the agnostic layer. A language keeps typed constructors over it, in a package that mappings may import | M4; rows R15 and R24 |
| Lowering and coverage | The target keeps both. It passes its language's table to the checker and adds its own coverage answers | M4; row R15, part 5. The accepted contract gives the checks to the target: "The planner invokes its checks against intended native forms and the completed program." (`docs/generation-spi.md:332-333`) |
| The condition "a Go defined type has a record body" (`model.go:105-110`) | It leaves the package. It is neither Go's rule nor the target's. Its owner is the contribution that supplies the declaration | M4; row R14. "alias or defined declaration \| The contribution that supplies the declaration: the default declaration policy or an exact mapping" (`docs/contribution-contract.md:203`) |
| The intent carrier | Only the native kind of an expression, a form with its data and its arguments, is intent. The planner's other kinds stay one unit with the planner | Row R15, part 3. "Parts that share a private invariant belong together. The planner and the plan types in `contracts/generation/go` are one unit, because only the planner may issue a plan." (`docs/package-boundaries.md:145-147`) |

**Why one checker.** The four packages on `main` are near copies: the same
issuer, the same refusal of a foreign form, the same arity checks and the
same record rule (`transform/syntax/{cpp,golang,rust,typescript}/model/go/model.go`).
The pairs' branch adds a fifth for horos. Decision
[O46 on #196](https://github.com/Bitspark/bittype/issues/196#issuecomment-6855097337)
lists them as AG-7, "five near copies", and leaves them to this evaluation.
The mechanism needs no knowledge of one language, so AC-11 gives it one owner,
and each language supplies its table.

**What the Go target keeps as coverage.** Three answers of
`CheckExpression` are statements about what today's generation builds and no
rules of a form: which primitive names exist (`model.go:63`), the refusal of
every external type (`:87-88`) and the refusal of an unknown form name
(`:89-90`). A model is not narrowed to what one operation supports (row R26;
[AC-06](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-06)).
Rust keeps its coverage answer for external types in the target already
(`transform/syntax/rust/model/go/model.go:57-60`).

**Two changes of category, made on purpose (row R24).** Failure meaning is
one of the guarantees that composition preserves (AC-05), so a changed
category is recorded and never a side effect.

1. A primitive name outside the two that the target lowers is invalid today
   (`model.go:63-64`). As the target's coverage answer it is unsupported.
2. A `defined` declaration over a body that is no record is invalid today
   (`model.go:107-108`). After the split no shared package calls it invalid.
   Row R24 records this as a change "from invalid to unsupported". Row R14
   says of the same condition that "as a contribution's answer it is a
   refusal". The two rows agree that it is no longer invalid and name the
   new category differently. PROPOSAL P8 below says how this page reads them.

**PROPOSAL P3: where the intent carrier is cut.** No source says whether the
native kind gets a type of its own on the model side.

- *Proposed:* the model side holds the form-table type and nothing that
  carries a value. `Form[T]`, `FormIssuer[T]` and `Expression[T]` stay whole
  in the generation contract package, with the planner. The reasons:
  `Form[T]` is created only by its issuer, which refuses an empty name
  (`contracts/generation/go/expression.go:18-23`), so sentence 1 of the test
  puts it on the operations side. `Expression[T]` has unexported fields, and
  one of its six kinds is issued by the planner alone (`expression.go:52-63`),
  so the accepted sentence on the planner's unit keeps it there. A native
  node's arguments can be any of the six kinds, so a tree of plain native
  nodes cannot hold a mapping's output. What a native node is made of is
  stated by the table type. The agnostic constructor is then the function
  that creates an `Expression[T]`, and Go keeps it in that type's package.
- *Alternative A:* a plain record of one native node on the model side, with
  a form name, data and arguments of a type parameter, which the sealed
  expression holds for its native kind. It gives intent the same pair as
  trees have, plain data on the model side and an admitted state on the
  operations side. It costs a second statement of a node's shape beside the
  table type, and a type that nothing outside the generation contract would
  read today. It can be added later without moving anything: its trigger
  would be the first package outside the generation contract that must hold
  native intent without a planner.
- *Alternative B:* `Expression[T]` moves to the model side with exported
  fields. It costs the guarantee that only the planner issues a declaration
  reference, and it contradicts the sentence quoted above.

**PROPOSAL P4: the form table, as a sketch that slice S2 makes final.** The
type names, the field names and the profile string are proposals. What the
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
| `primitive` | One datum, no argument, and the datum is `string` or `bool` (`:63`) | One datum that is not empty, no argument | Which names the target lowers. That is coverage (change 1 above). That the name is a predeclared type of Go is a constraint that exists only as its evaluator, the universe scope of the building toolchain (`transform/syntax/golang/symbol/go/symbol.go:34`) |
| `slice` | No datum, one argument (`:67`) | The same | |
| `field` | Two data, the first not empty, the second `true` or `false`, one argument (`:71`) | The same | That the name is a Go identifier. The Go checker refuses another name when the output is sealed (`check/syntax/golang/provider/go/check.go:27`). What `optional` means for Go: the target reads the name only (`transform/syntax/golang/target/go/provider.go:122-123`) |
| `record` | No datum, every argument is a `field`, no name twice (`:75-86`) | The same | |
| `external` | Refused as unsupported before its shape is read (`:87-88`). Its constructor writes two data and any arguments (`:49-51`) | Two data and any arguments, as the constructor states them. No rule for either datum | Whether Go builds it. That is coverage |
| `alias` | A declaration form (`:102`) | A declaration form without a rule for its body | |
| `defined` | A declaration form whose body is a `record` (`:105-110`) | A declaration form without a rule for its body (row R14) | The declaration policy's choice between `alias` and `defined` (`transform/syntax/horos/mapping/declarations/go/declarations.go:52-55`) |

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
target import. [Section 7](#7-the-layout) names its path (PROPOSAL P5) and
[section 8](#8-roles-per-root) its role (PROPOSAL P6).

**PROPOSAL P8: the category of two refusals after the split.**

- *The record body of Go's `defined`.* Proposed: no shared package refuses
  it. Go's table states no body rule, the Go target lowers a `defined`
  declaration over any body (`provider.go:96-106`), and the default
  declaration policy goes on choosing `defined` for a record and `alias`
  otherwise. A contribution that is asked for what its policy does not give
  answers with a refusal of its own class under the contribution contract.
  This is how row R14 reads. The alternative is row R24 read literally: the
  Go target refuses the case as unsupported. It costs a policy in the target
  that row R14 says is not the target's. Either way the change is recorded
  in the pull request of slice S2 with a test that names the old and the new
  answer.
- *An unknown form name.* Proposed: it stays unsupported, as today
  (`model.go:89-90`), and the answer is the one checker's. Row R24 records
  two changes of category, and a third would need its own decision. The
  alternative is invalid, by the rule of the Files contract that an input
  "malformed under the published vocabulary" is invalid
  (`docs/source.md:512-513`). It costs a third change of category and an
  answer to whether a later revision of a vocabulary is "a later slice" in
  the sense of that rule.

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

**PROPOSAL P9: how a grammar table names a predicate.** An entry of a
grammar table holds a function value today (`Attribute.Accept`,
`internal/grammar/rules.go:17`). Go's table binds 17 of them, so it is not
data, and row R7c makes the naming of predicates by identity the
precondition of its move.

- *Proposed:* an entry names its predicate by an owner-qualified name, as a
  semantic key names a fact (`semantics.Key`, `semantics/go/model.go:24`).
  The names that every language can use (`bool`, `nonempty`, `string`,
  `decimal`) are declared in `model/grammar/go`. A language's own, such as
  Go's `identifier`, are declared in its `model/lexical/<language>/go`. A set
  of allowed strings is data in the entry and no predicate. The checker that
  calls the interpreter supplies the evaluators by name, and a name without
  an evaluator fails the check and is never skipped. Where the evaluator is
  the building toolchain, as for Go's identifier (`check.go:13`), the name is
  all that the model holds, and the language's page says so (sentence 2 of
  the test).
- *Alternative A:* a closed enumeration of predicates in the package of the
  table type. It costs an edit of a shared package for each new predicate of
  one language
  ([AC-03](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-03)).
- *Alternative B:* the tables stay with their checkers and keep their
  function values. Row R4 refuses it: an operation does not own the
  declaration it interprets.

## 7. The layout

**Decided (M6, rows R21 and R22).**

| Point | Decision | Source |
| --- | --- | --- |
| The roots | Two: `model` and `ops`. No third root | The direction; row R8 |
| Below `model` | The sort of model first, then the language | M6. The owner's sketch is `model/{syntax,semantics,...}/go/...`. The standing rule is "Do not create a language-first packages tree" (`LAYOUT.md:134`, `AGENTS.md:236-237`). The advice's tree puts the language first and its own section 4.3 says the opposite: "organized around the kinds of models being declared, with language identity as a further dimension" |
| Below `ops` | The stage matrix as it is: `<stage>/<dimension>/<scope>/<component>/go` below the root | M6; row R22. The owner's sketch is `ops/{transform,...}/go/...` |
| Modules | One. "One public root Go module contains independently owned packages" (`LAYOUT.md:6`) | M6. Its cost: the module graph cannot carry rule 1. The import test and Go's rule for `internal` carry it ([section 9](#9-enforcement)) |
| `internal/` | It goes under `ops`. It stays at the module root until no future model package imports a helper | M6; row R21 |
| Integration and tools | Under `ops`: `ops/integration`, `ops/tools`, and `ops/verification` | Row R21. With the helpers at `ops/internal`, Go lets only code below `ops/` import them. `integration/go` imports `internal/horosfiles` and `internal/tree`, and the tools import `internal/grammar`, `internal/horosmodel` and `internal/lex` |
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

**PROPOSAL P1: the shape of a path below `model`.**

- *Proposed:* two shapes, as `semantics/go` and `semantics/<language>/go`
  have them today. `model/<sort>/go` is the agnostic part of a sort: the
  types in which every language's declaration of that sort is written.
  `model/<sort>/<scope>/go` is one scope's declaration of that sort, where a
  scope is one language or an ordered pair. Where one scope declares several
  independent things of one sort, each is `model/<sort>/<scope>/<name>/go`,
  the path form that #222 decided for several contributions
  (`LAYOUT.md:100-104`). The two tree carriers keep their own directory,
  `model/trees/ast/go` and `model/trees/cst/go`. This is M6 read literally:
  sort, then language, and a sort's agnostic part beside its languages.
- *Alternative A:* one shape with an explicit scope,
  `model/<sort>/<scope>/<part>/go`, where the scope is a language, a pair or
  `shared`, as the scope of an operations path is ("A scope identifies one
  language, an ordered pair or shared machinery", `LAYOUT.md:98`). The trees
  would be `model/syntax/shared/ast/go`, the graph records
  `model/semantics/shared/go`. One function reads a role from every path,
  and no directory holds a package beside language directories. It costs the
  name `shared` in every agnostic path, and it puts a model directory named
  `shared` next to the advice's `model/shared`, which collected whatever is
  hard to classify.
- *Alternative B:* the proposed shapes, with the grammar table and the word
  lists as parts of the syntax sort, `model/syntax/<language>/grammar/go` and
  `model/syntax/<language>/lexical/go`. It keeps one directory for all that
  says what can be written in a language, which is how M3 counts the models.
  It costs a home for the agnostic table type that is neither
  `model/syntax/go` nor a language's.

**PROPOSAL P2: the names of the sorts.** A name that the repository has is
kept. Each new name is a word that the pages already use for the thing.

| Sort | Directory | Holds | The word comes from |
| --- | --- | --- | --- |
| Language identity | `model/language` | `ID` and `Provider` | `contracts/language/go` today |
| Trees | `model/trees` | The generic abstract and concrete tree | `trees/` today |
| Syntax | `model/syntax` | A language's node kinds, roles, attributes and its syntax profile | `syntax/` today |
| Grammar | `model/grammar` | The table type, and each language's table of what a node kind may hold | `internal/grammar` today |
| Lexical | `model/lexical` | A language's word lists and character tables, and the names of its lexical predicates | "private lexical helpers" (`LAYOUT.md:73`) |
| Semantics | `model/semantics` | The plain records of the graph, and a language's keys, profiles and payload shapes | `semantics/` today |
| Files | `model/files` | The plain records of a file set | `files/` today |
| Source | `model/source` | The plain records of a Source: unit, group, relation | `source/` today |
| Intent | `model/intent` | The form-table type, and each language's form table | "native intent" (`docs/generation.md:77`). The alternatives are `forms`, after "intended native forms" (`LAYOUT.md:467`), and keeping `model`, which would give `model/model` |
| Symbols | `model/symbols` | A language's plain symbol record and its categories | `symbol` today. The alternative is the singular |
| Families | `model/families` | The profile identities and the family vocabularies of a source language | "family" (`docs/representation-families.md`). On the migration's branch |
| Contributions | `model/contributions` | The identity types of a contribution, and each pair's contribution vocabularies | "contribution" (`docs/contribution-contract.md`). On the migration's branch |

**PROPOSAL P5: the paths of the operations halves.** A package or a half
that this page does not name otherwise keeps its path below its root. That
is all that happens to a whole operations package (row R22). Four halves
need a name of their own.

| Half | Path | Why there | Alternative, and its cost |
| --- | --- | --- | --- |
| The contract identities, evidence requirements and settings of a language's analysis | `ops/semantics/<language>/go` | It is what remains at today's path once the names have their model package | `ops/analyze/files/<language>/contracts/go`. It costs a path under one stage for identities of two, resolution and analysis |
| The typed readers of a language's analysed graph | `ops/analyze/files/<language>/reader/go` | They open "an admitted analysis receipt" (`semantics/golang/go/model.go:41`), so they belong to the analyze stage of the Files dimension, beside the provider and the query service | Joining the query service, `analyze/files/<language>/queries/go`. A query service is a contribution, and a package that is no composition imports one only under an exception that names its symbols. Rust has none |
| The typed constructors of a language's intent, with its issuer | `ops/transform/syntax/<language>/intent/go` | It is today's `transform/syntax/<language>/model/go` without the table and the checks. The component is renamed, because a directory `model` below `ops` would say the opposite of the roots ([AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10) names a wrongly named responsibility with a misplaced one) | Keeping `model` as the component. It costs that confusion. Putting the constructors into the target: a mapping may not import a target (`AGENTS.md:190`) |
| The admission and the issuing of a language's symbols | `ops/transform/syntax/<language>/symbol/go` | It is today's path below the root | None is proposed |

### Every place named

The advice's tree leaves six places without a name (row R21). The
migration's branch adds four more that this page has to place (row R23).

| What | Path | Side | Standing | Slice |
| --- | --- | --- | --- | --- |
| The Source records `Unit`, `GroupID`, `Group` and `Relation` (`source/go/source.go:20-37`) | `model/source/go` | Model | Row R7d places them. The name is P2 | S6 |
| Language identity, `ID` and `Provider` (`contracts/language/go/language.go:11-25`) | `model/language/go` | Model | Row R3 places it. The name is P2 | S6 |
| The intent carrier | `ops/contracts/generation/go`, where it is | Operations | P3 | None |
| The form-table type | `model/intent/go` | Model | M4 places it. The name is P2 | S2 |
| The grammar-table type (`internal/grammar/rules.go:15-33`) | `model/grammar/go` | Model | Row R7c places it. The name is P2 | S3 |
| Lexical tables | `model/lexical/<language>/go` | Model | Row R7h places them. The name is P2 | S4 |
| The operations halves of split vocabulary packages | The four paths of P5 | Operations | P5 | S2, S5 |
| The identity types of a contribution (`ContributionIdentity`, `Meaning`, `ContributionReference`, `SchemaReference`, `EvidenceKind`; on the migration's branch `contracts/generation/go/contribution.go:22-46`, `:138`) | `model/contributions/go` | Model | Rows R7e and R23 place them. The name is P2 | S7 |
| A pair's contribution vocabulary (on the branch `transform/syntax/<source>-to-<target>/vocabulary/<name>/go`, 24 packages) | `model/contributions/<source>-to-<target>/<name>/go` | Model | Row R7g places it. The name is P2 | S7 |
| A family of one language (on the branch `transform/syntax/horos/family/<name>/go`, 6 packages), and the two profile identities of `transform/syntax/horos/mapping/go` | `model/families/<language>/<name>/go`, and `model/families/<language>/go` for the profile identities | Model | Row R23 asks for the place. The name is P2 | S7 |
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

It gives 129 at `0daee73`. "Whole" means that the package takes the prefix
of its root and nothing else changes. The paths in the column "After" use
the proposed names (P1, P2, P5). The side of each part is decided by the row
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
| `syntax/{cpp,golang,horos,rust,typescript}/go` | `model/syntax/<language>/go`, whole, once language identity is at `model/language/go` | Rows R3 and R7d |
| `semantics/horos/go` | **Splits.** The profiles, the keys, the descriptor prefix, the term encoding, the limits setting key and the prelude pin to `model/semantics/horos/go`. `ResolutionContract`, `AnalysisContract`, `AnalysisContractWithResolutionConfiguration` (`vocabulary.go:75-91`) and the reader `DescriptorConstructor` (`:55`) to `ops/semantics/horos/go` | Row R2 |
| `semantics/{golang,rust,typescript}/go` | **Splits in three.** The profiles, the key functions, the payload shapes and the plain records of a compiler run's inputs to `model/semantics/<language>/go`. The contracts, the evidence requirement, the requirements, the configuration and the encoding of inputs into settings to `ops/semantics/<language>/go`. `Model`, `Open` and the typed readers to `ops/analyze/files/<language>/reader/go` | Rows R2 and R7b |
| `transform/syntax/{cpp,golang,rust,typescript}/model/go` | **Splits.** The form names, the primitive names and the form table to `model/intent/<language>/go`. The issuer and the typed constructors to `ops/transform/syntax/<language>/intent/go`. The checks become the one checker in `ops/transform/syntax/shared/target/go`. The coverage answers go to `ops/transform/syntax/<language>/target/go`. Go's record-body condition goes to the declaration contribution | M4; [section 5](#5-intent) |
| `transform/syntax/{golang,rust,typescript}/symbol/go` | **Splits.** The plain record of a symbol, its categories and the names it declares to `model/symbols/<language>/go`. The refusing constructors, the issuer, `ExternalUse` and `ReadExternal` to `ops/transform/syntax/<language>/symbol/go` | Row R7a; M8, slice S5 |

**Contracts (13).**

| Today | After | Source |
| --- | --- | --- |
| `contracts/language/go` | `model/language/go`, whole | Row R3 |
| `contracts/{analysis,check,format,parse,pipeline,print,resolution,semantic,text,transform,verification}/go` | `ops/contracts/<name>/go`, whole | Rows R7e and R8 |
| `contracts/generation/go` | `ops/contracts/generation/go`. On `main` it is whole under P3. On the migration's branch the identity types of a contribution leave it for `model/contributions/go` | Rows R7e and R23; P3 |

**Kernel (7).** All whole, to `ops/<same path>` (row R7f): `engine/go`,
`transform/syntax/shared/generation/go`, `transform/syntax/shared/target/go`,
`check/syntax/shared/ast/go`, `check/syntax/shared/cst/go`,
`print/syntax/shared/provider/go`, `print/files/shared/provider/go`. The
kernel package `transform/syntax/shared/target/go` receives the one form
checker in slice S2.

**Stage providers and query services (28).**

| Today | After | Source |
| --- | --- | --- |
| `check/syntax/{cpp,golang,rust,typescript}/provider/go` | **Splits.** The grammar table to `model/grammar/<language>/go`. The call of the interpreter, the rules in code and the provider stay, at `ops/<same path>` | Row R7c |
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
| `internal/grammar` | **Splits.** The table types `Attribute`, `Role`, `Rule`, `Rules` (`rules.go:15-33`) to `model/grammar/go`. The interpreter `Check` (`:80`), the entry constructors and the predicates stay, at `ops/internal/grammar` | Row R7c |
| `internal/{cpp,horos,rust,typescript}` | **Splits.** The word lists and the character tables to `model/lexical/<language>/go`. The evaluators and the literal decoders stay, at `ops/internal/<language>` | Row R7h |
| `internal/declsyntax` | **Splits.** The horos grammar table (`check.go:67`) to `model/grammar/horos/go`. `Check` and the rules in code stay, at `ops/internal/declsyntax` | Row R7c |

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

**The sum.** By this table 8 packages of `main` are wholly model, 25 split
and 96 are wholly operations. M7 gives 11, 27 and 91, from the placement
table of verification report 2. The difference is in five places, each read
again at the source:

- `files/go` and `semantics/go` are counted there as wholly model. Row R7d
  splits them.
- `semantics/horos/go` is counted there as wholly model. It imports three
  contract packages and declares two contracts (`vocabulary.go:75-91`), so
  row R2 splits it like the three others.
- `trees/representation/go` is counted there as a split. Every type of it is
  sealed or the result of one of its operations.
- `contracts/semantic/go`, `contracts/resolution/go` and
  `contracts/analysis/go` are counted there as splits. Rows R2 and R7e keep
  them whole on the operations side, and the identities typed by them leave
  the vocabularies instead.
- `contracts/generation/go` is counted there as a split. On `main` it is
  whole under P3. It splits on the migration's branch, and on `main` only
  if alternative A of P3 is taken.

None of the five changes a decision. Every one of the 129 still gets a new
import path, and no move is a prefix for all of them (M7, row R31).

## 8. Roles per root

A package gets its import rule from its role, and its role from its path
(`integration/go/architecture_roles_test.go:157`). The roots give a first
cut. The roles stay, and they do not all stay whole (M8, row R27).

| Role today | Under `model` | Under `ops` |
| --- | --- | --- |
| Foundation | The agnostic model packages, `model/<sort>/go` and the two tree carriers. A closed table with one entry for each. A member imports only members of this table | The contract packages `ops/contracts/*/go` and the operations halves of the carriers: `ops/trees/representation/go`, `ops/files/go`, `ops/source/go`, `ops/semantics/go`. A closed table with one entry for each. A member imports members of both foundation tables and the kernel set, as today |
| Kernel | None | Whole (row R7f) |
| Shared private helper | None | `ops/internal/{concrete,grammar,lex,parseutil,tree}` |
| Vocabulary | `model/<sort>/<scope>/go`: the names of one language or one pair, and declared tables over them (P7). It imports the model-side foundations and the model packages of its own scope. It imports nothing under `ops` | `ops/semantics/<language>/go`: the identities of a language that are values of a contract package's type. It imports foundations and its own language's vocabulary, as today |
| Public part (new, P6) | None | A language's typed readers, typed constructors and symbol admission |
| Contribution of every kind, and common mapping part | None | Whole |
| Private part | None | `ops/internal/<language>`, `ops/internal/<language>sem` |
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

**PROPOSAL P6: the role of a language's typed readers, typed constructors
and symbol admission.** Row R27 asks that the reading half of a split
vocabulary package gets a named role. Verification report 2 found that no
role of today fits it or symbol admission.

- *Proposed:* one role, **public part**, the public sibling of the private
  part. The accepted term is: "**Private part.** Code that serves one owner
  and is never selected by a consumer, such as a lexical helper of one
  language." (`docs/package-boundaries.md:103-104`). A public part serves one
  language too and is never selected. It differs in who may import it: the
  contributions that name its language, among them mappings and
  verification laws that an independent library supplies. So it cannot live
  under `internal/`: "Private helpers are not an alternative public API."
  (`LAYOUT.md:133-134`). It imports foundations, its own language's
  vocabularies on both sides and its own language's lexical helper. It
  imports no contribution and no semantic helper. Its three shapes are the
  paths of P5.
- *Alternative A:* three roles, named reader, constructor and symbol
  admission. It costs three entries for one import rule.
- *Alternative B:* the three stay in vocabulary packages on the operations
  side. It costs the definition, which says "and nothing else".
- *Alternative C:* each joins a contribution of its language: the readers
  the query service, the constructors the target. A mapping may not import a
  target, and an importer of a query service needs an exception that names
  its symbols.

**PROPOSAL P7: a declared table in a vocabulary package.** A grammar table,
a form table and a word list are more than names: they declare what is
admissible.

- *Proposed:* the packages `model/<sort>/<scope>/go` keep the role
  vocabulary, and the definition gets a dated clarification: a vocabulary
  under `model` may hold, beside names, a table over those names that
  declares what is admissible, and it holds no function that evaluates the
  table. The import rule is the same for both, and row R27 asks to map onto
  the words the repository has before one is added.
- *Alternative:* a role of its own for a package that holds a table, named
  declaration. It leaves the definition as it is, to the letter. It costs a
  role whose import rule is the vocabulary's.

### The sentences that change

Each sentence below stays readable as it was accepted, and a dated
clarification stands beside it. The clarifications are made in the change
that adds this page. They state the target and say that the code is not
there. The column "After the slices" is what each sentence says once the
[order of work](#10-the-order-of-work) is done. It uses the proposed names.

| Place | The sentence today | After the slices | Why | Source |
| --- | --- | --- | --- | --- |
| `LAYOUT.md:3-4`, owner decision of 31 October 2026 | "organize by processing stage and enforce a uniform structure across languages." | The root splits into `model` and `ops`. Operations are organized by processing stage with a uniform structure across languages, below `ops`. Models are organized by sort and then by language, below `model` | The direction changes the rule at the root and keeps it below `ops` | The direction; M6 |
| `LAYOUT.md:67`, in the hierarchy | "native semantic forms and validation" | `model/intent/<native>/go` holds the declared table of a language's intended forms. Their validation is the one checker of the kernel | A model contains no operation. The definition of a vocabulary already said "and nothing else" | M4; rows R2 and R7i |
| `LAYOUT.md:145-146` | "Directory groups are not Go packages or an assertion that all members share the same dependency level." | Directory groups below a root are not Go packages or an assertion that all members share the same dependency level. The two roots are such an assertion in one direction: no package under `model` imports a package under `ops` | A root that means "a model imports no operation" is exactly such an assertion | Row R27 |
| `LAYOUT.md:155`, role table, Foundation | "`trees/*/go`, `files/go`, `source/go`, `semantics/go`, `contracts/*/go`" | Two rows. On the model side: the agnostic model packages and the tree carriers. On the operations side: `ops/contracts/*/go` and the operations halves of the carriers | The root cuts the role in two | Row R27 |
| `LAYOUT.md:158`, role table, Vocabulary | "Foundations and its own language's vocabulary. `semantics/<language>/go` and `symbol` may also import their language's lexical helper." | On the model side: the model-side foundations and the model packages of its own scope, and nothing under `ops`. On the operations side: as today. The lexical helper is imported by the public parts and by no vocabulary | Rule 1; every caller of a lexical evaluator is on the operations side | Rows R2, R7h and R27 |
| `LAYOUT.md:170-171`, one language | "A package that names a language imports no package of another language. A mapping names two." | The same, and a pair's contribution vocabulary under `model/contributions` names two as a mapping does | The model root gets packages that name a pair | Row R23 |
| `AGENTS.md:65-66` | "Identities live in vocabulary packages." | Identities live in vocabulary packages: a name on the model side, and an identity that is a value of a contract package's type on the operations side | The vocabulary role exists on both sides | Rows R2 and R27 |
| `AGENTS.md:179-181` | "Native semantic forms belong under transform/syntax/&lt;native&gt;/model/go, their AST interpretation under transform/syntax/&lt;native&gt;/target/go." | A language's intended forms are a declared table under `model/intent/<native>/go`. Their typed constructors are under `ops/transform/syntax/<native>/intent/go`, and their AST interpretation is under `ops/transform/syntax/<native>/target/go` | M4 | M4; P5 |
| `AGENTS.md:233-234` | "Follow the stage-first hierarchy in LAYOUT.md. Shared trees live under trees, shared contracts under contracts" | Follow the two roots of `LAYOUT.md`: models under `model`, by sort and then by language, and operations under `ops`, in the stage-first hierarchy. The generic trees live under `model/trees`, the shared contracts under `ops/contracts` | The direction | The direction; M6 |
| `docs/package-boundaries.md:98-100`, Terms | "**Vocabulary.** The names two parties must agree on and nothing else: syntax kinds, semantic keys, profile and contract identities, intended native forms. A vocabulary package imports no contribution." | The same, and: the role exists under both roots. A vocabulary under `model` holds names and declared tables over them and imports nothing under `ops`. A contract identity is a value of a contract package's type and lives in a vocabulary under `ops` | Contradiction 2 of the evaluation: the definition, the direction and row R8 cannot all hold for one package | Rows R2 and R27; P7 |
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
| S0 | This page, with the clarifications in `LAYOUT.md`, `AGENTS.md`, `docs/package-boundaries.md` and `docs/generation-composition.md`, and the case bittype/ACX-014 | `main` | Nothing | "the page merged after an independent review that checks each quoted rule against its source" |
| S1 | `CheckMember` leaves the four intended-forms packages, and the target states it | The migration's branch, pull request #308 | Nothing | "a test that fails while a check sits in a model package, red on the branch before the move" |
| S2 | The split of the intended-forms packages ([section 5](#5-intent)) | The migration's branch, as its own pull request | S0 | "the five intended-forms packages hold no function; one checker test per table, red by mutant; outputs equal by digest" |
| S3 | The grammar tables and their type on the model side, with predicates named by identity | `main` | S0 | "the tables are data in a model package; the checker test red by mutant of a table entry" |
| S4 | The lexical tables apart from their evaluators | `main` | S0 | "the evaluators import the tables; a test red by a changed table entry" |
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
is in the list of files that the branch changes. So the split on `main`
would be made against code that the migration deletes, and it would be made
again on the branch (the grade, M8). On the branch it is still its own
change, and it still comes before the gate.

`main` keeps the misplacement until the migration merges. That is a
deferral of a move, which
[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
permits only with a record: "Record such a deferral in the case with its
reason and the trigger that ends it; the violation stays open."
[bittype/ACX-014](architecture-cases/0014-models-hold-operations.md) records
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

1. **The half that leaves is created at its final path,** under `model/` or
   `ops/`. The half that stays keeps its old path until S8. Nothing moves
   twice: the script of S8 moves only what is still at an old path.
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
gives 2,121 in 365 files), the paths in the workflow files, the role
function and its tables, and the path prefixes that the documentation tool
reads. No identity of a model or of a contribution is a package path
([section 6](#6-declared-models)), so the move changes none. The numbers
say how large the step is. They are no reason for any choice on this page
(AC-09).

## 11. Decisions and deferrals

**The decisions of the grade.** Each is the research owner's, made on
13 November 2026 and recorded on
[#305](https://github.com/Bitspark/bittype/issues/305#issuecomment-6855741567).

| Decision | What it decides | Rows | Section |
| --- | --- | --- | --- |
| M1 | The test, with the two added sentences | R1 to R4, R7a to R7h, R8 | [2](#2-the-test) |
| M2 | Consistency and correctness, four levels, no new evidence concept, four gaps | R5, R6 | [3](#3-consistency-and-correctness) |
| M3 | Three related models of one language, with a stated relation and a test for each form | R9 to R13 | [4](#4-one-languages-models) |
| M4 | The intended-forms packages: a declared table, one checker, one constructor | R7i, R14, R15, R24, R26 | [5](#5-intent) |
| M5 | Declared models with the mechanisms that exist, in the order of AC-10, identified by a profile | R16 to R20 | [6](#6-declared-models) |
| M6 | The layout: two roots, sort first below `model`, the stage matrix below `ops`, one module | R21, R22 | [7](#7-the-layout) |
| M7 | No prefix move: each split creates its new half at its final path | R28, R31 | [10](#10-the-order-of-work) |
| M8 | The order, with the work in flight | R27, R29 | [8](#8-roles-per-root), [10](#10-the-order-of-work) |
| M9 | Enforcement: one prefix rule by imports, rule 2 in review | R30 | [9](#9-enforcement) |

Two rows belong to no single decision. Row R23 places the identities of a
contribution ([section 7](#every-place-named), slice S7). Row R25 takes one
principle from the advice's four outcomes and no more: a valid value that an
operation cannot process is not invalid. Each operation keeps its accepted
classes, which have a category for a limit and another meaning of
"rejected". The principle is applied where it bites, in the two changes of
[section 5](#5-intent).

**The proposals of this page.** None is decided. The research owner decides
them, and each then gets a number.

| Proposal | The choice | Section |
| --- | --- | --- |
| P1 | The shape of a path below `model` | [7](#names-below-the-roots) |
| P2 | The names of the sorts | [7](#names-below-the-roots) |
| P3 | Where the intent carrier is cut | [5](#5-intent) |
| P4 | The form-table type, as a sketch for slice S2 | [5](#5-intent) |
| P5 | The paths of the operations halves | [7](#names-below-the-roots) |
| P6 | The role "public part" | [8](#8-roles-per-root) |
| P7 | A declared table in a vocabulary package | [8](#8-roles-per-root) |
| P8 | The category of two refusals after the split | [5](#5-intent) |
| P9 | How a grammar table names a predicate | [6](#6-declared-models) |

**Deferred, with triggers.** A trigger is a state of a repository that
anyone can check.

| What is deferred | Why | The trigger that ends it | Source |
| --- | --- | --- | --- |
| A shared identity for single notions of intent and analysed semantics, such as Go's slice | No reader exists | A non-test file imports both a language's intent vocabulary and its semantic vocabulary. The command is in [section 4](#4-one-languages-models) | Row R13a |
| The move of the checks out of the four intended-forms packages on `main` | Every importer is a file that the migration's branch rewrites, so the move is made once, on the branch | The migration of #196 merges. The violation stays open until then, in [bittype/ACX-014](architecture-cases/0014-models-hold-operations.md) | M8; AC-10 |
| A relation from a form to what analysis reports, as data | No package would read it | The first package that must read such a relation without running the test of the form | Row R10 |
| A plain record of a native intent node on the model side | Nothing outside the generation contract would read it | The first package outside the generation contract that must hold native intent without a planner | P3, alternative A, if P3 is taken as proposed |

**What the consultation did not answer.** The grade and the table "Applied"
list these. Each is answered here, proposed here, or open.

| Point | State |
| --- | --- |
| Decision O45, which the advice does not mention | It stands as the research owner's. Silence is no confirmation ([section 4](#4-one-languages-models)) |
| A failure of a comparable system | None was named. Nothing on this page rests on one |
| When the root move itself happens | M8: slice S8, before the migration merges |
| What the split of a sealed Go type is | M1, sentence 1 |
| The six places that the advice's tree does not name | [Section 7](#every-place-named), with P1, P2, P3 and P5 |
| Whether what stays code, the rules between nodes and the predicates, is still model | M1, sentence 2: the model names it, the evaluator states it, and the declared model is an upper bound |
| The role of an operations-side package that is neither contribution nor private part | P6 |
| What `model/relations` is, against `model/contributions` | No directory for relations is created. P2 proposes `families` and `contributions` |
| Whether the payload of an analysed fact is declared as a Go table or as a horos type | **Open.** Decided with slice S9 |
| How a declared model relates to a rule that the building toolchain owns, such as Go's keywords and its universe scope | **Open.** The model names the rule (sentence 2). No profile records the toolchain's release today. Decided with slice S9, which adds the profiles of the native syntax vocabularies |
| Whether adding an intended form is an edit of a shared package in the sense of the gate's line 7 | **Open.** It belongs to the gate of [the contribution contract](contribution-contract.md#the-gate) and is decided there |
| What a family's meaning is for a source that is not horos | **Open**, by decision O46 on #196, item 6. Slice S7 places the family packages and does not answer it |
