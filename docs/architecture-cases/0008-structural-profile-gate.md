# bittype/ACX-008: The shared horos mapping policy admits one profile

**Status:** Identified; correction pending.
**Scope:** The common horos mapping policy,
[`transform/syntax/horos/mapping/go`](https://github.com/Bitspark/bittype/blob/ca5dc80/transform/syntax/horos/mapping/go/mapping.go),
at `ca5dc80`, and the four `horos-to-<native>` mapping packages that build
their libraries through it. Identified on 11 November 2026. Tracked in
[#216](https://github.com/Bitspark/bittype/issues/216).

## Rules involved

[AC-03](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-03)
requires that adding a contribution change only its owning implementation and
the consumer's composition. A representation profile is such a contribution:
the [accepted families design](../representation-families.md#domains) names the
lossless structural profile and "an explicitly named bounded profile" beside it.
The shared policy's `Library.Describe` returns `ErrUnsupported` for every
profile whose name is not the one constant `RepresentationProfile`. A second
profile that reuses the horos patterns and the declaration recipe therefore
needs an edit to this shared file.

[AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)
asks reviewers to inspect helpers outside the kernel. That is where this gate
sits. The kernel is not at fault: it treats a profile as an opaque identity and
selects only the mappings registered under the requested one
([registry](https://github.com/Bitspark/bittype/blob/ca5dc80/transform/syntax/shared/generation/go/registry.go#L185)).

The profile's name is a second, related fault, and not a rule violation by
itself. The identity `pagination-representation/1` names the first fixture, not
what the profile promises. It is corrected with the gate because the plan
records profile identities ([#195](https://github.com/Bitspark/bittype/issues/195)).

## Intended boundary

The kernel's selection by profile identity stays invariant. Profiles vary. Each
profile's owner supplies its mapping libraries under its own identity, and the
consumer's request selects one. The common horos policy supplies patterns and
the declaration recipe to any of them. An unsupported profile fails because no
mapping is registered for it, or because its owner refuses it, never because a
shared file lists the admitted names.

## Actual implementation and consequence

```go
const RepresentationProfile = "pagination-representation/1"
const EnforcementProfile = "pagination-enforcement/1"

func (p Library[T]) Describe(from, to representation.Language) (...) {
	if p.Profile != RepresentationProfile {
		return ..., g.ErrUnsupported
	}
```

Every rule constructor in the four native mapping packages takes a `profile`
argument, for example `Text(profile)` and `Libraries(profile)`. Exactly one
value of that argument works. The documented purpose of the check is to refuse
the enforcement profile, whose encoding does not exist
([GEN-02](../generation-gaps.md#gen-02)). It is written as an allowlist of one,
so it refuses every other profile as well.

Consequences:

- A bounded profile, or any second structural profile, cannot be registered
  through the common policy without changing it.
- Integers, maps, unions and every later family would be recorded in plan
  entries under an identity that says pagination. Renaming it after
  [#207](https://github.com/Bitspark/bittype/issues/207) writes the first entry
  changes an identity that binding backends pin.

This is architectural and not missing coverage: the missing profile is not
another contribution behind a sound contract. The shared helper must change to
admit it.

## What worked and what it did not establish

All generation tests select the one profile, and they pass. An independent
library can enter the kernel directly: an independent contribution in the
[composition test](../../integration/go/horos_generation_composition_test.go)
implements the mapping contract itself and sets the profile on its own
description. Those results establish that the kernel is open. They do not
establish that the common horos policy is, because no test registers a second
profile through it.

## Required correction

1. The common policy holds no list of admitted profile names.
2. A request for the enforcement profile still fails with `ErrUnsupported` and
   never yields representation output. A request for a profile with no
   registered mapping fails.
3. The structural and enforcement profiles get identities that state their
   promise and name no fixture, with no alias for the old names. The owner of
   research 0015 chooses them.
4. The native declarations generated for the existing fixtures stay
   byte-identical.

## Correction and qualification status

Nothing is implemented. The gate and both identities are unchanged at `ca5dc80`.
The correction is ordered before #207 under the owner's direction of 11 November
2026 that a misplaced or wrongly named part changes first and missing coverage
is added afterwards
([AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)).
Qualification needs a test that registers a second profile through the common
policy and fails on `ca5dc80`.

## Assessment lesson

The check was reviewed as the refusal of one unsupported profile, which is
correct behaviour. Its form went unexamined: a refusal written as "not the one
admitted name" is a closed catalogue of one. When a helper takes an identity as
a parameter, ask how many values of that parameter can succeed, and who would
have to edit what to add another.
