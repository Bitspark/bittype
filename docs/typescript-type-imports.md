# TypeScript type imports and file placement

**Design, 7 November 2026.** A representation library selects a public TypeScript module and named type export. The TypeScript Files formatter owns local file organization. Consumers supply one placement for owned declarations. That placement determines both the files returned by the printer and the relative type imports between them. See [Syntax and Files](source.md) and [generation composition](kernel-composition.md).

## From generation to files

1. A library opens against the registered TypeScript builder and requests an imported type by module specifier, exported identifier, represented type arguments and source origin. The target constructs an ordinary type application. It emits named type imports, binds each use to the real import binding and allocates local aliases deterministically. A request can reuse an imported symbol. Generation produces a logical TypeScript AST, with true local declaration and binder links and package references that survive formatting.
2. Capture the logical TypeScript unit with its text, checker and explicitly configured TypeScript resolution and analysis. Supply package declaration text under the exact public module specifiers chosen by the build. The compiler can consult captured source units, selected package exports and its pinned standard library. Dependency declarations belong to this semantic environment and are separate from files the application emits.
3. Call the existing TypeScript Files formatter with the established TypeScript graph, exact snapshot, analysis receipt and organization setting. An organization selects an entry file and assigns captured top-level declaration sites to file IDs and relative paths. Each declaration appears exactly once. The formatter preserves existing grouping when the caller supplies no organization. Declaration sites are relative to this input capture. Paths use the validation rules of `files.Files`.
4. For each use, follow its resolved symbol. A declaration placed in the same unit uses a local reference. An exported declaration placed elsewhere uses a named type import relative to the importing file. A selected dependency symbol uses its public package specifier. Allocate aliases per unit, avoiding lexical parameters and declarations. Coalesce repeated imports, rebuild reference links in the new ASTs, and format each AST directly to a CST.
5. Establish fresh semantic evidence for the organized units and check the relocated bindings against their intended destinations. Publish `source.Source` with the units, paths and local import relations. The shared Files printer emits these paths and texts. Failed placement, ambiguity, missing dependencies or failed semantic admission produce no partial Source.

The source language supplies the canonical Cell identity. The Cell mapping selects, for example, export `Cell` from `@bitspark/cell`. The TypeScript formatter derives relative imports between owned units; it uses the configured public specifier for a package export. It has no Cell rule. The shared generation contract and Files/Source carriers already express their respective operations. The implementation needs no universal layout registry.

## Example and scope

Assign the captured TypeScript declarations Page and Cells to file ID `pages` at `generated/pages/page.ts`, and Person to `model` at `generated/model/person.ts`. For a plain ESM .ts profile using .js relative specifiers, the pages unit contains:

```typescript
import type { Person } from "../model/person.js";
import type { Cell } from "@bitspark/cell";

export interface Page<T> { items: T[]; }
export type Cells = Page<Cell<Person>>;
```

Move Person into `pages` in that placement. The file now defines Person locally. Its import disappears and the Cell package import remains. Moving only the model file's path changes the relative specifier. Actual package declarations are explicitly supplied inputs, identified by public specifier in the analysis configuration. The consuming build selects them from its dependencies. This slice covers named type imports in declaration-only ESM .ts units. Namespace imports, runtime imports, re-exports, package discovery and project-wide resolution are separate capabilities.

Qualification exercises actual generated syntax, the Files formatter, fresh target receipts and the printer, and compiles both placements as modules against a separately supplied package. Missing and duplicate placements, unresolved exports, alias conflicts and stale receipts must fail or be resolved explicitly.
