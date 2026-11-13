# Syntax and Files

**4 November 2026. Implemented Go contracts and Horos composition.**

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
| Abstract view | CST exposes AST | Source exposes Semantics |

Source and Semantics have no Model suffix. Files and Syntax are processing
dimensions, not language IDs. File placement, grouping and concrete import routes
belong to Source; their construction belongs to Files formatting. There is no
separate layout carrier, layout provider or combined Source provider.

The [abstraction levels guide](abstraction-levels.md) follows one declaration
from bytes to established bindings. It distinguishes the level of a name or
binder from the file and scope used to interpret it.

Resolution and analysis operate on the semantic side of the Files dimension.
They remain independently callable on captured ASTs; the dimension is an
architectural grouping, not a requirement to open physical files. Source-owned
`bittype:dimension` doc comments on boundary types express this grouping for the
[generated provider matrix](providers.md). The shared language identity contract
uses `bittype:shared`. These annotations introduce no runtime methods or registry;
the generator reads declared operation boundaries, never storage fields or
package names, using the [metadata rules](../tools/provider-matrix/README.md#dimension-metadata).

## From generated declarations to files

**Design decision, 7 November 2026.** Type generation uses the source semantics
and selected mappings to build a target declaration and bind its type references.
Files formatting decides how to express
those bindings and the requested public surface in a particular arrangement of
native modules. It can do that entirely in memory, before anything is written.

| Owner | Decision |
| --- | --- |
| Selected mapping and generation target | Build declaration bodies and type expressions. Bind each named use to its intended target declaration, lexical parameter, or selected external public symbol. The plan records which declaration the caller requested to expose. A generated tree ID identifies syntax; it does not assign a file. |
| Consumer | Supply one placement for the declarations being emitted: output files and the native scopes they belong to. Supply the selected external dependencies in the target environment. There is no second table of imports. |
| Target-language Files formatter | Follow those bindings under that placement. Choose the local spelling, qualification or alias, and construct the necessary import, export or other module syntax. Realize the requested visibility under the language's rules and establish that uses still resolve to their intended targets. |
| Shared Files printer | Emit the paths and text already represented by `Source`. The caller decides whether and where to write them. |

For example, generate the alias `Cells = Page<Cell<Person>>`. `Page` and
`Person` refer to declarations in this output; the Cell library's selected
mapping identifies the public TypeScript export `Cell` from `@bitspark/cell`.
That identifies *which type* is needed. It specifies neither an import
statement nor the file inside that package. The generator needs no destination
path for any of the three names. A spelling alone is insufficient: two distinct
types named `Person` must not become interchangeable when placement changes.

**13 November 2026.** [Destinations](destinations.md) is the specification of
that binding: which nodes of each language are name uses, how a sealed output
records the destination of each, and what counts as the same destination. A
destination is a declaration's identity and never structural equality. An
external symbol names the selected public export and not the canonical type
underneath it. Provenance may be many to many. TypeScript, Go and Rust output
record a destination for every name use. Legacy C++ output records none
(row 36 of [research 0016](../research-docs/0016-native-files-placement.md)).

**The rule for every generated-Files formatter.** A generated-Files formatter
refuses a name use without a destination. It does not look the name up, and it
does not accept the binding that the compiler happens to find. The shared seal
cannot enforce this, because it does not know which nodes of a language use a
name, so each language's formatter does, and a
[test per target](../integration/go/destinations_test.go) asserts it over sealed
outputs. The TypeScript formatter implements the rule. Go Files and Rust Files
implement it with their first slice
([#161](https://github.com/Bitspark/bittype/issues/161),
[#159](https://github.com/Bitspark/bittype/issues/159)). The
[failure categories](#failure-categories) give the refusal its category.

| TypeScript placement for the same generated declarations | What Files puts in the module containing `Cells` |
| --- | --- |
| `Page` and `Cells` in `pages/page.ts`; `Person` in `model/person.ts` | A type import of `Person` from `../model/person.js` and a type import of `Cell` from `@bitspark/cell`. |
| Move `Person` into `pages/page.ts` | `Person` is in scope. Its import disappears; the Cell package import remains. |

Both placements use the same generated declarations, bindings and publication
intent. Files must preserve those bindings and the requested public surface;
moving a declaration cannot authorize exposing a private type or substituting
another type of the same name.
If it cannot express the placement with the supplied dependencies and visibility,
it must fail without returning a partial `Source`. Native scope matters: two
files in one Go package need no import between them, while two Rust modules can
occupy one file. A file boundary alone is never an import rule. The [native placement proposal](native-file-placement.md)
walks Page/Cell through Go, Rust and C++, including same-package and
same-file/different-scope cases; those Files implementations remain absent.

The distinction is about module visibility. A member's accessibility, when it
is part of the chosen type representation, still belongs to the declaration.

This is the boundary for **generating declarations**. Native ASTs can still
represent imports and exports when parsing existing source or constructing
organized units. The TypeScript Files stage constructs imports and export modifiers using the sealed generation output and its placement. The compiler checks the relocated declaration sites, external package symbols and actual module exports. [TypeScript implementation](typescript-type-imports.md) · [ACX-004](architecture-cases/0004-generation-file-syntax.md).
The generation design separates semantic mapping
from target AST construction. The planner follows the same nested reference
from Horos resolution through TypeScript file organization and binding checks.
[Generation stages](generation.md#semantic-mapping-and-target-ast-construction) ·
[Page and Cell walkthrough](generation-planner.md#pagination-on-read-and-write).

## The generated-Files contract

**13 November 2026. Implemented in `contracts/format/go`, with TypeScript as its
one language.** It delivers rows 1, 5, 6, 9, 10, 20 and 27 of
[research 0016](../research-docs/0016-native-files-placement.md) with
[#273](https://github.com/Bitspark/bittype/issues/273).

A placement type belongs to its language. TypeScript states files and the
declarations of each. Go will state packages, Rust a module tree and C++
namespaces that do not depend on files. One shared mechanism admits what every
placement has in common, and it runs in the shared formatter, before the
provider.

```go
// contracts/format/go
type GeneratedFilesInput struct {
    Output    generation.Output
    Placement GeneratedPlacement // issued by the language of the formatter
    Limits    FilesLimits
}
type GeneratedPlacementView struct {
    Files       map[files.ID]files.Path
    Definitions map[generation.DeclarationHandle]files.ID
}
type GeneratedFilesPrepared[L any] struct {
    Placement AdmittedPlacement
    Evidence  semantic.EvidenceRequirement
    Check     func(context.Context, source.Publication[L]) error
}
type GeneratedFilesContract[L any] struct {
    Print   func(context.Context, cst.CST[L]) (string, error)
    Prepare func(context.Context, GeneratedFilesPreparation) (GeneratedFilesPrepared[L], error)
}
```

| Step | Owner | What happens |
| --- | --- | --- |
| 1 | Consumer | Builds the language's placement from the handles of `Output.Inventory()` and has the language issue it. A handle carries its output, so a placement belongs to one run. No placement bytes cross the contract. The engine carries the value and does not read it. |
| 2 | Shared formatter | Checks the request: the language of the output, the limits, and the sealed output against the semantic limits. |
| 3 | Language, in `Prepare` | Reads its placement back. A value of another issuer and a request without a placement are invalid input. Derives the view and calls `Admit`. |
| 4 | Shared, `Admit` | Runs the common admission, once for the request. A refusal is final: every later call is refused, and the request fails with the refusal and its category, whatever the step returns. |
| 5 | Language, in `Prepare` | Decides its own rules on the admitted placement, validity before coverage. Returns the admitted placement, the evidence requirement of this request and the guard. |
| 6 | Shared formatter | Refuses a step that returns without the admitted placement of this request, without a guard or with a requirement that asks for no semantic graph. Calls the provider. A refused request never reaches it. |
| 7 | Provider | Realizes the admitted placement, has the fresh files analyzed and publishes. |
| 8 | Shared publisher and formatter | Match the receipt with the requirement of step 5, refuse a publication whose file identities and paths differ from the admitted files table, and call the language's guard. |

**The common admission.** It refuses, in this order:

1. more files than the file limit or the link limit;
2. a handle that the sealed output of the request did not issue;
3. an assignment to a file that the files table does not declare;
4. a declaration of the output's inventory without an assignment;
5. a files table that `files.Files.Capture` refuses: no file, an empty
   identity, a malformed path, a path used twice, a path that is also a
   directory of another, or identities and paths over the byte limit.

The number of files comes first, as [#271](https://github.com/Bitspark/bittype/issues/271)
fixed it. Each refusal is invalid input or an exceeded limit, with its
[category](#failure-categories). None is unsupported coverage: validity is
decided here, before a language refuses what it does not cover. The placement
of existing TypeScript source follows the same order with its own first two
steps: the number of files, then a file identity or a declaration that its
lists name twice, then steps 3 to 5 through `files.Cover`. A declaration that
the program does not have is refused there, as an assigned item that is not
part of the placed set.

**The step calls the admission while it runs.** Row 1 of research 0016 says
that the contract "receives the admitted view from one language-supplied step
and runs the common checks itself". The step hands its view to `Admit` and
returns the `AdmittedPlacement` that `Admit` produced, and the formatter
refuses anything else. So a language decides its own rules on an admitted
placement, and its refusal of coverage can follow the common validity checks.
This reading has three limits, and Go and Rust are written against them:

- *The order inside a step is the step's duty.* The shared side cannot see
  what a step refuses before it calls `Admit`. A step that refuses coverage
  first reports an input as unsupported that the admission would have found
  invalid.
- *The view is the language's word.* The shared side admits the view that the
  step derives, and it does not read the consumer's placement. A step that has
  another view admitted than its placement states is not caught by the shared
  side. The language's guard answers for that, under the rule for guards.
- *The error a step returns is the step's word.* When the admission has
  refused, the request fails with that refusal, and a step that passes it on
  gets it back as it is. A step that wraps the refusal in an error of another
  category returns a refusal that matches two categories. The shared side
  does not take a step's error apart, so one category per refusal is the
  step's duty too.

Both tables of the view are keyed. A language whose own placement type can
state a file twice or a declaration in two files refuses that before it
derives the view, as TypeScript does for the lists of existing source.

The admission interprets no scope. It accepts a file without a definition, which
a Rust crate root that holds only module declarations needs, and it accepts two
declarations of one name in one file. A language refuses those for itself:
TypeScript refuses the first as unsupported coverage and the second as invalid
input. No generic map of names exists in shared code. The admission speaks of
one definition assignment per owned declaration. It never states that a site is
one entity of a native compiler, which C++ would contradict.

The part of the admission that knows no declaration is
[`files.Cover`](../files/go/files.go): file identities, paths and any
comparable item. The placement of existing TypeScript source uses it with
declaration sites. The horos formatter keeps its own copy of these rules until
it moves onto the operation.

**Order.** No placement states an order. `AdmittedPlacement` hands out the
files in the order of their identities and the declarations in the order of
their sites in the sealed output. That is the order the accessors hand out,
and TypeScript realizes it: two TypeScript placements with the same
assignments give byte-identical files and the same unit order. The shared side
does not enforce an order of publication. Each language owns the order that it
realizes, since a language may have to schedule a declaration before its use.

**The entry.** No placement names an entry. `semantic.Snapshot.Entry` and
`source.Source.Entry` are an address. The entry names no root of a native
program, such as the packages of a Go build, the root file of a Rust crate, a
translation-unit context of C++ or the input files of a TypeScript compiler, and
it limits no check. Every capture needs one, so each formatter states the rule
by which it chooses the entry of its fresh snapshot. TypeScript chooses the file
of least identity. The horos formatter keeps an entry with a meaning of its own,
the module of the root pin.

**The evidence requirement belongs to one request.** The contract has no fixed
evidence. The step states the requirement for the request that it prepares, and
the shared side matches it with the receipt of the publication, as it did
before. It learns no manifest type. For TypeScript, and later for Go, the step
returns the same requirement for every request.

**Why the shared contract changed** ([AC-08](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-08)).

- *The missing capability.* The contract fixed the required evidence when a
  formatter was constructed, and it handed the placement to the provider
  unread. A language whose evidence depends on the placement could not state
  it. The rules that every placement has to satisfy existed as private copies,
  and a shared admission that shared code does not run is a convention.
- *Why this layer.* The shared formatter creates the publisher of each request
  and matches the receipt with the requirement there. Only this layer can
  guarantee that a refused placement never reaches a provider and that a
  publication is matched against the admitted files.
- *How it serves independently owned contributions.* Two concrete cases. C++
  headers are checked in translation-unit contexts, and the contexts follow
  from the headers that the placement creates. A Rust crate is compiled from
  its root file, and the root follows from the placement's module tree. In both
  a receipt for the right bytes under the wrong compilation root has to be
  refused, and only a requirement that follows from the placement can refuse
  it. Each language supplies its own step, and no shared code names a language.

**The rule for guards.** A guard compares the published value with the admitted
consumer input: where a definition was published against where the placement
puts it, and what the language realized against what the output intends.
Running the layout function a second time and comparing its result with the
publication shows that the function is deterministic. It is an integrity check
and no more: a layout function that ignores a part of the placement agrees with
itself. The rule binds the Go and Rust guards from their first slice. The
TypeScript guard compares the publication with the layout that its step derived
from the admitted placement, and the compiler's bindings and exports with the
output's intent. Strengthening it and the horos guard follows as coverage.

`NewGeneratedFilesFormatter` takes the step and the provider as separate
arguments, so a test can hand a language's own guard a publication that the
language's layout does not produce. TypeScript publishes its step through
`Stages.GeneratedContract`.
[Tests with a fixture language](../integration/go/generated_files_contract_test.go) ·
[Tests with TypeScript](../integration/go/typescript_generated_placement_test.go)

**What a later change may not alter:**

- the shape of the TypeScript placement: two keyed tables, no entry, no order;
- the shape of the view: `GeneratedPlacementView{Files, Definitions}`, one
  file per declaration;
- the input has three fields. A request carries no setting, so everything a
  language needs for one request is part of its placement type;
- the protocol: the step calls `Admit` while it runs, once, and returns that
  value with the evidence requirement and the guard;
- what the common admission checks, its order, and the category of each of
  its refusals;
- a publication has exactly the files of the admitted table. A formatter
  cannot publish a file that the view does not declare, so a file that a
  language derives, such as a module file, is part of its view;
- a defect of a step is `format.ErrInvalidFormatter` and carries no category;
- `Limits.Source.Files.MaxLinks` bounds the number of files;
- `Snapshot.Entry` and `Source.Entry` are an address.

**Allowed later, and not built:** a field of `GeneratedFilesPrepared`, an
accessor of `AdmittedPlacement`, and an accessor of `GeneratedFilesRequest`
that carries to the provider what the step derived. Today the guard is a
closure and has it, while the provider receives the placement and the admitted
placement and derives the language's layout a second time.

## Carriers

The existing `ast.AST[L]`, `ast.View[L]` and `cst.CST[L]` remain the Syntax carriers.
The [Files package](../files/go/files.go) holds application-supplied paths and text:

```go
type ID string
type Path string
type File struct { ID ID; Path Path; Text string }
type Files struct { Entries []File }
type SyntaxLink struct { File ID; Tree ast.TreeID; Source ast.SourceID }
type Context struct { Files Files; Links []SyntaxLink }
```

File IDs are inventory identities, separate from paths and semantic declarations.
Paths are relative, slash-separated names. IDs and paths are unique; absolute,
drive/UNC, empty, dot, dot-dot, backslash and NUL paths are rejected. A file cannot
also be a parent directory. The application selects its filesystem root and owns
host-specific path equivalence, reserved names, symlinks and overwrite policy.
Providers do no disk or network IO. An editor, archive or cube can supply Files.

`semantic.Input.Files` optionally captures Context together with its ASTs and
explicit entry. Capture owns the inventory and one-to-one file/tree/source links
before provider callbacks. Changes to files, paths, selected dependencies or entry
require a new capture; changed configuration requires a new semantic invocation.
Matching addresses or bytes cannot substitute for exact receipt association.

The [semantic graph](semantics.md) is language-parameterized:

```go
type Semantics[L any] struct { /* private immutable storage */ }
type View[L any] interface { SemanticsView() Semantics[L] }
func (Semantics[L]) SemanticsView() Semantics[L]
func (Semantics[L]) Data() Data
func (Semantics[L]) Entity(EntityID) (Entity, bool)
func (Semantics[L]) Term(TermID) (Term, bool)
```

Entities, terms, bindings, expressions, dependencies and subject facts remain
provider-owned meaning. Applications have a **callee term**, so the kernel can
retain indirect and higher-order expressions without interpreting them.
`DirectParameterBindings` is available only for an application whose callee is a
direct entity reference with matching effective arity. Subject facts are
multivalued relations on entities or terms; keys do not register behavior.

`semantic.Result.Semantics()` returns a `representation.Semantics` envelope.
`representation.UnpackSemantics[L]` checks the language ID, explicit profile and
actual Go marker before returning the typed graph. There is no conversion to
`Semantics[any]`. Graph integrity alone does not establish language validity.

The [Source package](../source/go/source.go) retains the concrete realization and
its successful semantic result together:

```go
type Unit[L any] struct {
    File files.ID
    Path files.Path
    Source ast.SourceID
    Syntax cst.CST[L]
}
type GroupID string
type Group struct {
    ID GroupID
    Kind semantics.Key
    Members []files.ID
    Sites []ast.NodeRef
}
type Relation struct {
    Kind semantics.Key
    From, To files.ID
    At *ast.NodeRef
}
type Publication[L any] struct {
    Snapshot semantic.Snapshot
    Analysis semantic.Result
    Units []Unit[L]
    Groups []Group
    Relations []Relation
}
type Source[L any] struct { /* private sealed storage */ }
func (Source[L]) Entry() ast.NodeRef
func (Source[L]) Units() []Unit[L]
func (Source[L]) Groups() []Group
func (Source[L]) Relations() []Relation
func (Source[L]) Snapshot() semantic.Snapshot
func (Source[L]) Evidence() semantic.Result
func (Source[L]) SemanticsView() semantics.Semantics[L]
```

Both Source and Semantics implement `semantics.View[L]`, just as CST and AST
implement `ast.View[L]`. Source derives its graph from Analysis; there is no
independent mutable Abstract field. Accessors copy mutable storage and retain
immutable receipt/snapshot handles. Zero Source is invalid.

## Symmetric role contracts

```go
// contracts/parse/go
type Provider[I, O any] interface {
    language.Provider
    Parse(context.Context, I) (O, error)
}
// contracts/format/go
type Provider[I, O any] interface {
    language.Provider
    Format(context.Context, I) (O, error)
}
// contracts/print/go
type Provider[I, O any] interface {
    language.Provider
    Print(context.Context, I) (O, error)
}
```

| Role | Syntax instantiation | Files instantiation |
| --- | --- | --- |
| parse.Provider | `[parse.Input, cst.CST[L]]` | `[parse.FilesRequest[L], source.Source[L]]` |
| format.Provider | `[ast.View[L], cst.CST[L]]` | `[format.FilesRequest[L], source.Source[L]]` |
| print.Provider | `[cst.CST[L], string]` | `[source.Source[L], files.Files]` |

All expose ID and Name. `pipeline.Stage[I,O]` also accepts a context. Input/output
shapes do not discover or register capabilities. The adapters admit exactly the
configured language/profile/marker. Syntax helpers and provider methods accept
context too.

Syntax adapters are `parse.SyntaxParser`, `format.SyntaxFormatter` and
`print.SyntaxPrinter`, constructed by `NewSyntaxParser[L]`,
`NewSyntaxFormatter[L]` and `NewSyntaxPrinter[L]`. Files adapters are typed
`parse.FilesParser[L]`, `format.FilesFormatter[L]` and `print.FilesPrinter[L]`.
Their provider-facing requests expose copied input and invocation-scoped Publish:

```go
// parse
type FilesInput struct {
    Files files.Files
    Entry files.ID
    Configuration []semantic.Setting
    Limits FilesLimits
}
type FilesRequest[L any] struct { /* private captured state */ }
func (FilesRequest[L]) Input() FilesInput
func (FilesRequest[L]) Publish(source.Publication[L]) (source.Source[L], error)

// format
type FilesInput[L any] struct {
    Abstract semantics.Semantics[L]
    Snapshot semantic.Snapshot
    Evidence semantic.Result
    Configuration []semantic.Setting
    Limits FilesLimits
}
type FilesRequest[L any] struct { /* private captured state */ }
func (FilesRequest[L]) Input() FilesInput[L]
func (FilesRequest[L]) Publish(source.Publication[L]) (source.Source[L], error)
```

Files contracts explicitly supply their required analysis profile/reports, Syntax
printer and language-owned admission/preservation callback. The shared publisher
requires exact capture/receipt and graph association; complete unit coverage;
full AST equality with the capture; matching file IDs, paths and SourceIDs; and
concrete text matching captured Files under that printer. Groups and relations
refer to current units/sites. Their meaning and CST/AST agreement belong to the
language guard. Returning Source from a previous invocation is rejected. Failed
or canceled operations return no usable complete output.

Format admission requires **the same graph identity** as the successful input
receipt, not merely equal facts. Output has a fresh capture and receipt. Both
Parse→Format and Format→Format retain what the next stage needs:

```go
parsed, err := parser.Parse(ctx, parse.FilesInput{Files: input, Entry: "people"})
// Handle err.
formatted, err := formatter.Format(ctx, format.FilesInput[L]{
    Abstract: parsed.SemanticsView(),
    Snapshot: parsed.Snapshot(),
    Evidence: parsed.Evidence(),
    Configuration: organization,
})
// Handle err before printing.
output, err := printer.Print(ctx, formatted)
```

The engine independently registers `RegisterSyntaxParser`, `RegisterSyntaxFormatter`,
`RegisterSyntaxPrinter`, `RegisterFilesParser`, `RegisterFilesFormatter` and
`RegisterFilesPrinter`. Check, resolve and analyze remain separate. Syntax engine
methods are `ParseSyntax`, `FormatSyntax`, `PrintSyntax`, `RenderSyntax`; Files
methods are `ParseFiles`, `FormatFiles`, `PrintFiles`. Callbacks run outside locks.
No Files capability falls back to Syntax when semantics is unavailable.

The heterogeneous engine uses `source.Envelope`, sealed by `source.Pack[L]` and
opened with `source.Unpack[L]`. `format.PackFilesInput` retains the actual marker
in `format.FilesValue`. Typed Files adapters expose Parse/Format/Print directly;
their `FilesRegistration` bridges supply the engine's heterogeneous calls.
Envelope also exposes Snapshot, Evidence and the dynamic SemanticsView for
agnostic composition. No Source import is added to tree representation.

## Failure categories

**13 November 2026.** A Files formatter refuses an input in one of four
categories. A caller tells them apart with `errors.Is`.

| Category | Value | Meaning |
| --- | --- | --- |
| Invalid input | `source.ErrInvalid` | No later slice can admit the same input. |
| Unsupported coverage | `source.ErrUnsupported` | A later slice can admit the same input. |
| Missing evidence | `semantic.ErrEvidenceUnavailable` | The analysis that the formatter requires is not established. |
| Exceeded limit | `source.ErrLimit` | The input or the fresh output exceeds a limit of the request. |

The rule for the two that are easy to confuse: a refusal is unsupported coverage
when a later slice can admit the same input, and it is invalid input when no
later slice can. The same input is judged under the same vocabulary. An input
that is malformed under the published vocabulary of a setting or a placement
type is invalid: an unknown field, an unknown setting key, a repeated member, a
path that is no path of the language. An input that is well formed under that
vocabulary and outside what the implementation covers is unsupported.
Publishing a new vocabulary, under a new key or type, is the one change that
may admit an input that is malformed today, and an input under that vocabulary
is a different input.

The category of a refusal is part of the Files contract. A later change may turn
an unsupported refusal into success. It may not turn an invalid refusal into
success, it may not move a refusal from one category to another, and it may not
accept and ignore what a formatter refuses as malformed today.

A name use without a destination is invalid input, for every generated-Files
formatter: a sealed output with such a use is malformed, and no later slice
admits it. So is an external use that the language's issuer cannot read, and
one whose node is no name use or does not name the symbol. An external
destination that the formatter's evidence cannot establish is unsupported
coverage, because a later slice can admit it. Rust's `bool` was one until the
Rust helper bound primitive types, on 13 November 2026. [Destinations](destinations.md#who-enforces-completeness)
states the rule. The table below has its row for TypeScript, the one formatter
that exists.

An input can have several defects. The category describes the defect that was
found. An unsupported refusal says nothing about the rest of the input: the
same input may hold a defect that no slice admits. The order in which a
formatter runs its checks is no part of the contract.

Two orders are guaranteed: the order of the common admission, stated with
[the generated-Files contract](#the-generated-files-contract), and that validity
is decided before coverage. A placement that
has a defect of both kinds is invalid input, whichever file has which.

- **A generated placement.** The [common admission](#the-generated-files-contract)
  runs first: the number of files against its limits, the handles, the file of
  every assignment, the completeness of the placement and the inventory of the
  files. Then TypeScript decides what is valid: what the sealed output
  carries, a path that is no TypeScript path, two declarations of one name in
  one file and a use of an unexposed declaration of another file. Only then it
  refuses a path that the profile excludes or a file without a declaration.
- **A placement of existing TypeScript source.** The form of the setting is
  decided before the captured program is read: its key, its JSON, its fields
  and its members. Then the files and the top-level declarations of the
  program are read, and the placement is judged against them: the number of
  files, its items, its files, its inventory, and a path that is no TypeScript
  path. Only then the graph is read, where a script is refused. Then the two
  refusals of relocation, and at last a path that the profile excludes or a
  file without a declaration.

The two refusals of relocation, for two declarations of one name in one file
and for an unexposed declaration in another file than its user, are decided in
a pass of their own, before a tree is constructed and before a refusal of
coverage. One order is not reached. The relocation refusals need the bindings
of the program, and a script is refused while the graph is read. A script
whose placement separates an unexposed declaration from its user, or puts two
declarations of one name into one file, is therefore reported as a script.
An import form outside the profile and a capture without files are refused
earlier still, while the declarations are read and before the placement is
judged. The horos placement has no coverage refusal, so it needs no such
order.

A refusal that a formatter decides matches exactly one category. When the
refusal carries the value of another package, that value is joined to the
category and `errors.Is` holds for both. One function says which category a
value of the `files` package and a limit of the semantic contract carry:
`source.WithCategory`. A `files.ErrLimit` and a `semantic.ErrLimit` are an
exceeded limit, a `files.ErrInvalid` is invalid input, and every other error is
returned as it is. The publisher calls it for its inventory check, and each
formatter calls it for its own inventory check and for the capture of its fresh
output. The generated-Files contract calls it for the check of the sealed
output against the semantic limits and for the answer of `files.Cover`, and
the TypeScript placement of existing source calls it for the answer of
`files.Cover`. A symbol that the TypeScript issuer cannot read keeps its
`generation` value, and a handle that the output did not issue keeps
`generation.ErrScope`. The horos formatter keeps the JSON decoder's error and
`semantic.ErrInvalidRequest` for its setting.

No part of a placement is accepted and ignored. A generated placement is a
typed value with two keyed tables, so it cannot name a member twice, state an
order or carry a field that its type does not have. The two formatters that
read a placement from a setting, TypeScript for existing source and horos,
refuse an unknown field at every level, a second JSON value and a member that
is named twice, in any case of its name, as invalid input. A field that the
vocabulary has, with a value outside coverage, is unsupported.

| Refusal | TypeScript | horos |
| --- | --- | --- |
| Another setting key, a second setting, malformed JSON, an empty value, a value of another type, an unknown field, a second JSON value, a member named twice | invalid, for existing source. A generated placement is no setting | invalid |
| A placement of existing source under the earlier key `declaration-placement-1`, whatever it holds. The vocabulary without an entry is `declaration-placement-2` since [#273](https://github.com/Bitspark/bittype/issues/273), and the earlier key is another setting key | invalid | does not apply |
| The member `Entry` in a placement of existing source under `declaration-placement-2`. It is an unknown field | invalid | does not apply |
| A request for generated declarations without a placement, or with a placement that another language issued | invalid | does not apply |
| A handle that the sealed output of the request did not issue, also the handle of another run of the same plan | invalid | does not apply |
| A file identity or a path used twice, an empty file identity, a path that is also a directory of another, a malformed path | invalid | invalid |
| A valid path that is no TypeScript source path, such as `person.txt` or `PERSON.TS` | invalid | does not apply |
| A path that TypeScript reads and this profile excludes: `.d.ts`, `.mts`, `.cts`, `.tsx`, or `#`, `%` or `?` in a TypeScript path | unsupported | does not apply |
| A file without an item, in a placement or among the captured files | unsupported | invalid |
| No file, an unknown, repeated or missing item, an item in a file that the placement does not declare. | invalid | invalid |
| Two declarations of one name in one file, a binding without a single target | invalid | does not apply |
| A declaration placed in another file than an unexposed declaration that it uses | invalid | does not apply |
| A generated tree with an import or an exported declaration, an exposed site that is no top-level declaration, a name use without a destination, an external use that the TypeScript issuer cannot read or that its site does not name | invalid | does not apply |
| An import form other than a named type import | unsupported | does not apply |
| A file that is a script, in the input or among the placed files | unsupported | does not apply |
| A capture without files | unsupported | admitted |
| The fresh output is not admitted by its analysis | missing evidence | missing evidence |
| Too many files, file identities and paths over the byte limit, placed trees over the node, depth or work limit | limit | limit |
| More files than the link limit of the files | limit | not applied |
| A sealed output over a semantic limit of the request | limit | does not apply |
| Fresh files over the Files byte limit, fresh paths over the input byte limit | limit | limit |
| Constructed imports over the node or byte limit | limit | does not apply |
| Fresh trees over the input node limit, fresh text over the input byte limit | limit | limit |
| More fresh units than the input limit, a capture of the fresh output over a semantic limit | limit | limit |
| Another error of the capture of the fresh output, such as the refusal of its Syntax checker | no category | no category |

Three rows need their reason.

- **A file without an item.** A TypeScript file without a declaration is well
  formed, and row 30a of [research 0016](../research-docs/0016-native-files-placement.md)
  defers the form that would admit it. A module that holds only imports is such
  a file, also in its original placement. A horos file without a module cannot
  exist: the horos grammar refuses a text without a module declaration
  (`parse/syntax/horos/provider/go/parse.go:77-79`).
- **TypeScript paths.** [The TypeScript document](typescript-type-imports.md)
  states the excluded extensions and characters as limits of its profile. A path
  that TypeScript does not read as a source file is no path of the language.
  `PERSON.TS` is one: the host that the bridge gives the compiler is
  case-sensitive (`internal/typescriptsem/compiler.cjs:78-79`), and the compiler
  reads no file with that extension.
- **An unexposed declaration in another file than its user.** It was examined
  and is invalid input. Files never exports a declaration without
  authorization, so no later slice can admit that placement for the same
  output. A caller who wants it has to change the request, which makes it a
  different input.

One [table test](../integration/go/files_failure_categories_test.go) feeds each
of these to the formatter that can receive it. The TypeScript bridge refuses a
namespace import and a capture without files before the formatter reads them,
and never publishes a binding without a single target, so a
[package test](../format/files/typescript/provider/go/placement_test.go) covers
those three.

The rule does not reach the following places. Each returns a value that the
table above does not give it:

- **The shared parse contract and the shared format contract for existing
  source.** A receipt or a graph of another snapshot and a setting key that is
  given twice return `semantic.ErrInvalidRequest`. An input over a semantic
  limit and a setting over the configuration limit return `semantic.ErrLimit`.
  That contract counts the input, so a limit of the request that the input
  exceeds carries no category, while the same limit carries the limit category
  when the capture of the fresh output exceeds it. The shared parse contract
  returns `files.ErrInvalid` and `files.ErrLimit` without a category. The
  generated-Files contract is reached since
  [#273](https://github.com/Bitspark/bittype/issues/273): a sealed output over
  a semantic limit is an exceeded limit, and every refusal of the common
  admission has its category from `source.WithCategory`. The existing-source
  contract was not changed with it.
- **The semantic stages.** A fresh output whose resolution, analysis or graph
  exceeds a limit returns the stage's own value, `semantic.ErrLimit` or
  `semantics.ErrLimit`, without a category. Only the capture of the fresh
  output is joined. That is pinned for the TypeScript resolver's graph limit
  and for the horos resolver's configuration limit. It is not established for
  a fresh horos graph over the graph limit: in the one attempt made on 13
  November 2026 that came back as missing evidence, and no row holds it.
- **The publisher.** It joins the `files` value of its inventory check. It
  passes on the values of the evidence admission, of the semantic and graph
  limits, of the representation and of the Syntax printer as they are.
- **`Limits.Source.Files.MaxLinks` in the horos formatter.** Each file of a
  placement is published with one syntax link, so `files.Cover` applies the
  link limit to the number of files, for a generated placement and for a
  placement of existing TypeScript source. The horos formatter is not on that
  operation: it validates the limit and never applies it. That ends when horos
  moves onto `files.Cover`.
- **The constructors.** A formatter that is constructed without the environment
  it needs returns `source.ErrInvalid` for TypeScript and, for horos, also
  values of the format, semantic and representation contracts. No rule says
  which of these a later slice could admit.
- **A defect of a preparation step.** A step that returns without the admitted
  placement of its request, without a guard or with a requirement that asks
  for no semantic graph fails the request with `format.ErrInvalidFormatter`,
  and a malformed requirement fails it with the value of the semantic
  contract. These are defects of a language's contribution and no category of
  the input. A second call of `Admit`, and a call after the step returned,
  answer `format.ErrInvalidFormatter` too. The refusal of a formatter without
  semantic evidence moved here
  from the constructor, because the requirement belongs to one request.
- **The publication guards.** A guard that finds the fresh output different
  from the placement returns `source.ErrInvalid`, and so does the shared
  comparison of the published files with the admitted files table. A receipt
  that does not match the requirement of its request returns
  `semantic.ErrInvalidRequest` from the publisher.
- **A type parameter that captures a declaration placed in its file.** The
  bridge reports the changed binding of the fresh output, and the caller
  receives `semantic.ErrInvalidRequest`.
  [#273](https://github.com/Bitspark/bittype/issues/273) decides that TypeScript
  Files refuses this placement as unsupported coverage. That is part of its
  second change, with the removal of the flat-name rule: until then generation
  refuses the declarations that make such a placement.
- **The horos Files parser.** A wrong entry setting returns
  `semantic.ErrInvalidRequest` and an input over the node limit returns
  `semantic.ErrLimit`.

A namespace import does not reach the line that refuses it as unsupported
coverage. The TypeScript bridge reads the file first, so the caller receives
missing evidence, and the diagnostic says that the syntax is outside the
profile.

## Horos composition and pagination

Files parsing invokes Syntax Parse for each file, retains distinct CSTs and
tree-qualified sites, captures their shared context, resolves the explicit entry
against the canonical module pool, analyzes that exact result, and publishes
Source. Horos now resolves several Program trees without concatenating text.
Imports select canonical pins; filenames do not grant visibility. Duplicate pins,
contradictory assertions, missing dependencies and rejected/incomplete semantics
cannot publish Source. The Files profile rejects supplied modules outside the
entry's reachable closure. A file may contain several whole modules; files cannot
implicitly contribute fragments of a canonical module.

The entry file is mandatory. Within it the first module section is selected by
the Syntax grammar, or `parse/files/horos.EntryModuleSetting()` explicitly selects
a supplied module pin. File order never selects the entry. A pagination library
and its consumer can therefore stay in separate files:

```text
lib/page.horos:  module example.pagination;
                type Page<T> = { items: list<T>; next?: text; };

people.horos:   module example.people;
                import pages from example.pagination pin "sha256:<actual-pin>";
                type Person = { name: text; };
                type People = pages.Page<Person>;
```

Files formatting uses the admitted graph and its captured syntax representation,
reconstructs fresh per-unit ASTs, runs Syntax Format directly, and explicitly
resolves/analyzes the resulting collection. It preserves the selected canonical
root and complete set of module identities. Current semantic sites are
re-established in the new trees. Old annotations, source spans and receipts are
not copied onto relocated syntax. No raw-source escape or cross-language type
choice is introduced. Its language guard compares canonical modules with the
stored text; that admission step is separate from AST-to-CST construction.

Default formatting preserves captured file grouping and paths. For a syntax-only
capture without Files context, the Horos formatting profile chooses one module
per `modules/<pin>.horos`. To choose another organization, pass
`format/files/horos.OrganizationSetting()` with UTF-8 JSON bytes:

```json
[{"ID":"bundle","Path":"generated/types.horos","Modules":["<people-pin>","<page-pin>"]}]
```

The list must assign every supplied canonical module exactly once. IDs/paths
must pass Files admission. Unknown settings and incomplete placements are rejected; ambient
filesystem input is never inferred. Pins here are canonical hex strings without
the source syntax's `sha256:` prefix.

Files printing invokes only the chosen Syntax printer and returns each recorded
path/text. Parse→Print preserves original bytes and trivia. Formatting can change
organization and spelling; semantic preservation does not require old tree IDs.

The caller supplies a language descriptor and the providers needed by the requested stage. For example, this constructs the Syntax printer and uses it to supply the Files printer.

```go
p, err := syntaxprint.New[horos.Language](horos.Language{})
if err != nil { return err }
syntaxPrinter, err := print.NewSyntaxPrinter[horos.Language](definition, p)
if err != nil { return err }
fp, err := filesprint.New[horos.Language](horos.Language{}, p)
if err != nil { return err }
filesPrinter, err := print.NewFilesPrinter[horos.Language](definition, fp)
if err != nil { return err }
```

The Horos Files parser and formatter accept explicit syntax/check/resolution/analysis
collaborators and semantic options. Shared printing at both dimensions uses
`Provider[L]` and `New[L]`; the Files constructor additionally takes the selected
Syntax printer. A native Files printer requires no semantic provider.
Go, TypeScript and Rust supply [bounded Files resolution and analysis](native-semantics.md).
The [TypeScript Files formatter](typescript-type-imports.md) organizes admitted native
declarations into modules. Its single placement supplies output paths and determines
local type imports. The [provider matrix](providers.md) lists the available stages.

## Bounds

Each parse/format `FilesLimits` combines `source.Limits` and `semantic.Limits`.
Zero selects defaults; partly invalid records are rejected. Files defaults allow
256 entries/links and 16 MiB of text/path/ID bytes. Source defaults allow 400,000
CST nodes, 1,600,000 links, 16 MiB concrete/organization metadata, depth 512 and
4,000,000 traversal steps. Semantic defaults separately bound input trees/nodes/
bytes, configuration, diagnostics and graph storage. Tighter invocation limits
also apply to retained capture and graph evidence. Capture bounds storage before
copying; publication bounds it before sealing. Language parsing and canonical
resolution retain their own bounded profiles.

## Multi-unit transformations and dependencies

`ASTTransformer.Transform` returns `representation.ASTSet{Entry, Trees}`.
`transform.NewTransformer[F,T]` takes a typed callback returning
`ast.Set[T]{Entry, Trees}`. The complete set has one language/profile/marker,
unique tree IDs and an explicit valid entry. One-unit mappings use one-element
sets. Each engine route leg admits every output tree, captures the whole collection
afresh and discards preceding semantic receipts. New evidence is explicit;
the engine never discovers an analyzer. No unit is silently dropped.

Generation returns a checked snapshot and binding intent. TypeScript Files takes that output and establishes native semantics for the organized files. Syntax formatting and printing are also available for individual trees.

The [stage-first layout](../LAYOUT.md) uses explicit dimensions:
`parse/syntax/<language>/provider/go` and `parse/files/<language>/provider/go`, likewise for
formatting. Shared printers live at `print/syntax/shared/provider/go` and `print/files/shared/provider/go`.
Files resolution and analysis live at `resolve/files/<language>/provider/go` and
`analyze/files/<language>/provider/go`. Only implemented capabilities have directories.
Files depends on AST addresses; Semantics remains independent of Files.
Source depends on the semantic capture/result contract to retain receipts.
Semantic capture may depend on Files but never Source. Tree representation admits
ASTSet and typed Semantics but never imports Source. Role contracts depend on these
carriers, not concrete languages; the engine remains agnostic.

Superseded single-parameter role contracts, unqualified Syntax adapter/engine
names and `semantics.Model` are removed without aliases. No runtime bindings,
behavior, serialization, network transport or implicit delivery workflow is added.
Function, Stream and Cell remain ordinary library contracts.
