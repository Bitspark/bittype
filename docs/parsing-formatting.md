# Native syntax, concrete trees and text

**Direction: 31 October 2026.** The owner requested native TypeScript, Go, Rust
and C++ parsers, formatters, checks and tests, then selected a shared CST alongside
the AST. The [Horos frontend](horos.md) now implements its [source grammar](horos-syntax.md)
through these same interfaces. The first [semantic pagination transformers](semantic-generation.md) now produce native ASTs; runtime adapters remain deferred.
The subsequent [syntax engine](engine.md) supplies registration,
route execution and composition of these native text stages.
The [Syntax and Files design](source.md) applies parse, format and print at both
levels: Syntax/CST/AST and Files/Source/Semantics. Files parsing explicitly
composes per-file Syntax parsing with resolution/analysis and publishes Source
only with established Semantics. No parser opens files or concatenates a project
into one text. The Files stages and multi-file Horos resolution are implemented;
Go, TypeScript and Rust have [bounded native semantic providers](native-semantics.md);
C++ semantics remain deferred. Only TypeScript currently has native Files formatting.

## Pipeline

```text
Parser:     source text -> CST[L]
Formatter:  AST view[L] -> CST[L]
Printer:    CST[L] -> source text

AST[L] and CST[L] both expose ASTView() AST[L].
```

A parser constructs an abstract view and a concrete tree together. Accessing the
view is a field-level operation; it does not parse, resolve names or invoke a
CST-to-AST provider. Native source remains native syntax, not the authoritative
canonical input for a future generator.

A formatter consumes the abstract view and constructs concrete token/trivia
nodes directly. The shared printer only traverses supplied leaves in order.
Escaping, punctuation, parentheses, indentation and declarator precedence belong
to each language formatter. Go uses its standard parser internally; production
formatting does not print Go source and reparse it into a CST.

## Carriers and ownership

[AST views](../trees/ast/go/view.go) expose an already constructed AST through the
language-parameterized `ast.View[L]` interface. The existing AST struct implements
it, as does the [CST carrier](../trees/cst/go/cst.go).

A CST owns its abstract arena, concrete arena and concrete root. Each concrete
node has a language-owned `Kind[L]`, optional leaf text, ordered children,
optional links into the stored abstract arena, and origin metadata. Non-nil
`Text` identifies a leaf; nil identifies a container. Empty containers represent
empty modules/layout groups. Zero-length lexical leaves are rejected.

Concrete kinds belong to the native syntax package. The generic CST and printer
do not contain a language list or interpret kind names. Grammar containers can
retain native rule names; layout groups and lexical/trivia kinds are also supplied
by the language. Concrete grouping follows source intervals. A grouped native
declaration can correspond to multiple abstract nodes, so `Syntax` is a list of
addresses. Links are optional; an implicit or normalized abstract node need not
have its own concrete spelling. This is a lossless concrete representation of the
admitted syntax, not a complete compiler CST for the whole language.

`ASTView()` returns borrowed, read-only storage. Its slices are not frozen by Go.
All stage inputs are read-only by contract. Formatter output receives a deep
snapshot of the input AST, including attributes, children, references and origin
slices. Later mutation of the caller's original AST must not change that output.
Parser results own their newly constructed storage.

Tree IDs and source IDs are caller-supplied opaque addresses. They are not
filesystem paths, hashes or declaration identities. Parser origins use half-open
UTF-8 byte offsets. Name resolution and cross-tree reference checking are separate
work; parsers do not invent resolved target addresses.

## Two different round trips

For admitted source, including comments and whitespace:

```text
Print(Parse(text)) == text
```

This retains trailing trivia, CRLF, accepted BOMs and original literal spelling.
Inserted Go parser semicolons are not concrete source tokens. Unsupported
preprocessor/directive forms are rejected rather than executed or discarded.

Formatting deliberately chooses a new layout from abstract syntax:

```text
Parse(Print(Format(ast))).ASTView()
```

