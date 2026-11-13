# Research: Models and operations: where a language's models live, what they are made of, and whether they are declared

**ID:** 0018
**Date:** 13 November 2026
**Status:** draft
**Owner:** claude-contract-chain-issues
**Issue:** https://github.com/Bitspark/bittype/issues/305

## Question

**In one sentence:** in a library of language tooling, should everything that
describes what exists in a language live apart from everything that processes
it, under the rule that a description never depends on a process; what belongs
on the describing side; should one language's descriptions share one
vocabulary of type-level notions; and which parts of them can be declared as
data instead of written as code?

**bittype** is a Go library. It parses, checks, formats and prints the type
syntax of five languages. It resolves and analyses four of them, three
through their real compilers. It generates native type declarations in TypeScript, Go, Rust and
C++ from canonical type declarations written in a language called horos.
Section 1 describes it.

Each language has four packages today that say what exists in it. We call
such a package a **vocabulary**. Two of them, the syntax vocabulary and the
semantic vocabulary, sit beside the processing stages. Two of them sit inside
the directory of one processing stage, generation: the package of intended
type forms, which the code calls the **model**, and the package of external
symbols. Each of the four states the type-level notions of its language in
its own words. Section 2.6 shows the four for Go.

The project's owner asked where these descriptions belong. The owner then
proposed a layout with two roots, `model/` and `ops/`, and one rule between
them. The owner also asked why a language's model needs code at all.

### The framing: the owner's seven messages

The owner asked the question on 13 November 2026 in another working session.
That session relayed the messages to the research owner, the seat that writes
this document. The research owner did not hear them directly. They are quoted
here verbatim and in the order they were asked. They are the framing of this
consultation.

1. > why can't we treat all languages the same, and transform can be from any source to any target? whether the transformation are implemented and exist is a different story, practically it would still be the way you describe, but not because it is baked in, but because it follows from the registered transforms' declarations

2. > Tell the fork we should add maybe a few other language-to-language transforms, not because we want to fully support this now, but just to make sure we don't silently assume this

3. > and this model for target/source language - where does it live?

4. > but why is this in transform? isn't that all just syntax-related? so the model, is it a syntax model, or a transform model?

5. > like, isn't 'transform' just one of many operations working on this model?

6. > But shouldn't we split the models from the operations? The model is independent of the operation, isn't it?

7. A proposal and a second question:

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

Messages 1 and 2 are already decided, by the research owner, as two recorded
decisions on the work item of the contract migration. Section 2.8 gives both.
They are quoted because the later messages follow from them: "this model for
target/source language" in message 3 is the model that a language would need
on either side of a transformation.

Message 7 is **the owner's proposal**. It is the direction this document
examines. The document states it exactly as proposed. It does not argue for
it or against it, and it does not answer its own questions.

### Four questions, as a working frame

The research owner separated the messages into four questions. The separation
is ours and is a working frame. It is no decision.

| Question | What it asks | Asked of the expert |
| --- | --- | --- |
| Grouping | Does a model package hold two things that do not belong together: a vocabulary, and a piece of one generation contribution? | No. Decided as a first step, see below |
| Placement | Do models live apart from operations, with the rule that a model never imports an operation? What is a model: only a language's vocabularies, or also the generic carriers and the value types that today sit with each operation's contract? | Yes, questions 1 and 4 |
| Unification | Do a language's intended forms and its analysed meaning share one vocabulary while keeping their different standing? | Yes, question 2 |
| Declaration | Which parts of a language's model are data, which are code and why, and can a model be declared instead of written? | Yes, question 3 |

**Why grouping is not asked.** Each model package holds form names,
constructors and two checks. Three answers inside those checks say what
today's generation covers. They do not say what a form is. The one production
caller of the checks is the generation target of the same language, and the
Rust target already holds its own coverage answer. The platform's rule on
package boundaries says that parts which do not need each other are separated
first, as a move. We therefore treat grouping as a step that needs no
consultation. Section 6.5 gives the facts and the reason, so that it need not
be answered again. It is our derivation, and if you see that it prejudges one
of the three open questions, please say so.

### How to read the standings

Every statement in this document has exactly one of six standings. Each
section says which standing its statements have, and section 5 lists the
fixed points with their standing.

