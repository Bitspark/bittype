# Abstraction levels and typed relations

**Design note · 8 November 2026.** A level of abstraction determines which
distinctions we can express. Its **data model** says what exists at that level.
A **context** identifies a particular document, tree or program together with
the information needed to interpret its entities. A **typed relation** says
which contexts or elements correspond, under which conditions and with what
cardinality. This guide describes the proposed architecture.

## The data model at each level

An element reference consists of **a context identity and a local address of
the appropriate kind**. A byte range belongs to one document. A type-use
occurrence belongs to one tree. The binding established for that occurrence
belongs to one analyzed program. The same tree can participate in a different
program with different dependencies. The abstraction level follows what the
model expresses. A JSON document and a table of nodes can both encode the same
abstract tree.

### Agnostic kernel and language-specific rules

The kernel owns the guarantees shared across languages. It invokes a supplied
interpretation when a result depends on a language or encoding. A grammar,
for example, can supply declarative rules. The matrix assigns these
responsibilities independently of package layout.

| Abstraction level | Context and entities | Agnostic kernel responsibility | What the language, codec or tool supplies |
| --- | --- | --- | --- |
| Encoding | One byte document, its length, offsets and ranges | Address the bytes in that document and check span bounds. | The selected codec interprets byte sequences. UTF-8 can serve many programming languages. |
| Text | One character sequence, its positions and its coordinate convention | Locate an occurrence in the specified document and units. Two identical letters at different positions remain distinct. | A consumer can specify UTF-16 code-unit offsets, as the TypeScript compiler does. Ordinary text addressing requires only that convention. |
| Concrete syntax (CST) | One concrete tree of groups, tokens and trivia, with links to text spans or to syntax being formatted | Retain occurrence identity, order and typed links to their contexts. | The language supplies its grammar and layout rules. In the `Page` declaration, TypeScript interprets `<` as opening the parameter list. |
| Abstract syntax (AST) | One tree of nodes, containment roles, type uses and references | Qualify addresses by tree and distinguish containment from reference. A node can link to its contributing syntax. | The language defines node kinds and the permitted children of each role. A TypeScript type parameter can have a default type. |
| Declaration semantics | A program with its AST units, selected dependencies and settings, and the entities, binders, terms and bindings established under them | Associate each fact with the program and its typed endpoints. Preserve the result and coverage of the language judgment. | TypeScript supplies the rules for lexical scope, type lookup and default arguments. Horos supplies its rules for pinned declarations and kind constraints. |

Context dependencies follow the operation. One semantic program can depend on
several AST units. A parsed concrete tree belongs to a text document, while a
formatted concrete tree can originate from an AST before any output text exists.
Within a single semantic program, two occurrences can also have different
lexical scopes. A file path locates text. A module and its scopes determine
which names are visible at a use. [Syntax and Files dimensions](source.md)

## Page across the five levels

Consider these TypeScript declarations, with no outer declaration also named
`Page` or `T`:

```typescript
type Page<T = string> = { first: T; last: T };
type TextPage = Page;
```

