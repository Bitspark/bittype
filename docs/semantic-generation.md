# Semantic native type generation

**Revision 4 implemented in Go source, 4 November 2026.** The four native
pagination mappings consume general kinds through explicit target
profiles, using public Horos Go 0.10.1. This
is types/interfaces generation; it adds no runtime adapters, application-value
validators, serialization or executable behavior.

## Explicit semantic input

Capture admitted Horos syntax once, resolve it explicitly, and analyze that exact
resolution result. Supply the same snapshot and explicit analysis evidence to
selected transformers. Equal bytes or matching tree addresses from another
capture do not substitute for that association. Configuration, dependency context,
operation/profile, required obligations and available typed Semantics are part of
the input contract.

The existing `Semantics[horos.Language]` graph and sealed representation envelope
remain in use. The bridge retains multi-Program captures and tree-qualified
sites. Files organization does not replace canonical module imports. Analysis
consumes resolution evidence; no mapping invokes a hidden resolver or analyzer.

```text
Horos AST collection + explicit resolution/analysis evidence
  -> target ASTSet
  -> target Syntax formatter -> CST -> Syntax printer -> native source
```

`ASTTransformer` retains FromID/ToID and returns `representation.ASTSet`.
The typed bridge carries a language-specific request and output set; it does not
convert syntax to `AST[any]` or semantics to `Semantics[any]`. Each engine route
leg admits its output and creates a new capture. Evidence from an earlier leg is
not automatically valid for changed syntax. A later semantic leg needs explicit
resolution/analysis of its own snapshot.

## Representation and enforcement profiles

Each of `transforms/horos-to-typescript/go`, `transforms/horos-to-go/go`,
`transforms/horos-to-rust/go` and `transforms/horos-to-cpp/go` exports:

```go
const RepresentationProfile = "pagination-representation/1"
const EnforcementProfile = "pagination-enforcement/1"
```

Their `Config.Profile` is required. `New(from, to, Config)` rejects an empty or
unknown profile with `transform.ErrInvalidTransformer`. `TreeID` remains explicit;
Go also takes its package configuration and C++ its namespace configuration.
The language/profile/marker of both endpoints must match the concrete mapping.
Applications register languages, checkers and transformer edges themselves.

| Selected profile | Promise and behavior |
| --- | --- |
| RepresentationProfile | Native declarations represent uses already admitted by Horos. The profile expressly permits erasing ordinary proper-kind contracts from target parameter declarations. |
| EnforcementProfile | Native clients would have to be prevented from violating the source constraint. No such target encoding is implemented; Transform returns `text.Unsupported` at the request entry and no output. |

The profile selection is part of the configured mapping, not inferred from a kind
name. Representation does not advertise arbitrary native instantiations as
Horos-admitted values. Native strings, records and generic parameters do not
prove every source invariant. An enforcement request never silently selects
representation instead.

The [native kind-constraint proposal](native-kind-constraints.md) describes a
future enforcing route using direct bounds or source type witnesses, with
language-specific limits. It is documentation only: the representation profile
still permits erasure and the enforcement profile still returns Unsupported.

## Pagination subset

```text
module example.pagination;
import h from horos.kinds pin "sha256:<actual-64-hex-pin>";
type Page<T> = { items: list<T>; next?: text; };
type Person = { name: text; active: bool; };
type PeoplePage = Page<Person>;
type NestedPage = Page<Page<Person>>;
type DataPage<D: h.Data>: h.Data = Page<D>;
```

The placeholder denotes the actual ordinary prelude module identity, not literal
source. Omitted type parameter/result kinds elaborate to `top(P)`. Horos checks
all source contracts through its general kind model. The generator preserves one
open Page and ordinary named applications; it does not materialize a separate
Page body for each caller.

The reader consumes proper-sort facts, checked bodies, named references,
applications and bindings. It does not inspect the names Type/Data or decode a
closed claim enum. KindDefinition entities are compile-time declarations and do
not become native type exports. All exported TypeDefinition entities in the
selected semantic closure remain in scope, including dependency modules. An
unsupported type definition is not pruned because a particular root ignores it.

| Horos form | TypeScript | Go | Rust | C++ |
| --- | --- | --- | --- | --- |
| UTF-8 text | string | string | String | std::string |
| Boolean | boolean | bool | bool | bool |
| list-of(T) | T[] | []T | Vec<T> | std::vector<T> |
| Named closed record | interface | struct | struct | struct |
| Optional field | field?: T | BittypeOptional[T] | Option<T> | std::optional<T> |
| Proper type parameter | native parameter | native parameter | native parameter | template parameter |
| Ordinary named application | F<T> | F[T] | F<T> | F<T> |

