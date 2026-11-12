# Research: Placing generated native declarations into files, packages and modules: the placement input, the unplaced output and the placement check

**ID:** 0016
**Date:** 12 November 2026
**Status:** draft
**Owner:** claude-contract-chain-issues
**Issue:** https://github.com/Bitspark/bittype/issues/161 (with #159)
**Code baseline:** bittype d35f1fad4f84bc754900fadefcba00c70858dc0d

## Question

**horos** is a declarative type language. **bittype** is a Go library that reads
horos declarations and generates native type declarations for TypeScript, Go,
Rust and C++.

bittype separates two stages:

- **Generation** builds the declarations: their bodies, their type expressions
  and the intended destination of every type reference. It must decide nothing
  about files, packages, modules, namespaces, imports or module-level
  visibility.
- A later stage, called **Files**, takes the generated declarations, one
  **placement** supplied by the consumer (which declaration goes into which
  file and which native scope) and the consumer's dependency environment. It
  writes the package clauses, imports, aliases and visibility modifiers the
  placement needs. It then has the target language's own compiler resolve the
  placed program, and compares what the compiler found with what generation
  intended.

This split is implemented for TypeScript only. Two languages still make a file
decision during generation:

- **Go** writes the package clause, taken from configuration.
- **Rust** writes `pub` on every generated struct and type alias. It also
  writes `pub` on every field; whether that is a misplaced decision at all is
  discussed in section 6.8, and our present answer is that it is not.

Two small changes are about to move these two decisions into new Go and Rust
Files stages. They are **slice 1** of two work items (a work item is a GitHub
issue): the **Go work item**, issue #161, "Add Files placement for generated
Go", and the **Rust work item**, issue #159, "Rust generation decides module
visibility before Files placement". Each slice admits only the simplest
placement (one Go package without imports, one Rust module) and refuses
everything else. Later slices add external types, imports, several packages and
several modules. C++ follows when it has compiler-backed semantic tooling.

**What we need advice on.** The project's rule for such a first move is that it
establishes the final structure: its inputs and identities already have their
intended shape, and later slices only add (section 2.1, rule AC-10). The
project's owner, one person who decides for every repository involved, directed
on 9 November 2026: "I want this properly done, so that we don't accumulate
technical debt and have to refactor later." Nothing in bittype is released yet,
so every interface can still change now, and none of it should have to change
later.

**What shape should the placement input, the unplaced generated output and the
placement check have, so that Go, Rust and later C++ are added without changing
what exists once the two slices have landed?**

TypeScript's present types may be reshaped now, as part of this decision. The
test is that nothing changes after the two slices land. We count as a later
change, and therefore as a failure of the first move:

- a change to the type a consumer writes to state a placement, to the key
  under which it travels, or to the meaning of a placement that was already
  accepted;
- a change to what generated output is before placement (its root, its names,
  what it records about intent);
- a change to the evidence a Files stage requires, or to the constructor a
  consumer calls to build one;
- a change to a name that appears in consumers' code.

Admitting a placement that is explicitly refused today is additive and not a
change.

**How to read the standing of our statements.** Section 5 lists every fixed
point with one of four statuses. Owner rules are fixed. The accepted decision
of 7 November 2026 is fixed unless you see that it contradicts an owner rule or
itself; if you do, please say so. Everything that is merely implemented today,
and everything we derived, is open, however firmly the code or our text states
it. In your answers, please name the option first and then the reason.

**The tensions**, in the order of the questions:

- **Question 1: scope differs per language.** A TypeScript module's scope is
  its file. A Go package is a directory of files and has a name that the path
  does not determine. A Rust file can hold several modules. A C++ header can
  hold several namespaces, and a namespace is independent of the path.
- **Question 2: names and scopes.** Today all generated names must be unique in
  one flat scope, and one naming policy, as we read it, exists to make that
  possible. Once the consumer decides scopes, it is unclear who owns
  collisions.
- **Question 3: unplaced Go is not a Go file.** Generated output is a
  syntax-checked tree of the target language (grammar only; no names
  resolved). Go's source file production requires a package clause.
- **Question 4: a check that cannot see the moved decision.** For the Go
  package name and for Rust item visibility, the compiler-established facts
  bittype publishes contain nothing to compare against. For Go's predeclared
  `string` and `bool` the compiler publishes a fact, but generation records no
  intended destination.

**Map.** Sections 1 and 2 give the system, the rules and the two governing
documents. Section 3 is a one-page map: who decides what today, per language,
and where section 4 shows it. Section 4 is the code. Section 5 lists the fixed
points and their standing. Section 6 gives the options and our leanings, all of
them open. Section 7 holds five questions.

Quotations from project documents are verbatim, including where they
capitalize project names; words in square brackets inside a quotation are
ours. In code excerpts a line `// …` marks lines we left out; everything else
is exact, with its file path and line numbers at the code baseline.

## 1. System overview

### 1.1 The components

| Component | What it is | Language(s) |
| --- | --- | --- |
| **ontos** | The value model under horos. A value is an atom (a byte string) or a tuple (an ordered sequence of values). It plays no further part in this document. | Go, Rust, TypeScript |
| **horos** | A type language whose types are ontos values. It has a declaration language with **modules**, generic declarations and imports. A module is identified by a **pin**, the SHA-256 digest of its canonical content, written as 64 hexadecimal digits. An import names a module by its pin. | Go, Rust, TypeScript |
| **bittype** | A Go library. It parses and analyzes horos source, and it generates native type declarations for TypeScript, Go, Rust and C++. Its native syntax covers type declarations only, not executable code. | Go; emits TypeScript, Go, Rust, C++ |
| **Consumers** | Applications that link bittype as a library, register the parts they want, and call generation and Files. Consumers' hand-written code then refers to the generated declarations by name. | Go |

A horos example used throughout, with the declarations a consumer asks for:

```text
Page<T> = { items: list<T>; next?: text; }       a generic record
Person  = { name: text; }                        a record
Cells   = Page<Cell<Person>>                     an alias; Cell is owned by a library
```

`Cell` stands for a type that a separate library owns in every target language
(`Cell` exported by the npm package `@bitspark/cell`, `cell.Cell[T]` in a Go
package, `Cell<T>` in a Rust crate). Generation must refer to it without
declaring it. In horos source the example writes `cells.Cell`, where `cells`
is the alias under which it imports the horos module that declares `Cell`.

### 1.2 Terms used throughout

Two house words first: a **carrier** is the data type that holds something, and
to **admit** an input is to validate and accept it, as opposed to refusing it.

