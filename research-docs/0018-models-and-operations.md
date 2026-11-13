# Research: Models and operations: where a language's models live, what they are made of, and whether they are declared

**ID:** 0018
**Date:** 13 November 2026
**Status:** delivering
**Run-ID:** run_36aec98c-ac83-446d-a057-f81b327ed237
**Reviewed:** https://github.com/Bitspark/bittype/issues/305#issuecomment-6854845540
**Owner:** claude-contract-chain-issues
**Issue:** https://github.com/Bitspark/bittype/issues/305

## Question

**In one sentence:** the project's owner has directed that a library of
language tooling splits at its root into models and operations, where a
model contains no operation and imports none and operations work on the
model; by what test is the line between the two drawn, what is one
language's model for every operation and on both sides of a transformation,
how far can such a model be data instead of code, and in what order is the
split carried out?

**bittype** is a Go library. It parses, checks, formats and prints the type
syntax of five languages. It resolves and analyses four of them, three
through their real compilers. It generates native type declarations in
TypeScript, Go, Rust and C++ from canonical type declarations written in a
language called horos. Section 1 describes it. It is one Go module with 129
packages and about 36,000 lines. Nothing is released and no consumer imports
it today. The present layout is two weeks old, and a move of every package
is a matter of days here (section 2.7).

Each language has four packages today that say what exists in it. We call
such a package a **vocabulary package**, and we say **vocabulary** for a set
of names. Two of the four, the syntax vocabulary package and the semantic
vocabulary package, sit beside the processing stages (a stage is one
processing step: parse, check, format, print, resolve, analyze, generate).
Two of them sit inside the directory of one stage, generation: the package
of the type forms that generation intends to build, and the package of
external symbols (the names that a generated declaration refers to and does
not itself declare, such as a predeclared type). Each of the four states the
type-level notions of its language in its own words. Section 3.1 shows the
four for Go.

**The word "model".** From here on "a model", unqualified, has the owner's
sense: what belongs under the model root. What that is, is asked. The
existing package of intended forms is always called the intended-forms
package.
The code calls that package `model`. Whether it is a model in the owner's
sense is part of what is asked. The code has a second thing named `Model`, a
reader type in the semantic vocabulary package, and we gloss it where it
appears.

### The framing: the owner's ten messages

The owner wrote these messages on 13 November 2026. The team works in
separate sessions that do not share a conversation, and the owner spoke to
one of them. That session relayed the messages to the research owner. The
research owner is the team member who is responsible for this line of work
and who writes this document; the header names it. The research owner did
not hear the messages directly. They are quoted verbatim and in order. They
are the framing of this consultation.

The owner and the code say "transform" for the operation that this document
calls generation: producing one language's type declarations from another
language's analysed declarations.

1. > why can't we treat all languages the same, and transform can be from any source to any target? whether the transformation are implemented and exist is a different story, practically it would still be the way you describe, but not because it is baked in, but because it follows from the registered transforms' declarations

   (By "the way you describe" the owner refers to the session's account of
   generation as going from horos to the four programming languages.)

2. > Tell the fork we should add maybe a few other language-to-language transforms, not because we want to fully support this now, but just to make sure we don't silently assume this

   ("The fork" is another working session of the team. The owner addressed
   the session that relayed the messages, and meant the session that does
   this work.)

3. > and this model for target/source language - where does it live?

4. > but why is this in transform? isn't that all just syntax-related? so the model, is it a syntax model, or a transform model?

5. > like, isn't 'transform' just one of many operations working on this model?

6. > But shouldn't we split the models from the operations? The model is independent of the operation, isn't it?

7. A layout and a second question:

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

   The block keeps the line breaks of the message as it was relayed.

8. The direction:

   ```text
   So the root split to me would be between model and ops.
   model must not contain any ops and must not import from ops either.
   ops work on the model
   ```

9. > ok, so let's look at ops then. you are talking about the correctness of a model, aren't you?

10. > or consistency?

**Messages 1 and 2** have been acted on by the research owner, in two
decisions on the work that replaces what a mapping rule returns (section
2.7). They are quoted because the later messages follow from them: "this
model for target/source language" in message 3 is what a language would
need on either side of a transformation.

**Messages 7 and 8 are the owner's direction.** We read them so, and the
reading of what is direction and what is a sketch is ours:

- *Direction.* The root of the repository splits into models and
  operations. **Rule 1 (no import):** a model imports no operation. **Rule 2
  (no containment):** a model contains no operation. Operations work on the
  model.
- *A sketch.* The names below the two roots in message 7,
  `model/{syntax,semantics,...}/go/...` and `ops/{transform,...}/go/...`.
  You may propose another order below a root, language first included. The
  trailing `go` is the segment of today's paths that names the
  implementation language.

**Messages 9 and 10 are questions of the owner.** They are no direction.
They answer our first question as another session had passed it on to the
owner: where a model's own well-formedness ends and an operation begins.
What we made of them is our derivation: a difference between consistency
and correctness, and four levels (section 6.1). Question 1 asks whether that
is the right axis.

**What is asked and what is not.** The split and its two rules are not
asked. What we want back is a test that we can apply to a package ourselves,
verdicts on named cases, patterns from comparable systems with what failed,
and an order of work. We do not want a review of the code shown.

**What is in flight** (section 2.7). A migration of the mapping contract is
being built on a branch. Two further pairs of languages are to be built on
that branch before it merges. Two Files formatters, for Go and for Rust, are
the next steps of two other work items. All of it is written for today's
paths.

**Which sections serve which question.**

| Question | Sections |
| --- | --- |
| 1. The line between model and operation | 6.1 and 6.2; the carriers in 2.3 and the contract packages in 2.5; code in 4.2 and 4.4 |
| 2. What one language's model is | Section 3; 6.3 |
| 3. Declared models | Section 4; 6.4 |
| 4. Carrying the split out | 2.4, 2.6 and 2.7; 6.5 and 6.6 |
| 5. What we are not seeing | The whole document |

A reader short of time can skip section 4 and return to it from section 6.

### Four questions, as a working frame

The research owner separated the messages into four questions. The
separation is ours and is a working frame.

| Working question | What it asks | Where it is asked |
| --- | --- | --- |
| Grouping | Does a package hold a model and an operation together? Rule 2 makes this general | Question 1 |
| Package placement | What is on the model side: only a language's vocabulary packages, or also the language-independent data types for trees, graphs and files (the carriers, section 2.3) and the input and result types declared beside each operation's interface (section 2.5)? And how is each root organized below? | Questions 1 and 4 |
| Unification | How do a language's syntax, what analysis reports and what generation intends relate, and should two of them share their names and differ only in status, intended or established? | Question 2 |
| Declaration | How far can a language's model be data, and what stays code? | Question 3 |

One small step is planned whatever the answer: section 6.5. Please say there
if it prejudges something.

### How to read the standings

Every statement in this document has exactly one of seven standings. Each
section says which standing its statements have, and section 5 lists the
fixed points with their standing.

| Standing | What it means | Weight in the design |
| --- | --- | --- |
| **Owner rule** | A rule or a decision of the project's owner. Quoted, with its date. Where the owner's own words are not on file, the rule is given in the wording of the page that records it, and we say so | Fixed |
| **Accepted design** | A design page merged on the main branch, or a row of an earlier expert consultation that was accepted. We name the page or the consultation | Fixed unless you see that it contradicts an owner rule, the owner's direction, another accepted design or itself. If so, please say so |
| **Implemented today** | What the code on the main branch does. Under the owner rule AC-09 the existing code has no weight in a design choice (the rules are numbered AC-01 to AC-11, section 2.1). It is given as a fact and never as a reason | None |
| **Our derivation** | A reading or a conclusion of the research owner or of the two studies behind this document. It is in no accepted source | Open. Please challenge it |
| **The owner's direction** | Messages 7 and 8 | Fixed: that the root splits into models and operations, and the two rules. Asked: how the line is drawn and in what order |
| **Decision of the research owner** | Made by the team member who writes this document, recorded on the work item, not reviewed by the owner | Open: say so if it is wrong |
| **Open** | Not decided anywhere we looked | Open |

"Decided" is said only of an owner rule, the owner's direction or an
accepted design. "Our reading" means our derivation, and "as facts" means
implemented today.

**Earlier consultations.** "Research N" is an earlier expert consultation of
this project. Its result is a table of rows, one per recommendation. Each
row has a verdict that the research owner of that consultation gave: HOLDS
(accepted as advised), NEEDS ADAPTATION (accepted with a stated change) or
DEFERRED. An accepted row is an accepted design in the row's own words.

Two consequences, stated once. First, the paths, package names and counts in
sections 2 to 4 describe the present code. None of them is a constraint on
the answer. Second, where pages contradict each other or the code, we say
so. Section 5.3 lists such places, and the answer should not rest on either
side of one.

**The two studies.** Two read-only studies were made for this document on 13
November 2026: one of the code, and one of the design pages and the
history. Numbers come from them, and we recounted the totals.

**Quotations.** Quotations of code give the file and the line at the code
baseline, which is bittype's main branch at commit a28febc
(a28febc44830d49ec3c2f22cd70c7ded9d0e5028) of 13 November 2026. All code is
Go unless marked. Quotations of documents give the file. Where a quoted
sentence held a Markdown link, the quotation keeps the link's text and drops
its address.

### Terms you need first

The glossary is section 8. These are needed from the first page.

| Term | In one line |
| --- | --- |
| Stage | One processing step: parse, check, format, print, resolve, analyze, generate |
| Generation, transform | Producing one language's type declarations from another language's analysed declarations. The owner and the code say "transform" |
| Kernel | The shared code that selects and composes. It names no language |
| Contribution | One thing a consumer registers or selects on its own: one language's parser, one language's generation target, one mapping |
| Carrier | A language-independent data type for one sort of artifact: a tree, a graph of analysed facts, a set of files |
| Contract package | The package of one operation's interface. It also declares the values that pass through the operation |
| Vocabulary package | A package that says what exists in one language. A vocabulary is a set of names |
| Target | A language that generation produces, and the contribution that builds its syntax trees |
| Mapping | A contribution that says, for one sort of source type and one pair of languages, what the target type is to be |
| Form, intent | A form is one sort of intended type construct of a target language, such as Go's `slice`. Intent is a tree of forms: a description of a type before any syntax exists |
| Intended-forms package | The package that issues one target language's forms. The code calls it `model` |
| Symbol | What a name in a generated declaration refers to when the referent is not itself generated |
| Native | The pages' word for the four programming languages as opposed to horos, and so also for the target side of generation |

## 1. System overview

Standing of this section: implemented today, except the two rules of 1.5,
which are owner rules.

### 1.1 The components

- **ontos** is a data model with a canonical binary encoding: a value is an
  atom (a byte string) or a tuple of values.
- **horos** is a declarative type language over ontos values: a horos type
  denotes a set of ontos values, declarations live in modules, and a module is
  identified by a hash of its content, its **pin**.
- **bittype** is a Go library, whatever language it processes. It has three
  groups of capabilities.
  - **Syntax.** For each of five languages (horos, TypeScript, Go, Rust and
    C++) it parses type syntax from text into a tree, checks a tree against
    the language's grammar, formats a tree back into concrete syntax and
    prints it. It covers type declarations and type expressions only. It has
    no statements and no executable bodies.
  - **Meaning.** For four of them (horos, TypeScript, Go and Rust) it resolves
    names and analyses declarations. It does not reimplement the languages'
    type checkers. Go is analysed by the type checker of the Go standard
    library. TypeScript is analysed by a helper script that runs the
    TypeScript compiler. Rust is analysed by a pinned helper executable built
    on the Rust compiler's driver. horos is analysed by the horos library.
    Each result is published as a graph of facts in one
    language-independent graph type. We call the Go code that runs a
    compiler and exports its facts a **bridge**.
  - **Generation.** From analysed horos declarations it generates native type
    declarations in TypeScript, Go, Rust and C++. The pages say "native" for
    the four programming languages as opposed to horos. Every rule that
    turns one sort of horos type into a native type is supplied from outside
    the shared mechanism, by a package that a consumer registers. For
    TypeScript a later step places the generated declarations into files and
    writes the imports. For Go and Rust that step is planned and not built.

  These three groups are not the two processing dimensions of section 2.4.
- A fourth, optional component relates a captured horos program to a captured
  Go program and checks stated representation laws through logos, a separate
  reasoning library of the organization. It appears in this document only as
  one more reader of a language's analysed meaning.
- **Consumers** are to be builds of the organization's own services and
  libraries. None imports bittype today (section 1.4).

The example that the rest of the document reuses is one horos record with one
text field:

```horos
module example.people;
type Person = { name: text; };
```

In horos, `module` names a module and `type` declares a named type.
`{ name: text; }` is a closed record with one required field, and a field
written `next?: text` would be optional. `text` is the built-in scalar type
of UTF-8 text, which the code spells `scalar` with `utf8-text`. `list<T>` is
a list of `T`.

For Go, generation produces this declaration (read from the code and followed
by hand in section 3.2; not run for this document):

```go
type Person struct {
	Name string
}
```

### 1.2 The shared mechanism and what is supplied to it

Two words recur. The **kernel** is the shared code that selects and composes:
an engine, a generation registry and planner, and the interfaces they are
written against. A **contribution** is one thing a consumer registers or
selects on its own: the parser of one language, the checker of one language,
the generation target of one language, one mapping rule. The kernel names no
language. A composition, which today is test code and a few tools, chooses
contributions and registers them.

The project's rules rest on this division. Section 2.1
gives the rules behind it, and section 2.6 gives the roles that the
repository's own architecture test enforces.

### 1.3 Size

bittype is one Go module with one `go.mod`. At the code baseline it has 129
production packages in 245 Go files with 35,943 lines, not counting tests.
It has 126 test files with 31,259 lines. The table gives the production
packages by top-level directory. We recounted it at the code baseline with a
script over the committed files.

| Directory | Packages | Lines | What it holds |
| --- | --- | --- | --- |
| `trees/` | 3 | 650 | The generic syntax trees and their type-erased envelope |
| `files/`, `source/` | 2 | 699 | The generic carriers for supplied files and for placed, analysed source |
| `semantics/` | 5 | 2,645 | The generic semantic graph (1 package, 560 lines) and four languages' semantic vocabularies |
| `syntax/` | 5 | 644 | Five languages' syntax vocabularies: constants only |
| `contracts/` | 13 | 5,030 | One package per operation: its interface and the value types it exchanges |
| `engine/` | 1 | 860 | Registration, route planning and execution |
| `parse/`, `check/`, `format/`, `print/` | 22 | 7,049 | The Syntax stages of each language, the structure checks of the generic trees, the shared printers, and the stages of horos and TypeScript in the Files dimension |
| `resolve/`, `analyze/` | 10 | 716 | The stage adapters around the bridges, and two query services |
| `transform/` | 39 | 3,584 | Generation: the shared kernel parts, one source reader, four targets, four intended-forms packages, three symbol packages, 21 mapping rules and four shared mapping parts |
| `verification/` | 7 | 1,676 | The optional verification component |
| `internal/` | 16 | 6,286 | Private helpers: lexing and parsing mechanics, the grammar table type, per-language spelling rules, the three compiler bridges, the horos adapters |
| `tools/` | 6 | 6,104 | A documentation generator that reads the source, and one horos command |
| **Total** | **129** | **35,943** | |

The terms of the last column are defined in sections 2.2 to 2.6. The table
is here for proportions: the five syntax vocabulary packages together have
644 lines, and the four intended-forms packages together have 467.

### 1.4 Who uses it

No consumer imports bittype today. Nothing is released, and the repository
has no version tag. Packages outside `internal/` have public import paths,
and the intended consumers are repositories of the same organization. The
study of the code found no other Go module that requires bittype in the
checkouts on one workstation.

The work is done by a small number of automated agent sessions that work in
parallel under the owner's rules. The present layout dates from 31 October
2026 and the directory of generation from 7 November 2026. A move of the
paths of the whole repository is a matter of days here.

### 1.5 Two rules of the owner

Both are rules of the project's owner. Both apply to every project of the
organization, and section 2.1 gives them in full.

- **Architecture is part of correctness** (owner rule of 6 November
  2026). A feature that works and passes its tests can still be incorrect,
  because a decision sits with the wrong owner. The rule AC-01 requires:
  "Satisfy the agreed ownership, dependency and composition contracts as well
  as observable behavior."
- **An interface is derived from the domain and never from the code that
  exists** (owner rule of 8 November 2026, AC-09). "The currently
  implemented surface must have **zero influence** on that choice.
  Implementation or migration effort has **zero weight**, including as a
  tiebreaker."

Our reading of the second rule for this document: that a new layout touches
every package is no objection to it.

## 2. Architectural context

Each subsection states the standing of what it holds.

### 2.1 Standing architecture rules (AC = architectural correctness)

