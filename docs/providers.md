# Language provider matrix

**4 November 2026. Current implemented source.** X means implemented and exposed
by the language provider; / means unavailable. Applications register capabilities
explicitly. Native support covers the documented type-declaration subset, not
complete programming languages.

Provider interfaces use their package to identify the role: `parse.Provider`,
`check.Provider`, `format.Provider`, `print.Provider`, `resolution.Provider` and
`analysis.Provider`. Each embeds `language.Provider` for ID and Name. The
parse/check/format/print contracts are independently usable, with no dependency
on a combined contract. The engine registers them separately with
`RegisterParser`, `RegisterChecker`, `RegisterFormatter` and `RegisterPrinter`.
Source input belongs to `parse.Input`. `text.Provider` and the engine's combined
provider registration are removed without compatibility aliases; `text` contains
shared diagnostics only.

| Provider / capability | Horos | TypeScript | Go | Rust | C++ |
| --- | :---: | :---: | :---: | :---: | :---: |
| language.Provider: ID and Name | X | X | X | X | X |
| check.Provider: AST syntax checking | X | X | X | X | X |
| parse.Provider: text to CST | X | X | X | X | X |
| format.Provider: AST to CST | X | X | X | X | X |
| print.Provider: CST to text | X | X | X | X | X |
| resolution.Provider: reference resolution | X | / | / | / | / |
| analysis.Provider: semantic analysis | X | / | / | / | / |

All five printer providers delegate to the same generic printer. Horos now
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

Horos uses one declaration grammar with explicit Type/Data parameter bounds.
Wire expressions do not create a second language or registered profile.

ASTTransformer relates two languages and consumes explicit semantic input.
The [pagination profile](semantic-generation.md) implements these directed edges:

| Source | TypeScript | Go | Rust | C++ |
| --- | :---: | :---: | :---: | :---: |
| Horos (admitted analysis) | X | X | X | X |

Reverse mappings are not implemented. Template lifting remains removed. A
same-language checked snapshot is not a generation or canonicalization mapping.

The proposed [layout.Provider](layout.md#3-independent-layoutprovider) arranges
same-language AST collections into files before formatting. It is not implemented
for any language and is not counted in the current matrix. The same proposal
describes per-file parsing and the remaining cross-file resolution work.

See the [checker contract](../contracts/check/go/provider.go),
[parser contract](../contracts/parse/go/provider.go),
[source input](../contracts/parse/go/input.go),
[formatter contract](../contracts/format/go/provider.go),
[printer contract](../contracts/print/go/provider.go),
[shared diagnostics](../contracts/text/go/diagnostic.go),
[Horos composition](../provider/horos/go/provider.go) and
[enforced repository layout](../LAYOUT.md).
