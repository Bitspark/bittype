# Horos declaration AST

**Implemented in Go source, 4 November 2026.** One AST-only frontend projects
Horos declaration revision 2. Type/Data bounds restrict generic slots within that
language; data and capability expressions share one module, declaration and
binder model. The old Horos/Wire packages and profiles are removed without aliases.
The authoritative [canonical format](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md)
lives in Horos. No readable parser, native generator or runtime binding is added.

The [dependency alignment](dependencies.md) and independent installed consumer
qualify the public module graph. Private projection and binding mechanics live
in `internal/horosmodel`; syntax framing lives in `internal/declsyntax`.
The agnostic kernel contains no Type/Data or Wire interpretation.

## API and ownership

- [syntax/horos/go](../syntax/horos/go/nodes.go) defines `Language`, `AST`, `Node`,
  `Kind` and the vocabulary. It imports only generic AST and language identity.
- [resolve/horos/go](../resolve/horos/go/resolve.go) provides
  `Resolve(root, dependencies, treeID)`, `Project(program, treeID)` and `Model(view)`.
  Resolve delegates semantics to Horos, then projects its checked model. Model
  checks the borrowed AST and returns a Horos program for inspection or explicit
  `Specialize(declarationID, closedArguments)`. No lookup performs network/file IO.
  The registered resolver uses [ResolveReferences](../resolve/horos/go/resolver.go)
  to export shared semantic facts. It verifies supplied semantic annotations
  against authoritative resolution and never repairs them implicitly.
- [check/horos/go](../check/horos/go/check.go) admits AST syntax: kinds, roles,
  cardinalities and lexical spellings. Unresolved or semantically invalid names
  can reach resolution. Check does not resolve bindings or validate module pins.
- [analyze/horos/go](../analyze/horos/go/analyze.go) consumes the explicit resolver
  result, reusing its model and adapting its evidence to
  [shared semantic facts and reports](analysis.md).
- [provider/horos/go](../provider/horos/go/provider.go) implements
  `check.Provider[horos.Language]`, `resolution.Provider[horos.Language]`
  and `analysis.Provider[horos.Language]`,
  with ID `horos` and name `Horos`.
  Applications explicitly register it through `representation.NewLanguage` and
  `Engine.RegisterLanguage`, using profile `horos-declarations-2/ast-1`.
  Syntax checking is separately bound with `check.NewChecker[h.Language](definition,
  provider)` and installed with `Engine.RegisterChecker`. Identity registration
  alone supplies no syntax or semantic capability.
  Resolution requires `resolution.NewResolver` and `Engine.RegisterResolver`.
  Analysis additionally requires `analysis.NewAnalyzer` and
  `Engine.RegisterAnalyzer`, and consumes the explicit resolution result.
  See [shared composition](resolution.md); implementing interfaces installs nothing.

Only the Horos resolver and its private bridge directly import Horos/Ontos.
The analyzer consumes the resolver contract and shared evidence.
The model dependency is pinned to a public immutable Go module revision in go.mod;
there is no local replacement or vendored private checkout. Generic infrastructure
and the four native text stages remain independent of the external model. The engine
can hold this AST without registering any transformer. Parse and Render report
missing text capabilities for this provider.

The syntax profile revision is a clean change: no old-profile alias or fallback checker
remains. `resolve.Model` is still explicitly semantic, while `Pack` and `Check`
establish only the new syntax profile. Existing callers that require resolution
must use Model or the analysis capability.

## Vocabulary and bindings

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
| Declaration | `name` | `parameter*`, one `body` | none |
| TypeParameter | `name`, canonical decimal `index` | one `bound` | none |
| Bound | `name` is `type` or `data` | none | none |
| Parameter | canonical decimal `index` | none | optional `target` TypeParameter assertion |
| Application | `scope` (`local` or `imported`), `name`, imported `alias` | `argument*` | optional `target` Declaration and imported `import` assertions |
| Descriptor | Horos `constructor` name without `horos/` prefix | `operand*` | none |
| Atom | raw bytes as lowercase `hex` | none | none |
| Tuple | none | `element*` | none |

Descriptor operands retain the canonical value layout. Only Horos-designated
type positions become Parameter, Application or Descriptor nodes. Other positions
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
opacity. Existing data-only fixtures now spell their Data bounds explicitly.
The [unified bound tests](../analyze/horos/go/bounds_test.go) add:

```text
Page<T> = list-of(T)
DataPage<D: Data> = Page<D>
People = Page<Person>
Live = Page<Wire<Person>>
```

The [installed consumer](../integration/go/testdata/installed/consumer_test.go)
uses the full closed record with required items and optional text next. Both
applications retain the same open Page. Horos Specialize produces a closed type
expression; DataDescriptor additionally proves data and returns the existing
Horos descriptor. A mixed result cannot be passed to ground WF/Conforms.
No specialization mutates the original generic.

`Forward<D: Data> = DataPage<D>` resolves; `Forward<T> = DataPage<T>` is a located
semantic bound violation even before specialization. Both are valid AST syntax.
Check validates the bound spelling/cardinality only; Horos supplies the meaning.
The resolver exports bounds through existing provider-owned semantic terms and
analysis consumes the same immutable receipt. No new kernel obligation is needed.

## Admission and qualification

Horos owns module-name/pin checking, arity, binding scope, recursion closure,
application-cycle rejection, symbolic generic bounds and specialization limits. Its errors retain module,
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

`analyze.Contract()` advertises `horos-declarations-2/analysis-1`. Its domain is
the selected entry and transitive pinned module closure, not every unused value
in the supplied pool. Unused pool modules contribute no established semantic
facts; their binding annotations are outside this analyzed domain. Every supplied
module still has its document-pin assertion checked during decoding.

| Obligation | Successful report |
| --- | --- |
| program-validity | Complete / Holds for the submitted entry and required context |
| declaration-index | Complete / Undetermined; module, declaration and parameter entities enumerated |
| occurrence-bindings | Complete / Undetermined; parameter/application targets enumerated |
| applications | Complete / Undetermined; references, effective applications and provider expressions exported |
| application-dependencies | Complete / Undetermined; every declaration application, including nested arguments, exported |

The resolver adds a sixth obligation, reference-resolution, requiring Complete /
Holds for its model. Its profile is horos-declarations-2/resolution-1. Analysis
requires that resolver evidence and preserves these five reports. All five
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
open parameters and zero-argument applications. Descriptor terms contain exact
Ontos codec-v1 bytes and authoritative operand paths. Parameter bounds are
provider-owned expression facts: `TermAt(parameterSite, resolve.BoundRole())`
returns a term with operator `horos/bound/type` or `horos/bound/data`, and the
codec-v1 bytes of that bound atom. Parameter entities retain sort `horos/type`;
bounds constrain arguments, not the kind of entity. No shared semantic field or
subtype rule is added. ProgramValidity includes Horos's symbolic bound checks. Four generic queries need
no Horos model import: `Resolve`, `TermAt`, `ParameterBindings` and
`DependenciesFrom`. Dependency edges describe declaration applications, not a
general import, inheritance or runtime-effect graph.

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