| Term | Meaning in this document |
| --- | --- |
| **Target language, native** | One of the four output languages. **Native** means "of the target language": a native type, the native compiler. |
| **Native scope** | The language construct inside which declarations see each other without import or qualification: a TypeScript module (its file), a Go package, a Rust module, a C++ namespace. |
| **Syntax tree (AST)** | bittype's own tree representation of a program in some language. One generic carrier serves every language: a flat array of nodes, each with a kind, text attributes, child edges with role names, and reference links (section 4.1). Each language supplies a **vocabulary**: its node kinds, roles and attribute names. |
| **Site** | The address of one node: a tree ID and a node ID. Declarations and uses are identified by site, never by spelling. |
| **Parse, check, format, print** | For one tree: **parse** turns text into a tree, **check** admits a tree against the vocabulary's grammar rules (the **syntax check**; it resolves no names), **format** turns a tree into concrete syntax, **print** turns concrete syntax into text. The project calls this the **Syntax** dimension, and a **Syntax formatter** does only that formatting step. |
| **Files** | The second dimension: a set of trees, each linked to a file ID and a relative path. A **Files formatter** produces such a set for a chosen organization. A **generated-Files formatter** is one whose input is generation output; "the Go Files stage", "Go Files" and "the Go generated-Files formatter" are the same thing. Despite the name it relocates declarations, writes imports and verifies bindings. |
| **Source** | The result a Files stage publishes: the units (file ID, path, concrete syntax) with the compiler-backed receipt for exactly those units. A Files printer turns a Source into paths and text. bittype writes nothing to disk. |
| **Snapshot, capture, entry** | A **snapshot** is an immutable, syntax-checked set of trees, optionally with the files they belong to. To **capture** is to create one. Every snapshot names one **entry** node, the root of one of its trees; the file of that tree is the **entry file**. The capture contract requires it. horos resolution starts from it. No Go, Rust or TypeScript compiler bridge reads it. |
| **Contribution, provider** | A **contribution** is one independently owned part a consumer registers: a mapping library, a native target, a resolution or analysis provider, a Files formatter. The project calls every explicitly supplied implementation of a contract a **provider**; in the excerpts a type named `Provider` is whichever contribution its file path says. Nothing is discovered or registered automatically. |
| **Resolution and analysis providers** | Per-language contributions that run the language's real compiler front end over a snapshot (Go's `go/types` library, the TypeScript compiler API, a helper linked against `rustc`) and publish what it established. |
| **Semantic graph, fact, binding** | What such a provider publishes: **entities** (declarations, packages, modules), **bindings** (this use site resolves to that entity) and **facts** (provider-owned relations such as "this module exports that declaration"). Each language owns its names in the graph. |
| **Receipt, report, evidence** | A **receipt** is the result of one resolution or analysis, bound to one snapshot and one configuration. It carries **reports** (a provider's verdict, holds or rejected, on one named obligation such as "the program is valid") and, when demanded, a semantic graph. A stage states an **evidence requirement** (operation, profile, configuration, required reports), and a receipt is admitted against it. "Compiler-established evidence" always means a receipt. From section 3 on we also use "evidence" in the ordinary sense of whatever a check rests on, including tree equality. The **evidence profile** of a formatter is the semantic profile and configuration its evidence requirement names. |
| **Compiler bridge** | The private code that converts bittype's tree into the compiler's own input and exports facts back: `internal/golangsem` for Go, `internal/typescriptsem` for TypeScript, a Rust executable (the **helper**) for Rust. |
| **Profile** | A named, bounded subset of a language: a syntax or semantic profile says which constructs are supported. |
| **Kernel, shared contract** | The **generation kernel** is bittype's language-agnostic code that plans and runs generation. **The shared Files contract** is the package `contracts/format/go` of section 4.4, which every generated-Files formatter goes through; code comments call it the Files kernel. Both are "shared mechanisms" in the sense of the rules in section 2.1. |
| **Mapping library, declaration recipe** | A mapping library is a contribution that says how a horos form (a built-in kind of horos type, such as text, list or record) or a horos declaration is represented natively: text becomes `string`, a list becomes a slice. The **declaration recipe** is the mapping part that turns a named horos declaration into a native declaration and chooses its name. |
| **Native model, form** | Each language's closed set of native type **forms** (for example primitive, sequence, field, record, external) in which mapping libraries state a representation. A form carries string data and child forms. It is intent, not syntax. |
| **Native target** | The per-language contribution that constructs the syntax tree from the closed plan the kernel hands it (closed: the requested declarations and everything they reference, transitively). Where we write "the Go target writes" or "target configuration", we mean this contribution, not the language. The quoted documents also say "generation target". What a reference points at is called its **destination**, except in code identifiers and quotations. |
| **Intent** | What generation recorded: for a reference, its destination (a link in the tree, or an external use); for exposure, the exposed sites. It is a claim. Only analysis of the placed program yields facts. |
| **Exposed site, exposure, visibility** | An **exposed site** is the site of a declaration the caller requested as a root of generation. **Exposure** is that request; **visibility** is its native realization (`export`, `pub`, a capitalized Go name). The documents also say "public root" and "publication". They do not say to whom an exposed declaration is public (section 4.2). |
| **External use, external symbol** | A use site that refers to a type outside the generated set (the library-owned `Cell`), recorded with an opaque **external symbol** that identifies the destination. The symbol is issued by the language's small symbol package, which code comments call the "target library"; both the native target and the Files formatter of that language import it. |
| **Seal** | To validate a native target's result and freeze it into an immutable `Output` (section 4.2). |
| **Placement, placement assignment** | The **placement** is the consumer's whole statement of where generated declarations go. One element of it, which assigns some declaration sites to one file (and, in our proposals, one native scope), is a **placement assignment**. The placement also names the entry file. |
| **Dependency environment** | What the consumer says about the world the output is compiled in: the module root, the dependency packages and their declarations, compiler settings. |
| **Limitation** | A declared gap in what a check establishes. "Limit" is reserved for resource bounds. |
| **ACX-nnn** | A numbered architecture case record in bittype's repository: a found violation of the rules in section 2.1, with its correction and limitations. ACX-004 is "generation emits file syntax before placement". |

### 1.3 The stages, in one picture

```text
horos files ──parse, resolve, analyze──▶ horos snapshot + receipt
                                              │
                    generation kernel: plan   │  mapping libraries choose representations
                                              ▼
                              closed native program (intent, no syntax yet)
                                              │
                    native target: construct  │  one syntax tree, references linked to their destinations
                    syntax check, seal        ▼
        generation.Output = snapshot of the unplaced tree + exposed sites + external uses
                                              │
      consumer: placement (per request)       │   consumer: dependency environment (per formatter)
                                              ▼
              generated-Files formatter of the target language
                 1. read intent from the output; admit the placement
                 2. copy declarations into one tree per output file
                 3. write the file-level syntax the placement needs
                 4. capture the new trees with their paths; run the native
                    resolution and analysis providers
                 5. compare the compiler's bindings and facts with the intent
                                              │
                                              ▼
              Source (units + receipt)  ──print──▶ paths and text
```

Generation runs once. The same output can be handed to the Files stage several
times with different placements.

### 1.4 Scale and direction

- **Release state.** Nothing is released. The first release is a tracked work
  item that follows this work.
- **Consumers, as the plan states them.** The pilot consumer is a content
  service of the surrounding platform. Its contract becomes TypeScript first;
  Go follows, for the services' own code. No consumer is planned yet for Rust
  or C++ output. Both exist as generation targets and must not be closed off.
- **Order of work.** TypeScript is corrected. Go slice 1 and Rust slice 1 are
  about to start. C++ waits for its semantic tooling.
- **Modules per output.** We expect an output to span a handful of horos
  modules, as the documented walkthrough does (three). Declarations from every
  module but one get a hash prefix in their name (section 4.9), so that prefix
  is common in practice, not a corner case.
- **Sizes.** We expect tens to a few hundred declarations per output. The
  Files dimension defaults to at most 256 files. The Go semantic profile
  accepts at most 16 captured files per analysis. The Rust semantic profile
  accepts one captured source file, which may contain inline modules.
- **Setting keys carry a version.** The TypeScript placement travels under a
  key ending in `-1`. A later `-2` would be a change to what a consumer
  writes, and counts as a failure in the sense of the Question section.
- **People.** Each repository is maintained by AI coding-agent sessions,
  several working in parallel, so clear ownership and additive composition
  are hard requirements.

## 2. Architectural context

### 2.1 Standing architecture rules (AC = architectural correctness)

The owner has made architecture part of correctness: a design that works and
passes its tests is still wrong if a responsibility sits with the wrong owner.
The rules are numbered. The questions rest on seven of them. Each is given in
our words, followed by its requirement text, quoted. For AC-09 the quoted text
is the owner's standing direction of 8 November 2026, which the rule adopts.

| Rule | In our words | Requirement text (verbatim) |
| --- | --- | --- |
| AC-02 | Each decision belongs to the component that defines its meaning. A shared mechanism owns invariants and never a contribution's meaning. | "Let independently owned implementations supply their domain interpretations, policies and laws. Shared mechanisms own their invariant guarantees and depend on extension contracts; consumer composition selects the concrete contributions." |
| AC-03 | Adding a language, a mapping or a formatter touches only that contribution and the consumer that composes it. | "Within the agreed extension domain, adding a contribution must change only its owning implementation and consumer composition, not shared orchestration, closed case catalogues or unrelated implementations." |
| AC-05 | Passing a result from one stage to the next must not lose identity, scope, evidence or failure meaning. | "Preserve identity, scope, ownership, authority, assumptions, evidence association, resource limits and failure meaning through composition, as applicable to the contract. Selecting a contribution does not establish the truth of its claims." |
| AC-07 | A check supports only what it can actually observe. A test that cannot fail for the property in question establishes nothing about it. | "Distinguish behavioral tests, compiler acceptance, conditional proofs and architectural assessments. State the supported scope and assumptions; do not use one kind of evidence to claim another guarantee." |
| AC-08 | A shared contract may change, but only for a general reason, stated. | "When a shared contract really must change, explain the missing invariant or capability, why its responsibility belongs to that layer and how the revised contract serves independently owned contributions." |
| AC-09 | Interfaces are designed from the domain. The existing API and the cost of changing it have no weight. | "Whenever we design, redesign or improve an interface, choose the best contract for the domain requirements, responsibilities, ownership, composition and guarantees. The currently implemented surface must have **zero influence** on that choice. Implementation or migration effort has **zero weight**, including as a tiebreaker. [...] A retained element needs the same independent justification as a new one." |
| AC-10 | A misplaced decision moves to its owner first, as a small change, and that change is already final in structure. Coverage is added afterwards and removes nothing. | "Separate a correction into what is in the wrong place and what is only absent. Move a misplaced or wrongly named responsibility to its owner first, as its own change, without waiting for the complete implementation of its destination. That first change establishes the final structure: its inputs and identities have their intended shape, it keeps every guarantee the destination's shared mechanism requires for the cases it admits, and it explicitly refuses the cases it cannot yet check. Missing coverage is then added without moving, renaming or removing what exists." |

AC-10 also names its violations: "Moving a decision without the guarantee its
new owner must give. Delivering a first change whose input or identity later
has to change." It rests on two statements of the owner of 11 November 2026:
"if something is placed incorrectly, it is incorrect and must change. Something
that is missing can be healed additively, i.e. there is nothing wrong requiring
removal." And: "no compromise on architecture, but practical deferral is
allowed, if the huge effort justifies it."

**Package roles.** The rules are also applied to Go packages. A shared
mechanism imports no contribution; only a consumer's composition imports a
contribution; a contribution imports shared contracts and the vocabularies of
its own language. A test forbids every other import. This bears on question 1
in one way only: a shared placement type would live in a shared package and
could therefore hold nothing that only one language understands.

### 2.2 The accepted decision of 7 November 2026

This is the accepted design. It is recorded in bittype's `docs/source.md`,
section "From generated declarations to files", and marked "**Design decision,
7 November 2026.**" Its ownership table, quoted in full:

| Owner | Decision |
| --- | --- |
| Selected mapping and generation target | "Build declaration bodies and type expressions. Bind each named use to its intended target declaration, lexical parameter, or selected external public symbol. The plan records which declaration the caller requested to expose. A generated tree ID identifies syntax; it does not assign a file." |
| Consumer | "Supply one placement for the declarations being emitted: output files and the native scopes they belong to. Supply the selected external dependencies in the target environment. There is no second table of imports." |
| Target-language Files formatter | "Follow those bindings under that placement. Choose the local spelling, qualification or alias, and construct the necessary import, export or other module syntax. Realize the requested visibility under the language's rules and establish that uses still resolve to their intended targets." |
| Shared Files printer | "Emit the paths and text already represented by `Source`. The caller decides whether and where to write them." |

Further sentences of the same decision that the questions rest on. On names:

> A spelling alone is insufficient: two distinct types named `Person` must not
> become interchangeable when placement changes.

Of its two example placements (Person in its own file; Person in Page's file):

> Both placements use the same generated declarations, bindings and publication
> intent. Files must preserve those bindings and the requested public surface;
> moving a declaration cannot authorize exposing a private type or substituting
> another type of the same name. If it cannot express the placement with the
> supplied dependencies and visibility, it must fail without returning a
> partial `Source`. Native scope matters: two files in one Go package need no
> import between them, while two Rust modules can occupy one file. A file
> boundary alone is never an import rule.

On which visibility Files owns:

> The distinction is about module visibility. A member's accessibility, when it
> is part of the chosen type representation, still belongs to the declaration.

On what the boundary covers:

> This is the boundary for **generating declarations**. Native ASTs can still
> represent imports and exports when parsing existing source or constructing
> organized units.

### 2.3 The proposal of 8 November 2026

`docs/native-file-placement.md` walks the Page, Person and Cells example
through Go, Rust and C++. **It is a proposal and not an accepted decision.**
It opens: "**Proposal, 8 November 2026.** These examples specify the intended
Files boundary." It says of its own tables: "The concrete Go/Rust/C++ settings
and checks still need specification and implementation. The following tables
are proposed inputs/output, not shipped Go structs or APIs." ACX-004 and both
work items nevertheless treat it as the specification of the boundary.

Its proposed inputs, quoted:

| Input | Proposed meaning |
| --- | --- |
| Owned declaration placement | "Output file and native scope for each emitted declaration/helper [helper: an auxiliary generated declaration, not the Rust helper of 1.2]; entries identify declarations, not just spellings" |
| Public roots | "Generation's exposed declaration sites; placement cannot make an inaccessible dependency public merely to obtain a valid file" |
| External destination | "Library-selected [the mapping library] public package/crate/header symbol and mapped arguments" |
| Dependency environment | "Consumer-selected module/crate/include roots, actual dependency artifacts and relevant native configuration" |

Its six required placements, condensed from its tables. Each must be produced
from one generated output without regenerating:

| Language | Placement | What Files must write in the unit containing `Cells` |
| --- | --- | --- |
| Go | Page and Cells in `pages/page.go`, package `example.org/app/pages`; Person in `people/person.go`, package `example.org/app/people` | Imports of `example.org/app/people` and of the Cell package; the use `Page[cell.Cell[people.Person]]` |
| Go | Person in `pages/person.go`, the same package as Page | No import of Person; the use `Page[cell.Cell[Person]]` |
| Rust | Page and Cells in `src/pages.rs`; Person in module `people` at `src/people.rs` | `use crate::people::Person;` and `use bitspark_cell::Cell;` |
| Rust | "Person's bytes move into an inline `mod people` in the same physical file" | "Person still belongs to another module; use a suitable module path/import and satisfy its visibility" |
| C++ | Page and Cells in namespace `app::pages`; Person in `model/person.hpp`, namespace `app::model` | Includes of the Person header and the Cell header; the qualified name `app::model::Person` |
| C++ | "Person moves into `pages/page.hpp`, retaining namespace `app::model`" | No include of Person; the qualified name stays |

Further sentences of the proposal that the questions rest on.

(All languages, on external types:)

> External-symbol support for a future target must be represented by that
> target's model and the shared output contract, not guessed from a name.

(Go:)

> A package import path and an output filename are different identities; the
> consumer's Go module layout relates them.

> Moving a declaration across packages also changes its native package identity
> and can introduce a dependency cycle. Files must check visibility, identity
> correspondence and package cycles; it cannot silently export an unexported
> type or treat same-package movement as a package import. Package aliases must
> avoid names in the use's lexical scope.

(Rust:)

> Files must construct `mod` declarations or inline modules consistently with
> the placement, and realize public-root visibility under the crate/module
> rules.

For Rust the proposal also assigns the selection and renaming of dependency
crates to the dependency environment. For C++ it assigns include directives,
guards and declaration order to the C++ Files formatter, and notes: "A forward
declaration is not automatically enough: required completeness depends on the
selected native use."

(All languages, on checks:)

> Check collisions, visibility, missing external destinations, dependency
> cycles/completeness and public roots. The shared contract associates sites
> and evidence; language integrations establish native scope and destination
> meaning. [...] Successful printing or compilation alone does not establish
> that every intended use bound to the right declaration.

### 2.4 The order of correction, and what the two work items say about slice 1

ACX-004 records four misplaced decisions and their order (11 November 2026):

| Instance | First step | State |
| --- | --- | --- |
| Go package clause | Generation decides no package. Go Files writes the clause for one package and no import | Can start |
| Rust visibility | Generation chooses no visibility. Rust Files decides it for one module | Can start |
| Rust external paths | The Rust target reports an external use and Rust Files decides the import | Deferred until Rust semantics cover a second crate |
| C++ includes and namespace | Generation emits neither. C++ Files decides both | Deferred until C++ has resolution and analysis providers |

The case states the rule for the first step: it "moves the misplaced decision to
the language's generated-Files formatter, for the placements that existing
semantic evidence can check, and refuses every other placement. The second step
adds coverage and changes nothing the first delivered."

The two work items define slice 1 as follows, quoted:

| | Go work item | Rust work item |
| --- | --- | --- |
| Generation | "The Go target takes no `Package`." | "The Rust target sets no `visibility` attribute on declarations or fields" |
| The formatter | "It reads the package from the consumer's placement setting and writes the package clause." | "It decides the visibility of each declaration and its fields from the exposures and the consumer's placement." |
| The check | "The existing Go Files resolution and analysis providers supply both [a check and semantic evidence] for what this slice admits." | "The existing single-crate Rust Files providers supply both for what this slice admits." |
| Admits and refuses | "It admits one package and no import. A placement into more than one package fails, and an output with an external use fails." | "It admits one module. A placement into more than one module fails." |
| The setting | "It assigns each declaration a file and a native scope, as TypeScript's `Organization` assigns declarations to files. Slice 2 then admits more placements without changing the setting." | "It assigns each declaration a file and a module, so slice 2 admits more placements without changing it." |
| Verification | "One generated output is placed into package `pages` and into package `people` without regenerating, and both results pass the Go semantic check." | "A declaration that is not exposed is not made public." The existing compilation with rustc "passes on the placed output." |

**Standing of the work items.** They were written on 10 and 11 November 2026
to make the work startable. The plan that orders them says of itself: "It
decides no design." So they are revisable. Our leanings in section 6 depart
from them in five places: which Go semantic profile the check uses (6.6), how
predeclared types relate to "an output with an external use fails" (6.7),
fields (6.8), whether Rust's existing providers supply a real check for item
visibility (6.8), and whether each language gets its own setting type (6.5).
"Files providers" in the quotes are the resolution and analysis providers that
accept a snapshot with file paths. The Go work item was written before the
second of Go's two semantic profiles (section 4.7) existed, and says that
slice 1 does not wait for it.

## 3. Map: who decides what today, and where the code shows it

The table follows `Cells = Page<Cell<Person>>`. It states the present code,
not the intended design. "Recipe" is the declaration recipe of section 1.2.

| Decision | TypeScript (corrected) | Go today | Rust today | Shown in |
| --- | --- | --- | --- | --- |
| Name of a declaration | Recipe, during generation | Recipe | Recipe | 4.9 |
| Destination of each use | Link in the tree, or external use | Link; none for `string`, `bool`; library-owned types refused | Link; none for three absolute paths | 4.2, 4.5, 4.7, 4.8 |
| Root of the unplaced tree | `source-file`, no file decision | `file`, with a package | `file`, no file decision except visibility | 4.5, 4.7, 4.8 |
| Which file | Placement | One tree, no placement | One tree, no placement | 4.5 |
| Native scope | The file | The Go target's configuration | One module, implied | 4.5, 4.7, 4.8 |
| Visibility of a declaration | Files writes `export` for exposed sites; checked against a compiler fact | The case of the name, from the recipe | The Rust target writes `pub` on everything | 4.5, 4.8, 4.9 |
| Imports and aliases | Files; checked through compiler bindings | None | None (absolute paths) | 4.5 |
| Name collisions | Flat rule during generation, and per file in Files | Flat rule | Flat rule | 4.5, 4.9 |
| Compiler evidence of the check | TypeScript import environment, required when the formatter is built | No Files stage | No Files stage | 4.5, 4.7, 4.8 |

What each part of section 4 is for:

| Part | The one thing to take from it | Feeds question |
| --- | --- | --- |
| 4.1 Carriers | A site is a node address; reference links carry intent; a setting is opaque bytes; bindings and facts are what a compiler publishes | all |
| 4.2 Generation output | What an output records, how a consumer finds a declaration's site, and that the same syntax checker seals it | 1, 3 |
| 4.3 Shared tree builder | Every target's tree is built by one function; the root is written last | 2, 3 |
| 4.4 Shared Files contract | It refuses a formatter without a check and compiler evidence, and it never reads a placement | 1, 4 |
| 4.5 TypeScript | The working model: what a consumer writes, what is admitted, and a check in three parts | 1, 2, 4 |
| 4.6 horos Files | A second implementation of placement admission, for a different kind of item | 1 |
| 4.7 Go | `file` requires a package; the package name is in no graph; `string` and `bool` carry no intent | 3, 4 |
| 4.8 Rust | Visibility is written by the target; the compiler facts cover fields only; "private" is an absent attribute | 3, 4 |
| 4.9 Names | Where names come from, and the two places that require them to be unique | 2 |
| 4.10 C++ | Why a language without a Files stage still constrains the shape | 1, 3 |

## 4. Relevant code

All code is Go unless marked, at the code baseline. In the excerpts, short
import aliases follow one pattern: `s` is the syntax vocabulary of the language
at hand, `n` its semantic vocabulary, `g` the generation contract package
(inside the grammar checkers, `g` is a small grammar helper), `m` the shared
tree builder, `f` the shared Files contract package and `model` the language's
native model. `semantic` is the shared semantic contract (snapshots, settings,
receipts); `semantics` is the graph model (keys, entities, bindings, facts).
`goast`, `types` and `token` are Go's `go/ast`, `go/types` and `go/token`;
`tr`, `native` and `imports` are private helpers. Inside a function a single
letter is sometimes a local variable instead. In section 4.1 declarations from
several places are listed together, without the comments between them.

### 4.1 Carriers

Four things matter here: `NodeRef` is the site; `References` carry intent;
a `Setting` is opaque bytes; `Binding` and `Fact` are what a compiler bridge
publishes. One tree carrier serves every language. `L` is a marker type owned
by the language, so trees of different languages cannot be mixed.

```go
// trees/ast/go/ast.go:13, 17, origin.go:7-10, node.go:28-31, 38-44
type TreeID string
type NodeID uint32
type NodeRef struct {
	Tree TreeID
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

A `NodeID` is a one-based index into its tree's node list; zero names no node.
Entity and term IDs in the graph are one-based in the same way. A tree
(`AST[L]`) is an ID, a root node ID and a flat list of such nodes. An
`Attribute` is a name and a text. A `Child` is a role and a node ID:
containment. `References` are resolved, non-containment links: a type reference
that names a generated declaration carries a link to that declaration's node.
An unresolved name is only an attribute. `Origin` records where a node came
from and plays no part here.

```go
// contracts/semantic/go/contract.go:82-85
type Setting struct {
	Key   semantics.Key
	Value []byte
}

// semantics/go/model.go:30-37, 77-86, 110-114
type Entity struct {
	Kind         Key
	Name         string
	Owner        EntityID
	Declarations []ast.NodeRef
	Parameters   []Parameter
	Body         TermID
}
type Occurrence struct {
	Site    ast.NodeRef
	Role    Key
	Ordinal uint32
}
type Binding struct {
	At     Occurrence
	Target EntityID
}
type Fact struct {
	Subject  Subject
	Relation Key
	Term     TermID
}
```

A `Key` is an owner and a name, for example `go/package`. A `Setting` is how
configuration, including a placement, travels through shared code that must
not interpret it. An `EvidenceRequirement` names an operation, a profile, a
configuration and the required reports, and says in its field `Semantics`
whether a semantic graph is demanded.
An `Occurrence` is one position within a node: its site, plus a role and an
ordinal that the provider owns. A term is a provider-owned expression in the
graph, for example a type; a fact's `Subject` is an entity or a term.

Files are a list of `File{ID, Path, Text}`. IDs and paths are unique within
one list; paths are relative and slash-separated; `Files.Capture(limits)`
admits such a list and refuses duplicates and invalid paths. The caller chooses
the directory the paths are relative to. `representation.Language` is a
language descriptor: a language ID, the name of its syntax profile and the
marker type. `a.Compatible(b)` holds only if all three are equal. `cst.CST[L]`
is the concrete syntax tree that formatting produces and printing consumes.

### 4.2 What generation produces

An output carries a snapshot, the exposed sites and the external uses, and
nothing else:

```go
// contracts/generation/go/output.go:14-21, 37-60
// SymbolIssuer is a target library's authority to issue and interpret immutable
// external destinations. Its key is opaque to the generation and Files kernels.
type SymbolIssuer struct{ owner *identity }
type ExternalSymbol struct {
	owner    *identity
	language representation.Language
	key      string
}
// …
type ExternalUse struct {
	At     semantics.Occurrence
	Symbol ExternalSymbol
}

// Assembly is a target construction record. Exposed contains the realized
// declaration site of each publication in the supplied publication order.
type Assembly struct {
	Syntax       representation.ASTSet
	Exposed      []ast.NodeRef
	ExternalUses []ExternalUse
}

// Output can only be produced by a successfully finished generation session.
// Sites and immutable target symbols always travel with their own checked capture.
type Output struct {
	snapshot semantic.Snapshot
	exposed  []ast.NodeRef
	external []ExternalUse
}

func (o Output) Snapshot() semantic.Snapshot { return o.snapshot }
func (o Output) Exposed() []ast.NodeRef      { return slices.Clone(o.exposed) }
func (o Output) ExternalUses() []ExternalUse { return slices.Clone(o.external) }
```

An issuer's `Issue(language, key)` creates a symbol and its `Read` returns the
key only to the issuer that created it. `representation.ASTSet` is a set of
trees of one language with an entry node. "Publication" in these comments means
an exposed declaration.

The native target returns an `Assembly`. The kernel seals it: it checks that
every exposed site and every external-use site lies inside the output's trees,
and then captures the trees with the target language's syntax checker
(`check.Checker`):

```go
// contracts/generation/go/output.go:62-66, 108-113
func seal(ctx context.Context, checker check.Checker, assembly Assembly, publications int, limits Limits) (Output, error) {
	set := assembly.Syntax
	if len(set.Trees) == 0 || len(set.Trees) > 256 || len(assembly.Exposed) != publications {
		return Output{}, ErrIncomplete
	}
	// …
	snapshot, err := semantic.Capture(ctx, checker, semantic.Input{Entry: set.Entry, Trees: set.Trees})
	if err != nil {
		return Output{}, err
	}
	return Output{snapshot, assembly.Exposed, assembly.ExternalUses}, nil
}
```

Intent for a reference to a generated declaration or to a type parameter is
the `Reference` link inside the tree. Intent for a reference to a library-owned
type is an `ExternalUse`. As the first quoted condition shows, a sealed output
may hold up to 256 trees. The request that starts generation names one tree
ID, and the kernel requires only that the output's entry tree has that ID.
Every native target today returns exactly one tree, because the shared tree
builder of section 4.3 builds one.

**How a consumer obtains a site.** A placement names declarations by site.
`Exposed()` returns the sites of the requested declarations, in the order of
the request. For every other declaration in the output (a dependency the
closure pulled in, or a generated support declaration such as Go's
`BittypeOptional`) the output has no accessor. The consumer must open the
snapshot's tree, walk the root's declaration children and recognize each by
its name attribute. The TypeScript acceptance test does exactly that. A
complete placement must assign those declarations too. No API maps a horos
declaration to its site. So today a consumer states identity by site but finds
an unexposed site by name, which works only because names are unique in one
output (section 4.9). If that uniqueness goes (option U2 in section 6.2), the
lookup by name becomes ambiguous.

**To whom an exposed declaration is public.** The documents do not say. They
define exposure as the caller's request and require that an unexposed
declaration never becomes public "merely to obtain a valid file". The only
realization, TypeScript's, is two-valued: an exposed declaration is exported
from its module, and an unexposed one can be used only inside its own file. A
placement that separates an unexposed declaration from its user is refused
(section 4.5).

### 4.3 The shared tree builder (`Compile`)

Every native target implements this interface:

```go
// contracts/generation/go/contracts.go:37-42
type Target[T any] interface {
	CheckExpression(context.Context, Expression[T]) error
	CheckDeclaration(context.Context, NativeDeclaration[T]) error
	CheckProgram(context.Context, Program[T]) error
	Build(context.Context, BuildInput[T]) (Assembly, error)
}
```

The generation kernel calls `Build` and seals whatever `Assembly` comes back.
All four targets implement `Build` by calling one shared helper function,
`Compile`, which walks the planned program. `Compile` is a helper the targets
chose, not a kernel invariant.

`CheckProgram` is called while planning; `Build` constructs the tree. The plan
a target receives is a `g.Program`: a list of planned declarations, each with a
numeric `ID` and a `Description` holding its `Name`, its type `Parameters` and
its body. `Compile` identifies every declaration and type parameter by a **plan
symbol**: the declaration's ID plus a parameter index. It is unrelated to the
external symbol of section 1.2. When the target creates a reference node it
calls `Link(node, planSymbol, role)`, and `Compile` later turns each link into
a `Reference` to the node the symbol ended up at. For a library-owned type the
target calls `External(node, use)` instead, and `Compile` collects those into
`Assembly.ExternalUses`. Finally the target's `Finish(declarations)` wraps the
declaration nodes in the root node. That call is where Go writes its package
and C++ its namespace.

In the target excerpts below, `b.Node(m.Node(kind, attributes, origin),
children…)` creates one node and returns a fragment, a handle to it.
`m.Attributes("name", x, …)` builds an attribute list from name and value
pairs, and `edges(role, fragments)` turns fragments into child edges with one
role. Where a target builds an expression from a native form, `name` is the
form's name, `data` its string data and `args` its already built children.

### 4.4 The shared Files contract

A language supplies a generated-Files formatter by giving the shared contract
an implementation (the `provider` below) and three things: an evidence
requirement, a printer and a check.

```go
// contracts/format/go/generated.go:13-28
// GeneratedFilesInput carries intended bindings from a sealed run and the
// consumer's placement settings. The target supplies its dependency environment.
type GeneratedFilesInput struct {
	Output        generation.Output
	Configuration []semantic.Setting
	Limits        FilesLimits
}
type GeneratedFilesContract[L any] struct {
	Evidence semantic.EvidenceRequirement
	Print    func(context.Context, cst.CST[L]) (string, error)
	Check    func(context.Context, GeneratedFilesInput, source.Publication[L]) error
}
type GeneratedFilesRequest[L any] struct {
	input     GeneratedFilesInput
	publisher source.Publisher[L]
}
```

Construction refuses a formatter that has no printer, no check, or an evidence
requirement that does not demand a semantic graph:

```go
// contracts/format/go/generated.go:55-57
	if representation.IsNil(provider) || provider.ID() != language.ID() || contract.Print == nil || contract.Check == nil || !contract.Evidence.Semantics {
		return GeneratedFilesFormatter[L]{}, ErrInvalidFormatter
	}
```

Each request requires the formatter's language to be compatible with the
language of the output's snapshot, passes the settings through unread, and
installs the language's check as a guard:

```go
// contracts/format/go/generated.go:77-80, 96-103
	snapshot := input.Output.Snapshot()
	if ctx == nil || !snapshot.Valid() || !p.definition.Compatible(snapshot.Language()) {
		return source.Source[L]{}, source.ErrInvalid
	}
// …
	policy := source.Policy[L]{Evidence: p.contract.Evidence, SemanticLimits: limits.Semantic, Limits: limits.Source, Print: p.contract.Print, Check: func(ctx context.Context, output source.Publication[L]) error {
		if snapshot.Same(output.Snapshot) {
			return source.ErrStale
		}
		guard := input
		guard.Configuration, _ = semantic.CaptureSettings(input.Configuration, limits.Semantic)
		return p.contract.Check(ctx, guard, output)
	}}
```

Here `input` is the `GeneratedFilesInput` of one request, `p.definition` is
the formatter's language, `limits` is the normalized form of `input.Limits`,
and `snapshot.Same(output.Snapshot)` refuses publishing the input snapshot
itself. The publisher of this request then works under this policy.

A formatter cannot return a `Source` it built itself. It hands a
`source.Publication` to `Publish`: a fresh snapshot of the placed units, their
analysis receipt, and the units. (This "publication" is the placed result, not
an exposed declaration.) A publication can also carry **groups**, each naming
files that form one native unit, and **relations** between files. TypeScript
states one group per file, of kind "module", and one relation per import
between output files. The shared code checks only that they refer to existing
files and sites. The shared publisher then:

1. admits the receipt against the contract's evidence requirement, for exactly
   the published snapshot;
2. requires the units to cover the snapshot completely, their trees to equal
   the captured trees, and their printed text to equal the captured file text;
3. runs the language's `Check` with the original input.

Three properties of this code matter for the questions. All three are
implemented today and open to change. The unplaced output and the placed
source have one language descriptor, so one syntax checker admits both. The
shared contract carries the placement as settings and never reads it, so each
formatter decides privately what a placement is and whether it is complete.
And the evidence requirement is fixed when the formatter is built, from the
dependency environment, while the placement arrives with each request.

### 4.5 TypeScript, the corrected path and working model

TypeScript was corrected on 8 November 2026. Before that, its native target
wrote import statements and `export` modifiers during generation, and the Files
stage discarded and rewrote them. ACX-004's record of that correction calls
this the generated syntax "doing double duty as the record of which dependency
was intended".

**What a consumer writes**, as the implementation document shows it. `tsfiles`
is the formatter's package; `page`, `cells` and `person` are sites in
`out.Snapshot()`; `stages` holds the collaborators and the dependency
environment. That environment is a list of npm package names, each with the
text of its type declarations:

```go
// docs/typescript-type-imports.md:25-42
out, err := eng.Generate(ctx, request)
if err != nil { return err }
organization := tsfiles.Organization{
    Entry: "pages",
    Units: []tsfiles.Placement{
        {File: "pages", Path: "generated/pages/page.ts", Declarations: []ast.NodeRef{page, cells}},
        {File: "model", Path: "generated/model/person.ts", Declarations: []ast.NodeRef{person}},
    },
}
setting, err := organization.Setting()
if err != nil { return err }
formatter, err := tsfiles.NewGenerated(stages)
if err != nil { return err }
organized, err := formatter.Format(ctx, format.GeneratedFilesInput{
    Output: out, Configuration: []semantic.Setting{setting},
})
if err != nil { return err }
emitted, err := filesPrinter.Print(ctx, organized)
```

Section 4.2 says where such sites come from: `Exposed()`, or a walk over the
generated tree that reads each declaration's name.

**Generation side: the TypeScript target writes no file decision.** Its
`Finish` creates the vocabulary's top node with declaration children only.
`string` and `boolean` are keyword nodes. For a library-owned type the target
creates a type-reference node without a link and attaches an external use,
which TypeScript's symbol package builds from a module and a name: "a named
export, or a global type when Module is empty". The use is created with an
empty site; `Compile` fills the site in when the node receives its ID. No new
node kind was needed, because the TypeScript grammar requires nothing on its
top node, and `exported` is an optional attribute of a declaration:

```go
// check/syntax/typescript/provider/go/check.go:19
	s.KindSourceFile:        {Roles: []g.Role[s.Language]{g.Many(s.RoleImport, 0, s.KindImport), g.Many(s.RoleDeclaration, 0, s.KindTypeAlias, s.KindInterface)}},
```

(`g.Many(role, min, kinds…)` is a repeated child role; the rule lists no
attribute.)

The formatter then works in the five steps of the picture in section 1.3.

**Step 1a: read intent, and refuse a file decision in generated syntax.** The
formatter rebuilds its view of the generated program from the sealed output
alone. A root child that is not a declaration, and a declaration that carries
an `exported` attribute, are both refused. Exposure is read only from
`Exposed()`. In this package `imports` is a private TypeScript helper
(`imports.Attr` reads an attribute, `imports.Export{Module, Name}` names an
export of a module), `ref(tree, node)` builds a site, and
`tssymbol.ReadExternal` asks TypeScript's symbol package to read an external
symbol back as `{Module, Name}`. The excerpt is the body of the loop over the
output's trees: `t` is one tree, `snap` the output's snapshot, and `p` the
formatter's view (`p.tops`: top-level declaration sites; `p.public`: exposed
ones; `p.bindings`: intended destinations):

```go
// format/files/typescript/provider/go/generated.go:70-115
		for _, c := range t.Nodes[t.Root-1].Children {
			if c.Role != s.RoleDeclaration {
				return nil, source.ErrInvalid
			}
			site := ref(t.ID, c.Node)
			if imports.Attr(t.Nodes[c.Node-1], s.AttributeExported) != "" {
				return nil, source.ErrInvalid
			}
			p.tops[site] = true
		}
		for i, node := range t.Nodes {
			// …
			if node.Kind == s.KindTypeReference && len(node.References) > 0 {
				if len(node.References) != 1 || node.References[0].Role != s.RoleTarget {
					return nil, source.ErrInvalid
				}
				p.bindings[ref(t.ID, ast.NodeID(i+1))] = symbol{site: node.References[0].Target}
			}
		}
	}
	for _, site := range output.Exposed() {
		if !p.tops[site] {
			return nil, source.ErrInvalid
		}
		p.public[site] = true
	}
	for _, use := range output.ExternalUses() {
		site := use.At.Site
		ext, err := tssymbol.ReadExternal(snap.Language(), use)
		if err != nil {
			return nil, err
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
	return p, nil
}
```

In this excerpt `symbol` is a private struct of the formatter, its record of a
destination: a local site, an external `{Module, Name}`, or a compiler builtin.
After this step every type reference has an intended destination, or
formatting fails.

**Step 1b: admit the placement.** The placement is a TypeScript-owned type in
the formatter's package. It has no scope field, because a TypeScript module's
scope is its file. It travels as JSON in one setting:

```go
// format/files/typescript/provider/go/placement.go:21-32
// Placement assigns declaration sites in one admitted input capture to a file.
type Placement struct {
	File         files.ID
	Path         files.Path
	Declarations []ast.NodeRef
}
type Organization struct {
	Entry files.ID
	Units []Placement
}

func OrganizationSetting() semantics.Key { return n.Key("declaration-placement-1") }
```

In this document's terms an `Organization` is the placement, a `Placement` is
one placement assignment, and `Entry` names the entry file (section 1.2). The
admission is private to the same package. `p.tops` is the set of top-level
declaration sites of the input, and `p.original` is the organization of an
already placed program (empty for generation output):

```go
// format/files/typescript/provider/go/placement.go:168, 182-210
func (p *program) organization(settings []semantic.Setting, limits source.Limits) (Organization, error) {
	// …
	if len(o.Units) > limits.Files.MaxFiles {
		return Organization{}, source.ErrLimit
	}
	inventory := files.Files{}
	assigned := map[ast.NodeRef]bool{}
	entry := false
	for _, u := range o.Units {
		if !imports.SourcePath(u.Path) || len(u.Declarations) == 0 {
			return Organization{}, source.ErrInvalid
		}
		inventory.Entries = append(inventory.Entries, files.File{ID: u.File, Path: u.Path})
		if u.File == o.Entry {
			entry = true
		}
		for _, site := range u.Declarations {
			if !p.tops[site] || assigned[site] {
				return Organization{}, fmt.Errorf("%w: duplicate or unknown declaration %v", source.ErrInvalid, site)
			}
			assigned[site] = true
		}
	}
	if !entry || len(assigned) != len(p.tops) {
		return Organization{}, fmt.Errorf("%w: incomplete placement or missing entry", source.ErrInvalid)
	}
	if _, err := inventory.Capture(limits.Files); err != nil {
		return Organization{}, err
	}
	return o, nil
}
```

The elided lines decode the one setting strictly as JSON, refusing unknown
fields. One line of the rest is about TypeScript: `imports.SourcePath` (a
`.ts` path that is not `.d.ts`). One assignment is one file, because the Files
list refuses a repeated file ID or path. The table at the end of this section
names the rules that are not TypeScript's.

**Steps 2 and 3: copy declarations, check collisions per file, write file
syntax.** Declarations are copied into one new tree per assignment, and a map
records every old site's new site. Names are checked per file:

```go
// format/files/typescript/provider/go/relocate.go:119-126
		unit.tree.Root = imports.AddNode(&unit.tree, s.Node{Kind: s.KindSourceFile})
		names := map[string]bool{}
		for _, site := range unit.placement.Declarations {
			name := imports.Attr(program.trees[site.Tree].Nodes[site.Node-1], s.AttributeName)
			if names[name] {
				return fmt.Errorf("%w: duplicate declaration %s in %s", source.ErrInvalid, name, unit.placement.Path)
			}
			names[name] = true
```

An exposed declaration gets `exported: true`. A use whose destination is in
the same file is re-spelled from the name of the declaration it links to,
never from its own text, and linked to the new site. A use of a declaration in
another file becomes a named type import with a relative path, and is refused
unless that declaration is exposed ("cross-unit type %s must be exported"). A
use of an external symbol becomes a named type import from its package. The
import helper picks an alias such as `Cell_1` when a local name collides. A
declaration is never renamed.

**Step 4: capture, and obtain fresh compiler evidence under the consumer's
environment.** The formatter cannot be built without the environment: its
constructor fails with "named TypeScript import environment required".

Its contract evidence is the analysis profile under exactly that environment,
and its check recomputes the layout. This is the template a Go or Rust
formatter would follow:

```go
// format/files/typescript/provider/go/generated.go:28-33, 44-50
	evidence, err := n.EvidenceFor(p.Configuration)
	if err != nil {
		return f.GeneratedFilesFormatter[s.Language]{}, err
	}
	return f.NewGeneratedFilesFormatter(p.Language, f.GeneratedFilesContract[s.Language]{Evidence: evidence, Print: p.Printer.Print, Check: checkGenerated}, GeneratedProvider{stages: p})
}
// …
func checkGenerated(ctx context.Context, input f.GeneratedFilesInput, publication source.Publication[s.Language]) error {
	layout, err := placeGenerated(ctx, input)
	if err != nil {
		return err
	}
	return checkLayout(ctx, layout, publication)
}
```

The formatter formats each new tree, captures the units with their paths, runs
the resolver and the analyzer, and hands the result to `Publish`.

**Step 5: the check.** `placeGenerated` is the formatter's own private layout
function (steps 1 to 3). The guard runs it a second time on the original input
and compares the result with what was published. `expected` is the recomputed
layout: `expected.moved` maps old sites to new sites, `expected.public` holds
the exposure of each new top-level site, and `u.expected` holds the intended
destination of each type reference in unit `u`. `actual.Source` is the ID of
the text a unit is printed to. `realized` is read from the fresh semantic graph
of the published snapshot:

```go
// format/files/typescript/provider/go/format.go:197-224
	realized, err := readProgram(ctx, out.Snapshot, graph)
	if err != nil {
		return err
	}
	for site, want := range expected.public {
		if !realized.tops[site] || realized.public[site] != want {
			return fmt.Errorf("%w: public exposure changed at %v", source.ErrInvalid, site)
		}
	}
	for i, u := range expected.units {
		actual := out.Units[i]
		if u.placement.File != actual.File || u.placement.Path != actual.Path || actual.Source != ast.SourceID(actual.File) || !reflect.DeepEqual(u.tree, actual.Syntax.Abstract) {
			return source.ErrInvalid
		}
		if u.placement.File == expected.entry && out.Snapshot.Entry() != ref(u.tree.ID, u.tree.Root) {
			return source.ErrInvalid
		}
		for site, original := range u.expected {
			want := original
			if want.site.Tree != "" {
				want.site = expected.moved[want.site]
			}
			if got, ok := realized.bindings[ref(u.tree.ID, site)]; !ok || got != want {
				return fmt.Errorf("%w: moved reference changed its compiler binding in %s", source.ErrInvalid, u.placement.Path)
			}
		}
	}
	return nil
```

So the check has three independent parts:

1. **Tree equality by reconstruction.** Each published tree equals the tree the
   guard recomputed from the input.
2. **Bindings, from the compiler.** `realized.bindings` is built from the
   graph's `Binding` records, classified as shown below. For every relocated
   type reference, the entity the compiler resolved it to must be the
   relocated intended declaration, the intended package export, or the
   intended builtin.
3. **Exposure, from the compiler.** `realized.public` is built from
   `named-export` facts, which TypeScript's bridge publishes for each module
   entity. A declaration is exported in the result exactly if it was exposed.

Parts 2 and 3 are facts the formatter did not produce itself. This is how a
compiler binding becomes a destination (`data` is the graph's content, `b` one
binding):

```go
// format/files/typescript/provider/go/placement.go:145-160
		entity := data.Entities[b.Target-1]
		var dest symbol
		switch {
		case len(entity.Declarations) == 1:
			dest.site = entity.Declarations[0]
		case entity.Kind == n.PackageTypeKind() && entity.Owner != 0:
			owner := data.Entities[entity.Owner-1]
			if owner.Kind != n.PackageKind() {
				return nil, source.ErrInvalid
			}
			dest.external = imports.Export{Module: owner.Name, Name: entity.Name}
		case entity.Kind == n.BuiltinKind():
			dest.builtin = entity.Name
		default:
			return nil, fmt.Errorf("%w: ambiguous native binding %s", source.ErrInvalid, entity.Name)
		}
```

The acceptance test for this path generates `Cells = Page<Cell<Person>>` once,
together with a local declaration also named `Cell`, and formats it twice.
Moving `Person` next to `Page` removes only its relative import; the package
import stays, aliased `Cell_1`. The test also asserts that the generated tree
contains no import node and no `exported` attribute.

**What in this path is TypeScript's, and what is not:**

| Part | TypeScript-specific | Not specific to TypeScript |
| --- | --- | --- |
| Refusing file decisions in generated syntax | Which nodes and attributes count (`import`, `exported`) | That generated syntax carries none |
| Placement | The path rule; that a scope is a file | Assignment by site, complete and exactly once; entry file; file IDs and paths; strict decoding; limits |
| Relocation | Import and alias syntax; relative paths | The map from old sites to new sites; re-spelling a use from its destination |
| Collisions | What a scope is | That names are unique per placed scope |
| Check | Which facts mean "exported" | Tree equality; binding comparison at every relocated use; fresh evidence under the environment |

### 4.6 The horos Files formatter's placement

The horos Files formatter places whole horos modules, identified by pin, into
files. It has its own placement type and its own admission, in a private
package:

```go
// internal/horosfiles/organization.go:18-24
// Placement is an explicit same-language organization decision. Modules are
// canonical module pins, never aliases selected by spelling or filesystem lookup.
type Placement struct {
	ID      files.ID
	Path    files.Path
	Modules []string
}
```

Its admission applies the file limit, admits the file IDs and paths through
the same Files list, refuses an empty assignment, refuses a pin that is
unknown or already assigned, and requires every module to be assigned. So it
enforces the same kind of rules as TypeScript's admission, for a different
kind of item: module pins, where TypeScript has declaration sites. Like
TypeScript's, it allows one assignment per file. The entry is not named by the
placement here; it is the file that receives the entry module.

### 4.7 Go as it is

Four things to take from this part:

1. The Go vocabulary's only root, `file`, requires a package name.
2. The Go target writes that name from its configuration, and writes `string`
   and `bool` with no recorded destination.
3. Go has two semantic profiles. In the base profile all files are one package
   under a constant path. Under an import environment, files form packages by
   directory.
4. In neither profile does the graph contain the package name.

#### What generation writes today

The checker requires a package identifier on `file` and accepts no other root
(`g.Attr(true, …)` marks a required attribute):

```go
// check/syntax/golang/provider/go/check.go:18, 45-49
	s.KindFile:              {Attributes: g.Attributes(s.AttributePackage, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.Many(s.RoleImport, 0, s.KindImport), g.Many(s.RoleDeclaration, 0, s.KindTypeDeclaration, s.KindTypeAlias)}},
// …
func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, s.KindFile, rules); err != nil {
		return err
	}
```

The vocabulary package describes itself as containing "no parser, builder,
grammar checker, generator or formatter". The Go Syntax formatter prints
`package <name>` first, from the same attribute.

The Go target takes the package name from its own configuration, refuses to
build without one, and writes it on the single file node. A reference to a
generated declaration is a named type linked to its destination. `string` and
`bool` are also named types, with no link and no external use:

```go
// transform/syntax/golang/target/go/provider.go:15, 27-33, 78-80, 88-89
type Provider struct{ Package string }
// …
func (p Provider) Build(ctx context.Context, in g.BuildInput[s.Language]) (g.Assembly, error) {
	if p.Package == "" || token.Lookup(p.Package).IsKeyword() {
		return g.Assembly{}, g.ErrInvalid
	}
	b := &Backend{Store: m.New[s.Language](ctx, in.Limits), config: p}
	return m.Compile(ctx, in, b.Store, b)
}
// …
func (b *Backend) Finish(declarations []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindFile, m.Attributes("package", b.config.Package), ast.NodeRef{}), edges(s.RoleDeclaration, declarations)...)
}
// …
	case model.FormPrimitive:
		return b.Node(m.Node(s.KindNamedType, m.Attributes("name", data[0]), origin))
