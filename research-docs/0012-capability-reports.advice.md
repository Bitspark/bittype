# Advice for consult
**Research doc:** doc_a5dee02f-b664-406e-a746-778b5ee5f7c8.md
**Document:** doc_a5dee02f-b664-406e-a746-778b5ee5f7c8
**Advisor:** ChatGPT, model "6 Pro" as labelled on the answer page
**Session:** nsess_e1a4ca3b-3193-4766-9f81-9dcc50e33a6b
**Run:** run_43b89675-3878-47ea-93fd-2452942f7f4d
**Chat:** https://chatgpt.com/c/6aedf281-8838-832a-9701-fb9709a814b8
**Verify:** nightfall consult verify run_43b89675-3878-47ea-93fd-2452942f7f4d --output <this file>

## Recommendation

**Proceed with the proposal, using committed static reports plus execution reports attached to individual CI runs.** That publication split fits the stated requirements and avoids making the matrix updater responsible for deciding whether another workflow’s results are still applicable.

I would refine the design before implementation in five places:

1. Make extraction **definition-centered**, with the host retaining all authority over provider, profile, and test associations.
2. Represent extraction limits at the individual field level, so useful structural facts survive alongside unknown or opaque behavior.
3. Preserve wrapper and delegate identities separately. The supplied `entry()` traversal establishes a static call relationship; it does not establish behavioral equivalence or transfer test execution evidence.
4. Attach observed outcomes to **test executions**, then show their relationship to authored cases.
5. Make the public host runner independently composable, and state resource guarantees for the entire command.

This review uses the baseline and code excerpts in research(9).md[research(9).md](sandbox:/workspace/scratch/344c0e738ef4/upload/research(9).md). The recommendations below address the five consultation questions.

## 1. Report vocabulary and model

### Keep the evidence layers independent

The report should answer several different questions without collapsing them into one support status.

| Layer | What it establishes | Useful presentation |
|---|---|---|
| Discovered implementation | A specific contribution satisfies discovery and construction requirements in the inspected build context | Implementation identity and available dependency alternatives |
| Authored scope | The owner states an intended boundary, profile, or exclusion | Links to the authoritative scope sections |
| Declared structure | An extractor interprets a particular source declaration according to its owning format | Forms, roles, attributes, cardinalities, and references |
| Extraction limits | Some information is unavailable, unsupported, or opaque | Localized limitations with source links |
| Declared case association | An author associates a test function with a target and an expected assertion | Case label, target, expectation, profile, and test declaration |
| Execution observation | A supplied invocation reports an outcome for a named test | Pass, fail, skip, or incomplete execution, with invocation context |
| Evidence applicability | The supplied observation can be joined to this declaration snapshot | Matched, mismatched, or insufficient context |
| Architectural case | A maintained human explanation identifies an architectural failure and its repair state | Independent ACX reference |

These are separate dimensions. A form can simultaneously have known structural roles, an opaque predicate, an authored exclusion, and a passing associated test.

Prefer labels that identify what was assessed:

- “No associated cases found in the selected test sources.”
- “Associated test function; no matching execution observation supplied.”
- “Predicate present; behavior not interpreted.”
- “Allowed-kind expression could not be decoded.”
- “Observed skip in invocation R.”
- “Authored exclusion; see scope.”

The first two labels matter especially. Missing associations do not establish that no tests exist, and missing observations do not establish that a test never ran.

### Describe source declarations precisely

For the grammar extractor, the routine claim should be:

> The selected initializer declares these forms and constraints under the grammar format’s interpretation.

Even when every initializer expression is decoded, the tool has not established that mutable variables retain those values during execution. That qualification applies to both `rules` and referenced slices such as `types`.

The supplied `mapped-type` entry supports a useful inventory:

| Facet | Declared information |
|---|---|
| `type-parameter` child | Exactly one |
| `name-type` child | Optional |
| `type` child | Exactly one |
| `readonly-modifier` attribute | Optional; finite declared values `add`, `remove` |
| `optional-modifier` attribute | Optional; finite declared values `add`, `remove` |
| Additional checker logic | Outside the table interpretation; link to the operation and authored scope |

