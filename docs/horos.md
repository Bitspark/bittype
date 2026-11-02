# Horos declaration AST

**Implemented in Go source, 2 November 2026.** This is an AST-only frontend for
Horos declaration revision 1, not a readable source grammar or native generator.
The authoritative value format and binding rules are in
[horos's declaration specification](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md).

The separate [contract-type proposal](hydrated-types.md) places Wire and mixed
data/capability declarations in a Bittype-owned language provider above Horos.
This Horos frontend stays data-only with revision 1; it does not acquire that
language's kinds, profile or meaning. The additional provider is not implemented.

## API and ownership

- [syntax/horos/go](../syntax/horos/go/nodes.go) defines `Language`, `AST`, `Node`,
  `Kind` and the vocabulary. It imports only generic AST and language identity.
- [resolve/horos/go](../resolve/horos/go/resolve.go) provides
  `Resolve(root, dependencies, treeID)`, `Project(program, treeID)` and `Model(view)`.
  Resolve delegates semantics to Horos, then projects its checked model. Model
  checks the borrowed AST and returns a Horos program for inspection or explicit
  `Specialize(declarationID, groundArguments)`. No lookup performs network/file IO.
- [check/horos/go](../check/horos/go/check.go) admits AST syntax: kinds, roles,
  cardinalities and lexical spellings. Unresolved or semantically invalid names
  can reach analysis. Check does not resolve bindings or validate module pins.
- [analyze/horos/go](../analyze/horos/go/analyze.go) adapts Horos's authoritative
  resolution to [shared semantic facts and reports](analysis.md). It verifies
  supplied semantic annotations and never repairs them implicitly.
- [provider/horos/go](../provider/horos/go/provider.go) implements
  `text.CheckerProvider[horos.Language]` and `analysis.Provider[horos.Language]`,
  with ID `horos` and name `Horos`.
  Applications explicitly register it through `representation.NewLanguage` and
  `Engine.RegisterLanguage`, using profile `horos-declarations-1/ast-2`.
  Semantic analysis requires an additional `analysis.NewAnalyzer` and
  `Engine.RegisterAnalyzer`; implementing the interface installs nothing.

Only the Horos resolver, analyzer and their private bridge depend on Horos/Ontos.
The model dependency is pinned to a public immutable Go module revision in go.mod;
there is no local replacement or vendored private checkout. Generic infrastructure
and the four native text stages are unchanged and model-independent. The engine
can hold this AST without registering any transformer. Parse and Render report
missing text capabilities for this provider.

The syntax profile revision is a clean change: no ast-1 alias or fallback checker
remains. `resolve.Model` is still explicitly semantic, while `Pack` and `Check`
establish only the new syntax profile. Existing callers that require resolution
must use Model or the analysis capability.

## Vocabulary and bindings

One Program tree supplies the entry module and dependency pool. Projection of a
checked Program emits its complete resolved import closure, ordered by module
digest; imports/declarations retain Horos's canonical map order. Parameters,
arguments and descriptor operands retain their specified order. Analysis selects
the Program root as its explicit entry; its required root reference selects the
entry module. Other entry kinds or multiple trees are unsupported request forms
in this first analysis profile.

| Kind | Attributes | Children | References |
| --- | --- | --- | --- |
| Program | none | `module*` | `root` to the entry module |
| Module | `name`, `pin` (lowercase hex) | `import*`, `declaration*` | none |
| Import | `alias`, expected `name`, `pin` | none | optional `target` module assertion |
| Declaration | `name` | `parameter*`, one `body` | none |
| TypeParameter | `name`, canonical decimal `index` | none | none |
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

Explicit semantic analysis reconstructs canonical module values, asks Horos to
resolve them, and compares the authoritative syntax and supplied binding edges
with the analyzed AST closure.
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
declaration and tuple-index path context. Structured codes and subjects distinguish
invalid source, missing dependencies, contradicted context, invalid options and
exhaustion. AST structure/profile errors are located at current AST nodes; no
source location is invented or copied from historical origins.

Projection checks generic AST node, depth, entry and text budgets while creating
nodes, including hexadecimal growth before copying atoms. Generic structure is
checked before the AST is decoded. Horos's default declaration budgets apply
when recovering a model; analysis can accept an explicitly supplied bounded
resolver configuration. Both boundaries return empty output on failure; a model
accepted under custom larger limits is re-admitted under the fixed Horos defaults
before projection, so successful projection respects `resolve.Model`'s fixed
canonical budgets. The syntax checker has separate AST budgets; analysis can
select explicit resolver limits. Canonical and AST depth are distinct, and tests
qualify their boundary.

## Analysis profile and generic queries

`analyze.Contract()` advertises `horos-declarations-1/analysis-1`. Its domain is
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

All five reports are model publication requirements. A semantic rejection or
missing dependency publishes diagnostics without a partial graph. Missing context
and resolver exhaustion are Incomplete / Undetermined. A contradictory context
can reject the submitted program/context obligation without asserting that the
entry source itself is malformed; diagnostics identify the subject. Export
exhaustion after successful resolution preserves Complete / Holds validity but
leaves graph enumeration incomplete and publishes no model. Unknown requested
obligations receive Unsupported / Undetermined.

`LimitsSetting()` names provider-owned JSON configuration containing
`[MaxDepth, MaxNodes, MaxBytes, MaxModules]`. All four limits are required; absence
uses Horos defaults. This profile requires positive limits with depth at most 512.
Invalid configuration is an API error before source decoding. The kernel copies and
bounds these bytes without interpreting them. Semantic graph budgets are separate.

The graph preserves module/declaration/parameter identity, current input sites,
open parameters and zero-argument applications. Descriptor terms contain exact
Ontos codec-v1 bytes and authoritative operand paths. Four generic queries need
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
