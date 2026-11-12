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
  about files, packages, modules, namespaces, imports or visibility.
- A later stage, called **Files**, takes the generated declarations, one
  **placement** supplied by the consumer (which declaration goes into which
  file and which native scope) and the consumer's dependency environment. It
  writes the package clauses, imports, aliases and visibility modifiers the
  placement needs. It then has the target language's own compiler resolve the
  placed program, and compares what the compiler found with what generation
  intended.

This split is implemented for TypeScript only. Two targets still make a file
decision during generation:

- **Go** writes the package clause, taken from target configuration.
- **Rust** writes `pub` on every generated struct, type alias and field.

Two small changes, called **slice 1** for Go and slice 1 for Rust, are about to
move these two decisions into new Go and Rust Files stages. Each slice admits
only the simplest placement (one Go package without imports, one Rust module)
and refuses everything else. Later slices add external types, imports, several
packages and several modules. C++ follows when it has compiler-backed semantic
tooling.

**What we need advice on.** The project's rule for such a first move is that it
establishes the final structure: its inputs and identities already have their
intended shape, and later slices only add (section 2.1, rule AC-10). The
owner's direction of 9 November 2026 is: "I want this properly done, so that we
don't accumulate technical debt and have to refactor later." Nothing in bittype
is released yet, so every interface can still change now, and none of it should
have to change later.

**What shape should the placement input, the unplaced generated output and the
placement check have, so that Go, Rust and later C++ are added without changing
what TypeScript, Go slice 1 and Rust slice 1 delivered?**

We count as a later change, and therefore as a failure of the first move:

- a change to the type a consumer writes to state a placement, or to the
  meaning of a placement that was already accepted;
- a change to what generated output is before placement (its root, its names,
  what it records about intent);
- a change to the evidence a Files stage requires, or to the constructor a
  consumer calls to build one;
- a change to a name that appears in consumers' code.

Admitting a placement that is explicitly refused today is additive and not a
change.

**The tensions**, each developed in section 6:

- **Scope differs per language.** A TypeScript module's scope is its file. A Go
  package is a directory of files and has a name that the path does not
  determine. A Rust file can hold several modules. A C++ header can hold
  several namespaces, and a namespace is independent of the path.
- **Unplaced Go is not a Go file.** Generated output is a checked syntax tree
  of the target language. Go's source file production requires a package
  clause.
- **Names and scopes.** Today all generated names must be unique in one flat
  scope, and a naming policy exists only to make that possible. Once the
  consumer decides scopes, it is unclear who owns collisions.
- **Evidence that cannot see the moved decision.** For the Go package name and
  for Rust item visibility, the compiler-established facts bittype publishes
  contain nothing to compare against.

Sections 1 and 2 give the system, the rules and the two governing documents.
Section 3 shows the code: the shared contracts, TypeScript as the working
model, then Go, Rust, name uniqueness and C++ as they are. Section 4 follows
one generated set through placement. Section 5 lists what is fixed. Section 6
gives the options and our leanings, marked as decided by an accepted document
or derived by us. Section 7 holds five questions.

Quotations from project documents are verbatim, including where they
capitalize project names. In code excerpts a line `// …` marks lines we left
out; everything else is exact, with its file path and line numbers at the code
baseline.

## 1. System overview

### 1.1 The components

| Component | What it is | Language(s) |
| --- | --- | --- |
| **ontos** | The value model under horos. A value is an atom (a byte string) or a tuple (an ordered sequence of values). It plays no further part in this document. | Go, Rust, TypeScript |
| **horos** | A type language whose types are ontos values. It has a declaration language with **modules**, generic declarations and imports. A module is identified by a **pin**, the SHA-256 digest of its canonical content, written as 64 hexadecimal digits. An import names a module by its pin. | Go, Rust, TypeScript |
| **bittype** | A Go library. It parses and analyzes horos source, and it generates native type declarations for TypeScript, Go, Rust and C++. Its native syntax covers type declarations only, not executable code. | Go; emits TypeScript, Go, Rust, C++ |
| **logos** | A separate proof language used elsewhere in the project. It appears here only as a name in a quoted sentence. | — |
| **Consumers** | Applications that link bittype as a library, register the parts they want, and call generation and Files. The first consumers are the build steps of four network services of the surrounding platform. | Go |

A horos example used throughout, with the declarations a consumer asks for:

