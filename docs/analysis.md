# Semantic analysis capability

**2 November 2026. Implemented in Go source.** Semantic analysis establishes what
syntax means in its supplied context. For `pagination.Page<Person>`, it can
establish which declarations Page and Person denote, whether the application has
the right number of arguments, and how Person binds to Page's parameter.

Bittype provides the common request, result and fact representation. Each language
provider supplies the rules that establish those facts. The engine selects an
explicitly registered provider; it starts empty and imports no concrete analyzer.
Horos, Go, TypeScript and Rust have production semantic providers. Native coverage is
bounded and compiler-owned; C++ currently has syntax tooling only.

This guide connects the concepts, package boundaries and a pagination example.
The [semantic model reference](semantics.md) specifies the fact records and their
invariants; the [Horos reference](horos.md) specifies that provider's vocabulary
and coverage. The [repository layout](../LAYOUT.md) is the full dependency policy.

## Concepts

There are three kinds of information: the syntax supplied by the caller, facts
established by the provider, and reports describing what the provider completed.
Keeping them separate lets a caller distinguish an unresolved name in an AST
from an established binding, or missing dependency context from invalid source.

| Concept | Meaning in this API | Pagination example |
| --- | --- | --- |
| Syntax envelope | An owned, structurally checked AST snapshot; Capture separately requires syntax admission | A Program tree containing the people module and its supplied pagination dependency |
| Entity | The identity of a module, declaration, parameter or provider-owned intrinsic | Page, Page's T, Person and PeoplePage are distinct entities |
| Term | An established expression: a Reference, Application or ProviderExpression | Applying Page to Person is an Application; using Page's T is a Reference |
| Occurrence and binding | A current AST node address plus a role and ordinal, and its established target entity | The target occurrence of `pagination.Page<Person>` binds to the imported Page declaration |
| Dependency | A provider-labeled directed edge between entities | PeoplePage applies Page and also applies Person as its argument |
| Obligation | A named analysis task for which the caller can request a report | Decide program validity, or enumerate application bindings |
| Domain | The part of the supplied context covered by an obligation | The selected Horos entry module and every pinned module transitively imported from it |
| Analysis profile | A versioned contract defining obligations, domains, vocabulary and requirements for a model | `horos-declarations-4/analysis-1` |
| Report | An obligation's completion status and, where applicable, its truth verdict | Validity is Complete / Holds; a completed binding enumeration is Complete / Undetermined |
| Semantics | The immutable graph of established facts, available when its publication requirements are met | Queryable entities, terms, occurrences and application dependencies |

Entity and term IDs belong to one model. Names are display metadata, so two
parameters both named T remain distinct. An occurrence uses a current tree/node
address; historical source origins are separate metadata. Provider vocabulary
uses `semantics.Key{Owner, Name}`: the kernel stores and compares these keys
without inventing their meaning.

A syntax profile specifies admitted AST forms. An analysis profile specifies
what the analyzer establishes about those forms. Neither profile is inferred
from a structure's shape or from a declaration name.

## Ownership and package layout

The layout follows processing stages. The final `go` denotes implementation
language: `analyze/files/horos/go` is Go code implementing analysis for Horos.