This preserves the text-representable native syntax after the documented
normalizations below. Text cannot recreate original node numbering, origins or
externally resolved references. These are retained in the formatter's stored
AST snapshot, but a fresh parse has new addresses and source origins.
Syntax comparison in tests is not a declaration identity algorithm, structural
fingerprint, or permission to substitute library contracts.

Normalization expands grouped Go declarations/fields/parameters into individual
abstract declarations, removes redundant type parentheses, and treats omitted
false/default flags as their documented defaults. Rust parenthesized single types
remain distinct from one-element tuples. C++ access sections become member
attributes, `enum struct` and `enum class` both use the scoped flag, and
`(void)` parameter lists become empty lists. TypeScript member separators and
quote choices may change during formatting. Generic binders and ordered type
applications remain explicit, without specialization or alias expansion.

## Provider contracts

[language.Provider](../contracts/language/go/language.go) is the shared base:

```go
type Provider interface {
    ID() ID
    Name() string
}
```

The former `language.Language` interface has been replaced without an alias.
Native marker types remain named `Language`; those are Go type parameters for
the native AST/CST and implement the shared identity contract.

The [checker contract](../contracts/check/go/provider.go) is independent of text
and CST. It admits AST structure and supported language syntax through a borrowed
view; success establishes neither reference resolution nor semantic validity:

```go
// package check
type Provider[L any] interface {
    language.Provider
    Check(ast.View[L]) error
}
```

The three other stage contracts also embed language identity and remain
independent of text composition. [Parsing](../contracts/parse/go/provider.go)
owns [source input](../contracts/parse/go/input.go):

```go
// package parse
type Input struct {
    SourceID ast.SourceID
    TreeID   ast.TreeID
    Text     string
}
type Provider[I, O any] interface {
    language.Provider
    Parse(context.Context, I) (O, error)
}
```

IDs remain caller-owned opaque addresses; parsing does not open paths or fetch
dependencies. [Formatting](../contracts/format/go/provider.go) consumes an AST view:

```go
// package format
type Provider[I, O any] interface {
    language.Provider
    Format(context.Context, I) (O, error)
}
```

[Printing](../contracts/print/go/provider.go) consumes a concrete tree:

```go
// package print
type Provider[I, O any] interface {
    language.Provider
    Print(context.Context, I) (O, error)
}
```

Syntax uses parse.Provider[parse.Input, cst.CST[L]],
format.Provider[ast.View[L], cst.CST[L]] and print.Provider[cst.CST[L], string].
Files uses the same roles with the captured requests described in [Source](source.md).

Each native parser, formatter and checker package supplies its specialized
provider. Each native stage supplies its own implementation. The shared Syntax printer lives at `print/syntax/shared/provider/go`. The separate
`contracts/print/go` package owns the printer interface. Applications can use
separate implementations for each capability.

The engine first registers identity/profile metadata from
`representation.NewLanguage[L](language.Provider, profile)`. It then independently
registers adapters made by `parse.NewSyntaxParser`, `check.NewChecker`,
`format.NewSyntaxFormatter` and `print.NewSyntaxPrinter`. Each adapter captures the same
ID/profile/Go-marker contract and owns its boundary snapshots. Registration of
one stage never installs any other stage; additional stages can be attached later.

Parse, Format and Print require only their respective capabilities. Render
requires a formatter and printer and preflights both before calling either.
Check, semantic Capture and transformation execution require explicit syntax
checkers. Generic Pack/PackCST operations validate carrier structure, not native
grammar; native stage implementations still enforce their own grammar contracts.

Horos exposes checking, resolution and analysis without any text capability.
`text.Provider` and the combined engine registration are removed without aliases.
The [text package](../contracts/text/go/diagnostic.go) retains shared diagnostics.