TypeScript uses strictNullChecks and exactOptionalPropertyTypes. Go's generated
`BittypeOptional` declaration has Present and Value fields; Value is ignored when
absent. It has no generated methods or runtime implementation. Optional presence
is not silently replaced by null or a zero value.

Canonical integers remain Unsupported. The existence of native integer syntax
does not choose an arbitrary-width representation or dependency. Open records,
Wire, tuples, sets/maps, unions/intersections, complements, exact/refinement forms
and ground recursion remain outside this profile. Broader representations need
an explicit target contract rather than a weaker fallback.

## Generic and naming limits

The complete exported type closure rejects residual constructor parameters,
constructor-valued results or aliases, higher-order arguments and computed
callees. It also rejects unused proper parameters. General Horos support for
these expressions does not establish a native encoding. A closed application
that could be reduced does not authorize hidden specialization.

Records must be named declarations. Ordinary nested applications such as
`Page<Page<Person>>` remain supported; anonymous record bodies inside fields or
arguments are Unsupported. Unsupported cases return source-node diagnostics and
no usable output.

Declarations and proper parameters capitalize their initial ASCII letter;
a leading underscore becomes `Type_`, and `Self` becomes `Self_`. Name collisions,
parameter shadowing of a declaration and unsupported field spellings fail.
Go fields also capitalize their initial letter; TypeScript retains property
names. Rust and C++ reject keyword field names. Names are not guessed from types.

Each mapping currently emits a one-tree ASTSet containing the checked type
closure. Root declarations keep their exported native names. Imported declarations
use `M<first-16-hex-of-SHA256(pin-text)>_<Name>`, with collisions rejected. This
product spelling is not semantic identity. The configured output tree ID and Go package
or C++ namespace remain separate from canonical module identity. This mapping
performs no implicit filesystem or network IO.

Go reserves `BittypeOptional` for its generated support declaration. C++ references
use the configured qualified namespace; that namespace cannot be `std` or another
reserved identifier. A field cannot shadow a C++ template parameter. Standard
library types are explicitly qualified in Rust and C++.

## Facts and required evidence

Horos-owned subject facts distinguish semantic sort, declared kind contract,
result sort/kind, checked symbolic body and individual established memberships.
The relation keys are provider vocabulary, not interpretation rules in the
shared graph. Multiple membership facts can concern one subject, and missing
facts assert nothing. The retained body is not a principal inferred kind.

Source membership records retain their actual parameter assumptions. A narrower
query such as Page preserving a scalar-derived category cannot create an
unconditional result fact under Page's original unrestricted parameter.
Arbitrary later Horos queries do not mutate the captured source evidence.

The mapping requires the current Horos analysis profile, its validity and
required enumeration reports, an available typed semantic graph, and compatible
configuration. The adapter admits capture/evidence association before mapping.
The source-specific reader decodes published terms; it does not rerun Horos
resolution, query a different Program or restore a Data-summary fallback.

The bundled prelude can have authoritative source-less entities when its exact
standard module identity is absent from the supplied syntax. No native source
origin is invented for those entities. Supplied prelude syntax corresponds
normally. The prelude's kind definitions remain compile-time declarations.

## Syntax and Files

Target ASTSet is syntax, not Source. The implemented [Syntax and Files
composition](source.md) keeps typed Semantics, file organization and exact
receipts separate from syntax output. Native pagination mappings do not establish
target-language Semantics. Native Files formatting therefore requires a native
semantic provider; target Syntax formatting and printing remain usable for every
output tree.

Mappings retain real declaration/parameter links and available current source
origins. Imported definitions are emitted once with dependencies before uses.
The source reader and bounded target arenas live in `internal/horosgen`; each
concrete mapping owns target AST construction. Shared contracts, engine and
printer gain no Horos/Ontos or behavioral-library imports.

Generation keeps the existing limits of 100,000 target nodes, traversal depth
128 and 800,000 visits, together with context cancellation. Formatter, carrier
and semantic-input limits still apply. Exhaustion returns no usable output.

## Library and runtime boundaries

Function, Stream, Cell and future contracts remain ordinary library definitions.
Their kind membership or expanded shape does not select a native binding.
Optional recipes belong to their owning libraries and consuming configuration,
selected by exact declaration identity and target profile. Recipe SPIs and richer
capability representations remain separate work.

This migration preserves types-only output. It does not create serializers,
application-value validators, executable adapters, protocol implementations or a
replacement hydration layer. Ground descriptor extraction remains a Horos
operation, separate from native representation.
