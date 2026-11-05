# Resolution and shared semantic composition

**2 November 2026. Implemented in Go source.** Resolution is an independently
registered capability. It establishes declaration identities and reference
bindings in a provider-defined domain. Successful resolution alone does not
establish complete language validity. See the [provider matrix](providers.md).

## Shared infrastructure and ownership

| Package | Responsibility |
| --- | --- |
| [semantics/go](../semantics/go/model.go) | Immutable facts, graph checks and queries |
| [contracts/semantic/go](../contracts/semantic/go/contract.go) | Captured input snapshots, operation requests/results, reports, diagnostics, budgets and publication |
| [contracts/resolution/go](../contracts/resolution/go/provider.go) | Resolver contract, dependency requirements and typed registration adapter |
| [contracts/analysis/go](../contracts/analysis/go/provider.go) | Analyzer and its optional explicit resolution input |
| [engine/go](../engine/go/resolution.go) | Explicit capture, registration, capability lookup and invocation |
| [resolve/files/horos/go](../resolve/files/horos/go/resolver.go) | Horos resolution, authoritative semantic export and existing canonical-value entry helpers |
| [analyze/files/horos/go](../analyze/files/horos/go/analyze.go) | Consume explicit Horos evidence and reuse its sealed model |

The common semantic package replaces the former analysis-owned construction
records. There are no aliases at their old paths. The semantic fact model remains
language agnostic. Resolution rules, vocabulary and optional additional evidence
belong to the selected provider. Function, Stream and Cell remain library-owned.

## Resolver contract

```go
type Provider[L any] interface {
    language.Provider
    Resolve(context.Context, semantic.Request[L]) (semantic.Result, error)
}

type Contract struct {
    Semantic semantic.Contract
    Bindings semantics.Key
}
```

Bindings names an advertised obligation whose domain identifies the references
requiring resolution. Complete / Holds means each such reference has an
established target. Complete / Rejected means resolution is definitively
unsuccessful, such as an absent or ambiguous required name. Missing context,
exhaustion or an unmet prerequisite yields Incomplete / Undetermined. Unknown
requested obligations receive Unsupported / Undetermined. The generic kernel
validates publication integrity; coverage remains a trusted provider claim.

Publishing a model requires Complete / Holds for Bindings, plus each additional
required Complete/verdict pair. Initially failed resolution returns diagnostics and reports
without a partial model. Entities and bindings can exist without expression
terms; the independent binding-only test provider demonstrates this. A provider
may perform and advertise additional semantic work when its language requires it.

The opaque resolver is constructed with
`resolution.NewResolver[L](syntax, contract, provider)`. RegisterResolver installs
it only for an already registered compatible language. Duplicate registration is
an error; registration does not install an analyzer or text stages. Resolver(id)
returns the registered capability and its copied contract for inspection.

## Capture once, invoke explicitly

```text
semantic.Input (Entry, Trees, capture Limits)
    -> Engine.Capture -> immutable semantic.Snapshot
                            |                  |
                            v                  |
                         Resolve               |
                            |                  |
                            v                  v
                     resolution Result ----> Analyze
                                               |
                                               v
                                         analysis Result
```

Capture requires an explicitly registered `check.Checker`; registering a language
identity, resolver or analyzer does not install it. Standalone capture uses
`semantic.Capture(ctx, checker, input)`. Generic Pack/PackCST establish only carrier
integrity. A semantic Snapshot additionally records successful syntax admission.

`Engine.Resolve(ctx, snapshot, semantic.Options)` receives requested obligations,
configuration and invocation limits separately from captured syntax. An entirely
zero limits value selects defaults; partial invalid values fail. Snapshot admission
bounds inputs before syntax callbacks, retains the original immutable envelopes,
and performs no file/network lookup. Invocation bounds are checked again before
typed inputs are copied. MaxConfigurationBytes is one aggregate budget for
requested keys, own settings, dependency operation/profile identifiers, dependency
settings and report requirements. MaxConfigurationEntries counts all settings;
MaxRequested counts requested keys plus all dependency report requirements.
Every dependency is bounded before copying any dependency metadata. Operations
call providers outside registry locks.