No provider registers itself. The syntax engine accepts shared contracts
without concrete imports or a name switch. Heterogeneous AST type erasure belongs
in trees/representation, consumed by explicit engine adapters. The typed generic
AST/CST carriers remain independent. Equal ID strings do not
establish equal Go marker types or compatible grammars.

For example:

```go
// tsparse/tsformat are the individual TypeScript stage packages.
// sharedprint is print/syntax/shared/provider/go.
p := tsparse.Provider{}
f := tsformat.Provider{}
printer, err := sharedprint.New[typescript.Language](typescript.Language{})
if err != nil { return err }
parsed, err := p.Parse(ctx, parse.Input{
    SourceID: "input", TreeID: "tree",
    Text: "type Names = string[];",
})
if err != nil { return err }
original, err := printer.Print(ctx, parsed) // retained concrete spelling
if err != nil { return err }
formatted, err := f.Format(ctx, parsed)
if err != nil { return err }
output, err := printer.Print(ctx, formatted) // newly selected layout
if err != nil { return err }
_, _ = original, output
```

These signatures fit the existing typed pipeline Stage contract. The engine
adapts them explicitly; filesystem IO, dependency retrieval and implicit package
discovery remain outside these stages.

## Checks and admitted syntax

The [native vocabulary](syntax.md) is the support boundary. Parsing consumes the
whole input. Unsupported constructs fail; no executable-body skipping, raw-source
escape node, arbitrary expression attribute or fallback implementation is used.

Generic checks cover valid roots and node addresses, containment ownership,
cycles, reachability, attributes, references and spans. Native checks cover
known kinds, allowed/required attributes, child roles and order/cardinality,
lexical atoms and representable contextual restrictions. They are syntax checks,
not application-value validators or complete compiler type checking.

Imports/includes are retained as syntax and never fetched. Checking cannot
establish that imported names exist, satisfy constraints or have compatible
visibility/layout. Native compiler fixtures provide separate evidence for the
qualified examples, not a proof for every external program.

The current source limit is 4 MiB and lexical/arena limits are 200,000 elements.
Each abstract/concrete arena has a 16 MiB aggregate string budget and a 1,600,000
entry budget for its child/attribute/reference/origin lists. Repeated string
references count per occurrence, before validation scans or printing allocate
expanded output. Nesting limits are also explicit resource bounds. Violations
return diagnostics rather than usable partial output. Diagnostics can identify
source spans or abstract node addresses.

Some deliberate boundaries deserve explicit mention:

