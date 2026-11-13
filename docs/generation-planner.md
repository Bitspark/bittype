# Planning and running type generation

**Implemented, 8 November 2026.** Planning resolves source meanings to intended
native types. A run constructs the resulting target ASTs. Files then organizes
and checks those declarations under a placement and dependency environment.

## Three SPI changes

1. **Carry generation intent into Files — implemented for TypeScript.**
   A sealed generation.Output holds checked declaration ASTs, an inventory of
   the declarations it owns, requested exposed sites and external use
   destinations. Files receives that output, one placement
   and the selected native dependencies, builds module syntax and checks its
   bindings before publishing Source.
2. **Separate semantic mapping from target construction — implemented.**
   Libraries register descriptions, declaration recipes and generic conventions.
   Planning resolves every reachable choice, source capture and helper into a
   closed native Program. Run receives that program, without source readers or
   mapping callbacks. The procedural operations and per-request mapping lists
   have been removed. [SPI](generation-spi.md)
3. **Read common structure from established semantics — implemented.**
   The kernel admits the graph from the evidence and walks declarations,
   references, applications and children itself. Sources keep language-owned
   identity, operand meaning and applicability judgments through a narrow
   interpretation contract and a reader for their mapping libraries.
   [SPI](generation-spi.md#semantic-traversal-and-source-interpretation) ·
   [GEN-01: correction and scope](generation-gaps.md#gen-01)

The C++ include and namespace, Rust visibility, Go package clause and Go naming cases remain open in
[ACX-004](architecture-cases/0004-generation-file-syntax.md). Language route
planning through ASTTransformer remains a separate operation; it is not the
source-bound generation plan documented here. [GEN-02](generation-gaps.md#gen-02)
separately records missing native enforcement, including the native parameter
contract limitation. Neither of those was closed by the declarative planning or
traversal changes.

## Caller interface

Implemented Go signatures, with compose denoting the shared generation registry
and generation denoting the public contract package:

```go
func (e *Engine) PlanGeneration(ctx context.Context, r compose.Request) (generation.GenerationPlan, error)
func (e *Engine) Generate(ctx context.Context, r compose.Request) (generation.Output, error)

type Request struct {
    Input   transform.Input
    To      representation.Language
    Roots   []ast.NodeRef
    Output  ast.TreeID
    Profile string
    Limits  generation.Limits
}

func (p GenerationPlan) Steps() []Step
func (p GenerationPlan) Run(context.Context) (Output, error)

type Step struct {
    Binding, Source, Declaration, Form string
    Expression string
}
```

The consumer registers its available source, target and mapping libraries once.
The kernel selects compatible descriptions and exact replacements for the
requested profile. The request carries source evidence and explicit public roots;
it does not choose an imperative pipeline. [Complete registration example](kernel-composition.md)

```go
plan, err := eng.PlanGeneration(ctx, compose.Request{
    Input: transform.Input{Source: snapshot, Evidence: []semantic.Result{horosAnalysis}},
    Roots: []ast.NodeRef{pageSite, personSite, cellsSite},
    To: typescriptLanguage, Output: "pagination-types",
    Profile: rules.RepresentationProfile, Limits: generation.DefaultLimits(),
})
if err != nil { return generation.Output{}, err }
return plan.Run(ctx)
```

This calling fragment assumes the registered engine and source sites from that
example. Steps returns copied inspection summaries. Execution uses the private
typed program, never parses these display strings. Each Run returns a fresh
capture. Registering another binding changes future plans, not an issued plan.

### The result of Run

Implemented accessors supplied by generation and consumed by Files:

```go
func (o Output) Snapshot() semantic.Snapshot
func (o Output) Exposed() []ast.NodeRef
func (o Output) ExternalUses() []ExternalUse

type ExternalUse struct {
    At     semantics.Occurrence
    Symbol ExternalSymbol
}
```

Snapshot owns syntax-checked target ASTs. References to generated declarations
and lexical parameters already identify their destination nodes. Exposed lists
the realized sites requested by the consumer, not every dependency or helper.
ExternalUses binds emitted occurrences to target-owned immutable public symbols.
TypeScript's symbol identifies a module and its named export, without choosing
a local import name. The sealed output associates all these addresses with one
capture; zero output represents failure. [Output contract](../contracts/generation/go/output.go)

### The declaration inventory and its handles

**Implemented, 13 November 2026** ([#272](https://github.com/Bitspark/bittype/issues/272),
rows 10, 12, 19 and 5 of [research 0016](../research-docs/0016-native-files-placement.md#applied)).
A sealed output lists every declaration it owns, and it issues one handle for
each of them:

```go
func (o Output) Inventory() []OwnedDeclaration
func (o Output) Resolve(DeclarationHandle) (OwnedDeclaration, error)

type OwnedDeclaration struct {
    Handle     DeclarationHandle
    Site       ast.NodeRef
    Role       DeclarationRole
    Provenance []ast.NodeRef
    Values     []IssuedValue
}

const (
    RequestedDeclaration DeclarationRole = iota + 1
    DependencyDeclaration
    HelperDeclaration
)
```

Inventory returns one record for each declaration of the snapshot and no other
record. A consumer therefore selects a declaration by what its record states.
It does not walk the tree and it reads no `name` attribute.

| Field | Meaning | Who states it |
| --- | --- | --- |
| `Handle` | One declaration of this one output. | The kernel, when it seals. |
| `Site` | The declaration's node in the snapshot. | The target's construction. The kernel checks that the site lies in the snapshot and is listed once. |
| `Role` | Requested: the request lists its source declaration as a root. Dependency: a planned declaration uses its source declaration and the request does not list it. Helper: a mapping demanded it by a helper key. | The planner. |
| `Provenance` | The declaring sites of the source declarations it represents. Empty for a helper, which represents none. | The planner, from the admitted semantic graph. |
| `Values` | What the owner of the declaration's description states about it. | A declaration recipe or a helper recipe. |

- **Order.** The records follow their sites: the trees in the snapshot's order,
  the nodes of one tree ascending. Research 0016 names the site as the order of
  the declarations of one file; the Files contract that uses it is
  [#273](https://github.com/Bitspark/bittype/issues/273). `Exposed` keeps the
  order of the request.
- **Requested records and `Exposed`.** The records with the role requested have
  exactly the sites of `Exposed`. The kernel derives both from the plan, so a
  target cannot make them disagree.
- **Provenance is a list.** The planner maps one source declaration to one
  native declaration today. The record does not assume it: a native declaration
  may later represent several source declarations, and one source declaration
  may yield several native ones.
- **A record speaks of ownership.** It says that the output owns a declaration
  at a site. It does not say that the site is one entity of a native compiler,
  and no bijection between sites and compiler entities is fixed.
- **No grouping field.** A consumer groups by provenance. A grouping that a
  later case needs enters as an issued value.

**A handle carries its output**, as an external symbol carries its issuer. A
site is a tree identity and a node number, and both repeat on every run of a
plan. Two handles are equal only if one output issued both for one declaration.
`Resolve` returns the record of a handle this output issued and fails with
`ErrScope` for every other handle. That includes a handle of a second run of
the same plan. A placement is built from handles, so it cannot be applied to an
output it was not written for.

**Issued values.** A naming intent and the access a helper requires are
meanings of a mapping and of a target language. The kernel does not interpret
them:

```go
func NewValueIssuer() ValueIssuer
func (i ValueIssuer) Issue(content string) (IssuedValue, error)
func (i ValueIssuer) Read(IssuedValue) (string, bool)
```

The owner of a meaning holds its issuer. A declaration recipe states values in
`DeclarationDescription.Values`, and a helper recipe states them in its
`NativeDeclaration.Values`. The planner checks that each value was issued and
charges its content to the output budget. It then carries the value to the
record of that declaration. Only the issuer reads the content back, and one
issuer is one meaning. A record without values is valid: no builtin mapping
states one until the naming design of
[#274](https://github.com/Bitspark/bittype/issues/274) lands.

**Why this is in the shared contract.** Platform's rule
[AC-08](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-08)
asks for the missing capability, its layer and whom it serves.

- **The capability.** Address every declaration an output owns without reading
  native syntax, and bind a placement to the output it was written for. Before,
  a consumer found an unexposed declaration by walking the tree and comparing
  names, and a placement written for one run was accepted with the output of
  another run.
- **The layer.** The generation kernel issues declaration identity. The planner
  alone knows which declarations a plan has, which of them the request named,
  which are helpers and which source declarations they represent. The kernel
  seals the capture that every site belongs to. No target and no Files stage
  can state these facts, and each of them would have to state them again.
- **The contributions it serves.** Each language's Files stage places
  declarations through the same records and handles. A mapping or a target
  language adds a meaning by issuing a value, and the kernel, the other targets
  and the other mappings do not change.
- **The boundary change.** It is not only an addition. `Assembly.Exposed` is
  replaced by `Assembly.Declarations`: a target now returns the site of every
  planned declaration, and the kernel derives the exposed sites from the plan.
  `PlannedDeclaration.Exposed` is replaced by `Role`, and it gains
  `Provenance`. `DeclarationDescription` and `NativeDeclaration` gain `Values`.
  The four targets build through the shared constructor and did not change. A
  target that fills its own `Assembly` has to list every declaration.

**What a later change does not alter:** the accessor `Inventory` and that it is
complete, the fields of a record and their meaning, the handle type and that it
identifies one declaration of one sealed output, and the key that orders sites
(see [execution guarantees](#execution-guarantees)). Adding a field to the
record is additive.

The typed Files boundary remains:

```go
func (p GeneratedFilesFormatter[L]) Format(ctx context.Context, input GeneratedFilesInput) (source.Source[L], error)
```

The application invokes it with the complete output, placement and environment.
The selected Files implementation relocates declarations, spells imports and
exports and compares fresh native bindings with the output's intent.
[TypeScript handoff](typescript-type-imports.md#from-generation-to-files)

## Semantic mapping and AST construction

| Stage | Input | Result and authority |
| --- | --- | --- |
| Read A | Files and selected source dependencies | Parsing, resolution and analysis establish captured syntax and Semantics[A]. |
| Map to B | Captured A facts, roots, registered descriptions and a profile | A closed native program: declarations, expressions, binders, complete applications and requested exposures. Mapping libraries own representation choices. |
| Construct B | The closed program and target construction rules | A syntax-checked generation.Output. The target realizes each native form and binds its generated uses. |
| Organize B | Output, one placement and target dependencies | Files supplies module syntax and local names. |
| Check and write B | Organized units in their native environment | Native resolution/analysis establishes Semantics[B]. Files verifies destinations and publication, then the printer returns paths and text. |

The planned native program is intent, not a fabricated semantic-analysis result.
Only analysis of the realized target program establishes native semantic facts.
Cross-language representation proofs use the separate verification contract.
[Abstraction tables](abstraction-levels.md#page-across-the-five-levels) ·
[Generation](generation.md#semantic-mapping-and-target-ast-construction)

### Pagination on read and write

Suppose the caller requests `Page`, `Person` and `Cells` as public TypeScript
declarations. The application supplies three Horos files with `people.horos`
as the entry. Each pin placeholder stands for the 64 hexadecimal digits of its
module's pin. The `Cell` declaration follows the independent library fixture.
[Horos Files reading](source.md#horos-composition-and-pagination) ·
[Cell mapping fixture](../integration/go/horos_generation_composition_test.go).

```horos
// pagination.horos
module example.pagination;
type Page<T> = { items: list<T>; next?: text; };

// cells.horos
module example.cells;
type Cell<V> = wire<V>;

// people.horos
module example.people;
import pages from example.pagination pin "sha256:<pagination-pin>";
import cells from example.cells pin "sha256:<cells-pin>";
type Person = { name: text; };
type Cells = pages.Page<cells.Cell<Person>>;
```

**Read Horos.** Syntax parsing retains the written names and nested applications.
The resolver follows the two import pins and resolves `Person` in
`example.people`. Analysis establishes that the inner application binds
`Cell.V` to that `Person`, and the outer one binds `Page.T` to the whole
`Cell<Person>` application. The `T` in `items` still belongs to the open
`Page` declaration. These facts concern the Horos input and its scopes.
[Pagination analysis](analysis.md#pagination-from-input-to-queries) ·
[Horos Files fixture](../integration/go/files_test.go).

**Map to TypeScript.** The structural mappings select an interface for each of
`Page` and `Person` and a type alias for `Cells`. The list mapping chooses an
array whose element is the target parameter of `Page`. The Cell library's
mapping matches the pinned Horos `Cell` declaration, selects the public type
export `Cell` from `@bitspark/cell`, and uses the mapped `Person` as its
argument. The outer application uses that entire result as the argument to
the generated `Page`. **Semantic mapping ends here.** The representation,
each reference destination and the three declarations to expose are fixed.

**Construct the TypeScript AST.** The builder now has enough information to
construct the interfaces, type parameter, properties and alias. The outline
shows containment and the intended destination of each reference. The
external package identity is an association on the use.

```text
interface Page
  type-parameter T
  property items → array-type → type-reference to Page.T
  optional property next → keyword-type string
interface Person
  property name → keyword-type string
type-alias Cells
  type-application
    constructor → type-reference to generated Page
    argument    → type-application
      constructor → type-reference to Cell exported by @bitspark/cell
      argument    → type-reference to generated Person
```

A syntax check admits these node kinds, roles and containment. Call the result
of `plan.Run(ctx)` for these roots `out`. The generated declaration names
below stand for sites in its target snapshot. [TypeScript AST vocabulary](syntax.md#typescript).

| Go accessor | Value for `Cells = Page<Cell<Person>>` |
| --- | --- |
| `out.Snapshot()` | The two interfaces and the alias shown above. The uses of `Page` and `Person` link to their declarations. The element `T` links to the parameter on `Page`. |
| `out.Exposed()` | The declaration sites of `Page`, `Person` and `Cells`. |
| `out.Inventory()` | Three records, one for each of these declarations, in the order of their sites. Each has the role requested, and its provenance is the horos site the request named. Had the request named only `Cells`, the records of `Page` and `Person` would have the role dependency. |
| `out.ExternalUses()` | One entry. Its occurrence is the constructor reference of the **inner** application in `Cells`. Its symbol identifies the public export `Cell` from `@bitspark/cell`. |

**Organize the output.** Run once and reuse that `out` for each placement.
Put `Page` and `Cells` in `pages/page.ts` and the public `Person` in
`model/person.ts`. Given the selected Cell package,
TypeScript Files produces the following module.

```typescript
import type { Person } from "../model/person.js";
import type { Cell } from "@bitspark/cell";

export interface Page<T> { items: T[]; next?: string; }
export type Cells = Page<Cell<Person>>;
```

Change only the placement to put `Person` beside `Page`. The pages module now
also contains `export interface Person { name: string; }`. Its relative import
disappears. The package import and the alias body stay the same.

**Check the bindings after placement.** As Files organizes the units, it keeps
the correspondence from each generated use and declaration to its new site.
For each layout it then resolves and analyzes the organized units in the
selected TypeScript environment. Each analysis has its own entity IDs. Files
compares its findings with the destinations and publication requests in `out`:

| Intended binding or publication | Check in the organized units |
| --- | --- |
| The outer constructor and inner argument refer to `Page` and `Person`. | Resolve both uses. They must reach the relocated declarations, whether `Person` is reached through the relative import or directly in the same module. |
| The element `T` in `items` refers to the binder on `Page`. | Resolve the element reference and compare its declaration with that specific parameter on the relocated `Page`. |
| The inner constructor refers to the selected package export. | If imported as `Cell as C`, the `C` in `Page<C<Person>>` must resolve to the public `Cell` from `@bitspark/cell`. |
| `Page`, `Person` and `Cells` were requested as public declarations. | Check the module exports against the corresponding relocated sites. An inaccessible dependency cannot silently become public to make a placement work. |

A missing export, an inaccessible target or a same-named declaration capturing
a use causes formatting to fail without publishing a partial `Source`. After
the checks pass, Files publishes the `Source` for that layout with its fresh
native evidence. The printer emits the paths and text.
[TypeScript placement and binding checks](typescript-type-imports.md#from-generation-to-files) ·
[Placement and alias fixtures](../integration/go/typescript_placement_test.go).

### Where the plan stops

The caller requests roots to expose, the plan identifies their target
declarations, and the run records the realized sites. Each named use retains the
identity of its owned declaration, lexical parameter or external public symbol.
A change of placement uses the same plan. Selecting a different external type
for Cell changes a mapping decision and requires a new plan.

The implemented TypeScript handoff is the [complete output of `Run`](#the-result-of-run).
Files takes the intended destinations and requested publications from that
output. Placement and the selected dependencies come from the consumer.
The target resolver and analyzer establish the actual destinations for each
placement.

The implemented TypeScript path consumes the whole generation output. Files constructs import and export syntax for the selected placement, then checks the relocated destinations and named exports before publishing Source. A separate path formats programs with established semantics. [Generated Files input](../contracts/format/go/generated.go) · [TypeScript example](typescript-type-imports.md)

## What a mapper declares

The [implemented SPI](generation-spi.md#mapping-descriptions) replaces the
description-plus-executable-authority proposal recorded in
[ACX-005](architecture-cases/0005-generation-interface-design.md).

For a fixed rule, a library supplies a source selector, authenticated captures
and an immutable target expression. This is the TypeScript list rule's body
(fragment inside Describe, using the actual public types; `patterns` is
`transform/syntax/horos/mapping/pattern/go`):

```go
element := generation.TypeSlot("E")
pattern, err := patterns.Element(from, element)
if err != nil { return description, err }
binding := generation.BindingDescription[horos.Language, typescript.Language, source.Reader]{
    ID: "list/typescript/1",
    Input: pattern,
    Output: model.Collection(generation.Mapped[typescript.Language](element)),
}
```

The Horos rule owns the constructor and element occurrence. The native model
owns the array form. The kernel owns capture identity, substitution and scope.
The target later constructs an array AST with the planned child. Neither the
kernel nor the native target learns a Cell kind: the Cell library selects an
external type form with its mapped argument through the same contract.

A variable-field record uses the implemented expansion boundary:

```go
Expand func(context.Context, generation.Expansion[F,T,S]) (generation.Expression[T], error)
```

Its library reads its own fields and calls the scoped request's Mapped method
to record child demands. The kernel authenticates those children, completes
them during planning and revokes the request after expansion. It never retains
that callback in the runnable plan. An exact rule that fails does not fall back
to structural generation. Structural declaration recipes are explicitly supplied
and apply only when no exact binding claims the declaration.

## What the plan contains

Planning admits exact snapshot evidence, authenticates roots and source
environments, resolves demanded rules and captures complete applications.
A single GenericConvention describes native parameters, the representation of
each source binder and arguments at each use. Extra or reordered native
parameters use the same convention; helpers have their own parameter scope.
Every source argument is admitted, even when native syntax omits it.
[Generic contract](generation-spi.md#source-and-representation-environments)

Owned declarations and helpers are resolved once by identity. References retain
their destination, not just a spelling. The current profile rejects cyclic
declaration dependencies. Missing demanded rules, ambiguous exact choices,
foreign slots, unsupported forms and scope failures prevent a runnable plan.

## Execution guarantees

The selected native target implements this exact construction boundary:

```go
Build(context.Context, generation.BuildInput[T]) (generation.Assembly, error)

type BuildInput[T any] struct {
    Program  Program[T]
    Language representation.Language
    Output   ast.TreeID
    Limits   Limits
}
```

The kernel invokes Build with the frozen native program. The target interprets
its own forms and uses fresh fragment/symbol storage for each run. The constructor
has neither a source reader nor mapping expansion code. The kernel then checks
syntax and seals the inventory, the exposed sites and the external-use records
into Output.

The assembly lists the realized site of every planned declaration, in planned
order. Sealing fails with `ErrIncomplete` when that list has another length
than the plan, and with `ErrScope` when a site lies outside the snapshot or is
listed twice. The kernel has no knowledge of a language's syntax. That each
site is the declaration's own node is the guarantee of the shared constructor,
which the four targets build through. It returns no assembly when a planned
declaration has no node in the unit.

**The order of sites.** The shared constructor emits each declaration after the
declarations it uses. The walk is seeded in the order of one comparator, and
the node numbers, and so the sites, follow the walk. The comparator's key is
the name, then the origin by tree and node, then the planned identity. Equal
names compare equal and the rest of the key decides, so the comparator is an
ordering and returns 0 only for one and the same planned declaration.
Limit: two declarations of equal name and equal origin are ordered by their
planned identity, which follows the order of the requested roots. Two helpers
of one name are such a pair, because a helper has no origin, and so are two
native declarations of one source declaration. The shared constructor still
refuses two declarations of one name, so no plan reaches the comparator with
such a pair. Removing that rule is [#273](https://github.com/Bitspark/bittype/issues/273).
[Constructor](../transform/syntax/shared/target/go/compile.go)

Changing mutable callback state after planning cannot change an array into a
ReadonlyArray at construction. Target implementations are trusted interpreters
of their own native forms; syntax admission is not a proof of faithful lowering
or a sandbox for arbitrary Go code. Target configuration must remain stable
after registration. Cancellation, limits, invalid target syntax and foreign
references return zero output. [Complete SPI and trust limits](generation-spi.md)

## Qualification

- Page<Cell<Person>> exercises a library-owned external type whose Wire body is
  never expanded. A two-argument Pair mapping reverses its native arguments
  inside Page through the same public contract.
- Plan inspection precedes any AST construction. Repeated and concurrent runs
  use fresh state; a later ambiguous registration affects only new plans.
- A non-Horos source fixture composes an extra native generic parameter and a
  helper, and its generated Go is checked by the Go compiler.
- Negative cases cover wrong receipts, roots, scopes, foreign slots/forms,
  expired expansion requests, missing helpers, budgets and cancellation.
- The existing Files fixture places one Page<Cell> output twice with a conflicting
  local Cell: moving Person removes only its relative import; the package
  Cell remains correctly bound through an alias.
- The independent fixture generates three roots that need one dependency and one
  helper. Its five declarations are selected through the inventory alone, by
  role, provenance and issued value. A handle of one run is refused by a second
  run of the same plan, and an assembly that omits, repeats or misplaces a
  declaration is not sealed. The pagination module gives a complete inventory
  in all four targets.

[Plan and independent mapping tests](../integration/go/generation_planning_test.go) ·
[Inventory and handles](../integration/go/generation_inventory_test.go) ·
[Independent source](../integration/go/generation_independent_test.go) ·
[Four-target pagination](../integration/go/horos_generation_pagination_test.go) ·
[Files placement](../integration/go/typescript_placement_test.go) ·
[Installed consumer](../integration/go/testdata/installed/generation_planning_test.go)
