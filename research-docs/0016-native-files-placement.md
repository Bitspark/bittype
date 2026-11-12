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
parts that call one another directly "instead of invoking the kernel". It also
names a fault of grouping: parts that "neither conceptually nor practically
depend on each other" sharing one package. bittype's own record of its package
boundaries lists, among the signs of wrong grouping: "**A policy next to a
shared mechanism.** A shared package also holds one case's rule."

For AC-09 the quoted text is the owner's standing direction of 8 November 2026,
which the rule adopts as its requirement.

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

The generation planner document, which describes the implemented TypeScript
path, uses the same terms: "A sealed generation.Output
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

## 3. Relevant code

All code is Go unless marked, at the code baseline. In the excerpts, short
import aliases follow one pattern: `s` is the syntax vocabulary of the language
at hand, `n` its semantic vocabulary, `g` the generation contract package
(inside the grammar checkers, `g` is a small grammar helper), `m` the shared
target constructor, `f` or `fm` the format contract package.

### 3.1 Carriers

One tree carrier serves every language. `L` is a marker type owned by the
language; it prevents mixing trees of different languages at compile time.

```go
// trees/ast/go/origin.go:7-10, ast.go:31-35, node.go:13-16, 20-23, 28-31, 38-44
type NodeRef struct {
	Tree TreeID
	Node NodeID
}
type AST[L any] struct {
	ID    TreeID
	Root  NodeID
	Nodes []Node[L]
}
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

`Children` are containment. `References` are resolved, non-containment links: a
type reference that names a generated declaration carries a link to that
declaration's node. An unresolved name is only an attribute. `Origin` records
where a node came from and plays no part here.

Files, settings and the semantic graph:

```go
// files/go, as summarized in docs/source.md:95-100
type ID string
type Path string
type File struct { ID ID; Path Path; Text string }
type Files struct { Entries []File }
type SyntaxLink struct { File ID; Tree ast.TreeID; Source ast.SourceID }
type Context struct { Files Files; Links []SyntaxLink }

// contracts/semantic/go/contract.go:82-85, evidence.go:14-20
type Setting struct {
	Key   semantics.Key
	Value []byte
}
type EvidenceRequirement struct {
	Operation     semantics.Key
	Profile       string
	Configuration []Setting
	Reports       []Requirement
	Semantics     bool
}

// semantics/go/model.go:24, 30-37, 77-86, 110-114
type Key struct{ Owner, Name string }
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

File IDs and paths are unique within one `Files`; paths are relative and
slash-separated; `Files.Capture(limits)` admits an inventory and refuses
duplicates and invalid paths. A `Setting` is an opaque key and bytes: it is how
configuration, including a placement, travels through shared code that must
not interpret it. `semantic.Snapshot` is the immutable capture (section 1.2),
and `semantic.Result` is a receipt. `representation.Language` is a language
descriptor: a language ID, a profile name and the marker type.
`a.Compatible(b)` holds only if all three are equal.

### 3.2 What generation produces

