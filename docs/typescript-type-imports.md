# TypeScript type imports and file placement

**Implemented, 7 November 2026.** A representation library selects a public TypeScript module and named type export. The TypeScript Files formatter owns local file organization. Consumers supply one placement for owned declarations. That placement determines both the files returned by the printer and the relative type imports between them. [Syntax and Files](source.md) · [Generation composition](kernel-composition.md)

## From generation to files

1. A library opens against the registered TypeScript builder and requests an imported type by module specifier, exported identifier, represented type arguments and source origin. The target constructs an ordinary type application. It emits named type imports, binds each use to the real import binding and allocates local aliases deterministically. A request can reuse an imported symbol. Generation produces a logical TypeScript AST, with true local declaration and binder links and package references that survive formatting.
2. Capture the logical TypeScript unit with its text, checker and explicitly configured TypeScript resolution and analysis. Supply package declaration text under the exact public module specifiers chosen by the build. The compiler can consult captured source units, selected package exports and its pinned standard library. Dependency declarations belong to this semantic environment and are separate from files the application emits.
3. Call the existing TypeScript Files formatter with the established TypeScript graph, exact snapshot, analysis receipt and organization setting. An organization selects an entry file and assigns captured top-level declaration sites to file IDs and relative paths. Each declaration appears exactly once. The formatter preserves existing grouping when the caller supplies no organization. Declaration sites are relative to this input capture. Paths use the validation rules of `files.Files`.
4. For each use, follow its resolved symbol. A declaration placed in the same unit uses a local reference. An exported declaration placed elsewhere uses a named type import relative to the importing file. A selected dependency symbol uses its public package specifier. Allocate aliases per unit, avoiding lexical parameters and declarations. Coalesce repeated imports, rebuild reference links in the new ASTs, and format each AST directly to a CST.
5. Establish fresh semantic evidence for the organized units and check the relocated bindings against their intended destinations. Publish `source.Source` with the units, paths and local import relations. The shared Files printer emits these paths and texts. Failed placement, ambiguity, missing dependencies or failed semantic admission produce no partial Source.

The source language supplies the canonical Cell identity. The Cell mapping selects export `Cell` from `@bitspark/cell` and passes its represented argument to the TypeScript builder.

```go
return builder.ImportedType("@bitspark/cell", "Cell", []generation.Fragment[typescript.Language]{arg}, origin)
```

The formatter reads each local use's compiler-resolved declaration site and each external use's package/export identity. The generation kernel and the shared Files printer operate on their existing contracts. [TypeScript builder](../transform/syntax/typescript/target/go/provider.go) · [Organization and relocation](../format/files/typescript/provider/go/placement.go)

## Consumer API

The build supplies the public package specifiers and their selected declaration text to `semantics/typescript.ImportConfiguration`. Pass those settings to `analyze.ContractFor` and to `format/files/typescript.Stages.Configuration`. Construct the stages with the TypeScript language, syntax formatter, printer, checker, resolver and analyzer. Format the generated AST into a logical CST, then call `stages.Capture(ctx, logicalUnits, entry, limits)`. The returned Files input owns the current TypeScript sites used below. The logical path supplies the initial compiler context. It does not prescribe the emitted path.

```go
organization := tsfiles.Organization{
    Entry: "pages",
    Units: []tsfiles.Placement{
        {File: "pages", Path: "generated/pages/page.ts", Declarations: []ast.NodeRef{page, cells}},
        {File: "model", Path: "generated/model/person.ts", Declarations: []ast.NodeRef{person}},
    },
}
setting, err := organization.Setting()
if err != nil { return err }
input.Configuration = []semantic.Setting{setting}
formatter, err := tsfiles.New(stages)
if err != nil { return err }
organized, err := formatter.Format(ctx, input)
if err != nil { return err }
emitted, err := filesPrinter.Print(ctx, organized)
```

`emitted.Entries` contains the output file IDs, paths and text. When using an already organized Source as input, use declaration sites in its new snapshot. Omitting the setting preserves the captured grouping and paths. The integration test includes the complete construction, both placement requests, native scope checks and compilation of a consuming TypeScript file. [Files provider](../format/files/typescript/provider/go/format.go) · [Executable example](../integration/go/typescript_placement_test.go)

## Example and scope

Assign the captured TypeScript declarations Page and Cells to file ID `pages` at `generated/pages/page.ts`, and Person to `model` at `generated/model/person.ts`. The generated ESM declarations include the following types.

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

An incomplete placement, an unexported cross-file target, a missing dependency, or a changed binding fails the operation. Package declarations must preserve native diagnostics. Named imports are deduplicated and aliased per file when required. Both arrangements compile against a standalone package fixture that exports Cell. [Integration tests](../integration/go/typescript_placement_test.go)

The Files stage coverage test registers this formatter alongside the Horos parser and formatter. [Bittype architectural case ACX-003](architecture-cases/0003-files-stage-coverage.md)
