# Horos canonical types and bittype generation

**Updated 4 November 2026. Status: Go declaration model and Horos AST implemented;
native pagination mappings implemented; broader native representations proposed.** This replaces the 31 October
proposal for an independent bittype value-type language. The generic AST/CST,
native text stages and syntax engine already exist. See the [Horos frontend](horos.md).

## Ownership and canonical input

| Component | Responsibility |
| --- | --- |
| ontos | Ground data values, structural identity, data embeddings and canonical encoding |
| horos | Ground descriptors/conformance; unified modules, kinds/interfaces, binding, Data certificates and specialization |
| optional ontos/wire | Go/TypeScript mixed value profile and native tuple structure in 0.16.0; ground Ontos stays unchanged |
| bittype | Source syntax, AST/CST projections, frontend tooling and generation |
| bitwire / bitruntime | Capability contracts and hydration above Horos |
| Contract libraries | Function, Stream, Cell and future behavioral declarations and realizations |

The semantic authority is the [horos descriptor specification](https://github.com/Bitspark/horos/blob/main/docs/spec/horos.md)
and the [canonical declaration specification](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md).
The latter specifies exact value grammar, immutable imports, indexed binders,
substitution, interfaces and declaration identity. Bittype must not maintain a
competing interpretation of those Horos data types or their generic substitution.

The separate [agnostic semantic model](semantics.md) implements shared facts:
identities, bindings, expressions, applications and dependencies. The
[analysis SPI](analysis.md) reports coverage and validity explicitly. Neither
supplies another canonical value-type algebra nor makes native meaning equivalent
to Horos. Horos supplies the first production analyzer.

The [unified declaration decision](unified-declarations.md) replaces separate
Horos and Horos/Wire languages with one declaration model. Its revision-3
[kind system](kinds.md) adds constructor interfaces and proper result claims.
Data is a refined kind of Type, backed by a compositional data-only certificate.
One `Page<T>` accepts any
admitted type; `DataPage<D: Data> = Page<D>` restricts arguments. Ground descriptors
are used directly, without a schema-level Data wrapper. The
[preservation map](wire-extension-preservation.md) distinguishes this type-language
change from unchanged Ontos value profiles and runtime responsibilities.

The [native representation proposal](type-representations.md) keeps optional
library-owned recipes selected by exact declaration identity. Native generation
retains open generics for the implemented pagination profile. Specialization
remains an independent operation;
no template lifting or copied module is required.

The Go model and one Bittype AST frontend implement revision 3. The
[dependency alignment](dependencies.md) qualifies public Horos/Ontos modules.
The first [native pagination transformers](semantic-generation.md) now consume this model. Runtime bindings remain unimplemented.

A declaration document is an ontos module value, encoded with ontos-codec-v1:

```text
Tuple(Atom(UTF8("horos/decl/module")),
      int(3), text(moduleName), importsMap, declarationsMap)
```

The maps have canonical text keys. Imports pin module bytes with SHA-256;
declarations contain ordered parameters with proper/constructor interfaces, an
explicit proper result claim and a body. The claim atom is `type` or `data`;
constructor interfaces have nonempty ordered inputs and a proper result claim.
Parameter occurrences store zero-based slot indices. References identify local
or imported declarations; applications contain an expression callee and a nonempty
ordered list of all arguments. Zero-parameter declarations use references.
See the [exact canonical grammar](kinds.md#exact-canonical-revision-3). Horos owns
the arities, name grammar, literal boundaries and resolution rules.

Existing horos descriptors are reused, allowing declaration expressions only
in child-descriptor positions. A parameter-shaped tuple inside an `exact`
literal remains data. `rec`/`ref` remain recursion syntax, not generic slots.

The [readable source grammar](horos-syntax.md) is implemented by bittype. The canonical
value format remains the authority for identity and semantics. Horos's value equations are
not JavaScript builders or executable input. Source spans, comments and native
naming preferences remain tooling metadata outside the canonical module value.

The old `bittype: 1` JSON proposal and its independent `string`, `bytes`, `i64`,
`u64`, `f64`, variant and nullable algebra are superseded. Source conveniences
need an explicit Horos meaning; they cannot create new canonical value domains.

## Pagination first

The complete value equations are in
[horos's pagination example](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md#5-complete-pagination-example).
Their readable meaning is:

```text
module example.pagination
  Page<T> = closed record {
    items: list-of(T)
    next?: scalar(utf8-text)
  }

  DataPage<D: Data> = Page<D>

module example.people
  import pagination = exact pinned example.pagination module
  Person = closed record { name: scalar(utf8-text) }
  PeoplePage = pagination.Page<Person>
```

This display is explanatory, not a frozen source grammar. `Page` is an ordinary
generic declaration. Its required items are an ontos list; optional next admits
text, not null. No pagination operation, cursor protocol or transport is implied.

Native generation preserves open `Page<T>` and applied `Page<Person>` in target
syntax. Separately, Horos specialization can produce a ground descriptor for
value checking. Users do not author that expanded description, and generation
does not require monomorphization as its only input.

Keep the imported Page declaration, the caller's Person binding, the named
PeoplePage alias and parameter origins. Same-spelled names in different modules
or declarations do not share a binding. The current declaration revision rejects
cyclic import/declaration-reference graphs, including ignored constructor
references; existing guarded `rec`/`ref` descriptions
remain available. Recursive named applications are deferred. Do not add a
productivity rule that rejects a well-formed but uninhabited Horos type.

Page's proper Type parameter accepts data or capabilities. DataPage's proper Data
parameter restricts its callers. `Page<Cell<Person>>` is valid for an ordinary
capability-bearing Cell declaration; `DataPage<Cell<Person>>` fails, including
when the capability is nested in an aggregate. `Forward<T> = DataPage<T>` fails
symbolically before specialization; `Forward<D: Data>` succeeds. Ignored type
arguments do not taint a data result, while data arguments do not imply a data
result: `wire<Person>` remains a capability type. Constructor parameters retain
those distinctions: known Page can satisfy `Data -> Data`, whereas an abstract
`F: Type -> Type` promises no Data result. The [kind examples](kinds.md#page-wrapped-and-compose)
show Wrapped/Compose and the corresponding rejected calls.

`Program.DataDescriptor` requires a data specialization before using WF/Conforms.
Ground `any`, complement, literals and recursion retain their meanings. The
complete constructor restrictions and bounded proof rules belong to Horos.

## Frontend and pipeline

```text
canonical Horos module values
    -> resolve/check using Horos declaration semantics
    -> AST[Horos]
         -> AST[TypeScript] -> CST[TypeScript] -> source
         -> AST[Go]         -> CST[Go]         -> source
         -> AST[Rust]       -> CST[Rust]       -> source
         -> AST[C++]        -> CST[C++]        -> source
```

The readable frontend adds text -> CST[Horos] and its AST view, using the same model. Native parsers continue to serve independent native inputs.
`AST[Horos]` is an ordinary language-specific use of the generic carrier. Its
provider and transformers are explicitly registered by the application; the
engine has no canonical preference or knowledge of descriptor constructors.
Transformers produce target ASTs, formatters construct CST, and the printer
alone emits text.

The [generation design](generation.md) and [layout](../LAYOUT.md) retain the
stage-first matrix. Planned canonical paths use `horos`, replacing the earlier
proposed `canonical` name. The AST, readable syntax and [native pagination mappings](semantic-generation.md)
are implemented. Broader native representations and reverse mappings remain future work.

## Native representation profiles

Horos classifies ontos values. Native types are selected representations of
those values; compilation does not prove conformance. Each transformer must
declare its supported subset, toolchain/configuration and representation limits.

| Horos case | Required mapping decision |
| --- | --- |
| `scalar(int)` | Arbitrary-width domain; no silent JS number or bounded machine integer |
| `scalar(utf8-text)` | Native string representation and its admissible-value limits |
| `list-of(T)` | Native sequence preserving the selected T representation and binding |
| Closed/open record | Field membership and native enforcement/extra-field limits |
| Optional field | Absence distinct from a present value |
| `exact(null())` and unions containing it | Explicit null distinct from absence |
| Unions, intersections, complements, exact values, refinements | Specified representation or unsupported-input diagnostic; no silent widening |
| Guarded recursion | Explicit native indirection/layout, without changing descriptor meaning |
| Named generics/applications | Open native parameters, declaration identity, interfaces and arguments |
| Proper Type/Data interfaces | Preserve source obligations; do not claim a native constraint proves Data without an explicit representation rule |
| Constructor parameters and residual higher kinds | Currently Unsupported in all four native profiles; no hidden specialization |

Pagination uses a generic interface in TypeScript, generic struct
in Go/Rust/C++, a sequence for items and presence representation for next. Exact
support types, imports and compiler settings accompany the [transformers](semantic-generation.md);
this document introduces no runtime or arbitrary-precision library dependency.

The example's next has two admitted states: absent or present text. An optional
union of text and `exact(null())` would have three. No new nullable constructor
or wrapper-normalization semantics are introduced. Current scalar tags are int,
utf8-text and bool; decimal has no scalar recognizer. A scalar extension requires
a concrete Horos model case and is unnecessary for the first pagination example.

## Dependencies, identity and forbidden shortcuts

Dependency direction is bittype's Horos tooling -> horos -> ontos. Horos must
not import bittype. Generic trees, contracts, engine and printer remain unaware
of concrete semantic dependencies. Native syntax and formatters stay independent
of Horos; its frontend and relevant transformers consume the semantic model.
Architecture tests allow the model dependency only in the Horos frontend and
classify its transitive packages. They retain the generic-package boundary and
the independence of native text stages.

There is one frontend and one explicit model dependency allowance. Kind/interface
meaning, Data certificates and constructor domains remain Horos-owned, not generic infrastructure.

Horos specifies document identity and pins. Bittype preserves them alongside
declaration/binder references. AST addresses and source spans are tooling
addresses/provenance, not module pins, ownership or authority proofs.

Function, Stream and Cell remain library-owned. This phase adds no live contract
semantics, adapters or binding registry. Native function/method-signature nodes
remain language syntax; they do not privilege a particular Function library.

Forbidden:

- Contract kinds, name switches, builtin imports or registrations in the general
  generator, carrier, engine, Wire, runtime or hydration.
- Structural fingerprints selecting behavioral ownership or bindings, including
  fingerprints hidden in a generated central dispatch table.
- Losing origins through expansion and guessing the owning library from shape.
- A competing data type algebra or native extraction as canonical authority.
- Adapters, RPC/streaming behavior, value validators or codecs in this types-only
  generation milestone.

## Implementation status and next scope

The owner authorized the first two stages on 2 November; they are implemented:

1. Qualify Horos declaration grammar, pins, resolution, slots and specialization
   with independently authored semantic vectors.
2. Add the AST-only Horos frontend in bittype, preserving origins, explicit
   binding references and located AST errors. Horos model errors retain their
   module/declaration/path context; no source location is invented.

The remaining stages are proposals, outside this change:

3. Implement four Horos-to-native ASTTransformer providers. First qualify open
   and applied pagination, then expand the supported subset explicitly.
4. Compose existing formatters/printer and compile generated declarations and
   use examples with pinned toolchains on shared workers or GitHub runners.
5. Later add partial native-to-Horos migration transformers. Missing native
   semantics cannot be reconstructed automatically.

Acceptance cases include imported pagination, nested applications, same-spelled
parameters/names in separate scopes, parameter-shaped literal data, absence/null,
bad imports, duplicate names, invalid arity/slots and recursive applications
outside the admitted revision. Exhaustion must not look like successful checking.
Generator tests inspect ASTs and origins; independent formatter fixtures remain
separate. Compiler success does not prove value conformance or transport behavior.