The procedural restrictions described in the document remain outside that structural inventory. The report should always explain that operation bodies are outside this extractor’s interpretation, even when it has found no particular custom predicate.

A further edge case deserves explicit treatment: an undecodable map key can potentially collide with a known key. In that situation, report known **source entries** and the unresolved portion; do not manufacture a complete effective map.

### Model uncertainty locally

Use a definition-level extraction disposition together with field-level knowledge.

At the definition level, distinguish:

- No definition attached.
- Definition attached, but no matching extractor.
- Selected declaration decoded within the supported format.
- Declaration partly decoded.
- Source unavailable.

At the field level, distinguish known values, known absence or lack of constraint, unknown information, and opaque interpretation.

For example:

- An unknown allowed-kind expression is different from an empty finite set.
- An unbounded maximum is different from an unknown maximum.
- An opaque predicate is different from having no predicate.
- A known optional attribute remains useful even when its custom value predicate is opaque.

Avoid a bare `Complete bool` whose meaning changes between callers. If the API includes an extraction extent, define it narrowly around the selected declaration and supported fact schema.

### Make unassessed areas navigable

The useful views are lists of:

- Extracted forms without associated cases.
- Associated cases without applicable observations.
- Opaque or undecoded fields.
- Attached definitions without extractors.
- Operations with scope documentation but no structural inventory.

The scope prose remains authoritative and need not become an enumerable feature catalogue. These views therefore expose gaps without supplying a language-completeness denominator.

A printer or semantic provider can have a useful report consisting of identity, scope, cases, and observations. Its lack of a grammar inventory is an explicit characteristic of the report.

## 2. Minimal extractor contract and ownership boundary

### Keep association authority in the host

The extractor normally needs only a definition root, a bounded source view, and limits. The host already knows which operation, language, and profile caused that definition to be selected.

A suitable public shape is:

```go
type Extractor interface {
    Selectors() []Selector
    Extract(context.Context, DefinitionRequest) (Inventory, error)
}

type DefinitionRequest struct {
    Root   DefinitionRef
    Source SourceView
    Limits Limits
}

type Inventory struct {
    Features    []Feature
    Limitations []Limitation
    Extent      ExtractionExtent
}
```

The supporting types should be concrete, documented records. This is an API sketch, not a proposed serialization protocol.

A selector identifies an exact definition symbol or named-type origin. The host resolves selectors against inspected type information. Selection should have three outcomes:

- No matching extractor: report missing extraction support.
- Exactly one matching extractor: invoke it.
- Multiple matching extractors: fail with an ambiguity diagnostic.

A selected extractor’s inability to decode an expression must produce an explicit limitation. It must not trigger registration-order fallback.

The host then validates the result and attaches it to the original operation/profile binding. The result should contain no destination provider IDs, test outcomes, Markdown, rendering callbacks, or arbitrary payload field.

This makes the authority boundary small: an extractor interprets a definition; the host owns identity, association, lifecycle, and presentation.

### Use a fixed structural result vocabulary

The initial result schema can cover:

- Feature keys.
- Child-role names.
- Minimum and maximum cardinalities.
- Allowed-kind constraints.
- Attribute presence constraints.
- Finite attribute-value constraints.
- Reference allowance.
- Opaque checks.
- Source references.
- Localized extraction limitations.

Each uncertain field needs an explicit knowledge state. Zero values must not accidentally mean “unconstrained” or “known absent.”

The host can validate record shape, source-handle validity, uniqueness, bounds, and association scope. The format owner remains responsible for the interpretation of its definitions.

That division also governs future changes. A new encoding of existing structural facts should require changes only to its extractor and consumer composition. A genuinely new family of reported facts may justify a general contract extension under AC-08. It should not enter through an unstructured escape hatch.

### Expose source without sharing mutable compiler state