This example uses TypeScript's language semantics. Bittype's current TypeScript
analysis profile reports generic defaults as unsupported.
[TypeScript profile](native-semantics.md#first-typescript-profile)

Read a row across for one occurrence or fact. An **empty cell** means that this
level has **no distinct counterpart at that site** in this model. It never means
that lookup failed or that nobody checked. Hexadecimal byte spans are for the
example encoded as UTF-8; `␠` makes a single ordinary space visible.

| Occurrence or fact | Encoding (hex) | Text | Concrete syntax | Abstract syntax | Declaration semantics |
| --- | --- | --- | --- | --- | --- |
| Name in the declaration of `Page` | `50 61 67 65` | `Page` | Identifier token | Declared name of type alias | Introduces the type alias `Page` |
| Opening delimiter of its parameter list | `3C` | `<` | Opening angle token | | |
| Name in the declaration of `T` | `54` | `T` | Identifier token | Name of formal type parameter | Introduces binder `Page.T` |
| Space between that `T` and `=` | `20` | `␠` | Whitespace trivia | | |
| Written default on the parameter | `73 74 72 69 6E 67` | `string` | Keyword token | Default type on the formal parameter | The declared default of `Page.T` is the string type |
| `T` in `first: T` | `54` | `T` | Identifier token in `first` | Named type use in the first field | Denotes `Page.T` |
| `T` in `last: T` | `54` | `T` | Identifier token in `last` | Named type use in the second field | Denotes the same `Page.T` |
| `Page` in the body of `TextPage` | `50 61 67 65` | `Page` | Identifier token | Named type use with no written arguments | Denotes `Page` |
| Effective argument at that use | | | | | `Page.T` receives the string type by default |

The delimiter has no AST node of its own, although the AST retains the
parameter list's structure. In the other direction, `TextPage` means
`Page<string>`, although there is no `string` token or argument AST node **at
the use**. The written default lives at the declaration. Likewise, the two
field types occupy different source positions, yet both bind to the same
parameter. `Page.T` names that binder here for explanation; it is not source
syntax. Writing the same `T` inside a different generic scope could denote a
different binder.

## Typed relations between levels

Each schematic signature names its two endpoint types. Its cardinality holds
**under the stated conditions**. The first row relates whole contexts. The
remaining rows relate elements addressed within those contexts.

| Relation | Cardinality for this example | Required context | What the relation establishes |
| --- | --- | --- | --- |
| `analyzedFrom(SemanticProgram, ASTUnit)` | One program can include several units, and one unit can participate in different programs. | The inputs of this analysis. | The AST containing `Page` contributed to this program's semantic facts. |
| `encodedBy(CharacterOccurrence, ByteSpan)` | Each Unicode scalar occurrence has one contiguous span of 1–4 bytes. | One document, decoded as valid UTF-8. | Which bytes encode that occurrence. The four letters of `Page` cover four such spans. |
| `spelledBy(ConcreteTokenOrTrivia, TextSpan)` | Each original token or trivia leaf has one contiguous span. | The parsed text and its coordinate system; this claim is about original leaves. | Where the token or space was written. A synthetic token need not have an original span. |
| `nameToken(SyntaxNameOccurrence, IdentifierToken)` | Each simple written name has one token. | The AST and CST for the same parse. | Which token spelled this particular name occurrence, irrespective of what the name denotes. |
| `introduces(TypeParameterDeclaration, SemanticBinder)` | This formal introduces one binder. | The captured declaration and TypeScript's binding rules. | The declaration of `T` introduces `Page.T`. A use of `T` does not introduce it again. |
| `denotes(TypeNameUse, SemanticEntity)` | Each successfully resolved simple use denotes one entity here; many uses may denote it. | The use's lexical scope, selected dependencies and type namespace in that capture. | Both field uses denote `Page.T`; their spelling alone cannot establish this. |
| `usesDefaultFrom(EffectiveArgumentOccurrence, DefaultTypeSyntax)` | This defaulted argument has one written default; many applications can use that default. | The resolved callee, its formal and the language's defaulting rules at the application. | The effective string argument at `TextPage` comes from the default on `Page.T`, at a different site. |

In the facts below, `element@context` qualifies a local address. Follow the
`T` in the first field through the byte document, text, concrete tree and
declaration tree. `program` is the semantic context that includes `pageAST`.

```text
analyzedFrom(program, pageAST)
encodedBy(firstFieldT@pageText, firstFieldBytes@pageBytes)
spelledBy(firstFieldToken@pageCST, firstFieldSpan@pageText)
nameToken(firstTypeUse@pageAST, firstFieldToken@pageCST)
introduces(formalT@pageAST, Page.T@program)
denotes(firstTypeUse@pageAST, Page.T@program)
denotes(lastTypeUse@pageAST, Page.T@program)
usesDefaultFrom(TextPageArgument@program, writtenStringDefault@pageAST)
```

`analyzedFrom` identifies a whole input unit. The spelling relations locate
the first field's occurrence. Both `denotes` facts point to one binder, so
following that relation back from the binder finds two uses. The last fact
connects the effective argument at `TextPage` to the declaration where its
default was written.

## Where the kernel invokes rules for Page

The kernel invokes supplied rules at transitions between contexts. These are
schematic input and result contracts for reading the example.

| Transition | Rule the kernel invokes | Result for `Page` |
| --- | --- | --- |
| `decode(byte document, encoding) → text` | The selected UTF-8 codec decodes the document. | Byte `54` at each field type spells one `T` at that position. |
| `parse(text, language) → CST with AST view` | TypeScript grammar identifies the formal, its default and both field types. | `nameToken` connects each written type use to its token. The space before `=` remains concrete trivia. |
| `resolve(AST units, environment) → binding result` | TypeScript name lookup uses the lexical scope and type namespace at each occurrence. | `introduces` identifies the formal's binder. Each field use has a `denotes` relation to that binder. The use in `TextPage` denotes the declared `Page`. |
| `analyze(AST units, environment, binding result) → semantic result` | TypeScript's generic application rules supply the declared default for an omitted argument. | `Page.T` receives the string type at `TextPage`. `usesDefaultFrom` identifies the written default. |

The kernel can check that the two field occurrences are in the supplied AST,
that the binder is in the resulting program and that each successful simple use
has one destination here. TypeScript establishes **which** destination it is.
With the same field syntax under `Other<T>`, TypeScript would establish a
binding to `Other.T`. Looking up a byte range or following an already
established binding can use the shared data model directly.

## Absent counterpart or unresolved binding?

For each question, first ask **whether this kind of element participates in
the relation**. If it does, ask **what was established in this capture**.

- **No counterpart by design.** The ordinary space has no AST element of its
  own; asking which declaration it denotes is inapplicable. The defaulted
  argument likewise has no written argument at the use. Neither is a failure.
- **Unresolved binding.** Replace the second line with
  `type TextPage = Missing;` in an environment with no type named `Missing`.
  The named type use exists and requires a destination. Lookup was performed
  and failed. Report the unbound name at that use; do not turn it into an
  empty cell or invent a target.
- **Not checked.** If only parsing has run, or the relevant lookup is outside
  the analysis's coverage, the absence of a recorded edge is not a diagnosis.
- **Ambiguous.** If lookup finds competing candidates and the language requires
  one, that is a different outcome from finding none.

Bittype's [semantic facts](semantics.md#occurrences-and-dependencies) store
established bindings. Consult the [analysis report](analysis.md) for coverage
and diagnostics; a missing graph edge alone cannot distinguish these cases.

## Reading and writing

On reading, a token establishes a spelling and the language interprets the
name in its captured scope. On writing `Page<Cell<Person>>`, the kernel
selects the Cell library's mapping, which names the intended TypeScript type
exported by `@bitspark/cell`. TypeScript supplies the syntax for the bound
use. Files formatting chooses its import spelling under the consumer's
placement. Moving `Person` into the file containing the alias removes its
relative import and preserves the Cell package import in this example.

The intended destination and the binding established for the organized
target program answer different questions. The latter must be checked against
the former under the actual target environment. Neither the spelling nor the
absence of an import establishes that the two agree. See the
[planner's pagination walkthrough](generation-planner.md#pagination-on-read-and-write)
and the [Files placement example](source.md#from-generated-declarations-to-files).

## Correspondence with the current implementation

The contexts above are a conceptual model, not five implemented context classes.
These current carriers implement parts of that model. Their identities and
relations must not be collapsed into filenames or same-spelled identifiers.

| Conceptual context/relation | Current carrier or contract | Established meaning and remaining limit |
| --- | --- | --- |
| Artifacts and supplied dependency environment | Application configuration and `files.Files` | Files carries selected paths/text; package-registry discovery and dependency fetching remain application-owned. |
| Text document and positions | `parse.Input`, `ast.SourceID`, source spans | Identifies supplied text and positions; a filename is not a semantic declaration identity. |
| Concrete tree context | `cst.CST[L]` with its retained abstract tree | Owns tokens/trivia and concrete/abstract links. Parsing retains source spelling; formatting selects new spelling. |
| Abstract tree collection | `ast.AST[L]`, `ast.Set[L]`, `representation.ASTSet` | Owns nodes and typed references. `ast.NodeRef` includes tree identity and local address; packing establishes carrier structure, not binding truth. |
| Captured program context | `semantic.Snapshot` | Retains the admitted tree collection, entry, language/profile and optional Files context. It is not itself a successful semantic judgment. |
| Established declaration semantics | `Semantics[L]` plus `semantic.Result` | Stores entities, terms, facts and bindings; result receipts associate coverage and judgments with the exact capture/settings. Language rules establish their meaning. |
| Concrete/abstract relation | CST abstract links and node origins | Structural correspondence and positions; arbitrary origins are not proof of semantic equivalence. Native compiler bridges check their required source/AST correspondence. |
| AST occurrence to semantic entity/term | Occurrence bindings and `TermAt`/`Resolve` queries | Only recorded answers are established. A missing required binding is unresolved/unsupported, not an intentional absent counterpart. |
| Source semantics to intended native declarations | `generation.GenerationPlan` and closed `generation.Program[L]` | Selected mappings and scoped substitutions fix native intent. This is not an analyzed target `Semantics[L]`. Common source navigation still needs [GEN-01](generation-gaps.md#gen-01). |
| Intended declaration/use to emitted syntax | `generation.Output` snapshot, exposed sites and external occurrences | Construction records actual sites. These are intention associations and checked syntax, not a cross-language correctness proof. |
| Generated syntax to placed native program | TypeScript generated Files formatting and fresh native evidence | Rebuilds imports/exports under placement and checks destinations and exposure. Other native Files stages remain [proposed](native-file-placement.md). |
| Semantic program and units to concrete output | `source.Source[L]`, its units, relations and receipts | Associates organization, CSTs and semantics; printing emits Files without choosing another scope or import policy. |
| Source/native claim association | `contracts/verification/go` and explicit Logos interpretation | Relates actual captures and sites under receiver requirements. Optional verification is distinct from target analysis and from generation. |

The language kernels invoke the relevant contracts at each boundary: parse/check
for syntax, resolution/analysis for meaning, declarative mapping for native intent,
Target.Build for syntax construction, and Files formatting for organization.
Their exact interfaces and suppliers are in [text stages](parsing-formatting.md),
[semantic stages](analysis.md), [generation SPI](generation-spi.md),
[Files handoff](typescript-type-imports.md) and [verification](verification.md#contribution-contracts).
There is no implemented universal cross-level relation registry or generic
"previous context" pointer. Cardinality and authority come from each operation's
contract. The absent/unresolved distinction in the tables remains a requirement
when extending these contracts, not a claim that one common relation-status enum
already exists.