```

`data[0]` is `"string"` or `"bool"`. In Go these are identifiers of the
universe scope, not keywords: a package-level declaration named `string` would
shadow them. TypeScript's `string` is a keyword node and cannot be shadowed.
The Go native model refuses library-owned types altogether today (its check
returns "unsupported" for the external form), so Go generation reports no
external use.

#### What compiler evidence can see

**The package name, the word after `package`, is a fact in no graph, in either
profile.** The rest of this part shows why.

The compiler bridge converts captured trees to Go's own `go/ast` and runs
`go/types`. It refuses a root that is not a `file`, and it reads the package
attribute. `c.imports` is true only under an import environment, so the base
profile refuses any import node:

```go
// internal/golangsem/bridge.go:360-366
func (c *converter) file() (*goast.File, error) {
	node := c.node(c.tree.Root)
	if node.Kind != s.KindFile || !c.imports && len(children(node, s.RoleImport)) != 0 {
		return nil, n.ErrUnsupported
	}
	f := &goast.File{Name: goast.NewIdent(attribute(node, s.AttributePackage)), Package: c.pos(c.tree.Root)}
	f.Name.NamePos = c.pos(c.tree.Root)
```

*The base profile.* All captured files are one package, checked under a
constant path, whatever their directories:

```go
// semantics/golang/go/vocabulary.go:16
const PackagePath = "captured"

// internal/golangsem/bridge.go:210-211
	conf := types.Config{GoVersion: n.LanguageVersion, Sizes: types.SizesFor("gc", "amd64")}
	pkg, err := conf.Check(n.PackagePath, fset, goFiles, info)
```

*The import environment.* It is Go's form of the dependency environment of
section 1.2, and TypeScript's counterpart is the environment of section 4.5.
It was merged on 12 November 2026, after the work items were written.
Supplying one as configuration selects the second profile. The consumer gives
a module path and the source text of every dependency package (an import path
and its files). Captured files then form packages by their directories:

```go
// semantics/golang/go/environment.go:16-24, 102-109
// Import environments bound the imported-packages extension of the first Go
// profile. A captured file's package import path is the module path joined with
// the file's directory. Every other imported package is supplied here as Go
// source; nothing is discovered from a module cache, GOPATH or the standard
// library. Package text participates in receipt identity and is never an output.
type ImportEnvironment struct {
	Module   string
	Packages []Package
}
// …
// CapturedImportPath is the import path of a captured file at filePath under module.
func CapturedImportPath(module, filePath string) string {
	dir := path.Dir(filePath)
	if dir == "." {
		return module
	}
	return module + "/" + dir
}
```

Under this profile the bridge groups captured files by that import path and
requires the files of one group to state one package name ("package %s;
expected package %s"). It rejects import cycles, imports of a `main` package
and uses of unexported names across packages, each as a failed validity
report.

In both profiles a package entity is named by its **path**, and a top-level
declaration's entity is owned by its package entity. A predeclared type gets
an entity of kind `builtin`:

```go
// internal/golangsem/export.go:129-141, 305-310
	if _, ok := o.(*types.TypeName); ok {
		kind = n.Key("builtin")
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
// …
	// Captured packages come first, in import-path order, then supplied packages.
	// The base profile has exactly one captured package, which is entity 1.
	for _, captured := range program.captured {
		e.data.Entities = append(e.data.Entities, semantics.Entity{Kind: n.PackageKind(), Name: captured.pkg.Path(), Declarations: slices.Clone(captured.roots)})
		e.packages[captured.pkg] = semantics.EntityID(len(e.data.Entities))
	}
```

| | Base profile | Import environment |
| --- | --- | --- |
| Package identity in the graph | The constant `captured` | Module path joined with the file's directory, for example `example.org/app/pages` |
| Meaning of a file's directory | None: all files are one package | Go's: one directory, one package |
| The package name | Read by the bridge, exported nowhere | Read by the bridge to compare the files of one directory and to refuse importing `main`; exported nowhere |
| `string`, `bool` | Bind to `builtin` entities | The same |

So in the base profile, placing one output into package `pages` and then into
package `people` yields two identical graphs.

**Evidence follows the configuration.** Go's semantic vocabulary package
exports two functions, `EvidenceFor(settings)` and
`AnalysisContractFor(settings)`. They return the evidence requirement, and the
contract a resolver and analyzer must be configured with, for exactly the
given settings: the base configuration, or the base configuration extended by
one normalized import environment. The reader that opens a Go receipt accepts
a receipt made under any valid configuration. A consumer that needs one
particular environment must therefore require `EvidenceFor(thatEnvironment)`.
A generated-Files formatter does this through its contract evidence, as
TypeScript's does.

**Limitations of the import environment today:** at most 16 captured files; a
module path and every supplied import path need a dot in their first element,
so a module named `app` is refused and `math/big` cannot be named; a supplied
package cannot import the standard library and cannot lie under the
consumer's module path; dot and blank imports are unsupported.

**A reader of the unplaced output that would change.** One integration test
resolves and analyzes the generation snapshot itself as Go. Once the unplaced
output is no longer a Go file, it must use the snapshot of the placed Source.

### 4.8 Rust as it is

**The Rust target writes visibility.** Every generated struct, every type
alias and every field is marked public. The root is the vocabulary's `file`
node with `item` children, on which the grammar requires nothing. The first
fragment below is the end of the function that builds a declaration (`children`
holds its generic parameters and, for a struct, its fields). The second is
from the function that builds an expression from a native form:

```go
// transform/syntax/rust/target/go/provider.go:91-99, 107-119
		children = append(children, fields...)
		return b.Node(m.Node(s.KindStruct, m.Attributes("name", name, "visibility", "public"), origin), children...)
	}
	children = append(children, m.Child[s.Language]{Role: s.RoleType, Fragment: body})
	return b.Node(m.Node(s.KindTypeAlias, m.Attributes("name", name, "visibility", "public"), origin), children...)
}
func (b *Backend) Finish(declarations []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindFile, nil, ast.NodeRef{}), edges(s.RoleItem, declarations)...)
}
// …
	case model.FormPrimitive:
		if data[0] == model.PrimitiveString {
			return b.Path([]string{"std", "string", "String"}, nil, origin, true)
		}
		return b.Node(m.Node(s.KindPrimitiveType, m.Attributes("name", "bool"), origin))
	case model.FormSequence:
		return b.Path([]string{"std", "vec", "Vec"}, args, origin, true)
	case model.FormField:
		return b.Node(m.Node(s.KindField, m.Attributes("name", data[0], "visibility", "public"), origin), m.Child[s.Language]{Role: s.RoleType, Fragment: args[0]})
	case model.FormRecord:
		return b.Node(m.Node(s.KindStruct, nil, origin), edges(s.RoleField, args)...)
	case model.FormExternal:
		return b.Path(append(strings.Split(data[0], "::"), data[1]), args, origin, true)