Standing: owner rules, except where a paragraph says otherwise. They are
written down in two pages of the organization's shared documentation
repository, called the platform. They apply to every project. bittype follows
them and keeps a numbered record of each violation it finds, a **case**
(bittype's cases are named ACX-001 to ACX-010 at the code baseline).

There are eleven numbered rules. The table gives each requirement in the
rule's own title.

| Rule | Requirement |
| --- | --- |
| AC-01 | Architecture is part of correctness. |
| AC-02 | Meanings remain with their semantic owners. |
| AC-03 | Extensions compose additively. |
| AC-04 | Contributions use the same public mechanism. |
| AC-05 | Composition preserves guarantees. |
| AC-06 | Limited coverage does not redefine the architecture. |
| AC-07 | Evidence supports only the claims it establishes. |
| AC-08 | General contract changes require a general justification. |
| AC-09 | Design interfaces independently of the implemented surface and the effort to change it. |
| AC-10 | Correct a misplaced responsibility first; add missing coverage afterwards. |
| AC-11 | Agnostic mechanisms belong in the agnostic layer. |

Six passages bear on this document.

**Semantic ownership through additive composition** (owner rule of 6 November 2026):

> A reusable mechanism owns the invariant rules of its operation. Domain-specific
> meanings and policies belong to independently supplied implementations. Within
> the declared extension boundary, adding a new case must require changes only
> to its owning implementation and the consumer's composition, not to shared
> orchestration or existing unrelated implementations.

The same page says what does not settle ownership: "Semantic ownership is
about who defines a meaning and how that meaning enters an operation. It is
not determined by a directory, package name, generic type parameter or the
word "plugin"."

**AC-09, no weight for what exists** (owner rule of 8 November 2026).
Section 1.5 quotes its core. Its application clause says what the existing
code may be used for: "First settle the intended contract. Then use the
existing implementation to plan migration, estimate the work and decide which
code fulfills that contract. Those estimates can guide delivery sequencing
only."

**AC-10, a misplaced responsibility moves first** (owner rule of 11
November 2026). The owner's words:

> if something is placed incorrectly, it is incorrect and must change. Something that is missing can be healed additively, i.e. there is nothing wrong requiring removal.

> no compromise on architecture, but practical deferral is allowed, if the huge effort justifies it.

The requirement derived from them: "Separate a correction into what is in the
wrong place and what is only absent. Move a misplaced or wrongly named
responsibility to its owner first, as its own change, without waiting for the
complete implementation of its destination. That first change establishes the
final structure: its inputs and identities have their intended shape, it
keeps every guarantee the destination's shared mechanism requires for the
cases it admits, and it explicitly refuses the cases it cannot yet check.
Missing coverage is then added without moving, renaming or removing what
exists." Its assessment asks: "After the first change, does any planned step
have to move, rename or remove what that change delivered?" A deferral of the
move itself "is permitted only when its destination needs a huge effort", and
then the violation stays open and recorded.

**AC-11, agnostic mechanisms** (owner rule of 13 November 2026). "A
language-agnostic mechanism belongs in the agnostic layer. Each language
supplies its own types, meanings and policies through the shared contract."
Its requirement: "Place a mechanism according to the knowledge its contract,
algorithm and guarantees require. When these are independent of a language or
domain, implement the mechanism once in the corresponding agnostic layer.
Supply the variation through typed parameters, data or independently owned
implementations." The rule is of the same day as this document. Three cases
are recorded under it so far (ACX-010 to ACX-012, the last two after the
code baseline). None applies it to the question asked here.

**Grouping before coupling, and what a dependency check is** (owner request
of 11 November 2026, with the platform's reading of it). The owner asked to
check for packages "that contain multiple implementations / functionalities
that either already are separate and neither conceptually nor practically
depend on each other, or which still depend on each other but shouldn't", and
to enforce the result "through forbidden dependencies that we could check for
and enforce". The platform page separates two faults:

| | Wrong grouping | Wrong coupling |
| --- | --- | --- |
| What it is | Parts that do not need each other share one package or module | A contribution names another contribution that the consumer is meant to choose |
| What is wrong | The boundary. No dependency between the parts is wrong | The dependency. The caller makes the shared mechanism's choice |
| Consequence | An importer of one part receives all of them, and a dependency check cannot tell them apart | The mechanism's selection, scope and failure guarantees do not apply to that call. Replacing the callee means editing the caller |
| Correction | A move | A demand that the shared mechanism resolves |

"The first hides the second. A dependency check reads a unit, usually a
package. Whatever is inside one unit can reach everything else in it unseen.
Correct the grouping first, as AC-10 orders, and the coupling becomes
visible." The page then says when a forbidden-dependency check can enforce a
decided boundary, namely "when each unit has one role", and gives the roles
that section 2.6 shows for bittype. It ends with a limit that question 4 of
this document rests on: "Such a check enforces a boundary that reasoning has
already decided. It does not contradict the paragraph on lint gates above,
and it does not replace that reasoning: it cannot see a decision made inline,
it cannot judge whether a unit holds one contribution, and a green run shows
only that the forbidden edges are absent."

**Deriving before extending, and unification.** Both come from the platform's
method page on semantic ownership and generalization. The page is recorded in
an owner decision of 6 October 2026. On 7 October 2026 the owner required
designers to derive a case from existing concepts before proposing an
extension. The page states that requirement so: "An extension needs a
concrete case that the existing contracts cannot express or uphold, together
with an explanation of why composition does not resolve it. Distinguish that
conceptual gap from missing code, an unspecified encoding or an unimplemented
connection between existing contracts."

On unification the page says:

> Unification claims that apparently separate concepts share one semantic
> contract. It therefore needs stronger evidence than overlapping fields or
> similar code.
>
> Compare identity, lifetime, authority, allowed operations, guarantees and
> failure behavior. A consumer using the proposed common contract must retain
> the same guarantees when one valid specialization is substituted for another.
> Differences must be represented explicitly rather than erased.

Its table of findings and responses:

| Finding | Appropriate response |
| --- | --- |
| Different names, same contract and compatible ownership | Unify terminology and the contract; plan the implementation change explicitly. |
| Shared laws with meaningful domain variations | Define a common contract with typed specializations or composition. |
| Similar mechanics but different contracts | Share implementation where useful while preserving the semantic distinction. |
| Independent identities or lifetimes | Keep separate concepts and define their relationship. |
| Overlapping roles of the same data | Model the roles or relations without forcing a single inheritance hierarchy. |
| Evidence does not settle the equivalence | Record the competing interpretations and the case needed to decide. |

We treat the derivation requirement as an owner rule. We treat the
unification test as an accepted principle of the platform and not as the
owner's own words: the page records the owner's request, and it does not
quote the owner on unification.

### 2.2 Generation

Standing: accepted design (the pages `docs/generation.md`, "Current
architecture, 8 November 2026", and `docs/kernel-composition.md`,
"Implemented, 8 November 2026") and implemented today.

```text
analysed source program (semantic graph + receipt)
   |  the kernel walks the graph; the source reader interprets it
   v
for each use of a type: one mapping is selected and states an intended form
   |  the kernel composes the forms, substitutes parameters, resolves references
   v
intent: a closed program of intended declarations      <- checked by the target
   |  the target builds one tree per program, through the shared tree builder
   v
target syntax trees, sealed with their declared symbols <- checked by the language's checker
   |  the Files formatter places declarations under a consumer's placement
   v
placed source, analysed again by the language's bridge  <- destinations compared
```

The participants and their words:

- **Source reader.** A contribution of the source language. It tells the
  kernel which entities are declarations and gives mappings a typed reading
  of a use, for example the fields of a record. Only horos has one.
- **Mapping.** A contribution for one ordered pair of languages and one sort
  of source type, for example horos records into Go. It states what the
  target type is to be. It uses the constructors of the target's
  intended-forms package and never builds a tree. There are 21, each in its
  own package of 25 to 78 lines.
- **Family.** One sort of source type: text, booleans, lists, closed
  records, optional fields, named declarations. A mapping is one family for
  one ordered pair of languages.
- **Form.** One sort of intended type construct of a target language,
  identified by a name that the language's intended-forms package issues,
  such as Go's `slice`, `record` and `field`. A constructor is the Go
  function that builds a piece of intent of one form.
- **Intent.** A tree of forms with string data and child expressions, with
  placeholders that the kernel fills: a slot for the representation of
  another use, a reference to a type parameter, a reference to another
  planned declaration. It describes types only.
- **Target.** A contribution of the target language. It checks intent, and it
  builds the syntax tree from it. A case record of the repository (ACX-005)
  says: "The native target remains the trusted interpreter of its own
  forms."
- **Output.** The trees, sealed together with an inventory of their
  declarations and the list of external symbols they use.
- **Files formatter.** A contribution of the target language in the Files
  dimension. It decides files, imports, qualification and visibility under a
  placement that the consumer supplies. It then has the placed program
  analysed and compares each name use with the symbol or declaration it was
  meant to reach. The code and the pages call this contribution a Files
  formatter. It exists for TypeScript. For Go and Rust it is planned.

The accepted sentence on intent, from `docs/generation.md`:

> Build receives no source reader or mapping callback. Native intent is not an
> analyzed `Semantics[T]`: target analysis establishes facts about the actual
> constructed/organized program.

The same page lists what each component may know. Its row for the "Native
generation model/target" reads: "Its native forms and syntax; no
Horos/library semantics". And `docs/kernel-composition.md`: "The native model
describes intended types; the target constructs their syntax. Analyzed native
facts remain under semantics."

The pages name the intended-forms package in five ways, among them "native
semantic forms and validation" and "native intent vocabulary". None calls it
a syntax model or a transform model, the two words of the owner's message 4.

### 2.3 The generic carriers

Standing: implemented today. Sections 4.1 and 4.3 inline the types.

A **carrier** is a generic data type that holds one sort of artifact for
every language. It knows no language. Three words recur in the table. A
**profile** is a named promise, identified by a string: of a language, which
revision of its syntax; of an analysis, which compiler version and which
fragment of the language; of generation, what the generated types retain. A
**provider** is the language-specific code that writes and reads a graph.
**Captured** means read once into an immutable snapshot. bittype has six carrier packages.

| Carrier | What it holds | Checked by |
| --- | --- | --- |
| **AST** (abstract syntax tree) | One tree as a flat list of nodes. A node is a kind, a list of attributes (a name and a text), a list of children (a role and a node index), a list of references (a role and the address of a node in some tree) and an origin (source spans and contributing nodes) | A structure check in another package: valid indices, no cycle, no node with two parents, bounds on size and depth |
| **CST** (concrete syntax tree) | One concrete tree. A node is a leaf with text or a container with children. The CST owns an AST and links concrete nodes to abstract nodes | A structure check in another package |
| **Envelope** | A type-erased, structurally checked, immutable snapshot of an AST, a CST or a semantic graph, and the descriptor of a language (its identifier, display name, a profile string and its Go marker type). The kernel carries envelopes because it names no language | The envelope calls the two structure checks when it packs a tree |
| **Semantic graph** | The analysed meaning of one captured program: entities, terms, bindings of occurrences to entities, expressions at occurrences, dependencies and facts. A term has one of three fixed variants: a reference to an entity, an application, or a **provider expression** (an operator, an encoding, bytes and operands, all owned by the language's provider). Every other name in the graph is a key, a pair of an owner string and a name | Its own constructor: graph integrity and bounds only. "Applicability, arity and saturation belong to the provider, not graph admission" (comment on the term type) |
| **Files** | Supplied text files with their identifiers and paths. No file system access | Its own constructor |
| **Source** | Placed concrete files together with the semantics established for them: units, groups, relations, and the receipt of the analysis | A publisher that runs under a policy with two callbacks, one for printing and one for a language's own check |

**How a language attaches to the carriers.** Three mechanisms exist.

1. *A marker type as a type argument.* Each language declares one empty Go
   struct type. The carriers are generic in it: the AST of Go is the generic
   AST instantiated with Go's marker. The Go compiler therefore refuses to
   pass a tree of one language where another's is expected. Node kinds are
   string constants of a type that carries the marker. Role names and
   attribute names are plain string constants.
2. *Keys with an owner.* In the semantic graph a language's names are keys
   whose owner is a string such as `go` or `rust`. Nothing registers a key.
   The graph stores it, and the language's own code compares it.
3. *Issued values.* In generation a language's intended forms and external
   symbols are values created by an **issuer**, an object that the language's
   package holds privately. Only the issuer can read a value it issued. The
   identity is the issuer's pointer, and the content is a string.

So for syntax the carrier is the whole data structure, and a language adds
names. No language declares a struct type for a node. The shape that a tree
of a language must have is neither in the carrier nor in the names. It is in
the language's checker (section 4.2).

For the semantic graph a language adds names and also Go types: handles that
read the graph in the language's own terms, for example the fields of a
struct type.

**Two carriers hold or call code of an operation.** The envelope package
calls the two structure checks, which live in a check directory
(`trees/representation/go/ast.go:63`, `cst.go:23`):

```go
	if err := check.Structure(input); err != nil {
...
	if err := cstcheck.Check(input); err != nil {
```

And the Source carrier admits a publication under a policy that holds two
callbacks into operations, a printer and a language's own check
(`source/go/source.go:121-130`, `:254`; `semantic` is the semantic contract
package of section 2.5):

```go
// Policy is an explicit language-owned admission contract. Check establishes
// CST/AST agreement, grouping and the operation's preservation obligations.
// Print supplies the configured Syntax printer; neither callback is discovered.
type Policy[L any] struct {
	Evidence       semantic.EvidenceRequirement
	SemanticLimits semantic.Limits
	Limits         Limits
	Print          func(context.Context, cst.CST[L]) (string, error)
	Check          func(context.Context, Publication[L]) error
}
...
func (p Publisher[L]) Publish(input Publication[L]) (Source[L], error)
```

### 2.4 Stages, the two dimensions and the path pattern

Standing: owner rules, in the wording of the pages that record them, and what
is implemented under them.

**Two dimensions.** The repository's contributor instructions attribute this
arrangement to the owner: "The owner's 4 November Syntax and Files symmetry".
The page `docs/source.md` describes it, in two parallel rows:

```text
Syntax : CST    : AST
Files  : Source : Semantics
```

| Aspect | Syntax | Files |
| --- | --- | --- |
| External form | Text | File set |
| Concrete representation | CST | Source |
| Abstract representation | AST | Semantics |
| Parse | Text → CST | Files → Source |
| Format | AST → CST | Semantics → Source |
| Print | CST → Text | Source → Files |
| Abstract view | CST exposes AST | Source exposes Semantics |

The page adds: "Files and Syntax are processing dimensions, not language
IDs." Resolution and analysis "operate on the semantic side of the Files
dimension". We call the last stage of generation, which places generated
declarations into files, the **Files formatter** of a target language.

**The stages** are parse, check, format and print in the Syntax dimension,
parse, format and print in the Files dimension, resolve and analyze, and
transform. The directory `transform/` holds generation. A separate directory
holds verification.

**The path pattern.** Two owner rules fix where an implementation lives. No
file that we read quotes the owner's own words of either day. The layout page `LAYOUT.md`
records them so:

> **Owner decision, 31 October 2026:** organize by processing stage and enforce a
> uniform structure across languages. The stage-by-language matrix is deliberate:
> the same responsibilities occupy the same paths, making omissions and drift
> visible.

> **Owner decision, 6 November 2026:** the accepted
> kernel-owned composition and uniform layout
> uses `<stage>/<dimension>/<scope>/<component>/go`, with generic factories in
> kernel/contracts and no per-pair assembly or presets. This hierarchy is implemented.

The same page explains the form: "Implementations use the same stage,
dimension, scope, component and implementation-language order. A scope
identifies one language, an ordered pair or shared machinery. Go packages use
their role name. Syntax vocabularies retain their language names." The final
`go` names the implementation language. The page also says: "Private helpers
are not an alternative public API. Do not create a language-first packages
tree, additional Go modules, compatibility aliases or placeholder source
directories."

Examples of the form at the code baseline:

| Path | What it is |
| --- | --- |
| `parse/syntax/rust/provider/go` | The Rust parser |
| `check/syntax/golang/provider/go` | The Go checker, which holds Go's grammar table |
| `check/syntax/shared/ast/go` | The structure check of the generic AST |
| `format/files/typescript/provider/go` | TypeScript's Files formatter |
| `analyze/files/golang/queries/go` | A query service over analysed Go |
| `transform/syntax/golang/target/go` | The Go generation target |
| `transform/syntax/golang/model/go` | Go's intended forms |
| `transform/syntax/horos-to-golang/mapping/record/go` | The mapping rule for horos records into Go |

**What is outside the form.** The carriers (`trees/*/go`, `files/go`,
`source/go`, `semantics/go`), the contract packages (`contracts/*/go`), the
engine, the private helpers under `internal/`, the tools, and two sorts of
vocabulary: `syntax/<language>/go` and `semantics/<language>/go`. The two
other sorts of vocabulary package, `model` and `symbol`, are inside the stage
directory `transform/`.

The four intended-forms packages were created inside `transform/` on 8
November 2026, by the change that separated the planning of intended types
from the construction of trees (before, a mapping built syntax trees
directly). No page gives a reason for that place against a place beside
`syntax/` and `semantics/`, and no consultation was asked.

### 2.5 The contract packages

Standing: implemented today. The column "value types" against "protocol" is
our derivation.

A **contract** here is the interface between one operation's shared mechanism
and its independently owned contributions. Each operation has one contract
package. It declares the Go interface that a contribution implements. It also
declares the values that pass through the operation: inputs, limits,
diagnostics, receipts. And in several cases it holds the code that admits a
request and seals a result, because only that code may create a value that
later stages trust.

By value types we mean plain data that crosses the interface. By protocol
we mean the interfaces and the functions that enforce who may create a value
and in what order. Four verbs recur. To **admit** a value is to validate it
and wrap it in a type that only the validating package can construct. To
**seal** or **freeze** it is to make it immutable after validation. To
**publish** is to hand out the sealed result. To **issue** is to create an
opaque value that only its creator can read.

| Contract package | Lines | Value types it declares | Protocol and admission code it holds |
| --- | --- | --- | --- |
| `contracts/language/go` | 25 | The identifier of a language and the two-method interface that returns it | None |
| `contracts/text/go` | 45 | A diagnostic: a code, a message, a span, a node | One constructor of a diagnostic at a node |
| `contracts/pipeline/go` | 18 | None | One function type. No production package imports it |
| `contracts/check/go` | 72 | None | The checker interface and its type-erasing adapter |
| `contracts/parse/go` | 258 | The parse input, the Files input and limits | The parser interfaces for Syntax and Files, and the request through which a Files parser publishes |
| `contracts/format/go` | 658 | The Files input, the placement of generated declarations, the placed declaration | The formatter interfaces, and the admission of a placement |
| `contracts/print/go` | 151 | None | The printer interfaces |
| `contracts/transform/go` | 151 | The transformation input and its contract description | The transformer interface and the admission of a request |
| `contracts/semantic/go` | 1,047 | Obligation, requirement, verdict, report, diagnostic, setting, limits, evidence requirement, and two handles: the **snapshot** (a captured, admitted input) and the **receipt** (the result of one semantic operation over a snapshot) | Capture, the operation runner, publication and the admission of evidence |
| `contracts/resolution/go` | 93 | The resolution contract description | The resolver interface |
| `contracts/analysis/go` | 83 | The analysis contract description | The analyzer interface |
| `contracts/generation/go` | 2,135 | The data types of intent (forms, expressions, declarations: about 140 lines), the issuers of symbols and values, the description types of a mapping, errors and limits. Together about 470 lines | The target and mapping interfaces, the planner, the freezing of descriptions, the sealing of an output, the walk over the source graph. About 1,650 lines |
| `contracts/verification/go` | 294 | Participant, requirement, site, limits and the association handle | Capture and the check of sites |

Three facts follow, and all three are implemented today.

- No contract package imports a vocabulary of any language.
- Every vocabulary of a language imports at least one contract package. The
  syntax vocabulary imports `contracts/language/go` to state its identifier.
  The semantic vocabulary imports the semantic, resolution and analysis
  contracts. The intended-forms package and the symbol package import the generation contract package,
  because every value they issue is a value of that contract.
- The data types of intent are generic like a carrier, but they are no
  carrier package. Forms and expressions are declared inside the generation
  contract package, beside the planner.

**The trusted handles, shown for the semantic contract package.** A snapshot
and a receipt are types with private state. Only functions of the contract
package create them, after they have validated what goes in
(`contracts/semantic/go/operation.go:26-28`, `:93-94`,
`request.go:91`, `:106-110`, signatures only):

```go
// Snapshot identifies one admitted input collection. Identical addresses or
// bytes captured again are a different snapshot. Storage is private and immutable.
type Snapshot struct{ state *snapshotState }

// Capture bounds and retains immutable envelopes before any checker callback.
func Capture(ctx context.Context, checker check.Checker, input Input) (Snapshot, error)

func (r Request[L]) Publish(publication Publication) (Result, error)

// Result is an immutable request-bound publication. A valid result can contain
// diagnostics and reports without a model. Valid does not mean semantic validity.
type Result struct{ state *resultState }
```

`Capture` runs the language's checker over the input trees and returns the
snapshot. `Publish` is called by a resolver or an analyzer at the end of its
work and returns the receipt, which the code names `Result`. A package that
names a snapshot or a receipt in its own types imports this contract
package. That the handles are trustworthy because only this code creates
them is how the code works today. We know no accepted page that requires
this pattern, and you may replace it.

### 2.6 Roles, and the test that derives a role from a path

Standing: accepted design. The page `docs/package-boundaries.md` states
"Implemented. State of 12 November 2026." and `LAYOUT.md` holds the role
table. What the test does is implemented today.

Every production package has one **role**. A role is an import rule. The
page defines five roles in its section "Terms":

> - **Kernel.** The shared mechanisms that select and compose: the engine, the
>   generation registry and planner, and the contracts they are written against.
> - **Contribution.** One thing a consumer registers or selects on its own: a
>   stage provider, a source, a target, one mapping library, one verification
>   law.
> - **Vocabulary.** The names two parties must agree on and nothing else: syntax
>   kinds, semantic keys, profile and contract identities, intended native forms.
>   A vocabulary package imports no contribution.
> - **Composition.** Code that chooses contributions and registers them: the
>   integration tests, the tools and a future command.
> - **Private part.** Code that serves one owner and is never selected by a
>   consumer, such as a lexical helper of one language.

The page's words map to this document's as follows: "a source" is a source
reader, "one mapping library" is a mapping, and "one verification law" is a
stated property that a Go type must have to represent a horos type.

A sixth role, **foundation**, has no definition in that section. The layout
page gives it by its members: the carriers and the contract packages. A
stage provider is the parser, checker, formatter, resolver or analyzer of one
language.

The role table of `LAYOUT.md`, abridged to the rows this document uses:

| Role | Paths | May import |
| --- | --- | --- |
| Foundation | `trees/*/go`, `files/go`, `source/go`, `semantics/go`, `contracts/*/go` | Exactly what its own entry in a closed table lists. Never a package outside that table |
| Kernel | `engine/go`, `transform/syntax/shared/*/go`, `check/syntax/shared/*/go`, `print/*/shared/provider/go` | The same table and the same rule. "So no foundation or kernel package imports a vocabulary, a contribution or a language's private part" |
| Shared private helper | `internal/tree`, `internal/grammar`, `internal/lex`, `internal/parseutil`, `internal/concrete` | The same table and the same rule |
| Vocabulary | `syntax/<language>/go`, `semantics/<language>/go`, `transform/syntax/<native>/model/go`, `transform/syntax/<native>/symbol/go` | "Foundations and its own language's vocabulary." The semantic vocabulary and the symbol package may also import their language's lexical helper. "Never a contribution." |
| Contribution: stage provider | `<stage>/<dimension>/<language>/provider/go` | Its own list per stage: foundations, shared private helpers and what belongs to its own language. "A checker gets no concrete tree, a formatter no lexer, and no Syntax stage a semantic contract." |
| Contribution: source reader | `transform/syntax/<language>/source/go` | Foundations and its own syntax and semantic vocabulary. Only horos has one |
| Contribution: native target | `transform/syntax/<native>/target/go` | The generation contract, the AST, the shared tree builder, and its own syntax vocabulary package, intended-forms package and lexical helper. No source reader and no mapping |
| Contribution: mapping | `transform/syntax/<source>-to-<native>/mapping/<name>/go` | The generation contract, the envelope, the syntax vocabulary of both languages, the source reader and the common mapping parts of its source language, and the intended-forms package of its target language. Never another mapping and never a target |
| Private part | `internal/<language>`, `internal/<language>sem` | A lexical helper imports the grammar helper only. A bridge imports foundations and its own syntax and semantic vocabulary |
| Composition | `integration/go`, `verification/integration/go`, `tools/**` | Chooses and registers contributions |

Two sentences of the layout page apply to every role today:

> **One language.** A package that names a language imports no package of
> another language. A mapping names two. A package belongs to a language by its
> path, or by the table `languageOwners` where its path gives no role. No entry
> and no extension lifts this rule.

> Directory groups are not Go packages or an assertion that all members share the
> same dependency level. A package gets its import rule from the role that its
> path gives.

The boundaries page says the same in one line: "A package's role follows from
its path, which the layout fixes."

**The architecture test.** A Go test reads the import block of every
production package and compares it with the package's rule. The rule comes
from one of two places.

- *From the path.* A function splits the package path into segments and
  matches their number and position. `syntax/<language>/go` is a syntax
  vocabulary. A path of five segments whose fourth is `provider` is a stage
  provider. `transform/syntax/<language>/model/go` is an intended-forms package. The language
  is read from the path and substituted into the rule's import patterns. A
  new package of a known shape has its rule without an edit to the test.
  The function is a switch over the number and position of the path's
  segments.
- *From a closed table.* Foundation, kernel and shared helper packages have
  one entry each (31 entries). Packages of which there is only one have a
  hand entry (12 entries). A table of **role extensions** names, per package,
  the imports that one language needs and its role does not give (12
  entries). A table of exceptions names, per importer, the exact symbols it
  may use from a contribution (14 entries).

So a path is not only a location. It is the input from which the import rule
is computed. Any move of a package changes its rule unless the function and
the tables change with it. That is a fact of today's test and no constraint.
Question 4 asks how the two rules of the owner's direction should be
enforced.

One entry each of three of the tables, as they stand in the test
(`integration/go/architecture_test.go:166`, `:282`,
`architecture_exceptions_test.go:23`):

```go
	"trees/representation/go": {"check/syntax/shared/ast/go", "check/syntax/shared/cst/go", "contracts/language/go", "internal/tree", "semantics/go", "trees/ast/go", "trees/cst/go"},
...
	"format/files/typescript/provider/go": {"contracts/generation/go", "files/go", "internal/typescriptimports", "transform/syntax/typescript/symbol/go"},
...
	{Importer: "parse/syntax/golang/provider/go", Provider: "check/syntax/golang/provider/go", Symbols: []string{"Check"}, Reason: admitsOwnTree},
```

The first is the envelope package's entry in the table of foundation and
kernel packages. The second is the role extension of TypeScript's Files
formatter. The third lets Go's parser use one function of Go's checker. A
package whose path gives no role and that has no entry fails the test: "A
production package with neither a role nor an entry fails as unclassified"
(`LAYOUT.md`). That is what a package under a new root meets first.

### 2.7 What is in flight

Standing: decisions of the research owner, except where a paragraph says
otherwise. Nothing in this subsection is on the main branch unless it says
so.

**The tempo.** The work is done by a small number of automated agent
sessions that work in parallel under the owner's rules. The present layout
is two weeks old. A move of the paths of the whole repository is a matter of
days here, not of a quarter. No consumer imports bittype today.

**The migration of the mapping contract.** A design of 13 November 2026 (the
page `docs/contribution-contract.md`, from research 0017) replaces what a
mapping registers and returns. Today a mapping returns intent. Under the new
contract a mapping returns one **decision record** each time the kernel
selects it for one use of a type: a typed value that holds every decision of
that selection, from which the intended form is computed. The type of a
decision record is declared in horos, and its Go type is generated by
bittype itself (section 4.7). The reason is a second reader: a code
generator in another repository must emit, from the same decisions, the
functions that convert between a horos value and the generated type.

The design is final only when its **gate** is met. The gate has seven lines,
each shown by named tests. The seventh is the one that matters here: a
materially different contribution is added, and no shared package changes.
The branch merges when the gate is met.

The migration is built on a branch. Its head is 8d10e00, nine commits ahead
of the code baseline and none behind. This document keeps the main branch as
its code baseline, because that is what is merged. What the branch changes
in the code that this document inlines:

- *Removed.* The mapping contract of the main branch is gone: the types
  `MappingLibrary`, `MappingDescription`, `BindingDescription`, `Pattern`,
  `Capture`, `DeclarationRecipe`, `HelperRecipe`, `Slot` and `HelperKey`,
  the constructors `Mapped` and `Helper`, and the expression kinds of a slot
  and of a helper.
- *Renamed.* In the shared generation packages, `NativeDeclaration` is
  `TargetDeclaration`, `NativeParameter` is `TargetParameter` and the
  constructor `Native` is `Apply`. The expression kinds are five. The
  snippets of section 4.4 that hold one of these names say so.
- *Added, per pair of languages and family.* A family is one sort of source
  type: text, booleans, lists, closed records, optional fields, named
  declarations. The table below lists the new packages.
- *Edited: the four intended-forms packages.* The branch moves none of them
  and edits all four. Each gains a third check and an identity
  (`transform/syntax/golang/model/go/model.go:114-127` on the branch; the
  words "row 7" and "row 13" point into the design):

```go
// EvidenceRecordForm is the evidence kind "the selected native type is this
// model's record form". It belongs to this vocabulary, because the fact
// concerns the native side of this target (row 7). The declaration
// contribution asks for it about a declaration's body (bittype#196, row 3).
var EvidenceRecordForm = g.EvidenceKind{Identity: "bittype/golang/record-form", Version: "1"}

// CheckMember is the member vocabulary of this model (row 13): a member of a
// record is one field form.
func CheckMember(ctx context.Context, e g.Expression[s.Language]) error {
	if name, ok := Name(e.Form()); !ok || name != FormField {
		return g.ErrInvalid
	}
	return CheckExpression(ctx, e)
}
```

| New package on the branch | What it holds |
| --- | --- |
| `transform/syntax/horos-to-<target>/vocabulary/<family>/go` | The identity of one contribution, and the schema and generated Go type of its decision record |
| `transform/syntax/horos-to-<target>/codec/<family>/go` | The hand-written encoding of that decision record |
| `transform/syntax/horos/family/<family>/go` | The identity of the family's meaning, the same for every target |
| `transform/syntax/horos/payload/go` | What every codec shares. The code calls a decision record a payload |

One package of the first sort and one of the third, in their essential lines
(`transform/syntax/horos-to-golang/vocabulary/declarations/go/vocabulary.go:30-38`
and `transform/syntax/horos/family/record/go/family.go:16-18` on the branch;
`rules.Publisher` is the publisher's name, and `family` is the family
package of the same contribution):

```go
var Reference = g.ContributionReference{
	Identity:       g.ContributionIdentity{Publisher: rules.Publisher, Key: "horos-to-golang/declarations"},
	Meaning:        family.Meaning,
	Representation: "golang/1",
}
...
var Schema = g.SchemaReference{Identity: "bittype.payload/horos-to-golang/declarations", Version: "1"}
```

```go
// Meaning is the meaning identity of the family and its revision. They are
// the same in every target.
var Meaning = g.Meaning{Identity: "bittype.structural/closed-record", Revision: "1"}
```

So the branch adds vocabulary packages of a new sort under the stage
directory `transform/`: names that belong to a relation between two
languages, or to one sort of source type for every target, and not to one
language. The layout page says of them: "This layout has no role for a facet
package of a mapping contribution, and none for a vocabulary package of one
family for one target: its role "Vocabulary" is per language." (A facet
package is one of the packages that together make one mapping contribution
on the branch.) These are case (g) of question 1.

**No direction in the shared contract** (decision of the research owner, 13
November 2026, on the owner's message 1). The shared generation contract
called the target side "native". The decision: "The shared part names two
roles of a transformation: source and target." The rename is done on the
branch. Which pairs of languages exist follows from what is registered.

**Two further pairs** (decision of the research owner, 13 November 2026, on
the owner's message 2). Two small pairs are to run through the new contract
before it is final, to show that no direction is assumed. They are built on
the migration's branch:

- `golang-to-typescript`: a source that is not horos. It needs a Go source
  reader over the Go analysis, which no language but horos has today.
- `typescript-to-horos`: horos as a target. It needs an intended-forms
  package and a target for horos, which do not exist on the main branch.

**Two Files formatters.** The Go and the Rust Files formatter are planned as
the next steps of two open work items. By the accepted specification of
destinations (`docs/destinations.md`, 13 November 2026) each is to read its
language's symbol package, as TypeScript's does.

**Written for the present places** (decision of the research owner, and
open). All four pieces of work, the migration, the two pairs, the Go Files
formatter and the Rust Files formatter, are written for the paths of today's
layout. Under today's path pattern a horos intended-forms package would be
`transform/syntax/horos/model/go`, a Go source reader
`transform/syntax/golang/source/go`, and each Files formatter a package
under `format/files/`. Your answer may reorder this work or pause it.
Question 4 asks.

## 3. One language's vocabularies and who uses them

Each subsection states the standing of what it holds. This section and
section 6.3 are what question 2 rests on.

### 3.1 The four vocabulary packages of one language, with Go as the example

Standing: implemented today.

| | Syntax vocabulary package | Semantic vocabulary package | Intended-forms package | Symbol package |
| --- | --- | --- | --- | --- |
| Path | `syntax/golang/go` | `semantics/golang/go` | `transform/syntax/golang/model/go` | `transform/syntax/golang/symbol/go` |
| Lines | 107 | 855 | 112 | 81 |
| What it states | What a Go type-syntax tree is made of: 23 node kinds, 15 child roles, 1 reference role, 9 attribute names, 4 kinds of concrete-tree node | What analysis of a Go program reports: kinds of entities, relations of facts, operators of terms, roles of bindings, all as keys with the owner `go` | What generation intends to construct: 7 form names and 2 primitive names | What a name in a generated declaration refers to when the referent is not itself generated: one category, a predeclared type such as `string` |
| What else it holds | Nothing. Two one-line methods return the language's identifier and display name | Profile strings, the identity values of the resolution and analysis contracts, the inputs that a compiler run is configured with, with their validation and encoding, and typed readers over an analysed graph | Five constructors that build forms, and two checks of a form | The rule that admits a name as a predeclared type, and two functions that record a symbol at a node and read it back |
| Imports of this module | The tree packages, `contracts/language/go` | The semantic graph, the tree, the envelope, the syntax vocabulary package, and the semantic, resolution and analysis contract packages | The generation contract package and the syntax vocabulary package | The generation contract package, the semantic graph, the tree, the envelope and the syntax vocabulary package |
| Who writes its names into an artifact | The parser and the generation target (into trees) | The Go bridge (into the graph) | The mappings into Go (into intent) | The generation target (into the sealed output) |
| Who reads them | The checker, the formatter and the bridge | The typed readers in the same package, a query service and the verification component (section 3.4) | The generation target | Only the generation target today. The planned Go Files formatter is to read them. For TypeScript, the Files formatter already reads its symbol package |
| Stages that import it | parse, check, format, resolve, analyze, generation | resolve, analyze, verification | generation only | generation only (for TypeScript also format in the Files dimension) |

Sections 3.2, 3.3 and 6.3 and question 2 speak of three sets of names:
node kinds, what analysis reports, and intended forms. The symbol package is
the fourth vocabulary package. It names destinations and no type-level
notions, so it has no column of its own in section 3.3.

Three things that today's roles do not count as vocabulary complete the
picture of Go.

- **The grammar table.** Which attributes and children each node kind has is
  stated in the Go checker, a contribution (section 4.2). The syntax
  vocabulary package holds the names only.
- **Lexical rules.** Go has no lexical helper of its own. The checker asks
  the Go standard library whether a word is a keyword, and the symbol package
  asks it which names Go's outermost scope, the universe scope, declares.
  TypeScript, Rust, C++ and horos each have a private package of spelling
  rules.
- **The bridge.** A private package of 1,241 lines runs the Go type checker
  over captured source and writes the graph. It is the writer of every name
  in the semantic vocabulary.

The other languages have the same four, with two gaps. C++ has no semantic
vocabulary package and no symbol package, because it has no analysis. horos
has no intended-forms package and no symbol package, because it is a source
of generation today and never a target. Its counterpart on the generation
side is a source reader, which reads horos's semantic vocabulary.

### 3.2 One record traced from horos to Go

Standing: implemented today for steps 1 to 7, read from the code and
followed by hand, not run. Steps 8 to 10 are planned for Go and are
described from the accepted specification of destinations
(`docs/destinations.md`) and from what the Go bridge writes today. For
TypeScript the same three steps are implemented (section 4.5).

The input is the module of section 1.1:

```horos
module example.people;
type Person = { name: text; };
```

A consumer's composition has registered the horos source reader, the Go
target, and the Go mappings for text, records and declarations. The record
mapping takes a helper key, the identity under which a helper declaration
for optional fields is demanded. The composition must create that key, also
when no field is optional.

| Step | What happens | Who does it (role, place) | Names written and read |
| --- | --- | --- | --- |
| 1 | The horos text is parsed into a concrete tree that owns an abstract tree. The tree is checked: the horos grammar table, then the horos rules in code | The horos parser and checker (contributions, `parse/syntax/horos`, `check/syntax/horos`) | **Writes the horos syntax vocabulary**: nodes of the kinds `program`, `module`, `type-definition`, `descriptor`. Reads it and the horos grammar table |
| 2 | The tree is captured: `semantic.Capture` runs the checker once more and makes a snapshot, the trusted handle of an admitted input (section 2.5). Resolution and analysis take the snapshot and no tree. The result is a semantic graph and two receipts | The semantic contract package, then the horos resolver and analyzer (contributions, `resolve/files/horos`, `analyze/files/horos`), over a private adapter to the horos library | **Writes the horos semantic vocabulary**: an entity of the kind `type-definition`, a provider expression with the operator `descriptor/record`, facts with the relations `sort` and `canonical-identity` |
| 3 | The planner walks the graph from the requested declaration `Person`. For each declaration and each use of a type it asks the source reader what it is | The planner (kernel, inside `contracts/generation/go`) and the horos source reader (contribution, `transform/syntax/horos/source`) | Reads the generic graph. The source reader reads the horos semantic vocabulary, and a few names of the horos syntax vocabulary |
| 4 | For each use one mapping is selected and states intent. The declaration mapping gives the name `Person` and the form `defined`, because the body is a record. The record mapping reads one field, asks the planner for the representation of its type, and returns a record form. The text mapping answers that request with a constant | Three mappings (contributions, `transform/syntax/horos-to-golang/mapping/{declarations,record,text}`) and the planner | **Writes Go's intended forms**: `defined`, `record`, `field` and `primitive`, and the primitive name `string`, through the constructors of the intended-forms package |
| 5 | The planner has every expression and every declaration checked as it composes them, and then the closed program | The Go target (contribution, `transform/syntax/golang/target`). It forwards the first two to the two checks of the intended-forms package and the third, `CheckProgram`, to the shared tree builder | Reads Go's intended forms |
| 6 | The target builds one tree from the closed program. It switches over the form names. `primitive` becomes a `named-type` node with the name `string`, and the use records its destination | The Go target and the shared tree builder (kernel, `transform/syntax/shared/target`) | Reads Go's intended forms. **Writes Go's syntax vocabulary**: `file`, `type-declaration`, `struct-type`, `field`, `named-type`. **Writes a Go symbol**: predeclared `string` |
| 7 | The output is sealed. The sealing code calls `semantic.Capture` on the built tree, which runs the Go checker and makes a second snapshot, of the output. The snapshot is kept with an inventory of the declarations and the list of external uses | The planner's sealing code (kernel) and the Go checker (contribution, `check/syntax/golang`) | Reads Go's syntax vocabulary and Go's grammar table |
| 8 | The declarations are placed into files of a package under the consumer's placement. Planned for Go | A Go Files formatter (contribution; under the present path pattern `format/files/golang`) | Reads Go's syntax vocabulary. Reads Go's symbol package to learn the destination of each external use |
| 9 | The placed files are captured and analysed by the Go type checker. Planned as part of step 8 | The Go resolver and analyzer over the Go bridge (private part, `internal/golangsem`) | **Writes Go's semantic vocabulary**: see the graph below |
| 10 | Each name use is compared with its destination: the entity that the compiler bound against the entity that the symbol identifies. Planned for Go | The Go Files formatter | Reads Go's semantic vocabulary (the binding, the analysed kind `builtin`) and Go's symbol package (the category `predeclared`) |

**Where each of Go's vocabulary packages is written and read in this
trace.**

| Go vocabulary package | Written in step | By which stage | Read in step | By which stage |
| --- | --- | --- | --- | --- |
| Intended forms (`transform/syntax/golang/model/go`) | 4 | generation (mappings) | 5, 6 | generation (target) |
| Syntax (`syntax/golang/go`) | 6 | generation (target) | 7, 8, 9 | check; format; resolve and analyze |
| Symbol (`transform/syntax/golang/symbol/go`) | 6 | generation (target) | 8, 10 | format in the Files dimension (planned) |
| Semantic (`semantics/golang/go`) | 9 | resolve and analyze (bridge) | 10 | format in the Files dimension (planned) |

Three observations. The first two are implemented today; the third is partly
our reading.

- The intended-forms package is written and read inside one stage. The
  symbol package is written in one stage and read in another. The syntax
  vocabulary package and the semantic vocabulary package are used by several
  stages.
- Step 10 is the one place in this trace where intent and analysed meaning
  of one language are compared. The comparison is of identities (a symbol
  against an entity), not of type structure. No production code today
  compares an intended `record` form with an analysed `struct` term.
- The source side has no intended-forms package. Generation reads horos
  through its analysed vocabulary. Our reading of why: a source is an
  analysed program, and a target program does not exist yet when it is
  described. No page states this asymmetry as a decision.

### 3.3 The same record, then the same notions, in three sets of names

Standing: the names are implemented today and verified against the code.
Three sets of names exist, one per sort of artifact: node kinds for trees,
analysed kinds for graphs (an analysed kind is the kind of an entity or the
operator of a term), and forms for intent. The tables lay them beside each
other. The pairing of rows is ours, and so is the last column of the second
table.

**The record `Person`.**

| | Intended (after step 4) | Syntax tree (after step 6) | Established by analysis (step 9, planned for Go) |
| --- | --- | --- | --- |
| The declaration | Declaration form `defined`, name `Person` | Node `type-declaration` with the attribute `name` | Entity of the analysed kind `defined`; a fact with the relation `underlying` |
| Its body | Form `record` | Node `struct-type` | Term with the operator `struct` |
| The field | Form `field` with the data `Name` and `false` | Node `field` with the attribute `name` | Entity of the analysed kind `field`; one JSON entry and one operand of the struct term |
| The field's type | Form `primitive` with the datum `string` | Node `named-type` with the name `string`, and the symbol predeclared `string` | Term with the operator `basic`; a binding with the role `name` to an entity of the analysed kind `builtin` |

The graph that the bridge would publish, in the words of the semantic
vocabulary. This is an illustration written by hand from the bridge's code
(section 4.3). A key is written owner/name. E is an entity and T a term:

```text
entity E1   kind go/defined   name "Person"   owner: the package entity   site: the type-declaration node
            body -> T1
term   T1   reference to E1
fact        subject T1, relation go/underlying -> T2
term   T2   operator go/struct, encoding go/json-1
            data [{"Name":"Name","Tag":"","Embedded":false,"Entity":E2}]
            operand 0 -> T3
term   T3   operator go/basic, data: the compiler's kind of string
entity E2   kind go/field     name "Name"
entity E3   kind go/builtin   name "string"   no owner, no site
binding     at the named-type node, role go/name -> E3
```

Two things in it matter later. The body of the entity `Person` is a
reference term to the entity itself, and the definition hangs on that term
as a fact. Every fact about a type has a term as its subject and never an
entity. The illustration is abridged: the display fact on every term, the
facts with the relation `identical`, the package entity with its
`package-name` fact, the bodies of E2 and E3 and the entries for type
expressions are left out.

**Go's type-level notions.**

| Notion | Node kind (`syntax/golang/go`, as the target writes it) | Established by analysis (`semantics/golang/go`, as the bridge writes it) | Intended form (the intended-forms package) | Our reading of the difference |
| --- | --- | --- | --- | --- |
| Predeclared basic type (`string`, `bool`) | `named-type` with the attribute `name`; the destination is recorded as a symbol `predeclared` | Entity kind `builtin`, named by the identifier, without owner and site; term operator `basic` with the compiler's basic kind as data | `primitive` with one datum, `string` or `bool` | Four words. Intent names the type it means. Analysis names the entity the compiler found, which a declaration named `string` can shadow |
| Slice | `slice-type` with one `element` child | Term operator `slice` with one operand | `slice` with one argument | The same word three times |
| Struct with fields | `struct-type` with `field` children; also `embedded-field` | Term operator `struct` with JSON entries (name, tag, embedded, entity) and one operand per field; entity kind `field`, an unpublished literal of the bridge | `record` whose arguments are `field` forms; a field has the data name and optional flag | Three words for the struct, one for the field. `record` is the word the four intended-forms packages share. The optional flag exists only in intent, and the Go target does not read it. Tag and embedding exist only in syntax and analysis: nothing generates them |
| Defined type declaration | `type-declaration` | Entity kind `defined`; fact relation `underlying` | Declaration form `defined`; the body must be a record | The same word. The restriction to a record body is a limit of today's generation |
| Alias declaration | `type-alias` | Entity kind `alias`; fact relation `alias-target` | Declaration form `alias` | The same word |
| Type of another package | `qualified-type` and `import`; the target writes neither | Entity kind `package` with a `package-name` fact (section 4.3); a `defined` or `alias` entity owned by it | `external` with the data module and name; refused as unsupported | Intent has a name and no meaning yet |
| Type parameter | `type-parameter` with a constraint child | Entity kind `parameter`; fact relation `constraint`, whose interface term carries a fact `empty` | No form. An expression kind of the shared contract package | The binder's identity is the kernel's on the intent side. The constraint exists only in syntax and analysis |
| Use of a declaration, generic application | `named-type` with a `target` reference; `type-application` | The graph's own reference and application terms, and a binding with the role `name` | No form. An expression kind of the shared contract package | Owned by the generic data types on both sides |
| Package | The attribute `package` on `file` | Entity kind `package` | Not a form. A setting of the target | |
| Only established by analysis | | Fact relations `identical` and `display`; operators `interface`, `union`, `signature`, `tuple`, `boolean` | | Type identity exists only for a checked program. The other operators are Go notions that nothing generates |
| Only in syntax | `pointer-type`, `array-type`, `map-type`, `channel-type`, `function-type`, `method`, `embedded-type`, `type-union`, `type-approximation` | The bridge refuses pointer, array, map and channel types as unsupported | | Notions that syntax can spell and that are neither analysed nor generated |

Summary of the pairing, which is ours: all seven forms of Go's
intended-forms package stand beside a notion that analysis also names. Three
do so under another word (`primitive`, `record`, `external`) and four under
the same word (`slice`, `field`, `alias`, `defined`). `field` is an
unpublished literal of the bridge. Every form also stands beside a node
kind.

Two structural differences are visible, and both are our reading. First, an
intended form is a tree of values with placeholders, without identity, site
or spelling. A node has a site in one tree. An analysed type is a term in a
graph with entities that have identity. Second, intent refers to a
destination before any program exists, and analysis refers to the entity
that a compiler bound in one placed program.

**TypeScript and Rust.** The picture of three sets of names holds for both,
checked on one notion each. TypeScript's collection: the form `array`, the
node kind `array-type`, and the operator `array` that the TypeScript bridge
writes as a literal in its script. Rust's nominal declaration: the
declaration form `struct`, the node kind `struct`, and the entity kind
`struct` that the Rust bridge writes as a literal in its Rust program. C++
has two of the three, because it has no analysis, and horos has two, because
it has no intended forms.

### 3.4 The other readers of the analysed vocabulary

Standing: implemented today, except the last item, which is work in flight.

In the trace Go's semantic vocabulary is read only in steps that are
planned. It has readers today, and one more is coming:

- **A query service** (`analyze/files/golang/queries/go`, a contribution)
  answers explicit questions about an opened graph of analysed Go.
- **The verification component** (`verification/files/horos-to-go/go` and
  one of its law packages) reads analysed Go types through the typed
  readers, to check a stated property that a Go type must have to represent
  a horos type. Generation's mappings state the same Go types as intended
  forms. The two do not share names (section 5.2, item 4).
- **TypeScript's Files formatter** reads TypeScript's semantic vocabulary
  package and its symbol package and compares identities (section 4.5).
- **A Go source reader**, to come with the pair of languages
  golang-to-typescript (section 2.7). There Go is the source of generation
  and is read through its semantic vocabulary, while Go as a target is
  written through its intended forms.

### 3.5 Where one analysed vocabulary could not be read as a source

Standing: the first paragraph is implemented today. The second is a finding
on a branch that is not merged. The third is a decision of the research
owner. It stands, and you may object to it.

**What the kernel walks, and what three analyses publish.** The generation
kernel walks a source's semantic graph from the body of an entity through
the operands of terms. It reads no fact
(`contracts/generation/go/occurrence.go:202`):

```go
body = &occurrence[F]{walk: w, term: e.Body, declaration: d.entity}
```

The horos analysis publishes a declaration's definition as the body of its
entity. The Go analysis publishes, as the body of a named type, a reference
to the entity itself, and the definition as the value of a fact that Go owns
(`internal/golangsem/export.go:158-159`, `:207`, `:286`):

```go
body := e.typ(o.Type())
e.data.Entities[id-1].Body = body
...
term = semantics.Term{Kind: semantics.Reference, Target: obj}
...
e.fact(id, n.UnderlyingRelation(), e.typ(t.Underlying()))
```

The TypeScript analysis does the same and says why
(`internal/typescriptsem/compiler.cjs:251-252`, `:316`). There the
definition is the value of the facts `properties` and `structure`:

```javascript
// References break recursive type graphs; definition structure is a fact,
// not a containment edge back into the named type.
...
graph.Entities[id - 1].Body = type;
```

Three lines of a dump of the Go graph for a struct type `Person`, made by
the lane named below: the entity, the term that is its body, and the fact.
The term of kind 1 is a reference, and its target is the entity itself:

```text
entity 7 kind=go/defined name="Person" owner=1 body=10 decls=[{0 12}] params=[]
term 10 kind=1 target=7 callee=0 args=[] op=/ data="" operands=null
fact subject(e0,t10) go/underlying -> t11
```

The typed readers of the two semantic vocabulary packages read the body in
that sense. Go's readers take `e.Body` as the type of the declared name in
two places (`semantics/golang/go/model.go:111`, `:253`), and TypeScript's
reader reaches the declared type through a fact with the relation
`declared-type` (`semantics/typescript/go/model.go:130`).

**What happened.** One working session built the two further pairs of
section 2.7 on the migration's branch. The first pair with Go as a source
was refused with an evidence cycle before any field of a struct became a
use of a type: the body of a declaration led back to the declaration. The analysed vocabulary of Go
and of TypeScript was built to be compared against, by verification and by
a Files formatter. It could not be walked as a source.

**The decision** (decision of the research owner, 13 November 2026). The
kernel stays as it is. The Go and the TypeScript analysis change: the body
of a declaration's entity becomes the term of its definition, and the link
from a declared name to its own type becomes a fact on the entity. The
decision rests on two accepted sentences. Of the kernel's walk: "No source
supplies a traversal" (`docs/generation-gaps.md`). Of the graph: "`Owner`
and `Body` are optional links, not a universal namespace or type-of rule."
(`docs/semantics.md`). Two alternatives were not taken. In the first, a
source names the term at which the kernel starts a declaration's body. In
the second, each such language gets a second analysis that publishes a
graph in the shape the kernel walks.

**Smaller facts from the same work,** each a place where one model of a
language for every operation is under strain. They are the lane's notes,
and its report is not finished.

- The definition is not the checked type. For Go's `type Temp Celsius` the
  written definition is a reference to `Celsius`, and the value of
  `underlying` is `float64`. In TypeScript an alias of an
  interface resolves to the interface's written-out properties.
- The Go graph shares one term among all occurrences of a type. The horos
  graph has a term per occurrence. A reader that needs the site of one use
  cannot get it from a shared term.
- A family (section 2.7) has an accepted meaning only for horos as the
  source: a set of canonical horos values. For a source that is not horos
  no page defines it.
- The contract's two declaration forms, alias and nominal, are the
  vocabulary of a native target. horos has one declaration form.
- For horos as a target, an intended-forms package had to be declared
  beside horos's syntax vocabulary and its semantic vocabulary. The typed
  readers of TypeScript's semantic vocabulary package start from a
  declaration or a type and give no handle for a term that the kernel hands
  over, so the source reader decodes the term itself.

## 4. Relevant code, as evidence

Standing of this section: implemented today. It is evidence of the present
state and never a reason for or against an option. Where a table classifies
code, the classification is our derivation, and the text says so.

A reader short of time can skip this section and return to it from section
6. Which part serves which question:

| Part | What it shows | Serves |
| --- | --- | --- |
| 4.1 | The generic tree, the names one language adds, and hand-written code that spells those names again | Questions 1 (d) and 3 |
| 4.2 | A language's grammar table, the rules in code beside it, and a lexical rule | Questions 1 (c), 1 (h) and 3 |
| 4.3 | The semantic graph, one language's names for it, and one fact written by a compiler bridge and read by a typed reader | Questions 1 (b), 2 and 3 |
| 4.4 | The data types of intent, Go's intended-forms package complete, and the target that reads it | Questions 1 (a), 1 (i), 2 and 3 |
| 4.5 | A symbol package | Questions 1 and 2 |
| 4.6 | The node patterns of horos kinds, a notation for tree shapes | Question 3 |
| 4.7 | The one generator of Go code from a declaration | Question 3 |

### 4.0 Reading guide

Every snippet is taken from the main branch at the code baseline, with its
path and line. Comments are the code's own unless a line starts with
`// [doc]`, which marks a comment added for this document. An ellipsis `...`
marks omitted lines. A function shown without a body is its signature only.

Import aliases recur. In every language's packages `s` is that language's
syntax vocabulary package, and `n` is its semantic vocabulary package where
both occur. `g` is the grammar helper `internal/grammar` inside a checker,
and the generation contract package `contracts/generation/go` inside an
intended-forms package, a symbol package and a target. The two never meet in
one file. `ast` is the generic tree package and `semantics` the generic
semantic graph. `text` is the diagnostics package, whose function
`text.AtNode(tree, id, message)` returns an error value that points at one
node. `tree.View` is a small read helper over one tree: `Node(id)` returns a
node, `Attr(id, name)` an attribute's text, `Has` and `Flag` test an
attribute, and `Children(id, role)` and `Child(id, role)` return the children
under one child role. Inside a checker `bad(message)` is a local shorthand
for `text.AtNode` at the current node. `representation.Language` is the
descriptor of a language from the envelope package `trees/representation/go`:
its identifier, display name, profile string and marker type. `types` and
`token` are packages of the Go standard library.

**Ten lines of Go, for a reader who does not write it.**

- `type AST[L any] struct { ... }` is a type that is generic in `L`; `any`
  means that `L` is unconstrained. `func F[L any]( ... )` is a generic
  function.
- `[]T` is a sequence of `T`, `map[K]V` a map, and `*T` a pointer, whose
  empty value `nil` means absent. `kinds ...T` is a variable number of
  arguments.
- `type A = B` is an alias: a second name for the same type.
- `func (t T) M()` is a method `M` of the type `T`.
- `struct{}` is a type without fields. `var _ I = T{}` asserts at compile
  time that `T` implements the interface `I`.
- `iota` numbers the constants of one block: 0, 1, 2.
- A name that starts with a capital letter is exported. A name that starts
  with a small letter is private to its package, and so is a struct field.
  The issuers and handles of this document rest on that: a type with private
  fields can be constructed only in its own package.
- A package below a directory named `internal` can be imported only from the
  tree around that directory (fixed point 21).

### 4.1 The generic tree, the names one language adds, and code that spells them again

The node of the generic abstract tree (`trees/ast/go/node.go:6-44`,
abridged). The tree itself is a flat list of nodes with a root index. A
`NodeID` is a one-based index into that list. A `NodeRef` is the address of
a node in some tree, and an `Origin` holds source spans and contributing
nodes:

```go
type Kind[L any] string

type Attribute struct {
	Name string
	Text string
}
type Child struct {
	Role string
	Node NodeID
}
type Reference struct {
	Role   string
	Target NodeRef
}
type Node[L any] struct {
	Kind       Kind[L]
	Attributes []Attribute
	Children   []Child
	References []Reference
	Origin     Origin
}
```

The package comment states what the tree type is not (`trees/ast/go/ast.go:7-8`):
"These are data types, not checked trees. Structural and grammatical checks
live in separate packages; the carrier performs no validation or
allocation." The comment on the tree type adds (`ast.go:27-28`): "Fields
remain mutable: callers must treat stage inputs as read-only and create
separate output storage. No immutability or validity is enforced here." The
comments on a
node's parts say who owns their meaning (`node.go:8-10`, `:18-19`,
`:25-27`): "Attribute is an uninterpreted textual syntax atom." "A grammar
defines the name, allowed text and meaning." "The grammar owns Role names
and cardinality; repeated roles can form a list." "A grammar owns Role and
the allowed target."

Go's syntax vocabulary package (`syntax/golang/go/nodes.go:1-95`, abridged:
the imports and some comments are left out, and the node kinds, roles and
attributes are shortened as marked). `language.Provider` is a two-method
interface of the contract package `contracts/language/go` that returns a
language's identifier and display name:

```go
// Package golang declares the Go type-syntax vocabulary over the generic AST.
// It contains no parser, builder, grammar checker, generator or formatter.
// The node shapes and this vocabulary's limits are documented in docs/syntax.md.
package golang

// Language identifies this syntax family and is its distinct Go AST marker.
type Language struct{}

func (Language) ID() language.ID { return "go" }
func (Language) Name() string    { return "Go" }

var _ language.Provider = Language{}

type AST = ast.AST[Language]
type Node = ast.Node[Language]
type Kind = ast.Kind[Language]

// Node kinds for file structure, type declarations and native type expressions.
const (
	KindFile              Kind = "file"
	KindTypeDeclaration   Kind = "type-declaration"
	KindNamedType         Kind = "named-type"
	KindStructType        Kind = "struct-type"
	KindField             Kind = "field"
	KindSliceType         Kind = "slice-type"
	// [doc] and 17 more, among them KindTypeAlias, KindTypeParameter,
	// [doc] KindTypeApplication, KindPointerType, KindArrayType, KindMapType
)

// [doc] 15 child role names (RoleType is "type"), 1 reference role name (RoleTarget is
// [doc] "target") and 9 attribute names (AttributeName is "name") follow, as untyped
// [doc] string constants
```

What this shows, as facts:

- A node kind is a typed constant. A child role name and an attribute name
  are untyped strings.
- The package states names and nothing about shape. It does not say that a
  `field` has a `name` attribute and one `type` child. It points to prose
  for that.
- The other four syntax vocabulary packages have the same shape. The number
  of node kinds is: horos 33, TypeScript 29, Go 23, Rust 36, C++ 23.
- Every language also declares the same four kinds of concrete-tree node:
  `layout`, `token`, `whitespace`, `comment`.

**Hand-written code spells the names again.** A piece of Go's formatter
(`format/syntax/golang/provider/go/format.go:100-107`; `f.Kind`, `f.Attr`
and `f.Children` read the tree, `f.T` makes a token, and `t(role)` formats
the child under one role):

```go
switch f.Kind(id) {
case "named-type":
	parts = []doc{f.T(f.Attr(id, "name"))}
case "qualified-type":
	parts = []doc{f.T(f.Attr(id, "qualifier")), f.T("."), f.T(f.Attr(id, "name"))}
case "type-application":
	var args []doc
	for _, a := range f.Children(id, "argument") {
	...
```

And one line of Go's parser, which converts the tree of the standard
library's parser (`parse/syntax/golang/provider/go/parse.go:67`; `tr.Child`
makes a child edge and `c.node` a node):

```go
imports = append(imports, tr.Child("import", c.node(s.KindImport, spec, attrs)))
```

A count over the non-test files, made for this document: the parser and the
formatter of Go, TypeScript, Rust and C++ are eight packages, and they use
the child role and attribute constants of their syntax vocabulary package 0
times. They write these names as string literals. The four formatters also
use no node kind constant and switch on literals, as above. The horos parser
uses the role and attribute constants 65 times and the horos formatter 53
times. The four targets write attribute names as literals too, for example
`m.Attributes("name", name)`. So the vocabulary package is not the one place
where these names are spelled.

### 4.2 The grammar table, the rules beside it, and a lexical rule

**The table type.** A private helper package declares what a grammar table
is and holds its one interpreter (`internal/grammar/rules.go:15-40`). The
code calls it the admission table:

```go
type Attribute struct {
	Required bool
	Accept   func(string) bool
}
type Role[L any] struct {
	Name     string
	Min, Max int
	Kinds    []ast.Kind[L]
}
type Rule[L any] struct {
	Attributes map[string]Attribute
	Roles      []Role[L]
	Reference  bool
}

// Rules is one language's admission table. Its keys and the kinds of its roles
// are that language's node kinds, so a table of one language cannot be applied
// to a tree of another. Role and attribute names are strings.
type Rules[L any] map[ast.Kind[L]]Rule[L]

func Attr(required bool, accept func(string) bool) Attribute { return Attribute{required, accept} }
func One[L any](name string, kinds ...ast.Kind[L]) Role[L]   { return Role[L]{name, 1, 1, kinds} }
func Opt[L any](name string, kinds ...ast.Kind[L]) Role[L]   { return Role[L]{name, 0, 1, kinds} }
func Many[L any](name string, min int, kinds ...ast.Kind[L]) Role[L] {
	return Role[L]{name, min, -1, kinds}
}
```

The same file supplies attribute predicates (`Set(values...)`, `Bool`,
`Nonempty`, `String`, `Decimal`, `Positive`, `Signed`, `Identifier`) and a
helper `Attributes(name, attribute, name, attribute, ...)` that builds the
attribute map. The interpreter is one function of 67 lines
(`internal/grammar/rules.go:80`):

```go
func Check[L any](tree ast.AST[L], root ast.Kind[L], rules Rules[L]) error
```

In order it checks: the structure of the generic tree, by a call into
`check/syntax/shared/ast/go`; the node kind of the root; and for every node,
that its node kind has a rule, that every attribute is declared and
accepted, that every required attribute is present, that every child's role
is declared and that children come in the declared order of roles, that each
child's node kind is in its role's list, that each role's count is within
its bounds, and that a reference is present only where the rule allows one
and only under the reference role `target`.

So one entry can express, per node kind: a set of attributes, each required
or optional, each with a predicate on its text; an ordered list of child
roles, each with a count range and a flat list of allowed node kinds; and
whether a reference may be attached.

**One table.** Go's is a package variable of the Go checker
(`check/syntax/golang/provider/go/check.go:13-41`; 7 of its 23 entries are
shown). Here `g` is the grammar helper:

```go
func identifier(v string) bool { return g.Identifier(v) && !token.Lookup(v).IsKeyword() }

var types = []s.Kind{s.KindNamedType, s.KindQualifiedType, s.KindTypeApplication, s.KindStructType, s.KindInterfaceType, s.KindFunctionType, s.KindPointerType, s.KindSliceType, s.KindArrayType, s.KindMapType, s.KindChannelType}
var constraints = append(append([]s.Kind{}, types...), s.KindTypeUnion, s.KindTypeApproximation)
var rules = g.Rules[s.Language]{
	s.KindFile:              {Attributes: g.Attributes(s.AttributePackage, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.Many(s.RoleImport, 0, s.KindImport), g.Many(s.RoleDeclaration, 0, s.KindTypeDeclaration, s.KindTypeAlias)}},
	...
	s.KindTypeDeclaration:   {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.Many(s.RoleTypeParameter, 0, s.KindTypeParameter), g.One(s.RoleType, types...)}},
	...
	s.KindNamedType:         {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier)), Reference: true},
	...
	s.KindField:             {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier), s.AttributeTag, g.Attr(false, g.String)), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	...
	s.KindArrayType:         {Attributes: g.Attributes(s.AttributeLength, g.Attr(true, g.Decimal)), Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	...
	s.KindChannelType:       {Attributes: g.Attributes(s.AttributeDirection, g.Attr(false, g.Set("both", "send", "receive"))), Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindTypeUnion:         {Roles: []g.Role[s.Language]{g.Many(s.RoleTerm, 2, append(types, s.KindTypeApproximation)...)}},
	...
}
```

Facts about the tables as data:

- Each language's table has one entry per node kind constant.
- The table is an unexported variable of a contribution, the checker. No
  other package can read it. The parser and the formatter of the language
  call the checker's exported function `Check` and nothing else of it.
- The table is a Go map literal that holds function values. In Go's table
  they are `identifier`, which asks the standard library whether a word is a
  keyword, and an inline function for the import alias. The other languages
  pass predicates of their lexical helpers.
- A documentation tool reads the tables as data. Because the value holds
  functions and is unexported, the tool does not import it. It parses the Go
  source text of the literal and evaluates the initializer syntactically. A
  comment marks the table for it: `// bittype:grammar rules`.

**The rules beside the table.** Each checker runs the table and then code
(`check/syntax/golang/provider/go/check.go:45-83`, the body left out):

```go
func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, s.KindFile, rules); err != nil {
		return err
	}
	...
}
```

Three rules from other checkers, each with what it expresses that the table
type cannot:

```go
// [doc] check/syntax/typescript/provider/go/check.go:60-72
// [doc] A condition on ancestors: an infer type is admitted only below the
// [doc] constraint child of a conditional type, at any depth.
var context func(ast.NodeID, bool) error
context = func(id ast.NodeID, inConstraint bool) error {
	n := v.Node(id)
	if n.Kind == s.KindInferType && !inConstraint {
		return text.AtNode(t, id, "infer is only supported within a conditional extends operand")
	}
	for _, child := range n.Children {
		if err := context(child.Node, inConstraint || (n.Kind == s.KindConditionalType && child.Role == s.RoleConstraint)); err != nil {
			return err
		}
	}
	return nil
}
```

```go
// [doc] check/syntax/rust/provider/go/check.go:100-107
// [doc] An order among the children of one role: the table lists allowed
// [doc] kinds per role, and no order between them.
seenType := false
for _, p := range v.Children(id, s.RoleGenericParameter) {
	if v.Node(p).Kind == s.KindTypeParameter {
		seenType = true
	} else if seenType {
		return bad("lifetime parameters must precede type parameters")
	}
}
```

```go
// [doc] check/syntax/cpp/provider/go/check.go:72-80
// [doc] A condition on the parent, and a relation between two attributes of
// [doc] one node. parents is a map from a node to its parent, built first.
if v.Has(id, s.AttributeAccess) && n.Kind != s.KindBase {
	parent := parents[id]
	if parent == 0 || v.Node(parent).Kind != s.KindRecord {
		return bad("access belongs to record members")
	}
}
if v.Flag(id, s.AttributePure) && !v.Flag(id, s.AttributeVirtual) {
	return bad("pure declaration must be virtual")
}
```

The count of such rules per checker, with the sorts of thing they express
(the grouping into sorts is our derivation):

| Language | Table entries | Rules in code | What the rules in code express |
| --- | --- | --- | --- |
| Go | 23 | 3 | Consistency among siblings of one role (all parameters named or none); position within a role (variadic last); a condition two levels down (the element of an embedded pointer) |
| TypeScript | 29 | 14 | A condition on ancestors; an attribute against another attribute, with the lexical rules of literals; an attribute against a child's node kind; order and position within a role (rest and optional elements); consistency among siblings (tuple labels) |
| Rust | 36 | 13 | A condition that depends on the parent's node kind and the child's role; order within a role (lifetimes before types, arguments before bindings); what a sort of declaration inherits (visibility, defaults) |
| C++ | 23 | 11 | A condition on the parent; an attribute against another attribute; which node kinds may nest (no array of functions, no reference to a reference) |
| horos | 33 | 11 | Which reference roles each node kind may carry (the table knows one role, `target`); uniqueness of names among siblings; contiguous indices; that an index refers to a binder in scope, through a walk that carries lexical depth |

The horos checker removes all references from a copy of the tree before it
runs the table, because "The shared grammar checker supports the native
target role. Horos has several annotation roles, checked below"
(`internal/declsyntax/check.go:127-128`, the private package that holds the
horos checker's code).

**What no table or rule states.** The concrete syntax. Four of the five
parsers are hand-written recursive descent over a shared scanner. The Go
parser calls the parser of the Go standard library and converts its tree.
Every formatter is a hand-written function per node kind. A parser calls the
checker after it has built its tree, and a formatter calls it before it
reads its input. That call is the only connection between the table and the
code that produces and consumes trees.

**A lexical rule, and a use of it from a vocabulary-side package.** A
lexical helper is a private package with one language's spelling rules. One
predicate of Rust's (`internal/rust/cratename.go:19-27`;
`crateNameCharacters` is a table of 840 ranges of code points, printed by a
program built with the pinned Rust compiler and pasted in by hand):

```go
func CrateName(name string) bool {
	for _, r := range name {
		at := sort.Search(len(crateNameCharacters), func(i int) bool { return crateNameCharacters[i][1] >= r })
		if at == len(crateNameCharacters) || crateNameCharacters[at][0] > r {
			return false
		}
	}
	return name != ""
}
```

The Rust checker and parser use this helper. So does Rust's symbol package,
which our test of section 6.1 puts on the model side
(`transform/syntax/rust/symbol/go/symbol.go:72-80`; `native` is the lexical
helper, `g` the generation contract package, and `Std()` the standard
library's crate):

```go
func Dependency(name string) (Crate, error) {
	if !native.CrateName(name) {
		return Crate{}, g.ErrInvalid
	}
	if name != Std().name {
		return Crate{}, g.ErrUnsupported
	}
	return Std(), nil
}
```

The function first applies the lexical rule and answers that a text is no
crate name. It then gives a second answer of another sort: a valid name that
today's generation does not support. Of the four lexical helpers only
`internal/rust` and `internal/typescript` are imported from both sides:
by stage providers, and by their language's semantic vocabulary package and
symbol package. `internal/cpp` and `internal/horos` are imported by stage
providers and a target only. Go has no lexical helper; its checker and its
symbol package ask the Go standard library.

### 4.3 The semantic graph, Go's names for it, and one fact written and read

**The graph.** One package holds the data type of analysed meaning for every
language (`semantics/go/model.go:19-146`, `:230`, abridged; the file is
named `model.go`, which is unrelated to the intended-forms package):

```go
// IDs are one-based addresses local to one Semantics[L], never persistent identities.
type EntityID uint32
type TermID uint32

// Key names vocabulary owned by a provider. It causes no registration or lookup.
type Key struct{ Owner, Name string }

type Entity struct {
	Kind         Key
	Name         string
	Owner        EntityID
	Declarations []ast.NodeRef
	Parameters   []Parameter
	Body         TermID
}
...
type TermKind uint8

const (
	Reference TermKind = iota + 1
	Application
	ProviderExpression
)

// Term has exactly one variant. ...
// Reference owns Target; Application owns Callee and Arguments. Applicability,
// arity and saturation belong to the provider, not graph admission.
// A provider expression's encoding defines its bytes and operand locators.
type Term struct {
	Kind      TermKind
	Target    EntityID
	Callee    TermID
	Arguments []TermID
	Operator  Key
	Encoding  Key
	Data      []byte
	Operands  []Operand
}
...
// Subject selects exactly one model-local entity or term.
type Subject struct {
	Entity EntityID
	Term   TermID
}

// Fact associates a subject with a provider-owned fact term. A relation may have
// several values; its meaning and cardinality belong to the provider.
type Fact struct {
	Subject  Subject
	Relation Key
	Term     TermID
}
...
// Semantics[L] is an immutable owned graph.
type Semantics[L any] struct{ state *Data }

func New[L any](data Data, limits Limits, contains func(ast.NodeRef) bool) (Semantics[L], error)
```

The omitted types are small records: a parameter of an entity with a sort;
an operand of a provider expression; an occurrence (a node, a role and an
ordinal); a binding of an occurrence to an entity; an expression at an
occurrence; a dependency between two entities; the storage type `Data` with
one list per sort of entry; and `Limits`.

The graph fixes three variants of a term and the shapes of its entries.
Everything else is a key that the language's own code owns: the kind of an
entity, the sort of a parameter, the role of an occurrence, the relation of
a dependency or a fact, the operator and the encoding of a provider
expression. The constructor `New` checks integrity and bounds. It does not
know what any key means, and it does not look into `Data` bytes.

**Go's names.** The semantic vocabulary package of Go declares its keys as
one-line functions, because a Go struct value cannot be a constant
(`semantics/golang/go/vocabulary.go:13-71`, abridged as marked):

```go
const ResolutionProfile = "go1.25/captured-types-1/resolution-1"
const AnalysisProfile = "go1.25/captured-types-1/analysis-1"
// [doc] also LanguageVersion, PackagePath and two error values

// Key returns a Go-owned vocabulary key; it performs no registration.
func Key(name string) semantics.Key { return semantics.Key{Owner: "go", Name: name} }
// [doc] also Bindings, Coverage and Validity: the obligations on which an analysis reports

// The graph names below are written by the compiler bridge internal/golangsem
// and read by name in model.go.

// Kinds of graph entities.
func PackageKind() semantics.Key   { return Key("package") }
func DefinedKind() semantics.Key   { return Key("defined") }
func AliasKind() semantics.Key     { return Key("alias") }
func ParameterKind() semantics.Key { return Key("parameter") }

// BuiltinKind is the kind of the entity of a predeclared basic type, such as
// string and bool. The entity's name is the identifier of the universe scope.
// It has no owner and no declaration site.
// ...
func BuiltinKind() semantics.Key { return Key("builtin") }

// NameRole is the role of the binding of an identifier that is used as a
// name: the binding's site is the node of the use, and its target is the
// entity that the compiler resolved the identifier to.
func NameRole() semantics.Key { return Key("name") }

// Relations of the facts about a native type.
func IdenticalRelation() semantics.Key   { return Key("identical") }
func UnderlyingRelation() semantics.Key  { return Key("underlying") }
func ConstraintRelation() semantics.Key  { return Key("constraint") }
func AliasTargetRelation() semantics.Key { return Key("alias-target") }
func EmptyRelation() semantics.Key       { return Key("empty") }
func DisplayRelation() semantics.Key     { return Key("display") }

// Relation of the fact about a package entity: the name that its package
// clauses declare. The entity's own name is the import path. The fact's term
// is a provider expression with PackageNameOperator and TextEncoding and no
// operand, whose data is the name. The name readers refuse another form.
func PackageNameRelation() semantics.Key { return Key("package-name") }
func PackageNameOperator() semantics.Key { return Key("package-name") }

// TextEncoding is the encoding of a term whose data is UTF-8 text.
func TextEncoding() semantics.Key { return Key("utf8") }

// Operators of the terms of a basic, a slice and a struct type.
func BasicOperator() semantics.Key  { return Key("basic") }
func SliceOperator() semantics.Key  { return Key("slice") }
func StructOperator() semantics.Key { return Key("struct") }
```

The comment on `PackageNameRelation` is the one place in this file where a
declaration's comment states the shape of a fact: which operator, which
encoding, no operand, and what the data is. For every other fact the shape
is stated in code only, as the next two snippets show.

The bridge writes more names than the vocabulary package publishes. It
spells these as literals in its own code: the entity kind `field`, the
parameter sort `type`, the operators `interface`, `boolean`, `union`,
`signature`, `tuple` and `display`, the binding role `import`, the
expression role `type-expression` and the encoding `json-1`
(`internal/golangsem/export.go:111`, `:135`, `:169`, `:253-288`, `:347`,
`:371`). The accepted criterion for when a name gets a published declaration
is: "A name needs a single owner when two parties must agree on it"
(`docs/package-boundaries.md`). No production package other than the bridge
reads one of these names by name.

**How the bridge chooses the kind and the identity of an entity**
(`internal/golangsem/export.go:128-153`, abridged; `o` is an object of the
Go type checker, `types` is the standard library's package `go/types`, and
`e.objects` remembers the entity of each object):

```go
func (e *exporter) object(o types.Object) semantics.EntityID {
	if id := e.objects[o]; id != 0 {
		return id
	}
	...
	kind := n.Key("field")
	owner := semantics.EntityID(0)
	if _, ok := o.(*types.TypeName); ok {
		kind = n.BuiltinKind()
		switch o.Type().(type) {
		case *types.Named:
			kind = n.DefinedKind()
		case *types.Alias:
			kind = n.AliasKind()
		case *types.TypeParam:
			kind = n.ParameterKind()
		}
		if o.Parent() != nil && o.Pkg() != nil && o.Parent() == o.Pkg().Scope() {
			owner = e.packages[o.Pkg()]
		}
	}
	id := semantics.EntityID(len(e.data.Entities) + 1)
	e.objects[o] = id
	entity := semantics.Entity{Kind: kind, Name: o.Name(), Owner: owner}
	...
}
```

So the analysed kind of an entity is chosen from what the compiler's object
is: a field, a predeclared type, a defined type, an alias or a type
parameter. Its identity is the compiler's object, and its owner is the
package whose scope declares it. Every fact about a type is attached to a
term and not to an entity (`export.go:117-119`):

```go
func (e *exporter) fact(id semantics.TermID, relation semantics.Key, value semantics.TermID) {
	e.data.Facts = append(e.data.Facts, semantics.Fact{Subject: semantics.Subject{Term: id}, Relation: relation, Term: value})
}
```

Two of the three bridges are not written in Go, and each spells the names
it writes in its own language. The TypeScript bridge is a script
(`internal/typescriptsem/compiler.cjs:8`, `:161`):

```javascript
const key = name => ({ Owner: "typescript", Name: name });
...
    graph.Entities.push({ Kind: key("package"), Name: specifier, Owner: 1, Declarations: [], Parameters: [] });
```

The Rust bridge is a Rust program that runs inside the compiler
(`internal/rustsem/helper/src/main.rs:232`):

```rust
            hir::ItemKind::Struct(..) => "struct",
```

**One fact, written and read.** The Go bridge writes the term of a slice
type and of a struct type (`internal/golangsem/export.go:233-243`; `e.typ`
exports a type and returns its term, `e.object` returns the entity of a
compiler object, and `e.expression` builds a provider expression whose data
is the JSON encoding of its second argument):

```go
case *types.Slice:
	term = e.expression(n.SliceOperator(), nil, e.typ(t.Elem()))
case *types.Struct:
	var fields []n.FieldInfo
	var operands []semantics.TermID
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		fields = append(fields, n.FieldInfo{Name: f.Name(), Tag: t.Tag(i), Embedded: f.Embedded(), Entity: e.object(f)})
		operands = append(operands, e.typ(f.Type()))
	}
	term = e.expression(n.StructOperator(), fields, operands...)
```

The typed reader in the vocabulary package reads it back
(`semantics/golang/go/model.go:20-39`, `:283-302`). `t.term()` returns the
graph's term of a type, and `t.Underlying()` follows the fact with the
relation `underlying`:

```go
type Type struct {
	owner *model
	id    semantics.TermID
}
type Field struct {
	Name     string
	Type     Type
	Site     ast.NodeRef
	Embedded bool
	Tag      string
}

// FieldInfo is the Go JSON field-metadata encoding. Operands hold the field types;
// Entity identifies the field's declaration in this same graph.
type FieldInfo struct {
	Name     string
	Tag      string
	Embedded bool
	Entity   semantics.EntityID
}
...
func (t Type) Fields() ([]Field, bool) {
	v := t.Underlying().term()
	if v.Operator != StructOperator() {
		return nil, false
	}
	var metadata []FieldInfo
	if json.Unmarshal(v.Data, &metadata) != nil || len(metadata) != len(v.Operands) {
		return nil, false
	}
	out := make([]Field, len(metadata))
	for i, f := range metadata {
		e, _ := t.owner.graph.Entity(f.Entity)
		var site ast.NodeRef
		if len(e.Declarations) > 0 {
			site = e.Declarations[0]
		}
		out[i] = Field{f.Name, Type{t.owner, v.Operands[i].Term}, site, f.Embedded, f.Tag}
	}
	return out, true
}
```

So the shape of a struct term (the operator `struct`, JSON entries for the
fields as data, one operand per field) is stated by the writer and by the
reader. The operator's name is declared. Its shape is not.

**What else the package holds.** A reader is opened from a receipt, the
value that the semantic contract package issues for one finished analysis
(section 2.5):

```go
// Open requires an admitted analysis receipt and never invokes a compiler.
func Open(result semantic.Result) (Model, error)
```

`Model` is the code's name for the handle of one opened graph. It is
unrelated to the intended-forms package. `Open` asks the semantic contract
package whether the receipt meets the requirement of this vocabulary
(`semantic.AdmitEvidence`) and then unpacks the graph from its envelope. The
package also builds the identity values of its resolution and analysis
contracts from types of three contract packages. That is why a semantic
vocabulary package imports three contract packages today. Go's has 855
lines: 99 are names, 40 contract identities, 400 the inputs that a compiler
run is configured with, with their validation and encoding, and 316 the
typed readers.

### 4.4 The data types of intent, Go's intended-forms package, and the target that reads it

**The data types of intent** are declared inside the generation contract
package, beside the planner (`contracts/generation/go/expression.go:9-93`,
abridged). `T` is the marker type of the target language:

On the migration's branch (section 2.7) three names of the next two blocks
are changed: the constructor `Native` is `Apply`, `NativeDeclaration` is
`TargetDeclaration` and `NativeParameter` is `TargetParameter`. The
expression kinds are five there: the kinds of a slot and of a helper are
gone with the types `Slot` and `HelperKey`, and a kind for a binder is
added. The blocks show the main branch.

```go
// [doc] main branch; the migration's branch renames Native to Apply and has five kinds
// Form is a target-owned native semantic constructor. Only its issuer can
// interpret its key. This is not an AST kind or a globally registered type kind.
type Form[T any] struct {
	owner *identity
	name  string
}
type FormIssuer[T any] struct{ owner *identity }

func NewFormIssuer[T any]() FormIssuer[T] { return FormIssuer[T]{fresh()} }
func (i FormIssuer[T]) Form(name string) Form[T] {
	if i.owner == nil || name == "" {
		return Form[T]{}
	}
	return Form[T]{i.owner, name}
}
func (i FormIssuer[T]) Read(f Form[T]) (string, bool) {
	return f.name, i.owner != nil && i.owner == f.owner && f.name != ""
}
...
type ExpressionKind uint8

const (
	NativeExpression ExpressionKind = iota + 1
	CapturedExpression
	ParameterExpression
	ArgumentExpression
	HelperExpression
	DeclarationExpression
)
...
// Expression is an immutable template. Parameter and Argument indices only have
// meaning in their declaration's single generic convention. DeclarationExpression
// is issued by the planner after it authenticates the source callee.
type Expression[T any] struct {
	origin      ast.NodeRef
	kind        ExpressionKind
	form        Form[T]
	data        []string
	args        []Expression[T]
	slot        Slot
	index       int
	helper      HelperKey
	declaration uint32
}
...
func Native[T any](form Form[T], data []string, args ...Expression[T]) Expression[T] {
	return Expression[T]{kind: NativeExpression, form: form, data: slices.Clone(data), args: slices.Clone(args)}
}
...
func (e Expression[T]) Form() Form[T]              { return e.form }
func (e Expression[T]) Data() []string             { return slices.Clone(e.data) }
func (e Expression[T]) Arguments() []Expression[T] { return slices.Clone(e.args) }
```

`identity` is a private type whose pointer is unique per issuer, so only the
issuer that made a form can read it. `Slot` and `HelperKey` are two more
issued identities of the same construction: a slot stands for the
representation of another use of a type, and a helper key names a
declaration that a mapping demands. An expression also has an accessor
`Origin()` for the source node it came from. Of the six expression kinds
only the first holds a form. The other five are the kernel's own: a slot, a
type parameter, a generic argument, a helper reference and a reference to a
planned declaration. So type parameters, references and applications are no
forms of any language. They belong to the shared contract package.

An intended declaration and the interface of a target
(`contracts/generation/go/native.go:13-24`, `contracts.go:13-42`, abridged):

```go
// [doc] main branch; the migration's branch says TargetDeclaration and TargetParameter
// NativeDeclaration is target intent after source mapping and substitution.
type NativeDeclaration[T any] struct {
	Name       string
	Form       Form[T]
	Origin     ast.NodeRef
	Parameters []NativeParameter
	Body       Expression[T]
	Values     []IssuedValue
}
type NativeParameter struct {
	Name   string
	Origin ast.NodeRef
}

var (
	ErrInvalid     = errors.New("invalid generation description")
	ErrUnsupported = errors.New("generation capability unavailable")
	ErrScope       = errors.New("incompatible generation scope or association")
	// [doc] and ErrLimit, ErrIncomplete, ErrAmbiguous, ErrRejected
)

// Target validates native intent during planning. Build receives only the closed
// native program; it has no source reader, selectors or mapping callbacks.
// Target implementations are trusted native semantic compilers.
type Target[T any] interface {
	CheckExpression(context.Context, Expression[T]) error
	CheckDeclaration(context.Context, NativeDeclaration[T]) error
	CheckProgram(context.Context, Program[T]) error
	Build(context.Context, BuildInput[T]) (Assembly, error)
}
```

`IssuedValue` is an opaque statement that a declaration's owner attaches and
only its issuer can read. `Program[T]` is the closed list of planned
declarations. `BuildInput[T]` holds the program, the language descriptor,
the identifier of the output tree and limits. `Assembly` is the built tree
with the nodes at which its declarations sit. Two errors matter below:
`ErrInvalid` says that a description is malformed, and `ErrUnsupported` says
that a capability is not available. The code's comments call a form a
"native semantic constructor" and a target a "trusted native semantic
compiler"; the accepted pages say that intent is not analysed meaning (fixed
point 8).

**Go's intended-forms package, complete**
(`transform/syntax/golang/model/go/model.go:1-112`). The code names the
package `model`. Here `g` is the generation contract package. The comments
marked `[doc]` classify each condition as **WF** (well-formedness: what a
form is) or **COV** (a coverage answer: what today's generation supports).
The classification is our derivation. It matters for sections 6.1 and 6.5.

```go
// Package model describes intended native types before AST construction.
package model

import (
	"context"
	"strconv"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
)

var forms = g.NewFormIssuer[s.Language]()

// The names of the forms this model issues, each declared here once. Form
// takes one of them and Name returns one. The first five are the expression
// forms that the constructors below build. FormAlias and FormDefined are the
// declaration forms.
const (
	FormPrimitive = "primitive"
	FormSlice     = "slice"
	FormField     = "field"
	FormRecord    = "record"
	FormExternal  = "external"
	FormAlias     = "alias"
	FormDefined   = "defined"
)

// The names a primitive form carries, each declared here once. Primitive
// takes one of them.
const (
	PrimitiveString = "string"
	PrimitiveBool   = "bool"
)

func Form(name string) g.Form[s.Language]      { return forms.Form(name) }
func Name(f g.Form[s.Language]) (string, bool) { return forms.Read(f) }
func Primitive(name string) g.Expression[s.Language] {
	return g.Native(Form(FormPrimitive), []string{name})
}
func Collection(element g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormSlice), nil, element)
}
func Field(name string, optional bool, ty g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormField), []string{name, strconv.FormatBool(optional)}, ty)
}
func Record(fields ...g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormRecord), nil, fields...)
}
func External(module, name string, args ...g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormExternal), []string{module, name}, args...)
}
func CheckExpression(ctx context.Context, e g.Expression[s.Language]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	name, ok := Name(e.Form())
	if !ok {
		return g.ErrScope // [doc] the form was not issued by this model
	}
	d, a := e.Data(), e.Arguments()
	switch name {
	case FormPrimitive:
		// [doc] WF: one datum, no argument.
		// [doc] COV, returned as invalid: the datum is one of two names.
		if len(d) != 1 || len(a) != 0 || (d[0] != PrimitiveString && d[0] != PrimitiveBool) {
			return g.ErrInvalid
		}
	case FormSlice:
		// [doc] WF: no data, one argument.
		if len(d) != 0 || len(a) != 1 {
			return g.ErrInvalid
		}
	case FormField:
		// [doc] WF: a name, a flag, one argument.
		if len(d) != 2 || d[0] == "" || (d[1] != "true" && d[1] != "false") || len(a) != 1 {
			return g.ErrInvalid
		}
	case FormRecord:
		// [doc] WF: no data; every argument is a field of this model; no name twice.
		if len(d) != 0 {
			return g.ErrInvalid
		}
		seen := map[string]bool{}
		for _, f := range a {
			n, ok := Name(f.Form())
			v := f.Data()
			if !ok || n != FormField || len(v) != 2 || seen[v[0]] {
				return g.ErrInvalid
			}
			seen[v[0]] = true
		}
	case FormExternal:
		return g.ErrUnsupported // [doc] COV: Go refuses every external type today
	default:
		return g.ErrUnsupported // [doc] COV: a form name this check does not know
	}
	return nil
}
func CheckDeclaration(ctx context.Context, d g.NativeDeclaration[s.Language]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	n, ok := Name(d.Form)
	if !ok {
		return g.ErrScope
	}
	// [doc] WF: a declaration has a name and one of the two declaration forms.
	if d.Name == "" || (n != FormAlias && n != FormDefined) {
		return g.ErrInvalid
	}
	if n == FormDefined {
		// [doc] A restriction of today's generation, returned as invalid: Go
		// [doc] itself admits a defined type over any type.
		body, ok := Name(d.Body.Form())
		if !ok || body != FormRecord {
			return g.ErrInvalid
		}
	}
	return nil
}
```

Facts about the four intended-forms packages:

- They are near copies. Each has 112 to 127 lines. They differ in the name
  of the collection form (`slice` in Go, `array` in TypeScript, `sequence` in
  Rust and C++), in the name of the nominal declaration form (`defined`,
  `interface`, `struct`, `struct`) and in the external form.
- No package is shared between them (the paragraph under fixed point 27).
- `Form(name)` issues a form of any name. The function does not check the
  name against the constants.
- The Go package exports a constructor, `External`, whose result its own
  check always refuses.
- No package has a form for a pointer, a map, an integer, a union or an
  optional type. The syntax vocabulary package of every language has node
  kinds for such types (section 4.1).
- Rust keeps one coverage answer outside this package: its intended-forms
  package checks only the shape of an external form, and the Rust target
  refuses every valid crate name but that of the standard library (the function `Dependency`
  of section 4.2).

**Who writes forms.** A mapping. The mapping for horos text into Go gives its
result as a constant: in effect `model.Primitive(model.PrimitiveString)`.
The record mapping uses a callback: it reads the fields, asks the kernel for
the representation of each field's type, and returns
`model.Record(model.Field(name, optional, type), ...)`. A mapping imports
the intended-forms package for its constructors and its names. No mapping
calls one of its checks.

**Who reads forms.** The target of the same language. Go's target forwards
the two checks to the package above, and it switches over the form names
when it builds (`transform/syntax/golang/target/go/provider.go:111-128`).
`m` is the shared tree builder: `m.Node(kind, attributes, origin)` makes a
node value, `m.Attributes` builds attributes from pairs of name and text,
and `m.Child` pairs a child role with a fragment. `b.Node` allocates one
node with its children and returns a handle to it, a `m.Fragment`.
`b.predeclared` writes a `named-type` node and records the predeclared
symbol of section 4.5 as its destination. `edges` pairs each fragment of a
list with one role:

```go
func (b *Backend) Expression(e g.Expression[s.Language], args []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	if err := model.CheckExpression(context.Background(), e); err != nil {
		return m.Fragment[s.Language]{}, err
	}
	name, _ := model.Name(e.Form())
	data, origin := e.Data(), e.Origin()
	switch name {
	case model.FormPrimitive:
		return b.predeclared(data[0], origin)
	case model.FormSlice:
		return b.Node(m.Node(s.KindSliceType, nil, origin), m.Child[s.Language]{Role: s.RoleElement, Fragment: args[0]})
	case model.FormField:
		return b.Node(m.Node(s.KindField, m.Attributes("name", data[0]), origin), m.Child[s.Language]{Role: s.RoleType, Fragment: args[0]})
	case model.FormRecord:
		return b.Node(m.Node(s.KindStructType, nil, origin), edges(s.RoleField, args)...)
	}
	return m.Fragment[s.Language]{}, g.ErrUnsupported
}
```

Three facts are visible here. The target is the one production caller of
the checks, and it calls no constructor. The target is where a form meets a
node kind: `slice` becomes `slice-type`, `record` becomes `struct-type`,
`field` becomes `field`. And the Go target reads `data[0]` of a field form
and ignores `data[1]`, the optional flag: Go has no optional field, and the
record mapping wraps the type of an optional field in a helper type instead.

### 4.5 A symbol package

A symbol is what a name use in generated output refers to when the referent
is outside the generated declarations: a predeclared type, an item of a
standard library, an export of a package. The generation contract package
supplies an issuer of the same construction as the form issuer. The value
it issues, an `ExternalSymbol`, holds a key that the generation kernel and
the Files formatters carry and never read. An `ExternalUse` is such a
symbol at one occurrence in a tree.

Go's symbol package (`transform/syntax/golang/symbol/go/symbol.go:15-42`,
one comment shortened as marked). Two further functions, `ExternalUse` and
`ReadExternal`, record a symbol at a node and read it back:

```go
// Symbol is a destination outside the declarations of a generated output. It
// states its category: a predeclared type is never a package symbol with an
// empty package. The zero Symbol is no destination. Symbols are comparable.
type Symbol struct{ category, name string }

// The one category of today. A type of another package is a further category
// and needs that package's semantics.
const predeclared = "predeclared"

// Predeclared identifies the type that name denotes in Go's universe scope.
// ... Only a basic type is issued today. error, comparable
// and any are refused with generation.ErrUnsupported: the Go bridge exports
// them with the kinds of a declaration, without a marker that tells them from
// one, so their destination cannot be established yet.
func Predeclared(name string) (Symbol, error) {
	object, ok := types.Universe.Lookup(name).(*types.TypeName)
	if !ok {
		return Symbol{}, g.ErrInvalid
	}
	if _, basic := object.Type().(*types.Basic); !basic {
		return Symbol{}, g.ErrUnsupported
	}
	return Symbol{predeclared, name}, nil
}
...
```

The package states which names are symbols by asking the Go standard library
(`types.Universe`). It holds no table. Its second refusal is a coverage
answer that depends on what the Go bridge can tell apart.

One Files formatter exists, TypeScript's, and it is where an intended
vocabulary and an analysed vocabulary of the same language meet today. It
reads each external use of the sealed output back through TypeScript's
symbol package, whose symbol is a module and a name. After it has placed
the declarations and had the placed program analysed, it reads the analysed
graph in the words of the semantic vocabulary package: a binding to an
entity with an analysed kind (`package-type`, `builtin`) and an owner. It
translates both into a private type of its own and compares. The accepted
specification of destinations requires that comparison to be one of
identities (fixed point 18). The same page says that nothing compares the
role under which a symbol is recorded with a role of the bridge, and the
boundaries page says that a symbol package does not import the semantic
vocabulary package.

### 4.6 The node patterns of horos kinds

horos has a second level above types. "A type describes values. A kind
describes types or type constructors." (`docs/kinds.md`). We say **horos
kind** for it, to keep it apart from a node kind. A horos kind is written in
a notation for the shape of a tree. bittype parses, checks and analyses this
notation as part of horos. Two examples and the grammar of child sequences,
from `docs/horos-syntax.md`:

```text
kind List<K> = node("list-of", {}, item(K));
kind ClosedRecord<K> = node("record", {"mode": "closed"}, star(item(K)));
```

```text
Q = epsilon
  | item(K)
  | concat(Q, Q, ...)
  | alternative(Q, Q, ...)
  | star(Q)
```

How to read them:

- `node(c, m, Q)` matches a tree node that was built by the constructor `c`,
  whose metadata meets the constraint `m`, and whose ordered children match
  the sequence pattern `Q`.
- `item(K)` matches one child that is of the horos kind `K`. `epsilon`
  matches no child. `star(Q)` matches any finite repetition of `Q`.
- `K` in `List<K>` is a parameter of the horos kind.
- A horos kind can also be a union `|`, an intersection `&`, or recursive, with
  greatest fixed points, written `nu { X = ...; } in X`.

The page says: "A node pattern matches the entire ordered list of observed
type children."

**What the notation describes today.** It describes horos type descriptors
and nothing else. A type descriptor is the structural description of a
horos type: a tree whose nodes are built from 22 built-in constructors such
as `list-of`, `record` and `scalar`. "Kinds classify **checked structural
type descriptions**." (`docs/kinds.md`). The first argument of `node` is one
of those 22 constructors. Metadata constraints exist for two of them:
`selector` on `scalar` and `mode` on `record`. The page puts outside its
present scope: "Public kind negation, arbitrary predicate callbacks and
subtree-equality constraints".

**Compared with the grammar table of section 4.2** (our derivation, from
reading the two notations side by side):

| | Grammar table | horos node pattern |
| --- | --- | --- |
| Node identity | Any node kind of the language | One of 22 descriptor constructors of horos |
| Attributes | Named, required or optional, each with a predicate function | Metadata constraints with string values, for two constructors |
| Children | Named child roles in a fixed order, each with a count range and a flat list of node kinds | An unnamed ordered sequence with concatenation, alternative and repetition |
| Recursion | By naming node kinds in the lists | Explicit fixed points |
| Parameters | None | Parameters of a horos kind |
| References | A flag: one reference role, `target` | None |
| Interpreter | One Go function of 67 lines | The horos library, which owns the meaning of horos kinds. An earlier consultation (research 0005) recorded the choice to "implement proper-kind inclusion using finite cuts of positive greatest-fixed-point hedge automata", that is, to decide whether one horos kind is contained in another |

No code and no page on the main branch applies a horos kind to a syntax
tree. No language's tree is described in this notation, and that includes
the tree of horos itself: its own shape is a Go table and 11 rules in code,
like the others.

### 4.7 The one generator

On the main branch no Go file is generated. No file holds a `go:generate`
directive or a "Code generated" header. One table, the character ranges of
section 4.2, is produced by a manual procedure. Three documentation pages
are generated from the source by a tool.

On the migration's branch one generator exists, `tools/payload-schema`. Its
input is a horos module that declares the type of a mapping's decision
record (section 2.7; the code calls a decision record a payload). An
example, complete
(`transform/syntax/horos-to-golang/vocabulary/declarations/go/payload.horos`
on the branch):

```horos
module bittype.payload.horos_to_golang_declarations;

type Payload = { name: text; nominal: bool; parameters: list<text>; };
```

From it the tool makes three files beside it: the canonical binary module,
a Go file with the module's pin, and a Go file with the type
(`payload_types.go`, complete):

```go
// Code generated by bittype from payload.horos. DO NOT EDIT.

package declarations

type Payload struct {
	Name       string
	Nominal    bool
	Parameters []string
}
```

The Go type is produced by bittype's own pipeline. The comment on the
function that makes it says: "the schema source goes through the horos
frontend and through generation with bittype's own structural Go
contributions", which are the Go mappings of section 2.2. A schema may use
what those mappings cover: text, booleans, lists, closed records with
required and optional fields, and named declarations. The branch holds 24
such schemas. A test regenerates each file and fails when the committed file
differs.

So the project already declares one sort of value type in horos and
generates its Go type: the decision record of a mapping. It uses horos types
for that and no horos kinds, and it describes a value and no tree.

## 5. Constraints

This section lists the owner's direction first, and then owner rules and
accepted designs. Two rows have another standing and say so. Section 5.2
says what is not decided, and section 5.3 lists five places where pages
contradict each other or the code. No fact of the present code is a
constraint. One row is a fact of the Go language (fixed point 21).

### 5.1 Fixed points and their standing

The Standing column uses the seven standings of the Question section. A rule given in
the page's wording means that the owner's own words are not on file and
the page that records the rule is quoted.

**The direction.**

| # | Statement | Standing | Shown in |
| --- | --- | --- | --- |
| DIR | "So the root split to me would be between model and ops. model must not contain any ops and must not import from ops either. ops work on the model" | The owner's direction (message 8, with the two roots and rule 1 of message 7) | Question, 6.1 |

The direction changes two older rules of the same owner at the root (fixed
point 24): under it the stage-by-language matrix lives below the operations
root. We derive no more than that. What stays of those two rules is asked in
question 4. The names below the roots in message 7 are a sketch (Question).

**About the answer's method.** Section 2.1 quotes each of these in full.

| # | Statement | Standing | Shown in |
| --- | --- | --- | --- |
| 1 | AC-09. The existing code and the effort of changing it have no weight in a design choice. Estimates may guide the order of delivery only | Owner rule of 8 November 2026 | 2.1 |
| 2 | AC-10. A misplaced or wrongly named responsibility moves to its owner first, as its own change. Later steps add and do not move, rename or remove what that change delivered | Owner rule of 11 November 2026 | 2.1 |
| 3 | Wrong grouping is corrected before wrong coupling. A dependency check enforces a boundary that reasoning has decided and does not replace that reasoning | Owner rule of 11 November 2026, in the platform page's wording | 2.1 |
| 4 | AC-11. A mechanism that needs no knowledge of one language is implemented once, in the agnostic layer. Each language supplies its own types, meanings and policies | Owner rule of 13 November 2026 | 2.1 |
| 5 | An extension needs a concrete case that the existing contracts cannot express or uphold | Owner rule of 7 October 2026, in the platform page's wording | 2.1 |
| 6 | Unification needs stronger evidence than overlapping fields or similar code. Identity, lifetime, authority, allowed operations, guarantees and failure behavior are compared, and differences are represented and not erased | Accepted design (a platform page; not the owner's own words) | 2.1 |
| 7 | Who owns a meaning is not determined by a directory, a package name or a generic type parameter | Owner rule of 6 November 2026, in the platform page's wording | 2.1 |

**About intent.**

| # | Statement | Standing | Shown in |
| --- | --- | --- | --- |
| 8 | Intent is not analysed meaning. "Native intent is not an analyzed `Semantics[T]`: target analysis establishes facts about the actual constructed/organized program." "The planned native program is intent, not a fabricated semantic-analysis result. Only analysis of the realized target program establishes native semantic facts." "Native semantic evidence belongs to the actual native capture. A generation plan or successful syntax check is not that evidence." | Accepted design (`docs/generation.md`, `docs/generation-planner.md`, `docs/type-system-architecture.md`) | 2.2 |
| 9 | What generation's sealed output guarantees before it is placed into files: "The five properties are adopted. A destination for every name use is required by each language's generated-Files formatter, which refuses a use without one as TypeScript does, and is asserted by a test per target over sealed outputs". The row adapts the recommendation that unplaced output needs no analysed meaning of its own | Accepted design (research 0016 on the placement of generated declarations, row 19, NEEDS ADAPTATION) | Quoted here only |
| 10 | "Native expressions contain target-owned forms, immutable scalar metadata and child expressions. The shared interpreter understands substitution and references, not the meanings of these forms. There is no AST kind, source string or arbitrary payload field in an expression." | Accepted design (`docs/generation-spi.md`; SPI is the interface that mappings and targets implement) | 4.4 |
| 11 | "Target implementations are trusted interpreters of their own native forms" | Accepted design (`docs/generation-planner.md`) | 4.4 |
| 12 | "bittype's native vocabularies stay types-only." The vocabularies of the four programming languages describe types and no functions, values or statements | Accepted design (`docs/binding-generation.md`, from research 0014 on generated conversion code) | Quoted here only |
| 13 | What two operations do not establish. Of generation: "It establishes neither native kind enforcement nor a proof that the chosen representation is correct." Of parsing: "success establishes neither reference resolution nor semantic validity" | Accepted design (`docs/generation.md`, `docs/parsing-formatting.md`) | 6.1 |
| 14 | "The kernel owns the guarantees shared across languages. It invokes a supplied interpretation when a result depends on a language or encoding. A grammar, for example, can supply declarative rules." | Accepted design (`docs/abstraction-levels.md`, which calls itself a design note) | 6.1 |

**About a shared model of types.**

| # | Statement | Standing | Shown in |
| --- | --- | --- | --- |
| 15 | The semantic graph "has no scope graph, substitution map, universal type algebra or conditional proof system". "A provider owns lookup, semantic rules, vocabulary, coverage and the interpretation of its payloads. The kernel validates graph integrity; it cannot prove a provider's claims or make native meanings equivalent to Horos." A provider here is the language-specific code that writes and reads a graph | Accepted design (`docs/semantics.md`; research 0002 on the semantic kernel, row 1, HOLDS: "Share established facts with provider-owned vocabulary, separate entity bodies and parameter ownership") | 2.3, 4.3 |
| 16 | "Keep the SPI experimental and avoid a richer generic semantic model." The expert's words behind the row, about the mechanism by which mappings then built target syntax: it "should not become a new language-neutral type calculus" | Accepted design (research 0007 on generator plugins, row 21, HOLDS) | Quoted here only |
| 17 | A model of a type system in the sense of terms, judgments and rules is content of logos packages: "These are requirements for language integrations and Logos semantic packages, not fields of a new shared Bittype model or implemented Go struct." Row 1, HOLDS: "Consume Logos semantic packages; no Bittype-owned type-system foundation" | Accepted design (`docs/type-system-model.md`; research 0010 on the architecture of type-system models, which was about verification) | Quoted here only |
| 18 | A destination is an identity. "A destination is a declaration's identity and never structural equality". "An external symbol names the selected public export and not the canonical type underneath it" | Accepted design (`docs/destinations.md`, specification of 13 November 2026, from research 0016) | 4.5 |

Fixed points 15 to 17 each rule out one thing: a universal type algebra in
the graph, a language-neutral calculus in generation, and a type-system
framework owned by bittype. None of them speaks about one language's own
vocabulary of its own types. That is our reading, and question 2 rests on
it.

**About vocabulary packages and the layout.**

| # | Statement | Standing | Shown in |
| --- | --- | --- | --- |
| 19 | A vocabulary is "The names two parties must agree on and nothing else: syntax kinds, semantic keys, profile and contract identities, intended native forms. A vocabulary package imports no contribution." | Accepted design (`docs/package-boundaries.md`, state of 12 November 2026). What is implemented holds more than names (section 5.3) | 2.6 |
| 20 | "**One language.** A package that names a language imports no package of another language. A mapping names two." | Accepted design (`LAYOUT.md`). The direction does not speak of it | 2.6 |
| 21 | Go's rule for a directory named `internal`: a package below it can be imported only by code in the tree rooted at the parent of that directory. Today `internal/` is at the module root, so every package can import the helpers below it. A helper below `ops/internal/` could not be imported from `model/`, and the reverse | A fact of the Go language. It is the one row outside the seven standings | 6.2 |
| 22 | "One public root Go module contains independently owned packages". "Do not create a language-first packages tree, additional Go modules, compatibility aliases or placeholder source directories." Here: no additional Go modules | Owner rule of 31 October 2026, in the layout page's wording. The direction does not speak of it. Fixed unless you see that it contradicts the direction. Please say what it costs for enforcing rule 1 | 2.4 |
| 23 | The first item of the sentence quoted in row 22: no language-first packages tree | Owner rule of the layout, in the layout page's wording. The direction's change at the root puts it in question for the model root, and the order below a root is asked (question 4) | 2.4 |
| 24 | "organize by processing stage and enforce a uniform structure across languages"; the path pattern `<stage>/<dimension>/<scope>/<component>/go` | Owner rules of 31 October and 6 November 2026, in the layout page's wording. The direction changes them at the root. What stays of them below the operations root is asked and is no limit on the answer | 2.4 |
| 25 | "A checker gets no concrete tree, a formatter no lexer, and no Syntax stage a semantic contract." | Accepted design (`LAYOUT.md`, role table). The direction does not speak of it | 2.6 |
| 26 | "Generic carriers, contracts, engine, printer and native syntax stages have no external Horos/Ontos imports" | Accepted design (`CHARTER.md`, the repository's founding page) | 6.4 |
| 27 | The design of the contract migration lists, among the shared packages that stay unchanged when a materially different contribution is added, "the native model and target of each language" | Decision of the research owner (`docs/contribution-contract.md`, 13 November 2026, final only when its gate is met, section 2.7) | 2.7 |

**One thing that is implemented today and open: the four intended-forms
packages share no package.** A change-log entry of 12 November 2026 records
it: "Each model owns its own forms, a form can be read only by its issuer,
and no package compares a form name of one language with another's or with a
literal." The four packages are near copies (section 4.4). Fixed point 4,
the owner rule AC-11, is one day younger. It says that a mechanism which
needs no knowledge of one language is implemented once. Whether the four
may or must share a mechanism is therefore open.

Out of scope: the design of the contract migration itself, the choice of a
native representation for any horos type, the design of the Files
formatters, and verification.

### 5.2 What is not decided

Standing: open. The study of the design pages read the layout page, the
charter, the contributor instructions, the change log, the design pages and
the result tables of the earlier consultations for each of these questions.
It found no accepted sentence that answers any of them. It lists what it
did not read, and the statements below are bounded by that.

1. **Whether the type-level notions of one language should have one set of
   names** in which generation states intent and analysis states what it
   established. Every accepted sentence (fixed point 8) separates a planned
   program from an analysed, placed program. None speaks about the names of
   the notions. No consultation asked. Question 2 asks.
2. **What the top-level directories `syntax/` and `semantics/` are.** No
   sentence says whether they are stages, dimensions, levels or homes of
   vocabulary packages. Put together from several sentences, each holds the
   per-language vocabulary of one abstract representation, the syntax tree
   and the semantic graph. That description is the study's.
3. **Whether a stage other than generation may import an intended-forms
   package.** No role lists it except the target and the mapping. No
   sentence forbids or permits it for a Files formatter, for verification
   or for a compiler bridge.
4. **Whether verification laws and generation mappings should speak about Go
   types in one vocabulary.** A verification law is a stated property that
   a Go type must have to represent a horos type. Today the laws use the
   semantic vocabulary and the mappings use the intended forms.
5. **Whether the source side of generation is meant to have a counterpart
   of the intended-forms package.** horos has a source reader and no such
   package (section 3.2). Question 2 asks what a language presents on both
   sides.
6. **Where a model's own consistency ends and an operation begins.** The
   direction says that a model contains no operation. No page says what
   counts as one. Section 6.1 gives our derivation, and question 1 asks.
7. **What a model is for a language without analysis.** C++ has no semantic
   vocabulary package and no symbol package.
8. **Whether the sameness of the four intended-forms packages is intended**
   (the paragraph under fixed point 27).

### 5.3 Five contradictions that the answer must not rest on

Standing: each is between two accepted pages, or between an accepted page
and what is implemented. We intend to correct each whatever the answer.
None is evidence for an option.

1. **The definition of a vocabulary against what three of the four sorts of
   vocabulary package hold.** The definition (fixed point 19) says "and
   nothing else". The layout page describes the intended-forms package as
   "native semantic forms and validation" and the Go semantic vocabulary
   package as "Go vocabulary, contract identities and immutable native type
   readers". The boundaries page itself admits a method into a vocabulary
   package with a criterion: "It reads no file and starts no process, so a
   vocabulary can hold it." In the code only the syntax vocabulary packages
   hold names and nothing else.
2. **The target's role.** The layout page gives a target "its own syntax,
   model and lexical helper". Three of the four targets also import their
   symbol package, through an extra entry of the architecture test.
3. **What a vocabulary package may import.** The pages say "Foundations and
   its own language's vocabulary". The architecture test permits each sort
   of vocabulary package only the syntax vocabulary package of its
   language. An intended-forms package may not import the semantic
   vocabulary package or the symbol package, and the reverse.
4. **Who supplies the intended forms.** One page says: "Libraries supply
   declarations, native models and mapping descriptions under their own
   ownership." A library there is a separately owned code base that
   contributes mappings. The code and two later documents treat the
   intended forms as the target's closed set: the issuer is private, a form
   of another issuer is refused, an unknown form name is unsupported, and
   the target switches over the names.
5. **The intended forms as fixed and shared, against sorts of source type
   that need new forms.** Fixed point 27 keeps the intended-forms package
   unchanged when a contribution is added. An earlier consultation recorded
   that for many sorts of source type "Only the models and mapping libraries
   are missing". A mapping that needs a form the package lacks, such as a
   map or an integer, must edit the package and the target of that
   language. No page says whether that is an edit of a shared package in
   the sense of the gate.

## 6. What we have considered

Standing of this section: the owner's direction is fixed (Question, fixed
point DIR). Everything else here is our derivation, or a fact of the code
that is marked as such, and is open. No option is ruled out because the
present code would have to change. The document does not choose among the
options. The labels are unique, so that an answer can cite them: P for
package placement, U for unification, D for declaration.

### 6.1 What the two rules reach, by our test; consistency and correctness

**A test, applied to every package.** To see what the direction reaches, the
study of the code applied one stated test to all 129 production packages. It
was made before the owner's eighth message. The test is ours, not the
owner's:

- A package is on the *model* side if it declares what exists (generic data
  types, vocabularies, forms, identities, value types with their
  well-formedness) and performs no processing of an input into an output.
- It is an *operation* if it parses, checks, formats, prints, resolves,
  analyzes, transforms, verifies or composes.
- It holds *both* if it holds a declaration and processing.
- It is *undecided* if the test gives no answer.

Validation at construction time was counted as well-formedness. Anything
that takes one artifact and returns another was counted as processing: a
tree from text, an encoding from an environment, a verdict on a form.
Because the test lists "composes", it files the shared mechanisms (the
engine, the generation registry, the planner) under operation. That is a
result of our test and is asked (case (f) of section 6.2).

**Result of our test.**

| Result | Packages | Lines | Which |
| --- | --- | --- | --- |
| On the model side, or both with a model-side core | 22 | 5,428 | The three tree packages, `files/go`, `source/go`, `semantics/go`, five syntax vocabulary packages, four semantic vocabulary packages, four intended-forms packages, three symbol packages |
| Contract packages | 13 | 5,030 | Section 2.5. By the test `language` and `text` are on the model side, `pipeline` is an interface only, and the other ten hold both |
| Undecided | 10 | 1,498 | Below |
| Operation | 84 | 23,987 | The stage providers and query services, the shared printers, the generation kernel packages, the source reader, four targets, 21 mappings and three shared mapping parts, the verification packages, the tools, the engine, and nine private helpers (scanning, parser state, tree construction, the three bridges, the two horos adapters and the import helper of TypeScript's Files formatter) |

**What rule 2 (no containment) reaches.** Twelve of the 22 model-side
packages hold both under the test, because each holds code whose input and
output are different artifacts:

| Package | The model-side part | The processing part |
| --- | --- | --- |
| `trees/representation/go` (the envelope package) | The language descriptor and the opaque envelopes | Packing and unpacking: a deep copy and a structure check, which it calls in a `check/` directory |
| `source/go` | Units, groups, relations, the publication | A publisher of 195 lines that admits a publication under a policy with two callbacks into operations (section 2.3) |
| `semantics/{golang,typescript,rust}/go` | Names, profiles, contract identities, the types of a compiler run's inputs | The typed readers over a graph, and the encoding of those inputs into settings and back |
| The four intended-forms packages | Form names and constructors | Two checks that return a verdict on a form (section 4.4). The migration's branch adds a third check to each (section 2.7) |
| The three symbol packages | A symbol type and its categories | Admission of a name (for Go by asking the standard library), and the encoding of a symbol into an issued key |

Rule 2 also reaches the contract packages from the other side. Ten of the
thirteen hold value types that the test counts as model-side, beside the
interface and the admission code of one operation (section 2.5). If those
value types are models, each of the ten is a package in which a model and an
operation sit together.

**The ten packages that the test does not decide.**

| Package | Why undecided |
| --- | --- |
| `check/syntax/shared/ast/go`, `check/syntax/shared/cst/go` | By their directory they are a check stage. By their content they are the well-formedness of the generic tree, which the test puts on the model side. The envelope package imports both |
| `internal/tree` | A read helper and a deep copy over the generic tree. It declares nothing and is no stage. The envelope package imports it |
| `internal/grammar` | It declares what a grammar table is (model-side) and holds the interpreter of 67 lines (processing) |
| `internal/declsyntax` | The horos grammar table, a literal, and the horos rules in code |
| `internal/cpp`, `internal/rust`, `internal/typescript`, `internal/horos` | Spelling rules. A predicate over strings states which strings exist, which is model-side, and it is written as code. Three of them also decode string literals, which is processing. `internal/rust` and `internal/typescript` are imported from both sides (section 4.2) |
| `transform/syntax/horos/mapping/go` | Two profile identities and one function type |

**What rule 1 (no import) reaches.** The 22 model-side packages have 33
direct imports of an operation package, a contract package or a helper. Four
of the 22 import only model-side packages. None imports a stage provider, a
query service, a target, a mapping, the engine or a bridge. Whether one of
the 33 is forbidden depends on where the imported package is placed, and the
direction leaves that open.

| Sort of import | Importing package | Imported package | Identifiers used | Imports |
| --- | --- | --- | --- | --- |
| A vocabulary package uses a contract package's value types | The five syntax vocabulary packages | `contracts/language/go` | `ID`, `Provider`: the two methods that state a language's identifier and name | 5 |
| | The four semantic vocabulary packages | `contracts/semantic/go` | `Setting`, `EvidenceRequirement`, `Result`, `Snapshot`, `Requirement`, `Obligation`, `Contract`, the verdict values `Holds` and `Undetermined`, default limits, three error values, and two functions: `AdmitEvidence`, `CaptureSettings` | 4 |
| | The four semantic vocabulary packages | `contracts/resolution/go`, `contracts/analysis/go` | `Contract`, `Requirement`, `Kind`: the types of the contract identities | 8 |
| | The four intended-forms packages | `contracts/generation/go` | `NewFormIssuer`, `Form`, `Expression`, `Native`, `NativeDeclaration`, and the errors `ErrInvalid`, `ErrScope`, `ErrUnsupported` | 4 |
| | The three symbol packages | `contracts/generation/go` | `NewSymbolIssuer`, `ExternalUse`, the same errors | 3 |
| A carrier uses a contract package's value types | `source/go` | `contracts/semantic/go` | `Snapshot`, `Result`, `EvidenceRequirement`, `Limits`, default limits, two error values, and the functions `AdmitEvidence`, `ValidateEvidenceRequirements` | 1 |
| | `trees/representation/go` | `contracts/language/go` | `ID`, `Provider` | 1 |
| The tree's own well-formedness lives in a check directory | `trees/representation/go` | `check/syntax/shared/ast/go`, `check/syntax/shared/cst/go` | `Structure`, `Check` | 2 |
| | `trees/representation/go` | `internal/tree` | `Clone` | 1 |
| Lexical rules needed on both sides | `semantics/rust/go`, the Rust symbol package | `internal/rust` | `CrateName`, `Identifier` | 2 |
| | `semantics/typescript/go`, the TypeScript symbol package | `internal/typescript` | `PackageSpecifier`, `TypeName` | 2 |
| | | | | **33** |

For a vocabulary package to import no contract package, three things would
have to be declared elsewhere: the data types of intent and the symbol
issuer (section 4.4), which sit beside the planner today; the receipt and
the snapshot, which only the contract package's admission code creates
(section 2.5); and the two declarations of a language's identity. In the
other direction the contract packages import the carriers, and `source/go`
imports a contract package that three other contract packages need in order
to import `source/go`: a cycle of directories, though not of packages.

**Consistency and correctness.** Standing of the rest of this subsection:
our derivation, except what is marked as implemented today or as an
accepted design. Our test left one thing unassigned: code that sits with a
description and decides whether a value is a well-formed instance of it.
The owner answered that with two questions (messages 9 and 10): whether we
were talking about the correctness of a model, or about its consistency. We
made the following of them. Question 1 asks whether it is the right axis.

- **Consistency:** a value agrees with its own model and with itself. It can
  be decided from the value and the model's declaration alone.
- **Correctness:** a value agrees with something outside itself: a source, a
  specification, a compiler's judgment. It is a relation between two things,
  and it needs evidence.

Four levels, as our hypothesis:

| Level | What is asked of a value | Decided today by |
| --- | --- | --- |
| 1 | The structure of the generic data type: identities in range, nothing dangling | A generic check |
| 2 | Conformance to the language's declared model: node kinds, child roles, attributes, rules between nodes. The repository calls this admission | The check stage, from a grammar table and from code |
| 3 | Validity in the language itself: names resolve, the program type-checks | The language's own judgment, through analysis |
| 4 | Correctness against another model: does this Go type represent that horos type | Verification |

Levels 1 to 3 are consistency at increasing depth inside one language. Level
4 crosses models. All four are decided by operations. What a model can state
is the declaration that levels 1 and 2 are checked against. No page states
these levels.

Three facts behind it, implemented today:

- The tree type is not checked by construction: "the carrier performs no
  validation or allocation" and "No immutability or validity is enforced
  here" (section 4.1). The evidence that a tree was checked is a value of
  another type, the envelope: "AST is an owned, structurally checked, opaque
  snapshot." (`trees/representation/go/ast.go:23`).
- An accepted page has the shape of a generic check with a supplied
  declaration (fixed point 14). The grammar interpreter of section 4.2 is
  that shape.
- Consistency is decided in three ways today: from a table, by a generic
  checker (the grammar tables); by code inside a check provider (3 to 14
  rules per language); and by a validating constructor on the model side
  (the semantic graph's `New`, the method `Capture` of the file data type,
  and the checks of an intended-forms package).

Two accepted sentences fit the levels (fixed point 13): generation
establishes no proof that a representation is correct, which is level 4, and
parsing establishes neither reference resolution nor semantic validity,
which is level 3.

**Three cases, which are levels 1 and 2 decided in three different places.**

*A constructor, and a separate check of the same condition* (level 2 for
intent). Go's intended-forms package builds a field form and, in another
function, refuses a malformed one
(`transform/syntax/golang/model/go/model.go:43-45`, `:70-73`; `g` is the
generation contract package):

```go
func Field(name string, optional bool, ty g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormField), []string{name, strconv.FormatBool(optional)}, ty)
}
	...
	case FormField:
		if len(d) != 2 || d[0] == "" || (d[1] != "true" && d[1] != "false") || len(a) != 1 {
			return g.ErrInvalid
		}
```

The constructor `Field` refuses nothing. It accepts an empty name. The
refusal is in `CheckExpression`, which the planner reaches through the
target. Other constructors do refuse: the form issuer returns an empty form
for an empty name, and the constructor of the semantic graph validates its
input. Our test counted validation at construction as well-formedness and
counted `CheckExpression` as processing. So the same condition, that a field
has a name, is a model's own if a constructor enforces it and an operation
if a check does.

*A typed reader over a graph* (it relies on level 2 for facts). Go's
semantic vocabulary package reads the element type of a slice
(`semantics/golang/go/model.go:276-282`):

```go
func (t Type) SliceElement() (Type, bool) {
	v := t.term()
	if v.Operator != SliceOperator() || len(v.Operands) != 1 {
		return Type{}, false
	}
	return Type{t.owner, v.Operands[0].Term}, true
}
```

The reader states the shape of a slice term, one operator and one operand.
It also takes one artifact, a term of the graph, and returns another, a
typed handle. Our test counted the readers as processing. One reader of the
same file also gives a coverage answer: it refuses every basic kind but two
as unsupported (`model.go:164-167`).

*The check that makes a list of nodes a tree* (level 1). It lives in a check
directory (`check/syntax/shared/ast/go/structure.go:66-76`, `:92-94`;
`MaxNodes` is a bound, `budget` counts storage, and `text.LimitExceeded` is
a diagnostic code). The check of the concrete tree calls this function
first:

```go
func Structure[L any](tree ast.AST[L]) error {
	bad := func(id ast.NodeID, message string) error { return text.AtNode(tree, id, message) }
	if tree.ID == "" || tree.Root == 0 || int(tree.Root) > len(tree.Nodes) {
		return bad(tree.Root, "nonempty tree ID and valid root are required")
	}
	if len(tree.Nodes) > MaxNodes {
		return &text.Diagnostic{Code: text.LimitExceeded, Message: "too many AST nodes"}
	}
	if err := budget(tree); err != nil {
		return err
	}
	...
		if seen[e.id-1] {
			return bad(e.id, "containment cycle or multiple parents")
		}
```

The envelope package, a carrier, calls it. Our test did not decide these two
packages. One level up the same holds for a language: the grammar table
says what a Go tree is, and its interpreter and the rules beside it are the
language's checker, a stage (section 4.2).

**The hard case.** A model that contains no operation cannot guarantee its
invariants by construction. For trees the repository already separates the
two: an unchecked data type, and a checked state that another value
witnesses. Whether "checked" then becomes a separate, evidenced state for
every sort of model, and what carries that evidence, is asked in question 1.

### 6.2 The line: nine cases, and four ways to place packages

**The cases of question 1.** Each is a piece of code that our test could
not assign, or assigned in a way that is asked.

- **(a)** A constructor that refuses a malformed value, against a separate
  check that returns a verdict on the same condition (section 6.1, the
  first case; section 4.4).
- **(b)** A typed reader over a graph (section 6.1, the second case;
  section 4.3).
- **(c)** The check that makes a list of nodes a tree, and a language's
  grammar table with the rules beside it (section 6.1, the third case;
  section 4.2).
- **(d)** The generic data types for trees, graphs and files, two of which
  hold or call code of an operation (section 2.3).
- **(e)** The value types and trusted handles (snapshot, receipt) that sit
  in each operation's contract package, and the contract package that is
  left (section 2.5).
- **(f)** The shared mechanisms that select and compose: the engine, the
  registry, the planner (sections 1.2 and 2.2).
- **(g)** What describes a relation between two languages: the identity and
  decision-record schema of one mapping, and the identity of one family's
  meaning (section 2.7).
- **(h)** A language's lexical rules, which both sides need (section 4.2).
- **(i)** The intended-forms package of a language (sections 3.1 and 4.4):
  a model of the language, or the working data of one operation.

**Four ways to place packages.** Standing: our derivation. Each is stated
with the same five points. None is chosen. What is below a root is open:
the names in the owner's message 7 are a sketch (Question), and fixed point
23 is in question for the model root.

**P1. The intended-forms packages and the symbol packages stand beside the
syntax and semantic vocabulary packages.**

- *On the model side:* the four sorts of vocabulary package of one language.
- *Rule 1:* not met. Every vocabulary package still imports a contract
  package.
- *Rule 2:* not met. Twelve packages still hold both.
- *Its risk:* it presupposes that the intended forms are a model of the
  language. If they are the working data of generation, they stay with that
  operation, and only the symbol packages, which two stages read, are in
  question.
- *It leaves open:* everything about carriers, contract packages, grammar
  tables and the packages of a pair of languages.

**P2. The carriers and the value types of the contract packages are on the
model side.** The operations root holds the stages, the kernel, what is left
of each contract package, the bridges and the tools.

- *On the model side:* the carriers, the vocabulary packages, and every
  value type that operations exchange.
- *Rule 1:* holds as one rule on a path prefix.
- *Rule 2:* every package that holds both is divided, and ten contract
  packages are split. Each split is a verdict on one of the cases (a) to
  (e). What is left of a contract package is an interface and the admission
  code of one operation.
- *Its risk:* the model root then holds every data type of the repository,
  and the word stops distinguishing anything. The handles that only
  admission code may create need a way for an operation to create a value
  of a model-side type.
- *It leaves open:* the cases (f), (g) and (h), and whether a grammar table
  is part of a language's model.

**P3. The contract packages stay with the operations, and the value types
that both sides need are in a third place.**

- *On the model side:* what one language says about itself, and possibly
  the carriers.
- *Rule 1:* a model imports no operation, and both sides import the third
  place. The direction names two roots. A third place is then a part of the
  model root that the rest of it may import, or a third root.
- *Rule 2:* the packages that hold both are divided as in P2. The value
  types leave the contract packages for the third place.
- *Its risk:* without a definition of its own the third place collects
  whatever is hard to classify.
- *It leaves open:* whether the carriers are model or third place, and the
  cases (f), (g) and (h).

**P4. A contract package is kept whole, as a third sort of thing beside
model and operation.** A contract package is then neither: it is the stated
interface between an operation and its contributions, with the values that
cross it.

- *On the model side:* the carriers and the vocabulary packages.
- *Rule 1:* a model imports no operation. Whether a model may import a
  contract package must be decided: today every vocabulary package does.
- *Rule 2:* the twelve packages that hold both are divided. The ten contract
  packages are not.
- *Its risk:* it answers the cases (e) and (f) by a third word and not by a
  test, and the direction names two roots.
- *It leaves open:* where contract packages live, and the cases (g) and (h).

**Where the hard packages go.** "Verdict" means that the place follows from
your answer on the named case.

| Package or part | Today | P1 | P2 | P3 | P4 |
| --- | --- | --- | --- | --- | --- |
| The tree data type, `trees/ast/go` | Top level, role foundation | As today | Model | Model or third place | Model |
| The envelope package, `trees/representation/go` | Foundation. Calls two checks | As today | Model, if the structure check counts as the tree's own consistency. Else it may not call the check | As P2 | As P2 |
| The structure check, `check/syntax/shared/ast/go` | Kernel, in a stage directory | As today | Verdict on (c) | Verdict on (c) | Verdict on (c) |
| The receipt, in `contracts/semantic/go` | In the contract package | As today | Its type on the model side, its creation in the operation | Its type in the third place | Stays in the contract package |
| The data types of intent and the planner, `contracts/generation/go` | One package | As today | The data types on the model side, the planner in the operations | The data types in the third place | One contract package |
| The names and the readers of Go's semantic vocabulary package | One package | As today | Names on the model side. Readers: verdict on (b) | As P2 | As P2 |
| The names and the checks of an intended-forms package | One package, inside `transform/` | Moves out whole | Verdict on (i), then on (a) | As P2 | As P2 |
| A symbol package | Inside `transform/`, read by two stages | Moves out whole | Model. Its admission rule: verdict on (a) | As P2 | As P2 |
| A lexical helper, `internal/rust` | Private, at the module root | As today | Importable from both roots (fixed point 21): verdict on (h) | Third place | Verdict on (h) |
| Go's grammar table | A private variable of the checker | As today | Verdict on (c) | Verdict on (c) | Verdict on (c) |
| A per-pair vocabulary package of the migration | Under `transform/`, by pair | As today | Verdict on (g) | Verdict on (g) | Verdict on (g) |

**The two rules and the package roles.** A root gives a first cut, model or
operation. The package roles of section 2.6 (foundation, kernel, vocabulary,
contribution, private part) still need a path pattern or a table entry
each. Our reading is that the two rules are added to the roles and do not
replace them. The alternative is that the two roots replace the roles
foundation and vocabulary, and only the roles inside the operations root
remain.

### 6.3 Unification: who would use a shared vocabulary, and three options

**Who would use a shared vocabulary.** Fixed point 5 asks for a concrete
case. Three candidates exist, and all three are facts of section 3:

- *Go as a source and as a target.* In the pair golang-to-typescript Go is
  read as a source through its semantic vocabulary. As a target Go is
  written through its intended forms. A mapping author meets two sets of
  names for Go's types. Section 3.5 shows what the first attempt met.
- *Verification laws against mappings.* A law reads an analysed Go type, and
  the mapping that produced the type stated it as intended forms (section
  3.4).
- *A Files formatter.* It compares what was intended with what the compiler
  bound. Today it compares identities only (section 4.5).

The source side has no intended-forms package (section 3.2), and no page
states that asymmetry as a decision.

Each option is stated with the same four points.

**U1. Three sets of names stay, and their relation is written down.**

- *What a user gains:* nothing new to learn, and a stated relation where
  today there is code.
- *What the unification test must show:* nothing, because nothing is
  unified. The platform's table has a row for the finding "Independent
  identities or lifetimes". Its response is: "Keep separate concepts and
  define their relationship."
- *What must be represented:* the relation from a form to the nodes it
  becomes, which today is the `switch` in each target, and the relation from
  a form to what analysis reports, which today exists nowhere.
- *Its cost and what it leaves open:* each notion keeps two or three names,
  and each shape stays stated in two or three places (section 6.4). A user
  on both sides of a transformation keeps translating.

**U2. Intended forms and analysed kinds share their names per language, and
differ in status, intended or established.** Go would have one set of names
for its type-level notions. Intent would use them as form names, and the
graph as operators and entity kinds. Whether a name states intent or an
established fact would be carried by the data structure it occurs in, and
not by the name. Node kinds would stay separate.

- *What a user gains:* one set of names for Go's types on both sides of a
  transformation and in verification.
- *What the unification test must show:* that a user who reads that a type
  is a slice of some element type keeps the same guarantees whether the
  statement is intended or established, and that every difference is
  represented and not erased (fixed point 6).
- *What must be represented:* (1) What a name refers to: intent refers to a
  destination, and analysis to the entity that a compiler bound. Fixed point
  18 makes those two identities. (2) Structure: intent is a tree of values
  with placeholders, and analysis a graph of entities and terms with facts.
  (3) Notions on one side only: the optional flag in intent; tags,
  embedding, type identity and the underlying type in analysis. (4) What is
  a restriction of today's coverage and no part of the notion: a defined
  type with a record body only. (5) Lifetime: intent lives within one
  generation request, and a graph is bound to the receipt of one captured
  program.
- *What it leaves open:* fixed points 8 and 9 separate the planned program
  from the analysed one. Whether that separation also requires two sets of
  names is said by neither. That is our reading.

**U3. One description per language, from which node kinds, analysed kinds
and forms are all derived.** A single description states, for each
type-level notion, its tree shape, its fact shape and its form shape.

- *What a user gains:* as U2, and one place that says what a language's
  types are.
- *What the unification test must show:* all of U2, and that the syntax
  vocabulary shares the contract.
- *What must be represented:* (1) Syntax has notions that are no type-level
  notions: files, imports, parameter lists. (2) A node has a site, an order
  among its siblings and a spelling. (3) The relation from a form to nodes
  is not one to one: for Rust and C++ one form becomes a path of several
  nodes. (4) Not every language has all three: C++ has no analysis.
- *What it leaves open:* in which notation such a description is written
  (section 6.4), and whether it is one description or three that refer to
  each other.

U2 and U3 are about one language. Whether the four near-copy intended-forms
packages share a mechanism across languages is a separate open matter (the
paragraph under fixed point 27).

### 6.4 Declaration: what is data, what is code, and four options

The owner's second question in message 7 is: "what is currently not
declarative about our models? like, don't we already have an ast / cst
model? why do we then need code for the model for the different languages
even?" The facts that answer its first half follow. Its second half is
question 3.

**What is already data** (Go literals):

- Node kinds, role names, attribute names, kinds of concrete-tree node and
  profile strings. For syntax the generic tree is the whole data structure,
  and a language's package is constants only (section 4.1).
- The grammar table of each language, except the predicates inside it. It
  is an unexported variable of a contribution and holds function values.
- Semantic names, as one-line functions that return a key.
- Form names and primitive names.
- Reserved-word lists, four regular expressions and one table of character
  ranges.
- The simple mappings: a pattern and a constant output.

**What is code, and what it expresses.**

| Code | What it expresses that no table of the repository states |
| --- | --- |
| 3 to 14 rules beside each grammar table | Relations between two attributes; between an attribute and a child; between a node and its parent or ancestors; order within one role; consistency and uniqueness among siblings; reference roles other than one |
| Attribute predicates | Which strings are identifiers, literals or module names of a language, including its reserved words |
| The structure check of the generic tree | Bounds, valid indices, no cycle, one parent |
| Parsers and formatters, entirely | Concrete spelling, precedence and layout. No table states the concrete syntax of any language |
| The writer of facts (a bridge) and the reader (a vocabulary package's typed readers) | The shape of each fact: which operator, how many operands, what the data bytes hold |
| The types of a compiler run's inputs | Normalisation, path rules, and the encoding into settings |
| The constructors and checks of an intended-forms package | The shape of each form: how many data strings, what they mean, how many arguments |
| Coverage answers and the admission of symbols | What today's generation supports; which names a compiler's library declares |

A shape is stated twice in code and, with one exception, in no declaration:
a fact's shape by the bridge that writes it and the reader that reads it
(the exception is the comment on one relation, section 4.3), and a form's
shape by its constructor and its check, with a third reading by position in
the target.

Two of the three bridges are not written in Go (section 4.3). Each spells
the names it writes in its own language. So a declared model would have to
be consumed from three implementation languages, or the names in two of them
stay hand-written.

**Four options,** each with the same four points.

**D0. Nothing more is declared.** Shapes stay code and are held by tests.

- *What is declared:* the names, and the grammar tables as they are.
- *What reads it:* the checker, and a documentation tool through the source
  text.
- *What stays code:* everything in the table above.
- *The two bridges that are not Go:* spell names by hand, as today.

**D1. A Go table type, extended.** The grammar table becomes an exported
value on the model side, with predicates named by identity in place of
function values. Tables of the same sort state fact shapes and form shapes.

- *What is declared:* tree shapes, fact shapes and form shapes, as Go
  values.
- *What reads it:* the generic checker, parsers, formatters, tools and
  tests, by import.
- *What stays code:* the rules between nodes, unless the table type grows
  the means to state them; concrete syntax; the bridges.
- *The two bridges that are not Go:* a generator would have to write their
  names from the Go tables.

**D2. A horos declaration that uses the node patterns of horos kinds**
(section 4.6).

- *What is declared:* tree shapes as horos kinds. The notation has
  sequences, alternatives, repetition, recursion and parameters. Today it
  ranges over 22 fixed descriptor constructors, has no child role names and
  no attribute predicates, and describes no language's tree, the tree of
  horos included.
- *What reads it:* the horos library, and bittype's own pipeline, which
  already generates Go from horos types for one sort of value type (section
  4.7).
- *What stays code:* what the notation cannot say; concrete syntax; the
  bridges.
- *The two bridges that are not Go:* bittype generates TypeScript and Rust
  declarations from horos, so their names could be generated too. Fixed
  point 26 bears on this option: the shared parts import no horos today.
  And horos is one of the five languages whose tree would be described.

**D3. A separate schema.** A description format of its own, from which the
Go constants and tables are generated.

- *What is declared:* tree shapes, fact shapes and form shapes, in a
  notation fitted to what sections 4.2 to 4.4 need.
- *What reads it:* a generator, at build time.
- *What stays code:* the generator, concrete syntax, the bridges.
- *The two bridges that are not Go:* the generator writes their names.
  Fixed point 5 asks, for a new notation, for a concrete case that the
  existing ones cannot express.

**One worked row: Go's `field`.** What each option would hold for it. The
cells for D1 to D3 are sketches of ours.

| | The tree shape of a field | The fact shape | The form shape |
| --- | --- | --- | --- |
| D0, today | A constant `KindField`, and an entry in the checker's private table: the attributes `name` (required, an identifier) and `tag`, one child under the role `type` | The literal `field` in the bridge, and a JSON entry (name, tag, embedded, entity) written by the bridge and read by a reader | The constant `FormField`, a constructor that writes two data strings and one argument, and a check that counts them |
| D1 | The same entry, exported, with the predicate named instead of held as a function | A table entry: the operator of the struct term, its data fields and one operand per field | A table entry: two data, their meaning, one argument |
| D2 | `node("field", ..., item(K))`, once the notation can name a child role and constrain an attribute | Not expressible today: the notation describes trees and no facts | Expressible as a tree shape, if a form is taken as a tree node |
| D3 | An entry of the schema | An entry of the schema | An entry of the schema |

**What must stay honest whatever is declared.** Five parsers, five
formatters, three bridges and four targets are hand-written and must agree
with any description. The checks of that agreement that we know of today are
three: a parser and a formatter call the checker; TypeScript's Files
formatter compares each destination with what the compiler bound; and a
documentation tool re-reads the grammar tables.

### 6.5 One step we intend to take first

Standing: decision of the research owner, recorded on the work item of this
consultation on 13 November 2026, before the owner's eighth message. It is
open: please say if it is wrong or if it prejudges something. Rule 2 now
makes it one instance of a general rule (section 6.1).

**The facts.** Each intended-forms package holds names, five constructors
and two checks (section 4.4). The mappings use the constructors and the
names and never the checks. The target uses the names and the checks and no
constructor. The one production caller of the checks is the target of the
same language. Most of the conditions in the checks are well-formedness of
a form. Three are coverage answers of today's generation: which primitives
exist, Go's refusal of every external type, and the refusal of an unknown
form name. For Go a fourth is the restriction of a defined type to a record
body. Rust already keeps its coverage answer for external types in the
target and the shape check in the package.

**The reason.** The platform's rule separates wrongly grouped parts first,
as a move (fixed points 2 and 3). The boundaries page lists the signs of
wrong grouping. Two of them hold here: "An importer that uses a slice of the
package" (both importers do) and "A policy next to a shared mechanism" (the
coverage answers sit beside the names). A coverage answer is a statement of
the target contribution about what it builds.

**The step.** The coverage answers leave the intended-forms packages for the
target of the same language, as a change of their own. The well-formedness
conditions stay where they are for now. Whether a model may hold them is
case (a) of question 1. The accepted definition of a vocabulary says "and
nothing else", and what is implemented holds more (section 5.3, item 1).

**What waits, and what the work in flight does meanwhile.** The step is
taken when the migration's branch no longer edits those files. The branch
edits all four packages and adds a third check to each (section 2.7). So
the work in flight adds a check to the packages from which this step takes
conditions away.

### 6.6 What a move would touch

Standing: implemented today, from a count over the committed files. Under
fixed point 1 these numbers have no weight in the choice between options.
The same rule lets them guide the order of delivery. They are stated once,
here, because question 4 is about order.

| What | Count |
| --- | --- |
| Production packages that get a new import path under the direction | 129, and 3 packages that hold only tests |
| Packages that rule 2 reaches, because they hold both under our test | 12 of the 22 model-side packages, and 10 of the 13 contract packages |
| Packages that move under P1 | 7 |
| Import strings of this module in Go files | 2,121 in 365 files: 959 in non-test files, 871 in test files, 291 in fixture files under test data |
| Places that compute a package role from a path | One function with its helpers, and five tables keyed by path with 31, 12, 5, 12 and 14 entries (section 2.6; the table of 5 names the language of packages whose path does not) |
| Quoted layout paths in test files | 1,166 in 17 files, 638 of them in one file of synthetic package sets for the architecture test |
| Another reader of paths | The documentation tool finds language markers under `syntax/` and contract packages under `contracts/` by path prefix |
| Mentions of a layout path in pages | About 1,150 in the design pages, the layout page and the change log. About 660 more are in dated research records, which are not edited |
| What does not change | The module path. No version tag pins the present import paths |

Work in flight that is written for the present paths: section 2.7.

## 7. Questions for the expert

The split and its two rules are not asked. What we want back is a test that
we can apply to a package ourselves, verdicts on named cases, patterns from
comparable systems with what failed, and an order of work. We do not want a
review of the code shown.

Owner rules and the owner's direction are fixed. Accepted designs are fixed
unless you see a contradiction, and section 5.3 names five that we see.
Everything else is open: what the code does today, every reading marked as
ours, every decision of the research owner, and the options of section 6,
to which the labels P1 to P4, U1 to U3 and D0 to D3 refer. Please name an
option, or state another, and then the reason.

1. **The line between model and operation.** The direction says three things: a model contains no operation, a model imports no operation, and operations work on the model. The owner then asked whether what we could not assign is the correctness of a model, or its consistency. How would you draw the line as a test that can be applied to one piece of code and that makes "the model is independent of the operation" checkable, and what does each of the three statements then require and forbid? Is the difference between consistency (a value agrees with its own model, decidable from the value and the model's declaration alone) and correctness (agreement with something outside the value, which needs evidence) the right axis, and would you keep, merge or replace our four levels? A model that contains no operation cannot guarantee its invariants by construction: should "checked" then be a separate, evidenced state for every sort of model, as it is for syntax trees today, and what carries that evidence? Please apply your test to the cases (a) to (i) of section 6.2, with one verdict and one reason each. Where your test gives neither answer, is that a third sort of thing beside model and operation, and how would you define it? If declaring a model as data (question 3) is itself the test, please say so.

   *What it rests on.* The direction is fixed point DIR. Section 6.1 gives
   our test, what each rule reaches by it, and our derivation from the
   owner's messages 9 and 10: consistency, correctness and the four levels.
   Section 6.2 lists the cases (a) to (i) with pointers to their code, and
   four ways to place packages (P1 to P4). The carriers are in section 2.3,
   the contract packages in 2.5, the shared mechanisms in 1.2 and 2.2. Fixed
   points 4, 13, 14 and 19 are the nearest rules and accepted sentences.

2. **What one language's model is, for every operation and on both sides of a transformation.** The owner asked whether the forms that generation states are a syntax model or a transform model, and whether transformation is not just one of many operations working on a language's model. If any language can be the source or the target of a transformation, what would you have one language present to all the operations on it: how do its abstract syntax, what analysis reports about a program, and what a transformation intends to construct relate, and are the intended forms a model of the language or the working data of one operation? How do systems you know arrange this, and what failed? What test tells you when two of these should share their names and differ only in status, intended or established, and when they stay apart with a stated relation? Please apply it to Go's rows in the side-by-side table.

   *What it rests on.* Section 3: the four vocabulary packages of Go (3.1),
   one record traced through them (3.2), the three sets of names beside each
   other (3.3), and the other readers of what analysis reports (3.4). Section
   3.5 is the case where it broke: the first pair with Go as a source could
   not walk Go's analysed graph. A decision of the research owner on that
   case stands, and you may object to it. Section 6.3 says who would use
   shared names and gives three options (U1 to U3). The fixed points on
   intent are 8 to 12, and those on a shared model of types 15 to 18.

3. **Declared models.** The owner asks why a language's model needs code at all, given one generic tree for every language. How far can a language's model be data, an instance of a generic description that shared code interprets, so that describing a language adds no code? What cannot be stated that way, and is that remainder still model? Would you have one description per language or one per sort of model (tree shapes, fact shapes, form shapes), and what is the smallest declared core worth building first? What gives a declared model its identity and version, and what keeps it honest against hand-written parsers, formatters and compiler bridges, two of which are not written in Go?

   *What it rests on.* Section 4 is what this question is for: names (4.1),
   a grammar table with the rules it cannot state and a lexical rule (4.2),
   a fact written and read (4.3), a form built and checked (4.4), the node
   patterns of horos kinds (4.6) and the one generator (4.7). Section 6.4
   says what is data and what is code today and gives four options (D0 to
   D3) with one worked row.

4. **Carrying the split out.** Below the two roots, what would you organize the model root by (sort of model, language, pair of languages), and what stays of the stage-by-language matrix below the operations root? A migration of the mapping contract is in flight, and two further language pairs and two Files formatters are about to be written against the present places; our rule is that a misplaced responsibility moves first, as its own change, and that later steps only add. In what order would you take the decisions and the moves so that nothing moves twice, what must be decided before the first move, and would you land the work in flight in the present places first or write it against the new ones? Which part of the line can a check on import paths carry, and what must remain a judgment made in review, given our rule that such a check enforces a boundary that reasoning has decided and does not replace it?

   *What it rests on.* The path pattern is in section 2.4, the package roles
   and the architecture test in 2.6, and the work in flight with its gate
   and the tempo of this project in 2.7. Section 6.2 says what is open below
   the roots, 6.5 gives the one step we intend to take first, and 6.6 says
   what a move would touch. The rules are fixed points 2 and 3, and the
   layout's fixed points are 20 to 25.

5. **What we are not seeing.** Which of the distinctions this document draws would you not draw, and which would you draw that it does not? What would you want to know, or have us measure or try, before the first move? And what would you ask in our place?

   *What it rests on.* The whole document, including its frame: the four
   working questions, our test of section 6.1, the four levels and the
   options.

## 8. Glossary

Where a word has more than one sense in this document, the entry lists
them. The list of terms in the Question section holds the dozen that
are needed from the first page.

| Term | Meaning here |
| --- | --- |
| AC-01 to AC-11 | The eleven numbered rules of architectural correctness, owner rules that apply to every project of the organization (section 2.1). AC-09: existing code has no weight in a design choice. AC-10: a misplaced responsibility moves first, as its own change. AC-11: a mechanism that needs no knowledge of one language is implemented once |
| Admit, seal, publish, issue | To admit a value is to validate it and wrap it in a type that only the validating package can construct. To seal it is to make it immutable. To publish is to hand out the sealed result. To issue is to create an opaque value that only its creator can read |
| Analysed kind | A name under which analysis reports something: the kind of an entity or the operator of a term, as a key |
| Architecture test | bittype's test that reads the imports of every package and compares them with the package's role (section 2.6) |
| Bridge | The code that runs one language's real compiler over captured source and writes the result into the semantic graph. Two of three are not written in Go |
| Carrier | A generic data type that holds one sort of artifact for every language: the trees, the envelope, the semantic graph, the files and the Source carrier. The data types of intent are generic too, but they are declared in a contract package |
| Consistency, correctness | Our derivation (section 6.1). Consistency: a value agrees with its own model and with itself, decidable from the value and the model's declaration alone. Correctness: a value agrees with something outside itself, which needs evidence |
| Contract, contract package | A contract package holds the interface between one operation's shared mechanism and its contributions, and the values that cross it (section 2.5). "The mapping contract" is what a mapping registers and returns; its migration is section 2.7. In the quoted rules "contract" also means an agreed obligation |
| Contribution | One thing a consumer registers or selects on its own: a stage provider, a source reader, a target, one mapping |
| Coverage answer | A refusal that means: not supported today. Opposed to a well-formedness answer, which means: this is not a valid form |
| Decision record | Under the contract migration, the typed value that a mapping returns for one selection. The code calls it a payload |
| Destination | What a name use in generated output is meant to refer to: a generated declaration, a type parameter or an external symbol. An identity, never a structure |
| Dimension | One of the two rows of stages. **Syntax**: text, concrete tree, abstract tree. **Files**: file set, source, semantics |
| Direction, the owner's | Messages 7 and 8: the repository splits at its root into models and operations. Rule 1 (no import): a model imports no operation. Rule 2 (no containment): a model contains no operation. Operations work on the model |
| Envelope | The type-erased, structurally checked snapshot of a tree or a graph, with the descriptor of its language. The package is `trees/representation/go` |
| Evidence | Three senses. In the rules: support for a claim. In the semantic contract package: a receipt offered as proof that an analysis with a given profile ran. In section 3: what analysis of an actual program established, as opposed to intent |
| Family | One sort of source type: text, booleans, lists, closed records, optional fields, named declarations. A mapping is one family for one ordered pair of languages |
| Files formatter | The contribution of a target language that places generated declarations into files, writes imports and visibility, has the result analysed and compares each name use with its destination. It exists for TypeScript and is planned for Go and Rust |
| Form | One sort of intended type construct of a target language, identified by a name that the language's intended-forms package issues. Go has seven. A constructor is the Go function that builds a piece of intent of one form |
| Gate | The condition under which the design of the contract migration becomes final: seven stated lines, each shown by tests (section 2.7) |
| Generation, transform | Producing one language's type declarations from another language's analysed declarations. The owner and the code say "transform" |
| Grammar table | One language's table of what each node kind may hold. A private variable of the language's checker. The code calls it the admission table (section 4.2) |
| horos | A declarative type language over ontos values, which are atoms and tuples in a canonical binary encoding. Its types denote sets of values |
| horos kind | A classifier of horos types, written with node patterns (section 4.6). Not a node kind |
| Intended-forms package | The package that issues one target language's forms, `transform/syntax/<language>/model/go`. The code calls it `model`. Whether it is a model in the owner's sense is asked |
| Intent | bittype's description of a target type before any syntax exists: a tree of forms with string data, arguments and placeholders. The pages also say "native intent" |
| Kernel | The shared mechanism that selects and composes: the engine, the generation registry and planner. It names no language |
| Mapping | A contribution for one ordered pair of languages and one family. It states the intended form for a use of a source type. The code also says "mapping library" |
| Model | Unqualified: the owner's sense, what belongs under the model root, which is what is asked. The code uses the word for two other things: the intended-forms package, and the reader type `Model` of a semantic vocabulary package |
| Native | The pages' word for the four programming languages as opposed to horos, and so also for the target side of generation. On the migration's branch the shared contract says "target" |
| Node kind | The discriminator of a syntax node, a typed string constant of one language |
| Operation | In the owner's direction: whatever works on a model; the owner writes "ops". In our test of section 6.1: a package that parses, checks, formats, prints, resolves, analyzes, transforms, verifies or composes |
| Package role | The import rule of a package: foundation, kernel, vocabulary, contribution, private part, composition. Computed from the package's path or read from a closed table (section 2.6). A child role and a reference role are something else: labels of edges in a tree |
| Path pattern | `<stage>/<dimension>/<scope>/<component>/go`. The scope is one language, an ordered pair or `shared`. The last segment names the implementation language |
| Profile | A named promise, identified by a string. Of a language: which revision of its syntax. Of generation: what the generated types retain. Of an analysis: which compiler version and which fragment of the language |
| Provider | The language-specific code that writes and reads a graph. A stage provider is the contribution that implements one stage for one language |
| Receipt | The value that the semantic contract package issues for one finished resolution or analysis over a snapshot. The code calls it `Result` |
| Semantic graph | The carrier of analysed meaning: entities, terms, bindings, expressions, dependencies and facts (section 4.3). The code calls it `Semantics` |
| Source reader | The contribution of a source language that interprets its analysed graph for the kernel and for mappings. Only horos has one on the main branch |
| Stage | One processing step: parse, check, format, print, resolve, analyze, generate |
| Standing | Which of seven sorts of statement something is (Question). For intent against analysis the document says status, intended or established |
| Symbol, symbol package | A symbol identifies a referent outside the generated declarations: a predeclared type, a library item, a package export. A symbol package issues one language's symbols |
| Target | A language that generation produces, and the contribution that checks intent for it and builds its trees. A target program is the program that will exist after generation |
| Vocabulary, vocabulary package | A vocabulary is a set of names. By the accepted definition it is "The names two parties must agree on and nothing else". A vocabulary package says what exists in one language. A language has four sorts today: the syntax vocabulary package, the semantic vocabulary package, the intended-forms package and the symbol package (section 3.1). The migration adds vocabulary packages of a pair and of a family |

## Applied

**Evaluated against:** `0018-models-and-operations.submitted.md` (sha256 `67dca75d6c7b5e41243347e10f6e75ab3bb18a4fa0c46271ad6722dc66201d51`), attempt 1, run `run_36aec98c-ac83-446d-a057-f81b327ed237`; code at bittype `3998db5` (`main`), which differs from the described baseline `a28febc` by two documentation commits; the migration's branch at `af1dcde`, the two other pairs at `447c44c`.
**Graded:** https://github.com/Bitspark/bittype/issues/305#issuecomment-6855741567 (by the research owner, claude-contract-chain-issues; bittype has no architecture seat). Slices S0 to S9 are defined there.

**How it was evaluated.** The advice was split into 39 recommendations (R1 to R31, with R7 in nine parts). Two agents that did not write the document checked each one's premises against the code, the owner rules and the accepted designs, by reading, and proposed a verdict with evidence: [rows R1 to R20](0018-models-and-operations.verify-1.md) and [rows R21 to R31](0018-models-and-operations.verify-2.md). The research owner decided each row. Where the decision differs from the verifier's proposal, the row says so.

**What the advice did not answer.** It states no objection to decision O45 (section 3.5) and does not mention it; that is no confirmation. It names no failure of a comparable system. It does not say when the root move itself happens, what a sealed Go type's split is, or where six things live that its tree does not name.

**Eight places where the advice meets a rule, an accepted page or itself**, found by the verifiers and decided in the rows named:

1. Two accepted pages give the intended-forms package "validation" and "validity" (`LAYOUT.md:67`, `docs/generation-composition.md:17`), against the same repository's definition of a vocabulary ("The names two parties must agree on and nothing else", `docs/package-boundaries.md:98`). The advice sides with the definition, and so does rule 2. The two sentences change (R2, R7i).
2. A vocabulary publishes identities that are values of contract-package types. Contract packages are operations, and a model imports none. So the vocabulary role exists on both sides (R2, R8, R27).
3. "The planner and the plan types in `contracts/generation/go` are one unit, because only the planner may issue a plan" (`docs/package-boundaries.md:144-146`). A model-side expression carrier would separate a planner-issued kind from its issuer (R15).
4. `docs/source.md` declares the Source as "private sealed storage" with two receipts. "Generic carriers: model" does not hold for it as a whole (R7d).
5. A content digest as a model's identity contradicts three accepted sentences, among them "Profiles are caller-owned compatibility promises, not hashes" (R19).
6. The advice orders declared schemas syntax, facts, intent. AC-10 and AC-11 order them the other way: the per-language copies are on the intent side, and fact schemas are additions (R17).
7. Language first below `model` contradicts "Do not create a language-first packages tree" (`LAYOUT.md:134`, `AGENTS.md:237`) and the owner's sketch, and the advice's own section 4.3 says sort first (R21).
8. The advice's "rejected" (a policy refuses) is not the accepted `Rejected` (the evidence refutes), and its four outcomes have no limit (R25).

| # | Recommendation | Verdict | Action | Owner | Link | Validation |
|---|---|---|---|---|---|---|
| R1 | The test: a responsibility declares what is true or admissible (model), or determines, interprets or produces (operation) | NEEDS ADAPTATION | Adopted as the test, with two added sentences. (1) A type whose non-zero values only a refusing function creates is the admitted state and lives with that function on the operations side; model-side types are plain data. Go allows nothing else. (2) A constraint that exists only as its evaluator is named by the model and stated by the evaluator; the declared model is then an upper bound of what is admitted, and the page says so (AC-05). Five kinds of code the bare test leaves open are listed in the design page with their side | claude-contract-chain-issues | #305, S0 | the page merged after an independent review that checks each quoted rule against its source |
| R2 | What a model package may and must not contain | NEEDS ADAPTATION | Adopted. Identities typed by a contract package go to an operations-side vocabulary. "Builders" is read with R1's first sentence | claude-contract-chain-issues | #305, S0, S5 | the page merged after an independent review that checks each quoted rule against its source |
| R3 | A value type that two operations need is owned independently of either | HOLDS | `language.ID` and `language.Provider` become a shared model-side value. A receipt of one operation that another consumes stays with its producer | claude-contract-chain-issues | #305, S6 | the import rule red by a seeded import from a model carrier into an operation |
| R4 | An operation takes the model as explicit input and does not own the declaration it interprets | HOLDS | Operations own declarations today in three places: the grammar tables, 11 literal keys of the Go bridge, the names spelled in the TypeScript and Rust bridges. Each moves to the model side with its split | claude-contract-chain-issues | #305, S3, S9 | the tables are data in a model package; the checker test red by mutant of a table entry |
| R5 | Four levels; language validity is external correctness; no guarantee is inferred from another | HOLDS | The document's own definition of correctness names a compiler's judgment, so its level 3 under consistency was wrong. The design page states the four levels; no page does today | claude-contract-chain-issues | #305, S0 | the page merged after an independent review that checks each quoted rule against its source |
| R6 | A checked state apart from the artifact for every sort of model: artifact, evidence, admitted artifact | NEEDS ADAPTATION | No new concept: the envelope, the snapshot, the receipt and the Source are that pattern, by accepted pages. Four gaps close: native syntax vocabularies declare their profile; a parsed or formatted envelope records that grammar admission ran; a fact-shape check exists apart from graph admission; captured Files get a type of their own | claude-contract-chain-issues | #305, S6, S9 | the import rule red by a seeded import from a model carrier into an operation |
| R7a | Constructors and checks: split | NEEDS ADAPTATION | A sealed type moves with its constructor to the operations side; they are not split. Constructors that refuse nothing are construction, operations side | claude-contract-chain-issues | #305, S2, S6 | the five intended-forms packages hold no function; one checker test per table, red by mutant; outputs equal by digest |
| R7b | Typed readers are operations; the fact schemas they read are model | HOLDS | The readers of `semantics/{golang,rust,typescript}/go` leave. Their role is named in the design page (they are neither contribution nor private part) | claude-contract-chain-issues | #305, S5 | no model-side package imports a contract package; the import rule red by a seeded import |
| R7c | Structure checks and grammar tables: split | HOLDS | The table type and the five tables go to the model side once the 17 predicate function values of a table are named by identity. The interpreter and the rules written in code stay operations | claude-contract-chain-issues | #305, S3 | the tables are data in a model package; the checker test red by mutant of a table entry |
| R7d | Generic carriers are model | NEEDS ADAPTATION | Of six packages: `trees/ast`, `trees/cst` wholly model; `semantics/go` and `files/go` split (plain records model, admission operations); `trees/representation` and `source` are operations except their plain records | claude-contract-chain-issues | #305, S6 | the import rule red by a seeded import from a model carrier into an operation |
| R7e | Contract packages split by responsibility | HOLDS | A model-side part exists in two of thirteen: language identity, and the identity types of the generation contract | claude-contract-chain-issues | #305, S6, S7 | the import rule red by a seeded import from a model carrier into an operation |
| R7f | Engine, registry, planner are operations | HOLDS | Fits the owner rule of 6 November 2026 | claude-contract-chain-issues | #305, S8 | the import rule of the two roots, red by a seeded import from model into ops; every test name passes as before the move |
| R7g | Mapping identities and decision schemas are model declarations | HOLDS | The vocabulary and family packages of the pairs are model once their identity types have a model-side owner. Codecs and mappings are operations | claude-contract-chain-issues | #305, S7 | no vocabulary or family package imports the generation contract; the plan digests unchanged |
| R7h | Lexical rules: split | HOLDS | Word lists and character tables are model; the evaluators are operations. Every caller of an evaluator ends on the operations side, so nothing is needed on both sides | claude-contract-chain-issues | #305, S4 | the evaluators import the tables; a test red by a changed table entry |
| R7i | Intended forms: split | NEEDS ADAPTATION | As R15. The two page sentences of point 1 change | claude-contract-chain-issues | #305, S2 | the five intended-forms packages hold no function; one checker test per table, red by mutant; outputs equal by digest |
| R8 | No third root; a contract is an operation's boundary; a value with independent domain meaning is model | HOLDS | 19 of 24 imports from vocabulary packages into contract packages end through the splits; 5 become model to model | claude-contract-chain-issues | #305, S0, S5 to S7 | the page merged after an independent review that checks each quoted rule against its source |
| R9 | A language has several related models | HOLDS | Syntax, analysed semantics, intent. Intent is a sort of model that the five abstraction levels of `docs/abstraction-levels.md` do not name; the design page says where it stands | claude-contract-chain-issues | #305, S0 | the page merged after an independent review that checks each quoted rule against its source |
| R10 | U1 first, with declared and testable relations | NEEDS ADAPTATION | U1. The relation from a form to nodes stays code, the target's (R18). "Declared" means a stated relation on the language's page and one test per form that builds, checks, analyses and compares. A relation as data waits for a reader | claude-contract-chain-issues | #305, S0, S9 | the page merged after an independent review that checks each quoted rule against its source |
| R11 | Shared concept identity only by criteria | HOLDS | By the platform's six criteria and its substitution condition (`semantic-ownership-and-generalization.md`), which are the binding wording | claude-contract-chain-issues | #305, S0 | the page merged after an independent review that checks each quoted rule against its source |
| R12 | Not one description for all three | HOLDS | | — | — | — |
| R13 | The Go rows | NEEDS ADAPTATION | Nine of eleven match the code. Type parameter: the binder is the kernel's, Go has nothing of its own to share. The analysed alias and defined kinds are wider than the declaration forms. Shared identities: see R13a | claude-contract-chain-issues | #305, S0 | the page merged after an independent review that checks each quoted rule against its source |
| R13a | A shared concept identity for slice, defined type, alias | DEFERRED | trigger: a non-test file imports both a language's intent vocabulary and its semantic vocabulary (`git grep -l` for both import paths over non-test files gives a hit). No reader exists today | claude-contract-chain-issues | #305 | — |
| R14 | "A Go defined type has a record body" is not Go's rule | NEEDS ADAPTATION | It leaves the intended-forms package. Its owner is the declaration contribution, not the target (`docs/contribution-contract.md:203`); the Go target lowers any body. Today it returns invalid; as a contribution's answer it is a refusal | claude-contract-chain-issues | #305, S2 | the five intended-forms packages hold no function; one checker test per table, red by mutant; outputs equal by digest |
| R15 | The intended-forms package becomes five things | NEEDS ADAPTATION | (1) and (2): each language declares its forms as a table, as data, with a declared identity. (3): only the native kind of an expression (form, data, arguments) is intent; the planner's other kinds stay one unit with the planner. (4): ONE checker and ONE constructor in the agnostic layer, which is the owner AC-11 asks for the five near copies of AG-7. (5): the target keeps lowering and coverage | claude-contract-chain-issues | #305, S2 | the five intended-forms packages hold no function; one checker test per table, red by mutant; outputs equal by digest |
| R16 | Declare with the mechanisms that exist; no new schema language; not the horos kind notation | HOLDS | Value-shaped schemas already have a mechanism the advice did not weigh: horos types with a generator and a regeneration test (the payload schemas of #196). The design page says which declaration uses which | claude-contract-chain-issues | #305, S0 | the page merged after an independent review that checks each quoted rule against its source |
| R17 | Three schema families, syntax then facts then intent | NEEDS ADAPTATION | The families hold. The order is by AC-10: the form table and the placement of the grammar tables are corrections and come first; fact schemas are additions and come after. Parameter and reference structures belong to the shared carrier. Versions are profiles | claude-contract-chain-issues | #305, S2, S3, S9 | the five intended-forms packages hold no function; one checker test per table, red by mutant; outputs equal by digest |
| R18 | What stays code | HOLDS | With R1's second sentence: the page says the declared grammar is an upper bound | claude-contract-chain-issues | #305, S0 | the page merged after an independent review that checks each quoted rule against its source |
| R19 | Explicit identity and revision, independent of package path and marker | NEEDS ADAPTATION | The identity is the profile, with the revision practice of #278. A digest is provenance, never identity. Missing and added: a declared profile for each native syntax vocabulary, an identity and revision for each form vocabulary | claude-contract-chain-issues | #305, S0, S2, S9 | the page merged after an independent review that checks each quoted rule against its source |
| R20 | A conformance suite; generated bindings for bridges not written in Go | NEEDS ADAPTATION | Three of its items exist. Added with the fact schemas: fact-shape validation, reader tests against schemas, generated names for the two bridges. "Targets consistent with intent schemas" gets a meaning only through R10's test per form | claude-contract-chain-issues | #305, S9 | per addition: its own red and green run |
| R21 | The layout | NEEDS ADAPTATION | Two roots. Below `model`: sort first, then language, as the owner's sketch and the standing rule have it. The six places the advice does not name are named in the design page before any move: the Source records, language identity, the intent carrier, the grammar-table type, lexical tables, the operations halves of split vocabulary packages. `internal/` stays at the module root until five imports from future model packages are split away; then it goes under `ops`, and Go's own rule enforces that a model imports no private part. Integration and tools go under `ops`. The trailing `go` stays. One module stays | claude-contract-chain-issues | #305, S0, S8 | the page merged after an independent review that checks each quoted rule against its source |
| R22 | The stage matrix stays below `ops` | HOLDS | 91 packages take the prefix and nothing else | claude-contract-chain-issues | #305, S8 | the import rule of the two roots, red by a seeded import from model into ops; every test name passes as before the move |
| R23 | Contribution identities and schemas under a model area | NEEDS ADAPTATION | No identity is a path today, so a move changes none. Their types leave the generation contract package for a model-side owner first. The design page says where a family of one language and where a pair's vocabulary live | claude-contract-chain-issues | #305, S7 | no vocabulary or family package imports the generation contract; the plan digests unchanged |
| R24 | The first move | NEEDS ADAPTATION | One agnostic checker over declared form tables, in the kernel package that already holds the single-scope check; the form table designed in this move; a stated home for issuer and constructors that mappings may import; the two coverage answers change category from invalid to unsupported on purpose, recorded. Made as its own change on the migration's branch before its gate is shown (R29) | claude-contract-chain-issues | #305, S1, S2 | a test that fails while a check sits in a model package, red on the branch before the move |
| R25 | Four outcomes for every operation | NEEDS ADAPTATION | The principle only: a valid value an operation cannot process is not invalid. Each operation keeps its accepted classes, which have a limit category and another meaning of rejected. Applied where it bites: R24's two answers and R14 | claude-contract-chain-issues | #305, S0, S2 | the page merged after an independent review that checks each quoted rule against its source |
| R26 | A model is not narrowed to what an operation supports | HOLDS | The two narrowing conditions leave with R24 | claude-contract-chain-issues | #305, S2 | the five intended-forms packages hold no function; one checker test per table, red by mutant; outputs equal by digest |
| R27 | Decide vocabulary, ownership and roles before any move | NEEDS ADAPTATION | A design page first. The roles do not all stay below a root: foundation is cut in two, and the reading half of a split vocabulary package gets a named role. Terms map onto the words the repository has before any is added. `LAYOUT.md:145-147` is amended | claude-contract-chain-issues | #305, S0 | the page merged after an independent review that checks each quoted rule against its source |
| R28 | Then misplaced responsibilities as independent changes | HOLDS | 92 of 129 packages are free of the migration's branch. Each split creates its new half at its final path | claude-contract-chain-issues | #305, S3 to S7 | the tables are data in a model package; the checker test red by mutant of a table entry |
| R29 | Do not land the work in flight on the present layout; retarget | NEEDS ADAPTATION | The order holds; three premises do not (58 of the migration's 98 packages change by a prefix only; no identity depends on a path; the pairs are built). Order: the design page; `CheckMember` leaves the model packages at once; the splits that touch no branch file on `main`; R24 on the migration's branch; the scripted root move before the migration merges, the same script on the branches; then the gate, on the lasting boundaries. **This differs from the verifier**, which would make R24 on `main` first: every importer of the four intended-forms packages is a file the branch rewrites, so that would do the split twice; made on the branch it is still its own change and precedes the gate | claude-contract-chain-issues | #305, S1, S2, S8 | a test that fails while a check sits in a model package, red on the branch before the move |
| R30 | Five rules for the import test | NEEDS ADAPTATION | Rules 1, 3 and 4 by imports, each with a test that breaks it. Rule 5 is Go's own. Rule 2, no operation under model, is a judgment in review and gets no check (owner rule: no keyword heuristics or lint gates in place of reasoning) | claude-contract-chain-issues | #305, S8 | the import rule of the two roots, red by a seeded import from model into ops; every test name passes as before the move |
| R31 | No move that only renames directories | HOLDS | 27 packages split, so a prefix-only move breaks rule 1 or rule 2. The scripted move is one step among the others | claude-contract-chain-issues | #305, S8 | the import rule of the two roots, red by a seeded import from model into ops; every test name passes as before the move |

**Noted on 13 November 2026, after the design page was written ([models and operations](../docs/models-and-operations.md), slice S0).** The rows above are not edited. Six things in this section are read with later decisions of the research owner on #305. (1) The page marked nine choices that no row settles, and they are decided as [M10 to M18](https://github.com/Bitspark/bittype/issues/305#issuecomment-6856947056). (2) The counts of M7, which rows R22 and R31 carry as 91 packages that take the prefix only and 27 that split, are corrected: of 129 packages 8 are wholly model, 23 split and 98 are wholly operations. The importers of the four intended-forms packages that the grade names are 9, 9, 7 and 6 files. (3) Row R15, part 3, is narrowed by M12: "only the native kind of an expression (form, data, arguments) is intent" says what intent is and places no type. The model side holds the form-table type and no value type, and the expression stays whole with the planner. (4) Row R24's change of category "from invalid to unsupported" applies to the primitive set only, by M17. Nothing refuses Go's record-body condition after the split, no code is written for it, and row R14 is the binding wording for it. (5) [M19](https://github.com/Bitspark/bittype/issues/305#issuecomment-6857366489) adjusts M11, after the page had checked M11 against the code: a language's lexical tables and the names of its predicates are a part below its syntax model, `model/syntax/<language>/lexical/go`, which imports no node kinds. Row R21's place for "lexical tables" and row R7h are read with it. (6) An independent reader then read the page, and its findings are decided as [M20 to M27](https://github.com/Bitspark/bittype/issues/305#issuecomment-6858493249). Row R7a's split of the symbol packages applies to TypeScript only, because the symbol types of Go and Rust are sealed (M20). Row R23's move of identities includes the publisher name (M21) and precedes the split of row R24 (M22). Row R29's order is a decision about where a correction is made and no deferral under AC-10, and row R28's splits are those of packages that the migration's branch does not change (M23, M24). Of row R4's three places, the eleven keys of the Go bridge belong to slice S5, and of row R6's four gaps the second and the fourth are designed in slice S9 (M26).
