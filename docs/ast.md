# Generic AST carrier

**Status, 31 October 2026:** the Go data types described here are implemented.
Language grammars, validation, traversal and builders are not implemented.
The canonical type language remains the [model proposal](model.md).

## Language parameter

The shared package is
`github.com/Bitspark/bittype/ast/go`, imported as `ast`.
It defines `AST[L]`, `Node[L]` and `Kind[L]`. A language package supplies
a distinct named Go type for `L` and owns its node vocabulary.

For example, these declarations belong to an independent consumer:

```go
type ExampleLanguage struct{}

const Declaration ast.Kind[ExampleLanguage] = "declaration"

var tree ast.AST[ExampleLanguage]
```

The marker needs no methods, registration or core change. Two distinct named
marker types produce distinct instantiated AST, node and kind types. A type
alias to an existing marker does not create another language. Converting or
rebuilding nodes explicitly is still possible; this is not a security boundary.

The core contains no canonical, TypeScript, Go or Rust kinds. Later language
packages will provide typed builders, views and grammar checks over these
carriers. A target language may describe native type syntax such as a function
signature without introducing a privileged behavioral contract in the core.
Function, Stream and Cell remain ordinary library contracts.

## Storage and links

The [public declarations](../ast/go/ast.go) use a flat node arena:

```go
type AST[L any] struct {
    ID    TreeID
    Root  NodeID
    Nodes []Node[L]
}
```

`TreeID` is an opaque string chosen by the owner of a tree collection.
It identifies that particular AST; it is not a semantic declaration identity,
content digest, import pin or filesystem location.

`NodeID` is a one-based index: node 1 is `Nodes[0]`. Zero means no node.
A populated AST has a nonempty tree ID and a root in `1..len(Nodes)`.
The zero AST is an empty carrier, not a checked empty language module.
The `uint32` index limits one carrier's capacity; it does not limit the
canonical language's value domains.

A [node](../ast/go/node.go) has:

| Field | Purpose |
| --- | --- |
| `Kind Kind[L]` | Language-owned syntax discriminator |
| `Attributes []Attribute` | Ordered, named textual syntax atoms |
| `Children []Child` | Ordered containment edges with language-owned roles |
| `References []Reference` | Resolved, non-containment edges |
| `Origin Origin` | Optional source and earlier-node provenance |

An attribute has a `Name` and `Text`. The grammar defines their meaning,
allowed spelling and uniqueness. Text preserves identifiers or literal
spellings without converting all numbers to floating point. It is not a raw
executable-source escape hatch or an application-value codec.

Each child has a `Role` and a local `NodeID`. Repeated roles can represent
lists. A well-formed containment tree has no child cycles, has one parent per
non-root node and makes every arena node reachable from its root. References
can form graphs, including recursive declarations; they do not create extra
containment parents.

A reference has a role and a `NodeRef{Tree, Node}`. The address is explicitly
tree-qualified, even for a reference within the same AST. An unresolved name
remains syntax; it is not represented by a fabricated resolved node address.
Later resolution must check that a target exists and has the language, kind
and declaration role expected by the referencing grammar.

These invariants are requirements for future checking, not guarantees made by
the exported fields. In particular, a plain `NodeID` cannot prove which
arena it came from. The carrier does not establish import authority or
declaration identity simply by storing an address.

## Provenance

[Origin metadata](../ast/go/origin.go) can contain multiple source spans and
multiple contributing node references. This supports a target declaration
derived from several canonical nodes, or several target nodes derived from
one canonical declaration.

A span is `SourceID` plus a half-open UTF-8 byte range `[Start, End)`.
Offsets must fit the caller-supplied source text; zero-length spans are valid.
An empty span slice means no source span. The AST package neither opens files
nor interprets source identifiers as URLs.

Provenance references deliberately have no language type parameter: a target
node can refer back to a canonical node. The owner of the associated tree
collection must preserve and resolve those addresses. Provenance is separate
from a declaration's semantic identity and does not establish a trust claim.

## Ownership and extension

These structs and slices are mutable. Pipeline stages must treat their input
and reachable storage as read-only and allocate independent mutable outputs.
Copying an AST struct alone does not copy its slices. There is no implemented
freeze operation, allocator, resolver or immutable builder.

Adding another language requires its own marker, grammar and stages, without
editing `ast` or `pipeline`. Adding an ordinary declaration such as
`Page<T>` requires neither a new language nor a new core kind. The canonical
grammar will express its declaration, generic binder, references and application
using the [proposed model](model.md).

The current [contract tests](../integration/go/contracts_test.go) type-check
a separately authored consumer against the actual public declarations. They
also require the Go checker to reject mixed language ASTs, nodes and kinds,
and incompatible stage inputs and outputs. These checks prove API separation,
not validity of a language grammar or correctness of a future generator.
