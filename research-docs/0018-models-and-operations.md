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
