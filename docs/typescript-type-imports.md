# TypeScript type imports and file placement

**Implemented, 8 November 2026.** A mapping selects a public TypeScript type. Generation records its destination and requested exposures in a checked output. Files owns placement, import aliases and exports. [Generation output](generation-planner.md#the-result-of-run) · [Architecture case](architecture-cases/0004-generation-file-syntax.md)

## From generation to files

1. The library supplies a declarative external-type expression with the public module, export and mapped arguments. The native target constructs the reference and records its external occurrence. Local uses retain their intended declaration or binder. No mapping receives an AST builder.
2. `Engine.Generate` returns a `generation.Output`. Its syntax snapshot, exposed sites and external uses belong to one run. The consumer supplies one placement and configures the target dependency environment.
3. Files chooses local names and imports, applies the requested visibility, and captures the organized units. It verifies the compiler’s binding at each relocated use and the module’s named exports before publishing `Source[TypeScript]`.

```go
// Mapping expression excerpt: model is transform/syntax/typescript/model/go.
return model.External("@bitspark/cell", "Cell", arg)
```

The [complete consumer walkthrough](page-cell-example.md) links the executable
fixture, exact compiler setup, registration and both placements. The
[other native placement proposals](native-file-placement.md) have separate status.

## Consumer API

Supply the selected packages and their declaration text through `semantics/typescript.ImportConfiguration` when constructing the TypeScript stages. In this example `page`, `cells` and `person` are handles from `out.Inventory()`, the list of the declarations that the output owns.

```go
out, err := eng.Generate(ctx, request)
if err != nil { return err }
placement := tsfiles.GeneratedPlacement{
    Files: map[files.ID]files.Path{
        "pages": "generated/pages/page.ts",
        "model": "generated/model/person.ts",
    },
    Definitions: map[generation.DeclarationHandle]files.ID{
        page: "pages", cells: "pages", person: "model",
    },
}
formatter, err := tsfiles.NewGenerated(stages)
if err != nil { return err }
organized, err := formatter.Format(ctx, format.GeneratedFilesInput{
    Output: out, Placement: placement.Issue(),
})
if err != nil { return err }
emitted, err := filesPrinter.Print(ctx, organized)
```

The placement is a typed value with two keyed tables: each file identity with its path, and each declaration with the file that holds it. It travels as no setting and as no bytes. A handle carries its output, so a placement that was built for one run is refused with another run of the same plan. The [shared contract](source.md#the-generated-files-contract) admits what every placement has in common before TypeScript reads the placement: every declaration of the output has one file, every file is declared, and the files are a valid inventory.

The placement states no order and names no entry. TypeScript Files owns the order: the files in the order of their identities, and the declarations of one file in the order of their sites in the sealed output. Two placements with the same assignments give byte-identical files. The entry of the fresh snapshot is the file of least identity. That entry is an address: it names no root of a TypeScript program and limits no check.

The engine also accepts this formatter through `RegisterGeneratedFilesFormatter` and `FormatGeneratedFiles`, and carries the placement unread. For an existing analyzed program, `tsfiles.New` accepts the graph, its snapshot and its analysis receipt. Its placement, `tsfiles.Organization`, names sites of that snapshot and travels as the setting `declaration-placement-1`; it names no entry either, and the same rule chooses the entry. The two operations share the relocation and the fresh native checking. [Generated Files SPI](../contracts/format/go/generated.go) · [Executable example](../integration/go/typescript_placement_test.go) · [Order, entry and guard tests](../integration/go/typescript_generated_placement_test.go)

## Example and scope

Assign the generated TypeScript declarations Page and Cells to file ID `pages` at `generated/pages/page.ts`, and Person to `model` at `generated/model/person.ts`. The generated ESM declarations include the following types.

```typescript
import type { Person } from "../model/person.js";
import type { Cell } from "@bitspark/cell";

export interface Page<T> { items: T[]; next?: string; }
export type Cells = Page<Cell<Person>>;
```

<!-- bittype:overview scope=example-and-scope part=features -->
Move Person into `pages` in that placement and remove the empty model unit. The pages file now defines Person locally. Its import disappears and the Cell package import remains. Changing the model file's path changes the relative specifier. The returned Source records a local import relation at its actual import node for each owned module dependency. External packages retain their public specifiers and are supplied by the consuming build.
<!-- bittype:overview:end -->

<!-- bittype:overview scope=example-and-scope part=exclusions -->
This profile covers named type imports in nonempty, declaration-only TypeScript ESM modules. Output paths are Files-relative `.ts` paths, excluding `.d.ts` and URL delimiters `?`, `#` and `%`. Imports between output units use relative `.js` specifiers. Package declaration inputs support type aliases, interfaces and named type imports. Namespace imports, re-exports, runtime imports, ambient scripts and discovery of package files on disk are outside this profile.
<!-- bittype:overview:end -->

An incomplete placement, an unexported cross-file target, a missing dependency, or a changed binding fails the operation. Package declarations must preserve native diagnostics. Named imports are deduplicated and aliased per file when required. Both arrangements compile against a standalone package fixture that exports Cell. When the output module also declares a local `Cell`, Files chooses the alias `Cell_1` for the external use and checks it against the package export. [Integration tests](../integration/go/typescript_placement_test.go)

The Files stage coverage test registers this formatter alongside the Horos parser and formatter. [Bittype architectural case ACX-003](architecture-cases/0003-files-stage-coverage.md)
