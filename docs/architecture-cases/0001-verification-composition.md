# bittype/ACX-001: Closed verification contributions

**Status:** Identified; correction pending.
**Recorded:** 6 November 2026, following the owner's architectural review.
**Scope:** Bittype's concrete Horos-to-Go static verification composition.

This is the same incident first recorded as
[platform/ACX-001](https://github.com/Bitspark/platform/blob/main/docs/concepts/architecture-cases/0001-bittype-verification-composition.md).
Bittype now maintains the detailed record under the
[project ownership rule](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#rule-scope-and-project-ownership).
The initial observation concerned
[a8062e3](https://github.com/Bitspark/bittype/tree/a8062e34fae156a604585fda7692b2d268a601a1/verification);
the code-level review below inspected
[87e9f15](https://github.com/Bitspark/bittype/tree/87e9f15930afec41cc830d45df2f85fe877a42d7/verification).
The [verification contract](../verification.md#known-architectural-violation)
and [type-system architecture](../type-system-architecture.md) are the owning
designs. This record does not claim the existing conditional derivations are invalid.

## Rules involved

| Rule | Actual violation |
| --- | --- |
| [AC-02: Semantic ownership](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02) | The reusable verification path owns the available mapping interpretations, law catalogue and claim forms instead of receiving owner-controlled contributions. |
| [AC-03: Additive extension](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-03) | A different representation requires editing the matcher, proof construction and claim path, or duplicating orchestration. Adding its implementation and consumer selection is insufficient. |
| [AC-04: Uniform public mechanism](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-04) | The initial cases issue sealed claims through private family/argument/counterexample constructors. Independent mappings have no equivalent public contribution contract. |

[AC-01](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-01),
[AC-06](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-06)
and [AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07)
explain the assessment error: bounded proof success and a clean lowest-level
kernel were treated as sufficient while the closed extension path was described
as future coverage. The agreed composition boundary applies from the first slice.
[AC-05](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-05)
states guarantees the replacement must preserve; this review does not establish
that current receipt, scope or Logos checks are unsound.

## Intended boundary

Shared verification associates captured artifacts and semantic receipts,
preserves identity, scope, receiver policy and resource limits, and invokes the
selected checker. Language integrations own their interpreted native/source
subjects and judgments. Mapping integrations own their selected representation
relations, law dependencies and candidate construction. Consumer composition
selects those contributions explicitly. Initial and independently added mappings
must use the same public mechanism.

Semantic ownership does not require a separate repository or Go module for every
mapping. Dedicated Bittype integration packages can own the mappings while using
the public contract. Moving a closed catalogue to such a package, or calling it a
plugin, does not establish the required boundary.

Bittype's agnostic association remains independent of Logos syntax. Logos owns
proof syntax and checking. Neither learns Page, TextBox, a mapping's child cases
or a behavioural contract's semantics.

## Actual implementation and consequence

The closure is distributed across the following cooperating sections. Fixing
only one of them leaves the contribution path closed.

| Inspected source | Misplaced responsibility and required replacement |
| --- | --- |
| [mapping.go:28, private pattern](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/mapping.go#L28), [scalar:75](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/mapping.go#L75), [match:118](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/mapping.go#L118) | Every field is reduced to the private parameter/text/bool/list vocabulary. The list branch recursively calls this same matcher. A selected list mapping must inspect its own source/target shapes and request element correspondence through the common contribution path. |
| [mapping.go:95, optionalValid](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/mapping.go#L95), [bindRecord:149](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/mapping.go#L149) | Record traversal directly selects the presence/value helper interpretation and then calls the closed matcher. Keep field completeness, field identity and binder checks with the record contribution; move helper-specific interpretation to its selected contribution. Parameter-count coverage limits are not by themselves the defect. |
| [record_logic.go:21, terms](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/record_logic.go#L21), [proof:37](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/record_logic.go#L37) | Two switches interpret the private pattern and choose child proof laws. Record proof construction must consume corresponding child evidence without interpreting the child's mapping vocabulary. |
| [logic.go:88, MappingDependencies](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/logic.go#L88), [mappingBasis:92](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/logic.go#L92) | A fixed catalogue constructs scalar, list/slice, optional and named laws. Each selected mapping must supply its exact law and declared justification through the common contract; the receiver controls permitted dependencies. Go-authored law construction is not itself a violation. |
| [joint.go:41, claim](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/joint.go#L41), [mapping.go:314, Application](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/mapping.go#L314) | The reusable claim retains private family, argument, theorem-depth and counterexample states; reuse checks private family identity. Preserve sealing and exact owner/scope association through a common contribution contract. Concrete family state remains with its implementation. |
| [joint.go:88, logicalOrigins](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/joint.go#L88) | Diagnostic associations are inferred from private symbol prefixes such as source., native. and law. The replacement needs explicit contribution-supplied locations validated against captured participants and statement positions. This is another dependency on the private vocabulary, not evidence of a false proof acceptance. |
| [counterexample.go:25, Counterexample](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/counterexample.go#L25), [CounterexampleCandidate:57](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/files/horos-to-go/go/counterexample.go#L57) | Both admission and proof construction fix the witness to bool. A counterexample contribution must consume a selected source/native witness pair and matching nonmembership, admission and representation evidence. Bool remains a fixture, not a reusable-path restriction. |

For example, replacing Go slices with an independently owned collection requires
new native interpretation and a different law. Today it also requires changing
the private matcher and proof switches before that law can compose with a record
proof. That additional change propagation is the architectural failure.

## What worked and what it did not establish

The [scoped integration tests](https://github.com/Bitspark/bittype/blob/87e9f15930afec41cc830d45df2f85fe877a42d7/verification/integration/go/joint_test.go)
exercise TextBox's admission counterexample and Page's conditional theorem,
including nested reuse, ordinary kinds, foreign handles and changed helpers.
The [agnostic association](../../contracts/verification/go/association.go) checks
actual captures and receipts; [joint checking](../../verification/files/horos-to-go/go/joint.go)
calls the public Logos checker. Those are useful preserved foundations. The
present review inspected source; it did not rerun executable tests or establish
an extensible replacement.

The current laws are constructed as explicit Logos assumptions. Logos checks
their composition and dependency closure, not their independent truth. The
selected relation describes static type/schema association; it does not establish
runtime value correspondence or universal native enforcement of Horos constraints.

TypeScript's provider already exports captured native facts, including arrays
and their element types. The illustrated TypeScript list-to-array law and dense
array value relation are not implemented verification contributions. A stored
semantic fact, a compiler query result and a representation law are distinct.
See the [native provider contract](../native-semantics.md#first-typescript-profile).

The independent proof fixture in [Logos PR3887](https://github.com/Bitspark/logos/pull/3887)
qualifies that checker's contribution boundary under its own stated scope. It
does not demonstrate that an independent mapping can enter Bittype's closed
verification path and cannot close this case.

## Required correction

Replace the closed matcher, pattern interpreter, law catalogue and claim issuance
path with owner-controlled contributions selected by consumer composition. The
first scalar/list/optional/record mappings must use the same public mechanism as
an independently supplied mapping. Preserve the current useful checks and remove
the superseded executable path without a compatibility branch.

A list contribution's preparation would inspect the exact captured list and
native collection, request a child representation judgment for their elements,
and bind its selected law and dependencies. Candidate construction would apply
that law to the matching child evidence. A record contribution would combine its
field evidence without learning the child's representation cases. These are
required responsibilities, not a finalized or implemented new SPI.

The replacement must preserve:

- Independent construction of the expected judgment from captured artifacts and
  receiver-selected interpretations; candidate production cannot replace it.
- Exact language/model, declaration, mapping configuration and binder identities,
  including nested scopes and original source constraints.
- Association of every child result and native/source query with the subjects
  requested, including selected helpers and declaration locations.
- Explicit admitted assumptions and independently checked evidence; a selected
  contribution cannot authorize its own truth or mint an acceptance outcome.
- Imported theorem checking as definitions, retaining their original dependencies,
  instead of promoting a prior accepted result to an unchecked assumption.
- Receiver policy, immutable capture, bounded work and distinct invalid,
  unsupported, exhausted and rejected-proof outcomes.
- A Logos-free agnostic association contract and explicit Logos-owned checking.

Replacing private fields with `any` or exposing arbitrary formula attachment does
not satisfy these obligations. Neither does adding another dedicated TypeScript
path alongside the closed Go path. Native factual extraction can remain in its
own providers; the defect concerns their composition with mapping interpretation
and checking.

## Correction and qualification status

**Pending.** [Bittype PR123](https://github.com/Bitspark/bittype/pull/123) recorded
the original correction requirement. This documentation update moves the owning
case here and makes the code-level findings and replacement obligations explicit.
It changes no executable verification code, finalizes no SPI and proves no
mapping law. Further supported language/type features are separate from repairing
this extension boundary.

Qualification must add a materially different mapping in an independently owned
package or fixture through the same public contract, compose it with an existing
mapping, and exercise nesting plus invalid/foreign/stale interactions. Adding it
must require only its implementation and consumer composition, with no shared
orchestration, private-case or unrelated-mapping edits. Explain the preserved
guarantees and remaining limits; another renamed Page is insufficient.

When the repair lands, update this same record with the actual implementation,
qualification and remaining limits. A partial repair remains pending. The
platform reference points here so correction status has one maintained home.

## Assessment lesson

Trace how an independent meaning enters the complete operation and what must
change. A clean kernel, dedicated package, generic type parameter, successful
proof or explicit coverage limit does not establish an open contribution boundary.
Missing contributions can be deferred; a mechanism that cannot admit them through
the promised contract requires architectural correction.