```

`b.Path(…, true)` writes an absolute path such as `::std::string::String`, with
no link and no external use. Three sites do this: the string type, the
sequence type and a library-owned type. An absolute path stays valid wherever
the declaration is placed, so no placement breaks today, but the output records
no intent for these uses. The Rust native model has no notion of visibility:
its field form carries a name and a flag that marks the field optional. So
today field visibility is written by the target for every record alike, and no
mapping library states it.

**The vocabulary's visibility values.** The checker admits three values. The
attribute is optional on items (`module`, `struct`, `type-alias` and so on) and
on fields. A `file` and a `module` both hold items, so one file can contain
several modules:

```go
// check/syntax/rust/provider/go/check.go:20, 23-25, 34-35
var visibility = g.Set("private", "public", "crate")
// …
func named() map[string]g.Attribute {
	return g.Attributes(s.AttributeName, g.Attr(true, native.Declaration), s.AttributeVisibility, g.Attr(false, visibility))
}
// …
	s.KindFile:                  {Roles: []g.Role[s.Language]{g.Many(s.RoleItem, 0, items...)}},
	s.KindModule:                {Attributes: named(), Roles: []g.Role[s.Language]{g.Many(s.RoleItem, 0, items...)}},
```

The parser produces only two of the values, and no attribute at all for a
private item:

```go
// parse/syntax/rust/provider/go/parse.go:51-61, 190-192
func (p *parser) visibility() string {
	if !p.Eat("pub") {
		return ""
	}
	if !p.Eat("(") {
		return "public"
	}
	p.Want("crate")
	p.Want(")")
	return "crate"
}
// …
	if vis != "" {
		attrs = append(attrs, tr.Attrs("visibility", vis)...)
	}
