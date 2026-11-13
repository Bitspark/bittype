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
