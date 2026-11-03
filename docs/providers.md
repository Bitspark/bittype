# Language provider matrix

**3 November 2026. Current implemented source.** X means implemented and exposed
by the language provider; / means unavailable. Applications register capabilities
explicitly. Native support covers the documented type-declaration subset, not
complete programming languages.

Provider interfaces use their package to identify the role: `language.Provider`,
`check.Provider`, `resolution.Provider`, `analysis.Provider` and `text.Provider`.
The individual text capabilities are `text.Parser`, `text.Formatter` and
`text.Printer`; each embeds `language.Provider`. The independent `check.Provider`
checks an AST view and embeds the same identity interface, with no text or CST
dependency. The combined `text.Provider` embeds the three text capabilities and
`check.Provider`. Superseded interface names have no compatibility aliases.

| Provider / capability | Horos | Horos/Wire | TypeScript | Go | Rust | C++ |
| --- | :---: | :---: | :---: | :---: | :---: | :---: |
| language.Provider: ID and Name | X | X | X | X | X | X |
| check.Provider: AST syntax checking | X | X | X | X | X | X |
| text.Parser: text to CST | / | / | X | X | X | X |
| text.Formatter: AST to CST | / | / | X | X | X | X |
| text.Printer: CST to text | / | / | X | X | X | X |
| text.Provider: combined text capabilities | / | / | X | X | X | X |
| resolution.Provider: reference resolution | X | X | / | / | / | / |
| analysis.Provider: semantic analysis | X | X | / | / | / | / |

The native printer providers delegate to the same generic printer. Horos remains
AST-only: its checker requires no parser, formatter or printer. Syntax admission
checks structure and supported syntax; success does not establish reference
resolution or semantic validity. `representation.NewLanguage` accepts this
independent checker for registration. Horos exposes its
[resolver and analyzer separately](resolution.md).
The [shared semantic model](semantics.md) and [analysis guide](analysis.md) explain
their results. The existing Horos canonical-value helpers are additional frontend
entry points, not a substitute for the registered resolution capability.

ASTTransformer relates two languages. The explicit [ground-to-optional partial
lift](horoswire.md#explicit-transformation) is implemented. Native mappings remain
unimplemented. The same-language checked
snapshot operation does not establish a generation or canonicalization mapping.
Actual mappings will need a separate source-by-target matrix.

See the [checker contract](../contracts/check/go/provider.go),
[text contracts](../contracts/text/go/provider.go),
[Horos composition](../provider/horos/go/provider.go),
[Horos/Wire composition](../provider/horoswire/go/provider.go) and
[enforced repository layout](../LAYOUT.md).
