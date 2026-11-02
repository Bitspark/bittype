# Horos declaration AST

**Implemented in Go source, 2 November 2026.** This is an AST-only frontend for
Horos declaration revision 1, not a readable source grammar or native generator.
The authoritative value format and binding rules are in
[horos's declaration specification](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md).

## API and ownership

- [syntax/horos/go](../syntax/horos/go/nodes.go) defines `Language`, `AST`, `Node`,
  `Kind` and the vocabulary. It imports only generic AST and language identity.
- [resolve/horos/go](../resolve/horos/go/resolve.go) provides
  `Resolve(root, dependencies, treeID)`, `Project(program, treeID)` and `Model(view)`.
  Resolve delegates semantics to Horos, then projects its checked model. Model
  checks the borrowed AST and returns a Horos program for inspection or explicit
  `Specialize(declarationID, groundArguments)`. No lookup performs network/file IO.
- [check/horos/go](../check/horos/go/check.go) validates the AST through the same
  Horos model. It checks binding edges and pins instead of trusting user-supplied
  reference fields. It does not check application values.
- [provider/horos/go](../provider/horos/go/provider.go) implements
  `text.CheckerProvider[horos.Language]`, with ID `horos` and name `Horos`.
  Applications explicitly register it through `representation.NewLanguage` and
  `Engine.RegisterLanguage`, using profile `horos-declarations-1/ast-1`.

Only the resolver and its private projection helper depend on Horos/Ontos.
The model dependency is pinned to a public immutable Go module revision in go.mod;
there is no local replacement or vendored private checkout. Generic infrastructure
and the four native text stages are unchanged and model-independent. The engine
can hold this AST without registering any transformer. Parse and Render report
missing text capabilities for this provider.

## Vocabulary and bindings

One tree contains the complete resolved import closure. Modules are ordered by
their digest; imports/declarations retain Horos's canonical map order. Parameters,
arguments and descriptor operands retain their specified order.

| Kind | Attributes | Children | Required references |
| --- | --- | --- | --- |
| Program | none | `module*` | `root` to the entry module |
| Module | `name`, `pin` (lowercase hex) | `import*`, `declaration*` | none |
| Import | `alias`, expected `name`, `pin` | none | `target` module |
| Declaration | `name` | `parameter*`, one `body` | none |
| TypeParameter | `name`, canonical decimal `index` | none | none |
| Parameter | canonical decimal `index` | none | `target` TypeParameter |
| Application | `scope` (`local` or `imported`), `name`, imported `alias` | `argument*` | `target` Declaration; imported `import` edge |
| Descriptor | Horos `constructor` name without `horos/` prefix | `operand*` | none |
| Atom | raw bytes as lowercase `hex` | none | none |
| Tuple | none | `element*` | none |

Descriptor operands retain the canonical value layout. Only Horos-designated
type positions become Parameter, Application or Descriptor nodes. Other positions
remain Atom/Tuple fragments, including field wrappers and everything inside an
`exact` literal. The frontend takes these positions from the checked Horos model;
it has no duplicate descriptor-constructor switch or record semantics.

Every binding reference is internal to the tree's closure and must identify the
correct module, import, declaration or parameter. No attributes, kinds or roles
outside this profile are admitted. Missing references fail instead of triggering
implicit binding repair. Arena addresses may change if every edge changes with
them; attribute/reference storage order is immaterial. Containment order remains
canonical. Origins are retained as tooling metadata but do not affect module pins.

The checker reconstructs canonical module values, asks Horos to resolve them,
and compares their projected syntax and binding edges with the supplied AST.
This checks one declared AST profile; it never fingerprints a shape to choose a
behavioral library. Invalid UTF-8 attributes are rejected before text embedding.
It does not reinterpret literal data as declarations or silently normalize an AST.

## Pagination and specialization

The [pagination tests](../resolve/horos/go/resolve_test.go) consume two modules:
`example.pagination.Page<T>` and `example.people.PeoplePage = Page<Person>`.
The Page body retains a Parameter node referencing its T declaration. The
PeoplePage application references the imported Page and caller-side Person;
a same-spelled Person in the imported module cannot capture that argument.

Projection preserves the open declaration. Calling Horos's Specialize explicitly
produces the ground closed-record descriptor with required list items and optional
text next. It does not mutate the AST or the checked model. Nested applications
retain distinct same-spelled binders; parameter/application-shaped literal tuples
remain opaque data. No pagination operation or cursor behavior is introduced.

## Admission and qualification

Horos owns module-name/pin checking, arity, binding scope, recursion closure,
application-cycle rejection and specialization limits. Its errors retain module,
declaration and tuple-index path context. AST structure/profile errors are located
at AST nodes; no source location is invented when none exists.

Projection checks generic AST node, depth, entry and text budgets while creating
nodes, including hexadecimal growth before copying atoms. Generic structure is
checked before the AST is decoded. Horos's default declaration budgets also apply
when recovering a model. Both boundaries return empty output on failure; a model
accepted under custom larger limits is re-admitted under the fixed Horos defaults
before projection, so successful projection cannot bypass the checker's canonical
budgets. Canonical and AST depth are distinct; tests qualify their boundary.

Remote tests cover canonical value round trips, imported and nested pagination,
specialization, wrong/missing references, invalid pins/indices/UTF-8, literal
boundaries, address renumbering, provenance, nil views and resource exhaustion.
Independent integration fixtures qualify explicit engine registration, immutable
envelopes and absent text stages. Architecture tests enforce the narrow external
dependency boundary, and static consumer tests distinguish Horos from every native
AST marker. Existing native-text/compiler gates remain required.

Native generators, readable Horos parsers/formatters, runtime adapters and
Function/Stream/Cell bindings are outside this change. These contracts retain
library ownership; none is a special core kind or builtin registration.
