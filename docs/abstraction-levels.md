# Abstraction levels and typed relations

**Design note · 8 November 2026.** A level of abstraction determines which
distinctions we can express. Its **data model** says what exists at that level.
A **typed relation** says how particular elements at different levels correspond,
in what context and with what cardinality. These are conceptual contracts, not
new Bittype APIs.

## The data model at each level

From the spelling of a document to the meaning of a declaration:

| Level | Elements in its data model | What the level can express |
| --- | --- | --- |
| Encoding | A byte sequence and spans in one captured document | Which bytes occupy which offsets. The character encoding is needed to interpret them as text. |
| Text | A sequence of characters and spans in that text | The spelling at a position. Two occurrences of the same letter have distinct positions. |
| Concrete syntax (CST) | Groups, tokens and trivia such as spaces and comments | How the spelling is grouped by the grammar, including punctuation and trivia. |
| Abstract syntax (AST) | Declarations, formal parameters, fields and type-use occurrences | That a field has a named type or a parameter has a default. The spelling alone does not identify the declaration a name denotes. |
| Declaration semantics | Entities, binders, terms and established bindings in a captured environment | Which binder a use denotes and which effective argument an application supplies, including an argument omitted in the text. |

These levels describe different observations of a program. Storing an AST as
JSON or as a table of nodes changes its encoding, not its abstraction level.
File paths, native modules and dependency selections also matter: a file path
locates text; a module and its lexical scopes determine which names are visible.
Those are different identities even in a language where file layout affects
module membership. Bittype's [Syntax and Files dimensions](source.md) describe
how whole units are organized; they are not extra columns that every token must
pass through.

## Page across the five levels

Consider these TypeScript declarations, with no outer declaration also named
`Page` or `T`:

```typescript
type Page<T = string> = { first: T; last: T };
type TextPage = Page;
```

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

Each signature below names its two endpoint types, in reading order. The
cardinality is for the indicated endpoints **under the stated conditions**.
These signatures describe facts, not Go methods or a universal lookup API.
Reading syntax establishes spelling and grammatical roles. TypeScript name
resolution and type analysis establish the binding and defaulting facts that
generation can consume.

| Relation | Cardinality for this example | Required context | What the relation establishes |
| --- | --- | --- | --- |
| `encodedBy(CharacterOccurrence, ByteSpan)` | Each Unicode scalar occurrence has one contiguous span of 1–4 bytes. | One document, decoded as valid UTF-8. | Which bytes encode that occurrence. The four letters of `Page` cover four such spans. |
| `spelledBy(ConcreteTokenOrTrivia, TextSpan)` | Each original token or trivia leaf has one contiguous span. | The parsed text and its coordinate system; this claim is about original leaves. | Where the token or space was written. A synthetic token need not have an original span. |
| `nameToken(SyntaxNameOccurrence, IdentifierToken)` | Each simple written name has one token. | The AST and CST for the same parse. | Which token spelled this particular name occurrence, irrespective of what the name denotes. |
| `introduces(TypeParameterDeclaration, SemanticBinder)` | This formal introduces one binder. | The captured declaration and TypeScript's binding rules. | The declaration of `T` introduces `Page.T`. A use of `T` does not introduce it again. |
| `denotes(TypeNameUse, SemanticEntity)` | Each successfully resolved simple use denotes one entity here; many uses may denote it. | The use's lexical scope, selected dependencies and type namespace in that capture. | Both field uses denote `Page.T`; their spelling alone cannot establish this. |
| `usesDefaultFrom(EffectiveArgumentOccurrence, DefaultTypeSyntax)` | This defaulted argument has one written default; many applications can use that default. | The resolved callee, its formal and the language's defaulting rules at the application. | The effective string argument at `TextPage` comes from the default on `Page.T`, at a different site. |

A relation can be followed in either direction, but it need not have a unique
inverse. Starting at the binder finds both field uses; starting at either
field use finds one binder. A relation is also allowed to cross sites: there
is no reason to fabricate an argument node at `TextPage` to make a row look
complete.

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

On reading, a token can establish a name's spelling; establishing its binding
requires the captured scope and environment. On writing, a selected mapping
can specify the *intended* target of a use
before its spelling or file has been chosen. For `Page<Cell<Person>>`, moving
`Person` into the file containing the alias can remove its import without
changing which `Person` the alias must denote. The external `Cell` still
needs its package import in that TypeScript example.

The intended destination and the binding established for the organized
target program answer different questions. The latter must be checked against
the former under the actual target environment. Neither the spelling nor the
absence of an import establishes that the two agree. See the
[planner's pagination walkthrough](generation-planner.md#pagination-on-read-and-write)
and the [Files placement example](source.md#from-generated-declarations-to-files).
