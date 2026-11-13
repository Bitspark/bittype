# Destinations: what a name use in generated output refers to

**Specification, 13 November 2026, from [research 0016](../research-docs/0016-native-files-placement.md) (rows 19, 25, 26, 34 and 36).**
It was written before Go and Rust issue symbols, as row 34 asks: "Written as a
specification before Go and Rust issue symbols". It is delivered with
[#276](https://github.com/Bitspark/bittype/issues/276). Go Files and Rust Files
do not exist yet. What they do with a destination is slice 1 of
[#161](https://github.com/Bitspark/bittype/issues/161) and of
[#159](https://github.com/Bitspark/bittype/issues/159).

The accepted decision of 7 November 2026 gives generation this task
([Syntax and Files](source.md#from-generated-declarations-to-files)): "Bind
each named use to its intended target declaration, lexical parameter, or
selected external public symbol." It gives the Files formatter of the language
this one: "establish that uses still resolve to their intended targets". This
page says what a name use is, what its destination is, and when a use that a
compiler has resolved reaches *the same* destination.

## Terms

**Name use.** An occurrence in a generated tree whose meaning depends on what a
name resolves to in the scope where Files places it. A declaration or a type
parameter of the same name can capture it. A keyword is no name use: nothing
can capture TypeScript's `string`. Each language states which nodes are name
uses ([below](#name-uses-per-language)).

**Destination.** What the use is meant to refer to. It is exactly one of three
things, and the sealed output records which:

| Destination | How the output records it |
| --- | --- |
| A declaration that the output owns | A reference on the use's node, with the role `target`, to the declaration's node |
| A type parameter of the enclosing declaration | A reference on the use's node, with the role `target`, to the parameter's node |
| An external symbol | An external use in `Output.ExternalUses()`, at the use's node, that carries a symbol of the language's issuer |

A use carries one of the two carriers and never both.

**External symbol.** A value that a language's symbol package issues and reads
back. The generation kernel and the shared Files contract carry it and never
read its key. Each language defines what a key identifies
([below](#symbols-per-language)).

## The same destination

Three statements define it. A later change may not alter them.

### 1. A destination is a declaration's identity and never structural equality

Two declarations with equal bodies are two destinations. Two aliases of one
type are two destinations. A declaration and a type that happens to be
identical to it are two destinations.

A check that a use reaches its destination therefore compares entities and not
types. It takes the entity that the compiler binds the use to and asks whether
that is the entity of the intended declaration, parameter or symbol. It does
not ask whether two types are identical, whether they print alike or whether
their structure is equal. The Go graph's `identical` fact and every `display`
fact answer another question and are not used for this.

[Syntax and Files](source.md#from-generated-declarations-to-files) states the
principle: "A spelling alone is insufficient: two distinct types named `Person`
must not become interchangeable when placement changes." The same holds for two
distinct types of one structure.

### 2. An external symbol names the selected public export and not the canonical type underneath it

The mapping, or the target for its own forms, selects one public endpoint. The
symbol identifies that endpoint. The canonical type that the endpoint leads to
is not the symbol's identity, and the same canonical type reached through
another endpoint is another destination:

- TypeScript: the export `Cell` of the module `@bitspark/cell`. If the package
  declares `Cell` as an alias of another exported type `Inner`, the destination
  is `Cell`, and `Inner` is another destination.
- Go: the identifier `string` of the universe scope. `byte` and `uint8` are two
  destinations although Go treats them as one type.
- Rust: the item at the path `string::String` of the crate `std`.
  `core::option::Option` and `std::option::Option` are two destinations
  although the second re-exports the first.

A semantic provider may have to establish the correspondence through an alias
or a re-export. It keeps which endpoint generation selected. How each bridge's
evidence relates to an endpoint is stated per language
[below](#symbols-per-language).

### 3. Provenance may be many to many

One source declaration may yield several native declarations, and one helper
may serve several source declarations or none. So a destination is never
identified by a source declaration. A reference identifies one native
declaration of the output. The inventory's `Provenance` is a list, and it is
not a key: two declarations may have the same provenance, and a helper has
none.

## Name uses per language

These are the nodes that carry a destination in a sealed output. A later change
may not alter which occurrences they are.

| Language | Name use | Not a name use |
| --- | --- | --- |
| TypeScript | Every `type-reference` node | `keyword-type` (`string`, `boolean`): a keyword. The TypeScript target writes no `qualified-type` |
| Go | Every `named-type` node. `string` and `bool` are names of the universe scope, and a declaration or a type parameter can capture them | The Go target writes no `qualified-type`. A type of another package is slice 2 of #161 |
| Rust | The last `path-segment` of every `path-type`, and every `primitive-type` node. `bool` is no keyword, and a declaration or a generic parameter of that name captures it | The segments before the last. They qualify the use and are part of how one use is spelled. The Rust target writes no `lifetime-reference` and no `qualified-path-type` |
| C++ | None is recorded. Row 36 of research 0016 excludes legacy C++ output until [#160](https://github.com/Bitspark/bittype/issues/160) | |

**Rust paths.** A path that refers to a declaration or a parameter of the
output has one segment, and that segment carries the reference. A path to an
external item is written today as an absolute path, `::std::string::String`,
and its last segment carries the external use. The first segment of an absolute
path names a crate and is resolved outside every module of the output, and each
later segment is looked up inside the item before it, so no declaration of the
output can capture them. Slice 1 of #159 changes what the fragment writes: one
segment that carries the symbol, and Rust Files writes the path
([decision 1 on #276](https://github.com/Bitspark/bittype/issues/276#issuecomment-6837794385)).
The occurrence that carries the destination is the same before and after: the
segment that names the item.

## Symbols per language

A later change may not alter what a symbol's key identifies. Each package can
gain a further category of symbol. That adds destinations and changes none.

Each package has two functions for a use: `ExternalUse` records a symbol at a
node, and `ReadExternal` reads it back and refuses a use of another issuer or
of another language.

**The role of an external use, one rule for the three languages.** The role
is the symbol package's own name. It is read back only through that package,
and nothing compares it with a role of the compiler bridge: a formatter takes
the site of a use from the symbol package and looks the compiler's binding up
under the bridge's own role. The keys are `typescript/name`, `go/target` and
`rust/target`. TypeScript's and Rust's happen to equal the role under which
their bridge publishes a binding. Go's does not: the Go bridge publishes
`go/name`. Nothing depends on either. Each package's test pins its key,
because the key is sealed into every output with the site.

### TypeScript

[`transform/syntax/typescript/symbol/go`](../transform/syntax/typescript/symbol/go/symbol.go),
unchanged. A symbol identifies a named export of a module, or a global type
when the module is empty. The bridge binds a use of an imported name to an
entity of the kind `package-type` whose name is the export's and whose owner is
the `package` entity of the module. For an export that is an alias of another
type the entity is the alias: the binding names `Cell` and not `Inner`. A
supplied package holds interfaces, aliases and named type imports, so it has no
re-export under another name.

### Go

`transform/syntax/golang/symbol/go` issues one category today.

**Predeclared.** `Predeclared(name)` identifies the type that the identifier
`name` denotes in Go's universe scope. The symbol states its category. It does
not encode "predeclared" as a package symbol with an empty package, so a
malformed package symbol can never be read as a predeclared one.

The destination is the predeclared entity itself, and not an entity that has
the name `string`. The Go bridge exports a predeclared basic type as an entity
of the kind `builtin`, with the identifier as its name, without an owner and
without a declaration site. A declaration named `string` is an entity of the
kind `defined` or `alias` with its package as owner and a site. A type
parameter named `string` is an entity of the kind `parameter`. A use that
either of them captures binds to another entity than the destination, and the
comparison of statement 1 sees it.

Limit: `builtin` covers the basic types. `error` and `comparable` are named
types of the universe scope and `any` is an alias, so the bridge exports them
with the kinds `defined` and `alias`. A predeclared symbol for one of them
needs "no owner and no site" as its test, or a published marker. Until then the
package refuses to issue it as unsupported.

A type of another package is no symbol yet. It needs that package's semantics:
slice 2 of #161. The Go model's external form stays unsupported.

Limit: `Predeclared` asks `go/types` of the toolchain that builds bittype
which names the universe scope declares. That is not the universe of the
pinned language version, Go 1.25. The basic types, which are all it issues,
are the same in both.

### Rust

`transform/syntax/rust/symbol/go` issues two categories today.

**Item.** `Item(crate, path...)` identifies the public item at `path` below the
root of a dependency. The dependency is separate from the path. A symbol is
bound to a dependency and not to a spelling: "rustc permits an external
dependency's locally visible name to differ from the compiled crate name"
(advice, row 26), so the symbol names the dependency that the environment
selects, and Rust Files chooses the local spelling. Today there is one
dependency, `Std()`: the crate `std` of the standard library of the pinned
compiler. A crate of a library is no symbol yet. It needs that crate's
semantics: [#170](https://github.com/Bitspark/bittype/issues/170).

A key holds a name and never a spelling. A raw identifier, `r#string`, is one
way to write the name `string`, so `Item` and `Dependency` refuse it as
invalid, and one item has one symbol. The underscore names nothing and is
refused too. Rust Files chooses a raw spelling where the language needs one.

How this relates to the Rust environment: the
[environment](native-semantics.md#rust-environment) refuses `std`, `core` and
`alloc` as the name of the *captured* crate, because the graph names an entity
by its crate and its path, and a captured crate of one of these names would
put two entities of one name into a graph. `Std()` is about another thing: the
crate that a use *names*. It is the one dependency a symbol may name. `core`
and `alloc` as the dependency of a symbol are unsupported: each of their items
is another export than the one of `std`, and no mapping selects one today.

Limit: `Dependency` judges a crate's name by the rule for a segment of a path.
So a keyword that the compiler takes as the name of a crate, such as `fn`, is
invalid here, where the compiler's own rule would make it a crate outside
coverage. That rule is on `main` in `semantics/rust/go`, which a symbol
package may not import. Where the one rule lives is open.

The Rust bridge names an entity by the crate and the path of its definition.
It names the entity of `::std::string::String` as `alloc::string::String`, and
that of `::std::option::Option` as `core::option::Option`. That name is the
canonical type underneath, and by statement 2 it is not compared with the
symbol. The bridge binds every segment of the written path: the first to the
root module of the crate `std`, each later one to what that name denotes in the
item before it. A use therefore reaches the destination `Item(Std(), p1, ...,
pn)` when it is written as the absolute path `::std::p1::...::pn`, its first
segment binds to the root module of `std`, and the compiler accepts the path.
That is why Rust Files realizes an external use as an absolute path in its
first slice and as nothing else.

**Primitive.** `Bool()` identifies Rust's primitive type `bool`. The helper
publishes bindings for path segments only, so a declaration or a generic
parameter named `bool` that captures a use is invisible today. `Bool()`
therefore records a destination that no evidence can establish yet.
[Decision 2 on #276](https://github.com/Bitspark/bittype/issues/276#issuecomment-6837794385)
gives the primitive node a binding from the helper. Rust Files cannot establish
this destination before the helper publishes it.

## Who enforces completeness

Row 19 of research 0016: "A destination for every name use is required by each
language's generated-Files formatter, which refuses a use without one as
TypeScript does, and is asserted by a test per target over sealed outputs".

- **The shared seal cannot.** It validates the external uses it is given: the
  site exists, the role is valid, no occurrence is used twice, the symbol has an
  issuer of the output's language. It does not know which node kinds of a
  language use a name.
- **A generated-Files formatter refuses a name use without a destination.**
  This is part of the Files contract and is stated with the
  [failure categories](source.md#failure-categories). The refusal is invalid
  input: a sealed output with such a use is malformed, and no later slice
  admits it. An external use that the language's issuer cannot read, or whose
  node is no name use or does not name the symbol, is invalid input too. An
  external destination that the formatter's evidence cannot establish is
  unsupported coverage: a later slice can admit it. Rust's `bool` is one as
  long as the helper publishes no binding for a primitive type.
- **A test per target** reads sealed outputs and asserts that every name use
  carries exactly one of a reference and an external use, and that every other
  node carries neither:
  [`destinations_test.go`](../integration/go/destinations_test.go). It
  classifies every node by its kind. Which kinds can carry a reference it asks
  of the language's syntax checker, whose grammar marks them: TypeScript
  `type-reference`, Go `named-type` and `qualified-type`, Rust `path-segment`
  and `lifetime-reference`. Rust's `primitive-type` uses a name and gets no
  reference from the grammar, so the test lists it. The kinds that use no name
  are listed per language. A node of a kind that is in no list and that the
  grammar gives no reference fails the test, so a new kind is classified
  before an output with it passes.
- **Which outputs it reads.** The pagination output of each of the three
  targets, and for TypeScript also the output with a type of a package beside
  a declaration of the same name and one whose lists are the global type
  `ReadonlyArray`. They hold every kind that the three native models can build
  today. The claim is shown for those kinds and those outputs, and not for
  every tree that a later model could build.

## What generation records, and what it does not decide

Recording a destination changes no generated text. The tree is the one it was,
and the output gains external uses. Generation still decides nothing about
placement: a recorded destination is no import, no qualification and no path.
Nobody places Go or Rust output yet.

## Limits

- Go Files and Rust Files do not exist. No publication is refused for a use
  that misses its destination, in either language.
- C++ records no destination for `::std::string`, `::std::vector`, the
  external form, the namespace segment or the includes.
- The Rust helper publishes no binding for a primitive type.
- A recorded use counts against three limits of an output. The seal refuses
  more external uses than `Limits.References`, which on `main` did not bound
  the text, boolean, list and optional uses of a Go or Rust output. Each use
  takes a second record of the target's store and the bytes of its key. On the
  pagination output the smallest limits that pass grew from 58 to 62 nodes and
  from 1558 to 1658 bytes for Go, with 4 uses, and from 64 to 70 nodes and from
  1743 to 1921 bytes for Rust, with 6 uses (measured by the independent review
  of #294). TypeScript is unchanged.
- The TypeScript symbol's encoding of a global type as an empty module is not
  changed (row 25 asks this of Go only).
- A destination says which entity is meant. It does not say that the entity is
  accessible from where Files places the use. Accessibility is the formatter's
  check under its environment.
