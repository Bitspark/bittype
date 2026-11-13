# Semantic native type generation

**Declarative generation SPI, 8 November 2026.** The four native
pagination mappings consume general kinds through explicit target
profiles, using public Horos Go 0.10.2. The [automatic composer](generation-composition.md)
now drives all four, with Horos mappings in Bittype's root Go module. This
is types/interfaces generation; it adds no runtime adapters, application-value
validators, serialization or executable behavior.

## Explicit semantic input

Capture admitted Horos syntax once, resolve it explicitly, and analyze that exact
resolution result. Supply the same snapshot and explicit analysis evidence to
the generation request. Equal bytes or matching tree addresses from another
capture do not substitute for that association. Configuration, dependency context,
operation/profile, required obligations and available typed Semantics are part of
the input contract.

The existing `Semantics[horos.Language]` graph and sealed representation envelope
remain in use. The bridge retains multi-Program captures and tree-qualified
sites. Files organization does not replace canonical module imports. Analysis
consumes resolution evidence; no mapping invokes a hidden resolver or analyzer.

```text
Horos AST collection + explicit resolution/analysis evidence
  -> PlanGeneration: closed native types and bound references
  -> Run -> generation.Output: checked target syntax, exposures and external uses
  -> TypeScript Files placement and native analysis -> Source -> files
```

`ASTTransformer` retains FromID/ToID and returns `representation.ASTSet`.
Syntax routes use that contract. Generation uses `Engine.Generate` and its sealed output. The typed bridge carries a language-specific request and output set; it does not
convert syntax to `AST[any]` or semantics to `Semantics[any]`. Each engine route
leg admits its output and creates a new capture. Evidence from an earlier leg is
not automatically valid for changed syntax. A later semantic leg needs explicit
resolution/analysis of its own snapshot.

## Representation and enforcement profiles

The [horos mapping package](../transform/syntax/horos/mapping/go/mapping.go) declares the two profile identities, RepresentationProfile and EnforcementProfile. One test beside the declaration pins both values. Every other test and fixture takes them from there. Each native mapping is its own package under `transform/syntax/horos-to-<native>/mapping/<name>/go` and exports `Library()`. No package lists the mappings of a pair. Each consumer registers a source, target and mapping libraries and passes the root declaration sites, output tree ID, receipts, profile and limits to Engine.Generate (or PlanGeneration then Run). Go's target provider takes Package and C++ takes Namespace. The [TypeScript example](kernel-composition.md#horos-to-typescript) shows registration and execution.

| Profile | Promise and behavior |
| --- | --- |
| RepresentationProfile, `structural-representation/1` | Native declarations represent uses already admitted by Horos. The profile expressly permits erasing ordinary proper-kind contracts from target parameter declarations. Every builtin mapping states this profile. |
| EnforcementProfile, `structural-enforcement/1` | Native clients would have to be prevented from violating the source constraint. No such target encoding is implemented, and no mapping in bittype states this profile. A request for it fails with generation.ErrUnsupported. |

The two identities state what the profile promises and name no fixture. The
version stays `/1`: nothing is released and the generated output did not change
with the names.

A profile is a promise about the native type. The kernel does not interpret
it. Three rules decide what a request receives:

- **A mapping states its own profile.** The profile is a field of the
  description that a mapping library returns. Each builtin mapping writes
  `RepresentationProfile` there, and its constructor takes no profile. A
  composition chooses which libraries it registers. It cannot choose what a
  builtin library promises.
- **A request selects by its profile.** The registry plans with the mappings
  that state the requested profile and have compatible languages. With none it
  fails with generation.ErrUnsupported. This is how the enforcement profile is
  refused, and any profile that no registered mapping states.
- **The registry does not judge the promise.** Selecting a contribution does
  not establish its claim. A contribution of another author that states a
  profile is selected under it, and that author answers for it.

No shared code lists, writes or refuses a profile. The horos mapping package
holds the two identities and the function type `Library`, which makes a
description function a mapping library and returns what the function returns. A
nil function fails with generation.ErrInvalid. A description that states no
profile is refused at registration by the generation contract, also with
generation.ErrInvalid.

