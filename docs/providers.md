# Language provider matrix

**2 November 2026. Current implemented source.** X means implemented and exposed
by the language provider; / means unavailable. Applications register capabilities
explicitly. Native support covers the documented type-declaration subset, not
complete programming languages.

Provider interfaces use their package to identify the role: `language.Provider`,
`resolution.Provider`, `analysis.Provider` and `text.Provider`. The individual
text capabilities are `text.Parser`, `text.Formatter`, `text.Checker` and
`text.Printer`; each embeds `language.Provider`. The combined `text.Provider`
embeds all four. Superseded interface names have no compatibility aliases.

| Provider / capability | Horos | TypeScript | Go | Rust | C++ |
| --- | :---: | :---: | :---: | :---: | :---: |
| language.Provider: ID and Name | X | X | X | X | X |
| text.Checker: syntax checking | X | X | X | X | X |
| text.Parser: text to CST | / | X | X | X | X |
| text.Formatter: AST to CST | / | X | X | X | X |
| text.Printer: CST to text | / | X | X | X | X |
| text.Provider: combined text capabilities | / | X | X | X | X |
| resolution.Provider: reference resolution | X | / | / | / | / |
| analysis.Provider: semantic analysis | X | / | / | / | / |

The native printer providers delegate to the same generic printer. Horos remains
AST-only and exposes its [resolver and analyzer separately](resolution.md).
The [shared semantic model](semantics.md) and [analysis guide](analysis.md) explain
their results. The existing Horos canonical-value helpers are additional frontend
entry points, not a substitute for the registered resolution capability.

ASTTransformer relates two languages. No production language-to-language mapping
exists yet; its interface and route engine exist. The same-language checked
snapshot operation does not establish a generation or canonicalization mapping.
Actual mappings will need a separate source-by-target matrix.

See the [provider contracts](../contracts/text/go/provider.go),
[Horos composition](../provider/horos/go/provider.go) and
[enforced repository layout](../LAYOUT.md).
