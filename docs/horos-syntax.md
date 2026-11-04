# Readable Horos declarations

**Revision-4 grammar implemented in Go source, 4 November 2026.**
Bittype supplies readable spelling for Horos's canonical general-kind language.
The [frontend](horos.md) owns its AST projection and stage composition;
[Horos](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md)
owns canonical meaning. This page does not claim that the migration or public
package delivery is complete.

## Modules and ordinary prelude imports

```text
module example.pagination;
import h from horos.kinds pin "sha256:<actual-64-hex-pin>";

kind Flat = h.Scalar | h.List<h.Scalar>;
kind Rows = h.Record<Flat>;
type Page<T> = { items: list<T>; next?: text; };
type DataPage<D: h.Data>: h.Data = Page<D>;
type Person = { name: text; active: bool; };
type PeoplePage = Page<Person>;
```

The placeholder is explanatory, not valid source. A real import contains exactly
64 hexadecimal digits after `sha256:` and an expected module name. Formatting
uses lowercase. The prelude uses this ordinary import form; there is no schema
header, implicit import or privileged binding for `Type` or `Data`.

An omitted type parameter annotation or result kind elaborates to `top(P)`.
It does not refer to a name `Type`. Kind-definition parameters and result sorts
default to `P`. These defaults are syntax elaboration, not kind inference.
All canonical fields remain explicit.

Imports precede declarations. One declaration namespace contains both `kind` and
`type` definitions; reference categories are determined by their expression
context. Every definition is exported. Parameter names are distinct within their
own binder. Names use `[A-Za-z_][A-Za-z0-9_]*`; module names are dot-separated
segments. `Type` and `Data` follow these ordinary naming rules. Backticks escape
reserved identifiers, not arbitrary text; record fields may use quoted Unicode
names. Import aliases and declaration names remain separate namespaces.

A bare name denotes the current binder when present, otherwise a local reference.
`self.Name` forces a local definition reference; `alias.Name` selects an imported
one. Parsing retains unresolved names and performs no network or filesystem IO.
Resolution determines visibility, declaration category and semantics.

A Syntax document may contain several `module name;` sections. The first section
selects the entry and later sections supply module values. They do not imply
imports or visibility. Each module retains its own identity; there is no bundle
identity. The [Files composition](source.md) also supports separate Program trees
without concatenating their source text.

## Definition forms and sorts

```text
kind Name<K: S, ...>: S = K;
type Name<T: K, ...>: K = E;
```

Here `S` denotes a sort, `K` a kind expression and `E` a type-level expression.
The metavariable letters are not source identifiers with special meaning.

Sorts have the spelling:

```text
P
P -> P
(P, P) -> P
(P -> P) -> P
P -> (P -> P)
```

`P` is the proper-type sort. An arrow sort has a nonempty ordered input list and
one result sort, which can itself be an arrow. Arrows associate to the right.
`(P, P) -> P` is a two-input constructor sort, not a curried unary sort.

Kind arguments and type arguments occupy different levels. A parameterized kind
definition consumes kinds with its declared sorts. A type definition consumes
type-level values satisfying its parameter kinds. Named kind applications are
exactly saturated; a bare kind reference names a zero-parameter kind definition.
Empty parameter and application lists are rejected.

## Kind expressions

| Source form | Meaning |
| --- | --- |
| `top(S)`, `bottom(S)` | All or no type-level values at sort S |
| `K`, `self.K`, `h.K` | A kind parameter or ordinary local/imported kind reference |
| `h.List<K>` | Saturated application of a named kind definition |
| `A \| B` | Kind union at one sort |
| `A & B` | Kind intersection at one sort |
| `A -> B`, `(A, B) -> C` | Unary or n-ary constructor contract |
| `node("form", selectors, Q)` | Primitive shape, metadata constraints and full child-sequence pattern |
| `nu { X = K; Y = K; } in X` | Guarded simultaneous recursive proper-kind group and selected entry |
| `(A \| B) at S` | Explicit sort on that Boolean kind expression |

Kind application binds most tightly, then `&`, then `|`, then right-associative
`->`. Parentheses preserve explicit grouping. A Boolean expression defaults to
sort `P`. An `at` annotation applies only to that Boolean node: it does not
propagate into nested Boolean nodes or infer sorts from names. For a Boolean
combination of constructor contracts, state its sort explicitly:

```text
kind Joint: P -> P =
    ((h.Scalar -> h.Scalar) & (h.Data -> h.Data)) at (P -> P);
```