```text
Page<T> = { items: list<T>; next?: text; }       a generic record
Person  = { name: text; }                        a record
Cells   = Page<cells.Cell<Person>>               an alias; Cell is owned by a library
```

`Cell` stands for a type that a separate library owns in every target language
(`Cell` exported by the npm package `@bitspark/cell`, `cell.Cell[T]` in a Go
package, `Cell<T>` in a Rust crate). Generation must refer to it without
declaring it.

### 1.2 Terms used throughout

| Term | Meaning in this document |
| --- | --- |
| **Target, native** | A **target** is one of the four output languages. **Native** means "of the target language": a native type, a native scope, the native compiler. |
| **Syntax tree (AST)** | bittype's own tree representation of a program in some language. One generic carrier serves every language: a flat array of nodes, each with a kind, text attributes, child edges with role names, and reference links (section 3.1). Each language supplies a **vocabulary**: its node kinds, roles and attribute names. |
| **Site** | The address of one node: a tree ID and a node ID. Declarations and uses are identified by site, never by spelling. |
| **Parse, check, format, print** | For one tree: **parse** turns text into a tree, **check** admits a tree against the vocabulary's grammar rules (the **syntax check**; it resolves no names), **format** turns a tree into concrete syntax, **print** turns concrete syntax into text. The project calls this the **Syntax** dimension. |
| **Files** | The second dimension: a set of trees, each linked to a file ID and a relative path. A **Files formatter** produces such a set for a chosen organization. "Files" with a capital letter always means this stage or dimension. |
| **Source** | The published result of a Files stage: the units (file ID, path, concrete syntax), together with the compiler-backed evidence for exactly those units. A Files printer turns a Source into paths and text. Nothing is written to disk by bittype. |
| **Snapshot, capture** | An immutable, syntax-checked set of trees, optionally with the files they belong to. To **capture** is to create one. Every later stage refers to a snapshot, never to mutable trees. |
| **Resolution and analysis providers** | Per-language components that run the language's real compiler front end over a snapshot (Go's `go/types` library, the TypeScript compiler API, a helper linked against `rustc`) and publish what it established. They are called "providers" because the consumer supplies them explicitly; nothing is discovered or registered automatically. |
| **Semantic graph, fact, binding** | What a provider publishes: **entities** (declarations, packages, modules), **bindings** (this use site resolves to that entity) and **facts** (provider-owned relations such as "this module exports that declaration"). Each language owns its names in the graph. |
| **Receipt, evidence** | The result of one resolution or analysis, bound to one snapshot and one configuration. A stage that needs compiler-established facts states an **evidence requirement** (operation, profile, configuration, required reports), and a receipt is **admitted** against it. A receipt for another snapshot or configuration is refused. |
| **Compiler bridge** | The private code that converts bittype's tree into the compiler's own input and exports facts back: `internal/golangsem` for Go, a Rust executable (the **helper**) for Rust. |
| **Profile** | A named, bounded subset of a language or a policy. A syntax or semantic profile says which constructs are supported. A **representation profile** says how horos forms become native types. |
| **Generation kernel** | bittype's language-agnostic code that plans and runs generation. The documents also call the shared, contribution-independent parts "the kernel" or "shared mechanisms". |
| **Mapping library** | A contribution that says how a horos form or declaration is represented natively (text becomes `string`, a list becomes a slice). The **declaration recipe** is the mapping part that turns a named horos declaration into a native declaration and chooses its name. |
| **Native target** | The per-language contribution that constructs the syntax tree from the closed plan the kernel hands it. |
| **Exposed site** | The site of a declaration the caller asked to be public. The documents also say "public root" and "publication". |
| **External use** | A use site that refers to a type outside the generated set (the library-owned `Cell`), recorded together with an opaque, target-issued **symbol** that identifies the destination. |
| **Placement** | The consumer's assignment of generated declarations to output files and native scopes. |
| **Dependency environment** | What the consumer says about the world the output is compiled in: the module root, the dependency packages and their declarations, compiler settings. |
| **Slice** | One delivered change. Slice 1 of an issue moves a misplaced decision; later slices add coverage. |
| **ACX-nnn** | A numbered architecture case record in bittype's repository: a found violation of the rules in section 2.1, with its correction and limits. ACX-004 is "generation emits file syntax before placement". |

### 1.3 The stages, in one picture