```

The Rust Syntax formatter prints `pub` and `pub(crate)` and nothing else. Four
limitations of this vocabulary bear on the questions:

- `pub(super)` and `pub(in path)` have no syntax.
- A `use` item cannot carry a visibility, so a re-export cannot be written.
- A module exists only in its inline form, `mod name { … }`. The form
  `mod name;`, which refers to another file, has no syntax. The proposal's
  first Rust placement (Person in `src/people.rs`) needs it.
- A private item and an item on which nobody has decided look the same: both
  have no attribute.

**The Rust bridge reparses and compares.** Rust's compiler needs source text.
The bridge parses the captured file text again with bittype's parser
(`expected`) and requires the result to equal the captured tree (`tree`) node
for node, including attributes:

```go
// internal/rustsem/source.go:82-87
		for i, node := range expected.Nodes {
			actual := tree.Nodes[i]
			if node.Kind != actual.Kind || !reflect.DeepEqual(node.Attributes, actual.Attributes) || !reflect.DeepEqual(node.Children, actual.Children) {
				return nil, fmt.Errorf("%w: captured source/AST mismatch at %s:%d", semantic.ErrInvalidRequest, tree.ID, i+1)
			}
		}
```

A tree whose private item carried `visibility: private` would print without
`pub`, reparse without the attribute, and fail this comparison. A Files
formatter must therefore express "private" as an absent attribute.

**What the helper publishes about visibility: fields only.** The helper is a
Rust program linked against the compiler's internals (Rust 1.90.0). It exports
entities of the kinds `module`, `struct`, `enum`, `variant`, `field`, `alias`,
`trait`, `impl` and `parameter`, each with its compiler path as name (for
example `captured::models::Person`), its captured site and its owner.
Visibility occurs once in the helper, in the field list of a struct or
variant, and the Go side reads it as a boolean `Public` per field:

```rust
// internal/rustsem/helper/src/main.rs:506
                        fields.push(json!({"Name":f.name.to_string(),"Public":self.tcx.visibility(f.did).is_public(),"Entity":entity}));
