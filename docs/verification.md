# Joint verification requests and results

**Contract and first implementation, 6 November 2026.** This defines the
association contract requested after the [Logos investigation](logos-verification.md)
and [architecture review](type-system-architecture.md). The
[Go surface](#implemented-go-surface) implements this boundary for the scoped
Horos/Go composition. The walkthroughs below specify its claim semantics;
[integration tests](../verification/integration/go/joint_test.go) exercise actual
captured inputs and the landed Logos checker. No runtime binding is added.

The first composition prepares one exact claim about actual captured artifacts,
then passes a separate candidate to Logos's typed Go checker. Bittype owns the
artifact association; language and mapping integrations own interpretation;
Logos owns statement syntax and checking. The
[Logos delivery status](logos-verification.md#concrete-typed-check-only-contract)
is separate from this contract's design and implementation.

## Existing contracts and the new boundary

Reuse these existing types and operations:

| Existing contract | Use here |
| --- | --- |
| [`semantic.Snapshot`](../contracts/semantic/go/operation.go) | An immutable, syntax-admitted collection, optionally including files and their syntax links |
| [`semantic.Result`](../contracts/semantic/go/request.go) | An immutable provider publication for one semantic invocation and its exact snapshot |
| [`semantic.AdmitEvidence`](../contracts/semantic/go/evidence.go) | Check a side's receipts against its snapshot and receiver-specified operation, profile, configuration, report and model requirements |
| [`source.Source[L]`](../source/go/source.go) | Supply its existing snapshot and evidence when a concrete Files realization is already available |
| [`ASTTransformer`](../contracts/transform/go/transformer.go) | Produce syntax; it does not issue verification results or run native analysis |
| Logos `typed.Statement`, `Expr`, `Policy`, `Limits`, `Outcome` | The logical request, separate candidate, receiver policy, checking budget and outcome |

A source receipt cannot be attached to a target capture. Even two equal-looking
trees captured separately are different snapshots. The new association retains
both sides and their dependencies; it does not create a mixed-language semantic
snapshot or publish everything as one language's `semantic.Result`.

One request asks one owner-defined question. A conjunction such as Page's paired
result can be one proposition. TextBox admission and field representation are
separate questions; there is no implicit all-obligations success flag.

## Request construction record

The following is the complete logical construction record. These are contract
fields for the concrete Logos composition, not a serialization schema or a new
generic type-system algebra. Only the agnostic participant/evidence/site
association is a candidate shared contract; the logical fields stay in the
concrete integration.

| Field | Required content and authority |
| --- | --- |
| Participants | A finite collection of existing `semantic.Snapshot` handles, with request-local references identifying the source, actual target and additional dependencies. One snapshot may contain several units. |
| Evidence | Ordered `semantic.Result` receipts for each participant. The selected integrations supply the exact `EvidenceRequirement` lists; candidate producers cannot weaken them. |
| Bound interpretation | Immutable, integration-owned source/native projections, original declaration and parameter identities, the selected mapping objects/configuration, actual helpers and relevant environment. Retain their input associations and supported domain. |
| Question | An owner-qualified `semantics.Key` for reporting, plus the owning integration's typed subject/argument selection, checking point and promised scope. The key does not discover code or establish meaning. |
| Basis | The exact authoritative procedure results and mapping laws admitted by the selected projections, with their assumptions. Their correspondence to artifacts belongs to those integrations. |
| Checking policy | Receiver-owned `typed.Policy`, including admitted foundation/checker contracts and permitted dependencies. The installed checker is selected explicitly by composition. |
| Resource limits | Explicit preparation limits and `typed.Limits`. The latter bounds the Logos call, not extraction, compiler work or construction of proof syntax. |

The bound interpretation is not `any`, a bag of flags or a universal byte
descriptor. Its concrete types belong to the chosen language/mapping integration.
The first private composition can retain those types directly. Shared association
code deals with existing envelopes, receipt requirements and participant-local
references; it never unpacks them as `AST[any]` or interprets native judgments.
Strings naming a mapping are insufficient without its captured selection and
applicable laws. An arbitrary mutable plugin object is not a captured selection.

Participant references are unique within the request; the source and target
references must each identify a supplied participant. They are not language IDs
or globally meaningful artifact names. One target participant can contain the
whole generated family and its helper units. Additional participants need not
share a language, and their own receipts keep their original ownership.

All semantically relevant environment inputs must be included: supplied modules,
dependencies, target language/compiler profile, relevant build configuration and
the permitted extension/consumer domain. An external helper can be a dependency
participant or be covered by an explicit trusted integration premise. Omitting
its meaning while claiming its representation is not a supported third option.

The initial claim covers captured ASTs and the explicitly modeled environment.
If the promise concerns emitted files, the input must additionally bind those
files and the admitted formatter/printer correspondence. Capturing file bytes
alone does not prove that correspondence. No stage reads an ambient filesystem,
fetches dependencies or invents a compiler profile.

## Preparation and checking

The first integration has two operations in its explicitly configured Horos/Go
composition. `Claim` is an opaque, owning-integration construction, not a formula
or verdict supplied by a proof producer. There is no generic checker-provider SPI:

```go
func (c Composition) Prepare(ctx context.Context, input Claim) (Request, error)
func (c Composition) Check(ctx context.Context, request Request, candidate typed.Expr) (Result, error)

type Request struct { /* private immutable state */ }
type Result  struct { /* private immutable state */ }
```

The concrete views expose `Request.Statement()`, `Policy()` and
`Limits()` for proof construction, and `Result.Request()`, `Candidate()` and
`Outcome()` for inspection. They return the retained Logos types and defensive
statement/policy lists. Association and diagnostic views retain their participant
qualification. Neither a proof producer nor a view can modify the sealed request.

`Input` is the construction record above, with concrete owner-specific projection
types selected by the composition. `Composition` retains the actual configured
projection/claim-construction operations and checker binding. An ID/Name list
does not instantiate it; no registry, discovery or fallback is added.

### Prepare

1. Admit every participant and its evidence using the existing snapshot/profile
   contracts. Validate the input's declared bounds before copying collections.
   Match receipts to their own snapshots, never just their entry addresses.
2. Require the explicitly supplied source/native projections and their procedure
   results to belong to those admitted inputs. Check declared domain and mapping
   compatibility. A graph's structural validity or a stored `Holds` flag alone
   is not a Logos proof or a guarantee of model adequacy.
3. The receiver-selected integration constructs the expected `typed.Statement`
   from the actual subjects, binders, fields, helpers and fixed interpretation.
   Build the declaration signature, ordered logical context and exact expected
   proposition. Check local references and record their diagnostic associations.
4. Seal a fresh request retaining that association, statement, policy, limits
   and the actual bound composition. No candidate is an input to preparation.

Prepared means associated and immutable, not logically admitted or proved.
Logos still checks the declaration signature, context and expected expression
when it checks the candidate; preparation does not invent a dummy proof to
pre-approve them.

Resolution, analysis, generation and any required authoritative Horos query run
explicitly before preparation. The owning integrations must provide their
input-associated results. Preparation does not silently run an analyzer, search
for a proof, invoke a native compiler or accept unproved source membership as
nonmembership. The implemented explicit query/projection bridges supply these
premises; the existing general Horos report alone does not contain them.

The expected statement is not caller-supplied proof metadata. An artifact
verification composition has no constructor that attaches an arbitrary formula
to arbitrary artifacts. Its selected projection/claim constructor owns that
translation. This is an explicit trust boundary: structural association checks
cannot prove that a trusted extractor or assumed representation law is sound.

A prepared request owns these values:

| Retained value | Contract |
| --- | --- |
| Association | Original immutable participant/receipt handles, owned collection storage, captured question and interpretation |
| Statement | Exact `typed.Statement`, including its supplied declaration/context order; clone mutable slices while sharing immutable Logos nodes |
| Basis and policy | Artifact-to-premise association and exact receiver policy, including unused allowances; actual used dependencies are determined by checking |
| Limits | The selected preparation and checking limits; no automatic budget increase or weaker retry |
| Origin map | Request-local correspondence from logical declarations/context/expected positions to zero or more participant-qualified sites and model subjects |
| Composition binding | The actual admitted operations used to prepare/check this request, retained rather than rediscovered by name |

Logical hypotheses, trusted signature assumptions and allowed dependencies remain
different things. `typed.Assume` declares a trusted premise and its dependencies;
`typed.Define` contains a value Logos must check. A closed generic theorem uses
binders inside its expected expression and has an empty external context.

### Check

`Check` rejects a zero request or a request from another bound composition. It
uses only that request's statement, policy and limits when calling:

```go
outcome := typed.Check(statement, candidate, policy, limits)
```

The candidate is an immutable Logos expression, not a request, assumption list,
mapping choice or claimed result. Proof syntax can be reused in a fresh check;
it carries no acceptance authority. The wrapper does not invoke proof search.

After the call, the wrapper seals a new result retaining the request, the exact
submitted candidate and the actual `typed.Outcome`. Only this checking path
constructs a joint result. There is no public `Publish(accepted: true)`,
`Attach(otherOutcome)` or constructor that retags an earlier result. Retain the
request even for nonaccepting outcomes, because they do not contain a Logos
acceptance snapshot. This gives every diagnostic an unambiguous input context.

## Identity, lifetime and diagnostics

Request identity is a private invocation handle, not a digest, UUID scheme,
module name or persistent artifact record. Re-preparing equal inputs creates a
new request. Results retain their request and remain about that request even if
the caller later edits its source. A new proof attempt against one request
returns its own result and retains its own candidate.

Immutable snapshot/result/Logos-node handles can be shared. Construction arrays,
configuration, origin maps and policy lists must be copied or transferred under
the owning integration's immutable capture contract. Accessors return immutable
handles or defensive copies. Zero handles are invalid; validity of a carrier
does not mean that the requested proposition was accepted.
Request/result accessors return value handles over private immutable state,
never a mutable outer wrapper whose whole-value replacement changes a retained
request or result.

Model references retain their owning semantic graph or Horos Program. Reusing
an entity number or declaration spelling in another model does not transfer its
meaning. Reconstruct owner-local handles when canonical inputs are resolved
again; re-prepare and recheck the resulting statement.

Diagnostic sites are participant-qualified: `(participant, ast.NodeRef)`, not
bare node IDs. The same node address may occur on both sides. Preparation checks
that every mapped site belongs to its participant. Model subject references also
retain their owning model. Logos's logical `Location` and code are retained even
when there is no source location. Never invent a source span for a generic
candidate-proof error. A diagnostic may relate several source/output sites.

The origin map supports explanation; it is not evidence that an extractor or
mapping is correct. Likewise, the source/native/proof binder correspondence is
supplied by the owning projection. A common name or node address is insufficient
to establish it.

## Result and failure contract

A joint result exposes its request, submitted candidate, unchanged Logos outcome
and participant-qualified diagnostics. It retains an immutable Logos acceptance
when present, including the exact statement/candidate/policy and actual used
dependency closure. It does not reconstruct that acceptance from its fields.

| Condition | Return and meaning |
| --- | --- |
| Malformed input, foreign receipt, invalid local reference or incompatible composition | Preparation/checking error with located context; no joint result and no logical verdict |
| Required semantic report is Unproved, incomplete, unsupported or lacks the required model | Preparation cannot seal a request; retain the underlying provider report and coverage/completion reason, using the existing evidence-unavailable/failure contracts where applicable |
| Preparation limit or cancellation | Error preserving the limit/cancellation cause; no partial prepared request or successful result |
| Logos `AcceptedStatus` | Joint result establishing the exact requested proposition relative to its admitted basis |
| Logos `InvalidStatus` | Joint result retaining invalid checking input or a candidate/type mismatch; the logical location distinguishes signature/context/expected/candidate failures, and the requested proposition remains undetermined |
| Logos `IncompleteStatus` | Joint result retaining the exact limit/capacity cause; no logical verdict |
| Logos `UnsupportedStatus` | Joint result retaining the unserved foundation or unadmitted dependency/checker cause; no logical verdict |
| Projection or adapter failure outside Logos checking | Go error; never convert it into a checked counterexample |

Failure to prepare is observable; it is not represented by a zero successful
result. Preserve existing semantic error causes, reports and sites instead of
collapsing them into a Boolean. Do not introduce a second proof outcome algebra.
An actual Logos invalid-input outcome remains visible even when its logical
location exposes a defect in the prepared signature or expected expression;
it must not be mislabeled as only a bad proof or as semantic refutation.

The wrapper checks cancellation before invoking the synchronous checker and
again before publishing its result. Cancellation during checking does not
preempt Logos; its own budget bounds the call. If cancellation is observed before
publication, return the context error without publishing a joint result. A later
cancellation does not revoke an already returned immutable result. It never
becomes source nonmembership.

The exact requested claim matters. An accepted counterexample is positive
evidence of that counterexample; it is not an `Invalid` outcome. Only the owning
integration can relate it to refutation of a separate preservation proposition.
Any summary using `semantic.Report` must name that proposition and preserve this
distinction. The generic association layer never maps every accepted result to
"generation preserves all constraints" or every invalid candidate to `Rejected`.

## TextBox walkthrough

Use the ordinary kind module imported as `h`, and the source body:

```text
type TextBox<T: h.Text> = { value: T; };
```

The actual generated Go declaration under examination is:

```go
type TextBox[T any] struct { Value T }
```

The source capture contains TextBox and its supplied kind dependencies. The
target capture contains this declaration and its native environment. The source
projection identifies the original parameter/bound; the native projection
identifies the actual `any` constraint, parameter and Value field. Use the fixed
Go fragment/profile from the [architecture example](type-system-architecture.md#statement-one-textbox).
The selected interpretation associates canonical text with Go string and
canonical bool with Go bool before considering evidence.

This request asks for an admission counterexample. Its expected proposition is
the owner-defined encoding of:

```text
exists T in H_C, N in G_C:
  NativeAdmits_C(actualGoTextBox, N)
  and Rep_C(T, N)
  and SourceNonmember_C(T, actualTextBound)
```

Here `C` fixes the actual declarations, interpretation, domain and environment;
`H_C` and `G_C` are model descriptor domains, not new Horos kinds. The native
admission point and source argument contract are fixed by the question. The
claim constructor encodes the existential and conjunction using Logos's existing
typed constructors and the integration's vocabulary.

| Stage | Concrete content |
| --- | --- |
| Explicit source judgment | Horos refutes canonical bool against this resolved Text bound; retain the exact query, Program association and authoritative procedure basis |
| Explicit native judgment | Within the admitted native profile, actual `TextBox[bool]` is accepted under its `any` constraint |
| Selected mapping premise | The fixed representation relates canonical bool and Go bool; the proof cannot change that choice |
| Preparation | Seal these precise premises, independently constructed expected counterexample, permitted basis and source/target sites |
| Candidate | Supply bool witnesses and the proof assembling the three premises; do not assume the entire counterexample |
| Accepted result | This exact request has a counterexample to argument-admission preservation under its stated basis |

The diagnostic can relate the source bound and target `any` node. The result says
nothing about preservation of Value's representation; that needs a separate
request. Nor does a single counterexample prove the complete behavior of every
native type in the domain.

If Horos returns Unproved rather than Refuted, preparation lacks the required
nonmembership premise. If the bool association is absent, coverage is missing.
Neither case is a counterexample. If the target changes to `interface{ string }`,
the old result remains about the old target. A new target capture, native
judgment and preparation are required; bool no longer supplies the native
acceptance premise in this fragment. Absence of that counterexample alone does
not prove universal preservation. `~string` and defined types require their own
domain/interpretation analysis.

## Page walkthrough

The source body and ordinary bound alias are:

```text
type Page<T> = { items: list<T>; next?: text; };
type DataPage<D: h.Data>: h.Data = Page<D>;
type Person = { name: text; active: bool; };
```

The target capture includes the actual named family and helper:

```go
type BittypeOptional[V any] struct {
    Present bool
    Value   V
}

type Page[N any] struct {
    Items []N
    Next  BittypeOptional[string]
}
```

This request asks about Page's conditional static representation and membership.
Fix the resolved ordinary kind `K = h.Data`; Page remains unconstrained. Prepare
the exact owner-defined statement:

```text
forall T in H_C, N in G_C:
  Member_C(T, K) -> Rep_C(T, N) ->
    (Member_C(PageH(T), K) and Rep_C(PageH(T), PageG(N)))
```

`PageH` and `PageG` refer to the actual captured named families. `C` includes
the chosen list, text, optional and record representations and the actual helper.
The theorem does not assert representation existence for every source type.

| Preparation input | What it must establish or explicitly assume |
| --- | --- |
| Source projection | Page's original declaration is valid, its binder is identified, and the authoritative constructor-kind query establishes the conditional K-to-K guarantee under its original contracts |
| Native projection | The actual Page body, formal parameter, field identities/types and optional helper in the admitted native environment |
| Selected representation laws | Separate text, list, optional, record and named-family laws applicable to those exact choices; a named Go type is not silently equated with its anonymous body |
| Binder correspondence | Source T, native N and logical variables are associated by the projection; same spelling is insufficient |
| Expected statement | Quantification over T and N and both premises, returning both membership and representation; empty external `Statement.Context` |
| Policy and locations | Permitted procedure/mapping/model dependencies and sites for the source family, target family, fields and helper |

The candidate uses real Logos abstraction and application to combine the source
guarantee with the representation laws. It cannot assert the final Page theorem
as a new trusted premise, replace universal binders with one chosen atom or drop
an original source bound. Logos checks that candidate against the independently
prepared expected expression and reports the actual dependency closure.

An accepted joint result establishes this conditional static theorem under its
basis. It does not establish native admission for all Go `Page[N]` applications:
`any` still admits native arguments without the source premises. The separate
[`(S,N,E)` admission claim](go-native-admission.md) remains separate. It also
does not establish a codec, function call, stream behavior or runtime value
correspondence.

### Person and nested application

Prepare a new request for Person using the same admitted Page/mapping scope and
actual Person source/native projections. Its independently constructed expected
result is both `Member(PageH(PersonH), K)` and
`Rep(PageH(PersonH), PageG(PersonG))`.

Retain the closed theorem's defining signature and add a `typed.Define` containing
its exact checked expected expression and candidate. Supply the Person membership
and representation premises, and check the application anew. The old acceptance
is not converted into an assumed theorem. The receiver's current policy must
admit the rechecked definition's actual dependencies.

For a nested request, the first application's membership **and** representation
results supply the two arguments required by the same theorem. Check the nested
application against its own exact expected pair. A membership proof alone cannot
serve as a representation proof. Reuse under another capture reconstructs the
owner-local bindings and requires fresh preparation/checking; the earlier joint
result cannot be retagged.

### Changed artifacts and valid proof reuse

Changing the helper, a field, the target environment, source bound or selected
mapping requires fresh projections and preparation. Keep the actual body and
helper meaning in the statement's premises even if the family name is unchanged.
For a changed optional helper outside the old law, the old proof must fail or
preparation must report missing coverage. Reusing the same `n.page` spelling
cannot preserve the old association.

Not every textual change invalidates a theorem. When the new projections and
admitted laws justify the same claim, the same candidate syntax may pass a fresh
check. The invariant is exact request/result association and rechecked meaning,
not mandatory failure on every edit. Results already returned remain valid only
for their retained inputs.

## Placement, qualification and remaining work

Follow the existing [stage-first roles](type-system-architecture.md#proposed-package-roles).
The agnostic association part retains captures, evidence and participant-qualified
references. Source/native/mapping interpretations stay with their owners. The
Logos-specific prepared request/result and direct checker call belong to the
concrete checking integration; general contracts do not import `typed.Expr` or
grow a proof enum. Start with the private composition, then promote only its
demonstrated agnostic association surface to `contracts/verification/go`.

No new LanguageProvider or checker registry is required: a joint operation has
multiple participant languages, and explicit composition already selects its
collaborators. Do not register Page, Data, TextBox, Function, Stream or Cell as
general verification cases. Type-owned mapping plugins contribute through the
same [library-owned selection design](type-generation-plugins.md).

The first implementation supplies the source/query, native and mapping
projections, bounded immutable capture and direct Logos checking described below.
The scoped native projection is not a registered general Go semantic provider.

Qualification exercises foreign receipts/models, changed configuration and
helpers, invalid/unsupported/exhausted checking, TextBox's counterexample and
missing-premise cases, Page's closed theorem and nested reuse, and mutation of
construction/accessor storage.
Use actual captured inputs for association tests; Logos's symbolic fixtures alone
cannot establish this layer. Verify that cancellation, invalid evidence and
missing source premises never become semantic refutation. Preserve the owner
rule that all Logos workloads run on qualified shared/CI workers.

## Implemented Go surface

The [agnostic association](../contracts/verification/go/association.go) owns
participant-qualified snapshots, receipt requirements, resource limits and site
membership. It imports no Logos or native model. `Capture` admits each receipt
against its own exact snapshot and the receiver's requirements, with aggregate
budgets across participants. It never invokes another semantic stage.

The concrete [Horos/Go composition](../verification/files/horos-to-go/go/joint.go)
has sealed value handles `Claim`, `Request` and `Result`. Only the owning mapping
constructs a Claim. There is no public formula attachment or outcome publication
API. The request retains the original projections and queries, mapping selection,
actual helpers, and, for an application, the original theorem result and argument
projection. Copying or replacing an outer handle does not change retained state.

| Explicit step | Current API and meaning |
| --- | --- |
| Source projection | [Horos `Inspect`](../verification/files/source/horos/go/model.go) admits the existing resolution receipt and reconstructs a Program from that captured collection, checking all present annotations against Horos. Queries run explicitly through `Satisfies`; their sealed results retain the Program, value, kind, resolution module and exact Horos verdict. |
| Native semantics | [Go providers](native-semantics.md) convert captured ASTs directly to `go/ast` and check them with `go/types` under the explicit Go 1.25 environment. Analysis reuses the exact resolution graph. Immutable native readers retain its declarations and type facts; a separately constructed query session checks generic admission. |
| Counterexample construction | `Counterexample` requires the original source bound, an authoritative bool nonmembership result, actual native bool admission and this integration's scalar association. It constructs the existential claim independently of `CounterexampleCandidate`. |
| Family selection | `BindFamily` takes exact source/native declarations, explicit field pairs, an optional-helper selection, a resolved ordinary kind K and an authoritative K-to-K constructor query. It checks every source and target field, the original parameter identities and the helper's actual presence/value fields. |
| Closed theorem | `Family.Theorem` constructs the conditional membership-and-representation claim. `Family.Candidate` constructs separate proof syntax using ordinary Logos constructors. No final family theorem is assumed. |
| Argument and nesting | `Family.BindArgument` captures an actual nullary record and its authoritative membership query. `Family.Application` retains the earlier checked closed theorem as `typed.Define`. `Argument.Candidate(depth)` applies it using both projections of each preceding pair. Each result is checked afresh. |
| Prepare and check | `New(policy, limits)` captures receiver configuration. `Prepare(ctx, claim)` seals a fresh association. `Check(ctx, request, candidate)` calls the landed `typed.Check` directly and retains its unchanged outcome. |

Logical dependencies are explicit integration contracts:
`horos-go/captured-model-1` for the interpreted domains and symbols,
`horos/declarations-4/procedure`, `go1.25/captured-types-1`, and the selected
scalar, list/slice, optional/presence, record-field and named-representation laws.
`MappingDependencies()` returns a new list for explicit receiver configuration;
it installs no policy. The source/native procedure and projection correspondence
are trusted premises, not proofs produced by Logos. Representation laws state
the chosen static meaning of those mappings. Logos checks their composition and
reports the actual dependency closure, including the rechecked theorem definition.

The family interpretation uses private model descriptor domains H and N, not a
new Horos type language. Record positions refer to the selected actual field
pairs. The optional symbol denotes the selected field-presence helper; it does
not add an optional proper-type constructor to Horos. Named native types remain
distinct from their bodies and use a separate representation law. The source
body premise is conditional on `Member(T,K)`, so its original contract is not
silently extended to inadmissible arguments.

### Supported domain and limits

This first mapping supports a single proper-type parameter, closed records,
text, bool, lists of supported operands and optional fields using an explicitly
selected two-field presence/value helper. It supports nullary record arguments
for application. It recognizes no Page, TextBox or Data name; a renamed family,
helper and ordinary user kind use the same path. Alias unfolding, arbitrary
Horos descriptors, higher-order parameters and alternate native representations
remain outside this mapping's coverage. Missing coverage returns `ErrCoverage`;
missing established semantic premises return the existing evidence-unavailable
cause. Neither is a source refutation.

The Go semantic provider accepts the implemented type-only syntax fragment: named types,
parameters, applications, slices, structs, constraint interfaces/unions and
approximation, and aliases during package checking. Mapping targets must be
actual defined named types. Imports, methods, executable bodies and other type
forms are rejected. Its environment is the supplied package plus the predeclared
Go universe; it reads no ambient files or dependency packages. This is a scoped
projection, not a general native resolution/analysis capability in the engine.

Source inspection uses existing semantic and Horos bounds. Native inspection
admits at most 16 units, 4096 AST nodes, 32768 entries, 1 MiB of text and conversion
depth 64. Mapping selection is bounded to 64 field pairs and operand depth 32;
the application helper accepts depth 1 through 4. Joint association separately
bounds aggregate input size, receipt/configuration metadata and diagnostic sites.
The receiver supplies Logos limits explicitly; insufficient budgets produce the
unchanged incomplete outcome. Preparation neither proves the statement nor raises
the budget to make a proof pass. Go/Horos synchronous procedures, like Logos,
check cancellation before and after their bounded calls.

Claims cover captured ASTs and the stated interpretation. They do not prove
formatter/printer correspondence, public field construction, codecs, native
`(S,N,E)` admission enforcement or runtime value behavior. The broader
[type-owned plugin design](type-generation-plugins.md) and native enforcement
proposal remain separate. These selected mapping laws are not automatic defaults
or registrations in the syntax engine.