```text
horos files ──parse, resolve, analyze──▶ horos snapshot + receipt
                                              │
                    generation kernel: plan   │  mapping libraries choose representations
                                              ▼
                              closed native program (intent, no syntax yet)
                                              │
                    native target: construct  │  one syntax tree, references linked to their targets
                    syntax check, seal        ▼
        generation.Output = snapshot of the unplaced tree + exposed sites + external uses
                                              │
      consumer: placement (per request)       │   consumer: dependency environment (per formatter)
                                              ▼
              generated-Files formatter of the target language
                 1. copy declarations into one tree per output file
                 2. write the file-level syntax the placement needs
                 3. capture the new trees with their paths
                 4. run the native resolution and analysis providers
                 5. compare the compiler's bindings and facts with the intent
                                              │
                                              ▼
              Source (units + receipt)  ──print──▶ paths and text
```

Generation runs once. The same output can be handed to the Files stage several
times with different placements.

### 1.4 Scale and direction

- **Release state.** Nothing is released. The first release is a tracked work
  item that follows this work. Interfaces, names and bytes can all change now.
- **Targets.** TypeScript first, Go next. Rust and C++ types are generated
  today; their Files stages must remain possible without reshaping the shared
  parts.
- **Sizes.** A generation output holds tens to a few hundred declarations. The
  Files dimension defaults to at most 256 files. The Go semantic profile
  accepts at most 16 captured files per analysis. The Rust semantic profile
  accepts one captured source file, which may contain inline modules.
- **People.** One human owner decides for all repositories involved. Each
  repository is maintained by AI coding-agent sessions, several working in
  parallel. Clear ownership and additive composition are therefore hard
  requirements.
- **Existing formatters.** The directory of Files formatters holds two: horos
  (it places whole horos modules into files) and TypeScript. TypeScript's is
  the only one that consumes generation output.

## 2. Architectural context

### 2.1 Standing architecture rules (AC = architectural correctness)

The owner has made architecture part of correctness: a design that works and
passes its tests is still wrong if a responsibility sits with the wrong owner.
The rules are numbered. The questions rest on seven of them. Each is given in
our words, followed by its requirement text, quoted.

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

**Package roles.** The rules are also applied to Go packages. Each package has
one role, and a test forbids imports the role does not allow (section 3.11):

- a **shared mechanism** imports no contribution and no vocabulary of one;
- a **vocabulary** (names two parties must agree on) imports no contribution;
- a **contribution** (one thing a consumer registers) is one package, and
  imports shared contracts, the vocabularies of its own language and its own
  private parts;
