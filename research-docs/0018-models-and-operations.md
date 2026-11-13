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