`analysis.Input` carries Snapshot, Options and an optional Resolution result.
`analysis.Contract` carries a Semantic contract and an optional resolution
Requirement: exact profile, expected configuration keys/bytes and required report
verdicts. A nil requirement means analysis without a resolution dependency.
Supplying an undeclared dependency or omitting a required one is an API error.
No operation silently runs another provider or selects a fallback.

The dependency must be a genuine resolution publication for the identical opaque
snapshot, matching the expected profile/configuration. Equal tree IDs, syntax
bytes or recapturing identical input do not establish snapshot identity. Matching
profiles are explicit compatibility promises, not proof of semantic equivalence.
Configuration comparison ignores entry order but preserves exact value bytes.

Failed dependency reports are valid diagnostic inputs. An analyzer may inspect
them and publish its own reports without facts. It cannot publish a model while
declared dependency requirements are unmet. Each operation receives a fresh
private invocation identity: returning an earlier result fails even when both
invocations share the same snapshot. Non-nil errors discard all output.

## Semantics reuse

`semantic.Publication` accepts either new Facts or ReuseSemanticsFrom, never both.
The latter must identify an exact admitted dependency receipt with an available
model. Merely producing another result for the same snapshot is insufficient.
Publication checks current output requirements, dependency requirements and the
receiving model budgets before retaining the same immutable model handle.

The new result belongs to the new invocation. The reused model retains its entity
and term IDs. New graphs have their own ID spaces even if they share a snapshot;
there is no automatic merging or correspondence inferred from equal numbers.
Semantics accessors continue to return independent mutable copies. Reuse does not
copy truth reports automatically or weaken stale-result protection.

## Horos composition

The Horos provider implements syntax checking, resolution.Provider and analysis.Provider.
Its profiles are `horos-declarations-4/ast-1`,
`horos-declarations-4/resolution-1` and `horos-declarations-4/analysis-1`.
Earlier declaration profiles are removed without a compatibility path.

Horos resolution accepts multiple separate Program trees and an explicit root entry, resolves its
reachable pinned module closure, verifies supplied annotations and exports exact
descriptor payloads and operand paths. Alongside `reference-resolution`, it
advertises the existing validity, declaration-index, occurrence-binding,
application, reference-dependency and type-fact evidence. General kind facts
separate sort, result sort, declared contract, result contract, checked symbolic
body and memberships with assumptions. Its owning resolver currently checks the
whole required program; a resolution-only request does not promise less work.

Unknown imports/declarations can reject reference resolution. Other semantic
failures retain their own validity diagnostics and may leave reference resolution
incomplete. They are not automatically evidence of an unbound name. Missing
dependencies and resource exhaustion remain incomplete outcomes, not rejected
types. Export exhaustion preserves established reference resolution and program
validity without a model.

The analyzer consumes that explicit result, preserves the advertised analysis
reports and reuses the same model without repeating authoritative resolution.
`analyze.ContractWithResolutionConfiguration(settings)` selects exact expected
resolver settings; `analyze.Contract()` selects the default empty configuration.
Resolver settings belong to Resolve's Options, not the analyzer's Options.

Pagination retains imported Page, the caller's plain Person reference,
distinct same-spelled binders and opaque literal data. Specialization remains an
explicit Horos operation. Native pagination generation consumes these facts;
no native semantic analyzer or runtime binding is added.

## Qualification

[Resolver tests](../integration/go/resolution_test.go) cover independent registration,
binding-only facts, same-snapshot stale results, foreign captures, profile/configuration
checks, unmet prerequisites, receipt-specific reuse, budgets, cancellation and
concurrent composition. Existing [analysis tests](../integration/go/analysis_test.go)
retain adversarial capture/publication coverage through the new API. The
[Horos tests](../analyze/files/horos/go/analyze_test.go) exercise explicit resolve/analyze
composition and the [independent consumer](../integration/go/horos_test.go) checks
model identity reuse and specialization. Architecture and negative type-check
fixtures enforce package and marker boundaries. Required checks run remotely.