A further profile, such as the bounded profile that the
[families design](representation-families.md#domains) names, enters as a
contribution that states its own identity in its own package, registered by the
consumer's composition. It needs no change to the horos mapping package or to
the registry. The tests in
[`horos_generation_profile_test.go`](../integration/go/horos_generation_profile_test.go)
show this with the fixture contribution
[`readonlyprofile`](../integration/go/testdata/readonlyprofile/library.go),
which maps a list to another native form under its own profile: one registry
holds both families, and each request receives its own. The same file shows
that every builtin library states the structural profile, and it covers the
two refusals. [bittype/ACX-008](architecture-cases/0008-structural-profile-gate.md)
records the earlier closed gate, its correction and its limits.

The profile is part of what a mapping describes, not inferred from a kind
name. Representation does not advertise arbitrary native instantiations as
Horos-admitted values. Native strings, records and generic parameters do not
prove every source invariant. A request for the enforcement profile is not
answered by a mapping that states the representation profile.
[GEN-02](generation-gaps.md#gen-02) records the missing native-plan
capability, target encodings and completion criteria.

The [native kind-constraint proposal](native-kind-constraints.md) describes a
future enforcing route using direct bounds or source type witnesses, with
language-specific limits. It is documentation only: the representation profile
still permits erasure, and no mapping states the enforcement profile.

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
not become native type exports. The caller requests explicit declaration sites. The pagination example requests the entry module’s type definitions. Imported types are demanded through selected mappings; unused
imports and bodies represented by an exact library mapping need no native
expansion. Source analysis still admits the complete original source.

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

Go's record mapping wraps an optional field in a helper, and Go's optional
mapping supplies that helper with the `BittypeOptional` declaration. Neither
package creates the helper's key and neither imports the other. The consumer
creates the key once and passes it to both:

```go
key := g.NewHelperKey("horos/optional")
libraries := []g.MappingLibrary[h.Language, s.Language, source.Reader]{
    record.Library(key),
    optional.Library(key),
}
```

Here g is `contracts/generation/go`, h and s are the horos and Go syntax
packages, and record and optional are the packages under
`transform/syntax/horos-to-golang/mapping/<name>/go`. An invalid key makes
registration of the record library fail with `ErrInvalid`. A valid key that no
registered library supplies makes planning fail with `ErrUnsupported`.
[`TestGolangRecordRequiresItsOptionalHelper`](../integration/go/horos_generation_golang_optional_test.go)
covers both. This form is interim:
[#196](https://github.com/Bitspark/bittype/issues/196), item 6, changes how a
record asks for the optional representation. Whether a record asks at all is
conditional since 13 November 2026: see
[optional fields](contribution-contract.md#optional-fields) in the
contribution contract. Decided later that day, after the prototype: a record
asks, because presence is delegated, unless an objection of the service
library's owner reopens the decision before the migration merges.

Canonical integers remain Unsupported. The existence of native integer syntax
does not choose an arbitrary-width representation or dependency. Open records,
Wire, tuples, sets/maps, unions/intersections, complements, exact/refinement forms
and ground recursion remain outside this profile. Broader representations need
an explicit target contract rather than a weaker fallback.

## Generic and naming limits

The default structural policy rejects residual constructor parameters,
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

Each mapping currently produces one declaration tree containing the checked type
closure. Declarations belonging to the entry module keep their native names (subject to
the naming rules above). Declarations from imported modules
use `M<first-16-hex-of-module-pin>_<Name>`, with collisions rejected. This
product spelling is not semantic identity. Here entry-module membership is
distinct from a consumer explicitly requesting a declaration as a public root. The configured output tree ID and Go package
or C++ namespace remain separate from canonical module identity. This mapping
performs no implicit filesystem or network IO.

**Naming design, 13 November 2026. Accepted, not implemented.** The two
paragraphs above describe what generation writes today. The
[naming design](naming.md) changes who decides what about a name. The mapping
spells it under a policy that the consumer's composition selects. Generation
records that as a naming intent: the chosen name, the module's identity, the
policy's identity and the spelling. Each language checks the spelling under
its own rules, and publishes it or refuses. No target has changed a name.
These statements change when [#274](https://github.com/Bitspark/bittype/issues/274) and the first Files slices of Go
and Rust ([#161](https://github.com/Bitspark/bittype/issues/161), [#159](https://github.com/Bitspark/bittype/issues/159)) deliver it:

- **The module-pin prefix is today's default, and nothing can select another
  rule.** It becomes one of two policies that the consumer's composition
  selects: qualified, which is today's rule, and unqualified. The policy
  becomes selectable first, with unchanged output. Unqualified becomes a
  target's default once that target's Files stage owns collisions. The
  reasons: a content change renames every prefixed public name, a truncated
  pin is no uniqueness guarantee, and the name of a declaration depends on
  which module the request enters through.
  [The qualification policy](naming.md#the-qualification-policy).
- **`Self` becomes `Self_` in all four targets, and only Rust reserves it.** A
  reserved word is the target's. Rust will refuse `Self`, and the other
  targets will keep it. The refusal of a name that contains `__`, which the
  shared rule makes for every target today, moves to the C++ target.
  [Reserved words](naming.md#reserved-words).
- **Every declaration name is capitalized, so every generated Go declaration
  is exported.** The capital stays the mapping's style for every target and
  decides no access. Go gets its own check: Go Files will refuse a
  declaration that requires public access and whose spelling is not exported,
  and an output with an unexposed dependency. It renames nothing.
  [Go](naming.md#go) ·
  [bittype/ACX-004](architecture-cases/0004-generation-file-syntax.md).
- **Name collisions and parameter shadowing fail for the whole output.** Both
  are the flat-name rule. It leaves a target when that target's Files stage
  checks collisions in a placed scope, and a type parameter that captures a
  use is then refused for that placement
  ([delivery order](naming.md#delivery-order-per-target)). The sentence below
  that "Go reserves `BittypeOptional`" is today an effect of the same rule.

The sentence above that the version stays `/1` because the generated output
did not change is about the renaming of the profile identities. A change of a
target's default policy changes names and no promise of a profile. The design
says when each lands
([What is fixed, and when](naming.md#what-is-fixed-and-when)).

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
configuration. The planner admits capture/evidence association before mapping.
The kernel walks the published graph itself
([GEN-01](generation-gaps.md#gen-01)). The horos interpretation and reader
decode published facts and terms; they do not rerun horos
resolution, query a different Program or restore a Data-summary fallback.

The bundled prelude can have authoritative source-less entities when its exact
standard module identity is absent from the supplied syntax. No native source
origin is invented for those entities. Supplied prelude syntax corresponds
normally. The prelude's kind definitions remain compile-time declarations.

## Syntax and Files

Generation outputs retain syntax-checked declaration sites and intended bindings. TypeScript Files applies placement and establishes target semantics before publishing Source. Syntax formatting and printing can also operate directly on individual output trees. [Generation handoff](typescript-type-imports.md#from-generation-to-files)

Mappings retain real declaration/parameter links and available current source
origins. Imported definitions are emitted once with dependencies before uses.
The horos interpretation and reader live at `transform/syntax/horos/source/go`; bounded native
storage lives at `transform/syntax/shared/target/go`, with native adapters under the
four language paths. The packages under `transform/syntax/horos-to-<native>/mapping/<name>/go` choose representations, one mapping per package. Shared contracts, engine and
printer gain no Horos/Ontos or behavioral-library imports.

Generation uses the public SPI defaults: 100,000 operations, depth 128, 4,096
plans/reservations/formals, 65,536 references, 200,000 output nodes and 32 MiB
output data, together with context cancellation. Formatter, carrier
and semantic-input limits still apply. Exhaustion returns no usable output.

## Library and runtime boundaries

Function, Stream, Cell and future contracts remain ordinary library definitions.
Their kind membership or expanded shape does not select a native binding.
Optional recipes belong to their owning libraries and consuming configuration,
selected by exact declaration identity and target profile. The
[type-owned generator design](type-generation-plugins.md) and
[composition API](generation-composition.md) are implemented. The imported Cell
fixture demonstrates extending the profile without changing Bittype or Wire.
Complete production contract libraries and their behavior remain separate work.

This migration preserves types-only output. It does not create serializers,
application-value validators, executable adapters, protocol implementations or a
replacement hydration layer. Ground descriptor extraction remains a Horos
operation, separate from native representation.