| Layer and package | Responsibility |
| --- | --- |
| [trees/ast/go](../trees/ast/go/ast.go) and [trees/representation/go](../trees/representation/go/ast.go) | Generic syntax storage, current node addresses and structurally checked snapshots |
| [semantics/go](../semantics/go/model.go) | Shared fact records, graph integrity, immutable storage and generic queries |
| [contracts/semantic/go](../contracts/semantic/go/contract.go) | Shared capture, requests, reports and publication |
| [contracts/resolution/go](../contracts/resolution/go/provider.go) | Independent resolution capability and requirements |
| [contracts/analysis/go](../contracts/analysis/go/provider.go) | Analysis capability and explicit resolution dependency |
| [engine/go](../engine/go/analysis.go) | Explicit registration and selection of a compatible analyzer |
| [syntax/horos/go](../syntax/horos/go/nodes.go) and [check/syntax/horos/go](../check/syntax/horos/go/check.go) | Horos AST vocabulary and syntax admission |
| [resolve/files/horos/go](../resolve/files/horos/go/resolver.go) | Registered resolution, authoritative fact export and explicit canonical-value helpers |
| [analyze/files/horos/go](../analyze/files/horos/go/analyze.go) | Consume explicit resolution evidence and reuse its immutable model |
| [internal/horosmodel](../internal/horosmodel) | Private decoding and correspondence between current AST sites and the authoritative model, used by the resolver |
| [provider/horos/go](../provider/horos/go/provider.go) | Compose Horos checking, resolution and analysis capabilities for explicit application registration |
| External horos and ontos packages | Horos owns declaration/type rules; Ontos owns values and encoding |

The engine depends on shared contracts and representations. The specialized contracts
depend on shared semantic infrastructure and the semantic model. Neither imports a concrete language. Only
the Horos resolver and private bridge directly consume the external Horos/Ontos
model. The analyzer consumes their shared semantic evidence. The application chooses the concrete provider and connects these layers.

The kernel checks that a published graph is well formed and belongs to the current
request. The provider is responsible for whether its semantic claims are true.
Generic facts do not supply universal rules for lookup, assignability or type
equivalence. Another language's analyzer can use the same API with different
rules. Function, Stream and Cell remain ordinary library declarations; their
names and behavior do not become kernel kinds or dispatch cases.

## Pagination: from input to queries

### 1. Supply the AST and its context