```

There is no fact for the visibility of a struct, an alias or a module. The
graph does carry compiler-resolved `target` bindings at path segments, so a
comparison of bindings works as in TypeScript.

**The Rust semantic profile.** Its configuration is one fixed setting list:
edition 2024, the compilation target `x86_64-unknown-linux-gnu`, a crate named
`captured`, library checking, dependencies from the compiler's own sysroot
only. There is no counterpart to Go's import environment: nothing a consumer
supplies enters the configuration. The supported fragment is "one captured
source file with inline modules"; "Imports and multi-file crates, dependency
closures" are outside it. The profile "is qualified on Linux x86-64", in the
project's sense of "accepted as tested", and the helper is built only in the
project's continuous integration on Linux. A Rust Files formatter's check
therefore runs only there.

### 4.9 Names, and the two places that require them to be unique

**Where names come from.** The horos declaration recipe chooses every
declaration's native name. It is one recipe for all four languages. The
**entry module** is the horos module that the horos snapshot's entry selects
(`decl.Root` in the code). A declaration from any other module gets a prefix
built from the first 16 hexadecimal digits of its module's pin:

```go
// transform/syntax/horos/mapping/declarations/go/declarations.go:39-45
			name, err := naming.Exported(decl.Name)
			if err != nil {
				return zero, src.Unsupported(open.Body, err.Error())
			}
			if !decl.Root {
				name = "M" + decl.ModulePin[:16] + "_" + name
			}
```

`naming.Exported` capitalizes the first ASCII letter, turns a leading
underscore into `Type_` and `Self` into `Self_`, and refuses a name containing
`__`. So `Person` in an imported module becomes, for example,
`M3f9a0c2d1e4b5a67_Person` in every language. The documentation says of this
spelling: "This product spelling is not semantic identity", and describes the
present state as "Each mapping currently produces one declaration tree
containing the checked type closure." No document states the prefix's purpose.
Nothing we read forbids two modules with the same module name and different
pins in one output; we did not test it.

Because every name is capitalized, every generated Go declaration is exported,
whether or not the caller exposed it. In Go, visibility is the case of the
name.

**The flat-scope rule, as each target's program rule.** All four targets call
one shared function from `CheckProgram`. Its comment marks it as one profile's
rule that a target selects:

```go
// transform/syntax/shared/target/go/check.go:9-32
// CheckSingleScope qualifies this target's current single-scope construction
// profile. Other target organizations can supply their own program rules.
func CheckSingleScope[T any](ctx context.Context, p g.Program[T]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	names := map[string]bool{}
	for _, d := range p.Declarations() {
		if names[d.Description.Name] {
			return g.ErrUnsupported
		}
		names[d.Description.Name] = true
	}
	for _, d := range p.Declarations() {
		parameters := map[string]bool{}
		for _, parameter := range d.Description.Parameters {
			if names[parameter.Name] || parameters[parameter.Name] {
				return g.ErrUnsupported
			}
			parameters[parameter.Name] = true
		}
	}
	return nil
}
```

**The same rule again, inside the shared tree builder.** `Compile`, which all
four targets call from `Build`, repeats both refusals unconditionally, with a
different error:

```go
// transform/syntax/shared/target/go/compile.go:18, 21-26, 112-118
	names := map[string]bool{}
	// …
	for _, d := range declarations {
		if d.ID == 0 || int(d.ID) > len(declarations) || names[d.Description.Name] {
			return fail(g.ErrScope)
		}
		names[d.Description.Name] = true
	}
// …
	order := slices.Clone(declarations)
	slices.SortFunc(order, func(a, b g.PlannedDeclaration[T]) int {
		if a.Description.Name < b.Description.Name {
			return -1
		}
		return 1
	})
