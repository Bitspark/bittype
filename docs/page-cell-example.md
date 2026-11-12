# Page and Cell: runnable consumer walkthrough

The complete executable example is
[`TestTypeScriptTypeImportPlacement`](../integration/go/typescript_placement_test.go).
It generates once, formats two placements, checks a conflicting local Cell name,
prints the files and compiles an independent TypeScript consumer. This guide
walks that maintained fixture rather than introducing a second implementation.
Go excerpts below show individual calling steps; fixture helpers and imports are
linked explicitly, so the excerpts are not presented as a standalone Go program.

## Run the complete example

On a configured development/CI worker with the repository's Go toolchain and Node,
provision the same TypeScript version used by the existing Go contracts workflow:

```sh
npm install --prefix "$RUNNER_TEMP/bittype-typescript" --ignore-scripts --no-audit --no-fund typescript@5.9.3 @types/node@24.0.0
export BITTYPE_TYPESCRIPT_NODE="$(command -v node)"
export BITTYPE_TYPESCRIPT_MODULE="$RUNNER_TEMP/bittype-typescript/node_modules/typescript/lib/typescript.js"
go test -count=1 -v -run '^TestTypeScriptTypeImportPlacement$' ./integration/go
```

The test requires the actual compiler; an unavailable compiler is not a passing
result. The [existing workflow](../.github/workflows/go.yml) provisions it and
runs this test in the contracts suite. This command needs no Rust or C++ helper.

## 1. Supply source declarations and their identity

The fixture constructs canonical Horos module values through ordinary declaration
constructors. Their readable equivalent is:

```horos
module cells;
type Cell<T> = wire<T>;

module consumer;
import cells from cells pin "sha256:<computed-cell-module-pin>";
type Cell = text;
type Page<T> = { items: list<T>; next?: text; };
type Person = { name: text; active: bool; };
type Cells = Page<cells.Cell<Person>>;
```

The placeholder is explanatory; the runnable fixture calls Horos `decl.Pin` on
the actual library module and uses that pin both in its import and its mapping
selector. No file path or same-spelled name supplies canonical identity. The
local `Cell = text` deliberately conflicts with the external export's spelling.
The library's `wire<T>` body is admitted at source level but is not structurally
mapped: its exact native mapping selects the external Cell type.

Cell here is a type-only fixture. Its TypeScript declaration is
`export interface Cell<T> { value: T }`; it implements no live cell behavior.

## 2. Capture, resolve and analyze explicitly

[`generationInput` and `generationSource`](../integration/go/horos_generation_pagination_test.go)
project the canonical modules to a Horos AST, register the Horos language and
independent checker/resolver/analyzer, then obtain the exact capture and evidence.
Implemented public call excerpts, after those registrations:

```go
snapshot, err := e.Capture(ctx, semantic.Input{
    Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root},
    Trees: []representation.AST{envelope},
})
if err != nil { return err }
resolved, err := e.Resolve(ctx, snapshot, semantic.Options{})
if err != nil { return err }
analyzed, err := e.Analyze(ctx, analysis.Input{
    Snapshot: snapshot, Resolution: &resolved,
})
if err != nil { return err }
input := transform.Input{Source: snapshot, Evidence: []semantic.Result{analyzed}}
```