The [canonical pagination example](model.md#pagination-first) has this meaning:

```text
module example.pagination
  Page<T> = closed record { items: list-of(T), next?: scalar(utf8-text) }

module example.people
  import pagination = exact pinned example.pagination module
  Person = closed record { name: scalar(utf8-text) }
  PeoplePage = pagination.Page<Person>
```

This is explanatory notation, not implemented Horos source syntax. The analyzer
receives an AST, not this text. Its Program node contains both modules and has a
root reference selecting `example.people`. Within the tree, Page has a parameter
declaration; its body's parameter use stores slot 0. PeoplePage's body is an
imported Application of Page whose argument is a local, zero-argument Application
of Person. Descriptor nodes retain record/list structure and their operands.

The application packs this tree with the Horos identity/profile descriptor.
Packing validates generic carrier structure. Capture invokes the independently
registered Horos checker to admit node kinds, spelling and syntax; optional
semantic target annotations may be absent. Neither packing nor syntax checking
proves that the names resolve or the pins agree.
`resolve.Resolve` is an existing way to obtain an already resolved AST projection
from canonical module values; analysis can also admit structurally valid ASTs
whose bindings have not yet been established.

For this profile, `semantic.Input.Trees` contains one or more separate Program envelopes. `Entry`
addresses the selected Program root. Capture produces the shared snapshot. Empty Requested
options select every obligation advertised by that operation. Dependencies are
supplied as whole modules in the Program trees; the engine fetches nothing.

### 2. Establish meaning through the owning provider

`Engine.Resolve` selects the registered Horos resolver for the captured snapshot.
The adapter reconstructs canonical module values, delegates resolution to Horos,
and verifies any supplied semantic annotations against that result.

Horos checks module pins, imports, declaration lookup, parameter slots and
application arity under its declaration rules. Page resolves in the imported
pagination module; Person resolves in the caller's people module. A same-spelled
Person in the pagination module cannot capture that argument. All supplied module
pin assertions are checked; semantic facts cover the entry's reachable import
closure, excluding unused pool modules.

### 3. Export facts without expanding the generic

The resolver publishes separate entities for the two modules, Page, Page's T,
Person and PeoplePage. The significant term links can be read as follows;
the notation labels records, rather than specifying a serialization:

```text
PeoplePage.Body -> Application(target: Page,
                              arguments: [Application(target: Person, arguments: [])])
Page.Parameters -> [Page.T]
Page.Body       -> ProviderExpression(record descriptor, located operands ...)
items operand  -> ProviderExpression(list-of descriptor,
                                    operand: Reference(target: Page.T))
```

The record/list meanings stay owned by Horos. ProviderExpression retains exact
encoded descriptor bytes and the operand locations supplied by Horos's model;
the kernel does not decode records or search literal bytes for possible bindings.

`DirectParameterBindings` pairs Page's parameter list with the application's arguments:
Page.T binds to the zero-argument Person application. The open Page declaration
is preserved. In `Outer<T> = Page<T>`, Page.T instead binds to a Reference to
Outer's distinct T entity. Analysis does not expand these declarations. Explicit
Horos `Specialize` can separately produce a ground descriptor; native generation
and runtime behavior are separate operations.

### 4. Read coverage before using the model

The application supplies the resolution result and the same snapshot to Analyze.
Horos analysis consumes that evidence and reuses its model without resolving again.
On success, `program-validity` reports Complete / Holds. The other four reports
are Complete / Undetermined: they enumerate declarations, occurrence bindings,
expressions/applications and application dependencies, without asserting a truth
verdict. All five are requirements for publishing this profile's model.

Expected outcomes differ even when the call returns a nil Go error:

| Situation | Validity report | Enumeration reports | Semantics |
| --- | --- | --- | --- |
| Valid pagination input, sufficient resources | Complete / Holds | Complete / Undetermined | Available |
| Page receives the wrong number of arguments | Complete / Rejected | Incomplete / Undetermined | Absent |
| The pinned pagination dependency is missing, or resolution exhausts its budget | Incomplete / Undetermined | Incomplete / Undetermined | Absent |
| Resolution and annotation checks succeed, but graph export exhausts its budget | Complete / Holds | Incomplete / Undetermined | Absent |

Inspect `Result.Reports()` or `Result.Report(key)`, read `Result.Diagnostics()`
for the reason and subject, and check the boolean returned by `Result.Semantics()`
before querying facts. `Result.Valid()` establishes successful publication, not
source validity. API errors such as invalid configuration return a non-nil Go
error and no usable result; they do not use the table's diagnostic-result path.

### 5. Query the published facts

At the current Page application node, the Horos profile exposes two occurrences:
`pageTarget` uses `resolve.TargetRole()` and `pageExpression` uses
`resolve.ExpressionRole()`, both with ordinal 0. The following query arguments
are symbolic labels for actual occurrences and IDs in this result:

| Consumer question | Generic model query | Answer for this example |
| --- | --- | --- |
| Which declaration does this use denote? | `Resolve(pageTarget)` | Page's entity ID |
| Which expression was established here? | `TermAt(pageExpression)` | The term ID of the Page application shown above |
| What fills each generic slot? | `DirectParameterBindings(pageApplicationTermID)` | One pair: Page.T's entity ID and the Person application's term ID |
| Which declarations does PeoplePage apply? | `DependenciesFrom(peoplePageID, resolve.ApplicationDependency())` | Edges to Page and Person, each with its current application occurrence |

`Entity(id)` and `Term(id)` inspect the referenced records. Consumers can also
enumerate the established occurrences in `Semantics.Data()` instead of constructing
provider-specific occurrences. A failed lookup means no recorded answer at that
occurrence; it does not establish that the source is invalid. Application
dependencies are explicit edges, not an inferred import or runtime-effect graph.

These queries need no external Horos model import and do not mutate the AST or
model. The [pagination analyzer tests](../analyze/files/horos/go/analyze_test.go) qualify
bindings, caller scope and reports; the [independent consumer](../integration/go/horos_test.go)
exercises registration and shared queries through the public API.

## Composition

```text
application selects syntax descriptor, checker, resolver and analyzer
    -> RegisterLanguage
    -> RegisterChecker (independent syntax admission)
    -> RegisterResolver / RegisterAnalyzer (independent capabilities)

syntax envelopes + explicit entry and dependency inputs
    -> Engine.Capture -> immutable Snapshot
    -> Engine.Resolve(snapshot, options) -> resolution Result
    -> Engine.Analyze(snapshot, resolution Result, options)
    -> fresh typed invocation on the same Snapshot
    -> language provider's Analyze
    -> Request.Publish(reused model, reports, diagnostics)
    -> immutable request-bound Result
```

`Check` admits syntax under its declared profile. It does not return a semantic
model or imply that names resolve. `Analyze` delegates semantic rules to the
provider for that language. `TransformTo` follows explicitly registered mappings;
analysis neither installs a transformer nor routes through a canonical language.

```go
type Provider[L any] interface {
    language.Provider
    Analyze(context.Context, Request[L]) (semantic.Result, error)
}
```

`analysis.NewAnalyzer[L](syntax, contract, provider)` creates an opaque capability.
It checks the actual Go marker, syntax/profile descriptor and provider metadata.
`Engine.RegisterAnalyzer` attaches exactly one capability to an already registered
compatible language. Duplicate analyzers, unknown languages and descriptor
mismatches are errors. There is no inferred registration or replacement.

The registry lock is released before syntax checks and provider callbacks.
Concurrent callers require a concurrency-safe provider. The engine does not
sandbox callbacks, recover arbitrary panics or certify semantic claims.

## Request and snapshot

`semantic.Input` supplies Entry, Trees and capture Limits to Engine.Capture.
Capture requires the registered check.Checker for the language; identity metadata
alone is insufficient. Outside an engine, use
`semantic.Capture(ctx, checker, input)` with an explicit configured checker.
`semantic.Options` supplies Requested obligation keys, Configuration and invocation
Limits. Empty Requested selects all advertised obligations; entirely zero limits
select defaults. `analysis.Input` then supplies:

| Field | Meaning |
| --- | --- |
| Snapshot | An immutable captured input collection admitted by the explicitly selected syntax checker |
| Options | This analysis invocation's requested obligations, configuration and limits |
| Resolution | Optional explicit dependency receipt, required when declared by the analyzer |

The [resolution contract](resolution.md) defines profile/configuration matching,
exact snapshot identity and model reuse. `analysis.Contract.Semantic` contains the
analysis profile and model requirements; its optional Resolution requirement
specifies prerequisite reports. Matching tree IDs never substitutes for the same
snapshot. Each invocation still has fresh private identity.

No step fetches imports, reads files or consults an ambient package environment.
The application supplies all dependency context explicitly. Configuration and
requested keys are copied. Request accessors provide independent copies of input
trees and mutable metadata; provider mutation cannot alter the captured request.
Callers must not mutate construction storage concurrently with capture.

Each request has fresh private identity, even when tree IDs, entry and options
are unchanged. Only its `Request.Publish` can create a result associated with that
invocation. Returning a retained result from an earlier call fails with
`ErrStaleResult`; the engine never stamps old output with new request identity.
This proves association, not freshness of arbitrary facts a dishonest provider
copies into a new publication. Providers remain trusted implementations.

Context cancellation is checked before and after invocation. Prompt interruption
inside a synchronous resolver requires provider cooperation. The Horos adapter
uses bounded resolution; it does not claim to preempt a running resolver.

## Coverage and validity

An analyzer advertises an analysis `Contract`: a versioned profile, obligations
with provider-owned domain keys, and requirements for publishing a model.
The profile documents each obligation's subject/domain relative to entry and inputs.
It is separate from the syntax profile and is not a computed fingerprint.

Each requested obligation must have exactly one report:

| Completion | Verdict | Meaning |
| --- | --- | --- |
| Complete | Holds or Rejected | A completed decision about the declared subject |
| Complete | Undetermined | A completed enumeration that asserts no truth verdict |
| Incomplete | Undetermined | Necessary context or resources were unavailable |
| Unsupported | Undetermined | This provider does not implement the requested obligation |

Zero-valued Unreported is never accepted as a report. Unknown requested keys must
be Unsupported. Additional reports may establish advertised model requirements.
`semantic.Contract.SemanticsRequirements` lists the exact Complete/verdict pairs that a facts
publication must establish. Missing or contrary reports prevent publication.

Consumers distinguish three questions: did the analysis complete, what verdict
was established, and is a model available? `nil` error does not mean Holds.
`Result.Valid()` means a properly published result, not valid source.

`Request.Publish` atomically validates coverage, current diagnostic sites and
graph integrity, then copies and seals output before the callback returns.
`Result` and its `semantics.Semantics[L]` expose no mutable stored slices or maps.
Non-nil Go errors discard all output, including a result returned alongside the
error. Expected source rejection or missing context is normally a diagnostic
result with no graph; infrastructure failures remain Go errors.

## Horos profile

The [Horos resolver and analyzer](horos.md) compose through shared evidence. Resolution
publishes only a completely resolved graph for its selected entry and reachable
pinned dependency closure. It distinguishes source rejection, missing context,
contradicted context, invalid configuration and resource exhaustion using Horos's
structured error codes and subjects, never diagnostic message matching.

<!-- bittype:overview scope=horos-profile part=features -->
The unified Horos profile exports general kinds as provider-owned terms and
facts. Semantic sort, declared contract, checked symbolic body and individual
memberships are separate relations. Type and kind definitions have distinct
provider-owned entity categories. No single fact is a principal inferred kind.
Membership facts retain their assumptions; missing facts assert nothing.
Program validity includes Horos's symbolic obligations on open declarations.
A syntactically valid `Forward<T> = DataPage<T>` can therefore be rejected before
specialization. Shared semantic carriers do not interpret these kind rules.
<!-- bittype:overview:end -->

An unlabelled dependency pool containing Q when import P is required is missing P;
it does not prove that Q was supplied as P. A present explicit assertion
contradicting the resolved pin/name or binding is a different failure.

<!-- bittype:overview scope=horos-profile part=exclusions -->
The Horos syntax profile is `horos-declarations-4/ast-1`: optional semantic binding
annotations may be absent, while supplied annotations in the analyzed closure
must agree with analysis. The required entry/root selection is structural.
Analysis does not repair input or turn stale references into new references
silently. The explicit `resolve.Model` helper remains semantic and shares the
same authoritative path.
<!-- bittype:overview:end -->

[Bounded Go, TypeScript and Rust semantics](native-semantics.md) use the same resolution and
analysis interfaces. Broader native coverage, C++ analysis, general
scope/relation APIs, partial checked graphs and incremental caches remain deferred.
The independent test-only language still qualifies the agnostic contract.

See the [provider matrix](providers.md) for the independently registered capabilities.

## Limits and evidence

Capture bounds input trees, nodes, metadata entries and bytes before syntax
callbacks. Each invocation bounds its requested keys, configuration and dependency
metadata together before provider invocation. Publication bounds diagnostics and semantic records, links,
payload bytes and containment depth. Horos resolver limits and semantic export
limits are separate: successful resolution does not authorize unbounded export.
Exhaustion cannot masquerade as an empty valid model or a rejected type.

The maintained remote gates exercise existing text/native compiler checks
alongside new semantic model, registration, ownership, request identity, coverage,
Horos pagination and foreign-provider tests. The earlier
[consultation delivery record](../research-docs/0002-semantic-analysis-kernel.md)
records the original semantic slice, not the subsequent resolver split; its frozen expert input and raw
advice remain separate from the implemented specification.

See [native semantic providers](native-semantics.md)
for the bounded Go, TypeScript and Rust resolvers/analyzers, immutable native facts and
explicit Go/TypeScript compiler queries. C++ semantics and broader Rust coverage remain planned.