```

The loop that follows the first one refuses a type parameter named like any
declaration. This second copy is a shared helper's rule, not a kernel
invariant: the kernel seals whatever a target returns (section 4.3), and its
planner has no name rule. `Compile` identifies a plan symbol by declaration ID
and parameter index, so references do not depend on unique names. The order of
the output does. Declarations are emitted in name order, each after the
declarations its body uses, and the comparison never reports two names as
equal. A declaration that depends on itself is refused.

**TypeScript has the rule three times.** Its target calls `CheckSingleScope`,
`Compile` repeats it, and its Files stage checks each placed file (section
4.5). The accepted TypeScript correction did not touch the first two.

### 4.10 C++, briefly

The C++ target takes one namespace from its configuration and, in its
`Finish`, wraps every declaration in one namespace node. It also writes one
include directive per standard header its declarations need. C++ has no
resolution or analysis provider. The shared Files contract refuses a formatter
without compiler evidence (section 4.4), so a C++ generated-Files formatter
cannot be built yet, and its move is deferred with that trigger. It still
constrains the shape now:

- a namespace is a native scope that is independent of the file path, can be
  nested, and can occur several times in one header and across headers;
- a use in another namespace is spelled with a qualified name, not through an
  import;
- includes are file-level, order matters, and whether a forward declaration
  suffices depends on how a type is used.

## 5. Fixed points and their standing

Every statement this document treats as given has one of four statuses:

- **Owner rule.** Fixed.
- **Accepted decision of 7 November 2026.** Fixed, unless you see that it
  contradicts an owner rule or itself. If so, please say so.
- **Implemented today.** A fact of the present code. It can change now, and
  rule AC-09 gives it no weight in the design.
- **Our assumption or derivation.** Open. This includes everything that only
  the proposal of 8 November 2026 or a work item states.

| # | Statement | Status | Shown in |
| --- | --- | --- | --- |
| 1 | The rules AC-02, 03, 05, 07, 08, 09 and 10. In particular: a first move is final in structure, keeps the guarantees its destination's shared mechanism requires, and refuses what it cannot check | Owner rule | 2.1 |
| 2 | "Properly done", with no refactor later; a misplaced decision moves first, missing coverage is added afterwards | Owner rule | Question, 2.1 |
| 3 | A shared mechanism holds no contribution's meaning | Owner rule | 2.1 |
| 4 | Generation builds declaration bodies and binds each named use to its intended destination; it records what the caller exposed; a tree ID assigns no file | Accepted decision of 7 November 2026 | 2.2 |
| 5 | The consumer supplies one placement (output files and the native scopes they belong to) and the dependency environment; there is no second table of imports | Accepted decision of 7 November 2026 | 2.2 |
| 6 | The Files formatter chooses the spelling, qualification or alias of uses, constructs import and export syntax, realizes the requested visibility, and establishes that uses still resolve to their intended destinations | Accepted decision of 7 November 2026 | 2.2 |
| 7 | A spelling is not identity. The same generated declarations serve every placement. No private type is exposed and no same-named type substituted to make a placement work. A placement that cannot be expressed fails without a partial Source | Accepted decision of 7 November 2026 | 2.2 |
| 8 | Module visibility is Files'. A member's accessibility, "when it is part of the chosen type representation", belongs to the declaration | Accepted decision of 7 November 2026 | 2.2 |
| 9 | An output is a snapshot, exposed sites and external uses, sealed with the language's syntax checker. It may hold up to 256 trees; the shared tree builder makes one | Implemented today | 4.2, 4.3 |
| 10 | Unplaced output and placed source have one language descriptor and one syntax checker | Implemented today | 4.4 |
| 11 | The shared Files contract refuses a formatter without a printer, a check, or evidence that demands a semantic graph. It never reads a placement. The evidence requirement is fixed when the formatter is built. (Because of statement 1, a first slice must satisfy whatever this contract requires at the time.) | Implemented today | 4.4 |
| 12 | A Go file's package import path is the module path joined with its directory, and the files of one directory form one package. Merged on 12 November 2026 | Implemented today | 4.7 |
| 13 | No discovery: providers do no disk or network IO; Go and TypeScript dependencies are supplied as text; Go's standard library is not available; Rust sees only its pinned compiler's sysroot | Implemented today | 4.5, 4.7, 4.8 |
| 14 | The flat-scope rule (in each target and in the shared tree builder), the module-pin prefix, and the capitalization of every name | Implemented today | 4.9 |
| 15 | TypeScript's placement type, its key ending in `-1`, and one assignment per file | Implemented today | 4.5 |
| 16 | The Rust helper builds on Linux only, so Rust checks run only in continuous integration | Implemented today | 4.8 |
| 17 | The import lists of the package-role test. Today the shared Files contract package may not import the package of file IDs and paths, and a Files formatter's role does not include the generation contract; TypeScript's formatter has it by a per-package exception | Implemented today | 2.1 |
| 18 | A declaration's name does not depend on placement, and Files never renames a declaration. We read this from statements 4 and 6; they do not literally say it. TypeScript implements it | Our assumption or derivation | 4.5 |
| 19 | Unplaced output needs no native semantics of its own; only the placed result is analyzed. bittype's generation document says "Native intent is not an analyzed `Semantics[T]` [the type of an analysis result]: target analysis establishes facts about the actual constructed/organized program." It speaks of the plan. Today one test does analyze unplaced Go output | Our assumption or derivation | 4.7 |
| 20 | Collisions are checked per placed scope by Files. Only the proposal says "Check collisions"; TypeScript's formatter does it | Our assumption or derivation | 2.3, 4.5 |
| 21 | Exposure is two-valued: exposed or not | Our assumption or derivation | 4.2 |
| 22 | The six required placements, and "output file and native scope for each emitted declaration" as the shape of the input. Stated only by the proposal | Our assumption or derivation | 2.3 |
| 23 | What slice 1 admits, refuses and verifies. Stated only by the work items | Our assumption or derivation | 2.4 |
| 24 | The list of what counts as a later change | Our assumption or derivation | Question |

Out of scope: conversion code between native values and horos values, runtime
behaviour, C++ modules (as opposed to headers), and the choice of native
representations.

## 6. What we have considered

Everything in this section is ours and open. "Our leaning" marks what we would
do today. No option is ruled out because the present code would have to
change.

### 6.1 The root of unplaced output

In TypeScript and Rust a declaration sequence is the top production. In Go it
is a part of `SourceFile = PackageClause { ImportDecl } { TopLevelDecl }`.

| Option | What speaks for and against |
| --- | --- |
| **A. A second root kind** in the Go vocabulary for a declaration sequence, with declaration children only | **Our leaning.** A checked `file` still always has a package. The bridge already refuses a non-`file` root, so resolution and analysis refuse unplaced output. Go Files requires this root and refuses a `file` from generation, as TypeScript refuses `exported`. Against: only Go then marks "unplaced" in its root |
| B. `file` with an optional package | One node kind stands for two things, and every reader of a `file` loses the guarantee that it is a Go file. The syntax checker has no notion of "before placement", so the attribute would be optional always. We think this repeats the "double duty" of section 4.5 in mirror image |
| C. No container: one tree per declaration | Not our leaning; open. The output contract already admits many trees, and references are sites, so only the shared tree builder would change. It makes "no file decision" structural in every language. But each tree's root is then a declaration, so the Go checker must still admit a second root kind; and a snapshot needs one entry node, which would be arbitrary |
| D. A distinct language profile, and so a distinct type, for unplaced output | Not our leaning; open. A fragment and a compilation unit would be different types, and no reader of files could receive a fragment by mistake. The Files stage would then map one language profile to another, and each language would need two checkers. We have not ruled this shape out |

Two riders on A. The kind should be named for the syntax (Go's sequence of
top-level declarations), not for generation, because the vocabulary package
holds no generation concept. And the Go Syntax formatter must print the new
root, as declarations without a clause; otherwise the statement in bittype's
generation document that "Syntax formatting and printing can also operate
directly on individual output trees" stops being true for Go.

Under A the three languages mark unplaced output differently. Go marks it in
the root. TypeScript marks it by absence, and its formatter can refuse a
generated `exported` or import. Rust cannot mark it at all: an undecided item
and a private item both have no attribute (section 4.8), so a Rust formatter
can refuse a generated `pub`, `use` or `mod` but cannot tell "not yet decided"
from "decided: private". In C++ the unplaced set is a part of a translation
unit, and placed declarations sit inside namespace nodes.

### 6.2 Who owns name collisions

| Option | What the consumer is guaranteed | Who owns what |
| --- | --- | --- |
| U1. Both: the flat rule stays in every target and in the shared tree builder; Files adds the rule per scope | Every output can be put into one scope, and a name identifies a declaration within one output. But an output with `Item` in two horos modules can never be placed as `pages.Item` and `people.Item`, although Go allows it | The target keeps a rule about scopes it no longer decides |
| **U2. Per placed scope only**: Files owns collisions; no target that has a Files stage and no shared helper requires uniqueness | Any placement the language permits is accepted. A collision fails at the stage that knows the scope | Files: scope rules. Recipe: names. Tree builder: identity only |
| U3. A target-selected profile: each target states its own program rule, as the comment on the shared function already says; the shared tree builder holds none | As U1 for a target that selects the flat rule, as U2 for one that does not | Each target; the shared helper holds no placement's policy |

**Our leaning is U2**, with C++ keeping the flat rule in its own target until
it has a Files stage (which is U3's mechanism). Reasons: "collision" is the
native scope rule applied to a placement, and only Files knows the scopes
(AC-02); and we found no justification of its own for a set-wide rule once a
Files stage exists (AC-09).

What speaks against U2:

- The accepted TypeScript correction kept the flat rule, and no document says
  it is wrong.
- By our own criterion U1 can become U2 later without a counted change,
  because it only admits what was refused. So U2 is not forced now.
- U2 removes the property that a name identifies a declaration in one output.
  Consumers rely on it today to find unexposed sites (section 4.2).
- The tree builder's output order needs a rule for equal names (section 4.9).

### 6.3 Two naming questions we cannot answer alone

**The prefix.** Our reading is that `M<16 hex digits>_Name` exists to make one
flat scope work. No document says so, and it may serve other purposes we do
not see: names that stay stable when an output grows, or two versions of one
module in one output. With packages and modules a consumer would write
`people.M3f9a0c2d1e4b5a67_Person`. Since outputs usually span several modules
(section 1.4), this is the common case.

**Go exposure.** The recipe capitalizes every name, so in Go a declaration the
caller did not expose is exported anyway. Go Files can check exposure and
cannot realize it without changing the name's case, which statement 18 of
section 5 forbids. So after Go slice 1 one visibility decision would still be
made during generation. Whether that is a misplaced decision, and who should
set the case of a Go name, is open.

Both must be settled before the first release fixes the names.

### 6.4 What one placement assignment carries

For Go, beyond TypeScript's file ID, path and declaration sites:

| Candidate | Carry it? | Reason, and what speaks against |
| --- | --- | --- |
| Directory | No | It is the directory of the path |
| Package import path | No (leaning) | The module path (environment) and the directory determine it; stating it again would be a second source for one fact. Against: the proposal calls an import path and an output filename "different identities"; output paths are relative to wherever the caller writes them, which need not be the module root; and the rule that joins them is implemented today, not accepted |
| Package name (the clause) | Yes | Not derivable: a package's name is not its last path element by rule (`main`, a path ending in `v3`). It is exactly the decision slice 1 moves. Note that the name is not the package's identity: two directories can both hold a package named `pages` |
| Imports, aliases, qualifiers | No | Files decides them after placement |
| Exportedness | No | In Go it is the case of the name (6.3) |

The package name could be stated per assignment (each file states its package,
as each Go file does; an inconsistent placement is representable and must
fail) or once per directory (inconsistency is unrepresentable, but a directory
without a file becomes a new failure). We lean to per assignment.

For Rust an assignment would give a file and a module path relative to the
crate root. Several assignments may name one file, because two modules can
occupy one file. The crate name and dependency crates are environment, as the
proposal says of Cargo package selection. `mod` items and their visibility are
not carried; Files constructs them. For C++ the scope would be a namespace
path, independent of the file path.

### 6.5 One shared placement type, or one per language

| Option | Shape |
| --- | --- |
| **S1. One shared type and its admission in the shared Files contract**, generic in the scope. Each language supplies its scope type, its setting key and the rules relating scope and path | One shape for the consumer in every language; the assignment rules have one implementation |
| S2. Per-language types; only the admission function is shared | The rules are shared; the consumer learns several shapes that differ in one field |
| S3. Per-language types and admission, as today | Nothing shared; Go and Rust would implement the rules of sections 4.5 and 4.6 a third and fourth time |

A sketch of S1, not code:

```text
Placement[Scope]
    Entry        file ID                        the entry file
    Assignments  list of
        File          file ID
        Path          relative path
        Scope         Scope                     TypeScript: none; Go: package name;
                                                Rust: module path; C++: namespace path
        Declarations  list of sites