A kind union is an alternative guarantee; an intersection requires all its
members. In particular, an abstract constructor satisfying a union of arrows
cannot be treated as satisfying both arrows. The [kind model](kinds.md) defines
application admission and result reasoning. Public kind negation, arbitrary
predicate callbacks and subtree-equality constraints are outside this revision.

### Structural patterns

```text
kind Text = node("scalar", {"selector": "utf8-text"}, epsilon);
kind List<K> = node("list-of", {}, item(K));
kind ClosedRecord<K> = node("record", {"mode": "closed"}, star(item(K)));
kind Pair<A, B> = node("tuple", {}, concat(item(A), item(B)));
```

Selector keys and values are strings. An omitted selector is a wildcard, not a
request to invent metadata. The initial selectors are `selector` on `scalar`
(`int`, `utf8-text`, `bool`) and `mode` on `record` (`open`, `closed`). Horos owns
which selectors apply to each primitive form. Other descriptor metadata remains
preserved and subject to its existing admission rules.

Child-sequence patterns are:

```text
Q = epsilon
  | item(K)
  | concat(Q, Q, ...)
  | alternative(Q, Q, ...)
  | star(Q)
```

`concat` and `alternative` require at least two patterns. `item` contains a proper
kind. A node pattern matches the entire ordered list of observed type children.
`star` admits arbitrary finite repetition, never an infinite tuple or record.
An open record exposes an implicit `any` child for additional fields, along with
all named fields, including optional ones. Labels and presence markers remain
descriptor metadata; a uniform field-kind query does not introduce row variables
or predicates over labels.

### Recursive kinds

```text
kind Tree = nu {
    X = h.Scalar | h.List<X> | h.Record<X>;
} in X;
```

Every group name is in scope throughout the group. Inner recursive groups shadow
outer names. Canonical recursive references use lexical group depth and slot;
formatting freshens names when needed to preserve those references. The group
and its entry have proper sort.

Recursive kinds use greatest fixed points. Every recursive cycle must cross a
node's type-child descent. Boolean operators and references alone are not guards.
Named kind-definition dependencies remain acyclic, and parameters stay fixed
within a recursive group. Ground descriptor `rec`/`ref` has a separate lexical
and runtime-value descent rule; kind syntax does not replace it.

## Type-level references and constructor results

```text
type ListType<T> = list<T>;
type Alias: top(P) -> top(P) = ListType;
type Later<F: top(P) -> top(P)>: top(P) -> top(P) = F;
type Texts = Later<Alias><text>;
```

`Alias` has no declared parameters but denotes a constructor. `Later` returns
one. Each postfix application supplies one complete, nonempty argument list to
its callee value; parameter references and earlier application results can be
callees. Partial application and type-level lambdas remain outside this revision.

`type Alias = ListType;` does not infer a constructor result: its omitted result
kind is `top(P)` and resolution rejects the mismatch. Empty application syntax
such as `Person<>` is rejected. A zero-parameter proper definition is used by
reference, `Person`. Descriptor forms such as `tuple<>` and `union<>` retain their
own empty-operand meanings and are not empty declaration applications.

Kinds and type definitions can share spellings only through normal scopes, not
through duplicate declarations in the same namespace. For example:

```text
type T = text;
type Example<T> = tuple<T, self.T>;
type `list`<`text`> = `text`;
type Escaped = `list`<text>;
```

There are no executable function bodies, value parameters or inferred missing
higher-order annotations.

## Type descriptors

| Source | Canonical descriptor |
| --- | --- |
| `any`, `atom` | Ground `any()`, `atom()` |
| `text`, `int`, `bool` | `scalar(utf8-text)`, `scalar(int)`, `scalar(bool)` |
| `tuple<A, B>` | `tuple(A, B)` |
| `compound(#"item", A, B)` | Compound with literal Atom label |
| `list<T>` | `list-of(T)`, the Ontos list embedding |
| `seq<D>`, `set<D>`, `map<K, V>` | Bare sequence, embedded set, embedded map |
| `{ name: A; next?: B; }` | Closed record with required/optional fields |
| `open { name: A; }` | Open record |
| `A \| B` or `union<A, B>` | Type union descriptor |
| `A & B` or `and<A, B>` | Type intersection descriptor |
| `not<D>` | Complement in the ground data domain |
| `wire<T>` | Sending-capability type with message type T |
| `rec(#"tree", E)`, `ref(#"tree")` | Ground recursive descriptor and lexical reference |

In type-expression context `&` binds more tightly than `|`. A chain has several
operands, while parentheses preserve nested union/intersection structure.
`union<>`, `union<A>`, `and<>` and `and<A>` retain empty/singleton forms. No
sorting, deduplication or semantic flattening occurs. Record-field order remains
part of the canonical module value.

