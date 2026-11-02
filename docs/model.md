# Horos canonical types and bittype generation

**2 November 2026. Status: ownership direction accepted; declaration format and
native mappings proposed, not implemented.** This replaces the 31 October
proposal for an independent bittype value-type language. The generic AST/CST,
native text stages and syntax engine already exist; this update changes docs only.

## Ownership and canonical input

| Component | Responsibility |
| --- | --- |
| ontos | Values, structural identity, data embeddings and canonical encoding |
| horos | Data descriptors/conformance; proposed module, generic and application semantics |
| bittype | Source syntax, AST/CST projections, frontend tooling, language transformations and source generation |
| Contract libraries | Function, Stream, Cell and future behavioral declarations and realizations |

The semantic authority is the [horos descriptor specification](https://github.com/Bitspark/horos/blob/v0.6.0/docs/spec/horos.md)
and the new [canonical declaration proposal](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md).
The latter specifies exact value grammar, immutable imports, indexed binders,
substitution and identity. Bittype must not maintain a competing record,
collection, union or generic-substitution specification.

Horos 0.6.0 implements ground descriptors and value checking, without a public
dependency on the older proof-integration packages. It does not yet implement
the declaration layer. Bittype does not yet implement its frontend or any
Horos-to-native transformation. No package release or migration is implied.

A declaration document is an ontos module value, encoded with ontos-codec-v1:

```text
Tuple(Atom(UTF8("horos/decl/module")),
      int(1), text(moduleName), importsMap, declarationsMap)
```

The maps have canonical text keys. Imports pin module bytes with SHA-256;
declarations contain ordered parameter names and type expressions. Parameter
occurrences store zero-based slot indices. Applications contain an explicit
local/imported target and an ordered list of every type argument. Horos owns
the arities, name grammar, literal boundaries and resolution rules.

Existing horos descriptors are reused, allowing declaration expressions only
in child-descriptor positions. A parameter-shaped tuple inside an `exact`
literal remains data. `rec`/`ref` remain recursion syntax, not generic slots.

Readable source notation is a future bittype frontend concern. The exact value
format does not imply an implemented text parser. Horos's value equations are
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
or declarations do not share a binding. The first declaration revision rejects
cyclic import/application graphs; existing guarded `rec`/`ref` descriptions
remain available. Recursive named applications are deferred. Do not add a
productivity rule that rejects a well-formed but uninhabited Horos type.

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

A future readable frontend adds text -> CST[Horos] and its AST view, using the
same model. Native parsers continue to serve independent native inputs.
`AST[Horos]` is an ordinary language-specific use of the generic carrier. Its
provider and transformers are explicitly registered by the application; the
engine has no canonical preference or knowledge of descriptor constructors.
Transformers produce target ASTs, formatters construct CST, and the printer
alone emits text.

The [generation design](generation.md) and [layout](../LAYOUT.md) retain the
stage-first matrix. Planned canonical paths use `horos`, replacing the earlier
proposed `canonical` name. No directories, interfaces or dependencies are added
by this documentation change.

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
| Named generics/applications | Open native parameters, declaration identity and arguments |

For pagination, propose a generic interface/record in TypeScript, generic struct
in Go/Rust/C++, a sequence for items and presence representation for next. Exact
support types, imports and compiler settings must accompany the transformers;
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
The current no-external-dependencies check remains true for current code. Future
implementation must add a narrow dependency allowance and qualification, not
weaken the generic-package boundary globally.

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

## Next implementation scope and evidence

No implementation is authorized by this document alone. Proposed stages are:

1. Qualify Horos declaration grammar, pins, resolution, slots and specialization
   with independently authored semantic vectors.
2. Add the Horos AST/frontend in bittype, preserving origins and located errors.
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
