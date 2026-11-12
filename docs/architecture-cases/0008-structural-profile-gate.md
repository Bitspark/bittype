# bittype/ACX-008: The shared horos mapping policy admitted one profile

**Status:** Corrected, with limits; state of 12 November 2026. The adapter of
the common horos mapping package holds no list of admitted profile names, and
the two identities are `structural-representation/1` and
`structural-enforcement/1`. No independent review has read the correction. The
[limits](#limits) say what it does not establish.
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
The shared policy's `Library.Describe` returned `ErrUnsupported` for every
profile whose name was not the one constant `RepresentationProfile`. A second
profile that reused the horos patterns and the declaration recipe therefore
needed an edit to this shared file.

[AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)
asks reviewers to inspect helpers outside the kernel. That is where this gate
sat. The kernel was not at fault: it treats a profile as an opaque identity and
selects only the mappings registered under the requested one
([registry](https://github.com/Bitspark/bittype/blob/ca5dc80/transform/syntax/shared/generation/go/registry.go#L185)).

The profile's name was a second, related fault, and not a rule violation by
itself. Both identities began with the word `pagination`, which names the
first fixture and not what the profile promises
([the two constants at `ca5dc80`](https://github.com/Bitspark/bittype/blob/ca5dc80/transform/syntax/horos/mapping/go/mapping.go#L15-L16)).
It was corrected with the gate because the plan records profile identities
([#195](https://github.com/Bitspark/bittype/issues/195)).

## Intended boundary

The kernel's selection by profile identity stays invariant. Profiles vary. Each
profile's owner supplies its mapping libraries under its own identity, and the
consumer's request selects one. The common horos policy supplies patterns and
the declaration recipe to any of them. An unsupported profile fails because no
mapping is registered for it, or because its owner refuses it, never because a
shared file lists the admitted names.

## Actual implementation and consequence

At `ca5dc80` the adapter began with this check:

```go
func (p Library[T]) Describe(from, to representation.Language) (...) {
	if p.Profile != RepresentationProfile {
		return ..., g.ErrUnsupported
	}
```

Every constructor of a native mapping library takes a `profile` argument, for
example `text.Library(profile)` and `record.Library(profile)` in the packages
under `transform/syntax/horos-to-<native>/mapping/<name>/go`. At `ca5dc80`
these were `Text(profile)` and `Libraries(profile)` in one package per pair;
[bittype/ACX-009](0009-package-boundaries.md) records that split. Exactly one
value of that argument worked. The documented purpose of the check was to
refuse the enforcement profile, whose encoding does not exist
([GEN-02](../generation-gaps.md#gen-02)). It was written as an allowlist of
one, so it refused every other profile as well.

Consequences:

- A bounded profile, or any second structural profile, could not be registered
  through the common policy without changing it.
- Integers, maps, unions and every later family would have been recorded in
  plan entries under an identity that says pagination. Renaming it after
  [#207](https://github.com/Bitspark/bittype/issues/207) writes the first entry
  would have changed an identity that binding backends pin.

This was architectural and not missing coverage: the missing profile was not
another contribution behind a sound contract. The shared helper had to change
to admit it.

## What worked and what it did not establish

All generation tests selected the one profile, and they passed. An independent
library could enter the kernel directly: an independent contribution in the
[composition test](../../integration/go/horos_generation_composition_test.go)
implements the mapping contract itself and sets the profile on its own
description. Those results established that the kernel is open. They did not
establish that the common horos policy was, because no test registered a second
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

Corrected on 12 November 2026 by pull request
[#265](https://github.com/Bitspark/bittype/pull/265), which closes
[#216](https://github.com/Bitspark/bittype/issues/216) and was made against
`main` at `9cfd195`. It is ordered before #207 under the owner's direction of
11 November 2026 that a misplaced or wrongly named part changes first and
missing coverage is added afterwards
([AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)).

### Correction

| Was | Is |
| --- | --- |
| Both identities began with `pagination` | `RepresentationProfile` is `structural-representation/1` and `EnforcementProfile` is `structural-enforcement/1`. The owner of research 0015 chose them and recorded the reason on #216. The Go names are unchanged. No alias exists |
| The adapter refused every profile but `RepresentationProfile` | The adapter copies the profile it is given into the library's description and compares it with no admitted name. The registry's selection by the requested profile is the only rule that admits one |
| The enforcement profile was refused as one of all other names | The adapter refuses to describe a library under `EnforcementProfile`, with `ErrUnsupported`, before it runs a rule. A request for it fails with `ErrUnsupported` in the registry, which finds no mapping |
| The integration fixture repeated the identity as a string | [`horosgeneration.Request`](../../integration/go/testdata/horosgeneration/request.go) takes it from the horos mapping package. No test or fixture spells either identity |

**Where each responsibility sits, and why.**

- **The identities** stay in `transform/syntax/horos/mapping/go`. They are the
  two promises of the horos mapping side, and every bittype mapping library is
  built through that package's adapter.
- **Admission of a profile** has no owner outside the kernel. A library states
  its profile in its own description, and the registry selects by the
  request's profile. A profile under which nothing is registered fails there
  with `ErrUnsupported`. That behaviour existed before and is unchanged; the
  kernel and `contracts/**` are not edited.
- **The refusal of the enforcement profile** sits with the package that
  declares that identity. Removing the gate alone would have left the refusal
  to the kernel's "nothing registered". A composition could then build the
  erasing libraries under the enforcement name, and a request for the
  enforcement profile would return representation output, which
  [GEN-02](../generation-gaps.md#gen-02) forbids. The fact behind the refusal
  is that nothing built on the common horos mapping policy enforces a kind
  constraint. The package that declares the promise and owns that policy is
  the one that knows it.

This refusal is not the closed list in a smaller form. Asked the question of
the [assessment lesson](#assessment-lesson): every value of the `profile`
argument succeeds except the one that the declaring package itself has not
implemented, and the owner of a new profile edits nothing here.

One edge changed. An adapter value with an empty profile was refused by the
gate with `ErrUnsupported`. The contract's `FreezeDescription` now refuses it
with `ErrInvalid`, which is the contract's error for a description without a
profile.

### Qualification

The tests are in
[`integration/go/horos_generation_profile_test.go`](../../integration/go/horos_generation_profile_test.go)
and beside the adapter in
[`mapping_test.go`](../../transform/syntax/horos/mapping/go/mapping_test.go).

- **A second profile, in all four targets.**
  `TestHorosMappingLibrariesAnswerUnderASecondProfile` builds the builtin
  libraries of each pair under an identity that only the test file names,
  registers them and generates. The text equals the text of the structural
  profile. A registry that holds only the second profile refuses a request for
  the structural one. On `9cfd195` the test fails in all four targets at
  registration, with one `ErrUnsupported` per library.
- **Two profiles in one registry.**
  `TestRequestSelectsTheLibrariesOfItsProfile` registers a second profile whose
  list mapping exists under that profile only and is built through the common
  adapter and the common patterns. Each request receives its own profile's
  list. It also fails on `9cfd195`.
- **The enforcement profile.** `TestEnforcementProfileIsRefused` requests it
  from a registry that has just generated representation output, and from a
  composition that offers every builtin library under the enforcement
  identity. Every library is refused with `ErrUnsupported`, and both requests
  fail with `ErrUnsupported` with no plan and no output, in all four targets.
- **An unregistered profile.** `TestUnregisteredProfileIsRefused` requests an
  identity under which nothing is registered, from an empty registry and from
  one that holds the structural profile. Both fail with `ErrUnsupported`.
- **Two seeded faults.** With the closed list put back into the adapter, the
  two second-profile tests and the adapter's own test fail. With the
  enforcement refusal removed, `TestEnforcementProfileIsRefused` fails in all
  four targets: every library registers, and the request returns the
  representation text.
- **Output.** The text of three fixtures in the four targets, twelve files,
  has the same SHA-256 before and after. So have all 91 texts that the existing
  integration tests format.
- **The whole suite.** `go test ./...` on the author's Windows machine gives
  the same outcome for every test that existed before. 28 outcomes are new and
  pass. The ten failing outcomes are the same before and after and come from that
  machine: the Rust semantic tests have no rustc helper there, and one test of
  `tools/provider-matrix` fails in a temporary directory on a network share.
- **Imports.** `TestProductionArchitecture` passes with no new entry and no
  new exception. The change adds no package and no production import.

### Limits

- **No independent review.** The author wrote the correction, the tests and
  this record.
- **The refusal moves when enforcement is implemented.** The statement "nothing
  built on the common policy enforces" holds for all libraries at once today.
  When [#162](https://github.com/Bitspark/bittype/issues/162) implements an
  enforcing encoding, some libraries will keep that promise and others will
  not, and the refusal has to be restated by whoever then knows which. The
  enforcement interface is not designed
  ([GEN-02](../generation-gaps.md#gen-02)), so this change does not anticipate
  it. Until then the owner of a library that enforces has to edit this
  package. That is the declaring package finishing its own promise, and not a
  second profile asking for admission.
- **The refusal covers the adapter, not the kernel.** A library that implements
  the mapping contract itself and never passes through the adapter can publish
  any description under any identity, including the enforcement one. The
  kernel treats a profile as opaque by design, and selecting a contribution
  does not establish its claim
  ([AC-05](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-05)).
  No test covers such a library.
- **The second profile of the tests reuses or replaces builtin mappings.** No
  bounded profile with a narrower domain exists yet. The tests show that the
  common policy admits one. They do not show that the declaration recipe, the
  patterns and the naming function suit it.
- **The naming rule is not decided.**
  [bittype/ACX-009](0009-package-boundaries.md#candidates) passed the question
  of the `Self` rule in the shared naming function to #216. This change does
  not decide it. What changed is its setting: the rule now binds the libraries
  that call the naming package, which are the common declaration recipe and
  Go's record mapping, and no longer every profile that can exist. It stays a
  candidate.
- **The names in plan entries.** No plan entry exists yet
  ([#207](https://github.com/Bitspark/bittype/issues/207)). That the new
  identities are what the entries record is for that work to show.
- **Verbatim records keep the old names.** The documents under `research-docs/`
  are submitted consultations and their replies and are unchanged.

## Assessment lesson

The check was reviewed as the refusal of one unsupported profile, which is
correct behaviour. Its form went unexamined: a refusal written as "not the one
admitted name" is a closed catalogue of one. When a helper takes an identity as
a parameter, ask how many values of that parameter can succeed, and who would
have to edit what to add another.

The correction adds the converse. Deleting a closed list is not the whole
repair when the list also carried a guarantee. Ask what the list refused on
purpose, and give that refusal to the owner of the refused thing, by name.
