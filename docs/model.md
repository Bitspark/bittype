# Horos canonical types and Bittype generation

**Revision 4 implemented in Go source, 4 November 2026.** This page describes
general kinds and their Bittype integration using public Horos Go 0.10.2.
Generic carriers, Syntax/Files composition and native pagination retain their
existing ownership boundaries.

## Ownership

| Component | Responsibility |
| --- | --- |
| ontos | Ground values, structural identity, data embeddings and encoding |
| horos | Ground descriptors/conformance; canonical modules, general kinds, binding, application and ground extraction |
| optional ontos/wire | Mixed values containing sending capabilities |
| bittype | Readable syntax, AST/CST projections, semantic providers and native type generation |
| bitwire / bitruntime | Capability transport and hydration |
| Contract libraries | Function, Stream, Cell and other behavioral declarations and realizations |

Horos's [descriptor specification](https://github.com/Bitspark/horos/blob/main/docs/spec/horos.md)
and [declaration specification](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md)
own canonical meaning. Bittype consumes that meaning through its
[Horos frontend](horos.md). Native extraction is not a competing authority.

The [general-kind design](kinds.md) replaces closed Type/Data bounds and the
dedicated Data evaluator. Users define, import, parameterize and combine kinds
through ordinary declarations. `Type`, `Data`, `Scalar`, `List` and `Record` are
ordinary names, not semantic switches in the checker, engine or generator.
Their standard definitions are maintained in Horos source and lowered through
Bittype's [ordinary module encoder](../tools/horos-module/README.md). Horos
packages the resulting canonical asset and loads it through the same declaration
machinery as consumer modules.

The [agnostic graph](semantics.md) retains entities, bindings, applications and
subject-associated facts. Typed `Semantics[L]` and its sealed envelope preserve
language/profile/marker identity. They supply storage and association contracts,
not Horos's kind algebra. [Analysis](analysis.md) consumes explicit resolution
evidence rather than invoking another resolver.

## Canonical revision 4

The module envelope remains:

```text
Tuple(Atom(UTF8("horos/decl/module")),
      int(4), text(moduleName), importsMap, declarationsMap)
```

There is no schema-digest field or separately managed formation artifact.
Imports retain an expected module name and exact module identity. One declarations
map contains kind definitions and type definitions. In Horos's notation:

```text
KindDefinition = D("kind", L(KindParameter, ...), resultSort, kindBody)
KindParameter  = D("kind-parameter", T(name), sort)
TypeDefinition = D("type", L(TypeParameter, ...), resultKind, typeBody)
TypeParameter  = D("parameter", T(name), kind)
```

`D` prefixes a tuple with `horos/decl/`; `T` and `L` are text and list embeddings.
Sorts distinguish proper types from constructor values. Kind expressions provide
top/bottom, references, applications, unions/intersections, structural node and
child-sequence patterns, recursive groups and arrow contracts. Horos owns their
canonical equations and formation rules.

The language revision fixes primitive metadata, type-child positions and operand
contracts. The bundled `horos.kinds` module supplies ordinary definitions used
by those contracts. Authors explicitly import it to use its names. A local kind
called `Data` cannot replace the identity selected by a primitive contract.

Parameters retain ordered binder slots. References preserve their defining
module; argument closures preserve caller bindings. Applications have expression
callees and fully supplied, nonempty argument lists. A type definition may return
a constructor, and a zero-parameter definition may name one. Its result kind must
express the constructor sort explicitly; parameter-list length does not determine
that result sort.

Only designated descriptor type positions admit declaration expressions. An
`exact` payload remains literal even when it resembles a parameter or kind form.
Ground `rec`/`ref` keep lexical scope and runtime-value descent. Kind recursion
has a separate structural descent rule.

## Pagination and general constraints

This [readable source](horos-syntax.md) uses an explicit prelude import. Replace
the placeholder with the supplied module's actual identity:

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

An omitted type parameter or result annotation elaborates to `top(P)`, not an
implicit reference to `Type`. Optional fields distinguish absence from a present
value and do not add null.