| Standing | What it means | Weight in the design |
| --- | --- | --- |
| **Owner rule** | A rule or a decision of the project's owner. Quoted, with its date | Fixed |
| **Accepted design** | A design page merged on the main branch, or a row of an earlier expert consultation that its research owner evaluated and recorded as holding. We name the page or the consultation | Fixed unless you see that it contradicts an owner rule, another accepted design or itself. If so, please say so |
| **Implemented today** | What the code on the main branch does. Under the project's rule AC-09 (section 2.1) the implemented surface has no weight in a design choice. It is given as a fact and never as a reason | None |
| **Our derivation** | A reading or a conclusion of the research owner or of the two investigations behind this document. It is in no accepted source | Open. Please challenge it |
| **The owner's proposal** | The layout of message 7. The direction to examine | Open. It is neither a rule nor an accepted design |
| **Open** | Not decided anywhere we looked | Open |

Two consequences, stated once. First, the paths, package names and counts in
sections 2 to 4 describe the present code. None of them is a constraint on
the answer. Second, where a design page and the code disagree, we say so, and
neither is a constraint: section 5.3 lists eight such disagreements.

Quotations of code give the file and the line at the code baseline, which is
bittype's main branch at commit a28febc
(a28febc44830d49ec3c2f22cd70c7ded9d0e5028) of 13 November 2026. All code is
Go unless marked. Quotations of documents give the file. Numbers come from
two investigation notes that read the main branch on 13 November 2026, and we
recounted the totals at the code baseline. A glossary is in section 8.

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
    Each result is published as a graph of facts in one shared carrier. We
    call the Go code that runs a compiler and exports its facts a **bridge**.
  - **Generation.** From analysed horos declarations it generates native type
    declarations in TypeScript, Go, Rust and C++. Every rule that turns one
    kind of horos type into a native type is supplied from outside the shared
    mechanism, by a package that a consumer registers. For TypeScript a later
    stage places the generated declarations into files and writes the
    imports. For Go and Rust that stage is planned and not built.
- A fourth, optional component relates a captured horos program to a captured
  Go program and checks stated representation laws through logos, a separate
  reasoning library of the organization. It appears in this document only as
  one more reader of a language's analysed meaning.
- **Consumers** are builds of the organization's own services and libraries.

The example that the rest of the document reuses is one horos record with one
text field:

```horos
module example.people;
type Person = { name: text; };
```

For Go, generation produces this declaration (read from the code and followed
by hand in section 4; not run for this document):

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

This division is the project's main architectural commitment. Section 2.1
gives the rules behind it, and section 2.5 gives the roles that the
repository's own dependency test enforces.

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
| `parse/`, `check/`, `format/`, `print/` | 22 | 7,049 | The Syntax stages of each language, the structure checks of the generic trees, the shared printers, and the Files stages of horos and TypeScript |
| `resolve/`, `analyze/` | 10 | 716 | The stage adapters around the bridges, and two query services |
| `transform/` | 39 | 3,584 | Generation: the shared kernel parts, one source reader, four targets, four models, three symbol packages, 21 mapping rules and four shared mapping parts |
| `verification/` | 7 | 1,676 | The optional verification component |
| `internal/` | 16 | 6,286 | Private helpers: lexing and parsing mechanics, the grammar table type, per-language spelling rules, the three compiler bridges, the horos adapters |
| `tools/` | 6 | 6,104 | A documentation generator that reads the source, and one horos command |
| **Total** | **129** | **35,943** | |

So the packages that only describe a language are small. The five syntax
vocabularies together have 644 lines, and the four models together have 467.

### 1.4 Who uses it

Nothing is released. The repository has no version tag. Packages outside
`internal/` have public import paths, and the consumers are repositories of
the same organization. An investigation of the checkouts on one workstation
found no other Go module that requires bittype. That is a statement about
those checkouts only.

### 1.5 The two rules that shape everything

Both are rules of the project's owner. Both apply to every project of the
organization, and section 2.1 gives them in full.

- **Architecture is part of correctness** (owner direction of 6 November
  2026). A feature that works and passes its tests can still be incorrect,
  because a decision sits with the wrong owner. The rule AC-01 requires:
  "Satisfy the agreed ownership, dependency and composition contracts as well
  as observable behavior."
- **An interface is derived from the domain and never from the code that
  exists** (owner direction of 8 November 2026, rule AC-09). "The currently
  implemented surface must have **zero influence** on that choice.
  Implementation or migration effort has **zero weight**, including as a
  tiebreaker."

The second rule is why this document can ask its question at all. A layout
that touches every package is a legitimate answer if it is the right one.