- TypeScript: type-only imports, aliases, interfaces and the listed type syntax,
  including [conditional and mapped types](syntax.md#typescript). Conditional
  operands, `infer` binders/constraints, mapped key binders, `as` remapping and
  modifier addition/removal have structured AST children or attributes.
  Formatting preserves precedence and builds CST directly; source CST retains
  comments and explicit plus-modifier spelling. Mapped types require explicit
  value annotations and cannot include ordinary members. Template-literal types,
  classes and executable declarations remain unsupported. Syntax admission does
  not perform conditional evaluation, inference or mapped-key checking.
- Go: imports, aliases/defined types and the listed type syntax. Computed array
  bounds and value/function declarations remain unsupported.
- Rust: the listed type declarations and trait signatures. Unit structs
  (`struct X;`) are rejected because the current AST distinguishes named-field
  and tuple structs but has no unit-struct form. Const generics, macros,
  attributes, executable implementations and higher-ranked types remain unsupported.
  [Qualified projections and associated-type equality bindings](syntax.md#rust)
  are structured AST nodes: `<T as Trait>::Item`, `<T>::Item` and
  `Trait<Item = U>`. They compose with nested paths, type/lifetime arguments,
  raw identifiers and generic associated-type uses. Associated-type argument
  bounds (`Item: Bound`) remain unsupported. Equality bindings require a single associated name with optional
  arguments and a type RHS; they follow positional arguments. Parsing retains
  source trivia and optional turbofish separators; formatting builds normalized
  CST directly from the structured AST, preserving references and origins.
  Bare primitive spellings in type-expression positions use PrimitiveType;
  ordinary multi-segment paths remain PathType, and `<...>::...` uses
  QualifiedPathType. These distinctions do not perform name resolution, prove
  what a shadowing declaration denotes or check associated-type equalities.
  Generic associated-type declarations and type-only trait impls support
  type/lifetime parameters, trait/lifetime bounds, relaxed bounds such as
  `?Sized`, and where clauses. Associated definitions reuse TypeAlias children
  within TraitImpl and format their where clauses after the RHS. The parser
  retains pre-equals where clauses in source CST and normalizes them during
  formatting. Methods, constants, inherent/unsafe/negative impls, associated-type
  defaults and generic defaults in these new contexts are rejected. Admission
  does not establish coherence, completeness or bound satisfaction; those remain
  native semantic questions, separate from the generation SPI.
- C++: the listed declarations and declarators. Function bodies, non-type
  template arguments, arbitrary preprocessing, constructors and expression
  initializers remain unsupported. Backslash-newline splicing is rejected even
  inside comments, before it could change declaration meaning. Qualified member-function-pointer signatures
  are rejected because FunctionType has no method CV/ref fields.

`Cell<T>`, `Function<T>`, `Stream<T>` and other library names are ordinary native
references/applications. None receives a kind or special treatment in these
foundations.

## Layout and dependencies

The [stage-first repository layout](../LAYOUT.md) defines the uniform native
matrix and directed dependency rules. Shared provider conformance, round-trip
and fuzz cases cover every language discovered in the stage directories.

```text
trees/ast/go/view.go                         AST-view interface
trees/cst/go/cst.go                          concrete carrier and stored AST view
contracts/parse/go/{input,provider}.go      source input and parser contract
contracts/check/go/provider.go             independent AST syntax admission
contracts/format/go/provider.go            formatter contract
contracts/print/go/provider.go             printer contract
contracts/text/go/diagnostic.go            shared diagnostics only
check/syntax/shared/ast/go/structure.go                  abstract carrier checks
check/syntax/shared/cst/go/check.go                  concrete carrier checks
check/syntax/<language>/provider/go/                   native syntax checks
parse/syntax/<language>/provider/go/                   text -> CST
format/syntax/<language>/provider/go/                  AST view -> CST
print/syntax/shared/provider/go/print.go                      shared CST -> text
internal/                             implementation helpers and native details
integration/go/                       independent contract/text/compiler tests
testdata/text/<language>/valid/        native source fixtures
```

The language directories are `typescript`, `golang`, `rust`, `cpp` and `horos`; the trailing
`go` denotes implementation language. Private parsing and concrete construction
helpers are not public builders.

AST, pipeline and language identity remain independent. CST depends on AST.
Native syntax depends on these generic carriers and identity, never another
language. Native parsers/formatters/checkers depend on their own syntax and shared
mechanics. The printer depends only on generic concrete structure/checks.
The native Syntax parser/checker/formatter stages have no foreign compiler or
external model dependency. This does not describe all production packages:
Horos tooling consumes Horos/Ontos, and native semantic stages use the dependencies
listed in [native semantics](native-semantics.md) and [dependencies](dependencies.md).

## Qualification

The maintained GitHub workflows run against exact committed revisions. Tests
cover source retention, independent AST-to-CST-to-text expectations, formatting
stability, typed provider identity/separation, malformed arenas, unsupported
syntax and resource bounds. Native compiler checks use pinned toolchains on
remote runners: Go 1.27.3, TypeScript 5.9.3 (ES2020), Rust 1.90.0 (edition 2024),
and GCC 14.2.0 (C++20). Workstation use is limited to source editing, Git and coordination.

The Horos declaration model and [text and semantic frontend](horos.md) are implemented
with the [readable Horos source grammar](horos-syntax.md). Native AST generators,
runtime bindings and package releases are separate work. Registration is
implemented by the separately described engine.