I recommend **immutable source snapshots plus narrow queries for existing Go type facts**.

The host should retain the exact source text used for discovery and expose:

- Request-bound file, declaration, expression, and symbol references.
- Immutable source strings for explicitly selected files.
- Physical source spans.
- Canonical symbol identities.
- Resolved callee and receiver identities.
- Named-type origins and concrete type arguments.
- Compile-time constant values.
- Access to selected referenced declarations.

Collections returned through this API should be copied or have clearly transferred ownership.

An extractor can parse its private source copy with the standard Go parser. It can then query host type facts using validated file/span references. This avoids both sharing the host’s mutable AST and requiring the host to recreate every Go AST shape in a new universal representation.

Go’s type checker already supplies identifier resolution, exact constant values, expression types, selections, and generic instantiations. Those are the factual products to expose through the narrow interface. [go.dev](https://go.dev/src/go/types/api.go?utm_source=chatgpt.com)

Do not expose the host’s `types.Info` maps, `*types.Package`, or discovery AST through an interface described as read-only. Also avoid deep-cloning the entire Go type graph merely to preserve that description.

Use physical loaded-file spans for source links. Adjusted positions from source line directives should not become identities or redirect links outside the selected source set.

### Recognize helpers statically

The grammar extractor should recognize the exact owning identities of `Rules`, `One`, `Opt`, `Many`, `Attr`, `Set`, and related definitions.

Consequently:

- An unrelated helper named `One` does not match.
- A shadowed identifier does not match.
- A true type alias follows the documented alias policy.
- A newly defined type with a similar underlying structure does not automatically match.
- Unsupported calls remain opaque.

The extractor need not import inspected runtime packages merely to identify them. Canonical identifiers resolved against inspected type information are sufficient.

Exact identity still depends on the owner maintaining the schema interpretation as its helpers evolve. Meaningful format tests should protect that relationship; executing inspected helpers is unnecessary.

### Make the host independently composable

There is a structural gap to settle in the proposed package layout: **a public extractor interface is insufficient if the usable host remains trapped in `package main`.**

Expose a tooling-only runner that accepts source selections, configured extractors, limits, and output options. Keep the command as a thin composition root.

An independent consumer must be able to import that runner and supply its extractor through the same mechanism used by the built-in grammar extractor. The host should not import the built-in extractor; their composition belongs in the command.

### Test the boundary with a different representation

The proposed struct-field fixture is a good choice. For example, an independent owner could define:

```go
type Forms struct {
    Conditional Conditional `form:"conditional-type"`
}

type Conditional struct {
    Check Node `role:"check" min:"1" max:"1"`
    True  Node `role:"true" min:"1" max:"1"`
}
```

Those tags are an illustrative fixture-owned schema.

Its extractor should select the exact definition, traverse struct declarations, interpret its own tags, and emit the same structural result records. Include an imported field type, an alias, and an unsupported representation.

The acceptance condition is that this works through the public runner without adding struct-tag interpretation to host orchestration, the grammar-table extractor, or the renderer. That demonstrates a materially different contribution within the promised boundary.

### State resource guarantees for the entire command

Separate three kinds of limits:

| Kind | Examples |
|---|---|
| Host-enforced data limits | Package/file counts, source bytes, query counts, accepted records, diagnostic bytes, subprocess output |
| Cooperative work limits | Reference depth, expansion counters, cycle detection, cancellation checks |
| External execution limits | Deadline covering loading, type checking, extraction, and rendering |

Replacing `exec.Command` with `CommandContext` is useful but incomplete. Output capture needs bounds, and process/I/O cleanup needs explicit handling. Go documents that `CommandContext` kills its process by default and leaves `WaitDelay` unset; `WaitDelay` addresses certain process and pipe waits. [Go Packages](https://pkg.go.dev/os/exec?utm_source=chatgpt.com)

Returning an oversized result can allocate memory inside a trusted extractor before host validation. Likewise, an in-process extractor or type checker cannot be forcibly cancelled through interface convention. If hard per-request resource isolation is promised, introduce a worker-process boundary with the corresponding enforcement.

For the initial version, a defensible contract is finite input/output and traversal limits, trusted extractors, and an outer remote-command deadline. Operational quota exhaustion should fail generation rather than quietly produce an apparently finished inventory.

## 3. Identities, wrappers, test loading, and directives

### Retain a contribution identity before rendering

The unit behind a matrix entry should be a discovered contribution, with identity retained before conversion to display text.

Conceptually:

```text
Contribution identity =
    full contract origin and concrete type arguments
    + public implementation or factory symbol
    + operation origin and instantiated receiver, where applicable
    + concrete language marker or ordered transformer endpoints

Association scope =
    contribution identity
    + optional exact profile label
```

Keep generic origins **and** concrete substitutions. Package short names, source positions, language display names, and rendered contract strings are presentation metadata.

Constructor dependency alternatives belong in the report as construction requirements. Adding an alternative should not change the contribution’s anchor.

Resolve each annotation to a candidate set of contributions, apply its qualifiers, and require an unambiguous association. A shared display cell must never cause associations to spread across all its implementations.

An omitted profile should be represented as actual absence, distinct from every string value. It cannot act as a wildcard, and equal profile labels cannot establish equal execution context.

### Correct the wrapper terminology and association rule

The supplied traversal accepts this shape:

```go
func Check(view View) error {
    return delegate.Check(rewrite(view))
}
```

It contains one return statement whose result is one statically resolved call, yet the wrapper transforms its argument.

Retain a **single-return static call chain**, including caller/callee identities, receiver selection, generic substitutions, call sites, and the reason traversal stopped. Use it for navigation and clearly attributed related information.

A wrapper report may show “cases declared against delegate X.” The actual test observation must remain attached to its declared target. Even a plain argument-preserving wrapper has not been observed executing merely because a delegate’s test passed.

Apply the same precision to transformer factories. The existing traversal establishes that typed direct constructions occur in a factory body. Preserve that discovery behavior, while avoiding claims about returned-value flow, runtime reachability, or composed routes.

### Namespace extracted features

Multiple roots can declare the same key with different constraints. Retain at least:

```text
Extracted feature identity =
    definition root + format owner + local feature key
```

A case’s `feature=identifier` can remain a valid operation-scoped authored label without being linked to an extracted item.

Link it to a specific inventory item only when the relationship is unambiguous. If a demonstrated use case requires stronger disambiguation, add a narrow definition-root qualifier. Do not merge roots or credit both merely because their strings match.

### Keep test source separate from production discovery

A separate test-declaration index is sufficient if it can resolve the supported reference syntax and validate test declarations. Analyze declarations and imports; no test-body interpretation is required.

Where full test typing is necessary, use an isolated graph and map resolved targets back to verified production declaration identities. Preserve test-file classification, package variant, and build selection. Test packages must not add providers to production discovery.

This distinction has concrete support in Go tooling: `go list -test` exposes synthetic test variants, and `go/packages` warns that type objects from separate load calls have inconsistent type identity and must not be mixed. [Go Packages](https://pkg.go.dev/cmd/go?utm_source=chatgpt.com)

There is also an overlooked event-matching collision: an internal test package and its external `_test` package can supply the same top-level test name. Go’s generated test registry stores the unqualified name separately from the package alias. [go.dev](https://go.dev/src/cmd/go/internal/load/test.go?utm_source=chatgpt.com)

Detect ambiguity in the mapping from `(tested package, test name)` to source declarations. Preserve both authored associations, but leave observations unassigned until the ambiguity is resolved. Do not infer ownership from event order or a presumed `_test` package suffix.

External dependency source and documentation must be explicitly selected. A qualified annotation reference should not trigger ambient downloads or arbitrary source inspection.

### Keep the directives, with two refinements

The three-directive design is reasonable. Use one raw-comment tokenizer with attachment-specific validation so the existing `dimension` and `shared` rules remain strict.

The parser should:

- Recognize standalone directive lines in declaration doc comments.
- Exclude trailing comments and fenced or indented examples.
- Validate `case` against an actual top-level test declaration.
- Parse quoted values using the documented Go-string subset.
- Reject unknown keys, repeated keys, and conflicting duplicate associations.
- Bound directive sizes and counts.
- Escape authored labels during rendering.

Using raw comment tokens is justified: `CommentGroup.Text()` removes directive-like comments and normalizes other formatting. [Go Packages](https://pkg.go.dev/go/ast?utm_source=chatgpt.com)

**First refinement: make expectations owner-defined.** `expects=accept|reject` is too narrow for formatting, printing, semantic judgments, and transformation assertions. Keep `expects`, but let it carry a nonempty owner-defined assertion label, documented by the linked scope or test.

Labels such as `semantic-invalid` and `operation-error` should remain distinguishable. The host does not interpret either as a Go test outcome.

A passing negative test is a passing test. An observed test failure does not reveal which operation-level result occurred.

**Second refinement: validate scope links against the owning module.** Resolve paths under that module’s root, validate the Markdown fragment, and reject lexical or symlink escapes. Prefer existing explicit anchors where available; otherwise use renderer-compatible heading rules, including duplicate headings.

### Generate stable anchors from canonical identities

Use explicit anchors derived from an injective encoding of the canonical identity tuple. A reversible encoding of length-delimited components satisfies the no-hash requirement.

Keep human-readable headings separate. Exclude source lines, revisions, case counts, outcomes, and discovery order from anchors. Preserve exact profile values and distinguish absence from a literal label such as `"unspecified"`.

The static report should also state its inspection build context. A source load under one platform/tag selection cannot establish availability under every build variant.

## 4. Test observations and publication

### Associate events with executions first

Use three separate records:

1. A test declaration.
2. Its authored case associations.
3. An observation from a particular test invocation.

The central join is:

```text
applicable invocation context
    + unambiguous test declaration mapping
    + exact emitted test name
    -> test observation
```

Then the report shows that observation alongside the authored associations.

If one function has four case annotations, its pass remains one function-level observation referenced by four relationships. It does not become four independently observed scenarios.

### Capture the actual Go JSON protocol

Use `go test -json` directly while preserving existing package selections and flags. Go recommends it for multipackage execution and distinguishes named test events from package events. Events can interleave. [Go Packages](https://pkg.go.dev/cmd/test2json?utm_source=chatgpt.com)

The importer also needs build-event support. The documented protocol includes interleaved `build-output` and `build-fail` records, whose `ImportPath` identity differs from `TestEvent.Package`. Serious early errors may still appear as non-JSON stderr. [Go Packages](https://pkg.go.dev/cmd/go?utm_source=chatgpt.com)

Therefore, decode stdout as the supported event protocol, capture stderr separately, and retain process exit status. Never derive outcomes from arbitrary `Output` text.

For repeated execution, preserve occurrence identity or explicitly reject unsupported repetition. “Last terminal event wins” can erase a failure.

### Preserve the distinctions that affect credit

| Situation | Report behavior |
|---|---|
| Named test passes | Show that function’s pass in this invocation |
| Named test skips | Show skip; award no successful assertion credit |
| Package passes without a matching named event | Show no named observation supplied |
| Test starts without a terminal event | Show incomplete observation |
| Parent passes while a child skips | Preserve both outcomes |
| Parent passes under a subtest filter | Show the selection; leave other scenarios unassessed |
| Build fails before tests run | Show build failure without inventing per-test failures |
| One test passes before another fails | Preserve the pass and the failed enclosing invocation |
| Names match but source/configuration does not | Show evidence mismatch |
| Event stream is malformed or truncated | Show capture/protocol failure; infer no successful completion |

Go’s subtest names and slash-aware filters make parent results particularly easy to overinterpret. A parent can run under a selection that excludes other children. [Go Packages](https://pkg.go.dev/testing?utm_source=chatgpt.com)

With dynamic-subtest declarations deferred, keep annotations at the top-level function. Child events can appear as subordinate execution details without becoming new cases.

Likewise, preserve fuzz runs as fuzz-invocation evidence unless an explicit association mechanism covers them. Fuzz duration does not supply a feature-scenario count.

Retain `-count=1` in the existing ordinary and race commands; Go documents it as the standard way to disable test-result caching. [Go Packages](https://pkg.go.dev/cmd/go?utm_source=chatgpt.com) For arbitrary imported results, freshness should remain unknown unless the captured invocation establishes it. The test-event schema has no authoritative cached-result boolean, and missing timestamps are only a convention. [Go Packages](https://pkg.go.dev/cmd/test2json?utm_source=chatgpt.com)

### Capture a small invocation context

The minimum useful context is:

| Context | Purpose |
|---|---|
| Actual checked-out revision and existing expected-revision check | Identify the tested checkout |
| Run ID, attempt, job, matrix variant, command identifier | Distinguish executions and reruns |
| Tested module/root and relevant module replacements | Identify the code selected by repository and consumer tests |
| Exact command and relevant build/test options | Preserve package selection, filters, race mode, tags, repetition, and similar choices |
| Actual Go version and target configuration | Explain source/build selection |
| Relevant native/helper setup and versions | Describe the available execution environment |
| Exit status and capture-completion state | Distinguish completed runs from interrupted capture |

Capture an allowlist of relevant settings, rather than the entire environment.

Record the actual checkout HEAD. GitHub’s triggering SHA is event-dependent, while its run-attempt value changes on reruns. [GitHub Docs](https://docs.github.com/en/actions/reference/workflows-and-actions/contexts?utm_source=chatgpt.com)

A SHA alone does not establish that generated files, overlays, or local replacements left the tested input unchanged. Use the controlled existing CI sequence to establish that relationship and disclose relevant deviations.

The independently installed consumer needs its own harness identity and actual selected Bittype module/revision. Its events must not be matched to same-named repository tests by intuition.

### Keep profile claims authored

The existing evidence can establish which test function reported an outcome and which environment the job prepared.

It cannot independently establish that:

- The function selected its annotated profile.
- A prepared compiler was actually invoked by that test.
- Every asserted branch executed.
- Equal profile labels imply equal collaborators or configuration.

Render the profile as an **authored test association**. Runtime profile observations would require an explicit owner-supplied mechanism; ordinary Go events do not provide one.

### Adopt the split publication model

For committed reports, include deterministic inventories, scope links, case associations, and the ordinary CI route.

For each test invocation:

1. Capture source/build context.
2. Run the existing command with JSON output.
3. Preserve events, stderr, exit status, and capture completeness.
4. Generate the observation report against the same declaration snapshot.
5. Upload the report, events, and context together.
6. Preserve the original test failure in the job result.

Avoid pipelines that accidentally retain only `tee`’s exit status.

Report generation and upload need explicit conditions after test failures because GitHub normally applies `success()` to steps. Bounded report work can use `!cancelled()`; log-style upload can use `always()` where cancellation behavior is intentional. [GitHub Docs](https://docs.github.com/en/actions/reference/evaluate-expressions-in-workflows-and-actions?utm_source=chatgpt.com)

Keep separate jobs separate. A native test can be skipped in the contracts job and pass in the native job. Both observations are useful. Missing native evidence should be described as not supplied unless an actual job-state observation establishes that it is pending.

Artifacts also have retention limits. Their later unavailability should not be interpreted as absence of past execution. [GitHub Docs](https://docs.github.com/en/actions/tutorials/store-and-share-data?utm_source=chatgpt.com)

### Avoid the documentation revision loop

The source-revision problem has a second consequence:

**Do not embed the updater’s current HEAD or generation timestamp in committed report content.**

Otherwise generation at S records S, the documentation commit creates D, and the next generation changes solely to record D.

Use stable repository-relative source links in committed reports. Put exact revisions and revision-pinned links in per-run artifacts. Git history already records static publication history.

There is no need to transfer S’s observations to D or to introduce a content-equivalence protocol for documentation-only commits.

## 5. Implementation sequence and acceptance tests

### Recommended sequence

**1. Preserve discovery while retaining identities.**  
Replace presentation-only locations with structured contribution records. Retain wrapper chains and generic substitutions. Verify that availability, dependency fixed points, and direct transformer directions remain unchanged.

**2. Introduce the shared directive parser and test-declaration index.**  
Resolve scope and target references, preserve profile absence, and diagnose ambiguity. Keep production discovery isolated from test variants.

**3. Expose the public tooling runner and source-view contract.**  
Implement both the grammar-table extractor and the independent struct-format fixture before committing to a large annotation rollout. This makes the extension boundary concrete early.

**4. Render complete static reports.**  
Create sections for all discovered implementations and factories, including those with no structural inventory. Add scope links for existing providers and decode the existing native grammar format across its users.

**5. Implement observation import and per-run rendering.**  
Exercise synthetic event streams first, then add JSON capture to the existing remote jobs without changing their intended test selections.

**6. Update publication for both generated files.**  
Render and validate both outputs before mutation. Expand diff checks, staging, and permitted generated paths together. A failed write must prevent the existing Git publication step.

**7. Add useful associations to current tests.**  
Associate existing tests where their owners can state the relationship accurately. Leave the remaining areas visibly unassessed.

Builds, tests, generation, and the independently installed fixture should run on the shared workers, consistent with the documented operating constraint.

### Highest-value counterexamples

These tests exercise the promised guarantees:

| Area | Counterexample | Required result |
|---|---|---|
| Format selection | Unrelated package defines `Rules` and `One` | No accidental recognition |
| Partial extraction | Role shape is known, predicate is custom | Preserve roles; mark predicate opaque |
| Unknown information | Allowed-kind expression cannot be decoded | Preserve unknown; do not emit an empty set |
| Multiple roots | Two roots declare `identifier` differently | Keep inventories separate |
| Resource limits | Cyclic references or oversized expansion | Bounded failure; no misleading finished report |
| Extension contract | Independent struct-tag extractor | Works through the public runner without host format branches |
| Generic identity | Same helper used for two language markers | Ambiguous association does not spread credit |
| Wrappers | `return delegate(rewrite(x))` | Retain relationship without transferring execution evidence |
| Test isolation | Test-only implementation satisfies a provider interface | It stays out of the production matrix |
| Test naming | Internal and external packages both define `TestFoo` | Diagnose ambiguous event matching |
| Execution | Package passes while annotated test skips | Preserve the skip |
| Subtests | Filtered child runs and parent passes | No credit for unselected scenarios |
| Failure handling | Pass followed by panic, truncation, or build failure | Preserve valid observations and enclosing failure |
| Applicability | Matching name, wrong revision or module replacement | Reject the association explicitly |
| Separate jobs | Native test skips in one job and passes in another | Preserve both contextual observations |
| Publication | Second generated-file write fails | Publish no generated commit |
| Determinism | Regeneration after a documentation-only commit | No change caused solely by the new HEAD |

Unknown syntax within a supported definition format can be represented as a localized extraction limitation. Identity ambiguity, invalid directives, invalid result records, loader failures, and operational quota exhaustion should fail generation by default.

If implementation review establishes an actual ownership or dependency violation, record that concrete failure and its repair state in the ACX corpus. A successful fixture or test run supplies evidence only for the behavior it exercises.

### Delivery boundary

The rollout is complete when every discovered implementation has a truthful detail section, existing providers have appropriate scope links, the current native grammar format is inventoried across its users, useful current tests are explicitly associated, and matching per-run results can be rendered through the same mechanism.

Remaining missing cases, opaque checks, absent extractors, and unavailable observations should remain explicit parts of that completed reporting system.