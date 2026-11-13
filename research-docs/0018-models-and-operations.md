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
Go unless marked. Quotations of documents give the file. Where a quoted
sentence held a Markdown link, the quotation keeps the link's text and drops
its address. Numbers come from
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

## 2. Architectural context

Each subsection states the standing of what it holds.

### 2.1 Standing architecture rules (AC = architectural correctness)

Standing: owner rules, except where a paragraph says otherwise. They are
written down in two pages of the organization's shared documentation
repository, called the platform. They apply to every project. bittype follows
them and keeps a numbered record of each violation it finds, a **case**
(bittype's cases are named ACX-001 to ACX-010).

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

**The central rule** (owner direction of 6 November 2026):

> A reusable mechanism owns the invariant rules of its operation. Domain-specific
> meanings and policies belong to independently supplied implementations. Within
> the declared extension boundary, adding a new case must require changes only
> to its owning implementation and the consumer's composition, not to shared
> orchestration or existing unrelated implementations.

The same page says what does not settle ownership: "Semantic ownership is
about who defines a meaning and how that meaning enters an operation. It is
not determined by a directory, package name, generic type parameter or the
word "plugin"."

**AC-09, no weight for what exists** (owner direction of 8 November 2026).
Section 1.5 quotes its core. Its application clause says what the existing
code may be used for: "First settle the intended contract. Then use the
existing implementation to plan migration, estimate the work and decide which
code fulfills that contract. Those estimates can guide delivery sequencing
only."

**AC-10, a misplaced responsibility moves first** (owner direction of 11
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

**AC-11, agnostic mechanisms** (owner clarification of 13 November 2026). "A
language-agnostic mechanism belongs in the agnostic layer. Each language
supplies its own types, meanings and policies through the shared contract."
Its requirement: "Place a mechanism according to the knowledge its contract,
algorithm and guarantees require. When these are independent of a language or
domain, implement the mechanism once in the corresponding agnostic layer.
Supply the variation through typed parameters, data or independently owned
implementations." This rule is of the same day as this document. bittype
has recorded one case under it, a forwarding adapter of the horos mapping
side. No page has applied the rule to the question asked here.

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
| Correction | A move | A demand that the shared mechanism resolves |

"The first hides the second. A dependency check reads a unit, usually a
package. Whatever is inside one unit can reach everything else in it unseen.
Correct the grouping first, as AC-10 orders, and the coupling becomes
visible." The page then says when a forbidden-dependency check can enforce a
decided boundary, namely "when each unit has one role", and gives the roles
that section 2.5 shows for bittype. It ends with a limit that question 4 of
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

### 2.2 The generic carriers

Standing: implemented today. Section 3.1 and 3.4 inline the types.

A **carrier** is a generic data type that holds one kind of artifact for
every language. It knows no language. bittype has six carrier packages.

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
the language's checker (section 3.2).

For the semantic graph a language adds names and also Go types: handles that
read the graph in the language's own terms, such as "the fields of this
struct type".

### 2.3 Stages, the two dimensions and the path form

Standing: owner decisions, as the repository's pages record them, and what
is implemented under them.

**Two dimensions.** The page `docs/source.md` records an owner decision of 4
November 2026. It arranges the stages in two parallel rows:

```text
Syntax : CST    : AST
Files  : Source : Semantics
```

| Role | Syntax | Files |
| --- | --- | --- |
| External form | Text | File set |
| Concrete representation | CST | Source |
| Abstract representation | AST | Semantics |
| Parse | Text → CST | Files → Source |
| Format | AST → CST | Semantics → Source |
| Print | CST → Text | Source → Files |

The page adds: "Files and Syntax are processing dimensions, not language
IDs." Resolution and analysis "operate on the semantic side of the Files
dimension". We call the last stage of generation, which places generated
declarations into files, the **Files stage** of a target language.

**The stages** are parse, check, format and print in the Syntax dimension,
parse, format and print in the Files dimension, resolve and analyze, and
transform. The directory `transform/` holds generation. A separate directory
holds verification.

**The path form.** Two owner decisions fix where an implementation lives. No
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
| `format/files/typescript/provider/go` | The TypeScript Files stage |
| `analyze/files/golang/queries/go` | A query service over analysed Go |
| `transform/syntax/golang/target/go` | The Go generation target |
| `transform/syntax/golang/model/go` | Go's intended forms |
| `transform/syntax/horos-to-golang/mapping/record/go` | The mapping rule for horos records into Go |

**What is outside the form.** The carriers (`trees/*/go`, `files/go`,
`source/go`, `semantics/go`), the contract packages (`contracts/*/go`), the
engine, the private helpers under `internal/`, the tools, and two kinds of
vocabulary: `syntax/<language>/go` and `semantics/<language>/go`. The two
other kinds of vocabulary, `model` and `symbol`, are inside the stage
directory `transform/`.

The history of that difference, as an investigation of the repository's
history found it (implemented today, and our reading of the history): the
syntax vocabulary was listed beside the stages from the first day of the
stage-first layout. The semantic graph came on 2 November 2026. The directory `transform/` came on 7 November 2026. The four
model packages were created inside it on 8 November 2026 by the change that
separated the planning of intended types from the construction of trees. That
commit has no message body. Section 5.2 says what no page decides about this.

The owner who decided the stage-first layout is the one who now proposes the
change.

### 2.4 The contract packages

Standing: implemented today. The column "value types" against "protocol" is
our derivation.

A **contract** here is the interface between one operation's shared mechanism
and its independently owned contributions. Each operation has one contract
package. It declares the Go interface that a contribution implements. It also
declares the values that pass through the operation: inputs, limits,
diagnostics, receipts. And in several cases it holds the code that admits a
request and seals a result, because only that code may create a value that
later stages trust.

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
| `contracts/generation/go` | 2,135 | The carrier of intent (forms, expressions, declarations: about 140 lines), the issuers of symbols and values, the description types of a mapping, errors and limits. Together about 470 lines | The target and mapping interfaces, the planner, the freezing of descriptions, the sealing of an output, the walk over the source graph. About 1,650 lines |
| `contracts/verification/go` | 294 | Participant, requirement, site, limits and the association handle | Capture and the check of sites |

Three facts follow, and all three are implemented today.

- No contract package imports a vocabulary of any language.
- Every vocabulary of a language imports at least one contract package. The
  syntax vocabulary imports `contracts/language/go` to state its identifier.
  The semantic vocabulary imports the semantic, resolution and analysis
  contracts. The model and the symbol package import the generation contract,
  because every value they issue is a value of that contract.
- The carrier of intent is not a carrier package. Forms and expressions are
  declared inside the generation contract, beside the planner.

### 2.5 Roles, and the test that derives a role from a path

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
| Contribution: native target | `transform/syntax/<native>/target/go` | The generation contract, the AST, the shared tree builder, and its own syntax vocabulary, model and lexical helper. No source reader and no mapping |
| Contribution: mapping | `transform/syntax/<source>-to-<native>/mapping/<name>/go` | The generation contract, the envelope, the syntax vocabulary of both languages, the source reader of its source language, and the model of its target language. Never another mapping and never a target |
| Private part | `internal/<language>`, `internal/<language>sem` | A lexical helper imports the grammar helper only. A bridge imports foundations and its own syntax and semantic vocabulary |
| Composition | `integration/go`, `verification/integration/go`, `tools/**` | Chooses and registers contributions |

Two sentences of the layout page bind every role:

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
  provider. `transform/syntax/<language>/model/go` is a model. The language
  is read from the path and substituted into the rule's import patterns. A
  new package of a known shape has its rule without an edit to the test.
  Section 3.8 inlines this function.
- *From a closed table.* Foundation, kernel and shared helper packages have
  one entry each (31 entries). Packages of which there is only one have a
  hand entry (12 entries). A table of **role extensions** names, per package,
  the imports that one language needs and its role does not give (12
  entries). A table of exceptions names, per importer, the exact symbols it
  may use from a contribution (14 entries).

So a path is not only a location. It is the input from which the import rule
is computed. Any move of a package changes its rule unless the function and
the tables change with it. That is a fact of today's test and no constraint.
Question 4 asks how such a rule should be enforced.

### 2.6 The four vocabularies of one language, with Go as the example

Standing: implemented today.

| | Syntax vocabulary | Semantic vocabulary | Model | Symbol package |
| --- | --- | --- | --- | --- |
| Path | `syntax/golang/go` | `semantics/golang/go` | `transform/syntax/golang/model/go` | `transform/syntax/golang/symbol/go` |
| Lines | 107 | 855 | 112 | 81 |
| What it states | What a Go type-syntax tree is made of: 23 node kinds, 15 containment roles, 1 reference role, 9 attribute names, 4 concrete kinds | What analysis of a Go program reports: kinds of entities, relations of facts, operators of terms, roles of bindings, all as keys with the owner `go` | What generation intends to construct: 7 form names and 2 primitive names | What a name use outside the generated declarations refers to: one category, a predeclared type |
| What else it holds | Nothing. Two one-line methods return the language's identifier and display name | Profile strings, the identity values of the resolution and analysis contracts, the types of a compiler environment with their validation and encoding, and typed readers over an analysed graph | Five constructors that build forms, and two checks of a form | The rule that admits a name as a predeclared type, and two functions that record a symbol at a node and read it back |
| Imports of this module | The AST, the CST, `contracts/language/go` | The semantic graph, the AST, the envelope, the syntax vocabulary, and the semantic, resolution and analysis contracts | The generation contract and the syntax vocabulary | The generation contract, the semantic graph, the AST, the envelope and the syntax vocabulary |
| Who writes its names into an artifact | The parser and the generation target (into trees) | The Go bridge (into the graph) | The mappings into Go (into intent) | The generation target (into the sealed output) |
| Who reads them | The checker, the formatter and the bridge | The typed readers in the same package, the query service, the verification component | The generation target | Only the generation target today. The planned Go Files stage is to read them. For TypeScript, the Files stage already reads its symbol package |
| Stages that import it | parse, check, format, resolve, analyze, transform | resolve, analyze, verification | transform only | transform only (for TypeScript also format/files) |

Three things that are not vocabularies complete the picture of Go.

- **The grammar table.** Which attributes and children each node kind has is
  stated in the Go checker, a contribution (section 3.2). The syntax
  vocabulary holds the names only.
- **Lexical rules.** Go has no lexical helper of its own. The checker asks
  the Go standard library whether a word is a keyword, and the symbol package
  asks it which names the universe scope declares. TypeScript, Rust, C++ and
  horos each have a private package of spelling rules.
- **The bridge.** A private package of 1,241 lines runs the Go type checker
  over captured source and writes the graph. It is the writer of every name
  in the semantic vocabulary.

The other languages have the same four, with two gaps. C++ has no semantic
vocabulary and no symbol package, because it has no analysis. horos has no
model and no symbol package, because it is a source of generation and never
a target. Its counterpart on the generation side is a source reader, which
reads horos's semantic vocabulary.

### 2.7 Generation

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
   |  the Files stage places declarations under a consumer's placement
   v
placed source, analysed again by the language's bridge  <- destinations compared
```

The participants and their words:

- **Source reader.** A contribution of the source language. It tells the
  kernel which entities are declarations and gives mappings a typed reading
  of a use, for example the fields of a record. Only horos has one.
- **Mapping.** A contribution for one ordered pair of languages and one kind
  of source type, for example horos records into Go. It states what the
  target type is to be. It uses the constructors of the target's model and
  never builds a tree. There are 21, each in its own package of 25 to 78
  lines.
- **Family.** One kind of source type together with its representation in a
  target: text, booleans, lists, closed records, optional fields, named
  declarations.
- **Form.** A constructor of intent that one target language issues, such as
  Go's `slice`, `record` and `field`.
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
- **Files stage.** A contribution of the target language in the Files
  dimension. It decides files, imports, qualification and visibility under a
  placement that the consumer supplies. It then has the placed program
  analysed and compares each name use with the symbol or declaration it was
  meant to reach. It exists for TypeScript. For Go and Rust it is planned.

The central sentence on intent, from `docs/generation.md`:

> Build receives no source reader or mapping callback. Native intent is not an
> analyzed `Semantics[T]`: target analysis establishes facts about the actual
> constructed/organized program.

The same page lists what each component may know. Its row for the "Native
generation model/target" reads: "Its native forms and syntax; no
Horos/library semantics". And `docs/kernel-composition.md`: "The native model
describes intended types; the target constructs their syntax. Analyzed native
facts remain under semantics."

The pages name the model package in five ways: "native semantic forms and
validation", "intended native forms", "native intent vocabulary", "native
semantic forms" and "Native form identity and validity, without AST syntax".
None calls it a syntax model or a transform model, which are the two words of
the owner's message 4.

### 2.8 What is in flight

Standing: accepted designs and decisions of the research owner, as marked.
Nothing in this subsection is on the main branch unless it says so.

**The migration of the mapping contract.** A design of 13 November 2026 (the
page `docs/contribution-contract.md`, from the expert consultation 0017,
decided by the research owner and final only when a stated gate is met)
replaces what a mapping registers and returns. Today a mapping returns
intent. Under the new contract a mapping returns one **decision record** per
selection: a typed value that holds every decision the selection contains,
from which the intended form is computed. The record's type is declared in
horos, and its Go type is generated by bittype itself (section 3.10). The
reason is a second reader: a separately owned backend must emit conversion
code from the same decisions.

The migration is being built on a branch. At its present head (c1e5069, four
commits ahead of the code baseline and three behind) it adds, per ordered
pair and family, two packages beside each mapping, and one more per family:

| New package on the branch | What it holds |
| --- | --- |
| `transform/syntax/horos-to-<target>/vocabulary/<family>/go` | The identity of one contribution, and the schema and generated Go type of its decision record |
| `transform/syntax/horos-to-<target>/codec/<family>/go` | The hand-written encoding of that record |
| `transform/syntax/horos/family/<family>/go` | The identity of the family's meaning, the same for every target. Its package comment calls it "the vocabulary of the family" |

So the branch adds vocabularies of a new kind under the stage directory
`transform/`: names that belong to a relation between two languages and not
to one language. The layout page says of them: "This layout has no role for a
facet package of a mapping contribution, and none for a vocabulary package of
one family for one target: its role "Vocabulary" is per language." The
migration moves no model package. It rewrites every mapping, and every
mapping calls the model's constructors.

**No direction in the shared contract** (decision O42 of the research owner,
13 November 2026, on the owner's message 1). The shared generation contract
called the target side "native". The decision: "The shared part names two
roles of a transformation: source and target." The rename is done on the
branch. Which pairs of languages exist follows from what is registered.

**Two further pairs** (decision O43 of the research owner, 13 November 2026,
on the owner's message 2). Two small pairs are to run through the new
contract before it is final, to show that no direction is assumed:

- `golang-to-typescript`: a source that is not horos. It needs a Go source
  reader over the Go analysis, which no language but horos has today.
- `typescript-to-horos`: horos as a target. It needs a model and a target for
  horos, which do not exist.

Both are to be written against the present places. Under the present path
form a horos model would be `transform/syntax/horos/model/go` and a Go source
reader `transform/syntax/golang/source/go`. That reading of the path form is
ours.

**Two Files stages.** The Go and the Rust Files stage are planned as the next
slices of two open work items. By the accepted specification of destinations
(the page `docs/destinations.md`, 13 November 2026) each is to read its
language's symbol package, as TypeScript's does. Under the present path form
each would be a package under `format/files/`.

These four pieces of work are why question 4 asks about order.

## 3. Relevant code

Standing of this section: implemented today. It is evidence of the present
state and never a reason for or against an option. Where a table classifies
code, the classification is our derivation, and the text says so.

### 3.0 Reading guide

Every snippet is taken from the main branch at the code baseline, with its
path and line. Comments are the code's own unless a line starts with
`// [doc]`, which marks a comment added for this document. An ellipsis `...`
marks omitted lines.

Import aliases recur. In every language's packages `s` is that language's
syntax vocabulary. `g` is the grammar helper `internal/grammar` inside a
checker, and the generation contract `contracts/generation/go` inside a
model, a symbol package, a target and a mapping. The two never meet in one
file. `ast` is the generic tree package. `semantics` is the generic semantic
graph. `text` is the diagnostics package, whose function
`text.AtNode(tree, id, message)` returns an error value that points at one
node. `tree.View` is a small read helper over one tree: `Node(id)` returns a
node, `Attr(id, name)` an attribute's text, `Has` and `Flag` test an
attribute, and `Children(id, role)` and `Child(id, role)` return the children
under one role.

### 3.1 The generic tree, and Go's syntax vocabulary

The abstract tree, complete except for comments (`trees/ast/go`, files
`ast.go`, `node.go`, `origin.go`, `view.go`):

```go
// trees/ast/go/ast.go:13-35, node.go:6-44, origin.go:7-32, view.go:7-11
type TreeID string
type NodeID uint32 // [doc] a one-based index into Nodes; zero names no node

type AST[L any] struct {
	ID    TreeID
	Root  NodeID
	Nodes []Node[L]
}

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

type NodeRef struct {
	Tree TreeID
	Node NodeID
}
type SourceID string
type Span struct {
	Source SourceID
	Start  uint64
	End    uint64
}
type Origin struct {
	Spans []Span
	Nodes []NodeRef
}

type View[L any] interface {
	ASTView() AST[L]
}
func (tree AST[L]) ASTView() AST[L] { return tree }
```

The package comment states what the carrier is and is not
(`trees/ast/go/ast.go:3-8`):

```go
// The language parameter separates ASTs and node kinds in Go's type system.
// A language package supplies its own named marker type, kinds and grammar;
// this package contains no language catalogue or grammar implementation.
//
// These are data types, not checked trees. Structural and grammatical checks
// live in separate packages; the carrier performs no validation or allocation.
```

And the comments on a node's parts say who owns their meaning
(`node.go:8-10`, `:18-19`, `:25-27`): "Attribute is an uninterpreted textual
syntax atom." "A grammar defines the name, allowed text and meaning." "The
grammar owns Role names and cardinality; repeated roles can form a list." "A
grammar owns Role and the allowed target."

The concrete tree, complete (`trees/cst/go/cst.go:7-30`):

```go
type NodeID uint32
type Kind[L any] string

// Node is either a leaf (Text != nil, no Children) or a container (Text == nil).
// Kind belongs to the language. Leaves include punctuation, comments and spaces.
// Syntax addresses zero or more nodes of this CST's stored abstract view.
type Node[L any] struct {
	Kind     Kind[L]
	Text     *string
	Children []NodeID
	Syntax   []ast.NodeID
	Origin   ast.Origin
}

// CST owns both arenas. Stage inputs and borrowed views must not be mutated.
type CST[L any] struct {
	Abstract ast.AST[L]
	Root     NodeID
	Nodes    []Node[L]
}

func (tree CST[L]) ASTView() ast.AST[L] { return tree.Abstract }
```

A language's identity is two methods of an interface in a contract package
(`contracts/language/go/language.go:11`, `:22-25`):

```go
type ID string

type Provider interface {
	ID() ID
	Name() string
}
```

Go's syntax vocabulary, with its full list of kinds
(`syntax/golang/go/nodes.go:1-95`, abridged: the imports and some comments
are left out, and the roles and attributes are shortened as marked):

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
	KindImport            Kind = "import"
	KindTypeDeclaration   Kind = "type-declaration"
	KindTypeAlias         Kind = "type-alias"
	KindTypeParameter     Kind = "type-parameter"
	KindNamedType         Kind = "named-type"
	KindQualifiedType     Kind = "qualified-type"
	KindTypeApplication   Kind = "type-application"
	KindStructType        Kind = "struct-type"
	KindField             Kind = "field"
	KindEmbeddedField     Kind = "embedded-field"
	KindInterfaceType     Kind = "interface-type"
	KindMethod            Kind = "method"
	KindEmbeddedType      Kind = "embedded-type"
	KindFunctionType      Kind = "function-type" // Native signature syntax; no library Function semantics.
	KindParameter         Kind = "parameter"
	KindPointerType       Kind = "pointer-type"
	KindSliceType         Kind = "slice-type"
	KindArrayType         Kind = "array-type"
	KindMapType           Kind = "map-type"
	KindChannelType       Kind = "channel-type"
	KindTypeUnion         Kind = "type-union" // A constraint union, not a value union.
	KindTypeApproximation Kind = "type-approximation"
)

// Containment roles. Parameters and results are separate ordered lists.
const (
	RoleImport        = "import"
	RoleDeclaration   = "declaration"
	RoleType          = "type"
	// [doc] and RoleTypeParameter, RoleConstraint, RoleConstructor, RoleArgument,
	// [doc] RoleField, RoleMember, RoleSignature, RoleParameter, RoleResult,
	// [doc] RoleElement, RoleKey, RoleTerm: 15 in all
)

// Reference roles. Target links name declarations, imports or lexical binders.
const (
	RoleTarget = "target"
)

// Attribute names. Package qualifiers are identifiers, not dot-joined type source.
const (
	AttributePackage = "package"
	AttributeName    = "name"
	// [doc] and AttributePath, AttributeAlias, AttributeQualifier, AttributeTag,
	// [doc] AttributeLength, AttributeDirection, AttributeVariadic: 9 in all
)
```

A second file of the package declares the concrete kinds
(`syntax/golang/go/concrete.go:7-12`). Every language declares the same four:

```go
const (
	ConcreteLayout     cst.Kind[Language] = "layout"
	ConcreteToken      cst.Kind[Language] = "token"
	ConcreteWhitespace cst.Kind[Language] = "whitespace"
	ConcreteComment    cst.Kind[Language] = "comment"
)
```

What this shows, as facts:

- A node kind is a typed constant. A role name and an attribute name are
  untyped strings.
- The package states names and nothing about shape. It does not say that a
  `field` has a `name` attribute and one `type` child. It points to prose
  for that.
- The other four syntax vocabularies have the same form. They declare 33
  (horos), 29 (TypeScript), 36 (Rust) and 23 (C++) node kinds.
- The four native parsers and formatters do not use the role and attribute
  constants. They write the names as string literals, and the formatters also
  switch on node kinds as string literals, for example
  `case "named-type":` (`format/syntax/golang/provider/go/format.go:101`).
  Only the horos parser and formatter use the constants throughout. So the
  vocabulary package is not the one place where these names are spelled.

### 3.2 The admission table, and the rules beside it

**The table type.** A private helper package declares what a grammar table
is and holds its one interpreter (`internal/grammar/rules.go:15-40`):

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
attribute map.

The interpreter is one function
(`internal/grammar/rules.go:80`):

```go
func Check[L any](tree ast.AST[L], root ast.Kind[L], rules Rules[L]) error
```

In order it checks: the structure of the carrier, by a call into
`check/syntax/shared/ast/go`; the kind of the root; and for every node, that
its kind has a rule, that every attribute is declared and accepted, that
every required attribute is present, that every child's role is declared and
that children come in the declared order of roles, that each child's kind is
in its role's list, that each role's count is within its bounds, and that a
reference is present only where the rule allows one and only under the role
`target`.

So one entry can express, per node kind: a set of attributes, each required
or optional, each with a predicate on its text; an ordered list of child
roles, each with a count range and a flat list of allowed child kinds; and
whether a reference may be attached.

**One table in full.** Go's is a package variable of the Go checker
(`check/syntax/golang/provider/go/check.go:13-41`). Here `g` is the grammar
helper:

```go
func identifier(v string) bool { return g.Identifier(v) && !token.Lookup(v).IsKeyword() }

var types = []s.Kind{s.KindNamedType, s.KindQualifiedType, s.KindTypeApplication, s.KindStructType, s.KindInterfaceType, s.KindFunctionType, s.KindPointerType, s.KindSliceType, s.KindArrayType, s.KindMapType, s.KindChannelType}
var constraints = append(append([]s.Kind{}, types...), s.KindTypeUnion, s.KindTypeApproximation)
var rules = g.Rules[s.Language]{
	s.KindFile:              {Attributes: g.Attributes(s.AttributePackage, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.Many(s.RoleImport, 0, s.KindImport), g.Many(s.RoleDeclaration, 0, s.KindTypeDeclaration, s.KindTypeAlias)}},
	s.KindImport:            {Attributes: g.Attributes(s.AttributePath, g.Attr(true, g.Nonempty), s.AttributeAlias, g.Attr(false, func(v string) bool { return v == "." || identifier(v) }))},
	s.KindTypeDeclaration:   {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.Many(s.RoleTypeParameter, 0, s.KindTypeParameter), g.One(s.RoleType, types...)}},
	s.KindTypeAlias:         {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.Many(s.RoleTypeParameter, 0, s.KindTypeParameter), g.One(s.RoleType, types...)}},
	s.KindTypeParameter:     {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.One(s.RoleConstraint, constraints...)}},
	s.KindNamedType:         {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier)), Reference: true},
	s.KindQualifiedType:     {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier), s.AttributeQualifier, g.Attr(true, identifier)), Reference: true},
	s.KindTypeApplication:   {Roles: []g.Role[s.Language]{g.One(s.RoleConstructor, s.KindNamedType, s.KindQualifiedType), g.Many(s.RoleArgument, 1, types...)}},
	s.KindStructType:        {Roles: []g.Role[s.Language]{g.Many(s.RoleField, 0, s.KindField, s.KindEmbeddedField)}},
	s.KindField:             {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier), s.AttributeTag, g.Attr(false, g.String)), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	s.KindEmbeddedField:     {Attributes: g.Attributes(s.AttributeTag, g.Attr(false, g.String)), Roles: []g.Role[s.Language]{g.One(s.RoleType, s.KindNamedType, s.KindQualifiedType, s.KindTypeApplication, s.KindPointerType)}},
	s.KindInterfaceType:     {Roles: []g.Role[s.Language]{g.Many(s.RoleMember, 0, s.KindMethod, s.KindEmbeddedType)}},
	s.KindMethod:            {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.One(s.RoleSignature, s.KindFunctionType)}},
	s.KindEmbeddedType:      {Roles: []g.Role[s.Language]{g.One(s.RoleType, constraints...)}},
	s.KindFunctionType:      {Roles: []g.Role[s.Language]{g.Many(s.RoleParameter, 0, s.KindParameter), g.Many(s.RoleResult, 0, s.KindParameter)}},
	s.KindParameter:         {Attributes: g.Attributes(s.AttributeName, g.Attr(false, identifier), s.AttributeVariadic, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	s.KindPointerType:       {Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindSliceType:         {Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindArrayType:         {Attributes: g.Attributes(s.AttributeLength, g.Attr(true, g.Decimal)), Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindMapType:           {Roles: []g.Role[s.Language]{g.One(s.RoleKey, types...), g.One(s.RoleElement, types...)}},
	s.KindChannelType:       {Attributes: g.Attributes(s.AttributeDirection, g.Attr(false, g.Set("both", "send", "receive"))), Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindTypeUnion:         {Roles: []g.Role[s.Language]{g.Many(s.RoleTerm, 2, append(types, s.KindTypeApproximation)...)}},
	s.KindTypeApproximation: {Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
}
```

Facts about the tables as data:

- Each language's table has one entry per node kind constant: 33, 29, 23, 36
  and 23 entries.
- The table is an unexported variable of a contribution, the checker. No
  other package can read it. The parser and the formatter of the language
  call the checker's exported function `Check` and nothing else of it.
- The table is a Go map literal that holds function values. In Go's table
  they are `identifier`, which asks the standard library whether a word is a
  keyword, and an inline function for the import alias. The other languages
  pass predicates of their private spelling packages.
- Sets of kinds are named Go slices (`types`, `constraints`). The Rust and
  C++ tables also build shared parts of entries with small functions.
- A documentation tool reads the tables as data. Because the value holds
  functions and is unexported, the tool does not import it. It parses the Go
  source of the literal and evaluates the initializer syntactically. A
  comment marks the table for it: `// bittype:grammar rules`.

**The rules beside the table.** Each checker runs the table and then code.
Go's code holds three rules (`check/syntax/golang/provider/go/check.go:45-83`):

```go
func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, s.KindFile, rules); err != nil {
		return err
	}
	v := tree.View[s.Language]{Tree: t}
	for i, n := range t.Nodes {
		id := ast.NodeID(i + 1)
		if n.Kind == s.KindFunctionType {
			for _, role := range []string{s.RoleParameter, s.RoleResult} {
				params := v.Children(id, role)
				named := -1
				for j, p := range params {
					mode := 0
					if v.Has(p, s.AttributeName) {
						mode = 1
					}
					if named >= 0 && named != mode {
						return text.AtNode(t, p, "cannot mix named and unnamed parameters")
					}
					named = mode
					if v.Flag(p, s.AttributeVariadic) && (role != s.RoleParameter || j != len(params)-1) {
						return text.AtNode(t, p, "variadic parameter must be final input")
					}
				}
			}
		}
		if n.Kind == s.KindEmbeddedField {
			e := v.Child(id, s.RoleType)
			if v.Node(e).Kind == s.KindPointerType {
				kind := v.Node(v.Child(e, s.RoleElement)).Kind
				if kind != s.KindNamedType && kind != s.KindQualifiedType && kind != s.KindTypeApplication {
					return text.AtNode(t, e, "embedded pointer requires a named type")
				}
			}
		}
	}
	return nil
}
```

Three rules from other checkers, each with what it expresses that the table
type cannot:

```go
// check/syntax/typescript/provider/go/check.go:60-72
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
// check/syntax/rust/provider/go/check.go:100-107
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
// check/syntax/cpp/provider/go/check.go:72-80
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

The count of such rules per checker, with the kinds of thing they express
(the grouping into kinds is our derivation):

| Language | Table entries | Rules in code | What the rules in code express |
| --- | --- | --- | --- |
| Go | 23 | 3 | Consistency among siblings of one role (all parameters named or none); position within a role (variadic last); a condition two levels down (the element of an embedded pointer) |
| TypeScript | 29 | 14 | A condition on ancestors; an attribute against another attribute, with the lexical rules of literals; an attribute against a child's kind; order and position within a role (rest and optional elements); consistency among siblings (tuple labels) |
| Rust | 36 | 13 | A condition that depends on the parent's kind and the child's role; order within a role (lifetimes before types, arguments before bindings); what a kind of declaration inherits (visibility, defaults) |
| C++ | 23 | 11 | A condition on the parent; an attribute against another attribute; which kinds may nest (no array of functions, no reference to a reference) |
| horos | 33 | 11 | Which reference roles each kind may carry (the table knows one role, `target`); uniqueness of names among siblings; contiguous indices; that an index refers to a binder in scope, through a walk that carries lexical depth |

The horos checker removes all references from a copy of the tree before it
runs the table, because "The shared grammar checker supports the native
target role. Horos has several annotation roles, checked below"
(`internal/declsyntax/check.go:127-128`).

**What no table or rule states.** The concrete syntax. Four of the five
parsers are hand-written recursive descent over a shared scanner. The Go
parser calls the parser of the Go standard library and converts its tree.
Every formatter is a hand-written function per node kind. A parser calls the
checker after it has built its tree, and a formatter calls it before it
reads its input. That call is the only connection between the table and the
code that produces and consumes trees.

### 3.3 The semantic graph, Go's semantic vocabulary and one typed reader

**The graph.** One package holds the carrier of analysed meaning for every
language (`semantics/go/model.go:19-125`, `:146`, `:230`; comments abridged):

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

// Parameter retains a real binder and its coarse provider-owned sort.
type Parameter struct {
	Entity EntityID
	Sort   Key
}

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

type Operand struct {
	Locator []uint32
	Term    TermID
}

// Occurrence identifies a semantic position within a current syntax node.
type Occurrence struct {
	Site    ast.NodeRef
	Role    Key
	Ordinal uint32
}

type Binding struct {
	At     Occurrence
	Target EntityID
}

type Expression struct {
	At   Occurrence
	Term TermID
}

type Dependency struct {
	From     EntityID
	To       EntityID
	Relation Key
	At       *Occurrence
}

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

// Data is construction storage. New bounds, validates and copies it before
// publication.
type Data struct {
	Entities     []Entity
	Terms        []Term
	Bindings     []Binding
	Expressions  []Expression
	Dependencies []Dependency
	Facts        []Fact
}

// Semantics[L] is an immutable owned graph.
type Semantics[L any] struct{ state *Data }

func New[L any](data Data, limits Limits, contains func(ast.NodeRef) bool) (Semantics[L], error)
```

The graph fixes three term variants and the record shapes. Everything else
is a key that a provider owns: the kind of an entity, the sort of a
parameter, the role of an occurrence, the relation of a dependency or a
fact, the operator and the encoding of a provider expression. The
constructor `New` checks integrity and bounds. It does not know what any key
means, and it does not look into `Data` bytes.

**Go's names.** The semantic vocabulary of Go declares its keys as one-line
functions, because a Go struct value cannot be a constant
(`semantics/golang/go/vocabulary.go:13-71`, abridged):

```go
const ResolutionProfile = "go1.25/captured-types-1/resolution-1"
const AnalysisProfile = "go1.25/captured-types-1/analysis-1"

// Key returns a Go-owned vocabulary key; it performs no registration.
func Key(name string) semantics.Key { return semantics.Key{Owner: "go", Name: name} }

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

// Operators of the terms of a basic, a slice and a struct type.
func BasicOperator() semantics.Key  { return Key("basic") }
func SliceOperator() semantics.Key  { return Key("slice") }
func StructOperator() semantics.Key { return Key("struct") }
```

The bridge writes more names than the vocabulary publishes. It spells these
as literals in its own code: the entity kind `field`, the parameter sort
`type`, the operators `interface`, `boolean`, `union`, `signature`, `tuple`
and `display`, the binding role `import`, the expression role
`type-expression` and the encoding `json-1`
(`internal/golangsem/export.go:111`, `:135`, `:169`, `:253-288`, `:347`,
`:371`). The accepted rule for when a name gets a published declaration is:
"A name needs a single owner when two parties must agree on it"
(`docs/package-boundaries.md`). No production package other than the bridge
reads one of these names by name.

**A fact's shape exists in two pieces of code and in no declaration.** The
bridge writes the term of a struct type
(`internal/golangsem/export.go:233-243`; `n` is the semantic vocabulary,
`e.typ` exports a type and returns its term, `e.expression` builds a
provider expression whose data is the JSON encoding of its second argument):

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
(`semantics/golang/go/model.go:20-39`, `:276-302`). `t.term()` returns the
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

func (t Type) SliceElement() (Type, bool) {
	v := t.term()
	if v.Operator != SliceOperator() || len(v.Operands) != 1 {
		return Type{}, false
	}
	return Type{t.owner, v.Operands[0].Term}, true
}
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

So "a struct term has the operator `struct`, JSON field records as data and
one operand per field" is stated once by the writer and once by the reader.
The operator's name is declared. Its shape is not.

**What else the package holds.** A reader is opened from a **receipt**, the
value that the semantic contract issues for one finished analysis
(`semantics/golang/go/model.go:41-56`):

```go
// Open requires an admitted analysis receipt and never invokes a compiler.
func Open(result semantic.Result) (Model, error) {
	requirement, err := EvidenceFor(result.Configuration())
	if err != nil {
		return Model{}, err
	}
	if _, err := semantic.AdmitEvidence(result.Snapshot(), []semantic.EvidenceRequirement{requirement}, []semantic.Result{result}); err != nil {
		return Model{}, err
	}
	envelope, _ := result.Semantics()
	graph, err := representation.UnpackSemantics[s.Language](result.Snapshot().Language(), envelope)
	if err != nil {
		return Model{}, err
	}
	return Model{&model{result, graph}}, nil
}
```

`semantic.Result` is the receipt. `semantic.Snapshot` is the captured input
that it belongs to. `semantic.EvidenceRequirement` says which operation,
profile, configuration and reports a consumer requires
(`contracts/semantic/go/evidence.go:14-20`):

```go
type EvidenceRequirement struct {
	Operation     semantics.Key
	Profile       string
	Configuration []Setting
	Reports       []Requirement
	Semantics     bool
}
```

with `Setting{Key semantics.Key; Value []byte}` and
`Requirement{Obligation semantics.Key; Verdict Verdict}` from the same
package. `semantic.AdmitEvidence` checks that the receipt meets the
requirement. `representation.UnpackSemantics` opens the type-erased
envelope.

The vocabulary package also builds the identity values of its contracts from
those types (`semantics/golang/go/contracts.go:10-21`):

```go
// The contract identities are what a consumer registers the Go providers under
// and what a dependent names to require them. They are published here, with the
// profiles and keys they are built from, so that naming one imports no provider.

func ResolutionContract() resolution.Contract {
	c := semantic.Contract{Profile: ResolutionProfile, SemanticsRequirements: Requirements()}
	for _, key := range []semantics.Key{Coverage(), Bindings(), Validity()} {
		c.Obligations = append(c.Obligations, semantic.Obligation{Key: key, Domain: Key("captured-package")})
	}
	return resolution.Contract{Semantic: c, Bindings: Bindings()}
}
```

These two snippets show why a semantic vocabulary imports three contract
packages today: its readers take a receipt, and its identities are values of
contract types. It holds 855 lines for Go, of which 99 are names, 40 contract
identities, 400 the types of the compiler environment with their validation
and encoding, and 316 the typed readers.

### 3.4 Intent: forms, expressions, and Go's whole model

**The carrier of intent** is declared inside the generation contract
(`contracts/generation/go/expression.go:9-106`, abridged to what forms and
native expressions need):

```go
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

type ExpressionKind uint8

const (
	NativeExpression ExpressionKind = iota + 1
	CapturedExpression
	ParameterExpression
	ArgumentExpression
	HelperExpression
	DeclarationExpression
)

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

func Native[T any](form Form[T], data []string, args ...Expression[T]) Expression[T] {
	return Expression[T]{kind: NativeExpression, form: form, data: slices.Clone(data), args: slices.Clone(args)}
}
func Mapped[T any](slot Slot) Expression[T] {
	return Expression[T]{kind: CapturedExpression, slot: slot}
}
func Parameter[T any](index int) Expression[T] {
	return Expression[T]{kind: ParameterExpression, index: index}
}
func (e Expression[T]) Form() Form[T]              { return e.form }
func (e Expression[T]) Data() []string             { return slices.Clone(e.data) }
func (e Expression[T]) Arguments() []Expression[T] { return slices.Clone(e.args) }

// identity is an allocation-unique owner for issuer-scoped values.
type identity struct{ value byte }

func fresh() *identity { return &identity{value: 1} }
```

`T` is the marker type of the target language. `Slot` and `HelperKey` are
two more issued identities of the same construction: a slot stands for the
representation of another use, and a helper key names a declaration that a
mapping demands. Of the six expression kinds only the first holds a form.
The other five are the kernel's own: a captured slot, a type parameter, a
generic argument, a helper reference and a reference to a planned
declaration. So type parameters, references and applications are not forms
of any language. They belong to the shared contract.

An intended declaration and the interface of a target
(`contracts/generation/go/native.go:13-24`, `contracts.go:13-42`):

```go
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
with its declared sites. The two errors that matter below are `ErrInvalid`,
which says a description is malformed, and `ErrUnsupported`, which says a
capability is not available. They are two different answers: the first is
about validity, the second about coverage.

**Go's model, complete** (`transform/syntax/golang/model/go/model.go:1-112`).
Here `g` is the generation contract. The comments marked `[doc]` classify
each check as **WF** (well-formedness: what a form is) or **COV** (a coverage
answer: what today's generation supports). The classification is our
derivation.

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

Facts about the four models:

- They are near copies. Each has 112 to 127 lines. They differ in the name
  of the collection form (`slice` in Go, `array` in TypeScript, `sequence` in
  Rust and C++), in the name of the nominal declaration form (`defined`,
  `interface`, `struct`, `struct`) and in the external form.
- No package is shared between them, by a decision recorded in the change
  log on 12 November 2026: "Each model owns its own forms, a form can be read
  only by its issuer, and no package compares a form name of one language
  with another's or with a literal."
- A form's shape is stated twice in code and in no table: once in the
  constructor that writes the data strings and arguments, and once in the
  check that counts them. A third reader is in the target, which reads the
  data by position (below).
- `Form(name)` issues a form of any name. The function does not check the
  name against the constants.
- The Go model exports a constructor, `External`, whose result its own check
  always refuses.
- No model has a form for a pointer, a map, an integer, a union or an
  optional type. The syntax vocabulary of every language has node kinds for
  such types (section 3.1).

**Who writes forms.** A mapping. The mapping for horos text into Go, complete
except for its imports
(`transform/syntax/horos-to-golang/mapping/text/go/text.go:16-30`). Here `h`
is the syntax vocabulary of horos, `rules.Library` is the function type of a
horos mapping, `patterns.Head` builds the selector of one built-in horos
constructor, and `source.Reader` is the reading interface of the horos
source reader:

```go
func primitive(sourceName, targetName string) rules.Library[s.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := patterns.Head(from, "scalar", sourceName)
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Profile: rules.RepresentationProfile, Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: sourceName + "/golang/1", Input: pattern, Output: model.Primitive(targetName)}}}, nil
	}
}

// Library returns the text mapping. It states the structural representation
// profile itself.
func Library() rules.Library[s.Language] {
	return primitive("utf8-text", model.PrimitiveString)
}
```

The description types it fills, from the generation contract
(`contracts/generation/go/description.go:21-31`, `:81-86`):

```go
type Pattern[F, R any] struct {
	Selector Selector
	Slots    []Slot
	Capture  func(context.Context, R, UseView[F]) ([]Capture[F], error)
}
type BindingDescription[F, T, R any] struct {
	ID     string
	Input  Pattern[F, R]
	Output Expression[T]
	Expand func(context.Context, Expansion[F, T, R]) (Expression[T], error)
}
type MappingDescription[F, T, R any] struct {
	Profile      string
	Bindings     []BindingDescription[F, T, R]
	Declarations []DeclarationRecipe[F, T, R]
	Helpers      []HelperRecipe[T]
}
```

`F` and `T` are the marker types of the source and the target language, and
`R` is the type of the source reader. A `Selector` is an identity that the
source language issues for one kind of use. A **profile** is a named promise
about the generated types; its identity is a string. A binding gives its
result either as a constant `Output` or through a callback `Expand`. The
record mapping uses the callback: it reads the fields, asks the kernel for
the representation of each field's type through `r.Mapped(f.Use)`, and
returns `model.Record(model.Field(name, optional, type), ...)`. A
`DeclarationRecipe` does the same for a named declaration, and a
`HelperRecipe` supplies a declaration under a helper key.

So a mapping imports the model for its constructors and its names. No
mapping calls a model's check.

**Who reads forms.** The target of the same language. Go's target forwards
the two checks and switches over the form names when it builds
(`transform/syntax/golang/target/go/provider.go:24-29`, `:78-95`,
`:111-128`). Here `m` is the shared tree builder: `b.Node` allocates one
node with its children and returns a handle to it, a `m.Fragment`;
`b.External` attaches an external use to a node; and `edges` pairs each
fragment of a list with one role:

```go
func (p Provider) CheckExpression(ctx context.Context, e g.Expression[s.Language]) error {
	return model.CheckExpression(ctx, e)
}
func (p Provider) CheckDeclaration(ctx context.Context, d g.NativeDeclaration[s.Language]) error {
	return model.CheckDeclaration(ctx, d)
}

// predeclared writes a use of a type of the universe scope. The identifier is
// a name that a declaration or a type parameter can capture, so the use
// records its destination: the predeclared symbol.
func (b *Backend) predeclared(name string, origin ast.NodeRef) (m.Fragment[s.Language], error) {
	destination, err := symbol.Predeclared(name)
	if err != nil {
		return m.Fragment[s.Language]{}, err
	}
	node, err := b.Node(m.Node(s.KindNamedType, m.Attributes("name", name), origin))
	if err != nil {
		return node, err
	}
	use, err := symbol.ExternalUse(b.language, ast.NodeRef{}, destination)
	if err != nil {
		return node, err
	}
	return b.External(node, use)
}

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
the model's checks, and it calls no constructor. The target is where a form
meets a node kind: `slice` becomes `slice-type`, `record` becomes
`struct-type`, `field` becomes `field`. And the Go target reads `data[0]` of
a field form and ignores `data[1]`, the optional flag: Go has no optional
field, and the record mapping wraps the type of an optional field in a
helper type instead.

**Where Rust keeps its coverage answer.** Rust's model checks only the shape
of an external form, and says so
(`transform/syntax/rust/model/go/model.go:56-64`, `:100-103`):

```go
// External is the public item at path below the root of the crate that is
// named crate, applied to args. The crate is separate from the path: the form
// carries the crate's name first and then the segments of the path. This
// model checks the shape of the form only. The Rust target admits a form
// whose crate is a dependency that a destination can name. That is Std today,
// and the target refuses every other crate.
func External(crate string, path []string, args ...g.Expression[s.Language]) g.Expression[s.Language] {
	return g.Native(Form(FormExternal), append([]string{crate}, path...), args...)
}
	...
	case FormExternal:
		if len(d) < 2 || slices.Contains(d, "") {
			return g.ErrInvalid
		}
```

The coverage answer is in the Rust target
(`transform/syntax/rust/target/go/provider.go:23-49`):

```go
// CheckExpression admits the forms of the Rust model. An external form is
// admitted when its crate is a dependency that a destination can name.
func (p Provider) CheckExpression(ctx context.Context, e g.Expression[s.Language]) error {
	if err := model.CheckExpression(ctx, e); err != nil {
		return err
	}
	if name, _ := model.Name(e.Form()); name == model.FormExternal {
		_, err := external(e.Data())
		return err
	}
	return nil
}

// external returns the destination that the data of an external form names:
// the crate first, then the path of the item. ...
func external(data []string) (symbol.Symbol, error) {
	crate := symbol.Std()
	if data[0] != model.Std {
		var err error
		if crate, err = symbol.Dependency(data[0]); err != nil {
			return symbol.Symbol{}, err
		}
	}
	return symbol.Item(crate, data[1:]...)
}
```

`symbol.Dependency` returns `ErrUnsupported` for every crate but the
standard library. So for the external form Rust has the split that Go lacks:
the model says what the form is, and the target says what is covered.
Section 6.5 builds on this.

### 3.5 A symbol package, and the Files stage that reads one

A **symbol** is what a name use in generated output refers to when the
referent is outside the generated declarations: a predeclared type, an item
of a standard library, an export of a package. The generation contract
supplies the issuer (`contracts/generation/go/output.go:17-43`):

```go
// SymbolIssuer is a target library's authority to issue and interpret immutable
// external destinations. Its key is opaque to the generation and Files kernels.
type SymbolIssuer struct{ owner *identity }
type ExternalSymbol struct {
	owner    *identity
	language representation.Language
	key      string
}

func NewSymbolIssuer() SymbolIssuer { return SymbolIssuer{fresh()} }
func (i SymbolIssuer) Issue(language representation.Language, key string) (ExternalSymbol, error)
func (i SymbolIssuer) Read(language representation.Language, s ExternalSymbol) (string, bool)

type ExternalUse struct {
	At     semantics.Occurrence
	Symbol ExternalSymbol
}
```

`representation.Language` is the descriptor of a language from the envelope
package: its identifier, display name, profile string and marker type.

Go's symbol package (`transform/syntax/golang/symbol/go/symbol.go:1-81`,
comments abridged):

```go
// Package symbol owns the Go external destinations shared by generation and Files.
package symbol

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

func (y Symbol) Predeclared() (string, bool) { return y.name, y.category == predeclared }

var issuer = g.NewSymbolIssuer()

// role is the role of an external use on its node: the role that a reference
// to a declaration of the output has on such a node.
var role = semantics.Key{Owner: string(s.Language{}.ID()), Name: s.RoleTarget}

// ExternalUse records symbol as the destination of the name use at site.
func ExternalUse(language representation.Language, site ast.NodeRef, symbol Symbol) (g.ExternalUse, error) {
	if representation.CheckMarker[s.Language](language) != nil || language.ID() != (s.Language{}).ID() {
		return g.ExternalUse{}, g.ErrScope
	}
	name, ok := symbol.Predeclared()
	if !ok {
		return g.ExternalUse{}, g.ErrInvalid
	}
	issued, err := issuer.Issue(language, predeclared+"\x00"+name)
	if err != nil {
		return g.ExternalUse{}, err
	}
	return g.ExternalUse{At: semantics.Occurrence{Site: site, Role: role}, Symbol: issued}, nil
}

// ReadExternal returns the destination that an external use of this package
// records. It refuses a use of another issuer, of another language or at
// another role with generation.ErrScope.
func ReadExternal(language representation.Language, use g.ExternalUse) (Symbol, error) {
	key, ok := issuer.Read(language, use.Symbol)
	if !ok || use.At.Role != role || use.At.Ordinal != 0 {
		return Symbol{}, g.ErrScope
	}
	return Predeclared(strings.TrimPrefix(key, predeclared+"\x00"))
}
```

The package states which names are symbols by asking the Go standard library
(`types.Universe`). It holds no table. Its second refusal is a coverage
answer that depends on what the Go bridge can tell apart.

**Where a Files stage uses a symbol package.** Only TypeScript has a Files
stage for generated output. It reads each external use of the sealed output
back through TypeScript's symbol package
(`format/files/typescript/provider/go/generated.go:163-180`; `tssymbol` is
the TypeScript symbol package, whose symbol type is
`Type{Module, Name string}`; `p.bindings` maps the node of a name use to
what the use is meant to reach):

```go
for _, use := range output.ExternalUses() {
	site := use.At.Site
	ext, err := tssymbol.ReadExternal(snap.Language(), use)
	if err != nil {
		// A use that the TypeScript issuer cannot read is invalid input. The
		// generation value stays joined to the Files category.
		return nil, fmt.Errorf("%w: external use at %v: %w", source.ErrInvalid, site, err)
	}
	node := p.trees[site.Tree].Nodes[site.Node-1]
	if node.Kind != s.KindTypeReference || len(node.References) != 0 || imports.Attr(node, s.AttributeName) != ext.Name {
		return nil, fmt.Errorf("%w: external destination at %v", source.ErrInvalid, site)
	}
	if ext.Module == "" {
		p.bindings[site] = symbol{builtin: ext.Name}
	} else {
		p.bindings[site] = symbol{external: imports.Export{Module: ext.Module, Name: ext.Name}}
	}
}
```

After it has placed the declarations and had the placed program analysed,
the same stage reads the analysed graph in the words of the semantic
vocabulary and builds the same kind of value from what the compiler bound
(`format/files/typescript/provider/go/placement.go:196-222`; `n` is
TypeScript's semantic vocabulary):

```go
for _, b := range data.Bindings {
	...
	if b.At.Role != n.NameRole() && b.At.Role != n.ImportRole() {
		continue
	}
	entity := data.Entities[b.Target-1]
	var dest symbol
	switch {
	case len(entity.Declarations) == 1:
		dest.site = entity.Declarations[0]
	case entity.Kind == n.PackageTypeKind() && entity.Owner != 0:
		owner := data.Entities[entity.Owner-1]
		if owner.Kind != n.PackageKind() {
			return source.ErrInvalid
		}
		dest.external = imports.Export{Module: owner.Name, Name: entity.Name}
	case entity.Kind == n.BuiltinKind():
		dest.builtin = entity.Name
	default:
		return fmt.Errorf("%w: ambiguous native binding %s", source.ErrInvalid, entity.Name)
	}
	...
	p.bindings[b.At.Site] = dest
}
```

So this one package is where an intended vocabulary and an analysed
vocabulary of the same language meet today. The intended destination comes
as a symbol (`Type{Module, Name}`). The established binding comes as an
entity with a kind (`package-type`, `builtin`) and an owner. The Files stage
translates both into a private type of its own and compares. The accepted
specification of destinations requires that comparison to be one of
identities: "A check that a use reaches its destination therefore compares
entities and not types" (`docs/destinations.md`). The same page decides that
nothing compares the role under which a symbol is recorded with a role of
the bridge, and the boundaries page decides that a symbol package does not
import the semantic vocabulary.

### 3.6 Go's three vocabularies side by side

The table lists the type-level notions of Go and what each of three
vocabularies calls them. The facts are verified against the code. The last
column is our derivation.

| Notion | Syntax node kind (`syntax/golang/go`, as the target writes it) | Analysed (`semantics/golang/go`, as the bridge writes it) | Intended form (`model`) | Our reading of the difference |
| --- | --- | --- | --- | --- |
| Predeclared basic type (`string`, `bool`) | `named-type` with the attribute `name`; the destination is recorded as a symbol `predeclared` | Entity kind `builtin`, named by the identifier, without owner and site; term operator `basic` with the compiler's basic kind as data | `primitive` with one datum, `string` or `bool` | Four words for one notion. Intent names the type it means. Evidence names the entity the compiler found, which a declaration named `string` can capture |
| Slice | `slice-type` with one `element` child | Term operator `slice` with one operand | `slice` with one argument | None in the notion. The same word three times |
| Struct with fields | `struct-type` with `field` children; also `embedded-field` | Term operator `struct` with JSON field records (name, tag, embedded, entity) and one operand per field; entity kind `field` | `record` whose arguments are `field` forms; a field has the data name and optional flag | Three words. `record` is the word the four models share. The optional flag exists only in intent, and the Go target does not read it. Tag and embedding exist only in syntax and evidence: nothing generates them |
| Defined type declaration | `type-declaration` | Entity kind `defined`; fact relation `underlying` | Declaration form `defined`; the body must be a record | The same word. The restriction to a record body is a limit of today's generation |
| Alias declaration | `type-alias` | Entity kind `alias`; fact relation `alias-target` | Declaration form `alias` | None |
| Type of another package | `qualified-type` and `import`; the target writes neither | Entity kind `package` with a `package-name` fact; a `defined` or `alias` entity owned by it | `external` with the data module and name; refused as unsupported | Intent has a name and no meaning yet |
| Type parameter | `type-parameter` with a constraint child | Entity kind `parameter`; fact relations `constraint` and `empty` | No form. The kernel's own expression kind and `NativeParameter` | The binder's identity is the kernel's on the intent side. The constraint exists only in syntax and evidence |
| Use of a declaration, generic application | `named-type` with a `target` reference; `type-application` | The graph's own reference and application terms, and a binding with the role `name` | No form. The kernel's own expression kind | Owned by the shared carriers on both sides |
| Package | The attribute `package` on `file` | Entity kind `package` | Not in the model. A setting of the target | |
| Only in evidence | | Fact relations `identical` and `display`; operators `interface`, `union`, `signature`, `tuple`, `boolean` | | Type identity exists only for a checked program. The other operators are Go notions that nothing generates |
| Only in syntax | `pointer-type`, `array-type`, `map-type`, `channel-type`, `function-type`, `method`, `embedded-type`, `type-union`, `type-approximation` | (some as unpublished operators, above) | | Notions that syntax can spell and nothing analyses by name or generates |

Summary, as facts: of the seven forms of Go's model, five name a notion that
analysis also names, three of them under another word
(`primitive`, `record`, `external`), and the two declaration forms under the
same word (`alias`, `defined`). Every form has a node kind that expresses
the same notion.

Two structural differences are visible, and both are our reading. First, an
intended form is a tree of values with placeholders, without identity, site
or spelling. A node has a site in one tree. An analysed type is a term in a
graph with entities that have identity. Second, intent refers to a
destination before any program exists, and evidence refers to the entity
that a compiler bound in one placed program.

### 3.7 How a role is computed from a path

The architecture test derives a package's import rule from its path
(`integration/go/architecture_roles_test.go:157-183`, `:336-346`; `role`
builds a rule from a shape, a kind, the language read from the path, and the
permitted import patterns, in which `<language>` is replaced):

```go
// packageRole derives a package's rule from its layout path. The language, the
// stage and the pair come from the path, so a new package of a known shape has
// its rule without an edit here.
func packageRole(relative string) (roleRule, bool) {
	s := strings.Split(relative, "/")
	switch {
	case len(s) == 2 && s[0] == "internal":
		return helperRole(s[1])
	case s[len(s)-1] != "go":
		return roleRule{}, false
	case len(s) == 3 && s[0] == "syntax" && isLanguageName(s[1]):
		// The syntax kinds and identity of one language.
		return role("syntax/<language>/go", roleVocabulary, []string{"language", s[1]},
			"contracts/language/go", "trees/ast/go", "trees/cst/go",
		), true
	case len(s) == 3 && s[0] == "semantics" && isLanguageName(s[1]):
		// The semantic names and contract identities of one language. It may use
		// its language's lexical helper, never internal/<language>sem, which
		// runs the compiler behind the names.
		rule := role("semantics/<language>/go", roleVocabulary, []string{"language", s[1]},
			"contracts/analysis/go", "contracts/resolution/go", "contracts/semantic/go", "semantics/go",
			"syntax/<language>/go", "trees/ast/go", "trees/representation/go", "internal/<language>",
		)
		rule.forbidden = func(dependency string) string {
			if isSemanticHelper(dependency) {
				return "a vocabulary never imports a semantic helper"
			}
			return ""
		}
		return rule, true
	case len(s) == 5 && s[3] == "provider" && isLanguageName(s[2]):
		return stageProviderRole(s[0], s[1], s[2])
	...
	case s[0] == "transform":
		return transformRole(relative)
	...
	}
	return roleRule{}, false
}

// [doc] inside transformRole, for transform/syntax/<language>/<component>/go:
	switch strings.Join(component, "/") {
	case "model":
		// The intended native forms. No tree is allocated here and no source
		// language is read.
		return role("transform/syntax/<language>/model/go", roleVocabulary, language, "contracts/generation/go", "syntax/<language>/go"), true
	case "symbol":
		// The public type symbols that generation and Files formatting of one
		// language agree on.
		return role("transform/syntax/<language>/symbol/go", roleVocabulary, language,
			"contracts/generation/go", "internal/<language>", "semantics/go", "syntax/<language>/go", "trees/ast/go", "trees/representation/go",
		), true
```

The three vocabulary roles side by side, as the test states them:

| Role by path | May import |
| --- | --- |
| `syntax/<language>/go` | `contracts/language/go`, `trees/ast/go`, `trees/cst/go` |
| `semantics/<language>/go` | `contracts/analysis/go`, `contracts/resolution/go`, `contracts/semantic/go`, `semantics/go`, `syntax/<language>/go`, `trees/ast/go`, `trees/representation/go`, `internal/<language>`. Never a bridge |
| `transform/syntax/<language>/model/go` | `contracts/generation/go`, `syntax/<language>/go` |
| `transform/syntax/<language>/symbol/go` | `contracts/generation/go`, `internal/<language>`, `semantics/go`, `syntax/<language>/go`, `trees/ast/go`, `trees/representation/go` |

Three facts. Every vocabulary role lists at least one contract package. No
vocabulary of a language may import another vocabulary of that language
except the syntax vocabulary. And the role of a Files stage does not list
the generation contract or a symbol package: TypeScript's Files stage has
both through an entry in the table of role extensions.

### 3.8 The node patterns of horos kinds

horos has a second level above types. "A type describes values. A kind
describes types or type constructors." (`docs/kinds.md`). A kind is written
in a notation for the shape of a tree. bittype parses, checks and analyses
this notation as part of horos. The page `docs/horos-syntax.md` gives it:

```text
kind Text = node("scalar", {"selector": "utf8-text"}, epsilon);
kind List<K> = node("list-of", {}, item(K));
kind ClosedRecord<K> = node("record", {"mode": "closed"}, star(item(K)));
kind Pair<A, B> = node("tuple", {}, concat(item(A), item(B)));
```

```text
Q = epsilon
  | item(K)
  | concat(Q, Q, ...)
  | alternative(Q, Q, ...)
  | star(Q)
```

A `node` pattern names a constructor, constrains its metadata through
selectors, and gives a pattern `Q` for its children. The page says: "A node
pattern matches the entire ordered list of observed type children." Kinds
also have union `|`, intersection `&`, parameters, and recursion with
greatest fixed points:

```text
kind Tree = nu {
    X = h.Scalar | h.List<X> | h.Record<X>;
} in X;
```

**What the notation describes today.** It describes horos type descriptors
and nothing else. "Kinds classify **checked structural type descriptions**."
(`docs/kinds.md`). The first argument of `node` is one of a closed set of 22
descriptor constructors of horos. The set is a literal list in the horos
admission table (`internal/declsyntax/check.go:91`): any, atom, atom-exact,
atom-len, atom-in, text-len, key-grammar, text-non-blank, lowercase-uuid,
scalar, exact, wire, seq, list-of, set-of, not, map-of, tuple, compound,
union, and, record. Selectors exist for two of them: `selector` on `scalar`
and `mode` on `record`. The page puts outside this revision: "Public kind
negation, arbitrary predicate callbacks and subtree-equality constraints".
It also says that "a uniform field-kind query does not introduce row
variables or predicates over labels".

**Compared with the admission table of section 3.2** (our derivation, from
reading the two notations side by side):

| | Admission table (`Rules`) | horos node pattern |
| --- | --- | --- |
| Node identity | Any node kind of the language | One of 22 descriptor constructors of horos |
| Attributes | Named, required or optional, each with a predicate function | Selectors with string values, for two constructors |
| Children | Named roles in a fixed order, each with a count range and a flat list of kinds | An unnamed ordered sequence with concatenation, alternative and repetition |
| Recursion | By naming kinds in the lists | Explicit fixed points |
| Parameters | None | Kind parameters |
| References | A flag: one role, `target` | None |
| Interpreter | One Go function of 67 lines | The horos library, which owns the meaning of kinds. An earlier consultation of this project (0005) recorded the choice to "implement proper-kind inclusion using finite cuts of positive greatest-fixed-point hedge automata" |

No code and no page on the main branch applies a horos kind to a syntax
tree. No language's tree is described in this notation, and that includes
the tree of horos itself: its own shape is a Go table and 11 rules in code,
like the others.

### 3.9 The one generator, and the absence of other generated code

On the main branch no Go file is generated. No file holds a `go:generate`
directive or a "Code generated" header. One table is produced by a manual
procedure: the character ranges that may form a Rust crate name, 840 ranges
printed by a program built with the pinned Rust compiler and pasted into
`internal/rust/cratename.go`. A test compares the table with the compiler
where the compiler is available.
Three documentation pages are generated from the source by a tool.

On the migration's branch one generator exists, `tools/payload-schema`. Its
input is a horos module that declares the type of a mapping's decision
record. An example, complete
(`transform/syntax/horos-to-golang/vocabulary/declarations/go/payload.horos`
at the branch head c1e5069):

```horos
module bittype.payload.horos_to_golang_declarations;

type Payload = { name: text; nominal: bool; parameters: list<text>; };
```

From it the tool makes three files beside it: the canonical binary module
`payload.ontos`, a Go file with the module's pin, and a Go file with the
type (`payload_types.go`, complete):

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
contributions". A schema may use what those contributions cover: text,
booleans, lists, closed records with required and optional fields, and
named declarations. The branch holds 24 such schemas. A test regenerates
each file and fails when the committed file differs.

So the project already declares one kind of model in horos and generates Go
from it: the decision record of a mapping. It uses horos types for that, not
kinds, and it describes a value, not a tree.

## 4. Data and control flow: one record from horos to Go

Standing: implemented today for steps 1 to 8, read from the code and
followed by hand, not run. Steps 9 to 11 are planned for Go and are
described from the accepted specification of destinations
(`docs/destinations.md`) and from what the Go bridge writes today. For
TypeScript the same three steps are implemented (section 3.5).

The input is the module of section 1.1:

```horos
module example.people;
type Person = { name: text; };
```

A consumer's composition has registered the horos source reader, the Go
target, and the Go mappings for text, records and declarations.

| Step | What happens | Who does it (role, place) | Vocabulary touched |
| --- | --- | --- | --- |
| 1 | The text is parsed into a concrete tree that owns an abstract tree | The horos parser (contribution, `parse/syntax/horos`) | **Writes the horos syntax vocabulary**: nodes of the kinds `program`, `module`, `type-definition`, `descriptor` |
| 2 | The tree is admitted: the horos table, then the horos rules in code | The horos checker (contribution, `check/syntax/horos`, with a private helper) | Reads the horos syntax vocabulary and the horos grammar table |
| 3 | Names are resolved and the declarations are analysed. The result is a semantic graph and two receipts | The horos resolver and analyzer (contributions, `resolve/files/horos`, `analyze/files/horos`), over a private adapter to the horos library | **Writes the horos semantic vocabulary**: an entity of the kind `type-definition`, a provider expression with the operator `descriptor/record`, facts with the relations `sort` and `canonical-identity` |
| 4 | The planner walks the graph from the requested declaration `Person`. For each declaration and each use of a type it asks the source reader what it is | The planner (kernel, inside `contracts/generation/go`) and the horos source reader (contribution, `transform/syntax/horos/source`) | Reads the generic graph. The source reader reads the horos semantic vocabulary, and a few names of the horos syntax vocabulary for record payloads |
| 5 | For each use one mapping is selected and states intent. The declaration mapping gives the name `Person` and the form `defined`, because the body is a record. The record mapping reads one field, asks the planner for the representation of its type, and returns a record form. The text mapping answers that request with a constant | Three mappings (contributions, `transform/syntax/horos-to-golang/mapping/{declarations,record,text}`) and the planner | **Writes Go's model**: the forms `defined`, `record`, `field` and `primitive`, and the primitive name `string`, through the model's constructors |
| 6 | The planner has every expression and every declaration checked as it composes them | The Go target (contribution, `transform/syntax/golang/target`), which forwards to the two checks in Go's model | Reads Go's model |
| 7 | The target builds one tree from the closed program. It switches over the form names. `primitive` becomes a `named-type` node with the name `string`, and the use records its destination | The Go target and the shared tree builder (kernel, `transform/syntax/shared/target`) | Reads Go's model. **Writes Go's syntax vocabulary**: `file`, `type-declaration`, `struct-type`, `field`, `named-type`. **Writes Go's symbol**: the predeclared symbol `string` |
| 8 | The tree is admitted and the output is sealed with its inventory and its external uses | The Go checker (contribution, `check/syntax/golang`), called by the planner's sealing code | Reads Go's syntax vocabulary and Go's grammar table |
| 9 | The declarations are placed into files of a package under the consumer's placement. Planned for Go | A Go Files stage (contribution; under the present path form `format/files/golang`) | Reads Go's syntax vocabulary. Reads Go's symbol package to learn the destination of each external use |
| 10 | The placed files are analysed by the Go type checker. Planned as part of step 9 | The Go resolver and analyzer over the Go bridge (private part, `internal/golangsem`) | **Writes Go's semantic vocabulary**: an entity `Person` of the kind `defined`, a fact `underlying` whose term has the operator `struct`, a field type whose term has the operator `basic`, and a binding with the role `name` from the node of the use `string` to an entity of the kind `builtin` |
| 11 | Each name use is compared with its destination: the entity that the compiler bound against the entity that the symbol identifies. Planned for Go | The Go Files stage | Reads Go's semantic vocabulary (the binding, the kind `builtin`) and Go's symbol (the category `predeclared`) |

The intent after step 5, written with the model's constructors:

```go
// [doc] the declaration: Name "Person", Form "defined", and as its body
model.Record(
	model.Field("Name", false, model.Primitive(model.PrimitiveString)),
)
```

The tree after step 7, as node kinds with their attributes and roles:

```text
file (package = <the target's configured package name>)
  declaration: type-declaration (name = Person)
    type: struct-type
      field: field (name = Name)
        type: named-type (name = string)     + external use: symbol predeclared "string"
```

The graph that the bridge would publish in step 10, in the words of the
semantic vocabulary. This is an illustration written by hand from the
bridge's code of section 3.3:

```text
entity  kind go/defined   name "Person"   owner: the package entity   site: the type-declaration node
  fact  go/underlying -> term  operator go/struct
                               data  [{"Name":"Name","Tag":"","Embedded":false,"Entity":<field entity>}]
                               operand 0 -> term operator go/basic (the string kind)
entity  kind go/builtin   name "string"   no owner, no site
binding at the named-type node, role go/name -> the builtin entity
```

**Where each of Go's vocabularies is written and read.**

| Go vocabulary | Written in step | By which stage | Read in step | By which stage |
| --- | --- | --- | --- | --- |
| Model (`transform/syntax/golang/model/go`) | 5 | transform (mappings) | 6, 7 | transform (target) |
| Syntax (`syntax/golang/go`) | 7 | transform (target) | 8, 9, 10 | check; format; resolve and analyze |
| Symbol (`transform/syntax/golang/symbol/go`) | 7 | transform (target) | 9, 11 | format in the Files dimension (planned) |
| Semantic (`semantics/golang/go`) | 10 | resolve and analyze (bridge) | 11 | format in the Files dimension (planned) |

Four observations, as facts.

- The same Go type, a struct with one string field, is stated three times
  along the way, each time in other words: as forms (`defined`, `record`,
  `field`, `primitive`), as nodes (`type-declaration`, `struct-type`,
  `field`, `named-type`) and as analysed facts (`defined`, `struct`, `basic`,
  `builtin`).
- The model is written and read inside one stage. The symbol is written in
  one stage and read in another. The syntax vocabulary and the semantic
  vocabulary are used by several stages.
- The source side has no model. Generation reads horos through its analysed
  vocabulary. A source is always an analysed program. A target does not exist
  yet when it is described. No page states this asymmetry as a decision.
- Step 11 is the one place where intent and analysed meaning of one language
  are compared. The comparison is of identities (a symbol against an
  entity), not of type structure. No production code today compares an
  intended `record` form with an analysed `struct` term.

## 5. Constraints

This section lists only owner rules and accepted designs. Section 5.2 then
says what is not decided, and section 5.3 lists eight places where pages
contradict each other or the code. No fact of the present code is a
constraint, with one exception that is a fact of the implementation language
(fixed point 21).

### 5.1 Fixed points and their standing

"Research N" is an earlier expert consultation of this project. A row of its
result table has a verdict that the research owner of that consultation
gave: HOLDS means accepted as advised.

**About the answer's method.**

| # | Statement | Standing | Shown in |
| --- | --- | --- | --- |
| 1 | "The currently implemented surface must have **zero influence** on that choice. Implementation or migration effort has **zero weight**, including as a tiebreaker." Estimates of the existing implementation "can guide delivery sequencing only" | Owner rule (AC-09, 8 November 2026) | 1.5, 2.1 |
| 2 | "Move a misplaced or wrongly named responsibility to its owner first, as its own change, without waiting for the complete implementation of its destination. That first change establishes the final structure". Afterwards coverage "is then added without moving, renaming or removing what exists" | Owner rule (AC-10, 11 November 2026) | 2.1 |
| 3 | Wrong grouping is corrected before wrong coupling: "Correct the grouping first, as AC-10 orders, and the coupling becomes visible." A dependency check "enforces a boundary that reasoning has already decided" and "does not replace that reasoning" | Owner request of 11 November 2026, with the platform's reading | 2.1 |
| 4 | "A language-agnostic mechanism belongs in the agnostic layer. Each language supplies its own types, meanings and policies through the shared contract." | Owner rule (AC-11, 13 November 2026) | 2.1 |
| 5 | "An extension needs a concrete case that the existing contracts cannot express or uphold, together with an explanation of why composition does not resolve it." | Owner rule (7 October 2026), in the platform's words | 2.1 |
| 6 | Unification "needs stronger evidence than overlapping fields or similar code." "Compare identity, lifetime, authority, allowed operations, guarantees and failure behavior." "Differences must be represented explicitly rather than erased." | Accepted principle of the platform | 2.1 |
| 7 | "Semantic ownership is about who defines a meaning and how that meaning enters an operation. It is not determined by a directory, package name, generic type parameter or the word "plugin"." | Owner rule (6 November 2026), in the platform's words | 2.1 |

**About intent.**

| # | Statement | Standing | Shown in |
| --- | --- | --- | --- |
| 8 | "Native intent is not an analyzed `Semantics[T]`: target analysis establishes facts about the actual constructed/organized program." | Accepted design (`docs/generation.md`, 8 November 2026) | 2.7 |
| 9 | "The planned native program is intent, not a fabricated semantic-analysis result. Only analysis of the realized target program establishes native semantic facts." | Accepted design (`docs/generation-planner.md`) | |
| 10 | "Native semantic evidence belongs to the actual native capture. A generation plan or successful syntax check is not that evidence." | Accepted design (`docs/type-system-architecture.md`) | |
| 11 | The unplaced output of generation "need not contain compiler-established facts about a hypothetical native program" | Accepted design (research 0016 on the placement of generated declarations, the expert's words; the corresponding row 19 is adopted with an adaptation) | |
| 12 | "Native expressions contain target-owned forms, immutable scalar metadata and child expressions. The shared interpreter understands substitution and references, not the meanings of these forms. There is no AST kind, source string or arbitrary payload field in an expression." | Accepted design (`docs/generation-spi.md`) | 3.4 |
| 13 | "The native target remains the trusted interpreter of its own forms." The four models share no package: "Each model owns its own forms, a form can be read only by its issuer, and no package compares a form name of one language with another's or with a literal." | Accepted design (case record ACX-005; change log entry of 12 November 2026) | 3.4 |
| 14 | "bittype's native vocabularies stay types-only." | Accepted design (`docs/binding-generation.md`, research 0014 on binding generation) | |

**About a shared model of types.**

| # | Statement | Standing | Shown in |
| --- | --- | --- | --- |
| 15 | The semantic graph "has no scope graph, substitution map, universal type algebra or conditional proof system". "A provider owns lookup, semantic rules, vocabulary, coverage and the interpretation of its payloads. The kernel validates graph integrity; it cannot prove a provider's claims or make native meanings equivalent to Horos." | Accepted design (`docs/semantics.md`; research 0002 on the semantic kernel, row 1, HOLDS: "Share established facts with provider-owned vocabulary, separate entity bodies and parameter ownership") | 2.2, 3.3 |
| 16 | The mechanism by which mappings state target types "should not become a new language-neutral type calculus". Row 21, HOLDS: "Keep the SPI experimental and avoid a richer generic semantic model." | Accepted design (research 0007 on generator plugins) | |
| 17 | A model of a type system in the sense of terms, judgments and rules is content of logos packages: "These are requirements for language integrations and Logos semantic packages, not fields of a new shared Bittype model or implemented Go struct." Row 1, HOLDS: "Consume Logos semantic packages; no Bittype-owned type-system foundation" | Accepted design (`docs/type-system-model.md`; research 0010 on the architecture of type-system models, which was about verification) | |
| 18 | A destination is an identity. "A destination is a declaration's identity and never structural equality". "An external symbol names the selected public export and not the canonical type underneath it". "Three statements define it. A later change may not alter them." | Accepted design (`docs/destinations.md`, specification of 13 November 2026, from research 0016) | 3.5 |

Fixed points 15 to 17 each rule out one thing: a universal type algebra in
the graph, a language-neutral calculus in generation, and a type-system
framework owned by bittype. None of them speaks about one language's own
vocabulary of its own types. That is our reading, and question 2 rests on
it.

**About vocabularies and the layout.**

| # | Statement | Standing | Shown in |
| --- | --- | --- | --- |
| 19 | A vocabulary is "The names two parties must agree on and nothing else: syntax kinds, semantic keys, profile and contract identities, intended native forms. A vocabulary package imports no contribution." | Accepted design (`docs/package-boundaries.md`, state of 12 November 2026). Section 5.3 names what contradicts "nothing else" | 2.5 |
| 20 | "**One language.** A package that names a language imports no package of another language. A mapping names two." | Accepted design (`LAYOUT.md`) | 2.5 |
| 21 | Go's rule for a directory named `internal`: a package below it can be imported only by code in the tree rooted at the parent of that directory. Today `internal/` is at the module root, so every package can import the helpers below it. A helper below `ops/internal/` could not be imported from `model/`, and the reverse | A fact of the implementation language. Not a design of ours | 6.2 |
| 22 | "One public root Go module contains independently owned packages". "Do not create a language-first packages tree, additional Go modules, compatibility aliases or placeholder source directories." | Owner decision of 31 October 2026 as the layout page records it, and the layout page | 2.3 |
| 23 | "organize by processing stage and enforce a uniform structure across languages"; `<stage>/<dimension>/<scope>/<component>/go` | Owner decisions of 31 October and 6 November 2026, as the layout page records them. The owner who made them now proposes a change to them. We therefore treat them as the thing under examination and not as a limit on the answer | 2.3 |
| 24 | "A checker gets no concrete tree, a formatter no lexer, and no Syntax stage a semantic contract." | Accepted design (`LAYOUT.md`, role table) | 2.5 |
| 25 | The design of the contract migration lists, among the shared packages that must stay unchanged when a materially different contribution is added, "the native model and target of each language". It is final only when its gate is met | Accepted design, provisional (`docs/contribution-contract.md`, 13 November 2026) | 2.8 |

Out of scope: the design of the contract migration itself, the choice of a
native representation for any horos type, the Files stage's design, and
verification.

### 5.2 What is not decided

An investigation read the layout page, the charter, the contributor
instructions, the change log, the design pages and the result tables of the
earlier consultations for each of these questions. It found no accepted
sentence that answers any of them. It lists what it did not read, and the
statements below are bounded by that.

1. **Why the model is inside `transform/`.** The pages state the place and
   the ordering rule (stage first). No sentence argues for the stage
   directory against a place beside `syntax/` and `semantics/`. The commit
   that created the four packages has no message body, and no consultation
   was asked. The two dependency statements that the pages give ("Mappings
   depend on their source's reader and target models, never AST construction
   packages. Targets depend on their model and shared constructor, never
   source readers or mapping libraries.") say who may import the model. They
   do not say where it lives.
2. **Why the symbol package is inside `transform/`,** although the layout
   page itself describes it as "for generation and Files" and a package of
   another stage imports it.
3. **Whether the type-level notions of one language should have one
   vocabulary** in which generation states intent and analysis states
   evidence. Every decided sentence (fixed points 8 to 11) separates a
   planned program from an analysed, placed program. None speaks about the
   names of the notions. No consultation asked.
4. **What the top-level directories `syntax/` and `semantics/` are.** No
   sentence says whether they are stages, dimensions, levels or homes of
   vocabularies. Put together from several sentences, each holds the
   per-language vocabulary of one abstract representation, the AST and the
   semantic graph. That description is the investigator's.
5. **Whether a stage other than `transform` may import a model.** No role
   lists it except the target and the mapping. No sentence forbids or
   permits it for a Files stage, for verification or for a bridge.
6. **Whether verification laws and generation mappings should speak about Go
   types in one vocabulary.** Today the laws use the semantic vocabulary and
   the mappings use the model.
7. **Whether the source side is meant to have a counterpart of the model.**
   horos has a source reader and no model (section 4).
8. **Whether the checks in a model belong to the vocabulary or to the
   target.** The scan of package boundaries of 11 November 2026 did not
   examine the model packages. We decided this as a first step (section
   6.5); that decision is ours.
9. **What a model is for a language without analysis.** C++ has no semantic
   vocabulary and no symbol package.
10. **Whether the sameness of the four models is intended.** The change log
    decides only that no package is shared.

### 5.3 Eight contradictions that the answer must not rest on

Each is between two pages, or between a page and the code. All are to be
corrected whatever is decided. None is evidence for an option.

1. **"Names and nothing else" against what three of the four kinds of
   vocabulary hold.** The definition (fixed point 19) says "nothing else".
   The layout page describes the model as "native semantic forms and
   validation" and the Go semantic vocabulary as "Go vocabulary, contract
   identities and immutable native type readers". The boundaries page itself
   admits a method into a vocabulary with a criterion: "It reads no file and
   starts no process, so a vocabulary can hold it." In the code only the
   syntax vocabularies hold names and nothing else.
2. **The number of role extensions.** Two pages say ten packages have one.
   The table in the test has twelve.
3. **The target's role.** The layout page gives a target "its own syntax,
   model and lexical helper". Three of the four targets also import their
   symbol package, through a role extension.
4. **"Its own language's vocabulary".** The pages say a vocabulary imports
   "Foundations and its own language's vocabulary". The test permits each
   kind of vocabulary only the syntax vocabulary of its language. A model
   may not import the semantic vocabulary or the symbol package, and the
   reverse.
5. **Two stale entries under "not examined".** Pages still list the form
   names and the readers' key literals as not examined. Both are corrected
   in the code.
6. **Who supplies a model.** One page says: "Libraries supply declarations,
   native models and mapping descriptions under their own ownership." The
   code and two later documents treat the model as the target's closed set:
   the issuer is private, a form of another issuer is refused, an unknown
   form name is unsupported, and the target switches over the model's names.
7. **The model as fixed and shared, against families that need new forms.**
   Fixed point 25 keeps the model unchanged when a contribution is added. An
   earlier consultation recorded that for many families "Only the models and
   mapping libraries are missing". A family that needs a form the model
   lacks, such as a map or an integer, must edit the model and the target of
   that language. No page says whether that is an edit of a shared package
   in the sense of the gate.
8. **The optional flag of Go's field form has no reader.** The model writes
   it and validates it. The Go target ignores it. The record mapping wraps
   the type instead.

## 6. What we have considered

Standing of this section: section 6.1 opens with the owner's proposal.
Everything else is our derivation, or a fact of the code that is marked as
such. All of it is open. No option is ruled out because the present code
would have to change. The document does not choose among the options. The
labels P, U and D are unique, so that an answer can cite them.

### 6.1 The owner's proposal, and what it means for the code as it is

**The proposal** (message 7, repeated here so that the section is complete):

```text
model/...
ops/...
model never imports from ops
model:
model/{syntax,semantics,...}/go/...
ops:
ops/{transform,...}/go/...
```

It names two roots and one rule. It gives two members of the first root,
`syntax` and `semantics`, and one of the second, `transform`. It leaves open
what the ellipses hold. In particular it does not say where the generic
carriers, the contract packages, the private helpers, the tools and the
verification component go. The owner's reason is in message 6: "The model
is independent of the operation, isn't it?"

**A test for "model", applied to every package.** To see what the proposal
means concretely, an investigation applied one stated test to all 129
production packages. The test is ours, not the owner's. In our words:

> A package is a *model* if it declares what exists (carriers, vocabularies
> of names, forms, identities, value types with their well-formedness) and
> performs no processing of an input into an output. It is an *operation* if
> it parses, checks, formats, prints, resolves, analyzes, transforms,
> verifies or composes. It is *both* if it holds a declaration and
> processing. It is *undecided* if the test gives no answer.

Validation at construction time was counted as well-formedness. Anything
that takes one artifact and returns another was counted as processing: a
tree from text, an encoding from an environment, a verdict on a form.

| Result | Packages | Lines | Which |
| --- | --- | --- | --- |
| Model, or both with a model core | 22 | 5,428 | The three tree packages, `files/go`, `source/go`, `semantics/go`, five syntax vocabularies, four semantic vocabularies, four models, three symbol packages |
| Contract packages | 13 | 5,030 | Section 2.4. By the test `language` and `text` are models, `pipeline` is protocol only, and the other ten are both |
| Undecided | 10 | 1,498 | Below |
| Operation | 84 | 23,987 | The stage providers and query services, the shared printers, the generation kernel packages, the source reader, four targets, 21 mappings and three shared mapping parts, the verification packages, the tools, the engine, and nine private helpers (scanning, parser state, tree construction, the three bridges, the two horos adapters and the import helper of TypeScript's Files stage) |

**The packages that hold both.** Twelve of the 22 are not plain models under
the test, because each holds code whose input and output are different
artifacts:

| Package | The model part | The processing part |
| --- | --- | --- |
| `trees/representation/go` | The language descriptor and the opaque envelopes | Packing and unpacking: a deep copy and a structure check, which it calls in a `check/` directory |
| `source/go` | Units, groups, relations, the publication | A publisher of 195 lines that admits a publication under a policy with two callbacks into operations |
| `semantics/{golang,typescript,rust}/go` | Names, profiles, contract identities, environment types | The typed readers over a graph, and the encoding of an environment into settings and back |
| The four models | Form names and constructors | Two checks that return a verdict on a form (section 3.4) |
| The three symbol packages | A symbol type and its categories | Admission of a name (for Go by asking the standard library), and the encoding of a symbol into an issued key |

**The ten that the test does not decide, and why.**

| Package | Why undecided |
| --- | --- |
| `check/syntax/shared/ast/go`, `check/syntax/shared/cst/go` | By their directory they are a check stage. By their content they are the well-formedness of the generic tree, which the test assigns to the model. The envelope package imports both |
| `internal/tree` | A read helper and a deep copy over the generic tree. It declares nothing and is no stage. The envelope package imports it |
| `internal/grammar` | It declares what an admission table is (model-like) and holds the interpreter of 67 lines (processing) |
| `internal/declsyntax` | The horos admission table, a literal, and the horos rules in code |
| `internal/cpp`, `internal/rust`, `internal/typescript`, `internal/horos` | Spelling rules. A predicate over strings states which strings exist, which is model-like, and it is written as code. Three of them also decode string literals, which is processing. Both sides import them: stage providers, and the semantic vocabularies and symbol packages |
| `transform/syntax/horos/mapping/go` | Two profile identities and one function type. The layout calls it a common mapping part, not a vocabulary |

**Every import that "model never imports from ops" would forbid today.**
The 22 model packages have 33 direct imports of an operation package, a
contract package or a helper. Four of the 22 import only models. No model
package imports a stage provider, a query service, a target, a mapping, the
engine or a bridge. Whether one of the 33 is forbidden depends on where its
target is placed, and that is exactly what the proposal leaves open.

| Kind | From | To | Symbols used | Edges |
| --- | --- | --- | --- | --- |
| A vocabulary uses a contract's value types | The five syntax vocabularies | `contracts/language/go` | `ID`, `Provider`: the two methods that state a language's identifier and name | 5 |
| | The four semantic vocabularies | `contracts/semantic/go` | `Setting`, `EvidenceRequirement`, `Result`, `Snapshot`, `Requirement`, `Obligation`, `Contract`, the verdict `Holds`, default limits, three error values, and two functions: `AdmitEvidence`, `CaptureSettings` | 4 |
| | The four semantic vocabularies | `contracts/resolution/go`, `contracts/analysis/go` | `Contract`, `Requirement`, `Kind`: the types of the contract identities | 8 |
| | The four models | `contracts/generation/go` | `NewFormIssuer`, `Form`, `Expression`, `Native`, `NativeDeclaration`, and the errors `ErrInvalid`, `ErrScope`, `ErrUnsupported` | 4 |
| | The three symbol packages | `contracts/generation/go` | `NewSymbolIssuer`, `ExternalUse`, the same errors | 3 |
| A carrier uses a contract's value types | `source/go` | `contracts/semantic/go` | `Snapshot`, `Result`, `EvidenceRequirement`, `Limits`, and the functions `AdmitEvidence`, `ValidateEvidenceRequirements` | 1 |
| | `trees/representation/go` | `contracts/language/go` | `ID`, `Provider` | 1 |
| The carrier's own well-formedness lives in a check directory | `trees/representation/go` | `check/syntax/shared/ast/go`, `check/syntax/shared/cst/go` | `Structure`, `Check` | 2 |
| | `trees/representation/go` | `internal/tree` | `Clone` | 1 |
| Lexical rules needed on both sides | `semantics/rust/go`, the Rust symbol package | `internal/rust` | `CrateName`, `Identifier` | 2 |
| | `semantics/typescript/go`, the TypeScript symbol package | `internal/typescript` | `PackageSpecifier`, `TypeName` | 2 |
| | | | | **33** |

What the first kind would require if a vocabulary were to import no
contract package, as facts of the code:

- For the models and symbol packages: `Form`, `FormIssuer`, `Expression`
  with everything it embeds (slots, helper keys, the six expression kinds
  and their constructors), `NativeDeclaration`, `NativeParameter`,
  `IssuedValue`, the symbol issuer with `ExternalSymbol` and `ExternalUse`,
  and three error values. That is the carrier of intent of section 3.4. It
  is declared in the same package as the planner, and one of its kinds,
  the reference to a planned declaration, "is issued by the planner after it
  authenticates the source callee" (comment on `Expression`).
- For the semantic vocabularies and `source/go`: the receipt and the
  snapshot. Both are handles with private state. Only the contract's
  admission code creates them, which is what makes a receipt trustworthy. In
  Go a type with unexported fields is constructed in its own package.
- For the syntax vocabularies: a 25-line package of two declarations.

In the other direction the contract packages import the carriers. There is
a cycle of directories, though not of packages: `source/go` imports
`contracts/semantic/go`, and the parse, print and format contracts import
`source/go`.

### 6.2 Placement: three options

Each option is stated with what it implies. None is chosen. For each we say
what it does to the accepted sentence "A package's role follows from its
path".

**P1. Only the two transform-side vocabularies move out, beside `syntax/`
and `semantics/`.** The four models and the three symbol packages get
top-level homes of their own kind. The stage-first layout otherwise stays.

- What it answers: the owner's messages 3 to 5 for the two vocabularies that
  are inside a stage today. All four kinds of vocabulary are then in one
  kind of place.
- What it implies: seven packages move. No root rule is introduced. A model
  still imports the generation contract. The Files stage's import of a
  symbol package no longer reaches into another stage's directory.
- What it leaves open: what the top-level directories are (section 5.2,
  item 4); whether a carrier, a contract's value types or a grammar table is
  a model; and the vocabularies that the migration adds for a pair of
  languages.
- The role-from-path rule: unchanged in kind. Two path shapes are replaced
  by two others.

**P2. The owner's two roots, with the generic carriers and the model-like
types of the contracts on the model side.** `model/` holds the carriers, the
vocabularies and the value types that operations exchange. `ops/` holds the
stages, the kernel, the protocol part of each contract, the bridges and the
tools.

- What it answers: message 6 and message 7 in full. "model never imports
  from ops" becomes one rule on a path prefix.
- What it implies: every package gets a new import path. Ten contract
  packages are split into a value part and a protocol part, and each split
  is a grouping decision of its own. The carrier of intent leaves the
  planner's package. The handles that only admission code may create need a
  home that lets the operation create them and the model name them. The
  structure checks of the generic tree are on the model side, or the
  envelope stops checking. The lexical helpers must be importable from both
  roots, so they cannot be below `ops/internal/` (fixed point 21).
- What it leaves open: on which side a thing belongs that describes a
  relation between two languages and is used by one operation only: a
  mapping's vocabulary, a family's meaning, a profile identity. And whether
  a grammar table, which today is a contribution's private variable, is part
  of a language's model.
- The role-from-path rule: the root gives a first cut, model or operation.
  The roles within a root (vocabulary, carrier, contribution, kernel,
  private part) still need a path shape or a table each. So the proposal's
  one rule is added to the roles and does not replace them.

**P3. Two roots where the contracts stay with the operations, and a model
may import contract value types through a third place.** The vocabularies
and possibly the carriers form `model/`. Stages, kernel and contracts form
`ops/`. A third place holds what both need: the identity of a language, the
diagnostics, the carrier of intent, the receipt and evidence types, and the
helpers that both sides import.

- What it answers: message 6, with a narrower meaning of "model": what one
  language says about itself.
- What it implies: the rule becomes "a model never imports an operation;
  both may import the third place". The third place needs a definition of
  its own that says what may be in it. Without one it collects whatever is
  hard to classify. The split of contract packages is the same as in P2,
  with another destination.
- What it leaves open: whether the generic carriers are model or third
  place, and the same relation-between-languages question as P2.
- The role-from-path rule: three roots instead of two, and roles within
  them as in P2.

Three things are the same under every option, as facts. The syntax
vocabulary is used by every stage of its language, the semantic vocabulary
by several, the model by one and the symbol package by two (section 2.6).
The model's values are values of one operation's contract. And the migration
in flight adds per-pair and per-family vocabularies that no option above
places without a further decision.
