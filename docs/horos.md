# Horos declaration AST

**Implemented in Go source, 4 November 2026.** One frontend projects
Horos declaration revision 3. [Kinds and interfaces](kinds.md) admit proper and
constructor parameters within that language; data and capability expressions share one module, declaration and
binder model. The old Horos/Wire packages and profiles are removed without aliases.
The authoritative [canonical format](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md)
lives in Horos. The [readable source grammar](horos-syntax.md) is implemented on this
same AST. Native generators and runtime bindings remain outside this frontend.

The [dependency alignment](dependencies.md) and independent installed consumer
qualify the public module graph. Private projection and binding mechanics live
in `internal/horosmodel`; syntax framing lives in `internal/declsyntax`.
The agnostic kernel contains no Type/Data or Wire interpretation.

## API and ownership

- [syntax/horos/go](../syntax/horos/go/nodes.go) defines `Language`, `AST`, `Node`,
  `Kind` and the vocabulary. It imports only generic AST/CST and language identity.
- [resolve/horos/go](../resolve/horos/go/resolve.go) provides
  `Resolve(root, dependencies, treeID)`, `Project(program, treeID)` and `Model(view)`.
  Resolve delegates semantics to Horos, then projects its checked model. Model
  checks the borrowed AST and returns a Horos program for inspection or explicit
  specialization with typed proper arguments and same-Program constructor handles.
  No lookup performs network/file IO.
  The registered resolver uses [ResolveReferences](../resolve/horos/go/resolver.go)
  to export shared semantic facts. It verifies supplied semantic annotations
  against authoritative resolution and never repairs them implicitly.
- [check/horos/go](../check/horos/go/check.go) admits AST syntax: kinds, roles,
  cardinalities and lexical spellings. Unresolved or semantically invalid names
  can reach resolution. Check does not resolve bindings or validate module pins.
- [analyze/horos/go](../analyze/horos/go/analyze.go) consumes the explicit resolver
  result, reusing its model and adapting its evidence to
  [shared semantic facts and reports](analysis.md).
- [parse/horos/go](../parse/horos/go/parse.go) constructs a lossless CST and stored AST
  from source; [format/horos/go](../format/horos/go/format.go) constructs CST directly
  from a borrowed AST. Printing delegates to the [generic printer](../print/go/print.go).
- [provider/horos/go](../provider/horos/go/provider.go) implements
  `parse.Provider[horos.Language]`, `format.Provider[horos.Language]`,
  `print.Provider[horos.Language]`, `check.Provider[horos.Language]`,
  `resolution.Provider[horos.Language]`
  and `analysis.Provider[horos.Language]`,
  with ID `horos` and name `Horos`.
  Applications explicitly register it through `representation.NewLanguage` and
  `Engine.RegisterLanguage`, using profile `horos-declarations-3/ast-1`.
  Syntax checking is separately bound with `check.NewChecker[h.Language](definition,
  provider)` and installed with `Engine.RegisterChecker`. Identity registration
  alone supplies no syntax or semantic capability.
  Resolution requires `resolution.NewResolver` and `Engine.RegisterResolver`.
  Analysis additionally requires `analysis.NewAnalyzer` and
  `Engine.RegisterAnalyzer`, and consumes the explicit resolution result.
  See [shared composition](resolution.md); implementing interfaces installs nothing.

Only Horos frontend stages and their private bridge directly import Horos/Ontos.
Text stages use Ontos for literal embeddings and the bridge for canonical module
encoding/pins; they do not resolve or specialize types.
The analyzer consumes the resolver contract and shared evidence.
The model dependency is pinned to a public immutable Go module revision in go.mod;
there is no local replacement or vendored private checkout. Generic infrastructure
and the four native text stages remain independent of the external model. The engine
can hold this AST without registering any transformer. Parse and Render use the
explicitly registered text stages; identity registration still installs nothing.

The syntax profile revision is a clean change: no old-profile alias or fallback checker
remains. `resolve.Model` is still explicitly semantic, while `Pack` and `Check`
establish only the new syntax profile. Existing callers that require resolution
must use Model or the analysis capability.

## Vocabulary and bindings

The syntax package's Kind is an AST node tag. Semantic kind shapes and interfaces
are Horos judgments projected separately; an AST tag does not establish a kind.

One Program tree supplies the entry module and dependency pool. Projection of a
checked Program emits its complete resolved import closure, ordered by module
digest; imports/declarations retain Horos's canonical map order. Parameters,
arguments and descriptor operands retain their specified order. Resolution selects
the Program root as its explicit entry; its required root reference selects the
entry module. Other entry kinds or multiple trees are unsupported request forms
in this first resolution profile.

