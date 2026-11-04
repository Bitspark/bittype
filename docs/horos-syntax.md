# Readable Horos declarations

**Owner direction, 4 November 2026.** Bittype implements this source spelling of
Horos revision 2 through `parse.Provider`, `format.Provider` and `print.Provider`.
The [Horos frontend](horos.md) retains the existing AST and semantic providers.
This is syntax for the existing canonical model, not a second type algebra.
Horos owns [descriptor meaning](https://github.com/Bitspark/horos/blob/v0.9.0/docs/spec/horos.md)
and [declarations, bounds and specialization](https://github.com/Bitspark/horos/blob/v0.9.0/docs/spec/declarations.md).

## Pagination

A `.horos` document normally contains one module:

```text
module example.pagination;

type Page<T> = {
    items: list<T>;
    next?: text;
};

type DataPage<D: Data> = Page<D>;
```

`T` defaults to bound `Type`; `<T: Type>` spells the same canonical parameter.
`Data` restricts an argument to a data-only type. Both bounds are explicit atoms
in the canonical value. The body of Page is one closed record, with required
items and an optional next field. Optional means absent or present; it does not
add null to the field's type. `DataPage` applies the same Page declaration.

```text
module example.people;

type Person = { name: text; age?: int; };
type Page<T> = { items: list<T>; next?: text; };
type PeoplePage = Page<Person>;
```

This source parses without specializing Page. Explicit resolution establishes
the applications; the existing Horos model can then specialize PeoplePage.
Neither source parsing nor formatting invokes that semantic operation.

## Modules, imports and names

```text
module example.people;
import pages from example.pagination pin "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef";
import cells from example.cell pin "sha256:abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789";

type Person = { name: text; };
type PeoplePage = pages.Page<Person>;
type EditablePeoplePage = pages.Page<cells.Cell<Person>>;
type StoredPeoplePage = pages.DataPage<Person>;
```

The example digests are illustrative, not published module identities. Actual
imports contain exactly 64 hexadecimal digits after `sha256:`. Formatting uses
lowercase. The alias, expected module name and digest all become canonical data.
Parsing does no file or network IO and requires no dependency to be present.
Imports precede declarations within each module. Every declaration is exported.
Duplicate import aliases, declaration names and parameter names are rejected.

Names use `[A-Za-z_][A-Za-z0-9_]*`; module names contain dot-separated segments.
Backticks escape reserved identifiers, not arbitrary text: `` `type` `` is a
name, but `` `some-name` `` is not. Fields may use quoted Unicode text as names.
Declaration and import-alias namespaces are separate. A bare name is a local
zero-argument application; `alias.Name<...>` is imported. Generic arguments are
ordered. Parameter uses bind to the current declaration's ordered slots.

```text
type T = text;
type Example<T> = tuple<T, self.T>;
type `list`<`text`> = `text`;
type Escaped = `list`<text>;
```

`self.Name` forces a local declaration application when a parameter shadows it.
A parameter cannot be applied as a constructor. Generic declarations cannot
have empty parameter lists, defaults, value parameters or higher-kinded slots.
An application may spell an explicit empty argument list; formatting omits it.
There are no function bodies or executable expressions.

An existing Horos AST Program can contain several supplied modules. To preserve
that Program without inventing another carrier, a source document can concatenate
module sections, each beginning `module name;`. The first section selects the
entry; later sections supply dependency documents. They do not create imports
or grant visibility. Formatter output puts the selected entry first and keeps
the other supplied sections. No module is silently dropped, fetched or inlined.
Each module has its own canonical pin; there is no extra bundle digest.

## Type expressions

| Source | Canonical descriptor |
| --- | --- |
| `any`, `atom` | `any()`, `atom()` |
| `text`, `int`, `bool` | `scalar(utf8-text)`, `scalar(int)`, `scalar(bool)` |
| `tuple<A, B>` | `tuple(A, B)`; `tuple<>` is the empty tuple type |
| `compound(#"item", A, B)` | labelled compound with literal Atom label |
| `list<T>` | `list-of(T)`, the Ontos list embedding |
| `seq<D>` | `seq(D)`, an arbitrary-length bare tuple |
| `set<D>` | `set-of(D)` |
| `map<K, V>` | `map-of(K, V)` |
| `{ name: A; next?: B; }` | closed record, required/optional fields in authored order |
| `open { name: A; }` | open record |
| `A \| B` or `union<A, B>` | `union(A, B)` |
| `A & B` or `and<A, B>` | `and(A, B)` |
| `not<D>` | complement within the data domain |
| `wire<T>` | capability carrying messages of type T |
| `rec(#"tree", E)` | guarded recursive data descriptor |
| `ref(#"tree")` | recursion reference, distinct from a declaration application |

`&` binds more tightly than `|`. Parentheses group expressions. A chain such as
`A | B | C` has three operands; `(A | B) | C` retains a nested union. Formatting
preserves this distinction with parentheses, including equal-precedence nesting.
`union<>`, `union<A>`, `and<>` and `and<A>` preserve empty and singleton forms.
No operand sorting, deduplication, union flattening or semantic normalization is
performed. Record-field order also remains part of the canonical document.

Literal-operand constructors retain their explicit spellings:

```text
type SmallBlob = atom-len(0, 4096);
type Ready = atom-exact(#"ready");
type Operation = atom-in(set(#"read", #"write"));
type ShortName = text & text-len(1, 80) & text-non-blank();
type Key = key-grammar();
type UUID = lowercase-uuid();
type NullableText = text | exact(null);
type Tree<D: Data> = rec(#"tree", {
    value: D;
    children: list<ref(#"tree")>;
});
```

The parser enforces the spelling and arity of these forms. Horos semantics checks
their operand domains, length ranges, recursion guards, generic arity and bounds.
For example, `map<wire<text>, text>` and `atom-len(-1, 2)` parse but fail semantic
validation. Open-record fields, seq/set/map operands, intersection/complement and
recursion retain their existing Data restrictions. `any` continues to mean ground
data, while bound Type also admits capability-containing types.

Function, Stream and Cell are ordinary declarations used through imports and
applications. They have no special parser kind, formatter rule or engine entry.
The text stages include no native generators or runtime bindings. Separate
[semantic pagination transformers](semantic-generation.md) consume their admitted ASTs.

## Literal values

Literal operands are Ontos data, not type expressions:

| Source | Value |
| --- | --- |
| `"hello"` | canonical Ontos text embedding |
| `#"hello"` | Atom containing the UTF-8 bytes |
| `#hex"00ff"` | Atom containing arbitrary bytes |
| `123`, `-123` | canonical arbitrary-width integer embedding |
| `true`, `false`, `null` | existing boolean/null embeddings |
| `tuple(v1, v2)` | raw Tuple; `tuple()` is empty |
| `set(v1, v2)` | canonical set embedding; duplicate members rejected |

Atoms and tuples can represent every ground Ontos value, including embeddings
without dedicated sugar. Literal tuple syntax uses parentheses; tuple type syntax
uses angle brackets. There are no implicit type literals outside `exact(...)`.
There is no decimal scalar added by this notation. Canonical sets use Ontos's
existing ordering, so their source member order does not affect their value.

Strings are double-quoted Unicode, with Go-style `\a`, `\b`, `\f`, `\n`, `\r`,
`\t`, `\v`, `\\`, `\"`, `\uNNNN` and `\UNNNNNNNN` escapes. Escapes must denote
valid Unicode scalar values; byte/octal escapes and surrogate values are rejected.
Use hexadecimal atoms for non-UTF-8 bytes. Integers are decimal without leading
zeroes, a plus sign, negative zero, separators or exponents. This frontend limits
one integer token to 4096 digits; the model's integer domain is unchanged.

```text
type Literal<T> = exact(tuple(#"horos/decl/param", 0));
```

The tuple inside exact remains literal data. It does not capture T or become a
parameter use. Formatters recognize canonical data embeddings only to choose
readable literal spellings; other values retain their raw atom/tuple structure.

## Grammar outline

```text
document    = module-section+
module-section = "module" module-name ";" import* declaration*
import      = "import" name "from" module-name "pin" string ";"
declaration = "type" name [ "<" parameter ("," parameter)* [","] ">" ] "=" E ";"
parameter   = name [ ":" ("Type" | "Data") ]
E           = intersection ("|" intersection)*
intersection = primary ("&" primary)*
primary     = application | parameter-use | "(" E ")" | constructor
application = (name | "self." name | name "." name) [ "<" [E ("," E)* [","]] ">" ]
record      = ["open"] "{" (field-name ["?"] ":" E ";")* "}"
```

Constructor forms and literal positions are the tables above. Comma-delimited
angle-bracket lists and literal tuples/sets admit a trailing comma; fixed-arity
parenthesized constructor calls do not. Compound accepts its literal label and
zero or more type arguments, without a trailing comma. Semicolons are required.
Whitespace, a leading UTF-8 BOM, `//` line comments and non-nesting `/* */` block
comments are admitted. Comments are trivia, never directives or canonical data.

## Pipeline, ownership and limits

```text
parse/horos/go:   text -> CST[Horos], with stored AST view and exact module pins
format/horos/go:  AST view[Horos] -> freshly constructed CST[Horos]
print/go:        CST[L] -> supplied leaf text
```

`provider/horos/go` composes these stages and the existing checker/resolver/analyzer.
Each capability is independently registered. No engine, shared contract or Wire
change is needed. The [provider matrix](providers.md) records availability.
The private canonical bridge reuses existing module encoding and Horos Pin;
hashing a document is not semantic validation. Source ASTs carry no invented
resolved target annotations. Resolution retains its explicit evidence contract.

`Print(Parse(source)) == source`, including comments, BOM and CRLF. Formatting
chooses four-space indentation and LF, drops source trivia, emits Type bounds
implicitly and preserves the canonical module values and pins. Module/import/
declaration map order is semantically immaterial; parameter, field and operand
order is not. `Parse(Print(Format(ast)))` cannot recreate original node IDs,
provenance or resolution receipts. The formatter's stored AST snapshot retains
that metadata, but callers must resolve the freshly parsed input explicitly.

Formatting constructs token/trivia nodes directly. It never emits source and
reparses it. The generic printer checks only carrier structure and knows no Horos
constructor. Unknown descriptor spellings, unrepresentable binder indices and
incorrect module-pin assertions fail formatting with no partial output.

Source is bounded to 4 MiB and 200,000 lexical tokens. Syntax nesting and AST/CST
storage obey existing carrier budgets; canonical module hashing obeys Horos's
default value budgets. All failures return empty output with diagnostics.
These are resource limits, not new canonical type semantics.

## Qualification

[Parser tests](../parse/horos/go/parse_test.go) cover the complete constructor
vocabulary, lexical rejection, comments/CRLF/BOM, nested operators, literal opacity,
limits and fuzzed source round trips. [Formatter tests](../format/horos/go/format_test.go)
start from independently authored canonical pagination and Cell modules, preserve
their exact closure pins, recover the model, specialize Page/DataPage and retain
the capability/data distinction. Engine and installed-consumer tests qualify the
public stages. Existing native compiler checks remain unchanged.