Literal-operand forms retain their spelling:

```text
type SmallBlob = atom-len(0, 4096);
type Ready = atom-exact(#"ready");
type Operation = atom-in(set(#"read", #"write"));
type ShortName = text & text-len(1, 80) & text-non-blank();
type Key = key-grammar();
type UUID = lowercase-uuid();
type NullableText = text | exact(null);
type Tree<D: h.Data> = rec(#"tree", {
    value: D;
    children: list<ref(#"tree")>;
});
```

The last example assumes the ordinary `h` prelude import. Horos establishes
primitive operand contracts using ordinary kinds. Existing seq/set/map,
open-record, intersection/complement and recursion restrictions remain in place;
no reserved Data flag decides them. `any` still describes ground data, whereas
`top(P)` admits all proper type descriptions, including mixed ones. Parsing and
formatting do not perform membership, inclusion or specialization.

## Literal values

| Source | Ontos value |
| --- | --- |
| `"hello"` | Text embedding |
| `#"hello"` | Atom containing UTF-8 bytes |
| `#hex"00ff"` | Arbitrary-byte atom |
| `123`, `-123` | Arbitrary-width integer embedding |
| `true`, `false`, `null` | Existing boolean/null embeddings |
| `tuple(v1, v2)` | Raw Tuple; `tuple()` is empty |
| `set(v1, v2)` | Canonical set embedding; duplicate values rejected |

Literal tuples use parentheses; tuple types use angle brackets. Literal atoms
and tuples can represent all ground Ontos values. There are no implicit type
literals outside literal-operand positions. No decimal scalar is introduced.

Strings use double quotes and Go-style `\a`, `\b`, `\f`, `\n`, `\r`, `\t`, `\v`,
`\\`, `\"`, `\uNNNN` and `\UNNNNNNNN` escapes for Unicode scalar values.
Byte/octal escapes and surrogate values are rejected; use hexadecimal atoms for
non-UTF-8 bytes. Integer spelling is decimal without leading zeros, plus signs,
negative zero, separators or exponents, with at most 4096 digits per token.

```text
type Literal<T> = exact(tuple(#"horos/decl/param", 0));
```

The tuple is literal and does not capture T. The same opacity applies to tuples
resembling kind, application or recursive-variable forms.

## Grammar outline

```text
document        = module-section+
module-section  = "module" module-name ";" import* definition*
import          = "import" name "from" module-name "pin" string ";"
definition      = kind-definition | type-definition
kind-definition = "kind" name ["<" kind-parameters ">"] [":" S] "=" K ";"
type-definition = "type" name ["<" type-parameters ">"] [":" K] "=" E ";"
kind-parameter  = name [":" S]
type-parameter  = name [":" K]
S               = "P" | "(" S ")" | S "->" S | "(" S "," S ... ")" "->" S
K               = kind-reference/application | top/bottom | kind-Boolean
                | kind-arrow | node-pattern | recursive-group | "(" K ")"
recursive-group = "nu" "{" (name "=" K ";")+ "}" "in" name
E               = type-Boolean | type-primary ("<" E ("," E)* [","] ">")*
```

This outline uses the precedence and arity rules above. Parameter lists are
nonempty. Type/kind argument and parameter lists, and selector maps, admit
trailing commas. Arrow domains and child-sequence pattern argument lists do not.
Fixed parenthesized descriptor calls retain their specific arities. Semicolons are
required. Whitespace, a leading UTF-8 BOM, line comments and non-nesting block
comments remain trivia, not directives or canonical data.

## Pipeline and limits

Syntax parsing constructs a lossless CST and AST view. Formatting constructs CST
from an AST view directly; the generic printer emits supplied leaves. The
provider composes independently registered Syntax/check/resolution/analysis
capabilities. [Files parsing and formatting](source.md) compose these stages over
several Program trees without creating a second language.

Direct parsing/printing retains original source, comments, BOM and line endings.
Formatting uses four-space indentation and LF, removes source trivia and preserves
canonical values. Default annotations may be omitted only when they elaborate to
the same canonical sort/kind. Recursive names can change while their canonical
lexical references remain fixed. A freshly parsed formatted document requires
new resolution; source text does not recreate old node IDs or receipts.

Source is bounded to 4 MiB and 200,000 lexical tokens. Syntax/AST/CST and canonical
operations retain their separate limits. Failures return no usable partial
output. Function, Stream and Cell remain ordinary declarations with no parser or engine
special cases.