| Kind | Attributes | Children | References |
| --- | --- | --- | --- |
| Program | none | `module*` | `root` to the entry module |
| Module | `name`, `pin` (lowercase hex) | `import*`, `declaration*` | none |
| Import | `alias`, expected `name`, `pin` | none | optional `target` module assertion |
| Declaration | `name` | `parameter*`, one proper `result` interface, one `body` | none |
| TypeParameter | `name`, canonical decimal `index` | one `interface` | none |
| Interface | `name` is `type`, `data` or `constructor` | none for proper interfaces; nonempty `input*` and one proper `result` for constructors | none |
| Parameter | canonical decimal `index` | none | optional `target` TypeParameter assertion |
| Reference | `scope` (`local` or `imported`), `name`, imported `alias` | none | optional `target` Declaration and imported `import` assertions |
| Application | none | one `callee`, nonempty `argument*` | none |
| Descriptor | Horos `constructor` name without `horos/` prefix | `operand*` | none |
| Atom | raw bytes as lowercase `hex` | none | none |
| Tuple | none | `element*` | none |

Descriptor operands retain the canonical value layout. Only Horos-designated
type positions become Parameter, Reference, Application or Descriptor nodes. Other positions
remain Atom/Tuple fragments, including field wrappers and everything inside an
`exact` literal. The frontend takes these positions from the checked Horos model;
it has no duplicate descriptor-constructor switch or record semantics.

Supplied binding annotations in the analyzed closure must identify the correct
module, import, declaration or parameter. No attributes, kinds or roles outside
this syntax profile are admitted. Optional annotations can be absent; missing
annotations do not trigger AST mutation. Arena addresses may change if every edge
changes with them; attribute/reference storage order is immaterial. Containment
order remains canonical in checked projections. Origins are retained as tooling
metadata but do not affect module pins.

Explicit reference resolution reconstructs canonical module values, asks Horos to
resolve them, and compares the authoritative syntax and supplied binding edges
with the supplied AST closure. Analysis consumes this established evidence.
This checks one declared AST profile; it never fingerprints a shape to choose a
behavioral library. Invalid UTF-8 attributes are rejected before text embedding.
It does not reinterpret literal data as declarations or silently normalize an AST.

## Pagination and specialization

The [binding tests](../resolve/horos/go/resolve_test.go) preserve imported Page
and caller-side Person identities, same-spelled distinct binders and literal
opacity. Data-only fixtures spell their proper Data interfaces explicitly.
The frontend retains these ordinary applications:

```text
Page<T> = list-of(T)
DataPage<D: Data> = Page<D>
People = Page<Person>
Live = Page<wire<Person>>
```

The [installed consumer](../integration/go/testdata/installed/consumer_test.go)
uses the full closed record with required items and optional text next. Both
applications retain the same open Page. Horos Specialize produces a closed type
expression; DataDescriptor additionally proves data and returns the existing
Horos descriptor. A mixed result cannot be passed to ground WF/Conforms.
No specialization mutates the original generic.

