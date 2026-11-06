# Language provider matrix

The tables below are generated from implemented Go providers and typed transformer
factories. A link leads to the implementation's source line; / means unavailable.
Applications compose and register capabilities explicitly. Availability of a stage
does not require every preceding stage to be available. Native support covers the
documented type-declaration subset, not complete programming languages.

TypeScript and Rust resolution require an explicitly supplied compiler and captured
source context; the zero-value text provider does not discover a compiler.
See the [TypeScript](native-semantics.md#first-typescript-profile) and
[Rust](native-semantics.md#first-rust-profile) semantic profiles.

The [matrix generator](../tools/provider-matrix/README.md) discovers language
markers, implemented provider contracts, generic constructor collaborators and
typed transformer endpoints. The [update workflow](../.github/workflows/provider-matrices.yml)
regenerates these regions after source changes and automatically lands a docs-only
pull request when their content changes. Keep explanations outside the generated
regions; do not edit cells or maintain a separate capability catalogue.

The Level column comes from `bittype:dimension` comments on declared operation
boundary types, not package names. `Shared` requires an explicit `bittype:shared`
contract annotation; absent dimension metadata appears as `Unclassified`, and
conflicting dimensions appear as `Mixed: A + B`. See the generator's
[metadata rules](../tools/provider-matrix/README.md#dimension-metadata).
Resolution and analysis belong to the Files dimension through their semantic
request/result types, even when called independently without physical files.

Provider interfaces use their package to identify the role: `parse.Provider`,
`check.Provider`, `format.Provider`, `print.Provider`, `resolution.Provider` and
`analysis.Provider`. Each embeds `language.Provider` for ID and Name. The
parse/check/format/print contracts are independently usable, with no dependency
on a combined contract. The engine registers them separately with
`RegisterSyntaxParser`, `RegisterChecker`, `RegisterSyntaxFormatter` and `RegisterSyntaxPrinter`.
Source input belongs to `parse.Input`. `text.Provider` and the engine's combined
provider registration are removed without compatibility aliases; `text` contains
shared diagnostics only.

Implementation packages use role names and explicit dimension paths, as specified
in [the repository layout](../LAYOUT.md). Every language composition package has
the same Syntax constructors and a Files printer constructor. Horos's additional
Files constructors expose its existing semantic capabilities. Each constructor
takes an explicit language/profile descriptor and registers nothing. The matrices
describe implemented capabilities; constructors do not create additional cells.

<!-- generated:providers:begin -->

| Level | Provider | C++ | Go | Horos | Rust | TypeScript |
| --- | --- | :---: | :---: | :---: | :---: | :---: |
| Files | `analysis.Provider[L]` | / | [Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/golang/go/analyze.go#L21) | [Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/horos/go/analyze.go#L36) | [Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/rust/go/analyze.go#L21) | [Analyze](https://github.com/Bitspark/bittype/blob/main/analyze/files/typescript/go/analyze.go#L21) |
| Files | `format.Provider[format.FilesRequest[L], source.Source[L]]` | / | / | [Format](https://github.com/Bitspark/bittype/blob/main/format/files/horos/go/format.go#L52) | / | / |
| Files | `parse.Provider[parse.FilesRequest[L], source.Source[L]]` | / | / | [Parse](https://github.com/Bitspark/bittype/blob/main/parse/files/horos/go/parse.go#L47) | / | / |
| Files | `print.Provider[source.Source[L], files.Files]` | [Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L12) | [Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L12) | [Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L12) | [Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L12) | [Print](https://github.com/Bitspark/bittype/blob/main/print/files/go/print.go#L12) |
| Files | `resolution.Provider[L]` | / | [Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/golang/go/resolver.go#L23) | [ResolveReferences](https://github.com/Bitspark/bittype/blob/main/resolve/files/horos/go/resolver.go#L48) | [Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/rust/go/resolver.go#L41) | [Resolve](https://github.com/Bitspark/bittype/blob/main/resolve/files/typescript/go/resolver.go#L41) |
| Shared | `language.Provider` | [Language](https://github.com/Bitspark/bittype/blob/main/syntax/cpp/go/nodes.go#L12) | [Language](https://github.com/Bitspark/bittype/blob/main/syntax/golang/go/nodes.go#L12) | [Language](https://github.com/Bitspark/bittype/blob/main/syntax/horos/go/nodes.go#L10) | [Language](https://github.com/Bitspark/bittype/blob/main/syntax/rust/go/nodes.go#L12) | [Language](https://github.com/Bitspark/bittype/blob/main/syntax/typescript/go/nodes.go#L12) |
| Syntax | `check.Provider[L]` | [Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/cpp/go/check.go#L55) | [Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/golang/go/check.go#L43) | [Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/horos/go/check.go#L11) | [Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/rust/go/check.go#L70) | [Check](https://github.com/Bitspark/bittype/blob/main/check/syntax/typescript/go/check.go#L50) |
| Syntax | `format.Provider[ast.View[L], cst.CST[L]]` | [Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/cpp/go/format.go#L22) | [Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/golang/go/format.go#L22) | [Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/horos/go/format.go#L41) | [Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/rust/go/format.go#L22) | [Format](https://github.com/Bitspark/bittype/blob/main/format/syntax/typescript/go/format.go#L23) |
| Syntax | `parse.Provider[parse.Input, cst.CST[L]]` | [Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/cpp/go/parse.go#L24) | [Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/golang/go/parse.go#L25) | [Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/horos/go/parse.go#L60) | [Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/rust/go/parse.go#L24) | [Parse](https://github.com/Bitspark/bittype/blob/main/parse/syntax/typescript/go/parse.go#L41) |
| Syntax | `print.Provider[cst.CST[L], string]` | [Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L13) | [Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L13) | [Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L13) | [Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L13) | [Print](https://github.com/Bitspark/bittype/blob/main/print/syntax/go/print.go#L13) |

<!-- generated:providers:end -->

The generic CST printer emits stored Syntax. Files printing has a reusable
implementation in `print/files/go`,
composed with the selected language's identity and Syntax printer. It consumes an
existing `Source[L]` and does not depend on a resolver or analyzer. Horos now
implements the [documented source syntax](horos-syntax.md). Its checker still
requires no parser, formatter or printer.
Each stage interface uses language identity and its own typed input/output.
Configured adapters bind those interfaces to representation envelopes. Syntax admission
checks structure and supported syntax; success does not establish reference
resolution or semantic validity. `representation.NewLanguage` captures identity, profile and Go marker only.
`check.NewChecker` binds the checker; `Engine.RegisterChecker` installs it
independently. Parser-only and formatter/printer-only engines are supported. Horos exposes its
[resolver and analyzer separately](resolution.md).
The [shared semantic model](semantics.md) and [analysis guide](analysis.md) explain
their results. The existing Horos canonical-value helpers are additional frontend
entry points, not a substitute for the registered resolution capability.

Horos uses one declaration grammar with [general structural kinds](kinds.md),
including ordinary parameterized and recursive kind definitions. Kind meanings
remain Horos-owned; the generic semantic graph carries their facts.
Wire expressions do not create a second language or registered profile.

ASTTransformer relates two languages and consumes explicit semantic input.
The following table derives direct edges from their implemented factories.
The [pagination profile](semantic-generation.md) explains the native representation
scope and its semantic requirements.

<!-- generated:transformers:begin -->

| Source → target | C++ | Go | Horos | Rust | TypeScript |
| --- | :---: | :---: | :---: | :---: | :---: |
| C++ | / | / | / | / | / |
| Go | / | / | / | / | / |
| Horos | [New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-cpp/go/transform.go#L14) | [New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-go/go/transform.go#L14) | / | [New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-rust/go/transform.go#L14) | [New](https://github.com/Bitspark/bittype/blob/main/transforms/horos-to-typescript/go/transform.go#L14) |
| Rust | / | / | / | / | / |
| TypeScript | / | / | / | / | / |

<!-- generated:transformers:end -->

Transformer cells show direct implemented edges, not paths inferred by the engine.
They do not imply support for every representation profile or every input type.

Reverse and same-language cells require actual transformer implementations; a
checked snapshot alone is not a generation or canonicalization mapping.

The [Syntax and Files design](source.md) is implemented: the same generic role
interfaces serve Syntax/CST/AST and Files/Source/Semantics. Methods take a context.
The engine registers Syntax and Files capabilities independently. Horos Files
compositions are constructed explicitly with `provider/horos.NewFilesParser`,
`NewFilesFormatter` and `NewFilesPrinter`; they do not install engine capabilities.
Go independently supplies `NewFilesResolver` and `NewFilesAnalyzer` under its
[bounded native profile](native-semantics.md); this does not add Files
parsing or formatting. Files parsing and formatting compositions require their language's semantic
collaborators. Missing collaborators do not disable Files printing: the generic
printer can serve an existing admitted Source when configured with its Syntax
printer. The matrix describes each capability independently, not whether a
complete Files pipeline is currently available for that language.

See the [checker contract](../contracts/check/go/provider.go),
[parser contract](../contracts/parse/go/provider.go),
[source input](../contracts/parse/go/input.go),
[formatter contract](../contracts/format/go/provider.go),
[printer contract](../contracts/print/go/provider.go),
[shared diagnostics](../contracts/text/go/diagnostic.go),
[Horos composition](../provider/horos/go/provider.go) and
[enforced repository layout](../LAYOUT.md).
