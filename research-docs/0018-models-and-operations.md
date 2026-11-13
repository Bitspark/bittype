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