`Forward<D: Data> = DataPage<D>` resolves; `Forward<T> = DataPage<T>` is a located
semantic interface violation even before specialization. Both are valid AST syntax.
Check validates interface structure and spelling only; Horos supplies the meaning.
Constructor interfaces support `Wrapped<F: Type -> Type, T>` and
`Compose<F: Type -> Type, G: Type -> Type, T>` without losing known Page/Ignore
result precision. See the [complete judgment table](kinds.md#page-wrapped-and-compose).
The resolver exports structured provider-owned facts, and analysis consumes the
same immutable receipt. The required type-facts report records enumeration
completeness without adding kind semantics to the kernel.

## Admission and qualification

Horos owns module-name/pin checking, arity, binding scope, recursion closure,
declaration-reference cycle rejection, symbolic interfaces/Data proofs and specialization limits. Its errors retain module,
declaration and tuple-index path context. Structured codes and subjects distinguish
invalid source, missing dependencies, contradicted context, invalid options and
exhaustion. AST structure/profile errors are located at current AST nodes; no
source location is invented or copied from historical origins.

Projection checks generic AST node, depth, entry and text budgets while creating
nodes, including hexadecimal growth before copying atoms. Generic structure is
checked before the AST is decoded. Horos's default declaration budgets apply
when recovering a model through the canonical-value helper; registered resolution
can accept an explicitly supplied bounded configuration. Both boundaries return empty output on failure; a model
accepted under custom larger limits is re-admitted under the fixed Horos defaults
before projection, so successful projection respects `resolve.Model`'s fixed
canonical budgets. The syntax checker has separate AST budgets; resolution can
select explicit resolver limits. Canonical and AST depth are distinct, and tests
qualify their boundary.

## Analysis profile and generic queries

`analyze.Contract()` advertises `horos-declarations-3/analysis-1`. Its domain is
the selected entry and transitive pinned module closure, not every unused value
in the supplied pool. Unused pool modules contribute no established semantic
facts; their binding annotations are outside this analyzed domain. Every supplied
module still has its document-pin assertion checked during decoding.

| Obligation | Successful report |
| --- | --- |
| program-validity | Complete / Holds for the submitted entry and required context |
| declaration-index | Complete / Undetermined; module, declaration and parameter entities enumerated |
| occurrence-bindings | Complete / Undetermined; parameter/declaration references enumerated |
| applications | Complete / Undetermined; references, effective applications and provider expressions exported |
| reference-dependencies | Complete / Undetermined; every named declaration reference, including unapplied and ignored constructor arguments, exported |
| type-facts | Complete / Undetermined; kinds, interfaces, Data summaries/statuses and canonical identities exported |

The resolver additionally reports reference-resolution, requiring Complete /
Holds for its model. Its profile is horos-declarations-3/resolution-1. Analysis
requires that resolver evidence and preserves these six reports. All six
reports are analysis model publication requirements. A semantic rejection or
missing dependency publishes diagnostics without a partial graph. Missing context
and resolver exhaustion are Incomplete / Undetermined. A contradictory context
can reject the submitted program/context obligation without asserting that the
entry source itself is malformed; diagnostics identify the subject. Export
exhaustion after successful resolution preserves Complete / Holds reference
resolution and validity but
leaves graph enumeration incomplete and publishes no model. Unknown requested
obligations receive Unsupported / Undetermined.

`resolve.LimitsSetting()` names provider-owned JSON configuration containing
`[MaxDepth, MaxNodes, MaxBytes, MaxModules]`. All four limits are required; absence
uses Horos defaults. This profile requires positive limits with depth at most 512.
Invalid configuration is an API error before source decoding. The kernel copies and
bounds these bytes without interpreting them. Semantic graph budgets are separate.

The graph preserves module/declaration/parameter identity, current input sites,
open proper/constructor parameters and declaration references. Descriptor terms contain exact
Ontos codec-v1 bytes and authoritative operand paths. Declared interfaces and
result claims replace the former bound-only representation. Application terms have Callee
and Arguments; their callee Reference carries the declaration Target. Parameter
callees remain References to the actual parameter entity, not invented
declarations with synthetic formals.

`Model.Facts(subject, relation)` retrieves provider-owned terms associated with
an entity or term. Horos relations include `kind`, `interface`, `data-status`,
`data-summary` and `canonical-identity`, all owned by `horos`. Kind and interface
terms retain input operands as explicit graph links; Data summaries retain
conditional parameter/application structure. Canonical identities retain exact
module/declaration/slot identity separately from model-local IDs. Parameter Sort
is the coarse key `horos/type` or `horos/constructor`; its meaning remains
provider-owned. No kernel subkinding rule or Type/Data
interpretation is added.

The resolver exposes these keys through `KindRelation`, `InterfaceRelation`,
`DataStatusRelation`, `DataSummaryRelation` and `IdentityRelation`.
All operators below have Owner `horos` and Encoding `ontos/codec-v1`.
Payload notation uses raw atoms, the standard text/integer embeddings and tuples.

| Relation | Operator Name | Payload and operand links |
| --- | --- | --- |
| kind | `kind/proper`, `kind/constructor` | Empty atom; constructor inputs at `[i]`, with a proper result implicit |
| interface | `interface/proper`, `interface/constructor` | Atom `type` or `data`; constructor inputs at `[i]` |
| data-status | `data-status/proven`, `data-status/failed`, `data-status/conditional` | Atom containing the same status; no operands |
| data-summary | `data-summary/true`, `data-summary/false`, `data-summary/need`, `data-summary/and`, `data-summary/probe` | Empty atom; Need links a parameter Reference at `[0]`, And child summaries at `[i]`, Probe retained checked expression syntax at `[0]` |
| canonical-identity | `identity/module` | Raw digest atom |
| canonical-identity | `identity/declaration` | Tuple of raw module digest atom and declaration-name text |
| canonical-identity | `identity/parameter` | Tuple of raw module digest atom, declaration-name text and zero-based slot integer |

Module entities receive identity facts. Declaration entities receive kind,
interface, identity and body Data-summary/status facts; parameter entities receive
kind, interface and identity. Proper expression terms receive kind, interface and
Data-summary/status facts; constructor expressions have no proper-result Data
facts. Probe export does not recursively expand that expression's summary or
invent an AST association. The
`declaration-reference` dependency relation is exposed by `ReferenceDependency`.

Generic queries need no Horos model import: `Resolve`, `TermAt`,
`DirectParameterBindings`, `DependenciesFrom` and `Facts`.
DirectParameterBindings succeeds only for a direct Reference callee with a
matching actual formal list; it is not a general substitution operation.
Dependency edges cover all declaration references, not just applications, and
do not describe runtime effects.

Remote tests cover canonical value round trips, imported and nested pagination,
specialization, wrong/missing references, invalid pins/indices/UTF-8, literal
boundaries, address renumbering, provenance, nil views and resource exhaustion.
Independent integration fixtures qualify explicit engine registration, immutable
envelopes and independent text-stage registration. The [source grammar](horos-syntax.md)
adds source and canonical round trips on the same AST. Architecture tests enforce the narrow external
dependency boundary, and static consumer tests distinguish Horos from every native
AST marker. Existing native-text/compiler gates remain required.

The separate [pagination mappings](semantic-generation.md) consume this model
and explicitly reject residual higher kinds. Runtime adapters and
Function/Stream/Cell bindings remain outside this frontend. These contracts retain
library ownership; none is a special core kind or builtin registration.
