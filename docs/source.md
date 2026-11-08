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
