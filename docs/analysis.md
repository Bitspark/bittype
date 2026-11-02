# Semantic analysis capability

**2 November 2026. Implemented in Go source.** Analysis is an optional capability
beside syntax checking, text stages and AST transformation. Applications register
it explicitly; the engine starts empty and imports no concrete analyzer.

The [semantic graph](semantics.md) contains shared facts. The public SPI is in
[contracts/analysis/go](../contracts/analysis/go/contract.go); registration and
selection are in [engine/go/analysis.go](../engine/go/analysis.go).

## Composition

```text
application selects syntax descriptor and analyzer
    -> RegisterLanguage / RegisterProvider
    -> RegisterAnalyzer

syntax envelope + explicit entry, inputs, options, requested obligations
    -> Engine.Analyze
    -> typed request capture
    -> language provider's Analyze
    -> Request.Publish(facts, reports, diagnostics)
    -> immutable request-bound Result
```

`Check` admits syntax under its declared profile. It does not return a semantic
model or imply that names resolve. `Analyze` delegates semantic rules to the
provider for that language. `TransformTo` follows explicitly registered mappings;
analysis neither installs a transformer nor routes through a canonical language.

```go
type Provider[L any] interface {
    language.LanguageProvider
    Analyze(context.Context, Request[L]) (Result, error)
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

`analysis.Input` supplies:

| Field | Meaning |
| --- | --- |
| Entry | An existing node in the captured tree collection; provider defines admitted entry kinds |
| Trees | Syntax envelopes with distinct tree IDs and the same language/profile/marker |
| Requested | Provider-owned obligation keys; empty requests all advertised obligations |
| Configuration | Owned keys and bounded uninterpreted bytes, decoded only by the provider |
| Limits | Capture/publication budgets; entirely zero selects defaults, partial invalid values fail |

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
`Contract.ModelRequirements` lists the exact Complete/verdict pairs that a facts
publication must establish. Missing or contrary reports prevent publication.

Consumers distinguish three questions: did the analysis complete, what verdict
was established, and is a model available? `nil` error does not mean Holds.
`Result.Valid()` means a properly published result, not valid source.

`Request.Publish` atomically validates coverage, current diagnostic sites and
graph integrity, then copies and seals output before the callback returns.
`Result` and its `semantics.Model` expose no mutable stored slices or maps.
Non-nil Go errors discard all output, including a result returned alongside the
error. Expected source rejection or missing context is normally a diagnostic
result with no graph; infrastructure failures remain Go errors.

## Horos profile

The [Horos analyzer](horos.md) adapts Horos's checked declaration Program. It
publishes only a completely resolved graph for its selected entry and reachable
pinned dependency closure. It distinguishes source rejection, missing context,
contradicted context, invalid configuration and resource exhaustion using Horos's
structured error codes and subjects, never diagnostic message matching.

An unlabelled dependency pool containing Q when import P is required is missing P;
it does not prove that Q was supplied as P. A present explicit assertion
contradicting the resolved pin/name or binding is a different failure.

The Horos syntax profile is `horos-declarations-1/ast-2`: optional semantic binding
annotations may be absent, while supplied annotations in the analyzed closure
must agree with analysis. The required entry/root selection is structural.
Analysis does not repair input or turn stale references into new references
silently. The explicit `resolve.Model` helper remains semantic and shares the
same authoritative path.

Native semantic analyzers, scope lookup, relation queries, partial checked graphs
and incremental caches remain deferred. A separate test-only unknown-language
provider qualifies the generic contract without adding a production native checker.

## Limits and evidence

Capture bounds inputs, nodes, metadata entries, bytes and configuration before
provider invocation. Publication bounds diagnostics and semantic records, links,
payload bytes and containment depth. Horos resolver limits and semantic export
limits are separate: successful resolution does not authorize unbounded export.
Exhaustion cannot masquerade as an empty valid model or a rejected type.

The maintained remote gates exercise existing text/native compiler checks
alongside new semantic model, registration, ownership, request identity, coverage,
Horos pagination and foreign-provider tests. The
[consultation delivery record](../research-docs/0002-semantic-analysis-kernel.md)
records exact revisions and validation evidence; the frozen expert input and raw
advice remain separate from the implemented specification.