The same retained Page body can establish requested contracts such as
`h.Data -> h.Data` and `h.Scalar -> Rows`. `Flat -> Rows` does not generally hold:
a list-valued input makes `items` a nested list. These are general kind queries,
not a dedicated Data summary or a principal inferred kind.

For abstract `F: Flat -> Rows` and `T: Flat`, `F<T>` establishes Rows.
`F: Flat -> Flat` cannot guarantee Rows because it may return text. A wrapper
`{ value: F<T>; }` establishes `h.Record<Rows>` without thereby establishing Rows.

Known constructors retain checked bodies and definition environments. Requested
arrow contracts must respect original admission before deriving results under
fresh rigid assumptions. Abstract constructors expose declared guarantees only.
Every open body and argument is admitted, including ignored arguments. A narrower
query cannot repair an invalid declaration.

## Checked values and ground extraction

Horos supplies sealed Program-owned kind handles and type-level values. Ordinary
membership and inclusion queries establish guarantees without changing declaration
identity. Application may return a constructor or a proper symbolic value;
materialization is separate. Foreign or zero handles confer no authority, even
when another Program contains identical modules.

Results distinguish Proven, Refuted, Unproved, Invalid and Exhausted. Missing
facts assert nothing. Required obligations must be established before publication
of a checked Program or Bittype semantic model. Exhaustion is incomplete work,
never nonmembership or permission to substitute a broader type.

Ground extraction requires a closed proper value, membership in the ordinary
ground kind, materialization and unchanged ground Horos well-formedness. `any`,
complement, literals and ground recursion retain their data-only meanings. Wire
remains a declaration-layer capability form. A mixed type does not become a
ground descriptor merely because it is uninhabited. The old Data-specific
extraction/proof route is replaced, not retained as a fallback.

## Syntax, Files and native output

```text
canonical Horos modules -> authoritative Horos resolution
  -> Horos ASTs + explicit typed Semantics
  -> selected mappings -> closed native intent -> checked generation.Output
  -> target Files organization where implemented -> CSTs -> native source
```

Readable Syntax parsing supplies CSTs and AST views. The implemented
[Files composition](source.md) retains multiple Program trees, an explicit entry,
file organization and exact semantic capture. It does not concatenate text to
grant visibility. Files formatting establishes new semantic sites and receipts
for the output collection.

The [native representation profile](semantic-generation.md) preserves open Page
and ordinary applications. It represents Horos-admitted uses and explicitly
permits erasing kind constraints from target declarations. An enforcement request
needs a specified native encoding; the current enforcement profile returns
Unsupported. Native compilation does not establish every source constraint for
arbitrary native instantiations.

The subset remains text, bool, lists, named closed records, optionals and ordinary
proper generic parameters/applications. Kind definitions are compile-time
declarations. The consumer explicitly selects public roots; native output follows their demanded
representation closure. Canonical integers,
residual constructor parameters/results/aliases, unused proper parameters, open
records and unsupported forms produce located diagnostics and no usable output.
There is no hidden specialization or omission of a failed requested root or
demanded representation. Source admission still checks the supplied reachable
module closure; a valid declaration outside the native demand closure need not
have a native mapping. See [selection and roots](kernel-composition.md#registration-and-request).

## Boundaries

Dependency direction remains Bittype's Horos tooling -> Horos -> Ontos. Generic
trees, contracts, engine and printer do not import Horos/Ontos. Native syntax and
formatters remain independent of canonical semantics. Source-owned transformation
readers consume the established facts.

Function, Stream and Cell remain ordinary libraries. Membership or equal expanded
syntax never selects their behavior or native bindings. Optional library-owned
representation recipes remain separately scoped and selected by exact declaration
identity. This types-only migration adds no runtime adapters, validators, codecs,
service implementations, transport or hydration behavior.

The [unified migration](unified-declarations.md) removes superseded readers,
Claim/Interface, DataSummary and the dedicated Data evaluator. It retains no
compatibility language ID, copied Page, template lifting or executable legacy path.