```go
// contracts/generation/go/output.go:14-60
// SymbolIssuer is a target library's authority to issue and interpret immutable
// external destinations. Its key is opaque to the generation and Files kernels.
type SymbolIssuer struct{ owner *identity }
type ExternalSymbol struct {
	owner    *identity
	language representation.Language
	key      string
}

func NewSymbolIssuer() SymbolIssuer { return SymbolIssuer{fresh()} }
func (i SymbolIssuer) Issue(language representation.Language, key string) (ExternalSymbol, error) {
	if i.owner == nil || language.Validate() != nil || key == "" {
		return ExternalSymbol{}, ErrInvalid
	}
	return ExternalSymbol{i.owner, language, key}, nil
}
func (i SymbolIssuer) Read(language representation.Language, s ExternalSymbol) (string, bool) {
	if i.owner == nil || i.owner != s.owner || !language.Compatible(s.language) {
		return "", false
	}
	return s.key, true
}

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

The native target returns an `Assembly`. The kernel seals it. Sealing checks
that every exposed site and external-use site lies inside the output's trees,
then captures the trees **with the target language's syntax checker**:

```go
// contracts/generation/go/output.go:62-66, 88-113
func seal(ctx context.Context, checker check.Checker, assembly Assembly, publications int, limits Limits) (Output, error) {
	set := assembly.Syntax
	if len(set.Trees) == 0 || len(set.Trees) > 256 || len(assembly.Exposed) != publications {
		return Output{}, ErrIncomplete
	}
	// …
	contains := func(ref ast.NodeRef) bool { return ref.Node != 0 && uint64(ref.Node) <= uint64(trees[ref.Tree]) }
	exposed := map[ast.NodeRef]bool{}
	for _, site := range assembly.Exposed {
		if !contains(site) || exposed[site] {
			return Output{}, ErrScope
		}
		exposed[site] = true
	}
	seen := map[semantics.Occurrence]bool{}
	for _, use := range assembly.ExternalUses {
		if !contains(use.At.Site) || !use.At.Role.Valid() || seen[use.At] || use.Symbol.owner == nil || !checker.Language().Compatible(use.Symbol.language) {
			return Output{}, ErrScope
		}
		// …
		seen[use.At] = true
	}
	snapshot, err := semantic.Capture(ctx, checker, semantic.Input{Entry: set.Entry, Trees: set.Trees})
	if err != nil {
		return Output{}, err
	}
	return Output{snapshot, assembly.Exposed, assembly.ExternalUses}, nil
}
```

So an output carries three things: a snapshot of syntax-checked trees, the
exposed sites, and the external uses. Intent for a reference to a generated
declaration or to a type parameter is the `Reference` link inside the tree.
Intent for a reference to a library-owned type is an `ExternalUse`, whose
symbol only the issuing target library can read. The request that starts
generation names **one** output tree ID (`Output ast.TreeID`), and every target
today returns one tree.

### 3.3 The shared constructor every target uses

Every native target implements this interface, and delegates the walk over the
planned program to one shared function, `Compile`:

```go
// contracts/generation/go, as shown in docs/generation.md:66-75
type Target[T any] interface {
    CheckExpression(context.Context, Expression[T]) error
    CheckDeclaration(context.Context, NativeDeclaration[T]) error
    CheckProgram(context.Context, Program[T]) error
    Build(context.Context, BuildInput[T]) (Assembly, error)
}
```

`CheckProgram` is called while planning; `Build` constructs the tree. `Compile`
identifies every declaration and type parameter by a **symbol**, which is the
declaration's numeric ID in the plan plus a parameter index. When the target
creates a reference node it calls `Link(node, symbol, role)`. After all nodes
exist, `Compile` turns each link into a `Reference` to the node that the symbol
ended up at. A target calls `External(node, use)` instead for a library-owned
type, and `Compile` collects those into `Assembly.ExternalUses`. Finally the
target's `Finish(declarations)` wraps the declaration nodes in the root node.
That `Finish` call is where Go writes its package and C++ its namespace
(sections 3.7 and 3.10). Section 3.9 shows the name rules inside `Compile`.

### 3.4 The shared generated-Files contract

A language supplies a generated-Files formatter by giving the shared contract a
provider and three things: an evidence requirement, a printer and a check.

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
// contracts/format/go/generated.go:51-63
func NewGeneratedFilesFormatter[L any](language representation.Language, contract GeneratedFilesContract[L], provider Provider[GeneratedFilesRequest[L], source.Source[L]]) (GeneratedFilesFormatter[L], error) {
	if err := representation.CheckMarker[L](language); err != nil {
		return GeneratedFilesFormatter[L]{}, err
	}
	if representation.IsNil(provider) || provider.ID() != language.ID() || contract.Print == nil || contract.Check == nil || !contract.Evidence.Semantics {
		return GeneratedFilesFormatter[L]{}, ErrInvalidFormatter
	}
	if err := semantic.ValidateEvidenceRequirements([]semantic.EvidenceRequirement{contract.Evidence}); err != nil {
		return GeneratedFilesFormatter[L]{}, err
	}
	contract.Evidence = contract.Evidence.Clone()
	return GeneratedFilesFormatter[L]{language, provider, contract}, nil
}
```