```

**Our leaning is S1.** The reasons are AC-02 (an invariant of the Files
boundary should have one owner) and AC-03 (adding a formatter should not copy
it). S1 changes a shared contract and reshapes TypeScript's `Organization`, so
it needs AC-08's justification, and under AC-10 it must be decided before
slice 1 fixes a Go type.

What speaks against S1:

- By our own failure criterion it is the riskiest choice. C++ has no semantic
  tooling yet, and a wrong guess about its scopes would change the type every
  language's consumers write. Per-language types confine such a mistake.
- "The same rules" is overstated. S1 needs several assignments per file, which
  both existing implementations forbid.
- S1 is generic in the scope only. The horos formatter assigns module pins, not
  sites, so it would stay separate unless the type were generic in the item
  too.
- Pairing a file with a scope in one assignment is itself an assumption. A Go
  package spans files and a C++ namespace is independent of files; two
  independent assignments (site to file, site to scope) are conceivable.

Open inside S1:

- **Flat or nested.** A flat list in which a file may repeat, or files that
  each hold a list of scoped groups.
- **Order.** The order of assignments and of sites is declaration order in
  TypeScript's output. In C++ order is semantically relevant. A placement that
  was accepted must not change meaning later.
- **The entry file.** No native compiler bridge reads it (section 1.2). A
  Go placement over several packages has no obvious one.
- **What the shared part may know.** It can check assignment, sites, files and
  the entry. It cannot know that a Go directory has one package name, or that
  Rust module paths form a tree. It also travels as opaque bytes today, and
  the shared contract package may not import the package of file IDs.

### 6.6 The evidence profile of the first slices

**Go.** Section 4.7 shows that the base profile cannot see which package a
declaration was placed into, while the import environment shows it as package
identity. The Go work item, written before the import environment existed,
says slice 1 "does not wait" for it.

**Our leaning:** Go slice 1 builds on the import environment and requires it
when the formatter is built, as TypeScript does. Otherwise a later slice
changes the evidence requirement, the constructor and the meaning of paths.
Its check would be: (1) a fresh receipt admitted against
`EvidenceFor(configuration)`; (2) every relocated use binds to the entity
declared at the relocated site; (3) every relocated declaration is owned by
the package entity of the expected import path; (4) the package name. We have
not run an import environment without supplied packages end to end; the code
bounds the number of supplied packages only from above.

For (4), the package name, three options:

| Option | What it establishes |
| --- | --- |
| P1. Tree equality with the recomputed layout, plus the compiler's rule that the files of a directory agree; stated as a limitation | That Files wrote the name the consumer stated. A later slice covers the name indirectly, when an import without an alias binds through it |
| P2. The Go bridge publishes the package name as a fact | An independent fact, as for TypeScript exports. The compiler already knows the name |
| P3. No check | Nothing |

We leaned to P1 for the package name and, in 6.8, against the equivalent
option for Rust item visibility. Our reason was hypothesis 4 of section 6.9:
the package name is copied from the consumer's input, while `pub` is derived
by Files and its effect is decided by the compiler. We are not sure this
holds. In a slice with one module, a declared `pub` is as syntactic as a
package clause. The different treatment is therefore unexplained, and part of
question 4.

**Rust.** There is nothing to choose between: the profile has one fixed
configuration and no consumer-supplied environment (section 4.8). A Rust
slice 1 formatter would be built without one. A later slice with dependency
crates would add one, and so change the evidence requirement and the
constructor, which we count as a failure. We have no leaning on how the first
slice should be shaped to avoid that.

### 6.7 Intent for predeclared and standard types

Go's `string` and `bool` are written with no recorded destination, and so are
Rust's three absolute paths. The proposal requires external symbols to be
represented, "not guessed from a name" (section 2.3). TypeScript records a
global type as an external use whose symbol has an empty module, and checks
it against a builtin entity.

**Our leaning for Go:** report a predeclared type as an external use with a
Go-issued symbol whose package is empty. Slice 1 accepts those, checks them
against `builtin` entities, and refuses a symbol with a package. Otherwise the
work item's rule "an output with an external use fails" and the binding check
contradict each other, and a later slice must change how the Go target
reports uses. Rust's three paths stay as they are until the deferred instance
of section 2.4 is taken up.

### 6.8 Rust visibility

**Fields: a decided rule, and our application of it.** The rule is statement 8
of section 5. The Rust work item extended the move to fields; we withdrew
that.

| Option | Guarantee | Owner |
| --- | --- | --- |
| **F1. Field visibility stays in the declaration** | A record's fields are accessible wherever the type is; an opaque type would be another mapping | Generation |
| F2. Files writes it from exposure and placement | The same result today, but Files must assume every struct is a transparent record, and the output has no carrier for another intent | Files, without intent |

We follow F1: `pub` on a field never widens exposure beyond the struct's own.
The application is ours, because the rule is conditional ("when it is part of
the chosen type representation") and today no mapping states field visibility:
the Rust target writes it for every record (section 4.8). We read that as a
representation choice made in the wrong component of generation, not as a
Files decision. The record form of the Rust native model should state it.

**Items.** Rust Files decides: an exposed declaration is `pub`, every other
one has no attribute. Two ways to check it:

1. add a compiler-established visibility fact for items to the helper and the
   Rust vocabulary, as TypeScript has `named-export`;
2. check by tree equality only and record the limitation.

**Our leaning is 1.** We are unsure which fact is the right one: the item's
declared visibility, its effective visibility from the crate root, or its
reachability from a using module.

### 6.9 Hypotheses we want attacked

Each is marked with its status from section 5.

1. "A declaration's name never depends on placement; Files chooses only how a
   use is spelled." (Our derivation, statement 18. Go exposure in 6.3 strains
   it.)
2. "Collision is a property of a placed scope, so generation needs no
   uniqueness rule once a language has a Files stage." (Our derivation.)
3. "A placement assignment is exactly a file, a native scope and declaration
   sites. Everything else is derived by Files or supplied by the
   environment." (Our derivation; the accepted decision says only "output
   files and the native scopes they belong to".)
4. "Equality with a recomputed tree is sound evidence only for what the
   formatter copies from the consumer's input. It is never evidence for a
   property whose meaning the compiler decides: binding, visibility, package
   identity." (Our assumption.)
5. "The unplaced set needs no semantics of its own. Only the placed result is
   analyzed." (Our derivation, statement 19.)

### 6.10 Surprises we already suspect

These are known to us, so they are not what question 5 asks for:

- **Hand-written code beside generated code.** A Go package's scope spans
  every file of its directory. The check sees only the captured files, and a
  supplied dependency cannot lie under the consumer's module path. If a
  consumer writes generated files into a directory that also holds its own
  files, collisions, the agreement on the package name and a declaration that
  shadows `string` are all outside the check.
- **Cycles and order.** Go package cycles between placed packages (the bridge
  rejects them); in C++, declaration order and whether a type is complete
  where it is used.
- **Visibility narrower than public, and re-exports.** The Rust vocabulary has
  neither (section 4.8).
- **Finding sites.** A consumer finds an unexposed declaration by its name
  (section 4.2).
- **The entry file** in layouts with several packages or crates.
- **Shared parts shaped by the first implementation.** The package-role list
  for Files formatters lacks what every generated-Files formatter needs
  (statement 17).

## 7. Questions for the expert

In section 5, owner rules are fixed, and the accepted decision is fixed unless
you see a contradiction. Everything else is open, including what the code does
today. Please name the option first and then the reason. Patterns from
comparable systems (multi-target code generators, compilers with separate
layout and emission phases, build systems that place generated code into
packages or modules) are welcome where they change the answer, including what
failed.

1. **The placement input.** A consumer assigns one generated set of
   declarations to files and native scopes: the file itself in TypeScript, a
   named package per directory in Go, a module path in Rust (several per
   file), a namespace in C++ (several per header, independent of the path).
   (a) How would you divide this input between a part all languages share and
   a part each language owns? (b) What must one placement assignment contain,
   and what should the order of assignments and the entry file mean? (c) How
   would you test the division before relying on it, for example with
   placements beyond the six of section 2.3? We lean to one shared type whose
   assignments pair a file, a scope and declaration sites, with a scope type
   per language (6.5, option S1; 6.4 for the fields); 6.5 also lists what
   speaks against it. Your answer decides S1, S2 or S3 and the content and
   order of an assignment, before Go slice 1 fixes a type.

2. **Names, collisions and exposure.** Generation chooses each declaration's
   name. Today it also requires all names to be unique in one flat scope, in
   two places, and prefixes names from all but one module with part of a
   module hash (4.9). (a) Once the consumer decides scopes, what would you
   have generation guarantee about names, and what would you give to the
   per-language Files stage? (b) What should become of the prefix? (c) In Go,
   where the case of a name is its visibility, who should set it? We lean to
   Files alone checking collisions, per placed scope (6.2, option U2). The
   accepted TypeScript correction kept the flat rule (U1), and today a
   consumer finds an unexposed declaration by its name (4.2). Your answer
   decides U1, U2 or U3, the prefix, and the owner of a Go name's case, before
   the first release fixes names.

3. **Unplaced output as syntax.** Generated output is a syntax-checked tree of
   the target language, and today one checker serves both it and the placed
   source (4.2, 4.4). A bare declaration sequence is a legal root in
   TypeScript and Rust. It is not one in Go, where a file requires a package
   clause, nor in C++, where placed declarations sit inside namespaces.
   (a) How do systems with typed syntax trees represent declarations that are
   not yet placed? (b) What should each reader (syntax checker, printer,
   compiler bridge) be able to rely on at a placed root and at an unplaced
   one? (c) Should every language mark unplaced output in its tree, given that
   Rust today cannot tell an undecided item from a private one? We lean to a
   second Go root kind for a declaration sequence (6.1, option A). Options B,
   C and D are open, including a distinct type for fragments. Your answer
   decides the Go root for slice 1 and whether Rust slice 1 needs a marker.

4. **Evidence for a moved decision.** Every move needs a check that the placed
   program means what generation intended. A check can rest on reconstruction
   from the consumer's input, on compiler acceptance, or on a fact the
   compiler bridge publishes; TypeScript uses all three (4.5). (a) Which would
   you require for each of these: the Go package name; Go package identity;
   Rust item visibility; and the destination of Go's predeclared `string` and
   `bool`, where the compiler publishes a fact but generation records no
   intent (4.7, 6.7)? (b) Where a fact is needed, which one exactly? For
   visibility: declared, effective from the crate root, or reachable from a
   using module? (c) What must a first slice that admits one package or one
   module already require and already refuse, so that its evidence
   requirement and constructor never change, given that Rust has no
   consumer-supplied environment today? We lean to the Go import environment
   from slice 1, tree equality with a stated limitation for the package name,
   and a new Rust visibility fact (6.6, 6.8); we cannot explain why we treat
   the last two differently. Your answer fixes each slice's evidence
   requirement, constructor and refusals.

5. **Open.** What should we know that these questions do not ask? In your
   experience of separating generation from placement, what is discovered late
   and then forces a change to the first delivered interface? Section 6.10
   lists the surprises we already suspect, so what lies outside that list is
   the most valuable to us. We will use your answer to add fields, refusals or
   checks to the two first slices before they merge.
