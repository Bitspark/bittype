# Horos canonical types and bittype generation

**Updated 4 November 2026. Status: Go declaration model and Horos AST implemented;
native mappings proposed, not implemented.** This replaces the 31 October
proposal for an independent bittype value-type language. The generic AST/CST,
native text stages and syntax engine already exist. See the [Horos frontend](horos.md).

## Ownership and canonical input

| Component | Responsibility |
| --- | --- |
| ontos | Ground data values, structural identity, data embeddings and canonical encoding |
| horos | Data descriptors/conformance; module, generic and application semantics for data declarations |
| optional ontos/wire | Go/TypeScript mixed value profile and native tuple structure in 0.16.0; ground Ontos stays unchanged |
| optional horos/wire | Implemented static Go mixed declarations and explicit closed ground embeddings |
| bittype | Source syntax, AST/CST projections, frontend tooling and generation |
| bitwire / bitruntime | Capability contracts and hydration above Horos |
| Contract libraries | Function, Stream, Cell and future behavioral declarations and realizations |

The semantic authority is the [horos descriptor specification](https://github.com/Bitspark/horos/blob/v0.8.0/docs/spec/horos.md)
and the [canonical declaration specification](https://github.com/Bitspark/horos/blob/v0.8.0/docs/spec/declarations.md).
The latter specifies exact value grammar, immutable imports, indexed binders,
substitution and identity for data declarations. Bittype must not maintain a
competing interpretation of those Horos data types or their generic substitution.

The separate [agnostic semantic model](semantics.md) implements shared facts:
identities, bindings, expressions, applications and dependencies. The
[analysis SPI](analysis.md) reports coverage and validity explicitly. Neither
supplies another canonical value-type algebra nor makes native meaning equivalent
to Horos. Horos supplies the first production analyzer.

The 3 November [optional-extension proposal](hydrated-types.md) puts the mixed
value profile in ontos/wire and its type/declaration model in horos/wire.
Ground Ontos/Horos remain data-only; hydration stays in existing Wire/runtime
layers. Bittype projects the optional Horos model instead of owning a second
contract language. Explicit Data embedding, Wire typing and ordinary library
Function/Stream/Cell compositions are retained. The
[preservation map](wire-extension-preservation.md) records the replacement of
both earlier unimplemented proposals without widening ground APIs.

The [native representation proposal](type-representations.md) specifies
per-language library recipes. Mixed-capable generics are authored directly in
horos/wire and support data and capability applications using the same identity.
Type generation preserves open generics. The consuming build selects recipes
by exact declaration identity. Template lifting and its transformer were removed
on 4 November 2026; native transformers and runtime bindings remain proposed.

Horos 0.8.0 implements ground descriptors and value checking and includes both
Go declaration profiles, without a public dependency on older proof-integration
packages. Bittype pins that release and Ontos 0.16.0 and implements their AST
projections/checkers. The [dependency alignment](dependencies.md) is a clean
upgrade with installed-consumer qualification. Neither package implements a
Horos-to-native transformation yet, and no service migration is implied.

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

Page's parameter here ranges over data types. A Wire cannot be supplied to that
ground declaration. When a generic must accept capabilities, author it directly
in horos/wire. That optional Page can accept both Data and Wire expressions
without a duplicate template or cross-profile transformation.

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
proposed `canonical` name. The AST-only stages are implemented; readable syntax
and directed transformations remain future work.

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
Architecture tests allow the model dependency only in the Horos frontend and
classify its transitive packages. They retain the generic-package boundary and
the independence of native text stages.

The horos/wire frontend has its own narrow dependency allowance to consume the
optional model and ground embeddings. That is not an allowance for generic
infrastructure. Horos owns both profiles' distinct semantics; Bittype projects
them without redefining their meaning.

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