Each request requires the formatter's language to be compatible with the
language of the output's snapshot, passes the settings through unread, and
installs the language's check as a guard on publication:

```go
// contracts/format/go/generated.go:73-80, 96-107
func (p GeneratedFilesFormatter[L]) Format(ctx context.Context, input GeneratedFilesInput) (source.Source[L], error) {
	if err := p.Validate(); err != nil {
		return source.Source[L]{}, err
	}
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
	return publishFiles(ctx, p.definition, policy, func(publisher source.Publisher[L]) (source.Source[L], error) {
		return p.provider.Format(ctx, GeneratedFilesRequest[L]{input, publisher})
	})
}
```

A provider cannot return a `Source` it built itself. It hands a
`source.Publication` (a fresh snapshot of the placed units, their analysis
receipt, the units, and optional groups and relations between files) to
`Publish`. The shared publisher then:

1. admits the receipt against the contract's evidence requirement, for exactly
   the published snapshot (`semantic.AdmitEvidence`, `source/go/source.go:212`);
2. requires the units to cover the snapshot completely, their trees to equal
   the captured trees, and their printed text to equal the captured file text;
3. runs the language's `Check` with the original input.

Three consequences matter for the questions:

- **One language, one checker.** The unplaced output and the placed source
  have the same language descriptor, so the same syntax checker admits both. A
  separate "unplaced Go" language or profile is excluded by the compatibility
  test above.
- **The shared contract knows nothing about a placement.** It carries
  `Configuration []semantic.Setting` and never reads it. What a placement is,
  and whether it is complete, is today decided privately by each formatter.
- **The evidence requirement is fixed at construction.** A formatter states it
  once, from the dependency environment it was built with. The placement
  arrives with each request.

### 3.5 TypeScript, the corrected path and working model

TypeScript was corrected on 8 November 2026. Before that, its target wrote
import statements and `export` modifiers during generation, and the Files stage
discarded and rewrote them. The case record calls this the generated syntax
"doing double duty as the record of which dependency was intended".

**The target writes no file decision.** Its root is the vocabulary's top node
with declaration children only. `string` and `boolean` are keyword nodes, and
a library-owned type is an unlinked reference recorded as an external use:

```go
// transform/syntax/typescript/target/go/provider.go:80-94, 115-117
func (b *Backend) externalType(module, export string, args []m.Fragment[s.Language], origin ast.NodeRef) (m.Fragment[s.Language], error) {
	ref, err := b.Node(m.Node(s.KindTypeReference, m.Attributes("name", export), origin))
	if err != nil {
		return ref, err
	}
	external, err := symbol.ExternalUse(b.language, ast.NodeRef{}, symbol.Type{Module: module, Name: export})
	if err != nil {
		return ref, err
	}
	ref, err = b.External(ref, external)
	if err != nil {
		return ref, err
	}
	return b.apply(ref, args, origin)
}
// …
func (b *Backend) Finish(declarations []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindSourceFile, nil, ast.NodeRef{}), edges(s.RoleDeclaration, declarations)...)
}
```

`symbol.Type{Module, Name}` is TypeScript's own symbol vocabulary: "a named
export, or a global type when Module is empty". This works without a new node
kind because the TypeScript grammar requires nothing on its top node: the
checker's rule for `KindSourceFile` lists no attribute, only optional `import`
and `declaration` children, and `exported` is an optional attribute of a
declaration.

**The placement setting** is a TypeScript-owned type in the formatter's
package. It has no scope field, because a TypeScript module's scope is its
file:

```go
// format/files/typescript/provider/go/placement.go:21-39
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
func (o Organization) Setting() (semantic.Setting, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return semantic.Setting{}, err
	}
	return semantic.Setting{Key: OrganizationSetting(), Value: data}, nil
}
```

**Its admission** is private to the same package. `p.tops` is the set of
top-level declaration sites of the input, and `p.original` is the organization
an already placed program has (empty for generation output):

```go
// format/files/typescript/provider/go/placement.go:168-210
func (p *program) organization(settings []semantic.Setting, limits source.Limits) (Organization, error) {
	o := p.original
	if len(settings) > 0 {
		if len(settings) != 1 || settings[0].Key != OrganizationSetting() {
			return Organization{}, source.ErrInvalid
		}
		o = Organization{}
		// A setting contains one strict object, including the explicit entry.
		decoder := json.NewDecoder(bytes.NewReader(settings[0].Value))
		decoder.DisallowUnknownFields()
		if !json.Valid(settings[0].Value) || decoder.Decode(&o) != nil {
			return Organization{}, source.ErrInvalid
		}
	}
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

Apart from `imports.SourcePath` (a `.ts` path that is not `.d.ts`), nothing
here is about TypeScript: strict decoding, a limit, no empty unit, an entry
file, every top-level declaration assigned exactly once and by site, and file
IDs and paths admitted by the Files inventory. One unit is one file, because
the inventory refuses a repeated file ID or path.

**Reading intent, and refusing a file decision in generated syntax.** The
formatter rebuilds its view of the generated program from the sealed output
alone. A root child that is not a declaration, and a declaration that carries
an `exported` attribute, are both refused. Publication intent is read only
from `Exposed()`:

```go
// format/files/typescript/provider/go/generated.go:61-115
func generatedProgram(ctx context.Context, output g.Output) (*program, error) {
	snap := output.Snapshot()
	p := &program{trees: map[ast.TreeID]s.AST{}, tops: map[ast.NodeRef]bool{}, public: map[ast.NodeRef]bool{}, bindings: map[ast.NodeRef]symbol{}}
	for _, env := range snap.Inputs() {
		t, err := representation.Unpack[s.Language](snap.Language(), env)
		if err != nil {
			return nil, err
		}
		p.trees[t.ID] = t
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

`symbol` is the formatter's record of a destination: a local site, an external
`{Module, Name}`, or a compiler builtin. After this, every type reference has
an intended destination, or formatting fails ("type reference lacks native
binding").

**Relocation, with a collision check per placed scope.** Declarations are
copied into one new tree per unit. `layout.moved` maps every old site to its
new site. Names are checked per unit, and a local use is re-spelled from the
declaration it links to, never from its own text. A declaration is never
renamed:

```go
// format/files/typescript/provider/go/relocate.go:119-126, 155-165
		unit.tree.Root = imports.AddNode(&unit.tree, s.Node{Kind: s.KindSourceFile})
		names := map[string]bool{}
		for _, site := range unit.placement.Declarations {
			name := imports.Attr(program.trees[site.Tree].Nodes[site.Node-1], s.AttributeName)
			if names[name] {
				return fmt.Errorf("%w: duplicate declaration %s in %s", source.ErrInvalid, name, unit.placement.Path)
			}
			names[name] = true
// …
		original := program.trees[dest.site.Tree].Nodes[dest.site.Node-1]
		name := imports.Attr(original, s.AttributeName)
		if moved.Tree == unit.tree.ID {
			node := &unit.tree.Nodes[id-1]
			imports.SetName(node, name)
			node.References = []ast.Reference{{Role: s.RoleTarget, Target: moved}}
			continue
		}
		if !program.tops[dest.site] || !program.public[dest.site] {
			return fmt.Errorf("%w: cross-unit type %s must be exported", source.ErrInvalid, name)
		}
```

An exposed declaration gets `exported: true`. A use of a declaration in another
unit, or of an external symbol, becomes a named type import; the import helper
picks an alias such as `Cell_1` when a local name collides.

**Fresh evidence, bound to the environment.** The formatter is constructed
with its collaborators and the dependency environment, and it refuses to exist
without the environment:

```go
// format/files/typescript/provider/go/format.go:23-34, 43-46
// Stages are the explicitly supplied collaborators for realizing and checking
// TypeScript units. Configuration is the native compiler's receipt-bound import
// environment. Organization is separate and belongs to each Format request.
type Stages struct {
	Language      representation.Language
	Syntax        f.Provider[ast.View[s.Language], cst.CST[s.Language]]
	Printer       print.Provider[cst.CST[s.Language], string]
	Checker       check.Checker
	Resolver      resolution.Resolver
	Analyzer      analysis.Analyzer
	Configuration []semantic.Setting
}
// …
	_, enabled, err := n.Imports(p.Configuration)
	if err != nil || !enabled {
		return Stages{}, fmt.Errorf("%w: named TypeScript import environment required", source.ErrInvalid)
	}
```

Its contract evidence is `n.EvidenceFor(p.Configuration)`: the analysis
profile under exactly that environment. After relocation the formatter formats
each tree, captures the units with their paths, runs the resolver and the
analyzer, and publishes.

**The check.** The guard recomputes the expected layout from the original
input (`placeGenerated` again) and compares it with the publication. `expected`
is the recomputed layout, `realized` is read from the fresh semantic graph of
the published snapshot:

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
   graph's `Binding` records. For every relocated type reference, the entity
   the compiler resolved it to must be the relocated intended declaration, the
   intended package export, or the intended builtin.
3. **Exposure, from the compiler.** `realized.public` is built from
   `named-export` facts, which the TypeScript bridge publishes for each module
   entity. A declaration is public in the result exactly if it was exposed.

Parts 2 and 3 are facts the formatter did not produce itself.

**Consumer composition**, as the implementation document shows it. `tsfiles`
is the formatter's package; `page`, `cells` and `person` are sites in
`out.Snapshot()`; `stages` holds the collaborators and the environment:

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

The qualification test generates `Cells = Page<Cell<Person>>` once, together
with a local declaration also named `Cell`, and formats it twice. Moving
`Person` next to `Page` removes only its relative import; the package import
stays, aliased `Cell_1`. The test also asserts that the generated tree contains
no import node and no `exported` attribute.

### 3.6 A second private copy of the same assignment rules

The horos Files formatter places whole horos modules, identified by pin, into
files. It has its own placement type and its own admission, in a private
package:

```go
// internal/horosfiles/organization.go:18-24, 100-109, 114-126, 170-172
// Placement is an explicit same-language organization decision. Modules are
// canonical module pins, never aliases selected by spelling or filesystem lookup.
type Placement struct {
	ID      files.ID
	Path    files.Path
	Modules []string
}
// …
	if len(placements) > limits.Files.MaxFiles {
		return fail(source.ErrLimit)
	}
	inventory := files.Files{}
	for _, placement := range placements {
		inventory.Entries = append(inventory.Entries, files.File{ID: placement.ID, Path: placement.Path})
	}
	if _, err := inventory.Capture(limits.Files); err != nil {
		return fail(err)
	}
// …
	for _, placement := range placements {
		if len(placement.Modules) == 0 {
			return fail(source.ErrInvalid)
		}
		// …
		for _, pin := range placement.Modules {
			module, exists := modules[pin]
			if !exists || seen[pin] {
				return fail(source.ErrInvalid)
			}
			seen[pin] = true
// …
	if len(seen) != len(modules) || entry == "" {
		return fail(source.ErrInvalid)
	}
```

These are the same rules for a different kind of item: a limit, admitted file
IDs and paths, no empty unit, each item assigned exactly once, a complete
assignment, an entry. Go and Rust formatters written like the TypeScript one
would hold the third and fourth copy.