- only a **composition** (the consumer's wiring) imports a contribution.

The owner's request of 11 November 2026 names the fault this guards against:
parts that call one another directly "instead of invoking the kernel". The
rule text adds two faults of grouping: parts that do not need each other
sharing a package, and "a policy next to a shared mechanism", where "a shared
package also holds one case's rule".

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

Further sentences of the same decision that the questions rest on:

> A spelling alone is insufficient: two distinct types named `Person` must not
> become interchangeable when placement changes.

> Both placements use the same generated declarations, bindings and publication
> intent. Files must preserve those bindings and the requested public surface;
> moving a declaration cannot authorize exposing a private type or substituting
> another type of the same name. If it cannot express the placement with the
> supplied dependencies and visibility, it must fail without returning a
> partial `Source`. Native scope matters: two files in one Go package need no
> import between them, while two Rust modules can occupy one file. A file
> boundary alone is never an import rule.

> The distinction is about module visibility. A member's accessibility, when it
> is part of the chosen type representation, still belongs to the declaration.

> This is the boundary for **generating declarations**. Native ASTs can still
> represent imports and exports when parsing existing source or constructing
> organized units.

Three other accepted documents describe the implemented TypeScript path in the
same terms. The generation planner document says: "A sealed generation.Output
holds checked declaration ASTs, requested exposed sites and external use
destinations. Files receives that output, one placement and the selected native
dependencies, builds module syntax and checks its bindings before publishing
Source." And: "The planned native program is intent, not a fabricated
semantic-analysis result. Only analysis of the realized target program
establishes native semantic facts."

### 2.3 The proposal of 8 November 2026

`docs/native-file-placement.md` walks the Page, Person and Cells example
through Go, Rust and C++. **It is a proposal and not an accepted decision.**
It opens: "**Proposal, 8 November 2026.** These examples specify the intended
Files boundary." It says of its own tables: "The concrete Go/Rust/C++ settings
and checks still need specification and implementation. The following tables
are proposed inputs/output, not shipped Go structs or APIs." The architecture
case and both work items nevertheless treat it as the specification of the
boundary. Where only this document speaks, we mark the point as open.

Its proposed inputs, quoted:

| Input | Proposed meaning |
| --- | --- |
| Owned declaration placement | "Output file and native scope for each emitted declaration/helper; entries identify declarations, not just spellings" |
| Public roots | "Generation's exposed declaration sites; placement cannot make an inaccessible dependency public merely to obtain a valid file" |
| External destination | "Library-selected public package/crate/header symbol and mapped arguments" |
| Dependency environment | "Consumer-selected module/crate/include roots, actual dependency artifacts and relevant native configuration" |

Its required placements, condensed from its tables. Each must be produced from
one generated output without regenerating:

| Language | Placement | What Files must write in the unit containing `Cells` |
| --- | --- | --- |
| Go | Page and Cells in `pages/page.go`, package `example.org/app/pages`; Person in `people/person.go`, package `example.org/app/people` | Imports of `example.org/app/people` and of the Cell package; the use `Page[cell.Cell[people.Person]]` |
| Go | Person in `pages/person.go`, the same package as Page | No import of Person; the use `Page[cell.Cell[Person]]` |
| Rust | Page and Cells in `src/pages.rs`; Person in module `people` at `src/people.rs` | `use crate::people::Person;` and `use bitspark_cell::Cell;` |
| Rust | "Person's bytes move into an inline `mod people` in the same physical file" | "Person still belongs to another module; use a suitable module path/import and satisfy its visibility" |
| C++ | Page and Cells in namespace `app::pages`; Person in `model/person.hpp`, namespace `app::model` | Includes of the Person header and the Cell header; the qualified name `app::model::Person` |
| C++ | "Person moves into `pages/page.hpp`, retaining namespace `app::model`" | No include of Person; the qualified name stays |

Further sentences of the proposal that the questions rest on:

> A package import path and an output filename are different identities; the
> consumer's Go module layout relates them.

> Moving a declaration across packages also changes its native package identity
> and can introduce a dependency cycle. Files must check visibility, identity
> correspondence and package cycles; it cannot silently export an unexported
> type or treat same-package movement as a package import. Package aliases must
> avoid names in the use's lexical scope.

> Files must construct `mod` declarations or inline modules consistently with
> the placement, and realize public-root visibility under the crate/module
> rules. Two files may contribute through one module arrangement; one file may
> contain several modules.

> Cargo package selection/renaming belongs to the dependency environment; it is
> not inferred from a source type's name.

> The Files profile owns include directives, guards, declaration order and the
> chosen relative/public-header policy. [...] A forward declaration is not
> automatically enough: required completeness depends on the selected native
> use.

> Check collisions, visibility, missing external destinations, dependency
> cycles/completeness and public roots. The shared contract associates sites
> and evidence; language integrations establish native scope and destination
> meaning. [...] Successful printing or compilation alone does not establish
> that every intended use bound to the right declaration.

### 2.4 The order of correction, and the two slices

The architecture case ACX-004 records four misplaced decisions and their order
(11 November 2026):

| Instance | First step | State |
| --- | --- | --- |
| Go package clause | Generation decides no package. Go Files writes the clause for one package and no import (slice 1 of the Go work item) | Can start |
| Rust visibility | Generation chooses no visibility. Rust Files decides it for one module (slice 1 of the Rust work item) | Can start |
| Rust external paths | The target reports an external use and Rust Files decides the import | Deferred until Rust semantics cover a second crate |
| C++ includes and namespace | Generation emits neither. C++ Files decides both | Deferred until C++ has resolution and analysis providers |

The case states the rule for the first step: it "moves the misplaced decision to
the language's generated-Files formatter, for the placements that existing
semantic evidence can check, and refuses every other placement. The second step
adds coverage and changes nothing the first delivered. The first step cannot
skip its check".

Both work items require of slice 1: "The placement setting has its final
shape." Before coding started, we read the accepted documents against the code
to find what they decide and what they leave open. The open parts are the
subject of this consultation. The plan that orders the work says of itself that
it "decides no design".