The source registration implements `SourceProvider` and the language stages
implement the independently registered resolution/analysis contracts. Their exact
SPIs are in [generation source admission](generation-spi.md#semantic-traversal-and-source-interpretation),
[resolution](resolution.md) and [analysis](analysis.md). The engine invokes the
selected stages; the generation reader admits their receipts. A nil error alone
does not establish complete semantics; the receiving generation contract checks
the required reports and model against this snapshot.

The established graph contains separate entities for Page, Page.T, the two Cell
declarations and Person. Cells' body applies a Reference to Page to an application
of the imported Cell Reference to Person. The open Page body refers to its own T.
Neither name equality nor expansion of Page supplies those bindings.

## 3. Register representations

The library implements the following public contract; registration invokes it:

```go
type MappingLibrary[F, T, R any] interface {
    Describe(from, to representation.Language) (MappingDescription[F,T,R], error)
}
```

[`cellMapping`](../integration/go/horos_generation_composition_test.go) supplies
this exact source selector and native expression. Here source is the Horos
generation integration, model is the TypeScript native model, h/ts are language
markers, and g is the generation contract package. This excerpt omits only the
fixture's callback invocation counter and enclosing Describe method:

```go
key, err := source.Named(from, pin, "Cell")
if err != nil { return g.MappingDescription[h.Language, ts.Language, source.Reader]{}, err }
slot := g.TypeSlot("V")
pattern := g.Pattern[h.Language, source.Reader]{
    Selector: key,
    Slots: []g.Slot{slot},
    Capture: func(ctx context.Context, src source.Reader, use g.UseView[h.Language]) ([]g.Capture[h.Language], error) {
        for _, a := range use.Actuals { // select Cell's own binder T, never a position
            if p, err := src.Parameter(a.Binder); err != nil {
                return nil, err
            } else if p.Name == "T" {
                return []g.Capture[h.Language]{{Slot: slot, Use: a.Use}}, nil
            }
        }
        return nil, g.ErrScope
    },
}
binding := g.BindingDescription[h.Language, ts.Language, source.Reader]{
    ID: "cells/typescript/1", Input: pattern,
    Output: model.External("@bitspark/cell", "Cell", g.Mapped[ts.Language](slot)),
}
```

Describe returns that binding in a MappingDescription that states the
`structural-representation/1` profile. The kernel pairs Cell's binder with its
actual from the graph, judges that actual, and accepts the capture only as a
child it minted for this use. Cell's body is never opened. The mapping chooses a
public symbol; it cannot build an import or invoke an AST builder.

[`selectedTS`](../integration/go/horos_generation_composition_test.go) and the
[request helper](../integration/go/testdata/horosgeneration/request.go) register
the source, target, builtin mappings, declaration policy and Cell contribution.
The same operations, using the full builtin library set, are:

```go
registry := e.Generation()
err := compose.RegisterSource[h.Language, source.Reader](registry, from, source.Provider{})
if err != nil { return err }
err = compose.RegisterTarget(registry, to, target.Provider{})
if err != nil { return err }
for _, library := range []g.MappingLibrary[h.Language, ts.Language, source.Reader]{
    text.Library(), boolean.Library(), list.Library(),
    record.Library(), declarations.Library(),
} {
    if _, err := compose.RegisterMapping(registry, from, to, library); err != nil { return err }
}
if _, err := compose.RegisterMapping(registry, from, to, cellLibrary); err != nil { return err }
```

`text`, `boolean`, `list`, `record` and `declarations` are the five packages
under `transform/syntax/horos-to-typescript/mapping/<name>/go`. Each exports
`Library()`, and the consumer names each one it registers. A library states
its profile itself, as the Cell contribution does. The request below names the
profile it wants.

The [registration reference](kernel-composition.md#registration-and-request)
shows complete signatures. The fixture requests Page, Person, Cells and local
Cell by their source declaration sites. An exact external-use mapping for imported
Cell does not request publication of that library declaration.

## 4. Generate once

With roots from that capture, the application calls:

```go
output, err := e.Generate(ctx, compose.Request{
    Input: input, To: to, Roots: roots, Output: "selected-ts",
    Profile: rules.RepresentationProfile, Limits: g.DefaultLimits(),
})
if err != nil { return err }
```

Generate calls PlanGeneration and runs the resulting plan. The planner selects
and composes native expressions; the target then constructs ASTs without calling
the mapping again. In the resulting Output, four sites are exposed and one
external occurrence identifies `@bitspark/cell`'s Cell. The TypeScript tree has
no imports or exported modifiers yet. The fixture asserts those properties.

The target SPI and output accessors are documented in
[generation-spi.md](generation-spi.md#target-construction). The kernel walks the
graph for this example; [GEN-01](generation-gaps.md#gen-01) records that
correction and its scope.

## 5. Supply the native environment and one placement

[`placementStages`](../integration/go/typescript_placement_test.go) composes the
independent TypeScript Syntax formatter/printer/checker and native resolver/analyzer.
It calls `ImportConfiguration` with the supplied Cell declaration text and uses
the same settings for native resolution and analysis. Node and the exact
TypeScript compiler path come from the two environment variables above.

The implemented Files configuration contracts are:

```go
type Placement struct {
    File         files.ID
    Path         files.Path
    Declarations []ast.NodeRef
}
type Organization struct {
    Entry files.ID
    Units []Placement
}
func (o Organization) Setting() (semantic.Setting, error)
func NewGenerated(stages Stages) (format.GeneratedFilesFormatter[typescript.Language], error)
```

The application supplies Organization. The TypeScript Files implementation
consumes it. Its [Stages definition](../format/files/typescript/provider/go/format.go)
names the independently supplied collaborators; there is no combined native
language provider or ambient package discovery. A package's public declaration
text is a dependency input, not a second table of imports.

After resolving the emitted declaration sites from output, the fixture uses:

```go
page := tsfiles.Placement{
    File: "pages", Path: "generated/pages/page.ts",
    Declarations: []ast.NodeRef{sites["Page"], sites["Cells"], sites["Cell"]},
}
person := tsfiles.Placement{
    File: "model", Path: "generated/model/person.ts",
    Declarations: []ast.NodeRef{sites["Person"]},
}
```

For the split placement Organization contains page and person. For the co-located
placement it contains page with Person added to its declaration sites. Each
iteration passes the **same** output to the Files formatter:

```go
setting, err := organization.Setting()
if err != nil { return err }
placed, err := formatter.Format(ctx, format.GeneratedFilesInput{
    Output: output, Configuration: []semantic.Setting{setting},
})
if err != nil { return err }
emitted, err := printer.Print(ctx, placed)
if err != nil { return err }
```

The complete fixture also invokes formatting through
`Engine.RegisterGeneratedFilesFormatter` / `Engine.FormatGeneratedFiles`, then
unpacks the returned Source envelope. The shared Files printer is constructed
from `print/files/shared/provider/go.New` and the Syntax printer, then adapted
with `contracts/print/go.NewFilesPrinter`. It returns paths and text; the test
harness writes them to its temporary consumer directory.

## 6. Check the realized bindings

The split result contains these equivalent declarations:

```typescript
import type { Person } from "../model/person.js";
import type { Cell as Cell_1 } from "@bitspark/cell";

export type Cell = string;
export interface Page<T> { items: T[]; next?: string; }
export type Cells = Page<Cell_1<Person>>;
```

After co-location, the relative Person import disappears and Person is declared
in this file. The package import still uses Cell_1; it must never bind to local
Cell. Files establishes fresh native evidence for each organized capture and
checks declaration, binder, external-symbol and exposure destinations before
publishing Source.

The runnable test checks both binding graphs, prints and compiles both layouts,
reformats the admitted Source, rejects a stale receipt and rejects a dependency
that exports Other instead of Cell. Compiler success is additional evidence;
the binding comparisons establish which declarations the uses denote.

This qualifies representation and placement. It does not enforce Horos kind
bounds on arbitrary native instantiations, prove arbitrary mapping laws, or
implement Cell runtime behavior. Those boundaries remain explicit in
[GEN-02](generation-gaps.md#gen-02) and [verification](verification.md).
