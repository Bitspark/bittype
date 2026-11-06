# Logos checking and the Bittype verification boundary

**Source inspection and design refinement, 6 November 2026. Documentation only.**
The owner requested a deeper Logos request trace before implementing verification.
This document records that inspection and narrows the
[verification SPI proposal](type-system-architecture.md#remaining-public-verification-spi-decisions).
It adds no provider, checker, dependency, generator or runtime binding.

The inspected sources are Logos
[`d94e8dc9`](https://github.com/Bitspark/logos/tree/d94e8dc996ebba99e2bd26cde518df4918e2354f),
Bittype
[`182adb3`](https://github.com/Bitspark/bittype/tree/182adb30133eeb9e0f3f7e95b8a563ae73e785ea)
and public Horos
[`ad7c6f4`](https://github.com/Bitspark/horos/tree/ad7c6f493d726a82e445875d42af319adce0e3ce).
The traces below follow source and test fixtures; no query, proof, build or test
was executed for this inspection. Fixture outputs are their specified expectations.

## What exists and what remains proposed

| Surface | Inspected status | Consequence for Bittype |
| --- | --- | --- |
| Logos Go ground contract | Implemented `Program`, `Goal`, `GroundLiteral`, `Derivation`, `CheckProfile` and `Check` | Can check supplied ground derivations through the existing API |
| Logos Go kernel | Implemented capability negotiation, candidate production, certification and terminal outcomes | Use when an integration explicitly needs search; checking supplied evidence does not require this pipeline |
| Logos typed semantic packages | Draft `Statement`, `JudgmentRequest`, package/rule bindings, typed evidence and checking contract | Use this ownership and contract direction; these are not production Go types |
| Logos typed foundation and replay | Test-private Rust experiment | Demonstrates scoped typed checking and composition; cannot be imported as an admitted production Go checker |
| Bittype semantic infrastructure | Implemented captured requests/results, provider-owned facts and explicit transformation inputs | Supplies artifact association, not a logical foundation or proof of every stored fact |
| Language semantics | Horos resolution/analysis implemented; native analysis deferred | Native text support does not supply Go, TypeScript, Rust or C++ semantic judgments |
| Cross-language verification | Proposed | Source/native projections, mapping laws, joint request and outcome association still need implementation |

The public Horos [history boundary](https://github.com/Bitspark/horos/blob/ad7c6f493d726a82e445875d42af319adce0e3ce/README.md#history-and-license)
explicitly leaves earlier proof-integration APIs in the private predecessor.
Its current [declaration API](https://github.com/Bitspark/horos/blob/ad7c6f493d726a82e445875d42af319adce0e3ce/core/go/declarations/README.md#guarantees-and-uncertainty)
establishes kind judgments with `Proven`, `Refuted` and `Unproved`, and distinguishes
invalid input and exhaustion. Those authoritative procedure results are not
automatically portable Logos proofs.

## Production request trace

The Go embedding [example](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/kernel/go/embed/example_test.go#L22)
supplies this finite program and query:

```text
open(T) :- tray(T), not sealed(T).
tray(t1). tray(t2). tray(t3). sealed(t2).
?- open(T)
```

Its fixture expects `open(t1)` and `open(t3)`, with the answer stream exhausted.
The actual source path is:

1. **Capture the caller's program and request.**
   [`embed.Query`](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/kernel/go/embed/embed.go#L115)
   resolves the input, partitions rules and extensional facts, binds a checking
   profile and constructs the runtime request. The example's
   [leaf verifier](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/kernel/go/embed/profile.go#L41)
   admits exactly the supplied facts and binds their complete tables for
   closed-world checking. The candidate producer does not choose this policy.
2. **Resolve and prepare an eligible engine.** The
   [Go drive](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/kernel/go/runtime/drive.go#L209)
   uses program analysis and the required profile to resolve a plan, retaining
   the actual admitted provider objects. Preparation and execution use those
   bindings; an identifier alone is not execution authority.
3. **Produce untrusted candidates.** The
   [engine SPI](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/kernel/go/seam/spi.go#L52)
   separates `Probe`, `Prepare` and `Solve`. The embedding's
   [WASM-backed instance](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/kernel/go/wasm/provider.go#L1016)
   executes the retained admission and returns candidate substitutions and
   derivations. An engine instance cannot certify its own answers.
4. **Reconstruct the expected claim.** The
   [certification gate](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/kernel/go/certify/certify.go#L813)
   applies the candidate substitution to the caller's goal and requires a ground
   result. It checks that exact claim, never a replacement conclusion chosen by
   the producer.
5. **Check the supplied derivation.**
   [`CheckWithCertContext`](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/contract/go/check.go#L663)
   checks structure and limits, any closed-world certificates, rule
   instantiations and premises, theory decisions and external leaves. The root
   must equal the expected claim structurally. Only an accepting check permits
   the gate to mint a certified result. The example's negation needs checked
   closed-world evidence; failure to find `sealed(t1)` is insufficient.
6. **Report answers and completion separately.**
   [`embed.Answers`](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/kernel/go/embed/answers.go#L9)
   preserves certified partial answers when a run reaches a limit. Its
   `complete` flag describes exhaustion of the answer stream. Even an exhausted
   empty stream asserts that no answer was certified, not that the goal is false.

This is the Go embedding's implemented route, not a claim that every Logos
language driver supports identical plan forms or that a complete answer-family
theorem follows from enumeration.

## Direct checking needs no search provider

For supplied all-positive evidence, the existing call is:

```go
outcome := contract.Check(program, expectedClaim, candidateDerivation, ctx, profile)
```

These are Logos-owned types and a consumer-owned immutable context. The
[check-only example](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/contract/go/example_check_test.go#L24)
constructs a derivation of `permit(s1)`, checks it, then shows that the same
derivation does not establish `permit(s2)`. No registry, planner, search engine
or WASM guest is needed for that call.

When evidence contains closed-world certificates,
[`certify.CheckWithDerivedContext`](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/kernel/go/certify/certify.go#L861)
derives the checking context from the caller's live program and profile. Copying
the candidate certificate's own context pins would not establish that it answers
the caller's request. These derived bindings also do not automatically identify
Bittype's source/output snapshots: the integration must bind the actual objects
and facts it claims to check.

The production extension points have different jobs:

| Contract | Responsibility | Bittype use |
| --- | --- | --- |
| `LeafVerifier[C]` | Admit external ground facts against evidence and an immutable context | Admit precise source/native procedure results for matching captured inputs under an explicitly trusted integration |
| `Theory` / completion-aware ground decision | Decide owned ground predicates under the selected theory | A selected decision procedure, when its meaning and coverage fit |
| `RuleAdmissibility[C]` | Apply consumer policy to clause metadata after structural checking | Optional rule policy; its signature does not supply general binder or assumption-discharge checking |
| `EngineProvider[C]` | Advertise, prepare and execute candidate-producing engines | Needed only if the composition explicitly requests proof search |

The [consumer hook contracts](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/contract/go/spi.go#L218)
require deterministic, total checking over immutable inputs. An arbitrary
registered callback is not thereby proved sound or resource-bounded.
Outcome-bearing source/native checks must complete before a binary leaf hook
admits their exact result; a producer's `verified` flag is not such evidence.

The low-level `Check` result is **Accept or Reject with a code**, including codes
for checking limits. `IsAccept()` alone loses distinctions the integration needs.
The kernel classifies those rejections; the
[completion algebra](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/contract/go/completion/completion.go#L1)
also distinguishes completed decisions, caller limits and noncompletion. The
adapter must preserve the selected API's actual outcomes, not assume every
current method already returns the proposed typed checker result.

## Typed statements: existing proposal and experiment

Logos [design 0014](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/docs/design/0014-typed-package-production.md#2-public-semantic-types-and-their-authority)
and [canonical supplement 0015](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/docs/design/0015-typed-package-canonical-contract.md)
already propose the richer semantic contract:

```text
JudgmentRequest<F>
  statement: Statement<F>
    foundation
    package instances
    ordered context
    bound typed objects
    requested judgment application
  policy
    assurance and permitted dependencies
    evidence, retention and verification requirements
    selection and budgets
```

Logical context assumptions are distinct from the receiver's permitted
dependencies. The expected judgment and bound objects come from the caller.
Portable declarations/rules and the receiver's executable checker bindings have
separate authority. Checking recomputes reachable dependencies and checks the
root against the expected judgment. A decision judgment requires evidence for
its positive or declared negative branch; an unchecked Boolean is insufficient.

The test-private Rust [foundation](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/conformance/rs/tests/heterogeneous/foundation.rs#L1)
implements a bounded dependent checking fragment. In the
[replay experiment](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/conformance/rs/tests/heterogeneous/replay.rs#L813),
the receiver validates specialized native evidence, constructs typed premises,
checks their assembly against the requested proposition, recomputes dependencies
and applies its dependency policy before accepting. This exercises the intended
composition; it neither proves all admitted assumptions nor promotes the
experimental API or decoder to production.

The documents explicitly leave production contract ratification and checker
admission outstanding. Bittype must not copy the experimental calculus, treat
its types as a shipped Go API, or create a silent ground-fact fallback for a
requested quantified proof. A narrow Logos-owned typed, check-only delivery is
the [proposed symbolic-certification prerequisite](type-system-architecture.md#the-logos-checking-prerequisite).

## Ownership map for the verification SPI

The earlier SPI gaps include both Bittype-specific association work and
Logos-owned checking contracts. Their separation is:

| Concern | Existing or proposed Logos contract | Bittype responsibility still to specify |
| --- | --- | --- |
| Expected claim and candidate evidence | Current `GroundLiteral` / `Derivation`; proposed typed `Statement` / evidence envelope | Language-owned projection and independent expected-claim construction; how concrete typed bindings cross the integration boundary |
| Joint source/output invocation | Current consumer context; proposed typed bound objects and receiver bindings | Capture both artifacts, their compatible semantic results, mappings, helpers and relevant environment; copying/lifetime and publication authority |
| Checking capability and compatibility | Current explicit `CheckProfile`; proposed foundation/package/service bindings | Select a supported route and check its required inputs without interpreting language judgments in the generic kernel |
| Outcome | Current checking codes, certification outcomes and completion algebra; proposed typed checked results | Associate the outcome and its assumptions with this invocation and source sites; preserve invalid evidence, refutation and noncompletion as distinct |

The minimal Bittype boundary therefore carries artifact association, delegates
semantic projection to owning integrations, invokes a supplied Logos checking
route, and publishes associated diagnostics. It does not define another proof
term algebra, generic theorem checker, foundation registry or automatic checker
discovery. No separate public checker-provider interface is required merely to
call `Check`; any shared SPI is promoted only for a demonstrated association or
composition requirement. Its concrete public Go shape remains undecided.

The existing [semantic request/result](../contracts/semantic/go/request.go)
contract seals a result against its capture. `Result.Valid()` means the result
exists as a valid carrier, not that every semantic claim is true. A `Holds`
report is not independently checked Logos evidence. The source integration must
admit the authoritative procedure result under its actual input/assumption
contract or supply checkable evidence. Native extraction and mapping laws need
the same explicit adequacy boundary.

## TextBox and Page exercise different capabilities

The [TextBox example](type-system-architecture.md#statement-one-textbox) compares
the source `TextBox<T: h.Text>` with an actual Go declaration using `T any`.
A counterexample requires Go acceptance of `bool`, the fixed association to
canonical `bool`, and Horos refutation of `bool : h.Text`. Current ground
checking can compose those exact premises into a counterexample once the
language integrations and rule are supplied. Failure to prove source membership
or rejection of malformed evidence is not a counterexample. Correct argument
admission also leaves the generated field-representation obligation separate.

The [Page theorem](type-system-architecture.md#statement-two-symbolic-page) is
universally quantified over qualifying source/native arguments:

```text
Member(T,K) and Rep_C(T,N)
  implies Member(PageHoros(T),K)
      and Rep_C(PageHoros(T),PageNative(N))
```

Horos already supports the source constructor-kind query. Native analysis must
describe the actual generated declaration, and the selected list, optional-field
and record mappings must supply their representation laws. The complete proof
needs scoped assumptions, binder correspondence, discharge and application under
an admitted checking contract. Enumerating concrete Page instances does not
establish that theorem; a ground leaf carrying its name does not check it.

Native admission remains the independent
[`(S,N,E)` obligation](go-native-admission.md): the source witness fixes its
intended canonical subject through `Denotes_C`, the native payload represents
that subject, and the evidence argument is admitted for that pair. No inverse
from payload N to source identity is assumed. The Go marker-method proposal
retains its explicit [emission scope extension](type-system-architecture.md#go-marker-method-scope-extension);
this inspection authorizes no method or runtime implementation.

## Missing work and next contract deliverable

Within Bittype and the owning language/mapping integrations, specify:

1. A joint immutable association for the actual source, output, semantic results,
   selected mappings/helpers and relevant language environment.
2. Source/native semantic projections and the expected-claim constructor, with
   real declaration/binder correspondence, supported domains and explicit
   procedure/model adequacy assumptions.
3. The selected Logos route, concrete argument/result bindings and outcome
   projection to Bittype diagnostics. Logical assumptions, permitted dependencies
   and operational limits remain distinct.
4. Representation laws owned by the selected mapping integrations. Mapping
   selection or successful compilation alone does not establish those laws.

For production symbolic certification, Logos still needs to ratify and expose
the required typed checking contract with an admitted implementation accessible
from the Go integration. A supported Go API or an explicitly specified process
client could supply that access; neither is selected here. General proof search,
a hosted service and the entire Logos roadmap are not prerequisites.

Walk through TextBox, Page, nesting, changed helpers, foreign captures,
incompatible checking contracts and exhausted checks before finalizing the public
association SPI. Those design exercises need not wait for complete implementations.
Keep [ASTTransformer](../contracts/transform/go/transformer.go) returning syntax,
resolution/analysis explicitly invoked, and the
[stage-first package roles](type-system-architecture.md#proposed-package-roles).
No general component learns named kinds or Function, Stream, Cell or Page.
